package load

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
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

// WaspSchedule calculates the load schedule based on the provided request per unit time and duration
// if multiple step durations are provided, it will calculate the schedule based on the step duration and
// corresponding request per unit time by matching the index of the request per unit time and step duration slice
func WaspSchedule(rps []int64, duration *config.Duration, steps []*config.Duration) []*wasp.Segment {
	var segments []*wasp.Segment
	var segmentDuration time.Duration

	if len(rps) > 1 {
		for i, req := range rps {
			duration := steps[i].Duration()
			segmentDuration += duration
			segments = append(segments, wasp.Plain(req, duration)...)
		}
		totalDuration := duration.Duration()
		repeatTimes := totalDuration.Seconds() / segmentDuration.Seconds()
		return wasp.CombineAndRepeat(int(math.Round(repeatTimes)), segments)
	}
	return wasp.Plain(rps[0], duration.Duration())
}

type LoadArgs struct {
	t                *testing.T
	Ctx              context.Context
	lggr             *zerolog.Logger
	RunnerWg         *errgroup.Group // to wait on individual load generators run
	LoadStarterWg    *sync.WaitGroup // waits for all the runners to start
	TestCfg          *testsetups.CCIPTestConfig
	TestSetupArgs    *testsetups.CCIPTestSetUpOutputs
	ChaosExps        []ChaosConfig
	LoadgenTearDowns []func()
	Labels           map[string]string
	pauseLoad        *atomic.Bool
}

func (l *LoadArgs) SetReportParams() {
	var qParams []string
	for k, v := range l.Labels {
		qParams = append(qParams, fmt.Sprintf("var-%s=%s", k, v))
	}
	// add one of the source and destination network to the grafana query params
	if len(l.TestSetupArgs.Lanes) > 0 {
		qParams = append(qParams, fmt.Sprintf("var-source_chain=%s", l.TestSetupArgs.Lanes[0].ForwardLane.SourceNetworkName))
		qParams = append(qParams, fmt.Sprintf("var-dest_chain=%s", l.TestSetupArgs.Lanes[0].ForwardLane.DestNetworkName))
	}
	err := l.TestSetupArgs.Reporter.AddToGrafanaDashboardQueryParams(qParams...)
	require.NoError(l.t, err, "failed to set grafana query params")
}

func (l *LoadArgs) Setup() {
	lggr := l.lggr
	existing := pointer.GetBool(l.TestCfg.TestGroupInput.ExistingDeployment)
	envName := "load-ccip"
	if existing {
		envName = "ccip-runner"
	}
	l.TestSetupArgs = testsetups.CCIPDefaultTestSetUp(l.TestCfg.Test, lggr, envName, nil, l.TestCfg)
	namespace := l.TestCfg.TestGroupInput.LoadProfile.TestRunName
	if l.TestSetupArgs.Env != nil && l.TestSetupArgs.Env.K8Env != nil && l.TestSetupArgs.Env.K8Env.Cfg != nil {
		namespace = l.TestSetupArgs.Env.K8Env.Cfg.Namespace
	}
	l.Labels = map[string]string{
		"test_group": "load",
		"test_id":    "ccip",
		"namespace":  namespace,
	}
	l.TestSetupArgs.Reporter.SetGrafanaURLProvider(l.TestCfg.EnvInput)
	l.SetReportParams()
}

func (l *LoadArgs) scheduleForDest(destNetworkName string) []*wasp.Segment {
	require.Greater(l.t, len(l.TestCfg.TestGroupInput.LoadProfile.RequestPerUnitTime), 0, "RequestPerUnitTime must be set")
	// try to locate if there is a frequency provided for the destination network
	// to locate the frequency, we check if the destination network name contains the network name in the frequency map
	// if found, use that frequency for the destination network
	// otherwise, use the default frequency
	if l.TestCfg.TestGroupInput.LoadProfile.FrequencyByDestination != nil {
		l.lggr.Debug().
			Interface("FrequencyByDestination", l.TestCfg.TestGroupInput.LoadProfile.FrequencyByDestination).
			Msg("LoadProfile provided")
		for networkName, freq := range l.TestCfg.TestGroupInput.LoadProfile.FrequencyByDestination {
			l.lggr.Debug().Str("Destination", destNetworkName).Str("NetworkName", networkName).Msg("Checking frequency for destination")
			if strings.EqualFold(destNetworkName, networkName) {
				l.lggr.Info().Str("Destination", destNetworkName).
					Ints64("RequestPerUnitTime", freq.RequestPerUnitTime).
					Msg("Using frequency for destination")
				return WaspSchedule(
					freq.RequestPerUnitTime,
					l.TestCfg.TestGroupInput.LoadProfile.TestDuration,
					freq.StepDuration)
			}
		}
	}

	return WaspSchedule(
		l.TestCfg.TestGroupInput.LoadProfile.RequestPerUnitTime,
		l.TestCfg.TestGroupInput.LoadProfile.TestDuration,
		l.TestCfg.TestGroupInput.LoadProfile.StepDuration)
}

