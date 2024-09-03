package load

import (
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
)

func TestLoadCCIPStableRPS(t *testing.T) {
	t.Parallel()
	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr)
	testArgs.Setup()
	// if the test runs on remote runner
	if len(testArgs.TestSetupArgs.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		log.Info().Msg("Tearing down the environment")
		require.NoError(t, testArgs.TestSetupArgs.TearDown())
	})
	testArgs.TriggerLoadByLane()
	testArgs.Wait()
}

// TestLoadCCIPWithUpgradeNodeVersion starts all nodes with a specific version, triggers load and then upgrades the node version as the load is running
func TestLoadCCIPWithUpgradeNodeVersion(t *testing.T) {
	t.Parallel()
	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr)
	testArgs.Setup()
	// if the test runs on remote runner
	if len(testArgs.TestSetupArgs.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		log.Info().Msg("Tearing down the environment")
		require.NoError(t, testArgs.TestSetupArgs.TearDown())
	})
	testArgs.TriggerLoadByLane()
	testArgs.lggr.Info().Msg("Waiting for load to start on all lanes")
	// wait for load runner to start
	testArgs.LoadStarterWg.Wait()
	// sleep for 30s to let load run for a while
	time.Sleep(30 * time.Second)
	// upgrade node version for few nodes
	err := testsetups.UpgradeNodes(testArgs.t, testArgs.lggr, testArgs.TestCfg, testArgs.TestSetupArgs.Env)
	require.NoError(t, err)
	// after upgrade send a request to all lanes as a sanity check
	testArgs.SanityCheck()
	// now wait for the load to finish
	testArgs.Wait()
}

func TestLoadCCIPStableRPSTriggerBySource(t *testing.T) {
	t.Parallel()
	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr)
	testArgs.TestCfg.TestGroupInput.MulticallInOneTx = ptr.Ptr(true)
	testArgs.Setup()
	// if the test runs on remote runner
	if len(testArgs.TestSetupArgs.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		log.Info().Msg("Tearing down the environment")
		testArgs.TearDown()
	})
	testArgs.TriggerLoadBySource()
	testArgs.Wait()
}

