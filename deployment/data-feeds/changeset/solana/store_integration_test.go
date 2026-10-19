package solana

import (
	"crypto/sha256"
	"encoding/binary"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/fees"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	df_store "github.com/smartcontractkit/chainlink-solana/contracts/generated/data_feeds_store"
	"github.com/smartcontractkit/chainlink-solana/contracts/generated/mock_forwarder"
	"github.com/smartcontractkit/chainlink-solana/contracts/generated/ocr_2"
	df1_store "github.com/smartcontractkit/chainlink-solana/contracts/generated/store"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
)

// storeProgramsPathEnv points at chainlink-solana's contracts/target/deploy, holding the
// data_feeds_store, store and mock_forwarder programs.
// TODO: load them from the chainlink-solana release artifacts once data-feeds-store is released.
const storeProgramsPathEnv = "DATA_FEEDS_STORE_PROGRAMS_PATH"

var (
	storeTestWorkflowOwner = [20]uint8{19: 11}
	storeTestWorkflowName  = [10]uint8{9: 10}
)

// TestDataFeedsStore runs the data feeds store changesets against the store, the DF1 store and a
// mock forwarder (the keystone forwarder without signature checks), through to a DF1 cutover and
// rollback, directly and through MCMS.
func TestDataFeedsStore(t *testing.T) {
	skipInCI(t)
	src := os.Getenv(storeProgramsPathEnv)
	if src == "" {
		t.Skipf("set %s to a directory with data_feeds_store.so, store.so and mock_forwarder.so", storeProgramsPathEnv)
	}

	selector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	programsPath := t.TempDir()
	_, programIDs := soltestutils.LoadMCMSPrograms(t, programsPath)
	for name, id := range map[string]solana.PublicKey{
		"data_feeds_store": df_store.ProgramID,
		"store":            df1_store.ProgramID,
		"mock_forwarder":   mock_forwarder.ProgramID,
	} {
		copyProgram(t, filepath.Join(src, name+".so"), filepath.Join(programsPath, name+".so"))
		programIDs[name] = id.String()
	}

	e, err := environment.New(t.Context(),
		environment.WithSolanaContainer(t, []uint64{selector}, programsPath, programIDs),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)
	chain := e.BlockChains.SolanaChains()[selector]
	deployer := chain.DeployerKey.PublicKey()

	// The store program is preloaded, so DeployStore only initializes a state
	ds := datastore.NewMemoryDataStore()
	soltestutils.RegisterMCMSPrograms(t, selector, ds)
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		Address:       df_store.ProgramID.String(),
		ChainSelector: selector,
		Type:          StoreContract,
		Version:       semver.MustParse("1.0.0"),
		Qualifier:     testQualifier,
	}))
	env := *e
	env.DataStore = ds.Seal()
	mcmsState, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolanaV2(env, ds, chain,
		cldfproposalutils.MCMSWithTimelockConfig{
			Canceller:        cldftesthelpers.SingleGroupMCMS(t),
			Proposer:         cldftesthelpers.SingleGroupMCMS(t),
			Bypasser:         cldftesthelpers.SingleGroupMCMS(t),
			TimelockMinDelay: big.NewInt(0),
		})
	require.NoError(t, err)
	soltestutils.FundSignerPDAs(t, chain, mcmsState)
	env.DataStore = ds.Seal()
	timelockSigner, err := helpers.FetchTimelockSigner(env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(selector)))
	require.NoError(t, err)
	mcmsCfg := &cldfproposalutils.TimelockConfig{MinDelay: time.Second}

	apply := func(t *testing.T, cs commonchangeset.ConfiguredChangeSet) error {
		t.Helper()
		next, _, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{cs})
		if err == nil {
			env = next
		}
		return err
	}
	target := StoreTarget{ChainSel: selector, Version: "1.0.0", Qualifier: testQualifier}
	targetMCMS := target
	targetMCMS.MCMS = mcmsCfg

	const (
		dataIDHex    = "0x0003aabbccddeeff0100000000000000"
		svrDataIDHex = "0x0003aabbccddeeff0100010000000000"
		proxyIDHex   = "0x0101010101010101010101010101010101010101010101010101010101010101"
	)
	dataID, err := ParseDataID(dataIDHex)
	require.NoError(t, err)
	svrDataID, err := ParseDataID(svrDataIDHex)
	require.NoError(t, err)

	// deploy and initialize, with a mock forwarder allowed to report
	forwarderState, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	require.NoError(t, apply(t, commonchangeset.Configure(DeployStore{}, &DeployStoreRequest{
		ChainSel:                  selector,
		Qualifier:                 testQualifier,
		Version:                   "1.0.0",
		FeedAdmins:                []solana.PublicKey{deployer},
		ForwarderProgramID:        mock_forwarder.ProgramID,
		DefaultMaxStoredReports:   16,
		MaxObservationForwardSecs: 3600,
	})))
	store, err := loadStore(env.DataStore, target.ref())
	require.NoError(t, err)
	state := readAccount(t, chain, store.State, df_store.ParseAccount_StoreState)
	require.Equal(t, deployer, state.Owner)
	require.Equal(t, uint64(16), state.DefaultMaxStoredReports)
	require.Equal(t, mock_forwarder.ProgramID, state.ForwarderId)

	confirm(t, chain, build(t)(mock_forwarder.NewInitializeInstruction(forwarderState.PublicKey(), deployer, solana.SystemProgramID)), forwarderState)

	t.Run("configure feeds", func(t *testing.T) {
		workflows := []StoreWorkflow{{
			Forwarder: Sender{ProgramID: mock_forwarder.ProgramID, StateID: forwarderState.PublicKey()},
			Owner:     storeTestWorkflowOwner,
			Name:      storeTestWorkflowName,
		}}
		require.NoError(t, apply(t, commonchangeset.Configure(SetStoreDefaultWorkflowMetadata{}, &SetStoreDefaultWorkflowMetadataRequest{
			StoreTarget: target, Workflows: workflows,
		})))
		state := readAccount(t, chain, store.State, df_store.ParseAccount_StoreState)
		require.Equal(t, uint64(1), state.DefaultWorkflowMetadata.Len)
		sender, err := deriveForwarderAuthority(forwarderState.PublicKey(), store.ProgramID, mock_forwarder.ProgramID)
		require.NoError(t, err)
		require.Equal(t, sender, state.DefaultWorkflowMetadata.Xs[0].AllowedSender)

		require.NoError(t, apply(t, commonchangeset.Configure(InitStoreFeeds{}, &InitStoreFeedsRequest{
			StoreTarget: target, DataIDs: []string{dataIDHex, svrDataIDHex},
		})))
		// idempotent
		require.NoError(t, apply(t, commonchangeset.Configure(InitStoreFeeds{}, &InitStoreFeedsRequest{
			StoreTarget: target, DataIDs: []string{dataIDHex},
		})))

		require.NoError(t, apply(t, commonchangeset.Configure(SetStoreMaxStoredReports{}, &SetStoreMaxStoredReportsRequest{
			StoreTarget: target, Feeds: []StoreMaxStoredReports{{DataID: dataIDHex, MaxStoredReports: 32}},
		})))
		require.NoError(t, apply(t, commonchangeset.Configure(SetStoreMinMaxAnswers{}, &SetStoreMinMaxAnswersRequest{
			StoreTarget: target, Feeds: []StoreMinMaxAnswer{{DataID: dataIDHex, Min: "1", Max: "1000000000000000000000000000"}},
		})))
		feedAddress, err := store.feed(dataID)
		require.NoError(t, err)
		feed := readAccount(t, chain, feedAddress, df_store.ParseAccount_FeedData)
		require.Equal(t, uint64(32), feed.MaxStoredReports)
		require.Equal(t, uint8(1), feed.MinMaxConfigured)
		require.Equal(t, "1000000000000000000000000000", feed.MaxAnswer.BigInt().String())

		require.NoError(t, apply(t, commonchangeset.Configure(SetStoreProxyConfigs{}, &SetStoreProxyConfigsRequest{
			StoreTarget: target,
			Proxies:     []StoreProxyConfig{{ProxyID: proxyIDHex, DataID: svrDataIDHex, Decimals: 8, FallbackTime: 10, MaxSyncIterations: 5}},
		})))
		proxyID, err := parseProxyID(proxyIDHex)
		require.NoError(t, err)
		proxyAddress, err := storePDA(store.ProgramID, store.State, proxySeed, proxyID[:])
		require.NoError(t, err)
		proxy := readAccount(t, chain, proxyAddress, df_store.ParseAccount_ProxyConfig)
		require.Equal(t, uint8(8), proxy.Decimals)

		require.NoError(t, apply(t, commonchangeset.Configure(SetStoreDescriptions{}, &SetStoreDescriptionsRequest{
			StoreTarget: target,
			Feeds:       []StoreFeedDescription{{DataID: svrDataIDHex, Description: "ETH / USD"}},
			SvrPaths:    []StoreSvrPathDescription{{SvrPathID: "0x0001", Description: "SVR path 1"}},
		})))

		require.NoError(t, apply(t, commonchangeset.Configure(SetStoreSvrDisabled{}, &SetStoreSvrDisabledRequest{
			StoreTarget: target, Paths: []StoreSvrDisabled{{DataID: svrDataIDHex, Disabled: true}},
		})))
		svrConfigAddress, err := store.dataIDConfig(svrDataID)
		require.NoError(t, err)
		require.Equal(t, uint8(1), readAccount(t, chain, svrConfigAddress, df_store.ParseAccount_DataIdConfig).SvrDisabled)

		// only SVR paths get a data ID config, the non-SVR settings live in the feed
		configAddress, err := store.dataIDConfig(dataID)
		require.NoError(t, err)
		_, err = getAccount(t.Context(), chain.Client, configAddress)
		require.Error(t, err)

		// feed-specific metadata: in the feed for the non-SVR data ID, in its config for the SVR path
		require.NoError(t, apply(t, commonchangeset.Configure(SetStoreSpecificWorkflowMetadata{}, &SetStoreSpecificWorkflowMetadataRequest{
			StoreTarget: target, DataIDs: []string{dataIDHex, svrDataIDHex}, Workflows: workflows,
		})))
		require.Equal(t, uint64(1), readAccount(t, chain, feedAddress, df_store.ParseAccount_FeedData).Settings.SpecificWorkflowMetadata.Len)
		require.Equal(t, uint64(1), readAccount(t, chain, svrConfigAddress, df_store.ParseAccount_DataIdConfig).Settings.SpecificWorkflowMetadata.Len)
		require.NoError(t, apply(t, commonchangeset.Configure(SetStoreSpecificWorkflowMetadata{}, &SetStoreSpecificWorkflowMetadataRequest{
			StoreTarget: target, DataIDs: []string{dataIDHex, svrDataIDHex},
		})))
		require.Equal(t, uint64(0), readAccount(t, chain, feedAddress, df_store.ParseAccount_FeedData).Settings.SpecificWorkflowMetadata.Len)
		require.Equal(t, uint64(0), readAccount(t, chain, svrConfigAddress, df_store.ParseAccount_DataIdConfig).Settings.SpecificWorkflowMetadata.Len)

		// rejected before anything is sent
		require.ErrorContains(t, apply(t, commonchangeset.Configure(SetStoreMinMaxAnswers{}, &SetStoreMinMaxAnswersRequest{
			StoreTarget: target, Feeds: []StoreMinMaxAnswer{{DataID: svrDataIDHex, Min: "1", Max: "2"}},
		})), "is an SVR data ID")
	})

	// DF1 feed at 8 decimals, written by the deployer standing in for OCR2
	legacyFeed, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	space := uint64(df1HeaderSize + df1TransmissionSize)
	rent, err := chain.Client.GetMinimumBalanceForRentExemption(t.Context(), space, rpc.CommitmentConfirmed)
	require.NoError(t, err)
	confirm(t, chain, build(t)(system.NewCreateAccountInstruction(rent, space, df1_store.ProgramID, deployer, legacyFeed.PublicKey()).ValidateAndBuild()), legacyFeed)
	confirm(t, chain, build(t)(df1_store.NewCreateFeedInstruction("ETH / USD", 8, 1, 1, legacyFeed.PublicKey(), deployer)))
	confirm(t, chain, build(t)(df1_store.NewSetWriterInstruction(deployer, legacyFeed.PublicKey(), deployer, deployer)))
	submitDF1 := func(t *testing.T, answer int64) {
		confirm(t, chain, build(t)(df1_store.NewSubmitInstruction(df1_store.NewTransmission{
			Timestamp: uint64(time.Now().Unix() - 3),
			Answer:    toInt128(big.NewInt(answer)),
		}, legacyFeed.PublicKey(), deployer)))
	}
	submitDF1(t, 330000000000) // 3300.00
	df1Feeds := []DF1Feed{{DataID: dataIDHex, Feed: legacyFeed.PublicKey()}}
	storeAuthority, err := store.authority()
	require.NoError(t, err)

	report := func(t *testing.T, answer *big.Int) {
		sendMockReport(t, chain, store, forwarderState.PublicKey(), dataID, answer, uint32(time.Now().Unix()-1), legacyFeed.PublicKey())
	}
	cutover := func(maxDeviationBps uint64, mcms *cldfproposalutils.TimelockConfig) commonchangeset.ConfiguredChangeSet {
		cutoverTarget := target
		cutoverTarget.MCMS = mcms
		return commonchangeset.Configure(CutoverDF1Feeds{}, &CutoverDF1FeedsRequest{
			StoreTarget:       cutoverTarget,
			DF1StoreProgramID: df1_store.ProgramID,
			Feeds:             df1Feeds,
			MaxStaleness:      5 * time.Minute,
			MaxDeviationBps:   maxDeviationBps,
		})
	}
	df1Writer := func(t *testing.T) solana.PublicKey {
		f, err := loadDF1Feed(t.Context(), chain.Client, df1_store.ProgramID, legacyFeed.PublicKey())
		require.NoError(t, err)
		return f.Header.Writer
	}
	df1Latest := func(t *testing.T) *big.Int {
		f, err := loadDF1Feed(t.Context(), chain.Client, df1_store.ProgramID, legacyFeed.PublicKey())
		require.NoError(t, err)
		latest, err := f.latest()
		require.NoError(t, err)
		return latest.Answer
	}
	answer18 := func(s string) *big.Int {
		v, ok := new(big.Int).SetString(s, 10)
		require.True(t, ok)
		return v
	}

	t.Run("cutover and rollback", func(t *testing.T) {
		require.ErrorContains(t, apply(t, cutover(10, nil)), "report callback program")

		require.NoError(t, apply(t, commonchangeset.Configure(ConfigureDF1ReportCallbacks{}, &ConfigureDF1ReportCallbacksRequest{
			StoreTarget: target, DF1StoreProgramID: df1_store.ProgramID, Feeds: df1Feeds,
		})))
		feedAddress, err := store.feed(dataID)
		require.NoError(t, err)
		settings := readAccount(t, chain, feedAddress, df_store.ParseAccount_FeedData).Settings
		require.Equal(t, df1_store.ProgramID, settings.CallbackProgram)
		require.Equal(t, uint8(1), settings.CallbackAccountsLen)
		require.Equal(t, legacyFeed.PublicKey(), settings.CallbackAccounts[0])
		require.Equal(t, uint8(1), settings.CallbackWritableMask)
		require.ErrorContains(t, apply(t, cutover(10, nil)), "the store has no report yet")

		// the store takes the report, the DF1 feed skips it as OCR2 still writes it
		report(t, answer18("3250120000000000000000")) // 3250.12
		require.Equal(t, uint64(1), readAccount(t, chain, feedAddress, df_store.ParseAccount_FeedData).RoundId)
		require.Equal(t, int64(330000000000), df1Latest(t).Int64())

		require.ErrorContains(t, apply(t, cutover(10, nil)), "deviates")

		submitDF1(t, 325012000000)
		require.NoError(t, apply(t, cutover(10, nil)))
		require.Equal(t, storeAuthority, df1Writer(t))
		require.NoError(t, apply(t, cutover(10, nil)), "already cut over feeds are skipped")

		// the store now writes the DF1 feed
		report(t, answer18("3251000000000000000000"))
		require.Equal(t, int64(325100000000), df1Latest(t).Int64())

		ocr2State := createAccount(t, chain, ocr_2.ProgramID)
		ocr2Authority, err := OCR2StoreAuthority(ocr_2.ProgramID, ocr2State)
		require.NoError(t, err)
		rollback := func(ocr2Program solana.PublicKey) commonchangeset.ConfiguredChangeSet {
			return commonchangeset.Configure(RollbackDF1Feeds{}, &RollbackDF1FeedsRequest{
				ChainSel:          selector,
				DF1StoreProgramID: df1_store.ProgramID,
				OCR2ProgramID:     ocr2Program,
				Feeds:             []RollbackDF1Feed{{Feed: legacyFeed.PublicKey(), OCR2State: ocr2State}},
			})
		}
		require.ErrorContains(t, apply(t, rollback(mock_forwarder.ProgramID)), "not the OCR2 program")
		require.NoError(t, apply(t, rollback(ocr_2.ProgramID)))
		require.Equal(t, ocr2Authority, df1Writer(t))
		require.NoError(t, apply(t, rollback(ocr_2.ProgramID)), "feeds already written by OCR2 are skipped")

		// back on OCR2, the DF1 feed skips store reports again
		report(t, answer18("3252000000000000000000"))
		require.Equal(t, int64(325100000000), df1Latest(t).Int64())
	})

	t.Run("with mcms", func(t *testing.T) {
		require.NoError(t, apply(t, commonchangeset.Configure(SetStoreFeedAdmins{}, &SetStoreFeedAdminsRequest{
			ChainSel: selector, Version: "1.0.0", Qualifier: testQualifier, Add: []solana.PublicKey{timelockSigner},
		})))
		require.NoError(t, apply(t, commonchangeset.Configure(SetStoreDescriptions{}, &SetStoreDescriptionsRequest{
			StoreTarget: targetMCMS,
			Feeds:       []StoreFeedDescription{{DataID: dataIDHex, Description: "ETH / USD (CRE)"}},
		})))
		descriptionAddress, err := storePDA(store.ProgramID, store.State, feedDescriptionSeed, func() []byte { id := indexedFeedID(dataID); return id[:] }())
		require.NoError(t, err)
		require.Contains(t, string(readAccount(t, chain, descriptionAddress, df_store.ParseAccount_FeedDescription).Description[:]), "ETH / USD (CRE)")

		// the timelock takes over the DF1 feed, then proposes the cutover
		confirm(t, chain, build(t)(df1_store.NewTransferFeedOwnershipInstruction(timelockSigner, legacyFeed.PublicKey(), deployer, deployer)))
		require.NoError(t, apply(t, changesetFunc(func(env cldf.Environment) (cldf.ChangesetOutput, error) {
			ix, err := storeInstruction(df1_store.ProgramID, nil)(df1_store.NewAcceptFeedOwnershipInstruction(legacyFeed.PublicKey(), timelockSigner, timelockSigner))
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			return sendInstructions(env, selector, mcmsCfg, "accept DF1 feed ownership", []solana.Instruction{ix})
		})))
		require.ErrorContains(t, apply(t, cutover(10, nil)), "not by the signer")

		report(t, answer18("3251000000000000000000")) // back in line with the DF1 feed
		require.NoError(t, apply(t, cutover(10, mcmsCfg)))
		require.Equal(t, storeAuthority, df1Writer(t))

		require.NoError(t, apply(t, commonchangeset.Configure(TransferOwnershipStore{}, &TransferOwnershipStoreRequest{
			ChainSel: selector, Version: "1.0.0", Qualifier: testQualifier, MCMSCfg: *mcmsCfg,
		})))
		require.Equal(t, timelockSigner, readAccount(t, chain, store.State, df_store.ParseAccount_StoreState).Owner)
	})
}

