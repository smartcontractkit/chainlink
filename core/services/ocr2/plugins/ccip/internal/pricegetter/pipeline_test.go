package pricegetter_test

import (
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	config2 "github.com/smartcontractkit/chainlink-common/pkg/config"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"

	pipelinemocks "github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"

	config "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestDataSource(t *testing.T) {
	ctx := testutils.Context(t)
	linkEth := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(`{"JuelsPerETH": "200000000000000000000"}`))
		require.NoError(t, err)
	}))
	defer linkEth.Close()
	usdcEth := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(`{"USDCWeiPerETH": "1000000000000000000000"}`)) // 1000 USDC / ETH
		require.NoError(t, err)
	}))
	defer usdcEth.Close()
	linkTokenAddress := ccipcalc.HexToAddress("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e05")
	usdcTokenAddress := ccipcalc.HexToAddress("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e10")
	source := fmt.Sprintf(`
	// Price 1
	link [type=http method=GET url="%s"];
	link_parse [type=jsonparse path="JuelsPerETH"];
	link->link_parse;
	// Price 2
	usdc [type=http method=GET url="%s"];
	usdc_parse [type=jsonparse path="USDCWeiPerETH"];
	usdc->usdc_parse;
	merge [type=merge left="{}" right="{\"%s\":$(link_parse), \"%s\":$(usdc_parse)}"];
`, linkEth.URL, usdcEth.URL, linkTokenAddress, usdcTokenAddress)

	priceGetter := newTestPipelineGetter(t, source)

	// Ask for all prices present in spec.
	prices, err := priceGetter.GetJobSpecTokenPricesUSD(ctx)
	require.NoError(t, err)
	assert.Equal(t, prices, map[cciptypes.Address]*big.Int{
		linkTokenAddress: big.NewInt(0).Mul(big.NewInt(200), big.NewInt(1000000000000000000)),
		usdcTokenAddress: big.NewInt(0).Mul(big.NewInt(1000), big.NewInt(1000000000000000000)),
	})

	// Specifically ask for all prices
	pricesWithInput, errWithInput := priceGetter.TokenPricesUSD(ctx, []cciptypes.Address{
		linkTokenAddress,
		usdcTokenAddress,
	})
	require.NoError(t, errWithInput)
	assert.Equal(t, pricesWithInput, map[cciptypes.Address]*big.Int{
		linkTokenAddress: big.NewInt(0).Mul(big.NewInt(200), big.NewInt(1000000000000000000)),
		usdcTokenAddress: big.NewInt(0).Mul(big.NewInt(1000), big.NewInt(1000000000000000000)),
	})

	// Ask a non-existent price.
	_, err = priceGetter.TokenPricesUSD(ctx, []cciptypes.Address{
		ccipcalc.HexToAddress("0x1591690b8638f5fb2dbec82ac741805ac5da8b45dc5263f4875b0496fdce4e11"),
	})
	require.Error(t, err)

	// Ask only one price
	prices, err = priceGetter.TokenPricesUSD(ctx, []cciptypes.Address{linkTokenAddress})
	require.NoError(t, err)
	assert.Equal(t, prices, map[cciptypes.Address]*big.Int{
		linkTokenAddress: big.NewInt(0).Mul(big.NewInt(200), big.NewInt(1000000000000000000)),
	})
}

func TestParsingDifferentFormats(t *testing.T) {
	tests := []struct {
		name          string
		inputValue    string
		expectedValue *big.Int
		expectedError bool
	}{
		{
			name:          "number as string",
			inputValue:    "\"200000000000000000000\"",
			expectedValue: new(big.Int).Mul(big.NewInt(200), big.NewInt(1e18)),
		},
		{
			name:          "number as big number",
			inputValue:    "500000000000000000000",
			expectedValue: new(big.Int).Mul(big.NewInt(500), big.NewInt(1e18)),
		},
		{
			name:          "number as int64",
			inputValue:    "150",
			expectedValue: big.NewInt(150),
		},
		{
			name:          "number in scientific notation",
			inputValue:    "3e22",
			expectedValue: new(big.Int).Mul(big.NewInt(30000), big.NewInt(1e18)),
		},
		{
			name:          "number as string in scientific notation returns error",
			inputValue:    "\"3e22\"",
			expectedError: true,
		},
		{
			name:          "invalid value should return error",
			inputValue:    "\"NaN\"",
			expectedError: true,
		},
		{
			name:          "null should return error",
			inputValue:    "null",
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := testutils.Context(t)
			token := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, err := fmt.Fprintf(w, `{"MyCoin": %s}`, tt.inputValue)
				require.NoError(t, err)
			}))
			defer token.Close()

			address := common.HexToAddress("0x94025780a1aB58868D9B2dBBB775f44b32e8E6e5")
			source := fmt.Sprintf(`
			// Price 1
			coin [type=http method=GET url="%s"];
			coin_parse [type=jsonparse path="MyCoin"];
			coin->coin_parse;
			merge [type=merge left="{}" right="{\"%s\":$(coin_parse)}"];
			`, token.URL, strings.ToLower(address.String()))

			prices, err := newTestPipelineGetter(t, source).
				TokenPricesUSD(ctx, []cciptypes.Address{ccipcalc.EvmAddrToGeneric(address)})

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, prices[ccipcalc.EvmAddrToGeneric(address)], tt.expectedValue)
			}
		})
	}
}

func newTestPipelineGetter(t *testing.T, source string) *pricegetter.PipelineGetter {
	lggr, _ := logger.NewLogger()
	cfg := pipelinemocks.NewConfig(t)
	cfg.On("MaxRunDuration").Return(time.Second)
	cfg.On("DefaultHTTPTimeout").Return(*config2.MustNewDuration(time.Second))
	cfg.On("DefaultHTTPLimit").Return(int64(1024 * 10))
	cfg.On("VerboseLogging").Return(true)
	db := pgtest.NewSqlxDB(t)
	bridgeORM := bridges.NewORM(db)
	runner := pipeline.NewRunner(pipeline.NewORM(db, lggr, config.NewTestGeneralConfig(t).JobPipeline().MaxSuccessfulRuns()),
		bridgeORM, cfg, nil, nil, nil, nil, lggr, &http.Client{}, &http.Client{})
	ds, err := pricegetter.NewPipelineGetter(source, runner, 1, uuid.New(), "test", lggr)
	require.NoError(t, err)
	return ds
}
