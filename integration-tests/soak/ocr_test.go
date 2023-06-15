package soak

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/networks"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestOCRSoak(t *testing.T) {
	l := utils.GetTestLogger(t)
	testEnvironment, network, testInputs := SetupOCRSoakEnv(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(network, testEnvironment)
	require.NoError(t, err, "Error connecting to network")

	ocrSoakTest := testsetups.NewOCRSoakTest(&testsetups.OCRSoakTestInputs{
		BlockchainClient:     chainClient,
		TestDuration:         testInputs.TestDuration,
		NumberOfContracts:    2,
		ChainlinkNodeFunding: big.NewFloat(testInputs.ChainlinkNodeFunding),
		ExpectedRoundTime:    time.Minute * 2,
		RoundTimeout:         time.Minute * 15,
		TimeBetweenRounds:    testInputs.TimeBetweenRounds,
		StartingAdapterValue: 5,
	})
	t.Cleanup(func() {
		if err := actions.TeardownRemoteSuite(ocrSoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	ocrSoakTest.Setup(t, testEnvironment)
	l.Info().Msg("Set up soak test")
	ocrSoakTest.Run(t)
}

func SetupOCRSoakEnv(t *testing.T) (*environment.Environment, blockchain.EVMNetwork, OcrSoakInputs) {
	var testInputs OcrSoakInputs
	err := envconfig.Process("OCR", &testInputs)
	require.NoError(t, err, "Error reading OCR soak test inputs")
	testInputs.setForRemoteRunner()
	network := networks.SelectedNetwork // Environment currently being used to soak test on

	baseEnvironmentConfig := &environment.Config{
		TTL: time.Hour * 720, // 30 days,
		NamespacePrefix: fmt.Sprintf(
			"soak-ocr-%s",
			strings.ReplaceAll(strings.ToLower(network.Name), " ", "-"),
		),
		Test: t,
	}

	// Use this variable to pass in any custom EVM specific TOML values to your Chainlink nodes
	customNetworkTOML := ``
	// Uncomment below for debugging TOML issues on the node
	// fmt.Println("Using Chainlink TOML\n---------------------")
	// fmt.Println(client.AddNetworkDetailedConfig(config.BaseOCRP2PV1Config, customNetworkTOML, network))
	// fmt.Println("---------------------")
	cd, err := chainlink.NewDeployment(6, map[string]any{
		"toml": client.AddNetworkDetailedConfig(config.BaseOCRP2PV1Config, customNetworkTOML, network),
		"db": map[string]any{
			"stateful": true, // stateful DB by default for soak tests
		},
	})
	require.NoError(t, err, "Error creating chainlink deployment")
	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})).
		AddHelmCharts(cd)
	err = testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")
	return testEnvironment, network, testInputs
}

type OcrSoakInputs struct {
	TestDuration         time.Duration `envconfig:"TEST_DURATION" default:"15m"`
	ChainlinkNodeFunding float64       `envconfig:"CHAINLINK_NODE_FUNDING" default:".1"`
	TimeBetweenRounds    time.Duration `envconfig:"TIME_BETWEEN_ROUNDS" default:"1m"`
}

func (i OcrSoakInputs) setForRemoteRunner() {
	os.Setenv("TEST_OCR_TEST_DURATION", i.TestDuration.String())
	os.Setenv("TEST_OCR_CHAINLINK_NODE_FUNDING", strconv.FormatFloat(i.ChainlinkNodeFunding, 'f', -1, 64))
	os.Setenv("TEST_OCR_TIME_BETWEEN_ROUNDS", i.TimeBetweenRounds.String())

	selectedNetworks := strings.Split(os.Getenv("SELECTED_NETWORKS"), ",")
	for _, networkPrefix := range selectedNetworks {
		urlEnv := fmt.Sprintf("%s_URLS", networkPrefix)
		httpEnv := fmt.Sprintf("%s_HTTP_URLS", networkPrefix)
		os.Setenv(fmt.Sprintf("TEST_%s", urlEnv), os.Getenv(urlEnv))
		os.Setenv(fmt.Sprintf("TEST_%s", httpEnv), os.Getenv(httpEnv))
	}
}
