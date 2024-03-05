package mercury_common_test

import (
	"context"
	"fmt"
	"math/big"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	mercury_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
)

var ServerFetcher = staticServerFetcher{}

type ServerFetcherEvaluator interface {
	mercury_types.ServerFetcher
	testtypes.Evaluator[mercury_types.ServerFetcher]
}

var _ ServerFetcherEvaluator = staticServerFetcher{}

type staticServerFetcher struct{}

var _ mercury_types.ServerFetcher = staticServerFetcher{}

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

func (staticServerFetcher) FetchInitialMaxFinalizedBlockNumber(context.Context) (*int64, error) {
	return &StaticServerFetcherFixtures.InitialMaxFinalizedBlockNumber, nil
}

func (staticServerFetcher) LatestPrice(ctx context.Context, feedID [32]byte) (*big.Int, error) {
	return StaticServerFetcherFixtures.LatestPrice, nil
}

func (staticServerFetcher) LatestTimestamp(context.Context) (int64, error) {
	return StaticServerFetcherFixtures.LatestTimestamp, nil
}

func (s staticServerFetcher) Evaluate(ctx context.Context, other mercury_types.ServerFetcher) error {
	gotInitialMaxFinalizedBlockNumber, err := other.FetchInitialMaxFinalizedBlockNumber(ctx)
	if err != nil {
		return err
	}
	if *gotInitialMaxFinalizedBlockNumber != StaticServerFetcherFixtures.InitialMaxFinalizedBlockNumber {
		return errMismatch(*gotInitialMaxFinalizedBlockNumber, StaticServerFetcherFixtures.InitialMaxFinalizedBlockNumber)
	}

	gotLatestPrice, err := other.LatestPrice(ctx, [32]byte{})
	if err != nil {
		return err
	}
	if gotLatestPrice.Cmp(StaticServerFetcherFixtures.LatestPrice) != 0 {
		return errMismatch(gotLatestPrice, StaticServerFetcherFixtures.LatestPrice)
	}

	gotLatestTimestamp, err := other.LatestTimestamp(ctx)
	if err != nil {
		return err
	}
	if gotLatestTimestamp != StaticServerFetcherFixtures.LatestTimestamp {
		return errMismatch(gotLatestTimestamp, StaticServerFetcherFixtures.LatestTimestamp)
	}

	return nil
}

func errMismatch(got, expected interface{}) error {
	return errExpected(expected, got)
}

func errExpected(expected, got interface{}) error {
	return &ErrExpected{expected, got}
}

type ErrExpected struct {
	Expected, Got interface{}
}

func (e *ErrExpected) Error() string {
	return fmt.Sprintf("expected %v but got %v", e.Expected, e.Got)
}
