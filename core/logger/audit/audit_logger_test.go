package audit_test

import (
	"flag"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

type MockedHTTPEvent struct {
	body string
}

type MockHTTPClient struct {
	audit.HTTPAuditLoggerInterface

	loggingChannel chan MockedHTTPEvent
}

func (mock *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	b, err := io.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	message := MockedHTTPEvent{
		body: string(b),
	}

	mock.loggingChannel <- message

	return &http.Response{}, nil
}

func getAuditLoggerConfig() *audit.AuditLoggerConfig {
	forwardToUrl := "empty"
	environment := "test"
	jsonWrapperKey := ""

	return &audit.AuditLoggerConfig{
		ForwardToUrl:   &forwardToUrl,
		Environment:    &environment,
		JsonWrapperKey: &jsonWrapperKey,
	}
}

func TestCheckLoginAuditLog(t *testing.T) {
	t.Parallel()

	// Create a channel that will be used instead of an HTTP client
	loggingChannel := make(chan MockedHTTPEvent, 2048)

	// Create the mock structure that will be used
	mockHTTPClient := MockHTTPClient{
		loggingChannel: loggingChannel,
	}

	// Set an environment variable to trick the system into starting the audit logger
	// and adding it to the subsystems. We are going to swap it out later so it's not
	// going to matter what we put here
	// os.Setenv("AUDIT_LOGGER_FORWARD_TO_URL", "http://test.local:9999")

	logger := logger.TestLogger(t)

	// Create new AuditLoggerService
	auditLogger, err := audit.NewAuditLogger(logger, getAuditLoggerConfig())
	assert.NoError(t, err)

	auditLoggerService, ok := auditLogger.(*audit.AuditLoggerService)
	assert.True(t, ok)

	auditLoggerService.SetLoggingClient(&mockHTTPClient)

	assert.NoError(t, auditLoggerService.Ready())

	app := cltest.NewApplication(t, logger, auditLogger)

	require.NoError(t, app.Start(testutils.Context(t)))

	enteredStrings := []string{cltest.APIEmailAdmin, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}
	client := app.NewAuthenticatingClient(prompter)

	set := flag.NewFlagSet("test", 0)
	set.Bool("bypass-version-check", true, "")
	set.String("admin-credentials-file", "", "")
	c := cli.NewContext(nil, set, nil)

	err = client.RemoteLogin(c)
	assert.NoError(t, err)

	select {
	case _ = <-loggingChannel:
		return
	case <-time.After(15 * time.Second):
	}

	assert.True(t, false)
}