func TestLoadCCIPStableRequestTriggeringWithNetworkChaos(t *testing.T) {
	t.Parallel()
	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr)
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
	if testArgs.TestCfg.TestGroupInput.LoadProfile.NetworkChaosDelay == nil {
		testArgs.TestCfg.TestGroupInput.LoadProfile.NetworkChaosDelay = config.MustNewDuration(200 * time.Millisecond)
	}
	chaosId, err := testEnv.K8Env.Chaos.Run(
		chaos.NewNetworkLatency(
			testEnv.K8Env.Cfg.Namespace, &chaos.Props{
				FromLabels:  &map[string]*string{"geth": ptr.Ptr(actions.ChaosGroupCCIPGeth)},
				ToLabels:    &map[string]*string{"app": ptr.Ptr("chainlink-0")},
				DurationStr: testArgs.TestCfg.TestGroupInput.LoadProfile.TestDuration.String(),
				Delay:       testArgs.TestCfg.TestGroupInput.LoadProfile.NetworkChaosDelay.Duration().String(),
			}))
	require.NoError(t, err)

	t.Cleanup(func() {
		if chaosId != "" {
			require.NoError(t, testEnv.K8Env.Chaos.Stop(chaosId))
		}
	})

	// now trigger the load
	testArgs.TriggerLoadByLane()
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
				LabelsSelector: &map[string]*string{actions.ChaosGroupCommitFaultyPlus: ptr.Ptr("1")},
				DurationStr:    "2m",
			},
		},
	}

	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr, inputs...)

	var allChaosDur time.Duration
	// to override the default duration of chaos with test input
	for i := range inputs {
		inputs[i].ChaosProps.DurationStr = testArgs.TestCfg.TestGroupInput.ChaosDuration.String()
		allChaosDur += testArgs.TestCfg.TestGroupInput.ChaosDuration.Duration()
		inputs[i].WaitBetweenChaos = testArgs.TestCfg.TestGroupInput.LoadProfile.WaitBetweenChaosDuringLoad.Duration()
		allChaosDur += inputs[i].WaitBetweenChaos
	}

	// the duration of load test should be greater than the duration of chaos
	if testArgs.TestCfg.TestGroupInput.LoadProfile.TestDuration.Duration() < allChaosDur+2*time.Minute {
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

	testArgs.TriggerLoadByLane()
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
				LabelsSelector: &map[string]*string{actions.ChaosGroupCommitFaulty: ptr.Ptr("1")},
				DurationStr:    "4m",
			},
		},
	}

	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr, inputs...)

	var allChaosDur time.Duration
	// to override the default duration of chaos with test input
	for i := range inputs {
		inputs[i].ChaosProps.DurationStr = testArgs.TestCfg.TestGroupInput.ChaosDuration.String()
		allChaosDur += testArgs.TestCfg.TestGroupInput.ChaosDuration.Duration()
		inputs[i].WaitBetweenChaos = testArgs.TestCfg.TestGroupInput.LoadProfile.WaitBetweenChaosDuringLoad.Duration()
		allChaosDur += inputs[i].WaitBetweenChaos
	}

	// the duration of load test should be greater than the duration of chaos
	if testArgs.TestCfg.TestGroupInput.LoadProfile.TestDuration.Duration() < allChaosDur+2*time.Minute {
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

	testArgs.TriggerLoadByLane()
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
				LabelsSelector: &map[string]*string{actions.ChaosGroupCommitFaultyPlus: ptr.Ptr("1")},
				DurationStr:    "2m",
			},
		},
		{
			ChaosName: "CCIP Execution works after majority of CL nodes are recovered from pod failure @pod-chaos",
			ChaosFunc: chaos.NewFailPods,
			ChaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupExecutionFaultyPlus: ptr.Ptr("1")},
				DurationStr:    "2m",
			},
		},
		{
			ChaosName: "CCIP Commit works while minority of CL nodes are in failed state for pod failure @pod-chaos",
			ChaosFunc: chaos.NewFailPods,
			ChaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupCommitFaulty: ptr.Ptr("1")},
				DurationStr:    "4m",
			},
		},
		{
			ChaosName: "CCIP Execution works while minority of CL nodes are in failed state for pod failure @pod-chaos",
			ChaosFunc: chaos.NewFailPods,
			ChaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupExecutionFaulty: ptr.Ptr("1")},
				DurationStr:    "4m",
			},
		},
	}
	for _, in := range inputs {
		in := in
		t.Run(in.ChaosName, func(t *testing.T) {
			t.Parallel()
			lggr := logging.GetTestLogger(t)
			testArgs := NewLoadArgs(t, lggr, in)
			testArgs.TestCfg.TestGroupInput.LoadProfile.TestDuration = config.MustNewDuration(5 * time.Minute)
			testArgs.TestCfg.TestGroupInput.LoadProfile.TimeUnit = config.MustNewDuration(1 * time.Second)
			testArgs.TestCfg.TestGroupInput.LoadProfile.RequestPerUnitTime = []int64{2}
			testArgs.TestCfg.TestGroupInput.PhaseTimeout = config.MustNewDuration(15 * time.Minute)

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
			testArgs.TriggerLoadByLane()
			testArgs.ApplyChaos()
			testArgs.Wait()
		})
	}
}

// TestLoadCCIPStableRPSAfterARMCurseAndUncurse validates that after ARM curse is lifted
// all pending requests get delivered.
// The test pauses loadgen while ARM is cursed and resumes it when curse is lifted.
// There is a known limitation of this test - if the test is run on remote-runner with high frequency
// the remote-runner pod gets evicted after the loadgen is resumed.
// The recommended frequency for this test 2req/min
func TestLoadCCIPStableRPSAfterARMCurseAndUncurse(t *testing.T) {
	t.Skipf("need to be enabled as part of CCIP-2277")
	t.Parallel()
	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr)
	testArgs.Setup()
	// if the test runs on remote runner
	if len(testArgs.TestSetupArgs.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		log.Info().Msg("Tearing down the environment")
		require.NoError(t, testArgs.TestSetupArgs.TearDown())
	})
	testArgs.TriggerLoadByLane()
	// wait for certain time so that few messages are sent
	time.Sleep(2 * time.Minute)
	// now validate the curse
	testArgs.ValidateCurseFollowedByUncurse()
	testArgs.Wait()
}
