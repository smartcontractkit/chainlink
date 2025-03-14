package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generation/generate/genwrapper"
)

// Multiple legacy automation classes named X_{digits}, while being in X_{digits} folder,
// drop the {digits} in the output wrapper go class. Once such classes are removed/renamed,
// we can drop the current wrap.go and switch to core/gethwrappers/generation/wrap.go.
func main() {
	rootDir := "../../contracts/solc/"
	project := "automation"
	inputClassName := os.Args[1]
	outputClassName := os.Args[2]
	pkgName := os.Args[3]

	abiPath := rootDir + project + "/" + inputClassName + "/" + inputClassName + ".sol/" + inputClassName + ".abi.json"
	binPath := rootDir + project + "/" + inputClassName + "/" + inputClassName + ".sol/" + inputClassName + ".bin"

	genwrapper.GenWrapper(abiPath, binPath, outputClassName, pkgName, "")
}
