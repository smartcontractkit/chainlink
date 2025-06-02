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
	chainsel "github.com/smartcontractkit/chain-selectors"
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
				require.Equal(t, tt.expectedValue, prices[ccipcalc.EvmAddrToGeneric(address)])
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
	sourceNative := ccipcalc.EvmAddrToGeneric(common.HexToAddress("0x"))
	sourceChain := chainsel.TEST_1000
	destChain := chainsel.TEST_1338
	ds, err := pricegetter.NewPipelineGetter(source, runner, 1, uuid.New(), "test",
		lggr, sourceNative, sourceChain.Selector, destChain.Selector)
	require.NoError(t, err)
	return ds
}
