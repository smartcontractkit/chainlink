package solana

import (
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/gagliardetto/solana-go"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/helpers"

	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestDeployCache(t *testing.T) {
	skipInCI(t)
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:     1,
		SolChains: 1,
	}

	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	solSel := env.BlockChains.ListChainSelectors(cldfchain.WithFamily(chain_selectors.FamilySolana))[0]

	chain := env.BlockChains.SolanaChains()[solSel]
	chain.ProgramsPath = getProgramsPath()
	env.BlockChains = cldfchain.NewBlockChains(map[uint64]cldfchain.BlockChain{solSel: chain})

	forwarderProgramID := solana.SystemProgramID // needs to be executable
	t.Run("should deploy cache", func(t *testing.T) {
		configuredChangeset := commonchangeset.Configure(DeployCache{},
			&DeployCacheRequest{
				ChainSel:  solSel,
				Qualifier: testQualifier,
				Version:   "1.0.0",
				BuildConfig: &helpers.BuildSolanaConfig{
					GitCommitSha:   "3305b4d55b5469e110133e5a36e5600aadf436fb",
					DestinationDir: getProgramsPath(),
					LocalBuild:     helpers.LocalBuildConfig{BuildLocally: true, CreateDestinationDir: true},
				},
				FeedAdmins:         []solana.PublicKey{chain.DeployerKey.PublicKey()},
				ForwarderProgramID: forwarderProgramID,
			},
		)

		var err error
		env, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
		require.NoError(t, err)

		// Check that the cache program and state addresses are present in the datastore
		ds := env.DataStore
		version := "1.0.0"
		cacheKey := datastore.NewAddressRefKey(solSel, CacheContract, mustParseVersion(version), testQualifier)
		cacheStateKey := datastore.NewAddressRefKey(solSel, CacheState, mustParseVersion(version), testQualifier)

		cacheAddr, err := ds.Addresses().Get(cacheKey)
		require.NoError(t, err)
		require.NotEmpty(t, cacheAddr.Address)

		cacheStateAddr, err := ds.Addresses().Get(cacheStateKey)
		require.NoError(t, err)
		require.NotEmpty(t, cacheStateAddr.Address)
	})

	t.Run("should pass upgrade authority", func(t *testing.T) {
		configuredChangeset := commonchangeset.Configure(SetCacheUpgradeAuthority{},
			&SetCacheUpgradeAuthorityRequest{
				ChainSel:            solSel,
				Qualifier:           testQualifier,
				Version:             "1.0.0",
				NewUpgradeAuthority: chain.DeployerKey.PublicKey().String(),
			},
		)

		var err error
		_, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
		require.NoError(t, err)
	})
}

