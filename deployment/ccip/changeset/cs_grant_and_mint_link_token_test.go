package changeset_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	tokencs "github.com/smartcontractkit/chainlink/deployment/tokens/changesets"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployLinktokenAndTransferOwnershipCS(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Bootstraps: 1,
		Chains:     1,
		Nodes:      4,
	})

	selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chainSelector := selectors[0]
	e, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(tokencs.DeployEVMLinkTokens, tokencs.DeployLinkTokensInput{
			ChainSelectors: []uint64{chainSelector},
		}))

	require.NoError(t, err)

	// Ensure the link token is deployed
	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)
	chain := e.BlockChains.EVMChains()[chainSelector]
	require.NotNil(t, state.Chains[chainSelector].LinkToken)
	linkToken := state.Chains[chainSelector].LinkToken

	recipientAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")

	cfg := changeset.GrantMintRoleAndMintConfig{
		Selector:  chainSelector,
		ToAddress: recipientAddr,
		Amount:    new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18)), // 10 LINK tokens
	}

	expectedMintAmount := new(big.Int).Set(cfg.Amount)
	err = changeset.GrantMintRoleAndMint.VerifyPreconditions(e, cfg)
	require.NoError(t, err)

	output, err := changeset.GrantMintRoleAndMint.Apply(e, cfg)
	require.NoError(t, err)
	require.NotNil(t, output)

	// Verify deployer no longer has mint role
	isMinter, err := linkToken.IsMinter(&bind.CallOpts{}, chain.DeployerKey.From)
	require.NoError(t, err)
	require.False(t, isMinter, "Deployer should not have mint role after changeset execution")

	// Verify recipient received the minted tokens
	balance, err := linkToken.BalanceOf(&bind.CallOpts{}, recipientAddr)
	require.NoError(t, err)
	require.GreaterOrEqual(t, balance.Cmp(expectedMintAmount), 0, "Recipient should have received at least the expected minted tokens")

	// Verify total supply increased
	totalSupply, err := linkToken.TotalSupply(&bind.CallOpts{})
	require.NoError(t, err)
	require.GreaterOrEqual(t, totalSupply.Cmp(expectedMintAmount), 0, "Total supply should be at least the minted amount")
}
