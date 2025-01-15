package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generation/generate/genwrapper"
)

const (
	rootDir = "../../../contracts/solc/"
)

func main() {
	project := os.Args[1]
	className := os.Args[2]
	pkgName := os.Args[3]

	abiPath := rootDir + project + "/" + className + "/" + className + ".sol/" + className + ".abi.json"
	binPath := rootDir + project + "/" + className + "/" + className + ".sol/" + className + ".bin"

	genwrapper.GenWrapper(abiPath, binPath, className, pkgName)
}
