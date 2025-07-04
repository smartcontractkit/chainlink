package keystone

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	feeds_consumer "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/feeds_consumer_1_0_0"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
)

// Test_runSecureMintWorkflow can be run with:
// `time SECURE_TRANSMITTER_HACK_DISABLED=true CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable go test -timeout 2m -run ^Test_runSecureMintWorkflow$ github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/keystone -v 2>&1 | tee all.log | awk '/DEBUG|INFO|WARN|ERROR/ { print > "node_logs.log"; next }; { print > "other.log" }'; tail all.log`
func Test_runSecureMintWorkflow(t *testing.T) {
	ctx := t.Context()
	lggr := logger.Test(t)
	chainID := chainSelector(1337)
	seqNr := uint64(1)

	// setup the trigger sink that will receive the trigger event in the securemint-specific format
	triggerSink := framework.NewTriggerSink(t, "securemint-trigger", "1.0.0")

	// setup the dons, the size is not important for this test
	workflowDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Workflow", NumNodes: 4, F: 1, AcceptsWorkflows: true})
	require.NoError(t, err)
	triggerDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Trigger", NumNodes: 4, F: 1})
	require.NoError(t, err)
	targetDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Target", NumNodes: 4, F: 1})
	require.NoError(t, err)

	workflowDon, consumer := setupKeystoneDons(ctx, t, lggr, workflowDonConfiguration, triggerDonConfiguration,
		targetDonConfiguration, triggerSink)

	// generate a wf job
	job := createSecureMintWorkflowJob(t, workflowOwnerID, int64(chainID), consumer.Address())
	err = workflowDon.AddJob(ctx, &job)
	require.NoError(t, err)

	// create the test trigger event in the format expected by the secure mint transmitter
	triggerEvent := createSecureMintTriggerEvent(t, chainID, seqNr)

	t.Logf("Sending triggerEvent: %+v", triggerEvent)

	// send the trigger event to the trigger sink and wait for the consumer to receive the feeds
	triggerSink.SendOutput(triggerEvent, "securemint-trigger")
	h := newSecureMintHandler([]secureMintUpdate{}, uint32(time.Now().Unix()))
	waitForConsumerReports(t, consumer, h)
}

type secureMintUpdate struct {
	feedID string
	price  decimal.Decimal
}

// chainSelector is mimicked after the por plugin, which mimics it from the chain-selectors repo
type chainSelector int64

// secureMintReport is mimicked after the report type of the por plugin, see its repo for more details
type secureMintReport struct {
	ConfigDigest ocr2types.ConfigDigest
	SeqNr        uint64
	Block        uint64
	Mintable     *big.Int
}

// createSecureMintTriggerEvent creates a secure mint trigger event in the format sent by the secure mint transmitter
// Excerpt from securemint/transmitter.go:
// ```
//
//	var report ocr3types.ReportWithInfo[por.ChainSelector]
//	outputs, err := values.NewMap(map[string]any{
//		"report":       report,
//		"sigs":         capSigs,
//		"seqNr":        seqNr,
//		"configDigest": cd,
//	})
//
//	event := capabilities.TriggerEvent{
//		TriggerType: t.CapabilityInfo.ID,
//		ID:          "securemint-trigger",
//		Outputs:     outputs,
//	}
//
//	triggerResponse := capabilities.TriggerResponse{
//		Event: event,
//	} // this is sent to trigger subscribers
//
// ```
func createSecureMintTriggerEvent(t *testing.T, chainID chainSelector, seqNr uint64) *values.Map {
	// Create mock signatures (in a real scenario, these would be actual OCR signatures)
	sigs := []commoncap.OCRAttributedOnchainSignature{
		{
			Signer:    0,
			Signature: []byte("mock-signature-1"),
		},
		{
			Signer:    1,
			Signature: []byte("mock-signature-2"),
		},
	}
	configDigest := []byte{0: 1, 31: 2}

	secureMintReport := &secureMintReport{
		ConfigDigest: ocr2types.ConfigDigest(configDigest),
		SeqNr:        seqNr,
		Block:        10,
		Mintable:     big.NewInt(99),
	}

	reportBytes, err := json.Marshal(secureMintReport)
	require.NoError(t, err)

	ocr3Report := &ocr3types.ReportWithInfo[chainSelector]{
		Report: ocr2types.Report(reportBytes),
		Info:   chainID,
	}

	jsonReport, err := json.Marshal(ocr3Report)
	require.NoError(t, err)

	outputs, err := values.NewMap(map[string]any{
		"report":       jsonReport,
		"sigs":         sigs,
		"seqNr":        seqNr,
		"configDigest": configDigest,
	})
	require.NoError(t, err)

	return outputs
}

// secureMintHandler is a handler for the received feeds
// produced by a workflow using the secure mint trigger and aggregator
type secureMintHandler struct {
	expected []secureMintUpdate
	ts       uint32 // unix timestamp in seconds
	found    map[string]struct{}
}

func newSecureMintHandler(expected []secureMintUpdate, ts uint32) *secureMintHandler {
	found := make(map[string]struct{})
	for _, update := range expected {
		found[update.feedID] = struct{}{}
	}
	return &secureMintHandler{
		expected: expected,
		ts:       ts,
		found:    found,
	}
}

// Implement the feedReceivedHandler interface
// to handle the received feeds
func (h *secureMintHandler) handleFeedReceived(t *testing.T, feed *feeds_consumer.KeystoneFeedsConsumerFeedReceived) (done bool) {
	t.Logf("handling event feedID %x", feed.FeedId[:])

	// Convert feed ID to string for comparison
	feedIDStr := fmt.Sprintf("0x%x", feed.FeedId[:])

	// Find the expected update for this feed ID
	var expectedUpdate *secureMintUpdate
	for _, update := range h.expected {
		if update.feedID == feedIDStr {
			expectedUpdate = &update
			break
		}
	}

	// TODO(gg): update the assertions to properly verify the decimal value

	require.NotNil(t, expectedUpdate, "feedID %s not found in expected updates", feedIDStr)

	// Verify the price (assuming 18 decimal places like in the original test)
	assert.Equal(t, expectedUpdate.price.Shift(18).BigInt(), feed.Price)
	assert.Equal(t, h.ts, feed.Timestamp)

	// Mark this feed as found
	delete(h.found, expectedUpdate.feedID)

	// Return true if all expected feeds have been found
	return len(h.found) == 0
}

func (h *secureMintHandler) handleDone(t *testing.T) {
	t.Logf("found %d of %d expected feeds", len(h.expected)-len(h.found), len(h.expected))
	require.Empty(t, h.found, "not all expected feeds were received")
}
