package changeset_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	tokencs "github.com/smartcontractkit/chainlink/deployment/tokens/changesets"
)

func TestDeployLinktokenAndTransferOwnershipCS(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(tokencs.DeployEVMLinkTokens, tokencs.DeployLinkTokensInput{
			ChainSelectors: []uint64{selector},
		}),
	)
	require.NoError(t, err)

	// Ensure the link token is deployed
	state, err := stateview.LoadOnchainState(rt.Environment())
	require.NoError(t, err)
	chain := rt.Environment().BlockChains.EVMChains()[selector]
	require.NotNil(t, state.Chains[selector].LinkToken)
	linkToken := state.Chains[selector].LinkToken

	recipientAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")

	cfg := changeset.GrantMintRoleAndMintConfig{
		Selector:  selector,
		ToAddress: recipientAddr,
		Amount:    new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18)), // 10 LINK tokens
	}

	expectedMintAmount := new(big.Int).Set(cfg.Amount)
	err = changeset.GrantMintRoleAndMint.VerifyPreconditions(rt.Environment(), cfg)
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(changeset.GrantMintRoleAndMint, cfg),
	)
	require.NoError(t, err)

	// Verify deployer no longer has mint role
	isMinter, err := linkToken.IsMinter(&bind.CallOpts{}, chain.DeployerKey.From)
	require.NoError(t, err)
	require.False(t, isMinter, "Deployer should not have mint role after changeset execution")

	// Verify deployer no longer has burn role
	isBurner, err := linkToken.IsBurner(&bind.CallOpts{}, chain.DeployerKey.From)
	require.NoError(t, err)
	require.False(t, isBurner, "Deployer should not have burn role after changeset execution")

	// Verify recipient received the minted tokens
	balance, err := linkToken.BalanceOf(&bind.CallOpts{}, recipientAddr)
	require.NoError(t, err)
	require.GreaterOrEqual(t, balance.Cmp(expectedMintAmount), 0, "Recipient should have received at least the expected minted tokens")

	// Verify total supply increased
	totalSupply, err := linkToken.TotalSupply(&bind.CallOpts{})
	require.NoError(t, err)
	require.GreaterOrEqual(t, totalSupply.Cmp(expectedMintAmount), 0, "Total supply should be at least the minted amount")
}

func TestDeployLinktokenAndGrantMintRole(t *testing.T) {
	t.Parallel()

	selector := chainselectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(tokencs.DeployEVMLinkTokens, tokencs.DeployLinkTokensInput{
			ChainSelectors: []uint64{selector},
		}),
	)
	require.NoError(t, err)

	// Ensure the link token is deployed
	state, err := stateview.LoadOnchainState(rt.Environment())
	require.NoError(t, err)
	require.NotNil(t, state.Chains[selector].LinkToken)
	linkToken := state.Chains[selector].LinkToken

	newMinter := common.HexToAddress("0x1234567890123456789012345678901234567890")

	cfg := changeset.GrantMintRoleInput{
		GrantMintRoleByChain: map[uint64]changeset.GrantMintRoleConfig{
			selector: {
				ToAddress: newMinter,
			},
		},
	}

	err = rt.Exec(
		runtime.ChangesetTask(changeset.GrantMintRole, cfg),
	)
	require.NoError(t, err)

	// Verify the newMinter has the mint role
	isMinter, err := linkToken.IsMinter(&bind.CallOpts{}, newMinter)
	require.NoError(t, err)
	require.True(t, isMinter, "New minter should have mint role after changeset execution")

	// Verify the newMinter has the burn role
	isBurner, err := linkToken.IsBurner(&bind.CallOpts{}, newMinter)
	require.NoError(t, err)
	require.True(t, isBurner, "New minter should have burn role after changeset execution")
}
