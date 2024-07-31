package src

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type fetchKeys struct {
	outputFile string
}

func NewFetchKeysCommand() *fetchKeys {
	return &fetchKeys{
		outputFile: "PublicKeys.json",
	}
}

func (g *fetchKeys) Name() string {
	return "fetch-keys"
}

func (g *fetchKeys) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	nodesFile := fs.String("nodes", "", "a file containing nodes urls, logins and passwords")
	chainID := fs.Int64("chainid", 80001, "chain id")
	outputFile := fs.String("outputfile", "", "Custom output file")

	if err := fs.Parse(args); err != nil || *nodesFile == "" || *chainID == 0 {
		fs.Usage()
		os.Exit(1)
	}

	if *outputFile != "" {
		fmt.Printf("Custom output file override flag detected, using custom path %s\n", *outputFile)
		g.outputFile = *outputFile
	}

	nodes := mustReadNodesList(*nodesFile)
	nca := mustFetchNodesKeys(*chainID, nodes)

	nodePublicKeys, err := json.MarshalIndent(nca, "", " ")
	if err != nil {
		panic(err)
	}
	filepath := g.outputFile
	err = os.WriteFile(filepath, nodePublicKeys, 0600)
	if err != nil {
		panic(err)
	}
	fmt.Println("Functions OCR2 public keys have been saved to:", filepath)
}
