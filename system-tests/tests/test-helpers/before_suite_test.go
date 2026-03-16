package helpers

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func TestLocalEnvironmentSatisfiesRequestedConfig(t *testing.T) {
	t.Run("returns true when saved state contains requested chain families and ids", func(t *testing.T) {
		ok, err := savedEnvironmentSatisfiesRequestedConfig(
			&envconfig.Config{
				Blockchains: []*blockchain.Input{
					{Type: blockchain.TypeAnvil, ChainID: "1337"},
					{Type: blockchain.TypeAptos, ChainID: "4"},
				},
			},
			&envconfig.Config{
				Blockchains: []*blockchain.Input{
					{Type: blockchain.TypeAnvil, ChainID: "1337", Out: &blockchain.Output{Family: blockchain.FamilyEVM, ChainID: "1337"}},
					{Type: blockchain.TypeAptos, ChainID: "4", Out: &blockchain.Output{Family: blockchain.FamilyAptos, ChainID: "4"}},
				},
			},
		)
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("returns false when saved state is missing a requested chain", func(t *testing.T) {
		ok, err := savedEnvironmentSatisfiesRequestedConfig(
			&envconfig.Config{
				Blockchains: []*blockchain.Input{
					{Type: blockchain.TypeAnvil, ChainID: "1337"},
					{Type: blockchain.TypeAptos, ChainID: "4"},
				},
			},
			&envconfig.Config{
				Blockchains: []*blockchain.Input{
					{Type: blockchain.TypeAnvil, ChainID: "1337", Out: &blockchain.Output{Family: blockchain.FamilyEVM, ChainID: "1337"}},
					{Type: blockchain.TypeAnvil, ChainID: "2337", Out: &blockchain.Output{Family: blockchain.FamilyEVM, ChainID: "2337"}},
				},
			},
		)
		require.NoError(t, err)
		require.False(t, ok)
	})
}

func TestChainKey(t *testing.T) {
	t.Run("derives family from type when output is absent", func(t *testing.T) {
		family, chainID, err := chainKey(&blockchain.Input{Type: blockchain.TypeAptos, ChainID: "4"})
		require.NoError(t, err)
		require.Equal(t, blockchain.FamilyAptos, family)
		require.Equal(t, "4", chainID)
	})

	t.Run("uses output family and chain id when available", func(t *testing.T) {
		family, chainID, err := chainKey(&blockchain.Input{
			Type:    blockchain.TypeAnvil,
			ChainID: "",
			Out: &blockchain.Output{
				Family:  blockchain.FamilyEVM,
				ChainID: "1337",
			},
		})
		require.NoError(t, err)
		require.Equal(t, blockchain.FamilyEVM, family)
		require.Equal(t, "1337", chainID)
	})
}

func TestParseRunningContainers(t *testing.T) {
	output := "123456789012.dkr.ecr.us-east-2.amazonaws.com/job-distributor:0.22.1\tjob-distributor\n" +
		"123456789012.dkr.ecr.us-east-2.amazonaws.com/chainlink-integration-tests:abc123\tworkflow-node0\n" +
		"ghcr.io/foundry-rs/foundry:stable\tanvil-1337\n"

	containers := parseRunningContainers(output)

	require.Len(t, containers, 3)
	require.Equal(t, "123456789012.dkr.ecr.us-east-2.amazonaws.com/job-distributor:0.22.1", containers[0].Image)
	require.Equal(t, "workflow-node0", containers[1].Name)
}

func TestFirstMatchingContainerImage(t *testing.T) {
	containers := []runningContainer{
		{Image: "123456789012.dkr.ecr.us-east-2.amazonaws.com/job-distributor:0.22.1", Name: "jd-ab123"},
		{Image: "chainlink:test", Name: "workflow-node0"},
	}

	jdImage := firstMatchingContainerImage(containers, func(name string) bool {
		return strings.HasPrefix(name, "job-distributor") || strings.HasPrefix(name, "jd-")
	})
	chainlinkImage := firstMatchingContainerImage(containers, func(name string) bool {
		return strings.HasPrefix(name, "workflow-node")
	})

	require.Equal(t, "123456789012.dkr.ecr.us-east-2.amazonaws.com/job-distributor:0.22.1", jdImage)
	require.Equal(t, "chainlink:test", chainlinkImage)
}

func TestHydrateRecreatedEnvironmentImageEnvFromConfig(t *testing.T) {
	t.Setenv("CTF_JD_IMAGE", "")
	t.Setenv("CTF_CHAINLINK_IMAGE", "")

	hydrateRecreatedEnvironmentImageEnvFromConfig(&envconfig.Config{
		JD: &jd.Input{Image: "123456789012.dkr.ecr.us-east-2.amazonaws.com/job-distributor:0.22.1"},
		NodeSets: []*cre.NodeSet{
			{
				NodeSpecs: []*cre.NodeSpecWithRole{
					{
						Input: &clnode.Input{
							Node: &clnode.NodeInput{Image: "123456789012.dkr.ecr.us-east-2.amazonaws.com/chainlink:abc123"},
						},
					},
				},
			},
		},
	})

	require.Equal(t, "123456789012.dkr.ecr.us-east-2.amazonaws.com/job-distributor:0.22.1", os.Getenv("CTF_JD_IMAGE"))
	require.Equal(t, "123456789012.dkr.ecr.us-east-2.amazonaws.com/chainlink:abc123", os.Getenv("CTF_CHAINLINK_IMAGE"))
}

func TestHydrateRecreatedEnvironmentImageEnvFromConfigDoesNotOverrideExistingEnv(t *testing.T) {
	t.Setenv("CTF_JD_IMAGE", "existing-jd")
	t.Setenv("CTF_CHAINLINK_IMAGE", "existing-cl")

	hydrateRecreatedEnvironmentImageEnvFromConfig(&envconfig.Config{
		JD: &jd.Input{Image: "new-jd"},
		NodeSets: []*cre.NodeSet{
			{
				NodeSpecs: []*cre.NodeSpecWithRole{
					{
						Input: &clnode.Input{
							Node: &clnode.NodeInput{Image: "new-cl"},
						},
					},
				},
			},
		},
	})

	require.Equal(t, "existing-jd", os.Getenv("CTF_JD_IMAGE"))
	require.Equal(t, "existing-cl", os.Getenv("CTF_CHAINLINK_IMAGE"))
}
