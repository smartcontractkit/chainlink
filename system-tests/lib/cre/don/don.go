package don

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/crypto"
)

func CreateJobs(testLogger zerolog.Logger, input cretypes.CreateJobsInput) error {
	if err := input.Validate(); err != nil {
		return errors.Wrap(err, "input validation failed")
	}

	for _, don := range input.DonTopology.DonsWithMetadata {
		if jobSpecs, ok := input.DonToJobSpecs[don.ID]; ok {
			createErr := jobs.Create(input.CldEnv.Offchain, don.DON, don.Flags, jobSpecs)
			if createErr != nil {
				return errors.Wrapf(createErr, "failed to create jobs for DON %d", don.ID)
			}
		} else {
			testLogger.Warn().Msgf("No job specs found for DON %d", don.ID)
		}
	}

	return nil
}

func BuildTopology(nodeSetInput []*cretypes.CapabilitiesAwareNodeSet) (*cretypes.Topology, error) {
	topology := &cretypes.Topology{}
	donsWithMetadata := make([]*cretypes.DonMetadata, len(nodeSetInput))

	// one DON to do everything
	if len(nodeSetInput) == 1 {
		flags, err := flags.NodeSetFlags(nodeSetInput[0])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get flags for nodeset %s", nodeSetInput[0].Name)
		}

		donsWithMetadata[0] = &cretypes.DonMetadata{
			ID:            1,
			Flags:         flags,
			NodesMetadata: make([]*cretypes.NodeMetadata, len(nodeSetInput[0].NodeSpecs)),
			Name:          nodeSetInput[0].Name,
		}
	} else {
		for i := range nodeSetInput {
			flags, err := flags.NodeSetFlags(nodeSetInput[i])
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get flags for nodeset %s", nodeSetInput[i].Name)
			}

			donsWithMetadata[i] = &cretypes.DonMetadata{
				ID:            libc.MustSafeUint32(i + 1),
				Flags:         flags,
				NodesMetadata: make([]*cretypes.NodeMetadata, len(nodeSetInput[i].NodeSpecs)),
				Name:          nodeSetInput[i].Name,
			}
		}
	}

	for i, donMetadata := range donsWithMetadata {
		for j := range donMetadata.NodesMetadata {
			nodeWithLabels := cretypes.NodeMetadata{}
			nodeType := cretypes.WorkerNode
			if nodeSetInput[i].BootstrapNodeIndex != -1 && j == nodeSetInput[i].BootstrapNodeIndex {
				nodeType = cretypes.BootstrapNode
			}
			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cretypes.Label{
				Key:   node.NodeTypeKey,
				Value: nodeType,
			})

			// TODO this will only work with Docker, for CRIB we need a different approach
			// that will need to be aware of namespace name and node naming pattern
			host := fmt.Sprintf("%s-node%d", donMetadata.Name, j)

			if nodeSetInput[i].GatewayNodeIndex != -1 && j == nodeSetInput[i].GatewayNodeIndex {
				nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cretypes.Label{
					Key:   node.ExtraRolesKey,
					Value: cretypes.GatewayNode,
				})

				topology.GatewayConnectorOutput = &cretypes.GatewayConnectorOutput{
					Path: "/node",
					Port: 5003,
					Host: host,
					// do not set gateway connector dons, they will be resolved automatically
				}
			}

			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cretypes.Label{
				Key:   node.IndexKey,
				Value: strconv.Itoa(j),
			})

			nodeWithLabels.Labels = append(nodeWithLabels.Labels, &cretypes.Label{
				Key:   node.HostLabelKey,
				Value: host,
			})

			donsWithMetadata[i].NodesMetadata[j] = &nodeWithLabels
		}
	}

	maybeID, err := flags.OneDonMetadataWithFlag(donsWithMetadata, cretypes.WorkflowDON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workflow DON ID")
	}

	topology.DonsMetadata = donsWithMetadata
	topology.WorkflowDONID = maybeID.ID

	return topology, nil
}

