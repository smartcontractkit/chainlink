// package main is a script for generating a geth golang contract wrappers for
// the LINK token contract.
//
//  Usage:
//
// With core/gethwrappers as your working directory, run
//
//  go run generation/generate_link/wrap.go
//
// This will output the generated file to
// generated/link_token_interface/link_token_interface.go

package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tidwall/gjson"

	gethwrappers2 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func main() {
	pkgName := "link_token_interface"
	fmt.Println("Generating", pkgName, "contract wrapper")
	className := "LinkToken"
	tmpDir, cleanup := gethwrappers2.TempDir(className)
	defer cleanup()
	linkDetails, err := os.ReadFile(filepath.Join(
		gethwrappers2.GetProjectRoot(), "contracts/LinkToken.json"))
	if err != nil {
		gethwrappers2.Exit("could not read LINK contract details", err)
	}
	if fmt.Sprintf("%x", sha256.Sum256(linkDetails)) !=
		"27c0e17a79553fccc63a4400c6bbe415ff710d9cc7c25757bff0f7580205c922" {
		gethwrappers2.Exit("LINK details should never change!", nil)
	}
	abi, err := utils.NormalizedJSON([]byte(
		gjson.Get(string(linkDetails), "abi").String()))
	if err != nil || abi == "" {
		gethwrappers2.Exit("could not extract LINK ABI", err)
	}
	abiPath := filepath.Join(tmpDir, "abi")
	if aErr := os.WriteFile(abiPath, []byte(abi), 0600); aErr != nil {
		gethwrappers2.Exit("could not write contract ABI to temp dir.", aErr)
	}
	bin := gjson.Get(string(linkDetails), "bytecode").String()
	if bin == "" {
		gethwrappers2.Exit("could not extract LINK bytecode", nil)
	}
	binPath := filepath.Join(tmpDir, "bin")
	if bErr := os.WriteFile(binPath, []byte(bin), 0600); bErr != nil {
		gethwrappers2.Exit("could not write contract binary to temp dir.", bErr)
	}
	cwd, err := os.Getwd()
	if err != nil {
		gethwrappers2.Exit("could not get working directory", nil)
	}
	if filepath.Base(cwd) != "gethwrappers" {
		gethwrappers2.Exit("must be run from gethwrappers directory", nil)
	}
	outDir := filepath.Join(cwd, "generated", pkgName)
	if err := os.MkdirAll(outDir, 0700); err != nil {
		gethwrappers2.Exit("failed to create wrapper dir", err)
	}
	gethwrappers2.Abigen(gethwrappers2.AbigenArgs{
		Bin:  binPath,
		ABI:  abiPath,
		Out:  filepath.Join(outDir, pkgName+".go"),
		Type: className,
		Pkg:  pkgName,
	})
}
