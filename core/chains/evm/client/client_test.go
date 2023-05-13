package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func mustNewClient(t *testing.T, wsURL string, sendonlys ...url.URL) evmclient.Client {
	return mustNewClientWithChainID(t, wsURL, testutils.FixtureChainID, sendonlys...)
}

func mustNewClientWithChainID(t *testing.T, wsURL string, chainID *big.Int, sendonlys ...url.URL) evmclient.Client {
	cfg := evmclient.TestNodeConfig{
		SelectionMode: evmclient.NodeSelectionMode_RoundRobin,
	}
	c, err := evmclient.NewClientWithTestNode(t, cfg, wsURL, nil, sendonlys, 42, chainID)
	require.NoError(t, err)
	return c
}

func TestEthClient_TransactionReceipt(t *testing.T) {
	t.Parallel()

	txHash := "0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"

	mustReadResult := func(t *testing.T, file string) []byte {
		response := cltest.MustReadFile(t, file)
		var resp struct {
			Result json.RawMessage `json:"result"`
		}
		err := json.Unmarshal(response, &resp)
		require.NoError(t, err)
		return resp.Result
	}

	t.Run("happy path", func(t *testing.T) {
		result := mustReadResult(t, "../../../testdata/jsonrpc/getTransactionReceipt.json")

		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			}
			if assert.Equal(t, "eth_getTransactionReceipt", method) && assert.True(t, params.IsArray()) &&
				assert.Equal(t, txHash, params.Array()[0].String()) {
				resp.Result = string(result)
			}
			return
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		hash := common.HexToHash(txHash)
		receipt, err := ethClient.TransactionReceipt(testutils.Context(t), hash)
		require.NoError(t, err)
		assert.Equal(t, hash, receipt.TxHash)
		assert.Equal(t, big.NewInt(11), receipt.BlockNumber)
	})

	t.Run("no tx hash, returns ethereum.NotFound", func(t *testing.T) {
		result := mustReadResult(t, "../../../testdata/jsonrpc/getTransactionReceipt_notFound.json")
		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			}
			if assert.Equal(t, "eth_getTransactionReceipt", method) && assert.True(t, params.IsArray()) &&
				assert.Equal(t, txHash, params.Array()[0].String()) {
				resp.Result = string(result)
			}
			return
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		hash := common.HexToHash(txHash)
		_, err = ethClient.TransactionReceipt(testutils.Context(t), hash)
		require.Equal(t, ethereum.NotFound, errors.Cause(err))
	})
}

func TestEthClient_PendingNonceAt(t *testing.T) {
	t.Parallel()

	address := testutils.NewAddress()

	wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
		switch method {
		case "eth_subscribe":
			resp.Result = `"0x00"`
			resp.Notify = headResult
			return
		case "eth_unsubscribe":
			resp.Result = "true"
			return
		}
		if !assert.Equal(t, "eth_getTransactionCount", method) || !assert.True(t, params.IsArray()) {
			return
		}
		arr := params.Array()
		if assert.Equal(t, strings.ToLower(address.Hex()), strings.ToLower(arr[0].String())) &&
			assert.Equal(t, "pending", arr[1].String()) {
			resp.Result = `"0x100"`
		}
		return
	})

	ethClient := mustNewClient(t, wsURL)
	err := ethClient.Dial(testutils.Context(t))
	require.NoError(t, err)

	result, err := ethClient.PendingNonceAt(testutils.Context(t), address)
	require.NoError(t, err)

	var expected uint64 = 256
	require.Equal(t, result, expected)
}

func TestEthClient_BalanceAt(t *testing.T) {
	t.Parallel()

	largeBalance, _ := big.NewInt(0).SetString("100000000000000000000", 10)
	address := testutils.NewAddress()

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
			wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
					resp.Result = `"0x00"`
					resp.Notify = headResult
					return
				case "eth_unsubscribe":
					resp.Result = "true"
					return
				}
				if assert.Equal(t, "eth_getBalance", method) && assert.True(t, params.IsArray()) &&
					assert.Equal(t, strings.ToLower(address.Hex()), strings.ToLower(params.Array()[0].String())) {
					resp.Result = `"` + hexutil.EncodeBig(test.balance) + `"`
				}
				return
			})

			ethClient := mustNewClient(t, wsURL)
			err := ethClient.Dial(testutils.Context(t))
			require.NoError(t, err)

			result, err := ethClient.BalanceAt(testutils.Context(t), address, nil)
			require.NoError(t, err)
			assert.Equal(t, test.balance, result)
		})
	}
}

