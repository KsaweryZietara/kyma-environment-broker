package memory

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/kyma-project/kyma-environment-broker/common/pagination"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/dberr"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/dbmodel"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/predicate"
	"github.com/pivotal-cf/brokerapi/v12/domain"
)

type instances struct {
	mu                      sync.Mutex
	instances               map[string]internal.Instance
	operationsStorage       *operations
	subaccountStatesStorage *SubaccountStates
}

func NewInstance(operations *operations, subaccountStates *SubaccountStates) *instances {
	return &instances{
		instances:               make(map[string]internal.Instance, 0),
		operationsStorage:       operations,
		subaccountStatesStorage: subaccountStates,
	}
}

func (s *instances) GetDistinctSubAccounts() ([]string, error) {
	//iterate over instances and return distinct subaccounts
	collectedSubAccounts := make(map[string]struct{})
	for _, v := range s.instances {
		collectedSubAccounts[v.SubAccountID] = struct{}{}
	}
	//convert map keys to slice
	var subAccounts []string
	for k := range collectedSubAccounts {
		subAccounts = append(subAccounts, k)
	}
	return subAccounts, nil
}

func (s *instances) UpdateInstanceLastOperation(instanceID, operationID string) error {
	return nil
}

func (s *instances) FindAllJoinedWithOperations(prct ...predicate.Predicate) ([]internal.InstanceWithOperation, error) {
	var instances []internal.InstanceWithOperation

	// simulate left join without grouping on column
	for id, v := range s.instances {
		dOps, dErr := s.operationsStorage.ListDeprovisioningOperationsByInstanceID(id)
		if dErr != nil && !dberr.IsNotFound(dErr) {
			return nil, dErr
		}
		pOps, pErr := s.operationsStorage.ListProvisioningOperationsByInstanceID(id)
		if pErr != nil && !dberr.IsNotFound(pErr) {
			return nil, pErr
		}

		if !dberr.IsNotFound(dErr) {
			for _, op := range dOps {
				instances = append(instances, internal.InstanceWithOperation{
					Instance:       v,
					Type:           sql.NullString{String: string(internal.OperationTypeDeprovision), Valid: true},
					State:          sql.NullString{String: string(op.State), Valid: true},
					Description:    sql.NullString{String: op.Description, Valid: true},
					OpCreatedAt:    op.CreatedAt,
					IsSuspensionOp: op.Temporary,
				})
			}
		}

		if !dberr.IsNotFound(pErr) {
			for _, op := range pOps {
				instances = append(instances, internal.InstanceWithOperation{
					Instance:       v,
					Type:           sql.NullString{String: string(internal.OperationTypeProvision), Valid: true},
					State:          sql.NullString{String: string(op.State), Valid: true},
					Description:    sql.NullString{String: op.Description, Valid: true},
					OpCreatedAt:    op.CreatedAt,
					IsSuspensionOp: false,
				})
			}
		}

		if dberr.IsNotFound(dErr) && dberr.IsNotFound(pErr) {
			instances = append(instances, internal.InstanceWithOperation{Instance: v})
		}
	}

	for _, p := range prct {
		p.ApplyToInMemory(instances)
	}

	return instances, nil
}

func (s *instances) FindAllInstancesForRuntimes(runtimeIdList []string) ([]internal.Instance, error) {
	var instances []internal.Instance

	for _, runtimeID := range runtimeIdList {
		for _, inst := range s.instances {
			if inst.RuntimeID == runtimeID {
				instances = append(instances, inst)
			}
		}
	}

	if len(instances) == 0 {
		return nil, dberr.NotFound("instances with runtime id from list %+q not exist", runtimeIdList)
	}

	return instances, nil
}

func (s *instances) FindAllInstancesForSubAccounts(subAccountslist []string) ([]internal.Instance, error) {
	var instances []internal.Instance

	for _, subAccount := range subAccountslist {
		for _, inst := range s.instances {
			if inst.SubAccountID == subAccount {
				instances = append(instances, inst)
			}
		}
	}

	return instances, nil
}

func (s *instances) GetNumberOfInstancesForGlobalAccountID(globalAccountID string) (int, error) {
	numberOfInstances := 0
	for _, inst := range s.instances {
		if inst.GlobalAccountID == globalAccountID && inst.DeletedAt.IsZero() {
			numberOfInstances++
		}
	}
	return numberOfInstances, nil
}

