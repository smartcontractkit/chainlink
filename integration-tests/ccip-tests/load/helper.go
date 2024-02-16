package load

import (
	"fmt"
	"math"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
)

type ChaosConfig struct {
	ChaosName        string
	ChaosFunc        chaos.ManifestFunc
	ChaosProps       *chaos.Props
	WaitBetweenChaos time.Duration
}

type LoadArgs struct {
	t                *testing.T
	lggr             zerolog.Logger
	schedules        []*wasp.Segment
	RunnerWg         *errgroup.Group // to wait on individual load generators run
	LoadStarterWg    *sync.WaitGroup // waits for all the runners to start
	TestCfg          *testsetups.CCIPTestConfig
	TestSetupArgs    *testsetups.CCIPTestSetUpOutputs
	ChaosExps        []ChaosConfig
	LoadgenTearDowns []func()
}

func (l *LoadArgs) Setup() {
	lggr := l.lggr
	existing := pointer.GetBool(l.TestCfg.TestGroupInput.ExistingDeployment)
	envName := "load-ccip"
	if existing {
		envName = "ccip-runner"
	}
	l.TestSetupArgs = testsetups.CCIPDefaultTestSetUp(l.TestCfg.Test, lggr, envName, nil, l.TestCfg)
}

func (l *LoadArgs) setSchedule() {
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

func (l *LoadArgs) SanityCheck() {
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

func (l *LoadArgs) TriggerLoadByLane() {
	l.setSchedule()
	l.TestSetupArgs.Reporter.SetDuration(l.TestCfg.TestGroupInput.TestDuration.Duration())
	namespace := l.TestCfg.TestGroupInput.TestRunName

	// start load for a lane
	startLoad := func(lane *actions.CCIPLane) {
		lane.Logger.Info().
			Str("Source Network", lane.SourceNetworkName).
			Str("Destination Network", lane.DestNetworkName).
			Msg("Starting load for lane")

		ccipLoad := NewCCIPLoad(l.TestCfg.Test, lane, l.TestCfg.TestGroupInput.PhaseTimeout.Duration(), 100000)
		ccipLoad.BeforeAllCall(l.TestCfg.TestGroupInput.MsgType, big.NewInt(*l.TestCfg.TestGroupInput.DestGasLimit))
		if lane.TestEnv != nil && lane.TestEnv.K8Env != nil && lane.TestEnv.K8Env.Cfg != nil {
			namespace = lane.TestEnv.K8Env.Cfg.Namespace
		}
		lokiConfig := l.TestCfg.EnvInput.Logging.Loki
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
			LokiConfig:            wasp.NewLokiConfig(lokiConfig.Endpoint, lokiConfig.TenantId, nil, nil),
			Labels: map[string]string{
				"test_group":   "load",
				"cluster":      "sdlc",
				"namespace":    namespace,
				"test_id":      "ccip",
				"source_chain": lane.SourceNetworkName,
				"dest_chain":   lane.DestNetworkName,
			},
			FailOnErr: true,
		})
		require.NoError(l.TestCfg.Test, err, "initiating loadgen for lane %s --> %s",
			lane.SourceNetworkName, lane.DestNetworkName)
		loadRunner.Run(false)
		l.AddToRunnerGroup(loadRunner)
	}
	for _, lane := range l.TestSetupArgs.Lanes {
		lane := lane
		l.LoadStarterWg.Add(1)
		go func() {
			defer l.LoadStarterWg.Done()
			startLoad(lane.ForwardLane)
		}()
		if pointer.GetBool(l.TestSetupArgs.Cfg.TestGroupInput.BiDirectionalLane) {
			l.LoadStarterWg.Add(1)
			go func() {
				defer l.LoadStarterWg.Done()
				startLoad(lane.ReverseLane)
			}()
		}
	}
}

func (l *LoadArgs) AddToRunnerGroup(gen *wasp.Generator) {
	l.RunnerWg.Go(func() error {
		_, failed := gen.Wait()
		if failed {
			return fmt.Errorf("load run is failed")
		}
		if len(gen.Errors()) > 0 {
			return fmt.Errorf("error in load sequence call %v", gen.Errors())
		}
		return nil
	})
}

func (l *LoadArgs) Wait() {
	l.lggr.Info().Msg("Waiting for load to start on all lanes")
	// wait for load runner to start
	l.LoadStarterWg.Wait()
	l.lggr.Info().Msg("Waiting for load to finish on all lanes")
	// wait for load runner to finish
	err := l.RunnerWg.Wait()
	require.NoError(l.t, err, "load run is failed")
	l.lggr.Info().Msg("Load finished on all lanes")
}

