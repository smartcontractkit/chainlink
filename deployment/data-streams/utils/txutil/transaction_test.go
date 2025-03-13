package txutil

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestTransaction(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})

	testChain := e.AllChainSelectors()[0]
	chain := e.Chains[testChain]

	recipient := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")

	initialBalance, err := chain.Client.BalanceAt(e.GetContext(), recipient, nil)
	require.NoError(t, err)

	// Basic ETH transfer transaction
	value := big.NewInt(1000000000000000000) // 1 ETH

	tx := gethtypes.NewTx(&gethtypes.DynamicFeeTx{
		To:    &recipient,
		Value: value,
		Gas:   21000,
	})

	preparedTx := &PreparedTx{
		Tx:            tx,
		ChainSelector: testChain,
		ContractType:  "ETHTransfer",
	}

	// Execute
	results, err := SignAndExecute(e, []*PreparedTx{preparedTx})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Greater(t, results[0].BlockNumber, uint64(0))

	// Verify
	newBalance, err := chain.Client.BalanceAt(e.GetContext(), recipient, nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0).Add(initialBalance, value), newBalance)
}
