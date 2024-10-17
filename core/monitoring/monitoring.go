package monitoring

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	beholderpb "github.com/smartcontractkit/chainlink-common/pkg/beholder/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
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

func sendLogAsCustomMessageW(msg string, labels map[string]string) error {
	// cast to map[string]any
	newLabels := map[string]any{}
	for k, v := range labels {
		newLabels[k] = v
	}

	m, err := values.NewMap(newLabels)
	if err != nil {
		return fmt.Errorf("could not wrap labels to map: %w", err)
	}

	// Define a custom protobuf payload to emit
	payload := &beholderpb.BaseMessage{
		Msg:    msg,
		Labels: values.ProtoMap(m),
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
