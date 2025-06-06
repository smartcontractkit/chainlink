package evm

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	gethtypes "github.com/ethereum/go-ethereum/core/types"

	evmtypes "github.com/smartcontractkit/chainlink-common/pkg/types/chains/evm"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmmocks "github.com/smartcontractkit/chainlink/v2/common/chains/mocks"
	lpmocks "github.com/smartcontractkit/chainlink/v2/common/logpoller/mocks"
	txmmocks "github.com/smartcontractkit/chainlink/v2/common/txmgr/mocks"

	"github.com/smartcontractkit/chainlink-evm/pkg/types"
)

func TestEVMService(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	chain := evmmocks.NewChain(t)
	txManager := txmmocks.NewMockEvmTxManager(t)
	evmClient := clienttest.NewClient(t)
	poller := lpmocks.NewLogPoller(t)
	ht := headstest.NewTracker[*types.Head](t)

	chain.On("TxManager").Return(txManager).Maybe()
	chain.On("LogPoller").Return(poller).Maybe()
	chain.On("HeadTracker").Return(ht).Maybe()
	chain.On("Client").Return(evmClient).Maybe()

	relayer := &Relayer{
		chain: chain,
	}

	t.Run("RegisterLogTracking", func(t *testing.T) {
		filter := evmtypes.LPFilterQuery{
			Name:         "filter-1",
			Retention:    time.Second,
			Addresses:    []evmtypes.Address{common.HexToAddress("0x123")},
			EventSigs:    []evmtypes.Hash{common.HexToHash("0x321")},
			Topic2:       []evmtypes.Hash{common.HexToHash("0x222")},
			Topic3:       []evmtypes.Hash{common.HexToHash("0x543")},
			Topic4:       []evmtypes.Hash{common.HexToHash("0x432")},
			MaxLogsKept:  100,
			LogsPerBlock: 10,
		}

		poller.On("HasFilter", mock.MatchedBy(func(fname string) bool {
			return fname == filter.Name
		})).Return(false)
		poller.On("RegisterFilter", ctx, mock.MatchedBy(func(f logpoller.Filter) bool {
			return f.LogsPerBlock == filter.LogsPerBlock &&
				f.Retention == filter.Retention &&
				f.Topic2[0] == filter.Topic2[0] &&
				f.Topic3[0] == filter.Topic3[0] &&
				f.Topic4[0] == filter.Topic4[0] &&
				f.EventSigs[0] == filter.EventSigs[0] &&
				f.MaxLogsKept == filter.MaxLogsKept &&
				f.Addresses[0] == filter.Addresses[0] &&
				f.Name == filter.Name
		})).Return(nil)

		err := relayer.RegisterLogTracking(ctx, filter)
		require.NoError(t, err)
	})

	t.Run("GetTransactionByHash", func(t *testing.T) {
		hash := common.HexToHash("0x123")
		nonce := uint64(1)
		to := common.HexToAddress("0x555")
		amount := big.NewInt(1)
		gasLimit := uint64(2)
		gasPrice := big.NewInt(2)
		data := []byte("kitties")

		transaction := gethtypes.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)
		evmClient.On("TransactionByHash", ctx, hash).Return(transaction, nil)
		tx, err := relayer.GetTransactionByHash(ctx, hash)
		require.NoError(t, err)
		require.Equal(t, transaction.Hash().Bytes(), tx.Hash[:])
		require.Equal(t, transaction.Nonce(), tx.Nonce)
		require.Equal(t, transaction.GasPrice(), tx.GasPrice)
		require.Equal(t, transaction.Data(), tx.Data)
		require.Equal(t, transaction.Gas(), tx.Gas)
		require.Equal(t, transaction.To().Bytes(), tx.To[:])
	})
}

func TestConverters(t *testing.T) {
	t.Parallel()

	t.Run("convert head", func(t *testing.T) {
		head := types.Head{
			Timestamp: time.Unix(100000, 100),
			Number:    100,
			Hash:      common.HexToHash("0x123"),
		}
		result := convertHead(&head)
		require.Equal(t, head.Hash.Bytes(), result.Hash[:])
	})

	t.Run("convert transaction", func(t *testing.T) {
		tx := gethtypes.NewTransaction(
			1,
			common.HexToAddress("0xabc123"),
			big.NewInt(1000),
			21000,
			big.NewInt(1e9),
			[]byte{1, 2, 3},
		)

		result := convertTransaction(tx)
		require.NotNil(t, result)
		require.Equal(t, tx.Hash().Bytes(), result.Hash[:])
		require.Equal(t, tx.Nonce(), result.Nonce)
		require.Equal(t, tx.Gas(), result.Gas)
		require.Equal(t, tx.GasPrice(), result.GasPrice)
		require.Equal(t, tx.Value(), result.Value)
		require.Equal(t, tx.To().Bytes(), result.To[:])
		require.Equal(t, tx.Data(), result.Data)
	})
}
