package mercury_common_test

import (
	"context"
	"math/big"

	mercury_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
)

type StaticServerFetcher struct{}

var _ mercury_types.ServerFetcher = StaticServerFetcher{}

type StaticServerFetcherValues struct {
	InitialMaxFinalizedBlockNumber int64
	LatestPrice                    *big.Int
	LatestTimestamp                int64
}

var StaticServerFetcherFixtures = StaticServerFetcherValues{
	InitialMaxFinalizedBlockNumber: 10,
	LatestPrice:                    big.NewInt(100),
	LatestTimestamp:                7,
}

func (StaticServerFetcher) FetchInitialMaxFinalizedBlockNumber(context.Context) (*int64, error) {
	return &StaticServerFetcherFixtures.InitialMaxFinalizedBlockNumber, nil
}

func (StaticServerFetcher) LatestPrice(ctx context.Context, feedID [32]byte) (*big.Int, error) {
	return StaticServerFetcherFixtures.LatestPrice, nil
}

func (StaticServerFetcher) LatestTimestamp(context.Context) (int64, error) {
	return StaticServerFetcherFixtures.LatestTimestamp, nil
}
