//go:build !windows
// +build !windows

package main

import (
	"io"
	"testing"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/static"
)

func init() {
	static.Version = "0.0.0"
	static.Sha = "exampleSHA"
}

func Run(args ...string) {
	t := &testing.T{}

	tc := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.DevMode = false
		foo := "foo"
		c.RootDir = &foo
	})
	lggr := logger.TestLogger(t)
	testClient := &cmd.Client{
		Renderer:               cmd.RendererTable{Writer: io.Discard},
		Config:                 tc,
		Logger:                 lggr,
		AppFactory:             cmd.ChainlinkAppFactory{},
		FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
		Runner:                 cmd.ChainlinkRunner{},
		HTTP:                   cltest.NewMockAuthenticatedHTTPClient(lggr, cmd.ClientOpts{}, "session"),
		ChangePasswordPrompter: cltest.MockChangePasswordPrompter{},
	}
	args = append([]string{""}, args...)
	run(testClient, args...)
}

func ExampleRun() {
	Run("--help")
	Run("--version")
	// Output:
	// NAME:
	//    core.test - CLI for Chainlink
	//
	// USAGE:
	//    core.test [global options] command [command options] [arguments...]
	//
	// VERSION:
	//    0.0.0@exampleSHA
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
	//    txs             Commands for handling transactions
	//    chains          Commands for handling chain configuration
	//    nodes           Commands for handling node configuration
	//    forwarders      Commands for managing forwarder addresses.
	//    help, h         Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --json, -j                     json output as opposed to table
	//    --admin-credentials-file FILE  optional, applies only in client mode when making remote API calls. If provided, FILE containing admin credentials will be used for logging in, allowing to avoid an additional login step. If `FILE` is missing, it will be ignored. Defaults to <RootDir>/apicredentials
	//    --remote-node-url URL          optional, applies only in client mode when making remote API calls. If provided, URL will be used as the remote Chainlink API endpoint (default: "http://localhost:6688")
	//    --insecure-skip-verify         optional, applies only in client mode when making remote API calls. If turned on, SSL certificate verification will be disabled. This is mostly useful for people who want to use Chainlink with a self-signed TLS certificate
	//    --config value, -c value       TOML configuration file(s) via flag, or raw TOML via env var. If used, legacy env vars must not be set. Multiple files can be used (-c configA.toml -c configB.toml), and they are applied in order with duplicated fields overriding any earlier values. If the 'CL_CONFIG' env var is specified, it is always processed last with the effect of being the final override. [$CL_CONFIG]
	//    --secrets value, -s value      TOML configuration file for secrets. Must be set if and only if config is set.
	//    --help, -h                     show help
	//    --version, -v                  print the version
	// core.test version 0.0.0@exampleSHA
}

