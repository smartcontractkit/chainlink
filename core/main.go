package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/recovery"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

func init() {
	// check version
	if static.Version == static.Unset {
		if v2.EnvDev.IsTrue() {
			return
		}
		log.Println(`Version was unset but dev mode is enabled. Chainlink should be built with static.Version set to a valid semver for production builds.`)
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
	if err := app.Run(args); err != nil {
		log.Fatalf("Error running app: %v\n", err)
	}
}

// newProductionClient configures an instance of the CLI to be used in production.
func newProductionClient() *cmd.Client {
	prompter := cmd.NewTerminalPrompter()
	return &cmd.Client{
		Renderer:                       cmd.RendererTable{Writer: os.Stdout},
		AppFactory:                     cmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          cmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         cmd.NewPromptingAPIInitializer(prompter),
		Runner:                         cmd.ChainlinkRunner{},
		PromptingSessionRequestBuilder: cmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         cmd.NewChangePasswordPrompter(),
		PasswordPrompter:               cmd.NewPasswordPrompter(),
	}
}
