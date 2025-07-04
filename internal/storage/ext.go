package storage

import (
	"time"

	"github.com/kyma-project/kyma-environment-broker/common/events"
	"github.com/kyma-project/kyma-environment-broker/common/runtime"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/dbmodel"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/predicate"
)

type Instances interface {
	FindAllJoinedWithOperations(prct ...predicate.Predicate) ([]internal.InstanceWithOperation, error)
	FindAllInstancesForRuntimes(runtimeIdList []string) ([]internal.Instance, error)
	FindAllInstancesForSubAccounts(subAccountslist []string) ([]internal.Instance, error)
	GetByID(instanceID string) (*internal.Instance, error)
	Insert(instance internal.Instance) error
	Update(instance internal.Instance) (*internal.Instance, error)
	Delete(instanceID string) error
	GetActiveInstanceStats() (internal.InstanceStats, error)
	GetERSContextStats() (internal.ERSContextStats, error)
	GetDistinctSubAccounts() ([]string, error)
	GetNumberOfInstancesForGlobalAccountID(globalAccountID string) (int, error)
	List(dbmodel.InstanceFilter) ([]internal.Instance, int, int, error)
	ListWithSubaccountState(dbmodel.InstanceFilter) ([]internal.InstanceWithSubaccountState, int, int, error)

	ListDeletedInstanceIDs(int) ([]string, error)

	DeletedInstancesStatistics() (internal.DeletedStats, error)

	UpdateInstanceLastOperation(instanceID, operationID string) error
}

type InstancesArchived interface {
	GetByInstanceID(instanceId string) (internal.InstanceArchived, error)
	Insert(instance internal.InstanceArchived) error
	TotalNumberOfInstancesArchived() (int, error)
	TotalNumberOfInstancesArchivedForGlobalAccountID(globalAccountID string, planID string) (int, error)
	List(filter dbmodel.InstanceFilter) ([]internal.InstanceArchived, int, int, error)
}

//go:generate mockery --name=Operations --output=automock --outpkg=mocks --case=underscore
type Operations interface {
	Provisioning
	Deprovisioning
	UpgradeCluster
	Updating

	GetLastOperation(instanceID string) (*internal.Operation, error)
	GetLastOperationByTypes(instanceID string, types []internal.OperationType) (*internal.Operation, error)
	GetOperationByID(operationID string) (*internal.Operation, error)
	GetNotFinishedOperationsByType(operationType internal.OperationType) ([]internal.Operation, error)
	GetOperationStatsByPlan() (map[string]internal.OperationStats, error)
	GetOperationStatsByPlanV2() ([]internal.OperationStatsV2, error)
	GetOperationsForIDs(operationIDList []string) ([]internal.Operation, error)
	ListOperations(filter dbmodel.OperationFilter) ([]internal.Operation, int, int, error)

	InsertOperation(operation internal.Operation) error
	UpdateOperation(operation internal.Operation) (*internal.Operation, error)
	ListOperationsByInstanceID(instanceID string) ([]internal.Operation, error)
	ListOperationsByInstanceIDGroupByType(instanceID string) (*internal.GroupedOperations, error)
	ListOperationsInTimeRange(from, to time.Time) ([]internal.Operation, error)

	DeleteByID(operationID string) error
	GetAllOperations() ([]internal.Operation, error)
}

type Provisioning interface {
	InsertProvisioningOperation(operation internal.ProvisioningOperation) error
	GetProvisioningOperationByID(operationID string) (*internal.ProvisioningOperation, error)
	GetProvisioningOperationByInstanceID(instanceID string) (*internal.ProvisioningOperation, error)
	UpdateProvisioningOperation(operation internal.ProvisioningOperation) (*internal.ProvisioningOperation, error)
	ListProvisioningOperationsByInstanceID(instanceID string) ([]internal.ProvisioningOperation, error)
}

type Deprovisioning interface {
	InsertDeprovisioningOperation(operation internal.DeprovisioningOperation) error
	GetDeprovisioningOperationByID(operationID string) (*internal.DeprovisioningOperation, error)
	GetDeprovisioningOperationByInstanceID(instanceID string) (*internal.DeprovisioningOperation, error)
	UpdateDeprovisioningOperation(operation internal.DeprovisioningOperation) (*internal.DeprovisioningOperation, error)
	ListDeprovisioningOperationsByInstanceID(instanceID string) ([]internal.DeprovisioningOperation, error)
	ListDeprovisioningOperations() ([]internal.DeprovisioningOperation, error)
}

type UpgradeCluster interface {
	InsertUpgradeClusterOperation(operation internal.UpgradeClusterOperation) error
	UpdateUpgradeClusterOperation(operation internal.UpgradeClusterOperation) (*internal.UpgradeClusterOperation, error)
	GetUpgradeClusterOperationByID(operationID string) (*internal.UpgradeClusterOperation, error)
	ListUpgradeClusterOperationsByInstanceID(instanceID string) ([]internal.UpgradeClusterOperation, error)
}

type Updating interface {
	InsertUpdatingOperation(operation internal.UpdatingOperation) error
	GetUpdatingOperationByID(operationID string) (*internal.UpdatingOperation, error)
	ListUpdatingOperationsByInstanceID(instanceID string) ([]internal.UpdatingOperation, error)
	UpdateUpdatingOperation(operation internal.UpdatingOperation) (*internal.UpdatingOperation, error)
}

type Events interface {
	InsertEvent(level events.EventLevel, message, instanceID, operationID string)
	ListEvents(filter events.EventFilter) ([]events.EventDTO, error)
}

type SubaccountStates interface {
	UpsertState(state internal.SubaccountState) error
	DeleteState(subaccountID string) error
	ListStates() ([]internal.SubaccountState, error)
}

type Bindings interface {
	Insert(binding *internal.Binding) error
	Update(binding *internal.Binding) error
	Get(instanceID string, bindingID string) (*internal.Binding, error)
	Delete(instanceID, bindingID string) error
	ListByInstanceID(instanceID string) ([]internal.Binding, error)
	ListExpired() ([]internal.Binding, error)
	GetStatistics() (internal.BindingStats, error)
}

type Actions interface {
	InsertAction(actionType runtime.ActionType, instanceID, message, oldValue, newValue string) error
	ListActionsByInstanceID(instanceID string) ([]runtime.Action, error)
}
