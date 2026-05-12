package changeset

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/workflow_registry/v2/changeset/operations/contracts"
)

func TestSetConfig(t *testing.T) {
	t.Parallel()

	t.Run("basic metadata config", func(t *testing.T) {
		fixture := setupTest(t)
		t.Log("Starting metadata config...")
		output, err := SetConfig{}.Apply(fixture.rt.Environment(), SetConfigInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			NameLen:                   32,
			TagLen:                    16,
			URLLen:                    128,
			AttrLen:                   256,
			ExpiryLen:                 604800,
			MCMSConfig:                nil,
		})
		t.Logf("Metadata config result: err=%v, output=%v", err, output)
		require.NoError(t, err, "metadata config should succeed")
		require.NotNil(t, output, "output should not be nil")
		t.Log("Metadata config completed successfully")
	})

	t.Run("metadata config with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Starting metadata config with MCMS...")
		output, err := SetConfig{}.Apply(fixture.rt.Environment(), SetConfigInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			NameLen:                   32,
			TagLen:                    16,
			URLLen:                    128,
			AttrLen:                   256,
			ExpiryLen:                 604800,
			MCMSConfig: &crecontracts.MCMSConfig{
				MinDelay: 30 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		t.Logf("MCMS metadata config result: err=%v, output=%v", err, output)
		require.NoError(t, err, "MCMS metadata config should succeed")
		require.NotNil(t, output, "output should not be nil")
		require.NotNil(t, output.MCMSTimelockProposals, "MCMS proposals should be created")
		t.Log("MCMS metadata config completed successfully")
	})
}

func TestUpdateAllowedSigners(t *testing.T) {
	t.Parallel()

	t.Run("update allowed signers", func(t *testing.T) {
		fixture := setupTest(t)

		t.Log("Starting update allowed signers...")
		output, err := UpdateAllowedSigners{}.Apply(fixture.rt.Environment(), UpdateAllowedSignersInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Signers: []common.Address{
				common.HexToAddress("0x1234567890123456789012345678901234567890"),
				common.HexToAddress("0x2234567890123456789012345678901234567890"),
			},
			Allowed:    true,
			MCMSConfig: nil,
		})
		t.Logf("Update allowed signers result: err=%v, output=%v", err, output)
		require.NoError(t, err, "update allowed signers should succeed")
		require.NotNil(t, output, "output should not be nil")
		t.Log("Update allowed signers completed successfully")
	})

	t.Run("update allowed signers with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Starting update allowed signers with MCMS...")
		output, err := UpdateAllowedSigners{}.Apply(fixture.rt.Environment(), UpdateAllowedSignersInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Signers: []common.Address{
				common.HexToAddress("0x1234567890123456789012345678901234567890"),
			},
			Allowed: true,
			MCMSConfig: &crecontracts.MCMSConfig{
				MinDelay: 30 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		t.Logf("MCMS update allowed signers result: err=%v, output=%v", err, output)
		require.NoError(t, err, "MCMS update allowed signers should succeed")
		require.NotNil(t, output, "output should not be nil")
		require.NotNil(t, output.MCMSTimelockProposals, "MCMS proposals should be created")
		t.Log("MCMS update allowed signers completed successfully")
	})
}

