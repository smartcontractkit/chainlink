package changeset

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

func TestAdminBatchPauseWorkflows(t *testing.T) {
	t.Parallel()

	testWorkflowID1 := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	testWorkflowID2 := [32]byte{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	t.Run("batch pause workflows - preconditions only", func(t *testing.T) {
		fixture := setupTest(t)

		t.Log("Testing admin batch pause workflows preconditions...")
		changeset := AdminBatchPauseWorkflows{}
		err := changeset.VerifyPreconditions(fixture.rt.Environment(), AdminBatchPauseWorkflowsInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			WorkflowIDs:               [][32]byte{testWorkflowID1, testWorkflowID2},
		})
		require.NoError(t, err, "preconditions should pass")
		t.Log("Admin batch pause workflows preconditions passed")
	})

	t.Run("batch pause with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Testing admin batch pause workflows with MCMS preconditions...")
		changeset := AdminBatchPauseWorkflows{}
		err := changeset.VerifyPreconditions(fixture.rt.Environment(), AdminBatchPauseWorkflowsInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			WorkflowIDs:               [][32]byte{testWorkflowID1},
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		require.NoError(t, err, "MCMS preconditions should pass")
		t.Log("Admin batch pause workflows with MCMS preconditions passed")

		csOutput, err := changeset.Apply(fixture.rt.Environment(), AdminBatchPauseWorkflowsInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			WorkflowIDs:               [][32]byte{testWorkflowID1},
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		require.NoError(t, err, "admin batch pause apply should pass")
		assert.NotNil(t, csOutput, "admin batch pause apply should pass")
		assert.NotNil(t, csOutput.Reports, "admin batch pause apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from admin batch pause")
	})
}

func TestAdminPauseWorkflow(t *testing.T) {
	t.Parallel()

	testWorkflowID := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

	t.Run("pause single workflow - preconditions only", func(t *testing.T) {
		fixture := setupTest(t)

		t.Log("Testing admin pause single workflow preconditions...")
		changeset := AdminPauseWorkflow{}
		err := changeset.VerifyPreconditions(fixture.rt.Environment(), AdminPauseWorkflowInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			WorkflowID:                testWorkflowID,
		})
		require.NoError(t, err, "preconditions should pass")
		t.Log("Admin pause single workflow preconditions passed")
	})

	t.Run("pause single workflow with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Testing admin pause single workflow with MCMS preconditions...")
		changeset := AdminPauseWorkflow{}
		err := changeset.VerifyPreconditions(fixture.rt.Environment(), AdminPauseWorkflowInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			WorkflowID:                testWorkflowID,
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		require.NoError(t, err, "MCMS preconditions should pass")
		t.Log("Admin pause single workflow with MCMS preconditions passed")

		csOutput, err := changeset.Apply(fixture.rt.Environment(), AdminPauseWorkflowInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			WorkflowID:                testWorkflowID,
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		require.NoError(t, err, "admin pause workflow apply should pass")
		assert.NotNil(t, csOutput, "admin pause workflow apply should pass")
		assert.NotNil(t, csOutput.Reports, "admin pause workflow apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from admin pause workflow")
	})
}

func TestAdminPauseAllByOwner(t *testing.T) {
	t.Parallel()

	testOwner := common.HexToAddress("0x1234567890123456789012345678901234567890")

	t.Run("pause all by owner", func(t *testing.T) {
		fixture := setupTest(t)

		t.Log("Testing admin pause all by owner preconditions...")
		changeset := AdminPauseAllByOwner{}
		err := changeset.VerifyPreconditions(fixture.rt.Environment(), AdminPauseAllByOwnerInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			Owner:                     testOwner,
		})
		require.NoError(t, err, "preconditions should pass")
		t.Log("Admin pause all by owner preconditions passed")

		csOutput, err := changeset.Apply(fixture.rt.Environment(), AdminPauseAllByOwnerInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			Owner:                     testOwner,
			Limit:                     big.NewInt(100),
		})
		require.NoError(t, err, "admin pause all by owner apply should pass")
		assert.NotNil(t, csOutput, "admin pause all by owner apply should pass")
		assert.NotNil(t, csOutput.Reports, "admin pause all by owner apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from admin pause all by owner")
	})

	t.Run("pause all by owner with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Testing admin pause all by owner with MCMS preconditions...")
		changeset := AdminPauseAllByOwner{}
		err := changeset.VerifyPreconditions(fixture.rt.Environment(), AdminPauseAllByOwnerInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			Owner:                     testOwner,
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 30 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		require.NoError(t, err, "MCMS preconditions should pass")
		t.Log("Admin pause all by owner with MCMS preconditions passed")

		csOutput, err := changeset.Apply(fixture.rt.Environment(), AdminPauseAllByOwnerInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			Owner:                     testOwner,
			Limit:                     big.NewInt(10),
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		require.NoError(t, err, "admin pause all by owner apply should pass with MCMS")
		assert.NotNil(t, csOutput, "admin pause all by owner apply should pass")
		assert.NotNil(t, csOutput.MCMSTimelockProposals, "admin pause all by owner apply should have MCMS timelock proposals")
		assert.Len(t, csOutput.MCMSTimelockProposals, 1, "expected one MCMS timelock proposal from admin pause all by owner")
	})
}