func TestConfigureCache(t *testing.T) {
	skipInCI(t)
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:     1,
		SolChains: 1,
	}

	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	solSel := env.BlockChains.ListChainSelectors(cldfchain.WithFamily(chain_selectors.FamilySolana))[0]

	chain := env.BlockChains.SolanaChains()[solSel]
	chain.ProgramsPath = getProgramsPath()
	env.BlockChains = cldfchain.NewBlockChains(map[uint64]cldfchain.BlockChain{solSel: chain})
	// Example array of DataIDs as [][16]uint8
	DataIDs := []string{
		"0x018e16c39e00032000000",
		"0x018e16c39e00032000001",
		"0x018e16c39e00032000002",
	}

	descriptions := [][32]uint8{
		[32]uint8{'B', 'i', 't', 'c', 'o', 'i', 'n', ' ', 'P', 'r', 'i', 'c', 'e', ' ', 'F', 'e', 'e', 'd'},
		[32]uint8{'E', 't', 'h', 'e', 'r', 'e', 'u', 'm', ' ', 'P', 'r', 'i', 'c', 'e', ' ', 'F', 'e', 'e', 'd'},
		[32]uint8{'S', 'o', 'l', 'a', 'n', 'a', ' ', 'P', 'r', 'i', 'c', 'e', ' ', 'F', 'e', 'e', 'd'},
	}

	// For AllowedSender (slice of solana.PublicKey)
	forwarderProgramID := []solana.PublicKey{
		solana.SystemProgramID, // should be executable
	}

	forwarderCacheID := []solana.PublicKey{
		solana.MustPublicKeyFromBase58("11111111111111111111111111111114"), // example public key
	}

	senderList := make([]Sender, len(forwarderProgramID))
	for i := range forwarderProgramID {
		senderList[i] = Sender{
			ProgramID: forwarderProgramID[i],
			StateID:   forwarderCacheID[i],
		}
	}

	// For AllowedWorkflowOwner (slice of [20]uint8 arrays)
	allowedWorkflowOwner := [][20]uint8{
		{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14},
	}

	// For AllowedWorkflowName (slice of [10]uint8 arrays)
	allowedWorkflowName := [][10]uint8{
		{0x74, 0x65, 0x73, 0x74, 0x5f, 0x77, 0x6f, 0x72, 0x6b, 0x00}, // "test_work" with null terminator
	}

	t.Run("should init cache decimal report without mcms", func(t *testing.T) {
		// First deploy the cache to get the program ID and state
		deployChangeset := commonchangeset.Configure(DeployCache{},
			&DeployCacheRequest{
				ChainSel:           solSel,
				Qualifier:          testQualifier,
				Version:            "1.0.0",
				FeedAdmins:         []solana.PublicKey{chain.DeployerKey.PublicKey()},
				ForwarderProgramID: forwarderProgramID[0],
			},
		)

		// Apply deploy changeset first to get the cache state and program ID
		_, _, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{deployChangeset})
		require.NoError(t, err)

		configuredChangeset := commonchangeset.Configure(InitCacheDecimalReport{},
			&InitCacheDecimalReportRequest{
				ChainSel:  solSel,
				Qualifier: testQualifier,
				Version:   "1.0.0",
				DataIDs:   DataIDs,
				FeedAdmin: chain.DeployerKey.PublicKey(),
			},
		)

		// Apply the init changeset
		out, _, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{deployChangeset, configuredChangeset})
		require.NoError(t, err)

		configuredChangeset = commonchangeset.Configure(ConfigureCacheDecimalReport{},
			&ConfigureCacheDecimalReportRequest{
				ChainSel:             solSel,
				Qualifier:            testQualifier,
				Version:              "1.0.0",
				SenderList:           senderList,
				AllowedWorkflowOwner: allowedWorkflowOwner,
				AllowedWorkflowName:  allowedWorkflowName,
				FeedAdmin:            chain.DeployerKey.PublicKey(),
				DataIDs:              DataIDs,
				Descriptions:         descriptions,
			},
		)

		// Apply the configure changeset
		_, _, err = commonchangeset.ApplyChangesets(t, out, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
		require.NoError(t, err)
	})

	t.Run("should set cache decimal report config without mcms", func(t *testing.T) {
		// First deploy the cache
		deployChangeset := commonchangeset.Configure(DeployCache{},
			&DeployCacheRequest{
				ChainSel:           solSel,
				Qualifier:          testQualifier,
				Version:            "1.0.0",
				FeedAdmins:         []solana.PublicKey{chain.DeployerKey.PublicKey()},
				ForwarderProgramID: forwarderProgramID[0],
			},
		)

		// Apply deploy changeset first to get the cache state and program ID
		out, _, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{deployChangeset})
		require.NoError(t, err)

		configuredChangeset := commonchangeset.Configure(ConfigureCacheDecimalReport{},
			&ConfigureCacheDecimalReportRequest{
				ChainSel:             solSel,
				Qualifier:            testQualifier,
				Version:              "1.0.0",
				SenderList:           senderList,
				AllowedWorkflowOwner: allowedWorkflowOwner,
				AllowedWorkflowName:  allowedWorkflowName,
				FeedAdmin:            chain.DeployerKey.PublicKey(),
				DataIDs:              DataIDs,
				Descriptions:         descriptions,
			},
		)

		// Apply the configure changeset
		_, _, err = commonchangeset.ApplyChangesets(t, out, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
		require.NoError(t, err)
	})

	t.Run("should set cache decimal report config with mcms", func(t *testing.T) {
		configuredChangeset := commonchangeset.Configure(ConfigureCacheDecimalReport{},
			&ConfigureCacheDecimalReportRequest{
				ChainSel:             solSel,
				Qualifier:            testQualifier,
				Version:              "1.0.0",
				SenderList:           senderList,
				AllowedWorkflowOwner: allowedWorkflowOwner,
				AllowedWorkflowName:  allowedWorkflowName,
				FeedAdmin:            chain.DeployerKey.PublicKey(),
				DataIDs:              DataIDs,
				Descriptions:         descriptions,
			},
		)

		deployChangeset := commonchangeset.Configure(DeployCache{},
			&DeployCacheRequest{
				ChainSel:           solSel,
				Qualifier:          testQualifier,
				Version:            "1.0.0",
				FeedAdmins:         []solana.PublicKey{chain.DeployerKey.PublicKey()},
				ForwarderProgramID: forwarderProgramID[0],
			},
		)

		ds := datastore.NewMemoryDataStore()

		// deploy mcms
		mcmsState, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolanaV2(env, ds, chain,
			commontypes.MCMSWithTimelockConfigV2{
				Canceller:        proposalutils.SingleGroupMCMSV2(t),
				Proposer:         proposalutils.SingleGroupMCMSV2(t),
				Bypasser:         proposalutils.SingleGroupMCMSV2(t),
				TimelockMinDelay: big.NewInt(0),
			},
		)
		require.NoError(t, err)

		ds.Seal()
		fundSignerPDAs(t, env, solSel, mcmsState)

		transferOwnershipChangeset := commonchangeset.Configure(TransferOwnershipCache{},
			&TransferOwnershipCacheRequest{
				ChainSel:  solSel,
				MCMSCfg:   proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
				Qualifier: testQualifier,
				Version:   "1.0.0",
			})

		_, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{deployChangeset, configuredChangeset, transferOwnershipChangeset})
		require.NoError(t, err)
	})
}

