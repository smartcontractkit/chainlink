package canonical

import (
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
