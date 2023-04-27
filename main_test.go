package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/smartcontractkit/chainlink/v2/core"
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"chainlink": core.Main,
	}))
}

func commonEnv(t *testing.T) func(env *testscript.Env) error {
	return func(env *testscript.Env) error {
		env.Setenv(string(v2.EnvDev), "true")
		env.Setenv("HOME", "$WORK/home")
		env.Setenv("VERSION", static.Version)
		env.Setenv("COMMIT_SHA", static.Sha)
		return nil
	}
}
func TestScripts(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir:   "testdata/scripts",
		Setup: commonEnv(t),
	})
}

func TestNodeDbScripts(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir:   "testdata/scripts/node/db",
		Setup: commonEnv(t),
	})
}
