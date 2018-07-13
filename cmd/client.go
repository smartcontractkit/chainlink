package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink/services"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/web"
	clipkg "github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
)

// Client is the shell for the node. It has fields for the Renderer,
// Config, AppFactory (the services application), Authenticator, and Runner.
type Client struct {
	Renderer
	Config          strpkg.Config
	AppFactory      AppFactory
	Auth            Authenticator
	UserInitializer UserInitializer
	Runner          Runner
	RemoteClient    RemoteClient
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
	server := web.Router(app.(*services.ChainlinkApplication))
	config := app.GetStore().Config
	var g errgroup.Group

	if config.Dev {
		g.Go(func() error { return server.Run(":" + config.Port) })
	} else {
		certFile := config.CertFile()
		keyFile := config.KeyFile()
		g.Go(func() error { return server.RunTLS(":"+config.Port, certFile, keyFile) })
	}
	return g.Wait()
}
