package deprovisioning

import (
	"fmt"
	"testing"
	"time"

	"github.com/kyma-project/kyma-environment-broker/common/hyperscaler"
	hyperscalerMocks "github.com/kyma-project/kyma-environment-broker/common/hyperscaler/automock"
	pkg "github.com/kyma-project/kyma-environment-broker/common/runtime"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/broker"
	"github.com/kyma-project/kyma-environment-broker/internal/fixture"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/pivotal-cf/brokerapi/v12/domain"
	"github.com/stretchr/testify/assert"
)

func TestReleaseSubscriptionStep_HappyPath(t *testing.T) {
	// given
	memoryStorage := storage.NewMemoryStorage()

	operation := fixDeprovisioningOperationWithPlanID(broker.GCPPlanID)
	instance := fixGCPInstance(operation.InstanceID)

	err := memoryStorage.Instances().Insert(instance)
	assert.NoError(t, err)

	accountProviderMock := &hyperscalerMocks.AccountProvider{}
	accountProviderMock.On("MarkUnusedGardenerSecretBindingAsDirty", hyperscaler.GCP("westeurope"), instance.GetSubscriptionGlobalAccoundID(), false).Return(nil)

	step := NewReleaseSubscriptionStep(memoryStorage, accountProviderMock)

	// when
	operation, repeat, err := step.Run(operation, fixLogger())

	assert.NoError(t, err)

	// then
	accountProviderMock.AssertNumberOfCalls(t, "MarkUnusedGardenerSecretBindingAsDirty", 1)
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), repeat)
	assert.Equal(t, domain.Succeeded, operation.State)
}

func TestReleaseSubscriptionStep_TrialPlan(t *testing.T) {
	// given
	memoryStorage := storage.NewMemoryStorage()

	operation := fixDeprovisioningOperationWithPlanID(broker.TrialPlanID)
	instance := fixGCPInstance(operation.InstanceID)

	err := memoryStorage.Instances().Insert(instance)
	assert.NoError(t, err)

	accountProviderMock := &hyperscalerMocks.AccountProvider{}
	accountProviderMock.On("MarkUnusedGardenerSecretBindingAsDirty", hyperscaler.GCP("westeurope"), instance.GetSubscriptionGlobalAccoundID(), false).Return(nil)

	step := NewReleaseSubscriptionStep(memoryStorage, accountProviderMock)

	// when
	operation, repeat, err := step.Run(operation, fixLogger())

	assert.NoError(t, err)

	// then
	accountProviderMock.AssertNumberOfCalls(t, "MarkUnusedGardenerSecretBindingAsDirty", 0)
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), repeat)
	assert.Equal(t, domain.Succeeded, operation.State)
}

func TestReleaseSubscriptionStep_OwnClusterPlan(t *testing.T) {
	// given
	memoryStorage := storage.NewMemoryStorage()

	operation := fixDeprovisioningOperationWithPlanID(broker.OwnClusterPlanID)
	instance := fixGCPInstance(operation.InstanceID)

	err := memoryStorage.Instances().Insert(instance)
	assert.NoError(t, err)

	accountProviderMock := &hyperscalerMocks.AccountProvider{}
	accountProviderMock.On("MarkUnusedGardenerSecretBindingAsDirty", hyperscaler.GCP("westeurope"), instance.GetSubscriptionGlobalAccoundID(), false).Return(nil)

	step := NewReleaseSubscriptionStep(memoryStorage, accountProviderMock)

	// when
	operation, repeat, err := step.Run(operation, fixLogger())

	assert.NoError(t, err)

	// then
	accountProviderMock.AssertNumberOfCalls(t, "MarkUnusedGardenerSecretBindingAsDirty", 0)
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), repeat)
	assert.Equal(t, domain.Succeeded, operation.State)
}

