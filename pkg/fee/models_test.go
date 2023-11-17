package fee

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// This test is based on EVM Fixed Fee Estimator.
func TestCalculateFee(t *testing.T) {
	t.Run("CalculateFee returns DefaultPrice when UserSpecifiedMaxFeePrice and MaxFeePriceConfigured are greater", func(t *testing.T) {
		cfg := newMockFeeEstimatorConfig(nil)
		assert.Equal(t, big.NewInt(42), CalculateFee(cfg.UserSpecifiedMaxFeePrice, cfg.DefaultPrice, cfg.MaxFeePriceConfigured))
	})

	t.Run("CalculateFee returns UserSpecifiedMaxFeePrice when it's lower than DefaultPrice and MaxFeePriceConfigured", func(t *testing.T) {
		cfg := newMockFeeEstimatorConfig(&mockFeeEstimatorConfig{
			UserSpecifiedMaxFeePrice: big.NewInt(30),
			DefaultPrice:             big.NewInt(42),
			MultiplierLimit:          1.1,
			MaxFeePriceConfigured:    big.NewInt(35),
		})
		assert.Equal(t, big.NewInt(30), CalculateFee(cfg.UserSpecifiedMaxFeePrice, cfg.DefaultPrice, cfg.MaxFeePriceConfigured))
	})

	t.Run("CalculateFee returns global maximum price", func(t *testing.T) {
		cfg := newMockFeeEstimatorConfig(&mockFeeEstimatorConfig{
			UserSpecifiedMaxFeePrice: big.NewInt(30),
			DefaultPrice:             big.NewInt(42),
			MultiplierLimit:          1.1,
			MaxFeePriceConfigured:    big.NewInt(20),
		})
		assert.Equal(t, big.NewInt(20), CalculateFee(cfg.UserSpecifiedMaxFeePrice, cfg.DefaultPrice, cfg.MaxFeePriceConfigured))
	})
}

func TestCalculateBumpedFee(t *testing.T) {
	lggr := logger.Sugared(logger.Test(t))
	// Create a mock config
	cfg := newMockFeeEstimatorConfig(&mockFeeEstimatorConfig{
		UserSpecifiedMaxFeePrice: big.NewInt(1000000),
		DefaultPrice:             big.NewInt(42),
		MaxFeePriceConfigured:    big.NewInt(1000000),
		MultiplierLimit:          1.1,
	})
	currentFeePrice := cfg.DefaultPrice
	originalFeePrice := big.NewInt(42)
	maxBumpPrice := big.NewInt(1000000)
	bumpMin := big.NewInt(150)
	bumpPercent := uint16(10)

	// Expected results
	expectedFeePrice := big.NewInt(192)

	actualFeePrice, err := CalculateBumpedFee(
		lggr,
		currentFeePrice,
		originalFeePrice,
		cfg.MaxFeePriceConfigured,
		maxBumpPrice,
		bumpMin,
		bumpPercent,
		toChainUnit,
	)
	require.NoError(t, err)

	assert.Equal(t, expectedFeePrice, actualFeePrice)
}

func TestApplyMultiplier(t *testing.T) {
	testCases := []struct {
		cfg   *mockFeeEstimatorConfig
		input int64
		want  int
	}{
		{
			cfg: newMockFeeEstimatorConfig(&mockFeeEstimatorConfig{
				UserSpecifiedMaxFeePrice: big.NewInt(2000000),
				DefaultPrice:             big.NewInt(84),
				MaxFeePriceConfigured:    big.NewInt(2000000),
				MultiplierLimit:          1.2,
			}),
			input: 100000,
			want:  120000,
		},
		{
			cfg:   newMockFeeEstimatorConfig(nil), // default config
			input: 100000,
			want:  110000,
		},
	}

	for _, tc := range testCases {
		got, err := ApplyMultiplier(uint32(tc.input), tc.cfg.MultiplierLimit)
		require.NoError(t, err)
		assert.Equal(t, tc.want, int(got))
	}
}

// type dummyFee big.Int

type mockFeeEstimatorConfig struct {
	UserSpecifiedMaxFeePrice *big.Int
	DefaultPrice             *big.Int
	MaxFeePriceConfigured    *big.Int
	MultiplierLimit          float32
}

// Currently the values are based on EVM Fixed Fee Estimator.
func newMockFeeEstimatorConfig(cfg *mockFeeEstimatorConfig) *mockFeeEstimatorConfig {
	if cfg == nil {
		return &mockFeeEstimatorConfig{
			UserSpecifiedMaxFeePrice: big.NewInt(1000000),
			DefaultPrice:             big.NewInt(42),
			MaxFeePriceConfigured:    big.NewInt(1000000),
			MultiplierLimit:          1.1,
		}
	}

	return cfg
}

func toChainUnit(fee *big.Int) string {
	return fmt.Sprintf("%d chain unit", fee)
}
