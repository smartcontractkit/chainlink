package src

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	clsessions "github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/urfave/cli"
	"go.uber.org/zap/zapcore"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	clcmd "github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Package-level cache and mutex
var (
	nodeAPICache = make(map[string]*nodeAPI)
	cacheMutex   = &sync.Mutex{}
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

func newApp(n NodeWthCreds, writer io.Writer) (*clcmd.Shell, *cli.App) {
	loggingCfg := logger.Config{
		LogLevel:    zapcore.InfoLevel,
		JsonConsole: true,
	}
	logger, closeLggr := loggingCfg.New()
	u, err := url.Parse(n.RemoteURL.String())
	PanicErr(err)

	clientOpts := clcmd.ClientOpts{RemoteNodeURL: *u, InsecureSkipVerify: true}
	sr := clsessions.SessionRequest{Email: n.APILogin, Password: n.APIPassword}

	// Set the log level to error for the HTTP client, we don't care about
	// the ssl warnings it emits for CRIB
	logger.SetLogLevel(zapcore.ErrorLevel)
	cookieAuth := cmd.NewSessionCookieAuthenticator(
		clientOpts,
		&cmd.MemoryCookieStore{},
		logger,
	)
	cookieAuth.Authenticate(context.Background(), sr)
	http := cmd.NewAuthenticatedHTTPClient(
		logger,
		clientOpts,
		cookieAuth,
		sr,
	)
	// Set the log level back to info for the shell
	logger.SetLogLevel(zapcore.InfoLevel)

	client := &clcmd.Shell{
		Logger:              logger,
		Renderer:            clcmd.RendererJSON{Writer: writer},
		AppFactory:          clcmd.ChainlinkAppFactory{},
		Runner:              clcmd.ChainlinkRunner{},
		HTTP:                http,
		CookieAuthenticator: cookieAuth,
		CloseLogger:         closeLggr,
	}
	app := clcmd.NewApp(client)
	return client, app
}

type nodeAPI struct {
	methods      *cmd.Shell
	app          *cli.App
	output       *bytes.Buffer
	fs           *flag.FlagSet
	clientMethod func(*cli.Context) error
}

func newNodeAPI(n NodeWthCreds) *nodeAPI {
	// Create a unique key for the cache
	key := n.RemoteURL.String()

	// Check if the nodeAPI exists in the cache
	cacheMutex.Lock()
	if api, exists := nodeAPICache[key]; exists {
		cacheMutex.Unlock()
		return api
	}
	cacheMutex.Unlock()

	output := &bytes.Buffer{}
	methods, app := newApp(n, output)

	api := &nodeAPI{
		output:  output,
		methods: methods,
		app:     app,
		fs:      flag.NewFlagSet("test", flag.ContinueOnError),
	}

	// Store the new nodeAPI in the cache
	cacheMutex.Lock()
	nodeAPICache[key] = api
	cacheMutex.Unlock()

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