func (s *instances) GetByID(instanceID string) (*internal.Instance, error) {
	inst, ok := s.instances[instanceID]
	if !ok {
		return nil, dberr.NotFound("instance with id %s not exist", instanceID)
	}

	// In database instance details are marshalled and kept as strings.
	// If marshaling is ommited below, fields with `json:"-"` are never cleared
	// when stored in memory db. Marshaling in the current contenxt allows for
	// memory db to behave similary to production env.
	marshaled, err := json.Marshal(inst)
	unmarshaledInstance := internal.Instance{}
	err = json.Unmarshal(marshaled, &unmarshaledInstance)

	op, err := s.operationsStorage.GetLastOperation(instanceID)
	if err != nil {
		if dberr.IsNotFound(err) {
			return &inst, nil
		}
		return nil, err
	}

	detailsMarshaled, err := json.Marshal(op.InstanceDetails)
	detailsUnmarshaled := internal.InstanceDetails{}
	err = json.Unmarshal(detailsMarshaled, &detailsUnmarshaled)
	unmarshaledInstance.InstanceDetails = detailsUnmarshaled

	return &unmarshaledInstance, nil
}

func (s *instances) Delete(instanceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.instances, instanceID)
	return nil
}

func (s *instances) Insert(instance internal.Instance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.instances[instance.InstanceID] = instance

	return nil
}

func (s *instances) Update(instance internal.Instance) (*internal.Instance, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	oldInst, exists := s.instances[instance.InstanceID]
	if !exists {
		return nil, dberr.NotFound("instance %s not found", instance.InstanceID)
	}
	if oldInst.Version != instance.Version {
		return nil, dberr.Conflict("unable to update instance %s - conflict", instance.InstanceID)
	}
	instance.Version = instance.Version + 1
	s.instances[instance.InstanceID] = instance

	return &instance, nil
}

func (s *instances) GetActiveInstanceStats() (internal.InstanceStats, error) {
	return internal.InstanceStats{}, fmt.Errorf("not implemented")
}

func (s *instances) GetERSContextStats() (internal.ERSContextStats, error) {
	return internal.ERSContextStats{}, fmt.Errorf("not implemented")
}

func (s *instances) List(filter dbmodel.InstanceFilter) ([]internal.Instance, int, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var toReturn []internal.Instance

	offset := pagination.ConvertPageAndPageSizeToOffset(filter.PageSize, filter.Page)

	instances := s.filterInstances(filter)
	sortInstancesByCreatedAt(instances)

	for i := offset; (filter.PageSize < 1 || i < offset+filter.PageSize) && i < len(instances); i++ {
		toReturn = append(toReturn, s.instances[instances[i].InstanceID])
	}

	return toReturn,
		len(toReturn),
		len(instances),
		nil
}

func (s *instances) ListWithSubaccountState(filter dbmodel.InstanceFilter) ([]internal.InstanceWithSubaccountState, int, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var toReturn []internal.InstanceWithSubaccountState

	offset := pagination.ConvertPageAndPageSizeToOffset(filter.PageSize, filter.Page)

	instances := s.filterInstances(filter)
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].CreatedAt.Before(instances[j].CreatedAt)
	})

	for i := offset; (filter.PageSize < 1 || i < offset+filter.PageSize) && i < len(instances); i++ {
		instanceToReturn := s.instances[instances[i].InstanceID]
		instanceWithSubaccountState := internal.InstanceWithSubaccountState{
			Instance: instanceToReturn,
		}
		if _, exists := s.subaccountStatesStorage.subaccountStates[instanceToReturn.SubAccountID]; exists {
			instanceWithSubaccountState.BetaEnabled = s.subaccountStatesStorage.subaccountStates[instanceToReturn.SubAccountID].BetaEnabled
			instanceWithSubaccountState.UsedForProduction = s.subaccountStatesStorage.subaccountStates[instanceToReturn.SubAccountID].UsedForProduction
		}
		toReturn = append(toReturn, instanceWithSubaccountState)
	}

	return toReturn,
		len(toReturn),
		len(instances),
		nil
}

func sortInstancesByCreatedAt(instances []internal.Instance) {
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].CreatedAt.Before(instances[j].CreatedAt)
	})
}