type changesetFunc func(cldf.Environment) (cldf.ChangesetOutput, error)

func (f changesetFunc) Apply(e cldf.Environment) (cldf.ChangesetOutput, error) { return f(e) }

// sendMockReport reports one answer for a non-SVR data ID through the mock forwarder, writing
// through to the DF1 feed configured as its callback.
func sendMockReport(t *testing.T, chain cldfsol.Chain, store storeAddresses, forwarderState solana.PublicKey,
	dataID [16]byte, answer *big.Int, observedAt uint32, df1Feed solana.PublicKey,
) {
	t.Helper()

	forwarderAuthority, err := deriveForwarderAuthority(forwarderState, store.ProgramID, mock_forwarder.ProgramID)
	require.NoError(t, err)
	storeAuthority, err := store.authority()
	require.NoError(t, err)
	feed, err := store.feed(dataID)
	require.NoError(t, err)
	// a non-SVR report passes only its feed, which holds its settings
	remaining := []*solana.AccountMeta{
		solana.Meta(store.State),
		solana.Meta(storeAuthority),
		solana.Meta(feed).WRITE(),
		solana.Meta(df1_store.ProgramID),
		solana.Meta(df1Feed).WRITE(),
	}

	// account hash over the receiver's accounts: forwarder state, forwarder authority, remaining
	hasher := sha256.New()
	hasher.Write(forwarderState.Bytes())
	hasher.Write(forwarderAuthority.Bytes())
	for _, m := range remaining {
		hasher.Write(m.PublicKey.Bytes())
	}

	// payload: Borsh Vec<ReceivedReport { data_id, timestamp: u32, answer: i128 }>
	payload := make([]byte, 4+36)
	binary.LittleEndian.PutUint32(payload, 1)
	copy(payload[4:], dataID[:])
	binary.LittleEndian.PutUint32(payload[20:], observedAt)
	v := toInt128(answer)
	binary.LittleEndian.PutUint64(payload[24:], v.Lo)
	binary.LittleEndian.PutUint64(payload[32:], v.Hi)

	// raw report: metadata (version | execution ID | timestamp | DON ID | config version |
	// workflow CID | workflow name | workflow owner | report ID) then the forwarder report
	metadata := make([]byte, 109)
	metadata[0] = 1
	binary.BigEndian.PutUint32(metadata[29:33], observedAt) // unique execution ID per report
	copy(metadata[77:87], storeTestWorkflowName[:])
	copy(metadata[87:107], storeTestWorkflowOwner[:])
	rawReport := append(metadata, hasher.Sum(nil)...)
	rawReport = binary.LittleEndian.AppendUint32(rawReport, uint32(len(payload)))
	rawReport = append(rawReport, payload...)

	// data: no signatures | raw report | report context
	data := append([]byte{0}, rawReport...)
	data = append(data, make([]byte, 96)...)

	ix := build(t)(mock_forwarder.NewReportInstruction(data, forwarderState, chain.DeployerKey.PublicKey(),
		forwarderAuthority, store.ProgramID, solana.SystemProgramID))
	generic := ix.(*solana.GenericInstruction)
	generic.AccountValues = append(generic.AccountValues, remaining...)
	require.NoError(t, chain.Confirm([]solana.Instruction{generic}, solanaUtils.AddComputeUnitLimit(fees.ComputeUnitLimit(1_400_000))))
}

