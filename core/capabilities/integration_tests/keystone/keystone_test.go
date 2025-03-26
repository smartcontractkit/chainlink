package keystone

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
	feeds_consumer "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/cre"
	reporttypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/types"
)

func Test_AllAtOnceTransmissionSchedule(t *testing.T) {
	testTransmissionSchedule(t, "2s", "allAtOnce")
}

func Test_OneAtATimeTransmissionSchedule(t *testing.T) {
	testTransmissionSchedule(t, "2s", "oneAtATime")
}

func testTransmissionSchedule(t *testing.T, deltaStage string, schedule string) {
	ctx, cancel := framework.Context(t)
	defer cancel()

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.InfoLevel)
	logsigners := func(donName string, c framework.DonConfiguration) {
		signers := make([]string, len(c.KeyBundles))
		for i, keyBundle := range c.KeyBundles {
			signers[i] = keyBundle.OnChainPublicKey()
		}
		lggr.Infow("created", "donName", donName, "signers", signers)
	}

	workflowDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Workflow", NumNodes: 4, F: 1, AcceptsWorkflows: true})
	require.NoError(t, err)
	logsigners("Workflow", workflowDonConfiguration)
	triggerDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Trigger", NumNodes: 4, F: 1})
	require.NoError(t, err)
	logsigners("Trigger", triggerDonConfiguration)
	targetDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Target", NumNodes: 4, F: 1})
	require.NoError(t, err)
	logsigners("Target", targetDonConfiguration)

	/*
		// mercury-style reports
		t.Run("Feeds v1.0.0", func(t *testing.T) {
			triggerSink := framework.NewTriggerSink(t, "streams-trigger", "1.0.0")
			workflowDon, consumer := setupKeystoneDons(ctx, t, lggr, workflowDonConfiguration, triggerDonConfiguration,
				targetDonConfiguration, triggerSink)

			feedCount := 3
			var feedIDs []string
			for i := 0; i < feedCount; i++ {
				feedIDs = append(feedIDs, newFeedID(t))
			}

			job := createKeystoneWorkflowJob(t, workflowName, workflowOwnerID, feedIDs, consumer.Address(), deltaStage, schedule)
			err = workflowDon.AddJob(ctx, &job)
			require.NoError(t, err)

			reports := []*datastreams.FeedReport{p
				createFeedReport(t, big.NewInt(1), 5, feedIDs[0], triggerDonConfiguration.KeyBundles),
				createFeedReport(t, big.NewInt(3), 7, feedIDs[1], triggerDonConfiguration.KeyBundles),
				createFeedReport(t, big.NewInt(2), 6, feedIDs[2], triggerDonConfiguration.KeyBundles),
			}

			wrappedReports, err := wrapReports(reports, 12, datastreams.Metadata{})
			require.NoError(t, err)

			triggerSink.SendOutput(wrappedReports)

			waitForConsumerReports(ctx, t, consumer, reports)
		})
	*/
	// llo-style reports
	t.Run("Streams v2.0.0", func(t *testing.T) {
		triggerSink := framework.NewTriggerSink(t, "streams-trigger", "2.0.0")
		workflowDon, consumer := setupKeystoneDons(ctx, t, lggr, workflowDonConfiguration, triggerDonConfiguration,
			targetDonConfiguration, triggerSink)
		feedCount := 3
		streamIDtoFeedID := make(map[uint32][]byte)
		for i := 1; i < feedCount; i++ {
			streamIDtoFeedID[uint32(i)] = newFeedIDBytes(t)
		}

		job := createLLOStreamWorkflowJob(t, workflowName, workflowOwnerID, streamIDtoFeedID, consumer.Address(), deltaStage, schedule)
		err = workflowDon.AddJob(ctx, &job)
		require.NoError(t, err)

		expected := []streamUpdate{
			{
				id:         1,
				remappedID: "0x" + hex.EncodeToString(streamIDtoFeedID[1]),
				price:      decimal.NewFromFloat(1250.427975),
			},
			{
				id:         2,
				remappedID: "0x" + hex.EncodeToString(streamIDtoFeedID[2]),
				price:      decimal.NewFromFloat(42.37),
			},
			{
				id:         3,
				remappedID: "0x" + hex.EncodeToString(streamIDtoFeedID[3]),
				price:      decimal.NewFromFloat(0.003),
			},
		}

		ts := uint64(time.Now().UnixNano())
		e := newLLoTriggerEvent(t, ts, expected)
		ocrTrigger, eventID := makeOCRTriggerEvent2(t, e, triggerDonConfiguration.KeyBundles)

		triggerSink.SendOCREvent(ocrTrigger, eventID)

		h := &handler{
			expected: expected,
			ts:       ts,
			found:    make(map[uint32]struct{}),
		}

		waitForConsumerReports2(ctx, t, consumer, h)
	})
}

type streamUpdate struct {
	id         uint32
	remappedID string //hex 0x
	price      decimal.Decimal
}

type handler struct {
	expected []streamUpdate
	ts       uint64

	found map[uint32]struct{}
}

// Implement the feedRecievedHandler interface
// to handle the received feeds
func (h *handler) handleFeedReceived(t *testing.T, feed *feeds_consumer.KeystoneFeedsConsumerFeedReceived) (done bool) {
	t.Logf("handling event feedID %x", feed.FeedId[:])
	var updated streamUpdate
	var found bool
	got := "0x" + hex.EncodeToString(feed.FeedId[:])
	for _, s := range h.expected {
		if got == s.remappedID {
			updated = s
			found = true
			break
		}
	}
	require.True(t, found, "streamID not found for feedID %s in %v", got, h.expected)

	// TODO cleanup api
	assert.Equal(t, updated.price.Shift(18).BigInt(), feed.Price)
	assert.Equal(t, uint32(h.ts), feed.Timestamp)
	h.found[updated.id] = struct{}{}
	return len(h.found) == len(h.expected)
}

