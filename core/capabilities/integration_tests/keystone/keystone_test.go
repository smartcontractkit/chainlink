package keystone

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/integration_tests/framework"
	feeds_consumer "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/feeds_consumer_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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

	workflowDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Workflow", NumNodes: 7, F: 2, AcceptsWorkflows: true})
	require.NoError(t, err)
	triggerDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Trigger", NumNodes: 7, F: 2})
	require.NoError(t, err)
	targetDonConfiguration, err := framework.NewDonConfiguration(framework.NewDonConfigurationParams{Name: "Target", NumNodes: 4, F: 1})
	require.NoError(t, err)

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

	reports := []*datastreams.FeedReport{
		createFeedReport(t, big.NewInt(1), 5, feedIDs[0], triggerDonConfiguration.KeyBundles),
		createFeedReport(t, big.NewInt(3), 7, feedIDs[1], triggerDonConfiguration.KeyBundles),
		createFeedReport(t, big.NewInt(2), 6, feedIDs[2], triggerDonConfiguration.KeyBundles),
	}

	wrappedReports, err := wrapReports(reports, 12, datastreams.Metadata{})
	require.NoError(t, err)

	triggerSink.SendOutput(wrappedReports)

	waitForConsumerReports(ctx, t, consumer, reports)
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

func newFeedID(t *testing.T) string {
	buf := [32]byte{}
	_, err := rand.Read(buf[:])
	require.NoError(t, err)
	return "0x" + hex.EncodeToString(buf[:])
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
