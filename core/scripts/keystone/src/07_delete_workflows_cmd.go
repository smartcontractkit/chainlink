package src

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type deleteWorkflows struct{}

func NewDeleteWorkflowsCommand() *deleteWorkflows {
	return &deleteWorkflows{}
}

func (g *deleteWorkflows) Name() string {
	return "delete-workflows"
}

func (g *deleteWorkflows) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ExitOnError)
	nodeList := fs.String("nodes", "", "Custom node list location")

	err := fs.Parse(args)
	if err != nil {
		fs.Usage()
		os.Exit(1)
	}

	if *nodeList == "" {
		*nodeList = defaultNodeList
	}

	nodes := downloadNodeAPICredentials(*nodeList)

	for _, node := range nodes {
		output := &bytes.Buffer{}
		client, app := newApp(node, output)

		fmt.Println("Logging in:", node.url)
		loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
		loginFs.Bool("bypass-version-check", true, "")
		loginCtx := cli.NewContext(app, loginFs, nil)
		err := client.RemoteLogin(loginCtx)
		helpers.PanicErr(err)
		output.Reset()

		fileFs := flag.NewFlagSet("test", flag.ExitOnError)
		err = client.ListJobs(cli.NewContext(app, fileFs, nil))
		helpers.PanicErr(err)

		var parsed []JobSpec
		err = json.Unmarshal(output.Bytes(), &parsed)
		helpers.PanicErr(err)

		for _, jobSpec := range parsed {
			if jobSpec.WorkflowSpec.WorkflowID != "" {
				fmt.Println("Deleting workflow job ID:", jobSpec.Id, "name:", jobSpec.Name)
				set := flag.NewFlagSet("test", flag.ExitOnError)
				err = set.Parse([]string{jobSpec.Id})
				helpers.PanicErr(err)
				err = client.DeleteJob(cli.NewContext(app, set, nil))
				helpers.PanicErr(err)
			}
		}

		output.Reset()
	}
}
