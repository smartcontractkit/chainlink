package logpoller

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type queryType string

const (
	create queryType = "create"
	read   queryType = "read"
	del    queryType = "delete"
)

var (
	sqlLatencyBuckets = []float64{
		float64(1 * time.Millisecond),
		float64(5 * time.Millisecond),
		float64(10 * time.Millisecond),
		float64(20 * time.Millisecond),
		float64(30 * time.Millisecond),
		float64(40 * time.Millisecond),
		float64(50 * time.Millisecond),
		float64(60 * time.Millisecond),
		float64(70 * time.Millisecond),
		float64(80 * time.Millisecond),
		float64(90 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(200 * time.Millisecond),
		float64(300 * time.Millisecond),
		float64(400 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(750 * time.Millisecond),
		float64(1 * time.Second),
		float64(2 * time.Second),
		float64(5 * time.Second),
	}
	lpQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "log_poller_query_duration",
		Help:    "Measures duration of Log Poller's queries fetching logs",
		Buckets: sqlLatencyBuckets,
	}, []string{"evmChainID", "query", "type"})
	lpQueryDataSets = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "log_poller_query_dataset_size",
		Help: "Measures size of the datasets returned by Log Poller's queries",
	}, []string{"evmChainID", "query", "type"})
	lpLogsInserted = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "log_poller_logs_inserted",
		Help: "Counter to track number of logs inserted by Log Poller",
	}, []string{"evmChainID"})
	lpBlockInserted = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "log_poller_blocks_inserted",
		Help: "Counter to track number of blocks inserted by Log Poller",
	}, []string{"evmChainID"})
)

// ObservedORM is a decorator layer for ORM used by LogPoller, responsible for pushing Prometheus metrics reporting duration and size of result set for the queries.
// It doesn't change internal logic, because all calls are delegated to the origin ORM
type ObservedORM struct {
	ORM
	queryDuration  *prometheus.HistogramVec
	datasetSize    *prometheus.GaugeVec
	logsInserted   *prometheus.CounterVec
	blocksInserted *prometheus.CounterVec
	chainId        string
}

// NewObservedORM creates an observed version of log poller's ORM created by NewORM
// Please see ObservedLogPoller for more details on how latencies are measured
func NewObservedORM(chainID *big.Int, db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) *ObservedORM {
	return &ObservedORM{
		ORM:            NewORM(chainID, db, lggr, cfg),
		queryDuration:  lpQueryDuration,
		datasetSize:    lpQueryDataSets,
		logsInserted:   lpLogsInserted,
		blocksInserted: lpBlockInserted,
		chainId:        chainID.String(),
	}
}

func (o *ObservedORM) InsertLogs(logs []Log, qopts ...pg.QOpt) error {
	err := withObservedExec(o, "InsertLogs", create, func() error {
		return o.ORM.InsertLogs(logs, qopts...)
	})
	trackInsertedLogsAndBlock(o, logs, nil, err)
	return err
}

func (o *ObservedORM) InsertLogsWithBlock(logs []Log, block LogPollerBlock, qopts ...pg.QOpt) error {
	err := withObservedExec(o, "InsertLogsWithBlock", create, func() error {
		return o.ORM.InsertLogsWithBlock(logs, block, qopts...)
	})
	trackInsertedLogsAndBlock(o, logs, &block, err)
	return err
}

func (o *ObservedORM) InsertFilter(filter Filter, qopts ...pg.QOpt) error {
	return withObservedExec(o, "InsertFilter", create, func() error {
		return o.ORM.InsertFilter(filter, qopts...)
	})
}

func (o *ObservedORM) LoadFilters(qopts ...pg.QOpt) (map[string]Filter, error) {
	return withObservedQuery(o, "LoadFilters", func() (map[string]Filter, error) {
		return o.ORM.LoadFilters(qopts...)
	})
}

func (o *ObservedORM) DeleteFilter(name string, qopts ...pg.QOpt) error {
	return withObservedExec(o, "DeleteFilter", del, func() error {
		return o.ORM.DeleteFilter(name, qopts...)
	})
}

func (o *ObservedORM) DeleteBlocksBefore(end int64, qopts ...pg.QOpt) error {
	return withObservedExec(o, "DeleteBlocksBefore", del, func() error {
		return o.ORM.DeleteBlocksBefore(end, qopts...)
	})
}

