package example_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/example"
)

// TestAddMintersBurnersLink tests the AddMintersBurnersLink changeset
func TestAddMintersBurnersLink(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	// Deploy Link Token and Timelock contracts and add addresses to environment
	env := setupLinkTransferTestEnv(t)

	chainSelector := env.AllChainSelectors()[0]
	chain := env.Chains[chainSelector]
	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)

	mcmsState, err := changeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)
	linkState, err := changeset.MaybeLoadLinkTokenChainState(chain, addrs)
	require.NoError(t, err)

	timelockAddress := mcmsState.Timelock.Address()

	// Mint some funds
	_, err = example.AddMintersBurnersLink(env, &example.AddMintersBurnersLinkConfig{
		ChainSelector: chainSelector,
		Minters:       []common.Address{timelockAddress},
		Burners:       []common.Address{timelockAddress},
	})
	require.NoError(t, err)

	// check timelock balance
	isMinter, err := linkState.LinkToken.IsMinter(&bind.CallOpts{Context: ctx}, timelockAddress)
	require.NoError(t, err)
	require.True(t, isMinter)
	isBurner, err := linkState.LinkToken.IsBurner(&bind.CallOpts{Context: ctx}, timelockAddress)
	require.NoError(t, err)
	require.True(t, isBurner)
}
