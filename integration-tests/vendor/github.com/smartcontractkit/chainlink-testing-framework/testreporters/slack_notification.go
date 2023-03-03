// Package testreporters holds all the tools necessary to report on tests that are run utilizing the testsetups package
package testreporters

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/smartcontractkit/chainlink-env/config"
)

// Common Slack Notification Helpers

// Values for reporters to use slack to notify user of test end
var (
	SlackAPIKey  = os.Getenv(config.EnvVarSlackKey)
	SlackChannel = os.Getenv(config.EnvVarSlackChannel)
	SlackUserID  = os.Getenv(config.EnvVarSlackUser)
)

// Uploads a slack file to the designated channel using the API key
func UploadSlackFile(slackClient *slack.Client, uploadParams slack.FileUploadParameters) error {
	log.Info().
		Str("Slack API Key", SlackAPIKey).
		Str("Slack Channel", SlackChannel).
		Str("User Id to Notify", SlackUserID).
		Str("File", uploadParams.File).
		Msg("Attempting to upload file")
	if SlackAPIKey == "" {
		return errors.New("Unable to upload file without a Slack API Key")
	}
	if SlackChannel == "" {
		return errors.New("Unable to upload file without a Slack Channel")
	}
	if uploadParams.Channels == nil || uploadParams.Channels[0] == "" {
		uploadParams.Channels = []string{SlackChannel}
	}
	if uploadParams.File != "" {
		if _, err := os.Stat(uploadParams.File); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("unable to upload file as it does not exist: %w", err)
		} else if err != nil {
			return err
		}
	}
	_, err := slackClient.UploadFile(uploadParams)
	return err
}

// Sends a slack message, and returns an error and the message timestamp
func SendSlackMessage(slackClient *slack.Client, msgOptions ...slack.MsgOption) (string, error) {
	log.Info().
		Str("Slack API Key", SlackAPIKey).
		Str("Slack Channel", SlackChannel).
		Msg("Attempting to send message")
	if SlackAPIKey == "" {
		return "", errors.New("Unable to send message without a Slack API Key")
	}
	if SlackChannel == "" {
		return "", errors.New("Unable to send message without a Slack Channel")
	}
	msgOptions = append(msgOptions, slack.MsgOptionAsUser(true))
	_, timeStamp, err := slackClient.PostMessage(SlackChannel, msgOptions...)
	return timeStamp, err
}

// creates a directory if it doesn't already exist
func MkdirIfNotExists(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to create directory: %s", dirName)
		}
	}
	return nil
}

func CommonSlackNotificationBlocks(
	t *testing.T,
	slackClient *slack.Client,
	headerText, namespace,
	reportCsvLocation,
	slackUserId string,
	testFailed bool,
) []slack.Block {
	notificationBlocks := []slack.Block{}
	notificationBlocks = append(notificationBlocks,
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, true, false)))
	notificationBlocks = append(notificationBlocks,
		slack.NewContextBlock("context_block", slack.NewTextBlockObject("plain_text", namespace, false, false)))
	notificationBlocks = append(notificationBlocks, slack.NewDividerBlock())
	notificationBlocks = append(notificationBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("Summary CSV created on _remote-test-runner_ at _%s_\nNotifying <@%s>",
			reportCsvLocation, SlackUserID), false, true), nil, nil))
	return notificationBlocks
}

// SlackNotifyBlocks creates a slack payload and writes into the specified json
func SlackNotifyBlocks(headerText string, msgtext []string, jsonFile *os.File) error {
	var notificationBlocks slack.Blocks
	notificationBlocks.BlockSet = append(notificationBlocks.BlockSet,
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, true, false)))
	notificationBlocks.BlockSet = append(notificationBlocks.BlockSet, slack.NewDividerBlock())
	msgtexts := ""
	for _, text := range msgtext {
		msgtexts = fmt.Sprintf("%s%s\n", msgtexts, text)
	}
	notificationBlocks.BlockSet = append(notificationBlocks.BlockSet, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		msgtexts, false, true), nil, nil))
	json, err := json.Marshal(slack.Msg{Blocks: notificationBlocks})
	if err != nil {
		return err
	}
	_, err = jsonFile.Write(json)
	return err
}
