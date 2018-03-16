package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/services"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/smartcontractkit/chainlink/web"
	"github.com/tidwall/gjson"
	clipkg "github.com/urfave/cli"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
)

// Client is the shell for the node. It has fields for the Renderer,
// Config, AppFactory (the services application), Authenticator, and Runner.
type Client struct {
	Renderer
	Config     strpkg.Config
	AppFactory AppFactory
	Auth       Authenticator
	Runner     Runner
}

// RunNode starts the Chainlink core.
func (cli *Client) RunNode(c *clipkg.Context) error {
	if c.Bool("debug") {
		cli.Config.LogLevel = strpkg.LogLevel{zapcore.DebugLevel}
	}
	logger.Infow("Starting Chainlink Node " + strpkg.Version + " at commit " + strpkg.Sha)
	app := cli.AppFactory.NewApplication(cli.Config)
	store := app.GetStore()
	cli.Auth.Authenticate(store, c.String("password"))
	if err := app.Start(); err != nil {
		return cli.errorOut(err)
	}
	defer app.Stop()
	logNodeBalance(store)
	return cli.errorOut(cli.Runner.Run(app))
}

func logNodeBalance(store *strpkg.Store) {
	balance, err := presenters.ShowEthBalance(store)
	logger.WarnIf(err)
	logger.Infow(balance)
	balance, err = presenters.ShowLinkBalance(store)
	logger.WarnIf(err)
	logger.Infow(balance)
}

// ShowJobSpec returns the status of the given JobID.
func (cli *Client) ShowJobSpec(c *clipkg.Context) error {
	cfg := cli.Config
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the job id to be shown"))
	}
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/specs/"+c.Args().First(),
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var job presenters.JobSpec
	return cli.deserializeResponse(resp, &job)
}

// GetJobSpecs returns all job specs.
func (cli *Client) GetJobSpecs(c *clipkg.Context) error {
	cfg := cli.Config
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/specs",
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var jobs []models.JobSpec
	return cli.deserializeResponse(resp, &jobs)
}

// CreateJobSpec creates job spec based on JSON input
func (cli *Client) CreateJobSpec(c *clipkg.Context) error {
	cfg := cli.Config
	errjs := errors.New("Must pass in JSON or filepath")
	if !c.Args().Present() {
		return cli.errorOut(errjs)
	}
	arg := c.Args().First()
	var buf *bytes.Buffer
	var err error
	if gjson.Valid(arg) {
		buf = bytes.NewBufferString(arg)
	} else if buf, err = fromFile(arg, cli); err != nil {
		return cli.errorOut(multierr.Append(errjs, err))
	}
	resp, err := utils.BasicAuthPost(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/specs",
		"application/json",
		buf,
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var jobs presenters.JobSpec
	return cli.deserializeResponse(resp, &jobs)
}

func fromFile(arg string, cli *Client) (*bytes.Buffer, error) {
	dir, err := homedir.Expand(arg)
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadFile(dir)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(file), nil
}

func (cli *Client) deserializeResponse(resp *http.Response, dst interface{}) error {
	if resp.StatusCode >= 400 {
		return cli.errorOut(errors.New(resp.Status))
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cli.errorOut(err)
	}
	if err = json.Unmarshal(b, &dst); err != nil {
		return cli.errorOut(err)
	}
	return cli.errorOut(cli.Render(dst))
}

func (cli *Client) errorOut(err error) error {
	if err != nil {
		return clipkg.NewExitError(err.Error(), 1)
	}
	return nil
}

// AppFactory implements the NewApplication method.
type AppFactory interface {
	NewApplication(strpkg.Config) services.Application
}

// ChainlinkAppFactory is used to create a new Application.
type ChainlinkAppFactory struct{}

// NewApplication returns a new instance of the node with the given config.
func (n ChainlinkAppFactory) NewApplication(config strpkg.Config) services.Application {
	return services.NewApplication(config)
}

// Runner implements the Run method.
type Runner interface {
	Run(services.Application) error
}

// ChainlinkRunner is used to run the node application.
type ChainlinkRunner struct{}

// Run sets the log level based on config and starts the web router to listen
// for input and return data.
func (n ChainlinkRunner) Run(app services.Application) error {
	gin.SetMode(app.GetStore().Config.LogLevel.ForGin())
	port := app.GetStore().Config.Port
	return web.Router(app.(*services.ChainlinkApplication)).Run(":" + port)
}