func TestEthClient_GetERC20Balance(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

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
			contractAddress := testutils.NewAddress()
			userAddress := testutils.NewAddress()
			functionSelector := evmtypes.HexToFunctionSelector("0x70a08231") // balanceOf(address)
			txData := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(userAddress.Bytes(), utils.EVMWordByteLen))

			wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
					resp.Result = `"0x00"`
					resp.Notify = headResult
					return
				case "eth_unsubscribe":
					resp.Result = "true"
					return
				}
				if !assert.Equal(t, "eth_call", method) || !assert.True(t, params.IsArray()) {
					return
				}
				arr := params.Array()
				callArgs := arr[0]
				if assert.True(t, callArgs.IsObject()) &&
					assert.Equal(t, strings.ToLower(contractAddress.Hex()), callArgs.Get("to").String()) &&
					assert.Equal(t, hexutil.Encode(txData), callArgs.Get("data").String()) &&
					assert.Equal(t, "latest", arr[1].String()) {

					resp.Result = `"` + hexutil.EncodeBig(test.balance) + `"`
				}
				return

			})

			ethClient := mustNewClient(t, wsURL)
			err := ethClient.Dial(testutils.Context(t))
			require.NoError(t, err)

			result, err := ethClient.TokenBalance(ctx, userAddress, contractAddress)
			require.NoError(t, err)
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
	t.Parallel()

	expectedBlockNum := big.NewInt(1)
	expectedBlockHash := "0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a"

	tests := []struct {
		name                  string
		expectedRequestBlock  *big.Int
		expectedResponseBlock int64
		error                 error
		rpcResp               string
	}{
		{"happy geth", expectedBlockNum, expectedBlockNum.Int64(), nil,
			`{"difficulty":"0xf3a00","extraData":"0xd883010503846765746887676f312e372e318664617277696e","gasLimit":"0xffc001","gasUsed":"0x0","hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd1aeb42885a43b72b518182ef893125814811048","mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","nonce":"0x0ece08ea8c49dfd9","number":"0x1","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x218","stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b","timestamp":"0x58318da2","totalDifficulty":"0x1f3a00","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}`},
		{"happy parity", expectedBlockNum, expectedBlockNum.Int64(), nil,
			`{"author":"0xd1aeb42885a43b72b518182ef893125814811048","difficulty":"0xf3a00","extraData":"0xd883010503846765746887676f312e372e318664617277696e","gasLimit":"0xffc001","gasUsed":"0x0","hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd1aeb42885a43b72b518182ef893125814811048","mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","nonce":"0x0ece08ea8c49dfd9","number":"0x1","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sealFields":["0xa00f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","0x880ece08ea8c49dfd9"],"sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x218","stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b","timestamp":"0x58318da2","totalDifficulty":"0x1f3a00","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}`},
		{"missing header", expectedBlockNum, 0, fmt.Errorf("no live nodes available for chain %s", cltest.FixtureChainID.String()),
			`null`},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
				switch method {
				case "eth_subscribe":
					resp.Result = `"0x00"`
					resp.Notify = headResult
					return
				case "eth_unsubscribe":
					resp.Result = "true"
					return
				}
				if !assert.Equal(t, "eth_getBlockByNumber", method) || !assert.True(t, params.IsArray()) {
					return
				}
				arr := params.Array()
				blockNumStr := arr[0].String()
				var blockNum hexutil.Big
				err := blockNum.UnmarshalText([]byte(blockNumStr))
				if assert.NoError(t, err) && assert.Equal(t, test.expectedRequestBlock, blockNum.ToInt()) &&
					assert.Equal(t, false, arr[1].Bool()) {
					resp.Result = test.rpcResp
				}
				return
			})

			ethClient := mustNewClient(t, wsURL)
			err := ethClient.Dial(testutils.Context(t))
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(testutils.Context(t), 5*time.Second)
			defer cancel()
			result, err := ethClient.HeadByNumber(ctx, expectedBlockNum)
			if test.error != nil {
				require.Error(t, err, test.error)
			} else {
				require.NoError(t, err)
				require.Equal(t, expectedBlockHash, result.Hash.Hex())
				require.Equal(t, test.expectedResponseBlock, result.Number)
				require.Zero(t, cltest.FixtureChainID.Cmp(result.EVMChainID.ToInt()))
			}
		})
	}
}

