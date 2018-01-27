package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/urfave/cli"
)

func main() {
	client := cmd.Client{cmd.RendererTable{os.Stdout}, store.NewConfig()}

	app := cli.NewApp()
	app.Usage = "CLI for Chainlink"
	app.Version = "0.2.0"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "json output as opposed to table",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("json") {
			client.Renderer = cmd.RendererJSON{os.Stdout}
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:    "node",
			Aliases: []string{"n"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "password, p",
					Usage: "password for the node's account",
				},
			},
			Usage:  "Run the chainlink node",
			Action: client.RunNode,
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
