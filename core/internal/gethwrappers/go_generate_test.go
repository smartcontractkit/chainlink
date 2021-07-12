// package gethwrappers_test verifies correct and up-to-date generation of golang wrappers
// for solidity contracts. See go_generate.go for the actual generation.
package gethwrappers

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	gethParams "github.com/ethereum/go-ethereum/params"
	"github.com/fatih/color"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const compileCommand = "../../../contracts/scripts/native_solc_compile_all"

// TestCheckContractHashesFromLastGoGenerate compares the abi and bytecode of the
// contract artifacts in contracts/solc with the abi and bytecode stored in the
// contract wrapper
func TestCheckContractHashesFromLastGoGenerate(t *testing.T) {
	versions, err := ReadVersionsDB()
	require.NoError(t, err)
	require.NotEmpty(t, versions.GethVersion, `version DB should have a "GETH_VERSION:" line`)

	wd, err := os.Getwd()
	if err != nil {
		wd = "<directory containing this test>"
	}
	require.Equal(t, versions.GethVersion, gethParams.Version,
		color.HiRedString(boxOutput("please re-run `go generate %s` and commit the"+
			"changes", wd)))

	for _, contractVersionInfo := range versions.ContractVersions {
		if isOCRContract(contractVersionInfo.AbiPath) {
			continue
		}
		compareCurrentCompilerArtifactAgainstRecordsAndSoliditySources(t, contractVersionInfo)
	}
	// Just check that LinkToken details haven't changed (they never ought to)
	linkDetails, err := ioutil.ReadFile(filepath.Join(getProjectRoot(t), "evm-test-helpers/src/LinkToken.json"))
	require.NoError(t, err, "could not read link contract details")
	require.Equal(t, fmt.Sprintf("%x", sha256.Sum256(linkDetails)),
		"27c0e17a79553fccc63a4400c6bbe415ff710d9cc7c25757bff0f7580205c922",
		"should never differ!")
}

func isOCRContract(fullpath string) bool {
	return strings.Contains(fullpath, "OffchainAggregator")
}

// rootDir is the local chainlink root working directory
var rootDir string

func init() { // compute rootDir
	var err error
	thisDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	rootDir, err = filepath.Abs(filepath.Join(thisDir, "../../.."))
	if err != nil {
		panic(err)
	}
}

// compareCurrentCompilerArtifactAgainstRecordsAndSoliditySources checks that
// the file at each ContractVersion.AbiPath and ContractVersion.BinaryPath hashes to its
// ContractVersion.Hash, and that the solidity source code recorded in the
// compiler artifact matches the current solidity contracts.
//
// Most of the compiler artifacts should contain output from sol-compiler, or
// "yarn compile". The relevant parts of its schema are
//
//    { "sourceCodes": { "<filePath>": "<code>", ... } }
//
// where <filePath> is the path to the contract, below the truffle contracts/
// directory, and <code> is the source code of the contract at the time the JSON
// file was generated.
func compareCurrentCompilerArtifactAgainstRecordsAndSoliditySources(
	t *testing.T, versionInfo ContractVersion,
) {
	hash := VersionHash(versionInfo.AbiPath, versionInfo.BinaryPath)
	recompileCommand := fmt.Sprintf("(cd %s; make go-solidity-wrappers)", rootDir)
	assert.Equal(t, versionInfo.Hash, hash,
		boxOutput(`compiled %s and/or %s has changed; please rerun
%s,
and commit the changes`, versionInfo.AbiPath, versionInfo.BinaryPath, recompileCommand))
}

// Ensure that solidity compiler artifacts are present before running this test,
// by compiling them if necessary.
func init() {
	db, err := versionsDBLineReader()
	if err != nil {
		panic(err)
	}
	var solidityArtifactsMissing []string
	for db.Scan() {
		line := strings.Fields(db.Text())
		if stripTrailingColon(line[0], "") != "GETH_VERSION" {
			if os.IsNotExist(utils.JustError(os.Stat(line[1]))) {
				solidityArtifactsMissing = append(solidityArtifactsMissing, line[1])
			}
		}
	}
	if len(solidityArtifactsMissing) == 0 {
		return
	}
	fmt.Printf("some solidity artifacts missing (%s); rebuilding...",
		solidityArtifactsMissing)
	// Don't want to run "make go-solidity-wrappers" here, because that would
	// result in an infinite loop
	cmd := exec.Command("bash", "-c", compileCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

// boxOutput formats its arguments as fmt.Printf, and encloses them in a box of
// arrows pointing at their content, in order to better highlight it. See
// ExampleBoxOutput
func boxOutput(errorMsgTemplate string, errorMsgValues ...interface{}) string {
	errorMsgTemplate = fmt.Sprintf(errorMsgTemplate, errorMsgValues...)
	lines := strings.Split(errorMsgTemplate, "\n")
	maxlen := 0
	for _, line := range lines {
		if len(line) > maxlen {
			maxlen = len(line)
		}
	}
	internalLength := maxlen + 4
	output := "↘" + strings.Repeat("↓", internalLength) + "↙\n" // top line
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	readme := strings.Repeat("README ", maxlen/7)
	output += "→  " + readme + strings.Repeat(" ", maxlen-len(readme)) + "  ←\n"
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	for _, line := range lines {
		output += "→  " + line + strings.Repeat(" ", maxlen-len(line)) + "  ←\n"
	}
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	output += "→  " + readme + strings.Repeat(" ", maxlen-len(readme)) + "  ←\n"
	output += "→  " + strings.Repeat(" ", maxlen) + "  ←\n"
	return "\n" + output + "↗" + strings.Repeat("↑", internalLength) + "↖" + // bottom line
		"\n\n"
}

func Example_boxOutput() {
	fmt.Println()
	fmt.Print(boxOutput("%s is %d", "foo", 17))
	// Output:
	// ↘↓↓↓↓↓↓↓↓↓↓↓↓↓↙
	// →             ←
	// →  README     ←
	// →             ←
	// →  foo is 17  ←
	// →             ←
	// →  README     ←
	// →             ←
	// ↗↑↑↑↑↑↑↑↑↑↑↑↑↑↖
}

// getProjectRoot returns the root of the chainlink project
func getProjectRoot(t *testing.T) (rootPath string) {
	root, err := os.Getwd()
	require.NoError(t, err, "could not get current working directory")
	for root != "/" { // Walk up path to find dir containing go.mod
		if _, err := os.Stat(filepath.Join(root, "go.mod")); os.IsNotExist(err) {
			root = filepath.Dir(root)
		} else {
			return root
		}
	}
	t.Fatal("could not find project root")
	panic("can't get here") // Appease staticcheck
}
