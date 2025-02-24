package don

import (
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func Configure(t *testing.T, testLogger zerolog.Logger, input types.ConfigureDonInput) (*types.ConfigureDonOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	for i, donTopology := range input.DonTopology.MetaDons {
		if configOverrides, ok := input.DonToConfigOverrides[donTopology.ID]; ok {
			for j, configOverride := range configOverrides {
				if len(donTopology.NodeInput.NodeSpecs)-1 < j {
					return nil, errors.Errorf("config override index out of bounds: %d", j)
				}
				donTopology.NodeInput.NodeSpecs[j].Node.TestConfigOverrides = configOverride
			}
			var setErr error
			input.DonTopology.MetaDons[i].NodeOutput, setErr = config.Set(t, donTopology.NodeInput, input.BlockchainOutput)
			if setErr != nil {
				return nil, errors.Wrap(setErr, "failed to set node output")
			}
		}
	}

	nodeOutputs := make([]*types.WrappedNodeOutput, 0, len(input.DonTopology.MetaDons))
	for i := range input.DonTopology.MetaDons {
		nodeOutputs = append(nodeOutputs, input.DonTopology.MetaDons[i].NodeOutput)
	}

	// after restarting the nodes, we need to reinitialize the JD clients otherwise
	// communication between JD and nodes will fail due to invalidated session cookie
	// TODO remove if our idea with pre-generating & importing keys works and we do not need to restart the nodes
	jdOutput, jdErr := jobs.ReinitialiseJDClients(input.CldEnv, input.JdOutput, nodeOutputs...)
	if jdErr != nil {
		return nil, errors.Wrap(jdErr, "failed to reinitialize JD clients")
	}
	for _, donTopology := range input.DonTopology.MetaDons {
		if jobSpecs, ok := input.DonToJobSpecs[donTopology.ID]; ok {
			createErr := jobs.Create(input.CldEnv.Offchain, donTopology.DON, donTopology.Flags, jobSpecs)
			if createErr != nil {
				return nil, errors.Wrapf(createErr, "failed to create jobs for DON %d", donTopology.ID)
			}
		} else {
			testLogger.Warn().Msgf("No job specs found for DON %d", donTopology.ID)
		}
	}

	return &types.ConfigureDonOutput{
		JdOutput: &jdOutput.Offchain,
	}, nil
}

func BuildDONTopology(dons []*devenv.DON, nodeSetInput []*types.CapabilitiesAwareNodeSet, nodeSetOutput []*types.WrappedNodeOutput) (*types.DonTopology, error) {
	donWithMeta := make([]*types.DonWithMetadata, len(dons))

	// one DON to do everything
	if len(dons) == 1 {
		flags, err := flags.NodeSetFlags(nodeSetInput[0])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert string flags to bitmap for nodeset %s", nodeSetInput[0].Name)
		}

		donWithMeta[0] = &types.DonWithMetadata{
			DON:        dons[0],
			NodeInput:  nodeSetInput[0],
			NodeOutput: nodeSetOutput[0],
			ID:         1,
			Flags:      flags,
		}
	} else {
		for i := range dons {
			flags, err := flags.NodeSetFlags(nodeSetInput[i])
			if err != nil {
				return nil, errors.Wrapf(err, "failed to convert string flags to bitmap for nodeset %s", nodeSetInput[i].Name)
			}

			donWithMeta[i] = &types.DonWithMetadata{
				DON:        dons[i],
				NodeInput:  nodeSetInput[i],
				NodeOutput: nodeSetOutput[i],
				ID:         libc.MustSafeUint32(i + 1),
				Flags:      flags,
			}
		}
	}

	maybeID, err := flags.OneDONTopologyWithFlag(donWithMeta, types.WorkflowDON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workflow DON ID")
	}

	return &types.DonTopology{
		MetaDons:      donWithMeta,
		WorkflowDONID: maybeID.ID,
	}, nil
}

// In order to whitelist host IP in the gateway, we need to resolve the host.docker.internal to the host IP,
// and since CL image doesn't have dig or nslookup, we need to use curl.
func ResolveHostDockerInternaIP(testLogger zerolog.Logger, nsOutput *ns.Output) (string, error) {
	containerName := nsOutput.CLNodes[0].Node.ContainerName
	cmd := []string{"curl", "-v", "http://host.docker.internal"}
	output, err := framework.ExecContainer(containerName, cmd)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`.*Trying ([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+).*`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		testLogger.Error().Msgf("failed to extract IP address from curl output:\n%s", output)
		return "", errors.New("failed to extract IP address from curl output")
	}

	testLogger.Info().Msgf("Resolved host.docker.internal to %s", matches[1])

	return matches[1], nil
}
