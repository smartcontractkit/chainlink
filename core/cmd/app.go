package cmd

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/urfave/cli"
)

// NewApp returns the command-line parser/function-router for the given client
func NewApp(client *Client) *cli.App {
	app := cli.NewApp()
	app.Usage = "CLI for Chainlink"
	app.Version = fmt.Sprintf("%v@%v", store.Version, store.Sha)
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "json output as opposed to table",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("json") {
			client.Renderer = RendererJSON{Writer: os.Stdout}
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "local",
			Usage: "Commands which are run locally",
			Subcommands: []cli.Command{
				{
					Name:    "node",
					Aliases: []string{"n"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "api, a",
							Usage: "text file holding the API email and password, each on a line",
						},
						cli.StringFlag{
							Name:  "password, p",
							Usage: "text file holding the password for the node's account",
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
					Name:   "deleteuser",
					Usage:  "Erase the *local node's* user and corresponding session to force recreation on next node launch. Does not work remotely over API.",
					Action: client.DeleteUser,
				},
				{
					Name:    "import",
					Aliases: []string{"i"},
					Usage:   "Import a key file to use with the node",
					Action:  client.ImportKey,
				},
			},
		},
		{
			Name:  "account",
			Usage: "Display the account related info for remote admin access",
			Subcommands: []cli.Command{
				{
					Name:   "show",
					Usage:  "Display the account address with its ETH & LINK balances",
					Action: client.DisplayAccountBalance,
				},
				{
					Name:   "login",
					Usage:  "Login to remote client by creating a session cookie",
					Action: client.RemoteLogin,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "file, f",
							Usage: "text file holding the API email and password needed to create a session cookie",
						},
					},
				},
				{
					Name:   "chpass",
					Usage:  "Change your password",
					Action: client.ChangePassword,
				},
			},
		},
		{
			Name:  "jobs",
			Usage: "Job related commands",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "List all jobs",
					Action: client.IndexJobSpecs,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "page",
							Usage: "page of results to display",
						},
					},
				},
				{
					Name:   "show",
					Usage:  "Show a specific job's details",
					Action: client.ShowJobSpec,
				},
			},
			{
				Name:   "create",
				Usage:  "Create job spec from JSON",
				Action: client.CreateJobSpec,
			},
			{
				Name:   "archive",
				Usage:  "Archive job and all associated runs",
				Action: client.ArchiveJobSpec,
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Create a new run for a job",
			Action:  client.CreateJobRun,
		},
		{
			Name:    "showrun",
			Aliases: []string{"sr"},
			Usage:   "Show a run for a specific ID",
			Action:  client.ShowJobRun,
		},
		{
			Name:    "runs",
			Aliases: []string{"lr"},
			Usage:   "List all runs",
			Action:  client.IndexJobRuns,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page",
					Usage: "page of results to display",
				},
				cli.StringFlag{
					Name:  "jobid",
					Usage: "filter all runs to match the given jobid",
				},
			},
		},
		{
			Name:   "addbridge",
			Usage:  "Create a new bridge to the node",
			Action: client.CreateBridge,
		},
		{
			Name:   "bridges",
			Usage:  "List all bridges added to the node",
			Action: client.IndexBridges,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page",
					Usage: "page of results to display",
				},
			},
		},
		{
			Name:   "bridge",
			Usage:  "Show a specific bridge",
			Action: client.ShowBridge,
		},
		{
			Name:   "removebridge",
			Usage:  "Removes a specific bridge",
			Action: client.RemoveBridge,
		},
		{
			Name:    "initiators",
			Aliases: []string{"exi"},
			Usage:   "Tasks for managing external initiators",
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Create an authentication key for a user of external initiators",
					Action: client.CreateExternalInitiator,
				},
				{
					Name:   "delete",
					Usage:  "Remove an authentication key",
					Action: client.DeleteExternalInitiator,
				},
			},
		},
		{
			Name:    "agree",
			Aliases: []string{"createsa"},
			Usage:   "Creates a service agreement",
			Action:  client.CreateServiceAgreement,
		},
		{
			Name:    "withdraw",
			Aliases: []string{"w"},
			Usage:   "Withdraw to <address>, <amount> units of LINK from the configured oracle",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "from-oracle-contract-address",
					Usage: "override the configured oracle address to withdraw from",
				},
			},
			Action: client.Withdraw,
		},
		{
			Name:  "sendether",
			Usage: "Send <amount> ETH from the node's ETH account to an <address>.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "from, f",
					Usage: "optional flag to specify which address should send the transaction",
				},
			},
			Action: client.SendEther,
		},
		{
			Name:   "setgasprice",
			Usage:  "Set the minimum gas price to use for outgoing transactions",
			Action: client.SetMinimumGasPrice,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "gwei",
					Usage: "Specify amount in gwei",
				},
			},
		},
		{
			Name:    "transactions",
			Aliases: []string{"txs"},
			Usage:   "List the transactions in descending order",
			Action:  client.IndexTransactions,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page",
					Usage: "page of results to display",
				},
			},
		},
		{
			Name:    "transaction",
			Aliases: []string{"tx"},
			Usage:   "get information on a specific transaction",
			Action:  client.ShowTransaction,
		},
		{
			Name:    "txattempts",
			Aliases: []string{"txas"},
			Usage:   "List the transaction attempts in descending order",
			Action:  client.IndexTxAttempts,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page",
					Usage: "page of results to display",
				},
			},
		},
	}

	if client.Config.Dev() {
		createextrakey := cli.Command{
			Name:   "createextrakey",
			Usage:  "Create a key in the node's keystore alongside the existing key; to create an original key, just run the node",
			Action: client.CreateExtraKey,
		}
		app.Commands = append(app.Commands, createextrakey)
	}

	return app
}
