package types

import (
	"context"
	"math/big"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type DataSourceORM interface {
	LatestReport(ctx context.Context, feedID [32]byte) (report []byte, err error)
}

type ReportCodec interface {
	BenchmarkPriceFromReport(ctx context.Context, report ocrtypes.Report) (*big.Int, error)
}

var (
	PriceFeedMissingCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_price_feed_missing",
		Help: "Running count of times mercury tried to query a price feed for billing from mercury server, but it was missing",
	},
		[]string{"queriedFeedID"},
	)
	PriceFeedErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mercury_price_feed_errors",
		Help: "Running count of times mercury tried to query a price feed for billing from mercury server, but got an error",
	},
		[]string{"queriedFeedID"},
	)
)
