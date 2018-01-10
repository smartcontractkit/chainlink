package main

import (
	"fmt"
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
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "Show a specific job",
			Action:  showJob,
		},
	}
	app.Run(os.Args)
}

func cliError(err error) error {
	if err != nil {
		fmt.Printf(err.Error())
	}
	return err
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
		return cliError(err)
	}
	defer resp.Body.Close()
	return cliError(utils.PrettyPrintJSON(resp.Body))
}

func showJob(c *cli.Context) error {
	cfg := store.NewConfig()
	if !c.Args().Present() {
		fmt.Println("Must pass the job id to be shown")
		return nil
	}
	resp, err := utils.BasicAuthGet(
		cfg.BasicAuthUsername,
		cfg.BasicAuthPassword,
		"http://localhost:8080/jobs/"+c.Args().First(),
	)
	if err != nil {
		return cliError(err)
	}
	defer resp.Body.Close()
	return cliError(utils.PrettyPrintJSON(resp.Body))
}
