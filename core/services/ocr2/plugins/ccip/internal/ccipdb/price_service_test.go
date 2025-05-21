package db

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	cciporm "github.com/smartcontractkit/chainlink/v2/core/services/ccip"
	ccipmocks "github.com/smartcontractkit/chainlink/v2/core/services/ccip/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

func TestPriceService_writeGasPrices(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChainSelector := uint64(12345)
	sourceChainSelector := uint64(67890)

	gasPrice := big.NewInt(1e18)

	expectedGasPriceUpdate := []cciporm.GasPrice{
		{
			SourceChainSelector: sourceChainSelector,
			GasPrice:            assets.NewWei(gasPrice),
		},
	}

	testCases := []struct {
		name          string
		gasPriceError bool
		expectedErr   bool
	}{
		{
			name:          "ORM called successfully",
			gasPriceError: false,
			expectedErr:   false,
		},
		{
			name:          "gasPrice clear failed",
			gasPriceError: true,
			expectedErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()

			var gasPricesError error
			if tc.gasPriceError {
				gasPricesError = errors.New("gas prices error")
			}

			mockOrm := ccipmocks.NewORM(t)
			mockOrm.On("UpsertGasPricesForDestChain", ctx, destChainSelector, expectedGasPriceUpdate).Return(int64(0), gasPricesError).Once()

			priceService := NewPriceService(
				lggr,
				mockOrm,
				jobId,
				destChainSelector,
				sourceChainSelector,
				"",
				nil,
				nil,
			).(*priceService)
			err := priceService.writeGasPricesToDB(ctx, gasPrice)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPriceService_writeTokenPrices(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChainSelector := uint64(12345)
	sourceChainSelector := uint64(67890)

	tokenPrices := map[cciptypes.Address]*big.Int{
		"0x123": big.NewInt(2e18),
		"0x234": big.NewInt(3e18),
	}

	expectedTokenPriceUpdate := []cciporm.TokenPrice{
		{
			TokenAddr:  "0x123",
			TokenPrice: assets.NewWei(big.NewInt(2e18)),
		},
		{
			TokenAddr:  "0x234",
			TokenPrice: assets.NewWei(big.NewInt(3e18)),
		},
	}

	testCases := []struct {
		name            string
		tokenPriceError bool
		expectedErr     bool
	}{
		{
			name:            "ORM called successfully",
			tokenPriceError: false,
			expectedErr:     false,
		},
		{
			name:            "tokenPrice clear failed",
			tokenPriceError: true,
			expectedErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()

			var tokenPricesError error
			if tc.tokenPriceError {
				tokenPricesError = errors.New("token prices error")
			}

			mockOrm := ccipmocks.NewORM(t)
			mockOrm.On("UpsertTokenPricesForDestChain", ctx, destChainSelector, expectedTokenPriceUpdate, tokenPriceUpdateInterval).
				Return(int64(len(expectedTokenPriceUpdate)), tokenPricesError).Once()

			priceService := NewPriceService(
				lggr,
				mockOrm,
				jobId,
				destChainSelector,
				sourceChainSelector,
				"",
				nil,
				nil,
			).(*priceService)
			err := priceService.writeTokenPricesToDB(ctx, tokenPrices)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPriceService_observeGasPriceUpdates(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChain := chainselectors.TEST_1338
	sourceChain := chainselectors.TEST_1000
	sourceNativeTokenID := ccipcommon.TokenID{
		TokenAddress:  ccipcalc.EvmAddrToGeneric(utils.RandomAddress()),
		ChainSelector: sourceChain.Selector,
	}

	testCases := []struct {
		name                 string
		priceGetterRespData  map[ccipcommon.TokenID]*big.Int
		priceGetterRespErr   error
		feeEstimatorRespFee  *big.Int
		feeEstimatorRespErr  error
		maxGasPrice          uint64
		expSourceGasPriceUSD *big.Int
		expErr               bool
	}{
		{
			name: "base",
			priceGetterRespData: map[ccipcommon.TokenID]*big.Int{
				sourceNativeTokenID: val1e18(100),
			},
			priceGetterRespErr:   nil,
			feeEstimatorRespFee:  big.NewInt(10),
			feeEstimatorRespErr:  nil,
			maxGasPrice:          1e18,
			expSourceGasPriceUSD: big.NewInt(1000),
			expErr:               false,
		},
		{
			name:                "price getter returned an error",
			priceGetterRespData: nil,
			priceGetterRespErr:  errors.New("some random network error"),
			expErr:              true,
		},
		{
			name: "price getter did not return source native gas price",
			priceGetterRespData: map[ccipcommon.TokenID]*big.Int{
				ccipcommon.TokenID{
					TokenAddress:  sourceNativeTokenID.TokenAddress,
					ChainSelector: destChain.Selector, // the chain selector is from the dest
				}: val1e18(100),
			},
			priceGetterRespErr: nil,
			expErr:             true,
		},
		{
			name: "dynamic fee cap overrides legacy",
			priceGetterRespData: map[ccipcommon.TokenID]*big.Int{
				sourceNativeTokenID: val1e18(100),
			},
			priceGetterRespErr:   nil,
			feeEstimatorRespFee:  big.NewInt(20),
			feeEstimatorRespErr:  nil,
			maxGasPrice:          1e18,
			expSourceGasPriceUSD: big.NewInt(2000),
			expErr:               false,
		},
		{
			name: "nil gas price",
			priceGetterRespData: map[ccipcommon.TokenID]*big.Int{
				sourceNativeTokenID: val1e18(100),
			},
			feeEstimatorRespFee: nil,
			maxGasPrice:         1e18,
			expErr:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			priceGetter := pricegetter.NewMockAllTokensPriceGetter(t)
			defer priceGetter.AssertExpectations(t)

			gasPriceEstimator := prices.NewMockGasPriceEstimatorCommit(t)
			defer gasPriceEstimator.AssertExpectations(t)

			priceGetter.EXPECT().GetTokenPricesUSD(mock.Anything, []ccipcommon.TokenID{sourceNativeTokenID}).
				Return(tc.priceGetterRespData, tc.priceGetterRespErr)

			if tc.maxGasPrice > 0 {
				gasPriceEstimator.On("GetGasPrice", mock.Anything).Return(tc.feeEstimatorRespFee, tc.feeEstimatorRespErr)
				if tc.feeEstimatorRespFee != nil {
					pUSD := ccipcalc.CalculateUsdPerUnitGas(tc.feeEstimatorRespFee, tc.priceGetterRespData[sourceNativeTokenID])
					gasPriceEstimator.On("DenoteInUSD", mock.Anything, mock.Anything, mock.Anything).Return(pUSD, nil)
				}
			}

			priceService := NewPriceService(
				lggr,
				nil,
				jobId,
				destChain.Selector,
				sourceChain.Selector,
				sourceNativeTokenID.TokenAddress,
				priceGetter,
				nil,
			).(*priceService)
			priceService.gasPriceEstimator = gasPriceEstimator

			sourceGasPriceUSD, err := priceService.observeGasPriceUpdates(context.Background(), lggr)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expSourceGasPriceUSD.Cmp(sourceGasPriceUSD), 0)
		})
	}
}

func TestPriceService_observeTokenPriceUpdates(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChain := chainselectors.TEST_1338
	sourceChain := chainselectors.TEST_1000

	cntAddrs := 0
	genTokenAddrDeterministic := func() cciptypes.Address { // ordered token address generator
		cntAddrs++
		return cciptypes.Address(fmt.Sprintf("0x%04d", cntAddrs))
	}

	sourceNativeTokenID := ccipcommon.TokenID{
		TokenAddress:  genTokenAddrDeterministic(),
		ChainSelector: sourceChain.Selector,
	}

	const nTokens = 10
	destTokenIDs := make([]ccipcommon.TokenID, nTokens)
	for i := range destTokenIDs {
		destTokenIDs[i] = ccipcommon.TokenID{
			TokenAddress:  genTokenAddrDeterministic(),
			ChainSelector: destChain.Selector,
		}
	}

	testCases := []struct {
		name                string
		tokenDecimalsParams []cciptypes.Address
		tokenDecimalsResps  []uint8
		sourceNativeToken   ccipcommon.TokenID
		filterOutTokens     []ccipcommon.TokenID
		priceGetterRespData map[ccipcommon.TokenID]*big.Int
		priceGetterRespErr  error
		expTokenPricesUSD   map[cciptypes.Address]*big.Int
		expErr              bool
		expDecimalErr       bool
	}{
		{
			name:                "base case with src native token not equals to dest token address",
			tokenDecimalsParams: []cciptypes.Address{destTokenIDs[1].TokenAddress, destTokenIDs[2].TokenAddress},
			tokenDecimalsResps:  []uint8{18, 12},
			sourceNativeToken:   sourceNativeTokenID,
			priceGetterRespData: map[ccipcommon.TokenID]*big.Int{ // should return all tokens (including source native token)
				sourceNativeTokenID: val1e18(100),
				destTokenIDs[1]:     val1e18(200),
				destTokenIDs[2]:     val1e18(300),
			},
			priceGetterRespErr: nil,
			expTokenPricesUSD: map[cciptypes.Address]*big.Int{ // should only return the tokens in destination chain
				destTokenIDs[1].TokenAddress: val1e18(200),
				destTokenIDs[2].TokenAddress: val1e18(300 * 1e6),
			},
			expErr: false,
		},
		{
			name:                "base case with src native token address equal to a dest token address",
			tokenDecimalsParams: []cciptypes.Address{sourceNativeTokenID.TokenAddress, destTokenIDs[1].TokenAddress},
			tokenDecimalsResps:  []uint8{18, 12},
			sourceNativeToken:   sourceNativeTokenID,
			priceGetterRespData: map[ccipcommon.TokenID]*big.Int{
				ccipcommon.TokenID{
					TokenAddress:  sourceNativeTokenID.TokenAddress,
					ChainSelector: destChain.Selector,
				}: val1e18(100),
				destTokenIDs[1]: val1e18(200),
			},
			priceGetterRespErr: nil,
			expTokenPricesUSD: map[cciptypes.Address]*big.Int{
				sourceNativeTokenID.TokenAddress: val1e18(100),
				destTokenIDs[1].TokenAddress:     val1e18(200 * 1e6),
			},
			expErr: false,
		},
		{
			name:                "price getter returned an error",
			tokenDecimalsParams: []cciptypes.Address{destTokenIDs[0].TokenAddress, destTokenIDs[1].TokenAddress},
			tokenDecimalsResps:  []uint8{18, 18},
			sourceNativeToken:   destTokenIDs[0],
			priceGetterRespData: nil,
			priceGetterRespErr:  errors.New("some random network error"),
			expErr:              true,
		},
		{
			name:                "price getter returns more prices for non-dest tokens",
			tokenDecimalsParams: []cciptypes.Address{destTokenIDs[1].TokenAddress, destTokenIDs[2].TokenAddress, destTokenIDs[3].TokenAddress},
			tokenDecimalsResps:  []uint8{18, 12, 18},
			sourceNativeToken:   sourceNativeTokenID,
			priceGetterRespData: map[ccipcommon.TokenID]*big.Int{
				sourceNativeTokenID: val1e18(100), // <-- not a dest token
				destTokenIDs[1]:     val1e18(200),
				destTokenIDs[2]:     val1e18(300),
				destTokenIDs[3]:     val1e18(400),
			},
			expTokenPricesUSD: map[cciptypes.Address]*big.Int{
				destTokenIDs[1].TokenAddress: val1e18(200),
				destTokenIDs[2].TokenAddress: val1e18(300 * 1e6),
				destTokenIDs[3].TokenAddress: val1e18(400),
			},
		},
		{
			name:                "a token price is nil",
			tokenDecimalsParams: []cciptypes.Address{destTokenIDs[0].TokenAddress, destTokenIDs[1].TokenAddress, destTokenIDs[2].TokenAddress},
			tokenDecimalsResps:  []uint8{18, 18, 18},
			sourceNativeToken:   sourceNativeTokenID,
			filterOutTokens:     []ccipcommon.TokenID{destTokenIDs[2]},
			priceGetterRespData: map[ccipcommon.TokenID]*big.Int{
				destTokenIDs[0]: nil,
				destTokenIDs[1]: val1e18(200),
				destTokenIDs[2]: val1e18(300),
			},
			expErr: true,
		},
		{
			name:                "decimals call errored",
			tokenDecimalsParams: []cciptypes.Address{destTokenIDs[1].TokenAddress, destTokenIDs[2].TokenAddress},
			tokenDecimalsResps:  []uint8{18, 12},
			sourceNativeToken:   sourceNativeTokenID,
			priceGetterRespData: map[ccipcommon.TokenID]*big.Int{ // should return all tokens (including source native token)
				sourceNativeTokenID: val1e18(100),
				destTokenIDs[1]:     val1e18(200),
				destTokenIDs[2]:     val1e18(300),
			},
			priceGetterRespErr: nil,
			expTokenPricesUSD: map[cciptypes.Address]*big.Int{ // should only return the tokens in destination chain
				destTokenIDs[1].TokenAddress: val1e18(200),
				destTokenIDs[2].TokenAddress: val1e18(300 * 1e6),
			},
			expErr:        true,
			expDecimalErr: true,
		},
		{
			name:                "decimals call returned invalid results",
			tokenDecimalsParams: []cciptypes.Address{destTokenIDs[1].TokenAddress, destTokenIDs[2].TokenAddress},
			tokenDecimalsResps:  []uint8{18, 12, 145}, // <-- one extra result
			sourceNativeToken:   sourceNativeTokenID,
			priceGetterRespData: map[ccipcommon.TokenID]*big.Int{ // should return all tokens (including source native token)
				sourceNativeTokenID: val1e18(100),
				destTokenIDs[1]:     val1e18(200),
				destTokenIDs[2]:     val1e18(300),
			},
			priceGetterRespErr: nil,
			expTokenPricesUSD: map[cciptypes.Address]*big.Int{ // should only return the tokens in destination chain
				destTokenIDs[1].TokenAddress: val1e18(200),
				destTokenIDs[2].TokenAddress: val1e18(300 * 1e6),
			},
			expErr:        true,
			expDecimalErr: false,
		},
		{
			name: "src native token address equals dest token address and dest token price missing",
			tokenDecimalsParams: []cciptypes.Address{
				sourceNativeTokenID.TokenAddress, // marks it as a dest token
				destTokenIDs[1].TokenAddress,
				destTokenIDs[2].TokenAddress,
			},
			tokenDecimalsResps: []uint8{12, 18, 12},
			sourceNativeToken:  sourceNativeTokenID,
			priceGetterRespData: map[ccipcommon.TokenID]*big.Int{ // should return all tokens (including source native token)
				sourceNativeTokenID: val1e18(100),
				destTokenIDs[1]:     val1e18(200),
				destTokenIDs[2]:     val1e18(300),
			},
			priceGetterRespErr: nil,
			expTokenPricesUSD: map[cciptypes.Address]*big.Int{ // should only return the tokens in destination chain
				destTokenIDs[1].TokenAddress:     val1e18(200),
				destTokenIDs[2].TokenAddress:     val1e18(300 * 1e6),
				sourceNativeTokenID.TokenAddress: val1e18(100 * 1e6),
			},
			expErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			priceGetter := pricegetter.NewMockAllTokensPriceGetter(t)
			offRampReader := ccipdatamocks.NewOffRampReader(t)
			destPriceReg := ccipdatamocks.NewPriceRegistryReader(t)

			destTokens := make([]cciptypes.Address, len(tc.tokenDecimalsParams))
			copy(destTokens, tc.tokenDecimalsParams)

			offRampReader.EXPECT().GetTokens(mock.Anything).Return(cciptypes.OffRampTokens{
				DestinationTokens: destTokens,
			}, nil).Maybe()
			destPriceReg.EXPECT().GetFeeTokens(mock.Anything).Return(destTokens, nil).Maybe()

			priceGetter.EXPECT().GetJobSpecTokenPricesUSD(mock.Anything).
				Return(tc.priceGetterRespData, tc.priceGetterRespErr)

			if tc.expDecimalErr {
				destPriceReg.EXPECT().GetTokensDecimals(mock.Anything, tc.tokenDecimalsParams).
					Return([]uint8{}, errors.New("token not found")).Maybe()
			} else {
				destPriceReg.EXPECT().GetTokensDecimals(mock.Anything, tc.tokenDecimalsParams).
					Return(tc.tokenDecimalsResps, nil).Maybe()
			}

			priceService := NewPriceService(
				lggr,
				nil,
				jobId,
				destChain.Selector,
				sourceChain.Selector,
				tc.sourceNativeToken.TokenAddress,
				priceGetter,
				offRampReader,
			).(*priceService)
			priceService.destPriceRegistryReader = destPriceReg

			tokenPricesUSD, err := priceService.observeTokenPriceUpdates(context.Background(), lggr)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expTokenPricesUSD, tokenPricesUSD)
		})
	}
}

func TestPriceService_calculateUsdPer1e18TokenAmount(t *testing.T) {
	testCases := []struct {
		name       string
		price      *big.Int
		decimal    uint8
		wantResult *big.Int
	}{
		{
			name:       "18-decimal token, $6.5 per token",
			price:      big.NewInt(65e17),
			decimal:    18,
			wantResult: big.NewInt(65e17),
		},
		{
			name:       "6-decimal token, $1 per token",
			price:      big.NewInt(1e18),
			decimal:    6,
			wantResult: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e12)), // 1e30
		},
		{
			name:       "0-decimal token, $1 per token",
			price:      big.NewInt(1e18),
			decimal:    0,
			wantResult: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e18)), // 1e36
		},
		{
			name:       "36-decimal token, $1 per token",
			price:      big.NewInt(1e18),
			decimal:    36,
			wantResult: big.NewInt(1),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateUsdPer1e18TokenAmount(tt.price, tt.decimal)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}

