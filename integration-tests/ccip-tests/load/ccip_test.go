package load

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/chaos"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/ccip/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestLoadCCIPStableRPS(t *testing.T) {
	t.Parallel()
	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr, context.Background())
	testArgs.Setup()
	// if the test runs on remote runner
	if len(testArgs.TestSetupArgs.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		log.Info().Msg("Tearing down the environment")
		require.NoError(t, testArgs.TestSetupArgs.TearDown())
	})
	testArgs.TriggerLoad()
	testArgs.Wait()
}

func TestLoadCCIPSequentialLaneAdd(t *testing.T) {
	t.Parallel()
	t.Skipf("test needs maintenance")
	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr, context.Background())
	testArgs.TestCfg.TestGroupInput.SequentialLaneAddition = utils.Ptr(true)
	if len(testArgs.TestCfg.NetworkPairs) <= 1 {
		t.Skip("Skipping the test as there are not enough network pairs to run the test")
	}
	testArgs.Setup()
	// if the test runs on remote runner
	if len(testArgs.TestSetupArgs.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		log.Info().Msg("Tearing down the environment")
		require.NoError(t, testArgs.TestSetupArgs.TearDown())
	})
	testArgs.TriggerLoad()
	testArgs.AddMoreLanesToRun()
	testArgs.Wait()
}

func TestLoadCCIPStableRequestTriggeringWithNetworkChaos(t *testing.T) {
	t.Parallel()
	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr, context.Background())
	testArgs.Setup()
	// if the test runs on remote runner
	if len(testArgs.TestSetupArgs.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		log.Info().Msg("Tearing down the environment")
		require.NoError(t, testArgs.TestSetupArgs.TearDown())
	})
	testEnv := testArgs.TestSetupArgs.Env
	require.NotNil(t, testEnv)
	require.NotNil(t, testEnv.K8Env)

	// apply network chaos so that chainlink's RPC calls are affected by some network delay for the duration of the test
	var gethNetworksLabels []string
	for _, net := range testArgs.TestCfg.SelectedNetworks {
		gethNetworksLabels = append(gethNetworksLabels, actions.GethLabel(net.Name))
	}
	testEnv.ChaosLabelForAllGeth(t, gethNetworksLabels)
	chaosId, err := testEnv.K8Env.Chaos.Run(
		chaos.NewNetworkLatency(
			testEnv.K8Env.Cfg.Namespace, &chaos.Props{
				FromLabels:  &map[string]*string{"geth": a.Str(actions.ChaosGroupCCIPGeth)},
				ToLabels:    &map[string]*string{"app": a.Str("chainlink-0")},
				DurationStr: testArgs.TestCfg.TestGroupInput.TestDuration.String(),
				Delay:       "300ms",
			}))
	require.NoError(t, err)

	t.Cleanup(func() {
		if chaosId != "" {
			require.NoError(t, testEnv.K8Env.Chaos.Stop(chaosId))
		}
	})

	// now trigger the load
	testArgs.TriggerLoad()
	testArgs.Wait()
}

// This test applies pod chaos to the CL nodes asynchronously and sequentially while the load is running
// the pod chaos is applied at a regular interval throughout the test duration
// this test needs to be run for a longer duration to see the effects of pod chaos
// in this test commit and execution are set up to be on the same node
func TestLoadCCIPStableWithMajorityNodeFailure(t *testing.T) {
	t.Parallel()

	inputs := []ChaosConfig{
		{
			ChaosName: "CCIP works after majority of CL nodes are recovered from pod failure @pod-chaos",
			ChaosFunc: chaos.NewFailPods,
			ChaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupCommitFaultyPlus: a.Str("1")},
				DurationStr:    "2m",
			},
		},
	}

	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr, context.Background(), inputs...)

	var allChaosDur time.Duration
	// to override the default duration of chaos with test input
	for i := range inputs {
		inputs[i].ChaosProps.DurationStr = testArgs.TestCfg.TestGroupInput.ChaosDuration.String()
		allChaosDur += testArgs.TestCfg.TestGroupInput.ChaosDuration.Duration()
		inputs[i].WaitBetweenChaos = testArgs.TestCfg.TestGroupInput.WaitBetweenChaosDuringLoad.Duration()
		allChaosDur += inputs[i].WaitBetweenChaos
	}

	// the duration of load test should be greater than the duration of chaos
	if testArgs.TestCfg.TestGroupInput.TestDuration.Duration() < allChaosDur+2*time.Minute {
		t.Fatalf("Skipping the test as the test duration is less than the chaos duration")
	}

	testArgs.Setup()
	// if the test runs on remote runner
	if len(testArgs.TestSetupArgs.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		log.Info().Msg("Tearing down the environment")
		require.NoError(t, testArgs.TestSetupArgs.TearDown())
	})

	testEnv := testArgs.TestSetupArgs.Env
	require.NotNil(t, testEnv)
	require.NotNil(t, testEnv.K8Env)

	testArgs.TriggerLoad()
	testArgs.ApplyChaos()
	testArgs.Wait()
}

