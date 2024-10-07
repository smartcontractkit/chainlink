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

// Observability keys
const (
	cIDKey  = "capabilityID"
	ceIDKey = "capabilityExecutionID"
	tIDKey  = "triggerID"
	wIDKey  = "workflowID"
	eIDKey  = "workflowExecutionID"
	wnKey   = "workflowName"
	woIDKey = "workflowOwner"
	sIDKey  = "stepID"
	sRKey   = "stepRef"
)

var OrderedKeystoneLabels = []string{wIDKey, eIDKey}

var KeystoneLabelsMap = make(map[string]interface{})

func init() {
	for _, label := range OrderedKeystoneLabels {
		KeystoneLabelsMap[label] = interface{}(0)
	}
}

var registerTriggerFailureCounter metric.Int64Counter

func initMonitoringResources() (err error) {
	registerTriggerFailureCounter, err = beholder.GetMeter().Int64Counter("RegisterTriggerFailure")
	if err != nil {
		return fmt.Errorf("failed to register trigger failure: %s", err)
	}
	return nil
}

func sendLogAsCustomMessageF(labels map[string]string, format string, values ...interface{}) error {
	return sendLogAsCustomMessage(fmt.Sprintf(format, values...), labels)
}

func sendCtxLogAsCustomMessageF(ctx context.Context, format string, values ...interface{}) error {
	return sendCtxLogAsCustomMessage(ctx, fmt.Sprintf(format, values...))
}

func sendCtxLogAsCustomMessage(ctx context.Context, msg string) error {
	labelsStruct, oerr := GetKeystoneLabelsFromContext(ctx)
	if oerr != nil {
		return oerr
	}
	labels := labelsStruct.ToMap()

	return sendLogAsCustomMessageF(labels, msg)
}

func sendLogAsCustomMessage(msg string, labels map[string]string) error {
	// Define a custom protobuf payload to emit
	payload := &keystone.KeystoneCustomMessage{
		Msg: msg,
		Labels: map[string]*keystone.Value{
			wIDKey: {Value: &keystone.Value_StrValue{StrValue: labels[wIDKey]}},
			eIDKey: {Value: &keystone.Value_StrValue{StrValue: labels[eIDKey]}},
		},
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return fmt.Errorf("sending custom message failed to marshal protobuf: %s", err)
	}

	err = beholder.GetEmitter().Emit(context.Background(), payloadBytes,
		"beholder_data_schema", "/keystone-custom-message/versions/1", // required
		"beholder_data_type", "custom_message",
	)
	if err != nil {
		return fmt.Errorf("sending custom message failed on emit: %s", err)
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
