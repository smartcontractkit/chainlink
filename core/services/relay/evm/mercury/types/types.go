package types

import (
	"context"
	"math/big"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

//go:generate mockery --quiet --name ChainHeadTracker --output ../mocks/ --case=underscore
type ChainHeadTracker interface {
	Client() evmclient.Client
	HeadTracker() httypes.HeadTracker
}

type DataSourceORM interface {
	LatestReport(ctx context.Context, feedID [32]byte, qopts ...pg.QOpt) (report []byte, err error)
}

type ReportCodec interface {
	BenchmarkPriceFromReport(report ocrtypes.Report) (*big.Int, error)
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
