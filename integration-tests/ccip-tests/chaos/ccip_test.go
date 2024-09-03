package chaos_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
)

/* @network-chaos and @pod-chaos are split intentionally into 2 parallel groups
we can't use chaos.NewNetworkPartition and chaos.NewFailPods in parallel
because of jsii runtime bug, see Makefile and please use those targets to run tests
In .github/workflows/ccip-chaos-tests.yml we use these tags to run these tests separately
*/

func TestChaosCCIP(t *testing.T) {
	inputs := []struct {
		testName             string
		chaosFunc            chaos.ManifestFunc
		chaosProps           *chaos.Props
		waitForChaosRecovery bool
	}{
		{
			testName:  "CCIP works after rpc is down for NetworkA @network-chaos",
			chaosFunc: chaos.NewNetworkPartition,
			chaosProps: &chaos.Props{
				FromLabels: &map[string]*string{actions.ChaosGroupNetworkACCIPGeth: ptr.Ptr("1")},
				// chainlink-0 is default label set for all cll nodes
				ToLabels:    &map[string]*string{"app": ptr.Ptr("chainlink-0")},
				DurationStr: "1m",
			},
			waitForChaosRecovery: true,
		},
		{
			testName:  "CCIP works after rpc is down for NetworkB @network-chaos",
			chaosFunc: chaos.NewNetworkPartition,
			chaosProps: &chaos.Props{
				FromLabels:  &map[string]*string{actions.ChaosGroupNetworkBCCIPGeth: ptr.Ptr("1")},
				ToLabels:    &map[string]*string{"app": ptr.Ptr("chainlink-0")},
				DurationStr: "1m",
			},
			waitForChaosRecovery: true,
		},
		{
			testName:  "CCIP works after 2 rpc's are down for all cll nodes @network-chaos",
			chaosFunc: chaos.NewNetworkPartition,
			chaosProps: &chaos.Props{
				FromLabels:  &map[string]*string{"geth": ptr.Ptr(actions.ChaosGroupCCIPGeth)},
				ToLabels:    &map[string]*string{"app": ptr.Ptr("chainlink-0")},
				DurationStr: "1m",
			},
			waitForChaosRecovery: true,
		},
		{
			testName:  "CCIP Commit works after majority of CL nodes are recovered from pod failure @pod-chaos",
			chaosFunc: chaos.NewFailPods,
			chaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupCommitFaultyPlus: ptr.Ptr("1")},
				DurationStr:    "1m",
			},
			waitForChaosRecovery: true,
		},
		{
			testName:  "CCIP Execution works after majority of CL nodes are recovered from pod failure @pod-chaos",
			chaosFunc: chaos.NewFailPods,
			chaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupExecutionFaultyPlus: ptr.Ptr("1")},
				DurationStr:    "1m",
			},
			waitForChaosRecovery: true,
		},
		{
			testName:  "CCIP Commit works while minority of CL nodes are in failed state for pod failure @pod-chaos",
			chaosFunc: chaos.NewFailPods,
			chaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupCommitFaulty: ptr.Ptr("1")},
				DurationStr:    "90s",
			},
			waitForChaosRecovery: false,
		},
		{
			testName:  "CCIP Execution works while minority of CL nodes are in failed state for pod failure @pod-chaos",
			chaosFunc: chaos.NewFailPods,
			chaosProps: &chaos.Props{
				LabelsSelector: &map[string]*string{actions.ChaosGroupExecutionFaulty: ptr.Ptr("1")},
				DurationStr:    "90s",
			},
			waitForChaosRecovery: false,
		},
	}

	for _, in := range inputs {
		in := in
		t.Run(in.testName, func(t *testing.T) {
			t.Parallel()
			l := logging.GetTestLogger(t)
			testCfg := testsetups.NewCCIPTestConfig(t, l, testconfig.Chaos)
			var numOfRequests = 3

			setUpArgs := testsetups.CCIPDefaultTestSetUp(
				t, &l, "chaos-ccip", nil, testCfg)

			if len(setUpArgs.Lanes) == 0 {
				return
			}

			lane := setUpArgs.Lanes[0].ForwardLane

			tearDown := setUpArgs.TearDown
			testEnvironment := setUpArgs.Env.K8Env
			testSetup := setUpArgs.Env

			testSetup.ChaosLabelForGeth(t, lane.SourceChain.GetNetworkName(), lane.DestChain.GetNetworkName())
			testSetup.ChaosLabelForCLNodes(t)

			lane.RecordStateBeforeTransfer()
			// Send the ccip-request and verify ocr2 is running
			err := lane.SendRequests(1, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err)
			lane.ValidateRequests(nil)

			// apply chaos
			chaosId, err := testEnvironment.Chaos.Run(in.chaosFunc(testEnvironment.Cfg.Namespace, in.chaosProps))
			require.NoError(t, err)
			t.Cleanup(func() {
				if chaosId != "" {
					require.NoError(t, testEnvironment.Chaos.Stop(chaosId))
				}
				require.NoError(t, tearDown())
			})
			lane.RecordStateBeforeTransfer()
			// Now send the ccip-request while the chaos is at play
			err = lane.SendRequests(numOfRequests, big.NewInt(actions.DefaultDestinationGasLimit))
			require.NoError(t, err)
			if in.waitForChaosRecovery {
				// wait for chaos to be recovered before further validation
				require.NoError(t, testEnvironment.Chaos.WaitForAllRecovered(chaosId, 1*time.Minute))
			} else {
				l.Info().Msg("proceeding without waiting for chaos recovery")
			}
			lane.ValidateRequests(nil)
		})
	}
}
