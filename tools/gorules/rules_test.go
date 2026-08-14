package gorules_test

import (
	"path/filepath"
	"testing"

	"github.com/quasilyte/go-ruleguard/analyzer"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRules(t *testing.T) {
	t.Parallel()

	rulesPath, err := filepath.Abs("rules.go")
	require.NoError(t, err, "resolve rules path")

	err = analyzer.Analyzer.Flags.Set("rules", rulesPath)
	require.NoError(t, err, "set rules path")

	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "target")
}
