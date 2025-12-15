package cre

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/solwrite/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	"github.com/stretchr/testify/require"
)

var (
	workflowFileLocation = "123"
)

func ExecuteSolanaWriteTest(t *testing.T, tenv *configuration.TestEnvironment) {
	creEnvironment := tenv.CreEnvironment
	bcs := tenv.CreEnvironment.Blockchains
	ds := creEnvironment.CldfEnvironment.DataStore
	// prevalidate environment
	forwarders := creEnvironment.CldfEnvironment.DataStore.Addresses().Filter(
		datastore.AddressRefByQualifier(ks_sol.DefaultForwarderQualifier),
		datastore.AddressRefByType(ks_sol.ForwarderContract))
	require.Len(t, forwarders, 1)
	forwarderStates := creEnvironment.CldfEnvironment.DataStore.Addresses().Filter(
		datastore.AddressRefByQualifier(ks_sol.DefaultForwarderQualifier),
		datastore.AddressRefByType(ks_sol.ForwarderState))
	require.Len(t, forwarderStates, 1)

	// 1. Get solana chain
	var s setup
	solChain := getSolChain(t, ds, &s, bcs)
	require.False(t, s.ForwarderProgramID.IsZero(), "failed to receive forwarder program id from blockchains output")
	s.Selector = solChain.ChainSelector()
	// 2. Deploy data-feeds cache
	framework.L.Info().Msg("Deploy and configure data-feeds cache programs...")
	deployAndConfigureCache(t, &s, *creEnvironment.CldfEnvironment, solChain)
	testLogger := tenv.Logger
	framework.L.Info().Msg("Successfully deployed and configured")
	// 3. Compile and deploy workflow
	var workflowConfig config.Config
	workflowConfig.Receiver = s.CacheProgramID
	workflowConfig.ForwarderState = s.ForwarderState
	workflowConfig.ForwarderProgramID = s.ForwarderProgramID

	const workflowFileLocation = "./solana/solwrite/main.go"
	//const workflowFileLocation = "../../../../core/scripts/cre/environment/examples/workflows/v2/proof-of-reserve/cron-based/main.go"
	workflowName := fmt.Sprintf("sol-write-workflow--%04d", rand.Intn(10000))
	t_helpers.CompileAndDeployWorkflow(t,
		tenv, testLogger, workflowName, &workflowConfig,
		workflowFileLocation)
	time.Sleep(time.Second * 90)
	// 4. Confirm report submission
}

func getSolChain(t *testing.T, ds datastore.DataStore, s *setup, bcs []blockchains.Blockchain) *solana.Blockchain {
	var solChain *solana.Blockchain
	for _, w := range bcs {
		if !w.IsFamily(chainselectors.FamilySolana) {
			continue
		}
		require.IsType(t, &solana.Blockchain{}, solChain, "expected Solana blockchain type")
		solChain = w.(*solana.Blockchain)
		s.ForwarderProgramID = mustGetContract(t, ds, solChain.ChainSelector(), ks_sol.ForwarderContract)
		s.ForwarderState = mustGetContract(t, ds, solChain.ChainSelector(), ks_sol.ForwarderState)
		// we assume we always have just 1 solana chain
		break
	}

	return solChain
}
