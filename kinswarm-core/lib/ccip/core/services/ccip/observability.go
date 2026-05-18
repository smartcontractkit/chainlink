package ccip

import (
	"context"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var (
	sqlLatencyBuckets = []float64{
		float64(10 * time.Millisecond),
		float64(20 * time.Millisecond),
		float64(30 * time.Millisecond),
		float64(40 * time.Millisecond),
		float64(50 * time.Millisecond),
		float64(70 * time.Millisecond),
		float64(90 * time.Millisecond),
		float64(100 * time.Millisecond),
		float64(200 * time.Millisecond),
		float64(300 * time.Millisecond),
		float64(400 * time.Millisecond),
		float64(500 * time.Millisecond),
		float64(750 * time.Millisecond),
		float64(1 * time.Second),
	}
	ccipQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "ccip_orm_query_duration",
		Buckets: sqlLatencyBuckets,
	}, []string{"query", "destChainSelector"})
	ccipQueryDatasets = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ccip_orm_dataset_size",
	}, []string{"query", "destChainSelector"})
)

type observedORM struct {
	ORM
	queryDuration *prometheus.HistogramVec
	datasetSize   *prometheus.GaugeVec
}

var _ ORM = (*observedORM)(nil)

func NewObservedORM(ds sqlutil.DataSource, lggr logger.Logger) (*observedORM, error) {
	delegate, err := NewORM(ds, lggr)
	if err != nil {
		return nil, err
	}

	return &observedORM{
		ORM:           delegate,
		queryDuration: ccipQueryDuration,
		datasetSize:   ccipQueryDatasets,
	}, nil
}

func (o *observedORM) GetGasPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]GasPrice, error) {
	return withObservedQueryAndResults(o, "GetGasPricesByDestChain", destChainSelector, func() ([]GasPrice, error) {
		return o.ORM.GetGasPricesByDestChain(ctx, destChainSelector)
	})
}

func (o *observedORM) GetTokenPricesByDestChain(ctx context.Context, destChainSelector uint64) ([]TokenPrice, error) {
	return withObservedQueryAndResults(o, "GetTokenPricesByDestChain", destChainSelector, func() ([]TokenPrice, error) {
		return o.ORM.GetTokenPricesByDestChain(ctx, destChainSelector)
	})
}

func (o *observedORM) UpsertGasPricesForDestChain(ctx context.Context, destChainSelector uint64, gasPrices []GasPrice) (int64, error) {
	return withObservedQueryAndRowsAffected(o, "UpsertGasPricesForDestChain", destChainSelector, func() (int64, error) {
		return o.ORM.UpsertGasPricesForDestChain(ctx, destChainSelector, gasPrices)
	})
}

func (o *observedORM) UpsertTokenPricesForDestChain(ctx context.Context, destChainSelector uint64, tokenPrices []TokenPrice, interval time.Duration) (int64, error) {
	return withObservedQueryAndRowsAffected(o, "UpsertTokenPricesForDestChain", destChainSelector, func() (int64, error) {
		return o.ORM.UpsertTokenPricesForDestChain(ctx, destChainSelector, tokenPrices, interval)
	})
}

func withObservedQueryAndRowsAffected(o *observedORM, queryName string, chainSelector uint64, query func() (int64, error)) (int64, error) {
	rowsAffected, err := withObservedQuery(o, queryName, chainSelector, query)
	if err == nil {
		o.datasetSize.
			WithLabelValues(queryName, strconv.FormatUint(chainSelector, 10)).
			Set(float64(rowsAffected))
	}
	return rowsAffected, err
}

func withObservedQueryAndResults[T any](o *observedORM, queryName string, chainSelector uint64, query func() ([]T, error)) ([]T, error) {
	results, err := withObservedQuery(o, queryName, chainSelector, query)
	if err == nil {
		o.datasetSize.
			WithLabelValues(queryName, strconv.FormatUint(chainSelector, 10)).
			Set(float64(len(results)))
	}
	return results, err
}

func withObservedQuery[T any](o *observedORM, queryName string, chainSelector uint64, query func() (T, error)) (T, error) {
	queryStarted := time.Now()
	defer func() {
		o.queryDuration.
			WithLabelValues(queryName, strconv.FormatUint(chainSelector, 10)).
			Observe(float64(time.Since(queryStarted)))
	}()
	return query()
}
