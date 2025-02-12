package genwrapper

import (
	"fmt"
	"os"
	"path/filepath"

	gethParams "github.com/ethereum/go-ethereum/params"

	gethwrappers2 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers"
)

// GenWrapper generates a contract wrapper for the given contract.
//
// abiPath is the path to the contract's ABI JSON file.
//
// binPath is the path to the contract's binary file, typically with .bin extension.
//
// className is the name of the generated contract class.
//
// pkgName is the name of the package the contract will be generated in. Try
// to follow idiomatic Go package naming conventions where possible.
//
// outDirSuffixInput is the directory suffix to generate the wrapper in. If not provided, the
// wrapper will be generated in the default location. The default location is
// <project>/generated/<pkgName>/<pkgName>.go. The suffix will take place after
// the <project>/generated, so the overridden location would be
// <project>/generated/<outDirSuffixInput>/<pkgName>/<pkgName>.go.
func GenWrapper(abiPath, binPath, className, pkgName, outDirSuffixInput string) {
	fmt.Println("Generating", pkgName, "contract wrapper")

	cwd, err := os.Getwd() // gethwrappers directory
	if err != nil {
		gethwrappers2.Exit("could not get working directory", err)
	}
	outDir := filepath.Join(cwd, "generated", outDirSuffixInput, pkgName)
	if mkdErr := os.MkdirAll(outDir, 0700); err != nil {
		gethwrappers2.Exit(
			fmt.Sprintf("failed to create wrapper dir, outDirSuffixInput: %s (could be empty)", outDirSuffixInput),
			mkdErr)
	}
	outPath := filepath.Join(outDir, pkgName+".go")

	gethwrappers2.Abigen(gethwrappers2.AbigenArgs{
		Bin: binPath, ABI: abiPath, Out: outPath, Type: className, Pkg: pkgName,
	})

	// Build succeeded, so update the versions db with the new contract data
	versions, err := gethwrappers2.ReadVersionsDB()
	if err != nil {
		gethwrappers2.Exit("could not read current versions database", err)
	}
	versions.GethVersion = gethParams.Version
	versions.ContractVersions[pkgName] = gethwrappers2.ContractVersion{
		Hash:       gethwrappers2.VersionHash(abiPath, binPath),
		AbiPath:    abiPath,
		BinaryPath: binPath,
	}
	if err := gethwrappers2.WriteVersionsDB(versions); err != nil {
		gethwrappers2.Exit("could not save versions db", err)
	}
}
