package cre

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	solgo "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

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
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/solwrite/config"
	t_helpers "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers"
	"github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
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
	const workflowFileLocation = "./solana/solwrite/main.go"

	t_helpers.CompileAndDeployWorkflow(t,
		tenv, testLogger, workflowName, &workflowConfig,
		workflowFileLocation)

	waitForFeedUpdate(t, solChain.SolClient, &s)
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

func waitForFeedUpdate(t *testing.T, solclient *rpc.Client, s *setup) {
	tt := time.NewTicker(time.Second * 2)
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

func getDecimalReportAccount(t *testing.T, s *setup) solgo.PublicKey {
	dataID, _ := new(big.Int).SetString(s.FeedID, 0)
	var data [16]byte
	copy(data[:], dataID.Bytes())
	decimalReportSeeds := [][]byte{
		[]byte("decimal_report"),
		s.CacheState.Bytes(),
		data[:],
	}
	decimalReportKey, _, err := solgo.FindProgramAddress(decimalReportSeeds, s.CacheProgramID)
	require.NoError(t, err, "failed to derive decimal report key")
	return decimalReportKey
}

type setup struct {
	Selector           uint64
	ForwarderProgramID solgo.PublicKey
	ForwarderState     solgo.PublicKey
	CacheProgramID     solgo.PublicKey
	CacheState         solgo.PublicKey

	FeedID       string
	Descriptions [][32]byte
	WFOwner      [20]byte
	WFName       string
}

var (
	feedID        = [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	wFDescription = "securemint test"
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
		FeedAdmins:         []solgo.PublicKey{solChain.PrivateKey.PublicKey()},
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

func mustGetContract(t *testing.T, ds datastore.DataStore, sel uint64, ctype datastore.ContractType) solgo.PublicKey {
	key := datastore.NewAddressRefKey(
		sel,
		ctype,
		semver.MustParse("1.0.0"),
		ks_sol.DefaultForwarderQualifier,
	)
	contract, err := ds.Addresses().Get(key)

	require.NoError(t, err)

	return solgo.MustPublicKeyFromBase58(contract.Address)
}
