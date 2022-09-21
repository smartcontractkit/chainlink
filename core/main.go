package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/recovery"
	"github.com/smartcontractkit/chainlink/core/static"
)

func init() {
	// check version
	if static.Version == static.Unset {
		if os.Getenv("CHAINLINK_DEV") == "true" {
			return
		}
		log.Println(`Version was unset but CHAINLINK_DEV was not set to "true". Chainlink should be built with static.Version set to a valid semver for production builds.`)
	} else if _, err := semver.NewVersion(static.Version); err != nil {
		panic(fmt.Sprintf("Version invalid: %q is not valid semver", static.Version))
	}
}

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
