package workflows

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder/pb/keystone"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/proto"
)

type customMessageLabeler struct {
	labels map[string]string
}

func NewCustomMessageLabeler() customMessageLabeler {
	return customMessageLabeler{labels: make(map[string]string)}
}

// with adds multiple key-value pairs to the customMessageLabeler for transmission with sendLogAsCustomMessage
func (c customMessageLabeler) with(keyValues ...string) customMessageLabeler {
	newCustomMessageLabeler := NewCustomMessageLabeler()

	if len(keyValues)%2 != 0 {
		// If an odd number of key-value arguments is passed, return the original customMessageLabeler unchanged
		return c
	}

	// Copy existing labels from the current agent
	for k, v := range c.labels {
		newCustomMessageLabeler.labels[k] = v
	}

	// Add new key-value pairs
	for i := 0; i < len(keyValues); i += 2 {
		key := keyValues[i]
		value := keyValues[i+1]
		newCustomMessageLabeler.labels[key] = value
	}

	return newCustomMessageLabeler
}

// sendLogAsCustomMessage emits a KeystoneCustomMessage with msg and labels as data.
// any key in labels that is not part of orderedLabelKeys will not be transmitted
func (c customMessageLabeler) sendLogAsCustomMessage(msg string) error {
	return sendLogAsCustomMessageW(msg, c.labels)
}

type customMetricsLabeler struct {
	labels map[string]string
}

func NewCustomMetricsLabeler() customMetricsLabeler {
	return customMetricsLabeler{labels: make(map[string]string)}
}

// with adds multiple key-value pairs to the customMessageLabeler for transmission with sendLogAsCustomMessage
func (c customMetricsLabeler) with(keyValues ...string) customMetricsLabeler {
	newCustomMetricsLabeler := NewCustomMetricsLabeler()

	if len(keyValues)%2 != 0 {
		// If an odd number of key-value arguments is passed, return the original customMessageLabeler unchanged
		return c
	}

	// Copy existing labels from the current agent
	for k, v := range c.labels {
		newCustomMetricsLabeler.labels[k] = v
	}

	// Add new key-value pairs
	for i := 0; i < len(keyValues); i += 2 {
		key := keyValues[i]
		value := keyValues[i+1]
		newCustomMetricsLabeler.labels[key] = value
	}

	return newCustomMetricsLabeler
}

func (c customMetricsLabeler) incrementRegisterTriggerFailureCounter(ctx context.Context) {
	otelLabels := kvMapToOtelAttributes(c.labels)
	registerTriggerFailureCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c customMetricsLabeler) incrementCapabilityInvocationCounter(ctx context.Context) {
	otelLabels := kvMapToOtelAttributes(c.labels)
	capabilityInvocationCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c customMetricsLabeler) updateWorkflowExecutionLatencyGauge(ctx context.Context, val int64) {
	otelLabels := kvMapToOtelAttributes(c.labels)
	workflowExecutionLatencyGauge.Record(ctx, val, metric.WithAttributes(otelLabels...))
}

func (c customMetricsLabeler) incrementTotalWorkflowStepErrorsCounter(ctx context.Context) {
	otelLabels := kvMapToOtelAttributes(c.labels)
	workflowStepErrorCounter.Add(ctx, 1, metric.WithAttributes(otelLabels...))
}

func (c customMetricsLabeler) updateTotalWorkflowsGauge(ctx context.Context, val int64) {
	otelLabels := kvMapToOtelAttributes(c.labels)
	workflowsRunningGauge.Record(ctx, val, metric.WithAttributes(otelLabels...))
}

// Observability keys
const (
	cIDKey  = "capabilityID"
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
var workflowsRunningGauge metric.Int64Gauge
var capabilityInvocationCounter metric.Int64Counter
var workflowExecutionLatencyGauge metric.Int64Gauge //ms
var workflowStepErrorCounter metric.Int64Counter

func initMonitoringResources() (err error) {
	registerTriggerFailureCounter, err = beholder.GetMeter().Int64Counter("RegisterTriggerFailure")
	if err != nil {
		return fmt.Errorf("failed to register trigger failure: %s", err)
	}

	workflowsRunningGauge, err = beholder.GetMeter().Int64Gauge("WorkflowsRunning")
	if err != nil {
		return fmt.Errorf("failed to register workflows running: %s", err)
	}

	capabilityInvocationCounter, err = beholder.GetMeter().Int64Counter("CapabilityInvocation")
	if err != nil {
		return fmt.Errorf("failed to register capability invocation: %s", err)
	}

	workflowExecutionLatencyGauge, err = beholder.GetMeter().Int64Gauge("WorkflowExecutionLatency")
	if err != nil {
		return fmt.Errorf("failed to register workflow execution latency: %s", err)
	}

	workflowStepErrorCounter, err = beholder.GetMeter().Int64Counter("WorkflowStepError")
	if err != nil {
		return fmt.Errorf("failed to register workflow step error: %s", err)
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
