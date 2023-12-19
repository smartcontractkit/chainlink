package automationv2_1

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/slack-go/slack"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s/config"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

func getEnv(key, fallback string) string {
	if inputs, ok := os.LookupEnv("TEST_INPUTS"); ok {
		values := strings.Split(inputs, ",")
		for _, value := range values {
			if strings.Contains(value, key) {
				return strings.Split(value, "=")[1]
			}
		}
	}
	return fallback
}

func extraBlockWithText(text string) slack.Block {
	return slack.NewSectionBlock(slack.NewTextBlockObject(
		"mrkdwn", text, false, false), nil, nil)
}

func sendSlackNotification(header string, l zerolog.Logger, namespace string, numberOfNodes,
	startingTime string, endingTime string, extraBlocks []slack.Block, msgOption slack.MsgOption) (string, error) {
	slackClient := slack.New(reportModel.SlackAPIKey)

	headerText := ":chainlink-keepers: Automation Load Test " + header + " :white_check_mark:"

	formattedDashboardUrl := fmt.Sprintf("%s?orgId=1&from=%s&to=%s&var-namespace=%s&var-number_of_nodes=%s", testreporters.DashboardUrl, startingTime, endingTime, namespace, numberOfNodes)
	l.Info().Str("Dashboard", formattedDashboardUrl).Msg("Dashboard URL")

	pyroscopeServer := os.Getenv(config.EnvVarPyroscopeServer)
	pyroscopeEnvironment := os.Getenv(config.EnvVarPyroscopeEnvironment)

	formattedPyroscopeUrl := fmt.Sprintf("%s/?query=chainlink-node.cpu{Environment=\"%s\"}&from=%s&to=%s", pyroscopeServer, pyroscopeEnvironment, startingTime, endingTime)

	var notificationBlocks []slack.Block

	notificationBlocks = append(notificationBlocks,
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, true, false)))
	notificationBlocks = append(notificationBlocks,
		slack.NewContextBlock("context_block", slack.NewTextBlockObject("plain_text", namespace, false, false)))
	notificationBlocks = append(notificationBlocks, slack.NewDividerBlock())
	if pyroscopeServer != "" {
		l.Info().Str("Pyroscope", formattedPyroscopeUrl).Msg("Dashboard URL")
		notificationBlocks = append(notificationBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
			fmt.Sprintf("<%s|Pyroscope>",
				formattedPyroscopeUrl), false, true), nil, nil))
	}
	notificationBlocks = append(notificationBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("<%s|Test Dashboard> \nNotifying <@%s>",
			formattedDashboardUrl, reportModel.SlackUserID), false, true), nil, nil))

	if len(extraBlocks) > 0 {
		notificationBlocks = append(notificationBlocks, extraBlocks...)
	}

	ts, err := reportModel.SendSlackMessage(slackClient, slack.MsgOptionBlocks(notificationBlocks...), msgOption)
	l.Info().Str("ts", ts).Msg("Sent Slack Message")
	return ts, err
}
