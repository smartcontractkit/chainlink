package monitoring

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics is a thin interface on top of the prometheus API.
// As such there should be little logic in the implementation of these methods.
type Metrics interface {
	SetHeadTrackerCurrentHead(blockNumber float64, networkName, chainID, networkID string)
	SetFeedContractMetadata(chainID, contractAddress, feedID, contractStatus, contractType, feedName, feedPath, networkID, networkName, symbol string)
	SetFeedContractLinkBalance(balance float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string)
	SetFeedContractTransmissionsSucceeded(numSucceeded float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string)
	SetFeedContractTransmissionsFailed(numFailed float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string)
	SetNodeMetadata(chainID, networkID, networkName, oracleName, sender string)
	// Deprecated: use SetOffchainAggregatorAnswers
	SetOffchainAggregatorAnswersRaw(answer float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string)
	SetOffchainAggregatorAnswers(answer float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string)
	IncOffchainAggregatorAnswersTotal(contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string)
	// Deprecated: use SetOffchainAggregatorJuelsPerFeeCoin
	SetOffchainAggregatorJuelsPerFeeCoinRaw(juelsPerFeeCoin float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string)
	SetOffchainAggregatorJuelsPerFeeCoin(juelsPerFeeCoin float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string)
	SetOffchainAggregatorSubmissionReceivedValues(value float64, contractAddress, feedID, sender, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string)
	SetOffchainAggregatorAnswerStalled(isSet bool, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string)
	SetOffchainAggregatorRoundID(aggregatorRoundID float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string)
	// Cleanup deletes all the metrics
	Cleanup(networkName, networkID, chainID, oracleName, sender, feedName, feedPath, symbol, contractType, contractStatus, contractAddress, feedID string)
	// Exposes the accumulated metrics to HTTP.
	HTTPHandler() http.Handler
}

var (
	headTrackerCurrentHead = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "head_tracker_current_head",
			Help: "Tracks the current block height that the monitoring instance has processed.",
		},
		[]string{"network_name", "chain_id", "network_id"},
	)
	feedContractMetadata = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "feed_contract_metadata",
			Help: "Exposes metadata for individual feeds. It should simply be set to 1, as the relevant info is in the labels.",
		},
		[]string{"chain_id", "contract_address", "feed_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name", "symbol"},
	)
	feedContractLinkBalance = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "feed_contract_link_balance",
		},
		[]string{"contract_address", "feed_id", "chain_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name"},
	)
	feedContractTransmissionsSucceeded = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "feed_contract_transmissions_succeeded",
		},
		[]string{"contract_address", "feed_id", "chain_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name"},
	)
	feedContractTransmissionsFailed = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "feed_contract_transmissions_failed",
		},
		[]string{"contract_address", "feed_id", "chain_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name"},
	)
	nodeMetadata = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_metadata",
			Help: "Exposes metadata for node operators. It should simply be set to 1, as the relevant info is in the labels.",
		},
		[]string{"chain_id", "network_id", "network_name", "oracle_name", "sender"},
	)
	offchainAggregatorAnswersRaw = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "offchain_aggregator_answers_raw",
			Help: "Reports the latest answer for a contract.",
		},
		[]string{"contract_address", "feed_id", "chain_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name"},
	)
	offchainAggregatorAnswers = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "offchain_aggregator_answers",
			Help: "Reports the latest answer for a contract divided by the feed's Multiply parameter.",
		},
		[]string{"contract_address", "feed_id", "chain_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name"},
	)
	offchainAggregatorAnswersTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "offchain_aggregator_answers_total",
			Help: "Bump this metric every time there is a transmission on chain.",
		},
		[]string{"contract_address", "feed_id", "chain_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name"},
	)
	offchainAggregatorJuelsPerFeeCoinRaw = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "offchain_aggregator_juels_per_fee_coin_raw",
			Help: "Reports the latest raw answer for juels/fee_coin.",
		},
		[]string{"contract_address", "feed_id", "chain_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name"},
	)
	offchainAggregatorJuelsPerFeeCoin = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "offchain_aggregator_juels_per_fee_coin",
			Help: "Reports the latest raw answer for juels/fee_coin divided by the feed's multiplier.",
		},
		[]string{"contract_address", "feed_id", "chain_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name"},
	)
	offchainAggregatorSubmissionReceivedValues = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "offchain_aggregator_submission_received_values",
			Help: "Report individual node observations for the latest transmission on chain. (Should be 1 time series per node per contract). The values are divided by the feed's multiplier config.",
		},
		[]string{"contract_address", "feed_id", "sender", "chain_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name"},
	)
	offchainAggregatorAnswerStalled = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "offchain_aggregator_answer_stalled",
			Help: "Set to 1 if the heartbeat interval has passed on a feed without a transmission. Set to 0 otherwise.",
		},
		[]string{"contract_address", "feed_id", "chain_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name"},
	)
	offchainAggregatorRoundID = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "offchain_aggregator_round_id",
			Help: "Sets the aggregator contract's round id.",
		},
		[]string{"contract_address", "feed_id", "chain_id", "contract_status", "contract_type", "feed_name", "feed_path", "network_id", "network_name"},
	)
)