func (o *ObservedORM) DeleteLogsAndBlocksAfter(start int64, qopts ...pg.QOpt) error {
	return withObservedExec(o, "DeleteLogsAndBlocksAfter", del, func() error {
		return o.ORM.DeleteLogsAndBlocksAfter(start, qopts...)
	})
}

func (o *ObservedORM) DeleteExpiredLogs(qopts ...pg.QOpt) error {
	return withObservedExec(o, "DeleteExpiredLogs", del, func() error {
		return o.ORM.DeleteExpiredLogs(qopts...)
	})
}

func (o *ObservedORM) SelectBlockByNumber(n int64, qopts ...pg.QOpt) (*LogPollerBlock, error) {
	return withObservedQuery(o, "SelectBlockByNumber", func() (*LogPollerBlock, error) {
		return o.ORM.SelectBlockByNumber(n, qopts...)
	})
}

func (o *ObservedORM) SelectLatestBlock(qopts ...pg.QOpt) (*LogPollerBlock, error) {
	return withObservedQuery(o, "SelectLatestBlock", func() (*LogPollerBlock, error) {
		return o.ORM.SelectLatestBlock(qopts...)
	})
}

func (o *ObservedORM) SelectLatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs Confirmations, qopts ...pg.QOpt) (*Log, error) {
	return withObservedQuery(o, "SelectLatestLogByEventSigWithConfs", func() (*Log, error) {
		return o.ORM.SelectLatestLogByEventSigWithConfs(eventSig, address, confs, qopts...)
	})
}

func (o *ObservedORM) SelectLogsWithSigs(start, end int64, address common.Address, eventSigs []common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectLogsWithSigs", func() ([]Log, error) {
		return o.ORM.SelectLogsWithSigs(start, end, address, eventSigs, qopts...)
	})
}

func (o *ObservedORM) SelectLogsCreatedAfter(address common.Address, eventSig common.Hash, after time.Time, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectLogsCreatedAfter", func() ([]Log, error) {
		return o.ORM.SelectLogsCreatedAfter(address, eventSig, after, confs, qopts...)
	})
}

func (o *ObservedORM) SelectIndexedLogs(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectIndexedLogs", func() ([]Log, error) {
		return o.ORM.SelectIndexedLogs(address, eventSig, topicIndex, topicValues, confs, qopts...)
	})
}

func (o *ObservedORM) SelectIndexedLogsByBlockRange(start, end int64, address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectIndexedLogsByBlockRange", func() ([]Log, error) {
		return o.ORM.SelectIndexedLogsByBlockRange(start, end, address, eventSig, topicIndex, topicValues, qopts...)
	})
}

func (o *ObservedORM) SelectIndexedLogsCreatedAfter(address common.Address, eventSig common.Hash, topicIndex int, topicValues []common.Hash, after time.Time, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectIndexedLogsCreatedAfter", func() ([]Log, error) {
		return o.ORM.SelectIndexedLogsCreatedAfter(address, eventSig, topicIndex, topicValues, after, confs, qopts...)
	})
}

func (o *ObservedORM) SelectIndexedLogsWithSigsExcluding(sigA, sigB common.Hash, topicIndex int, address common.Address, startBlock, endBlock int64, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectIndexedLogsWithSigsExcluding", func() ([]Log, error) {
		return o.ORM.SelectIndexedLogsWithSigsExcluding(sigA, sigB, topicIndex, address, startBlock, endBlock, confs, qopts...)
	})
}

func (o *ObservedORM) SelectLogs(start, end int64, address common.Address, eventSig common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectLogs", func() ([]Log, error) {
		return o.ORM.SelectLogs(start, end, address, eventSig, qopts...)
	})
}

func (o *ObservedORM) SelectIndexedLogsByTxHash(address common.Address, eventSig common.Hash, txHash common.Hash, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectIndexedLogsByTxHash", func() ([]Log, error) {
		return o.ORM.SelectIndexedLogsByTxHash(address, eventSig, txHash, qopts...)
	})
}

