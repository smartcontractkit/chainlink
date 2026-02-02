package cre

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	solanago "github.com/gagliardetto/solana-go"
	solgo "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df_sol "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/evm"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	logtrigger_config "github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/sollogtrigger/config"
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

	waitForFeedUpdate(t, solChain.SolClient, &s)
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

func waitForFeedUpdate(t *testing.T, solclient *rpc.Client, s *setup) {
	tt := time.NewTicker(time.Second * 5)
	defer tt.Stop()
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute*4)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			require.FailNow(t, "The feed failed to update before timeout expired")
		case <-tt.C:
			reportAcc := getDecimalReportAccount(t, s)

			decimalReportAccount, err := solclient.GetAccountInfoWithOpts(t.Context(), reportAcc, &rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentProcessed})
			if errors.Is(err, rpc.ErrNotFound) {
				continue
			}
			require.NoError(t, err, "failed to receive decimal report account")
			// that's how report is stored on chain
			type report struct {
				timestamp uint32   // 4 byte
				answer    *big.Int // 16 byte
			}
			var r report
			data := decimalReportAccount.Value.Data.GetBinary()
			descriminatorLen := 8
			expectedLen := descriminatorLen + 4 + 16
			require.GreaterOrEqual(t, len(data), expectedLen)
			offset := descriminatorLen
			r.timestamp = binary.LittleEndian.Uint32(data[offset : offset+4])
			offset += 4
			answerLE := data[offset : offset+16]
			amount, _, _ := parsePackedU128([16]byte(answerLE))
			r.answer = amount

			if r.answer.Uint64() == 0 {
				framework.L.Info().Msgf("Feed not updated yet.. Retrying...")
				continue
			}
			framework.L.Info().Msg("Feed is updated. Asserting results...")
			require.Equal(t, Mintable.String(), r.answer.String(), "onchain answer value is not equal to sent value")
			return
		}
	}
}

// u128 layout (MSB..LSB): [1 unused][36 block][91 amount]
func parsePackedU128(le [16]byte) (amount *big.Int, block uint64, unused uint8) {
	// Convert LE -> big.Int (big-endian expected by SetBytes)
	be := make([]byte, 16)
	for i := range 16 {
		be[15-i] = le[i]
	}
	x := new(big.Int).SetBytes(be)

	// Masks
	amountMask := new(big.Int).Lsh(big.NewInt(1), 91)
	amountMask.Sub(amountMask, big.NewInt(1)) // (1<<91)-1
	blockMask := new(big.Int).Lsh(big.NewInt(1), 36)
	blockMask.Sub(blockMask, big.NewInt(1)) // (1<<36)-1

	// amount = x & ((1<<91)-1)
	amount = new(big.Int).And(x, amountMask)

	// block = (x >> 91) & ((1<<36)-1)
	blockInt := new(big.Int).Rsh(new(big.Int).Set(x), 91)
	blockInt.And(blockInt, blockMask)
	block = blockInt.Uint64()

	// unused = (x >> 127) & 1
	top := new(big.Int).Rsh(x, 127)
	if top.BitLen() > 0 && top.Bit(0) == 1 {
		unused = 1
	}
	return
}

func getDecimalReportAccount(t *testing.T, s *setup) solanago.PublicKey {
	dataID, _ := new(big.Int).SetString(s.FeedID, 0)
	var data [16]byte
	copy(data[:], dataID.Bytes())
	decimalReportSeeds := [][]byte{
		[]byte("decimal_report"),
		s.CacheState.Bytes(),
		data[:],
	}
	decimalReportKey, _, err := solanago.FindProgramAddress(decimalReportSeeds, s.CacheProgramID)
	require.NoError(t, err, "failed to derive decimal report key")
	return decimalReportKey
}

type setup struct {
	Selector           uint64
	ForwarderProgramID solanago.PublicKey
	ForwarderState     solanago.PublicKey
	CacheProgramID     solanago.PublicKey
	CacheState         solanago.PublicKey

	FeedID       string
	Descriptions [][32]byte
	WFOwner      [20]byte
	WFName       string
}

