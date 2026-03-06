//go:build wasip1

package main

import (
	_ "embed"
	"fmt"
)

// ballast.dat must exist before building. Use `make build` for a minimal
// placeholder or `make build-max` for a ~14 MB file that pushes the
// compressed WASM binary toward the 20 MB limit.
//
//go:embed ballast.dat
var ballastData []byte

func init() {
	// Reference ballastData so the linker retains it.
	_ = fmt.Sprintf("ballast: %d bytes", len(ballastData))
}
