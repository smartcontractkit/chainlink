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

type deleteJobs struct {
	NodeList  string
	Artefacts string
}

type OCRSpec struct {
	ContractID string
}

type BootSpec struct {
	ContractID string
}

type WorkflowSpec struct {
	WorkflowID string
}

type JobSpec struct {
	Id                           string
	Name                         string
	BootstrapSpec                BootSpec
	OffChainReporting2OracleSpec OCRSpec
	WorkflowSpec                 WorkflowSpec
}

func NewDeleteJobsCommand() *deleteJobs {
	return &deleteJobs{
		NodeList:  ".cache/NodeList.txt",
		Artefacts: artefactsDir,
	}
}

func (g *deleteJobs) Name() string {
	return "delete-ocr3-jobs"
}

func (g *deleteJobs) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	customNodeList := fs.String("nodes", "", "Custom node list location")
	customArtefacts := fs.String("artefacts", "", "Custom artefacts directory location")

	err := fs.Parse(args)
	if err != nil {
		fs.Usage()
		os.Exit(1)
	}

	if *customArtefacts != "" {
		fmt.Printf("Custom  artefacts folder flag detected, using custom path %s", *customArtefacts)
		g.Artefacts = *customArtefacts
	}

	if *customNodeList != "" {
		fmt.Printf("Custom node file override flag detected, using custom node file path %s", *customNodeList)
		g.NodeList = *customNodeList
	}

	deployedContracts, err := LoadDeployedContracts(g.Artefacts)
	helpers.PanicErr(err)
	nodes := downloadNodeAPICredentialsDefault(g.NodeList)

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
			if jobSpec.BootstrapSpec.ContractID == deployedContracts.OCRContract.String() ||
				jobSpec.OffChainReporting2OracleSpec.ContractID == deployedContracts.OCRContract.String() {
				fmt.Println("Deleting OCR3 job ID:", jobSpec.Id, "name:", jobSpec.Name)
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
