package v2

import (
	"context"
	"errors"
	"sync/atomic"
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
	SecretsFetcher

	// secretsCalled records whether GetSecrets was invoked at least once
	// through this helper, so callers (e.g. Subscribe) can report a metric.
	secretsCalled *atomic.Bool
}

func NewDisallowedExecutionHelper(lggr logger.Logger, userLogChan chan<- *protoevents.LogLine, timeProvider TimeProvider, secretsFetcher SecretsFetcher) *disallowedExecutionHelper {
	secretsCalled := &atomic.Bool{}
	return &disallowedExecutionHelper{
		lggr:           lggr,
		UserLogChan:    userLogChan,
		TimeProvider:   timeProvider,
		SecretsFetcher: &secretsCallTrackingFetcher{SecretsFetcher: secretsFetcher, called: secretsCalled},
		secretsCalled:  secretsCalled,
	}
}

// SecretsCalled reports whether GetSecrets was invoked at least once through
// this helper.
func (d disallowedExecutionHelper) SecretsCalled() bool {
	return d.secretsCalled.Load()
}

// secretsCallTrackingFetcher wraps a SecretsFetcher to record whether it was
// ever called, without altering its behavior.
type secretsCallTrackingFetcher struct {
	SecretsFetcher
	called *atomic.Bool
}

func (f *secretsCallTrackingFetcher) GetSecrets(ctx context.Context, request *sdkpb.GetSecretsRequest) ([]*sdkpb.SecretResponse, error) {
	f.called.Store(true)
	return f.SecretsFetcher.GetSecrets(ctx, request)
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

func (d disallowedExecutionHelper) EmitUserMetric(_ context.Context, _ *eventsv2.WorkflowUserMetric) error {
	return errors.New("metric emission is not allowed during this execution")
}
