package cre

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// ExecuteAptosTest runs the Aptos CRE suite.
func ExecuteAptosTest(t *testing.T, tenv *configuration.TestEnvironment) {
	ExecuteAptosReadOnlyTest(t, tenv)
}

func ExecuteAptosReadOnlyTest(t *testing.T, tenv *configuration.TestEnvironment) {
	creEnv := tenv.CreEnvironment
	require.NotEmpty(t, creEnv.Blockchains, "Aptos suite expects at least one blockchain in the environment")

	var aptosChain blockchains.Blockchain
	for _, bc := range creEnv.Blockchains {
		if bc.IsFamily(blockchain.FamilyAptos) {
			aptosChain = bc
			break
		}
	}
	require.NotNil(t, aptosChain, "Aptos suite expects an Aptos chain in the environment (use config workflow-gateway-don-aptos.toml)")

	lggr := framework.L
	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)

	server := t_helpers.StartChipTestSink(t, t_helpers.GetLoggingPublishFn(lggr, userLogsCh, baseMessageCh, "./logs/aptos_capability_workflow_test.log"))
	t.Cleanup(func() {
		server.Shutdown(t.Context())
		close(userLogsCh)
		close(baseMessageCh)
	})

	t.Run("Aptos Read", func(t *testing.T) {
		ExecuteAptosReadTest(t, tenv, aptosChain, userLogsCh, baseMessageCh, lggr)
	})
}

// ExecuteAptosReadTest deploys a workflow that reads 0x1::coin::name() on Aptos local devnet
// in a consensus read step and asserts the expected value.
func ExecuteAptosReadTest(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	aptosChain blockchains.Blockchain,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
	lggr zerolog.Logger,
) {
	const workflowName = "aptos-read-workflow"
	workflowConfig := t_helpers.AptosReadWorkflowConfig{
		ChainSelector:    aptosChain.ChainSelector(),
		WorkflowName:     workflowName,
		ExpectedCoinName: "Aptos Coin",
	}

	const workflowFileLocation = "./aptos/aptosread/main.go"
	t_helpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &workflowConfig, workflowFileLocation)

	expectedLog := "Aptos read consensus succeeded"
	t_helpers.WatchWorkflowLogs(t, lggr, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectedLog, 4*time.Minute)
	lggr.Info().Str("expected_log", expectedLog).Msg("Aptos read capability test passed")
}
