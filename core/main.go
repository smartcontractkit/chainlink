package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/store/config"
)

func main() {
	Run(NewProductionClient(), os.Args...)
}

// Run runs the CLI, providing further command instructions by default.
func Run(client *cmd.Client, args ...string) {
	app := cmd.NewApp(client)
	logger.Default.WarnIf(app.Run(args), "Error running app")
}

// NewProductionClient configures an instance of the CLI to be used
// in production.
func NewProductionClient() *cmd.Client {
	cfg := config.NewGeneralConfig()

	prompter := cmd.NewTerminalPrompter()
	cookieAuth := cmd.NewSessionCookieAuthenticator(cfg, cmd.DiskCookieStore{Config: cfg})
	sr := sessions.SessionRequest{}
	sessionRequestBuilder := cmd.NewFileSessionRequestBuilder()
	if credentialsFile := cfg.AdminCredentialsFile(); credentialsFile != "" {
		var err error
		sr, err = sessionRequestBuilder.Build(credentialsFile)
		if err != nil && errors.Cause(err) != cmd.ErrNoCredentialFile && !os.IsNotExist(err) {
			logger.Fatalw("Error loading API credentials", "error", err, "credentialsFile", credentialsFile)
		}
	}
	return &cmd.Client{
		Renderer:                       cmd.RendererTable{Writer: os.Stdout},
		Config:                         cfg,
		AppFactory:                     cmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          cmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         cmd.NewPromptingAPIInitializer(prompter),
		Runner:                         cmd.ChainlinkRunner{},
		HTTP:                           cmd.NewAuthenticatedHTTPClient(cfg, cookieAuth, sr),
		CookieAuthenticator:            cookieAuth,
		FileSessionRequestBuilder:      sessionRequestBuilder,
		PromptingSessionRequestBuilder: cmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         cmd.NewChangePasswordPrompter(),
		PasswordPrompter:               cmd.NewPasswordPrompter(),
	}
}
