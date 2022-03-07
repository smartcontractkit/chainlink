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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.uber.org/atomic"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func mustNewClient(t *testing.T, wsURL string, sendonlys ...url.URL) evmclient.Client {
	return mustNewClientWithChainID(t, wsURL, testutils.FixtureChainID, sendonlys...)
}

func mustNewClientWithChainID(t *testing.T, wsURL string, chainID *big.Int, sendonlys ...url.URL) evmclient.Client {
	cfg := evmclient.TestNodeConfig{}
	c, err := evmclient.NewClientWithTestNode(cfg, logger.TestLogger(t), wsURL, nil, sendonlys, 42, chainID)
	require.NoError(t, err)
	return c
}

func TestEthClient_TransactionReceipt(t *testing.T) {
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

		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (string, string) {
			require.Equal(t, "eth_getTransactionReceipt", method)
			require.True(t, params.IsArray())
			require.Equal(t, txHash, params.Array()[0].String())
			return string(result), ""
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(context.Background())
		require.NoError(t, err)

		hash := common.HexToHash(txHash)
		receipt, err := ethClient.TransactionReceipt(context.Background(), hash)
		require.NoError(t, err)
		assert.Equal(t, hash, receipt.TxHash)
		assert.Equal(t, big.NewInt(11), receipt.BlockNumber)
	})

	t.Run("no tx hash, returns ethereum.NotFound", func(t *testing.T) {
		result := mustReadResult(t, "../../../testdata/jsonrpc/getTransactionReceipt_notFound.json")
		wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (string, string) {
			require.Equal(t, "eth_getTransactionReceipt", method)
			require.True(t, params.IsArray())
			require.Equal(t, txHash, params.Array()[0].String())
			return string(result), ""
		})

		ethClient := mustNewClient(t, wsURL)
		err := ethClient.Dial(context.Background())
		require.NoError(t, err)

		hash := common.HexToHash(txHash)
		_, err = ethClient.TransactionReceipt(context.Background(), hash)
		require.Equal(t, ethereum.NotFound, errors.Cause(err))
	})
}

func TestEthClient_PendingNonceAt(t *testing.T) {
	t.Parallel()

	address := testutils.NewAddress()

	wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (string, string) {
		require.Equal(t, "eth_getTransactionCount", method)
		require.True(t, params.IsArray())
		arr := params.Array()
		require.Equal(t, strings.ToLower(address.Hex()), strings.ToLower(arr[0].String()))
		require.Equal(t, "pending", arr[1].String())
		return `"0x100"`, ""
	})

	ethClient := mustNewClient(t, wsURL)
	err := ethClient.Dial(context.Background())
	require.NoError(t, err)

	result, err := ethClient.PendingNonceAt(context.Background(), address)
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
			wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (string, string) {
				require.Equal(t, "eth_getBalance", method)
				require.True(t, params.IsArray())
				require.Equal(t, strings.ToLower(address.Hex()), strings.ToLower(params.Array()[0].String()))
				return `"` + hexutil.EncodeBig(test.balance) + `"`, ""
			})

			ethClient := mustNewClient(t, wsURL)
			err := ethClient.Dial(context.Background())
			require.NoError(t, err)

			result, err := ethClient.BalanceAt(context.Background(), address, nil)
			require.NoError(t, err)
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
			contractAddress := testutils.NewAddress()
			userAddress := testutils.NewAddress()
			functionSelector := evmtypes.HexToFunctionSelector("0x70a08231") // balanceOf(address)
			txData := utils.ConcatBytes(functionSelector.Bytes(), common.LeftPadBytes(userAddress.Bytes(), utils.EVMWordByteLen))

			wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (string, string) {
				require.Equal(t, "eth_call", method)
				require.True(t, params.IsArray())
				arr := params.Array()
				callArgs := arr[0]
				require.True(t, callArgs.IsObject())
				require.Equal(t, strings.ToLower(contractAddress.Hex()), callArgs.Get("to").String())
				require.Equal(t, hexutil.Encode(txData), callArgs.Get("data").String())

				require.Equal(t, "latest", arr[1].String())
				return `"` + hexutil.EncodeBig(test.balance) + `"`, ""
			})

			ethClient := mustNewClient(t, wsURL)
			err := ethClient.Dial(context.Background())
			require.NoError(t, err)

			result, err := ethClient.GetERC20Balance(userAddress, contractAddress)
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
		{"missing header", expectedBlockNum, 0, ethereum.NotFound,
			`null`},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (string, string) {
				require.Equal(t, "eth_getBlockByNumber", method)
				require.True(t, params.IsArray())
				arr := params.Array()
				blockNumStr := arr[0].String()
				var blockNum hexutil.Big
				err := blockNum.UnmarshalText([]byte(blockNumStr))
				require.NoError(t, err)
				require.Equal(t, test.expectedRequestBlock, blockNum.ToInt())

				require.Equal(t, false, arr[1].Bool())
				return test.rpcResp, ""
			})

			ethClient := mustNewClient(t, wsURL)
			err := ethClient.Dial(context.Background())
			require.NoError(t, err)
			defer ethClient.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			result, err := ethClient.HeadByNumber(ctx, expectedBlockNum)
			if test.error != nil {
				require.Equal(t, test.error, errors.Cause(err))
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

	wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (string, string) {
		require.Equal(t, "eth_sendRawTransaction", method)
		return `"` + tx.Hash().Hex() + `"`, ""
	})

	ethClient := mustNewClient(t, wsURL)
	err := ethClient.Dial(context.Background())
	require.NoError(t, err)

	err = ethClient.SendTransaction(context.Background(), tx)
	assert.NoError(t, err)
}

