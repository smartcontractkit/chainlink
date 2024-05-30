package automationv2_1

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pelletier/go-toml/v2"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/require"

	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/wasp"

	ocr2keepers30config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/wiremock"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	seth_utils "github.com/smartcontractkit/chainlink-testing-framework/utils/seth"

	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/config"

	gowiremock "github.com/wiremock/go-wiremock"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/automationv2"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	contractseth "github.com/smartcontractkit/chainlink/integration-tests/contracts/ethereum"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	aconfig "github.com/smartcontractkit/chainlink/integration-tests/testconfig/automation"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
	ac "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_compatible_utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/simple_log_upkeep_counter_wrapper"
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
	secretsTOML = `[Mercury.Credentials.%s]
LegacyURL = '%s'
URL = '%s'
Username = '%s'
Password = '%s'`

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
				"cpu":    "4000m",
				"memory": "4Gi",
			},
			"limits": map[string]interface{}{
				"cpu":    "4000m",
				"memory": "4Gi",
			},
		},
		"stateful": true,
		"capacity": "20Gi",
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

	recDbSpec = minimumDbSpec

	gethNodeSpec = map[string]interface{}{
		"requests": map[string]interface{}{
			"cpu":    "8000m",
			"memory": "8Gi",
		},
		"limits": map[string]interface{}{
			"cpu":    "16000m",
			"memory": "16Gi",
		},
	}
)

func setUpDataStreamsWireMock(url string) error {
	wm := gowiremock.NewClient(url)
	rule200 := gowiremock.Get(gowiremock.URLPathEqualTo("/api/v1/reports/bulk")).
		WithQueryParam("feedIDs", gowiremock.EqualTo("0x000200")).
		WillReturnResponse(gowiremock.NewResponse().
			WithBody(`{"reports":[{"feedID":"0x000200","validFromTimestamp":0,"observationsTimestamp":0,"fullReport":"0x000abc"}]}`).
			WithStatus(200))
	err := wm.StubFor(rule200)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/__admin/mappings/save", url), "application/json", nil)
	if err != nil {
		return errors.New("error saving wiremock mappings")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return errors.New("error saving wiremock mappings")
	}

	return nil
}

