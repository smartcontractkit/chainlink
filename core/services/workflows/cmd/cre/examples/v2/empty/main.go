//go:build wasip1

package main

import (
	"github.com/smartcontractkit/cre-sdk-go/sdk"
	"github.com/smartcontractkit/cre-sdk-go/sdk/wasm"
)

func RunEmptyWorkflow(_ *sdk.Environment[struct{}]) (sdk.Workflow[struct{}], error) {
	return sdk.Workflow[struct{}]{}, nil
}

func main() {
	wasm.NewRunner(func(_ []byte) (struct{}, error) { return struct{}{}, nil }).Run(RunEmptyWorkflow)
}
