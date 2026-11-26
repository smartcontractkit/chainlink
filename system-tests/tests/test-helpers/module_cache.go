package helpers

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	dc "github.com/moby/moby/client"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

// Node defaults for [Capabilities.WorkflowRegistry.ModuleCache], see core/config/docs/core.toml.
const (
	defaultModuleCacheIdleEviction = true
	defaultModuleCacheIdleTimeout  = 10 * time.Minute
	defaultModuleCacheMaxLoaded    = 200
)

// nodeConfigFiles are the TOML files a CTF node container starts with, in load order (later wins).
var nodeConfigFiles = []string{"/config/config", "/config/overrides", "/config/user-overrides"}

// ModuleCacheStatus is the workflow module cache configuration a node runs with.
type ModuleCacheStatus struct {
	Node         string
	Enabled      bool
	IdleEviction bool
	IdleTimeout  time.Duration
	MaxLoaded    int
}

// moduleCacheSettings holds the [Capabilities.WorkflowRegistry.ModuleCache] fields a TOML file
// sets explicitly; nil means unset.
type moduleCacheSettings struct {
	Enabled      *bool   `toml:"Enabled"`
	IdleEviction *bool   `toml:"IdleEviction"`
	IdleTimeout  *string `toml:"IdleTimeout"`
	MaxLoaded    *int    `toml:"MaxLoaded"`
}

type moduleCacheOverrides struct {
	Capabilities struct {
		WorkflowRegistry struct {
			ModuleCache moduleCacheSettings `toml:"ModuleCache"`
		} `toml:"WorkflowRegistry"`
	} `toml:"Capabilities"`
}

// RequireModuleCacheAsConfigured logs the module cache configuration of every node in nodesetName,
// read from the config files the node container started with, and fails the test unless it matches
// what the topology asks for. It catches a reused environment that was started from another topology.
func RequireModuleCacheAsConfigured(t *testing.T, testEnv *ttypes.TestEnvironment, nodesetName string) []ModuleCacheStatus {
	t.Helper()

	topologyFiles := requestedTopologyFiles(testEnv.TestConfig)
	wantBySpec := moduleCacheSettingsFromTopology(t, topologyFiles, nodesetName)

	cli := newDockerClient(t)
	containers := nodesetContainerNames(t, testEnv, nodesetName)
	statuses := make([]ModuleCacheStatus, 0, len(containers))
	for _, container := range containers {
		configs := make([]string, 0, len(nodeConfigFiles))
		for _, path := range nodeConfigFiles {
			content, err := readContainerFile(t.Context(), cli, container, path)
			require.NoError(t, err, "failed to read %s from %s", path, container)
			configs = append(configs, content)
		}
		got, err := moduleCacheStatusFromConfigs(container, configs)
		require.NoError(t, err)
		statuses = append(statuses, got)
		framework.L.Info().
			Str("node", got.Node).
			Bool("enabled", got.Enabled).
			Int("maxLoaded", got.MaxLoaded).
			Bool("idleEviction", got.IdleEviction).
			Dur("idleTimeout", got.IdleTimeout).
			Msg("Module cache status")

		want := wantBySpec[0]
		if len(wantBySpec) > 1 { // override_mode = "each": one node spec per node
			idx, err := strconv.Atoi(strings.TrimPrefix(container, nodesetName+"-node"))
			require.NoError(t, err, "cannot derive node index from container name %q", container)
			require.Less(t, idx, len(wantBySpec), "no node spec for %q in %v", container, topologyFiles)
			want = wantBySpec[idx]
		}
		if err := moduleCacheMismatch(got, want); err != nil {
			require.FailNowf(t, "module cache not as configured", "%s: %v (topology %v). If an environment was already running, "+
				"the test reused it instead of starting this topology: run `go run . env stop` in core/scripts/cre/environment and re-run",
				container, err, topologyFiles)
		}
	}
	return statuses
}

// moduleCacheStatusFromConfigs resolves the cache configuration from node config files in load order.
func moduleCacheStatusFromConfigs(node string, configs []string) (ModuleCacheStatus, error) {
	status := ModuleCacheStatus{
		Node:         node,
		IdleEviction: defaultModuleCacheIdleEviction,
		IdleTimeout:  defaultModuleCacheIdleTimeout,
		MaxLoaded:    defaultModuleCacheMaxLoaded,
	}
	for i, content := range configs {
		var overrides moduleCacheOverrides
		if err := toml.Unmarshal([]byte(content), &overrides); err != nil {
			return ModuleCacheStatus{}, fmt.Errorf("%s: parse config file %d: %w", node, i, err)
		}
		s := overrides.Capabilities.WorkflowRegistry.ModuleCache
		if s.Enabled != nil {
			status.Enabled = *s.Enabled
		}
		if s.IdleEviction != nil {
			status.IdleEviction = *s.IdleEviction
		}
		if s.MaxLoaded != nil {
			status.MaxLoaded = *s.MaxLoaded
		}
		if s.IdleTimeout != nil {
			d, err := time.ParseDuration(*s.IdleTimeout)
			if err != nil {
				return ModuleCacheStatus{}, fmt.Errorf("%s: invalid ModuleCache.IdleTimeout %q: %w", node, *s.IdleTimeout, err)
			}
			status.IdleTimeout = d
		}
	}
	return status, nil
}

