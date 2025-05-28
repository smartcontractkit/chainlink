package aptos_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

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
		HomeChainSel:    env.AllChainSelectors()[0],
		RemoteChainSels: env.AllChainSelectorsAptos(),
		MCMS: &proposalutils.TimelockConfig{
			MinDelay:     time.Duration(1) * time.Second,
			MCMSAction:   mcmstypes.TimelockActionSchedule,
			OverrideRoot: false,
		},
		CCIPHomeConfigType: globals.ConfigTypeActive, // TODO: investigate why this is not being used, might be a bug
	}
	env, _, err := commonchangeset.ApplyChangesetsV2(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(aptoscs.SetOCR3Offramp{}, cfg),
	})
	require.NoError(t, err)

	// Load onchain state
	state, err := stateview.LoadOnchainState(env)
	require.NoError(t, err, "must load onchain state")

	// bind ccip aptos
	aptosCCIPAddr := state.AptosChains[env.AllChainSelectorsAptos()[0]].CCIPAddress
	aptosOffRamp := ccip_offramp.Bind(aptosCCIPAddr, env.BlockChains.AptosChains()[env.AllChainSelectorsAptos()[0]].Client)
	ocr3Commit, err := aptosOffRamp.Offramp().LatestConfigDetails(nil, uint8(types.PluginTypeCCIPCommit))
	require.NoError(t, err)
	require.Len(t, ocr3Commit.Signers, 4)
	ocr3Exec, err := aptosOffRamp.Offramp().LatestConfigDetails(nil, uint8(types.PluginTypeCCIPExec))
	require.NoError(t, err)
	require.Len(t, ocr3Exec.Transmitters, 4)
}
