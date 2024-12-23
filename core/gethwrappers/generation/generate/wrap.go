package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generation/generate/gen_wrapper"
)

func main() {
	abiPath := os.Args[1]
	binPath := os.Args[2]
	className := os.Args[3]
	pkgName := os.Args[4]

	gen_wrapper.GenWrapper(abiPath, binPath, className, pkgName)
}
