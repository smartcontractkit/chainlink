//go:build wasip1

package main

import (
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

func RunEmptyWorkflow(_ *cre.Environment[struct{}]) (cre.Workflow[struct{}], error) {
	return cre.Workflow[struct{}]{}, nil
}

func main() {
	wasm.NewRunner(func(_ []byte) (struct{}, error) { return struct{}{}, nil }).Run(RunEmptyWorkflow)
}
