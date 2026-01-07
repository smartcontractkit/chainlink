package ocr2

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/go-resty/resty/v2"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/products/ocr2"

	f "github.com/smartcontractkit/chainlink-testing-framework/framework"
)

var (
	L          = ocr2.L
	BlockEvery = 1 * time.Second

	TotalRoundsPerTestCount = int64(0)
	LatestRound             = int64(0)
	LatestRoundAnswer       = int64(0)
)

type chaosSettings struct {
	command          string
	recoveryWaitTime time.Duration
}

type gasSettings struct {
	gasPriceStart  *big.Int
	gasPriceBump   *big.Int
	rampSeconds    int
	holdSeconds    int
	releaseSeconds int
}

type roundSettings struct {
	value int
	gas   *gasSettings
	chaos *chaosSettings
}

type testcase struct {
	name               string
	roundCheckInterval time.Duration
	roundTimeout       time.Duration
	repeat             int
	roundSettings      []*roundSettings
	cfg                *ocr2.OCRv2SetConfigOptions
}

// simulateGasSpike is changing next block gas base fee in 3 steps: ramp, hold and release simulating a gas spike
func simulateGasSpike(t *testing.T, r *rpc.RPCClient, g *gasSettings) {
	currentGasPrice := g.gasPriceStart
	for i := 0; i < g.rampSeconds; i++ {
		err := r.PrintBlockBaseFee()
		require.NoError(t, err)
		t.Logf("Setting block base fee: %d", currentGasPrice)
		err = r.AnvilSetNextBlockBaseFeePerGas(currentGasPrice)
		require.NoError(t, err)
		currentGasPrice = currentGasPrice.Add(currentGasPrice, g.gasPriceBump)
		time.Sleep(BlockEvery)
	}
	for i := 0; i < g.holdSeconds; i++ {
		err := r.PrintBlockBaseFee()
		require.NoError(t, err)
		time.Sleep(BlockEvery)
		t.Logf("Setting block base fee: %d", currentGasPrice)
		err = r.AnvilSetNextBlockBaseFeePerGas(currentGasPrice)
		require.NoError(t, err)
	}
	for i := 0; i < g.releaseSeconds; i++ {
		err := r.PrintBlockBaseFee()
		require.NoError(t, err)
		time.Sleep(BlockEvery)
	}
}

// verifyRounds is a main test loop that applies EA deviations, chaos and verifier that eventually next round is still published on-chain
func verifyRounds(t *testing.T, in *de.Cfg, o2 *ocr2aggregator.OCR2Aggregator, tc testcase, c *rpc.RPCClient) {
	roundTicker := time.NewTicker(tc.roundCheckInterval)
	defer roundTicker.Stop()

	rounds := make([]struct {
		RoundId         *big.Int //nolint:revive // we can't change this field in generated binding
		Answer          *big.Int
		StartedAt       *big.Int
		UpdatedAt       *big.Int
		AnsweredInRound *big.Int
	}, 0)
	defer func() { TotalRoundsPerTestCount = 0 }()

	for {
		select {
		case <-time.After(tc.roundTimeout):
			L.Warn().Msgf("timeout reached, goal of %d rounds is not complete!", len(tc.roundSettings))
			return
		case <-roundTicker.C:
			L.Trace().
				Msg("checking for new rounds")
			currentRoundSettings := tc.roundSettings[TotalRoundsPerTestCount]

			rd, err := o2.LatestRoundData(&bind.CallOpts{})
			require.NoError(t, err)

			if rd.Answer.Int64() != LatestRoundAnswer {
				LatestRound = rd.RoundId.Int64()
				LatestRoundAnswer = rd.Answer.Int64()
				rounds = append(rounds, rd)
				L.Info().
					Int64("RoundID", rd.RoundId.Int64()).
					Int64("Answer", rd.Answer.Int64()).
					Msg("New round data")

				// apply next value deviation
				L.Info().
					Int("Value", currentRoundSettings.value).
					Msg("Settings new value for EA")
				r := resty.New().SetBaseURL(in.FakeServer.Out.BaseURLHost)
				_, err := r.R().Post(
					fmt.Sprintf(
						`/trigger_deviation?result=%d`, currentRoundSettings.value,
					),
				)
				// apply varios chaos experiments for next round
				if currentRoundSettings.gas != nil {
					L.Info().Msg("Creating gas spike")
					simulateGasSpike(t, c, currentRoundSettings.gas)
				}
				if currentRoundSettings.chaos != nil {
					L.Info().Msg("Executing chaos action")
					_, err = chaos.ExecPumba(
						currentRoundSettings.chaos.command,
						currentRoundSettings.chaos.recoveryWaitTime,
					)
					require.NoError(t, err)
				}
				require.NoError(t, err)
				TotalRoundsPerTestCount++
			}
			if len(rounds) == len(tc.roundSettings) {
				L.Info().
					Int64("LatestRound", LatestRound).
					Int("RequiredRounds", len(tc.roundSettings)).
					Int64("TotalRounds", TotalRoundsPerTestCount).
					Msg("All rounds are complete")
				return
			}
		}
	}
}

// checkResourceConsumption checks if resource consumption during tests is acceptable
func checkResourceConsumption(t *testing.T, in *de.Cfg, start, end time.Time, maxCPUTotalPercentage float64, maxMem int) {
	pc := f.NewPrometheusQueryClient(f.LocalPrometheusBaseURL)
	cpuResp, err := pc.Query("sum(rate(container_cpu_usage_seconds_total{name=~\".*don.*\"}[5m])) by (name) *100", end)
	require.NoError(t, err)
	cpu := f.ToLabelsMap(cpuResp)
	for i := 0; i < in.NodeSets[0].Nodes; i++ {
		nodeLabel := fmt.Sprintf("name:don-node%d", i)
		nodeCPU, cpuErr := strconv.ParseFloat(cpu[nodeLabel][0].(string), 64)
		L.Info().Int("Node", i).Float64("CPU", nodeCPU).Msg("CPU usage percentage")
		require.NoError(t, cpuErr)
		require.LessOrEqual(t, nodeCPU, maxCPUTotalPercentage)
	}
	memoryResp, err := pc.Query("sum(container_memory_rss{name=~\".*don.*\"}) by (name)", end)
	require.NoError(t, err)
	mem := f.ToLabelsMap(memoryResp)
	for i := 0; i < in.NodeSets[0].Nodes; i++ {
		nodeLabel := fmt.Sprintf("name:don-node%d", i)
		nodeMem, err := strconv.Atoi(mem[nodeLabel][0].(string))
		L.Info().Int("Node", i).Int("Memory", nodeMem).Msg("Total memory")
		require.NoError(t, err)
		require.LessOrEqual(t, nodeMem, maxMem)
	}
}
