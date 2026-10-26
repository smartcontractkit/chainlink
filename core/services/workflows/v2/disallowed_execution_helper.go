package v2

import (
	"context"
	"errors"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	protoevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
)

type disallowedExecutionHelper struct {
	lggr        logger.Logger
	UserLogChan chan<- *protoevents.LogLine
	TimeProvider
}

func NewDisallowedExecutionHelper(lggr logger.Logger, userLogChan chan<- *protoevents.LogLine, timeProvider TimeProvider) *disallowedExecutionHelper {
	return &disallowedExecutionHelper{
		lggr:         lggr,
		UserLogChan:  userLogChan,
		TimeProvider: timeProvider,
	}
}

var _ host.ExecutionHelper = &disallowedExecutionHelper{}

var (
	// ErrCapabilityCallDuringSubscription is returned to a guest that attempts
	// a capability call during the trigger subscription phase. Capability
	// calls are only allowed during workflow executions.
	ErrCapabilityCallDuringSubscription = errors.New("capability calls cannot be made during trigger subscription")
	// ErrSecretsCallDuringSubscription is returned to a guest that attempts a
	// secrets call during the trigger subscription phase. Secrets calls are
	// only allowed during workflow executions.
	ErrSecretsCallDuringSubscription = errors.New("secrets calls cannot be made during trigger subscription")
)

func (d disallowedExecutionHelper) CallCapability(_ context.Context, _ *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	return nil, ErrCapabilityCallDuringSubscription
}

func (d disallowedExecutionHelper) GetSecrets(_ context.Context, _ *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	return nil, ErrSecretsCallDuringSubscription
}

func (d disallowedExecutionHelper) GetWorkflowExecutionID() string {
	return ""
}

func (d disallowedExecutionHelper) EmitUserLog(msg string) error {
	select {
	case d.UserLogChan <- &protoevents.LogLine{
		NodeTimestamp: time.Now().Format(time.RFC3339Nano),
		Message:       msg,
	}:
		// Successfully sent to channel
	default:
		d.lggr.Warnw("Exceeded max allowed user log messages, dropping")
	}
	return nil
}

func (d disallowedExecutionHelper) EmitUserMetric(_ context.Context, _ *eventsv2.WorkflowUserMetric) error {
	return errors.New("metric emission is not allowed during this execution")
}
