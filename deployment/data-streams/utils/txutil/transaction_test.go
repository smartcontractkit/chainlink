package txutil

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestSignAndExecute_ETHTransfer(t *testing.T) {
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

func TestSignAndExecute_ContractInteraction(t *testing.T) {
	ctx := context.Background()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})

	testChain := e.AllChainSelectors()[0]
	chain := e.Chains[testChain]

	e, err := commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			deployment.CreateLegacyChangeSet(commonchangeset.DeployLinkToken),
			[]uint64{testChain},
		),
	)
	require.NoError(t, err)

	addresses, err := e.ExistingAddresses.AddressesForChain(testChain)
	require.NoError(t, err)

	linkState, err := commonchangeset.MaybeLoadLinkTokenChainState(chain, addresses)
	require.NoError(t, err)

	// Mint some funds
	// grant minter permissions
	tx, err := linkState.LinkToken.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	recipient := common.HexToAddress("0x0fd8b81e3d1143ec7f1ce474827ab93c43523ea2")

	tx, err = linkState.LinkToken.Mint(chain.DeployerKey, chain.DeployerKey.From, big.NewInt(750))
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	endBalance, err := linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, chain.DeployerKey.From)
	require.NoError(t, err)
	expectedBalance := big.NewInt(750)
	require.Equal(t, expectedBalance, endBalance)

	// This is the key part - generate a transaction with call data / arguments but do not send it
	tx, err = linkState.LinkToken.Transfer(deployment.SimTransactOpts(), recipient, big.NewInt(500))
	require.NoError(t, err)

	// Ensure that the transaction is not sent to the chain
	endBalance, err = linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, recipient)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0).Int64(), endBalance.Int64())

	preparedTx := &PreparedTx{
		Tx:            tx,
		ChainSelector: testChain,
		ContractType:  types.LinkToken.String(),
	}

	// Execute
	results, err := SignAndExecute(e, []*PreparedTx{preparedTx})
	require.NoError(t, err)
	require.Len(t, results, 1)

	// Verify
	endBalance, err = linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, recipient)
	require.NoError(t, err)
	expectedBalance = big.NewInt(500)
	require.Equal(t, expectedBalance, endBalance)
}