func (o *ObservedORM) GetBlocksRange(start int64, end int64, qopts ...pg.QOpt) ([]LogPollerBlock, error) {
	return withObservedQueryAndResults(o, "GetBlocksRange", func() ([]LogPollerBlock, error) {
		return o.ORM.GetBlocksRange(start, end, qopts...)
	})
}

func (o *ObservedORM) SelectLatestLogEventSigsAddrsWithConfs(fromBlock int64, addresses []common.Address, eventSigs []common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectLatestLogEventSigsAddrsWithConfs", func() ([]Log, error) {
		return o.ORM.SelectLatestLogEventSigsAddrsWithConfs(fromBlock, addresses, eventSigs, confs, qopts...)
	})
}

func (o *ObservedORM) SelectLatestBlockByEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs Confirmations, qopts ...pg.QOpt) (int64, error) {
	return withObservedQuery(o, "SelectLatestBlockByEventSigsAddrsWithConfs", func() (int64, error) {
		return o.ORM.SelectLatestBlockByEventSigsAddrsWithConfs(fromBlock, eventSigs, addresses, confs, qopts...)
	})
}

func (o *ObservedORM) SelectLogsDataWordRange(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin, wordValueMax common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectLogsDataWordRange", func() ([]Log, error) {
		return o.ORM.SelectLogsDataWordRange(address, eventSig, wordIndex, wordValueMin, wordValueMax, confs, qopts...)
	})
}

func (o *ObservedORM) SelectLogsDataWordGreaterThan(address common.Address, eventSig common.Hash, wordIndex int, wordValueMin common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectLogsDataWordGreaterThan", func() ([]Log, error) {
		return o.ORM.SelectLogsDataWordGreaterThan(address, eventSig, wordIndex, wordValueMin, confs, qopts...)
	})
}

func (o *ObservedORM) SelectLogsDataWordBetween(address common.Address, eventSig common.Hash, wordIndexMin int, wordIndexMax int, wordValue common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectLogsDataWordBetween", func() ([]Log, error) {
		return o.ORM.SelectLogsDataWordBetween(address, eventSig, wordIndexMin, wordIndexMax, wordValue, confs, qopts...)
	})
}

func (o *ObservedORM) SelectIndexedLogsTopicGreaterThan(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectIndexedLogsTopicGreaterThan", func() ([]Log, error) {
		return o.ORM.SelectIndexedLogsTopicGreaterThan(address, eventSig, topicIndex, topicValueMin, confs, qopts...)
	})
}

func (o *ObservedORM) SelectIndexedLogsTopicRange(address common.Address, eventSig common.Hash, topicIndex int, topicValueMin, topicValueMax common.Hash, confs Confirmations, qopts ...pg.QOpt) ([]Log, error) {
	return withObservedQueryAndResults(o, "SelectIndexedLogsTopicRange", func() ([]Log, error) {
		return o.ORM.SelectIndexedLogsTopicRange(address, eventSig, topicIndex, topicValueMin, topicValueMax, confs, qopts...)
	})
}

func withObservedQueryAndResults[T any](o *ObservedORM, queryName string, query func() ([]T, error)) ([]T, error) {
	results, err := withObservedQuery(o, queryName, query)
	if err == nil {
		o.datasetSize.
			WithLabelValues(o.chainId, queryName, string(read)).
			Set(float64(len(results)))
	}
	return results, err
}

func withObservedQuery[T any](o *ObservedORM, queryName string, query func() (T, error)) (T, error) {
	queryStarted := time.Now()
	defer func() {
		o.queryDuration.
			WithLabelValues(o.chainId, queryName, string(read)).
			Observe(float64(time.Since(queryStarted)))
	}()
	return query()
}

func withObservedExec(o *ObservedORM, query string, queryType queryType, exec func() error) error {
	queryStarted := time.Now()
	defer func() {
		o.queryDuration.
			WithLabelValues(o.chainId, query, string(queryType)).
			Observe(float64(time.Since(queryStarted)))
	}()
	return exec()
}

func trackInsertedLogsAndBlock(o *ObservedORM, logs []Log, block *LogPollerBlock, err error) {
	if err != nil {
		return
	}
	o.logsInserted.
		WithLabelValues(o.chainId).
		Add(float64(len(logs)))

	if block != nil {
		o.blocksInserted.
			WithLabelValues(o.chainId).
			Inc()
	}
}
