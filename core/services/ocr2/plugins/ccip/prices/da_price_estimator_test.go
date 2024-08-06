package prices

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/rollups/mocks"
)

func encodeGasPrice(daPrice, execPrice *big.Int) *big.Int {
	return new(big.Int).Add(new(big.Int).Lsh(daPrice, daGasPriceEncodingLength), execPrice)
}

func TestDAPriceEstimator_GetGasPrice(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name         string
		daGasPrice   *big.Int
		execGasPrice *big.Int
		expPrice     *big.Int
		expErr       bool
	}{
		{
			name:         "base",
			daGasPrice:   big.NewInt(1),
			execGasPrice: big.NewInt(0),
			expPrice:     encodeGasPrice(big.NewInt(1), big.NewInt(0)),
			expErr:       false,
		},
		{
			name:         "large values",
			daGasPrice:   big.NewInt(1e9),   // 1 gwei
			execGasPrice: big.NewInt(200e9), // 200 gwei
			expPrice:     encodeGasPrice(big.NewInt(1e9), big.NewInt(200e9)),
			expErr:       false,
		},
		{
			name:         "zero DA price",
			daGasPrice:   big.NewInt(0),
			execGasPrice: big.NewInt(200e9),
			expPrice:     encodeGasPrice(big.NewInt(0), big.NewInt(200e9)),
			expErr:       false,
		},
		{
			name:         "zero exec price",
			daGasPrice:   big.NewInt(1e9),
			execGasPrice: big.NewInt(0),
			expPrice:     encodeGasPrice(big.NewInt(1e9), big.NewInt(0)),
			expErr:       false,
		},
		{
			name:         "price out of bounds",
			daGasPrice:   new(big.Int).Lsh(big.NewInt(1), daGasPriceEncodingLength),
			execGasPrice: big.NewInt(1),
			expPrice:     nil,
			expErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			execEstimator := NewMockGasPriceEstimator(t)
			execEstimator.On("GetGasPrice", ctx).Return(tc.execGasPrice, nil)

			l1Oracle := mocks.NewL1Oracle(t)
			l1Oracle.On("GasPrice", ctx).Return(assets.NewWei(tc.daGasPrice), nil)

			g := DAGasPriceEstimator{
				execEstimator:       execEstimator,
				l1Oracle:            l1Oracle,
				priceEncodingLength: daGasPriceEncodingLength,
			}

			gasPrice, err := g.GetGasPrice(ctx)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expPrice, gasPrice)
		})
	}

	t.Run("nil L1 oracle", func(t *testing.T) {
		expPrice := big.NewInt(1)

		execEstimator := NewMockGasPriceEstimator(t)
		execEstimator.On("GetGasPrice", ctx).Return(expPrice, nil)

		g := DAGasPriceEstimator{
			execEstimator:       execEstimator,
			l1Oracle:            nil,
			priceEncodingLength: daGasPriceEncodingLength,
		}

		gasPrice, err := g.GetGasPrice(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expPrice, gasPrice)
	})
}

func TestDAPriceEstimator_DenoteInUSD(t *testing.T) {
	val1e18 := func(val int64) *big.Int { return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(val)) }

	testCases := []struct {
		name        string
		gasPrice    *big.Int
		nativePrice *big.Int
		expPrice    *big.Int
	}{
		{
			name:        "base",
			gasPrice:    encodeGasPrice(big.NewInt(1e9), big.NewInt(10e9)),
			nativePrice: val1e18(2_000),
			expPrice:    encodeGasPrice(big.NewInt(2000e9), big.NewInt(20000e9)),
		},
		{
			name:        "low price truncates to 0",
			gasPrice:    encodeGasPrice(big.NewInt(1e9), big.NewInt(10e9)),
			nativePrice: big.NewInt(1),
			expPrice:    big.NewInt(0),
		},
		{
			name:        "high price",
			gasPrice:    encodeGasPrice(val1e18(1), val1e18(10)),
			nativePrice: val1e18(2000),
			expPrice:    encodeGasPrice(val1e18(2_000), val1e18(20_000)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := DAGasPriceEstimator{
				priceEncodingLength: daGasPriceEncodingLength,
			}

			gasPrice, err := g.DenoteInUSD(tc.gasPrice, tc.nativePrice)
			assert.NoError(t, err)
			assert.True(t, tc.expPrice.Cmp(gasPrice) == 0)
		})
	}
}

