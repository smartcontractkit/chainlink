package helpers

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"

	cldlogger "github.com/smartcontractkit/chainlink/deployment/logger"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
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

	return &ttypes.TestEnvironment{
		Config:         in,
		TestConfig:     tconf,
		Logger:         framework.L,
		CreEnvironment: creEnvironment,
		Dons:           dons,
	}
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

	recreateErr := recreateEnvironmentIfIncompatible(t.Context(), testConfig.RelativePathToRepoRoot, testConfig.EnvironmentDirPath, testConfig.EnvironmentConfigPath, flags...)
	require.NoError(t, recreateErr, "failed to recreate incompatible environment")

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

		if err := startEnvironment(ctx, environmentDir, flags...); err != nil {
			return err
		}
	}

	return nil
}

func recreateEnvironmentIfIncompatible(ctx context.Context, relativePathToRepoRoot, environmentDir, requestedConfigPath string, flags ...string) error {
	if !envconfig.LocalCREStateFileExists(relativePathToRepoRoot) {
		return nil
	}

	compatible, err := localEnvironmentSatisfiesRequestedConfig(relativePathToRepoRoot, requestedConfigPath)
	if err != nil {
		return errors.Wrap(err, "failed to compare requested config with saved environment")
	}
	if compatible {
		return nil
	}

	framework.L.Info().
		Str("requested_config", requestedConfigPath).
		Str("local CRE state file", envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)).
		Msg("Saved local CRE state is incompatible with requested test config, recreating environment")

	if err := hydrateRecreatedEnvironmentImageEnv(ctx, relativePathToRepoRoot); err != nil {
		framework.L.Warn().Err(err).Msg("failed to hydrate recreated environment image env vars from running containers")
	}

	stopCmd := exec.CommandContext(ctx, "go", "run", ".", "env", "stop")
	stopCmd.Dir = environmentDir
	stopCmd.Stdout = os.Stdout
	stopCmd.Stderr = os.Stderr
	if err := stopCmd.Run(); err != nil {
		return errors.Wrap(err, "failed to stop incompatible environment")
	}

	previousCTFConfigs := os.Getenv("CTF_CONFIGS")
	if err := os.Setenv("CTF_CONFIGS", requestedConfigPath); err != nil {
		return errors.Wrap(err, "failed to set requested CTF_CONFIGS before environment recreation")
	}
	defer func() {
		_ = os.Setenv("CTF_CONFIGS", previousCTFConfigs)
	}()

	return startEnvironment(ctx, environmentDir, flags...)
}

func startEnvironment(ctx context.Context, environmentDir string, flags ...string) error {
	args := []string{"run", ".", "env", "start"}
	args = append(args, flags...)

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = environmentDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to start environment")
	}

	return nil
}

func hydrateRecreatedEnvironmentImageEnv(ctx context.Context, relativePathToRepoRoot string) error {
	if err := hydrateRecreatedEnvironmentImageEnvFromSavedState(relativePathToRepoRoot); err != nil {
		return err
	}

	containers, err := runningContainers(ctx)
	if err != nil {
		return err
	}

	setEnvIfMissing("CTF_JD_IMAGE", firstMatchingContainerImage(containers, func(name string) bool {
		return strings.HasPrefix(name, "job-distributor")
	}))
	setEnvIfMissing("CTF_CHAINLINK_IMAGE", firstMatchingContainerImage(containers, func(name string) bool {
		return strings.HasPrefix(name, "workflow-node") ||
			strings.HasPrefix(name, "capabilities-node") ||
			strings.HasPrefix(name, "bootstrap-gateway")
	}))
	setEnvIfMissing("CHIP_INGRESS_IMAGE", firstMatchingContainerImage(containers, func(name string) bool {
		return strings.Contains(name, "chip-ingress")
	}))
	setEnvIfMissing("BILLING_PLATFORM_SERVICE_IMAGE", firstMatchingContainerImage(containers, func(name string) bool {
		return strings.Contains(name, "billing-platform-service")
	}))

	return nil
}

