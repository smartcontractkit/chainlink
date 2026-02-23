package helpers

import (
	"context"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

func SetupTestEnvironmentWithConfig(t *testing.T, tconf *ttypes.TestConfig, flags ...string) *ttypes.TestEnvironment {
	t.Helper()

	createEnvironment(t, tconf, flags...)
	in := getEnvironmentConfig(t)
	creEnvironment, dons, err := environment.BuildFromSavedState(t.Context(), cldlogger.NewSingleFileLogger(t), in)
	require.NoError(t, err, "failed to load environment")

	testEnv := &ttypes.TestEnvironment{
		Config:         in,
		TestConfig:     tconf,
		Logger:         framework.L,
		CreEnvironment: creEnvironment,
		Dons:           dons,
	}
	ensureMixedModeComponentRelays(t, testEnv)

	t.Cleanup(func() {
		if t.Failed() {
			framework.L.Warn().Msg("Test failed - checking for panics in Docker containers...")
			foundPanics := infra.CheckContainersForPanics(framework.L, 100)
			if !foundPanics {
				var lastLines uint64 = 30
				framework.L.Warn().Msgf("No panic patterns detected in Docker container logs. Displaying last %d lines of logs for debugging:", lastLines)
				infra.PrintFailedContainerLogs(framework.L, lastLines)
			}
		}
	})

	return testEnv
}

func GetDefaultTestConfig(t *testing.T) *ttypes.TestConfig {
	t.Helper()

	return GetTestConfig(t, "/configs/workflow-gateway-don.toml")
}

func GetTestConfig(t *testing.T, configPath string) *ttypes.TestConfig {
	relativePathToRepoRoot := "../../../../"
	environmentDirPath := filepath.Join(relativePathToRepoRoot, "core/scripts/cre/environment")

	return &ttypes.TestConfig{
		RelativePathToRepoRoot: relativePathToRepoRoot,
		EnvironmentDirPath:     environmentDirPath,
		EnvironmentConfigPath:  filepath.Join(environmentDirPath, configPath), // change to your desired config, if you want to use another topology
		EnvironmentStateFile:   filepath.Join(environmentDirPath, envconfig.StateDirname, envconfig.LocalCREStateFilename),
		ChipIngressGRPCPort:    chipingressset.DEFAULT_CHIP_INGRESS_GRPC_PORT,
	}
}

func getEnvironmentConfig(t *testing.T) *envconfig.Config {
	t.Helper()

	// we call our own Load function because it executes a couple of crucial extra input transformations
	in := &envconfig.Config{}
	err := in.Load(os.Getenv("CTF_CONFIGS"))
	require.NoError(t, err, "couldn't load environment state")
	return in
}

func createEnvironment(t *testing.T, testConfig *ttypes.TestConfig, flags ...string) {
	t.Helper()

	confErr := setConfigurationIfMissing(testConfig.EnvironmentConfigPath)
	require.NoError(t, confErr, "failed to set configuration")

	createErr := createEnvironmentIfNotExists(t.Context(), testConfig.RelativePathToRepoRoot, testConfig.EnvironmentDirPath, flags...)
	require.NoError(t, createErr, "failed to create environment")

	setErr := os.Setenv("CTF_CONFIGS", envconfig.MustLocalCREStateFileAbsPath(testConfig.RelativePathToRepoRoot))
	require.NoError(t, setErr, "failed to set CTF_CONFIGS env var")
}

func setConfigurationIfMissing(configName string) error {
	if os.Getenv("CTF_CONFIGS") == "" {
		err := os.Setenv("CTF_CONFIGS", configName)
		if err != nil {
			return errors.Wrap(err, "failed to set CTF_CONFIGS env var")
		}
	}

	return environment.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey)
}

func createEnvironmentIfNotExists(ctx context.Context, relativePathToRepoRoot, environmentDir string, flags ...string) error {
	if !envconfig.LocalCREStateFileExists(relativePathToRepoRoot) {
		framework.L.Info().Str("CTF_CONFIGS", os.Getenv("CTF_CONFIGS")).Str("local CRE state file", envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)).Msg("Local CRE state file does not exist, starting environment...")

		args := []string{"run", ".", "env", "start"}
		args = append(args, flags...)

		cmd := exec.CommandContext(ctx, "go", args...)
		cmd.Dir = environmentDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			return errors.Wrap(cmdErr, "failed to start environment")
		}
	}

	return nil
}