func TestDAPriceEstimator_Median(t *testing.T) {
	val1e18 := func(val int64) *big.Int { return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(val)) }

	testCases := []struct {
		name      string
		gasPrices []*big.Int
		expMedian *big.Int
	}{
		{
			name: "base",
			gasPrices: []*big.Int{
				encodeGasPrice(big.NewInt(1), big.NewInt(1)),
				encodeGasPrice(big.NewInt(2), big.NewInt(2)),
				encodeGasPrice(big.NewInt(3), big.NewInt(3)),
			},
			expMedian: encodeGasPrice(big.NewInt(2), big.NewInt(2)),
		},
		{
			name: "median 2",
			gasPrices: []*big.Int{
				encodeGasPrice(big.NewInt(1), big.NewInt(1)),
				encodeGasPrice(big.NewInt(2), big.NewInt(2)),
			},
			expMedian: encodeGasPrice(big.NewInt(2), big.NewInt(2)),
		},
		{
			name: "large values",
			gasPrices: []*big.Int{
				encodeGasPrice(val1e18(5), val1e18(5)),
				encodeGasPrice(val1e18(4), val1e18(4)),
				encodeGasPrice(val1e18(3), val1e18(3)),
				encodeGasPrice(val1e18(2), val1e18(2)),
				encodeGasPrice(val1e18(1), val1e18(1)),
			},
			expMedian: encodeGasPrice(val1e18(3), val1e18(3)),
		},
		{
			name:      "zeros",
			gasPrices: []*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0)},
			expMedian: big.NewInt(0),
		},
		{
			name: "picks median of each price component individually",
			gasPrices: []*big.Int{
				encodeGasPrice(val1e18(1), val1e18(3)),
				encodeGasPrice(val1e18(2), val1e18(2)),
				encodeGasPrice(val1e18(3), val1e18(1)),
			},
			expMedian: encodeGasPrice(val1e18(2), val1e18(2)),
		},
		{
			name: "unsorted even number of price components",
			gasPrices: []*big.Int{
				encodeGasPrice(val1e18(1), val1e18(22)),
				encodeGasPrice(val1e18(4), val1e18(33)),
				encodeGasPrice(val1e18(2), val1e18(44)),
				encodeGasPrice(val1e18(3), val1e18(11)),
			},
			expMedian: encodeGasPrice(val1e18(3), val1e18(33)),
		},
		{
			name: "equal DA price components",
			gasPrices: []*big.Int{
				encodeGasPrice(val1e18(2), val1e18(22)),
				encodeGasPrice(val1e18(2), val1e18(33)),
				encodeGasPrice(val1e18(2), val1e18(44)),
				encodeGasPrice(val1e18(2), val1e18(11)),
			},
			expMedian: encodeGasPrice(val1e18(2), val1e18(33)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := DAGasPriceEstimator{
				priceEncodingLength: daGasPriceEncodingLength,
			}

			gasPrice, err := g.Median(tc.gasPrices)
			assert.NoError(t, err)
			assert.True(t, tc.expMedian.Cmp(gasPrice) == 0)
		})
	}
}

func TestDAPriceEstimator_Deviates(t *testing.T) {
	testCases := []struct {
		name             string
		gasPrice1        *big.Int
		gasPrice2        *big.Int
		daDeviationPPB   int64
		execDeviationPPB int64
		expDeviates      bool
	}{
		{
			name:             "base",
			gasPrice1:        encodeGasPrice(big.NewInt(100e8), big.NewInt(100e8)),
			gasPrice2:        encodeGasPrice(big.NewInt(79e8), big.NewInt(79e8)),
			daDeviationPPB:   2e8,
			execDeviationPPB: 2e8,
			expDeviates:      true,
		},
		{
			name:             "negative difference also deviates",
			gasPrice1:        encodeGasPrice(big.NewInt(100e8), big.NewInt(100e8)),
			gasPrice2:        encodeGasPrice(big.NewInt(121e8), big.NewInt(121e8)),
			daDeviationPPB:   2e8,
			execDeviationPPB: 2e8,
			expDeviates:      true,
		},
		{
			name:             "only DA component deviates",
			gasPrice1:        encodeGasPrice(big.NewInt(100e8), big.NewInt(100e8)),
			gasPrice2:        encodeGasPrice(big.NewInt(150e8), big.NewInt(110e8)),
			daDeviationPPB:   2e8,
			execDeviationPPB: 2e8,
			expDeviates:      true,
		},
		{
			name:             "only exec component deviates",
			gasPrice1:        encodeGasPrice(big.NewInt(100e8), big.NewInt(100e8)),
			gasPrice2:        encodeGasPrice(big.NewInt(110e8), big.NewInt(150e8)),
			daDeviationPPB:   2e8,
			execDeviationPPB: 2e8,
			expDeviates:      true,
		},
		{
			name:             "both do not deviate",
			gasPrice1:        encodeGasPrice(big.NewInt(100e8), big.NewInt(100e8)),
			gasPrice2:        encodeGasPrice(big.NewInt(110e8), big.NewInt(110e8)),
			daDeviationPPB:   2e8,
			execDeviationPPB: 2e8,
			expDeviates:      false,
		},
		{
			name:             "zero DA price and exec deviates",
			gasPrice1:        encodeGasPrice(big.NewInt(0), big.NewInt(100e8)),
			gasPrice2:        encodeGasPrice(big.NewInt(0), big.NewInt(121e8)),
			daDeviationPPB:   2e8,
			execDeviationPPB: 2e8,
			expDeviates:      true,
		},
		{
			name:             "zero DA price and exec does not deviate",
			gasPrice1:        encodeGasPrice(big.NewInt(0), big.NewInt(100e8)),
			gasPrice2:        encodeGasPrice(big.NewInt(0), big.NewInt(110e8)),
			daDeviationPPB:   2e8,
			execDeviationPPB: 2e8,
			expDeviates:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := DAGasPriceEstimator{
				execEstimator: ExecGasPriceEstimator{
					deviationPPB: tc.execDeviationPPB,
				},
				daDeviationPPB:      tc.daDeviationPPB,
				priceEncodingLength: daGasPriceEncodingLength,
			}

			deviated, err := g.Deviates(tc.gasPrice1, tc.gasPrice2)
			assert.NoError(t, err)
			if tc.expDeviates {
				assert.True(t, deviated)
			} else {
				assert.False(t, deviated)
			}
		})
	}
}

