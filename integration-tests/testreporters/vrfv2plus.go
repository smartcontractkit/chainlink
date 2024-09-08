package testreporters

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/types"

	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
)

type VRFV2PlusTestReporter struct {
	TestType            string
	LoadTestMetrics     VRFLoadTestMetrics
	VRFv2PlusTestConfig types.VRFv2PlusTestConfig
}

// todo - fix import cycle to avoid struct duplicate
type VRFLoadTestMetrics struct {
	RequestCount                         *big.Int
	FulfilmentCount                      *big.Int
	AverageFulfillmentInMillions         *big.Int
	SlowestFulfillment                   *big.Int
	FastestFulfillment                   *big.Int
	P90FulfillmentBlockTime              float64
	P95FulfillmentBlockTime              float64
	AverageResponseTimeInSecondsMillions *big.Int
	SlowestResponseTimeInSeconds         *big.Int
	FastestResponseTimeInSeconds         *big.Int
}

func (o *VRFV2PlusTestReporter) SetReportData(
	testType string,
	metrics VRFLoadTestMetrics,
	testConfig types.VRFv2PlusTestConfig,
) {
	o.TestType = testType
	o.LoadTestMetrics = metrics
	o.VRFv2PlusTestConfig = testConfig
}

// SendSlackNotification sends a slack message to a slack webhook
func (o *VRFV2PlusTestReporter) SendSlackNotification(t *testing.T, slackClient *slack.Client, vtfv2PlusTestConfig types.VRFv2PlusTestConfig) error {
	if slackClient == nil {
		slackClient = slack.New(testreporters.SlackAPIKey)
	}

	testFailed := t.Failed()
	headerText := fmt.Sprintf(":white_check_mark: VRF V2 Plus %s Test PASSED :white_check_mark:", o.TestType)
	if testFailed {
		headerText = fmt.Sprintf(":x: VRF V2 Plus %s Test FAILED :x:", o.TestType)
	}

	perfCfg := o.VRFv2PlusTestConfig.GetVRFv2PlusConfig().Performance
	messageBlocks := testreporters.SlackNotifyBlocks(headerText, strings.Join(vtfv2PlusTestConfig.GetNetworkConfig().SelectedNetworks, ","), []string{
		fmt.Sprintf(
			"Summary\n"+
				"Perf Test Type: %s\n"+
				"Test Duration set in parameters: %s\n"+
				"Use Existing Env: %t\n"+
				"Request Count: %s\n"+
				"Fulfilment Count: %s\n"+
				"AverageFulfillmentInMillions (blocks): %s\n"+
				"Slowest Fulfillment (blocks): %s\n"+
				"P90 Fulfillment (blocks): %f\n"+
				"P95 Fulfillment (blocks): %f\n"+
				"Fastest Fulfillment (blocks): %s \n"+
				"AverageFulfillmentInMillions (seconds): %s\n"+
				"Slowest Fulfillment (seconds): %s\n"+
				"Fastest Fulfillment (seconds): %s \n"+
				"RPS: %d\n"+
				"RateLimitUnitDuration: %s\n"+
				"RandomnessRequestCountPerRequest: %d\n"+
				"RandomnessRequestCountPerRequestDeviation: %d\n",
			o.TestType,
			perfCfg.TestDuration.Duration.Truncate(time.Second).String(),
			*o.VRFv2PlusTestConfig.GetVRFv2PlusConfig().General.UseExistingEnv,
			o.LoadTestMetrics.RequestCount.String(),
			o.LoadTestMetrics.FulfilmentCount.String(),
			o.LoadTestMetrics.AverageFulfillmentInMillions.String(),
			o.LoadTestMetrics.SlowestFulfillment.String(),
			o.LoadTestMetrics.P90FulfillmentBlockTime,
			o.LoadTestMetrics.P95FulfillmentBlockTime,
			o.LoadTestMetrics.FastestFulfillment.String(),
			o.LoadTestMetrics.AverageResponseTimeInSecondsMillions.String(),
			o.LoadTestMetrics.SlowestResponseTimeInSeconds.String(),
			o.LoadTestMetrics.FastestResponseTimeInSeconds.String(),
			*perfCfg.RPS,
			perfCfg.RateLimitUnitDuration.String(),
			*o.VRFv2PlusTestConfig.GetVRFv2PlusConfig().General.RandomnessRequestCountPerRequest,
			*o.VRFv2PlusTestConfig.GetVRFv2PlusConfig().General.RandomnessRequestCountPerRequestDeviation,
		),
	})

	_, err := testreporters.SendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return err
	}
	return nil
}
