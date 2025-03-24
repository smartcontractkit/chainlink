package testutil

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
)

func MustFundAddressWithLink(t *testing.T, e deployment.Environment, chain deployment.Chain, to common.Address, amount int64) {

	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	require.NoError(t, err)

	linkState, err := commonstate.MaybeLoadLinkTokenChainState(chain, addresses)
	require.NoError(t, err)
	require.NotNil(t, linkState.LinkToken)

	// grant minter permissions - only owner can call this function
	fmt.Println("Granting minter permissions for chain", chain.DeployerKey)
	tx, err := linkState.LinkToken.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	// Mint 'To' address some tokens
	tx, err = linkState.LinkToken.Mint(chain.DeployerKey, to, big.NewInt(amount))
	require.NoError(t, err)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	// 'To' address should have the tokens
	ctx := e.GetContext()
	endBalance, err := linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, to)
	require.NoError(t, err)
	expectedBalance := big.NewInt(amount)
	require.Equal(t, expectedBalance, endBalance)
}

func MustGetLinkBalance(t *testing.T, e deployment.Environment, chain deployment.Chain, addr common.Address) *big.Int {
	addresses, err := e.ExistingAddresses.AddressesForChain(chain.Selector)
	require.NoError(t, err)

	linkState, err := commonstate.MaybeLoadLinkTokenChainState(chain, addresses)
	require.NoError(t, err)

	ctx := e.GetContext()
	endBalance, err := linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, addr)
	require.NoError(t, err)
	return endBalance
}
