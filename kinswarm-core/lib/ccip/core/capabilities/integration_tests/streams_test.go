package integration_tests

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	reporttypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/types"
)

func Test_AllAtOnceTransmissionSchedule(t *testing.T) {
	ctx := testutils.Context(t)

	// The don IDs set in the below calls are inferred from the order in which the dons are added to the capabilities registry
	// in the setupCapabilitiesRegistryContract function, should this order change the don IDs will need updating.
	workflowDonInfo := createDonInfo(t, don{id: 1, numNodes: 7, f: 2})
	triggerDonInfo := createDonInfo(t, don{id: 2, numNodes: 7, f: 2})
	targetDonInfo := createDonInfo(t, don{id: 3, numNodes: 4, f: 1})

	consumer, feedIDs, triggerSink := setupStreamDonsWithTransmissionSchedule(ctx, t, workflowDonInfo, triggerDonInfo, targetDonInfo, 3,
		"2s", "allAtOnce")

	reports := []*datastreams.FeedReport{
		createFeedReport(t, big.NewInt(1), 5, feedIDs[0], triggerDonInfo.keyBundles),
		createFeedReport(t, big.NewInt(3), 7, feedIDs[1], triggerDonInfo.keyBundles),
		createFeedReport(t, big.NewInt(2), 6, feedIDs[2], triggerDonInfo.keyBundles),
	}

	triggerSink.sendReports(reports)

	waitForConsumerReports(ctx, t, consumer, reports)
}

func Test_OneAtATimeTransmissionSchedule(t *testing.T) {
	ctx := testutils.Context(t)

	// The don IDs set in the below calls are inferred from the order in which the dons are added to the capabilities registry
	// in the setupCapabilitiesRegistryContract function, should this order change the don IDs will need updating.
	workflowDonInfo := createDonInfo(t, don{id: 1, numNodes: 7, f: 2})
	triggerDonInfo := createDonInfo(t, don{id: 2, numNodes: 7, f: 2})
	targetDonInfo := createDonInfo(t, don{id: 3, numNodes: 4, f: 1})

	consumer, feedIDs, triggerSink := setupStreamDonsWithTransmissionSchedule(ctx, t, workflowDonInfo, triggerDonInfo, targetDonInfo, 3,
		"2s", "oneAtATime")

	reports := []*datastreams.FeedReport{
		createFeedReport(t, big.NewInt(1), 5, feedIDs[0], triggerDonInfo.keyBundles),
		createFeedReport(t, big.NewInt(3), 7, feedIDs[1], triggerDonInfo.keyBundles),
		createFeedReport(t, big.NewInt(2), 6, feedIDs[2], triggerDonInfo.keyBundles),
	}

	triggerSink.sendReports(reports)

	waitForConsumerReports(ctx, t, consumer, reports)
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
			t.Fatalf("timed out waiting for feed reports, expected %d, received %d", len(triggerFeedReports), feedCount)
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
