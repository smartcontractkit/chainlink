package load

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/rs/zerolog"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/chaos"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
)

type laneLoadCfg struct {
	lane *actions.CCIPLane
}

type ChaosConfig struct {
	ChaosName        string
	ChaosFunc        chaos.ManifestFunc
	ChaosProps       *chaos.Props
	WaitBetweenChaos time.Duration
}

type loadArgs struct {
	t             *testing.T
	lggr          zerolog.Logger
	ctx           context.Context
	ccipLoad      []*CCIPE2ELoad
	schedules     []*wasp.Segment
	loadRunner    []*wasp.Generator
	LaneLoadCfg   chan laneLoadCfg
	RunnerWg      *errgroup.Group // to wait on individual load generators run
	LoadStarterWg *sync.WaitGroup // waits for all the runners to start
	TestCfg       *testsetups.CCIPTestConfig
	TestSetupArgs *testsetups.CCIPTestSetUpOutputs
	ChaosExps     []ChaosConfig
}

func (l *loadArgs) Setup() {
	lggr := l.lggr
	existing := pointer.GetBool(l.TestCfg.TestGroupInput.ExistingDeployment)
	envName := "load-ccip"
	if existing {
		envName = "ccip-runner"
	}
	l.TestSetupArgs = testsetups.CCIPDefaultTestSetUp(l.TestCfg.Test, lggr, envName, nil, l.TestCfg)
}

func (l *loadArgs) setSchedule() {
	var segments []*wasp.Segment
	var segmentDuration time.Duration
	require.Greater(l.t, len(l.TestCfg.TestGroupInput.RequestPerUnitTime), 0, "RequestPerUnitTime must be set")
	if len(l.TestCfg.TestGroupInput.RequestPerUnitTime) > 1 {
		for i, req := range l.TestCfg.TestGroupInput.RequestPerUnitTime {
			duration := l.TestCfg.TestGroupInput.StepDuration[i].Duration()
			segmentDuration += duration
			segments = append(segments, wasp.Plain(req, duration)...)
		}
		totalDuration := l.TestCfg.TestGroupInput.TestDuration.Duration()
		repeatTimes := totalDuration.Seconds() / segmentDuration.Seconds()
		l.schedules = wasp.CombineAndRepeat(int(math.Round(repeatTimes)), segments)
	} else {
		l.schedules = wasp.Plain(l.TestCfg.TestGroupInput.RequestPerUnitTime[0], l.TestCfg.TestGroupInput.TestDuration.Duration())
	}
}

func (l *loadArgs) SanityCheck() {
	for _, lane := range l.TestSetupArgs.Lanes {
		lane.ForwardLane.RecordStateBeforeTransfer()
		err := lane.ForwardLane.SendRequests(1, l.TestCfg.TestGroupInput.MsgType, big.NewInt(600_000))
		require.NoError(l.t, err)
		lane.ForwardLane.ValidateRequests(true)
		lane.ReverseLane.RecordStateBeforeTransfer()
		err = lane.ReverseLane.SendRequests(1, l.TestCfg.TestGroupInput.MsgType, big.NewInt(600_000))
		require.NoError(l.t, err)
		lane.ReverseLane.ValidateRequests(true)
	}
}

func (l *loadArgs) TriggerLoad(schedule ...*wasp.Segment) {
	l.Start()
	if len(schedule) == 0 {
		l.setSchedule()
	} else {
		l.schedules = schedule
	}
	for _, lane := range l.TestSetupArgs.Lanes {
		if lane.LaneDeployed {
			l.LaneLoadCfg <- laneLoadCfg{
				lane: lane.ForwardLane,
			}
			if lane.ReverseLane != nil {
				l.LaneLoadCfg <- laneLoadCfg{
					lane: lane.ReverseLane,
				}
			}
		}
	}
	l.TestSetupArgs.Reporter.SetDuration(l.TestCfg.TestGroupInput.TestDuration.Duration())
}