func (h *handler) handleDone(t *testing.T) {
	t.Logf("found (%v) %d of %d", h.found, len(h.found), len(h.expected))
}
func toPayload(m []streamUpdate) []*datastreams.LLOStreamDecimal {
	var result []*datastreams.LLOStreamDecimal
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

func wrapReports(reportList []*datastreams.FeedReport,
	timestamp int64, meta datastreams.Metadata) (*values.Map, error) {
	var rl []datastreams.FeedReport
	for _, r := range reportList {
		rl = append(rl, *r)
	}

	return values.WrapMap(datastreams.StreamsTriggerEvent{
		Payload:   rl,
		Metadata:  meta,
		Timestamp: timestamp,
	})
}

func makeOCRTriggerEvent2(t *testing.T, reports *datastreams.LLOStreamsTriggerEvent, keyBundles []ocr2key.KeyBundle) (event *commoncap.OCRTriggerEvent, eventID string) {
	// Create the report codec with a don ID (using 1 for testing)
	lggr := logger.TestLogger(t)
	reportCodec := cre.NewReportCodecCapabilityTrigger(lggr, 1)

	// Convert LLOStreamsTriggerEvent to datastreamsllo.Report
	values := make([]datastreamsllo.StreamValue, len(reports.Payload))
	for i, payload := range reports.Payload {
		// Create decimal stream value
		dec := &datastreamsllo.Decimal{}
		err := dec.UnmarshalBinary(payload.Decimal)
		require.NoError(t, err)
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
	reportBytes, err := reportCodec.Encode(context.Background(), report, channelDef)
	require.NoError(t, err)
	eventID = reportCodec.EventID(report)
	// Create report context
	//reportCtx := ocrTypes.ReportContext{}

	// Create OCR trigger event
	event = &commoncap.OCRTriggerEvent{
		ConfigDigest: []byte{0: 1, 31: 2},
		SeqNr:        0,
		Report:       reportBytes,
		Sigs:         make([]commoncap.OCRAttributedOnchainSignature, 0, len(keyBundles)),
	}

	// Create signature for the report
	//reportHash := ocr2key.ReportToSigData3(ocrTypes.ConfigDigest(event.ConfigDigest), event.SeqNr, reportBytes)

	// Sign the report with each key bundle
	for i, key := range keyBundles {
		sig, err := key.Sign3(ocrTypes.ConfigDigest(event.ConfigDigest), event.SeqNr, reportBytes)
		require.NoError(t, err)
		event.Sigs = append(event.Sigs, commoncap.OCRAttributedOnchainSignature{
			Signer:    uint32(i),
			Signature: sig,
		})
	}

	return event, eventID
}

func newFeedID(t *testing.T) string {
	b := newFeedIDBytes(t)
	return "0x" + hex.EncodeToString(b[:])
}

func newFeedIDBytes(t *testing.T) []byte {
	buf := [32]byte{}
	_, err := rand.Read(buf[:])
	require.NoError(t, err)
	return buf[:]
}

func feedIDToBytes(t *testing.T, feedID string) [32]byte {
	bytes, err := hex.DecodeString(feedID[2:])
	require.NoError(t, err)
	var feedIDBytes [32]byte
	copy(feedIDBytes[:], bytes)
	return feedIDBytes
}

func waitForConsumerReports(ctx context.Context, t *testing.T, consumer *feeds_consumer.KeystoneFeedsConsumer, triggerFeedReports []*datastreams.FeedReport) {
	feedsReceived := make(chan *feeds_consumer.KeystoneFeedsConsumerFeedReceived, 1000)
	feedsSub, err := consumer.WatchFeedReceived(&bind.WatchOpts{}, feedsReceived, nil)
	require.NoError(t, err)

	feedToReport := map[string]*datastreams.FeedReport{}
	for _, report := range triggerFeedReports {
		feedToReport[report.FeedID] = report
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	feedCount := 0
	for {
		select {
		case <-ctxWithTimeout.Done():
			t.Fatalf("timed out waiting for feeds reports, expected %d, received %d", len(triggerFeedReports), feedCount)
		case err := <-feedsSub.Err():
			require.NoError(t, err)
		case feed := <-feedsReceived:
			feedID := "0x" + hex.EncodeToString(feed.FeedId[:])
			report := feedToReport[feedID]
			decodedReport, err := reporttypes.Decode(report.FullReport)
			require.NoError(t, err)
			assert.Equal(t, decodedReport.BenchmarkPrice, feed.Price)
			assert.Equal(t, decodedReport.ObservationsTimestamp, feed.Timestamp)

			feedCount++
			if feedCount == len(triggerFeedReports) {
				return
			}
		}
	}
}

type feedRecievedHandler interface {
	handleFeedReceived(t *testing.T, feed *feeds_consumer.KeystoneFeedsConsumerFeedReceived) (done bool)
	handleDone(t *testing.T)
}

func waitForConsumerReports2(ctx context.Context, t *testing.T, consumer *feeds_consumer.KeystoneFeedsConsumer, h feedRecievedHandler) {
	feedsReceived := make(chan *feeds_consumer.KeystoneFeedsConsumerFeedReceived, 1000)
	feedsSub, err := consumer.WatchFeedReceived(&bind.WatchOpts{}, feedsReceived, nil)
	require.NoError(t, err)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctxWithTimeout.Done():
			h.handleDone(t)
			t.Fatalf("timed out waiting for feeds reports")
		case err := <-feedsSub.Err():
			require.NoError(t, err)
		case feed := <-feedsReceived:
			done := h.handleFeedReceived(t, feed)
			if done {
				return
			}
		}
	}
}
