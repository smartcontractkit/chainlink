package generate_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/v2/tools/githooks/internal/generate"
)

const rootMockeryConfig = `dir: "{{ .InterfaceDir }}/mocks"
mockname: "{{ .InterfaceName }}"
outpkg: mocks
filename: "{{ .InterfaceName | snakecase }}.go"
fail-on-missing: true
packages:
  github.com/smartcontractkit/chainlink/v2/core/bridges:
    interfaces:
      ORM:
  github.com/smartcontractkit/chainlink/v2/core/services/cron:
    config:
      dir: cron_mocks
    interfaces:
      Cron:
`

const deploymentMockeryConfig = `dir: "{{ .InterfaceDir }}/mocks"
mockname: "{{ .InterfaceName }}"
outpkg: mocks
filename: "{{ .InterfaceName | snakecase }}.go"
packages:
  github.com/smartcontractkit/chainlink/deployment/environment:
    interfaces:
      Env:
`

type recordRunner struct {
	mu   sync.Mutex
	runs [][]string
	dirs []string
	cfgs [][]byte
}

func (r *recordRunner) Run(ctx context.Context, dir string, args ...string) error {
	argsCopy := append([]string(nil), args...)
	var cfgContent []byte
	if len(argsCopy) == 3 && argsCopy[0] == "mockery" && argsCopy[1] == "--config" {
		raw, err := os.ReadFile(argsCopy[2])
		if err != nil {
			return err
		}
		cfgContent = raw
	}
	r.mu.Lock()
	r.runs = append(r.runs, argsCopy)
	r.dirs = append(r.dirs, dir)
	if cfgContent != nil {
		r.cfgs = append(r.cfgs, cfgContent)
	}
	r.mu.Unlock()
	return nil
}

func (r *recordRunner) runArgs(name string) [][]string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var out [][]string
	for _, run := range r.runs {
		if len(run) > 0 && run[0] == name {
			out = append(out, run)
		}
	}
	return out
}

func writeTree(t *testing.T) string {
	t.Helper()
	fsys, err := os.OpenRoot(t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, fsys.Close()) })

	require.NoError(t, fsys.WriteFile("go.mod", []byte("module github.com/smartcontractkit/chainlink/v2\n"), 0o600))
	require.NoError(t, fsys.WriteFile(".mockery.yaml", []byte(rootMockeryConfig), 0o600))
	require.NoError(t, fsys.MkdirAll("tools/githooks", 0o750))
	require.NoError(t, fsys.WriteFile("tools/githooks/go.mod", []byte("module github.com/smartcontractkit/chainlink/v2/tools/githooks\n"), 0o600))
	require.NoError(t, fsys.MkdirAll("deployment", 0o750))
	require.NoError(t, fsys.WriteFile("deployment/go.mod", []byte("module github.com/smartcontractkit/chainlink/deployment\n"), 0o600))
	require.NoError(t, fsys.WriteFile("deployment/.mockery.yaml", []byte(deploymentMockeryConfig), 0o600))
	return fsys.Name()
}

func parseScopedConfig(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var doc map[string]any
	require.NoError(t, yaml.Unmarshal(raw, &doc))
	return doc
}

func readMockeryScopedPackages(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	doc := parseScopedConfig(t, raw)
	packages, ok := doc["packages"].(map[string]any)
	require.True(t, ok, "packages should be a mapping, got %#v", doc["packages"])
	return packages
}

func TestRun_MockeryScopedForChangedPackage(t *testing.T) {
	t.Parallel()

	root := writeTree(t)
	rec := &recordRunner{}
	cfg := generate.Config{Runner: rec.Run}

	err := generate.Run(t.Context(), root, []string{"core/bridges/bridge.go"}, cfg)
	require.NoError(t, err)

	mockeryRuns := rec.runArgs("mockery")
	require.Len(t, mockeryRuns, 1)
	run := mockeryRuns[0]
	require.Len(t, run, 3)
	assert.Equal(t, "--config", run[1])

	rec.mu.Lock()
	raw := rec.cfgs[0]
	rec.mu.Unlock()

	packages := readMockeryScopedPackages(t, raw)
	require.Contains(t, packages, "github.com/smartcontractkit/chainlink/v2/core/bridges")
	assert.NotContains(t, packages, "github.com/smartcontractkit/chainlink/v2/core/services/cron")

	doc := parseScopedConfig(t, raw)
	assert.Equal(t, "{{ .InterfaceDir }}/mocks", doc["dir"])
	assert.Equal(t, "{{ .InterfaceName }}", doc["mockname"])
	assert.Equal(t, true, doc["fail-on-missing"])
}

