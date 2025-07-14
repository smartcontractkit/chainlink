package crib

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/crib-sdk/crib"
	nodev1 "github.com/smartcontractkit/crib-sdk/crib/composite/chainlink/node/v1"
)

// todo: after it's done it will replace crib.DeployDons
func DeployDonsWithCribSDK(input *types.DeployCribDonsInput) ([]*types.CapabilitiesAwareNodeSet, error) {
	if input == nil {
		return nil, errors.New("DeployCribDonsInput is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, errors.Wrap(valErr, "input validation failed")
	}

	componentFuncs := make([]crib.ComponentFunc, 0)

	for j, donMetadata := range input.Topology.DonsMetadata {

		imageName, imageTag, err := imageNameAndTag(input, j)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get image name and tag for %s", donMetadata.Name)
		}

		for i, nodeMetadata := range donMetadata.NodesMetadata {
			configToml, secrets, confSecretsErr := getConfigAndSecretsForNode(nodeMetadata, j, input, donMetadata)
			if confSecretsErr != nil {
				return nil, confSecretsErr
			}
			cFunc := nodev1.Component(&nodev1.Props{
				Namespace:       input.Namespace,
				Image:           fmt.Sprintf("%s:%s", imageName, imageTag),
				AppInstanceName: fmt.Sprintf("%s-%d", donMetadata.Name, i),
				// passing as config not as override
				Config: *configToml,
				SecretsOverrides: map[string]string{
					"overrides": *secrets,
				},
			})
			componentFuncs = append(componentFuncs, cFunc)
		}
	}

	plan := crib.NewPlan(
		"nodesets",
		crib.Namespace(input.Namespace),
		crib.ComponentSet(
			componentFuncs...,
		),
	)

	planState, err := plan.Apply(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply plan")
	}

	// Setting outputs based on the Plan Results
	nodeComponents := planState.ComponentByName(nodev1.ComponentName)

	var nodeResults []nodev1.Result

	for component := range nodeComponents {
		res := crib.ComponentState[nodev1.Result](component)
		nodeResults = append(nodeResults, res)
		fmt.Printf("Node API URL: %s\n", res.APIUrl())
		fmt.Printf("API Credentials: username: %s , password: %s\n", res.APICredentials.UserName, res.APICredentials.Password)
	}

	// setting outputs in a similar way as in func ReadNodeSetURL
	for j := range input.Topology.DonsMetadata {
		out := &ns.Output{
			// UseCache: true will disable deploying docker containers via CTF
			UseCache: true,
			CLNodes:  []*clnode.Output{},
		}
		// todo: for now this is hardcoded for a single don, we need to group results for each don
		for _, res := range nodeResults {
			out.CLNodes = append(out.CLNodes, &clnode.Output{
				// UseCache: true will disable deploying docker containers via CTF
				UseCache: true,
				Node: &clnode.NodeOut{
					APIAuthUser:     res.APICredentials.UserName,
					APIAuthPassword: res.APICredentials.Password,
					ExternalURL:     res.APIUrl(),
					InternalURL:     res.APIUrl(),
					// todo: this should be simplified in the CTF types, we should just pass P2P port
					InternalP2PUrl: fmt.Sprintf("http://%s:%d", res.HostName(), res.P2PPort),
					InternalIP:     res.HostName(),
				},
			})
		}
		input.NodeSetInputs[j].Out = out
	}

	return input.NodeSetInputs, nil
}
