package changeset_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

// TestLinkTransferTimelock tests the LinkTransferTimelock changeset by
func TestLinkTransferTimelock(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	lggr := logger.TestLogger(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelector := env.AllChainSelectors()[0]
	chain := env.Chains[chainSelector]
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
			Changeset: changeset.WrapChangeSet(changeset.LinkTransferTimelock),
			Config: &changeset.LinkTransferTimelockRequest{
				Transfers: map[uint64][]changeset.LinkTransfer{
					chainSelector: {
						{
							To:    chain.DeployerKey.From,
							Value: big.NewInt(500),
						},
					},
				},
				ValidUntil:   4131638958,
				MinDelay:     0,
				OverrideRoot: true,
				StartingOpCount: map[uint64]uint64{
					chainSelector: 0,
				},
				UseMCMS: true,
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
