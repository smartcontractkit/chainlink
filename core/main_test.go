//go:build !windows
// +build !windows

package main

import (
	"io/ioutil"
	"testing"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func run(args ...string) {
	t := &testing.T{}
	tc := cltest.NewTestGeneralConfig(t)
	tc.Overrides.Dev = null.BoolFrom(false)
	testClient := &cmd.Client{
		Renderer:               cmd.RendererTable{Writer: ioutil.Discard},
		Config:                 tc,
		Logger:                 logger.TestLogger(t),
		AppFactory:             cmd.ChainlinkAppFactory{},
		FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
		Runner:                 cmd.ChainlinkRunner{},
		HTTP:                   cltest.NewMockAuthenticatedHTTPClient(tc, "session"),
		ChangePasswordPrompter: cltest.MockChangePasswordPrompter{},
	}
	args = append([]string{""}, args...)
	Run(testClient, args...)
}

func ExampleRun() {
	run("--help")
	run("--version")
	// Output:
	// NAME:
	//    core.test - CLI for Chainlink
	//
	// USAGE:
	//    core.test [global options] command [command options] [arguments...]
	//
	// VERSION:
	//    unset@unset
	//
	// COMMANDS:
	//    admin           Commands for remotely taking admin related actions
	//    attempts, txas  Commands for managing Ethereum Transaction Attempts
	//    blocks          Commands for managing blocks
	//    bridges         Commands for Bridges communicating with External Adapters
	//    config          Commands for the node's configuration
	//    jobs            Commands for managing Jobs
	//    keys            Commands for managing various types of keys used by the Chainlink node
	//    node, local     Commands for admin actions that must be run locally
	//    txs             Commands for handling Ethereum transactions
	//    chains          Commands for handling chain configuration
	//    nodes           Commands for handling node configuration
	//    help, h         Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --json, -j     json output as opposed to table
	//    --help, -h     show help
	//    --version, -v  print the version
	// core.test version unset@unset
}

func ExampleRun_admin() {
	run("admin", "--help")
	// Output:
	// NAME:
	//    core.test admin - Commands for remotely taking admin related actions
	//
	// USAGE:
	//    core.test admin command [command options] [arguments...]
	//
	// COMMANDS:
	//    chpass  Change your API password remotely
	//    login   Login to remote client by creating a session cookie
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_attempts() {
	run("attempts", "--help")
	// Output:
	// NAME:
	//    core.test attempts - Commands for managing Ethereum Transaction Attempts
	//
	// USAGE:
	//    core.test attempts command [command options] [arguments...]
	//
	// COMMANDS:
	//    list  List the Transaction Attempts in descending order
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_blocks() {
	run("blocks", "--help")
	// Output:
	// NAME:
	//    core.test blocks - Commands for managing blocks
	//
	// USAGE:
	//    core.test blocks command [command options] [arguments...]
	//
	// COMMANDS:
	//    replay  Replays block data from the given number
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_bridges() {
	run("bridges", "--help")
	// Output:
	// NAME:
	//    core.test bridges - Commands for Bridges communicating with External Adapters
	//
	// USAGE:
	//    core.test bridges command [command options] [arguments...]
	//
	// COMMANDS:
	//    create   Create a new Bridge to an External Adapter
	//    destroy  Destroys the Bridge for an External Adapter
	//    list     List all Bridges to External Adapters
	//    show     Show an Bridge's details
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_config() {
	run("config", "--help")
	// Output:
	// NAME:
	//    core.test config - Commands for the node's configuration
	//
	// USAGE:
	//    core.test config command [command options] [arguments...]
	//
	// COMMANDS:
	//    list         Show the node's environment variables
	//    setgasprice  Set the default gas price to use for outgoing transactions
	//    loglevel     Set log level
	//    logpkg       Set package specific logging
	//    logsql       Enable/disable sql statement logging
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_jobs() {
	run("jobs", "--help")
	// Output:
	// NAME:
	//    core.test jobs - Commands for managing Jobs
	//
	// USAGE:
	//    core.test jobs command [command options] [arguments...]
	//
	// COMMANDS:
	//    list    List all jobs
	//    show    Show a job
	//    create  Create a job
	//    delete  Delete a job
	//    run     Trigger a job run
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys() {
	run("keys", "--help")
	// Output:
	// NAME:
	//    core.test keys - Commands for managing various types of keys used by the Chainlink node
	//
	// USAGE:
	//    core.test keys command [command options] [arguments...]
	//
	// COMMANDS:
	//    eth  Remote commands for administering the node's Ethereum keys
	//    p2p  Remote commands for administering the node's p2p keys
	//    csa  Remote commands for administering the node's CSA keys
	//    ocr  Remote commands for administering the node's off chain reporting keys
	//    vrf  Remote commands for administering the node's vrf keys
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys_eth() {
	run("keys", "eth", "--help")
	// Output:
	// NAME:
	//    core.test keys eth - Remote commands for administering the node's Ethereum keys
	//
	// USAGE:
	//    core.test keys eth command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create an key in the node's keystore alongside the existing key; to create an original key, just run the node
	//    update  Update the existing key's parameters
	//    list    List available Ethereum accounts with their ETH & LINK balances, nonces, and other metadata
	//    delete  Delete the ETH key by address
	//    import  Import an ETH key from a JSON file
	//    export  Exports an ETH key to a JSON file
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys_p2p() {
	run("keys", "p2p", "--help")
	// Output:
	// NAME:
	//    core.test keys p2p - Remote commands for administering the node's p2p keys
	//
	// USAGE:
	//    core.test keys p2p command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a p2p key, encrypted with password from the password file, and store it in the database.
	//    delete  Delete the encrypted P2P key by id
	//    list    List available P2P keys
	//    import  Imports a P2P key from a JSON file
	//    export  Exports a P2P key to a JSON file
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys_csa() {
	run("keys", "csa", "--help")
	// Output:
	// NAME:
	//    core.test keys csa - Remote commands for administering the node's CSA keys
	//
	// USAGE:
	//    core.test keys csa command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a CSA key, encrypted with password from the password file, and store it in the database.
	//    list    List available CSA keys
	//    import  Imports a CSA key from a JSON file.
	//    export  Exports an existing CSA key by its ID.
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys_ocr() {
	run("keys", "ocr", "--help")
	// Output:
	// NAME:
	//    core.test keys ocr - Remote commands for administering the node's off chain reporting keys
	//
	// USAGE:
	//    core.test keys ocr command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create an OCR key bundle, encrypted with password from the password file, and store it in the database
	//    delete  Deletes the encrypted OCR key bundle matching the given ID
	//    list    List available OCR key bundles
	//    import  Imports an OCR key bundle from a JSON file
	//    export  Exports an OCR key bundle to a JSON file
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys_vrf() {
	run("keys", "vrf", "--help")
	// Output:
	// NAME:
	//    core.test keys vrf - Remote commands for administering the node's vrf keys
	//
	// USAGE:
	//    core.test keys vrf command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a VRF key
	//    import  Import VRF key from keyfile
	//    export  Export VRF key to keyfile
	//    delete  Archive or delete VRF key from memory and the database, if present. Note that jobs referencing the removed key will also be removed.
	//    list    List the VRF keys
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_node() {
	run("node", "--help")
	// Output:
	// NAME:
	//    core.test node - Commands can only be run from on the same machine as the Chainlink node.
	//
	// USAGE:
	//    core.test node command [command options] [arguments...]
	//
	// COMMANDS:
	//    deleteuser                Erase the *local node's* user and corresponding session to force recreation on next node launch.
	//    setnextnonce              Manually set the next nonce for a key. This should NEVER be necessary during normal operation. USE WITH CAUTION: Setting this incorrectly can break your node.
	//    start, node, n            Run the chainlink node
	//    rebroadcast-transactions  Manually rebroadcast txs matching nonce range with the specified gas price. This is useful in emergencies e.g. high gas prices and/or network congestion to forcibly clear out the pending TX queue
	//    status                    Displays the health of various services running inside the node.
	//    db                        Commands for managing the database.
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_node_start() {
	run("node", "start", "--help")
	// Output:
	// NAME:
	//    core.test node start - Run the chainlink node
	//
	// USAGE:
	//    core.test node start [command options] [arguments...]
	//
	// OPTIONS:
	//    --api value, -a value            text file holding the API email and password, each on a line
	//    --debug, -d                      set logger level to debug
	//    --password value, -p value       text file holding the password for the node's account
	//    --vrfpassword value, --vp value  textfile holding the password for the vrf keys; enables chainlink VRF oracle
}

func ExampleRun_node_db() {
	run("node", "db", "--help")
	// Output:
	// NAME:
	//    core.test node db - Potentially destructive commands for managing the database.
	//
	// USAGE:
	//    core.test node db command [command options] [arguments...]
	//
	// COMMANDS:
	//    version   Display the current database version.
	//    status    Display the current database migration status.
	//    migrate   Migrate the database to the latest version.
	//    rollback  Roll back the database to a previous <version>. Rolls back a single migration if no version specified.
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_txs() {
	run("txs", "--help")
	// Output:
	// NAME:
	//    core.test txs - Commands for handling Ethereum transactions
	//
	// USAGE:
	//    core.test txs command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Send <amount> Eth from node ETH account <fromAddress> to destination <toAddress>.
	//    list    List the Ethereum Transactions in descending order
	//    show    get information on a specific Ethereum Transaction
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_chains() {
	run("chains", "--help")
	// Output:
	// NAME:
	//    core.test chains - Commands for handling chain configuration
	//
	// USAGE:
	//    core.test chains command [command options] [arguments...]
	//
	// COMMANDS:
	//    evm  Commands for handling EVM chains
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_chains_evm() {
	run("chains", "evm", "--help")
	// Output:
	// NAME:
	//    core.test chains evm - Commands for handling EVM chains
	//
	// USAGE:
	//    core.test chains evm command [command options] [arguments...]
	//
	// COMMANDS:
	//    create     Create a new EVM chain
	//    delete     Delete an EVM chain
	//    list       List all chains
	//    configure  Configure an EVM chain
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_nodes() {
	run("nodes", "--help")
	// Output:
	// NAME:
	//    core.test nodes - Commands for handling node configuration
	//
	// USAGE:
	//    core.test nodes command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a new node
	//    delete  Delete a node
	//    list    List all nodes
	//
	// OPTIONS:
	//    --help, -h  show help
}
