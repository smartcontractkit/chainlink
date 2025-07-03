package crib

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	libtypes "github.com/smartcontractkit/chainlink/system-tests/lib/types"
	"github.com/smartcontractkit/crib-sdk/crib"
	anvilv1 "github.com/smartcontractkit/crib-sdk/crib/composite/anvil/v1"
	jdv1 "github.com/smartcontractkit/crib-sdk/crib/composite/chainlink/jd/v1"
	nodev1 "github.com/smartcontractkit/crib-sdk/crib/composite/chainlink/node/v1"
	namespacev1 "github.com/smartcontractkit/crib-sdk/crib/scalar/k8s/namespace/v1"
	"path/filepath"
)

func Bootstrap(infraInput *libtypes.InfraInput) error {
	plan := crib.NewPlan(
		"namespace",
		crib.Namespace(infraInput.CRIB.Namespace),
		crib.ComponentSet(
			namespacev1.Component(infraInput.CRIB.Namespace),
			// todo: add telepresence install here for now
		),
	)
	_, err := plan.Apply(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to apply plan")
	}

	return nil
}

func DeployBlockchain(input *types.DeployCribBlockchainInput) (*blockchain.Output, error) {
	err := input.Validate()
	if err != nil {
		return nil, errors.Wrapf(err, "invalid input for deploying blockchain")
	}

	ctx := context.Background()

	anvil := anvilv1.Component(&anvilv1.Props{
		Namespace: input.Namespace,
		ChainID:   input.BlockchainInput.ChainID,
	})

	plan := crib.NewPlan(
		"anvilv1",
		crib.Namespace(input.Namespace),
		crib.ComponentSet(
			anvil,
		),
	)

	result, err := plan.Apply(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply a plan")
	}

	anvilComponents := result.ComponentByName(anvilv1.ComponentName)

	for component := range anvilComponents {
		res := crib.ComponentState[anvilv1.Result](component)

		return &blockchain.Output{
			Type:    input.BlockchainInput.Type,
			Family:  "evm",
			ChainID: input.BlockchainInput.ChainID,
			Nodes: []*blockchain.Node{
				{
					InternalWSUrl:   res.RPCWebsocketURL(),
					ExternalWSUrl:   res.RPCWebsocketURL(),
					InternalHTTPUrl: res.RPCHTTPURL(),
					ExternalHTTPUrl: res.RPCHTTPURL(),
				},
			},
		}, nil
	}

	return nil, errors.New("failed to find a valid component")
}

func DeployJdWithCRIBSDK(input *types.DeployCribJdInput) (*jd.Output, error) {
	if input == nil {
		return nil, errors.New("DeployCribJdInput is nil")
	}

	if valErr := input.Validate(); valErr != nil {
		return nil, errors.Wrap(valErr, "input validation failed")
	}

	imgTagIndex, err := dockerImageTag(input.JDInput.Image)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get image tag")
	}

	jdv1.Component(&jdv1.Props{
		JD: jdv1.JDProps{
			Image:            input.JDInput.Image,
			CSAEncryptionKey: input.JDInput.CSAEncryptionKey,
		},
		WaitForRollout: true,
	})

	jdEnvVars := map[string]string{
		"JOB_DISTRIBUTOR_IMAGE_TAG": imgTagIndex,
	}
	_, err = input.NixShell.RunCommandWithEnvVars("devspace run deploy-jd --no-warn", jdEnvVars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run devspace run deploy-jd")
	}

	jdOut, err := infra.ReadJdURL(filepath.Join(".", input.CribConfigsDir))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read JD URL from file")
	}

	return jdOut, nil
}

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
