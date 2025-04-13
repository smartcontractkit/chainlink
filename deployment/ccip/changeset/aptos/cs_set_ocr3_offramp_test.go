package aptos_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	aptoscs "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/stretchr/testify/require"
)

func TestSetOCR3Offramp_Apply(t *testing.T) {
	// Setup environment and config
	deployedEnvironment, _ := testhelpers.NewMemoryEnvironment(
		t,
		testhelpers.WithAptosChains(1),
	)
	env := deployedEnvironment.Env

	cfg := v1_6.SetOCR3OffRampConfig{
		HomeChainSel:       env.AllChainSelectors()[0],
		RemoteChainSels:    env.AllChainSelectorsAptos(),
		CCIPHomeConfigType: globals.ConfigTypeActive,
	}
	env, err := commonchangeset.ApplyChangesetsV2(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(aptoscs.SetOCR3Offramp{}, cfg),
	})
	require.NoError(t, err)

	// Load onchain state
	state, err := changeset.LoadOnchainState(env)
	require.NoError(t, err, "must load onchain state")

	// bind ccip aptos
	aptosCCIPAddr := state.AptosChains[env.AllChainSelectorsAptos()[0]].CCIPAddress
	aptosOffRamp := ccip_offramp.Bind(aptosCCIPAddr, env.AptosChains[env.AllChainSelectorsAptos()[0]].Client)
	ocr3Commit, err := aptosOffRamp.Offramp().LatestConfigDetails(nil, uint8(types.PluginTypeCCIPCommit))
	require.NoError(t, err)
	require.Len(t, ocr3Commit.Signers, 4)
	ocr3Exec, err := aptosOffRamp.Offramp().LatestConfigDetails(nil, uint8(types.PluginTypeCCIPExec))
	require.NoError(t, err)
	require.Len(t, ocr3Exec.Transmitters, 4)
}
