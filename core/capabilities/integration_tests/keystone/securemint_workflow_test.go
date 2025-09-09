package keystone

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	data_feeds_cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
)

// Test_runSecureMintWorkflow can be run with:
// `CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable go test -timeout 2m -run ^Test_runSecureMintWorkflow$ github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/keystone -v 2>&1 | tee all.log | awk '/DEBUG|INFO|WARN|ERROR/ { print > "node_logs.log"; next }; { print > "other.log" }'; tail all.log`
func Test_runSecureMintWorkflow(t *testing.T) {
	ctx := t.Context()
	lggr := logger.Test(t)
	chainID := chainSelector(16015286601757825753)
	seqNr := uint64(123456789)

	// setup the trigger sink that will receive the trigger event in the securemint-specific format
	triggerSink := framework.NewTriggerSink(t, "securemint-trigger", "1.0.0")

	// setup the dons, the size is not important for this test
	workflowDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Workflow", NumNodes: 4, F: 1, AcceptsWorkflows: true})
	require.NoError(t, err)
	triggerDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Trigger", NumNodes: 4, F: 1})
	require.NoError(t, err)
	targetDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Target", NumNodes: 4, F: 1})
	require.NoError(t, err)

	workflowDon, _, dataFeedsCache, forwarder := setupKeystoneDons(ctx, t, lggr, workflowDonConfiguration, triggerDonConfiguration,
		targetDonConfiguration, triggerSink)
	t.Logf("Forwarder contract address: %s", forwarder.Address().String())
	t.Logf("DataFeedsCache contract address: %s", dataFeedsCache.Address().String())

	// make sure we know about forwarder errors in case they happen
	trackErrorsOnForwarder(t, forwarder, dataFeedsCache.Address())

	// track invalid update permission events on the DF Cache
	trackInvalidPermissionEventsOnDFCache(t, dataFeedsCache)

	// generate a wf job
	job := createSecureMintWorkflowJob(t, workflowName, workflowOwnerID, uint64(chainID), secureMintFeedDataID, dataFeedsCache.Address())
	err = workflowDon.AddJob(ctx, &job)
	require.NoError(t, err)

	// create the test trigger event in the format expected by the secure mint transmitter
	mintableAmount := big.NewInt(99)
	blockNumber := big.NewInt(10)
	triggerEvent := createSecureMintTriggerEvent(t, chainID, seqNr, mintableAmount, blockNumber)

	t.Logf("Sending triggerEvent: %+v", triggerEvent)

	// send the trigger event to the trigger sink and wait for the consumer to receive the feeds
	triggerSink.SendOutput(triggerEvent, "securemint-trigger")

	// The workflow is configured to use feed ID "1020001001" and should generate a feed
	// with a price derived from the Mintable value (99) in the trigger event
	// The data ID is generated from the chain selector (1337) as bytes: 0x04 + chain selector as bytes + right padded with 0s
	// The price is packed from Mintable (99) and block number (10)
	expectedUpdates := []secureMintUpdate{
		{
			dataID:         secureMintFeedDataID,
			mintableAmount: mintableAmount,
			blockNumber:    blockNumber,
		},
	}
	h := newSecureMintHandler(expectedUpdates, big.NewInt(int64(seqNr))) // currently the secure mint aggregator uses the reports' sequence number as timestamp
	waitForDataFeedsCacheReports(t, dataFeedsCache, h)
}

type secureMintUpdate struct {
	dataID         string
	mintableAmount *big.Int
	blockNumber    *big.Int
}

// chainSelector is mimicked after the por plugin, which mimics it from the chain-selectors repo
type chainSelector uint64

// secureMintReport is mimicked after the report type of the por plugin, see its repo for more details
type secureMintReport struct {
	ConfigDigest ocr2types.ConfigDigest
	SeqNr        uint64
	Block        uint64
	Mintable     *big.Int
}

// createSecureMintTriggerEvent creates a secure mint trigger event in the format sent by the secure mint transmitter
// Excerpt from securemint/transmitter.go:
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
//	}
func createSecureMintTriggerEvent(t *testing.T, chainID chainSelector, seqNr uint64, mintable *big.Int, blockNumber *big.Int) *values.Map {
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
		Block:        blockNumber.Uint64(),
		Mintable:     mintable,
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
	ts       *big.Int
	found    map[string]struct{}
}

