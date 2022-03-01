package monitoring

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Chain-level Metrics
	newFeedConfigsDetected = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "new_feed_configs_detected",
			Help: "set the number of feeds monitored every time the set of feed configs read from the RDD has changed",
		},
		[]string{"network_name", "network_id", "chain_id"},
	)
	sendMessageToKafkaFailed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "send_message_to_kafka_failed",
			Help: "number of failed writes to Kafka",
		},
		[]string{"topic", "network_name", "network_id", "chain_id"},
	)
	sendMessageToKafkaSucceeded = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "send_message_to_kafka_succeeded",
			Help: "number of successful writes to Kafka",
		},
		[]string{"topic", "network_name", "network_id", "chain_id"},
	)

	// Feed-level Metrics

	fetchFromSourceFailed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fetch_from_source_failed",
			Help: "number of failed reads from the chain",
		},
		[]string{"source_name", "feed_id", "feed_name", "contract_status", "contract_type", "network_name", "network_id", "chain_id"},
	)
	fetchFromSourceSucceeded = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fetch_from_source_succeeded",
			Help: "number of successful reads from the chain",
		},
		[]string{"source_name", "feed_id", "feed_name", "contract_status", "contract_type", "network_name", "network_id", "chain_id"},
	)
	fetchFromSourceDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "fetch_from_source_duration",
			Help: "time it takes for a source to read the chain and return data",
			Buckets: []float64{
				float64(100 * time.Millisecond),
				float64(500 * time.Millisecond),
				float64(1 * time.Second),
				float64(2 * time.Second),
				float64(5 * time.Second),
				float64(10 * time.Second),
				float64(20 * time.Second),
			},
		},
		[]string{"source_name", "feed_id", "feed_name", "contract_status", "contract_type", "network_name", "network_id", "chain_id"},
	)
)

type ChainMetrics interface {
	SetNewFeedConfigsDetected(numFeeds float64)

	IncSendMessageToKafkaFailed(topic string)
	IncSendMessageToKafkaSucceeded(topic string)
}

func NewChainMetrics(chainConfig ChainConfig) ChainMetrics {
	return &chainMetrics{chainConfig}
}

type chainMetrics struct {
	chainConfig ChainConfig
}

func (c *chainMetrics) SetNewFeedConfigsDetected(numFeeds float64) {
	newFeedConfigsDetected.With(prometheus.Labels{
		"network_name": c.chainConfig.GetNetworkName(),
		"network_id":   c.chainConfig.GetNetworkID(),
		"chain_id":     c.chainConfig.GetChainID(),
	}).Set(numFeeds)
}

func (c *chainMetrics) IncSendMessageToKafkaFailed(topic string) {
	sendMessageToKafkaFailed.With(prometheus.Labels{
		"topic":        topic,
		"network_name": c.chainConfig.GetNetworkName(),
		"network_id":   c.chainConfig.GetNetworkID(),
		"chain_id":     c.chainConfig.GetChainID(),
	}).Inc()
}

func (c *chainMetrics) IncSendMessageToKafkaSucceeded(topic string) {
	sendMessageToKafkaSucceeded.With(prometheus.Labels{
		"topic":        topic,
		"network_name": c.chainConfig.GetNetworkName(),
		"network_id":   c.chainConfig.GetNetworkID(),
		"chain_id":     c.chainConfig.GetChainID(),
	}).Inc()
}

type FeedMetrics interface {
	IncFetchFromSourceFailed(sourceName string)
	IncFetchFromSourceSucceeded(sourceName string)
	ObserveFetchFromSourceDuraction(duration time.Duration, sourceName string)
}

func NewFeedMetrics(chainConfig ChainConfig, feedConfig FeedConfig) FeedMetrics {
	return &feedMetrics{chainConfig, feedConfig}
}

type feedMetrics struct {
	chainConfig ChainConfig
	feedConfig  FeedConfig
}

func (f *feedMetrics) IncFetchFromSourceFailed(sourceName string) {
	fetchFromSourceFailed.With(prometheus.Labels{
		"source_name":     sourceName,
		"feed_id":         f.feedConfig.GetID(),
		"feed_name":       f.feedConfig.GetName(),
		"contract_status": f.feedConfig.GetContractStatus(),
		"contract_type":   f.feedConfig.GetContractType(),
		"network_name":    f.chainConfig.GetNetworkName(),
		"network_id":      f.chainConfig.GetNetworkID(),
		"chain_id":        f.chainConfig.GetChainID(),
	}).Inc()
}

func (f *feedMetrics) IncFetchFromSourceSucceeded(sourceName string) {
	fetchFromSourceSucceeded.With(prometheus.Labels{
		"source_name":     sourceName,
		"feed_id":         f.feedConfig.GetID(),
		"feed_name":       f.feedConfig.GetName(),
		"contract_status": f.feedConfig.GetContractStatus(),
		"contract_type":   f.feedConfig.GetContractType(),
		"network_name":    f.chainConfig.GetNetworkName(),
		"network_id":      f.chainConfig.GetNetworkID(),
		"chain_id":        f.chainConfig.GetChainID(),
	}).Inc()
}

func (f *feedMetrics) ObserveFetchFromSourceDuraction(duration time.Duration, sourceName string) {
	fetchFromSourceDuration.With(prometheus.Labels{
		"source_name":     sourceName,
		"feed_id":         f.feedConfig.GetID(),
		"feed_name":       f.feedConfig.GetName(),
		"contract_status": f.feedConfig.GetContractStatus(),
		"contract_type":   f.feedConfig.GetContractType(),
		"network_name":    f.chainConfig.GetNetworkName(),
		"network_id":      f.chainConfig.GetNetworkID(),
		"chain_id":        f.chainConfig.GetChainID(),
	}).Observe(float64(duration))
}
