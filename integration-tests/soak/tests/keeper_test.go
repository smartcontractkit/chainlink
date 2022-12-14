package soak

//revive:disable:dot-imports
import (
	"math/big"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

func TestKeeperSoak(t *testing.T) {
	soakNetwork := blockchain.LoadNetworkFromEnvironment()
	testEnvironment := environment.New(&environment.Config{InsideK8s: true})
	err := testEnvironment.
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
	require.NoError(t, err, "Error deploying soak environment")
	log.Info().Str("Namespace", testEnvironment.Cfg.Namespace).Msg("Connected to Soak Environment")

	chainClient, err := blockchain.NewEVMClient(soakNetwork, testEnvironment)
	require.NoError(t, err, "Connecting to blockchain nodes shouldn't fail")
	keeperBlockTimeTest := testsetups.NewKeeperBlockTimeTest(
		testsetups.KeeperBlockTimeTestInputs{
			BlockchainClient:  chainClient,
			NumberOfContracts: 5,
			KeeperRegistrySettings: &contracts.KeeperRegistrySettings{
				PaymentPremiumPPB:    uint32(200000000),
				FlatFeeMicroLINK:     uint32(0),
				BlockCountPerTurn:    big.NewInt(3),
				CheckGasLimit:        uint32(2500000),
				StalenessSeconds:     big.NewInt(90000),
				GasCeilingMultiplier: uint16(1),
				MinUpkeepSpend:       big.NewInt(0),
				MaxPerformGas:        uint32(5000000),
				FallbackGasPrice:     big.NewInt(2e11),
				FallbackLinkPrice:    big.NewInt(2e18),
			},
			CheckGasToBurn:       1,
			PerformGasToBurn:     1,
			BlockRange:           1000,
			BlockInterval:        50,
			ChainlinkNodeFunding: big.NewFloat(1),
		},
	)
	t.Cleanup(func() {
		if err := actions.TeardownRemoteSuite(keeperBlockTimeTest.TearDownVals(t)); err != nil {
			log.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	keeperBlockTimeTest.Setup(t, testEnvironment)
	keeperBlockTimeTest.Run(t)
}