func (l *loadArgs) AddMoreLanesToRun() {
	require.Len(l.t, l.TestSetupArgs.Lanes, 1, "lane for first network pair should be deployed already")
	if len(l.TestSetupArgs.Lanes) == len(l.TestCfg.NetworkPairs) {
		l.lggr.Info().Msg("All lanes are already deployed, no need to add more lanes")
		return
	}
	transferAmounts := []*big.Int{big.NewInt(1)}
	// set the ticker duration based on number of network pairs and the total test duration
	noOfPair := int64(len(l.TestCfg.NetworkPairs))
	step := l.TestCfg.TestGroupInput.TestDuration.Duration().Nanoseconds() / noOfPair
	ticker := time.NewTicker(time.Duration(step))
	l.setSchedule()
	// Lane for the first network pair is already deployed
	netIndex := 1
	for {
		select {
		case <-ticker.C:
			n := l.TestCfg.NetworkPairs[netIndex]
			l.lggr.Info().
				Str("Network 1", n.NetworkA.Name).
				Str("Network 2", n.NetworkB.Name).
				Msg("Adding lanes for network pair")
			err := l.TestSetupArgs.AddLanesForNetworkPair(
				l.lggr, n.NetworkA, n.NetworkB,
				n.ChainClientA, n.ChainClientB,
				transferAmounts, 5, true,
				true)
			assert.NoError(l.t, err)
			l.LaneLoadCfg <- laneLoadCfg{
				lane: l.TestSetupArgs.Lanes[netIndex].ForwardLane,
			}
			if l.TestSetupArgs.Lanes[netIndex].ReverseLane != nil {
				l.LaneLoadCfg <- laneLoadCfg{
					lane: l.TestSetupArgs.Lanes[netIndex].ReverseLane,
				}
			}
			netIndex++
			if netIndex >= len(l.TestCfg.NetworkPairs) {
				ticker.Stop()
				return
			}
		}
	}
}

// Start polls the LaneLoadCfg channel for new lanes and starts the load runner.
// LaneLoadCfg channel should receive a lane whenever the deployment is complete.
func (l *loadArgs) Start() {
	l.LoadStarterWg.Add(1)
	waitForLoadRun := func(gen *wasp.Generator, ccipLoad *CCIPE2ELoad) error {
		_, failed := gen.Wait()
		if failed {
			return fmt.Errorf("load run is failed")
		}
		if len(gen.Errors()) > 0 {
			return fmt.Errorf("error in load sequence call %v", gen.Errors())
		}
		return nil
	}
	go func() {
		defer l.LoadStarterWg.Done()
		loadCount := 0
		namespace := l.TestCfg.TestGroupInput.ExistingEnv
		for {
			select {
			case cfg := <-l.LaneLoadCfg:
				loadCount++
				lane := cfg.lane
				l.lggr.Info().
					Str("Source Network", lane.SourceNetworkName).
					Str("Destination Network", lane.DestNetworkName).
					Msg("Starting load for lane")

				ccipLoad := NewCCIPLoad(l.TestCfg.Test, lane, l.TestCfg.TestGroupInput.PhaseTimeout.Duration(), 100000)
				ccipLoad.BeforeAllCall(l.TestCfg.TestGroupInput.MsgType)
				if lane.TestEnv != nil && lane.TestEnv.K8Env != nil && lane.TestEnv.K8Env.Cfg != nil {
					namespace = lane.TestEnv.K8Env.Cfg.Namespace
				}

				loadRunner, err := wasp.NewGenerator(&wasp.Config{
					T:                     l.TestCfg.Test,
					GenName:               fmt.Sprintf("lane %s-> %s", lane.SourceNetworkName, lane.DestNetworkName),
					Schedule:              l.schedules,
					LoadType:              wasp.RPS,
					RateLimitUnitDuration: l.TestCfg.TestGroupInput.TimeUnit.Duration(),
					CallResultBufLen:      10, // we keep the last 10 call results for each generator, as the detailed report is generated at the end of the test
					CallTimeout:           (l.TestCfg.TestGroupInput.PhaseTimeout.Duration()) * 5,
					Gun:                   ccipLoad,
					Logger:                ccipLoad.Lane.Logger,
					SharedData:            l.TestCfg.TestGroupInput.MsgType,
					LokiConfig:            wasp.NewEnvLokiConfig(),
					Labels: map[string]string{
						"test_group":   "load",
						"cluster":      "sdlc",
						"namespace":    namespace,
						"test_id":      "ccip",
						"source_chain": lane.SourceNetworkName,
						"dest_chain":   lane.DestNetworkName,
					},
				})
				require.NoError(l.TestCfg.Test, err, "initiating loadgen for lane %s --> %s",
					lane.SourceNetworkName, lane.DestNetworkName)
				loadRunner.Run(false)
				l.ccipLoad = append(l.ccipLoad, ccipLoad)
				l.loadRunner = append(l.loadRunner, loadRunner)
				l.RunnerWg.Go(func() error {
					return waitForLoadRun(loadRunner, ccipLoad)
				})
				if loadCount == len(l.TestCfg.NetworkPairs)*2 {
					l.lggr.Info().Msg("load is running for all lanes now")
					return
				}
			}
		}
	}()
}