func TestEthClient_SendTransaction_NoSecondaryURL(t *testing.T) {
	t.Parallel()

	tx := types.NewTransaction(uint64(42), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})

	wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
		switch method {
		case "eth_subscribe":
			resp.Result = `"0x00"`
			resp.Notify = headResult
			return
		case "eth_unsubscribe":
			resp.Result = "true"
			return
		}
		if !assert.Equal(t, "eth_sendRawTransaction", method) {
			return
		}
		resp.Result = `"` + tx.Hash().Hex() + `"`
		return
	})

	ethClient := mustNewClient(t, wsURL)
	err := ethClient.Dial(testutils.Context(t))
	require.NoError(t, err)

	err = ethClient.SendTransaction(testutils.Context(t), tx)
	assert.NoError(t, err)
}

func TestEthClient_SendTransaction_WithSecondaryURLs(t *testing.T) {
	t.Parallel()

	tx := types.NewTransaction(uint64(42), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})

	wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
		switch method {
		case "eth_subscribe":
			resp.Result = `"0x00"`
			resp.Notify = headResult
			return
		case "eth_unsubscribe":
			resp.Result = "true"
			return
		case "eth_sendRawTransaction":
			resp.Result = `"` + tx.Hash().Hex() + `"`
		}
		return
	})

	rpcSrv := rpc.NewServer()
	t.Cleanup(rpcSrv.Stop)
	service := sendTxService{chainID: &cltest.FixtureChainID}
	rpcSrv.RegisterName("eth", &service)
	ts := httptest.NewServer(rpcSrv)
	t.Cleanup(ts.Close)

	sendonlyURL := *cltest.MustParseURL(t, ts.URL)
	ethClient := mustNewClient(t, wsURL, sendonlyURL, sendonlyURL)
	err := ethClient.Dial(testutils.Context(t))
	require.NoError(t, err)

	err = ethClient.SendTransaction(testutils.Context(t), tx)
	require.NoError(t, err)

	// Unfortunately it's a bit tricky to test this, since there is no
	// synchronization. We have to rely on timing instead.
	require.Eventually(t, func() bool { return service.sentCount.Load() == int32(2) }, testutils.WaitTimeout(t), 500*time.Millisecond)
}

