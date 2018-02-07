package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/smartcontractkit/chainlink/web"
	clipkg "github.com/urfave/cli"
	"go.uber.org/zap/zapcore"
)

type Client struct {
	Renderer
	Config     store.Config
	AppFactory AppFactory
	Auth       Authenticator
	Runner     Runner
}

func (cli *Client) RunNode(c *clipkg.Context) error {
	if c.Bool("debug") {
		cli.Config.LogLevel = store.LogLevel{zapcore.DebugLevel}
	}
	app := cli.AppFactory.NewApplication(cli.Config)
	cli.Auth.Authenticate(app.GetStore(), c.String("password"))

	if err := app.Start(); err != nil {
		return cli.errorOut(err)
	}
	defer app.Stop()
	return cli.errorOut(cli.Runner.Run(app))
}

func (cli *Client) ShowJob(c *clipkg.Context) error {
	cfg := cli.Config
	if !c.Args().Present() {
		return cli.errorOut(errors.New("Must pass the job id to be shown"))
	}
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/jobs/"+c.Args().First(),
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()
	var job presenters.Job
	return cli.deserializeResponse(resp, &job)
}

func (cli *Client) GetJobs(c *clipkg.Context) error {
	cfg := cli.Config
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		cfg.ClientNodeURL+"/v2/jobs",
	)
	if err != nil {
		return cli.errorOut(err)
	}
	defer resp.Body.Close()

	var jobs []models.Job
	return cli.deserializeResponse(resp, &jobs)
}

func (cli *Client) deserializeResponse(resp *http.Response, dst interface{}) error {
	if resp.StatusCode >= 300 {
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

type AppFactory interface {
	NewApplication(store.Config) services.Application
}

type ChainlinkAppFactory struct{}

func (n ChainlinkAppFactory) NewApplication(config store.Config) services.Application {
	return services.NewApplication(config)
}

type Runner interface {
	Run(services.Application) error
}

type NodeRunner struct{}

func (n NodeRunner) Run(app services.Application) error {
	gin.SetMode(app.GetStore().Config.LogLevel.ForGin())
	return web.Router(app.(*services.ChainlinkApplication)).Run()
}
