package canonical

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/integration-tests/smoke/ccip/canonical/execctxs"
	"github.com/smartcontractkit/chainlink/integration-tests/smoke/ccip/canonical/scenarios"
)

func Test_MessagingToEOA_EVM2EVM(t *testing.T) {
	// the “where”
	exec := execctxs.NewEVM2EVMCtx(t)
	// the “what”
	scenario := scenarios.NewMessagingToEOAScenario(t, scenarios.ValidationTypeExec)

	scenario.Run(t, exec)
}

func Test_MessagingToReceiver_EVM2Solana(t *testing.T) {
	exec := execctxs.NewEVM2SolanaCtx(t)
	scenario := scenarios.NewMessagingToCCIPReceiverScenario(t, scenarios.ValidationTypeExec)

	scenario.Run(t, exec)
}

func Test_MessagingToReceiver_Solana2EVM(t *testing.T) {
	exec := execctxs.NewSolana2EVM(t)
	scenario := scenarios.NewMessagingToCCIPReceiverScenario(t, scenarios.ValidationTypeExec)

	scenario.Run(t, exec)
}

func Test_MessagingToReceiver_Matrix(t *testing.T) {
	execContexts := execctxs.AllOneToOneExecContexts(t)
	for _, execCtx := range execContexts {
		t.Run(fmt.Sprintf("messaging_%s", execCtx.Name()), func(t *testing.T) {
			scenario := scenarios.NewMessagingToCCIPReceiverScenario(t, scenarios.ValidationTypeExec)
			scenario.Run(t, execCtx)
		})
	}
}

func Test_MessagingToReceiver_Matrix_ListTests(t *testing.T) {
	execCtxNames := execctxs.AllOneToOneExecContextNames()
	for _, execCtx := range execCtxNames {
		t.Run(fmt.Sprintf("messaging_%s", execCtx), func(t *testing.T) {
			// do nothing, just to list the tests.
		})
	}
}

func Test_AllMessagingTests_ListTests(t *testing.T) {
	scenarioNames := scenarios.AllMessagingScenariosNames()
	execCtxNames := execctxs.AllOneToOneExecContextNames()
	for _, execCtxName := range execCtxNames {
		for _, scenarioName := range scenarioNames {
			t.Run(fmt.Sprintf("%s_%s", scenarioName, execCtxName), func(t *testing.T) {
				// do nothing, just to list the tests.
			})
		}
	}
}
