//go:build wasip1 && maxbinary

package main

import (
	_ "embed"
	"fmt"
)

// ballast.dat must be generated before building with -tags maxbinary:
//
//	dd if=/dev/urandom of=ballast.dat bs=1024 count=14336
//
// This produces a ~14 MB file that inflates the compressed WASM binary
// to approach the WASMCompressedBinarySizeLimit (20 MB).
//
//go:embed ballast.dat
var ballastData []byte

func init() {
	// Reference ballastData so the linker retains it.
	_ = fmt.Sprintf("ballast: %d bytes", len(ballastData))
}
