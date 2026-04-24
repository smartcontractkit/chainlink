package flux

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/flux"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSmoke(t *testing.T) {
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	productCfg, err := products.LoadOutput[flux.Configurator](outputFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	c, _, _, err := products.ETHClient(
		t.Context(),
		in.Blockchains[0].Out.Nodes[0].ExternalWSUrl,
		productCfg.Config[0].GasSettings.FeeCapMultiplier,
		productCfg.Config[0].GasSettings.TipCapMultiplier,
	)
	require.NoError(t, err)
	fluxAggregatorWrapper, err := flux_aggregator_wrapper.NewFluxAggregator(
		common.HexToAddress(productCfg.Config[0].DeployedContracts.FluxAggregator),
		c,
	)
	require.NoError(t, err)
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		lrd, err := fluxAggregatorWrapper.LatestRoundData(&bind.CallOpts{})
		assert.NoError(c, err)
		assert.Equal(c, int64(200), lrd.Answer.Int64())
	}, 2*time.Minute, 2*time.Second)
}