func TestPriceService_GetGasAndTokenPrices(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChainSelector := uint64(12345)
	sourceChainSelector := uint64(67890)

	token1 := ccipcalc.HexToAddress("0x123")
	token2 := ccipcalc.HexToAddress("0x234")

	gasPrice := big.NewInt(1e18)
	tokenPrices := map[cciptypes.Address]*big.Int{
		token1: big.NewInt(2e18),
		token2: big.NewInt(3e18),
	}

	testCases := []struct {
		name                 string
		ormGasPricesResult   []cciporm.GasPrice
		ormTokenPricesResult []cciporm.TokenPrice

		expectedGasPrices   map[uint64]*big.Int
		expectedTokenPrices map[cciptypes.Address]*big.Int

		gasPriceError   bool
		tokenPriceError bool
		expectedErr     bool
	}{
		{
			name: "ORM called successfully",
			ormGasPricesResult: []cciporm.GasPrice{
				{
					SourceChainSelector: sourceChainSelector,
					GasPrice:            assets.NewWei(gasPrice),
				},
			},
			ormTokenPricesResult: []cciporm.TokenPrice{
				{
					TokenAddr:  string(token1),
					TokenPrice: assets.NewWei(tokenPrices[token1]),
				},
				{
					TokenAddr:  string(token2),
					TokenPrice: assets.NewWei(tokenPrices[token2]),
				},
			},
			expectedGasPrices: map[uint64]*big.Int{
				sourceChainSelector: gasPrice,
			},
			expectedTokenPrices: tokenPrices,
			gasPriceError:       false,
			tokenPriceError:     false,
			expectedErr:         false,
		},
		{
			name: "multiple gas prices with nil token price",
			ormGasPricesResult: []cciporm.GasPrice{
				{
					SourceChainSelector: sourceChainSelector,
					GasPrice:            assets.NewWei(gasPrice),
				},
				{
					SourceChainSelector: sourceChainSelector + 1,
					GasPrice:            assets.NewWei(big.NewInt(200)),
				},
				{
					SourceChainSelector: sourceChainSelector + 2,
					GasPrice:            assets.NewWei(big.NewInt(300)),
				},
			},
			ormTokenPricesResult: nil,
			expectedGasPrices: map[uint64]*big.Int{
				sourceChainSelector:     gasPrice,
				sourceChainSelector + 1: big.NewInt(200),
				sourceChainSelector + 2: big.NewInt(300),
			},
			expectedTokenPrices: map[cciptypes.Address]*big.Int{},
			gasPriceError:       false,
			tokenPriceError:     false,
			expectedErr:         false,
		},
		{
			name:               "multiple token prices with nil gas price",
			ormGasPricesResult: nil,
			ormTokenPricesResult: []cciporm.TokenPrice{
				{
					TokenAddr:  string(token1),
					TokenPrice: assets.NewWei(tokenPrices[token1]),
				},
				{
					TokenAddr:  string(token2),
					TokenPrice: assets.NewWei(tokenPrices[token2]),
				},
			},
			expectedGasPrices:   map[uint64]*big.Int{},
			expectedTokenPrices: tokenPrices,
			gasPriceError:       false,
			tokenPriceError:     false,
			expectedErr:         false,
		},
		{
			name: "nil prices filtered out",
			ormGasPricesResult: []cciporm.GasPrice{
				{
					SourceChainSelector: sourceChainSelector,
					GasPrice:            nil,
				},
				{
					SourceChainSelector: sourceChainSelector + 1,
					GasPrice:            assets.NewWei(gasPrice),
				},
			},
			ormTokenPricesResult: []cciporm.TokenPrice{
				{
					TokenAddr:  string(token1),
					TokenPrice: assets.NewWei(tokenPrices[token1]),
				},
				{
					TokenAddr:  string(token2),
					TokenPrice: nil,
				},
			},
			expectedGasPrices: map[uint64]*big.Int{
				sourceChainSelector + 1: gasPrice,
			},
			expectedTokenPrices: map[cciptypes.Address]*big.Int{
				token1: tokenPrices[token1],
			},
			gasPriceError:   false,
			tokenPriceError: false,
			expectedErr:     false,
		},
		{
			name:            "gasPrice clear failed",
			gasPriceError:   true,
			tokenPriceError: false,
			expectedErr:     true,
		},
		{
			name:            "tokenPrice clear failed",
			gasPriceError:   false,
			tokenPriceError: true,
			expectedErr:     true,
		},
		{
			name:            "both ORM calls failed",
			gasPriceError:   true,
			tokenPriceError: true,
			expectedErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()

			mockOrm := ccipmocks.NewORM(t)
			if tc.gasPriceError {
				mockOrm.On("GetGasPricesByDestChain", ctx, destChainSelector).Return(nil, errors.New("gas prices error")).Once()
			} else {
				mockOrm.On("GetGasPricesByDestChain", ctx, destChainSelector).Return(tc.ormGasPricesResult, nil).Once()
			}
			if tc.tokenPriceError {
				mockOrm.On("GetTokenPricesByDestChain", ctx, destChainSelector).Return(nil, errors.New("token prices error")).Once()
			} else {
				mockOrm.On("GetTokenPricesByDestChain", ctx, destChainSelector).Return(tc.ormTokenPricesResult, nil).Once()
			}

			priceService := NewPriceService(
				lggr,
				mockOrm,
				jobId,
				destChainSelector,
				sourceChainSelector,
				"",
				nil,
				nil,
			).(*priceService)
			gasPricesResult, tokenPricesResult, err := priceService.GetGasAndTokenPrices(ctx, destChainSelector)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedGasPrices, gasPricesResult)
				assert.Equal(t, tc.expectedTokenPrices, tokenPricesResult)
			}
		})
	}
}

