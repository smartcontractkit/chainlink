package main

import (
	"os"

	"github.com/smartcontractkit/chainlink-go/logger"
	"github.com/smartcontractkit/chainlink-go/services"
	"github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/smartcontractkit/chainlink-go/web"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "CLI for Chainlink"
	app.Commands = []cli.Command{
		{
			Name:    "node",
			Aliases: []string{"n"},
			Usage:   "Run the chainlink node",
			Action:  runNode,
		},
		{
			Name:    "jobs",
			Aliases: []string{"j"},
			Usage:   "Get all jobs",
			Action:  getJobs,
		},
	}
	app.Run(os.Args)
}

func runNode(c *cli.Context) error {
	cl := services.NewApplication(store.NewConfig())
	services.Authenticate(cl.Store)
	r := web.Router(cl)

	if err := cl.Start(); err != nil {
		logger.Fatal(err)
	}
	defer cl.Stop()
	logger.Fatal(r.Run())
	return nil
}

func getJobs(c *cli.Context) error {
	cfg := store.NewConfig()
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		"http://localhost:8080/jobs",
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()
	if err = utils.PrettyPrintJSON(resp.Body); err != nil {
		logger.Fatal(err)
	}
	return nil
}
