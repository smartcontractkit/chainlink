package main

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/tools/txtar"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"chainlink": core.Main,
	}))
}

func TestScripts(t *testing.T) {
	t.Parallel()

	visitor := txtar.NewDirVisitor("testdata/scripts", txtar.Recurse, func(path string) error {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			testscript.Run(t, testscript.Params{
				Dir:   path,
				Setup: commonEnv,
			})
		})
		return nil
	})

	require.NoError(t, visitor.Walk())
}

func commonEnv(env *testscript.Env) error {
	env.Setenv("HOME", "$WORK/home")
	env.Setenv("VERSION", static.Version)
	env.Setenv("COMMIT_SHA", static.Sha)
	return nil
}

// BenchmarkAdd is a dummy benchmark test to monitor for regressions in overal benchmark CI behavior.
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dummyBenchmarkAdd(1, 2)
	}
}

func dummyBenchmarkAdd(i int, j int) int {
	return i + j
}