func build(t *testing.T) func(solana.Instruction, error) solana.Instruction {
	return func(ix solana.Instruction, err error) solana.Instruction {
		t.Helper()
		require.NoError(t, err)
		return ix
	}
}

func confirm(t *testing.T, chain cldfsol.Chain, ix solana.Instruction, signers ...solana.PrivateKey) {
	t.Helper()
	require.NoError(t, chain.Confirm([]solana.Instruction{ix}, solanaUtils.AddSigners(signers...)))
}

func readAccount[T any](t *testing.T, chain cldfsol.Chain, address solana.PublicKey, parse func([]byte) (*T, error)) *T {
	t.Helper()
	info, err := getAccount(t.Context(), chain.Client, address)
	require.NoError(t, err)
	v, err := parse(info.Data.GetBinary())
	require.NoError(t, err)
	return v
}

// createAccount creates an empty account owned by the given program.
func createAccount(t *testing.T, chain cldfsol.Chain, owner solana.PublicKey) solana.PublicKey {
	t.Helper()
	key, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	rent, err := chain.Client.GetMinimumBalanceForRentExemption(t.Context(), 8, rpc.CommitmentConfirmed)
	require.NoError(t, err)
	confirm(t, chain, build(t)(system.NewCreateAccountInstruction(rent, 8, owner, chain.DeployerKey.PublicKey(), key.PublicKey()).ValidateAndBuild()), key)
	return key.PublicKey()
}

func copyProgram(t *testing.T, src, dst string) {
	t.Helper()
	in, err := os.Open(src)
	require.NoError(t, err)
	defer in.Close()
	out, err := os.Create(dst)
	require.NoError(t, err)
	defer out.Close()
	_, err = io.Copy(out, in)
	require.NoError(t, err)
}
