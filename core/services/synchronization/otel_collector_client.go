package synchronization

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
)

type otelCollectorClient struct {
	clientCtx      context.Context
	clientCancelFn context.CancelCauseFunc

	client *beholder.Client

	pipe chan telemetryWrapper
}

type telemetryWrapper struct {
	telemetry  []byte
	contractID string
	telemType  TelemetryType
}

func NewOpenTelemetryClient(cfg beholder.Config) (TelemetryService, error) {
	ctx, cancel := context.WithCancelCause(context.TODO())
	client, err := beholder.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	pipe := make(chan telemetryWrapper)
	return &otelCollectorClient{ctx, cancel, client, pipe}, nil
}

func (o *otelCollectorClient) Start(localCtx context.Context) error {
	defer o.client.Close()
	go func() {
		for {
			select {
			case wrapper := <-o.pipe:
				err := o.client.Emitter.Emit(localCtx, wrapper.telemetry, wrapper.contractID, wrapper.telemType)
				if err != nil {
					o.clientCancelFn(err)
				}
			case <-o.clientCtx.Done():
			case <-localCtx.Done():
				if err := localCtx.Err(); err != nil {
					o.clientCancelFn(err)
				}
			}
		}
	}()
	return nil
}

var errCancelledByClose = fmt.Errorf("OTel client terminated by Close()")

func (o *otelCollectorClient) Close() error {
	o.clientCancelFn(errCancelledByClose)
	return nil
}

func (o *otelCollectorClient) Send(localCtx context.Context, telemetry []byte, contractID string, telemType TelemetryType) {
	select {
	case o.pipe <- telemetryWrapper{telemetry, contractID, telemType}:
	case <-localCtx.Done():
	case <-o.clientCtx.Done():
		if err := localCtx.Err(); err != nil {
			o.clientCancelFn(err)
		}
	}
}

func (o *otelCollectorClient) Ready() error {
	return o.clientCtx.Err()
}

func (o *otelCollectorClient) HealthReport() map[string]error {
	return map[string]error{
		o.Name(): o.clientCtx.Err(),
	}
}

func (o *otelCollectorClient) Name() string {
	return "otel-clinet"
}
