package gethwrappers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// ContractVersion records information about the solidity compiler artifact a
// golang contract wrapper package depends on.
type ContractVersion struct {
	// Hash of the artifact at the timem the wrapper was last generated
	Hash string
	// Path to compiled abi file
	AbiPath string
	// Path to compiled bin file (if exists, this can be empty)
	BinaryPath string
}

// IntegratedVersion carries the full versioning information checked in this test
type IntegratedVersion struct {
	// Version of geth last used to generate the wrappers
	GethVersion string
	// { golang-pkg-name: version_info }
	ContractVersions map[string]ContractVersion
}

func dbPath() (path string, err error) {
	dirOfThisTest, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dBBasename := "generated-wrapper-dependency-versions-do-not-edit.txt"
	return filepath.Join(dirOfThisTest, "generation", dBBasename), nil
}

func versionsDBLineReader() (*bufio.Scanner, error) {
	versionsDBPath, err := dbPath()
	if err != nil {
		return nil, errors.Wrapf(err, "could not construct versions DB path")
	}
	versionsDBFile, err := os.Open(versionsDBPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open versions database")
	}
	return bufio.NewScanner(versionsDBFile), nil

}

// readVersionsDB populates an IntegratedVersion with all the info in the
// versions DB
func ReadVersionsDB() (*IntegratedVersion, error) {
	rv := IntegratedVersion{}
	rv.ContractVersions = make(map[string]ContractVersion)
	db, err := versionsDBLineReader()
	if err != nil {
		return nil, err
	}
	for db.Scan() {
		line := strings.Fields(db.Text())
		if !strings.HasSuffix(line[0], ":") {
			return nil, errors.Errorf(
				`each line in versions.txt should start with "$TOPIC:"`)
		}
		topic := stripTrailingColon(line[0], "")
		if topic == "GETH_VERSION" {
			if len(line) != 2 {
				return nil, errors.Errorf("GETH_VERSION line should contain geth "+
					"version, and only that: %s", line)
			}
			if rv.GethVersion != "" {
				return nil, errors.Errorf("more than one geth version")
			}
			rv.GethVersion = line[1]
		} else { // It's a wrapper from a compiler artifact
			if len(line) != 4 {
				return nil, errors.Errorf(`"%s" should have four elements `+
					`"<pkgname>: <abi-path> <bin-path> <hash>"`,
					db.Text())
			}
			_, alreadyExists := rv.ContractVersions[topic]
			if alreadyExists {
				return nil, errors.Errorf(`topic "%s" already mentioned!`, topic)
			}
			rv.ContractVersions[topic] = ContractVersion{
				AbiPath: line[1], BinaryPath: line[2], Hash: line[3],
			}
		}
	}
	return &rv, nil
}

var stripTrailingColon = regexp.MustCompile(":$").ReplaceAllString

func WriteVersionsDB(db *IntegratedVersion) (err error) {
	versionsDBPath, err := dbPath()
	if err != nil {
		return errors.Wrap(err, "could not construct path to versions DB")
	}
	f, err := os.Create(versionsDBPath)
	if err != nil {
		return errors.Wrapf(err, "while opening %s", versionsDBPath)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()
	gethLine := "GETH_VERSION: " + db.GethVersion + "\n"
	n, err := f.WriteString(gethLine)
	if err != nil {
		return errors.Wrapf(err, "while recording geth version line")
	}
	if n != len(gethLine) {
		return errors.Errorf("failed to write entire geth version line, %s", gethLine)
	}
	var pkgNames []string
	for name := range db.ContractVersions {
		pkgNames = append(pkgNames, name)
	}
	sort.Strings(pkgNames)
	for _, name := range pkgNames {
		vinfo := db.ContractVersions[name]
		versionLine := fmt.Sprintf("%s: %s %s %s\n", name,
			vinfo.AbiPath, vinfo.BinaryPath, vinfo.Hash)
		n, err = f.WriteString(versionLine)
		if err != nil {
			return errors.Wrapf(err, "while recording %s version line", name)
		}
		if n != len(versionLine) {
			return errors.Errorf("failed to write entire version line %s", versionLine)
		}
	}
	return nil
}
