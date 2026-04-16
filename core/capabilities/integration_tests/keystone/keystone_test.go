package keystone

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	feeds_consumer "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/feeds_consumer_1_0_0"
	reporttypes "github.com/smartcontractkit/chainlink-evm/pkg/mercury/v3/types"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
)

func Test_AllAtOnceTransmissionSchedule(t *testing.T) {
	testTransmissionSchedule(t, "2s", "allAtOnce")
}

func Test_OneAtATimeTransmissionSchedule(t *testing.T) {
	testTransmissionSchedule(t, "2s", "oneAtATime")
}

func Test_AllAtOnceZeroDelta_SubmitsDuplicateForwarderTxs(t *testing.T) {
	rawDBURL, ok := os.LookupEnv("CL_DATABASE_URL")
	if !ok {
		t.Skip("CL_DATABASE_URL is required for this integration test")
	}
	parsedDBURL, err := url.Parse(rawDBURL)
	if err != nil || parsedDBURL.Path == "" {
		t.Skip("CL_DATABASE_URL must include a database path for this integration test")
	}

	ctx := t.Context()
	lggr := logger.Test(t)

	workflowDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{
		Name:             "Workflow",
		NumNodes:         7,
		F:                2,
		AcceptsWorkflows: true,
	})
	require.NoError(t, err)
	triggerDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{
		Name:     "Trigger",
		NumNodes: 7,
		F:        2,
	})
	require.NoError(t, err)
	targetDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{
		Name:     "Target",
		NumNodes: 7,
		F:        2,
	})
	require.NoError(t, err)

	triggerSink := framework.NewTriggerSink(t, "streams-trigger", "1.0.0")
	donContext := framework.CreateDonContext(ctx, t)

	workflowDon := createKeystoneWorkflowDon(ctx, t, lggr, workflowDonConfiguration, triggerDonConfiguration, targetDonConfiguration, donContext)
	forwarderAddr, _ := SetupForwarderContract(t, workflowDon, donContext.EthBlockchain)
	_, consumer := SetupConsumerContract(t, donContext.EthBlockchain, forwarderAddr, workflowOwnerID, workflowName)
	writeTargetDon := createKeystoneWriteTargetDon(ctx, t, lggr, targetDonConfiguration, donContext, forwarderAddr)
	triggerDon := createKeystoneTriggerDon(ctx, t, lggr, triggerDonConfiguration, donContext, triggerSink)

	servicetest.Run(t, workflowDon)
	servicetest.Run(t, triggerDon)
	servicetest.Run(t, writeTargetDon)
	donContext.WaitForCapabilitiesToBeExposed(t, workflowDon, triggerDon, writeTargetDon)

	feedID := newFeedID(t)
	job := createKeystoneWorkflowJob(
		t,
		workflowName,
		workflowOwnerID,
		[]string{feedID},
		consumer.Address(),
		"0s",
		"allAtOnce",
	)
	err = workflowDon.AddJob(ctx, &job)
	require.NoError(t, err)

	report := createFeedReport(t, big.NewInt(1), 7, feedID, triggerDonConfiguration.KeyBundles)
	wrappedReports, err := wrapReports([]*datastreams.FeedReport{report}, 12, datastreams.Metadata{})
	require.NoError(t, err)

	startBlock, err := donContext.EthBlockchain.Client().BlockNumber(ctx)
	require.NoError(t, err)

	triggerSink.SendOutput(wrappedReports, uuid.New().String())
	waitForConsumerReports(t, consumer, newStreamsV1Handler([]*datastreams.FeedReport{report}))

	var observedTxCount atomic.Uint64
	require.Eventually(t, func() bool {
		latestBlock, blockErr := donContext.EthBlockchain.Client().BlockNumber(ctx)
		if blockErr != nil || latestBlock <= startBlock {
			return false
		}
		txCount, countErr := countTransactionsToAddress(ctx, donContext.EthBlockchain, forwarderAddr, startBlock+1, latestBlock)
		if countErr != nil {
			return false
		}
		observedTxCount.Store(txCount)
		return txCount > 1
	}, 20*time.Second, 500*time.Millisecond, "expected duplicate forwarder tx submissions for one execution")

	// TODO: @ilija42 Expected behavior after fix: exactly one forwarder transaction should be submitted for a single execution.
	// require.EqualValues(
	// 	t,
	// 	1,
	// 	observedTxCount.Load(),
	// 	"TODO: enforce single forwarder submission per execution under allAtOnce + low deltaStage",
	// )
}

