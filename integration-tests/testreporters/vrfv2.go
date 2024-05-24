package testreporters

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/types"
)

type VRFV2TestReporter struct {
	TestType        string
	LoadTestMetrics VRFLoadTestMetrics
	VRFv2TestConfig types.VRFv2TestConfig
}

func (o *VRFV2TestReporter) SetReportData(
	testType string,
	metrics VRFLoadTestMetrics,
	vrfv2TestConfig types.VRFv2TestConfig,
) {
	o.TestType = testType
	o.LoadTestMetrics = metrics
	o.VRFv2TestConfig = vrfv2TestConfig
}

// SendSlackNotification sends a slack message to a slack webhook
func (o *VRFV2TestReporter) SendSlackNotification(t *testing.T, slackClient *slack.Client) error {
	if slackClient == nil {
		slackClient = slack.New(testreporters.SlackAPIKey)
	}

	testFailed := t.Failed()
	headerText := fmt.Sprintf(":white_check_mark: VRF V2 %s Test PASSED :white_check_mark:", o.TestType)
	if testFailed {
		headerText = fmt.Sprintf(":x: VRF V2 %s Test FAILED :x:", o.TestType)
	}

	perfCfg := o.VRFv2TestConfig.GetVRFv2Config().Performance
	messageBlocks := testreporters.SlackNotifyBlocks(headerText, strings.Join(o.VRFv2TestConfig.GetNetworkConfig().SelectedNetworks, ","), []string{
		fmt.Sprintf(
			"Summary\n"+
				"Perf Test Type: %s\n"+
				"Test Duration set in parameters: %s\n"+
				"Use Existing Env: %t\n"+
				"Request Count: %s\n"+
				"Fulfilment Count: %s\n"+
				"AverageFulfillmentInMillions: %s\n"+
				"Slowest Fulfillment: %s\n"+
				"Fastest Fulfillment: %s \n"+
				"RPS: %d\n"+
				"RateLimitUnitDuration: %s\n"+
				"RandomnessRequestCountPerRequest: %d\n"+
				"RandomnessRequestCountPerRequestDeviation: %d\n",
			o.TestType,
			perfCfg.TestDuration.Duration.Truncate(time.Second).String(),
			*o.VRFv2TestConfig.GetVRFv2Config().General.UseExistingEnv,
			o.LoadTestMetrics.RequestCount.String(),
			o.LoadTestMetrics.FulfilmentCount.String(),
			o.LoadTestMetrics.AverageFulfillmentInMillions.String(),
			o.LoadTestMetrics.SlowestFulfillment.String(),
			o.LoadTestMetrics.FastestFulfillment.String(),
			*perfCfg.RPS,
			perfCfg.RateLimitUnitDuration.String(),
			*o.VRFv2TestConfig.GetVRFv2Config().General.RandomnessRequestCountPerRequest,
			*o.VRFv2TestConfig.GetVRFv2Config().General.RandomnessRequestCountPerRequestDeviation,
		),
	})

	_, err := testreporters.SendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return err
	}
	return nil
}