func TestDAPriceEstimator_EstimateMsgCostUSD(t *testing.T) {
	execCostUSD := big.NewInt(100_000)

	testCases := []struct {
		name               string
		gasPrice           *big.Int
		wrappedNativePrice *big.Int
		msg                cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta
		daOverheadGas      int64
		gasPerDAByte       int64
		daMultiplier       int64
		expUSD             *big.Int
	}{
		{
			name:               "only DA overhead",
			gasPrice:           encodeGasPrice(big.NewInt(1e9), big.NewInt(0)), // 1 gwei DA price, 0 exec price
			wrappedNativePrice: big.NewInt(1e18),                               // $1
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					Data:            []byte{},
					TokenAmounts:    []cciptypes.TokenAmount{},
					SourceTokenData: [][]byte{},
				},
			},
			daOverheadGas: 100_000,
			gasPerDAByte:  0,
			daMultiplier:  10_000, // 1x multiplier
			expUSD:        new(big.Int).Add(execCostUSD, big.NewInt(100_000e9)),
		},
		{
			name:               "include message data gas",
			gasPrice:           encodeGasPrice(big.NewInt(1e9), big.NewInt(0)), // 1 gwei DA price, 0 exec price
			wrappedNativePrice: big.NewInt(1e18),                               // $1
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					Data:         make([]byte, 1_000),
					TokenAmounts: make([]cciptypes.TokenAmount, 5),
					SourceTokenData: [][]byte{
						make([]byte, 10), make([]byte, 10), make([]byte, 10), make([]byte, 10), make([]byte, 10),
					},
				},
			},
			daOverheadGas: 100_000,
			gasPerDAByte:  16,
			daMultiplier:  10_000, // 1x multiplier
			expUSD:        new(big.Int).Add(execCostUSD, big.NewInt(134_208e9)),
		},
		{
			name:               "zero DA price",
			gasPrice:           big.NewInt(0),    // 1 gwei DA price, 0 exec price
			wrappedNativePrice: big.NewInt(1e18), // $1
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					Data:            []byte{},
					TokenAmounts:    []cciptypes.TokenAmount{},
					SourceTokenData: [][]byte{},
				},
			},
			daOverheadGas: 100_000,
			gasPerDAByte:  16,
			daMultiplier:  10_000, // 1x multiplier
			expUSD:        execCostUSD,
		},
		{
			name:               "double native price",
			gasPrice:           encodeGasPrice(big.NewInt(1e9), big.NewInt(0)), // 1 gwei DA price, 0 exec price
			wrappedNativePrice: big.NewInt(2e18),                               // $1
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					Data:            []byte{},
					TokenAmounts:    []cciptypes.TokenAmount{},
					SourceTokenData: [][]byte{},
				},
			},
			daOverheadGas: 100_000,
			gasPerDAByte:  0,
			daMultiplier:  10_000, // 1x multiplier
			expUSD:        new(big.Int).Add(execCostUSD, big.NewInt(200_000e9)),
		},
		{
			name:               "half multiplier",
			gasPrice:           encodeGasPrice(big.NewInt(1e9), big.NewInt(0)), // 1 gwei DA price, 0 exec price
			wrappedNativePrice: big.NewInt(1e18),                               // $1
			msg: cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					Data:            []byte{},
					TokenAmounts:    []cciptypes.TokenAmount{},
					SourceTokenData: [][]byte{},
				},
			},
			daOverheadGas: 100_000,
			gasPerDAByte:  0,
			daMultiplier:  5_000, // 0.5x multiplier
			expUSD:        new(big.Int).Add(execCostUSD, big.NewInt(50_000e9)),
		},
	}

	for _, tc := range testCases {
		execEstimator := NewMockGasPriceEstimator(t)
		execEstimator.On("EstimateMsgCostUSD", mock.Anything, tc.wrappedNativePrice, tc.msg).Return(execCostUSD, nil)

		t.Run(tc.name, func(t *testing.T) {
			g := DAGasPriceEstimator{
				execEstimator:       execEstimator,
				l1Oracle:            nil,
				priceEncodingLength: daGasPriceEncodingLength,
				daOverheadGas:       tc.daOverheadGas,
				gasPerDAByte:        tc.gasPerDAByte,
				daMultiplier:        tc.daMultiplier,
			}

			costUSD, err := g.EstimateMsgCostUSD(tc.gasPrice, tc.wrappedNativePrice, tc.msg)
			assert.NoError(t, err)
			assert.Equal(t, tc.expUSD, costUSD)
		})
	}
}
