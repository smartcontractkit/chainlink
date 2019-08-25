// +build !windows

package main

import (
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
)

func ExampleRun() {
	t := &testing.T{}
	tc, cleanup := cltest.NewConfig(t)
	defer cltest.WipePostgresDatabase(t, tc.Config)
	defer cleanup()
	tc.Config.Set("CHAINLINK_DEV", false)
	testClient := &cmd.Client{
		Renderer:               cmd.RendererTable{Writer: ioutil.Discard},
		Config:                 tc.Config,
		AppFactory:             cmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:  cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{}},
		FallbackAPIInitializer: &cltest.MockAPIInitializer{},
		Runner:                 cmd.ChainlinkRunner{},
		HTTP:                   cltest.NewMockAuthenticatedHTTPClient(tc.Config),
		ChangePasswordPrompter: cltest.MockChangePasswordPrompter{},
	}

	Run(testClient, "core.test", "--help")
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
	//    local                      Commands which are run locally
	//    login                      Login to remote client by creating a session cookie
	//    account, a                 Display the account address with its ETH & LINK balances
	//    jobs, js, specs, jobspecs  List all jobs
	//    job, j                     Show a specific job's details
	//    create, c                  Create job spec from JSON
	//    archivejob                 Archive job and all associated runs
	//    run, r                     Create a new run for a job
	//    showrun, sr                Show a run for a specific ID
	//    runs, lr                   List all runs
	//    addbridge                  Create a new bridge to the node
	//    bridges                    List all bridges added to the node
	//    bridge                     Show a specific bridge
	//    removebridge               Removes a specific bridge
	//    initiators, exi            Tasks for managing external initiators
	//    agree, createsa            Creates a service agreement
	//    withdraw, w                Withdraw, to an authorized Ethereum <address>, <amount> units of LINK. Withdraws from the configured oracle contract by default, or from contract optionally specified by a third command-line argument --from-oracle-contract-address=<contract address>. Address inputs must be in EIP55-compliant capitalization.
	//    sendether                  Send <amount> ETH from the node's ETH account to an <address>.
	//    chpass                     Change your password
	//    setgasprice                Set the minimum gas price to use for outgoing transactions
	//    transactions, txs          List the transactions in descending order
	//    transaction, tx            get information on a specific transaction
	//    txattempts, txas           List the transaction attempts in descending order
	//    help, h                    Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --json, -j     json output as opposed to table
	//    --help, -h     show help
	//    --version, -v  print the version
}

func ExampleVersion() {
	t := &testing.T{}
	tc, cleanup := cltest.NewConfig(t)
	defer cltest.WipePostgresDatabase(t, tc.Config)
	defer cleanup()
	testClient := &cmd.Client{
		Renderer:               cmd.RendererTable{Writer: ioutil.Discard},
		Config:                 tc.Config,
		AppFactory:             cmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:  cmd.TerminalKeyStoreAuthenticator{Prompter: &cltest.MockCountingPrompter{}},
		FallbackAPIInitializer: &cltest.MockAPIInitializer{},
		Runner:                 cmd.ChainlinkRunner{},
		HTTP:                   cltest.NewMockAuthenticatedHTTPClient(tc.Config),
	}

	Run(testClient, "core.test", "--version")
	// Output:
	// core.test version unset@unset
}