func NewMetrics(log Logger) Metrics {
	return &defaultMetrics{log}
}

type defaultMetrics struct {
	log Logger
}

func (d *defaultMetrics) SetHeadTrackerCurrentHead(blockNumber float64, networkName, chainID, networkID string) {
	headTrackerCurrentHead.With(prometheus.Labels{
		"network_name": networkName,
		"chain_id":     chainID,
		"network_id":   networkID,
	}).Set(blockNumber)
}

func (d *defaultMetrics) SetFeedContractMetadata(chainID, contractAddress, feedID, contractStatus, contractType, feedName, feedPath, networkID, networkName, symbol string) {
	feedContractMetadata.With(prometheus.Labels{
		"chain_id":         chainID,
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
		"symbol":           symbol,
	}).Set(1)
}

func (d *defaultMetrics) SetFeedContractLinkBalance(balance float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	feedContractLinkBalance.With(prometheus.Labels{
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"chain_id":         chainID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
	}).Set(balance)
}

func (d *defaultMetrics) SetFeedContractTransmissionsSucceeded(numSucceeded float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	feedContractTransmissionsSucceeded.With(prometheus.Labels{
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"chain_id":         chainID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
	}).Set(numSucceeded)
}

func (d *defaultMetrics) SetFeedContractTransmissionsFailed(numFailed float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	feedContractTransmissionsFailed.With(prometheus.Labels{
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"chain_id":         chainID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
	}).Set(numFailed)
}

func (d *defaultMetrics) SetNodeMetadata(chainID, networkID, networkName, oracleName, sender string) {
	nodeMetadata.With(prometheus.Labels{
		"chain_id":     chainID,
		"network_id":   networkID,
		"network_name": networkName,
		"oracle_name":  oracleName,
		"sender":       sender,
	}).Set(1)
}

func (d *defaultMetrics) SetOffchainAggregatorAnswersRaw(answer float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	offchainAggregatorAnswersRaw.With(prometheus.Labels{
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"chain_id":         chainID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
	}).Set(answer)
}

func (d *defaultMetrics) SetOffchainAggregatorAnswers(answer float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	offchainAggregatorAnswers.With(prometheus.Labels{
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"chain_id":         chainID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
	}).Set(answer)
}

func (d *defaultMetrics) IncOffchainAggregatorAnswersTotal(contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	offchainAggregatorAnswersTotal.With(prometheus.Labels{
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"chain_id":         chainID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
	}).Inc()
}

func (d *defaultMetrics) SetOffchainAggregatorJuelsPerFeeCoinRaw(juelsPerFeeCoin float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	offchainAggregatorJuelsPerFeeCoinRaw.With(prometheus.Labels{
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"chain_id":         chainID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
	}).Set(juelsPerFeeCoin)
}

func (d *defaultMetrics) SetOffchainAggregatorJuelsPerFeeCoin(juelsPerFeeCoin float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	offchainAggregatorJuelsPerFeeCoin.With(prometheus.Labels{
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"chain_id":         chainID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
	}).Set(juelsPerFeeCoin)
}

func (d *defaultMetrics) SetOffchainAggregatorSubmissionReceivedValues(value float64, contractAddress, feedID, sender, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	offchainAggregatorSubmissionReceivedValues.With(prometheus.Labels{
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"sender":           sender,
		"chain_id":         chainID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
	}).Set(value)
}

func (d *defaultMetrics) SetOffchainAggregatorAnswerStalled(isSet bool, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	var value float64 = 0
	if isSet {
		value = 1
	}
	offchainAggregatorAnswerStalled.With(prometheus.Labels{
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"chain_id":         chainID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
	}).Set(value)
}