func (l *LoadArgs) ApplyChaos() {
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

func (l *LoadArgs) TearDown() {
	for _, tearDn := range l.LoadgenTearDowns {
		tearDn()
	}
	if l.TestSetupArgs.TearDown != nil {
		require.NoError(l.t, l.TestSetupArgs.TearDown())
	}
}

func (l *LoadArgs) TriggerLoadBySource() {
	require.NotNil(l.t, l.TestCfg.TestGroupInput.TestDuration, "test duration input is nil")
	require.GreaterOrEqual(l.t, 1, len(l.TestCfg.TestGroupInput.RequestPerUnitTime), "time unit input must be specified")
	l.TestSetupArgs.Reporter.SetDuration(l.TestCfg.TestGroupInput.TestDuration.Duration())
	namespace := l.TestCfg.TestGroupInput.TestRunName

	var laneBySource = make(map[string][]*actions.CCIPLane)
	for _, lane := range l.TestSetupArgs.Lanes {
		laneBySource[lane.ForwardLane.SourceNetworkName] = append(laneBySource[lane.ForwardLane.SourceNetworkName], lane.ForwardLane)
		if lane.ReverseLane != nil {
			laneBySource[lane.ReverseLane.SourceNetworkName] = append(laneBySource[lane.ReverseLane.SourceNetworkName], lane.ReverseLane)
		}
	}
	for source, lanes := range laneBySource {
		source := source
		lanes := lanes
		l.LoadStarterWg.Add(1)
		go func() {
			defer l.LoadStarterWg.Done()
			l.lggr.Info().
				Str("Source Network", source).
				Msg("Starting load for source")
			if lanes[0].TestEnv != nil && lanes[0].TestEnv.K8Env != nil && lanes[0].TestEnv.K8Env.Cfg != nil {
				namespace = lanes[0].TestEnv.K8Env.Cfg.Namespace
			}
			allLabels := map[string]string{
				"test_group":   "load",
				"cluster":      "sdlc",
				"namespace":    namespace,
				"test_id":      "ccip",
				"source_chain": source,
			}
			multiCallGen, err := NewMultiCallLoadGenerator(l.TestCfg, lanes, l.TestCfg.TestGroupInput.RequestPerUnitTime[0], allLabels)
			require.NoError(l.t, err)
			lokiConfig := l.TestCfg.EnvInput.Logging.Loki
			loadRunner, err := wasp.NewGenerator(&wasp.Config{
				T:                     l.TestCfg.Test,
				GenName:               fmt.Sprintf("Source %s", source),
				Schedule:              wasp.Plain(1, l.TestCfg.TestGroupInput.TestDuration.Duration()), // hardcoded request per unit time to 1 as we are using multiCallGen
				LoadType:              wasp.RPS,
				RateLimitUnitDuration: l.TestCfg.TestGroupInput.TimeUnit.Duration(),
				CallResultBufLen:      10, // we keep the last 10 call results for each generator, as the detailed report is generated at the end of the test
				CallTimeout:           (l.TestCfg.TestGroupInput.PhaseTimeout.Duration()) * 5,
				Gun:                   multiCallGen,
				Logger:                multiCallGen.logger,
				LokiConfig:            wasp.NewLokiConfig(lokiConfig.Endpoint, lokiConfig.TenantId, nil, nil),
				Labels:                allLabels,
				FailOnErr:             true,
			})
			require.NoError(l.TestCfg.Test, err, "initiating loadgen for source %s", source)
			loadRunner.Run(false)
			l.AddToRunnerGroup(loadRunner)
			l.LoadgenTearDowns = append(l.LoadgenTearDowns, func() {
				require.NoError(l.t, multiCallGen.Stop())
			})
		}()
	}
}

func NewLoadArgs(t *testing.T, lggr zerolog.Logger, chaosExps ...ChaosConfig) *LoadArgs {
	wg, _ := errgroup.WithContext(testcontext.Get(t))
	return &LoadArgs{
		t:             t,
		lggr:          lggr,
		RunnerWg:      wg,
		TestCfg:       testsetups.NewCCIPTestConfig(t, lggr, testconfig.Load),
		ChaosExps:     chaosExps,
		LoadStarterWg: &sync.WaitGroup{},
	}
}
