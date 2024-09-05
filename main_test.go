package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
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
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"chainlink": core.Main,
	}))
}

// TestScripts walks through the testdata/scripts directory and runs all tests that end in
// .txt or .txtar with the testscripts library. To run an individual test, specify it in the
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

			testscript.Run(t, testscript.Params{
				Dir:             path,
				Setup:           commonEnv,
				ContinueOnError: true,
				// UpdateScripts:   true, // uncomment to update golden files
			})
		})
		return nil
	})

	require.NoError(t, visitor.Walk())
}

func commonEnv(te *testscript.Env) error {
	te.Setenv("HOME", "$WORK/home")
	te.Setenv("VERSION", static.Version)
	te.Setenv("COMMIT_SHA", static.Sha)

	b, err := os.ReadFile(filepath.Join(te.WorkDir, testPortName))
	if err != nil && !os.IsNotExist(err) {
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
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read file %s: %w", testDBName, err)
	} else if err == nil {
		envVarName := strings.TrimSpace(string(b))
		te.T().Log("test database requested:", envVarName)

		u2, err2 := initDB()
		if err2 != nil {
			return err2
		}

		te.Setenv(envVarName, u2)
	}
	return nil
}

func takeFreePort() (int, func(), error) {
	ports, err := freeport.Take(1)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get free port: %w", err)
	}
	return ports[0], func() { freeport.Return(ports) }, nil
}

func initDB() (string, error) {
	u, err := url.Parse(string(env.DatabaseURL.Get()))
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	name := strings.ReplaceAll(uuid.NewString(), "-", "_") + "_test"
	u2, err := testdb.CreateOrReplace(*u, name, true)
	if err != nil {
		return "", fmt.Errorf("failed to create DB: %w", err)
	}
	return u2, nil
}
