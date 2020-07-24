// package main is a script for generating a geth golang contract wrappers for
// the LINK token contract.
//
//  Usage:
//
// With core/internal/gethwrappers as your working directory, run
//
//  go run generation/generate_link/wrap.go
//
// This will output the generated file to
// generated/link_token_interface/link_token_interface.go

package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func main() {
	pkgName := "link_token_interface"
	fmt.Println("Generating", pkgName, "contract wrapper")
	className := "LinkToken"
	tmpDir, cleanup := gethwrappers.TempDir(className)
	defer cleanup()
	linkDetails, err := ioutil.ReadFile(filepath.Join(
		gethwrappers.GetProjectRoot(), "evm-test-helpers/src/LinkToken.json"))
	if err != nil {
		gethwrappers.Exit("could not read LINK contract details", err)
	}
	if fmt.Sprintf("%x", sha256.Sum256(linkDetails)) !=
		"27c0e17a79553fccc63a4400c6bbe415ff710d9cc7c25757bff0f7580205c922" {
		gethwrappers.Exit("LINK details should never change!", nil)
	}
	abi, err := utils.NormalizedJSON([]byte(
		gjson.Get(string(linkDetails), "abi").String()))
	if err != nil || abi == "" {
		gethwrappers.Exit("could not extract LINK ABI", err)
	}
	abiPath := filepath.Join(tmpDir, "abi")
	if aErr := ioutil.WriteFile(abiPath, []byte(abi), 0600); aErr != nil {
		gethwrappers.Exit("could not write contract ABI to temp dir.", aErr)
	}
	bin := gjson.Get(string(linkDetails), "bytecode").String()
	if bin == "" {
		gethwrappers.Exit("could not extract LINK bytecode", nil)
	}
	binPath := filepath.Join(tmpDir, "bin")
	if bErr := ioutil.WriteFile(binPath, []byte(bin), 0600); bErr != nil {
		gethwrappers.Exit("could not write contract binary to temp dir.", bErr)
	}
	cwd, err := os.Getwd()
	if err != nil {
		gethwrappers.Exit("could not get working directory", nil)
	}
	if filepath.Base(cwd) != "gethwrappers" {
		gethwrappers.Exit("must be run from gethwrappers directory", nil)
	}
	outDir := filepath.Join(cwd, "generated", pkgName)
	if err := os.MkdirAll(outDir, 0700); err != nil {
		gethwrappers.Exit("failed to create wrapper dir", err)
	}
	gethwrappers.Abigen(gethwrappers.AbigenArgs{
		Bin:  binPath,
		ABI:  abiPath,
		Out:  filepath.Join(outDir, pkgName+".go"),
		Type: className,
		Pkg:  pkgName,
	})
}
