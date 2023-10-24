package gethwrappers

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Test running an abigen command to generate a contract wrapper.
// Using as test example:
// go run ./generation/generate/wrap.go ../../contracts/solc/v0.6/VRF.abi ../../contracts/solc/v0.6/VRF.bin VRF solidity_vrf_wrapper
func TestAbigen(t *testing.T) {
	abiPath := "../../contracts/solc/v0.6/VRF.abi"
	binPath := "../../contracts/solc/v0.6/VRF.bin"
	className := "VRF"
	pkgName := "solidity_vrf_wrapper"
	fmt.Println("Generating", pkgName, "contract wrapper")

	cwd, err := os.Getwd() // gethwrappers directory
	if err != nil {
		Exit("could not get working directory", err)
	}
	outDir := filepath.Join(cwd, "generated", pkgName)
	if mkdErr := os.MkdirAll(outDir, 0700); err != nil {
		Exit("failed to create wrapper dir", mkdErr)
	}
	outPath := filepath.Join(outDir, pkgName+".go")

	Abigen(AbigenArgs{
		Bin: binPath, ABI: abiPath, Out: outPath, Type: className, Pkg: pkgName,
	})
}