func ExampleRun_admin() {
	Run("admin", "--help")
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
	//    logout  Delete any local sessions
	//    users   Create, edit permissions, or delete API users
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_attempts() {
	Run("attempts", "--help")
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
	Run("blocks", "--help")
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
	Run("bridges", "--help")
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
	//    show     Show a Bridge's details
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_config() {
	Run("config", "--help")
	// Output:
	// NAME:
	//    core.test config - Commands for the node's configuration
	//
	// USAGE:
	//    core.test config command [command options] [arguments...]
	//
	// COMMANDS:
	//    dump         Dump prints V2 TOML that is equivalent to the current environment and database configuration [Not supported with TOML]
	//    list         Show the node's environment variables [Not supported with TOML]
	//    show         Show the application configuration [Only supported with TOML]
	//    setgasprice  Set the default gas price to use for outgoing transactions [Not supported with TOML]
	//    loglevel     Set log level
	//    logsql       Enable/disable sql statement logging
	//    validate     Validate provided TOML config file, and print the full effective configuration, with defaults included [Only supported with TOML]
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_jobs() {
	Run("jobs", "--help")
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
	Run("keys", "--help")
	// Output:
	// NAME:
	//    core.test keys - Commands for managing various types of keys used by the Chainlink node
	//
	// USAGE:
	//    core.test keys command [command options] [arguments...]
	//
	// COMMANDS:
	//    eth         Remote commands for administering the node's Ethereum keys
	//    p2p         Remote commands for administering the node's p2p keys
	//    csa         Remote commands for administering the node's CSA keys
	//    ocr         Remote commands for administering the node's legacy off chain reporting keys
	//    ocr2        Remote commands for administering the node's off chain reporting keys
	//    solana      Remote commands for administering the node's Solana keys
	//    starknet    Remote commands for administering the node's StarkNet keys
	//    dkgsign     Remote commands for administering the node's DKGSign keys
	//    dkgencrypt  Remote commands for administering the node's DKGEncrypt keys
	//    vrf         Remote commands for administering the node's vrf keys
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys_eth() {
	Run("keys", "eth", "--help")
	// Output:
	// NAME:
	//    core.test keys eth - Remote commands for administering the node's Ethereum keys
	//
	// USAGE:
	//    core.test keys eth command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a key in the node's keystore alongside the existing key; to create an original key, just run the node
	//    update  Update the existing key's parameters
	//    list    List available Ethereum accounts with their ETH & LINK balances, nonces, and other metadata
	//    delete  Delete the ETH key by address (irreversible!)
	//    import  Import an ETH key from a JSON file
	//    export  Exports an ETH key to a JSON file
	//    chain   Update an EVM key for the given chain
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys_p2p() {
	Run("keys", "p2p", "--help")
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
	Run("keys", "csa", "--help")
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

func ExampleRun_keys_ocr_legacy() {
	Run("keys", "ocr", "--help")
	// Output:
	// NAME:
	//    core.test keys ocr - Remote commands for administering the node's legacy off chain reporting keys
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

func ExampleRun_keys_ocr() {
	Run("keys", "ocr2", "--help")
	// Output:
	// NAME:
	//    core.test keys ocr2 - Remote commands for administering the node's off chain reporting keys
	//
	// USAGE:
	//    core.test keys ocr2 command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create an OCR2 key bundle, encrypted with password from the password file, and store it in the database
	//    delete  Deletes the encrypted OCR2 key bundle matching the given ID
	//    list    List available OCR2 key bundles
	//    import  Imports an OCR2 key bundle from a JSON file
	//    export  Exports an OCR2 key bundle to a JSON file
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys_solana() {
	Run("keys", "solana", "--help")
	// Output:
	// NAME:
	//    core.test keys solana - Remote commands for administering the node's Solana keys
	//
	// USAGE:
	//    core.test keys solana command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a Solana key
	//    import  Import Solana key from keyfile
	//    export  Export Solana key to keyfile
	//    delete  Delete Solana key if present
	//    list    List the Solana keys
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys_starknet() {
	Run("keys", "starknet", "--help")
	// Output:
	// NAME:
	//    core.test keys starknet - Remote commands for administering the node's StarkNet keys
	//
	// USAGE:
	//    core.test keys starknet command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a StarkNet key
	//    import  Import StarkNet key from keyfile
	//    export  Export StarkNet key to keyfile
	//    delete  Delete StarkNet key if present
	//    list    List the StarkNet keys
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys_vrf() {
	Run("keys", "vrf", "--help")
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

func ExampleRun_keys_dkgsign() {
	Run("keys", "dkgsign", "--help")
	// Output:
	// NAME:
	//    core.test keys dkgsign - Remote commands for administering the node's DKGSign keys
	//
	// USAGE:
	//    core.test keys dkgsign command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a DKGSign key
	//    import  Import DKGSign key from keyfile
	//    export  Export DKGSign key to keyfile
	//    delete  Delete DKGSign key if present
	//    list    List the DKGSign keys
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_keys_dkgencrypt() {
	Run("keys", "dkgencrypt", "--help")
	// Output:
	// NAME:
	//    core.test keys dkgencrypt - Remote commands for administering the node's DKGEncrypt keys
	//
	// USAGE:
	//    core.test keys dkgencrypt command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a DKGEncrypt key
	//    import  Import DKGEncrypt key from keyfile
	//    export  Export DKGEncrypt key to keyfile
	//    delete  Delete DKGEncrypt key if present
	//    list    List the DKGEncrypt keys
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_node() {
	Run("node", "--help")
	// Output:
	// NAME:
	//    core.test node - Commands can only be run from on the same machine as the Chainlink node.
	//
	// USAGE:
	//    core.test node command [command options] [arguments...]
	//
	// COMMANDS:
	//    start, node, n            Run the Chainlink node
	//    rebroadcast-transactions  Manually rebroadcast txs matching nonce range with the specified gas price. This is useful in emergencies e.g. high gas prices and/or network congestion to forcibly clear out the pending TX queue
	//    status                    Displays the health of various services running inside the node.
	//    profile                   Collects profile metrics from the node.
	//    db                        Commands for managing the database.
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_node_start() {
	Run("node", "start", "--help")
	// Output:
	// NAME:
	//    core.test node start - Run the Chainlink node
	//
	// USAGE:
	//    core.test node start [command options] [arguments...]
	//
	// OPTIONS:
	//    --api value, -a value            text file holding the API email and password, each on a line
	//    --debug, -d                      set logger level to debug
	//    --password value, -p value       text file holding the password for the node's account
	//    --vrfpassword value, --vp value  text file holding the password for the vrf keys; enables Chainlink VRF oracle
}

func ExampleRun_node_db() {
	Run("node", "db", "--help")
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

func ExampleRun_node_profile() {
	Run("node", "profile", "--help")
	// Output:
	// NAME:
	//    core.test node profile - Collects profile metrics from the node.
	//
	// USAGE:
	//    core.test node profile [command options] [arguments...]
	//
	// OPTIONS:
	//    --seconds value, -s value     duration of profile capture (default: 8)
	//    --output_dir value, -o value  output directory of the captured profile (default: "/tmp/")
}

func ExampleRun_txs() {
	Run("txs", "--help")
	// Output:
	// NAME:
	//    core.test txs - Commands for handling transactions
	//
	// USAGE:
	//    core.test txs command [command options] [arguments...]
	//
	// COMMANDS:
	//    evm     Commands for handling EVM transactions
	//    solana  Commands for handling Solana transactions
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_txs_evm() {
	Run("txs", "evm", "--help")
	// Output:
	// NAME:
	//    core.test txs evm - Commands for handling EVM transactions
	//
	// USAGE:
	//    core.test txs evm command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Send <amount> ETH (or wei) from node ETH account <fromAddress> to destination <toAddress>.
	//    list    List the Ethereum Transactions in descending order
	//    show    get information on a specific Ethereum Transaction
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_txs_solana() {
	Run("txs", "solana", "--help")
	// Output:
	// NAME:
	//    core.test txs solana - Commands for handling Solana transactions
	//
	// USAGE:
	//    core.test txs solana command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Send <amount> lamports from node Solana account <fromAddress> to destination <toAddress>.
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_chains() {
	Run("chains", "--help")
	// Output:
	// NAME:
	//    core.test chains - Commands for handling chain configuration
	//
	// USAGE:
	//    core.test chains command [command options] [arguments...]
	//
	// COMMANDS:
	//    evm       Commands for handling EVM chains
	//    solana    Commands for handling Solana chains
	//    starknet  Commands for handling StarkNet chains
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_chains_evm() {
	Run("chains", "evm", "--help")
	// Output:
	// NAME:
	//    core.test chains evm - Commands for handling EVM chains
	//
	// USAGE:
	//    core.test chains evm command [command options] [arguments...]
	//
	// COMMANDS:
	//    create     Create a new EVM chain
	//    delete     Delete an existing EVM chain
	//    list       List all existing EVM chains
	//    configure  Configure an existing EVM chain
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_chains_solana() {
	Run("chains", "solana", "--help")
	// Output:
	// NAME:
	//    core.test chains solana - Commands for handling Solana chains
	//
	// USAGE:
	//    core.test chains solana command [command options] [arguments...]
	//
	// COMMANDS:
	//    create     Create a new Solana chain
	//    delete     Delete an existing Solana chain
	//    list       List all existing Solana chains
	//    configure  Configure an existing Solana chain
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_chains_starknet() {
	Run("chains", "starknet", "--help")
	// Output:
	// NAME:
	//    core.test chains starknet - Commands for handling StarkNet chains
	//
	// USAGE:
	//    core.test chains starknet command [command options] [arguments...]
	//
	// COMMANDS:
	//    create     Create a new StarkNet chain
	//    delete     Delete an existing StarkNet chain
	//    list       List all existing StarkNet chains
	//    configure  Configure an existing StarkNet chain
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_nodes() {
	Run("nodes", "--help")
	// Output:
	// NAME:
	//    core.test nodes - Commands for handling node configuration
	//
	// USAGE:
	//    core.test nodes command [command options] [arguments...]
	//
	// COMMANDS:
	//    evm       Commands for handling EVM node configuration
	//    solana    Commands for handling Solana node configuration
	//    starknet  Commands for handling StarkNet node configuration
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_nodes_evm() {
	Run("nodes", "evm", "--help")
	// Output:
	// NAME:
	//    core.test nodes evm - Commands for handling EVM node configuration
	//
	// USAGE:
	//    core.test nodes evm command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a new EVM node
	//    delete  Delete an existing EVM node
	//    list    List all existing EVM nodes
	//
	// OPTIONS:
	//    --help, -h  show help
}

// TODO https://app.shortcut.com/chainlinklabs/story/37975/chains-nodes-should-be-read-from-the-config-interface
//func ExampleRun_nodes_evm_list() {
//	runExTOML(`[[EVM]]
//	ChainID = 1
//	[[EVM.Node]]
//	Name = 'foo'
//	WSURL = 'wss://example.com
//	`, "nodes", "evm", "list")
//
//	// Output:
//	// [{...
//}

func ExampleRun_nodes_solana() {
	Run("nodes", "solana", "--help")
	// Output:
	// NAME:
	//    core.test nodes solana - Commands for handling Solana node configuration
	//
	// USAGE:
	//    core.test nodes solana command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a new Solana node
	//    delete  Delete an existing Solana node
	//    list    List all existing Solana nodes
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_nodes_starknet() {
	Run("nodes", "starknet", "--help")
	// Output:
	// NAME:
	//    core.test nodes starknet - Commands for handling StarkNet node configuration
	//
	// USAGE:
	//    core.test nodes starknet command [command options] [arguments...]
	//
	// COMMANDS:
	//    create  Create a new StarkNet node
	//    delete  Delete an existing StarkNet node
	//    list    List all existing StarkNet nodes
	//
	// OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_forwarders() {
	Run("forwarders", "--help")
	// Output:
	// NAME:
	//    core.test forwarders - Commands for managing forwarder addresses.
	//
	// USAGE:
	//    core.test forwarders command [command options] [arguments...]
	//
	// COMMANDS:
	//    list    List all stored forwarders addresses
	//    track   Track a new forwarder
	//    delete  Delete a forwarder address
	//
	// OPTIONS:
	//    --help, -h  show help
}
