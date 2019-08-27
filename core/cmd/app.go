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
			Name:  "admin",
			Usage: "Commands for remotely taking admin related actions",
			Subcommands: []cli.Command{
				{
					Name:   "chpass",
					Usage:  "Change your account password remotely",
					Action: client.ChangePassword,
				},
				{
					Name:   "info",
					Usage:  "Display the Account's address with its ETH & LINK balances",
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
					Name:        "withdraw",
					Usage:       "Withdraw to <address>, <amount> units of LINK from the configured Oracle Contract",
					Description: "Only works if the Chainlink node is the owner of the contract being withdrawn from",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "from",
							Usage: "override the configured oracle address to withdraw from",
						},
					},
					Action: client.Withdraw,
				},
			},
		},

		{
			Name:  "bridges",
			Usage: "Commands for Bridges communicating with External Adapters",
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Create a new Bridge to an External Adapter",
					Action: client.CreateBridge,
				},
				{
					Name:   "destroy",
					Usage:  "Destroys the Bridge for an External Adapter",
					Action: client.RemoveBridge,
				},
				{
					Name:   "list",
					Usage:  "List all Bridges to External Adapters",
					Action: client.IndexBridges,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "page",
							Usage: "page of results to display",
						},
					},
				},
				{
					Name:   "show",
					Usage:  "Show an Bridge's details",
					Action: client.ShowBridge,
				},
			},
		},

		{
			Name:  "config",
			Usage: "Commands for the node's configuration",
			Subcommands: []cli.Command{
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
			},
		},

		{
			Name:  "jobs",
			Usage: "Commands for managing Jobs",
			Subcommands: []cli.Command{
				{
					Name:   "archive",
					Usage:  "Archive a Job and all its associated Runs",
					Action: client.ArchiveJobSpec,
				},
				{
					Name:   "create",
					Usage:  "Create Job from a Job Specification JSON",
					Action: client.CreateJobSpec,
				},
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
					Usage:  "Show a specific Job's details",
					Action: client.ShowJobSpec,
				},
			},
		},

		{
			Name:        "node",
			Usage:       "Commands for admin actions that must be run locally",
			Description: "Commands can only be run from on the same machine as the Chainlink node.",
			Subcommands: []cli.Command{
				{
					Name:    "start",
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
					Name:        "deleteuser",
					Usage:       "Erase the *local node's* user and corresponding session to force recreation on next node launch.",
					Description: "Does not work remotely over API.",
					Action:      client.DeleteUser,
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
			Name:  "runs",
			Usage: "Commands for managing Runs",
			Subcommands: []cli.Command{
				{
					Name:        "create",
					Aliases:     []string{"r"},
					Usage:       "Create a new Run for a Job given an Job ID and optional JSON body",
					Description: "Takes a Job ID and a JSON string or path to a JSON file",
					Action:      client.CreateJobRun,
				},
				{
					Name:   "list",
					Usage:  "List all Runs",
					Action: client.IndexJobRuns,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "page",
							Usage: "page of results to display",
						},
						cli.StringFlag{
							Name:  "jobid",
							Usage: "filter all Runs to match the given jobid",
						},
					},
				},
				{
					Name:    "show",
					Aliases: []string{"sr"},
					Usage:   "Show a Run for a specific ID",
					Action:  client.ShowJobRun,
				},
			},
		},

		{
			Name:  "txs",
			Usage: "Commands for handling Ethereum transactions",
			Subcommands: []cli.Command{
				{
					Name:  "create",
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
					Name:   "list",
					Usage:  "List the Ethereum Transactions in descending order",
					Action: client.IndexTransactions,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "page",
							Usage: "page of results to display",
						},
					},
				},
				{
					Name:   "show",
					Usage:  "get information on a specific Ethereum Transaction",
					Action: client.ShowTransaction,
				},
			},
		},
	}

	if client.Config.Dev() {
		app.Commands = append(app.Commands, cli.Command{
			Name:    "agreements",
			Aliases: []string{"agree"},
			Usage:   "Commands for handling service agreements",
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Creates a Service Agreement",
					Action: client.CreateServiceAgreement,
				},
			},
		},

			cli.Command{
				Name:  "attempts",
				Usage: "Commands for managing Ethereum Transaction Attempts",
				Subcommands: []cli.Command{
					{
						Name:    "list",
						Aliases: []string{"txas"},
						Usage:   "List the Transaction Attempts in descending order",
						Action:  client.IndexTxAttempts,
						Flags: []cli.Flag{
							cli.IntFlag{
								Name:  "page",
								Usage: "page of results to display",
							},
						},
					},
				},
			},

			cli.Command{
				Name:   "createextrakey",
				Usage:  "Create a key in the node's keystore alongside the existing key; to create an original key, just run the node",
				Action: client.CreateExtraKey,
			},

			cli.Command{
				Name:  "initiators",
				Usage: "Commands for managing External Initiators",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create an authentication key for a user of External Initiators",
						Action: client.CreateExternalInitiator,
					},
					{
						Name:   "destroy",
						Usage:  "Remove an authentication key",
						Action: client.DeleteExternalInitiator,
					},
				},
			})
	}

	return app
}