func hydrateRecreatedEnvironmentImageEnvFromSavedState(relativePathToRepoRoot string) error {
	if !envconfig.LocalCREStateFileExists(relativePathToRepoRoot) {
		return nil
	}

	current := &envconfig.Config{}
	if err := current.Load(envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)); err != nil {
		return errors.Wrap(err, "failed to load saved local CRE state while hydrating image env vars")
	}

	hydrateRecreatedEnvironmentImageEnvFromConfig(current)
	return nil
}

func hydrateRecreatedEnvironmentImageEnvFromConfig(current *envconfig.Config) {
	if current == nil {
		return
	}
	if current.JD != nil {
		setEnvIfMissing("CTF_JD_IMAGE", current.JD.Image)
	}
	setEnvIfMissing("CTF_CHAINLINK_IMAGE", firstConfiguredNodeImage(current.NodeSets))
}

type runningContainer struct {
	Image string
	Name  string
}

func runningContainers(ctx context.Context) ([]runningContainer, error) {
	cmd := exec.CommandContext(ctx, "docker", "ps", "--format", "{{.Image}}\t{{.Names}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list running docker containers")
	}
	return parseRunningContainers(string(out)), nil
}

func parseRunningContainers(output string) []runningContainer {
	var containers []runningContainer
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		containers = append(containers, runningContainer{
			Image: strings.TrimSpace(parts[0]),
			Name:  strings.TrimSpace(parts[1]),
		})
	}
	return containers
}

func firstMatchingContainerImage(containers []runningContainer, match func(name string) bool) string {
	for _, container := range containers {
		if container.Image == "" || container.Name == "" {
			continue
		}
		if match(container.Name) {
			return container.Image
		}
	}
	return ""
}

func firstConfiguredNodeImage(nodeSets []*cre.NodeSet) string {
	for _, nodeSet := range nodeSets {
		if nodeSet == nil {
			continue
		}
		for _, nodeSpec := range nodeSet.NodeSpecs {
			if nodeSpec == nil || nodeSpec.Input == nil || nodeSpec.Node == nil {
				continue
			}
			if nodeSpec.Node.Image != "" {
				return nodeSpec.Node.Image
			}
		}
	}
	return ""
}

func setEnvIfMissing(key, value string) {
	if value == "" || os.Getenv(key) != "" {
		return
	}
	_ = os.Setenv(key, value)
}

func localEnvironmentSatisfiesRequestedConfig(relativePathToRepoRoot, requestedConfigPath string) (bool, error) {
	requested := &envconfig.Config{}
	if err := requested.Load(requestedConfigPath); err != nil {
		return false, errors.Wrap(err, "failed to load requested environment config")
	}

	current := &envconfig.Config{}
	if err := current.Load(envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)); err != nil {
		return false, errors.Wrap(err, "failed to load saved local CRE state")
	}

	return savedEnvironmentSatisfiesRequestedConfig(requested, current)
}

func savedEnvironmentSatisfiesRequestedConfig(requested, current *envconfig.Config) (bool, error) {
	required, err := requiredChainKeys(requested)
	if err != nil {
		return false, err
	}
	available, err := requiredChainKeys(current)
	if err != nil {
		return false, err
	}

	for _, req := range required {
		if !slices.Contains(available, req) {
			return false, nil
		}
	}

	return true, nil
}

func requiredChainKeys(cfg *envconfig.Config) ([]string, error) {
	keys := make([]string, 0, len(cfg.Blockchains))
	for _, in := range cfg.Blockchains {
		if in == nil {
			continue
		}

		family, chainID, err := chainKey(in)
		if err != nil {
			return nil, err
		}
		keys = append(keys, family+":"+chainID)
	}

	return keys, nil
}

func chainKey(in *blockchain.Input) (string, string, error) {
	if in == nil {
		return "", "", errors.New("nil blockchain input")
	}

	family := ""
	if in.Out != nil && in.Out.Family != "" {
		family = in.Out.Family
	} else {
		derived, err := blockchain.TypeToFamily(in.Type)
		if err != nil {
			return "", "", errors.Wrap(err, "failed to derive blockchain family")
		}
		family = string(derived)
	}

	chainID := in.ChainID
	if chainID == "" && in.Out != nil {
		chainID = in.Out.ChainID
	}
	if chainID == "" {
		return "", "", errors.New("blockchain chain ID is empty")
	}

	return family, chainID, nil
}
