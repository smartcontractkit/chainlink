package example_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/example"
)

// TestMintLink tests the MintLink changeset
func TestMintLink(t *testing.T) {
	t.Parallel()
	env := setupLinkTransferTestEnv(t)
	ctx := env.GetContext()
	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]
	chain := env.BlockChains.EVMChains()[chainSelector]

	addrs, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)

	mcmsState, err := changeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)
	linkState, err := changeset.MaybeLoadLinkTokenChainState(chain, addrs)
	require.NoError(t, err)

	_, err = changeset.Apply(t, env,
		changeset.Configure(
			cldf.CreateLegacyChangeSet(example.AddMintersBurnersLink),
			&example.AddMintersBurnersLinkConfig{
				ChainSelector: chainSelector,
				Minters:       []common.Address{chain.DeployerKey.From},
			},
		),
	)
	require.NoError(t, err)

	timelockAddress := mcmsState.Timelock.Address()

	// Mint some funds
	_, err = example.MintLink(env, &example.MintLinkConfig{
		ChainSelector: chainSelector,
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
