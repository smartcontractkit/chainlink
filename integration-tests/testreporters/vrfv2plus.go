package testreporters

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2plus/vrfv2plus_config"

	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
)

type VRFV2PlusTestReporter struct {
	TestType                     string
	RequestCount                 *big.Int
	FulfilmentCount              *big.Int
	AverageFulfillmentInMillions *big.Int
	SlowestFulfillment           *big.Int
	FastestFulfillment           *big.Int
	Vrfv2PlusConfig              *vrfv2plus_config.VRFV2PlusConfig
}

func (o *VRFV2PlusTestReporter) SetReportData(
	testType string,
	RequestCount *big.Int,
	FulfilmentCount *big.Int,
	AverageFulfillmentInMillions *big.Int,
	SlowestFulfillment *big.Int,
	FastestFulfillment *big.Int,
	vrfv2PlusConfig vrfv2plus_config.VRFV2PlusConfig,
) {
	o.TestType = testType
	o.RequestCount = RequestCount
	o.FulfilmentCount = FulfilmentCount
	o.AverageFulfillmentInMillions = AverageFulfillmentInMillions
	o.SlowestFulfillment = SlowestFulfillment
	o.FastestFulfillment = FastestFulfillment
	o.Vrfv2PlusConfig = &vrfv2PlusConfig
}

// SendSlackNotification sends a slack message to a slack webhook
func (o *VRFV2PlusTestReporter) SendSlackNotification(t *testing.T, slackClient *slack.Client) error {
	if slackClient == nil {
		slackClient = slack.New(testreporters.SlackAPIKey)
	}

	testFailed := t.Failed()
	headerText := fmt.Sprintf(":white_check_mark: VRF V2 Plus %s Test PASSED :white_check_mark:", o.TestType)
	if testFailed {
		headerText = fmt.Sprintf(":x: VRF V2 Plus %s Test FAILED :x:", o.TestType)
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
			o.Vrfv2PlusConfig.TestDuration.Truncate(time.Second).String(),
			o.Vrfv2PlusConfig.UseExistingEnv,
			o.RequestCount.String(),
			o.FulfilmentCount.String(),
			o.AverageFulfillmentInMillions.String(),
			o.SlowestFulfillment.String(),
			o.FastestFulfillment.String(),
			o.Vrfv2PlusConfig.RPS,
			o.Vrfv2PlusConfig.RateLimitUnitDuration.String(),
			o.Vrfv2PlusConfig.RandomnessRequestCountPerRequest,
			o.Vrfv2PlusConfig.RandomnessRequestCountPerRequestDeviation,
		),
	})

	_, err := testreporters.SendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return err
	}
	return nil
}
