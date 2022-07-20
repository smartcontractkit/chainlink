package testreporters

import "github.com/slack-go/slack"

// TestReporter is a general interface for all test reporters
type TestReporter interface {
	WriteReport(folderLocation string) error
	SendSlackNotification(slackClient *slack.Client) error
	SetNamespace(namespace string)
}
