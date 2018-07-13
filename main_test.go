// +build !windows

package main

import (
	"io/ioutil"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
)

func ExampleRun() {
	tc, cleanup := cltest.NewConfig()
	defer cleanup()
	testClient := &cmd.Client{
		Renderer:        cmd.RendererTable{Writer: ioutil.Discard},
		Config:          tc.Config,
		AppFactory:      cmd.ChainlinkAppFactory{},
		Auth:            cmd.TerminalAuthenticator{Prompter: &cltest.MockCountingPrompt{}},
		UserInitializer: cltest.MockUserInitializer{},
		Runner:          cmd.ChainlinkRunner{},
	}

	Run(testClient, "chainlink.test", "--help")
	// Output:
	// NAME:
	//    chainlink.test - CLI for Chainlink
	//
	// USAGE:
	//    chainlink.test [global options] command [command options] [arguments...]
	//
	// VERSION:
	//    unset@unset
	//
	// COMMANDS:
	//      node, n                   Run the chainlink node
	//      account, a                Display the account address with its ETH & LINK balances
	//      jobspecs, jobs, j, specs  Get all jobs
	//      show, s                   Show a specific job
	//      create, c                 Create job spec from JSON
	//      run, r                    Begin job run for specid
	//      backup                    Backup the database of the running node
	//      import, i                 Import a key file to use with the node
	//      bridge                    Add a new bridge to the node
	//      getbridges                List all bridges added to the node
	//      showbridge                Show a specific bridge
	//      removebridge              Removes a specific bridge
	//      help, h                   Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --json, -j     json output as opposed to table
	//    --help, -h     show help
	//    --version, -v  print the version
}

func ExampleVersion() {
	tc, cleanup := cltest.NewConfig()
	defer cleanup()
	testClient := &cmd.Client{
		Renderer:        cmd.RendererTable{Writer: ioutil.Discard},
		Config:          tc.Config,
		AppFactory:      cmd.ChainlinkAppFactory{},
		Auth:            cmd.TerminalAuthenticator{Prompter: &cltest.MockCountingPrompt{}},
		UserInitializer: cltest.MockUserInitializer{},
		Runner:          cmd.ChainlinkRunner{},
	}

	Run(testClient, "chainlink.test", "--version")
	// Output:
	// chainlink.test version unset@unset
}