func TestSetWorkflowOwnerConfig(t *testing.T) {
	t.Parallel()

	t.Run("set workflow owner config", func(t *testing.T) {
		fixture := setupTest(t)

		t.Log("Starting set workflow owner config...")
		output, err := SetWorkflowOwnerConfig{}.Apply(fixture.rt.Environment(), SetWorkflowOwnerConfigInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Owner:                     common.HexToAddress("0x1234567890123456789012345678901234567890"),
			Config:                    []byte("test config data"),
			MCMSConfig:                nil,
		})
		t.Logf("Set workflow owner config result: err=%v, output=%v", err, output)
		require.NoError(t, err, "set workflow owner config should succeed")
		require.NotNil(t, output, "output should not be nil")
		t.Log("Set workflow owner config completed successfully")
	})

	t.Run("set workflow owner config with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Starting set workflow owner config with MCMS...")
		output, err := SetWorkflowOwnerConfig{}.Apply(fixture.rt.Environment(), SetWorkflowOwnerConfigInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Owner:                     common.HexToAddress("0x1234567890123456789012345678901234567890"),
			Config:                    []byte("test config data"),
			MCMSConfig: &crecontracts.MCMSConfig{
				MinDelay: 30 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		t.Logf("MCMS set workflow owner config result: err=%v, output=%v", err, output)
		require.NoError(t, err, "MCMS set workflow owner config should succeed")
		require.NotNil(t, output, "output should not be nil")
		require.NotNil(t, output.MCMSTimelockProposals, "MCMS proposals should be created")
		t.Log("MCMS set workflow owner config completed successfully")
	})
}

func TestSetDONLimit(t *testing.T) {
	t.Parallel()

	t.Run("set DON limit", func(t *testing.T) {
		fixture := setupTest(t)

		t.Log("Starting set DON limit...")
		output, err := SetDONLimit{}.Apply(fixture.rt.Environment(), SetDONLimitInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			DONFamily:                 "test-don-family",
			DONLimit:                  10,
			UserDefaultLimit:          5,
			MCMSConfig:                nil,
		})
		t.Logf("Set DON limit result: err=%v, output=%v", err, output)
		require.NoError(t, err, "set DON limit should succeed")
		require.NotNil(t, output, "output should not be nil")
		t.Log("Set DON limit completed successfully")
	})

	t.Run("set DON limit with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Starting set DON limit with MCMS...")
		output, err := SetDONLimit{}.Apply(fixture.rt.Environment(), SetDONLimitInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			DONFamily:                 "test-don-family",
			DONLimit:                  10,
			UserDefaultLimit:          5,
			MCMSConfig: &crecontracts.MCMSConfig{
				MinDelay: 30 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		t.Logf("MCMS set DON limit result: err=%v, output=%v", err, output)
		require.NoError(t, err, "MCMS set DON limit should succeed")
		require.NotNil(t, output, "output should not be nil")
		require.NotNil(t, output.MCMSTimelockProposals, "MCMS proposals should be created")
		t.Log("MCMS set DON limit completed successfully")
	})
}

func TestSetUserDONOverride(t *testing.T) {
	t.Parallel()

	t.Run("set user DON override", func(t *testing.T) {
		fixture := setupTest(t)

		// set DON limit first
		_, err := SetDONLimit{}.Apply(fixture.rt.Environment(), SetDONLimitInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			DONFamily:                 "test-don-family",
			DONLimit:                  10,
			UserDefaultLimit:          5,
			MCMSConfig:                nil,
		})
		require.NoError(t, err, "set DON limit should succeed")

		t.Log("Starting set user DON override...")
		output, err := SetUserDONOverride{}.Apply(fixture.rt.Environment(), SetUserDONOverrideInput{
			ChainSelector:             fixture.selector,
			User:                      common.HexToAddress("0x1234567890123456789012345678901234567890"),
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			DONFamily:                 "test-don-family",
			Limit:                     5,
			Enabled:                   true,
			MCMSConfig:                nil,
		})
		t.Logf("Set user DON override result: err=%v, output=%v", err, output)
		require.NoError(t, err, "set user DON override should succeed")
		require.NotNil(t, output, "output should not be nil")
		t.Log("Set user DON override completed successfully")
	})

	t.Run("set user DON override with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		// set DON limit first
		_, err := SetDONLimit{}.Apply(fixture.rt.Environment(), SetDONLimitInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			DONFamily:                 "test-don-family",
			DONLimit:                  10,
			UserDefaultLimit:          5,
			MCMSConfig:                nil,
		})
		require.NoError(t, err, "set DON limit should succeed")

		t.Log("Starting set user DON override with MCMS...")
		output, err := SetUserDONOverride{}.Apply(fixture.rt.Environment(), SetUserDONOverrideInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			User:                      common.HexToAddress("0x1234567890123456789012345678901234567890"),
			DONFamily:                 "test-don-family",
			Limit:                     5,
			Enabled:                   true,
			MCMSConfig: &crecontracts.MCMSConfig{
				MinDelay: 30 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		t.Logf("MCMS set user DON override result: err=%v, output=%v", err, output)
		require.NoError(t, err, "MCMS set user DON override should succeed")
		require.NotNil(t, output, "output should not be nil")
		require.NotNil(t, output.MCMSTimelockProposals, "MCMS proposals should be created")
		t.Log("MCMS set user DON override completed successfully")
	})
}

