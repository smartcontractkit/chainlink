package automationv2_1

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/slack-go/slack"

	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func extraBlockWithText(text string) slack.Block {
	return slack.NewSectionBlock(slack.NewTextBlockObject(
		"mrkdwn", text, false, false), nil, nil)
}

func sendSlackNotification(header string, l zerolog.Logger, config *tc.TestConfig, namespace string, numberOfNodes,
	startingTime string, endingTime string, extraBlocks []slack.Block, msgOption slack.MsgOption) (string, error) {
	slackClient := slack.New(reportModel.SlackAPIKey)

	headerText := ":chainlink-keepers: Automation Load Test " + header + " :white_check_mark:"

	grafanaUrl, err := config.GetGrafanaBaseURL()
	if err != nil {
		return "", err
	}

	dashboardUrl, err := config.GetGrafanaDashboardURL()
	if err != nil {
		return "", err
	}

	formattedDashboardUrl := fmt.Sprintf("%s%s?orgId=1&from=%s&to=%s&var-namespace=%s&var-number_of_nodes=%s", grafanaUrl, dashboardUrl, startingTime, endingTime, namespace, numberOfNodes)
	l.Info().Str("Dashboard", formattedDashboardUrl).Msg("Dashboard URL")

	var notificationBlocks []slack.Block

	notificationBlocks = append(notificationBlocks,
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, true, false)))
	notificationBlocks = append(notificationBlocks,
		slack.NewContextBlock("context_block", slack.NewTextBlockObject("plain_text", namespace, false, false)))
	notificationBlocks = append(notificationBlocks, slack.NewDividerBlock())
	if *config.Pyroscope.Enabled {
		pyroscopeServer := *config.Pyroscope.ServerUrl
		pyroscopeEnvironment := *config.Pyroscope.Environment

		formattedPyroscopeUrl := fmt.Sprintf("%s/?query=chainlink-node.cpu{Environment=\"%s\"}&from=%s&to=%s", pyroscopeServer, pyroscopeEnvironment, startingTime, endingTime)
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
