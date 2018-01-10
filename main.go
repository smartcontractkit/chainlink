package main

import (
	"os"

	"github.com/smartcontractkit/chainlink-go/commands"
	"github.com/urfave/cli"
)

func main() {
	client := commands.Client{os.Stdout}
	app := cli.NewApp()
	app.Usage = "CLI for Chainlink"
	app.Commands = []cli.Command{
		{
			Name:    "node",
			Aliases: []string{"n"},
			Usage:   "Run the chainlink node",
			Action:  client.RunNode,
		},
		{
			Name:    "jobs",
			Aliases: []string{"j"},
			Usage:   "Get all jobs",
			Action:  client.GetJobs,
		},
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "Show a specific job",
			Action:  client.ShowJob,
		},
	}
	app.Run(os.Args)
}
