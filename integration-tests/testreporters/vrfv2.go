package testreporters

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_config"

	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
)

type VRFV2TestReporter struct {
	TestType                     string
	RequestCount                 *big.Int
	FulfilmentCount              *big.Int
	AverageFulfillmentInMillions *big.Int
	SlowestFulfillment           *big.Int
	FastestFulfillment           *big.Int
	Vrfv2Config                  *vrfv2_config.VRFV2Config
}

func (o *VRFV2TestReporter) SetReportData(
	testType string,
	RequestCount *big.Int,
	FulfilmentCount *big.Int,
	AverageFulfillmentInMillions *big.Int,
	SlowestFulfillment *big.Int,
	FastestFulfillment *big.Int,
	vrfv2Config vrfv2_config.VRFV2Config,
) {
	o.TestType = testType
	o.RequestCount = RequestCount
	o.FulfilmentCount = FulfilmentCount
	o.AverageFulfillmentInMillions = AverageFulfillmentInMillions
	o.SlowestFulfillment = SlowestFulfillment
	o.FastestFulfillment = FastestFulfillment
	o.Vrfv2Config = &vrfv2Config
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

	messageBlocks := testreporters.SlackNotifyBlocks(headerText, os.Getenv("SELECTED_NETWORKS"), []string{
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
			o.Vrfv2Config.TestDuration.Truncate(time.Second).String(),
			o.Vrfv2Config.UseExistingEnv,
			o.RequestCount.String(),
			o.FulfilmentCount.String(),
			o.AverageFulfillmentInMillions.String(),
			o.SlowestFulfillment.String(),
			o.FastestFulfillment.String(),
			o.Vrfv2Config.RPS,
			o.Vrfv2Config.RateLimitUnitDuration.String(),
			o.Vrfv2Config.RandomnessRequestCountPerRequest,
			o.Vrfv2Config.RandomnessRequestCountPerRequestDeviation,
		),
	})

	_, err := testreporters.SendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return err
	}
	return nil
}
