package experiments

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
)

func TestGasExperiment(t *testing.T) {
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig("Soak", tc.OCR)
	require.NoError(t, err, "Error getting config")

	network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]
	readSethCfg := config.GetSethConfig()
	require.NotNil(t, readSethCfg, "Seth config shouldn't be nil")

	sethCfg, err := utils.MergeSethAndEvmNetworkConfigs(network, *readSethCfg)
	require.NoError(t, err, "Error merging seth and evm network configs")
	err = utils.ValidateSethNetworkConfig(sethCfg.Network)
	require.NoError(t, err, "Error validating seth network config")

	seth, err := seth.NewClientWithConfig(&sethCfg)
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