func TestBatchSetUserDONOverride(t *testing.T) {
	t.Parallel()

	// A small but non-trivial set of overrides reused across subtests.
	makeOverrides := func() []contracts.SetUserDONOverrideEntry {
		return []contracts.SetUserDONOverrideEntry{
			{User: common.HexToAddress("0x1111111111111111111111111111111111111111"), DONFamily: "test-don-family", Limit: 5, Enabled: true},
			{User: common.HexToAddress("0x2222222222222222222222222222222222222222"), DONFamily: "test-don-family", Limit: 7, Enabled: true},
			{User: common.HexToAddress("0x3333333333333333333333333333333333333333"), DONFamily: "test-don-family", Limit: 9, Enabled: false},
		}
	}

	t.Run("verify preconditions: empty overrides rejected", func(t *testing.T) {
		fixture := setupTest(t)
		err := BatchSetUserDONOverride{}.VerifyPreconditions(fixture.rt.Environment(), BatchSetUserDONOverrideInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Overrides:                 nil,
		})
		require.Error(t, err)
	})

	t.Run("verify preconditions: zero user address rejected", func(t *testing.T) {
		fixture := setupTest(t)
		err := BatchSetUserDONOverride{}.VerifyPreconditions(fixture.rt.Environment(), BatchSetUserDONOverrideInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Overrides: []contracts.SetUserDONOverrideEntry{
				{User: common.Address{}, DONFamily: "test-don-family", Limit: 1, Enabled: true},
			},
		})
		require.Error(t, err)
	})

	t.Run("verify preconditions: empty donFamily rejected", func(t *testing.T) {
		fixture := setupTest(t)
		err := BatchSetUserDONOverride{}.VerifyPreconditions(fixture.rt.Environment(), BatchSetUserDONOverrideInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Overrides: []contracts.SetUserDONOverrideEntry{
				{User: common.HexToAddress("0x1111111111111111111111111111111111111111"), DONFamily: "", Limit: 1, Enabled: true},
			},
		})
		require.Error(t, err)
	})

	t.Run("verify preconditions: duplicate (user, donFamily) rejected", func(t *testing.T) {
		fixture := setupTest(t)
		dup := common.HexToAddress("0x1111111111111111111111111111111111111111")
		err := BatchSetUserDONOverride{}.VerifyPreconditions(fixture.rt.Environment(), BatchSetUserDONOverrideInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Overrides: []contracts.SetUserDONOverrideEntry{
				{User: dup, DONFamily: "test-don-family", Limit: 1, Enabled: true},
				{User: dup, DONFamily: "test-don-family", Limit: 2, Enabled: false},
			},
		})
		require.Error(t, err)
	})

	t.Run("batch set user DON override (no MCMS)", func(t *testing.T) {
		fixture := setupTest(t)

		// set DON limit first
		_, err := SetDONLimit{}.Apply(fixture.rt.Environment(), SetDONLimitInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			DONFamily:                 "test-don-family",
			DONLimit:                  100,
			UserDefaultLimit:          5,
			MCMSConfig:                nil,
		})
		require.NoError(t, err, "set DON limit should succeed")

		overrides := makeOverrides()
		t.Log("Starting batch set user DON override...")
		output, err := BatchSetUserDONOverride{}.Apply(fixture.rt.Environment(), BatchSetUserDONOverrideInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Overrides:                 overrides,
			MCMSConfig:                nil,
		})
		t.Logf("Batch set user DON override result: err=%v, output=%v", err, output)
		require.NoError(t, err, "batch set user DON override should succeed")
		require.NotNil(t, output, "output should not be nil")
		// Non-MCMS path: each override was confirmed on-chain individually; no proposals expected.
		require.Empty(t, output.MCMSTimelockProposals, "no MCMS proposals expected without MCMSConfig")
		t.Log("Batch set user DON override completed successfully")
	})

	t.Run("batch set user DON override with MCMS produces single BatchOperation", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		// set DON limit first
		_, err := SetDONLimit{}.Apply(fixture.rt.Environment(), SetDONLimitInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			DONFamily:                 "test-don-family",
			DONLimit:                  100,
			UserDefaultLimit:          5,
			MCMSConfig:                nil,
		})
		require.NoError(t, err, "set DON limit should succeed")

		overrides := makeOverrides()
		t.Log("Starting batch set user DON override with MCMS...")
		output, err := BatchSetUserDONOverride{}.Apply(fixture.rt.Environment(), BatchSetUserDONOverrideInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Overrides:                 overrides,
			MCMSConfig: &crecontracts.MCMSConfig{
				MinDelay: 30 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		t.Logf("MCMS batch set user DON override result: err=%v, output=%v", err, output)
		require.NoError(t, err, "MCMS batch set user DON override should succeed")
		require.NotNil(t, output, "output should not be nil")

		// Core invariant: 1 proposal, 1 batch operation, N transactions inside.
		require.Len(t, output.MCMSTimelockProposals, 1, "expected exactly one MCMS proposal")
		require.Len(t, output.MCMSTimelockProposals[0].Operations, 1, "expected exactly one BatchOperation")
		require.Len(t, output.MCMSTimelockProposals[0].Operations[0].Transactions, len(overrides), "BatchOperation should contain one Transaction per override")
		t.Log("MCMS batch set user DON override completed successfully")
	})

	t.Run("op rejects input.ChainSelector that disagrees with deps.Chain.Selector", func(t *testing.T) {
		// The defensive guard inside BatchSetUserDONOverrideOp prevents a programmatic caller
		// (sequence, other changeset, future refactor) from generating a proposal whose target
		// chain selector disagrees with the chain used to build the calldata. The public
		// Apply() path can't trigger this — both values are derived from the same source — so
		// we exercise the Op directly.
		fixture := setupTest(t)
		env := fixture.rt.Environment()

		deps, err := prepareWorkflowRegistryDeps(env, fixture.selector, fixture.workflowRegistryQualifier, nil, contracts.BatchSetUserDONOverrideDescription)
		require.NoError(t, err, "prepareWorkflowRegistryDeps should succeed for the test fixture chain")

		_, err = operations.ExecuteOperation(env.OperationsBundle, contracts.BatchSetUserDONOverrideOp, deps, contracts.BatchSetUserDONOverrideOpInput{
			ChainSelector: fixture.selector + 1, // deliberately wrong; must not match deps.Chain.Selector
			Qualifier:     fixture.workflowRegistryQualifier,
			Overrides:     makeOverrides(),
		})
		require.Error(t, err, "op must reject mismatched input.ChainSelector vs deps.Chain.Selector")
		require.ErrorContains(t, err, "does not match deps.Chain.Selector")
	})

	t.Run("op works in MCMS path with nil deps.Chain", func(t *testing.T) {
		// MCMS-only callers (e.g. a sequence assembling proposals) may legitimately have
		// deps.Registry + strategy but no chain pointer. In that case the op should still
		// build calldata via SimTransactOpts and produce a BatchOperation; input.ChainSelector
		// is the authoritative source.
		fixture := setupTestWithMCMS(t)
		env := fixture.rt.Environment()

		mcmsConfig := &crecontracts.MCMSConfig{
			MinDelay: 30 * time.Second,
			TimelockQualifierPerChain: map[uint64]string{
				fixture.selector: "",
			},
		}
		deps, err := prepareWorkflowRegistryDeps(env, fixture.selector, fixture.workflowRegistryQualifier, mcmsConfig, contracts.BatchSetUserDONOverrideDescription)
		require.NoError(t, err, "prepareWorkflowRegistryDeps should succeed")
		deps.Chain = nil // simulate an MCMS-only caller that didn't pass a chain pointer

		report, err := operations.ExecuteOperation(env.OperationsBundle, contracts.BatchSetUserDONOverrideOp, deps, contracts.BatchSetUserDONOverrideOpInput{
			ChainSelector: fixture.selector,
			Qualifier:     fixture.workflowRegistryQualifier,
			Overrides:     makeOverrides(),
			MCMSConfig:    mcmsConfig,
		})
		require.NoError(t, err, "op should run in MCMS path without deps.Chain")
		require.NotNil(t, report.Output.MCMSOperation, "MCMS BatchOperation should be produced")
		require.Len(t, report.Output.MCMSOperation.Transactions, len(makeOverrides()), "BatchOperation should contain one Transaction per override")
		require.Equal(t, mcmstypes.ChainSelector(fixture.selector), report.Output.MCMSOperation.ChainSelector, "BatchOperation chain selector must equal input.ChainSelector when deps.Chain is nil")
	})

	t.Run("op requires deps.Chain in non-MCMS path", func(t *testing.T) {
		fixture := setupTest(t)
		env := fixture.rt.Environment()

		deps, err := prepareWorkflowRegistryDeps(env, fixture.selector, fixture.workflowRegistryQualifier, nil, contracts.BatchSetUserDONOverrideDescription)
		require.NoError(t, err, "prepareWorkflowRegistryDeps should succeed")
		deps.Chain = nil // non-MCMS path needs DeployerKey + Confirm, so this must fail loudly

		_, err = operations.ExecuteOperation(env.OperationsBundle, contracts.BatchSetUserDONOverrideOp, deps, contracts.BatchSetUserDONOverrideOpInput{
			ChainSelector: fixture.selector,
			Qualifier:     fixture.workflowRegistryQualifier,
			Overrides:     makeOverrides(),
		})
		require.Error(t, err, "op must reject nil deps.Chain in non-MCMS path")
		require.ErrorContains(t, err, "deps.Chain is required when MCMSConfig is nil")
	})
}

