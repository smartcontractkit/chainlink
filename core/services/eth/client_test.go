package eth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEthClient_TransactionReceipt(t *testing.T) {
	txHash := "0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"

	t.Run("happy path", func(t *testing.T) {
		response := cltest.MustReadFile(t, "../../testdata/jsonrpc/getTransactionReceipt.json")
		_, wsUrl := cltest.NewWSServer(t, string(response), func(data []byte) {
			resp := cltest.ParseJSON(t, bytes.NewReader(data))
			require.Equal(t, "eth_getTransactionReceipt", resp.Get("method").String())
			require.True(t, resp.Get("params").IsArray())
			require.Equal(t, txHash, resp.Get("params").Get("0").String())
		})

		ethClient, err := eth.NewClient(logger.TestLogger(t), wsUrl, nil, []url.URL{}, nil)
		require.NoError(t, err)
		err = ethClient.Dial(context.Background())
		require.NoError(t, err)

		hash := common.HexToHash(txHash)
		receipt, err := ethClient.TransactionReceipt(context.Background(), hash)
		assert.NoError(t, err)
		assert.Equal(t, hash, receipt.TxHash)
		assert.Equal(t, big.NewInt(11), receipt.BlockNumber)
	})

	t.Run("no tx hash, returns ethereum.NotFound", func(t *testing.T) {
		response := cltest.MustReadFile(t, "../../testdata/jsonrpc/getTransactionReceipt_notFound.json")
		_, wsUrl := cltest.NewWSServer(t, string(response), func(data []byte) {
			resp := cltest.ParseJSON(t, bytes.NewReader(data))
			require.Equal(t, "eth_getTransactionReceipt", resp.Get("method").String())
			require.True(t, resp.Get("params").IsArray())
			require.Equal(t, txHash, resp.Get("params").Get("0").String())
		})

		ethClient, err := eth.NewClient(logger.TestLogger(t), wsUrl, nil, nil, nil)
		require.NoError(t, err)
		err = ethClient.Dial(context.Background())
		require.NoError(t, err)

		hash := common.HexToHash(txHash)
		_, err = ethClient.TransactionReceipt(context.Background(), hash)
		require.Equal(t, ethereum.NotFound, errors.Cause(err))
	})
}