func (l *loadArgs) Wait() {
	// wait for load runner to start on all lanes
	l.LoadStarterWg.Wait()
	l.lggr.Info().Msg("Waiting for load to finish on all lanes")
	// wait for load runner to finish
	err := l.RunnerWg.Wait()
	require.NoError(l.t, err, "load run is failed")
}

func (l *loadArgs) ApplyChaos() {
	testEnv := l.TestSetupArgs.Env
	if testEnv == nil || testEnv.K8Env == nil {
		l.lggr.Warn().Msg("test environment is nil, skipping chaos")
		return
	}
	testEnv.ChaosLabelForCLNodes(l.TestCfg.Test)

	for _, exp := range l.ChaosExps {
		if exp.WaitBetweenChaos > 0 {
			l.lggr.Info().Msgf("sleeping for %s after chaos %s", exp.WaitBetweenChaos, exp.ChaosName)
			time.Sleep(exp.WaitBetweenChaos)
		}
		l.lggr.Info().Msgf("Starting to apply chaos %s at %s", exp.ChaosName, time.Now().UTC())
		// apply chaos
		chaosId, err := testEnv.K8Env.Chaos.Run(exp.ChaosFunc(testEnv.K8Env.Cfg.Namespace, exp.ChaosProps))
		require.NoError(l.t, err)
		if chaosId != "" {
			chaosDur, err := time.ParseDuration(exp.ChaosProps.DurationStr)
			require.NoError(l.t, err)
			err = testEnv.K8Env.Chaos.WaitForAllRecovered(chaosId, chaosDur+1*time.Minute)
			require.NoError(l.t, err)
			l.lggr.Info().Msgf("chaos %s is recovered at %s", exp.ChaosName, time.Now().UTC())
			err = testEnv.K8Env.Chaos.Stop(chaosId)
			require.NoError(l.t, err)
			l.lggr.Info().Msgf("stopped chaos %s at %s", exp.ChaosName, time.Now().UTC())
		}
	}
}

func (l *loadArgs) TearDown() {
	if l.TestSetupArgs.TearDown != nil {
		require.NoError(l.t, l.TestSetupArgs.TearDown())
	}
}

func NewLoadArgs(t *testing.T, lggr zerolog.Logger, parent context.Context, chaosExps ...ChaosConfig) *loadArgs {
	wg, ctx := errgroup.WithContext(parent)
	return &loadArgs{
		t:             t,
		lggr:          lggr,
		RunnerWg:      wg,
		ctx:           ctx,
		TestCfg:       testsetups.NewCCIPTestConfig(t, lggr, testconfig.Load),
		LaneLoadCfg:   make(chan laneLoadCfg),
		LoadStarterWg: &sync.WaitGroup{},
		ChaosExps:     chaosExps,
	}
}
