package keystone

import (
	"encoding/hex"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	feeds_consumer "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/feeds_consumer_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/cre"
)

func Test_runLLOWorkflow(t *testing.T) {
	ctx := t.Context()

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.InfoLevel)

	// setup the trigger sink that will receive the trigger event in the llo-specific format, per v2.0.0
	triggerSink := framework.NewTriggerSink(t, "streams-trigger:don_16nodes", "2.0.0") // note the label {"don": "16nodes"} to ensure that we can support labelled capabilities; it must match the llo wf spec in [workflow.go]. the label nor the value are important for this test

	// setup the dons, the size is not important for this test
	workflowDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Workflow", NumNodes: 4, F: 1, AcceptsWorkflows: true})
	require.NoError(t, err)
	triggerDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Trigger", NumNodes: 4, F: 1})
	require.NoError(t, err)
	targetDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Target", NumNodes: 4, F: 1})
	require.NoError(t, err)

	workflowDon, consumer := setupKeystoneDons(ctx, t, lggr, workflowDonConfiguration, triggerDonConfiguration,
		targetDonConfiguration, triggerSink)

	// generate a wf job with 10 feeds
	feedCount := 10
	updates := generateSteamUpdates(t, feedCount)
	job := createLLOStreamWorkflowJob(t, workflowName, workflowOwnerID, streamIDToRemappedID(updates), consumer.Address())
	err = workflowDon.AddJob(ctx, &job)
	require.NoError(t, err)

	// create the test trigger event in the same format as the llo asset don
	ts := time.Now()
	tsUnixNano := uint64(ts.UnixNano()) //nolint: gosec // G115
	e := newLLoTriggerEvent(t, tsUnixNano, updates)
	ocrTrigger, eventID, err := MakeOCRTriggerEvent(lggr, e, triggerDonConfiguration.KeyBundles)
	require.NoError(t, err)
	triggerOutput, err := ocrTrigger.ToMap()
	require.NoError(t, err)

	// send the trigger event to the trigger sink and wait for the consumer to receive the feeds
	triggerSink.SendOutput(triggerOutput, eventID)
	h := newStreamsV2Handler(updates, uint32(ts.Unix())) //nolint: gosec // G115 use the timestamp in seconds for the feed received events
	waitForConsumerReports(t, consumer, h)
}

type streamUpdate struct {
	id         uint32
	remappedID string // hex starting with 0x
	price      decimal.Decimal
}

// MakeOCRTriggerEvent creates an OCR trigger event from the given LLOStreamsTriggerEvent and key bundles
// It can be used to create the underlying data of [capabilities.TriggerEvent.Outputs] in the llo-specific format
func MakeOCRTriggerEvent(lggr logger.Logger, reports *datastreams.LLOStreamsTriggerEvent, keyBundles []ocr2key.KeyBundle) (event *commoncap.OCRTriggerEvent, eventID string, err error) {
	reportCodec := cre.NewReportCodecCapabilityTrigger(lggr, 1 /*donID, unused*/)

	// Convert LLOStreamsTriggerEvent to datastreamsllo.Report
	values := make([]datastreamsllo.StreamValue, len(reports.Payload))
	for i, payload := range reports.Payload {
		// Create decimal stream value
		dec := &datastreamsllo.Decimal{}
		err2 := dec.UnmarshalBinary(payload.Decimal)
		if err2 != nil {
			return nil, "", fmt.Errorf("failed to unmarshal decimal: %w", err2)
		}
		values[i] = dec
	}

	// Create the report
	report := datastreamsllo.Report{
		ObservationTimestampNanoseconds: reports.ObservationTimestampNanoseconds,
		Values:                          values,
	}

	// Create simple channel definition to match streams
	streams := make([]llotypes.Stream, len(reports.Payload))
	for i, payload := range reports.Payload {
		streams[i] = llotypes.Stream{
			StreamID: payload.StreamID,
		}
	}
	channelDef := llotypes.ChannelDefinition{
		Streams: streams,
	}

	// Encode the report to bytes
	reportBytes, err := reportCodec.Encode(report, channelDef)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode report: %w", err)
	}
	eventID = reportCodec.EventID(report)
	// Create report context

	// Create OCR trigger event
	event = &commoncap.OCRTriggerEvent{
		ConfigDigest: []byte{0: 1, 31: 2},
		SeqNr:        0,
		Report:       reportBytes,
		Sigs:         make([]commoncap.OCRAttributedOnchainSignature, 0, len(keyBundles)),
	}

	// Sign the report with each key bundle
	for i, key := range keyBundles {
		sig, err := key.Sign3(ocrTypes.ConfigDigest(event.ConfigDigest), event.SeqNr, reportBytes)
		if err != nil {
			return nil, "", fmt.Errorf("failed to sign report with key %s: %w", key, err)
		}
		event.Sigs = append(event.Sigs, commoncap.OCRAttributedOnchainSignature{
			Signer:    uint32(i), //nolint:gosec // G115
			Signature: sig,
		})
	}

	return event, eventID, nil
}

