package ccip

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type skipReason string

const (
	// reasonNotBlessed describes when a report is skipped due to not being blessed.
	reasonNotBlessed skipReason = "not blessed"

	// reasonAllExecuted describes when a report is skipped due to messages being all executed.
	reasonAllExecuted skipReason = "all executed"
)

var (
	execPluginLabels          = []string{"configDigest"}
	execPluginDurationBuckets = []float64{
		float64(10 * time.Millisecond),
		float64(20 * time.Millisecond),
		float64(50 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(200 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(1 * time.Second),
		float64(2 * time.Second),
		float64(5 * time.Second),
		float64(10 * time.Second),
	}
	metricReportSkipped = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ccip_unexpired_report_skipped",
		Help: "Times report is skipped for the possible reasons",
	}, []string{"reason"})
	execPluginReportsCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_execution_observation_reports_count",
		Help: "Number of reports that are being processed by Execution Plugin during single observation",
	}, execPluginLabels)
	execPluginObservationBuildDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_execution_observation_build_duration",
		Help:    "Duration of generating Observation in Execution Plugin",
		Buckets: execPluginDurationBuckets,
	}, execPluginLabels)
	execPluginBatchBuildDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_execution_build_single_batch",
		Help:    "Duration of building single batch in Execution Plugin",
		Buckets: execPluginDurationBuckets,
	}, execPluginLabels)
	execPluginReportsIterationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_execution_reports_iteration_build_batch",
		Help:    "Duration of iterating over all unexpired reports in Execution Plugin",
		Buckets: execPluginDurationBuckets,
	}, execPluginLabels)
)

func measureExecPluginDuration(histogram *prometheus.HistogramVec, timestamp types.ReportTimestamp, duration time.Duration) {
	histogram.
		WithLabelValues(timestampToLabels(timestamp)...).
		Observe(float64(duration))
}

func measureObservationBuildDuration(timestamp types.ReportTimestamp, duration time.Duration) {
	measureExecPluginDuration(execPluginObservationBuildDuration, timestamp, duration)
}

func measureBatchBuildDuration(timestamp types.ReportTimestamp, duration time.Duration) {
	measureExecPluginDuration(execPluginBatchBuildDuration, timestamp, duration)
}

func measureReportsIterationDuration(timestamp types.ReportTimestamp, duration time.Duration) {
	measureExecPluginDuration(execPluginReportsIterationDuration, timestamp, duration)
}

func measureNumberOfReportsProcessed(timestamp types.ReportTimestamp, count int) {
	execPluginReportsCount.
		WithLabelValues(timestampToLabels(timestamp)...).
		Set(float64(count))
}

func incSkippedRequests(reason skipReason) {
	metricReportSkipped.WithLabelValues(string(reason)).Inc()
}

func timestampToLabels(t types.ReportTimestamp) []string {
	return []string{t.ConfigDigest.Hex()}
}

// ChainName returns the name of the EVM network based on its chainID
func ChainName(chainID int64) string {
	switch chainID {
	case 1:
		return "ethereum-mainnet"
	case 4:
		return "ethereum-testnet-rinkeby"
	case 5:
		return "ethereum-testnet-goerli"
	case 10:
		return "ethereum-mainnet-optimism-1"
	case 42:
		return "ethereum-testnet-kovan"
	case 56:
		return "binance_smart_chain-mainnet"
	case 97:
		return "binance_smart_chain-testnet"
	case 137:
		return "polygon-mainnet"
	case 420:
		return "ethereum-testnet-goerli-optimism-1"
	case 1111:
		return "wemix-mainnet"
	case 1112:
		return "wemix-testnet"
	case 255:
		return "ethereum-mainnet-kroma-1"
	case 2358:
		return "ethereum-testnet-sepolia-kroma-1"
	case 4002:
		return "fantom-testnet"
	case 8453:
		return "ethereum-mainnet-base-1"
	case 84531:
		return "ethereum-testnet-goerli-base-1"
	case 84532:
		return "ethereum-testnet-sepolia-base-1"
	case 42161:
		return "ethereum-mainnet-arbitrum-1"
	case 421613:
		return "ethereum-testnet-goerli-arbitrum-1"
	case 421614:
		return "ethereum-testnet-sepolia-arbitrum-1"
	case 43113:
		return "avalanche-testnet-fuji"
	case 43114:
		return "avalanche-mainnet"
	case 76578:
		return "avalanche-testnet-anz-subnet"
	case 80001:
		return "polygon-testnet-mumbai"
	case 11155111:
		return "ethereum-testnet-sepolia"
	case 11155420:
		return "ethereum-testnet-sepolia-optimism-1"
	default: // Unknown chain, return chainID as string
		return strconv.FormatInt(chainID, 10)
	}
}
