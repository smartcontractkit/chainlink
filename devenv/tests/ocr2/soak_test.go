package ocr2

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/leak"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/ocr2"
)

func TestOCR2Soak(t *testing.T) {
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

	testCases := []testcase{
		{
			name:               "clean",
			roundCheckInterval: 5 * time.Second,
			roundTimeout:       2 * time.Minute,
			repeat:             60,
			cfg:                DefaultProductionOCR2Config,
			roundSettings: []*roundSettings{
				{value: 1},
				{value: 1e3},
				{value: 1e5},
				{value: 1e7},
				{value: 1e9},
			},
		},
		{
			name:               "gas spikes",
			roundCheckInterval: 5 * time.Second,
			roundTimeout:       2 * time.Minute,
			repeat:             2,
			roundSettings: []*roundSettings{
				{
					value: 1,
				},
				{
					value: 1e3,
					gas: &gasSettings{
						gasPriceStart:  big.NewInt(2e9),
						gasPriceBump:   big.NewInt(1e9),
						rampSeconds:    2,
						holdSeconds:    5,
						releaseSeconds: 2,
					},
				},
				{
					value: 1e5,
					gas: &gasSettings{
						gasPriceStart:  big.NewInt(2e9),
						gasPriceBump:   big.NewInt(5e9),
						rampSeconds:    2,
						holdSeconds:    5,
						releaseSeconds: 2,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			o2, err := ocr2aggregator.NewOCR2Aggregator(common.HexToAddress(pdConfig.Config[0].DeployedContracts.OCRv2AggregatorAddr), c)
			require.NoError(t, err)
			L.Info().Any("Config", tc.cfg).Msg("Applying new OCR2 configuration")
			err = ocr2.UpdateOCR2ConfigOffChainValues(t.Context(), in.Blockchains[0], pdConfig.Config[0], o2, clNodes, tc.cfg)
			require.NoError(t, err)
			for range tc.repeat {
				verifyRounds(t, in, o2, tc, anvilClient)
			}

			l, err := leak.NewCLNodesLeakDetector(leak.NewResourceLeakChecker())
			require.NoError(t, err)
			errs := l.Check(&leak.CLNodesCheck{
				// since the test is stable we assert absolute values
				// no more than 25% CPU and 350Mb (last 5m)
				ComparisonMode:  leak.ComparisonModeAbsolute,
				NumNodes:        in.NodeSets[0].Nodes,
				Start:           start,
				End:             time.Now(),
				WarmUpDuration:  30 * time.Minute,
				CPUThreshold:    25.0,
				MemoryThreshold: 350.0,
			})
			require.NoError(t, errs)
		})
	}
}