func ParseSemver(v string) *semver.Version {
	ver, err := semver.NewVersion(v)
	if err != nil {
		panic(err)
	}
	return ver
}

func mustParseVersion(v string) *semver.Version {
	return ParseSemver(v)
}

func getProgramsPath() string {
	// Get the directory of the current file (environment.go)
	_, currentFile, _, _ := runtime.Caller(0)
	// Go up to the root of the deployment package
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	// Construct the absolute path
	return filepath.Join(rootDir, "changeset/solana", "solana_contracts")
}

func skipInCI(t *testing.T) {
	ci := os.Getenv("CI") == "true"
	if ci {
		t.Skip("Skipping in CI")
	}
}

func fundSignerPDAs(
	t *testing.T, env cldf.Environment, chainSelector uint64, chainState *state.MCMSWithTimelockStateSolana,
) {
	t.Helper()
	solChain := env.BlockChains.SolanaChains()[chainSelector]
	timelockSignerPDA := state.GetTimelockSignerPDA(chainState.TimelockProgram, chainState.TimelockSeed)
	mcmSignerPDA := state.GetMCMSignerPDA(chainState.McmProgram, chainState.ProposerMcmSeed)
	signerPDAs := []solana.PublicKey{timelockSignerPDA, mcmSignerPDA}
	err := memory.FundSolanaAccounts(env.GetContext(), signerPDAs, 1, solChain.Client)
	require.NoError(t, err)
}

const (
	testQualifier = "test-deploy"
)
