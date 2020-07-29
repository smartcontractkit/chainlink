package eth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"math/big"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEthClient_TransactionReceipt(t *testing.T) {
	txHash := "0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"

	t.Run("happy path", func(t *testing.T) {
		response := cltest.MustReadFile(t, "testdata/getTransactionReceipt.json")
		_, wsUrl, wsCleanup := cltest.NewWSServer(string(response), func(data []byte) {
			resp := cltest.ParseJSON(t, bytes.NewReader(data))
			require.Equal(t, "eth_getTransactionReceipt", resp.Get("method").String())
			require.True(t, resp.Get("params").IsArray())
			require.Equal(t, txHash, resp.Get("params").Get("0").String())
		})
		defer wsCleanup()

		ethClient := eth.NewClient(wsUrl)
		err := ethClient.Dial(context.Background())
		require.NoError(t, err)

		hash := common.HexToHash(txHash)
		receipt, err := ethClient.TransactionReceipt(context.Background(), hash)
		assert.NoError(t, err)
		assert.Equal(t, hash, receipt.TxHash)
		assert.Equal(t, big.NewInt(11), receipt.BlockNumber)
	})

	t.Run("no tx hash, returns ethereum.NotFound", func(t *testing.T) {
		response := cltest.MustReadFile(t, "testdata/getTransactionReceipt_notFound.json")
		_, wsUrl, wsCleanup := cltest.NewWSServer(string(response), func(data []byte) {
			resp := cltest.ParseJSON(t, bytes.NewReader(data))
			require.Equal(t, "eth_getTransactionReceipt", resp.Get("method").String())
			require.True(t, resp.Get("params").IsArray())
			require.Equal(t, txHash, resp.Get("params").Get("0").String())
		})
		defer wsCleanup()

		ethClient := eth.NewClient(wsUrl)
		err := ethClient.Dial(context.Background())
		require.NoError(t, err)

		hash := common.HexToHash(txHash)
		_, err = ethClient.TransactionReceipt(context.Background(), hash)
		require.Equal(t, ethereum.NotFound, err)
	})
}

func TestEthClient_PendingNonceAt(t *testing.T) {
	t.Parallel()

	address := cltest.NewAddress()

	_, url, cleanup := cltest.NewWSServer(`{
      "id": 1,
      "jsonrpc": "2.0",
      "result": "0x100"
    }`, func(data []byte) {
		resp := cltest.ParseJSON(t, bytes.NewReader(data))
		require.Equal(t, "eth_getTransactionCount", resp.Get("method").String())
		require.True(t, resp.Get("params").IsArray())
		require.Equal(t, strings.ToLower(address.Hex()), strings.ToLower(resp.Get("params").Get("0").String()))
		require.Equal(t, "pending", resp.Get("params").Get("1").String())
	})
	defer cleanup()

	ethClient := eth.NewClient(url)
	err := ethClient.Dial(context.Background())
	require.NoError(t, err)

	result, err := ethClient.PendingNonceAt(context.Background(), address)
	require.NoError(t, err)

	var expected uint64 = 256
	require.Equal(t, result, expected)
}

func TestEthClient_SendRawTx(t *testing.T) {
	t.Parallel()

	txData := "0xdeadbeef"

	returnedHash := cltest.NewHash()
	_, url, cleanup := cltest.NewWSServer(`{
      "id": 1,
      "jsonrpc": "2.0",
      "result": "`+returnedHash.Hex()+`"
    }`, func(data []byte) {
		resp := cltest.ParseJSON(t, bytes.NewReader(data))
		require.Equal(t, "eth_sendRawTransaction", resp.Get("method").String())
		require.True(t, resp.Get("params").IsArray())
		require.Equal(t, txData, resp.Get("params").Get("0").String())
	})
	defer cleanup()

	ethClient := eth.NewClient(url)
	err := ethClient.Dial(context.Background())
	require.NoError(t, err)

	result, err := ethClient.SendRawTx(hexutil.MustDecode(txData))
	assert.NoError(t, err)
	assert.Equal(t, result, returnedHash)
}

func TestEthClient_BalanceAt(t *testing.T) {
	t.Parallel()

	largeBalance, _ := big.NewInt(0).SetString("100000000000000000000", 10)
	address := cltest.NewAddress()

	tests := []struct {
		name    string
		balance *big.Int
	}{
		{"basic", big.NewInt(256)},
		{"larger than signed 64 bit integer", largeBalance},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			_, url, cleanup := cltest.NewWSServer(`{
              "id": 1,
              "jsonrpc": "2.0",
              "result": "`+hexutil.EncodeBig(test.balance)+`"
            }`, func(data []byte) {
				resp := cltest.ParseJSON(t, bytes.NewReader(data))
				require.Equal(t, "eth_getBalance", resp.Get("method").String())
				require.True(t, resp.Get("params").IsArray())
				require.Equal(t, strings.ToLower(address.Hex()), strings.ToLower(resp.Get("params").Get("0").String()))
			})
			defer cleanup()

			ethClient := eth.NewClient(url)
			err := ethClient.Dial(context.Background())
			require.NoError(t, err)

			result, err := ethClient.BalanceAt(context.Background(), address, nil)
			assert.NoError(t, err)
			assert.Equal(t, test.balance, result)
		})
	}
}

func TestEthClient_GetERC20Balance(t *testing.T) {
	t.Parallel()

	expectedBig, _ := big.NewInt(0).SetString("100000000000000000000000000000000000000", 10)

	tests := []struct {
		name    string
		balance *big.Int
	}{
		{"small", big.NewInt(256)},
		{"big", expectedBig},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			contractAddress := cltest.NewAddress()
			userAddress := cltest.NewAddress()
			functionSelector := models.HexToFunctionSelector("0x70a08231") // balanceOf(address)
			txData := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(userAddress.Bytes(), utils.EVMWordByteLen))

			_, url, cleanup := cltest.NewWSServer(`{
              "id": 1,
              "jsonrpc": "2.0",
              "result": "`+hexutil.EncodeBig(test.balance)+`"
            }`, func(data []byte) {
				resp := cltest.ParseJSON(t, bytes.NewReader(data))
				require.Equal(t, "eth_call", resp.Get("method").String())
				require.True(t, resp.Get("params").IsArray())

				callArgs := resp.Get("params").Get("0")
				require.True(t, callArgs.IsObject())
				require.Equal(t, strings.ToLower(contractAddress.Hex()), callArgs.Get("to").String())
				require.Equal(t, hexutil.Encode(txData), callArgs.Get("data").String())

				require.Equal(t, "latest", resp.Get("params").Get("1").String())
			})
			defer cleanup()

			ethClient := eth.NewClient(url)
			err := ethClient.Dial(context.Background())
			require.NoError(t, err)

			result, err := ethClient.GetERC20Balance(userAddress, contractAddress)
			assert.NoError(t, err)
			assert.NoError(t, err)
			assert.Equal(t, test.balance, result)
		})
	}
}

func TestReceipt_UnmarshalEmptyBlockHash(t *testing.T) {
	t.Parallel()

	input := `{
        "transactionHash": "0x444172bef57ad978655171a8af2cfd89baa02a97fcb773067aef7794d6913374",
        "gasUsed": "0x1",
        "cumulativeGasUsed": "0x1",
        "logs": [],
        "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "blockNumber": "0x8bf99b",
        "blockHash": null
    }`

	var receipt types.Receipt
	err := json.Unmarshal([]byte(input), &receipt)
	require.NoError(t, err)
}