func TestEthClient_SendTransactionReturnCode(t *testing.T) {
	t.Parallel()

	fromAddress := testutils.NewAddress()
	tx := types.NewTransaction(uint64(42), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})

	t.Run("returns Fatal error type when error message is fatal", func(t *testing.T) {
		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_sendRawTransaction":
				resp.Result = `"` + tx.Hash().Hex() + `"`
				resp.Error.Message = "invalid sender"
			}
			return
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		errType, err := ethClient.SendTransactionReturnCode(testutils.Context(t), tx, fromAddress)
		assert.Error(t, err)
		assert.Equal(t, errType, clienttypes.Fatal)
	})

	t.Run("returns TransactionAlreadyKnown error type when error message is nonce too low", func(t *testing.T) {
		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_sendRawTransaction":
				resp.Result = `"` + tx.Hash().Hex() + `"`
				resp.Error.Message = "nonce too low"
			}
			return
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		errType, err := ethClient.SendTransactionReturnCode(testutils.Context(t), tx, fromAddress)
		assert.Error(t, err)
		assert.Equal(t, errType, clienttypes.TransactionAlreadyKnown)
	})

	t.Run("returns Successful error type when there is no error message", func(t *testing.T) {
		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_sendRawTransaction":
				resp.Result = `"` + tx.Hash().Hex() + `"`
			}
			return
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		errType, err := ethClient.SendTransactionReturnCode(testutils.Context(t), tx, fromAddress)
		assert.NoError(t, err)
		assert.Equal(t, errType, clienttypes.Successful)
	})

	t.Run("returns Underpriced error type when transaction is terminally underpriced", func(t *testing.T) {
		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_sendRawTransaction":
				resp.Result = `"` + tx.Hash().Hex() + `"`
				resp.Error.Message = "transaction underpriced"
			}
			return
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		errType, err := ethClient.SendTransactionReturnCode(testutils.Context(t), tx, fromAddress)
		assert.Error(t, err)
		assert.Equal(t, errType, clienttypes.Underpriced)
	})

	t.Run("returns Unsupported error type when error message is queue full", func(t *testing.T) {
		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_sendRawTransaction":
				resp.Result = `"` + tx.Hash().Hex() + `"`
				resp.Error.Message = "queue full"
			}
			return
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		errType, err := ethClient.SendTransactionReturnCode(testutils.Context(t), tx, fromAddress)
		assert.Error(t, err)
		assert.Equal(t, errType, clienttypes.Unsupported)
	})

	t.Run("returns Retryable error type when there is a transaction gap", func(t *testing.T) {
		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_sendRawTransaction":
				resp.Result = `"` + tx.Hash().Hex() + `"`
				resp.Error.Message = "NonceGap"
			}
			return
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		errType, err := ethClient.SendTransactionReturnCode(testutils.Context(t), tx, fromAddress)
		assert.Error(t, err)
		assert.Equal(t, errType, clienttypes.Retryable)
	})

	t.Run("returns InsufficientFunds error type when the sender address doesn't have enough funds", func(t *testing.T) {
		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_sendRawTransaction":
				resp.Result = `"` + tx.Hash().Hex() + `"`
				resp.Error.Message = "insufficient funds for transfer"
			}
			return
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		errType, err := ethClient.SendTransactionReturnCode(testutils.Context(t), tx, fromAddress)
		assert.Error(t, err)
		assert.Equal(t, errType, clienttypes.InsufficientFunds)
	})

	t.Run("returns ExceedsFeeCap error type when gas price is too high for the node", func(t *testing.T) {
		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_sendRawTransaction":
				resp.Result = `"` + tx.Hash().Hex() + `"`
				resp.Error.Message = "Transaction fee cap exceeded"
			}
			return
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		errType, err := ethClient.SendTransactionReturnCode(testutils.Context(t), tx, fromAddress)
		assert.Error(t, err)
		assert.Equal(t, errType, clienttypes.ExceedsMaxFee)
	})

	t.Run("returns Unknown error type when the error can't be categorized", func(t *testing.T) {
		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
			switch method {
			case "eth_subscribe":
				resp.Result = `"0x00"`
				resp.Notify = headResult
				return
			case "eth_unsubscribe":
				resp.Result = "true"
				return
			case "eth_sendRawTransaction":
				resp.Result = `"` + tx.Hash().Hex() + `"`
				resp.Error.Message = "some random error"
			}
			return
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		errType, err := ethClient.SendTransactionReturnCode(testutils.Context(t), tx, fromAddress)
		assert.Error(t, err)
		assert.Equal(t, errType, clienttypes.Unknown)
	})
}

type sendTxService struct {
	chainID   *big.Int
	sentCount atomic.Int32
}

func (x *sendTxService) ChainId(ctx context.Context) (*hexutil.Big, error) {
	return (*hexutil.Big)(x.chainID), nil
}

func (x *sendTxService) SendRawTransaction(ctx context.Context, signRawTx hexutil.Bytes) error {
	x.sentCount.Add(1)
	return nil
}

func TestEthClient_SubscribeNewHead(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testutils.Context(t), testutils.WaitTimeout(t))
	defer cancel()

	chainId := big.NewInt(123456)
	wsURL := cltest.NewWSServer(t, chainId, func(method string, params gjson.Result) (resp testutils.JSONRPCResponse) {
		if method == "eth_unsubscribe" {
			resp.Result = "true"
			return
		}
		assert.Equal(t, "eth_subscribe", method)
		if assert.True(t, params.IsArray()) && assert.Equal(t, "newHeads", params.Array()[0].String()) {
			resp.Result = `"0x00"`
			resp.Notify = headResult
		}
		return
	})

	ethClient := mustNewClientWithChainID(t, wsURL, chainId)
	err := ethClient.Dial(testutils.Context(t))
	require.NoError(t, err)

	headCh := make(chan *evmtypes.Head)
	sub, err := ethClient.SubscribeNewHead(ctx, headCh)
	require.NoError(t, err)
	defer sub.Unsubscribe()

	select {
	case err := <-sub.Err():
		t.Fatal(err)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	case h := <-headCh:
		require.NotNil(t, h.EVMChainID)
		require.Zero(t, chainId.Cmp(h.EVMChainID.ToInt()))
	}
}

const headResult = evmclient.HeadResult
