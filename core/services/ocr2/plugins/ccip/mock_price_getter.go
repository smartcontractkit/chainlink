package ccip

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
)

var _ PriceGetter = &mockPriceGetter{}

type mockPriceGetter struct {
	mock.Mock
}

func newMockPriceGetter() *mockPriceGetter {
	return &mockPriceGetter{}
}

func (g *mockPriceGetter) TokenPricesUSD(_ context.Context, tokens []common.Address) (map[common.Address]*big.Int, error) {
	args := g.Called(tokens)
	return args.Get(0).(map[common.Address]*big.Int), args.Error(1)
}
