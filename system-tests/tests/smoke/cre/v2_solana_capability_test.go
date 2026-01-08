package cre

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	solgo "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/solwrite/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
	"github.com/stretchr/testify/require"
)

var _ rpc.Client
var _ solgo.Message

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
	workflowName := fmt.Sprintf("sol-write-workflow--%04d", 3411)
	b, _ := hex.DecodeString("35386530643935613437")
	s.WFName = string(b)
	s.WFOwner = tenv.CreEnvironment.Blockchains[0].(*evm.Blockchain).SethClient.Addresses[0]
	deployAndConfigureCache(t, &s, *creEnvironment.CldfEnvironment, solChain)
	testLogger := tenv.Logger
	framework.L.Info().Msg("Successfully deployed and configured")
	// 3. Compile and deploy workflow
	var err error
	var workflowConfig config.Config
	workflowConfig.Receiver = s.CacheProgramID
	workflowConfig.ForwarderState = s.ForwarderState
	workflowConfig.ForwarderProgramID = s.ForwarderProgramID
	workflowConfig.ReceiverState = s.CacheState
	workflowConfig.FeedID, err = dataIDtoBytes(s.FeedID)
	require.NoError(t, err)
	copy(workflowConfig.WFName[:], b)
	workflowConfig.WFOwner = s.WFOwner
	authority, err := deriveForwarderAuthority(workflowConfig.ForwarderState, workflowConfig.Receiver, workflowConfig.ForwarderProgramID)
	hash := createReportHash(workflowConfig.FeedID[:], authority.Bytes(), workflowConfig.WFOwner[:], workflowConfig.WFName[:])
	require.NoError(t, err)
	log.Printf("~~ repHash:%x dataID: %x,sender:%x  owner: %x, name: %x", hash, workflowConfig.FeedID, authority, workflowConfig.WFOwner, workflowConfig.WFName)
	writeFlagSeeds := [][]byte{
		[]byte("permission_flag"),
		workflowConfig.ReceiverState.Bytes(),
		hash[:],
	}

	writeFlagKey, _, err := solgo.FindProgramAddress(writeFlagSeeds, workflowConfig.Receiver)
	log.Printf("~~~write flag: %x name: %x length: %d", writeFlagKey, b, len(b))
	const workflowFileLocation = "./solana/solwrite/main.go"

	t_helpers.CompileAndDeployWorkflow(t,
		tenv, testLogger, workflowName, &workflowConfig,
		workflowFileLocation)
	time.Sleep(time.Second * 90)
	// 4. Confirm report submission
}

func createReportHash(dataID []byte, forwarderAuthority []byte, workflowOwner []byte, workflowName []byte) [32]byte {
	var data []byte
	data = append(data, dataID...)
	data = append(data, forwarderAuthority...)
	data = append(data, workflowOwner...)
	data = append(data, workflowName...)

	return sha256.Sum256(data)
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

func dataIDtoBytes(dataID string) ([16]byte, error) {
	var out [16]byte
	bigID, ok := new(big.Int).SetString(dataID, 0)
	if !ok {
		return out, fmt.Errorf("invalid data_id: %v", dataID)
	}
	if bigID.BitLen() > 128 {
		return out, fmt.Errorf("data_id is too long: %d", bigID.BitLen())
	}

	copy(out[:], bigID.Bytes())
	return out, nil
}

func deriveForwarderAuthority(forwarderState solgo.PublicKey, receiverProgram solgo.PublicKey, forwarderProgram solgo.PublicKey) (solgo.PublicKey, error) {
	seeds := [][]byte{
		[]byte("forwarder"),
		forwarderState[:],
		receiverProgram[:],
	}
	ret, _, err := solgo.FindProgramAddress(seeds, forwarderProgram)

	return ret, err
}
