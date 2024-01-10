package util

import (
	"bytes"
	"flag"
	"fmt"
	"io"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/scripts/common/vrf/model"
	clcmd "github.com/smartcontractkit/chainlink/v2/core/cmd"
)

func MapToSendingKeyArr(nodeSendingKeys []string) []model.SendingKey {
	var sendingKeys []model.SendingKey

	for _, key := range nodeSendingKeys {
		sendingKeys = append(sendingKeys, model.SendingKey{Address: key})
	}
	return sendingKeys
}

func MapToAddressArr(sendingKeys []model.SendingKey) []string {
	var sendingKeysString []string
	for _, sendingKey := range sendingKeys {
		sendingKeysString = append(sendingKeysString, sendingKey.Address)
	}
	return sendingKeysString
}

func ConnectToNode(nodeURL *string, output *bytes.Buffer, credFile *string) (*clcmd.Shell, *cli.App) {
	client, app := newApp(*nodeURL, output)
	// login first to establish the session
	fmt.Println("logging in to:", *nodeURL)
	loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
	if credFile != nil {
		loginFs.String("file", *credFile, "")
	}
	loginFs.Bool("bypass-version-check", true, "")
	loginCtx := cli.NewContext(app, loginFs, nil)
	err := client.RemoteLogin(loginCtx)
	helpers.PanicErr(err)
	output.Reset()
	fmt.Println()
	return client, app
}

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
