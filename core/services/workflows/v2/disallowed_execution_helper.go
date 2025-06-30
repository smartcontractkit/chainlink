package v2

import (
	"context"
	"errors"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	protoevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
)

type DisallowedExecutionHelper struct {
	lggr        logger.Logger
	UserLogChan chan<- *protoevents.LogLine
	TimeProvider
	SecretsFetcher
}

var _ host.ExecutionHelper = &DisallowedExecutionHelper{}

func (d DisallowedExecutionHelper) CallCapability(_ context.Context, _ *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	return nil, errors.New("capability calls cannot be made during this execution")
}

func (d DisallowedExecutionHelper) GetWorkflowExecutionID() string {
	return ""
}

func (d DisallowedExecutionHelper) EmitUserLog(msg string) error {
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
