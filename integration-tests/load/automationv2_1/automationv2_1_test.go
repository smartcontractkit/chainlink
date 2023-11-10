package automationv2_1

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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
	cltypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	registrar21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_registrar_wrapper2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
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
	numberOfUpkeeps := 2
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

	//linkFeed, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
	//require.NoError(t, err, "Error deploying link feed contract")
	//
	//gasFeed, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
	//require.NoError(t, err, "Error deploying gas feed contract")

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

	var utilsABI = cltypes.MustGetABI(automation_utils_2_1.AutomationUtilsABI)
	var registrarABI = cltypes.MustGetABI(registrar21.AutomationRegistrarABI)
	var emitterABI = cltypes.MustGetABI(log_emitter.LogEmitterABI)
	var bytes0 = [32]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	registrationTxHashes := make([]common.Hash, 0)
	upkeepIds := make([]*big.Int, 0)

	//counterContract, err := contractDeployer.DeployUpkeepCounter(big.NewInt(99999), big.NewInt(1))
	//require.NoError(t, err, "Error deploying upkeep counter contract")
	//err = chainClient.WaitForEvents()
	//require.NoError(t, err, "Failed waiting for contracts to deploy")
	//
	//registrationRequest, err := registrarABI.Pack(
	//	"register",
	//	"UpkeepCounter",
	//	[]byte("test@mail.com"),
	//	common.HexToAddress(counterContract.Address()),
	//	automationDefaultUpkeepGasLimit,
	//	common.HexToAddress(chainClient.GetDefaultWallet().Address()),
	//	uint8(0),
	//	[]byte("0"),
	//	[]byte("0"),
	//	[]byte("0"),
	//	automationDefaultLinkFunds,
	//	common.HexToAddress(chainClient.GetDefaultWallet().Address()),
	//)
	//require.NoError(t, err, "Error encoding upkeep registration request")
	//tx, err := linkToken.TransferAndCall(registrar.Address(), automationDefaultLinkFunds, registrationRequest)
	//require.NoError(t, err, "Error sending upkeep registration request")
	//registrationTxHashes = append(registrationTxHashes, tx.Hash())

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
	time.Sleep(time.Second * 90)

	p := wasp.NewProfile()

	for i, triggerContract := range triggerContracts {
		g, err := wasp.NewGenerator(&wasp.Config{
			T:           t,
			LoadType:    wasp.RPS,
			GenName:     fmt.Sprintf("log_trigger_gen_%s", triggerContract.Address().String()),
			CallTimeout: time.Minute * 3,
			Schedule: wasp.Plain(
				1,
				time.Second*30,
			),
			Gun: NewLogTriggerUser(
				&triggerContract,
				&consumerContracts[i],
				l,
			),
		})
		p.Add(g, err)
	}

	_, err = p.Run(true)
	require.NoError(t, err, "Error running load generators")

	t.Cleanup(func() {
		if err = actions.TeardownRemoteSuite(t, testEnvironment.Cfg.Namespace, chainlinkNodes, nil, chainClient); err != nil {
			l.Error().Err(err).Msg("Error when tearing down remote suite")
		}
	})

}
