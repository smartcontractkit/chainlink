package src

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type deleteJobs struct {
}

type OCRSpec struct {
	ContractID string
}

type BootSpec struct {
	ContractID string
}

type JobSpec struct {
	Id                           string
	Name                         string
	BootstrapSpec                BootSpec
	OffChainReporting2OracleSpec OCRSpec
}

func NewDeleteJobsCommand() *deleteJobs {
	return &deleteJobs{}
}

func (g *deleteJobs) Name() string {
	return "delete-ocr3-jobs"
}

func (g *deleteJobs) Run(args []string) {
	deployedContracts, err := LoadDeployedContracts()
	helpers.PanicErr(err)
	nodes := downloadNodeAPICredentialsDefault()

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
