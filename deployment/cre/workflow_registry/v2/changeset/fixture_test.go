package changeset

import (
	"testing"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
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
		selector  = chainselectors.TEST_90000001.Selector
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
