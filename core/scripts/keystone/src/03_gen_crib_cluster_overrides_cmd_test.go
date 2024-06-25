package src

import (
	"strings"
	"testing"

	"os"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCribConfig(t *testing.T) {
	chainID := int64(11155111)
	templatesDir := "../templates"
	forwarderAddress := "0x1234567890abcdef"
	publicKeysPath := "./testdata/PublicKeys.json"

	lines := generateCribConfig(publicKeysPath, &chainID, templatesDir, forwarderAddress)

	snaps.MatchSnapshot(t, strings.Join(lines, "\n"))
}

func TestGetCribRepoPath(t *testing.T) {
	t.Run("default path test", func(t *testing.T) {

		// create a test dir that represents both the core repository and the crib repository at the same file depth
		dir, err := os.MkdirTemp("", "crib-test-**")
		assert.NoError(t, err)

		// This represents the core/scripts/keystone path in this repo, but a tempdir mock version
		scriptTmpCwd := dir + "/chainlink/core/scripts/keystone"
		err = os.MkdirAll(scriptTmpCwd, 0755)
		assert.NoError(t, err)

		// Execute this script as if we are in the core/scripts/keystone directory
		os.Chdir(scriptTmpCwd)

		_, err = getCribRepoPath("")
		assert.Error(t, err, "it should return an error if the crib repo path does not exist")

		// This represents the crib repo, but a temp mock version
		err = os.MkdirAll(dir+"/crib", 0755)
		assert.NoError(t, err)

		_, err = getCribRepoPath("")
		assert.Error(t, err, "it should return an error if the crib path is doesnt contain a .git dir")

		// create a .git dir
		err = os.Mkdir(dir+"/crib/.git", 0755)

		_, err = getCribRepoPath("")
		assert.NoError(t, err)

	})
	t.Run("custom path test", func(t *testing.T) {
		cribRepoPath := "../../../../crib"
		// create a test dir that represents both the core repository and the crib repository at the same file depth
		dir, err := os.MkdirTemp("", "crib-test-**")
		assert.NoError(t, err)

		// This represents the core/scripts/keystone path in this repo, but a tempdir mock version
		scriptTmpCwd := dir + "/chainlink/core/scripts/keystone"
		err = os.MkdirAll(scriptTmpCwd, 0755)
		assert.NoError(t, err)

		// Execute this script as if we are in the core/scripts/keystone directory
		os.Chdir(scriptTmpCwd)

		_, err = getCribRepoPath(cribRepoPath)
		assert.Error(t, err, "it should return an error if the crib repo path does not exist")

		// This represents the crib repo, but a temp mock version
		err = os.MkdirAll(dir+"/crib", 0755)
		assert.NoError(t, err)

		_, err = getCribRepoPath(cribRepoPath)
		assert.Error(t, err, "it should return an error if the crib path is doesnt contain a .git dir")

		// create a .git dir
		err = os.Mkdir(dir+"/crib/.git", 0755)

		_, err = getCribRepoPath(cribRepoPath)
		assert.NoError(t, err)

		_, err = getCribRepoPath(cribRepoPath + "idonotexist")
		assert.Error(t, err)
	})
}
