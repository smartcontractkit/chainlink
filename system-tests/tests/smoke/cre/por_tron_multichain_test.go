package cre

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	computecap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/compute"
	consensuscap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/consensus"
	croncap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/cron"
	webapicap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/webapi"
	writeevmcap "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/writeevm"
)

func TestCRE_OCR3_PoR_Workflow_MultiChain_Tron_MockedPrice(t *testing.T) {
	configErr := setConfigurationIfMissing("environment-multichain-tron.toml")
	require.NoError(t, configErr, "failed to set CTF config")
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t)
	require.Len(t, in.NodeSets, 1, "expected 1 node set in the test config")
	require.Len(t, in.Blockchains, 2, "expected 2 blockchains in the test config (Anvil + Tron)")

	// Start Tron node for the secondary chain
	privateKeyHex := DefaultTronPrivateKey
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err)

	genesisAddress := address.PubkeyToAddress(privateKey.PublicKey)
	testLogger.Info().Str("genesis address", genesisAddress.String()).Msg("Using genesis account")

	// err = StartTronNode(genesisAddress.String())
	// require.NoError(t, err, "failed to start Tron node")
	// testLogger.Info().Msg("Tron node started")
	// defer ShutdownTronBackend(t)

	// Set Tron private key environment variable
	err = SetTronPrivateKeyEnv()
	require.NoError(t, err, "failed to set Tron private key")

	testLogger.Info().Msg("Using multi-chain configuration (Anvil + Tron)")
	testLogger.Info().Msgf("Blockchain count: %d", len(in.Blockchains))

	var anvilChain, tronChain *cre.WrappedBlockchainInput

	for i, bc := range in.Blockchains {
		testLogger.Info().Msgf("Blockchain %d: ChainID=%s, Type=%s", i, bc.ChainID, bc.Type)

		if bc.Type == "anvil" {
			anvilChain = bc
			testLogger.Info().Msg("Found Anvil chain - will be used as home chain for Keystone contracts")
		} else if bc.Type == "tron" {
			tronChain = bc
			testLogger.Info().Msg("Found Tron chain - will be used as secondary chain for write target")
		}
	}

	require.NotNil(t, anvilChain, "Anvil chain not found in configuration")
	require.NotNil(t, tronChain, "Tron chain not found in configuration")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*cre.CapabilitiesAwareNodeSet {
		return []*cre.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       SinglePoRDonCapabilitiesFlags,
				DONTypes:           []string{cre.WorkflowDON, cre.GatewayDON},
				BootstrapNodeIndex: 0,
				GatewayNodeIndex:   0,
			},
		}
	}

	feedIDs := make([]string, 0, len(in.WorkflowConfigs))
	for _, wc := range in.WorkflowConfigs {
		feedIDs = append(feedIDs, wc.FeedID)
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake, AuthorizationKey, feedIDs)
	require.NoError(t, priceErr, "failed to create fake price provider")

	// Convert Tron chain ID to uint64 for capability factory
	tronChainIDInt, chainErr := strconv.Atoi(tronChain.ChainID)
	require.NoError(t, chainErr, "failed to convert Tron chain ID to int")
	tronChainIDUint64 := libc.MustSafeUint64(int64(tronChainIDInt))

	capabilityFactoryFns := []cre.DONCapabilityWithConfigFactoryFn{
		webapicap.WebAPITriggerCapabilityFactoryFn,
		webapicap.WebAPITargetCapabilityFactoryFn,
		computecap.ComputeCapabilityFactoryFn,
		consensuscap.OCR3CapabilityFactoryFn,
		croncap.CronCapabilityFactoryFn,
		writeevmcap.WriteEVMCapabilityFactory(tronChainIDUint64), // Use Tron chain ID for write target
	}

	// Use the existing setupPoRTestEnvironment with the multi-chain config
	setupOutput := setupPoRTestEnvironment(
		t,
		testLogger,
		in,
		priceProvider,
		mustSetCapabilitiesFn,
		capabilityFactoryFns,
	)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			testLogger.Info().Msg("Test failed - logging multi-chain debug information...")
			testLogger.Info().Msgf("Anvil Chain ID: %s", anvilChain.ChainID)
			testLogger.Info().Msgf("Tron Chain ID: %s", tronChain.ChainID)
			testLogger.Info().Msgf("Anvil Type: %s", anvilChain.Type)
			testLogger.Info().Msgf("Tron Type: %s", tronChain.Type)
		}
		debugTest(t, testLogger, setupOutput, in)
	})

	testLogger.Info().Msg("Starting multi-chain PoR workflow test - waiting for feed updates...")
	testLogger.Info().Msgf("Home chain (Anvil): %s", anvilChain.ChainID)
	testLogger.Info().Msgf("Secondary chain (Tron): %s", tronChain.ChainID)

	waitForFeedUpdate(t, testLogger, priceProvider, setupOutput, 5*time.Minute)
	testLogger.Info().Msg("✅ Multi-chain PoR workflow test completed successfully!")
}
