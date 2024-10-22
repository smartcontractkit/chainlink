package src

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type deployWorkflows struct{}

func NewDeployWorkflowsCommand() *deployWorkflows {
	return &deployWorkflows{}
}

func (g *deployWorkflows) Name() string {
	return "deploy-workflows"
}

func (g *deployWorkflows) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	workflowFile := fs.String("workflow", "workflow.yml", "path to workflow file")
	keylessNodeSetsPath := fs.String("nodes", "", "Custom keyless node sets location")
	nodeSetSize := fs.Int("nodeSetSize", 5, "number of nodes in a nodeset")
	err := fs.Parse(args)
	if err != nil || workflowFile == nil || *workflowFile == "" || nodeSetSize == nil || *nodeSetSize == 0 {
		fs.Usage()
		os.Exit(1)
	}
	if *keylessNodeSetsPath == "" {
		*keylessNodeSetsPath = defaultKeylessNodeSetsPath
	}
	fmt.Println("Deploying workflows")

	nodes := downloadKeylessNodeSets(*keylessNodeSetsPath, *nodeSetSize).Workflow.Nodes

	if _, err = os.Stat(*workflowFile); err != nil {
		PanicErr(errors.New("toml file does not exist"))
	}

	for i, n := range nodes {
		if i == 0 {
			continue // skip bootstrap node
		}
		output := &bytes.Buffer{}
		client, app := newApp(n, output)
		fmt.Println("Logging in:", n.RemoteURL)
		loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
		loginFs.Bool("bypass-version-check", true, "")
		loginCtx := cli.NewContext(app, loginFs, nil)
		err := client.RemoteLogin(loginCtx)
		helpers.PanicErr(err)
		output.Reset()

		fmt.Printf("Deploying workflow\n... \n")
		fs := flag.NewFlagSet("test", flag.ExitOnError)
		err = fs.Parse([]string{*workflowFile})

		helpers.PanicErr(err)
		err = client.CreateJob(cli.NewContext(app, fs, nil))
		if err != nil {
			fmt.Println("Failed to deploy workflow:", "Error:", err)
		}
		output.Reset()
	}
}