func testTransmissionSchedule(t *testing.T, deltaStage string, schedule string) {
	ctx := t.Context()

	lggr := logger.Test(t)

	workflowDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Workflow", NumNodes: 7, F: 2, AcceptsWorkflows: true})
	require.NoError(t, err)
	triggerDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Trigger", NumNodes: 7, F: 2})
	require.NoError(t, err)
	targetDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Target", NumNodes: 4, F: 1})
	require.NoError(t, err)

	// mercury-style reports
	triggerSink := framework.NewTriggerSink(t, "streams-trigger", "1.0.0")
	workflowDon, consumer := setupKeystoneDons(ctx, t, lggr, workflowDonConfiguration, triggerDonConfiguration,
		targetDonConfiguration, triggerSink)

	feedCount := 3
	feedIDs := make([]string, 0, feedCount)
	for range feedCount {
		feedIDs = append(feedIDs, newFeedID(t))
	}

	job := createKeystoneWorkflowJob(t, workflowName, workflowOwnerID, feedIDs, consumer.Address(), deltaStage, schedule)
	err = workflowDon.AddJob(ctx, &job)
	require.NoError(t, err)

	reports := []*datastreams.FeedReport{
		createFeedReport(t, big.NewInt(1), 5, feedIDs[0], triggerDonConfiguration.KeyBundles),
		createFeedReport(t, big.NewInt(3), 7, feedIDs[1], triggerDonConfiguration.KeyBundles),
		createFeedReport(t, big.NewInt(2), 6, feedIDs[2], triggerDonConfiguration.KeyBundles),
	}

	wrappedReports, err := wrapReports(reports, 12, datastreams.Metadata{})
	require.NoError(t, err)

	triggerSink.SendOutput(wrappedReports, uuid.New().String())
	h := newStreamsV1Handler(reports)

	waitForConsumerReports(t, consumer, h)
}

func wrapReports(reportList []*datastreams.FeedReport,
	timestamp int64, meta datastreams.Metadata) (*values.Map, error) {
	rl := make([]datastreams.FeedReport, 0, len(reportList))
	for _, r := range reportList {
		rl = append(rl, *r)
	}

	return values.WrapMap(datastreams.StreamsTriggerEvent{
		Payload:   rl,
		Metadata:  meta,
		Timestamp: timestamp,
	})
}

func newFeedID(t *testing.T) string {
	buf := [32]byte{}
	_, err := rand.Read(buf[:])
	require.NoError(t, err)
	return "0x" + hex.EncodeToString(buf[:])
}

type feedReceivedHandler interface {
	handleFeedReceived(t *testing.T, feed *feeds_consumer.KeystoneFeedsConsumerFeedReceived) (done bool)
	handleDone(t *testing.T)
}

func waitForConsumerReports(t *testing.T, consumer *feeds_consumer.KeystoneFeedsConsumer, h feedReceivedHandler) {
	feedsReceived := make(chan *feeds_consumer.KeystoneFeedsConsumerFeedReceived, 1000)
	feedsSub, err := consumer.WatchFeedReceived(&bind.WatchOpts{}, feedsReceived, nil)
	require.NoError(t, err)
	ctx := t.Context()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
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

type streamsV1Handler struct {
	mu       sync.Mutex
	expected map[string]*datastreams.FeedReport
	found    map[string]struct{}
}

func newStreamsV1Handler(expected []*datastreams.FeedReport) *streamsV1Handler {
	h := &streamsV1Handler{
		expected: make(map[string]*datastreams.FeedReport),
		found:    make(map[string]struct{}),
	}
	for _, report := range expected {
		h.expected[report.FeedID] = report
	}
	return h
}

func (h *streamsV1Handler) handleFeedReceived(t *testing.T, feed *feeds_consumer.KeystoneFeedsConsumerFeedReceived) (done bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	feedID := "0x" + hex.EncodeToString(feed.FeedId[:])
	report := h.expected[feedID]
	require.NotNil(t, report, "unexpected feedID %s", feedID)

	decodedReport, err := reporttypes.Decode(report.FullReport)
	require.NoError(t, err)

	assert.Equal(t, decodedReport.BenchmarkPrice, feed.Price)
	assert.Equal(t, decodedReport.ObservationsTimestamp, feed.Timestamp)

	h.found[feedID] = struct{}{}
	return len(h.found) == len(h.expected)
}

func (h *streamsV1Handler) handleDone(t *testing.T) {
	h.mu.Lock()
	defer h.mu.Unlock()
	t.Logf("found (%v) %d of %d", h.found, len(h.found), len(h.expected))
}

func countTransactionsToAddress(
	ctx context.Context,
	chain *framework.EthBlockchain,
	target common.Address,
	fromBlock uint64,
	toBlock uint64,
) (uint64, error) {
	var count uint64
	for blockNum := fromBlock; blockNum <= toBlock; blockNum++ {
		block, err := chain.Client().BlockByNumber(ctx, new(big.Int).SetUint64(blockNum))
		if err != nil {
			return 0, err
		}

		for _, tx := range block.Transactions() {
			if tx.To() != nil && *tx.To() == target {
				count++
			}
		}
	}

	return count, nil
}
