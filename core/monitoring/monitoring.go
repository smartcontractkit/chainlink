package monitoring

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	beholderpb "github.com/smartcontractkit/chainlink-common/pkg/beholder/pb"
	valuespb "github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	"google.golang.org/protobuf/proto"
)

type CustomMessageLabeler struct {
	labels map[string]string
}

func NewCustomMessageLabeler() CustomMessageLabeler {
	return CustomMessageLabeler{labels: make(map[string]string)}
}

// With adds multiple key-value pairs to the CustomMessageLabeler for transmission With SendLogAsCustomMessage
func (c CustomMessageLabeler) With(keyValues ...string) CustomMessageLabeler {
	newCustomMessageLabeler := NewCustomMessageLabeler()

	if len(keyValues)%2 != 0 {
		// If an odd number of key-value arguments is passed, return the original CustomMessageLabeler unchanged
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

// SendLogAsCustomMessage emits a BaseMessage With msg and labels as data.
// any key in labels that is not part of orderedLabelKeys will not be transmitted
func (c CustomMessageLabeler) SendLogAsCustomMessage(msg string) error {
	return sendLogAsCustomMessageW(msg, c.labels)
}

type MetricsLabeler struct {
	Labels map[string]string
}

func NewMetricsLabeler() MetricsLabeler {
	return MetricsLabeler{Labels: make(map[string]string)}
}

// With adds multiple key-value pairs to the CustomMessageLabeler for transmission With SendLogAsCustomMessage
func (c MetricsLabeler) With(keyValues ...string) MetricsLabeler {
	newCustomMetricsLabeler := NewMetricsLabeler()

	if len(keyValues)%2 != 0 {
		// If an odd number of key-value arguments is passed, return the original CustomMessageLabeler unchanged
		return c
	}

	// Copy existing labels from the current agent
	for k, v := range c.Labels {
		newCustomMetricsLabeler.Labels[k] = v
	}

	// Add new key-value pairs
	for i := 0; i < len(keyValues); i += 2 {
		key := keyValues[i]
		value := keyValues[i+1]
		newCustomMetricsLabeler.Labels[key] = value
	}

	return newCustomMetricsLabeler
}

// sendLogAsCustomMessageF formats into a msg to be consumed by sendLogAsCustomMessageW
func sendLogAsCustomMessageF(labels map[string]string, format string, values ...any) error {
	return sendLogAsCustomMessageW(fmt.Sprintf(format, values...), labels)
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
	protoLabels := make(map[string]*valuespb.Value)
	for _, l := range labels {
		protoLabels[l] = &valuespb.Value{Value: &valuespb.Value_StringValue{StringValue: labels[l]}}
	}
	// Define a custom protobuf payload to emit
	payload := &beholderpb.BaseMessage{
		Msg:    msg,
		Labels: protoLabels,
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return fmt.Errorf("sending custom message failed to marshal protobuf: %w", err)
	}

	err = beholder.GetEmitter().Emit(context.Background(), payloadBytes,
		"beholder_data_schema", "/beholder-base-message/versions/1", // required
		"beholder_data_type", "custom_message",
	)
	if err != nil {
		return fmt.Errorf("sending custom message failed on emit: %w", err)
	}

	return nil
}
