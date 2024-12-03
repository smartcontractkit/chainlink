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

	require.NoError(t, err)

	addrs, err := respTimelock.AddressBook.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 4)

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

	// Transfer some funds
	tx, err := linkContract.Transfer(chain.DeployerKey, common.HexToAddress(timelockAddress), big.NewInt(15000))
	require.NoError(t, err)
	err = chain.Client.SendTransaction(ctx, tx)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)
	tx, err = linkContract.Transfer(chain.DeployerKey, common.HexToAddress(mcmsAddress), big.NewInt(15000))
	require.NoError(t, err)
	err = chain.Client.SendTransaction(context.Background(), tx)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	timelockContract, err := gethwrappers.NewRBACTimelock(common.HexToAddress(timelockAddress), chain.Client)
	require.NoError(t, err)
	timelocks := map[uint64]*gethwrappers.RBACTimelock{
		chainSelector: timelockContract,
	}
	// Apply the changest
	_, err = changeset.ApplyChangesets(t, env, timelocks, []changeset.ChangesetApplication{
		// this produces proposals, ApplyChangesets will sign & execute them.
		// in practice, signing and executing are separated processes.
		{
			Changeset: changeset.WrapChangeSet(changeset.LinkTransferTimelock),
			Config: changeset.LinkTransferTimelockRequest{
				Transfers: []changeset.LinkTransfer{
					{
						From:  common.HexToAddress(mcmsAddress),
						To:    chain.DeployerKey.From,
						Value: *big.NewInt(10000),
					},
					{
						From:  common.HexToAddress(timelockAddress),
						To:    chain.DeployerKey.From,
						Value: *big.NewInt(10000),
					},
				},
			},
		},
	})
	require.NoError(t, err)

	// Check new balances
	endBalance, err := linkContract.BalanceOf(&bind.CallOpts{Context: ctx}, chain.DeployerKey.From)
	require.NoError(t, err)
	expectedBalance := big.NewInt(20000)
	require.Equal(t, expectedBalance, endBalance)
}
