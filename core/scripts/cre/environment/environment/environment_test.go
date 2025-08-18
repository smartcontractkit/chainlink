package environment

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/stretchr/testify/require"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities/sets"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

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

	// defer func() {
	// 	cmd := exec.Command("go", "run", ".", "env", "stop")
	// 	cmd.Stdout = os.Stdout
	// 	cmd.Stderr = os.Stderr
	// 	err := cmd.Run()
	// 	require.NoError(t, err, "failed to stop environment")
	// }()

	homeChainIDInt, chainErr := strconv.Atoi(in.Blockchains[0].ChainID)
	require.NoError(t, chainErr, "failed to convert chain ID to int")

	defaultCapabilities, defaultCapabilitiesErr := sets.NewDefaultSet(libc.MustSafeUint64FromInt(homeChainIDInt), []int{in.Fake.Port}, []string{}, []string{"0.0.0.0/0"})
	require.NoError(t, defaultCapabilitiesErr, "failed to create default capabilities")

	capabilityFlagsProvider := flags.NewDefaultCapabilityFlagsProvider()

	// Call StartCLIEnvironment with defaults - SET BREAKPOINTS HERE
	output, err := StartCLIEnvironment(
		ctx,
		in,
		"workflow", // topology
		"",         // no custom docker image
		defaultCapabilities,
		capabilityFlagsProvider,
	)

	require.NoError(t, err, "StartCLIEnvironment should succeed")
	require.NotNil(t, output, "output should not be nil")

	t.Logf("Environment created successfully!")
	t.Logf("Blockchains: %d, Node sets: %d", len(output.BlockchainOutput), len(output.NodeOutput))
}