func TestEthClient_SendTransaction_WithSecondaryURLs(t *testing.T) {
	t.Parallel()

	tx := types.NewTransaction(uint64(42), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})

	wsURL := cltest.NewWSServer(t, &cltest.FixtureChainID, func(method string, params gjson.Result) (string, string) {
		require.Equal(t, "eth_sendRawTransaction", method)
		return `"` + tx.Hash().Hex() + `"`, ""
	})

	rpcSrv := rpc.NewServer()
	t.Cleanup(rpcSrv.Stop)
	service := sendTxService{chainID: &cltest.FixtureChainID}
	rpcSrv.RegisterName("eth", &service)
	ts := httptest.NewServer(rpcSrv)
	t.Cleanup(ts.Close)

	sendonlyURL := *cltest.MustParseURL(t, ts.URL)
	ethClient := mustNewClient(t, wsURL, sendonlyURL, sendonlyURL)
	defer ethClient.Close()
	err := ethClient.Dial(context.Background())
	require.NoError(t, err)

	err = ethClient.SendTransaction(context.Background(), tx)
	require.NoError(t, err)

	// Unfortunately it's a bit tricky to test this, since there is no
	// synchronization. We have to rely on timing instead.
	require.Eventually(t, func() bool { return service.sentCount.Load() == int32(2) }, cltest.WaitTimeout(t), 500*time.Millisecond)
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

	ctx, cancel := context.WithTimeout(context.Background(), cltest.WaitTimeout(t))
	defer cancel()

	chainId := big.NewInt(123456)
	const headResult = `{"difficulty":"0xf3a00","extraData":"0xd883010503846765746887676f312e372e318664617277696e","gasLimit":"0xffc001","gasUsed":"0x0","hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd1aeb42885a43b72b518182ef893125814811048","mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","nonce":"0x0ece08ea8c49dfd9","number":"0x1","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x218","stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b","timestamp":"0x58318da2","totalDifficulty":"0x1f3a00","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}`
	wsURL := cltest.NewWSServer(t, chainId, func(method string, params gjson.Result) (string, string) {
		if method == "eth_unsubscribe" {
			return "true", ""
		}
		assert.Equal(t, "eth_subscribe", method)
		if assert.True(t, params.IsArray()) {
			require.Equal(t, "newHeads", params.Array()[0].String())
		}
		return `"0x00"`, headResult
	})

	ethClient := mustNewClientWithChainID(t, wsURL, chainId)
	err := ethClient.Dial(context.Background())
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
