package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/urfave/cli"
)

func removeHidden(cmds ...cli.Command) []cli.Command {
	var ret []cli.Command
	for _, cmd := range cmds {
		if cmd.Hidden {
			continue
		}
		ret = append(ret, cmd)
	}
	return ret
}

// NewApp returns the command-line parser/function-router for the given client
func NewApp(client *Client) *cli.App {
	app := cli.NewApp()
	app.Usage = "CLI for Chainlink"
	app.Version = fmt.Sprintf("%v@%v", static.Version, static.Sha)
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
	app.Commands = removeHidden([]cli.Command{
		{
			Name:  "admin",
			Usage: "Commands for remotely taking admin related actions",
			Subcommands: []cli.Command{
				{
					Name:   "chpass",
					Usage:  "Change your API password remotely",
					Action: client.ChangePassword,
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
			},
		},

		{
			Name:    "agreements",
			Aliases: []string{"agree"},
			Usage:   "Commands for handling service agreements",
			Hidden:  !client.Config.Dev(),
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Creates a Service Agreement",
					Action: client.CreateServiceAgreement,
				},
			},
		},

		{
			Name:    "attempts",
			Aliases: []string{"txas"},
			Usage:   "Commands for managing Ethereum Transaction Attempts",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "List the Transaction Attempts in descending order",
					Action: client.IndexTxAttempts,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "page",
							Usage: "page of results to display",
						},
					},
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
					Name:   "list",
					Usage:  "Show the node's environment variables",
					Action: client.GetConfiguration,
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
					Name:   "loglevel",
					Usage:  "Set log level",
					Action: client.SetLogLevel,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "level",
							Usage: "set log level for node (debug||info||warn||error)",
						},
					},
				},
				{
					Name:   "logpkg",
					Usage:  "Set package specific logging",
					Action: client.SetLogPkg,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "pkg",
							Usage: "set log filter for package specific logging",
						},
						cli.StringFlag{
							Name:  "level",
							Usage: "set log level for specified pkg",
						},
					},
				},
				{
					Name:   "logsql",
					Usage:  "Enable/disable sql statement logging",
					Action: client.SetLogSQL,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "enable",
							Usage: "enable sql logging",
						},
						cli.BoolFlag{
							Name:  "disable",
							Usage: "disable sql logging",
						},
					},
				},
			},
		},

		{
			Name:  "job_specs",
			Usage: "Commands for managing Job Specs (jobs V1)",
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
					Usage:  "Show a specific Job's details",
					Action: client.ShowJobSpec,
				},
				{
					Name:   "create",
					Usage:  "Create Job from a Job Specification JSON",
					Action: client.CreateJobSpec,
				},
				{
					Name:   "archive",
					Usage:  "Archive a Job and all its associated Runs",
					Action: client.ArchiveJobSpec,
				},
			},
		},
		{
			Name:  "jobs",
			Usage: "Commands for managing Jobs (V2)",
			Subcommands: []cli.Command{
				{
					Name:   "list",
					Usage:  "List all V2 jobs",
					Action: client.ListJobsV2,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "page",
							Usage: "page of results to display",
						},
					},
				},
				{
					Name:   "create",
					Usage:  "Create a V2 job",
					Action: client.CreateJobV2,
				},
				{
					Name:   "delete",
					Usage:  "Delete a V2 job",
					Action: client.DeleteJobV2,
				},
				{
					Name:   "run",
					Usage:  "Trigger a V2 job run",
					Action: client.TriggerPipelineRun,
				},
				{
					Name:   "migrate",
					Usage:  "Migrate a V1 job (JSON) to a V2 job (TOML)",
					Action: client.Migrate,
				},
			},
		},
		{
			Name:  "keys",
			Usage: "Commands for managing various types of keys used by the Chainlink node",
			Subcommands: []cli.Command{
				{
					Name:  "eth",
					Usage: "Local commands for administering the node's Ethereum keys",
					Subcommands: cli.Commands{
						{
							Name:   "create",
							Usage:  "Create an key in the node's keystore alongside the existing key; to create an original key, just run the node",
							Action: client.CreateETHKey,
						},
						{
							Name:   "list",
							Usage:  "List available Ethereum accounts with their ETH & LINK balances, nonces, and other metadata",
							Action: client.ListETHKeys,
						},
						{
							Name:  "delete",
							Usage: format(`Delete the ETH key by address`),
							Flags: []cli.Flag{
								cli.BoolFlag{
									Name:  "yes, y",
									Usage: "skip the confirmation prompt",
								},
								cli.BoolFlag{
									Name:  "hard",
									Usage: "hard-delete the key instead of archiving (irreversible!)",
								},
							},
							Action: client.DeleteETHKey,
						},
						{
							Name:  "import",
							Usage: format(`Import an ETH key from a JSON file`),
							Flags: []cli.Flag{
								cli.StringFlag{
									Name:  "oldpassword, p",
									Usage: "`FILE` containing the password used to encrypt the key in the JSON file",
								},
							},
							Action: client.ImportETHKey,
						},
						{
							Name:  "export",
							Usage: format(`Exports an ETH key to a JSON file`),
							Flags: []cli.Flag{
								cli.StringFlag{
									Name:  "newpassword, p",
									Usage: "`FILE` containing the password to encrypt the key (required)",
								},
								cli.StringFlag{
									Name:  "output, o",
									Usage: "Path where the JSON file will be saved (required)",
								},
							},
							Action: client.ExportETHKey,
						},
					},
				},

				{
					Name:  "p2p",
					Usage: "Remote commands for administering the node's p2p keys",
					Subcommands: cli.Commands{
						{
							Name:   "create",
							Usage:  format(`Create a p2p key, encrypted with password from the password file, and store it in the database.`),
							Action: client.CreateP2PKey,
						},
						{
							Name:  "delete",
							Usage: format(`Delete the encrypted P2P key by id`),
							Flags: []cli.Flag{
								cli.BoolFlag{
									Name:  "yes, y",
									Usage: "skip the confirmation prompt",
								},
								cli.BoolFlag{
									Name:  "hard",
									Usage: "hard-delete the key instead of archiving (irreversible!)",
								},
							},
							Action: client.DeleteP2PKey,
						},
						{
							Name:   "list",
							Usage:  format(`List available P2P keys`),
							Action: client.ListP2PKeys,
						},
						{
							Name:  "import",
							Usage: format(`Imports a P2P key from a JSON file`),
							Flags: []cli.Flag{
								cli.StringFlag{
									Name:  "oldpassword, p",
									Usage: "`FILE` containing the password used to encrypt the key in the JSON file",
								},
							},
							Action: client.ImportP2PKey,
						},
						{
							Name:  "export",
							Usage: format(`Exports a P2P key to a JSON file`),
							Flags: []cli.Flag{
								cli.StringFlag{
									Name:  "newpassword, p",
									Usage: "`FILE` containing the password to encrypt the key (required)",
								},
								cli.StringFlag{
									Name:  "output, o",
									Usage: "`FILE` where the JSON file will be saved (required)",
								},
							},
							Action: client.ExportP2PKey,
						},
					},
				},

				{
					Name:  "ocr",
					Usage: "Remote commands for administering the node's off chain reporting keys",
					Subcommands: cli.Commands{
						{
							Name:   "create",
							Usage:  format(`Create an OCR key bundle, encrypted with password from the password file, and store it in the database`),
							Action: client.CreateOCRKeyBundle,
						},
						{
							Name:  "delete",
							Usage: format(`Deletes the encrypted OCR key bundle matching the given ID`),
							Flags: []cli.Flag{
								cli.BoolFlag{
									Name:  "yes, y",
									Usage: "skip the confirmation prompt",
								},
								cli.BoolFlag{
									Name:  "hard",
									Usage: "hard-delete the key instead of archiving (irreversible!)",
								},
							},
							Action: client.DeleteOCRKeyBundle,
						},
						{
							Name:   "list",
							Usage:  format(`List available OCR key bundles`),
							Action: client.ListOCRKeyBundles,
						},
						{
							Name:  "import",
							Usage: format(`Imports an OCR key bundle from a JSON file`),
							Flags: []cli.Flag{
								cli.StringFlag{
									Name:  "oldpassword, p",
									Usage: "`FILE` containing the password used to encrypt the key in the JSON file",
								},
							},
							Action: client.ImportOCRKey,
						},
						{
							Name:  "export",
							Usage: format(`Exports an OCR key bundle to a JSON file`),
							Flags: []cli.Flag{
								cli.StringFlag{
									Name:  "newpassword, p",
									Usage: "`FILE` containing the password to encrypt the key (required)",
								},
								cli.StringFlag{
									Name:  "output, o",
									Usage: "`FILE` where the JSON file will be saved (required)",
								},
							},
							Action: client.ExportOCRKey,
						},
					},
				},

				{
					Name: "vrf",
					Usage: format(`Local commands for administering the database of VRF proof
           keys. These commands will not affect the extant in-memory keys of
           any live node.`),
					Subcommands: cli.Commands{
						{
							Name: "create",
							Usage: format(`Create a VRF key, encrypted with password from the
               password file, and store it in the database.`),
							Flags:  flags("password, p"),
							Action: client.CreateVRFKey,
						},
						{
							Name:   "import",
							Usage:  "Import key from keyfile.",
							Flags:  append(flags("password, p"), flags("file, f")...),
							Action: client.ImportVRFKey,
						},
						{
							Name:   "export",
							Usage:  "Export key to keyfile.",
							Flags:  append(flags("file, f"), flags("publicKey, pk")...),
							Action: client.ExportVRFKey,
						},
						{
							Name:  "delete",
							Usage: "Remove key from database, if present",
							Flags: []cli.Flag{
								cli.StringFlag{Name: "publicKey, pk"},
								cli.BoolFlag{
									Name:  "yes, y",
									Usage: "skip the confirmation prompt",
								},
								cli.BoolFlag{
									Name:  "hard",
									Usage: "hard-delete the key instead of archiving (irreversible!)",
								},
							},
							Action: client.DeleteVRFKey,
						},
						{
							Name: "list", Usage: "List the public keys in the db",
							Action: client.ListVRFKeys,
						},
						{
							Name: "",
						},
						{
							Name: "xxxCreateWeakKeyPeriodYesIReallyKnowWhatIAmDoingAndDoNotCareAboutThisKeyMaterialFallingIntoTheWrongHandsExclamationPointExclamationPointExclamationPointExclamationPointIAmAMasochistExclamationPointExclamationPointExclamationPointExclamationPointExclamationPoint",
							Usage: format(`
                               For testing purposes ONLY! DO NOT USE FOR ANY OTHER PURPOSE!

                               Creates a key with weak key-derivation-function parameters, so that it can be
                               decrypted quickly during tests. As a result, it would be cheap to brute-force
                               the encryption password for the key, if the ciphertext fell into the wrong
                               hands!`),
							Flags:  append(flags("password, p"), flags("file, f")...),
							Action: client.CreateAndExportWeakVRFKey,
							Hidden: !client.Config.Dev(), // For when this suite gets promoted out of dev mode
						},
					},
				},
			},
		},
		{
			Name:        "node",
			Aliases:     []string{"local"},
			Usage:       "Commands for admin actions that must be run locally",
			Description: "Commands can only be run from on the same machine as the Chainlink node.",
			Subcommands: []cli.Command{
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
				{
					Name:   "setnextnonce",
					Usage:  "Manually set the next nonce for a key. This should NEVER be necessary during normal operation. USE WITH CAUTION: Setting this incorrectly can break your node.",
					Action: client.SetNextNonce,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "address",
							Usage: "address of the key for which to set the nonce",
						},
						cli.Uint64Flag{
							Name:  "nextNonce",
							Usage: "the next nonce in the sequence",
						},
					},
				},
				{
					Name:    "start",
					Aliases: []string{"node", "n"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "api, a",
							Usage: "text file holding the API email and password, each on a line",
						},
						cli.BoolFlag{
							Name:  "debug, d",
							Usage: "set logger level to debug",
						},
						cli.StringFlag{
							Name:  "password, p",
							Usage: "text file holding the password for the node's account",
						},
						cli.StringFlag{
							Name:  "vrfpassword, vp",
							Usage: "textfile holding the password for the vrf keys; enables chainlink VRF oracle",
						},
						cli.Int64Flag{
							Name:  "replay-from-block, r",
							Usage: "historical block height from which to replay log-initiated jobs",
							Value: -1,
						},
					},
					Usage:  "Run the chainlink node",
					Action: client.RunNode,
				},
				{
					Name:   "rebroadcast-transactions",
					Usage:  "Manually rebroadcast txs matching nonce range with the specified gas price. This is useful in emergencies e.g. high gas prices and/or network congestion to forcibly clear out the pending TX queue",
					Action: client.RebroadcastTransactions,
					Flags: []cli.Flag{
						cli.Uint64Flag{
							Name:  "beginningNonce, b",
							Usage: "beginning of nonce range to rebroadcast",
						},
						cli.Uint64Flag{
							Name:  "endingNonce, e",
							Usage: "end of nonce range to rebroadcast (inclusive)",
						},
						cli.Uint64Flag{
							Name:  "gasPriceWei, g",
							Usage: "gas price (in Wei) to rebroadcast transactions at",
						},
						cli.StringFlag{
							Name:  "password, p",
							Usage: "text file holding the password for the node's account",
						},
						cli.StringFlag{
							Name:  "address, a",
							Usage: "The address (in hex format) for the key which we want to rebroadcast transactions",
						},
						cli.Uint64Flag{
							Name:  "gasLimit",
							Usage: "OPTIONAL: gas limit to use for each transaction ",
						},
					},
				},
				{
					Name:   "hard-reset",
					Usage:  "Removes unstarted transactions, cancels pending transactions as well as deletes job runs. Use with caution, this command cannot be reverted! Only execute when the node is not started!",
					Action: client.HardReset,
					Flags:  []cli.Flag{},
				},
				{
					Name:   "status",
					Usage:  "Displays the health of various services running inside the node.",
					Action: client.Status,
					Flags:  []cli.Flag{},
				},
				{
					Name:        "db",
					Usage:       "Commands for managing the database.",
					Description: "Potentially destructive commands for managing the database.",
					Subcommands: []cli.Command{
						{
							Name:   "reset",
							Usage:  "Drop, create and migrate database. Useful for setting up the database in order to run tests or resetting the dev database. WARNING: This will ERASE ALL DATA for the specified DATABASE_URL.",
							Hidden: !client.Config.Dev(),
							Action: client.ResetDatabase,
							Flags: []cli.Flag{
								cli.BoolFlag{
									Name:  "dangerWillRobinson",
									Usage: "set to true to enable dropping non-test databases",
								},
							},
						},
						{
							Name:   "preparetest",
							Usage:  "Reset database and load fixtures.",
							Hidden: !client.Config.Dev(),
							Action: client.PrepareTestDatabase,
							Flags:  []cli.Flag{},
						},
						{
							Name:   "version",
							Usage:  "Display the current database version.",
							Action: client.VersionDatabase,
							Flags:  []cli.Flag{},
						},
						{
							Name:   "migrate",
							Usage:  "Migrate the database to the latest version.",
							Action: client.MigrateDatabase,
							Flags:  []cli.Flag{},
						},
					},
				},
			},
		},

		{
			Name:   "initiators",
			Usage:  "Commands for managing External Initiators",
			Hidden: !client.Config.Dev() && !client.Config.FeatureExternalInitiators(),
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Create an authentication key for a user of External Initiators",
					Action: client.CreateExternalInitiator,
				},
				{
					Name:   "destroy",
					Usage:  "Remove an authentication key by name",
					Action: client.DeleteExternalInitiator,
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
				{
					Name:   "cancel",
					Usage:  "Cancel a Run with a specified ID",
					Action: client.CancelJobRun,
				},
			},
		},

		{
			Name:  "txs",
			Usage: "Commands for handling Ethereum transactions",
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Send <amount> Eth from node ETH account <fromAddress> to destination <toAddress>.",
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
	}...)
	return app
}

var whitespace = regexp.MustCompile(`\s+`)

// format returns result of replacing all whitespace in s with a single space
func format(s string) string {
	return string(whitespace.ReplaceAll([]byte(s), []byte(" ")))
}

// flags is an abbreviated way to express a CLI flag
func flags(s string) []cli.Flag { return []cli.Flag{cli.StringFlag{Name: s}} }
