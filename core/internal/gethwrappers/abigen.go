package gethwrappers

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"

	gethParams "github.com/ethereum/go-ethereum/params"
)

// AbigenArgs is the arguments to the abigen executable. E.g., Bin is the -bin
// arg.
type AbigenArgs struct {
	Bin, ABI, Out, Type, Pkg string
}

// Abigen calls Abigen  with the given arguments
//
// It might seem like a shame, to shell out to another golang program like
// this, but the abigen executable is the stable public interface to the
// geth contract-wrapper machinery.
//
// Check whether native abigen is installed, and has correct version
func Abigen(a AbigenArgs) {
	var versionResponse bytes.Buffer
	abigenExecutablePath := filepath.Join(GetProjectRoot(), "tools/bin/abigen")
	abigenVersionCheck := exec.Command(abigenExecutablePath, "--version")
	abigenVersionCheck.Stdout = &versionResponse
	if err := abigenVersionCheck.Run(); err != nil {
		Exit("no native abigen; you must install it (`make abigen` in the "+
			"chainlink root dir)", err)
	}
	version := string(regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`).Find(
		versionResponse.Bytes()))
	if version != gethParams.Version {
		Exit(fmt.Sprintf("wrong version (%s) of abigen; install the correct one "+
			"(%s) with `make abigen` in the chainlink root dir", version,
			gethParams.Version),
			nil)
	}
	buildCommand := exec.Command(
		abigenExecutablePath,
		"-bin", a.Bin,
		"-abi", a.ABI,
		"-out", a.Out,
		"-type", a.Type,
		"-pkg", a.Pkg,
	)
	var buildResponse bytes.Buffer
	buildCommand.Stderr = &buildResponse
	if err := buildCommand.Run(); err != nil {
		Exit("failure while building "+a.Pkg+" wrapper, stderr: "+
			buildResponse.String(), err)
	}
}
