package cre

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	texttmpl "text/template"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	writetarget "github.com/smartcontractkit/chainlink-solana/pkg/solana/write_target"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	df_sol "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	ks_sol "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	mock_capability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"

	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func ExecuteSecureMintTest(t *testing.T, tenv *ttypes.TestEnvironment) {
	creEnvironment := tenv.CreEnvironment
	wrappedBlockchainOutputs := tenv.WrappedBlockchainOutputs
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

	var s setup
	var solChain *cre.WrappedBlockchainOutput
	for _, w := range wrappedBlockchainOutputs {
		if w.BlockchainOutput.Type != blockchain.FamilySolana {
			continue
		}
		s.ForwarderProgramID = mustGetContract(t, ds, w.SolChain.ChainSelector, ks_sol.ForwarderContract)
		s.ForwarderState = mustGetContract(t, ds, w.SolChain.ChainSelector, ks_sol.ForwarderState)
		// we assume we always have just 1 solana chain
		solChain = w
		break
	}
	require.False(t, s.ForwarderProgramID.IsZero(), "failed to receive forwarder program id from blockchains output")
	s.Selector = solChain.SolChain.ChainSelector

	// configure cache program
	framework.L.Info().Msg("Deploy and configure data-feeds cache programs...")
	deployAndConfigureCache(t, &s, *creEnvironment.CldfEnvironment, solChain)
	framework.L.Info().Msg("Successfully deployed and configured")

	// deploy workflow
	framework.L.Info().Msg("Generate and propose secure mint job...")
	jobSpec := createSecureMintWorkflowJobSpec(t, &s, solChain)
	proposeSecureMintJob(t, creEnvironment.CldfEnvironment.Offchain, creEnvironment.DonTopology, jobSpec)
	framework.L.Info().Msgf("Secure mint job is successfully posted. Job spec:\n %v", jobSpec)

	// trigger workflow
	trigger := createFakeTrigger(t, &s, creEnvironment.DonTopology)
	ctx, cancel := context.WithCancel(t.Context())
	eg := &errgroup.Group{}
	eg.Go(func() error {
		return trigger.run(ctx)
	})

	// wait for price update
	waitForFeedUpdate(t, solChain.SolClient, &s)
	cancel()
	require.NoError(t, eg.Wait(), "failed while waiting for feed update")
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
			require.Equal(t, uint32(SeqNr), r.timestamp) // #nosec G115 - we defined seqnr above
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

func getDecimalReportAccount(t *testing.T, s *setup) solana.PublicKey {
	dataID, _ := new(big.Int).SetString(s.FeedID, 0)
	var data [16]byte
	copy(data[:], dataID.Bytes())
	decimalReportSeeds := [][]byte{
		[]byte("decimal_report"),
		s.CacheState.Bytes(),
		data[:],
	}
	decimalReportKey, _, err := solana.FindProgramAddress(decimalReportSeeds, s.CacheProgramID)
	require.NoError(t, err, "failed to derive decimal report key")
	return decimalReportKey
}

type setup struct {
	Selector           uint64
	ForwarderProgramID solana.PublicKey
	ForwarderState     solana.PublicKey
	CacheProgramID     solana.PublicKey
	CacheState         solana.PublicKey

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

func deployAndConfigureCache(t *testing.T, s *setup, env cldf.Environment, solChain *cre.WrappedBlockchainOutput) {
	var d [32]byte
	copy(d[:], []byte(wFDescription))
	s.Descriptions = append(s.Descriptions, d)
	s.WFName = wFName
	s.WFOwner = wFOwner
	s.FeedID = new(big.Int).SetBytes(feedID[:]).String()
	var wfname [10]byte
	copy(wfname[:], []byte(s.WFName))

	ds := datastore.NewMemoryDataStore()
	populateContracts := map[string]datastore.ContractType{
		deployment.DataFeedsCacheProgramName: df_sol.CacheContract,
	}
	err := memory.PopulateDatastore(ds.AddressRefStore, populateContracts, semver.MustParse("1.0.0"), ks_sol.DefaultForwarderQualifier, solChain.SolChain.ChainSelector)
	require.NoError(t, err, "failed to populate datastore")
	env.DataStore = ds.Seal()

	s.CacheProgramID = mustGetContract(t, env.DataStore, solChain.SolChain.ChainSelector, df_sol.CacheContract)
	// deploy df cache
	deployCS := commonchangeset.Configure(df_sol.DeployCache{}, &df_sol.DeployCacheRequest{
		ChainSel:           solChain.SolChain.ChainSelector,
		Qualifier:          ks_sol.DefaultForwarderQualifier,
		Version:            "1.0.0",
		FeedAdmins:         []solana.PublicKey{solChain.SolChain.PrivateKey.PublicKey()},
		ForwarderProgramID: s.ForwarderProgramID,
	})

	// init decimal report
	initCS := commonchangeset.Configure(df_sol.InitCacheDecimalReport{},
		&df_sol.InitCacheDecimalReportRequest{
			ChainSel:  solChain.SolChain.ChainSelector,
			Qualifier: ks_sol.DefaultForwarderQualifier,
			Version:   "1.0.0",
			FeedAdmin: solChain.SolChain.PrivateKey.PublicKey(),
			DataIDs:   []string{s.FeedID},
		})

	// configure decimal report
	configureCS := commonchangeset.Configure(df_sol.ConfigureCacheDecimalReport{},
		&df_sol.ConfigureCacheDecimalReportRequest{
			ChainSel:  solChain.SolChain.ChainSelector,
			Qualifier: ks_sol.DefaultForwarderQualifier,
			Version:   "1.0.0",
			SenderList: []df_sol.Sender{
				{
					ProgramID: s.ForwarderProgramID,
					StateID:   s.ForwarderState,
				},
			},
			FeedAdmin:            solChain.SolChain.PrivateKey.PublicKey(),
			DataIDs:              []string{s.FeedID},
			AllowedWorkflowOwner: [][20]byte{s.WFOwner},
			AllowedWorkflowName:  [][10]byte{wfname},
			Descriptions:         s.Descriptions,
		})
	env, _, cacheErr := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{deployCS, initCS, configureCS})
	require.NoError(t, cacheErr)
	s.CacheProgramID = mustGetContract(t, env.DataStore, solChain.SolChain.ChainSelector, df_sol.CacheContract)
	s.CacheState = mustGetContract(t, env.DataStore, solChain.SolChain.ChainSelector, df_sol.CacheState)
}

const reportSchema = `{
      "kind": "struct",
      "fields": [
        { "name": "payload", "type": { "vec": { "defined": "DecimalReport" } } }
      ]
    }`
const definedTypes = `
     [
      {
        "name":"DecimalReport",
         "type":{
          "kind":"struct",
          "fields":[
            { "name":"timestamp", "type":"u32" },
            { "name":"answer",    "type":"u128" },
            { "name": "dataId",   "type": {"array": ["u8",16]}}
          ]
        }
      }
    ]`

const secureMintWorkflowTemplate = `
name: "{{.WorkflowName}}"
owner: "{{.WorkflowOwner}}"
triggers:
  - id: "securemint-trigger@1.0.0" #currently mocked
    config:
      maxFrequencyMs: 5000
actions:
  - id: "{{.DeriveID}}"
    ref: "solana_data_feeds_cache_accounts"
    inputs:
      trigger_output: $(trigger.outputs) # don't really need it, but without inputs can't pass wf validation
    config:
      Receiver: "{{.DFCacheAddr}}"
      State: "{{.CacheStateID}}"
      FeedIDs: ["{{.FeedID}}"]
consensus:
  - id: "offchain_reporting@1.0.0"
    ref: "secure-mint-consensus"
    inputs:
      observations:
        - event: $(trigger.outputs)
          solana: $(solana_data_feeds_cache_accounts.outputs.remaining_accounts)
    config:
      report_id: "0003"
      key_id: "solana"
      aggregation_method: "secure_mint"
      aggregation_config:
        targetChainSelector: "{{.ChainSelector}}" # CHAIN_ID_FOR_WRITE_TARGET: NEW Param, to match write target
        dataID: "{{.DataID}}"
      encoder: "Borsh"
      encoder_config:
        report_schema: |
          {
            "kind": "struct",
            "fields": [
              { "name": "payload", "type": { "vec": { "defined": "DecimalReport" } } }
            ]
          }
        defined_types: |
          [
            {
              "name": "DecimalReport",
              "type": {
                "kind": "struct",
                "fields": [
                  { "name": "timestamp", "type": "u32" },
                  { "name": "answer",    "type": "u128" },
                  { "name": "dataId",    "type": { "array": ["u8", 16] } }
                ]
              }
            }
          ]

targets:
  - id: "{{.SolanaWriteTargetID}}"
    inputs:
      signed_report: $(secure-mint-consensus.outputs)
      remaining_accounts: $(solana_data_feeds_cache_accounts.outputs.remaining_accounts)
    config:
      address: "{{.DFCacheAddr}}"
      params: ["$(report)"]
      deltaStage: 1s
      schedule: oneAtATime
`

func createSecureMintWorkflowJobSpec(t *testing.T, s *setup, solChain *cre.WrappedBlockchainOutput) string {
	tmpl, err := texttmpl.New("secureMintWorkflow").Parse(secureMintWorkflowTemplate)
	require.NoError(t, err)

	chainID, err := solChain.SolClient.GetGenesisHash(context.Background())
	require.NoError(t, err, "failed to receive genesis hash")

	deriveCapabilityID := writetarget.GenerateDeriveRemainingName(chainID.String())
	writeCapabilityID := writetarget.GenerateWriteTargetName(chainID.String())
	owner := hex.EncodeToString(s.WFOwner[:])
	d, _ := new(big.Int).SetString(s.FeedID, 0)
	data := map[string]any{
		"WorkflowName":        s.WFName,
		"WorkflowOwner":       "0x" + owner,
		"ChainSelector":       s.Selector,
		"DFCacheAddr":         s.CacheProgramID.String(),
		"CacheStateID":        s.CacheState.String(),
		"SolanaWriteTargetID": writeCapabilityID,
		"DeriveID":            deriveCapabilityID,
		"FeedID":              s.FeedID,
		"DataID":              hex.EncodeToString(d.Bytes()),
		"ReportSchema":        reportSchema,
		"DefinedTypes":        definedTypes,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err)

	spec := buf.String()
	workflowJobSpec := testspecs.GenerateWorkflowJobSpec(t, spec)
	return fmt.Sprintf(
		`
		externaljobid   		 	=  "123e4567-e89b-12d3-a456-426655440002"
		%s
	`, workflowJobSpec.Toml())
}

func proposeSecureMintJob(t *testing.T, offchain offchain.Client, donTopology *cre.DonTopology, jobSpec string) {
	workerNodes, err := offchain.ListNodes(t.Context(), &node.ListNodesRequest{
		Filter: &node.ListNodesRequest_Filter{
			Selectors: []*ptypes.Selector{{
				Key:   cre.LabelNodeTypeKey,
				Value: ptr.Ptr(cre.LabelNodeTypeValuePlugin),
				Op:    ptypes.SelectorOp_EQ,
			},
			},
		},
	})
	require.NoError(t, err, "failed to get list nodes")
	var specs cre.DonJobs
	for _, n := range workerNodes.GetNodes() {
		specs = append(specs, &job.ProposeJobRequest{
			Spec:   jobSpec,
			NodeId: n.Id,
		})
	}
	err = jobs.Create(t.Context(), offchain, donTopology, specs)
	if err != nil && strings.Contains(err.Error(), "is already approved") {
		return
	}
	require.NoError(t, err, "failed to propose jobs")
}

type fakeTrigger struct {
	triggerCap *mock_capability.Controller
	setup      *setup
	triggerID  string
}

func (f *fakeTrigger) run(ctx context.Context) error {
	tt := time.NewTicker(time.Second * 21)
	defer tt.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tt.C:
			err := f.Call(context.Background())
			if err != nil {
				return fmt.Errorf("failed call fake trigger: %w", err)
			}
		}
	}
}

