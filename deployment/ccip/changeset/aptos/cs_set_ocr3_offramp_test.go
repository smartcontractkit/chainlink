package aptos_test

import (
	"testing"
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
)

func TestSetOCR3Offramp_Apply(t *testing.T) {
	// Setup environment and config
	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithAptosChains(1),
	)
	env := deployedEnvironment.Env

	cfg := v1_6.SetOCR3OffRampConfig{
		HomeChainSel:    env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0],
		RemoteChainSels: env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyAptos)),
		MCMS: &proposalutils.TimelockConfig{
			MinDelay:     time.Duration(1) * time.Second,
			MCMSAction:   mcmstypes.TimelockActionSchedule,
			OverrideRoot: false,
		},
		CCIPHomeConfigType: globals.ConfigTypeActive, // TODO: investigate why this is not being used, might be a bug
	}

	env, _, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(aptoscs.SetOCR3Offramp{}, cfg),
	})
	require.NoError(t, err)

	// Load onchain state
	state, err := stateview.LoadOnchainState(env)
	require.NoError(t, err, "must load onchain state")

	// bind ccip aptos
	aptosCCIPAddr := state.AptosChains[env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyAptos))[0]].CCIPAddress
	aptosOffRamp := ccip_offramp.Bind(aptosCCIPAddr, env.BlockChains.AptosChains()[env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyAptos))[0]].Client)
	ocr3Commit, err := aptosOffRamp.Offramp().LatestConfigDetails(nil, uint8(types.PluginTypeCCIPCommit))
	require.NoError(t, err)
	require.Len(t, ocr3Commit.Signers, 4)
	require.NotEmpty(t, ocr3Commit.ConfigInfo.ConfigDigest)
	ocr3Exec, err := aptosOffRamp.Offramp().LatestConfigDetails(nil, uint8(types.PluginTypeCCIPExec))
	require.NoError(t, err)
	require.Len(t, ocr3Exec.Transmitters, 4)
	require.NotEmpty(t, ocr3Exec.ConfigInfo.ConfigDigest)
}
