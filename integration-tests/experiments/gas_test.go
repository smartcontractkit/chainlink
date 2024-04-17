package experiments

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestGasExperiment(t *testing.T) {
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig("Soak", tc.OCR)
	require.NoError(t, err, "Error getting config")

	network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]
	seth, err := actions_seth.GetChainClient(&config, network)
	require.NoError(t, err, "Error creating seth client")

	_, err = actions_seth.SendFunds(l, seth, actions_seth.FundsToSendPayload{
		ToAddress:  seth.Addresses[0],
		Amount:     big.NewInt(10_000_000),
		PrivateKey: seth.PrivateKeys[0],
	})
	require.NoError(t, err, "Error sending funds")

	for i := 0; i < 1; i++ {
		_, err = contracts.DeployLinkTokenContract(l, seth)
		require.NoError(t, err, "Error deploying LINK contract")
		time.Sleep(2 * time.Second)
	}
}
