package soak

import (
	"math/big"
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

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

type OcrSoakInputs struct {
	TestDuration         time.Duration `envconfig:"TEST_DURATION" default:"15m"`
	ChainlinkNodeFunding float64       `envconfig:"CHAINLINK_NODE_FUNDING" default:".1"`
	TimeBetweenRounds    time.Duration `envconfig:"TIME_BETWEEN_ROUNDS" default:"1m"`
}

func TestOCRSoak(t *testing.T) {
	soakNetwork := blockchain.LoadNetworkFromEnvironment()
	testEnvironment := environment.New(&environment.Config{InsideK8s: true})
	err := testEnvironment.
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: soakNetwork.Name,
			Simulated:   soakNetwork.Simulated,
			WsURLs:      soakNetwork.URLs,
		})).
		AddHelm(chainlink.New(0, nil)).
		AddHelm(chainlink.New(1, nil)).
		AddHelm(chainlink.New(2, nil)).
		AddHelm(chainlink.New(3, nil)).
		AddHelm(chainlink.New(4, nil)).
		AddHelm(chainlink.New(5, nil)).
		Run()
	require.NoError(t, err, "Error running soak environment")
	log.Info().Str("Namespace", testEnvironment.Cfg.Namespace).Msg("Connected to Soak Environment")

	chainClient, err := blockchain.NewEVMClient(soakNetwork, testEnvironment)
	require.NoError(t, err, "Error connecting to network")

	var testInputs OcrSoakInputs
	err = envconfig.Process("OCR", &testInputs)
	require.NoError(t, err, "Error reading OCR soak test inputs")

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
