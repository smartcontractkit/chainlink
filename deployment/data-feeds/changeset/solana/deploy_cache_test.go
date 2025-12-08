package solana

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
)

const (
	testQualifier = "test-deploy"
)

// TODO: This test is not working as expected, however this has been modified to work with the new
// test engine environment.
func TestDeployCache(t *testing.T) {
	skipInCI(t)
	t.Parallel()

	selector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	e, err := environment.New(t.Context(),
		environment.WithSolanaContainer(t, []uint64{selector}, t.TempDir(), map[string]string{}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	chain := e.BlockChains.SolanaChains()[selector]

	// Reassign the environment value to keep the test the same as the original test
	env := *e

	forwarderProgramID := solana.SystemProgramID // needs to be executable
	t.Run("should deploy cache", func(t *testing.T) {
		configuredChangeset := commonchangeset.Configure(DeployCache{},
			&DeployCacheRequest{
				ChainSel:  selector,
				Qualifier: testQualifier,
				Version:   "1.0.0",
				BuildConfig: &helpers.BuildSolanaConfig{
					GitCommitSha:   "cd449e02f649fab782739685b57b373394e6f3e8",
					DestinationDir: chain.ProgramsPath,
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
		cacheKey := datastore.NewAddressRefKey(selector, CacheContract, semver.MustParse(version), testQualifier)
		cacheStateKey := datastore.NewAddressRefKey(selector, CacheState, semver.MustParse(version), testQualifier)

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
				ChainSel:            selector,
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

// TODO: This test is not working as expected, however this has been modified to work with the new
// test engine environment.
func TestConfigureCache(t *testing.T) {
	skipInCI(t)
	t.Parallel()

	selector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	programsPath := t.TempDir()
	soltestutils.LoadDataFeedsPrograms(t, programsPath)

	e, err := environment.New(t.Context(),
		environment.WithSolanaContainer(t, []uint64{selector}, programsPath, map[string]string{}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	chain := e.BlockChains.SolanaChains()[selector]

	files, err := os.ReadDir(programsPath)
	require.NoError(t, err)
	for _, file := range files {
		absPath := filepath.Join(programsPath, file.Name())
		t.Logf("Program file: %s", absPath)
	}

	env := *e

	// Example array of DataIDs as [][16]uint8
	DataIDs := []string{
		"0x018e16c39e00032000000",
		"0x018e16c39e00032000001",
		"0x018e16c39e00032000002",
	}

	descriptions := [][32]uint8{
		{'B', 'i', 't', 'c', 'o', 'i', 'n', ' ', 'P', 'r', 'i', 'c', 'e', ' ', 'F', 'e', 'e', 'd'},
		{'E', 't', 'h', 'e', 'r', 'e', 'u', 'm', ' ', 'P', 'r', 'i', 'c', 'e', ' ', 'F', 'e', 'e', 'd'},
		{'S', 'o', 'l', 'a', 'n', 'a', ' ', 'P', 'r', 'i', 'c', 'e', ' ', 'F', 'e', 'e', 'd'},
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
				ChainSel:           selector,
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
				ChainSel:  selector,
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
				ChainSel:             selector,
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
				ChainSel:           selector,
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
				ChainSel:             selector,
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
				ChainSel:             selector,
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
				ChainSel:           selector,
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

		soltestutils.FundSignerPDAs(t, chain, mcmsState)

		transferOwnershipChangeset := commonchangeset.Configure(TransferOwnershipCache{},
			&TransferOwnershipCacheRequest{
				ChainSel:  selector,
				MCMSCfg:   proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
				Qualifier: testQualifier,
				Version:   "1.0.0",
			})

		_, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{deployChangeset, configuredChangeset, transferOwnershipChangeset})
		require.NoError(t, err)
	})
}

func skipInCI(t *testing.T) {
	ci := os.Getenv("CI") == "true"
	if ci {
		t.Skip("Skipping in CI")
	}
}
