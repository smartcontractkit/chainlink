package crib

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/crib-sdk/crib"
	nodev1 "github.com/smartcontractkit/crib-sdk/crib/composite/chainlink/node/v1"
)

// todo: after it's done it will replace DeployDonsCrib
func DeployDonsWithCribSDK(input *types.DeployCribDonsInput) ([]*types.CapabilitiesAwareNodeSet, error) {
	if input == nil {
		return nil, errors.New("DeployCribDonsInput is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, errors.Wrap(valErr, "input validation failed")
	}

	// component funcs with all nodes from all nodesets
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
			component := nodev1.Component(&nodev1.Props{
				Image:       fmt.Sprintf("%s:%s", imageName, imageTag),
				ReleaseName: fmt.Sprintf("%s-%d", donMetadata.Name, i),
				// passing as config not as override
				Config: *configToml,
				SecretsOverrides: map[string]string{
					"overrides": *secrets,
				},
			})
			componentFuncs = append(componentFuncs, component)
		}
	}

	plan := crib.NewPlan(
		"dons",
		crib.Namespace(input.Namespace),
		crib.ComponentSet(
			componentFuncs...,
		),
	)

	planState, err := plan.Apply(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply plan")
	}

	// todo: set outputs basd on the plan results
	nodeComponents := planState.ComponentByName(nodev1.ComponentName)

	for component := range nodeComponents {
		res := crib.ComponentState[nodev1.Result](component)
		fmt.Printf("result: %v\n", res)
	}

	return input.NodeSetInputs, nil
}