func AddKeysToTopology(topology *cretypes.Topology, keys *cretypes.GenerateKeysOutput) (*cretypes.Topology, error) {
	if topology == nil {
		return nil, errors.New("topology is nil")
	}

	if keys == nil {
		return nil, errors.New("keys is nil")
	}

	for _, donMetadata := range topology.DonsMetadata {
		if p2pKeys, ok := keys.P2PKeys[donMetadata.ID]; ok {
			for idx, nodeMetadata := range donMetadata.NodesMetadata {
				nodeMetadata.Labels = append(nodeMetadata.Labels, &cretypes.Label{
					Key:   node.NodeP2PIDKey,
					Value: p2pKeys.PeerIDs[idx],
				})
			}
		}

		if evmKeys, ok := keys.EVMKeys[donMetadata.ID]; ok {
			for idx, nodeMetadata := range donMetadata.NodesMetadata {
				nodeMetadata.Labels = append(nodeMetadata.Labels, &cretypes.Label{
					Key:   node.EthAddressKey,
					Value: evmKeys.PublicAddresses[idx].Hex(),
				})
			}
		}
	}

	return topology, nil
}

func GenereteKeys(input *cretypes.GenerateKeysInput) (*cretypes.GenerateKeysOutput, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	output := &cretypes.GenerateKeysOutput{
		EVMKeys: make(cretypes.DonsToEVMKeys),
		P2PKeys: make(cretypes.DonsToP2PKeys),
	}

	for _, donMetadata := range input.Topology.DonsMetadata {
		if input.GenerateP2PKeys {
			p2pKeys, err := crypto.GenerateP2PKeys(input.Password, len(donMetadata.NodesMetadata))
			if err != nil {
				return nil, errors.Wrap(err, "failed to generate P2P keys")
			}
			output.P2PKeys[donMetadata.ID] = p2pKeys
		}

		if len(input.GenerateEVMKeysForChainIDs) > 0 {
			evmKeys, err := crypto.GenerateEVMKeys(input.Password, len(donMetadata.NodesMetadata))
			if err != nil {
				return nil, errors.Wrap(err, "failed to generate EVM keys")
			}
			evmKeys.ChainIDs = append(evmKeys.ChainIDs, input.GenerateEVMKeysForChainIDs...)

			output.EVMKeys[donMetadata.ID] = evmKeys
		}
	}

	return output, nil
}

// In order to whitelist host IP in the gateway, we need to resolve the host.docker.internal to the host IP,
// and since CL image doesn't have dig or nslookup, we need to use curl.
func ResolveHostDockerInternaIP(testLogger zerolog.Logger, containerName string) (string, error) {
	if isCurlInstalled(containerName) {
		return resolveDockerHostWithCurl(containerName)
	} else if isNsLookupInstalled(containerName) {
		return resolveDockerHostWithNsLookup(containerName)
	}

	return "", errors.New("neither curl nor nslookup is installed")
}

func isNsLookupInstalled(containerName string) bool {
	cmd := []string{"which", "nslookup"}
	output, err := framework.ExecContainer(containerName, cmd)

	if err != nil || output == "" {
		return false
	}

	return true
}

func resolveDockerHostWithNsLookup(containerName string) (string, error) {
	cmd := []string{"nslookup", "host.docker.internal"}
	output, err := framework.ExecContainer(containerName, cmd)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`host.docker.internal(\n|\r)Address:\s+([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", errors.New("failed to extract IP address from curl output")
	}

	return matches[2], nil
}

func isCurlInstalled(containerName string) bool {
	cmd := []string{"which", "curl"}
	output, err := framework.ExecContainer(containerName, cmd)

	if err != nil || output == "" {
		return false
	}

	return true
}

func resolveDockerHostWithCurl(containerName string) (string, error) {
	cmd := []string{"curl", "-v", "http://host.docker.internal"}
	output, err := framework.ExecContainer(containerName, cmd)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`.*Trying ([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+).*`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", errors.New("failed to extract IP address from curl output")
	}

	return matches[1], nil
}