func (f *fakeTrigger) Call(ctx context.Context) error {
	outputs, err := f.createReport()
	if err != nil {
		return fmt.Errorf("failed to create fake report: %w", err)
	}

	outputsBytes, err := mock_capability.MapToBytes(outputs)
	if err != nil {
		return fmt.Errorf("failed to convert map to bytes: %w", err)
	}

	message := pb.SendTriggerEventRequest{
		TriggerID: f.triggerID,
		ID:        uuid.New().String(),
		Outputs:   outputsBytes,
	}

	err = f.triggerCap.SendTrigger(ctx, &message)
	if err != nil {
		return fmt.Errorf("failed to send trigger event: %w", err)
	}

	return nil
}

func (f *fakeTrigger) createReport() (*values.Map, error) {
	type secureMintReport struct {
		ConfigDigest ocr2types.ConfigDigest
		SeqNr        uint64
		Block        uint64
		Mintable     *big.Int
	}

	configDigest, _ := hex.DecodeString("000eb2d48aa4727bab3d60885ed3ab7be6e9d6b5855f706b4b01086797ac7730")
	report := &secureMintReport{
		ConfigDigest: ocr2types.ConfigDigest(configDigest),
		SeqNr:        uint64(SeqNr), // #nosec G115 - const conversion
		Block:        uint64(Block), // #nosec G115 - const conversion
		Mintable:     Mintable,
	}

	reportBytes, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}

	ocr3Report := &ocr3types.ReportWithInfo[uint64]{
		Report: ocr2types.Report(reportBytes),
		Info:   f.setup.Selector,
	}

	jsonReport, err := json.Marshal(ocr3Report)
	if err != nil {
		return nil, err
	}

	event, err := values.NewMap(map[string]any{
		"ConfigDigest": configDigest,
		"SeqNr":        SeqNr,
		"Report":       jsonReport,
	})
	if err != nil {
		return nil, err
	}

	return event, nil
}

