package main

import (
	"embed"
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

//go:embed testdata/scripts
var scripts embed.FS

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/scripts",
		Setup: func(env *testscript.Env) error {
			env.Setenv(string(v2.EnvDev), "true")
			env.Setenv("HOME", "$WORK/home")
			env.Setenv("VERSION", static.Version)
			env.Setenv("COMMIT_SHA", static.Sha)
			return nil
		},
	})
}
