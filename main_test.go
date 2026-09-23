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

			matches, err := filepath.Glob(filepath.Join(path, "*.txtar"))
			require.NoError(t, err)
			if len(matches) == 0 {
				t.Skip("no scripts found")
			}

			testscript.Run(t, testscript.Params{
				Files:               matches,
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

			// Env.T implements testing.TB when started via testscript.Run
			// (it wraps the per-script *testing.T), so its Cleanup runs when
			// this script ends — not when the whole TestScripts suite finishes.
			// Do not switch this suite to testscript.RunT.
			// https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript#Env.T
			tb := te.T().(testing.TB)
			u2 := testdb.New(tb, true).String()

			te.Setenv(envVarName, u2)
		}
		return nil
	}
}

func takeFreePort() (int, func(), error) {
	ports, err := freeport.Take(1)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get free port: %w", err)
	}
	return ports[0], func() { freeport.Return(ports) }, nil
}
