package soak

import (
	"math/big"
	"testing"
	"time"

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
		if err := actions.TeardownRemoteSuite(ocrSoakTest.TearDownVals(t)); err != nil {
			log.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	ocrSoakTest.Setup(t, testEnvironment)
	log.Info().Msg("Set up soak test")
	ocrSoakTest.Run(t)
}
