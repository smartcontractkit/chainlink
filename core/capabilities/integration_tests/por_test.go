package integration_tests

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	scheduleEverySecond = "* * * * * *"
)

func Test_Cron_OneAtATimeTransmissionSchedule(t *testing.T) {
	ctx := testutils.Context(t)

	// The don IDs set in the below calls are inferred from the order in which the dons are added to the capabilities registry
	// in the setupCapabilitiesRegistryContract function, should this order change the don IDs will need updating.
	workflowDonInfo := createDonInfo(t, don{id: 1, numNodes: 7, f: 2})
	triggerDonInfo := createDonInfo(t, don{id: 2, numNodes: 7, f: 2})
	targetDonInfo := createDonInfo(t, don{id: 3, numNodes: 4, f: 1})

	lggr := setupDonsWithTransmissionSchedulePoR(ctx, t, workflowDonInfo, triggerDonInfo, targetDonInfo, scheduleEverySecond, "2s", "oneAtATime")

	waitForLogs(ctx, t, lggr, 5)
}

func waitForLogs(ctx context.Context, t *testing.T, lggr logger.SugaredLogger, expectedNumRuns int) {
	// feedsReceived := make(chan *feeds_consumer.KeystoneFeedsConsumerFeedReceived, 1000)
	// feedsSub, err := consumer.WatchFeedReceived(&bind.WatchOpts{}, feedsReceived, nil)
	// require.NoError(t, err)

	// feedToReport := map[string]*datastreams.FeedReport{}
	// for _, report := range triggerFeedReports {
	// 	feedToReport[report.FeedID] = report
	// }

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	feedCount := 0
	for {
		select {
		case <-ctxWithTimeout.Done():
			t.Fatalf("timed out waiting for runs, expected %d, received %d", expectedNumRuns, feedCount)
			// case err := <-feedsSub.Err():
			// 	require.NoError(t, err)
			// case feed := <-feedsReceived:
			// 	feedID := "0x" + hex.EncodeToString(feed.FeedId[:])
			// 	report := feedToReport[feedID]
			// 	decodedReport, err := reporttypes.Decode(report.FullReport)
			// 	require.NoError(t, err)
			// 	assert.Equal(t, decodedReport.BenchmarkPrice, feed.Price)
			// 	assert.Equal(t, decodedReport.ObservationsTimestamp, feed.Timestamp)

			// 	feedCount++
			// 	if feedCount == len(triggerFeedReports) {
			// 		return
			// 	}
		}
	}
}
