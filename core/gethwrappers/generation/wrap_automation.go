package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generation/generate/genwrapper"
)

func main() {
	rootDir := "../../contracts/solc/"
	project := "automation"
	dirName := ""
	className := ""
	pkgName := ""
	if len(os.Args) == 3 {
		dirName = os.Args[1]
		className = dirName
		pkgName = os.Args[2]
	} else if len(os.Args) == 4 {
		dirName = os.Args[1]
		className = os.Args[2]
		pkgName = os.Args[3]
	} else {
		panic("Unsupported number of args")
	}

	abiPath := rootDir + project + "/" + dirName + "/" + dirName + ".sol/" + dirName + ".abi.json"
	binPath := rootDir + project + "/" + dirName + "/" + dirName + ".sol/" + dirName + ".bin"

	genwrapper.GenWrapper(abiPath, binPath, className, pkgName, "")
}
