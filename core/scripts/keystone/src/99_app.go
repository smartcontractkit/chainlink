package src

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/zap/zapcore"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	clcmd "github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	clsessions "github.com/smartcontractkit/chainlink/v2/core/sessions"
)

// Package-level cache and mutex
var (
	nodeAPICache = make(map[string]*nodeAPI)
	cacheMutex   = &sync.Mutex{}
)

func newApp(n NodeWithCreds, writer io.Writer) (*clcmd.Shell, *cli.App) {
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

	http := NewRetryableAuthenticatedHTTPClient(logger, clientOpts, cookieAuth, sr)
	// Set the log level back to info for the shell
	logger.SetLogLevel(zapcore.InfoLevel)
	client := &clcmd.Shell{
		Logger:     logger,
		Renderer:   clcmd.RendererJSON{Writer: writer},
		AppFactory: clcmd.ChainlinkAppFactory{},
		Runner:     clcmd.ChainlinkRunner{},
		HTTP:       http,

		CloseLogger: closeLggr,
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
	logger       logger.Logger
}

func newNodeAPI(n NodeWithCreds) *nodeAPI {
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
		logger:  methods.Logger.Named("NodeAPI"),
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

	defer c.output.Reset()
	defer func() {
		c.fs = flag.NewFlagSet("test", flag.ContinueOnError)
		c.clientMethod = nil
	}()

	if c.clientMethod == nil {
		c.clientMethod = clientMethod[0]
	}

	retryCount := 3
	for i := 0; i < retryCount; i++ {
		c.logger.Tracew("Attempting API request", "attempt", i+1, "maxAttempts", retryCount)
		c.output.Reset()
		ctx := cli.NewContext(c.app, c.fs, nil)
		err := c.clientMethod(ctx)

		if err == nil {
			c.logger.Tracew("API request completed successfully", "attempt", i+1)
			return c.output.Bytes(), nil
		}

		if !strings.Contains(err.Error(), "invalid character '<' looking for beginning of value") {
			c.logger.Tracew("API request failed with non-retriable error",
				"attempt", i+1,
				"err", err,
			)
			return nil, err
		}

		c.logger.Warnw("Encountered 504 gateway error during API request, retrying",
			"attempt", i+1,
			"maxAttempts", retryCount,
			"err", err,
		)

		if i == retryCount-1 {
			c.logger.Error("Failed to complete API request after all retry attempts")
			return nil, err
		}

		c.logger.Tracew("Waiting before retry attempt",
			"attempt", i+1,
			"waitTime", "1s",
		)
		time.Sleep(3 * time.Second)
	}

	return nil, errors.New("API request failed after retries")
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

type retryableAuthenticatedHTTPClient struct {
	client cmd.HTTPClient
	logger logger.Logger
}

func NewRetryableAuthenticatedHTTPClient(lggr logger.Logger, clientOpts clcmd.ClientOpts, cookieAuth cmd.CookieAuthenticator, sessionRequest clsessions.SessionRequest) cmd.HTTPClient {
	return &retryableAuthenticatedHTTPClient{
		client: cmd.NewAuthenticatedHTTPClient(lggr, clientOpts, cookieAuth, sessionRequest),
		logger: lggr.Named("RetryableAuthenticatedHTTPClient"),
	}
}

func logBody(body io.Reader) (string, io.Reader) {
	if body == nil {
		return "", nil
	}

	var buf bytes.Buffer
	tee := io.TeeReader(body, &buf)
	bodyBytes, _ := io.ReadAll(tee)
	return string(bodyBytes), bytes.NewReader(buf.Bytes())
}

func logResponse(logger logger.Logger, resp *http.Response) {
	if resp == nil {
		logger.Trace("Response was nil")
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorw("Failed to read response body for logging", "err", err)
		return
	}
	// Replace the body so it can be read again by the caller
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	logger.Tracew("Response details",
		"statusCode", resp.StatusCode,
		"status", resp.Status,
		"headers", resp.Header,
		"body", string(bodyBytes),
	)
}

func (h *retryableAuthenticatedHTTPClient) Get(ctx context.Context, path string, headers ...map[string]string) (*http.Response, error) {
	h.logger.Tracew("Making GET request",
		"path", path,
		"headers", headers,
	)
	return h.doRequestWithRetry(ctx, func() (*http.Response, error) {
		return h.client.Get(ctx, path, headers...)
	})
}

func (h *retryableAuthenticatedHTTPClient) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	bodyStr, newBody := logBody(body)
	h.logger.Tracew("Making POST request",
		"path", path,
		"body", bodyStr,
	)
	return h.doRequestWithRetry(ctx, func() (*http.Response, error) {
		return h.client.Post(ctx, path, newBody)
	})
}

func (h *retryableAuthenticatedHTTPClient) Put(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	bodyStr, newBody := logBody(body)
	h.logger.Tracew("Making PUT request",
		"path", path,
		"body", bodyStr,
	)
	return h.doRequestWithRetry(ctx, func() (*http.Response, error) {
		return h.client.Put(ctx, path, newBody)
	})
}

func (h *retryableAuthenticatedHTTPClient) Patch(ctx context.Context, path string, body io.Reader, headers ...map[string]string) (*http.Response, error) {
	bodyStr, newBody := logBody(body)
	h.logger.Tracew("Making PATCH request",
		"path", path,
		"headers", headers,
		"body", bodyStr,
	)
	return h.doRequestWithRetry(ctx, func() (*http.Response, error) {
		return h.client.Patch(ctx, path, newBody, headers...)
	})
}

func (h *retryableAuthenticatedHTTPClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	h.logger.Tracew("Making DELETE request",
		"path", path,
	)
	return h.doRequestWithRetry(ctx, func() (*http.Response, error) {
		return h.client.Delete(ctx, path)
	})
}

func (h *retryableAuthenticatedHTTPClient) doRequestWithRetry(_ context.Context, req func() (*http.Response, error)) (*http.Response, error) {
	retryCount := 3
	for i := 0; i < retryCount; i++ {
		h.logger.Tracew("Attempting request", "attempt", i+1, "maxAttempts", retryCount)

		response, err := req()
		logResponse(h.logger, response)

		if err == nil || !strings.Contains(err.Error(), "invalid character '<' looking for beginning of value") {
			if err != nil {
				h.logger.Warn("Request completed with error",
					"attempt", i+1,
					"err", err,
				)
			} else {
				h.logger.Tracew("Request completed successfully",
					"attempt", i+1,
					"statusCode", response.StatusCode,
				)
			}
			return response, err
		}

		h.logger.Warnw("Encountered 504 error during request, retrying",
			"attempt", i+1,
			"maxAttempts", retryCount,
			"err", err,
		)

		if i == retryCount-1 {
			h.logger.Error("Failed to complete request after all retry attempts")
			return response, err
		}

		h.logger.Tracew("Waiting before retry attempt",
			"attempt", i+1,
			"waitTime", "1s",
		)
		time.Sleep(1 * time.Second)
	}
	return nil, errors.New("request failed after retries")
}
