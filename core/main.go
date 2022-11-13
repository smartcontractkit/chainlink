package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/recovery"
	"github.com/smartcontractkit/chainlink/core/static"
)

func init() {
	// check version
	if static.Version == static.Unset {
		if isDevMode() {
			return
		}
		log.Println(`Version was unset but dev mode is enabled. Chainlink should be built with static.Version set to a valid semver for production builds.`)
	} else if _, err := semver.NewVersion(static.Version); err != nil {
		panic(fmt.Sprintf("Version invalid: %q is not valid semver", static.Version))
	}
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func isDevMode() bool {
	var clDev string
	v1, v2 := os.Getenv("CHAINLINK_DEV"), os.Getenv("CL_DEV")
	if v1 != "" && v2 != "" {
		if v1 != v2 {
			panic("you may only set one of CHAINLINK_DEV and CL_DEV environment variables, not both")
		}
	} else if v1 == "" {
		clDev = v2
	} else if v2 == "" {
		clDev = v1
	}
	return strings.ToLower(clDev) == "true"
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
