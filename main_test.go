package main_test

import (
	"io/ioutil"
	"os"

	"github.com/smartcontractkit/chainlink"
	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
)

func ExampleMain_Help() {
	tc, cleanup := cltest.NewConfig()
	defer cleanup()
	testClient := &cmd.Client{
		cmd.RendererTable{ioutil.Discard},
		tc.Config,
		cmd.ChainlinkAppFactory{},
		cmd.TerminalAuthenticator{&cltest.MockCountingPrompt{}, os.Exit},
		cmd.NodeRunner{},
	}

	main.Run(testClient, "chainlink.test --help")
	// Output:
	// NAME:
	//    chainlink.test - CLI for Chainlink
	//
	// USAGE:
	//    chainlink.test [global options] command [command options] [arguments...]
	//
	// VERSION:
	//    0.2.0
	//
	// COMMANDS:
	//      node, n  Run the chainlink node
	//      jobs, j  Get all jobs
	//      show, s  Show a specific job
	//      help, h  Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --json, -j     json output as opposed to table
	//    --help, -h     show help
	//    --version, -v  print the version
}
