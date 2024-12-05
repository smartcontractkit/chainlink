package changeset_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

// TestLinkTransferTimelock tests the LinkTransferTimelock changeset by
func TestLinkTransferTimelock(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelector := env.AllChainSelectors()[0]
	chain := env.Chains[chainSelector]
	// Deploy Value Token
	resp, err := changeset.DeployLinkToken(env, chainSelector)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NoError(t, env.ExistingAddresses.Merge(resp.AddressBook))

	// Deploy MCMS and Timelock
	config := changeset.SingleGroupMCMS(t)
	respTimelock, err := changeset.DeployMCMSWithTimelock(env, map[uint64]types.MCMSWithTimelockConfig{
		chainSelector: {
			Canceller:         config,
			Bypasser:          config,
			Proposer:          config,
			TimelockExecutors: []common.Address{chain.DeployerKey.From},
			TimelockMinDelay:  big.NewInt(0),
		},
	})
	require.NoError(t, env.ExistingAddresses.Merge(respTimelock.AddressBook))
	require.NoError(t, err)

	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 5)

	timelockAddress := ""
	linkAddress := ""
	mcmsAddress := ""
	for key, typeAndVer := range addrs {
		if typeAndVer.Type == types.LinkToken {
			linkAddress = key
		}
		if typeAndVer.Type == types.RBACTimelock {
			timelockAddress = key
		}
		if typeAndVer.Type == types.ProposerManyChainMultisig {
			mcmsAddress = key
		}
	}

	linkContract, err := link_token.NewLinkToken(common.HexToAddress(linkAddress), chain.Client)
	require.NoError(t, err)

	// Mint some funds
	// grant minter permissions
	tx, err := linkContract.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tx, err = linkContract.Mint(chain.DeployerKey, common.HexToAddress(timelockAddress), big.NewInt(750))
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	timelockContract, err := gethwrappers.NewRBACTimelock(common.HexToAddress(timelockAddress), chain.Client)
	require.NoError(t, err)
	timelocks := map[uint64]*gethwrappers.RBACTimelock{
		chainSelector: timelockContract,
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
							Value: *big.NewInt(500),
						},
					},
				},
				LinkTokenAddress: map[uint64]common.Address{
					chainSelector: common.HexToAddress(linkAddress),
				},
				TimelockAddress: map[uint64]common.Address{
					chainSelector: common.HexToAddress(timelockAddress),
				},
				MCMSAddress: map[uint64]common.Address{
					chainSelector: common.HexToAddress(mcmsAddress),
				},
				ValidUntil:   4131638958,
				MinDelay:     0,
				OverrideRoot: true,
				StartingOpCount: map[uint64]uint64{
					chainSelector: 0,
				},
			},
		},
	})
	require.NoError(t, err)

	// Check new balances
	endBalance, err := linkContract.BalanceOf(&bind.CallOpts{Context: ctx}, chain.DeployerKey.From)
	require.NoError(t, err)
	expectedBalance := big.NewInt(500)
	require.Equal(t, expectedBalance, endBalance)

	// check timelock balance
	endBalance, err = linkContract.BalanceOf(&bind.CallOpts{Context: ctx}, common.HexToAddress(timelockAddress))
	require.NoError(t, err)
	expectedBalance = big.NewInt(250)
	require.Equal(t, expectedBalance, endBalance)
}
