package changeset_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

// TestMintLinkTimelock tests the MintLinkTimelock changeset
func TestMintLinkTimelock(t *testing.T) {
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

	// Deploy Link Token
	resp, err := changeset.DeployLinkToken(env, []uint64{chainSelector})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NoError(t, env.ExistingAddresses.Merge(resp.AddressBook))

	// Deploy MCMS and Timelock
	config := changeset.SingleGroupMCMS(t)
	respTimelock, err := changeset.DeployMCMSWithTimelock(env, map[uint64]types.MCMSWithTimelockConfig{
		chainSelector: {
			Canceller:        config,
			Bypasser:         config,
			Proposer:         config,
			TimelockMinDelay: big.NewInt(0),
		},
	})
	require.NoError(t, env.ExistingAddresses.Merge(respTimelock.AddressBook))
	require.NoError(t, err)

	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)

	mcmsState, err := changeset.MaybeLoadMCMSWithTimelockState(chain, addrs)
	require.NoError(t, err)
	linkState, err := changeset.MaybeLoadLinkTokenState(chain, addrs)
	require.NoError(t, err)
	linkAddress := linkState.LinkToken.Address()
	timelockAddress := mcmsState.Timelock.Address()

	linkContract, err := link_token.NewLinkToken(linkAddress, chain.Client)
	require.NoError(t, err)

	// Mint some funds
	_, err = changeset.MintLinkTimelock(env, &changeset.MintLinkTimelockConfig{
		ChainSelector: chainSelector,
		Amount:        big.NewInt(7568),
	})
	require.NoError(t, err)

	// check timelock balance
	endBalance, err := linkContract.BalanceOf(&bind.CallOpts{Context: ctx}, timelockAddress)
	require.NoError(t, err)
	expectedBalance := big.NewInt(7568)
	require.Equal(t, expectedBalance, endBalance)
}
