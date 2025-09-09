package keystone

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	data_feeds_cache "github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	feeds_consumer "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/feeds_consumer_1_0_0"
	fwd "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
	reporttypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/types"
)

func Test_AllAtOnceTransmissionSchedule(t *testing.T) {
	testTransmissionSchedule(t, "2s", "allAtOnce")
}

func Test_OneAtATimeTransmissionSchedule(t *testing.T) {
	testTransmissionSchedule(t, "2s", "oneAtATime")
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
	workflowDon, consumer, _, _ := setupKeystoneDons(ctx, t, lggr, workflowDonConfiguration, triggerDonConfiguration,
		targetDonConfiguration, triggerSink)

	feedCount := 3
	var feedIDs []string
	for i := 0; i < feedCount; i++ {
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

// trackErrorsOnForwarder watches the forwarder contract for report processed events and fails the test if the report is not forwarded to the consumer
func trackErrorsOnForwarder(t *testing.T, forwarder *fwd.KeystoneForwarder, dfCacheAddress common.Address) {
	t.Helper()

	reportsProcessed := make(chan *fwd.KeystoneForwarderReportProcessed, 1000)
	reportsSub, err := forwarder.WatchReportProcessed(nil, reportsProcessed, nil, nil, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	closeFunc := func() {
		cancel()
		<-done
	}
	t.Cleanup(closeFunc)

	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-reportsSub.Err():
				assert.NoError(t, err)
				return
			case report := <-reportsProcessed:
				t.Logf("Forwarder received report: %+v", report)

				transmissionInfo, err := forwarder.GetTransmissionInfo(nil, dfCacheAddress, report.WorkflowExecutionId, report.ReportId)
				assert.NoError(t, err)
				if !report.Result { // if the report is not forwarded to the consumer, we need to get the transmission info to see why
					t.Errorf("Report not forwarded to DataFeeds Cache: %+v", transmissionInfo)
				} else {
					t.Logf("Report successfully forwarded to DataFeeds Cache: %+v", transmissionInfo)
				}
			}
		}
	}()
}

// trackInvalidPermissionEventsOnDFCache watches the DF Cache contract for invalid permission events
func trackInvalidPermissionEventsOnDFCache(t *testing.T, dataFeedsCache *data_feeds_cache.DataFeedsCache) {
	t.Helper()

	invalidUpdatePermissionEvents := make(chan *data_feeds_cache.DataFeedsCacheInvalidUpdatePermission, 1000)
	invalidUpdatePermissionSub, err := dataFeedsCache.WatchInvalidUpdatePermission(nil, invalidUpdatePermissionEvents, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	closeFunc := func() {
		cancel()
		<-done
	}
	t.Cleanup(closeFunc)

	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-invalidUpdatePermissionSub.Err():
				assert.NoError(t, err)
				return
			case evt := <-invalidUpdatePermissionEvents:
				t.Logf("DF Cache received invalid update permission event: %+v", evt)
			}
		}
	}()
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
