package src

import (
	"errors"
	"flag"
	"fmt"
	"os"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)
// Could be useful https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/scripts/chaincli/handler/scrape_node_config.go#L102
type deployJobSpecs struct{}

func NewDeployJobSpecsCommand() *deployJobSpecs {
	return &deployJobSpecs{}
}

func (g *deployJobSpecs) Name() string {
	return "deploy-jobspecs"
}

func (g *deployJobSpecs) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	chainID := fs.Int64("chainid", 1337, "chain id")
	p2pPort := fs.Int64("p2pport", 6690, "p2p port")
	onlyReplay := fs.Bool("onlyreplay", false, "only replay the block from the OCR3 contract setConfig transaction")
	templatesLocation := fs.String("templates", "", "Custom templates location")
	nodeList := fs.String("nodes", "", "Custom node list location")
	publicKeys := fs.String("publickeys", "", "Custom public keys json location")
	artefactsDir := fs.String("artefacts", "", "Custom artefacts directory location")

	err := fs.Parse(args)
	if err != nil || chainID == nil || *chainID == 0 || p2pPort == nil || *p2pPort == 0 || onlyReplay == nil {
		fs.Usage()
		os.Exit(1)
	}
	if *onlyReplay {
		fmt.Println("Only replaying OCR3 contract setConfig transaction")
	} else {
		fmt.Println("Deploying OCR3 job specs")
	}

	if *artefactsDir == "" {
		*artefactsDir = defaultArtefactsDir
	}
	if *publicKeys == "" {
		*publicKeys = defaultPublicKeys
	}
	if *nodeList == "" {
		*nodeList = defaultNodeList
	}
	if *templatesLocation == "" {
		*templatesLocation = "templates"
	}

	nodes := downloadNodeAPICredentials(*nodeList)
	deployedContracts, err := LoadDeployedContracts(*artefactsDir)
	PanicErr(err)

	jobspecs := genSpecs(
		*publicKeys,
		*nodeList,
		*templatesLocation,
		*chainID, *p2pPort, deployedContracts.OCRContract.Hex(),
	)
	flattenedSpecs := []hostSpec{jobspecs.bootstrap}
	flattenedSpecs = append(flattenedSpecs, jobspecs.oracles...)

	// sanity check arr lengths
	if len(nodes) != len(flattenedSpecs) {
		PanicErr(errors.New("Mismatched node and job spec lengths"))
	}

	for i, n := range nodes {
		api := newNodeAPI(n)
		if !*onlyReplay {
			specToDeploy := flattenedSpecs[i].spec.ToString()
			specFragment := flattenedSpecs[i].spec[0:1]
			fmt.Printf("Deploying jobspec: %s\n... \n", specFragment)

			_, err := api.withArg(specToDeploy).exec(api.methods.CreateJob)
			if err != nil {
				fmt.Println("Failed to deploy job spec:", specFragment, "Error:", err)
			}
		}

		fmt.Printf("Replaying from block: %d\n", deployedContracts.SetConfigTxBlock)
		fmt.Printf("EVM Chain ID: %d\n\n", *chainID)
		api.withFlags(api.methods.ReplayFromBlock, func(fs *flag.FlagSet) {
			err = fs.Set("block-number", fmt.Sprint(deployedContracts.SetConfigTxBlock))
			helpers.PanicErr(err)
			err = fs.Set("evm-chain-id", fmt.Sprint(*chainID))
			helpers.PanicErr(err)
		}).mustExec()
	}
}
