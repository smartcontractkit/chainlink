package congestion

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

type Config struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
	Seth        *seth.Config      `toml:"Seth"`
	Congestion  *Input            `toml:"congestion"`
}

func TestSimulator(t *testing.T) {
	t.Parallel()
	cfg, err := framework.Load[Config](t)
	require.NoError(t, err)
	bc, err := blockchain.NewBlockchainNetwork(cfg.BlockchainA)
	require.NoError(t, err)
	require.NotEmpty(t, bc.Nodes[0].HostHTTPUrl)
	cfg.Seth.Network.URLs = []string{bc.Nodes[0].HostWSUrl}
	client, err := seth.NewClientWithConfig(cfg.Seth)
	require.NoError(t, err)
	anvilClient := rpc.New(bc.Nodes[0].HostHTTPUrl, map[string][]string{})
	lggr, _ := logger.TestObserved(t, zapcore.DebugLevel)
	sim, err := NewSimulator(t, *cfg.Congestion, client, anvilClient, logger.Sugared(lggr))
	require.NoError(t, err)
	observations := map[phaseName][]chainState{}
	sim.observationsHandler = func(state chainState) {
		observations[state.PhaseName] = append(observations[state.PhaseName], state) // handle is not called concurrently
	}
	err = sim.Start(tests.Context(t))
	require.NoError(t, err)
	// wait for 2 minutes to collect observations through all phases
	time.Sleep(time.Minute * 2)
	require.NoError(t, sim.Close())
	require.Len(t, observations, 4, "expected all 4 phases to be observed")
	requireCongestionWeakEqual(t, "inactive", observations[phaseInactive], 0) // during inactivity of simulation fees might fluctuate according to eip1559, so we can't check them
	requirePhaseObservations(t, "rampUp", observations[phaseRampUp], cfg.Congestion.RampUp.Congestion, 0, cfg.Congestion.FeesIncreasePercent)
	requirePhaseObservations(t, "plateau", observations[phasePlateau], cfg.Congestion.Plateau.Congestion, cfg.Congestion.FeesIncreasePercent, cfg.Congestion.FeesIncreasePercent)
	requirePhaseObservations(t, "cool down", observations[phaseCoolDown], cfg.Congestion.CoolDown.Congestion, 0, cfg.Congestion.FeesIncreasePercent)
}

func requirePhaseObservations(t *testing.T, name string, observations []chainState, expectedCongestion, minDelta, maxDelta float64) {
	requireCongestionWeakEqual(t, name, observations, expectedCongestion)
	requireValueInRange(t, name+" tipPercent", observations, func(state chainState) float64 {
		return state.TipDeltaPercent
	}, minDelta, maxDelta)
	requireValueInRange(t, name+" baseFeePercent", observations, func(state chainState) float64 {
		return state.BaseFeeDeltaPercent
	}, minDelta, maxDelta)
}

func requireCongestionWeakEqual(t *testing.T, name string, observations []chainState, target float64) {
	vals := sanitizeValues(t, observations, func(state chainState) float64 {
		return state.Congestion
	})

	total := float64(0)
	for _, val := range vals {
		total += val
	}
	avg := total / float64(len(vals))
	require.GreaterOrEqual(t, avg, target*0.95, "%s avg congestion is not within expected range %v", name, vals)
	require.LessOrEqual(t, avg, target*1.05, "%s avg congestion is not within expected range %v", name, vals)
}

func requireValueInRange(t *testing.T, name string, observations []chainState, extract func(state chainState) float64, min, max float64) {
	vals := sanitizeValues(t, observations, extract)
	require.GreaterOrEqual(t, vals[0], min, "%s is not within expected range", name)
	require.LessOrEqual(t, vals[len(vals)-1], max, "%s is not within expected range", name)
}

func sanitizeValues(t *testing.T, observations []chainState, extract func(state chainState) float64) []float64 {
	vals := make([]float64, 0, len(observations))
	for _, obs := range observations {
		vals = append(vals, extract(obs))
	}

	sort.Float64s(vals)
	if len(vals) < 10 {
		t.Fatalf("expected at least 10 observations for each simulation state, got %d", len(vals))
	}

	// trim 20% to cut outliers
	trim := len(vals) / 10
	vals = vals[trim : len(vals)-trim]
	return vals
}
