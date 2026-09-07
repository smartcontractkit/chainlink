package cre

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crelib "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	stellchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/stellar"
	stellarfeature "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/stellar"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const (
	stellarReadWorkflowFile  = "../../smoke/cre/stellar/stellarread/main.go"
	stellarWriteWorkflowFile = "../../smoke/cre/stellar/stellarwrite/main.go"
	stellarConfigPath        = "/configs/workflow-gateway-don-stellar.toml"

	expectLogReadContractBatchOK = "Stellar ReadContract batch passed"
	expectLogStellarWriteFailure = "Stellar write failure observed as expected"
	undeployedContractID         = "CDLZFC3SYJYDZT7K67VZ75HPJVIEUVNIXF47ZG2FB2RMQQVU2HHGCYSC"
)

// Test_CRE_V2_Stellar_Regression runs all Stellar regression scenarios as subtests under a single shared environment.
//
//nolint:paralleltest // subtests share a single Local CRE env
func Test_CRE_V2_Stellar_Regression(t *testing.T) {
	testEnv := t_helpers.SetupTestEnvironmentWithConfig(t, t_helpers.GetTestConfig(t, stellarConfigPath))
	chain := t_helpers.MustStellarChainInEnv(t, testEnv)
	reqSigs := requiredStellarSignatures(t, testEnv)

	t.Run("ReadContractNegative", func(t *testing.T) {
		fixtureID := t_helpers.MustDeployStellarReadFixture(t, chain)
		negatives := []t_helpers.StellarReadContractStep{
			{Name: "trap", ContractID: fixtureID, Function: "trap", ExpectError: true},
			{Name: "unknown_function", ContractID: fixtureID, Function: "no_such_fn", ExpectError: true},
			{Name: "wrong_arity", ContractID: fixtureID, Function: "sum_u32", ArgMode: "echo_u32", ArgU32A: 1, ExpectError: true},
			{Name: "nonexistent_contract", ContractID: undeployedContractID, Function: "get_u32", ExpectError: true},
		}

		workflowName := t_helpers.UniqueStellarWorkflowName("stellar-rc-negative-batch")
		workflowConfig := t_helpers.StellarReadWorkflowConfig{
			ChainSelector: chain.ChainSelector(),
			WorkflowName:  workflowName,
			CronSchedule:  t_helpers.StellarTestCronSchedule(),
			ReadKind:      t_helpers.StellarReadKindReadContract,
			Cases:         negatives,
		}
		userLogsCh, baseMessageCh := t_helpers.StartChipTestSinkWithLogging(t, t_helpers.LogFilePath("stellar_regression", t.Name()))
		workflowID := t_helpers.CompileAndDeployWorkflow(t, testEnv, framework.L, workflowName, &workflowConfig, stellarReadWorkflowFile)
		t_helpers.WatchWorkflowLogs(t, framework.L, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectLogReadContractBatchOK, t_helpers.StellarWorkflowTimeout, t_helpers.WithUserLogWorkflowID(workflowID))
	})

	t.Run("WriteInvalidReceiver", func(t *testing.T) {
		runStellarWriteFailure(t, testEnv, chain, reqSigs, undeployedContractID)
	})

	t.Run("WriteWrongSymbolReceiver", func(t *testing.T) {
		runStellarWriteFailure(t, testEnv, chain, reqSigs, t_helpers.MustDeployStellarReadFixture(t, chain))
	})

	t.Run("WriteFailingReceiver", func(t *testing.T) {
		receiverID, err := stellarfeature.DeployStellarRejectingReceiver(context.Background(), chain)
		require.NoError(t, err, "failed to deploy Stellar rejecting receiver")
		runStellarWriteFailure(t, testEnv, chain, reqSigs, receiverID)
	})
}

// ─── shared helpers ───────────────────────────────────────────────────────────

func requiredStellarSignatures(t *testing.T, tenv *ttypes.TestEnvironment) int {
	t.Helper()
	require.NotNil(t, tenv.Dons, "test environment DON metadata is required")
	dons := tenv.Dons.DonsWithFlag(crelib.StellarCapability)
	require.NotEmpty(t, dons, "could not find a DON hosting the Stellar capability")
	workers, err := dons[0].Workers()
	require.NoError(t, err, "failed to list Stellar DON workers")
	require.NotEmpty(t, workers, "Stellar DON has no worker nodes")
	return (len(workers)-1)/3 + 1
}

// runStellarWriteFailure deploys the stellarwrite workflow targeting the given
// receiver with ExpectFailure=true and waits for the expected failure log.
func runStellarWriteFailure(t *testing.T, env *ttypes.TestEnvironment, chain *stellchain.Blockchain, reqSigs int, receiverContractID string) {
	t.Helper()
	lggr := framework.L
	workflowConfig := t_helpers.StellarWriteWorkflowConfig{
		ChainSelector:      chain.ChainSelector(),
		WorkflowName:       t_helpers.UniqueStellarWorkflowName("stellar-write-fail"),
		ReceiverContractID: receiverContractID,
		RequiredSignatures: reqSigs,
		ExpectFailure:      true,
	}
	logPath := t_helpers.LogFilePath("stellar_regression", t.Name())
	userLogsCh, baseMessageCh := t_helpers.StartChipTestSinkWithLogging(t, logPath)
	workflowID := t_helpers.CompileAndDeployWorkflow(t, env, lggr, workflowConfig.WorkflowName, &workflowConfig, stellarWriteWorkflowFile)
	t_helpers.WatchWorkflowLogs(t, lggr, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectLogStellarWriteFailure, t_helpers.StellarWorkflowTimeout, t_helpers.WithUserLogWorkflowID(workflowID))
}
