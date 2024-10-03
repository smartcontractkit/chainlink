package workflows

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var registerTriggerFailureCounter metric.Int64Counter

func initMonitoringResources(_ context.Context) (err error) {
	registerTriggerFailureCounter, err = beholder.GetMeter().Int64Counter("RegisterTriggerFailure")
	if err != nil {
		return fmt.Errorf("failed to register trigger failure: %s", err)
	}
	return nil
}

func incrementRegisterTriggerFailureCounter(ctx context.Context, lggr logger.Logger, labels ...attribute.KeyValue) {
	ctxLabels, oerr := getOtelAttributesFromCtx(ctx)
	if oerr != nil {
		// custom messages require this extracting of values from the context
		// but if i set them in the proto, then could use
		// lggr.With for logs, metric.WithAttributes, and set the labels directly in the proto for custom messages
		lggr.Errorf("failed to get otel attributes from context: %s", oerr)
	}

	labels = append(labels, ctxLabels...)
	// TODO add duplicate labels check?
	registerTriggerFailureCounter.Add(ctx, 1, metric.WithAttributes(labels...))
}
