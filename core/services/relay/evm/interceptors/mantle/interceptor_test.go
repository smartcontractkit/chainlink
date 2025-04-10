package mantle

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

var randomAddress = types.MustEIP55Address(utils.RandomAddress().String())

func TestInterceptor(t *testing.T) {
	ethClient := clienttest.NewClient(t)
	ctx := context.Background()

	_, err := NewInterceptor(ctx, ethClient, nil)
	require.Error(t, err)

	tokenRatio := big.NewInt(10)
	interceptor, err := NewInterceptor(ctx, ethClient, &randomAddress)
	require.NoError(t, err)

	// request token ratio
	ethClient.On("CallContract", ctx, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).
		Return(common.BigToHash(tokenRatio).Bytes(), nil).Once()

	modExecGasPrice, modDAGasPrice, err := interceptor.ModifyGasPriceComponents(ctx, big.NewInt(1), big.NewInt(1))
	require.NoError(t, err)
	require.Equal(t, int64(10), modExecGasPrice.Int64())
	require.Equal(t, int64(10), modDAGasPrice.Int64())

	// second call won't invoke eth client
	modExecGasPrice, modDAGasPrice, err = interceptor.ModifyGasPriceComponents(ctx, big.NewInt(2), big.NewInt(1))
	require.NoError(t, err)
	require.Equal(t, int64(20), modExecGasPrice.Int64())
	require.Equal(t, int64(10), modDAGasPrice.Int64())
}

func TestModifyGasPriceComponents(t *testing.T) {
	testCases := map[string]struct {
		execGasPrice       *big.Int
		daGasPrice         *big.Int
		tokenRatio         *big.Int
		resultExecGasPrice *big.Int
		resultDAGasPrice   *big.Int
	}{
		"regular": {
			execGasPrice:       big.NewInt(1000),
			daGasPrice:         big.NewInt(100),
			resultExecGasPrice: big.NewInt(2000),
			resultDAGasPrice:   big.NewInt(200),
			tokenRatio:         big.NewInt(2),
		},
		"zero DAGasPrice": {
			execGasPrice:       big.NewInt(1000),
			daGasPrice:         big.NewInt(0),
			resultExecGasPrice: big.NewInt(5000),
			resultDAGasPrice:   big.NewInt(0),
			tokenRatio:         big.NewInt(5),
		},
		"zero ExecGasPrice": {
			execGasPrice:       big.NewInt(0),
			daGasPrice:         big.NewInt(10),
			resultExecGasPrice: big.NewInt(0),
			resultDAGasPrice:   big.NewInt(50),
			tokenRatio:         big.NewInt(5),
		},
		"zero token ratio": {
			execGasPrice:       big.NewInt(15),
			daGasPrice:         big.NewInt(10),
			resultExecGasPrice: big.NewInt(0),
			resultDAGasPrice:   big.NewInt(0),
			tokenRatio:         big.NewInt(0),
		},
	}

	for tcName, tc := range testCases {
		t.Run(tcName, func(t *testing.T) {
			ethClient := clienttest.NewClient(t)
			ctx := context.Background()

			interceptor, err := NewInterceptor(ctx, ethClient, &randomAddress)
			require.NoError(t, err)

			// request token ratio
			ethClient.On("CallContract", ctx, mock.IsType(ethereum.CallMsg{}), mock.IsType(&big.Int{})).
				Return(common.BigToHash(tc.tokenRatio).Bytes(), nil).Once()

			modExecGasPrice, modDAGasPrice, err := interceptor.ModifyGasPriceComponents(ctx, tc.execGasPrice, tc.daGasPrice)
			require.NoError(t, err)
			require.Equal(t, tc.resultExecGasPrice.Int64(), modExecGasPrice.Int64())
			require.Equal(t, tc.resultDAGasPrice.Int64(), modDAGasPrice.Int64())
		})
	}
}