func TestAdminPauseAllByDON(t *testing.T) {
	t.Parallel()

	testDONFamily := "test-don-family"

	t.Run("pause all by DON", func(t *testing.T) {
		fixture := setupTest(t)

		t.Log("Testing admin pause all by DON preconditions...")
		changeset := AdminPauseAllByDON{}
		err := changeset.VerifyPreconditions(fixture.rt.Environment(), AdminPauseAllByDONInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			DONFamily:                 testDONFamily,
		})
		require.NoError(t, err, "preconditions should pass")
		t.Log("Admin pause all by DON preconditions passed")

		csOutput, err := changeset.Apply(fixture.rt.Environment(), AdminPauseAllByDONInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			DONFamily:                 testDONFamily,
			Limit:                     big.NewInt(10),
		})
		require.NoError(t, err, "admin pause all by DON apply should pass")
		assert.NotNil(t, csOutput, "admin pause all by DON apply should pass")
		assert.NotNil(t, csOutput.Reports, "admin pause all by DON apply should have reports")
		assert.Len(t, csOutput.Reports, 1, "expected one report from admin pause all by DON")
	})

	t.Run("pause all by DON with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Testing admin pause all by DON with MCMS preconditions...")
		changeset := AdminPauseAllByDON{}
		err := changeset.VerifyPreconditions(fixture.rt.Environment(), AdminPauseAllByDONInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			DONFamily:                 testDONFamily,
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		require.NoError(t, err, "MCMS preconditions should pass")
		t.Log("Admin pause all by DON with MCMS preconditions passed")

		csOutput, err := changeset.Apply(fixture.rt.Environment(), AdminPauseAllByDONInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: "test-workflow-registry-v2",
			DONFamily:                 testDONFamily,
			Limit:                     big.NewInt(10),
			MCMSConfig: &contracts.MCMSConfig{
				MinDelay: 1 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		require.NoError(t, err, "admin pause all by DON apply should pass with MCMS")
		assert.NotNil(t, csOutput, "admin pause all by DON apply should pass")
		assert.NotNil(t, csOutput.MCMSTimelockProposals, "admin pause all by DON apply should have MCMS timelock proposals")
		assert.Len(t, csOutput.MCMSTimelockProposals, 1, "expected one MCMS timelock proposal from admin pause all by DON")
	})
}

func TestAdminBatchPauseWorkflowsValidation(t *testing.T) {
	t.Parallel()

	fixture := setupTest(t)

	testWorkflowID1 := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

	tests := []struct {
		name        string
		input       AdminBatchPauseWorkflowsInput
		expectError bool
	}{
		{
			name: "valid input",
			input: AdminBatchPauseWorkflowsInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
				WorkflowIDs:               [][32]byte{testWorkflowID1},
			},
			expectError: false,
		},
		{
			name: "empty workflow IDs",
			input: AdminBatchPauseWorkflowsInput{
				ChainSelector:             fixture.selector,
				WorkflowRegistryQualifier: "test-workflow-registry-v2",
				WorkflowIDs:               [][32]byte{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changeset := AdminBatchPauseWorkflows{}
			err := changeset.VerifyPreconditions(fixture.rt.Environment(), tt.input)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