func (l *LoadArgs) SanityCheck() {
	var allLanes []*actions.CCIPLane
	for _, lane := range l.TestSetupArgs.Lanes {
		allLanes = append(allLanes, lane.ForwardLane)
		if lane.ReverseLane != nil {
			allLanes = append(allLanes, lane.ReverseLane)
		}
	}
	for _, lane := range allLanes {
		ccipLoad := NewCCIPLoad(
			l.TestCfg.Test, lane,
			l.TestCfg.TestGroupInput.PhaseTimeout.Duration(),
			1, l.TestCfg.TestGroupInput.LoadProfile.MsgProfile,
			0, nil,
		)
		ccipLoad.BeforeAllCall()
		resp := ccipLoad.Call(nil)
		require.False(l.t, resp.Failed, "request failed in sanity check")
	}
}

// ValidateCurseFollowedByUncurse assumes the lanes under test are bi-directional.
// It assumes requests in both direction are in flight when this is called.
// It assumes the ARM is not already cursed, it will fail the test if it is in cursed state.
// It curses source ARM for forward lanes so that destination curse is also validated for reverse lanes.
// It waits for 2 minutes for curse to be seen by ccip plugins and contracts.
// It captures the curse timestamp to verify no execution state changed event is emitted after the cure is applied.
// It uncurses the source ARM at the end so that it can be verified that rest of the requests are processed as expected.
// Validates that even after uncursing the lane should not function for 30 more minutes.
func (l *LoadArgs) ValidateCurseFollowedByUncurse() {
	var lanes []*actions.CCIPLane
	for _, lane := range l.TestSetupArgs.Lanes {
		lanes = append(lanes, lane.ForwardLane)
	}
	// check if source is already cursed
	for _, lane := range lanes {
		cursed, err := lane.Source.Common.IsCursed()
		require.NoError(l.t, err, "cannot get cursed state")
		if cursed {
			require.Fail(l.t, "test will not work if ARM is already cursed")
		}
	}
	// before cursing set pause
	l.pauseLoad.Store(true)
	// wait for some time for pause to be active in wasp
	l.lggr.Info().Msg("Waiting for 1 minute after applying pause on load")
	time.Sleep(1 * time.Minute)
	curseTimeStamps := make(map[string]time.Time)
	for _, lane := range lanes {
		if _, exists := curseTimeStamps[lane.SourceNetworkName]; exists {
			continue
		}
		curseTx, err := lane.Source.Common.CurseARM()
		require.NoError(l.t, err, "error in cursing arm")
		require.NotNil(l.t, curseTx, "invalid cursetx")
		receipt, err := lane.Source.Common.ChainClient.GetTxReceipt(curseTx.Hash())
		require.NoError(l.t, err)
		hdr, err := lane.Source.Common.ChainClient.HeaderByNumber(context.Background(), receipt.BlockNumber)
		require.NoError(l.t, err)
		curseTimeStamps[lane.SourceNetworkName] = hdr.Timestamp
		l.lggr.Info().Str("Source", lane.SourceNetworkName).Msg("Curse is applied on source")
		l.lggr.Info().Str("Destination", lane.SourceNetworkName).Msg("Curse is applied on destination")
	}

	l.lggr.Info().Msg("Curse is applied on all lanes. Waiting for 2 minutes")
	time.Sleep(2 * time.Minute)

	for _, lane := range lanes {
		// try to send requests on lanes on which curse is applied on source RMN and the request should revert
		// data-only transfer is sufficient
		lane.Source.TransferAmount = []*big.Int{}
		failedTx, _, _, _, err := lane.Source.SendRequest(lane.Dest.ReceiverDapp.EthAddress, big.NewInt(actions.DefaultDestinationGasLimit))
		if lane.Source.Common.ChainClient.GetNetworkConfig().MinimumConfirmations > 0 {
			require.Error(l.t, err)
		} else {
			require.NoError(l.t, err)
		}
		errReason, v, err := lane.Source.Common.ChainClient.RevertReasonFromTx(failedTx, router.RouterABI)
		require.NoError(l.t, err)
		require.Equal(l.t, "BadARMSignal", errReason)
		lane.Logger.Info().
			Str("Revert Reason", errReason).
			Interface("Args", v).
			Str("FailedTx", failedTx.Hex()).
			Msg("Msg sent while source ARM is cursed")
	}

	// now uncurse all
	for _, lane := range lanes {
		require.NoError(l.t, lane.Source.Common.UnvoteToCurseARM(), "error to unvote in cursing arm")
	}
	l.lggr.Info().Msg("Curse is lifted on all lanes")
	// lift the pause on load test
	l.pauseLoad.Store(false)

	// now add the reverse lanes so that destination curse is also verified
	// we add the reverse lanes now to verify absence of commit and execution for the reverse lanes
	for _, lane := range l.TestSetupArgs.Lanes {
		lanes = append(lanes, lane.ReverseLane)
	}

	// verify that even after uncursing the lane should not function for 30 more minutes,
	// i.e no execution state changed or commit report accepted event is generated
	errGrp := &errgroup.Group{}
	for _, lane := range lanes {
		lane := lane
		curseTimeStamp, exists := curseTimeStamps[lane.SourceNetworkName]
		// if curse timestamp does not exist for source, it will exist for destination
		if !exists {
			curseTimeStamp, exists = curseTimeStamps[lane.DestNetworkName]
			require.Truef(l.t, exists, "did not find curse time stamp for lane %s->%s", lane.SourceNetworkName, lane.DestNetworkName)
		}
		errGrp.Go(func() error {
			lane.Logger.Info().Msg("Validating no CommitReportAccepted event is received for 29 minutes")
			// we allow additional 1 minute after curse timestamp for curse to be visible by plugin
			return lane.Dest.AssertNoReportAcceptedEventReceived(lane.Logger, 25*time.Minute, curseTimeStamp.Add(1*time.Minute))
		})
		errGrp.Go(func() error {
			lane.Logger.Info().Msg("Validating no ExecutionStateChanged event is received for 25 minutes")
			// we allow additional 1 minute after curse timestamp for curse to be visible by plugin
			return lane.Dest.AssertNoExecutionStateChangedEventReceived(lane.Logger, 25*time.Minute, curseTimeStamp.Add(1*time.Minute))
		})
	}
	l.lggr.Info().Msg("waiting for no commit/execution validation")
	err := errGrp.Wait()
	require.NoError(l.t, err, "error received to validate no commit/execution is generated after lane is cursed")
}

