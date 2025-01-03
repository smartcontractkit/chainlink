package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generation/generate/genwrapper"
)

func main() {
	abiPath := os.Args[1]
	binPath := os.Args[2]
	className := os.Args[3]
	pkgName := os.Args[4]

	genwrapper.GenWrapper(abiPath, binPath, className, pkgName)
}
