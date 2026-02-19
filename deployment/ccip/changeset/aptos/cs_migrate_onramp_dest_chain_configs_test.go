package aptos_test

import (
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

func TestMigrateOnRampDestChainConfigsToV2_Apply(t *testing.T) {
	t.Skip("skipping - no need to run these tests in CI")
	t.Parallel()

	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithAptosChains(1),
	)
	env := deployedEnvironment.Env

	chainSelector := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyAptos))[0]

	cfg := config.MigrateOnRampDestChainConfigsToV2Config{
		ChainSelector:         chainSelector,
		DestChainSelectors:    []uint64{chain_selectors.ETHEREUM_MAINNET.Selector},
		RouterModuleAddresses: []aptos.AccountAddress{{}},
		MCMS: &proposalutils.TimelockConfig{
			MinDelay:     time.Second,
			MCMSAction:   mcmstypes.TimelockActionSchedule,
			OverrideRoot: false,
		},
	}

	_, out, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(aptoscs.MigrateOnRampDestChainConfigsToV2{}, cfg),
	})

	require.NoError(t, err)
	proposals := out[0].MCMSTimelockProposals
	require.Len(t, proposals, 1)
	require.Len(t, proposals[0].Operations, 1)
}
