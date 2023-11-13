package automationv2_1

import (
	"context"
	"fmt"
	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	contractseth "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	registrar21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registrar_wrapper2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/simple_log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"
	"math/big"
	"strings"
	"testing"
	"time"
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
)

func TestLogTrigger(t *testing.T) {
	l := logging.GetTestLogger(t)

	l.Info().Msg("Starting basic log trigger test")

	testNetwork := networks.MustGetSelectedNetworksFromEnv()[0]
	testType := "load"
	numberofNodes := 6
	networkDetailTOML := `MinIncomingConfirmations = 1`
	blockTime := "1"
	numberOfUpkeeps := 500
	const durationInSeconds = 300
	automationDefaultLinkFunds := big.NewInt(int64(9e18))
	automationDefaultUpkeepGasLimit := uint32(2500000)

	registrySettings := &contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    uint32(0),
		FlatFeeMicroLINK:     uint32(40000),
		BlockCountPerTurn:    big.NewInt(100),
		CheckGasLimit:        uint32(45_000_000), //45M
		StalenessSeconds:     big.NewInt(90_000),
		GasCeilingMultiplier: uint16(2),
		MaxPerformGas:        uint32(5000000),
		MinUpkeepSpend:       big.NewInt(0),
		FallbackGasPrice:     big.NewInt(2e11),
		FallbackLinkPrice:    big.NewInt(2e18),
		MaxCheckDataSize:     uint32(5_000),
		MaxPerformDataSize:   uint32(5_000),
		RegistryVersion:      contractseth.RegistryVersion_2_1,
	}

	testEnvironment := environment.New(&environment.Config{
		TTL: time.Hour * 720, // 30 days,
		NamespacePrefix: fmt.Sprintf(
			"automation-%s-%s",
			testType,
			strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-"),
		),
		Test:               t,
		PreventPodEviction: true,
	})

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
						"cpu":    "4000m",
						"memory": "4Gi",
					},
				},
				"geth": map[string]interface{}{
					"blocktime": blockTime,
				},
			},
		}))

	err := testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")

	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	for i := 0; i < numberofNodes; i++ {
		testEnvironment.AddHelm(chainlink.New(i, map[string]any{
			"toml": client.AddNetworkDetailedConfig(baseTOML, networkDetailTOML, testNetwork),
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

	err = actions.FundChainlinkNodesAddress(chainlinkNodes[1:], chainClient, big.NewFloat(1.5), 0)
	require.NoError(t, err, "Error funding chainlink nodes")

	actions.CreateOCRKeeperJobs(
		t,
		chainlinkNodes,
		registry.Address(),
		chainClient.GetChainID().Int64(),
		0,
		contractseth.RegistryVersion_2_1,
	)

	ocrConfig, err := actions.BuildAutoOCR2ConfigVars(t, chainlinkNodes[1:], *registrySettings, registrar.Address(), time.Second*15)
	require.NoError(t, err, "Error building OCR config vars")

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

	for i := 0; i < numberOfUpkeeps; i++ {
		consumerContract, err := contractDeployer.DeployAutomationSimpleLogTriggerConsumer()
		require.NoError(t, err, "Error deploying automation consumer contract")
		consumerContracts = append(consumerContracts, consumerContract)
		l.Debug().
			Str("Contract Address", consumerContract.Address()).
			Int("Number", i+1).
			Int("Out Of", numberOfUpkeeps).
			Msg("Deployed Automation Log Trigger Consumer Contract")

		triggerContract, err := contractDeployer.DeployLogEmitterContract()
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
	time.Sleep(time.Second * 30)

	startingBlock, err := chainClient.LatestBlockNumber(context.Background())
	require.NoError(t, err, "Error getting latest block number")

	p := wasp.NewProfile()

	for i, triggerContract := range triggerContracts {
		g, err := wasp.NewGenerator(&wasp.Config{
			T:           t,
			LoadType:    wasp.RPS,
			GenName:     fmt.Sprintf("log_trigger_gen_%s", triggerContract.Address().String()),
			CallTimeout: time.Minute * 3,
			Schedule: wasp.Plain(
				1,
				time.Second*durationInSeconds,
			),
			Gun: NewLogTriggerUser(
				triggerContract,
				consumerContracts[i],
				l,
			),
		})
		p.Add(g, err)
	}

	l.Info().Msg("Starting load generators")
	_, err = p.Run(true)
	require.NoError(t, err, "Error running load generators")

	l.Info().Msg("Finished load generators")
	l.Info().Msg("Waiting for upkeeps to be performed")
	time.Sleep(time.Second * 60)
	l.Info().Msg("Finished waiting for upkeeps to be performed")

	upkeepCounters := make([]int64, 0)
	upkeepDelays := make([][]int64, 0)

	for i, consumerContract := range consumerContracts {
		count, err := consumerContract.Counter(nil)
		require.NoError(t, err, "Error getting counter value")
		upkeepCounters = append(upkeepCounters, count.Int64())
		l.Debug().
			Int("Count", int(count.Int64())).
			Int("Number", i+1).
			Int("Out Of", numberOfUpkeeps).
			Msg("Counter Value")
		//assert.GreaterOrEqual(t, count.Int64(), int64(durationInSeconds+1), "Counter should be greater than 2")
	}

	for _, consumerContract := range consumerContracts {
		var (
			logs        []types.Log
			address     = common.HexToAddress(consumerContract.Address())
			timeout     = 5 * time.Second
			filterQuery = geth.FilterQuery{
				Addresses: []common.Address{address},
				FromBlock: big.NewInt(0).SetUint64(startingBlock),
				Topics:    [][]common.Hash{{consumerABI.Events["PerformingUpkeep"].ID}},
			}
		)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		logs, err = chainClient.FilterLogs(ctx, filterQuery)
		cancel()
		if err != nil {
			l.Error().Err(err).
				Interface("FilterQuery", filterQuery).
				Str("Contract Address", consumerContract.Address()).
				Str("Timeout", timeout.String()).
				Msg("Error getting logs")
		} else {
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
	}

	l.Info().Interface("Upkeep Counters", upkeepCounters).Msg("Upkeep Counters")
	l.Info().Interface("Upkeep Delays", upkeepDelays).Msg("Upkeep Delays")

	t.Cleanup(func() {
		if err = actions.TeardownRemoteSuite(t, testEnvironment.Cfg.Namespace, chainlinkNodes, nil, chainClient); err != nil {
			l.Error().Err(err).Msg("Error when tearing down remote suite")
		}
	})

}