func (s *instances) filterInstances(filter dbmodel.InstanceFilter) []internal.Instance {
	inst := make([]internal.Instance, 0, len(s.instances))
	var ok bool
	equal := func(a, b string) bool {
		return a == b
	}
	shootMatch := func(shootName, filter string) bool {
		return shootName == filter
	}

	for _, v := range s.instances {
		if ok = matchFilter(v.InstanceID, filter.InstanceIDs, equal); !ok {
			continue
		}
		if ok = matchFilter(v.GlobalAccountID, filter.GlobalAccountIDs, equal); !ok {
			continue
		}
		if ok = matchFilter(v.SubscriptionGlobalAccountID, filter.SubscriptionGlobalAccountIDs, equal); !ok {
			continue
		}
		if ok = matchFilter(v.SubAccountID, filter.SubAccountIDs, equal); !ok {
			continue
		}
		if ok = matchFilter(v.RuntimeID, filter.RuntimeIDs, equal); !ok {
			continue
		}
		if ok = matchFilter(v.ServicePlanName, filter.Plans, equal); !ok {
			continue
		}
		if ok = matchFilter(v.ServicePlanID, filter.PlanIDs, equal); !ok {
			continue
		}
		if ok = matchFilter(v.ProviderRegion, filter.Regions, equal); !ok {
			continue
		}
		if len(filter.Shoots) > 0 {
			if ok = matchFilter(v.InstanceDetails.ShootName, filter.Shoots, shootMatch); !ok {
				continue
			}
		}
		if ok = s.matchInstanceState(v.InstanceID, filter.States); !ok {
			continue
		}

		inst = append(inst, v)
	}

	return inst
}

func matchFilter(value string, filters []string, match func(string, string) bool) bool {
	if len(filters) == 0 {
		return true
	}
	for _, f := range filters {
		if match(value, f) {
			return true
		}
	}
	return false
}

func (s *instances) matchInstanceState(instanceID string, states []dbmodel.InstanceState) bool {
	if len(states) == 0 {
		return true
	}
	op, err := s.operationsStorage.GetLastOperation(instanceID)
	if err != nil {
		// To support instance test cases without any operations
		return true
	}

	for _, s := range states {
		switch s {
		case dbmodel.InstanceSucceeded:
			if op.State == domain.Succeeded && op.Type != internal.OperationTypeDeprovision {
				return true
			}
		case dbmodel.InstanceFailed:
			if op.State == domain.Failed && (op.Type == internal.OperationTypeProvision || op.Type == internal.OperationTypeDeprovision) {
				return true
			}
		case dbmodel.InstanceError:
			if op.State == domain.Failed && op.Type != internal.OperationTypeProvision && op.Type != internal.OperationTypeDeprovision {
				return true
			}
		case dbmodel.InstanceProvisioning:
			if op.Type == internal.OperationTypeProvision && op.State == domain.InProgress {
				return true
			}
		case dbmodel.InstanceDeprovisioning:
			if op.Type == internal.OperationTypeDeprovision && op.State == domain.InProgress {
				return true
			}
		case dbmodel.InstanceUpgrading:
			if op.Type == internal.OperationTypeUpgradeCluster && op.State == domain.InProgress {
				return true
			}
		case dbmodel.InstanceUpdating:
			if op.Type == internal.OperationTypeUpdate && op.State == domain.InProgress {
				return true
			}
		case dbmodel.InstanceDeprovisioned:
			if op.State == domain.Succeeded && op.Type == internal.OperationTypeDeprovision {
				return true
			}
		case dbmodel.InstanceNotDeprovisioned:
			if !(op.State == domain.Succeeded && op.Type == internal.OperationTypeDeprovision) {
				return true
			}
		}
	}

	return false
}

func (s *instances) ListDeletedInstanceIDs(int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	resultMap := make(map[string]struct{})
	for _, op := range s.operationsStorage.operations {
		if _, exists := s.instances[op.InstanceID]; !exists {
			resultMap[op.InstanceID] = struct{}{}
		}
	}
	var result []string
	for k := range resultMap {
		result = append(result, k)
	}
	return result, nil
}

func (s *instances) DeletedInstancesStatistics() (internal.DeletedStats, error) {
	panic("not implemented")
}
