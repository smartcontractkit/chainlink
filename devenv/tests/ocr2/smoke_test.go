package ocr2

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/ocr2"
)

func TestSmoke(t *testing.T) {
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	pdConfig, err := products.LoadOutput[ocr2.Configurator](outputFile)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})
	c, _, _, err := products.ETHClient(t.Context(), in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pdConfig.Config[0].GasSettings.FeeCapMultiplier, pdConfig.Config[0].GasSettings.TipCapMultiplier)
	require.NoError(t, err)
	clNodes, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	anvilClient := rpc.New(in.Blockchains[0].Out.Nodes[0].ExternalHTTPUrl, nil)

	o2, err := ocr2aggregator.NewOCR2Aggregator(common.HexToAddress(pdConfig.Config[0].DeployedContracts.OCRv2AggregatorAddr), c)
	require.NoError(t, err)
	L.Info().Any("Config", DefaultProductionOCR2Config).Msg("Applying new OCR2 configuration")
	err = ocr2.UpdateOCR2ConfigOffChainValues(
		context.Background(),
		in.Blockchains[0],
		pdConfig.Config[0],
		o2,
		clNodes,
		DefaultProductionOCR2Config,
	)
	require.NoError(t, err)
	verifyRounds(t, in, o2, testcase{
		name:               "rounds",
		roundCheckInterval: 5 * time.Second,
		roundTimeout:       2 * time.Minute,
		cfg:                DefaultProductionOCR2Config,
		roundSettings: []*roundSettings{
			{value: 1},
			{value: 1e3},
			{value: 1e5},
		},
	}, anvilClient)
}
