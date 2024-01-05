package testreporters

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

type VRFV2TestReporter struct {
	TestType                     string
	RequestCount                 *big.Int
	FulfilmentCount              *big.Int
	AverageFulfillmentInMillions *big.Int
	SlowestFulfillment           *big.Int
	FastestFulfillment           *big.Int
	TestConfig                   *tc.TestConfig
}

func (o *VRFV2TestReporter) SetReportData(
	testType string,
	RequestCount *big.Int,
	FulfilmentCount *big.Int,
	AverageFulfillmentInMillions *big.Int,
	SlowestFulfillment *big.Int,
	FastestFulfillment *big.Int,
	testConfig *tc.TestConfig,
) {
	o.TestType = testType
	o.RequestCount = RequestCount
	o.FulfilmentCount = FulfilmentCount
	o.AverageFulfillmentInMillions = AverageFulfillmentInMillions
	o.SlowestFulfillment = SlowestFulfillment
	o.FastestFulfillment = FastestFulfillment
	o.TestConfig = testConfig
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

	cfg := o.TestConfig.VRFv2.Performance
	var sb strings.Builder
	for _, n := range o.TestConfig.Network.SelectedNetworks {
		sb.WriteString(n)
		sb.WriteString(", ")
	}

	messageBlocks := testreporters.SlackNotifyBlocks(headerText, sb.String(), []string{
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
			cfg.TestDuration.Duration.Truncate(time.Second).String(),
			*cfg.UseExistingEnv,
			o.RequestCount.String(),
			o.FulfilmentCount.String(),
			o.AverageFulfillmentInMillions.String(),
			o.SlowestFulfillment.String(),
			o.FastestFulfillment.String(),
			*cfg.RPS,
			cfg.RateLimitUnitDuration.String(),
			*o.TestConfig.VRFv2.General.RandomnessRequestCountPerRequest,
			*o.TestConfig.VRFv2.General.RandomnessRequestCountPerRequestDeviation,
		),
	})

	_, err := testreporters.SendSlackMessage(slackClient, slack.MsgOptionBlocks(messageBlocks...))
	if err != nil {
		return err
	}
	return nil
}
