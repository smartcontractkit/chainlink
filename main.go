package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/v2/core"
)

//go:generate make modgraph
func main() {
	os.Exit(core.Main())
}
