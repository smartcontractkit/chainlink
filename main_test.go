package main

import (
	"embed"
	"io/fs"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core"
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/static"
)

//go:embed testdata/**/*txtar
var testFs embed.FS

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"chainlink": core.Main,
	}))
}

func TestScripts(t *testing.T) {
	t.Parallel()
	testDataRootDir := "testdata"
	visitFn := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && hasScripts(t, path) {
			t.Run(path, func(t *testing.T) {
				t.Parallel()
				testscript.Run(t, testscript.Params{
					Dir:   path,
					Setup: commonEnv(t),
				})
			})
		}
		return nil
	}

	require.NoError(t, fs.WalkDir(testFs, testDataRootDir, visitFn))
}

func hasScripts(t *testing.T, dir string) bool {
	t.Helper()
	matches, err := fs.Glob(os.DirFS(dir), "*txtar")
	require.NoError(t, err)
	return len(matches) > 0
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
