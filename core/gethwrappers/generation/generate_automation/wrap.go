package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generation/generate/genwrapper"
)

// Multiple legacy automation classes named X_{digits}, while being in X_{digits} folder,
// drop the {digits} in the output wrapper go class. Once such classes are removed, we can drop the current wrap.go
// and switch to core/gethwrappers/generation/wrap.go.
func main() {
	rootDir := "../../contracts/solc/"
	project := "automation"
	inputClassName := ""
	outClassName := ""
	pkgName := ""
	switch len(os.Args) {
	case 3:
		inputClassName = os.Args[1]
		outClassName = inputClassName
		pkgName = os.Args[2]
	case 4:
		inputClassName = os.Args[1]
		outClassName = os.Args[2]
		pkgName = os.Args[3]
	default:
		panic("Unsupported number of args")
	}

	abiPath := rootDir + project + "/" + inputClassName + "/" + inputClassName + ".sol/" + inputClassName + ".abi.json"
	binPath := rootDir + project + "/" + inputClassName + "/" + inputClassName + ".sol/" + inputClassName + ".bin"

	genwrapper.GenWrapper(abiPath, binPath, outClassName, pkgName, "")
}
