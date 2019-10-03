// +build !windows

package main

import (
	"io/ioutil"
	"testing"

	"chainlink/core/cmd"
	"chainlink/core/internal/cltest"
)

func ExampleRun() {
	t := &testing.T{}
	tc, cleanup := cltest.NewConfig(t)
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
	// COMMANDS:
	//    admin        Commands for remotely taking admin related actions
	//    bridges      Commands for Bridges communicating with External Adapters
	//    config       Commands for the node's configuration
	//    jobs         Commands for managing Jobs
	//    node, local  Commands for admin actions that must be run locally
	//    runs         Commands for managing Runs
	//    txs          Commands for handling Ethereum transactions
	//    help, h      Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --json, -j     json output as opposed to table
	//    --help, -h     show help
	//    --version, -v  print the version
}

func ExampleVersion() {
	t := &testing.T{}
	tc, cleanup := cltest.NewConfig(t)
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
