package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

const (
	SignerAddress    = "0x0469abeD5aE5f30cb558eE7286b469B6fF35AB34"
	SignerPrivateKey = "3888c3e07120e200d085e15f1bd4132695aa722427d692951af3b5c2c6b1ddf9"
)

type testFixture struct {
	rt                        *runtime.Runtime
	selector                  uint64
	workflowRegistryAddress   string
	workflowRegistryQualifier string
}

// setupTest sets up an runtime with a single EVM simulated chain and a deployed WorkflowRegistry.
func setupTest(t *testing.T) *testFixture {
	var (
		qualifier = "test-workflow-registry-v2"
		selector  = chainselectors.GETH_TESTNET.Selector
	)

	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(DeployWorkflowRegistry{}, DeployWorkflowRegistryInput{
			ChainSelector: selector,
			Qualifier:     qualifier,
		}),
	)
	require.NoError(t, err, "failed to deploy WorkflowRegistry")

	workflowRegistryAddress := rt.State().DataStore.Addresses().Filter(
		datastore.AddressRefByQualifier("test-workflow-registry-v2"),
	)[0].Address

	// UpdateAllowedSigners to include the allowed owner key
	err = rt.Exec(
		runtime.ChangesetTask(UpdateAllowedSigners{}, UpdateAllowedSignersInput{
			ChainSelector:             selector,
			WorkflowRegistryQualifier: qualifier,
			Signers: []common.Address{
				common.HexToAddress(SignerAddress),
			},
			Allowed: true,
		}),
	)
	require.NoError(t, err, "failed to update allowed signers")
	t.Logf("Upated allowed signers to include key %s", SignerAddress)

	// SetDONLimit to allow workflows to be created
	err = rt.Exec(
		runtime.ChangesetTask(SetDONLimit{}, SetDONLimitInput{
			ChainSelector:             selector,
			WorkflowRegistryQualifier: qualifier,
			DONFamily:                 "zone-a",
			DONLimit:                  100,
			UserDefaultLimit:          100,
			MCMSConfig:                nil,
		},
		),
	)
	require.NoError(t, err, "failed to set DON limit")

	return &testFixture{
		rt:                        rt,
		selector:                  selector,
		workflowRegistryAddress:   workflowRegistryAddress,
		workflowRegistryQualifier: qualifier,
	}
}

func setupTestWithMCMS(t *testing.T) *testFixture {
	fixture := setupTest(t)

	err := fixture.rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), map[uint64]commontypes.MCMSWithTimelockConfigV2{
			fixture.selector: proposalutils.SingleGroupTimelockConfigV2(t),
		}),
	)
	require.NoError(t, err, "failed to deploy MCMS")

	return fixture
}
