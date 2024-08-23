package src

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type fetchKeys struct {
}

func NewFetchKeysCommand() *fetchKeys {
	return &fetchKeys{}
}

func (g *fetchKeys) Name() string {
	return "fetch-keys"
}

func (g *fetchKeys) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	nodesFile := fs.String("nodes", "", "a file containing nodes urls, logins and passwords")
	chainID := fs.Int64("chainid", 80001, "chain id")
	if err := fs.Parse(args); err != nil || *nodesFile == "" || *chainID == 0 {
		fs.Usage()
		os.Exit(1)
	}

	nodes := mustReadNodesList(*nodesFile)
	nca := mustFetchNodesKeys(*chainID, nodes)

	nodePublicKeys, err := json.MarshalIndent(nca, "", " ")
	if err != nil {
		panic(err)
	}
	filepath := "PublicKeys.json"
	err = os.WriteFile(filepath, nodePublicKeys, 0600)
	if err != nil {
		panic(err)
	}
	fmt.Println("Functions OCR2 public keys have been saved to:", filepath)
}
