package process

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/dberr"
	"github.com/pivotal-cf/brokerapi/v12/domain"
)

type UpgradeClusterOperationManager struct {
	storage storage.UpgradeCluster
}

func NewUpgradeClusterOperationManager(storage storage.Operations) *UpgradeClusterOperationManager {
	return &UpgradeClusterOperationManager{storage: storage}
}

// OperationSucceeded marks the operation as succeeded and only repeats it if there is a storage error
func (om *UpgradeClusterOperationManager) OperationSucceeded(operation internal.UpgradeClusterOperation, description string, log *slog.Logger) (internal.UpgradeClusterOperation, time.Duration, error) {
	updatedOperation, repeat, _ := om.update(operation, internal.OperationStateSucceeded, description, log)
	// repeat in case of storage error
	if repeat != 0 {
		return updatedOperation, repeat, nil
	}

	return updatedOperation, 0, nil
}

// OperationFailed marks the operation as failed and only repeats it if there is a storage error
func (om *UpgradeClusterOperationManager) OperationFailed(operation internal.UpgradeClusterOperation, description string, err error, log *slog.Logger) (internal.UpgradeClusterOperation, time.Duration, error) {
	updatedOperation, repeat, _ := om.update(operation, internal.OperationStateFailed, description, log)
	// repeat in case of storage error
	if repeat != 0 {
		return updatedOperation, repeat, nil
	}

	var retErr error
	if err == nil {
		// no exact err passed in
		retErr = fmt.Errorf("%s", description)
	} else {
		// keep the original err object for error categorizer
		retErr = fmt.Errorf("%s: %w", description, err)
	}

	return updatedOperation, 0, retErr
}

// OperationSucceeded marks the operation as succeeded and only repeats it if there is a storage error
func (om *UpgradeClusterOperationManager) OperationCanceled(operation internal.UpgradeClusterOperation, description string, log *slog.Logger) (internal.UpgradeClusterOperation, time.Duration, error) {
	updatedOperation, repeat, _ := om.update(operation, internal.OperationStateCanceled, description, log)
	if repeat != 0 {
		return updatedOperation, repeat, nil
	}

	return updatedOperation, 0, nil
}

// RetryOperation retries an operation for at maxTime in retryInterval steps and fails the operation if retrying failed
func (om *UpgradeClusterOperationManager) RetryOperation(operation internal.UpgradeClusterOperation, errorMessage string, err error, retryInterval time.Duration, maxTime time.Duration, log *slog.Logger) (internal.UpgradeClusterOperation, time.Duration, error) {
	since := time.Since(operation.UpdatedAt)

	log.Info(fmt.Sprintf("Retry Operation was triggered with message: %s", errorMessage))
	log.Info(fmt.Sprintf("Retrying for %s in %s steps", maxTime.String(), retryInterval.String()))
	if since < maxTime {
		return operation, retryInterval, nil
	}
	log.Error(fmt.Sprintf("Aborting after %s of failing retries", maxTime.String()))
	return om.OperationFailed(operation, errorMessage, err, log)
}

// UpdateOperation updates a given operation
func (om *UpgradeClusterOperationManager) UpdateOperation(operation internal.UpgradeClusterOperation, update func(operation *internal.UpgradeClusterOperation), log *slog.Logger) (internal.UpgradeClusterOperation, time.Duration, error) {
	update(&operation)
	updatedOperation, err := om.storage.UpdateUpgradeClusterOperation(operation)
	switch {
	case dberr.IsConflict(err):
		{
			op, err := om.storage.GetUpgradeClusterOperationByID(operation.Operation.ID)
			if err != nil {
				log.Error(fmt.Sprintf("while getting operation: %v", err))
				return operation, 1 * time.Minute, err
			}
			op.Merge(&operation.Operation)
			update(op)
			updatedOperation, err = om.storage.UpdateUpgradeClusterOperation(*op)
			if err != nil {
				log.Error(fmt.Sprintf("while updating operation after conflict: %v", err))
				return operation, 1 * time.Minute, err
			}
		}
	case err != nil:
		log.Error(fmt.Sprintf("while updating operation: %v", err))
		return operation, 1 * time.Minute, err
	}
	return *updatedOperation, 0, nil
}

// Deprecated: SimpleUpdateOperation updates a given operation without handling conflicts. Should be used when operation's data mutations are not clear
func (om *UpgradeClusterOperationManager) SimpleUpdateOperation(operation internal.UpgradeClusterOperation) (internal.UpgradeClusterOperation, time.Duration) {
	updatedOperation, err := om.storage.UpdateUpgradeClusterOperation(operation)
	if err != nil {
		slog.With("instanceID", operation.InstanceID).
			Error(fmt.Sprintf("Update upgradeCluster operation failed: %s", err.Error()))
		return operation, 1 * time.Minute
	}
	return *updatedOperation, 0
}

// RetryOperationWithoutFail retries an operation for at maxTime in retryInterval steps and omits the operation if retrying failed
func (om *UpgradeClusterOperationManager) RetryOperationWithoutFail(operation internal.UpgradeClusterOperation, description string, retryInterval, maxTime time.Duration, log *slog.Logger) (internal.UpgradeClusterOperation, time.Duration, error) {
	since := time.Since(operation.UpdatedAt)

	log.Info(fmt.Sprintf("Retry Operation was triggered with message: %s", description))
	log.Info(fmt.Sprintf("Retrying for %s in %s steps", maxTime.String(), retryInterval.String()))
	if since < maxTime {
		return operation, retryInterval, nil
	}
	// update description to track failed steps
	updatedOperation, repeat, _ := om.update(operation, domain.InProgress, description, log)
	if repeat != 0 {
		return updatedOperation, repeat, nil
	}

	log.Error(fmt.Sprintf("Omitting after %s of failing retries", maxTime.String()))
	return updatedOperation, 0, nil
}

func (om *UpgradeClusterOperationManager) update(operation internal.UpgradeClusterOperation, state domain.LastOperationState, description string, log *slog.Logger) (internal.UpgradeClusterOperation, time.Duration, error) {
	return om.UpdateOperation(operation, func(operation *internal.UpgradeClusterOperation) {
		operation.State = state
		operation.Description = description
	}, log)
}
