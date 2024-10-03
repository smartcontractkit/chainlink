package workflows

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/pb/keystone"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/proto"
)

var registerTriggerFailureCounter metric.Int64Counter

func initMonitoringResources() (err error) {
	registerTriggerFailureCounter, err = beholder.GetMeter().Int64Counter("RegisterTriggerFailure")
	if err != nil {
		return fmt.Errorf("failed to register trigger failure: %s", err)
	}
	return nil
}

func sendLogAsCustomMessageF(ctx context.Context, lggr logger.Logger, format string, values ...interface{}) error {
	return sendLogAsCustomMessage(ctx, lggr, fmt.Sprintf(format, values...))
}

func sendLogAsCustomMessage(ctx context.Context, lggr logger.Logger, msg string) error {
	labelsStruct, oerr := GetKeystoneLabelsFromContext(ctx)
	if oerr != nil {
		return oerr
	}
	labels := labelsStruct.ToMap()

	// Define a custom protobuf payload to emit
	payload := &keystone.KeystoneCustomMessage{
		Msg: msg,
		Labels: map[string]*keystone.Value{
			WorkflowID:          {Value: &keystone.Value_StrValue{StrValue: labels[WorkflowID]}},
			WorkflowExecutionID: {Value: &keystone.Value_StrValue{StrValue: labels[WorkflowExecutionID]}},
		},
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		lggr.Fatalf("Failed to marshal protobuf: %s", err)
	}

	err = beholder.GetEmitter().Emit(context.Background(), payloadBytes,
		"beholder_data_schema", "/keystone-custom-message/versions/1", // required
		"beholder_data_type", "custom_message",
	)
	if err != nil {
		lggr.Criticalf("Error emitting message: %s", err)
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
