package soak

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

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
	"github.com/smartcontractkit/chainlink/integration-tests/networks"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestMain(m *testing.M) {
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

	baseEnvironmentConfig := &environment.Config{
		TTL: time.Hour * 720, // 30 days,
		NamespacePrefix: fmt.Sprintf(
			"soak-forwarder-ocr-%s",
			strings.ReplaceAll(strings.ToLower(network.Name), " ", "-"),
		),
		Test: t,
	}

	cd, err := chainlink.NewDeployment(6, map[string]interface{}{
		"toml": client.AddNetworkDetailedConfig(ocrForwarderBaseTOML, ocrForwarderNetworkDetailTOML, network),
	})
	require.NoError(t, err, "Error creating chainlink deployment")
	// Values you want each node to have the exact same of (e.g. eth_chain_id)
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
	return testEnvironment, network
}