func (d *defaultMetrics) SetOffchainAggregatorRoundID(aggregatorRoundID float64, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	offchainAggregatorRoundID.With(prometheus.Labels{
		"contract_address": contractAddress,
		"feed_id":          feedID,
		"chain_id":         chainID,
		"contract_status":  contractStatus,
		"contract_type":    contractType,
		"feed_name":        feedName,
		"feed_path":        feedPath,
		"network_id":       networkID,
		"network_name":     networkName,
	}).Set(aggregatorRoundID)
}

func (d *defaultMetrics) Cleanup(
	networkName, networkID, chainID, oracleName, sender string,
	feedName, feedPath, symbol, contractType, contractStatus string,
	contractAddress, feedID string,
) {
	for _, metric := range []struct {
		name   string
		vec    *prometheus.MetricVec
		labels prometheus.Labels
	}{
		{
			"head_tracker_current_head",
			headTrackerCurrentHead.MetricVec,
			prometheus.Labels{
				"network_name": networkName,
				"chain_id":     chainID,
				"network_id":   networkID,
			},
		},
		{
			"feed_contract_metadata",
			feedContractMetadata.MetricVec,
			prometheus.Labels{
				"chain_id":         chainID,
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
				"symbol":           symbol,
			},
		},
		{
			"feed_contract_link_balance",
			feedContractLinkBalance.MetricVec,
			prometheus.Labels{
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"chain_id":         chainID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
			},
		},
		{
			"feed_contract_transmissions_succeeded",
			feedContractTransmissionsSucceeded.MetricVec,
			prometheus.Labels{
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"chain_id":         chainID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
			},
		},
		{
			"feed_contract_transmissions_failed",
			feedContractTransmissionsSucceeded.MetricVec,
			prometheus.Labels{
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"chain_id":         chainID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
			},
		},
		{
			"metric node_metadata",
			nodeMetadata.MetricVec,
			prometheus.Labels{
				"chain_id":     chainID,
				"network_id":   networkID,
				"network_name": networkName,
				"oracle_name":  oracleName,
				"sender":       sender,
			},
		},
		{
			"offchain_aggregator_answers_raw",
			offchainAggregatorAnswersRaw.MetricVec,
			prometheus.Labels{
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"chain_id":         chainID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
			},
		},
		{
			"offchain_aggregator_answers",
			offchainAggregatorAnswers.MetricVec,
			prometheus.Labels{
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"chain_id":         chainID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
			},
		},
		{
			"offchain_aggregator_answers_total",
			offchainAggregatorAnswersTotal.MetricVec,
			prometheus.Labels{
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"chain_id":         chainID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
			},
		},
		{
			"offchain_aggregator_juels_per_fee_coin_raw",
			offchainAggregatorJuelsPerFeeCoinRaw.MetricVec,
			prometheus.Labels{
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"chain_id":         chainID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
			},
		},
		{
			"offchain_aggregator_juels_per_fee_coin",
			offchainAggregatorJuelsPerFeeCoin.MetricVec,
			prometheus.Labels{
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"chain_id":         chainID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
			},
		},
		{
			"offchain_aggregator_submission_received_values",
			offchainAggregatorSubmissionReceivedValues.MetricVec,
			prometheus.Labels{
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"sender":           sender,
				"chain_id":         chainID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
			},
		},
		{
			"offchain_aggregator_answer_stalled",
			offchainAggregatorAnswerStalled.MetricVec,
			prometheus.Labels{
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"chain_id":         chainID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
			},
		},
		{
			"offchain_aggregator_round_id",
			offchainAggregatorRoundID.MetricVec,
			prometheus.Labels{
				"contract_address": contractAddress,
				"feed_id":          feedID,
				"chain_id":         chainID,
				"contract_status":  contractStatus,
				"contract_type":    contractType,
				"feed_name":        feedName,
				"feed_path":        feedPath,
				"network_id":       networkID,
				"network_name":     networkName,
			},
		},
	} {
		if !metric.vec.Delete(metric.labels) {
			errArgs := []interface{}{}
			for key, value := range metric.labels {
				errArgs = append(errArgs, key, value)
			}
			d.log.Errorw(fmt.Sprintf("unable to delete metric '%s'", metric.name), errArgs...)
		}
	}
}

func (d *defaultMetrics) HTTPHandler() http.Handler {
	return promhttp.Handler()
}