func ensureMixedModeComponentRelays(t *testing.T, testEnv *ttypes.TestEnvironment) {
	t.Helper()
	if testEnv == nil || testEnv.Config == nil || !hasRemoteNodeSets(testEnv.Config) {
		return
	}
	nodeSetTargetsByName := map[string]string{}
	for _, nsCfg := range testEnv.Config.NodeSets {
		if nsCfg == nil {
			continue
		}
		name := strings.TrimSpace(nsCfg.Name)
		if name == "" {
			continue
		}
		nodeSetTargetsByName[name] = strings.TrimSpace(nsCfg.Target)
	}

	// Local blockchain endpoints used by remote nodesets.
	for idx, bcCfg := range testEnv.Config.Blockchains {
		if bcCfg == nil || strings.TrimSpace(string(bcCfg.Target)) != string(envconfig.TargetLocal) {
			continue
		}
		if idx >= len(testEnv.CreEnvironment.Blockchains) || testEnv.CreEnvironment.Blockchains[idx] == nil {
			continue
		}
		for nodeIdx, node := range testEnv.CreEnvironment.Blockchains[idx].CtfOutput().Nodes {
			if node == nil {
				continue
			}
			if p, ok := extractPort(node.ExternalHTTPUrl); ok {
				EnsureFixtureRelayForPort(t, testEnv, "blockchain-http-"+strconv.Itoa(idx)+"-"+strconv.Itoa(nodeIdx), p)
			}
			if p, ok := extractPort(node.ExternalWSUrl); ok {
				EnsureFixtureRelayForPort(t, testEnv, "blockchain-ws-"+strconv.Itoa(idx)+"-"+strconv.Itoa(nodeIdx), p)
			}
		}
	}

	// Local JD endpoints used by remote nodesets.
	if testEnv.Config.JD != nil && strings.TrimSpace(string(testEnv.Config.JD.Target)) == string(envconfig.TargetLocal) && testEnv.Config.JD.Out != nil {
		if p, ok := extractPort(testEnv.Config.JD.Out.ExternalGRPCUrl); ok {
			EnsureFixtureRelayForPort(t, testEnv, "jd-grpc", p)
		}
		if p, ok := extractPort(testEnv.Config.JD.Out.ExternalWSRPCUrl); ok {
			EnsureFixtureRelayForPort(t, testEnv, "jd-wsrpc", p)
		}
	}

	// Local gateway incoming ports used by remote workflow nodesets.
	if testEnv.Dons != nil && testEnv.Dons.GatewayConnectors != nil {
		for _, cfg := range testEnv.Dons.GatewayConnectors.Configurations {
			if cfg == nil || cfg.GatewayConfiguration == nil {
				continue
			}
			node, found := testEnv.Dons.NodeWithUUID(cfg.NodeUUID)
			if !found || node == nil || node.DON == nil {
				continue
			}
			donName := strings.TrimSpace(node.DON.Name)
			target := nodeSetTargetsByName[donName]
			if target != string(envconfig.TargetLocal) {
				continue
			}
			if cfg.Incoming.ExternalPort > 0 {
				EnsureFixtureRelayForPort(t, testEnv, "gateway-"+cfg.AuthGatewayID, cfg.Incoming.ExternalPort)
			}
		}
	}
}

func extractPort(raw string) (int, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, false
	}
	if strings.Contains(trimmed, "://") {
		parsed, err := url.Parse(trimmed)
		if err != nil || parsed.Port() == "" {
			return 0, false
		}
		port, convErr := strconv.Atoi(parsed.Port())
		if convErr != nil || port <= 0 || port > 65535 {
			return 0, false
		}
		return port, true
	}
	_, portRaw, err := net.SplitHostPort(trimmed)
	if err != nil {
		return 0, false
	}
	port, convErr := strconv.Atoi(portRaw)
	if convErr != nil || port <= 0 || port > 65535 {
		return 0, false
	}
	return port, true
}
