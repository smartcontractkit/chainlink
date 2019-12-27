// package vrf_test verifies correct-up-to-date generation of golang wrappers
// for solidity contracts. See go_generate.go for the actual generation.
package vrf_test

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	gethParams "github.com/ethereum/go-ethereum/params"

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

// TestCheckContractHashesFromLastGoGenerate compares the metadata in recorded
// by record_versions.sh, and fails if it indicates that the corresponding
// golang wrappers are out of date with respect to the solidty contracts they
// wrap. See record_versions.sh for description of file format.
func TestCheckContractHashesFromLastGoGenerate(t *testing.T) {
	versions := readVersionsDB(t)
	require.NotEmpty(t, versions.gethVersion, `version DB should have a "GETH_VERSION:" line`)
	require.Equal(t, versions.gethVersion, gethParams.Version)
	for _, contractVersionInfo := range versions.contractVersions {
		compareCurrentCompilerAritfactAgainstRecordsAndSoliditySources(t, contractVersionInfo)
	}
}

// compareCurrentCompilerAritfactAgainstRecordsAndSoliditySources checks that
// the file at path hashes to hash, and that the solidity source code recorded
// at path match the current solidity contracts.
//
// The contents of the file at path should contain output from sol-compiler, or
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
	path := versionInfo.compilerArtifactPath

	// check the compiler artifact hasn't changed
	compilerJSON, err := ioutil.ReadFile(path)
	require.NoError(t, err, "failed to read JSON compiler artifact %s", path)
	hasher := md5.New()
	_, err = io.WriteString(hasher, string(compilerJSON))
	require.NoError(t, err, "failed to hash compiler artifact %s", path)
	recompileCommand := "`yarn workspace chainlinkv0.5 compile; go generate`"
	require.Equal(t, fmt.Sprintf("%x", hasher.Sum(nil)), versionInfo.hash,
		"compiler artifact %s has changed; please rerun %s for the vrf package",
		path, recompileCommand)

	var artifact struct {
		Sources map[string]string `json:"sourceCodes"`
	}
	require.NoError(t, json.Unmarshal(compilerJSON, &artifact),
		"could not read compiler artifact %s", path)

	// Check that each of the contract source codes hasn't changed
	soliditySourceRoot := filepath.Dir(filepath.Dir(filepath.Dir(path)))
	contractPath := filepath.Join(soliditySourceRoot, "contracts")
	for sourcePath, sourceCode := range artifact.Sources { // compare to current source
		sourcePath = filepath.Join(contractPath, sourcePath)
		actualSource, err := ioutil.ReadFile(sourcePath)
		require.NoError(t, err, "could not read "+sourcePath)
		require.Equal(t, string(actualSource), sourceCode,
			"%s has changed; please rerun %sfor the vrf package",
			sourcePath, recompileCommand)
	}
}

var stripTrailingColon = regexp.MustCompile(":$").ReplaceAllString

// readVersionsDB populates an integratedVersion with all the info in the
// versions DB
func readVersionsDB(t *testing.T) integratedVersion {
	rv := integratedVersion{}
	rv.contractVersions = make(map[string]contractVersion)

	dirOfThisTest, err := os.Getwd()
	require.NoError(t, err)
	dBPath := "generated-wrapper-dependency-versions-do-not-edit.txt"
	versionsDBFile, err := os.Open(filepath.Join(dirOfThisTest, "generation", dBPath))
	require.NoError(t, err, "could not open versions database")
	versionsDBLineReader := bufio.NewScanner(versionsDBFile)

	for versionsDBLineReader.Scan() {
		line := strings.Fields(versionsDBLineReader.Text())
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
				versionsDBLineReader.Text())
			_, alreadyExists := rv.contractVersions[topic]
			require.False(t, alreadyExists, `topic "%s" already mentioned!`, topic)
			rv.contractVersions[topic] = contractVersion{
				compilerArtifactPath: line[1], hash: line[2],
			}
		}
	}
	return rv
}
