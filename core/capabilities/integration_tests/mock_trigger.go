package integration_tests

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const triggerID = "streams-trigger@1.0.0"

func mockMercuryTrigger(t *testing.T, reports []datastreams.FeedReport) capabilities.TriggerCapability {
	mt := &mockTriggerCapability{
		CapabilityInfo: capabilities.MustNewCapabilityInfo(
			triggerID,
			capabilities.CapabilityTypeTrigger,
			"issues a trigger when a mercury report is received.",
		),
		ch: make(chan capabilities.CapabilityResponse, 10),
	}

	resp, err := wrapReports(reports, "1", 12, datastreams.SignersMetadata{})
	require.NoError(t, err)
	mt.triggerEvent = &resp
	return mt
}

type mockTriggerCapability struct {
	capabilities.CapabilityInfo
	triggerEvent *capabilities.CapabilityResponse
	ch           chan capabilities.CapabilityResponse
}

var _ capabilities.TriggerCapability = (*mockTriggerCapability)(nil)

func (m *mockTriggerCapability) RegisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) (<-chan capabilities.CapabilityResponse, error) {
	if m.triggerEvent != nil {
		m.ch <- *m.triggerEvent
	}
	return m.ch, nil
}

func (m *mockTriggerCapability) UnregisterTrigger(ctx context.Context, req capabilities.CapabilityRequest) error {
	return nil
}

func wrapReports(reportList []datastreams.FeedReport, eventID string, timestamp int64, meta datastreams.SignersMetadata) (capabilities.CapabilityResponse, error) {
	val, err := values.Wrap(reportList)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	metaVal, err := values.Wrap(meta)
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}

	triggerEvent := capabilities.TriggerEvent{
		TriggerType: triggerID,
		ID:          eventID,
		Timestamp:   strconv.FormatInt(timestamp, 10),
		Metadata:    metaVal,
		Payload:     val,
	}

	//triggerEventMap := map[string]interface{}{"event": triggerEvent}

	triggerEventMapValue, err := values.WrapMap(triggerEvent)

	// Create a new CapabilityResponse with the MercuryTriggerEvent
	return capabilities.CapabilityResponse{
		Value: triggerEventMapValue,
	}, nil
}
