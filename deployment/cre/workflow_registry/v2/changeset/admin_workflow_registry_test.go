package changeset

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

type adminTestFixture struct {
	env                     cldf.Environment
	chainSelector           uint64
	workflowRegistryAddress string
}

func setupWorkflowRegistryAdminTest(t *testing.T) *adminTestFixture {
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  0,
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	chainSelectors := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	require.NotEmpty(t, chainSelectors, "should have at least one EVM chain")
	chainSelector := chainSelectors[0]

	// Deploy the WorkflowRegistry (same pattern as capabilities registry)
	t.Log("Deploying WorkflowRegistry...")
	deployOutput, err := DeployWorkflowRegistry{}.Apply(env, DeployWorkflowRegistryInput{
		ChainSelector: chainSelector,
		Qualifier:     "test-workflow-registry-v2",
	})
	require.NoError(t, err, "failed to deploy WorkflowRegistry")
	require.NotNil(t, deployOutput, "deployment output should not be nil")
	t.Log("WorkflowRegistry deployed successfully")

	// Merge the deployment datastore into the environment (following ApplyChangesets pattern)
	if deployOutput.DataStore != nil {
		ds1 := datastore.NewMemoryDataStore()
		err = ds1.Merge(deployOutput.DataStore.Seal())
		require.NoError(t, err, "failed to merge new addresses into datastore")
		err = ds1.Merge(env.DataStore)
		require.NoError(t, err, "failed to merge current addresses into datastore")
		env.DataStore = ds1.Seal()
	}

	workflowRegistryAddress := deployOutput.DataStore.Addresses().Filter(datastore.AddressRefByQualifier("test-workflow-registry-v2"))[0].Address

	return &adminTestFixture{
		env:                     env,
		chainSelector:           chainSelector,
		workflowRegistryAddress: workflowRegistryAddress,
	}
}

