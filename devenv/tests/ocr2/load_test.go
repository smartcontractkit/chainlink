package ocr2

import (
	"context"
	"fmt"
	"math/big"
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

func TestLoad(t *testing.T) {
	ctx := context.Background()
	outputFile := "../../env-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err)
	pdConfig, err := products.LoadOutput[ocr2.Configurator](outputFile)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})
	c, _, _, err := ocr2.ETHClient(ctx, in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, pdConfig.OCR2.GasSettings.FeeCapMultiplier, pdConfig.OCR2.GasSettings.TipCapMultiplier)
	require.NoError(t, err)
	clNodes, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err)

	anvilClient := rpc.New(in.Blockchains[0].Out.Nodes[0].ExternalHTTPUrl, nil)

	// this config must be as close to production as possible
	productionCfg := &ocr2.OCRv2SetConfigOptions{
		RMax:                                    3,
		DeltaProgress:                           20 * time.Second,
		DeltaResend:                             20 * time.Second,
		DeltaStage:                              15 * time.Second,
		MaxDurationInitialization:               5 * time.Second,
		MaxDurationQuery:                        5 * time.Second,
		MaxDurationObservation:                  5 * time.Second,
		MaxDurationReport:                       5 * time.Second,
		MaxDurationShouldAcceptFinalizedReport:  5 * time.Second,
		MaxDurationShouldTransmitAcceptedReport: 5 * time.Second,
	}

	testCases := []testcase{
		{
			name:               "clean",
			roundCheckInterval: 5 * time.Second,
			roundTimeout:       2 * time.Minute,
			repeat:             2,
			cfg:                productionCfg,
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
		{
			name:               "chaos",
			roundCheckInterval: 5 * time.Second,
			roundTimeout:       2 * time.Minute,
			repeat:             2,
			roundSettings: []*roundSettings{
				// these are just Pumba tool commands, read more here https://github.com/alexei-led/pumba
				{
					value: 1,
					chaos: &chaosSettings{
						command:          "stop --duration=10s --restart re2:don-node0",
						recoveryWaitTime: 10 * time.Second,
					},
				},
				{
					value: 1e3,
					chaos: &chaosSettings{
						command:          "netem --tc-image=gaiadocker/iproute2 --duration=10s delay --time=1000 re2:don-node.*",
						recoveryWaitTime: 10 * time.Second,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			o2, err := ocr2aggregator.NewOCR2Aggregator(common.HexToAddress(pdConfig.OCR2.DeployedContracts.OCRv2AggregatorAddr), c)
			require.NoError(t, err)
			L.Info().Any("Config", tc.cfg).Msg("Applying new OCR2 configuration")
			err = ocr2.UpdateOCR2ConfigOffChainValues(context.Background(), in.Blockchains[0], pdConfig.OCR2, o2, clNodes, tc.cfg)
			require.NoError(t, err)
			for range tc.repeat {
				verifyRounds(t, in, o2, tc, anvilClient)
			}
			checkResourceConsumption(t, in, start, time.Now(), 10.0, 400e6)
		})
	}
}
