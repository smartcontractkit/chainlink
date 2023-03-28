package juelsfeecoin

import (
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/mocks"
)

func Test_JuelsPerFeeCoin(t *testing.T) {
	t.Parallel()

	mockAggregator := mocks.NewAggregatorV3Interface(t)

	t.Run("returns juels per fee coin", func(t *testing.T) {
		latestRoundData := aggregator_v3_interface.LatestRoundData{Answer: big.NewInt(10000)}
		mockAggregator.On("LatestRoundData", mock.Anything).Return(latestRoundData, nil).Once()
		p := linkEthPriceProvider{aggregator: mockAggregator}
		price, err := p.JuelsPerFeeCoin()

		require.NoError(t, err)
		assert.Equal(t, int64(10000), price.Int64())
	})

	t.Run("returns error when contract call fails", func(t *testing.T) {
		latestRoundData := aggregator_v3_interface.LatestRoundData{}
		mockAggregator.On("LatestRoundData", mock.Anything).Return(latestRoundData, errors.New("network failure")).Once()
		p := linkEthPriceProvider{aggregator: mockAggregator}
		price, err := p.JuelsPerFeeCoin()

		require.Error(t, err)
		assert.Nil(t, price)
		assert.Equal(t, "get aggregator latest answer: network failure", err.Error())
	})

	t.Run("returns juels per fee coin", func(t *testing.T) {
		latestRoundData := aggregator_v3_interface.LatestRoundData{Answer: big.NewInt(10000)}
		mockAggregator.On("LatestRoundData", mock.Anything).Return(latestRoundData, nil)
		p := linkEthPriceProvider{aggregator: mockAggregator, stubbed: true}
		price, err := p.JuelsPerFeeCoin()

		require.NoError(t, err)
		assert.Equal(t, int64(0), price.Int64())
		mockAggregator.AssertNotCalled(t, "LatestRoundData")
	})

	t.Run("returns error when contract call fails", func(t *testing.T) {
		latestRoundData := aggregator_v3_interface.LatestRoundData{}
		mockAggregator.On("LatestRoundData", mock.Anything).Return(latestRoundData, errors.New("network failure"))
		p := linkEthPriceProvider{aggregator: mockAggregator, stubbed: true}
		price, err := p.JuelsPerFeeCoin()

		require.NoError(t, err)
		assert.Zero(t, price.Int64())
		mockAggregator.AssertNotCalled(t, "LatestRoundData")
	})
}