// This test applies pod chaos to the CL nodes asynchronously and sequentially while the load is running
// the pod chaos is applied at a regular interval throughout the test duration
// this test needs to be run for a longer duration to see the effects of pod chaos
// in this test commit and execution are set up to be on the same node
func TestLoadCCIPStableWithMinorityNodeFailure(t *testing.T) {
	t.Parallel()

	inputs := []ChaosConfig{
		{
			ChaosName: "CCIP works while minority of CL nodes are in failed state for pod failure @pod-chaos",
			ChaosFunc: chaos.NewFailPods,
			ChaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupCommitFaulty: a.Str("1")},
				DurationStr:    "4m",
			},
		},
	}

	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr, context.Background(), inputs...)

	var allChaosDur time.Duration
	// to override the default duration of chaos with test input
	for i := range inputs {
		inputs[i].ChaosProps.DurationStr = testArgs.TestCfg.TestGroupInput.ChaosDuration.String()
		allChaosDur += testArgs.TestCfg.TestGroupInput.ChaosDuration.Duration()
		inputs[i].WaitBetweenChaos = testArgs.TestCfg.TestGroupInput.WaitBetweenChaosDuringLoad.Duration()
		allChaosDur += inputs[i].WaitBetweenChaos
	}

	// the duration of load test should be greater than the duration of chaos
	if testArgs.TestCfg.TestGroupInput.TestDuration.Duration() < allChaosDur+2*time.Minute {
		t.Fatalf("Skipping the test as the test duration is less than the chaos duration")
	}

	testArgs.Setup()
	// if the test runs on remote runner
	if len(testArgs.TestSetupArgs.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		log.Info().Msg("Tearing down the environment")
		require.NoError(t, testArgs.TestSetupArgs.TearDown())
	})

	testEnv := testArgs.TestSetupArgs.Env
	require.NotNil(t, testEnv)
	require.NotNil(t, testEnv.K8Env)

	testArgs.TriggerLoad()
	testArgs.ApplyChaos()
	testArgs.Wait()
}

// This test applies pod chaos to the CL nodes asynchronously and sequentially  while the load is running
// the pod chaos is applied at a regular interval throughout the test duration
// in this test commit and execution are set up to be on different node
func TestLoadCCIPStableWithPodChaosDiffCommitAndExec(t *testing.T) {
	t.Parallel()
	inputs := []ChaosConfig{
		{
			ChaosName: "CCIP Commit works after majority of CL nodes are recovered from pod failure @pod-chaos",
			ChaosFunc: chaos.NewFailPods,
			ChaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupCommitFaultyPlus: a.Str("1")},
				DurationStr:    "2m",
			},
		},
		{
			ChaosName: "CCIP Execution works after majority of CL nodes are recovered from pod failure @pod-chaos",
			ChaosFunc: chaos.NewFailPods,
			ChaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupExecutionFaultyPlus: a.Str("1")},
				DurationStr:    "2m",
			},
		},
		{
			ChaosName: "CCIP Commit works while minority of CL nodes are in failed state for pod failure @pod-chaos",
			ChaosFunc: chaos.NewFailPods,
			ChaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupCommitFaulty: a.Str("1")},
				DurationStr:    "4m",
			},
		},
		{
			ChaosName: "CCIP Execution works while minority of CL nodes are in failed state for pod failure @pod-chaos",
			ChaosFunc: chaos.NewFailPods,
			ChaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupExecutionFaulty: a.Str("1")},
				DurationStr:    "4m",
			},
		},
	}
	for _, in := range inputs {
		in := in
		t.Run(in.ChaosName, func(t *testing.T) {
			t.Parallel()
			lggr := logging.GetTestLogger(t)
			testArgs := NewLoadArgs(t, lggr, context.Background(), in)
			testArgs.TestCfg.TestGroupInput.TestDuration = models.MustNewDuration(5 * time.Minute)
			testArgs.TestCfg.TestGroupInput.TimeUnit = models.MustNewDuration(1 * time.Second)
			testArgs.TestCfg.TestGroupInput.RequestPerUnitTime = []int64{2}

			testArgs.Setup()
			// if the test runs on remote runner
			if len(testArgs.TestSetupArgs.Lanes) == 0 {
				return
			}
			t.Cleanup(func() {
				log.Info().Msg("Tearing down the environment")
				require.NoError(t, testArgs.TestSetupArgs.TearDown())
			})
			testArgs.SanityCheck()
			testArgs.TriggerLoad()
			testArgs.ApplyChaos()
			testArgs.Wait()
		})
	}
}
