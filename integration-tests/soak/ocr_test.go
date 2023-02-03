package soak

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
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
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestOCRSoak(t *testing.T) {
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
			log.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	ocrSoakTest.Setup(t, testEnvironment)
	log.Info().Msg("Set up soak test")
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

	replicas := 6
	baseTOML := `[OCR]
Enabled = true

[P2P]
[P2P.V1]
Enabled = true
ListenIP = '0.0.0.0'
ListenPort = 6690`
	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		}))
	for i := 0; i < replicas; i++ {
		testEnvironment.AddHelm(chainlink.New(i, map[string]interface{}{
			"toml": client.AddNetworksConfig(baseTOML, network),
		}))
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
	os.Setenv("TEST_OCR_CHAINLINK_NODE_FUNDING", fmt.Sprintf("%f", i.ChainlinkNodeFunding))
	os.Setenv("TEST_OCR_TIME_BETWEEN_ROUNDS", i.TimeBetweenRounds.String())
}
