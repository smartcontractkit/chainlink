package juelsfeecoin

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/aggregator_v3_interface"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/mocks"
)

func Test_JuelsPerFeeCoin(t *testing.T) {
	t.Parallel()

	t.Run("returns juels per fee coin", func(t *testing.T) {
		mockAggregator := mocks.NewAggregatorV3Interface(t)

		// Start linkEthPriceProvider.
		var p *linkEthPriceProvider
		evmClient := evmclimocks.NewClient(t)
		pInterface, err := NewLinkEthPriceProvider(common.Address{}, evmClient, time.Second/2, time.Second, logger.TestLogger(t))
		require.NoError(t, err)
		p = pInterface.(*linkEthPriceProvider)

		// Assert correct initial price.
		price, err := p.JuelsPerFeeCoin()
		require.NoError(t, err)
		assert.Equal(t, int64(0), price.Int64())

		// Mock its aggregator contract.
		latestRoundData := aggregator_v3_interface.LatestRoundData{Answer: big.NewInt(10000)}
		mockAggregator.On("LatestRoundData", mock.Anything).Return(latestRoundData, nil)
		p.aggregatorLock.Lock()
		p.aggregator = mockAggregator
		p.aggregatorLock.Unlock()

		// Wait until the price is updated.
		time.Sleep(2 * time.Second)

		// Ensure the correct price is returned.
		price, err = p.JuelsPerFeeCoin()
		require.NoError(t, err)
		assert.Equal(t, int64(10000), price.Int64())
	})

	t.Run("returns juels per fee coin (error updating)", func(t *testing.T) {
		mockAggregator := mocks.NewAggregatorV3Interface(t)

		// Start linkEthPriceProvider.
		var p *linkEthPriceProvider
		evmClient := evmclimocks.NewClient(t)
		pInterface, err := NewLinkEthPriceProvider(common.Address{}, evmClient, time.Second/2, time.Second, logger.TestLogger(t))
		require.NoError(t, err)
		p = pInterface.(*linkEthPriceProvider)

		// Assert correct initial price.
		price, err := p.JuelsPerFeeCoin()
		require.NoError(t, err)
		assert.Equal(t, int64(0), price.Int64())

		// Mock its aggregator contract. Error out when fetching latest round data.
		mockAggregator.On("LatestRoundData", mock.Anything).Return(*new(aggregator_v3_interface.LatestRoundData), errors.New("could not fetch"))
		p.aggregatorLock.Lock()
		p.aggregator = mockAggregator
		p.aggregatorLock.Unlock()

		// Wait until the price is updated.
		time.Sleep(2 * time.Second)

		// Ensure the correct price is returned.
		price, err = p.JuelsPerFeeCoin()
		require.NoError(t, err)
		assert.Equal(t, int64(0), price.Int64())
	})

	t.Run("errors out for timeout >= interval", func(t *testing.T) {
		evmClient := evmclimocks.NewClient(t)
		_, err := NewLinkEthPriceProvider(common.Address{}, evmClient, time.Second, time.Second, logger.TestLogger(t))
		require.Error(t, err)
		require.Equal(t, "timeout must be less than interval", err.Error())
	})
}
