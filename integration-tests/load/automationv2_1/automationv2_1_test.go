package automationv2_1

import (
	"context"
	"encoding/json"
	"fmt"
	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/slack-go/slack"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	contractseth "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
	registrar21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registrar_wrapper2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/simple_log_upkeep_counter_wrapper"
	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2keepers30config "github.com/smartcontractkit/ocr2keepers/pkg/v3/config"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	StartupWaitTime = 30 * time.Second
	StopWaitTime    = 60 * time.Second
)

var (
	baseTOML = `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
Enabled = true
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`

	minimumNodeSpec = map[string]interface{}{
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "4Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "4Gi",
			},
		},
	}

	minimumDbSpec = map[string]interface{}{
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "1000m",
				"memory": "1Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "1000m",
				"memory": "1Gi",
			},
		},
		"stateful": true,
		"capacity": "5Gi",
	}

	recNodeSpec = map[string]interface{}{
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "4000m",
				"memory": "8Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "4000m",
				"memory": "8Gi",
			},
		},
	}

	recDbSpec = map[string]interface{}{
		"resources": map[string]interface{}{
			"requests": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "2000m",
				"memory": "2Gi",
			},
		},
		"stateful": true,
		"capacity": "10Gi",
	}
)

var (
	numberofNodes, _   = strconv.Atoi(getEnv("NUMBEROFNODES", "6"))      // Number of nodes in the DON
	numberOfUpkeeps, _ = strconv.Atoi(getEnv("NUMBEROFUPKEEPS", "100"))  // Number of log triggered upkeeps
	duration, _        = strconv.Atoi(getEnv("DURATION", "900"))         // Test duration in seconds
	blockTime, _       = strconv.Atoi(getEnv("BLOCKTIME", "1"))          // Block time in seconds for geth simulated dev network
	numberOfEvents, _  = strconv.Atoi(getEnv("NUMBEROFEVENTS", "1"))     // Number of events to emit per trigger
	specType           = getEnv("SPECTYPE", "minimum")                   // minimum, recommended, local specs for the test
	logLevel           = getEnv("LOGLEVEL", "info")                      // log level for the chainlink nodes
	pyroscope, _       = strconv.ParseBool(getEnv("PYROSCOPE", "false")) // enable pyroscope for the chainlink nodes
)

