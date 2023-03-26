package soak

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/logging"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestMain(m *testing.M) {
	logging.Init()
	os.Exit(m.Run())
}

func TestForwarderOCRSoak(t *testing.T) {
	l := utils.GetTestLogger(t)
	testEnvironment, network := SetupForwarderOCRSoakEnv(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(network, testEnvironment)
	require.NoError(t, err, "Error connecting to blockchain")
	ocrSoakTest := testsetups.NewOCRSoakTest(&testsetups.OCRSoakTestInputs{
		BlockchainClient:     chainClient,
		TestDuration:         time.Minute * 15,
		NumberOfContracts:    2,
		ChainlinkNodeFunding: big.NewFloat(.1),
		ExpectedRoundTime:    time.Minute * 2,
		RoundTimeout:         time.Minute * 15,
		TimeBetweenRounds:    time.Minute * 1,
		StartingAdapterValue: 5,
	})
	t.Cleanup(func() {
		if err = actions.TeardownRemoteSuite(ocrSoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error when tearing down remote suite")
		}
	})
	ocrSoakTest.OperatorForwarderFlow = true
	ocrSoakTest.Setup(t, testEnvironment)
	l.Info().Msg("Setup soak test")
	ocrSoakTest.Run(t)
}

func SetupForwarderOCRSoakEnv(t *testing.T) (*environment.Environment, blockchain.EVMNetwork) {
	var (
		ocrForwarderEnvVars = map[string]any{
			"FEATURE_LOG_POLLER": "true",
			"ETH_USE_FORWARDERS": "true",
			"P2P_LISTEN_IP":      "0.0.0.0",
			"P2P_LISTEN_PORT":    "6690",
		}

		ocrForwarderBaseTOML = `[OCR]
	Enabled = true
	
	[Feature]
	LogPoller = true
	
	[P2P]
	[P2P.V1]
	Enabled = true
	ListenIP = '0.0.0.0'
	ListenPort = 6690`

		ocrForwarderNetworkDetailTOML = `[EVM.Transactions]
	ForwardersEnabled = true`
	)
	network := networks.SelectedNetwork // Environment currently being used to soak test on
	ocrForwarderEnvVars["ETH_URL"] = network.URLs[0]
	ocrForwarderEnvVars["ETH_HTTP_URL"] = network.HTTPURLs[0]
	ocrForwarderEnvVars["ETH_CHAIN_ID"] = fmt.Sprint(network.ChainID)

	baseEnvironmentConfig := &environment.Config{
		TTL: time.Hour * 720, // 30 days,
		NamespacePrefix: fmt.Sprintf(
			"soak-forwarder-ocr-%s",
			strings.ReplaceAll(strings.ToLower(network.Name), " ", "-"),
		),
		Test: t,
	}

	replicas := 6
	// Values you want each node to have the exact same of (e.g. eth_chain_id)
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
				"env": ocrForwarderEnvVars,
			}))
		} else {
			testEnvironment.AddHelm(chainlink.New(i, map[string]interface{}{
				"toml": client.AddNetworkDetailedConfig(ocrForwarderBaseTOML, ocrForwarderNetworkDetailTOML, network),
			}))
		}
	}

	err := testEnvironment.Run()
	require.NoError(t, err, "Error launching test environment")
	return testEnvironment, network

}
