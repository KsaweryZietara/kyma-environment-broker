package broker

import (
	"context"
	"fmt"

	"github.com/pivotal-cf/brokerapi/v11/domain"
	"github.com/sirupsen/logrus"
)

type LastBindingOperationEndpoint struct {
	log logrus.FieldLogger
}

func NewLastBindingOperation(log logrus.FieldLogger) *LastBindingOperationEndpoint {
	return &LastBindingOperationEndpoint{log: log.WithField("service", "LastBindingOperationEndpoint")}
}

// LastBindingOperation fetches last operation state for a service binding
//
//	GET /v2/service_instances/{instance_id}/service_bindings/{binding_id}/last_operation
func (b *LastBindingOperationEndpoint) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details domain.PollDetails) (domain.LastOperation, error) {
	b.log.Infof("LastBindingOperation instanceID: %s", instanceID)
	b.log.Infof("LastBindingOperation bindingID: %s", bindingID)
	b.log.Infof("LastBindingOperation details: %+v", details)

	return domain.LastOperation{}, fmt.Errorf("not supported")
}
