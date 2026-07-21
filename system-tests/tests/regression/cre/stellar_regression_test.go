package cre

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
)

// Paths are relative to system-tests/tests/regression/cre (go test package cwd).
const (
	stellarReadWorkflowFile = "../../smoke/cre/stellar/stellarread/main.go"
	stellarConfigPath       = "/configs/workflow-gateway-don-stellar.toml"
)

// expectLogReadContractBatchOK must match the stellarread workflow's aggregate success log.
// For the negative suite every case sets ExpectError, so a "batch passed" line means each
// case failed in exactly the expected way.
const expectLogReadContractBatchOK = "Stellar ReadContract batch passed"

// undeployedContractID is a valid StrKey C-address not deployed on the local network.
const undeployedContractID = "CDLZFC3SYJYDZT7K67VZ75HPJVIEUVNIXF47ZG2FB2RMQQVU2HHGCYSC"

// Test_CRE_V2_Stellar_ReadContract_Regression covers negative ReadContract paths
// (trap, unknown fn, wrong arity, nonexistent contract) in a single workflow trigger.
// Requires a Local CRE started with workflow-gateway-don-stellar.toml. Contract WASM is
// built at runtime via stellar CLI (STELLAR_CONTRACTS_SOURCE_DIR or auto go-list module Dir).
//
//nolint:paralleltest // would flake
func Test_CRE_V2_Stellar_ReadContract_Regression(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, stellarConfigPath))
	t.Run("StellarReadContractNegative", func(t *testing.T) {
		lggr := framework.L

		scenarioEnv := t_helpers.SetupTestEnvironmentWithPerTestKeys(t, testEnv.TestConfig)
		stellarChain := t_helpers.MustStellarChainInEnv(t, scenarioEnv)
		logPath := t_helpers.LogFilePath("stellar_regression", t.Name())
		userLogsCh, baseMessageCh := t_helpers.StartChipTestSinkWithLogging(t, logPath)
		lggr.Info().Str("log_file", logPath).Msg("Starting Stellar ReadContract regression")

		fixtureID := t_helpers.MustDeployStellarReadFixture(t, stellarChain)
		negatives := []t_helpers.StellarReadContractStep{
			{Name: "trap", ContractID: fixtureID, Function: "trap", ExpectError: true},
			{Name: "unknown_function", ContractID: fixtureID, Function: "no_such_fn", ExpectError: true},
			// sum_u32 expects two u32s; echo_u32 supplies only one.
			{Name: "wrong_arity", ContractID: fixtureID, Function: "sum_u32", ArgMode: "echo_u32", ArgU32A: 1, ExpectError: true},
			{Name: "nonexistent_contract", ContractID: undeployedContractID, Function: "get_u32", ExpectError: true},
		}

		workflowName := t_helpers.UniqueStellarWorkflowName("stellar-rc-negative-batch")
		workflowConfig := t_helpers.StellarReadWorkflowConfig{
			ChainSelector: stellarChain.ChainSelector(),
			WorkflowName:  workflowName,
			CronSchedule:  t_helpers.StellarTestCronSchedule(),
			ReadKind:      t_helpers.StellarReadKindReadContract,
			Cases:         negatives,
		}
		workflowID := t_helpers.CompileAndDeployWorkflow(t, scenarioEnv, lggr, workflowName, &workflowConfig, stellarReadWorkflowFile)
		t_helpers.WatchWorkflowLogs(t, lggr, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectLogReadContractBatchOK, t_helpers.StellarWorkflowTimeout, t_helpers.WithUserLogWorkflowID(workflowID))
		lggr.Info().Int("cases", len(negatives)).Str("expected_log", expectLogReadContractBatchOK).Msg("Stellar ReadContract negative regression passed")
	})
}