func TestReleaseSubscriptionStep_InstanceNotFound(t *testing.T) {
	// given
	memoryStorage := storage.NewMemoryStorage()

	operation := fixDeprovisioningOperationWithPlanID(broker.GCPPlanID)
	instance := fixGCPInstance(operation.InstanceID)

	accountProviderMock := &hyperscalerMocks.AccountProvider{}
	accountProviderMock.On("MarkUnusedGardenerSecretBindingAsDirty", hyperscaler.GCP("westeurope"), instance.GetSubscriptionGlobalAccoundID(), false).Return(nil)

	step := NewReleaseSubscriptionStep(memoryStorage, accountProviderMock)
	err := memoryStorage.Operations().InsertOperation(operation)
	assert.NoError(t, err)

	// when
	operation, repeat, err := step.Run(operation, fixLogger())

	assert.NoError(t, err)

	// then
	accountProviderMock.AssertNotCalled(t, "MarkUnusedGardenerSecretBindingAsDirty")
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), repeat)
	assert.Equal(t, domain.Succeeded, operation.State)
}

func TestReleaseSubscriptionStep_ProviderNotFound(t *testing.T) {
	// given
	memoryStorage := storage.NewMemoryStorage()

	operation := fixDeprovisioningOperationWithPlanID(broker.GCPPlanID)
	instance := fixGCPInstance(operation.InstanceID)
	instance.Provider = "unknown"

	err := memoryStorage.Instances().Insert(instance)
	assert.NoError(t, err)

	accountProviderMock := &hyperscalerMocks.AccountProvider{}
	accountProviderMock.On("MarkUnusedGardenerSecretBindingAsDirty", hyperscaler.GCP("westeurope"), instance.GetSubscriptionGlobalAccoundID(), false).Return(nil)

	step := NewReleaseSubscriptionStep(memoryStorage, accountProviderMock)
	err = memoryStorage.Operations().InsertOperation(operation)
	assert.NoError(t, err)

	// when
	operation, repeat, err := step.Run(operation, fixLogger())

	assert.NoError(t, err)

	// then
	accountProviderMock.AssertNotCalled(t, "MarkUnusedGardenerSecretBindingAsDirty")
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), repeat)
	assert.Equal(t, domain.Succeeded, operation.State)
}

func TestReleaseSubscriptionStepGardener_CallFails(t *testing.T) {
	// given
	memoryStorage := storage.NewMemoryStorage()

	operation := fixDeprovisioningOperationWithPlanID(broker.GCPPlanID)
	instance := fixGCPInstance(operation.InstanceID)

	err := memoryStorage.Instances().Insert(instance)
	assert.NoError(t, err)

	accountProviderMock := &hyperscalerMocks.AccountProvider{}
	accountProviderMock.On("MarkUnusedGardenerSecretBindingAsDirty", hyperscaler.GCP("westeurope"), instance.GetSubscriptionGlobalAccoundID(), false).Return(fmt.Errorf("failed to release subscription for tenant. Gardener Account pool is not configured"))

	step := NewReleaseSubscriptionStep(memoryStorage, accountProviderMock)

	// when
	operation, repeat, err := step.Run(operation, fixLogger())

	assert.NoError(t, err)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 10*time.Second, repeat)
	assert.Equal(t, domain.Succeeded, operation.State)
}

func fixGCPInstance(instanceID string) internal.Instance {
	instance := fixture.FixInstance(instanceID)
	instance.Provider = pkg.GCP
	return instance
}
func fixDeprovisioningOperationWithPlanID(planID string) internal.Operation {
	deprovisioningOperation := fixture.FixDeprovisioningOperationAsOperation(testOperationID, testInstanceID)
	deprovisioningOperation.ProvisioningParameters.PlanID = planID
	deprovisioningOperation.ProvisioningParameters.ErsContext.GlobalAccountID = testGlobalAccountID
	deprovisioningOperation.ProvisioningParameters.ErsContext.SubAccountID = testSubAccountID
	return deprovisioningOperation
}
