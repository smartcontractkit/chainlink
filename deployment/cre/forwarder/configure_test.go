package forwarder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	"github.com/smartcontractkit/chainlink/deployment/cre/test"
)

func TestConfigureForwardersSeq(t *testing.T) {
	// Setup test environment with a single DON and capabilities registry
	envWrapper := test.SetupEnvV2(t, false)
	env := envWrapper.Env
	registryChainSel := envWrapper.RegistrySelector

	// Deploy Keystone Forwarder contracts to the test chains using the forwarder.DeploySequence

	b := optest.NewBundle(t)
	deps := forwarder.DeploySequenceDeps{
		Env: env,
	}
	input := forwarder.DeploySequenceInput{
		Targets:   []uint64{registryChainSel},
		Qualifier: "test-forwarder",
	}

	got, err := operations.ExecuteSequence(b, forwarder.DeploySequence, deps, input)
	require.NoError(t, err)
	// Check that the output has the address
	addrRefs, err := got.Output.Addresses.Fetch()
	require.NoError(t, err)
	require.Len(t, addrRefs, len(input.Targets))
	require.NotEmpty(t, got.Output.Datastore)

	env.DataStore = got.Output.Datastore

	// Create test DON configuration
	// Using test nodes from the environment wrapper
	testNodeIDs := []string{
		"p2p_12D3KooWBCF1xt5Wi8FzfgNCqRL76Swv8TRU3TiD4QiJm8NMNX7N",
		"p2p_12D3KooWG1AyvwmCpZ93J8pBQUE1SuzrjDXnT4BeouncHR3jWLCG",
		"p2p_12D3KooWGeUKZBRMbx27FUTgBwZa9Ap9Ym92E9orQ8suaJuKYa8g",
		"p2p_12D3KooWH3uVF6wv47WnArKHk5p6cvgCJEb74UTmxztmQDc298wb",
	}

	donConfig := forwarder.DonConfiguration{
		Name:    "testDON",
		ID:      1,
		F:       1,
		Version: 1,
		NodeIDs: testNodeIDs,
	}

	// Setup dependencies for ConfigureForwardersSeq
	deps2 := forwarder.ConfigureForwardersSeqDeps{
		Env: env,
	}

	// Setup input for ConfigureForwardersSeq
	input2 := forwarder.ConfigureForwardersSeqInput{
		DON: donConfig,
		// Not using MCMS for this test
		MCMSConfig: nil,
		// Empty chains means run for all available chains
		Chains: map[uint64]struct{}{},
	}

	// Create operations bundle for testing
	b2 := optest.NewBundle(t)

	// Execute the ConfigureForwardersSeq operation
	output, err := operations.ExecuteSequence(b2, forwarder.ConfigureForwardersSeq, deps2, input2)
	require.NoError(t, err, "ConfigureForwardersSeq should execute successfully")
	require.NotNil(t, output, "ConfigureForwardersSeq should return output")

}
