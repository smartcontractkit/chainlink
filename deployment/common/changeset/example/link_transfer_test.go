package example_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/example"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

// setupLinkTransferContracts deploys all required contracts for the link transfer tests and returns the updated env.
func setupLinkTransferTestEnv(t *testing.T) deployment.Environment {

	lggr := logger.TestLogger(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelector := env.AllChainSelectors()[0]
	config := changeset.SingleGroupMCMS(t)

	// Deploy MCMS and Timelock
	env, err := changeset.ApplyChangesets(t, env, nil, []changeset.ChangesetApplication{
		{
			Changeset: changeset.WrapChangeSet(changeset.DeployLinkToken),
			Config:    []uint64{chainSelector},
		},
		{
			Changeset: changeset.WrapChangeSet(changeset.DeployMCMSWithTimelock),
			Config: map[uint64]types.MCMSWithTimelockConfig{
				chainSelector: {
					Canceller:        config,
					Bypasser:         config,
					Proposer:         config,
					TimelockMinDelay: big.NewInt(0),
				},
			},
		},
	})
	require.NoError(t, err)
	return env
}

// TestLinkTransferMCMS tests the LinkTransfer changeset by sending LINK from a timelock contract
// to the deployer key via mcms proposal.
func TestLinkTransferMCMS(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	env := setupLinkTransferTestEnv(t)
	chainSelector := env.AllChainSelectors()[0]
	chain := env.Chains[chainSelector]
	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)

	mcmsState, err := changeset.MaybeLoadMCMSWithTimelockState(chain, addrs)
	require.NoError(t, err)
	linkState, err := changeset.MaybeLoadLinkTokenState(chain, addrs)
	require.NoError(t, err)
	timelockAddress := mcmsState.Timelock.Address()

	// Mint some funds
	// grant minter permissions
	tx, err := linkState.LinkToken.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tx, err = linkState.LinkToken.Mint(chain.DeployerKey, timelockAddress, big.NewInt(750))
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	timelocks := map[uint64]*changeset.TimelockExecutionContracts{
		chainSelector: {
			Timelock:  mcmsState.Timelock,
			CallProxy: mcmsState.CallProxy,
		},
	}
	// Apply the changeset
	_, err = changeset.ApplyChangesets(t, env, timelocks, []changeset.ChangesetApplication{
		// the changeset produces proposals, ApplyChangesets will sign & execute them.
		// in practice, signing and executing are separated processes.
		{
			Changeset: changeset.WrapChangeSet(example.LinkTransfer),
			Config: &example.LinkTransferConfig{
				From: timelockAddress,
				Transfers: map[uint64][]example.TransferConfig{
					chainSelector: {
						{
							To:    chain.DeployerKey.From,
							Value: big.NewInt(500),
						},
					},
				},
				McmsConfig: &example.MCMSConfig{
					ValidUntil:   4131638958,
					MinDelay:     0,
					OverrideRoot: true,
				},
			},
		},
	})
	require.NoError(t, err)

	// Check new balances
	endBalance, err := linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, chain.DeployerKey.From)
	require.NoError(t, err)
	expectedBalance := big.NewInt(500)
	require.Equal(t, expectedBalance, endBalance)

	// check timelock balance
	endBalance, err = linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, timelockAddress)
	require.NoError(t, err)
	expectedBalance = big.NewInt(250)
	require.Equal(t, expectedBalance, endBalance)
}

// TestLinkTransfer tests the LinkTransfer changeset by sending LINK from a timelock contract to the deployer key.
func TestLinkTransfer(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	env := setupLinkTransferTestEnv(t)
	chainSelector := env.AllChainSelectors()[0]
	chain := env.Chains[chainSelector]
	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)

	mcmsState, err := changeset.MaybeLoadMCMSWithTimelockState(chain, addrs)
	require.NoError(t, err)
	linkState, err := changeset.MaybeLoadLinkTokenState(chain, addrs)
	require.NoError(t, err)
	timelockAddress := mcmsState.Timelock.Address()

	// Mint some funds
	// grant minter permissions
	tx, err := linkState.LinkToken.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tx, err = linkState.LinkToken.Mint(chain.DeployerKey, chain.DeployerKey.From, big.NewInt(750))
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	timelocks := map[uint64]*changeset.TimelockExecutionContracts{
		chainSelector: {
			Timelock:  mcmsState.Timelock,
			CallProxy: mcmsState.CallProxy,
		},
	}

	// Apply the changeset
	_, err = changeset.ApplyChangesets(t, env, timelocks, []changeset.ChangesetApplication{
		// the changeset produces proposals, ApplyChangesets will sign & execute them.
		// in practice, signing and executing are separated processes.
		{
			Changeset: changeset.WrapChangeSet(example.LinkTransfer),
			Config: &example.LinkTransferConfig{
				From: chain.DeployerKey.From,
				Transfers: map[uint64][]example.TransferConfig{
					chainSelector: {
						{
							To:    timelockAddress,
							Value: big.NewInt(500),
						},
					},
				},
				// No MCMSConfig here means we'll execute the txs directly.
			},
		},
	})
	require.NoError(t, err)

	// Check new balances
	endBalance, err := linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, chain.DeployerKey.From)
	require.NoError(t, err)
	expectedBalance := big.NewInt(250)
	require.Equal(t, expectedBalance, endBalance)

	// check timelock balance
	endBalance, err = linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, timelockAddress)
	require.NoError(t, err)
	expectedBalance = big.NewInt(500)
	require.Equal(t, expectedBalance, endBalance)
}