func TestLogTrigger(t *testing.T) {
	l := logging.GetTestLogger(t)

	l.Info().Msg("Starting automation v2.1 log trigger load test")
	l.Info().Str("TEST_INPUTS", os.Getenv("TEST_INPUTS")).Int("Number of Nodes", numberofNodes).
		Int("Number of Upkeeps", numberOfUpkeeps).
		Int("Duration", duration).
		Int("Block Time", blockTime).
		Int("Number of Events", numberOfEvents).
		Str("Spec Type", specType).
		Str("Log Level", logLevel).
		Str("Image", os.Getenv(config.EnvVarCLImage)).
		Str("Tag", os.Getenv(config.EnvVarCLTag)).
		Msg("Test Config")

	testConfig := fmt.Sprintf("Number of Nodes: %d\nNumber of Upkeeps: %d\nDuration: %d\nBlock Time: %d\n"+
		"Number of Events: %d\nSpec Type: %s\nLog Level: %s\nImage: %s\nTag: %s\n", numberofNodes, numberOfUpkeeps, duration,
		blockTime, numberOfEvents, specType, logLevel, os.Getenv(config.EnvVarCLImage), os.Getenv(config.EnvVarCLTag))

	testNetwork := networks.MustGetSelectedNetworksFromEnv()[0]
	testType := "load"
	loadDuration := time.Duration(duration) * time.Second
	automationDefaultLinkFunds := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(10000))) //10000 LINK
	automationDefaultUpkeepGasLimit := uint32(1_000_000)

	registrySettings := &contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    uint32(0),
		FlatFeeMicroLINK:     uint32(40_000),
		BlockCountPerTurn:    big.NewInt(100),
		CheckGasLimit:        uint32(45_000_000), //45M
		StalenessSeconds:     big.NewInt(90_000),
		GasCeilingMultiplier: uint16(2),
		MaxPerformGas:        uint32(5_000_000),
		MinUpkeepSpend:       big.NewInt(0),
		FallbackGasPrice:     big.NewInt(2e11),
		FallbackLinkPrice:    big.NewInt(2e18),
		MaxCheckDataSize:     uint32(5_000),
		MaxPerformDataSize:   uint32(5_000),
		RegistryVersion:      contractseth.RegistryVersion_2_1,
	}

	testEnvironment := environment.New(&environment.Config{
		TTL: time.Hour * 24, // 1 day,
		NamespacePrefix: fmt.Sprintf(
			"automation-%s-%s",
			testType,
			strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-"),
		),
		Test:               t,
		PreventPodEviction: true,
	})

	if testEnvironment.WillUseRemoteRunner() {
		key := "TEST_INPUTS"
		err := os.Setenv(fmt.Sprintf("TEST_%s", key), os.Getenv(key))
		require.NoError(t, err, "failed to set the environment variable TEST_INPUTS for remote runner")

		key = config.EnvVarPyroscopeServer
		err = os.Setenv(fmt.Sprintf("TEST_%s", key), os.Getenv(key))
		require.NoError(t, err, "failed to set the environment variable PYROSCOPE_SERVER for remote runner")

		key = config.EnvVarPyroscopeKey
		err = os.Setenv(fmt.Sprintf("TEST_%s", key), os.Getenv(key))
		require.NoError(t, err, "failed to set the environment variable PYROSCOPE_KEY for remote runner")
	}

	testEnvironment.
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
			Values: map[string]interface{}{
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "4000m",
						"memory": "4Gi",
					},
					"limits": map[string]interface{}{
						"cpu":    "8000m",
						"memory": "8Gi",
					},
				},
				"geth": map[string]interface{}{
					"blocktime": blockTime,
					"capacity":  "10Gi",
				},
			},
		}))

	err := testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")

	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	var (
		nodeSpec = minimumNodeSpec
		dbSpec   = minimumDbSpec
	)

	switch specType {
	case "recommended":
		nodeSpec = recNodeSpec
		dbSpec = recDbSpec
	case "local":
		nodeSpec = map[string]interface{}{}
		dbSpec = map[string]interface{}{"stateful": true}
	default:
		// minimum:

	}

	if !pyroscope {
		err = os.Setenv(config.EnvVarPyroscopeServer, "")
		require.NoError(t, err, "Error setting pyroscope server env var")
	}

	err = os.Setenv(config.EnvVarPyroscopeEnvironment, testEnvironment.Cfg.Namespace)
	require.NoError(t, err, "Error setting pyroscope environment env var")

	for i := 0; i < numberofNodes+1; i++ { // +1 for the OCR boot node
		nodeTOML := baseTOML
		if i == 1 || i == 3 {
			nodeTOML = fmt.Sprintf("%s\n\n[Log]\nLevel = \"%s\"", baseTOML, logLevel)
		} else {
			nodeTOML = fmt.Sprintf("%s\n\n[Log]\nLevel = \"info\"", baseTOML)
		}
		nodeTOML = client.AddNetworksConfig(nodeTOML, testNetwork)
		testEnvironment.AddHelm(chainlink.New(i, map[string]any{
			"toml":      nodeTOML,
			"chainlink": nodeSpec,
			"db":        dbSpec,
		}))
	}

	err = testEnvironment.Run()
	require.NoError(t, err, "Error running chainlink DON")

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment, l)
	require.NoError(t, err, "Error building chain client")

	contractDeployer, err := contracts.NewContractDeployer(chainClient, l)
	require.NoError(t, err, "Error building contract deployer")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to chainlink nodes")

	chainClient.ParallelTransactions(true)

	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Error deploying link token contract")

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Failed waiting for contracts to deploy")

	registry, registrar := actions.DeployAutoOCRRegistryAndRegistrar(
		t, contractseth.RegistryVersion_2_1, *registrySettings, linkToken, contractDeployer, chainClient,
	)

	err = actions.FundChainlinkNodesAddress(chainlinkNodes[1:], chainClient, big.NewFloat(100), 0)
	require.NoError(t, err, "Error funding chainlink nodes")

	actions.CreateOCRKeeperJobs(
		t,
		chainlinkNodes,
		registry.Address(),
		chainClient.GetChainID().Int64(),
		0,
		contractseth.RegistryVersion_2_1,
	)

	S, oracleIdentities, err := actions.GetOracleIdentities(chainlinkNodes)
	require.NoError(t, err, "Error getting oracle identities")
	offC, err := json.Marshal(ocr2keepers30config.OffchainConfig{
		TargetProbability:    "0.999",
		TargetInRounds:       1,
		PerformLockoutWindow: 80_000, // Copied from arbitrum mainnet prod value
		GasLimitPerReport:    10_300_000,
		GasOverheadPerUpkeep: 300_000,
		MinConfirmations:     0,
		MaxUpkeepBatchSize:   10,
	})
	require.NoError(t, err, "Error marshalling offchain config")

	signerOnchainPublicKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err := ocr3.ContractSetConfigArgsForTests(
		10*time.Second,        // deltaProgress time.Duration,
		15*time.Second,        // deltaResend time.Duration,
		500*time.Millisecond,  // deltaInitial time.Duration,
		1000*time.Millisecond, // deltaRound time.Duration,
		200*time.Millisecond,  // deltaGrace time.Duration,
		300*time.Millisecond,  // deltaCertifiedCommitRequest time.Duration
		15*time.Second,        // deltaStage time.Duration,
		24,                    // rMax uint64,
		S,                     // s []int,
		oracleIdentities,      // oracles []OracleIdentityExtra,
		offC,                  // reportingPluginConfig []byte,
		20*time.Millisecond,   // maxDurationQuery time.Duration,
		20*time.Millisecond,   // maxDurationObservation time.Duration, // good to here
		1200*time.Millisecond, // maxDurationShouldAcceptAttestedReport time.Duration,
		20*time.Millisecond,   // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,                     // f int,
		nil,                   // onchainConfig []byte,
	)
	require.NoError(t, err, "Error setting OCR config vars")

	var signers []common.Address
	for _, signer := range signerOnchainPublicKeys {
		require.Equal(t, 20, len(signer), "OnChainPublicKey '%v' has wrong length for address", signer)
		signers = append(signers, common.BytesToAddress(signer))
	}

	var transmitters []common.Address
	for _, transmitter := range transmitterAccounts {
		require.True(t, common.IsHexAddress(string(transmitter)), "TransmitAccount '%s' is not a valid Ethereum address", string(transmitter))
		transmitters = append(transmitters, common.HexToAddress(string(transmitter)))
	}

	onchainConfig, err := registrySettings.EncodeOnChainConfig(registrar.Address(), common.HexToAddress(chainClient.GetDefaultWallet().Address()))
	require.NoError(t, err, "Error encoding onchain config")
	l.Info().Msg("Done building OCR config")
	ocrConfig := contracts.OCRv2Config{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}

	err = registry.SetConfig(*registrySettings, ocrConfig)
	require.NoError(t, err, "Error setting registry config")

	consumerContracts := make([]contracts.KeeperConsumer, 0)
	triggerContracts := make([]contracts.LogEmitter, 0)

	utilsABI, err := automation_utils_2_1.AutomationUtilsMetaData.GetAbi()
	require.NoError(t, err, "Error getting automation utils abi")
	registrarABI, err := registrar21.AutomationRegistrarMetaData.GetAbi()
	require.NoError(t, err, "Error getting automation registrar abi")
	emitterABI, err := log_emitter.LogEmitterMetaData.GetAbi()
	require.NoError(t, err, "Error getting log emitter abi")
	consumerABI, err := simple_log_upkeep_counter_wrapper.SimpleLogUpkeepCounterMetaData.GetAbi()
	require.NoError(t, err, "Error getting consumer abi")

	var bytes0 = [32]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	registrationTxHashes := make([]common.Hash, 0)
	upkeepIds := make([]*big.Int, 0)
	evmClients := make([]*blockchain.EVMClient, 0)

	for i := 0; i < numberOfUpkeeps; i++ {
		consumerContract, err := contractDeployer.DeployAutomationSimpleLogTriggerConsumer()
		require.NoError(t, err, "Error deploying automation consumer contract")
		consumerContracts = append(consumerContracts, consumerContract)
		l.Debug().
			Str("Contract Address", consumerContract.Address()).
			Int("Number", i+1).
			Int("Out Of", numberOfUpkeeps).
			Msg("Deployed Automation Log Trigger Consumer Contract")

		cEVMClient, err := blockchain.ConcurrentEVMClient(testNetwork, testEnvironment, chainClient, l)
		require.NoError(t, err, "Error building concurrent chain client")

		cContractDeployer, err := contracts.NewContractDeployer(cEVMClient, l)
		require.NoError(t, err, "Error building concurrent contract deployer")

		evmClients = append(evmClients, &cEVMClient)

		triggerContract, err := cContractDeployer.DeployLogEmitterContract()
		require.NoError(t, err, "Error deploying log emitter contract")
		triggerContracts = append(triggerContracts, triggerContract)
		l.Debug().
			Str("Contract Address", triggerContract.Address().Hex()).
			Int("Number", i+1).
			Int("Out Of", numberOfUpkeeps).
			Msg("Deployed Automation Log Trigger Emitter Contract")
	}

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Failed waiting for contracts to deploy")

	for i, consumerContract := range consumerContracts {
		logTriggerConfigStruct := automation_utils_2_1.LogTriggerConfig{
			ContractAddress: triggerContracts[i].Address(),
			FilterSelector:  0,
			Topic0:          emitterABI.Events["Log1"].ID,
			Topic1:          bytes0,
			Topic2:          bytes0,
			Topic3:          bytes0,
		}
		encodedLogTriggerConfig, err := utilsABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
		require.NoError(t, err, "Error encoding log trigger config")
		l.Debug().Bytes("Encoded Log Trigger Config", encodedLogTriggerConfig).Msg("Encoded Log Trigger Config")

		registrationRequest, err := registrarABI.Pack(
			"register",
			fmt.Sprintf("LogTriggerUpkeep-%d", i),
			[]byte("test@mail.com"),
			common.HexToAddress(consumerContract.Address()),
			automationDefaultUpkeepGasLimit,
			common.HexToAddress(chainClient.GetDefaultWallet().Address()),
			uint8(1),
			[]byte("0"),
			encodedLogTriggerConfig,
			[]byte("0"),
			automationDefaultLinkFunds,
			common.HexToAddress(chainClient.GetDefaultWallet().Address()),
		)
		require.NoError(t, err, "Error encoding upkeep registration request")
		tx, err := linkToken.TransferAndCall(registrar.Address(), automationDefaultLinkFunds, registrationRequest)
		require.NoError(t, err, "Error sending upkeep registration request")
		registrationTxHashes = append(registrationTxHashes, tx.Hash())
	}

	err = chainClient.WaitForEvents()
	require.NoError(t, err, "Failed waiting for upkeeps to be registered")

	for _, txHash := range registrationTxHashes {
		receipt, err := chainClient.GetTxReceipt(txHash)
		require.NoError(t, err, "Registration tx should be completed")
		var upkeepId *big.Int
		for _, rawLog := range receipt.Logs {
			parsedUpkeepId, err := registry.ParseUpkeepIdFromRegisteredLog(rawLog)
			if err == nil {
				upkeepId = parsedUpkeepId
				break
			}
		}
		require.NotNil(t, upkeepId, "Upkeep ID should be found after registration")
		l.Debug().
			Str("TxHash", txHash.String()).
			Str("Upkeep ID", upkeepId.String()).
			Msg("Found upkeepId in tx hash")
		upkeepIds = append(upkeepIds, upkeepId)
	}
	l.Info().Msg("Successfully registered all Automation Consumer Contracts")
	l.Info().Interface("Upkeep IDs", upkeepIds).Msg("Upkeep IDs")
	l.Info().Str("STARTUP_WAIT_TIME", StartupWaitTime.String()).Msg("Waiting for plugin to start")
	time.Sleep(StartupWaitTime)

	startBlock, err := chainClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error getting latest block number")

	p := wasp.NewProfile()

	for i, triggerContract := range triggerContracts {
		g, err := wasp.NewGenerator(&wasp.Config{
			T:           t,
			LoadType:    wasp.RPS,
			GenName:     fmt.Sprintf("log_trigger_gen_%s", triggerContract.Address().String()),
			CallTimeout: time.Second * 10,
			Schedule: wasp.Plain(
				1,
				loadDuration,
			),
			Gun: NewLogTriggerUser(
				triggerContract,
				consumerContracts[i],
				l,
				numberOfEvents,
			),
			CallResultBufLen: 1000000,
		})
		p.Add(g, err)
	}

	l.Info().Msg("Starting load generators")
	startTime := time.Now()
	err = sendSlackNotification("Started", l, testEnvironment.Cfg.Namespace, strconv.Itoa(numberofNodes),
		strconv.FormatInt(startTime.UnixMilli(), 10), "now",
		[]slack.Block{extraBlockWithText("\bTest Config\b\n```" + testConfig + "```")})
	_, err = p.Run(true)
	require.NoError(t, err, "Error running load generators")

	l.Info().Msg("Finished load generators")
	l.Info().Str("STOP_WAIT_TIME", StopWaitTime.String()).Msg("Waiting for upkeeps to be performed")
	time.Sleep(StopWaitTime)
	l.Info().Msg("Finished waiting 60s for upkeeps to be performed")
	endTime := time.Now()
	testDuration := endTime.Sub(startTime)
	l.Info().Str("Duration", testDuration.String()).Msg("Test Duration")
	endBlock, err := chainClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error getting latest block number")
	l.Info().Uint64("Starting Block", startBlock).Uint64("Ending Block", endBlock).Msg("Test Block Range")

	upkeepDelays := make([][]int64, 0)
	var numberOfEventsEmitted int
	var batchSize = 500

	for _, gen := range p.Generators {
		numberOfEventsEmitted += len(gen.GetData().OKData.Data)
	}
	numberOfEventsEmitted = numberOfEventsEmitted * numberOfEvents
	l.Info().Int("Number of Events Emitted", numberOfEventsEmitted).Msg("Number of Events Emitted")

	if endBlock-startBlock < uint64(batchSize) {
		batchSize = int(endBlock - startBlock)
	}

	for cIter, consumerContract := range consumerContracts {
		var (
			logs    []types.Log
			address = common.HexToAddress(consumerContract.Address())
			timeout = 5 * time.Second
		)
		for fromBlock := startBlock; fromBlock < endBlock; fromBlock += uint64(batchSize) + 1 {
			var (
				filterQuery = geth.FilterQuery{
					Addresses: []common.Address{address},
					FromBlock: big.NewInt(0).SetUint64(fromBlock),
					ToBlock:   big.NewInt(0).SetUint64(fromBlock + uint64(batchSize)),
					Topics:    [][]common.Hash{{consumerABI.Events["PerformingUpkeep"].ID}},
				}
			)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			logsInBatch, err := chainClient.FilterLogs(ctx, filterQuery)
			cancel()
			if err != nil {
				l.Error().Err(err).
					Interface("FilterQuery", filterQuery).
					Str("Contract Address", consumerContract.Address()).
					Str("Timeout", timeout.String()).
					Msg("Error getting logs")
			}
			logs = append(logs, logsInBatch...)
			time.Sleep(time.Millisecond * 500)
		}

		if len(logs) > 0 {
			delay := make([]int64, 0)
			for _, log := range logs {
				eventDetails, err := consumerABI.EventByID(log.Topics[0])
				require.NoError(t, err, "Error getting event details")
				consumer, err := simple_log_upkeep_counter_wrapper.NewSimpleLogUpkeepCounter(
					address, chainClient.Backend(),
				)
				require.NoError(t, err, "Error getting consumer contract")
				if eventDetails.Name == "PerformingUpkeep" {
					parsedLog, err := consumer.ParsePerformingUpkeep(log)
					require.NoError(t, err, "Error parsing log")
					delay = append(delay, parsedLog.TimeToPerform.Int64())
				}
			}
			upkeepDelays = append(upkeepDelays, delay)
		}
		if (cIter+1)%batchSize == 0 {
			time.Sleep(time.Millisecond * 500)
		}
	}

	l.Info().Interface("Upkeep Delays", upkeepDelays).Msg("Upkeep Delays")

	var allUpkeepDelays []int64

	for _, upkeepDelay := range upkeepDelays {
		allUpkeepDelays = append(allUpkeepDelays, upkeepDelay...)
	}

	avg, median, ninetyPct, ninetyNinePct, maximum := testreporters.IntListStats(allUpkeepDelays)
	l.Info().
		Float64("Average", avg).Int64("Median", median).
		Int64("90th Percentile", ninetyPct).Int64("99th Percentile", ninetyNinePct).
		Int64("Max", maximum).Msg("Upkeep Delays in seconds")

	l.Info().
		Int("Total Perform Count", len(allUpkeepDelays)).
		Int("Total Events Emitted", numberOfEventsEmitted).
		Int("Total Events Missed", numberOfEventsEmitted-len(allUpkeepDelays)).
		Msg("Test completed")

	testReport := fmt.Sprintf("Upkeep Delays in seconds\nAverage: %f\nMedian: %d\n90th Percentile: %d\n"+
		"99th Percentile: %d\nMax: %d\nTotal Perform Count: %d\n\nTotal Events Emitted: %d\nTotal Events Missed: %d\nTest Duration: %s\n",
		avg, median, ninetyPct, ninetyNinePct, maximum, len(allUpkeepDelays), numberOfEventsEmitted,
		numberOfEventsEmitted-len(allUpkeepDelays), testDuration.String())

	err = sendSlackNotification("Finished", l, testEnvironment.Cfg.Namespace, strconv.Itoa(numberofNodes),
		strconv.FormatInt(startTime.UnixMilli(), 10), strconv.FormatInt(endTime.UnixMilli(), 10),
		[]slack.Block{extraBlockWithText("\bTest Report\b\n```" + testReport + "```")})
	require.NoError(t, err, "Error sending slack notification")

	t.Cleanup(func() {
		if err = actions.TeardownRemoteSuite(t, testEnvironment.Cfg.Namespace, chainlinkNodes, nil, chainClient); err != nil {
			l.Error().Err(err).Msg("Error when tearing down remote suite")
		}
	})

}