func createFakeTrigger(t *testing.T, s *setup, dons *cre.DonTopology) *fakeTrigger {
	client := createMockClient(t)
	framework.L.Info().Msg("Successfully exported ocr2 keys")

	return &fakeTrigger{
		triggerCap: client,
		setup:      s,
		triggerID:  "securemint-trigger@1.0.0",
	}
}

func createMockClient(t *testing.T) *mock_capability.Controller {
	in, err := framework.Load[envconfig.Config](nil)
	require.NoError(t, err, "couldn't load environment state")
	mockClientsAddress := make([]string, 0)
	for _, nodeSet := range in.NodeSets {
		for i, n := range nodeSet.NodeSpecs {
			if i == 0 {
				continue
			}
			if len(n.Node.CustomPorts) == 0 {
				panic("no custom port specified, mock capability running in kind must have a custom port in order to connect")
			}
			ports := strings.Split(n.Node.CustomPorts[0], ":")
			mockClientsAddress = append(mockClientsAddress, "127.0.0.1:"+ports[0])
		}
	}

	mocksClient := mock_capability.NewMockCapabilityController(framework.L)
	require.NoError(t, mocksClient.ConnectAll(mockClientsAddress, true, true), " failed to connect mock client")

	return mocksClient
}

func mustGetContract(t *testing.T, ds datastore.DataStore, sel uint64, ctype datastore.ContractType) solana.PublicKey {
	key := datastore.NewAddressRefKey(
		sel,
		ctype,
		semver.MustParse("1.0.0"),
		ks_sol.DefaultForwarderQualifier,
	)
	contract, err := ds.Addresses().Get(key)

	require.NoError(t, err)

	return solana.MustPublicKeyFromBase58(contract.Address)
}