func TestSetCapabilitiesRegistry(t *testing.T) {
	t.Parallel()

	// Test data for DON registry configuration
	donRegistryAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
	donChainSelector := uint64(11155111) // Sepolia chain selector for testing

	t.Run("single DON registry configuration", func(t *testing.T) {
		fixture := setupTest(t)

		t.Log("Starting DON registry configuration...")
		configureOutput, err := SetCapabilitiesRegistry{}.Apply(fixture.rt.Environment(), SetCapabilitiesRegistryInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Registry:                  donRegistryAddress,
			ChainSelectorDON:          donChainSelector,
			MCMSConfig:                nil,
		})
		t.Logf("DON registry configuration result: err=%v, output=%v", err, configureOutput)
		require.NoError(t, err, "DON registry configuration should succeed")
		require.NotNil(t, configureOutput, "configuration output should not be nil")
		require.NotNil(t, configureOutput.Reports, "reports should be present")
		require.Len(t, configureOutput.Reports, 1, "should have exactly one report")
		t.Logf("DON registry configured successfully")
	})

	t.Run("idempotency test - double DON registry configuration", func(t *testing.T) {
		fixture := setupTest(t)

		input := SetCapabilitiesRegistryInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Registry:                  donRegistryAddress,
			ChainSelectorDON:          donChainSelector,
			MCMSConfig:                nil,
		}

		t.Log("Starting first DON registry configuration...")
		configureOutput1, err := SetCapabilitiesRegistry{}.Apply(fixture.rt.Environment(), input)
		require.NoError(t, err, "first configuration should succeed")
		require.NotNil(t, configureOutput1, "first configuration output should not be nil")
		t.Logf("First DON registry configuration completed successfully")

		t.Log("Starting second DON registry configuration (idempotency test)...")
		configureOutput2, err := SetCapabilitiesRegistry{}.Apply(fixture.rt.Environment(), input)
		require.NoError(t, err, "second configuration should succeed (idempotent)")
		require.NotNil(t, configureOutput2, "second configuration output should not be nil")
		t.Logf("Second DON registry configuration completed successfully - idempotency verified")
	})

	t.Run("DON registry configuration with MCMS", func(t *testing.T) {
		fixture := setupTestWithMCMS(t)

		t.Log("Starting DON registry configuration with MCMS...")
		configureOutput, err := SetCapabilitiesRegistry{}.Apply(fixture.rt.Environment(), SetCapabilitiesRegistryInput{
			ChainSelector:             fixture.selector,
			WorkflowRegistryQualifier: fixture.workflowRegistryQualifier,
			Registry:                  donRegistryAddress,
			ChainSelectorDON:          donChainSelector,
			MCMSConfig: &crecontracts.MCMSConfig{
				MinDelay: 30 * time.Second,
				TimelockQualifierPerChain: map[uint64]string{
					fixture.selector: "",
				},
			},
		})
		t.Logf("MCMS DON registry configuration result: err=%v, output=%v", err, configureOutput)
		require.NoError(t, err, "MCMS DON registry configuration should succeed")
		require.NotNil(t, configureOutput, "configuration output should not be nil")
		require.NotNil(t, configureOutput.MCMSTimelockProposals, "MCMS proposals should be created")
		require.NotEmpty(t, configureOutput.MCMSTimelockProposals, "should have at least one MCMS proposal")

		t.Logf("MCMS DON registry configuration completed successfully")
		t.Logf("Created %d MCMS proposals for DON registry configuration", len(configureOutput.MCMSTimelockProposals))

		// Verify proposal content
		for i, proposal := range configureOutput.MCMSTimelockProposals {
			require.NotEmpty(t, proposal.Operations, "proposal %d should have operations", i)
			require.Greater(t, proposal.Delay.Seconds(), float64(0), "proposal %d should have a minimum delay", i)

			for j, op := range proposal.Operations {
				require.NotEmpty(t, op.Transactions, "proposal %d operation %d should have transactions", i, j)
				t.Logf("MCMSOperation %d MCMSOperation %d: %d transactions", i, j, len(op.Transactions))
			}

			t.Logf("MCMSOperation %d: %d operations, delay: %v", i, len(proposal.Operations), proposal.Delay)
		}

		// Verify timelock addresses are set correctly
		for i, proposal := range configureOutput.MCMSTimelockProposals {
			require.NotEmpty(t, proposal.TimelockAddresses, "proposal %d should have timelock addresses", i)
			t.Logf("MCMSOperation %d timelock addresses: %v", i, proposal.TimelockAddresses)
		}

		t.Logf("MCMS DON registry configuration test completed successfully")
		t.Logf("MCMS proposals created and ready for execution through governance")
	})
}

func TestConfigureWorkflowRegistryValidation(t *testing.T) {
	t.Parallel()

	fixture := setupTest(t)

	t.Run("validate SetConfig input", func(t *testing.T) {
		tests := []struct {
			name        string
			input       SetConfigInput
			expectError bool
		}{
			{
				name: "valid input",
				input: SetConfigInput{
					ChainSelector: fixture.selector,
					NameLen:       32,
					TagLen:        16,
					URLLen:        128,
					AttrLen:       256,
					ExpiryLen:     604800,
					MCMSConfig:    nil,
				},
				expectError: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				changeset := SetConfig{}
				err := changeset.VerifyPreconditions(fixture.rt.Environment(), tt.input)
				if tt.expectError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})
		}
	})
}