func newSecureMintHandler(expected []secureMintUpdate, ts *big.Int) *secureMintHandler {
	found := make(map[string]struct{})
	for _, update := range expected {
		found[update.dataID] = struct{}{}
	}
	return &secureMintHandler{
		expected: expected,
		ts:       ts,
		found:    found,
	}
}

// Implement the dataFeedsCacheHandler interface
// to handle the received feeds from DataFeedsCache
func (h *secureMintHandler) handleDecimalReportUpdated(t *testing.T, event *data_feeds_cache.DataFeedsCacheDecimalReportUpdated) (done bool) {
	t.Logf("handling event for dataID %x: %+v", event.DataId[:], event)

	// Convert data ID to string for comparison (DataFeedsCache uses bytes16 dataId instead of bytes32 feedId)
	dataIDStr := fmt.Sprintf("0x%x", event.DataId[:])

	// Find the expected update for this data ID
	var expectedUpdate *secureMintUpdate
	for _, update := range h.expected {
		if update.dataID == dataIDStr {
			expectedUpdate = &update
			break
		}
	}

	require.NotNil(t, expectedUpdate, "dataID %s not found in expected updates", dataIDStr)

	mintableMask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))
	extractedMintable := new(big.Int).And(event.Answer, mintableMask)
	t.Logf("extractedMintable: %d", extractedMintable)
	assert.Equalf(t, expectedUpdate.mintableAmount, extractedMintable, "mintable amount mismatch: expected %d, got %d", expectedUpdate.mintableAmount, extractedMintable)

	// Extract block number from bits 128-191
	blockNumberMask := new(big.Int).Lsh(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 64), big.NewInt(1)), 128)
	extractedBlockNumber := new(big.Int).And(event.Answer, blockNumberMask)
	extractedBlockNumber = new(big.Int).Rsh(extractedBlockNumber, 128)
	t.Logf("extractedBlockNumber: %d", extractedBlockNumber)
	assert.Equalf(t, expectedUpdate.blockNumber, extractedBlockNumber, "block number mismatch: expected %d, got %d", expectedUpdate.blockNumber, extractedBlockNumber.Uint64())

	assert.Equalf(t, h.ts, event.Timestamp, "timestamp mismatch: expected %d, got %d", h.ts, event.Timestamp)

	// Mark this feed as found
	delete(h.found, expectedUpdate.dataID)

	// Return true if all expected feeds have been found
	t.Logf("found %d of %d expected feeds, left: %+v, expected: %+v", len(h.expected)-len(h.found), len(h.expected), h.found, h.expected)
	return len(h.found) == 0
}

func (h *secureMintHandler) handleDone(t *testing.T) {
	t.Logf("found %d of %d expected feeds", len(h.expected)-len(h.found), len(h.expected))
	require.Empty(t, h.found, "not all expected feeds were received")
}

// Interface for DataFeedsCache event handling
type dataFeedsCacheHandler interface {
	handleDecimalReportUpdated(t *testing.T, event *data_feeds_cache.DataFeedsCacheDecimalReportUpdated) (done bool)
	handleDone(t *testing.T)
}

// waitForDataFeedsCacheReports waits for DecimalReportUpdated events from DataFeedsCache contract
func waitForDataFeedsCacheReports(t *testing.T, dataFeedsCache *data_feeds_cache.DataFeedsCache, h dataFeedsCacheHandler) {
	reportsReceived := make(chan *data_feeds_cache.DataFeedsCacheDecimalReportUpdated, 1000)
	reportsSub, err := dataFeedsCache.WatchDecimalReportUpdated(&bind.WatchOpts{}, reportsReceived, nil, nil, nil)
	require.NoError(t, err)
	ctx := t.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			h.handleDone(t)
			t.Fatalf("timed out waiting for data feeds cache reports")
		case err := <-reportsSub.Err():
			require.NoError(t, err)
		case report := <-reportsReceived:
			done := h.handleDecimalReportUpdated(t, report)
			if done {
				return
			}
		}
	}
}
