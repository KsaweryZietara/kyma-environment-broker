// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	internal "github.com/kyma-project/kyma-environment-broker/internal"
	dbmodel "github.com/kyma-project/kyma-environment-broker/internal/storage/dbmodel"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Operations is an autogenerated mock type for the Operations type
type Operations struct {
	mock.Mock
}

// DeleteByID provides a mock function with given fields: operationID
func (_m *Operations) DeleteByID(operationID string) error {
	ret := _m.Called(operationID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(operationID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDeprovisioningOperationByID provides a mock function with given fields: operationID
func (_m *Operations) GetDeprovisioningOperationByID(operationID string) (*internal.DeprovisioningOperation, error) {
	ret := _m.Called(operationID)

	if len(ret) == 0 {
		panic("no return value specified for GetDeprovisioningOperationByID")
	}

	var r0 *internal.DeprovisioningOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*internal.DeprovisioningOperation, error)); ok {
		return rf(operationID)
	}
	if rf, ok := ret.Get(0).(func(string) *internal.DeprovisioningOperation); ok {
		r0 = rf(operationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.DeprovisioningOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(operationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeprovisioningOperationByInstanceID provides a mock function with given fields: instanceID
func (_m *Operations) GetDeprovisioningOperationByInstanceID(instanceID string) (*internal.DeprovisioningOperation, error) {
	ret := _m.Called(instanceID)

	if len(ret) == 0 {
		panic("no return value specified for GetDeprovisioningOperationByInstanceID")
	}

	var r0 *internal.DeprovisioningOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*internal.DeprovisioningOperation, error)); ok {
		return rf(instanceID)
	}
	if rf, ok := ret.Get(0).(func(string) *internal.DeprovisioningOperation); ok {
		r0 = rf(instanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.DeprovisioningOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(instanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastOperation provides a mock function with given fields: instanceID
func (_m *Operations) GetLastOperation(instanceID string) (*internal.Operation, error) {
	ret := _m.Called(instanceID)

	if len(ret) == 0 {
		panic("no return value specified for GetLastOperation")
	}

	var r0 *internal.Operation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*internal.Operation, error)); ok {
		return rf(instanceID)
	}
	if rf, ok := ret.Get(0).(func(string) *internal.Operation); ok {
		r0 = rf(instanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.Operation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(instanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNotFinishedOperationsByType provides a mock function with given fields: operationType
func (_m *Operations) GetNotFinishedOperationsByType(operationType internal.OperationType) ([]internal.Operation, error) {
	ret := _m.Called(operationType)

	if len(ret) == 0 {
		panic("no return value specified for GetNotFinishedOperationsByType")
	}

	var r0 []internal.Operation
	var r1 error
	if rf, ok := ret.Get(0).(func(internal.OperationType) ([]internal.Operation, error)); ok {
		return rf(operationType)
	}
	if rf, ok := ret.Get(0).(func(internal.OperationType) []internal.Operation); ok {
		r0 = rf(operationType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.Operation)
		}
	}

	if rf, ok := ret.Get(1).(func(internal.OperationType) error); ok {
		r1 = rf(operationType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOperationByID provides a mock function with given fields: operationID
func (_m *Operations) GetOperationByID(operationID string) (*internal.Operation, error) {
	ret := _m.Called(operationID)

	if len(ret) == 0 {
		panic("no return value specified for GetOperationByID")
	}

	var r0 *internal.Operation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*internal.Operation, error)); ok {
		return rf(operationID)
	}
	if rf, ok := ret.Get(0).(func(string) *internal.Operation); ok {
		r0 = rf(operationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.Operation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(operationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOperationByInstanceID provides a mock function with given fields: instanceID
func (_m *Operations) GetOperationByInstanceID(instanceID string) (*internal.Operation, error) {
	ret := _m.Called(instanceID)

	if len(ret) == 0 {
		panic("no return value specified for GetOperationByInstanceID")
	}

	var r0 *internal.Operation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*internal.Operation, error)); ok {
		return rf(instanceID)
	}
	if rf, ok := ret.Get(0).(func(string) *internal.Operation); ok {
		r0 = rf(instanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.Operation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(instanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOperationStatsByPlan provides a mock function with given fields:
func (_m *Operations) GetOperationStatsByPlan() (map[string]internal.OperationStats, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetOperationStatsByPlan")
	}

	var r0 map[string]internal.OperationStats
	var r1 error
	if rf, ok := ret.Get(0).(func() (map[string]internal.OperationStats, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() map[string]internal.OperationStats); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]internal.OperationStats)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOperationStatsByPlanV2 provides a mock function with given fields:
func (_m *Operations) GetOperationStatsByPlanV2() ([]internal.OperationStatsV2, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetOperationStatsByPlanV2")
	}

	var r0 []internal.OperationStatsV2
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]internal.OperationStatsV2, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []internal.OperationStatsV2); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.OperationStatsV2)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOperationStatsForOrchestration provides a mock function with given fields: orchestrationID
func (_m *Operations) GetOperationStatsForOrchestration(orchestrationID string) (map[string]int, error) {
	ret := _m.Called(orchestrationID)

	if len(ret) == 0 {
		panic("no return value specified for GetOperationStatsForOrchestration")
	}

	var r0 map[string]int
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (map[string]int, error)); ok {
		return rf(orchestrationID)
	}
	if rf, ok := ret.Get(0).(func(string) map[string]int); ok {
		r0 = rf(orchestrationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]int)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(orchestrationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOperationsForIDs provides a mock function with given fields: operationIDList
func (_m *Operations) GetOperationsForIDs(operationIDList []string) ([]internal.Operation, error) {
	ret := _m.Called(operationIDList)

	if len(ret) == 0 {
		panic("no return value specified for GetOperationsForIDs")
	}

	var r0 []internal.Operation
	var r1 error
	if rf, ok := ret.Get(0).(func([]string) ([]internal.Operation, error)); ok {
		return rf(operationIDList)
	}
	if rf, ok := ret.Get(0).(func([]string) []internal.Operation); ok {
		r0 = rf(operationIDList)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.Operation)
		}
	}

	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(operationIDList)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProvisioningOperationByID provides a mock function with given fields: operationID
func (_m *Operations) GetProvisioningOperationByID(operationID string) (*internal.ProvisioningOperation, error) {
	ret := _m.Called(operationID)

	if len(ret) == 0 {
		panic("no return value specified for GetProvisioningOperationByID")
	}

	var r0 *internal.ProvisioningOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*internal.ProvisioningOperation, error)); ok {
		return rf(operationID)
	}
	if rf, ok := ret.Get(0).(func(string) *internal.ProvisioningOperation); ok {
		r0 = rf(operationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.ProvisioningOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(operationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProvisioningOperationByInstanceID provides a mock function with given fields: instanceID
func (_m *Operations) GetProvisioningOperationByInstanceID(instanceID string) (*internal.ProvisioningOperation, error) {
	ret := _m.Called(instanceID)

	if len(ret) == 0 {
		panic("no return value specified for GetProvisioningOperationByInstanceID")
	}

	var r0 *internal.ProvisioningOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*internal.ProvisioningOperation, error)); ok {
		return rf(instanceID)
	}
	if rf, ok := ret.Get(0).(func(string) *internal.ProvisioningOperation); ok {
		r0 = rf(instanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.ProvisioningOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(instanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUpdatingOperationByID provides a mock function with given fields: operationID
func (_m *Operations) GetUpdatingOperationByID(operationID string) (*internal.UpdatingOperation, error) {
	ret := _m.Called(operationID)

	if len(ret) == 0 {
		panic("no return value specified for GetUpdatingOperationByID")
	}

	var r0 *internal.UpdatingOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*internal.UpdatingOperation, error)); ok {
		return rf(operationID)
	}
	if rf, ok := ret.Get(0).(func(string) *internal.UpdatingOperation); ok {
		r0 = rf(operationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.UpdatingOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(operationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUpgradeClusterOperationByID provides a mock function with given fields: operationID
func (_m *Operations) GetUpgradeClusterOperationByID(operationID string) (*internal.UpgradeClusterOperation, error) {
	ret := _m.Called(operationID)

	if len(ret) == 0 {
		panic("no return value specified for GetUpgradeClusterOperationByID")
	}

	var r0 *internal.UpgradeClusterOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*internal.UpgradeClusterOperation, error)); ok {
		return rf(operationID)
	}
	if rf, ok := ret.Get(0).(func(string) *internal.UpgradeClusterOperation); ok {
		r0 = rf(operationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.UpgradeClusterOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(operationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertDeprovisioningOperation provides a mock function with given fields: operation
func (_m *Operations) InsertDeprovisioningOperation(operation internal.DeprovisioningOperation) error {
	ret := _m.Called(operation)

	if len(ret) == 0 {
		panic("no return value specified for InsertDeprovisioningOperation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(internal.DeprovisioningOperation) error); ok {
		r0 = rf(operation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertOperation provides a mock function with given fields: operation
func (_m *Operations) InsertOperation(operation internal.Operation) error {
	ret := _m.Called(operation)

	if len(ret) == 0 {
		panic("no return value specified for InsertOperation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(internal.Operation) error); ok {
		r0 = rf(operation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertProvisioningOperation provides a mock function with given fields: operation
func (_m *Operations) InsertProvisioningOperation(operation internal.ProvisioningOperation) error {
	ret := _m.Called(operation)

	if len(ret) == 0 {
		panic("no return value specified for InsertProvisioningOperation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(internal.ProvisioningOperation) error); ok {
		r0 = rf(operation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertUpdatingOperation provides a mock function with given fields: operation
func (_m *Operations) InsertUpdatingOperation(operation internal.UpdatingOperation) error {
	ret := _m.Called(operation)

	if len(ret) == 0 {
		panic("no return value specified for InsertUpdatingOperation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(internal.UpdatingOperation) error); ok {
		r0 = rf(operation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertUpgradeClusterOperation provides a mock function with given fields: operation
func (_m *Operations) InsertUpgradeClusterOperation(operation internal.UpgradeClusterOperation) error {
	ret := _m.Called(operation)

	if len(ret) == 0 {
		panic("no return value specified for InsertUpgradeClusterOperation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(internal.UpgradeClusterOperation) error); ok {
		r0 = rf(operation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListDeprovisioningOperations provides a mock function with given fields:
func (_m *Operations) ListDeprovisioningOperations() ([]internal.DeprovisioningOperation, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ListDeprovisioningOperations")
	}

	var r0 []internal.DeprovisioningOperation
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]internal.DeprovisioningOperation, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []internal.DeprovisioningOperation); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.DeprovisioningOperation)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDeprovisioningOperationsByInstanceID provides a mock function with given fields: instanceID
func (_m *Operations) ListDeprovisioningOperationsByInstanceID(instanceID string) ([]internal.DeprovisioningOperation, error) {
	ret := _m.Called(instanceID)

	if len(ret) == 0 {
		panic("no return value specified for ListDeprovisioningOperationsByInstanceID")
	}

	var r0 []internal.DeprovisioningOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]internal.DeprovisioningOperation, error)); ok {
		return rf(instanceID)
	}
	if rf, ok := ret.Get(0).(func(string) []internal.DeprovisioningOperation); ok {
		r0 = rf(instanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.DeprovisioningOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(instanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOperations provides a mock function with given fields: filter
func (_m *Operations) ListOperations(filter dbmodel.OperationFilter) ([]internal.Operation, int, int, error) {
	ret := _m.Called(filter)

	if len(ret) == 0 {
		panic("no return value specified for ListOperations")
	}

	var r0 []internal.Operation
	var r1 int
	var r2 int
	var r3 error
	if rf, ok := ret.Get(0).(func(dbmodel.OperationFilter) ([]internal.Operation, int, int, error)); ok {
		return rf(filter)
	}
	if rf, ok := ret.Get(0).(func(dbmodel.OperationFilter) []internal.Operation); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.Operation)
		}
	}

	if rf, ok := ret.Get(1).(func(dbmodel.OperationFilter) int); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(dbmodel.OperationFilter) int); ok {
		r2 = rf(filter)
	} else {
		r2 = ret.Get(2).(int)
	}

	if rf, ok := ret.Get(3).(func(dbmodel.OperationFilter) error); ok {
		r3 = rf(filter)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// ListOperationsByInstanceID provides a mock function with given fields: instanceID
func (_m *Operations) ListOperationsByInstanceID(instanceID string) ([]internal.Operation, error) {
	ret := _m.Called(instanceID)

	if len(ret) == 0 {
		panic("no return value specified for ListOperationsByInstanceID")
	}

	var r0 []internal.Operation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]internal.Operation, error)); ok {
		return rf(instanceID)
	}
	if rf, ok := ret.Get(0).(func(string) []internal.Operation); ok {
		r0 = rf(instanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.Operation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(instanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOperationsByInstanceIDGroupByType provides a mock function with given fields: instanceID
func (_m *Operations) ListOperationsByInstanceIDGroupByType(instanceID string) (*internal.GroupedOperations, error) {
	ret := _m.Called(instanceID)

	if len(ret) == 0 {
		panic("no return value specified for ListOperationsByInstanceIDGroupByType")
	}

	var r0 *internal.GroupedOperations
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*internal.GroupedOperations, error)); ok {
		return rf(instanceID)
	}
	if rf, ok := ret.Get(0).(func(string) *internal.GroupedOperations); ok {
		r0 = rf(instanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.GroupedOperations)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(instanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOperationsByOrchestrationID provides a mock function with given fields: orchestrationID, filter
func (_m *Operations) ListOperationsByOrchestrationID(orchestrationID string, filter dbmodel.OperationFilter) ([]internal.Operation, int, int, error) {
	ret := _m.Called(orchestrationID, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListOperationsByOrchestrationID")
	}

	var r0 []internal.Operation
	var r1 int
	var r2 int
	var r3 error
	if rf, ok := ret.Get(0).(func(string, dbmodel.OperationFilter) ([]internal.Operation, int, int, error)); ok {
		return rf(orchestrationID, filter)
	}
	if rf, ok := ret.Get(0).(func(string, dbmodel.OperationFilter) []internal.Operation); ok {
		r0 = rf(orchestrationID, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.Operation)
		}
	}

	if rf, ok := ret.Get(1).(func(string, dbmodel.OperationFilter) int); ok {
		r1 = rf(orchestrationID, filter)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(string, dbmodel.OperationFilter) int); ok {
		r2 = rf(orchestrationID, filter)
	} else {
		r2 = ret.Get(2).(int)
	}

	if rf, ok := ret.Get(3).(func(string, dbmodel.OperationFilter) error); ok {
		r3 = rf(orchestrationID, filter)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// ListOperationsInTimeRange provides a mock function with given fields: from, to
func (_m *Operations) ListOperationsInTimeRange(from time.Time, to time.Time) ([]internal.Operation, error) {
	ret := _m.Called(from, to)

	if len(ret) == 0 {
		panic("no return value specified for ListOperationsInTimeRange")
	}

	var r0 []internal.Operation
	var r1 error
	if rf, ok := ret.Get(0).(func(time.Time, time.Time) ([]internal.Operation, error)); ok {
		return rf(from, to)
	}
	if rf, ok := ret.Get(0).(func(time.Time, time.Time) []internal.Operation); ok {
		r0 = rf(from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.Operation)
		}
	}

	if rf, ok := ret.Get(1).(func(time.Time, time.Time) error); ok {
		r1 = rf(from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProvisioningOperationsByInstanceID provides a mock function with given fields: instanceID
func (_m *Operations) ListProvisioningOperationsByInstanceID(instanceID string) ([]internal.ProvisioningOperation, error) {
	ret := _m.Called(instanceID)

	if len(ret) == 0 {
		panic("no return value specified for ListProvisioningOperationsByInstanceID")
	}

	var r0 []internal.ProvisioningOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]internal.ProvisioningOperation, error)); ok {
		return rf(instanceID)
	}
	if rf, ok := ret.Get(0).(func(string) []internal.ProvisioningOperation); ok {
		r0 = rf(instanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.ProvisioningOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(instanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUpdatingOperationsByInstanceID provides a mock function with given fields: instanceID
func (_m *Operations) ListUpdatingOperationsByInstanceID(instanceID string) ([]internal.UpdatingOperation, error) {
	ret := _m.Called(instanceID)

	if len(ret) == 0 {
		panic("no return value specified for ListUpdatingOperationsByInstanceID")
	}

	var r0 []internal.UpdatingOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]internal.UpdatingOperation, error)); ok {
		return rf(instanceID)
	}
	if rf, ok := ret.Get(0).(func(string) []internal.UpdatingOperation); ok {
		r0 = rf(instanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.UpdatingOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(instanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUpgradeClusterOperationsByInstanceID provides a mock function with given fields: instanceID
func (_m *Operations) ListUpgradeClusterOperationsByInstanceID(instanceID string) ([]internal.UpgradeClusterOperation, error) {
	ret := _m.Called(instanceID)

	if len(ret) == 0 {
		panic("no return value specified for ListUpgradeClusterOperationsByInstanceID")
	}

	var r0 []internal.UpgradeClusterOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]internal.UpgradeClusterOperation, error)); ok {
		return rf(instanceID)
	}
	if rf, ok := ret.Get(0).(func(string) []internal.UpgradeClusterOperation); ok {
		r0 = rf(instanceID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.UpgradeClusterOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(instanceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUpgradeClusterOperationsByOrchestrationID provides a mock function with given fields: orchestrationID, filter
func (_m *Operations) ListUpgradeClusterOperationsByOrchestrationID(orchestrationID string, filter dbmodel.OperationFilter) ([]internal.UpgradeClusterOperation, int, int, error) {
	ret := _m.Called(orchestrationID, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListUpgradeClusterOperationsByOrchestrationID")
	}

	var r0 []internal.UpgradeClusterOperation
	var r1 int
	var r2 int
	var r3 error
	if rf, ok := ret.Get(0).(func(string, dbmodel.OperationFilter) ([]internal.UpgradeClusterOperation, int, int, error)); ok {
		return rf(orchestrationID, filter)
	}
	if rf, ok := ret.Get(0).(func(string, dbmodel.OperationFilter) []internal.UpgradeClusterOperation); ok {
		r0 = rf(orchestrationID, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]internal.UpgradeClusterOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(string, dbmodel.OperationFilter) int); ok {
		r1 = rf(orchestrationID, filter)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(string, dbmodel.OperationFilter) int); ok {
		r2 = rf(orchestrationID, filter)
	} else {
		r2 = ret.Get(2).(int)
	}

	if rf, ok := ret.Get(3).(func(string, dbmodel.OperationFilter) error); ok {
		r3 = rf(orchestrationID, filter)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// UpdateDeprovisioningOperation provides a mock function with given fields: operation
func (_m *Operations) UpdateDeprovisioningOperation(operation internal.DeprovisioningOperation) (*internal.DeprovisioningOperation, error) {
	ret := _m.Called(operation)

	if len(ret) == 0 {
		panic("no return value specified for UpdateDeprovisioningOperation")
	}

	var r0 *internal.DeprovisioningOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(internal.DeprovisioningOperation) (*internal.DeprovisioningOperation, error)); ok {
		return rf(operation)
	}
	if rf, ok := ret.Get(0).(func(internal.DeprovisioningOperation) *internal.DeprovisioningOperation); ok {
		r0 = rf(operation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.DeprovisioningOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(internal.DeprovisioningOperation) error); ok {
		r1 = rf(operation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOperation provides a mock function with given fields: operation
func (_m *Operations) UpdateOperation(operation internal.Operation) (*internal.Operation, error) {
	ret := _m.Called(operation)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOperation")
	}

	var r0 *internal.Operation
	var r1 error
	if rf, ok := ret.Get(0).(func(internal.Operation) (*internal.Operation, error)); ok {
		return rf(operation)
	}
	if rf, ok := ret.Get(0).(func(internal.Operation) *internal.Operation); ok {
		r0 = rf(operation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.Operation)
		}
	}

	if rf, ok := ret.Get(1).(func(internal.Operation) error); ok {
		r1 = rf(operation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProvisioningOperation provides a mock function with given fields: operation
func (_m *Operations) UpdateProvisioningOperation(operation internal.ProvisioningOperation) (*internal.ProvisioningOperation, error) {
	ret := _m.Called(operation)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProvisioningOperation")
	}

	var r0 *internal.ProvisioningOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(internal.ProvisioningOperation) (*internal.ProvisioningOperation, error)); ok {
		return rf(operation)
	}
	if rf, ok := ret.Get(0).(func(internal.ProvisioningOperation) *internal.ProvisioningOperation); ok {
		r0 = rf(operation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.ProvisioningOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(internal.ProvisioningOperation) error); ok {
		r1 = rf(operation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUpdatingOperation provides a mock function with given fields: operation
func (_m *Operations) UpdateUpdatingOperation(operation internal.UpdatingOperation) (*internal.UpdatingOperation, error) {
	ret := _m.Called(operation)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUpdatingOperation")
	}

	var r0 *internal.UpdatingOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(internal.UpdatingOperation) (*internal.UpdatingOperation, error)); ok {
		return rf(operation)
	}
	if rf, ok := ret.Get(0).(func(internal.UpdatingOperation) *internal.UpdatingOperation); ok {
		r0 = rf(operation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.UpdatingOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(internal.UpdatingOperation) error); ok {
		r1 = rf(operation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUpgradeClusterOperation provides a mock function with given fields: operation
func (_m *Operations) UpdateUpgradeClusterOperation(operation internal.UpgradeClusterOperation) (*internal.UpgradeClusterOperation, error) {
	ret := _m.Called(operation)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUpgradeClusterOperation")
	}

	var r0 *internal.UpgradeClusterOperation
	var r1 error
	if rf, ok := ret.Get(0).(func(internal.UpgradeClusterOperation) (*internal.UpgradeClusterOperation, error)); ok {
		return rf(operation)
	}
	if rf, ok := ret.Get(0).(func(internal.UpgradeClusterOperation) *internal.UpgradeClusterOperation); ok {
		r0 = rf(operation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal.UpgradeClusterOperation)
		}
	}

	if rf, ok := ret.Get(1).(func(internal.UpgradeClusterOperation) error); ok {
		r1 = rf(operation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOperations creates a new instance of Operations. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOperations(t interface {
	mock.TestingT
	Cleanup(func())
}) *Operations {
	mock := &Operations{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
