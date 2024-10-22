package src

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type deleteJobs struct{}

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
	return &deleteJobs{}
}

func (g *deleteJobs) Name() string {
	return "delete-ocr3-jobs"
}

func (g *deleteJobs) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	keylessNodeSetsPath := fs.String("nodes", "", "Custom keyless node sets location")
	artefactsDir := fs.String("artefacts", "", "Custom artefacts directory location")
	nodeSetSize := fs.Int("nodeSetSize", 5, "number of nodes in a nodeset")

	err := fs.Parse(args)
	if err != nil {
		fs.Usage()
		os.Exit(1)
	}

	if *artefactsDir == "" {
		*artefactsDir = defaultArtefactsDir
	}
	if *keylessNodeSetsPath == "" {
		*keylessNodeSetsPath = defaultKeylessNodeSetsPath
	}

	deployedContracts, err := LoadDeployedContracts(*artefactsDir)
	if err != nil {
		fmt.Println("Error loading deployed contracts, skipping:", err)
		return
	}
	nodes := downloadKeylessNodeSets(*keylessNodeSetsPath, *nodeSetSize).Workflow.Nodes

	for _, node := range nodes {
		api := newNodeAPI(node)
		jobsb := api.mustExec(api.methods.ListJobs)

		var parsed []JobSpec
		err = json.Unmarshal(jobsb, &parsed)
		helpers.PanicErr(err)

		for _, jobSpec := range parsed {
			if jobSpec.BootstrapSpec.ContractID == deployedContracts.OCRContract.String() ||
				jobSpec.OffChainReporting2OracleSpec.ContractID == deployedContracts.OCRContract.String() {
				fmt.Println("Deleting OCR3 job ID:", jobSpec.Id, "name:", jobSpec.Name)
				api.withArg(jobSpec.Id).mustExec(api.methods.DeleteJob)
			}
		}
	}
}