func (l *LoadArgs) TriggerLoadByLane() {
	l.TestSetupArgs.Reporter.SetDuration(l.TestCfg.TestGroupInput.LoadProfile.TestDuration.Duration())

	// start load for a lane
	startLoad := func(lane *actions.CCIPLane) {
		lane.Logger.Info().
			Str("Source Network", lane.SourceNetworkName).
			Str("Destination Network", lane.DestNetworkName).
			Msg("Starting load for lane")
		sendMaxData := pointer.GetInt64(l.TestCfg.TestGroupInput.LoadProfile.SendMaxDataInEveryMsgCount)
		ccipLoad := NewCCIPLoad(
			l.TestCfg.Test, lane, l.TestCfg.TestGroupInput.PhaseTimeout.Duration(),
			100000, l.TestCfg.TestGroupInput.LoadProfile.MsgProfile, sendMaxData,
			l.TestCfg.TestGroupInput.SkipRequestIfAnotherRequestTriggeredWithin,
		)
		ccipLoad.BeforeAllCall()
		// if it's not multicall set the tokens to nil to free up some space,
		// we have already formed the msg to be sent in load, there is no need to store the bridge tokens anymore
		// In case of multicall we still need the BridgeTokens to transfer amount from mutlicall to owner
		if !lane.Source.Common.MulticallEnabled {
			lane.Source.Common.BridgeTokens = nil
			lane.Dest.Common.BridgeTokens = nil
		}
		// no need for price registry in load
		lane.Source.Common.PriceRegistry = nil
		lane.Dest.Common.PriceRegistry = nil
		lokiConfig := l.TestCfg.EnvInput.Logging.Loki
		labels := make(map[string]string)
		for k, v := range l.Labels {
			labels[k] = v
		}
		labels["source_chain"] = lane.SourceNetworkName
		labels["dest_chain"] = lane.DestNetworkName
		waspCfg := &wasp.Config{
			T:                     l.TestCfg.Test,
			GenName:               fmt.Sprintf("lane %s-> %s", lane.SourceNetworkName, lane.DestNetworkName),
			Schedule:              l.scheduleForDest(lane.DestNetworkName),
			LoadType:              wasp.RPS,
			RateLimitUnitDuration: l.TestCfg.TestGroupInput.LoadProfile.TimeUnit.Duration(),
			CallResultBufLen:      10, // we keep the last 10 call results for each generator, as the detailed report is generated at the end of the test
			CallTimeout:           (l.TestCfg.TestGroupInput.PhaseTimeout.Duration()) * 5,
			Gun:                   ccipLoad,
			Logger:                *ccipLoad.Lane.Logger,
			LokiConfig:            wasp.NewLokiConfig(lokiConfig.Endpoint, lokiConfig.TenantId, nil, nil),
			Labels:                labels,
			FailOnErr:             pointer.GetBool(l.TestCfg.TestGroupInput.LoadProfile.FailOnFirstErrorInLoad),
		}
		waspCfg.LokiConfig.Timeout = time.Minute
		loadRunner, err := wasp.NewGenerator(waspCfg)
		require.NoError(l.TestCfg.Test, err, "error while initiating loadgen for lane %s --> %s",
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
	// watch for pause signal
	go func(gen *wasp.Generator) {
		ticker := time.NewTicker(time.Second)
		pausedOnce := false
		resumedAlready := false
		for {
			select {
			case <-ticker.C:
				if l.pauseLoad.Load() && !pausedOnce {
					gen.Pause()
					pausedOnce = true
					continue
				}
				if pausedOnce && !resumedAlready && !l.pauseLoad.Load() {
					gen.Resume()
					resumedAlready = true
				}
			case <-l.Ctx.Done():
				return
			}
		}
	}(gen)
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
	require.NotNil(l.t, l.TestCfg.TestGroupInput.LoadProfile.TestDuration, "test duration input is nil")
	require.GreaterOrEqual(l.t, 1, len(l.TestCfg.TestGroupInput.LoadProfile.RequestPerUnitTime), "time unit input must be specified")
	l.TestSetupArgs.Reporter.SetDuration(l.TestCfg.TestGroupInput.LoadProfile.TestDuration.Duration())
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
			allLabels := make(map[string]string)
			for k, v := range l.Labels {
				allLabels[k] = v
			}
			allLabels["source_chain"] = source
			multiCallGen, err := NewMultiCallLoadGenerator(l.TestCfg, lanes, l.TestCfg.TestGroupInput.LoadProfile.RequestPerUnitTime[0], allLabels)
			require.NoError(l.t, err)
			lokiConfig := l.TestCfg.EnvInput.Logging.Loki
			loadRunner, err := wasp.NewGenerator(&wasp.Config{
				T:                     l.TestCfg.Test,
				GenName:               fmt.Sprintf("Source %s", source),
				Schedule:              wasp.Plain(1, l.TestCfg.TestGroupInput.LoadProfile.TestDuration.Duration()), // hardcoded request per unit time to 1 as we are using multiCallGen
				LoadType:              wasp.RPS,
				RateLimitUnitDuration: l.TestCfg.TestGroupInput.LoadProfile.TimeUnit.Duration(),
				CallResultBufLen:      10, // we keep the last 10 call results for each generator, as the detailed report is generated at the end of the test
				CallTimeout:           (l.TestCfg.TestGroupInput.PhaseTimeout.Duration()) * 5,
				Gun:                   multiCallGen,
				Logger:                multiCallGen.logger,
				LokiConfig:            wasp.NewLokiConfig(lokiConfig.Endpoint, lokiConfig.TenantId, nil, nil),
				Labels:                allLabels,
				FailOnErr:             pointer.GetBool(l.TestCfg.TestGroupInput.LoadProfile.FailOnFirstErrorInLoad),
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
	ctx := testcontext.Get(t)
	return &LoadArgs{
		t:             t,
		Ctx:           ctx,
		lggr:          &lggr,
		RunnerWg:      wg,
		TestCfg:       testsetups.NewCCIPTestConfig(t, lggr, testconfig.Load),
		ChaosExps:     chaosExps,
		LoadStarterWg: &sync.WaitGroup{},
		pauseLoad:     atomic.NewBool(false),
	}
}
