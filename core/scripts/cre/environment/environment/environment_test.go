package environment

import (
	"context"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/sets"
	gateway "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

// isDockerRunningOnPort checks if any docker container is running and exposing the specified port
func isDockerRunningOnPort(port string) (bool, error) {
	cmd := exec.Command("docker", "ps", "--format", "table {{.Ports}}")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	// Check if the output contains the specified port
	return strings.Contains(string(output), port), nil
}

// temp test to debug environment creation
func TestEnvironment(t *testing.T) {
	// Skip in CI to avoid resource usage - this is for local debugging only
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping debug test in CI")
	}

	// Change to the parent directory so config paths work correctly
	originalDir, err := os.Getwd()
	require.NoError(t, err, "failed to get working directory")

	err = os.Chdir("..")
	require.NoError(t, err, "failed to change to parent directory")

	// Check if docker containers are running on port 8545 before attempting to stop
	dockerRunning, dockerErr := isDockerRunningOnPort("8545")
	if dockerErr != nil {
		t.Logf("Warning: failed to check docker status: %v", dockerErr)
	}

	if dockerRunning {
		t.Logf("Docker containers found on port 8545, stopping environment...")
		cmd := exec.Command("go", "run", ".", "env", "stop")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		require.NoError(t, err, "failed to stop environment")
	} else {
		t.Logf("No docker containers found on port 8545, skipping stop command")
	}

	defer func() {
		_ = os.Chdir(originalDir)
	}()

	// Set default config if not provided (include capability defaults like CLI does)
	if os.Getenv("CTF_CONFIGS") == "" {
		err := os.Setenv("CTF_CONFIGS", "configs/capability_defaults.toml,configs/workflow-don-tron.toml")
		require.NoError(t, err, "failed to set CTF_CONFIGS")
	}

	// Load configuration
	in, err := framework.Load[envconfig.Config](nil)
	require.NoError(t, err, "failed to load config")

	for _, nodeSet := range in.NodeSets {
		err := nodeSet.ParseChainCapabilities()
		require.NoError(t, err, "failed to parse chain capabilities")

		err = nodeSet.ValidateChainCapabilities(in.Blockchains)
		require.NoError(t, err, "failed to validate chain capabilities")
	}

	// Set environment variables that CLI normally sets
	err = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	require.NoError(t, err)

	err = creenv.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey)
	require.NoError(t, err)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	t.Logf("Starting environment with defaults...")

	homeChainIDInt, chainErr := strconv.Atoi(in.Blockchains[0].ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	defaultCapabilities, defaultCapabilitiesErr := sets.NewDefaultSet(libc.MustSafeUint64FromInt(homeChainIDInt))
	require.NoError(t, defaultCapabilitiesErr, "failed to create default capabilities")

	// Create environment dependencies like the CLI does
	envDependencies := cre.NewEnvironmentDependencies(
		flags.NewDefaultCapabilityFlagsProvider(),
		cre.NewContractVersionsProvider(envconfig.DefaultContractSet(false)), // withV2Registries = false
		cre.NewCLIFlagsProvider(false),                                       // withV2Registries = false
	)

	// Create extra job spec functions like the CLI does
	extraJobSpecFunctions := []cre.JobSpecFn{
		gateway.JobSpec([]int{in.Fake.Port}, []string{}, []string{"0.0.0.0/0"}),
	}

	// Call StartCLIEnvironment with all required parameters - SET BREAKPOINTS HERE
	output, err := StartCLIEnvironment(
		ctx,
		".", // relativePathToRepoRoot - we're in the right directory after chdir
		in,
		"workflow", // topology
		"",         // withPluginsDockerImageFlag - no custom docker image
		defaultCapabilities,
		extraJobSpecFunctions,
		envDependencies,
	)

	require.NoError(t, err, "StartCLIEnvironment should succeed")
	require.NotNil(t, output, "output should not be nil")

	t.Logf("Environment created successfully!")
	t.Logf("Blockchains: %d, Node sets: %d", len(output.BlockchainOutput), len(output.NodeOutput))
}
