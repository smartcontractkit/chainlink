package main

import (
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/cmd/cre/examples/v2/e2e/pkg"
)

func main() {
	pkg.InitWorkflow(wasm.NewDonRunner())
}
