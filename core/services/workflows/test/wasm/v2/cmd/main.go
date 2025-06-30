//go:build wasip1

package main

import (
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/testhelpers/v2"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2"
)

func main() {
	testhelpers.RunTestWorkflow(wasm.NewRunner(func(b []byte) (string, error) { return string(b), nil }))
}
