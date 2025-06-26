//go:build wasip1

package main

import (
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2"
)

func RunEmptyWorkflow(_ *sdk.Environment[struct{}]) (sdk.Workflow[struct{}], error) {
	return sdk.Workflow[struct{}]{}, nil
}

func main() {
	wasm.NewRunner(func(_ []byte) (struct{}, error) { return struct{}{}, nil }).Run(RunEmptyWorkflow)
}