func TestAdminBatchPauseWorkflows(t *testing.T) {
	fixture := setupWorkflowRegistryAdminTest(t)

	testWorkflowID1 := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	testWorkflowID2 := [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	t.Run("batch pause workflows - preconditions only", func(t *testing.T) {
		t.Log("Testing admin batch pause workflows preconditions...")
		changeset := AdminBatchPauseWorkflows{}
		err := changeset.VerifyPreconditions(fixture.env, AdminBatchPauseWorkflowsInput{
			ChainSelector:             fixture.chainSelector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			WorkflowIDs:               [][32]byte{testWorkflowID1, testWorkflowID2},
		})
		require.NoError(t, err, "preconditions should pass")
		t.Log("Admin batch pause workflows preconditions passed")
	})

	t.Run("batch pause with MCMS - preconditions only", func(t *testing.T) {
		t.Log("Testing admin batch pause workflows with MCMS preconditions...")
		changeset := AdminBatchPauseWorkflows{}
		err := changeset.VerifyPreconditions(fixture.env, AdminBatchPauseWorkflowsInput{
			ChainSelector:             fixture.chainSelector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			WorkflowIDs:               [][32]byte{testWorkflowID1},
			MCMSConfig: &ocr3.MCMSConfig{
				MinDuration: 30 * time.Second,
			},
		})
		require.NoError(t, err, "MCMS preconditions should pass")
		t.Log("Admin batch pause workflows with MCMS preconditions passed")
	})
}

func TestAdminPauseWorkflow(t *testing.T) {
	fixture := setupWorkflowRegistryAdminTest(t)

	testWorkflowID := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

	t.Run("pause single workflow - preconditions only", func(t *testing.T) {
		t.Log("Testing admin pause single workflow preconditions...")
		changeset := AdminPauseWorkflow{}
		err := changeset.VerifyPreconditions(fixture.env, AdminPauseWorkflowInput{
			ChainSelector:             fixture.chainSelector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			WorkflowID:                testWorkflowID,
		})
		require.NoError(t, err, "preconditions should pass")
		t.Log("Admin pause single workflow preconditions passed")
	})

	t.Run("pause single workflow with MCMS - preconditions only", func(t *testing.T) {
		t.Log("Testing admin pause single workflow with MCMS preconditions...")
		changeset := AdminPauseWorkflow{}
		err := changeset.VerifyPreconditions(fixture.env, AdminPauseWorkflowInput{
			ChainSelector:             fixture.chainSelector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			WorkflowID:                testWorkflowID,
			MCMSConfig: &ocr3.MCMSConfig{
				MinDuration: 30 * time.Second,
			},
		})
		require.NoError(t, err, "MCMS preconditions should pass")
		t.Log("Admin pause single workflow with MCMS preconditions passed")
	})
}

func TestAdminPauseAllByOwner(t *testing.T) {
	fixture := setupWorkflowRegistryAdminTest(t)

	testOwner := common.HexToAddress("0x1234567890123456789012345678901234567890")

	t.Run("pause all by owner - preconditions only", func(t *testing.T) {
		t.Log("Testing admin pause all by owner preconditions...")
		changeset := AdminPauseAllByOwner{}
		err := changeset.VerifyPreconditions(fixture.env, AdminPauseAllByOwnerInput{
			ChainSelector:             fixture.chainSelector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			Owner:                     testOwner,
		})
		require.NoError(t, err, "preconditions should pass")
		t.Log("Admin pause all by owner preconditions passed")
	})

	t.Run("pause all by owner with MCMS - preconditions only", func(t *testing.T) {
		t.Log("Testing admin pause all by owner with MCMS preconditions...")
		changeset := AdminPauseAllByOwner{}
		err := changeset.VerifyPreconditions(fixture.env, AdminPauseAllByOwnerInput{
			ChainSelector:             fixture.chainSelector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			Owner:                     testOwner,
			MCMSConfig: &ocr3.MCMSConfig{
				MinDuration: 30 * time.Second,
			},
		})
		require.NoError(t, err, "MCMS preconditions should pass")
		t.Log("Admin pause all by owner with MCMS preconditions passed")
	})
}

func TestAdminPauseAllByDON(t *testing.T) {
	fixture := setupWorkflowRegistryAdminTest(t)

	testDONFamily := "test-don-family"

	t.Run("pause all by DON - preconditions only", func(t *testing.T) {
		t.Log("Testing admin pause all by DON preconditions...")
		changeset := AdminPauseAllByDON{}
		err := changeset.VerifyPreconditions(fixture.env, AdminPauseAllByDONInput{
			ChainSelector:             fixture.chainSelector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			DONFamily:                 testDONFamily,
		})
		require.NoError(t, err, "preconditions should pass")
		t.Log("Admin pause all by DON preconditions passed")
	})

	t.Run("pause all by DON with MCMS - preconditions only", func(t *testing.T) {
		t.Log("Testing admin pause all by DON with MCMS preconditions...")
		changeset := AdminPauseAllByDON{}
		err := changeset.VerifyPreconditions(fixture.env, AdminPauseAllByDONInput{
			ChainSelector:             fixture.chainSelector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			DONFamily:                 testDONFamily,
			MCMSConfig: &ocr3.MCMSConfig{
				MinDuration: 30 * time.Second,
			},
		})
		require.NoError(t, err, "MCMS preconditions should pass")
		t.Log("Admin pause all by DON with MCMS preconditions passed")
	})
}

func TestAdminBatchPauseWorkflowsValidation(t *testing.T) {
	fixture := setupWorkflowRegistryAdminTest(t)

	testWorkflowID1 := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

	tests := []struct {
		name        string
		input       AdminBatchPauseWorkflowsInput
		expectError bool
	}{
		{
			name: "valid input",
			input: AdminBatchPauseWorkflowsInput{
				ChainSelector:             fixture.chainSelector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
				WorkflowIDs:               [][32]byte{testWorkflowID1},
			},
			expectError: false,
		},
		{
			name: "empty workflow IDs",
			input: AdminBatchPauseWorkflowsInput{
				ChainSelector:             fixture.chainSelector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
				WorkflowIDs:               [][32]byte{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changeset := AdminBatchPauseWorkflows{}
			err := changeset.VerifyPreconditions(fixture.env, tt.input)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