// moduleCacheMismatch compares a node's status with the settings its topology sets explicitly.
func moduleCacheMismatch(got ModuleCacheStatus, want moduleCacheSettings) error {
	wantEnabled := want.Enabled != nil && *want.Enabled
	if got.Enabled != wantEnabled {
		return fmt.Errorf("module cache enabled=%t, topology wants enabled=%t", got.Enabled, wantEnabled)
	}
	if !wantEnabled {
		return nil
	}
	if want.MaxLoaded != nil && got.MaxLoaded != *want.MaxLoaded {
		return fmt.Errorf("MaxLoaded=%d, topology wants %d", got.MaxLoaded, *want.MaxLoaded)
	}
	if want.IdleEviction != nil && got.IdleEviction != *want.IdleEviction {
		return fmt.Errorf("IdleEviction=%t, topology wants %t", got.IdleEviction, *want.IdleEviction)
	}
	if want.IdleTimeout != nil {
		wantTimeout, err := time.ParseDuration(*want.IdleTimeout)
		if err != nil {
			return fmt.Errorf("invalid IdleTimeout %q in topology: %w", *want.IdleTimeout, err)
		}
		if got.IdleTimeout != wantTimeout {
			return fmt.Errorf("IdleTimeout=%s, topology wants %s", got.IdleTimeout, wantTimeout)
		}
	}
	return nil
}

func readContainerFile(ctx context.Context, cli *dc.Client, container, path string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := cli.CopyFromContainer(ctx, container, dc.CopyFromContainerOptions{SourcePath: path})
	if err != nil {
		return "", err
	}
	defer res.Content.Close()

	tr := tar.NewReader(res.Content)
	if _, err = tr.Next(); err != nil {
		return "", fmt.Errorf("read tar header of %s: %w", path, err)
	}
	content, err := io.ReadAll(tr)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	return string(content), nil
}

// requestedTopologyFiles returns the topology files `env start` uses for this test: CTF_CONFIGS
// when set (comma-separated, relative to the environment dir), otherwise the test's own topology.
func requestedTopologyFiles(tconf *ttypes.TestConfig) []string {
	fromEnv := os.Getenv("CTF_CONFIGS")
	if fromEnv == "" {
		return []string{tconf.EnvironmentConfigPath}
	}
	var files []string
	for f := range strings.SplitSeq(fromEnv, ",") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if !filepath.IsAbs(f) {
			f = filepath.Join(tconf.EnvironmentDirPath, f)
		}
		files = append(files, f)
	}
	return files
}

// moduleCacheSettingsFromTopology returns the ModuleCache settings of each node spec of nodesetName.
// Later files override earlier ones, as when `env start` merges CTF_CONFIGS.
func moduleCacheSettingsFromTopology(t *testing.T, files []string, nodesetName string) []moduleCacheSettings {
	t.Helper()

	type topology struct {
		NodeSets []struct {
			Name      string `toml:"name"`
			NodeSpecs []struct {
				Node struct {
					UserConfigOverrides string `toml:"user_config_overrides"`
				} `toml:"node"`
			} `toml:"node_specs"`
		} `toml:"nodesets"`
	}

	var specs []moduleCacheSettings
	for _, file := range files {
		raw, err := os.ReadFile(file)
		require.NoError(t, err, "failed to read topology %s", file)
		var topo topology
		require.NoError(t, toml.Unmarshal(raw, &topo), "failed to parse topology %s", file)

		for _, ns := range topo.NodeSets {
			if ns.Name != nodesetName {
				continue
			}
			specs = specs[:0]
			for _, spec := range ns.NodeSpecs {
				var overrides moduleCacheOverrides
				require.NoError(t, toml.Unmarshal([]byte(spec.Node.UserConfigOverrides), &overrides),
					"failed to parse user_config_overrides of nodeset %q in %s", nodesetName, file)
				specs = append(specs, overrides.Capabilities.WorkflowRegistry.ModuleCache)
			}
		}
	}
	require.NotEmpty(t, specs, "nodeset %q not found in topology %v", nodesetName, files)
	return specs
}