func generateSteamUpdates(t *testing.T, count int) []streamUpdate {
	var result []streamUpdate
	for i := 1; i <= count; i++ {
		result = append(result, streamUpdate{
			id:         uint32(i), //nolint:gosec // G115
			remappedID: newFeedID(t),
			price:      decimal.NewFromFloat(float64(i)),
		})
	}
	return result
}

func streamIDToRemappedID(updates []streamUpdate) map[uint32]string {
	result := make(map[uint32]string)
	for _, u := range updates {
		result[u.id] = u.remappedID
	}
	return result
}

// streamsV2Handler is a handler for the received feeds
// produced by a workflow using the llo streams trigger and llo aggregator
type streamsV2Handler struct {
	mu       sync.Mutex
	expected []streamUpdate
	ts       uint32 // unix timestamp in seconds

	found map[uint32]struct{}
}

func newStreamsV2Handler(expected []streamUpdate, ts uint32) *streamsV2Handler {
	return &streamsV2Handler{
		expected: expected,
		ts:       ts,
		found:    make(map[uint32]struct{}),
	}
}

// Implement the feedReceivedHandler interface
// to handle the received feeds
func (h *streamsV2Handler) handleFeedReceived(t *testing.T, feed *feeds_consumer.KeystoneFeedsConsumerFeedReceived) (done bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	t.Logf("handling event feedID %x", feed.FeedId[:])
	var updated streamUpdate
	var found bool
	got := "0x" + hex.EncodeToString(feed.FeedId[:])
	for _, s := range h.expected {
		if got == s.remappedID {
			updated = s
			t.Logf("found streamID %d for feedID %s price %s", s.id, got, feed.Price.String())
			found = true
			break
		}
	}
	require.True(t, found, "streamID not found for feedID %s in %v", got, h.expected)

	// TODO cleanup api: we happen to know here that the LLO conversion from decimal to big.Int is 18 decimal places
	assert.Equal(t, updated.price.Shift(18).BigInt(), feed.Price)
	assert.Equal(t, h.ts, feed.Timestamp)
	h.found[updated.id] = struct{}{}
	return len(h.found) == len(h.expected)
}

func (h *streamsV2Handler) handleDone(t *testing.T) {
	t.Logf("found (%v) %d of %d", h.found, len(h.found), len(h.expected))
}

func toPayload(m []streamUpdate) []*datastreams.LLOStreamDecimal {
	result := make([]*datastreams.LLOStreamDecimal, 0, len(m))
	for _, v := range m {
		b, err := v.price.MarshalBinary()
		if err != nil {
			panic(err)
		}
		result = append(result, &datastreams.LLOStreamDecimal{
			StreamID: v.id,
			Decimal:  b,
		})
	}
	return result
}

func newLLoTriggerEvent(t *testing.T, observationTimestamp uint64,
	expected []streamUpdate) *datastreams.LLOStreamsTriggerEvent {
	event := &datastreams.LLOStreamsTriggerEvent{
		ObservationTimestampNanoseconds: observationTimestamp,
		Payload:                         toPayload(expected),
	}
	return event
}
