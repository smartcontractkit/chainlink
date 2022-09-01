package main

import (
	"log"
	"os"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/recovery"
)

func main() {
	recovery.ReportPanics(func() {
		run(newProductionClient(), os.Args...)
	})
}

// run the CLI, providing further command instructions by default.
func run(client *cmd.Client, args ...string) {
	app := cmd.NewApp(client)
	client.Logger.ErrorIf(app.Run(args), "Error running app")
	if err := client.CloseLogger(); err != nil {
		log.Fatal(err)
	}
}

// newProductionClient configures an instance of the CLI to be used in production.
func newProductionClient() *cmd.Client {
	lggr, closeLggr := logger.NewLogger()
	prompter := cmd.NewTerminalPrompter()
	return &cmd.Client{
		Renderer:                       cmd.RendererTable{Writer: os.Stdout},
		Logger:                         lggr,
		CloseLogger:                    closeLggr,
		AppFactory:                     cmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          cmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         cmd.NewPromptingAPIInitializer(prompter, lggr),
		Runner:                         cmd.ChainlinkRunner{},
		PromptingSessionRequestBuilder: cmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         cmd.NewChangePasswordPrompter(),
		PasswordPrompter:               cmd.NewPasswordPrompter(),
	}
}