var (
	feedID        = [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	wFName        = "testwf1234"
	wFDescription = "securemint test"
	wFOwner       = [20]byte{1, 2, 3}
	SeqNr         = 5
	Block         = 10
	Mintable      = big.NewInt(15)
)

func deployAndConfigureCache(t *testing.T, s *setup, env cldf.Environment, solChain *solana.Blockchain) {
	var d [32]byte
	copy(d[:], []byte(wFDescription))
	s.Descriptions = append(s.Descriptions, d)
	s.FeedID = new(big.Int).SetBytes(feedID[:]).String()
	var wfname [10]byte
	copy(wfname[:], []byte(s.WFName))

	ds := datastore.NewMemoryDataStore()

	err := ds.AddressRefStore.Add(datastore.AddressRef{
		Address:       solutils.GetProgramID(solutils.ProgDataFeedsCache),
		ChainSelector: solChain.ChainSelector(),
		Type:          df_sol.CacheContract,
		Version:       semver.MustParse("1.0.0"),
		Qualifier:     ks_sol.DefaultForwarderQualifier,
	})
	require.NoError(t, err, "failed to populate datastore")

	env.DataStore = ds.Seal()

	s.CacheProgramID = mustGetContract(t, env.DataStore, solChain.ChainSelector(), df_sol.CacheContract)
	// deploy df cache
	deployCS := commonchangeset.Configure(df_sol.DeployCache{}, &df_sol.DeployCacheRequest{
		ChainSel:           solChain.ChainSelector(),
		Qualifier:          ks_sol.DefaultForwarderQualifier,
		Version:            "1.0.0",
		FeedAdmins:         []solanago.PublicKey{solChain.PrivateKey.PublicKey()},
		ForwarderProgramID: s.ForwarderProgramID,
	})

	// init decimal report
	initCS := commonchangeset.Configure(df_sol.InitCacheDecimalReport{},
		&df_sol.InitCacheDecimalReportRequest{
			ChainSel:  solChain.ChainSelector(),
			Qualifier: ks_sol.DefaultForwarderQualifier,
			Version:   "1.0.0",
			FeedAdmin: solChain.PrivateKey.PublicKey(),
			DataIDs:   []string{s.FeedID},
		})

	// configure decimal report
	configureCS := commonchangeset.Configure(df_sol.ConfigureCacheDecimalReport{},
		&df_sol.ConfigureCacheDecimalReportRequest{
			ChainSel:  solChain.ChainSelector(),
			Qualifier: ks_sol.DefaultForwarderQualifier,
			Version:   "1.0.0",
			SenderList: []df_sol.Sender{
				{
					ProgramID: s.ForwarderProgramID,
					StateID:   s.ForwarderState,
				},
			},
			FeedAdmin:            solChain.PrivateKey.PublicKey(),
			DataIDs:              []string{s.FeedID},
			AllowedWorkflowOwner: [][20]byte{s.WFOwner},
			AllowedWorkflowName:  [][10]byte{wfname},
			Descriptions:         s.Descriptions,
		})
	env, _, cacheErr := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{deployCS, initCS, configureCS})
	require.NoError(t, cacheErr)
	s.CacheProgramID = mustGetContract(t, env.DataStore, solChain.ChainSelector(), df_sol.CacheContract)
	s.CacheState = mustGetContract(t, env.DataStore, solChain.ChainSelector(), df_sol.CacheState)
}

func mustGetContract(t *testing.T, ds datastore.DataStore, sel uint64, ctype datastore.ContractType) solanago.PublicKey {
	key := datastore.NewAddressRefKey(
		sel,
		ctype,
		semver.MustParse("1.0.0"),
		ks_sol.DefaultForwarderQualifier,
	)
	contract, err := ds.Addresses().Get(key)

	require.NoError(t, err)

	return solanago.MustPublicKeyFromBase58(contract.Address)
}

func ExecuteSolanaLogTriggerTest(t *testing.T, tenv *configuration.TestEnvironment) {
	bcs := tenv.CreEnvironment.Blockchains
	testLogger := tenv.Logger

	var solChain *solana.Blockchain
	for _, w := range bcs {
		if !w.IsFamily(chainselectors.FamilySolana) {
			continue
		}
		require.IsType(t, &solana.Blockchain{}, w, "expected Solana blockchain type")
		solChain = w.(*solana.Blockchain)
		break
	}
	require.NotNil(t, solChain, "Solana blockchain not found in test environment")

	logReadTestProgramID := solanago.MustPublicKeyFromBase58("J1zQwrBNBngz26jRPNWsUSZMHJwBwpkoDitXRV95LdK4")

	workflowName := fmt.Sprintf("sol-logtrigger-wf--%04d", 1234)
	var workflowConfig logtrigger_config.Config
	workflowConfig.LogReadTestProgramID = logReadTestProgramID

	const workflowFileLocation = "./solana/sollogtrigger/main.go"

	listenerCtx, messageChan, kafkaErrChan := t_helpers.StartBeholder(t, testLogger, tenv)

	t_helpers.CompileAndDeployWorkflow(t,
		tenv, testLogger, workflowName, &workflowConfig,
		workflowFileLocation)

	workflowInitMessage := "RunSolLogTriggerWorkflow called"
	err := t_helpers.AssertBeholderMessage(listenerCtx, t, workflowInitMessage, testLogger, messageChan, kafkaErrChan, 2*time.Minute)
	require.NoError(t, err, "Workflow should have initialized")

	err = triggerLogReadTestEvent(t.Context(), solChain, logReadTestProgramID, 42)
	require.NoError(t, err, "failed to trigger log_read_test event")

	timeout := 5 * time.Minute
	expectedLogTriggerMessage := "TestEvent received!"

	err = t_helpers.AssertBeholderMessage(listenerCtx, t, expectedLogTriggerMessage, testLogger, messageChan, kafkaErrChan, timeout)
	require.NoError(t, err, "Log trigger should have received TestEvent")
}

func triggerLogReadTestEvent(ctx context.Context, solChain *solana.Blockchain, programID solanago.PublicKey, value uint64) error {
	discriminator := getCreateLogDiscriminator()

	var instructionData []byte
	instructionData = append(instructionData, discriminator[:]...)

	valueBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(valueBytes, value)
	instructionData = append(instructionData, valueBytes...)

	instruction := solanago.NewInstruction(
		programID,
		solanago.AccountMetaSlice{
			{PublicKey: solChain.PrivateKey.PublicKey(), IsSigner: true, IsWritable: true},
			{PublicKey: solanago.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		instructionData,
	)

	_, err := solCommonUtil.SendAndConfirm(
		ctx,
		solChain.SolClient,
		[]solanago.Instruction{instruction},
		solChain.PrivateKey,
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return fmt.Errorf("failed to send create_log transaction: %w", err)
	}

	return nil
}

func getCreateLogDiscriminator() [8]byte {
	hash := sha256.Sum256([]byte("global:create_log"))
	var discriminator [8]byte
	copy(discriminator[:], hash[:8])
	return discriminator
}
