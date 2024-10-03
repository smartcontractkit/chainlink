package workflows

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/pb"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/proto"
	"log"
)

var registerTriggerFailureCounter metric.Int64Counter

func initMonitoringResources(_ context.Context) (err error) {
	registerTriggerFailureCounter, err = beholder.GetMeter().Int64Counter("RegisterTriggerFailure")
	if err != nil {
		return fmt.Errorf("failed to register trigger failure: %s", err)
	}
	return nil
}

// TODO: tradeoff between having labels be variadic vs values be variadic?
func sendLogAsCustomMessageWithLabels(ctx context.Context, labels map[string]string, format string, values ...interface{}) (err error) {
	for labelKey, labelValue := range labels {
		ctx, err = KeystoneContextWithLabel(ctx, labelKey, labelValue)
		if err != nil {
			return err
		}
	}

	return sendLogAsCustomMessage(ctx, format, values...)

}

func sendLogAsCustomMessage(ctx context.Context, format string, values ...interface{}) error {
	msg, err := composeLabeledMsg(ctx, format, values...)
	if err != nil {
		return fmt.Errorf("sendLogAsCustomMessage failed: %w", err)
	}

	labelsStruct, oerr := GetKeystoneLabelsFromContext(ctx)
	if oerr != nil {
		return oerr
	}

	labels := labelsStruct.ToMap()

	// Define a custom protobuf payload to emit
	payload := &pb.KeystoneCustomMessage{
		Msg:                 msg,
		WorkflowID:          labels[WorkflowID],
		WorkflowExecutionID: labels[WorkflowExecutionID],
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal protobuf")
	}

	err = beholder.GetEmitter().Emit(context.Background(), payloadBytes,
		"beholder_data_schema", "/keystone-custom-message/versions/1", // required
		"beholder_data_type", "custom_message",
	)
	if err != nil {
		log.Printf("Error emitting message: %v", err)
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
