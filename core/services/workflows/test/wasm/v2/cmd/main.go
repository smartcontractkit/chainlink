//go:build wasip1

package main

import (
	"github.com/smartcontractkit/cre-sdk-go/sdk/testutils"
	"github.com/smartcontractkit/cre-sdk-go/sdk/wasm"
)

func main() {
	testutils.RunTestWorkflow(wasm.NewRunner(func(b []byte) (string, error) { return string(b), nil }))
}
