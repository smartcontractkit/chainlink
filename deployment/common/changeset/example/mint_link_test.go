package example_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/example"
)

// TestMintLink tests the MintLink changeset
func TestMintLink(t *testing.T) {
	t.Parallel()

	var (
		ctx          = t.Context()
		rt, selector = setupLinkTransferRuntime(t) // Deploy Link Token and Timelock contracts and add addresses to environment
	)

	chain := rt.Environment().BlockChains.EVMChains()[selector]

	addrs, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)

	mcmsState, err := changeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)
	linkState, err := changeset.MaybeLoadLinkTokenChainState(chain, addrs)
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(example.AddMintersBurnersLink), &example.AddMintersBurnersLinkConfig{
			ChainSelector: selector,
			Minters:       []common.Address{chain.DeployerKey.From},
		}),
	)
	require.NoError(t, err)

	timelockAddress := mcmsState.Timelock.Address()

	// Mint some funds
	_, err = example.MintLink(rt.Environment(), &example.MintLinkConfig{
		ChainSelector: selector,
		To:            timelockAddress,
		Amount:        big.NewInt(7568),
	})
	require.NoError(t, err)

	// check timelock balance
	endBalance, err := linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, timelockAddress)
	require.NoError(t, err)
	expectedBalance := big.NewInt(7568)
	require.Equal(t, expectedBalance, endBalance)
}
