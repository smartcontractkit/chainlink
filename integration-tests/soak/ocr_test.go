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
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
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
	var ocrEnvVars = map[string]any{
		"P2P_LISTEN_IP":   "0.0.0.0",
		"P2P_LISTEN_PORT": "6690",
	}
	// For if we end up using env vars
	ocrEnvVars["ETH_URL"] = network.URLs[0]
	ocrEnvVars["ETH_HTTP_URL"] = network.HTTPURLs[0]
	ocrEnvVars["ETH_CHAIN_ID"] = fmt.Sprint(network.ChainID)

	baseEnvironmentConfig := &environment.Config{
		TTL: time.Hour * 720, // 30 days,
		NamespacePrefix: fmt.Sprintf(
			"soak-ocr-%s",
			strings.ReplaceAll(strings.ToLower(network.Name), " ", "-"),
		),
		Test: t,
	}

	replicas := 6
	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		}))
	for i := 0; i < replicas; i++ {
		useEnvVars := strings.ToLower(os.Getenv("TEST_USE_ENV_VAR_CONFIG"))
		if useEnvVars == "true" {
			testEnvironment.AddHelm(chainlink.NewVersioned(i, "0.0.11", map[string]any{
				"env": ocrEnvVars,
			}))
		} else {
			testEnvironment.AddHelm(chainlink.New(i, map[string]any{
				"toml": client.AddNetworksConfig(config.BaseOCRP2PV1Config, network),
			}))
		}
	}
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
