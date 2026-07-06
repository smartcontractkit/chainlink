package cre

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	commonevents "github.com/smartcontractkit/chainlink-protos/workflows/go/common"
	workflowevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

const stellarWorkflowTimeout = 4 * time.Minute

// Bounds for proving the local Stellar chain is actively closing ledgers before
// the consensus read runs. stellar/quickstart standalone closes a ledger roughly
// every ~5s, so 30s comfortably covers several closes.
const (
	stellarLedgerAdvanceTimeout = 30 * time.Second
	stellarLedgerAdvancePoll    = 2 * time.Second
)

// ExecuteStellarTest runs the Stellar read CRE smoke scenario: it stands up a
// Chip test sink to capture user logs, then deploys and waits on the read
// workflow.
func ExecuteStellarTest(t *testing.T, tenv *configuration.TestEnvironment) {
	stellarChain := mustStellarChainInEnv(t, tenv)
	lggr := framework.L

	userLogsCh := make(chan *workflowevents.UserLogs, 1000)
	baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
	server := t_helpers.StartChipTestSink(t, t_helpers.GetPublishFn(lggr, userLogsCh, baseMessageCh))
	t.Cleanup(func() {
		// can't use t.Context() here because it will have been cancelled before the cleanup function is called
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t_helpers.ShutdownChipSinkWithDrain(ctx, server, userLogsCh, baseMessageCh)
	})

	ExecuteStellarReadTest(t, tenv, stellarChain, userLogsCh, baseMessageCh)
}

// ExecuteStellarReadTest deploys a workflow that reads the latest ledger on the
// local Stellar standalone network in a consensus read step and asserts the
// returned sequence is past the configured minimum.
func ExecuteStellarReadTest(
	t *testing.T,
	tenv *configuration.TestEnvironment,
	stellarChain blockchains.Blockchain,
	userLogsCh <-chan *workflowevents.UserLogs,
	baseMessageCh <-chan *commonevents.BaseMessage,
) {
	lggr := framework.L

	// Fixed name so re-runs against the same DON overwrite the same workflow.
	const workflowName = "stellar-read-workflow"
	workflowConfig := t_helpers.StellarReadWorkflowConfig{
		ChainSelector:     stellarChain.ChainSelector(),
		WorkflowName:      workflowName,
		MinLedgerSequence: uint64(freshnessFloor),
	}

	const workflowFileLocation = "./stellar/stellarread/main.go"
	workflowID := t_helpers.CompileAndDeployWorkflow(t, tenv, lggr, workflowName, &workflowConfig, workflowFileLocation)

	expectedLog := "Stellar read consensus succeeded"
	t_helpers.WatchWorkflowLogs(t, lggr, userLogsCh, baseMessageCh, t_helpers.WorkflowEngineInitErrorLog, expectedLog, stellarWorkflowTimeout, t_helpers.WithUserLogWorkflowID(workflowID))
	lggr.Info().Str("expected_log", expectedLog).Msg("Stellar read capability test passed")
}

func mustStellarChainInEnv(t *testing.T, tenv *configuration.TestEnvironment) blockchains.Blockchain {
	t.Helper()

	require.NotNil(t, tenv, "Stellar suite requires a test environment")
	require.NotNil(t, tenv.CreEnvironment, "Stellar suite requires a CRE environment")
	require.NotEmpty(t, tenv.CreEnvironment.Blockchains, "Stellar suite expects at least one blockchain in the environment")

	for _, bc := range tenv.CreEnvironment.Blockchains {
		if bc.IsFamily(blockchain.FamilyStellar) {
			return bc
		}
	}

	require.FailNow(t, "Stellar suite expects a Stellar chain in the environment (use config workflow-gateway-don-stellar.toml)")
	return nil
}
