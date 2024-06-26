package reader

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/smartcontractkit/ccipocr3/internal/mocks"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	EthAcc = ocr2types.Account("ETH")
	OpAcc  = ocr2types.Account("OP")
	ArbAcc = ocr2types.Account("ARB")
)

var (
	EthPrice = big.NewInt(100)
	OpPrice  = big.NewInt(10)
	ArbPrice = big.NewInt(1)
)

func TestOnchainTokenPricesReader_GetTokenPricesUSD(t *testing.T) {
	testCases := []struct {
		name          string
		staticPrices  map[ocr2types.Account]big.Int
		inputTokens   []ocr2types.Account
		mockPrices    map[ocr2types.Account]*big.Int
		want          []*big.Int
		errorAccounts []ocr2types.Account
		wantErr       bool
	}{
		{
			name:         "Static price only",
			staticPrices: map[ocr2types.Account]big.Int{EthAcc: *EthPrice, OpAcc: *OpPrice},
			inputTokens:  []ocr2types.Account{EthAcc, OpAcc},
			mockPrices:   map[ocr2types.Account]*big.Int{},
			want:         []*big.Int{EthPrice, OpPrice},
		},
		{
			name:         "On-chain price only",
			staticPrices: map[ocr2types.Account]big.Int{},
			inputTokens:  []ocr2types.Account{ArbAcc, OpAcc, EthAcc},
			mockPrices:   map[ocr2types.Account]*big.Int{OpAcc: OpPrice, ArbAcc: ArbPrice, EthAcc: EthPrice},
			want:         []*big.Int{ArbPrice, OpPrice, EthPrice},
		},
		{
			name:         "Mix of static price and onchain price",
			staticPrices: map[ocr2types.Account]big.Int{EthAcc: *EthPrice},
			inputTokens:  []ocr2types.Account{EthAcc, OpAcc, ArbAcc},
			mockPrices:   map[ocr2types.Account]*big.Int{ArbAcc: ArbPrice, OpAcc: OpPrice},
			want:         []*big.Int{EthPrice, OpPrice, ArbPrice},
		},
		{
			name:          "Missing price should error",
			staticPrices:  map[ocr2types.Account]big.Int{},
			inputTokens:   []ocr2types.Account{ArbAcc, OpAcc, EthAcc},
			mockPrices:    map[ocr2types.Account]*big.Int{OpAcc: OpPrice, ArbAcc: ArbPrice},
			errorAccounts: []ocr2types.Account{EthAcc},
			want:          nil,
			wantErr:       true,
		},
	}

	for _, tc := range testCases {
		contractReader := createMockReader(tc.mockPrices, tc.errorAccounts)
		tokenPricesReader := OnchainTokenPricesReader{
			TokenPriceConfig: TokenPriceConfig{StaticPrices: tc.staticPrices},
			ContractReader:   contractReader,
		}
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := tokenPricesReader.GetTokenPricesUSD(ctx, tc.inputTokens)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, result)
		})
	}

}

func createMockReader(
	mockPrices map[ocr2types.Account]*big.Int, errorAccounts []ocr2types.Account,
) *mocks.ContractReaderMock {
	reader := mocks.NewContractReaderMock()
	for _, acc := range errorAccounts {
		acc := acc
		reader.On(
			"GetLatestValue", mock.Anything, "PriceAggregator", "getTokenPrice", acc, mock.Anything,
		).Return(fmt.Errorf("error"))
	}
	for acc, price := range mockPrices {
		acc := acc
		price := price
		reader.On("GetLatestValue", mock.Anything, "PriceAggregator", "getTokenPrice", acc, mock.Anything).Run(
			func(args mock.Arguments) {
				arg := args.Get(4).(*big.Int)
				arg.Set(price)
			}).Return(nil)
	}
	return reader
}
