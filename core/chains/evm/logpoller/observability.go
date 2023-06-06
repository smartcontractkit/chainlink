package logpoller

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

var (
	sqlLatencyBuckets = []float64{
		float64(1 * time.Millisecond),
		float64(5 * time.Millisecond),
		float64(10 * time.Millisecond),
		float64(25 * time.Millisecond),
		float64(50 * time.Millisecond),
		float64(75 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(250 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(750 * time.Millisecond),
		float64(1 * time.Second),
		float64(2 * time.Second),
		float64(5 * time.Second),
		float64(7 * time.Second),
		float64(10 * time.Second),
	}
	lpQueryHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "log_poller_query_duration",
		Help:    "Measures duration of Log Poller's queries fetching logs",
		Buckets: sqlLatencyBuckets,
	}, []string{"query", "address"})
)

// ObservedLogPoller is a decorator layer for LogPoller, responsible for adding pushing histogram metrics for some of the queries.
// It doesn't change internal logic, because all calls are delegated to the origin LogPoller
type ObservedLogPoller struct {
	LogPoller
	histogram *prometheus.HistogramVec
}

// NewObservedLogPoller creates an observed version of log poller created by NewLogPoller
// Please see ObservedLogPoller for more details on how latencies are measured
func NewObservedLogPoller(orm *ORM, ec Client, lggr logger.Logger, pollPeriod time.Duration,
	finalityDepth int64, backfillBatchSize int64, rpcBatchSize int64, keepBlocksDepth int64) LogPoller {

	return &ObservedLogPoller{
		LogPoller: NewLogPoller(orm, ec, lggr, pollPeriod, finalityDepth, backfillBatchSize, rpcBatchSize, keepBlocksDepth),
		histogram: lpQueryHistogram,
	}
}

func (o *ObservedLogPoller) LogsCreatedAfter(eventSig common.Hash, address common.Address, after time.Time, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQuery(o.histogram, "LogsCreatedAfter", address, func() ([]Log, error) {
		return o.LogPoller.LogsCreatedAfter(eventSig, address, after, confs, qopts...)
	})
}

func (o *ObservedLogPoller) LatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs int, qopts ...pg.QOpt) (*Log, error) {
	return withObservedQuery(o.histogram, "LatestLogByEventSigWithConfs", common.Address{}, func() (*Log, error) {
		return o.LogPoller.LatestLogByEventSigWithConfs(eventSig, address, confs, qopts...)
	})
}

func (o *ObservedLogPoller) LatestLogEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQuery(o.histogram, "LatestLogEventSigsAddrsWithConfs", common.Address{}, func() ([]Log, error) {
		return o.LogPoller.LatestLogEventSigsAddrsWithConfs(fromBlock, eventSigs, addresses, confs, qopts...)
	})
}

func (o *ObservedLogPoller) IndexedLogs(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQuery(o.histogram, "IndexedLogs", address, func() ([]Log, error) {
		return o.LogPoller.IndexedLogs(eventSig, address, topicIndex, topicValues, confs, qopts...)
	})
}

func (o *ObservedLogPoller) IndexedLogsByBlockRange(start, end int64, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQuery(o.histogram, "IndexedLogsByBlockRange", address, func() ([]Log, error) {
		return o.LogPoller.IndexedLogsByBlockRange(start, end, eventSig, address, topicIndex, topicValues, qopts...)
	})
}

func (o *ObservedLogPoller) IndexedLogsCreatedAfter(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, after time.Time, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQuery(o.histogram, "IndexedLogsCreatedAfter", address, func() ([]Log, error) {
		return o.LogPoller.IndexedLogsCreatedAfter(eventSig, address, topicIndex, topicValues, after, confs, qopts...)
	})
}

func (o *ObservedLogPoller) IndexedLogsTopicGreaterThan(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQuery(o.histogram, "IndexedLogsTopicGreaterThan", address, func() ([]Log, error) {
		return o.LogPoller.IndexedLogsTopicGreaterThan(eventSig, address, topicIndex, topicValueMin, confs, qopts...)
	})
}

func (o *ObservedLogPoller) IndexedLogsTopicRange(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQuery(o.histogram, "IndexedLogsTopicRange", address, func() ([]Log, error) {
		return o.LogPoller.IndexedLogsTopicRange(eventSig, address, topicIndex, topicValueMin, topicValueMax, confs, qopts...)
	})
}

func (o *ObservedLogPoller) IndexedLogsWithSigsExcluding(address common.Address, eventSigA, eventSigB common.Hash, topicIndex int, fromBlock, toBlock int64, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQuery(o.histogram, "IndexedLogsWithSigsExcluding", address, func() ([]Log, error) {
		return o.LogPoller.IndexedLogsWithSigsExcluding(address, eventSigA, eventSigB, topicIndex, fromBlock, toBlock, confs, qopts...)
	})
}

func (o *ObservedLogPoller) LogsDataWordRange(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin, wordValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQuery(o.histogram, "LogsDataWordRange", address, func() ([]Log, error) {
		return o.LogPoller.LogsDataWordRange(eventSig, address, wordIndex, wordValueMin, wordValueMax, confs, qopts...)
	})
}

func (o *ObservedLogPoller) LogsDataWordGreaterThan(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQuery(o.histogram, "LogsDataWordGreaterThan", address, func() ([]Log, error) {
		return o.LogPoller.LogsDataWordGreaterThan(eventSig, address, wordIndex, wordValueMin, confs, qopts...)
	})
}

func withObservedQuery[T any](histogram *prometheus.HistogramVec, queryName string, address common.Address, query func() (T, error)) (T, error) {
	queryStarted := time.Now()
	defer func() {
		histogram.
			WithLabelValues(queryName, address.String()).
			Observe(float64(time.Since(queryStarted)))
	}()
	return query()
}
