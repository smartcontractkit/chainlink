package cre

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"
	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libcontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	crecompute "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/compute"
	creconsensus "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
	crecron "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/cron"
	cregateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	libcrecli "github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

type multiPorSetupOutput struct {
	priceProvider           PriceProvider
	dataFeedsCacheAddresses []common.Address
	forwarderAddresses      []common.Address
	sethClients             []*seth.Client
	blockchainsOutput       []*blockchain.Output
	donTopology             *keystonetypes.DonTopology
	nodeOutput              []*keystonetypes.WrappedNodeOutput
}

func setupPoRMultiChainTestEnvironment(
	t *testing.T,
	testLogger zerolog.Logger,
	in *TestConfig,
	priceProvider PriceProvider,
	mustSetCapabilitiesFn func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet,
	capabilityFactoryFns []func([]string) []keystone_changeset.DONCapabilityWithConfig,
) *multiPorSetupOutput {
	extraAllowedPorts := []int{}
	if _, ok := priceProvider.(*FakePriceProvider); ok {
		extraAllowedPorts = append(extraAllowedPorts, in.Fake.Port)
	}

	customBinariesPaths := map[string]string{}
	containerPath, pathErr := capabilities.DefaultContainerDirectory(in.Infra.InfraType)
	require.NoError(t, pathErr, "failed to get default container directory")
	var cronBinaryPathInTheContainer string
	if in.WorkflowConfig.DependenciesConfig.CronCapabilityBinaryPath != "" {
		// where cron binary is located in the container
		cronBinaryPathInTheContainer = filepath.Join(containerPath, filepath.Base(in.WorkflowConfig.DependenciesConfig.CronCapabilityBinaryPath))
		// where cron binary is located on the host
		customBinariesPaths[keystonetypes.CronCapability] = in.WorkflowConfig.DependenciesConfig.CronCapabilityBinaryPath
	} else {
		// assume that if cron binary is already in the image it is in the default location and has default name
		cronBinaryPathInTheContainer = filepath.Join(containerPath, "cron")
	}

	firstBlockchain := in.Blockchains[0]

	chainIDInt, err := strconv.Atoi(firstBlockchain.ChainID)
	require.NoError(t, err, "failed to convert chain ID to int")
	chainIDUint64 := libc.MustSafeUint64(int64(chainIDInt))

	universalSetupInput := creenv.SetupInput{
		CapabilitiesAwareNodeSets:            mustSetCapabilitiesFn(in.NodeSets),
		CapabilitiesContractFactoryFunctions: capabilityFactoryFns,
		BlockchainsInput:                     in.Blockchains,
		JdInput:                              *in.JD,
		InfraInput:                           *in.Infra,
		CustomBinariesPaths:                  customBinariesPaths,
		ExtraAllowedPorts:                    extraAllowedPorts,
		JobSpecFactoryFunctions: []keystonetypes.JobSpecFactoryFn{
			creconsensus.ConsensusJobSpecFactoryFn(chainIDUint64),
			crecron.CronJobSpecFactoryFn(cronBinaryPathInTheContainer),
			cregateway.GatewayJobSpecFactoryFn(extraAllowedPorts, []string{}, []string{"0.0.0.0/0"}),
			crecompute.ComputeJobSpecFactoryFn,
		},
	}

	universalSetupOutput, setupErr := creenv.SetupTestEnvironment(testcontext.Get(t), testLogger, cldlogger.NewSingleFileLogger(t), universalSetupInput)
	require.NoError(t, setupErr, "failed to setup test environment")
	homeChainOutput := universalSetupOutput.BlockchainOutput[0]

	if in.CustomAnvilMiner != nil {
		for _, bi := range universalSetupInput.BlockchainsInput {
			if bi.Type == blockchain.TypeAnvil {
				require.NotContains(t, bi.DockerCmdParamsOverrides, "-b", "custom_anvil_miner was specified but Anvil has '-b' key set, remove that parameter from 'docker_cmd_params' to run deployments instantly or remove custom_anvil_miner key from TOML config")
			}
		}
		for _, bo := range universalSetupOutput.BlockchainOutput {
			if bo.BlockchainOutput.Type == blockchain.TypeAnvil {
				miner := rpc.NewRemoteAnvilMiner(bo.BlockchainOutput.Nodes[0].ExternalHTTPUrl, nil)
				miner.MinePeriodically(time.Duration(in.CustomAnvilMiner.BlockSpeedSeconds) * time.Second)
			}
		}
	}

	var dataFeedsCacheAddresses []common.Address

	for idx, bo := range universalSetupOutput.BlockchainOutput {
		workflowName := in.WorkflowConfig.WorkflowName + "-" + fmt.Sprint(bo.ChainID)
		deployDataFeedsInput := &keystonetypes.DeployDataFeedsCacheInput{
			ChainSelector: bo.ChainSelector,
			CldEnv:        universalSetupOutput.CldEnvironment,
		}
		deployDataFeedsCacheOutput, dfErr := libcontracts.DeployDataFeedsCache(testLogger, deployDataFeedsInput)
		require.NoError(t, dfErr, "failed to deploy data feeds cache")

		dataFeedsCacheAddresses = append(dataFeedsCacheAddresses, deployDataFeedsCacheOutput.DataFeedsCacheAddress)
		var creCLIAbsPath string
		var creCLISettingsFile *os.File
		if in.WorkflowConfig.UseCRECLI {
			// make sure that path is indeed absolute
			var pathErr error
			creCLIAbsPath, pathErr = filepath.Abs(in.WorkflowConfig.DependenciesConfig.CRECLIBinaryPath)
			require.NoError(t, pathErr, "failed to get absolute path for CRE CLI")

			// create CRE CLI settings file
			var settingsErr error
			creCLISettingsFile, settingsErr = libcrecli.PrepareCRECLISettingsFile(
				bo.SethClient.MustGetRootKeyAddress(),
				universalSetupOutput.CldEnvironment.ExistingAddresses,
				universalSetupOutput.DonTopology.WorkflowDonID,
				homeChainOutput.ChainSelector,
				map[uint64]string{
					homeChainOutput.ChainSelector: homeChainOutput.BlockchainOutput.Nodes[0].ExternalHTTPUrl,
					bo.ChainSelector:              bo.BlockchainOutput.Nodes[0].ExternalHTTPUrl,
				},
			)
			require.NoError(t, settingsErr, "failed to create CRE CLI settings file")
		}

		chainAddrBook, addrErr := universalSetupOutput.CldEnvironment.ExistingAddresses.AddressesForChain(bo.ChainSelector)
		require.NoError(t, addrErr, "failed to get existing addresses for chain %d", bo.ChainSelector)

		var forwarderAddr common.Address
		for addrStr, tv := range chainAddrBook {
			if strings.Contains(tv.String(), "KeystoneForwarder") {
				forwarderAddr = common.HexToAddress(addrStr)
				break
			}
		}

		dfConfigInput := &configureDataFeedsCacheInput{
			useCRECLI:             in.WorkflowConfig.UseCRECLI,
			chainSelector:         bo.ChainSelector,
			fullCldEnvironment:    universalSetupOutput.CldEnvironment,
			forwarderAddress:      forwarderAddr,
			dataFeedsCacheAddress: deployDataFeedsCacheOutput.DataFeedsCacheAddress,
			workflowName:          workflowName,
			feedID:                in.WorkflowConfig.FeedIDs[idx],
			sethClient:            bo.SethClient,
			blockchain:            bo.BlockchainOutput,
			creCLIAbsPath:         creCLIAbsPath,
			settingsFile:          creCLISettingsFile,
			deployerPrivateKey:    bo.DeployerPrivateKey,
		}
		dfConfigErr := configureDataFeedsCacheContract(testLogger, dfConfigInput)
		require.NoError(t, dfConfigErr, "failed to configure data feeds cache")

		registerInput := registerPoRWorkflowInput{
			WorkflowConfig:          in.WorkflowConfig,
			workflowName:            workflowName,
			chainSelector:           bo.ChainSelector,
			workflowDonID:           universalSetupOutput.DonTopology.WorkflowDonID,
			feedID:                  in.WorkflowConfig.FeedIDs[idx],
			workflowRegistryAddress: universalSetupOutput.KeystoneContractsOutput.WorkflowRegistryAddress,
			dataFeedsCacheAddress:   deployDataFeedsCacheOutput.DataFeedsCacheAddress,
			priceProvider:           priceProvider,
			sethClient:              bo.SethClient,
			deployerPrivateKey:      bo.DeployerPrivateKey,
			creCLIAbsPath:           creCLIAbsPath,
			creCLIsettingsFile:      creCLISettingsFile,
			writeTargetName:         corevm.GenerateWriteTargetName(bo.ChainID),
		}

		workflowErr := registerPoRWorkflow(registerInput)
		require.NoError(t, workflowErr, "failed to register PoR workflow")
	}
	// Workflow-specific configuration -- END

	// Set inputs in the test config, so that they can be saved
	//TODO this should be a map
	// in.KeystoneContracts = &keystonetypes.KeystoneContractsInput{
	// 	Out: universalSetupOutput.KeystoneContractsOutput,
	// }
	// in.DataFeedsCacheContract = &keystonetypes.DeployDataFeedsCacheInput{
	// 	Out: &keystonetypes.DeployDataFeedsCacheOutput{
	// 		DataFeedsCacheAddress: deployDataFeedsCacheOutput.DataFeedsCacheAddress,
	// 	},
	// }
	// in.WorkflowRegistryConfiguration = &keystonetypes.WorkflowRegistryInput{
	// 	Out: universalSetupOutput.WorkflowRegistryConfigurationOutput,
	// }

	ret := &multiPorSetupOutput{
		priceProvider:           priceProvider,
		dataFeedsCacheAddresses: dataFeedsCacheAddresses,
		donTopology:             universalSetupOutput.DonTopology,
		nodeOutput:              universalSetupOutput.NodeOutput,
	}

	for _, bo := range universalSetupOutput.BlockchainOutput {
		ret.forwarderAddresses = append(ret.forwarderAddresses, universalSetupOutput.KeystoneContractsOutput.ForwarderAddress)
		ret.sethClients = append(ret.sethClients, bo.SethClient)
		ret.blockchainsOutput = append(ret.blockchainsOutput, bo.BlockchainOutput)
	}

	return ret
}

