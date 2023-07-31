package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	clcmd "github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func newApp(remoteNodeURL string, writer io.Writer) (*clcmd.Shell, *cli.App) {
	prompter := clcmd.NewTerminalPrompter()
	client := &clcmd.Shell{
		Renderer:                       clcmd.RendererJSON{Writer: writer},
		AppFactory:                     clcmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          clcmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         clcmd.NewPromptingAPIInitializer(prompter),
		Runner:                         clcmd.ChainlinkRunner{},
		PromptingSessionRequestBuilder: clcmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         clcmd.NewChangePasswordPrompter(),
		PasswordPrompter:               clcmd.NewPasswordPrompter(),
	}
	app := clcmd.NewApp(client)
	fs := flag.NewFlagSet("blah", flag.ContinueOnError)
	fs.Bool("json", true, "")
	fs.String("remote-node-url", remoteNodeURL, "")
	helpers.PanicErr(app.Before(cli.NewContext(nil, fs, nil)))
	// overwrite renderer since it's set to stdout after Before() is called
	client.Renderer = clcmd.RendererJSON{Writer: writer}
	return client, app
}

var (
	remoteNodeURLs = flag.String("remote-node-urls", "", "remote node URL")
	credsFile      = flag.String("creds-file", "", "Creds to authenticate to the node")

	checkMarkEmoji = "✅"
	xEmoji         = "❌"
	infoEmoji      = "ℹ️"
)

func main() {
	flag.Parse()

	if remoteNodeURLs == nil {
		fmt.Println("flag -remote-node-urls required")
		os.Exit(1)
	}

	urls := strings.Split(*remoteNodeURLs, ",")
	file := *credsFile

	var (
		allETHKeys []string
	)
	for _, remoteNodeURL := range urls {
		output := &bytes.Buffer{}
		client, app := newApp(remoteNodeURL, output)

		// login first to establish the session
		fmt.Println("logging in to:", remoteNodeURL)
		loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
		loginFs.String("file", file, "")
		loginFs.Bool("bypass-version-check", true, "")
		loginCtx := cli.NewContext(app, loginFs, nil)
		err := client.RemoteLogin(loginCtx)
		helpers.PanicErr(err)
		output.Reset()
		fmt.Println()

		// check for ETH keys
		err = client.ListETHKeys(&cli.Context{
			App: app,
		})
		helpers.PanicErr(err)
		var ethKeys []presenters.ETHKeyResource
		helpers.PanicErr(json.Unmarshal(output.Bytes(), &ethKeys))
		switch {
		case len(ethKeys) >= 5:
			fmt.Println(checkMarkEmoji, "found", len(ethKeys), "eth keys on", remoteNodeURL)
		case len(ethKeys) < 5:
			fmt.Println(xEmoji, "found only", len(ethKeys), "eth keys on", remoteNodeURL, ", consider creating more")
			output.Reset()
			var newETHKeys []presenters.ETHKeyResource

			err = client.CreateETHKey(&cli.Context{
				App: app,
			})
			helpers.PanicErr(err)
			helpers.PanicErr(json.Unmarshal(output.Bytes(), &newETHKeys))
			fmt.Println("NEW ETH KEY:", newETHKeys)

		}
		output.Reset()
		fmt.Println()

		for _, ethKey := range ethKeys {
			allETHKeys = append(allETHKeys, ethKey.Address)
		}

	}

	fmt.Println("------------- NODE INFORMATION -------------")
	fmt.Println("ETH addresses:", strings.Join(allETHKeys, ","))
}
