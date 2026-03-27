package ocr2

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/ocr2"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
)

func TestOCR2Chaos(t *testing.T) {
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

	anvilClient := rpc.New(in.Blockchains[0].Out.Nodes[0].ExternalHTTPUrl, nil)

	dtc, err := chaos.NewDockerChaos(t.Context())
	require.NoError(t, err)

	roundCheckInterval := 5 * time.Second
	roundTimeout := 2 * time.Minute
	chaosActionDuration := 30 * time.Second
	eaChaosDuration := 30 * time.Second
	defaultTwoRounds := []*roundSettings{{value: 1}, {value: 1e3}}
	anvilContainerName := "anvil-1337"

	testCases := []testcase{
		{
			name:               "rpc pause",
			roundCheckInterval: roundCheckInterval,
			roundTimeout:       roundTimeout,
			roundSettings:      defaultTwoRounds,
			repeat:             1,
			chaos: func() {
				err := dtc.Chaos(anvilContainerName, chaos.CmdPause, "")
				require.NoError(t, err)
				time.Sleep(chaosActionDuration)
				err = dtc.RemoveAll()
				require.NoError(t, err)
			},
		},
		{
			name:               "rpc latency spike",
			roundCheckInterval: roundCheckInterval,
			roundTimeout:       roundTimeout,
			roundSettings:      defaultTwoRounds,
			repeat:             1,
			chaos: func() {
				err := dtc.Chaos(anvilContainerName, chaos.CmdDelay, "3s")
				require.NoError(t, err)
				time.Sleep(chaosActionDuration)
				err = dtc.RemoveAll()
				require.NoError(t, err)
			},
		},
		{
			name:               "nodes mixed",
			roundCheckInterval: roundCheckInterval,
			roundTimeout:       roundTimeout,
			roundSettings:      defaultTwoRounds,
			repeat:             1,
			chaos: func() {
				err := dtc.Chaos("don-node1", chaos.CmdDelay, "1s")
				require.NoError(t, err)
				err = dtc.Chaos("don-node2", chaos.CmdLoss, "30%")
				require.NoError(t, err)
				err = dtc.Chaos("don-node3", chaos.CmdCorrupt, "30%")
				require.NoError(t, err)
				err = dtc.Chaos("don-node4", chaos.CmdDuplicate, "30%")
				require.NoError(t, err)
				time.Sleep(chaosActionDuration)
				err = dtc.RemoveAll()
				require.NoError(t, err)
			},
		},
		{
			name:               "nodes pause minority",
			roundCheckInterval: roundCheckInterval,
			roundTimeout:       roundTimeout,
			roundSettings:      defaultTwoRounds,
			repeat:             1,
			chaos: func() {
				err := dtc.Chaos("don-node1", chaos.CmdPause, "")
				require.NoError(t, err)
				err = dtc.Chaos("don-node2", chaos.CmdPause, "")
				require.NoError(t, err)
				time.Sleep(chaosActionDuration)
				err = dtc.RemoveAll()
				require.NoError(t, err)
			},
		},
		{
			name:               "nodes pause majority",
			roundCheckInterval: roundCheckInterval,
			roundTimeout:       roundTimeout,
			roundSettings:      defaultTwoRounds,
			repeat:             1,
			chaos: func() {
				err := dtc.Chaos("don-node1", chaos.CmdPause, "")
				require.NoError(t, err)
				err = dtc.Chaos("don-node2", chaos.CmdPause, "")
				require.NoError(t, err)
				err = dtc.Chaos("don-node3", chaos.CmdPause, "")
				require.NoError(t, err)
				time.Sleep(chaosActionDuration)
				err = dtc.RemoveAll()
				require.NoError(t, err)
			},
		},

		{
			name:               "pause ea",
			roundCheckInterval: roundCheckInterval,
			roundTimeout:       roundTimeout,
			roundSettings:      defaultTwoRounds,
			repeat:             1,
			chaos: func() {
				err := dtc.Chaos("fake", chaos.CmdPause, "")
				require.NoError(t, err)
				time.Sleep(eaChaosDuration)
				err = dtc.RemoveAll()
				require.NoError(t, err)
			},
		},
		{
			name:               "slow ea",
			roundCheckInterval: roundCheckInterval,
			roundSettings:      defaultTwoRounds,
			roundTimeout:       roundTimeout,
			chaos: func() {
				err := dtc.Chaos("fake", chaos.CmdDelay, "5s")
				require.NoError(t, err)
				time.Sleep(eaChaosDuration)
				err = dtc.RemoveAll()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o2, err := ocr2aggregator.NewOCR2Aggregator(common.HexToAddress(pdConfig.Config[0].DeployedContracts.OCRv2AggregatorAddr), c)
			require.NoError(t, err)
			tc.chaos()
			for range tc.repeat {
				verifyRounds(t, in, o2, tc, anvilClient)
			}
		})
	}
}
