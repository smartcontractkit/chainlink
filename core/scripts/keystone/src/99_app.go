package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/jpillora/backoff"
	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	clcmd "github.com/smartcontractkit/chainlink/v2/core/cmd"
)

// NewRedialBackoff is a standard backoff to use for redialling or reconnecting to
// unreachable network endpoints
func NewRedialBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}
}

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
	fs.Bool("insecure-skip-verify", true, "")
	helpers.PanicErr(app.Before(cli.NewContext(nil, fs, nil)))
	// overwrite renderer since it's set to stdout after Before() is called
	client.Renderer = clcmd.RendererJSON{Writer: writer}
	return client, app
}

type nodeAPI struct {
	methods      *cmd.Shell
	app          *cli.App
	output       *bytes.Buffer
	fs           *flag.FlagSet
	clientMethod func(*cli.Context) error
}

func newNodeAPI(n *node) *nodeAPI {
	output := &bytes.Buffer{}
	methods, app := newApp(n, output)

	api := &nodeAPI{
		output:  output,
		methods: methods,
		app:     app,
		fs:      flag.NewFlagSet("test", flag.ContinueOnError),
	}

	fmt.Println("Logging in:", n.url)
	loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
	loginFs.Bool("bypass-version-check", true, "")
	loginCtx := cli.NewContext(app, loginFs, nil)

	redial := NewRedialBackoff()

	for {
		err := methods.RemoteLogin(loginCtx)
		if err == nil {
			break
		}

		fmt.Println("Error logging in:", err)
		if strings.Contains(err.Error(), "invalid character '<' looking for beginning of value") {
			fmt.Println("Likely a transient network error, retrying...")
		} else {
			helpers.PanicErr(err)
		}

		time.Sleep(redial.Duration())
	}
	output.Reset()

	return api
}

func (c *nodeAPI) withArg(arg string) *nodeAPI {
	err := c.fs.Parse([]string{arg})
	helpers.PanicErr(err)

	return c
}

func (c *nodeAPI) withArgs(args ...string) *nodeAPI {
	err := c.fs.Parse(args)
	helpers.PanicErr(err)

	return c
}

func (c *nodeAPI) withFlags(clientMethod func(*cli.Context) error, applyFlags func(*flag.FlagSet)) *nodeAPI {
	flagSetApplyFromAction(clientMethod, c.fs, "")
	applyFlags(c.fs)

	c.clientMethod = clientMethod

	return c
}

func (c *nodeAPI) exec(clientMethod ...func(*cli.Context) error) ([]byte, error) {
	if len(clientMethod) > 1 {
		PanicErr(errors.New("Only one client method allowed"))
	}

	c.output.Reset()
	defer c.output.Reset()
	defer func() {
		c.fs = flag.NewFlagSet("test", flag.ContinueOnError)
		c.clientMethod = nil
	}()

	if c.clientMethod == nil {
		c.clientMethod = clientMethod[0]
	}
	ctx := cli.NewContext(c.app, c.fs, nil)
	err := c.clientMethod(ctx)
	if err != nil {
		return nil, err
	}

	return c.output.Bytes(), nil
}

func (c *nodeAPI) mustExec(clientMethod ...func(*cli.Context) error) []byte {
	bytes, err := c.exec(clientMethod...)
	helpers.PanicErr(err)
	return bytes
}

// flagSetApplyFromAction applies the flags from action to the flagSet.
//
// `parentCommand` will filter the app commands and only applies the flags if the command/subcommand has a parent with that name, if left empty no filtering is done
//
// Taken from: https://github.com/smartcontractkit/chainlink/blob/develop/core/cmd/shell_test.go#L590
func flagSetApplyFromAction(action interface{}, flagSet *flag.FlagSet, parentCommand string) {
	cliApp := cmd.Shell{}
	app := cmd.NewApp(&cliApp)

	foundName := parentCommand == ""
	actionFuncName := getFuncName(action)

	for _, command := range app.Commands {
		flags := recursiveFindFlagsWithName(actionFuncName, command, parentCommand, foundName)

		for _, flag := range flags {
			flag.Apply(flagSet)
		}
	}
}

func recursiveFindFlagsWithName(actionFuncName string, command cli.Command, parent string, foundName bool) []cli.Flag {
	if command.Action != nil {
		if actionFuncName == getFuncName(command.Action) && foundName {
			return command.Flags
		}
	}

	for _, subcommand := range command.Subcommands {
		if !foundName {
			foundName = strings.EqualFold(subcommand.Name, parent)
		}

		found := recursiveFindFlagsWithName(actionFuncName, subcommand, parent, foundName)
		if found != nil {
			return found
		}
	}
	return nil
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func mustJSON[T any](bytes []byte) *T {
	typedPayload := new(T)
	err := json.Unmarshal(bytes, typedPayload)
	if err != nil {
		PanicErr(err)
	}
	return typedPayload
}
