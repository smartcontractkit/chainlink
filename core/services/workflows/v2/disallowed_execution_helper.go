package v2

import (
	"context"
	"errors"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	protoevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
)

type disallowedExecutionHelper struct {
	lggr        logger.Logger
	UserLogChan chan<- *protoevents.LogLine
	TimeProvider
	SecretsFetcher
}

func NewDisallowedExecutionHelper(lggr logger.Logger, userLogChan chan<- *protoevents.LogLine, timeProvider TimeProvider, secretsFetcher SecretsFetcher) *disallowedExecutionHelper {
	return &disallowedExecutionHelper{
		lggr:           lggr,
		UserLogChan:    userLogChan,
		TimeProvider:   timeProvider,
		SecretsFetcher: secretsFetcher,
	}
}

var _ host.ExecutionHelper = &disallowedExecutionHelper{}

func (d disallowedExecutionHelper) CallCapability(_ context.Context, _ *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	return nil, errors.New("capability calls cannot be made during this execution")
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
