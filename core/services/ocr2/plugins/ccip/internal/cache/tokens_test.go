package cache

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func Test_tokenToDecimals(t *testing.T) {
	tokenDecimalsMapping := map[common.Address]uint8{
		common.HexToAddress("0xA"): 10,
		common.HexToAddress("0xB"): 5,
		common.HexToAddress("0xC"): 2,
	}

	tests := []struct {
		name       string
		destTokens []common.Address
		feeTokens  []common.Address
		want       map[common.Address]uint8
		wantErr    bool
	}{
		{
			name:       "empty map for empty tokens from origin",
			destTokens: []common.Address{},
			feeTokens:  []common.Address{},
			want:       map[common.Address]uint8{},
		},
		{
			name:       "separate destination and fee tokens",
			destTokens: []common.Address{common.HexToAddress("0xC")},
			feeTokens:  []common.Address{common.HexToAddress("0xB")},
			want: map[common.Address]uint8{
				common.HexToAddress("0xC"): 2,
				common.HexToAddress("0xB"): 5,
			},
		},
		{
			name:       "fee tokens and dest tokens are overlapping",
			destTokens: []common.Address{common.HexToAddress("0xA")},
			feeTokens:  []common.Address{common.HexToAddress("0xA")},
			want: map[common.Address]uint8{
				common.HexToAddress("0xA"): 10,
			},
		},
		{
			name:       "only fee tokens are returned",
			destTokens: []common.Address{},
			feeTokens:  []common.Address{common.HexToAddress("0xA"), common.HexToAddress("0xC")},
			want: map[common.Address]uint8{
				common.HexToAddress("0xA"): 10,
				common.HexToAddress("0xC"): 2,
			},
		},
		{
			name:       "error on invalid token",
			destTokens: []common.Address{},
			feeTokens:  []common.Address{common.HexToAddress("0xD")},
			want:       map[common.Address]uint8{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offRampReader := ccipdatamocks.NewOffRampReader(t)
			offRampReader.On("GetDestinationTokens", mock.Anything).Return(tt.destTokens, nil)

			decimalsQueryTokens := make([]common.Address, 0)
			tokenDecimals := make([]uint8, 0)
			var queryErr error
			for i := range tt.destTokens {
				decimals, exists := tokenDecimalsMapping[tt.destTokens[i]]
				if !exists {
					queryErr = fmt.Errorf("decimals not found")
				}
				tokenDecimals = append(tokenDecimals, decimals)
				decimalsQueryTokens = append(decimalsQueryTokens, tt.destTokens[i])
			}
			for i := range tt.feeTokens {
				if slices.Contains(decimalsQueryTokens, tt.feeTokens[i]) {
					continue
				}
				decimals, exists := tokenDecimalsMapping[tt.feeTokens[i]]
				if !exists {
					queryErr = fmt.Errorf("decimals not found")
				}
				tokenDecimals = append(tokenDecimals, decimals)
				decimalsQueryTokens = append(decimalsQueryTokens, tt.feeTokens[i])
			}

			priceRegistryReader := ccipdatamocks.NewPriceRegistryReader(t)
			priceRegistryReader.On("GetFeeTokens", mock.Anything).Return(tt.feeTokens, nil)
			if len(decimalsQueryTokens) > 0 {
				priceRegistryReader.On("GetTokensDecimals", mock.Anything, decimalsQueryTokens).Return(tokenDecimals, queryErr).Once()
			}

			tokenToDecimal := &tokenToDecimals{
				lggr:                logger.TestLogger(t),
				offRamp:             offRampReader,
				priceRegistryReader: priceRegistryReader,
			}

			got, err := tokenToDecimal.CallOrigin(testutils.Context(t))
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)

			// we don't expect rpc call to be made, decimals should be in cache.
			got, err = tokenToDecimal.CallOrigin(testutils.Context(t))
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCallOrigin(t *testing.T) {
	src1 := common.HexToAddress("10")
	dst1 := common.HexToAddress("11")
	src2 := common.HexToAddress("20")
	dst2 := common.HexToAddress("21")

	testCases := []struct {
		name       string
		srcToDst   map[common.Address]common.Address
		offRampErr error
		expErr     bool
	}{
		{
			name: "base",
			srcToDst: map[common.Address]common.Address{
				src1: dst1,
				src2: dst2,
			},
			expErr: false,
		},
		{
			name: "off ramp returned an error",
			srcToDst: map[common.Address]common.Address{
				src1: dst1,
				src2: dst2,
			},
			offRampErr: errors.New("some err"),
			expErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			offRampReader := ccipdatamocks.NewOffRampReader(t)
			srcTks := make([]common.Address, 0, len(tc.srcToDst))
			destTks := make([]common.Address, 0, len(tc.srcToDst))

			for sourceTk, destTk := range tc.srcToDst {
				srcTks = append(srcTks, sourceTk)
				destTks = append(destTks, destTk)
			}
			offRampReader.On("GetDestinationTokensFromSourceTokens", mock.Anything, srcTks).Return(destTks, tc.offRampErr)
			offRampReader.On("GetSupportedTokens", mock.Anything).Return(srcTks, nil)
			o := supportedTokensOrigin{offRamp: offRampReader}
			srcToDst, err := o.CallOrigin(context.Background())

			if tc.expErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			for src, dst := range tc.srcToDst {
				assert.Equal(t, dst, srcToDst[src])
			}
		})
	}
}

func Test_copyArray(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		a := []common.Address{common.HexToAddress("1"), common.HexToAddress("2")}
		b := copyArray(a)
		assert.Equal(t, a, b)
		b[0] = common.HexToAddress("3")
		assert.NotEqual(t, a, b)
	})

	t.Run("empty", func(t *testing.T) {
		b := copyArray([]common.Address{})
		assert.Empty(t, b)
	})
}

func Test_copyMap(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		val := map[string]int{"a": 100, "b": 50}
		cp := copyMap(val)
		assert.Len(t, val, 2)
		assert.Equal(t, 100, cp["a"])
		assert.Equal(t, 50, cp["b"])
		val["b"] = 10
		assert.Equal(t, 50, cp["b"])
	})

	t.Run("pointer val", func(t *testing.T) {
		val := map[string]*big.Int{"a": big.NewInt(100), "b": big.NewInt(50)}
		cp := copyMap(val)
		val["a"] = big.NewInt(20)
		assert.Equal(t, int64(100), cp["a"].Int64())
	})
}

func Test_cachedDecimals(t *testing.T) {
	tokenDecimalsCache := &tokenToDecimals{}
	addr := utils.RandomAddress()

	decimals, exists := tokenDecimalsCache.getCachedDecimals(addr)
	assert.Zero(t, decimals)
	assert.False(t, exists)

	tokenDecimalsCache.setCachedDecimals(addr, 123)
	decimals, exists = tokenDecimalsCache.getCachedDecimals(addr)
	assert.Equal(t, uint8(123), decimals)
	assert.True(t, exists)
}
