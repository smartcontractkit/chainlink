package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/freeport"

	"github.com/smartcontractkit/chainlink/v2/core"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/internal/testdb"
	"github.com/smartcontractkit/chainlink/v2/tools/txtar"
)

// special files can be included to allocate additional test resources
const (
	// testDBName triggers initializing of a test database.
	// The URL will be set as the value of an env var named by the file.
	//
	//	-- testdb.txt --
	//	CL_DATABASE_URL
	testDBName = "testdb.txt"
	// testPortName triggers injection of a free port as the value of an env var named by the file.
	//
	//	-- testport.txt --
	//	PORT
	testPortName = "testport.txt"
	// integrationBuildName acts like a build tag: //go:build integration
	integrationBuildName = "go:build.integration"
)

// updateScripts enables updating testscript golden files, like `go test . -update`
var updateScripts = flag.Bool("update", false, "update testscript golden files")

func TestMain(m *testing.M) {
	// keep GOTMPDIR short: osx default is too long for go-plugin sockets.
	// Not removed afterwards because testscript.Main never returns (os.Exit).
	tmp, err := os.MkdirTemp("", "chainlink-testscripts")
	if err != nil {
		log.Fatalf("failed to create temp dir: %v", err)
	}
	os.Setenv("GOTMPDIR", tmp)

	testscript.Main(m, map[string]func(){
		"chainlink": func() { os.Exit(core.Main()) },
	})
}

// Temporary workaround for skipping flaky tests as we improve our tracking process
var skipFlakyTests = map[string]string{ // test name: issue number
	// "TestScripts/nodes/evm/list/list":       "https://smartcontract-it.atlassian.net/browse/DX-107",
	// "TestScripts/keys/eth/list/unavailable": "https://smartcontract-it.atlassian.net/browse/DX-110",
}

// TestScripts walks through the testdata/scripts directory and runs all .txtar
// files with the testscripts library. To run an individual test, specify it in the
// -run param of go test without the txtar or txt suffix, like so:
// go test . -run TestScripts/node/validate/default
func TestScripts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testscript")
	}

	t.Parallel()

	visitor := txtar.NewDirVisitor("testdata/scripts", txtar.Recurse, func(path string) error {
		t.Run(strings.TrimPrefix(path, "testdata/scripts/"), func(t *testing.T) {
			t.Parallel()

			// Check each .txtar file against skipFlakyTests
			matches, err := filepath.Glob(filepath.Join(path, "*.txtar"))
			require.NoError(t, err)

			var filesToRun []string
			for _, match := range matches {
				scriptName := strings.TrimSuffix(filepath.Base(match), ".txtar")
				fullTestName := t.Name() + "/" + scriptName

				if message, shouldSkip := skipFlakyTests[fullTestName]; shouldSkip {
					t.Logf("Skipping Flaky Test: %s - %s", fullTestName, message)
					continue
				}
				filesToRun = append(filesToRun, match)
			}

			if len(filesToRun) == 0 {
				t.Skip("all scripts in directory skipped")
			}

			testscript.Run(t, testscript.Params{
				Files:               filesToRun,
				Setup:               commonEnv(),
				ContinueOnError:     true,
				RequireExplicitExec: true,
				UpdateScripts:       *updateScripts,
			})
		})
		return nil
	})

	require.NoError(t, visitor.Walk())
}

// isIntegrationBuild is toggled true by a func init() with a //go:build integration gate
var isIntegrationBuild = false

func commonEnv() func(*testscript.Env) error {
	return func(te *testscript.Env) error {
		if _, err := os.Stat(filepath.Join(te.WorkDir, integrationBuildName)); err == nil && !isIntegrationBuild {
			te.T().Skip("integration test")
			return nil
		}

		home := filepath.Join(te.WorkDir, "home")
		if err := os.MkdirAll(home, 0o777); err != nil {
			return fmt.Errorf("failed to create home dir %s: %w", home, err)
		}
		te.Setenv("HOME", home)
		te.Setenv("VERSION", static.Version)
		te.Setenv("VERSION_TAG", static.VersionTag)
		te.Setenv("COMMIT_SHA", static.Sha)
		te.Setenv("TMPDIR", "/tmp") // osx default is too long for go-plugin sockets

		b, err := os.ReadFile(filepath.Join(te.WorkDir, testPortName))
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("failed to read file %s: %w", testPortName, err)
		} else if err == nil {
			envVarName := strings.TrimSpace(string(b))
			te.T().Log("test port requested:", envVarName)

			port, ret, err2 := takeFreePort()
			if err2 != nil {
				return err2
			}
			te.Defer(ret)

			te.Setenv(envVarName, strconv.Itoa(port))
		}

		b, err = os.ReadFile(filepath.Join(te.WorkDir, testDBName))
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("failed to read file %s: %w", testDBName, err)
		} else if err == nil {
			envVarName := strings.TrimSpace(string(b))
			te.T().Log("test database requested:", envVarName)

			// use a script-scoped TB so the DB is dropped when the script ends,
			// not when the whole TestScripts suite finishes
			u2 := testdb.New(scriptTB{te: te}, true).String()

			te.Setenv(envVarName, u2)
		}
		return nil
	}
}

// scriptTB adapts a testscript Env to testing.TB. Only the methods used by
// testdb (via pgtestdb and testify) are implemented; Cleanup is scoped to the
// current script via Env.Defer. Unimplemented methods panic via the nil-
// embedded interface.
type scriptTB struct {
	testing.TB
	te *testscript.Env
}

func (s scriptTB) Cleanup(f func()) { s.te.Defer(f) }
func (s scriptTB) Helper()          {}
func (s scriptTB) Failed() bool     { return false }
func (s scriptTB) FailNow()         { s.te.T().FailNow() }
func (s scriptTB) Fatal(args ...any) {
	s.te.T().Fatal(args...)
}

func (s scriptTB) Fatalf(format string, args ...any) {
	s.te.T().Fatal(fmt.Sprintf(format, args...))
}

func (s scriptTB) Error(args ...any) {
	s.te.T().Log(args...)
	s.te.T().FailNow()
}

func (s scriptTB) Errorf(format string, args ...any) {
	s.te.T().Log(fmt.Sprintf(format, args...))
	s.te.T().FailNow()
}
func (s scriptTB) Log(args ...any) { s.te.T().Log(args...) }
func (s scriptTB) Logf(format string, args ...any) {
	s.te.T().Log(fmt.Sprintf(format, args...))
}
func (s scriptTB) Name() string { return s.te.WorkDir }

func takeFreePort() (int, func(), error) {
	ports, err := freeport.Take(1)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get free port: %w", err)
	}
	return ports[0], func() { freeport.Return(ports) }, nil
}
