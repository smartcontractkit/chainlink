package main

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/require"

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

	testDataRootDir := "testdata"
	testFs := os.DirFS(".")

	visitFn := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			matches, err := fs.Glob(os.DirFS(path), "*txtar")
			if err != nil {
				return err
			}
			if len(matches) > 0 {
				t.Run(fmt.Sprintf("test scripts @ %s", path), func(t *testing.T) {
					t.Parallel()
					testscript.Run(t, testscript.Params{
						Dir:   path,
						Setup: commonEnv(t),
					})
				})
			}
		}
		return nil
	}

	require.NoError(t, fs.WalkDir(testFs, testDataRootDir, visitFn))

}