func TestEthClient_PendingNonceAt(t *testing.T) {
	t.Parallel()

	address := cltest.NewAddress()

	_, url := cltest.NewWSServer(t, `{
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

	ethClient, err := eth.NewClient(logger.TestLogger(t), url, nil, nil, nil)
	require.NoError(t, err)
	err = ethClient.Dial(context.Background())
	require.NoError(t, err)

	result, err := ethClient.PendingNonceAt(context.Background(), address)
	require.NoError(t, err)

	var expected uint64 = 256
	require.Equal(t, result, expected)
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
			_, url := cltest.NewWSServer(t, `{
              "id": 1,
              "jsonrpc": "2.0",
              "result": "`+hexutil.EncodeBig(test.balance)+`"
            }`, func(data []byte) {
				resp := cltest.ParseJSON(t, bytes.NewReader(data))
				require.Equal(t, "eth_getBalance", resp.Get("method").String())
				require.True(t, resp.Get("params").IsArray())
				require.Equal(t, strings.ToLower(address.Hex()), strings.ToLower(resp.Get("params").Get("0").String()))
			})

			ethClient, err := eth.NewClient(logger.TestLogger(t), url, nil, nil, nil)
			require.NoError(t, err)
			err = ethClient.Dial(context.Background())
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
			functionSelector := eth.HexToFunctionSelector("0x70a08231") // balanceOf(address)
			txData := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(userAddress.Bytes(), utils.EVMWordByteLen))

			_, url := cltest.NewWSServer(t, `{
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

			ethClient, err := eth.NewClient(logger.TestLogger(t), url, nil, nil, nil)
			require.NoError(t, err)
			err = ethClient.Dial(context.Background())
			require.NoError(t, err)

			result, err := ethClient.GetERC20Balance(userAddress, contractAddress)
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

func TestEthClient_HeaderByNumber(t *testing.T) {
	expectedBlockNum := big.NewInt(1)
	expectedBlockHash := "0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a"

	tests := []struct {
		name                  string
		expectedRequestBlock  *big.Int
		expectedResponseBlock int64
		error                 error
		rpcResp               string
	}{
		{"happy geth", expectedBlockNum, expectedBlockNum.Int64(), nil, `{"jsonrpc":"2.0","id":1,"result":{"difficulty":"0xf3a00","extraData":"0xd883010503846765746887676f312e372e318664617277696e","gasLimit":"0xffc001","gasUsed":"0x0","hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd1aeb42885a43b72b518182ef893125814811048","mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","nonce":"0x0ece08ea8c49dfd9","number":"0x1","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x218","stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b","timestamp":"0x58318da2","totalDifficulty":"0x1f3a00","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}}`},
		{"happy parity", expectedBlockNum, expectedBlockNum.Int64(), nil, `{"jsonrpc":"2.0","result":{"author":"0xd1aeb42885a43b72b518182ef893125814811048","difficulty":"0xf3a00","extraData":"0xd883010503846765746887676f312e372e318664617277696e","gasLimit":"0xffc001","gasUsed":"0x0","hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd1aeb42885a43b72b518182ef893125814811048","mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","nonce":"0x0ece08ea8c49dfd9","number":"0x1","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sealFields":["0xa00f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","0x880ece08ea8c49dfd9"],"sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x218","stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b","timestamp":"0x58318da2","totalDifficulty":"0x1f3a00","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]},"id":1}`},
		{"missing header", expectedBlockNum, 0, ethereum.NotFound, `{"jsonrpc":"2.0","id":1,"result":null}`},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			_, url := cltest.NewWSServer(t, test.rpcResp, func(data []byte) {
				req := cltest.ParseJSON(t, bytes.NewReader(data))

				require.True(t, req.IsObject())

				require.Equal(t, "eth_getBlockByNumber", req.Get("method").String())
				require.True(t, req.Get("params").IsArray())

				blockNumStr := req.Get("params").Get("0").String()
				var blockNum hexutil.Big
				err := blockNum.UnmarshalText([]byte(blockNumStr))
				require.NoError(t, err)
				require.Equal(t, test.expectedRequestBlock, blockNum.ToInt())

				require.Equal(t, false, req.Get("params").Get("1").Bool())
			})

			ethClient, err := eth.NewClient(logger.TestLogger(t), url, nil, nil, nil)
			require.NoError(t, err)
			err = ethClient.Dial(context.Background())
			require.NoError(t, err)
			defer ethClient.Close()

			result, err := ethClient.HeadByNumber(context.Background(), expectedBlockNum)
			if test.error != nil {
				require.Equal(t, test.error, errors.Cause(err))
			} else {
				require.NoError(t, err)
				require.Equal(t, expectedBlockHash, result.Hash.Hex())
				require.Equal(t, test.expectedResponseBlock, result.Number)
			}
		})
	}
}

func TestEthClient_SendTransaction_NoSecondaryURL(t *testing.T) {
	t.Parallel()

	tx := types.NewTransaction(uint64(42), cltest.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})

	_, url := cltest.NewWSServer(t, `{
  "id": 1,
  "jsonrpc": "2.0",
  "result": "`+tx.Hash().Hex()+`"
}`, func(data []byte) {
		resp := cltest.ParseJSON(t, bytes.NewReader(data))
		require.Equal(t, "eth_sendRawTransaction", resp.Get("method").String())
		require.True(t, resp.Get("params").IsArray())
	})

	ethClient, err := eth.NewClient(logger.TestLogger(t), url, nil, nil, nil)
	require.NoError(t, err)
	err = ethClient.Dial(context.Background())
	require.NoError(t, err)

	err = ethClient.SendTransaction(context.Background(), tx)
	assert.NoError(t, err)
}

func TestEthClient_SendTransaction_WithSecondaryURLs(t *testing.T) {
	t.Parallel()

	tx := types.NewTransaction(uint64(42), cltest.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})

	response := `{
  "id": 1,
  "jsonrpc": "2.0",
  "result": "` + tx.Hash().Hex() + `"
}`

	_, wsUrl := cltest.NewWSServer(t, response, func(data []byte) {
		resp := cltest.ParseJSON(t, bytes.NewReader(data))
		require.Equal(t, "eth_sendRawTransaction", resp.Get("method").String())
		require.True(t, resp.Get("params").IsArray())
	})

	requests := make(chan struct{}, 2)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(response))
		require.NoError(t, err)
		requests <- struct{}{}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	sendonlyUrl := *cltest.MustParseURL(t, server.URL)
	ethClient, err := eth.NewClient(logger.TestLogger(t), wsUrl, nil, []url.URL{sendonlyUrl, sendonlyUrl}, nil)
	require.NoError(t, err)
	err = ethClient.Dial(context.Background())
	require.NoError(t, err)

	err = ethClient.SendTransaction(context.Background(), tx)
	assert.NoError(t, err)

	cltest.NewGomegaWithT(t).Eventually(func() int {
		return len(requests)
	}).Should(gomega.Equal(2))
}
