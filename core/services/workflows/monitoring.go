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

type customMessageAgent struct {
	labels map[string]string
}

func NewCustomMessageAgent() customMessageAgent {
	return customMessageAgent{labels: make(map[string]string)}
}

// with adds multiple key-value pairs to the customMessageAgent for transmission with sendLogAsCustomMessage
func (c customMessageAgent) with(keyValues ...string) customMessageAgent {
	newCustomMessageAgent := NewCustomMessageAgent()

	if len(keyValues)%2 != 0 {
		// If an odd number of key-value arguments is passed, return the original customMessageAgent unchanged
		return c
	}

	// Copy existing labels from the current agent
	for k, v := range c.labels {
		newCustomMessageAgent.labels[k] = v
	}

	// Add new key-value pairs
	for i := 0; i < len(keyValues); i += 2 {
		key := keyValues[i]
		value := keyValues[i+1]
		newCustomMessageAgent.labels[key] = value
	}

	return newCustomMessageAgent
}

// sendLogAsCustomMessage emits a KeystoneCustomMessage with msg and labels as data.
// any key in labels that is not part of orderedLabelKeys will not be transmitted
func (c customMessageAgent) sendLogAsCustomMessage(msg string) error {
	return sendLogAsCustomMessageW(msg, c.labels)
}

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

var orderedLabelKeys = []string{sRKey, sIDKey, tIDKey, cIDKey, eIDKey, wIDKey}

var labelsMap = make(map[string]interface{})

func init() {
	for _, label := range orderedLabelKeys {
		labelsMap[label] = interface{}(0)
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

// sendLogAsCustomMessageF formats into a msg to be consumed by sendLogAsCustomMessageW
func sendLogAsCustomMessageF(labels map[string]string, format string, values ...interface{}) error {
	return sendLogAsCustomMessageW(fmt.Sprintf(format, values...), labels)
}

// sendCtxLogAsCustomMessageF formats into a msg to be consumed by sendCtxLogAsCustomMessage
func sendCtxLogAsCustomMessageF(ctx context.Context, format string, values ...interface{}) error {
	return sendCtxLogAsCustomMessage(ctx, fmt.Sprintf(format, values...))
}

// sendCtxLogAsCustomMessage emits a KeystoneCustomMessage with msg and labels extracted from ctx as data.
func sendCtxLogAsCustomMessage(ctx context.Context, msg string) error {
	labelsStruct, oerr := GetKeystoneLabelsFromContext(ctx)
	if oerr != nil {
		return oerr
	}
	labels := labelsStruct.ToMap()

	return sendLogAsCustomMessageF(labels, msg)
}

// sendLogAsCustomMessageV allows the consumer to pass in variable number of label key value pairs
func sendLogAsCustomMessageV(msg string, labelKVs ...string) error {
	if len(labelKVs)%2 != 0 {
		return fmt.Errorf("labelKVs must be provided in key-value pairs")
	}

	labels := make(map[string]string)
	for i := 0; i < len(labelKVs); i += 2 {
		key := labelKVs[i]
		value := labelKVs[i+1]
		labels[key] = value
	}

	return sendLogAsCustomMessageF(labels, msg)
}

func sendLogAsCustomMessageW(msg string, labels map[string]string) error {
	protoLabels := make(map[string]*keystone.Value)
	for _, k := range orderedLabelKeys {
		if _, ok := labels[k]; ok {
			protoLabels[k] = &keystone.Value{Value: &keystone.Value_StrValue{StrValue: labels[k]}}
		}
	}
	// Define a custom protobuf payload to emit
	payload := &keystone.KeystoneCustomMessage{
		Msg:    msg,
		Labels: protoLabels,
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
		// lggr.with for logs, metric.WithAttributes, and set the labels directly in the proto for custom messages
		lggr.Errorf("failed to get otel attributes from context: %s", oerr)
	}

	labels = append(labels, ctxLabels...)
	// TODO add duplicate labels check?
	registerTriggerFailureCounter.Add(ctx, 1, metric.WithAttributes(labels...))
}
