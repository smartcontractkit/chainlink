package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/urfave/cli"
)

func main() {
	Run(NewProductionClient(), os.Args...)
}

func Run(client *cmd.Client, args ...string) {
	app := cli.NewApp()
	app.Usage = "CLI for Chainlink"
	app.Version = store.Version
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
					EnvVar: "PASSWORD",
				},
				cli.BoolFlag{
					Name:  "debug, d",
					Usage: "set logger level to debug",
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
	app.Run(args)
}

func NewProductionClient() *cmd.Client {
	return &cmd.Client{
		cmd.RendererTable{os.Stdout},
		store.NewConfig(),
		cmd.ChainlinkAppFactory{},
		cmd.TerminalAuthenticator{cmd.PasswordPrompter{}, os.Exit},
		cmd.ChainlinkRunner{},
	}
}