func val1e18(val int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(val))
}

func setupORM(t *testing.T) cciporm.ORM {
	t.Helper()

	db := pgtest.NewSqlxDB(t)
	orm, err := cciporm.NewORM(db, logger.TestLogger(t))

	require.NoError(t, err)

	return orm
}

func checkResultLen(t *testing.T, priceService PriceService, destChainSelector uint64, gasCount int, tokenCount int) error {
	ctx := t.Context()
	dbGasResult, dbTokenResult, err := priceService.GetGasAndTokenPrices(ctx, destChainSelector)
	if err != nil {
		return nil
	}
	if len(dbGasResult) != gasCount {
		return fmt.Errorf("expected %d gas prices, got %d", gasCount, len(dbGasResult))
	}
	if len(dbTokenResult) != tokenCount {
		return fmt.Errorf("expected %d token prices, got %d", tokenCount, len(dbTokenResult))
	}
	return nil
}

func TestPriceService_priceWriteInBackground(t *testing.T) {
	lggr := logger.TestLogger(t)
	jobId := int32(1)
	destChain := chainselectors.TEST_1338
	sourceChain := chainselectors.TEST_1000
	ctx := t.Context()

	sourceNative := ccipcommon.TokenID{
		TokenAddress:  ccipcalc.EvmAddrToGeneric(common.HexToAddress("0x001")),
		ChainSelector: sourceChain.Selector,
	}

	destToken1 := ccipcommon.TokenID{
		TokenAddress:  ccipcalc.EvmAddrToGeneric(common.HexToAddress("0x001")), // <-- same addr as source native
		ChainSelector: destChain.Selector,
	}

	destToken2 := ccipcommon.TokenID{
		TokenAddress:  ccipcalc.EvmAddrToGeneric(common.HexToAddress("0x003")),
		ChainSelector: destChain.Selector,
	}

	destToken3 := ccipcommon.TokenID{
		TokenAddress:  ccipcalc.EvmAddrToGeneric(common.HexToAddress("0x004")),
		ChainSelector: destChain.Selector,
	}

	tokens := []ccipcommon.TokenID{sourceNative, destToken1, destToken2, destToken3}
	tokenPrices := []int64{2, 3, 4, 5}

	destTokenAddrs := []cciptypes.Address{
		destToken1.TokenAddress,
		destToken2.TokenAddress,
		destToken3.TokenAddress,
	}
	tokenDecimals := []uint8{18, 18, 18}

	gasPrice := big.NewInt(10)

	orm := setupORM(t)

	priceGetter := pricegetter.NewMockAllTokensPriceGetter(t)
	defer priceGetter.AssertExpectations(t)

	gasPriceEstimator := prices.NewMockGasPriceEstimatorCommit(t)
	defer gasPriceEstimator.AssertExpectations(t)

	priceGetter.EXPECT().GetTokenPricesUSD(mock.Anything, []ccipcommon.TokenID{sourceNative}).
		Return(map[ccipcommon.TokenID]*big.Int{sourceNative: val1e18(tokenPrices[0])}, nil)

	priceGetter.EXPECT().GetJobSpecTokenPricesUSD(mock.Anything).Return(map[ccipcommon.TokenID]*big.Int{
		tokens[0]: val1e18(tokenPrices[0]),
		tokens[1]: val1e18(tokenPrices[1]),
		tokens[2]: val1e18(tokenPrices[2]),
		tokens[3]: val1e18(tokenPrices[3]),
	}, nil)

	gasPriceEstimator.On("GetGasPrice", mock.Anything).Return(gasPrice, nil)
	pUSD := ccipcalc.CalculateUsdPerUnitGas(gasPrice, val1e18(tokenPrices[0]))
	gasPriceEstimator.On("DenoteInUSD", mock.Anything, mock.Anything, mock.Anything).Return(pUSD, nil)

	destPriceReg := ccipdatamocks.NewPriceRegistryReader(t)
	destPriceReg.On("GetTokensDecimals", mock.Anything, destTokenAddrs).Return(tokenDecimals, nil).Maybe()

	priceService := NewPriceService(
		lggr,
		orm,
		jobId,
		destChain.Selector,
		sourceChain.Selector,
		tokens[0].TokenAddress,
		priceGetter,
		nil,
	).(*priceService)

	gasUpdateInterval := 2000 * time.Millisecond
	tokenUpdateInterval := 5000 * time.Millisecond

	// run gas price task every 2 second
	priceService.gasUpdateInterval = gasUpdateInterval
	// run token price task every 5 second
	priceService.tokenUpdateInterval = tokenUpdateInterval

	// initially, db is empty
	assert.NoError(t, checkResultLen(t, priceService, destChain.Selector, 0, 0))

	// starts PriceService in the background
	assert.NoError(t, priceService.Start(ctx))

	// setting dynamicConfig triggers initial price update
	err := priceService.UpdateDynamicConfig(ctx, gasPriceEstimator, destPriceReg)
	assert.NoError(t, err)
	assert.NoError(t, checkResultLen(t, priceService, destChain.Selector, 1, len(destTokenAddrs)))

	assert.NoError(t, priceService.Close())
}
