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
			Name:    "account",
			Aliases: []string{"a"},
			Usage:   "Display the account address with its ETH & LINK balances",
			Action:  client.DisplayAccountBalance,
		},
		{
			Name:    "jobspecs",
			Aliases: []string{"jobs", "j", "specs"},
			Usage:   "Get all jobs",
			Action:  client.GetJobSpecs,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page",
					Usage: "page of results to display",
				},
			},
		},
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "Show a specific job",
			Action:  client.ShowJobSpec,
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create job spec from JSON",
			Action:  client.CreateJobSpec,
		},
		{
			Name:   "archivejob",
			Usage:  "Archive job and all associated runs",
			Action: client.ArchiveJobSpec,
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Begin job run for specid",
			Action:  client.CreateJobRun,
		},
		{
			Name:    "showrun",
			Aliases: []string{"sr"},
			Usage:   "Show a job run for a RunID",
			Action:  client.ShowJobRun,
		},
		{
			Name:    "listruns",
			Aliases: []string{"lr"},
			Usage:   "List all job runs",
			Action:  client.GetJobRuns,
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
			Name:   "bridge",
			Usage:  "Add a new bridge to the node",
			Action: client.AddBridge,
		},
		{
			Name:   "getbridges",
			Usage:  "List all bridges added to the node",
			Action: client.GetBridges,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page",
					Usage: "page of results to display",
				},
			},
		},
		{
			Name:   "showbridge",
			Usage:  "Show a specific bridge",
			Action: client.ShowBridge,
		},
		{
			Name:   "removebridge",
			Usage:  "Removes a specific bridge",
			Action: client.RemoveBridge,
		},
		{
			Name:    "externalinitiators",
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
			Usage:   "Withdraw, to an authorized Ethereum <address>, <amount> units of LINK. Withdraws from the configured oracle contract by default, or from contract optionally specified by a third command-line argument --from-oracle-contract-address=<contract address>. Address inputs must be in EIP55-compliant capitalization.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "from-oracle-contract-address",
					Usage: "address of Oracle contract to withdraw from (will use node default if unspecified)",
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
			Name:   "chpass",
			Usage:  "Change your password",
			Action: client.ChangePassword,
		},
		{
			Name:   "setgasprice",
			Usage:  "Set the minimum gas price to use for outgoing transactions",
			Action: client.SetMinimumGasPrice,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "gwei",
					Usage: "Specify amount in gwei",
				},
			},
		},
		{
			Name:   "transactions",
			Usage:  "List the transactions in descending order",
			Action: client.GetTransactions,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "page",
					Usage: "page of results to display",
				},
			},
		},
		{
			Name:   "txattempts",
			Usage:  "List the transaction attempts in descending order",
			Action: client.GetTxAttempts,
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
