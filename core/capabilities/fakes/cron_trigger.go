package fakes

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron"
	crontriggermock "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/triggers/cron/cron_triggermock"
	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type cronCapability struct {
	doneCh chan bool
	stopCh chan bool
	lggr   logger.Logger
	*testutils.CapabilityWrapper
}

func (c *cronCapability) RegisterTrigger(ctx context.Context, request capabilities.TriggerRegistrationRequest) (<-chan capabilities.TriggerResponse, error) {
	ch := make(chan capabilities.TriggerResponse, 1)

	go func() {
		defer func() {
			c.lggr.Debug("shutting down cron capability")
			close(ch)
			close(c.doneCh)
		}()

		for {
			select {
			case <-time.After(3 * time.Second):
				response := capabilities.TriggerResponse{}

				trigger, err := c.InvokeTrigger(ctx, &sdkpb.TriggerSubscription{
					ExecId:  request.Metadata.WorkflowExecutionID,
					Id:      request.TriggerID,
					Payload: request.Payload,
					Method:  request.Method,
				})
				if err != nil {
					response.Err = err
				}

				if trigger == nil {
					return
				}

				response.Event = capabilities.TriggerEvent{
					TriggerType: request.TriggerID,
					ID:          uuid.NewString(),
					Payload:     trigger.Payload,
				}

				select {
				case ch <- response:
				case <-c.stopCh:
					return
				}
			case <-c.stopCh:
				return
			}
		}
	}()

	return ch, nil
}

type fakeCronTrigger struct {
	wrapped *cronCapability
}

func (f *fakeCronTrigger) Capability() capabilities.BaseCapability {
	return f.wrapped
}

func (f *fakeCronTrigger) Start(_ context.Context) error {
	return nil
}

func (f *fakeCronTrigger) Close() error {
	close(f.wrapped.stopCh)
	<-f.wrapped.doneCh
	return nil
}

func (f *fakeCronTrigger) HealthReport() map[string]error {
	return nil
}

func (f *fakeCronTrigger) Ready() error { return nil }

func (f *fakeCronTrigger) Name() string { return "fake-cron-trigger-server" }

func NewFakeCronTrigger(lggr logger.Logger) *fakeCronTrigger {
	capMock := &crontriggermock.CronCapability{}
	capMock.Trigger = func(ctx context.Context, input *cron.Config) (*cron.Payload, error) {
		return &cron.Payload{
			ScheduledExecutionTime: time.Now().String(),
		}, nil
	}
	stopCh := make(chan bool)
	doneCh := make(chan bool)
	return &fakeCronTrigger{
		wrapped: &cronCapability{
			stopCh: stopCh,
			doneCh: doneCh,
			lggr:   lggr,
			CapabilityWrapper: &testutils.CapabilityWrapper{
				Capability: capMock,
			},
		},
	}
}
