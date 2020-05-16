// package gethwrappers_test verifies correct and up-to-date generation of golang wrappers
// for solidity contracts. See go_generate.go for the actual generation.
package gethwrappers_test

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/utils"

	gethParams "github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/fatih/color"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// contractVersion records information about the solidity compiler artifact a
// golang contract wrapper package depends on.
type contractVersion struct {
	// path to compiler artifact used by generate.sh to create wrapper package
	compilerArtifactPath string
	// hash of the artifact at the timem the wrapper was last generated
	hash string
}

// integratedVersion carries the full versioning information checked in this test
type integratedVersion struct {
	// Version of geth last used to generate the wrappers
	gethVersion string
	// { golang-pkg-name: version_info }
	contractVersions map[string]contractVersion
}

// TestCheckContractHashesFromLastGoGenerate compares the metadata recorded by
// record_versions.sh, and fails if it indicates that the corresponding golang
// wrappers are out of date with respect to the solidty contracts they wrap. See
// record_versions.sh for description of file format.
func TestCheckContractHashesFromLastGoGenerate(t *testing.T) {
	versions := readVersionsDB(t)
	require.NotEmpty(t, versions.gethVersion,
		`version DB should have a "GETH_VERSION:" line`)
	require.Equal(t, versions.gethVersion, gethParams.Version,
		color.HiRedString(boxOutput("please re-run `go generate ./core/services/vrf` and commit the changes")))
	for _, contractVersionInfo := range versions.contractVersions {
		compareCurrentCompilerAritfactAgainstRecordsAndSoliditySources(
			t, contractVersionInfo)
	}
}

// compareCurrentCompilerAritfactAgainstRecordsAndSoliditySources checks that
// the file at each contractVersion.compilerArtifactPath hashes to its
// contractVersion.hash, and that the solidity source code recorded in the
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
func compareCurrentCompilerAritfactAgainstRecordsAndSoliditySources(
	t *testing.T, versionInfo contractVersion,
) {
	apath := versionInfo.compilerArtifactPath
	// check the compiler outputs (abi and bytecode object) haven't changed
	compilerJSON, err := ioutil.ReadFile(apath)
	require.NoError(t, err, "failed to read JSON compiler artifact %s", apath)
	abiPath := "compilerOutput.abi"
	binPath := "compilerOutput.evm.bytecode.object"
	isLINKCompilerOutput :=
		path.Base(versionInfo.compilerArtifactPath) == "LinkToken.json"
	if isLINKCompilerOutput {
		abiPath = "abi"
		binPath = "bytecode"
	}
	// Normalize the whitespace in the ABI JSON
	abiBytes, err := utils.NormalizedJSON(
		[]byte(gjson.GetBytes(compilerJSON, abiPath).String()))
	assert.NoError(t, err, "failed to normalize JSON %s", compilerJSON)
	binBytes := gjson.GetBytes(compilerJSON, binPath).String()
	if !isLINKCompilerOutput {
		// Remove the varying contract metadata, as in ./generation/generate.sh
		binBytes = binBytes[:len(binBytes)-106]
	}
	hasher := sha256.New()
	hashMsg := string(abiBytes+binBytes) + "\n" // newline from <<< in record_versions.sh
	_, err = io.WriteString(hasher, hashMsg)
	require.NoError(t, err, "failed to hash compiler artifact %s", apath)
	thisDir, err := os.Getwd()
	if err != nil {
		thisDir = "<could not get absolute path to gethwrappers package>"
	}
	recompileCommand := color.HiRedString(fmt.Sprintf("`%s && go generate %s`",
		compileCommand(t), thisDir))
	assert.Equal(t, versionInfo.hash, fmt.Sprintf("%x", hasher.Sum(nil)),
		boxOutput(`compiler artifact %s has changed; please rerun
%s
and commit the changes`, apath, recompileCommand))

	var artifact struct {
		Sources map[string]string `json:"sourceCodes"`
	}
	require.NoError(t, json.Unmarshal(compilerJSON, &artifact),
		"could not read compiler artifact %s", apath)

	if !isLINKCompilerOutput { // No need to check contract source for LINK token
		// Check that each of the contract source codes hasn't changed
		soliditySourceRoot := filepath.Dir(filepath.Dir(filepath.Dir(apath)))
		contractPath := filepath.Join(soliditySourceRoot, "src", "v0.6")
		for sourcePath, sourceCode := range artifact.Sources { // compare to current source
			sourcePath = filepath.Join(contractPath, sourcePath)
			actualSource, err := ioutil.ReadFile(sourcePath)
			require.NoError(t, err, "could not read "+sourcePath)
			// These outputs are huge, so silence them by assert.True on explicit equality
			assert.True(t, string(actualSource) == sourceCode,
				boxOutput(`Change detected in %s,
which is a dependency of %s.

For the gethwrappers package, please rerun

%s

and commit the changes`,
					sourcePath, versionInfo.compilerArtifactPath, recompileCommand))
		}
	}
}

func versionsDBLineReader() (*bufio.Scanner, error) {
	dirOfThisTest, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dBBasename := "generated-wrapper-dependency-versions-do-not-edit.txt"
	dbPath := filepath.Join(dirOfThisTest, "generation", dBBasename)
	versionsDBFile, err := os.Open(dbPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open versions database")
	}
	return bufio.NewScanner(versionsDBFile), nil

}

// readVersionsDB populates an integratedVersion with all the info in the
// versions DB
func readVersionsDB(t *testing.T) integratedVersion {
	rv := integratedVersion{}
	rv.contractVersions = make(map[string]contractVersion)
	db, err := versionsDBLineReader()
	require.NoError(t, err)
	for db.Scan() {
		line := strings.Fields(db.Text())
		require.True(t, strings.HasSuffix(line[0], ":"),
			`each line in versions.txt should start with "$TOPIC:"`)
		topic := stripTrailingColon(line[0], "")
		if topic == "GETH_VERSION" {
			require.Len(t, line, 2,
				"GETH_VERSION line should contain geth version, and only that")
			require.Empty(t, rv.gethVersion, "more than one geth version")
			rv.gethVersion = line[1]
		} else { // It's a wrapper from a json compiler artifact
			require.Len(t, line, 3,
				`"%s" should have three elements "<pkgname>: <compiler-artifact-path> <compiler-artifact-hash>"`,
				db.Text())
			_, alreadyExists := rv.contractVersions[topic]
			require.False(t, alreadyExists, `topic "%s" already mentioned!`, topic)
			rv.contractVersions[topic] = contractVersion{
				compilerArtifactPath: line[1], hash: line[2],
			}
		}
	}
	return rv
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
	cmd := exec.Command("bash", "-c", compileCommand(nil))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

var stripTrailingColon = regexp.MustCompile(":$").ReplaceAllString

// compileCommand() is a shell command which compiles chainlink's solidity
// contracts.
func compileCommand(t *testing.T) string {
	cmd, err := ioutil.ReadFile("./generation/compile_command.txt")
	if err != nil {
		if t != nil {
			t.Fatal(err)
		}
		panic(err)
	}
	return strings.Trim(string(cmd), "\n")
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
