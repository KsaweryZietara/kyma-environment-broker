package broker

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/broker/automock"
	"github.com/kyma-project/kyma-environment-broker/internal/fixture"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/pivotal-cf/brokerapi/v12/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	instanceID  = "instance-001"
	planID      = "awesome-plan"
	operationID = "1234"
)

func TestDeprovisionEndpoint_DeprovisionNotExistingInstance(t *testing.T) {
	// given
	memoryStorage := storage.NewMemoryStorage()
	queue := &automock.Queue{}
	queue.On("Add", mock.AnythingOfType("string"))

	svc := NewDeprovision(memoryStorage.Instances(), memoryStorage.Operations(), queue, fixLogger())

	// when
	_, err := svc.Deprovision(context.TODO(), "inst-0001", domain.DeprovisionDetails{}, true)

	// then
	assert.NoError(t, err)
}

func TestDeprovisionEndpoint_DeprovisionExistingInstance(t *testing.T) {
	// given
	memoryStorage := storage.NewMemoryStorage()
	err := memoryStorage.Instances().Insert(fixInstance())
	require.NoError(t, err)

	queue := &automock.Queue{}
	queue.On("Add", mock.AnythingOfType("string"))

	svc := NewDeprovision(memoryStorage.Instances(), memoryStorage.Operations(), queue, fixLogger())

	// when
	_, err = svc.Deprovision(context.TODO(), instanceID, domain.DeprovisionDetails{}, true)

	// then
	require.NoError(t, err)
	operation, err := memoryStorage.Operations().GetDeprovisioningOperationByInstanceID(instanceID)
	require.NoError(t, err)
	assert.Equal(t, domain.LastOperationState("pending"), operation.State)
}

func TestDeprovisionEndpoint_DeprovisionExistingOperationInProgress(t *testing.T) {
	// given
	memoryStorage := storage.NewMemoryStorage()
	err := memoryStorage.Instances().Insert(fixInstance())
	require.NoError(t, err)

	err = memoryStorage.Operations().InsertDeprovisioningOperation(fixDeprovisioningOperation(domain.InProgress))
	require.NoError(t, err)

	queue := &automock.Queue{}
	queue.On("Add", mock.AnythingOfType("string"))

	svc := NewDeprovision(memoryStorage.Instances(), memoryStorage.Operations(), queue, fixLogger())

	// when
	res, err := svc.Deprovision(context.TODO(), instanceID, domain.DeprovisionDetails{}, true)

	// then
	require.NoError(t, err)
	assert.Equal(t, operationID, res.OperationData)

	operation, err := memoryStorage.Operations().GetDeprovisioningOperationByInstanceID(instanceID)
	require.NoError(t, err)
	assert.Equal(t, domain.InProgress, operation.State)
	assert.Equal(t, "", operation.ProvisionerOperationID)
}

func TestDeprovisionEndpoint_DeprovisionExistingOperationFailed(t *testing.T) {
	// given
	memoryStorage := storage.NewMemoryStorage()
	err := memoryStorage.Instances().Insert(fixInstance())
	require.NoError(t, err)

	err = memoryStorage.Operations().InsertDeprovisioningOperation(fixDeprovisioningOperation(domain.Failed))
	require.NoError(t, err)

	queue := &automock.Queue{}
	queue.On("Add", mock.Anything)

	svc := NewDeprovision(memoryStorage.Instances(), memoryStorage.Operations(), queue, fixLogger())

	// when
	res, err := svc.Deprovision(context.TODO(), instanceID, domain.DeprovisionDetails{}, true)

	// then
	require.NoError(t, err)
	assert.NotEqual(t, operationID, res.OperationData)

	operation, err := memoryStorage.Operations().GetDeprovisioningOperationByInstanceID(instanceID)
	require.NoError(t, err)
	assert.Equal(t, domain.LastOperationState("pending"), operation.State)
}

func fixDeprovisioningOperation(state domain.LastOperationState) internal.DeprovisioningOperation {
	deprovisioningOperation := fixture.FixDeprovisioningOperation(operationID, instanceID)
	deprovisioningOperation.State = state

	return deprovisioningOperation
}

func fixInstance() internal.Instance {
	instance := fixture.FixInstance(instanceID)
	instance.ServicePlanID = planID

	return instance
}

func fixLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
