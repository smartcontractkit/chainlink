package core

import (
	"fmt"
	"log"
	"os"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/recovery"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

func init() {
	// check version
	if static.Version == static.Unset {
		if !build.IsProd() {
			return
		}
		log.Println(`Version was unset on production build. Chainlink should be built with static.Version set to a valid semver for production builds.`)
	} else if _, err := semver.NewVersion(static.Version); err != nil {
		panic(fmt.Sprintf("Version invalid: %q is not valid semver", static.Version))
	}
}

func Main() (code int) {
	recovery.ReportPanics(func() {
		app := cmd.NewApp(newProductionClient())
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
			code = 1
		}
	})
	return
}

// newProductionClient configures an instance of the CLI to be used in production.
func newProductionClient() *cmd.Shell {
	prompter := cmd.NewTerminalPrompter()
	return &cmd.Shell{
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