func TestLogTrigger(t *testing.T) {
	ctx := tests.Context(t)
	l := logging.GetTestLogger(t)

	loadedTestConfig, err := tc.GetConfig("Load", tc.Automation)
	if err != nil {
		t.Fatal(err)
	}
	l.Info().Interface("loadedTestConfig", loadedTestConfig).Msg("Loaded Test Config")

	version := *loadedTestConfig.ChainlinkImage.Version
	image := *loadedTestConfig.ChainlinkImage.Image

	l.Info().Msg("Starting automation v2.1 log trigger load test")
	l.Info().
		Int("Number of Nodes", *loadedTestConfig.Automation.General.NumberOfNodes).
		Int("Duration", *loadedTestConfig.Automation.General.Duration).
		Int("Block Time", *loadedTestConfig.Automation.General.BlockTime).
		Str("Spec Type", *loadedTestConfig.Automation.General.SpecType).
		Str("Log Level", *loadedTestConfig.Automation.General.ChainlinkNodeLogLevel).
		Str("Image", image).
		Str("Tag", version).
		Msg("Test Config")

	testConfigFormat := `Number of Nodes: %d
Duration: %d
Block Time: %d
Spec Type: %s
Log Level: %s
Image: %s
Tag: %s

Load Config:
%s`

	prettyLoadConfig, err := toml.Marshal(loadedTestConfig.Automation.Load)
	require.NoError(t, err, "Error marshalling load config")

	testConfig := fmt.Sprintf(testConfigFormat, *loadedTestConfig.Automation.General.NumberOfNodes, *loadedTestConfig.Automation.General.Duration,
		*loadedTestConfig.Automation.General.BlockTime, *loadedTestConfig.Automation.General.SpecType, *loadedTestConfig.Automation.General.ChainlinkNodeLogLevel, image, version, string(prettyLoadConfig))
	l.Info().Str("testConfig", testConfig).Msg("Test Config")

	testNetwork := networks.MustGetSelectedNetworkConfig(loadedTestConfig.Network)[0]
	testType := "load"
	loadDuration := time.Duration(*loadedTestConfig.Automation.General.Duration) * time.Second
	automationDefaultLinkFunds := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(int64(10000))) //10000 LINK

	testEnvironment := environment.New(&environment.Config{
		TTL: loadDuration.Round(time.Hour) + time.Hour,
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
				"resources": gethNodeSpec,
				"geth": map[string]interface{}{
					"blocktime":      *loadedTestConfig.Automation.General.BlockTime,
					"capacity":       "20Gi",
					"startGaslimit":  "20000000",
					"targetGasLimit": "30000000",
				},
			},
		}))

	err = testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")

	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	var (
		nodeSpec = minimumNodeSpec
		dbSpec   = minimumDbSpec
	)

	switch *loadedTestConfig.Automation.General.SpecType {
	case "recommended":
		nodeSpec = recNodeSpec
		dbSpec = recDbSpec
	case "local":
		nodeSpec = map[string]interface{}{}
		dbSpec = map[string]interface{}{"stateful": true}
	default:
		// minimum:

	}

	if *loadedTestConfig.Pyroscope.Enabled {
		loadedTestConfig.Pyroscope.Environment = &testEnvironment.Cfg.Namespace
	}

	if *loadedTestConfig.Automation.DataStreams.Enabled {
		if loadedTestConfig.Automation.DataStreams.URL == nil || *loadedTestConfig.Automation.DataStreams.URL == "" {
			testEnvironment.AddHelm(wiremock.New(nil))
			err := testEnvironment.Run()
			require.NoError(t, err, "Error running wiremock server")
			wiremockURL := testEnvironment.URLs[wiremock.InternalURLsKey][0]
			secretsTOML = fmt.Sprintf(
				secretsTOML, "cred1",
				wiremockURL, wiremockURL,
				"username", "password",
			)
			if !testEnvironment.Cfg.InsideK8s {
				wiremockURL = testEnvironment.URLs[wiremock.LocalURLsKey][0]
			}
			err = setUpDataStreamsWireMock(wiremockURL)
			require.NoError(t, err, "Error setting up wiremock server")
		} else {
			secretsTOML = fmt.Sprintf(
				secretsTOML, "cred1",
				*loadedTestConfig.Automation.DataStreams.URL, *loadedTestConfig.Automation.DataStreams.URL,
				*loadedTestConfig.Automation.DataStreams.Username, *loadedTestConfig.Automation.DataStreams.Password,
			)
		}
	} else {
		secretsTOML = ""
	}

	numberOfUpkeeps := *loadedTestConfig.Automation.General.NumberOfNodes

	for i := 0; i < numberOfUpkeeps+1; i++ { // +1 for the OCR boot node
		var nodeTOML string
		if i == 1 || i == 3 {
			nodeTOML = fmt.Sprintf("%s\n\n[Log]\nLevel = \"%s\"", baseTOML, *loadedTestConfig.Automation.General.ChainlinkNodeLogLevel)
		} else {
			nodeTOML = fmt.Sprintf("%s\n\n[Log]\nLevel = \"info\"", baseTOML)
		}
		nodeTOML = networks.AddNetworksConfig(nodeTOML, loadedTestConfig.Pyroscope, testNetwork)

		var overrideFn = func(_ interface{}, target interface{}) {
			ctfconfig.MustConfigOverrideChainlinkVersion(loadedTestConfig.GetChainlinkImageConfig(), target)
			ctfconfig.MightConfigOverridePyroscopeKey(loadedTestConfig.GetPyroscopeConfig(), target)
		}

		cd := chainlink.NewWithOverride(i, map[string]any{
			"toml":        nodeTOML,
			"chainlink":   nodeSpec,
			"db":          dbSpec,
			"prometheus":  *loadedTestConfig.Automation.General.UsePrometheus,
			"secretsToml": secretsTOML,
		}, loadedTestConfig.ChainlinkImage, overrideFn)
		testEnvironment.AddHelm(cd)
	}

	err = testEnvironment.Run()
	require.NoError(t, err, "Error running chainlink DON")

	testNetwork = seth_utils.MustReplaceSimulatedNetworkUrlWithK8(l, testNetwork, *testEnvironment)

	chainClient, err := actions_seth.GetChainClientWithConfigFunction(loadedTestConfig, testNetwork, actions_seth.OneEphemeralKeysLiveTestnetCheckFn)
	require.NoError(t, err, "Error creating seth client")

	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to chainlink nodes")

	multicallAddress, err := contracts.DeployMultiCallContract(chainClient)
	require.NoError(t, err, "Error deploying multicall contract")

	a := automationv2.NewAutomationTestK8s(l, chainClient, chainlinkNodes)
	conf := loadedTestConfig.Automation.AutomationConfig
	a.RegistrySettings = contracts.KeeperRegistrySettings{
		PaymentPremiumPPB:    *conf.RegistrySettings.PaymentPremiumPPB,
		FlatFeeMicroLINK:     *conf.RegistrySettings.FlatFeeMicroLINK,
		CheckGasLimit:        *conf.RegistrySettings.CheckGasLimit,
		StalenessSeconds:     conf.RegistrySettings.StalenessSeconds,
		GasCeilingMultiplier: *conf.RegistrySettings.GasCeilingMultiplier,
		MaxPerformGas:        *conf.RegistrySettings.MaxPerformGas,
		MinUpkeepSpend:       conf.RegistrySettings.MinUpkeepSpend,
		FallbackGasPrice:     conf.RegistrySettings.FallbackGasPrice,
		FallbackLinkPrice:    conf.RegistrySettings.FallbackLinkPrice,
		MaxCheckDataSize:     *conf.RegistrySettings.MaxCheckDataSize,
		MaxPerformDataSize:   *conf.RegistrySettings.MaxPerformDataSize,
		MaxRevertDataSize:    *conf.RegistrySettings.MaxRevertDataSize,
		RegistryVersion:      contractseth.RegistryVersion_2_1,
	}
	a.RegistrarSettings = contracts.KeeperRegistrarSettings{
		AutoApproveConfigType: uint8(2),
		AutoApproveMaxAllowed: math.MaxUint16,
		MinLinkJuels:          big.NewInt(0),
	}
	a.PluginConfig = ocr2keepers30config.OffchainConfig{
		TargetProbability:    *conf.PluginConfig.TargetProbability,
		TargetInRounds:       *conf.PluginConfig.TargetInRounds,
		PerformLockoutWindow: *conf.PluginConfig.PerformLockoutWindow,
		GasLimitPerReport:    *conf.PluginConfig.GasLimitPerReport,
		GasOverheadPerUpkeep: *conf.PluginConfig.GasOverheadPerUpkeep,
		MinConfirmations:     *conf.PluginConfig.MinConfirmations,
		MaxUpkeepBatchSize:   *conf.PluginConfig.MaxUpkeepBatchSize,
		LogProviderConfig: ocr2keepers30config.LogProviderConfig{
			BlockRate: *conf.PluginConfig.LogProviderConfig.BlockRate,
			LogLimit:  *conf.PluginConfig.LogProviderConfig.LogLimit,
		},
	}
	a.PublicConfig = ocr3.PublicConfig{
		DeltaProgress:                           *conf.PublicConfig.DeltaProgress,
		DeltaResend:                             *conf.PublicConfig.DeltaResend,
		DeltaInitial:                            *conf.PublicConfig.DeltaInitial,
		DeltaRound:                              *conf.PublicConfig.DeltaRound,
		DeltaGrace:                              *conf.PublicConfig.DeltaGrace,
		DeltaCertifiedCommitRequest:             *conf.PublicConfig.DeltaCertifiedCommitRequest,
		DeltaStage:                              *conf.PublicConfig.DeltaStage,
		RMax:                                    *conf.PublicConfig.RMax,
		MaxDurationQuery:                        *conf.PublicConfig.MaxDurationQuery,
		MaxDurationObservation:                  *conf.PublicConfig.MaxDurationObservation,
		MaxDurationShouldAcceptAttestedReport:   *conf.PublicConfig.MaxDurationShouldAcceptAttestedReport,
		MaxDurationShouldTransmitAcceptedReport: *conf.PublicConfig.MaxDurationShouldTransmitAcceptedReport,
		F:                                       *conf.PublicConfig.F,
	}

	if *loadedTestConfig.Automation.DataStreams.Enabled {
		a.SetMercuryCredentialName("cred1")
	}

	if *conf.UseLogBufferV1 {
		a.SetUseLogBufferV1(true)
	}

	startTimeTestSetup := time.Now()
	l.Info().Str("START_TIME", startTimeTestSetup.String()).Msg("Test setup started")

	a.SetupAutomationDeployment(t)

	err = actions_seth.FundChainlinkNodesFromRootAddress(l, a.ChainClient, contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes[1:]), big.NewFloat(*loadedTestConfig.Common.ChainlinkNodeFunding))
	require.NoError(t, err, "Error funding chainlink nodes")

	consumerContracts := make([]contracts.KeeperConsumer, 0)
	triggerContracts := make([]contracts.LogEmitter, 0)
	triggerAddresses := make([]common.Address, 0)

	convenienceABI, err := ac.AutomationCompatibleUtilsMetaData.GetAbi()
	require.NoError(t, err, "Error getting automation utils abi")
	emitterABI, err := log_emitter.LogEmitterMetaData.GetAbi()
	require.NoError(t, err, "Error getting log emitter abi")
	consumerABI, err := simple_log_upkeep_counter_wrapper.SimpleLogUpkeepCounterMetaData.GetAbi()
	require.NoError(t, err, "Error getting consumer abi")

	var bytes0 = [32]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	var bytes1 = [32]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
	}

	upkeepConfigs := make([]automationv2.UpkeepConfig, 0)
	loadConfigs := make([]aconfig.Load, 0)

	expectedTotalUpkeepCount := 0
	for _, u := range loadedTestConfig.Automation.Load {
		expectedTotalUpkeepCount += *u.NumberOfUpkeeps
	}

	maxDeploymentConcurrency := 100

	for _, u := range loadedTestConfig.Automation.Load {
		deploymentData, err := deployConsumerAndTriggerContracts(l, u, a.ChainClient, multicallAddress, maxDeploymentConcurrency, automationDefaultLinkFunds, a.LinkToken)
		require.NoError(t, err, "Error deploying consumer and trigger contracts")

		consumerContracts = append(consumerContracts, deploymentData.ConsumerContracts...)
		triggerContracts = append(triggerContracts, deploymentData.TriggerContracts...)
		triggerAddresses = append(triggerAddresses, deploymentData.TriggerAddresses...)
		loadConfigs = append(loadConfigs, deploymentData.LoadConfigs...)
	}

	require.Equal(t, expectedTotalUpkeepCount, len(consumerContracts), "Incorrect number of consumer/trigger contracts deployed")

	for i, consumerContract := range consumerContracts {
		logTriggerConfigStruct := ac.IAutomationV21PlusCommonLogTriggerConfig{
			ContractAddress: triggerAddresses[i],
			FilterSelector:  1,
			Topic0:          emitterABI.Events["Log4"].ID,
			Topic1:          bytes1,
			Topic2:          bytes0,
			Topic3:          bytes0,
		}
		encodedLogTriggerConfig, err := convenienceABI.Methods["_logTriggerConfig"].Inputs.Pack(&logTriggerConfigStruct)
		require.NoError(t, err, "Error encoding log trigger config")
		l.Debug().
			Interface("logTriggerConfigStruct", logTriggerConfigStruct).
			Str("Encoded Log Trigger Config", hex.EncodeToString(encodedLogTriggerConfig)).Msg("Encoded Log Trigger Config")

		checkDataStruct := simple_log_upkeep_counter_wrapper.CheckData{
			CheckBurnAmount:   loadConfigs[i].CheckBurnAmount,
			PerformBurnAmount: loadConfigs[i].PerformBurnAmount,
			EventSig:          bytes1,
			Feeds:             loadConfigs[i].Feeds,
		}

		encodedCheckDataStruct, err := consumerABI.Methods["_checkDataConfig"].Inputs.Pack(&checkDataStruct)
		require.NoError(t, err, "Error encoding check data struct")
		l.Debug().
			Interface("checkDataStruct", checkDataStruct).
			Str("Encoded Check Data Struct", hex.EncodeToString(encodedCheckDataStruct)).Msg("Encoded Check Data Struct")

		upkeepConfig := automationv2.UpkeepConfig{
			UpkeepName:     fmt.Sprintf("LogTriggerUpkeep-%d", i),
			EncryptedEmail: []byte("test@mail.com"),
			UpkeepContract: common.HexToAddress(consumerContract.Address()),
			GasLimit:       *loadConfigs[i].UpkeepGasLimit,
			AdminAddress:   chainClient.MustGetRootKeyAddress(),
			TriggerType:    uint8(1),
			CheckData:      encodedCheckDataStruct,
			TriggerConfig:  encodedLogTriggerConfig,
			OffchainConfig: []byte(""),
			FundingAmount:  automationDefaultLinkFunds,
		}
		l.Debug().Interface("Upkeep Config", upkeepConfig).Msg("Upkeep Config")
		upkeepConfigs = append(upkeepConfigs, upkeepConfig)
	}

	require.Equal(t, expectedTotalUpkeepCount, len(upkeepConfigs), "Incorrect number of upkeep configs created")
	registrationTxHashes, err := a.RegisterUpkeeps(upkeepConfigs, maxDeploymentConcurrency)
	require.NoError(t, err, "Error registering upkeeps")

	upkeepIds, err := a.ConfirmUpkeepsRegistered(registrationTxHashes, maxDeploymentConcurrency)
	require.NoError(t, err, "Error confirming upkeeps registered")
	require.Equal(t, expectedTotalUpkeepCount, len(upkeepIds), "Incorrect number of upkeeps registered")

	l.Info().Msg("Successfully registered all Automation Upkeeps")
	l.Info().Interface("Upkeep IDs", upkeepIds).Msg("Upkeeps Registered")
	l.Info().Str("STARTUP_WAIT_TIME", StartupWaitTime.String()).Msg("Waiting for plugin to start")
	time.Sleep(StartupWaitTime)

	startBlock, err := a.ChainClient.Client.BlockNumber(ctx)
	require.NoError(t, err, "Error getting latest block number")

	p := wasp.NewProfile()

	configs := make([]LogTriggerConfig, 0)
	var numberOfEventsEmitted int64
	var numberOfEventsEmittedPerSec int64

	for i, triggerContract := range triggerContracts {
		c := LogTriggerConfig{
			Address:                       triggerContract.Address().String(),
			NumberOfEvents:                int64(*loadConfigs[i].NumberOfEvents),
			NumberOfSpamMatchingEvents:    int64(*loadConfigs[i].NumberOfSpamMatchingEvents),
			NumberOfSpamNonMatchingEvents: int64(*loadConfigs[i].NumberOfSpamNonMatchingEvents),
		}
		numberOfEventsEmittedPerSec = numberOfEventsEmittedPerSec + int64(*loadConfigs[i].NumberOfEvents)
		configs = append(configs, c)
	}

	endTimeTestSetup := time.Now()
	testSetupDuration := endTimeTestSetup.Sub(startTimeTestSetup)
	l.Info().
		Str("END_TIME", endTimeTestSetup.String()).
		Str("Duration", testSetupDuration.String()).
		Msg("Test setup ended")

	ts, err := sendSlackNotification("Started :white_check_mark:", l, &loadedTestConfig, testEnvironment.Cfg.Namespace, strconv.Itoa(*loadedTestConfig.Automation.General.NumberOfNodes),
		strconv.FormatInt(startTimeTestSetup.UnixMilli(), 10), "now",
		[]slack.Block{extraBlockWithText("\bTest Config\b\n```" + testConfig + "```")}, slack.MsgOptionBlocks())
	if err != nil {
		l.Error().Err(err).Msg("Error sending slack notification")
	}

	g, err := wasp.NewGenerator(&wasp.Config{
		T:           t,
		LoadType:    wasp.RPS,
		GenName:     "log_trigger_gen",
		CallTimeout: time.Minute * 3,
		Schedule: wasp.Plain(
			1,
			loadDuration,
		),
		Gun: NewLogTriggerUser(
			l,
			configs,
			a.ChainClient,
			multicallAddress.Hex(),
		),
		CallResultBufLen: 1000,
	})
	p.Add(g, err)

	startTimeTestEx := time.Now()
	l.Info().Str("START_TIME", startTimeTestEx.String()).Msg("Test execution started")

	l.Info().Msg("Starting load generators")
	_, err = p.Run(true)
	require.NoError(t, err, "Error running load generators")

	l.Info().Msg("Finished load generators")
	l.Info().Str("STOP_WAIT_TIME", StopWaitTime.String()).Msg("Waiting for upkeeps to be performed")
	time.Sleep(StopWaitTime)
	l.Info().Msg("Finished waiting 60s for upkeeps to be performed")
	endTimeTestEx := time.Now()
	testExDuration := endTimeTestEx.Sub(startTimeTestEx)
	l.Info().
		Str("END_TIME", endTimeTestEx.String()).
		Str("Duration", testExDuration.String()).
		Msg("Test execution ended")

	l.Info().Str("Duration", testExDuration.String()).Msg("Test Execution Duration")
	endBlock, err := chainClient.Client.BlockNumber(ctx)
	require.NoError(t, err, "Error getting latest block number")
	l.Info().Uint64("Starting Block", startBlock).Uint64("Ending Block", endBlock).Msg("Test Block Range")

	startTimeTestReport := time.Now()
	l.Info().Str("START_TIME", startTimeTestReport.String()).Msg("Test reporting started")

	for _, gen := range p.Generators {
		if len(gen.Errors()) != 0 {
			l.Error().Strs("Errors", gen.Errors()).Msg("Error in load gen")
			t.Fail()
		}
	}

	upkeepDelaysFast := make([][]int64, 0)
	upkeepDelaysRecovery := make([][]int64, 0)

	var batchSize uint64 = 500

	if endBlock-startBlock < batchSize {
		batchSize = endBlock - startBlock
	}

	for _, consumerContract := range consumerContracts {
		var (
			logs    []types.Log
			address = common.HexToAddress(consumerContract.Address())
			timeout = 5 * time.Second
		)
		for fromBlock := startBlock; fromBlock < endBlock; fromBlock += batchSize + 1 {
			filterQuery := geth.FilterQuery{
				Addresses: []common.Address{address},
				FromBlock: big.NewInt(0).SetUint64(fromBlock),
				ToBlock:   big.NewInt(0).SetUint64(fromBlock + batchSize),
				Topics:    [][]common.Hash{{consumerABI.Events["PerformingUpkeep"].ID}},
			}
			err = fmt.Errorf("initial error") // to ensure our for loop runs at least once
			for err != nil {
				var (
					logsInBatch []types.Log
				)
				ctx2, cancel := context.WithTimeout(ctx, timeout)
				logsInBatch, err = a.ChainClient.Client.FilterLogs(ctx2, filterQuery)
				cancel()
				if err != nil {
					l.Error().Err(err).
						Interface("FilterQuery", filterQuery).
						Str("Contract Address", consumerContract.Address()).
						Str("Timeout", timeout.String()).
						Msg("Error getting consumer contract logs")
					timeout = time.Duration(math.Min(float64(timeout)*2, float64(2*time.Minute)))
					continue
				}
				l.Debug().
					Interface("FilterQuery", filterQuery).
					Str("Contract Address", consumerContract.Address()).
					Str("Timeout", timeout.String()).
					Int("Number of Logs", len(logsInBatch)).
					Msg("Collected consumer contract logs")
				logs = append(logs, logsInBatch...)
			}
		}

		if len(logs) > 0 {
			delayFast := make([]int64, 0)
			delayRecovery := make([]int64, 0)
			for _, log := range logs {
				eventDetails, err := consumerABI.EventByID(log.Topics[0])
				require.NoError(t, err, "Error getting event details")
				consumer, err := simple_log_upkeep_counter_wrapper.NewSimpleLogUpkeepCounter(
					address, a.ChainClient.Client,
				)
				require.NoError(t, err, "Error getting consumer contract")
				if eventDetails.Name == "PerformingUpkeep" {
					parsedLog, err := consumer.ParsePerformingUpkeep(log)
					require.NoError(t, err, "Error parsing log")
					if parsedLog.IsRecovered {
						delayRecovery = append(delayRecovery, parsedLog.TimeToPerform.Int64())
					} else {
						delayFast = append(delayFast, parsedLog.TimeToPerform.Int64())
					}
				}
			}
			upkeepDelaysFast = append(upkeepDelaysFast, delayFast)
			upkeepDelaysRecovery = append(upkeepDelaysRecovery, delayRecovery)
		}
	}

	for _, triggerContract := range triggerContracts {
		var (
			logs    []types.Log
			address = triggerContract.Address()
			timeout = 5 * time.Second
		)
		for fromBlock := startBlock; fromBlock < endBlock; fromBlock += batchSize + 1 {
			filterQuery := geth.FilterQuery{
				Addresses: []common.Address{address},
				FromBlock: big.NewInt(0).SetUint64(fromBlock),
				ToBlock:   big.NewInt(0).SetUint64(fromBlock + batchSize),
				Topics:    [][]common.Hash{{emitterABI.Events["Log4"].ID}, {bytes1}, {bytes1}},
			}
			err = fmt.Errorf("initial error") // to ensure our for loop runs at least once
			for err != nil {
				var (
					logsInBatch []types.Log
				)
				ctx2, cancel := context.WithTimeout(ctx, timeout)
				logsInBatch, err = chainClient.Client.FilterLogs(ctx2, filterQuery)
				cancel()
				if err != nil {
					l.Error().Err(err).
						Interface("FilterQuery", filterQuery).
						Str("Contract Address", address.Hex()).
						Str("Timeout", timeout.String()).
						Msg("Error getting trigger contract logs")
					timeout = time.Duration(math.Min(float64(timeout)*2, float64(2*time.Minute)))
					continue
				}
				l.Debug().
					Interface("FilterQuery", filterQuery).
					Str("Contract Address", address.Hex()).
					Str("Timeout", timeout.String()).
					Int("Number of Logs", len(logsInBatch)).
					Msg("Collected trigger contract logs")
				logs = append(logs, logsInBatch...)
			}
		}
		numberOfEventsEmitted = numberOfEventsEmitted + int64(len(logs))
	}

	l.Info().Int64("Number of Events Emitted", numberOfEventsEmitted).Msg("Number of Events Emitted")

	l.Info().
		Interface("Upkeep Delays Fast", upkeepDelaysFast).
		Interface("Upkeep Delays Recovered", upkeepDelaysRecovery).
		Msg("Upkeep Delays")

	var allUpkeepDelays []int64
	var allUpkeepDelaysFast []int64
	var allUpkeepDelaysRecovery []int64

	for _, upkeepDelay := range upkeepDelaysFast {
		allUpkeepDelays = append(allUpkeepDelays, upkeepDelay...)
		allUpkeepDelaysFast = append(allUpkeepDelaysFast, upkeepDelay...)
	}

	for _, upkeepDelay := range upkeepDelaysRecovery {
		allUpkeepDelays = append(allUpkeepDelays, upkeepDelay...)
		allUpkeepDelaysRecovery = append(allUpkeepDelaysRecovery, upkeepDelay...)
	}

	avgF, medianF, ninetyPctF, ninetyNinePctF, maximumF := testreporters.IntListStats(allUpkeepDelaysFast)
	avgR, medianR, ninetyPctR, ninetyNinePctR, maximumR := testreporters.IntListStats(allUpkeepDelaysRecovery)
	eventsMissed := (numberOfEventsEmitted) - int64(len(allUpkeepDelays))
	percentMissed := float64(eventsMissed) / float64(numberOfEventsEmitted) * 100
	l.Info().
		Float64("Average", avgF).Int64("Median", medianF).
		Int64("90th Percentile", ninetyPctF).Int64("99th Percentile", ninetyNinePctF).
		Int64("Max", maximumF).Msg("Upkeep Delays Fast Execution in seconds")
	l.Info().
		Float64("Average", avgR).Int64("Median", medianR).
		Int64("90th Percentile", ninetyPctR).Int64("99th Percentile", ninetyNinePctR).
		Int64("Max", maximumR).Msg("Upkeep Delays Recovery Execution in seconds")
	l.Info().
		Int("Total Perform Count", len(allUpkeepDelays)).
		Int("Perform Count Fast Execution", len(allUpkeepDelaysFast)).
		Int("Perform Count Recovery Execution", len(allUpkeepDelaysRecovery)).
		Int64("Total Events Emitted", numberOfEventsEmitted).
		Int64("Total Events Missed", eventsMissed).
		Float64("Percent Missed", percentMissed).
		Msg("Test completed")

	testReportFormat := `Upkeep Delays in seconds - Fast Execution
Average: %f
Median: %d
90th Percentile: %d
99th Percentile: %d
Max: %d

Upkeep Delays in seconds - Recovery Execution
Average: %f
Median: %d
90th Percentile: %d
99th Percentile: %d
Max: %d

Total Perform Count: %d
Perform Count Fast Execution: %d
Perform Count Recovery Execution: %d
Total Expected Log Triggering Events: %d
Total Log Triggering Events Emitted: %d
Total Events Missed: %d
Percent Missed: %f
Test Duration: %s`

	endTimeTestReport := time.Now()
	testReDuration := endTimeTestReport.Sub(startTimeTestReport)
	l.Info().
		Str("END_TIME", endTimeTestReport.String()).
		Str("Duration", testReDuration.String()).
		Msg("Test reporting ended")

	numberOfExpectedEvents := numberOfEventsEmittedPerSec * int64(loadDuration.Seconds())
	if numberOfEventsEmitted < numberOfExpectedEvents {
		l.Error().Msg("Number of events emitted is less than expected")
		t.Fail()
	}
	testReport := fmt.Sprintf(testReportFormat, avgF, medianF, ninetyPctF, ninetyNinePctF, maximumF,
		avgR, medianR, ninetyPctR, ninetyNinePctR, maximumR, len(allUpkeepDelays), len(allUpkeepDelaysFast),
		len(allUpkeepDelaysRecovery), numberOfExpectedEvents, numberOfEventsEmitted, eventsMissed, percentMissed, testExDuration.String())
	l.Info().Str("Test Report", testReport).Msg("Test Report prepared")

	testStatus := "Failed :x:"
	if !t.Failed() {
		testStatus = "Finished :white_check_mark:"
	}

	_, err = sendSlackNotification(testStatus, l, &loadedTestConfig, testEnvironment.Cfg.Namespace, strconv.Itoa(*loadedTestConfig.Automation.General.NumberOfNodes),
		strconv.FormatInt(startTimeTestSetup.UnixMilli(), 10), strconv.FormatInt(time.Now().UnixMilli(), 10),
		[]slack.Block{extraBlockWithText("\bTest Report\b\n```" + testReport + "```")}, slack.MsgOptionTS(ts))
	if err != nil {
		l.Error().Err(err).Msg("Error sending slack notification")
	}

	t.Cleanup(func() {
		if err = actions_seth.TeardownRemoteSuite(t, chainClient, testEnvironment.Cfg.Namespace, chainlinkNodes, nil, &loadedTestConfig); err != nil {
			l.Error().Err(err).Msg("Error when tearing down remote suite")
			testEnvironment.Cfg.TTL += time.Hour * 48
			err := testEnvironment.Run()
			if err != nil {
				l.Error().Err(err).Msg("Error increasing TTL of namespace")
			}
		} else if chainClient.Cfg.IsSimulatedNetwork() && *loadedTestConfig.Automation.General.RemoveNamespace {
			err := testEnvironment.Client.RemoveNamespace(testEnvironment.Cfg.Namespace)
			if err != nil {
				l.Error().Err(err).Msg("Error removing namespace")
			}
		}
	})

}
