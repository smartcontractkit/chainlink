package aptos_test

import (
	"math/big"
	"testing"
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestDeployRegulatedToken_Apply(t *testing.T) {
	t.Parallel()
	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithAptosChains(1))
	env := deployedEnvironment.Env

	aptosSelectors := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyAptos))
	require.Len(t, aptosSelectors, 1)
	aptosSelector := aptosSelectors[0]

	cfg := config.DeployRegulatedTokenConfig{
		ChainSelector: aptosSelector,
		TokenParams: config.TokenParams{
			MaxSupply: big.NewInt(1_000_000),
			Name:      "RegDeploy",
			Symbol:    "RDT",
			Decimals:  8,
			Icon:      "",
			Project:   "",
		},
		MCMSConfig: &cldfproposalutils.TimelockConfig{
			MinDelay:     time.Second,
			MCMSAction:   mcmstypes.TimelockActionSchedule,
			OverrideRoot: false,
		},
	}

	_, outputs, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(aptoscs.DeployRegulatedToken{}, cfg),
	})
	require.NoError(t, err)
	require.Len(t, outputs, 1)
	output := outputs[0]
	require.Len(t, output.MCMSTimelockProposals, 1)
	require.NotEmpty(t, output.Reports)
	ops := output.MCMSTimelockProposals[0].Operations
	require.NotEmpty(t, ops)
	require.Len(t, ops[0].Transactions, 2)
}