func TestRun_MockeryGroupsMultipleChangedPackages(t *testing.T) {
	t.Parallel()

	root := writeTree(t)
	rec := &recordRunner{}
	cfg := generate.Config{Runner: rec.Run}

	files := []string{
		"core/bridges/bridge.go",
		"core/services/cron/cron.go",
	}
	err := generate.Run(t.Context(), root, files, cfg)
	require.NoError(t, err)

	mockeryRuns := rec.runArgs("mockery")
	require.Len(t, mockeryRuns, 1)
	run := mockeryRuns[0]
	require.Len(t, run, 3)
	assert.Equal(t, "--config", run[1])

	rec.mu.Lock()
	raw := rec.cfgs[0]
	rec.mu.Unlock()

	packages := readMockeryScopedPackages(t, raw)
	require.Len(t, packages, 2)
	assert.Contains(t, packages, "github.com/smartcontractkit/chainlink/v2/core/bridges")
	cronCfg, ok := packages["github.com/smartcontractkit/chainlink/v2/core/services/cron"].(map[string]any)
	require.True(t, ok)
	inner, ok := cronCfg["config"].(map[string]any)
	require.True(t, ok, "package-level config should be preserved, got %#v", cronCfg)
	assert.Equal(t, "cron_mocks", inner["dir"])
}

func TestRun_MockeryFullOnConfigChange(t *testing.T) {
	t.Parallel()

	root := writeTree(t)
	rec := &recordRunner{}
	cfg := generate.Config{Runner: rec.Run}

	err := generate.Run(t.Context(), root, []string{".mockery.yaml"}, cfg)
	require.NoError(t, err)

	mockeryRuns := rec.runArgs("mockery")
	require.Len(t, mockeryRuns, 1)
	assert.Equal(t, []string{"mockery"}, mockeryRuns[0])
}

func TestRun_MockeryFullOnCoveredModuleGoModChange(t *testing.T) {
	t.Parallel()

	root := writeTree(t)
	rec := &recordRunner{}
	cfg := generate.Config{Runner: rec.Run}

	err := generate.Run(t.Context(), root, []string{"go.mod"}, cfg)
	require.NoError(t, err)

	mockeryRuns := rec.runArgs("mockery")
	require.Len(t, mockeryRuns, 1)
	assert.Equal(t, []string{"mockery"}, mockeryRuns[0])

	modgraphRuns := rec.runArgs("modgraph")
	require.Len(t, modgraphRuns, 1)
}

func TestRun_MockerySkipsUncoveredModuleGoModChange(t *testing.T) {
	t.Parallel()

	root := writeTree(t)
	rec := &recordRunner{}
	cfg := generate.Config{Runner: rec.Run}

	err := generate.Run(t.Context(), root, []string{"tools/githooks/go.mod"}, cfg)
	require.NoError(t, err)

	require.Empty(t, rec.runArgs("mockery"))
	assert.Len(t, rec.runArgs("modgraph"), 1)
}

func TestRun_MockerySkipsUncoveredPackage(t *testing.T) {
	t.Parallel()

	root := writeTree(t)
	rec := &recordRunner{}
	cfg := generate.Config{Runner: rec.Run}

	err := generate.Run(t.Context(), root, []string{"core/logger/logger.go"}, cfg)
	require.NoError(t, err)

	require.Empty(t, rec.runArgs("mockery"))
}

func TestRun_MockeryNestedConfigWalkUp(t *testing.T) {
	t.Parallel()

	root := writeTree(t)
	rec := &recordRunner{}
	cfg := generate.Config{Runner: rec.Run}

	err := generate.Run(t.Context(), root, []string{"deployment/environment/env.go"}, cfg)
	require.NoError(t, err)

	mockeryRuns := rec.runArgs("mockery")
	require.Len(t, mockeryRuns, 1)
	run := mockeryRuns[0]
	require.Len(t, run, 3)
	assert.Equal(t, "--config", run[1])

	rec.mu.Lock()
	raw := rec.cfgs[0]
	rec.mu.Unlock()

	packages := readMockeryScopedPackages(t, raw)
	require.Contains(t, packages, "github.com/smartcontractkit/chainlink/deployment/environment")

	rec.mu.Lock()
	dir := rec.dirs[0]
	rec.mu.Unlock()
	assert.Equal(t, filepath.Join(root, "deployment"), dir)
}

func TestRun_MockeryMalformedConfigReturnsError(t *testing.T) {
	t.Parallel()

	root := writeTree(t)
	require.NoError(t, os.WriteFile(filepath.Join(root, ".mockery.yaml"), []byte("packages: [unclosed\n"), 0o600))
	rec := &recordRunner{}
	cfg := generate.Config{Runner: rec.Run}

	err := generate.Run(t.Context(), root, []string{"core/bridges/bridge.go"}, cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mockery")
}
