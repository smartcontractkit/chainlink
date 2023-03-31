package audit_test

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type MockedHTTPEvent struct {
	body string
}

type MockHTTPClient struct {
	audit.HTTPAuditLoggerInterface

	loggingChannel chan MockedHTTPEvent
}

type LoginData struct {
	Email string `json:"email"`
}

type LoginLogItem struct {
	EventID string    `json:"eventID"`
	Env     string    `json:"env"`
	Data    LoginData `json:"data"`
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

type Config struct{}

func (c Config) AuditLoggerEnabled() bool {
	return true
}

func (c Config) AuditLoggerEnvironment() string {
	return "test"
}

func (c Config) AuditLoggerForwardToUrl() (models.URL, error) {
	url, err := models.ParseURL("http://localhost:9898")
	if err != nil {
		return models.URL{}, err
	}
	return *url, nil
}

func (c Config) AuditLoggerHeaders() (audit.ServiceHeaders, error) {
	return make(audit.ServiceHeaders, 0), nil
}

func (c Config) AuditLoggerJsonWrapperKey() string {
	return ""
}

func TestCheckLoginAuditLog(t *testing.T) {
	t.Parallel()

	// Create a channel that will be used instead of an HTTP client
	loggingChannel := make(chan MockedHTTPEvent, 2048)

	// Create the mock structure that will be used
	mockHTTPClient := MockHTTPClient{
		loggingChannel: loggingChannel,
	}

	// Create a test logger because the audit logger relies on this logger
	// as well
	logger := logger.TestLogger(t)

	auditLoggerTestConfig := Config{}

	// Create new AuditLoggerService
	auditLogger, err := audit.NewAuditLogger(logger.Named("AuditLogger"), &auditLoggerTestConfig)
	assert.NoError(t, err)

	// Cast to concrete type so we can swap out the internals
	auditLoggerService, ok := auditLogger.(*audit.AuditLoggerService)
	assert.True(t, ok)

	// Swap the internals with a testing handler
	auditLoggerService.SetLoggingClient(&mockHTTPClient)
	assert.NoError(t, auditLoggerService.Ready())

	// Create a new chainlink test application passing in our test logger
	// and audit logger
	app := cltest.NewApplication(t, logger, auditLogger)
	require.NoError(t, app.Start(testutils.Context(t)))

	enteredStrings := []string{cltest.APIEmailAdmin, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}
	client := app.NewAuthenticatingClient(prompter)

	set := flag.NewFlagSet("test", 0)
	set.Bool("bypass-version-check", true, "")
	set.String("admin-credentials-file", "", "")
	c := cli.NewContext(nil, set, nil)

	// Login
	err = client.RemoteLogin(c)
	assert.NoError(t, err)

	select {
	case event := <-loggingChannel:
		deserialized := &LoginLogItem{}
		assert.NoError(t, json.Unmarshal([]byte(event.body), deserialized))

		assert.Equal(t, deserialized.Data.Email, cltest.APIEmailAdmin)
		assert.Equal(t, deserialized.Env, "test")

		assert.Equal(t, deserialized.EventID, "AUTH_LOGIN_SUCCESS_NO_2FA")
		return
	case <-time.After(5 * time.Second):
	}

	assert.True(t, false)
}
