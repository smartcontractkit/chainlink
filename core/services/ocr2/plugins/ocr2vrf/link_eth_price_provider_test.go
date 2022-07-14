package ocr2vrf

import (
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/core/services/vrf/mocks"
)

func Test_JuelsPerFeeCoin(t *testing.T) {
	mockAggregator := mocks.NewAggregatorV3Interface(t)

	t.Run("returns juels per fee coin", func(t *testing.T) {
		latestRoundData := aggregator_v3_interface.LatestRoundData{Answer: big.NewInt(10000)}
		mockAggregator.On("LatestRoundData", mock.Anything).Return(latestRoundData, nil).Once()
		p := linkEthPriceProvider{aggregator: mockAggregator}
		price, err := p.JuelsPerFeeCoin()

		require.NoError(t, err)
		assert.Equal(t, int64(10000), price.Int64())
		mockAggregator.AssertExpectations(t)
	})

	t.Run("returns error when contract call fails", func(t *testing.T) {
		latestRoundData := aggregator_v3_interface.LatestRoundData{}
		mockAggregator.On("LatestRoundData", mock.Anything).Return(latestRoundData, errors.New("network failure")).Once()
		p := linkEthPriceProvider{aggregator: mockAggregator}
		price, err := p.JuelsPerFeeCoin()

		require.Error(t, err)
		assert.Nil(t, price)
		assert.Equal(t, "get aggregator latest answer: network failure", err.Error())
		mockAggregator.AssertExpectations(t)
	})
}
