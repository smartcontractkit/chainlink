package ccip

import (
	"math/big"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_MetricsAreTrackedForAllMethods(t *testing.T) {
	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	ccipORM, err := NewObservedORM(db, logger.TestLogger(t))
	require.NoError(t, err)

	tokenPrices := []TokenPrice{
		{
			TokenAddr:  "0xA",
			TokenPrice: assets.NewWei(big.NewInt(1e18)),
		},
		{
			TokenAddr:  "0xB",
			TokenPrice: assets.NewWei(big.NewInt(1e18)),
		},
	}
	tokensUpserted, err := ccipORM.UpsertTokenPricesForDestChain(ctx, 100, tokenPrices, time.Second)
	require.NoError(t, err)
	assert.Equal(t, len(tokenPrices), int(tokensUpserted))
	assert.Equal(t, len(tokenPrices), counterFromGaugeByLabels(ccipORM.datasetSize, "UpsertTokenPricesForDestChain", "100"))
	assert.Equal(t, 0, counterFromGaugeByLabels(ccipORM.datasetSize, "UpsertTokenPricesForDestChain", "200"))

	tokens, err := ccipORM.GetTokenPricesByDestChain(ctx, 100)
	require.NoError(t, err)
	assert.Equal(t, len(tokenPrices), len(tokens))
	assert.Equal(t, len(tokenPrices), counterFromGaugeByLabels(ccipORM.datasetSize, "GetTokenPricesByDestChain", "100"))
	assert.Equal(t, 1, counterFromHistogramByLabels(t, ccipORM.queryDuration, "GetTokenPricesByDestChain", "100"))

	gasPrices := []GasPrice{
		{
			SourceChainSelector: 200,
			GasPrice:            assets.NewWei(big.NewInt(1e18)),
		},
		{
			SourceChainSelector: 201,
			GasPrice:            assets.NewWei(big.NewInt(1e18)),
		},
		{
			SourceChainSelector: 202,
			GasPrice:            assets.NewWei(big.NewInt(1e18)),
		},
	}
	gasUpserted, err := ccipORM.UpsertGasPricesForDestChain(ctx, 100, gasPrices)
	require.NoError(t, err)
	assert.Equal(t, len(gasPrices), int(gasUpserted))
	assert.Equal(t, len(gasPrices), counterFromGaugeByLabels(ccipORM.datasetSize, "UpsertGasPricesForDestChain", "100"))
	assert.Equal(t, 0, counterFromGaugeByLabels(ccipORM.datasetSize, "UpsertGasPricesForDestChain", "200"))

	gas, err := ccipORM.GetGasPricesByDestChain(ctx, 100)
	require.NoError(t, err)
	assert.Equal(t, len(gasPrices), len(gas))
	assert.Equal(t, len(gasPrices), counterFromGaugeByLabels(ccipORM.datasetSize, "GetGasPricesByDestChain", "100"))
	assert.Equal(t, 1, counterFromHistogramByLabels(t, ccipORM.queryDuration, "GetGasPricesByDestChain", "100"))
}

func counterFromHistogramByLabels(t *testing.T, histogramVec *prometheus.HistogramVec, labels ...string) int {
	observer, err := histogramVec.GetMetricWithLabelValues(labels...)
	require.NoError(t, err)

	metricCh := make(chan prometheus.Metric, 1)
	observer.(prometheus.Histogram).Collect(metricCh)
	close(metricCh)

	metric := <-metricCh
	pb := &io_prometheus_client.Metric{}
	err = metric.Write(pb)
	require.NoError(t, err)

	return int(pb.GetHistogram().GetSampleCount())
}

func counterFromGaugeByLabels(gaugeVec *prometheus.GaugeVec, labels ...string) int {
	value := testutil.ToFloat64(gaugeVec.WithLabelValues(labels...))
	return int(value)
}
