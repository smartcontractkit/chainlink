package client_test

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func mustNewChainClient(t *testing.T, wsURL string, sendonlys ...url.URL) *evmclient.ChainClient {
	return mustNewChainClientWithChainID(t, wsURL, testutils.FixtureChainID, sendonlys...)
}

func mustNewChainClientWithChainID(t *testing.T, wsURL string, chainID *big.Int, sendonlys ...url.URL) *evmclient.ChainClient {
	cfg := evmclient.TestNodePoolConfig{
		NodeSelectionMode: evmclient.NodeSelectionMode_RoundRobin,
	}
	c, err := evmclient.NewChainClientWithTestNode(t, cfg, time.Second*0, wsURL, nil, sendonlys, 42, chainID)
	require.NoError(t, err)
	return c
}

func TestEthChainClient_TransactionReceipt(t *testing.T) {
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

		chainClient := mustNewChainClient(t, wsURL)
		err := chainClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		hash := common.HexToHash(txHash)
		receipt, err := chainClient.TransactionReceipt(testutils.Context(t), hash)
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

		chainClient := mustNewChainClient(t, wsURL)
		err := chainClient.Dial(testutils.Context(t))
		require.NoError(t, err)

		hash := common.HexToHash(txHash)
		_, err = chainClient.TransactionReceipt(testutils.Context(t), hash)
		require.Equal(t, ethereum.NotFound, errors.Cause(err))
	})
}

func TestEthChainClient_PendingSequenceAt(t *testing.T) {
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

	chainClient := mustNewChainClient(t, wsURL)
	err := chainClient.Dial(testutils.Context(t))
	require.NoError(t, err)

	result, err := chainClient.PendingSequenceAt(testutils.Context(t), address)
	require.NoError(t, err)

	var expected int64 = 256
	require.Equal(t, result.Int64(), expected)
}

func TestEthChainClient_BalanceAt(t *testing.T) {
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

			chainClient := mustNewChainClient(t, wsURL)
			err := chainClient.Dial(testutils.Context(t))
			require.NoError(t, err)

			result, err := chainClient.BalanceAt(testutils.Context(t), address, nil)
			require.NoError(t, err)
			assert.Equal(t, test.balance, result)
		})
	}
}

func TestEthChainClient_LatestBlockHeight(t *testing.T) {
	t.Parallel()

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
		if !assert.Equal(t, "eth_blockNumber", method) {
			return
		}
		resp.Result = `"0x100"`
		return
	})

	chainClient := mustNewChainClient(t, wsURL)
	err := chainClient.Dial(testutils.Context(t))
	require.NoError(t, err)

	result, err := chainClient.LatestBlockHeight(testutils.Context(t))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(256), result)
}

func TestEthChainClient_GetERC20Balance(t *testing.T) {
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
			functionSelector := evmtypes.HexToFunctionSelector(evmclient.BALANCE_OF_ADDRESS_FUNCTION_SELECTOR) // balanceOf(address)
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

			chainClient := mustNewChainClient(t, wsURL)
			err := chainClient.Dial(testutils.Context(t))
			require.NoError(t, err)

			result, err := chainClient.TokenBalance(ctx, userAddress, contractAddress)
			require.NoError(t, err)
			assert.Equal(t, test.balance, result)
		})
	}
}

func TestEthChainClient_SendTransaction_NoSecondaryURL(t *testing.T) {
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

	chainClient := mustNewChainClient(t, wsURL)
	err := chainClient.Dial(testutils.Context(t))
	require.NoError(t, err)

	err = chainClient.SendTransaction(testutils.Context(t), tx)
	assert.NoError(t, err)
}

func TestEthChainClient_SendTransaction_WithSecondaryURLs(t *testing.T) {
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
	err := rpcSrv.RegisterName("eth", &service)
	require.NoError(t, err)
	ts := httptest.NewServer(rpcSrv)
	t.Cleanup(ts.Close)

	sendonlyURL := *cltest.MustParseURL(t, ts.URL)
	chainClient := mustNewChainClient(t, wsURL, sendonlyURL, sendonlyURL)
	err = chainClient.Dial(testutils.Context(t))
	require.NoError(t, err)

	err = chainClient.SendTransaction(testutils.Context(t), tx)
	require.NoError(t, err)

	// Unfortunately it's a bit tricky to test this, since there is no
	// synchronization. We have to rely on timing instead.
	require.Eventually(t, func() bool { return service.sentCount.Load() == int32(2) }, testutils.WaitTimeout(t), 500*time.Millisecond)
}

func TestEthChainClient_SendTransactionReturnCode(t *testing.T) {
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

func TestEthChainClient_SubscribeNewHead(t *testing.T) {
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
