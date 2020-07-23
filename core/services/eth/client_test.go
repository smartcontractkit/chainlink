package eth_test

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"math/big"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEthClient_TransactionReceipt(t *testing.T) {
	response := cltest.MustReadFile(t, "testdata/getTransactionReceipt.json")
	mockServer, wsCleanup := cltest.NewWSServer(string(response))
	defer wsCleanup()
	config := cltest.NewConfigWithWSServer(t, mockServer)
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	ec := store.TxManager.(*strpkg.EthTxManager).Client
	err := ec.Dial(context.Background())
	require.NoError(t, err)

	hash := common.HexToHash("0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238")
	receipt, err := ec.TransactionReceipt(context.Background(), hash)
	assert.NoError(t, err)
	assert.Equal(t, hash, receipt.TxHash)
	assert.Equal(t, big.NewInt(11), receipt.BlockNumber)
}

func TestTxReceipt_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantLogLen int
	}{
		{"basic", "testdata/getTransactionReceipt.json", 0},
		{"runlog request", "testdata/runlogReceipt.json", 4},
		{"runlog response", "testdata/responseReceipt.json", 2},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			jsonStr := cltest.JSONFromFixture(t, test.path).Get("result").String()
			var receipt models.TxReceipt
			err := json.Unmarshal([]byte(jsonStr), &receipt)
			require.NoError(t, err)

			assert.Equal(t, test.wantLogLen, len(receipt.Logs))
		})
	}
}

func TestEthClient_PendingNonceAt(t *testing.T) {
	t.Parallel()

	server, cleanup := cltest.NewWSServer(`{
      "id": 1,
      "jsonrpc": "2.0",
      "result": "0x100"
    }`)
	defer cleanup()
	u, err := url.Parse(server.URL)
	require.NoError(t, err)
	u.Scheme = "ws"

	ethClient := eth.NewClient(u.String())
	err = ethClient.Dial(context.Background())
	require.NoError(t, err)

	address := cltest.NewAddress()

	result, err := ethClient.PendingNonceAt(context.Background(), address)
	require.NoError(t, err)

	var expected uint64 = 256
	require.Equal(t, result, expected)
}

func TestEthClient_SendRawTx(t *testing.T) {
	t.Parallel()

	rpcClient := new(mocks.RPCClient)
	ethClient := eth.NewClientWith(rpcClient, nil)
	txData := hexutil.MustDecode("0xdeadbeef")
	returnedHash := cltest.NewHash()

	rpcClient.On("Call", mock.Anything, "eth_sendRawTransaction", "0xdeadbeef").
		Return(nil).
		Run(func(args mock.Arguments) {
			res := args.Get(0).(*common.Hash)
			*res = returnedHash
		})

	result, err := ethClient.SendRawTx(txData)
	assert.NoError(t, err)
	assert.Equal(t, result, returnedHash)
}

func TestEthClient_GetEthBalance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic", "0x0100", "0.000000000000000256"},
		{"larger than signed 64 bit integer", "0x4b3b4ca85a86c47a098a224000000000", "100000000000000000000.000000000000000000"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {

			rpcClient := new(mocks.RPCClient)
			ethClient := eth.NewClientWith(rpcClient, nil)

			rpcClient.On("Call", mock.Anything, "eth_getBalance", mock.Anything, "latest").
				Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0).(*string)
					*res = test.input
				})

			result, err := ethClient.GetEthBalance(cltest.NewAddress())
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result.String())
		})
	}
}

func TestEthClient_GetERC20Balance(t *testing.T) {
	t.Parallel()

	expectedBig, _ := big.NewInt(0).SetString("100000000000000000000000000000000000000", 10)

	tests := []struct {
		name     string
		input    string
		expected *big.Int
	}{
		{"small", "0x0100", big.NewInt(256)},
		{"big", "0x4b3b4ca85a86c47a098a224000000000", expectedBig},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {

			rpcClient := new(mocks.RPCClient)
			ethClient := eth.NewClientWith(rpcClient, nil)

			contractAddress := cltest.NewAddress()
			userAddress := cltest.NewAddress()

			functionSelector := models.HexToFunctionSelector("0x70a08231") // balanceOf(address)
			data := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(userAddress.Bytes(), utils.EVMWordByteLen))
			callArgs := eth.CallArgs{
				To:   contractAddress,
				Data: data,
			}

			rpcClient.On("Call", mock.Anything, "eth_call", callArgs, "latest").
				Return(nil).
				Run(func(args mock.Arguments) {
					res := args.Get(0).(*string)
					*res = test.input
				})

			result, err := ethClient.GetERC20Balance(userAddress, contractAddress)
			assert.NoError(t, err)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}
