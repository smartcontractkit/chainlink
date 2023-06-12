package src

import (
	"flag"
	"io"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	clcmd "github.com/smartcontractkit/chainlink/v2/core/cmd"
)

func newApp(n *node, writer io.Writer) (*clcmd.Shell, *cli.App) {
	client := &clcmd.Shell{
		Renderer:                       clcmd.RendererJSON{Writer: writer},
		AppFactory:                     clcmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          clcmd.TerminalKeyStoreAuthenticator{Prompter: n},
		FallbackAPIInitializer:         clcmd.NewPromptingAPIInitializer(n),
		Runner:                         clcmd.ChainlinkRunner{},
		PromptingSessionRequestBuilder: clcmd.NewPromptingSessionRequestBuilder(n),
		ChangePasswordPrompter:         clcmd.NewChangePasswordPrompter(),
		PasswordPrompter:               clcmd.NewPasswordPrompter(),
	}
	app := clcmd.NewApp(client)
	fs := flag.NewFlagSet("blah", flag.ContinueOnError)
	fs.String("remote-node-url", n.url.String(), "")
	helpers.PanicErr(app.Before(cli.NewContext(nil, fs, nil)))
	// overwrite renderer since it's set to stdout after Before() is called
	client.Renderer = clcmd.RendererJSON{Writer: writer}
	return client, app
}
