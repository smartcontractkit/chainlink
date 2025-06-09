package main

import (
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/cmd/download"
	"github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/cmd/root"
)

func init() {
	root.RootCmd.AddCommand(download.DownloadCmd)
}

func main() {
	if err := root.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
