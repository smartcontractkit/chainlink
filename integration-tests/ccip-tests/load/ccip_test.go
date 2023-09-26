package load

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/chaos"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
)

func TestLoadCCIPStableRPS(t *testing.T) {
	t.Parallel()
	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr, context.Background())
	testArgs.Setup(true)
	// if the test runs on remote runner
	if len(testArgs.TestSetupArgs.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		log.Info().Msg("Tearing down the environment")
		testArgs.TestSetupArgs.TearDown()
	})
	testArgs.TriggerLoad()
	testArgs.Wait()
}

func TestLoadCCIPSequentialLaneAdd(t *testing.T) {
	t.Parallel()
	t.Skipf("test needs maintenance")
	lggr := logging.GetTestLogger(t)
	testArgs := NewLoadArgs(t, lggr, context.Background())
	testArgs.TestCfg.SequentialLaneAddition = true
	if len(testArgs.TestCfg.NetworkPairs) <= 1 {
		t.Skip("Skipping the test as there are not enough network pairs to run the test")
	}
	testArgs.Setup(true)
	// if the test runs on remote runner
	if len(testArgs.TestSetupArgs.Lanes) == 0 {
		return
	}
	t.Cleanup(func() {
		log.Info().Msg("Tearing down the environment")
		testArgs.TestSetupArgs.TearDown()
	})
	testArgs.TriggerLoad()
	testArgs.AddMoreLanesToRun()
	testArgs.Wait()
}

// This test applies pod chaos to the CL nodes asynchronously and sequentially  while the load is running
// the pod chaos is applied at a regular interval throughout the test duration
func TestLoadCCIPStableRequestTriggeringWithPodChaos(t *testing.T) {
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
			testArgs.TestCfg.TestDuration = 5 * time.Minute
			testArgs.TestCfg.Load.TimeUnit = 1 * time.Second
			testArgs.TestCfg.Load.RequestPerUnitTime = []int64{2}

			testArgs.Setup(false)
			// if the test runs on remote runner
			if len(testArgs.TestSetupArgs.Lanes) == 0 {
				return
			}
			t.Cleanup(func() {
				log.Info().Msg("Tearing down the environment")
				testArgs.TestSetupArgs.TearDown()
			})
			testArgs.SanityCheck()
			testArgs.TriggerLoad()
			testArgs.ApplyChaos()
			testArgs.Wait()
		})
	}
}