// config file to use: environment-multichain-one-don.toml
func TestCRE_OCR3_PoR_Workflow_SingleDon_MultipleWriters_MockedPrice(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 1, "expected 1 node set in the test config")

	// Assign all capabilities to the single node set
	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       SinglePoRDonCapabilitiesFlags,
				DONTypes:           []string{keystonetypes.WorkflowDON, keystonetypes.GatewayDON},
				BootstrapNodeIndex: 0, // not required, but set to make the configuration explicit
				GatewayNodeIndex:   0, // not required, but set to make the configuration explicit
			},
		}
	}

	priceProvider, priceErr := NewFakePriceProvider(testLogger, in.Fake, AuthorizationKey, in.WorkflowConfig.FeedIDs)
	require.NoError(t, priceErr, "failed to create fake price provider")

	homeChain := in.Blockchains[0]
	targetChain := in.Blockchains[1]
	homeChainID, chainErr := strconv.Atoi(homeChain.ChainID)
	require.NoError(t, chainErr, "failed to convert home chain ID to int")
	targetChainID, chainErr := strconv.Atoi(targetChain.ChainID)
	require.NoError(t, chainErr, "failed to convert target chain ID to int")

	setupOutput := setupPoRMultiChainTestEnvironment(
		t,
		testLogger,
		in,
		priceProvider,
		mustSetCapabilitiesFn,
		[]keystonetypes.DONCapabilityWithConfigFactoryFn{
			libcontracts.DefaultCapabilityFactoryFn,
			libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(homeChainID))),
			libcontracts.ChainWriterCapabilityFactory(libc.MustSafeUint64(int64(targetChainID))),
		},
	)

	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			for idx, feedID := range in.WorkflowConfig.FeedIDs {
				logTestInfo(testLogger, feedID, in.WorkflowConfig.WorkflowName, setupOutput.dataFeedsCacheAddresses[idx].Hex(), setupOutput.forwarderAddresses[idx].Hex())

				// log scanning is not supported for CRIB
				if in.Infra.InfraType == libtypes.CRIB {
					return
				}

				logDir := fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name())

				removeErr := os.RemoveAll(logDir)
				if removeErr != nil {
					testLogger.Error().Err(removeErr).Msg("failed to remove log directory")
					return
				}

				_, saveErr := framework.SaveContainerLogs(logDir)
				if saveErr != nil {
					testLogger.Error().Err(saveErr).Msg("failed to save container logs")
					return
				}

				debugDons := make([]*keystonetypes.DebugDon, 0, len(setupOutput.donTopology.DonsWithMetadata))
				for i, donWithMetadata := range setupOutput.donTopology.DonsWithMetadata {
					containerNames := make([]string, 0, len(donWithMetadata.NodesMetadata))
					for _, output := range setupOutput.nodeOutput[i].Output.CLNodes {
						containerNames = append(containerNames, output.Node.ContainerName)
					}
					debugDons = append(debugDons, &keystonetypes.DebugDon{
						NodesMetadata:  donWithMetadata.NodesMetadata,
						Flags:          donWithMetadata.Flags,
						ContainerNames: containerNames,
					})
				}

				debugInput := keystonetypes.DebugInput{
					DebugDons:        debugDons,
					BlockchainOutput: setupOutput.blockchainsOutput[idx],
					InfraInput:       in.Infra,
				}
				lidebug.PrintTestDebug(t.Name(), testLogger, debugInput)
			}
		}
	})

	for idx, feedID := range in.WorkflowConfig.FeedIDs {
		testLogger.Info().Msgf("Waiting for feed %s to update...", feedID)
		timeout := 5 * time.Minute // It can take a while before the first report is produced, particularly on CI.

		dataFeedsCacheInstance, instanceErr := data_feeds_cache.NewDataFeedsCache(setupOutput.dataFeedsCacheAddresses[idx], setupOutput.sethClients[idx].Client)
		require.NoError(t, instanceErr, "failed to create data feeds cache instance")

		startTime := time.Now()
		assert.Eventually(t, func() bool {
			elapsed := time.Since(startTime).Round(time.Second)
			price, err := dataFeedsCacheInstance.GetLatestAnswer(setupOutput.sethClients[idx].NewCallOpts(), [16]byte(common.Hex2Bytes(feedID)))
			require.NoError(t, err, "failed to get price from Data Feeds Cache contract")

			// if there are no more prices to be found, we can stop waiting
			return !setupOutput.priceProvider.NextPrice(feedID, price, elapsed)
		}, timeout, 10*time.Second, "feed %s did not update, timeout after: %s", feedID, timeout)

		require.EqualValues(t, priceProvider.ExpectedPrices(feedID), priceProvider.ActualPrices(feedID), "prices do not match")
		testLogger.Info().Msgf("All %d prices were found in the feed %s", len(priceProvider.ExpectedPrices(feedID)), feedID)
	}
}
