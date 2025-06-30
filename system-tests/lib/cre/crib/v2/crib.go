package cribv2

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/crib-sdk/crib"
	anvilv1 "github.com/smartcontractkit/crib-sdk/crib/composite/anvil/v1"
)

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

	anvilComponents := result.ComponentByName("sdk.AnvilCompositeV1")

	for component := range anvilComponents {
		res := crib.ComponentState[*anvilv1.Result](component)
		fmt.Printf("The args used: %v\n", res.InternalWSUrl())

		return &blockchain.Output{
			Type:    input.BlockchainInput.Type,
			Family:  "evm",
			ChainID: input.BlockchainInput.ChainID,
			Nodes: []*blockchain.Node{
				{
					ExternalWSUrl:   "todo, requires telepresence",
					ExternalHTTPUrl: "todo, requires telepresence",
					InternalWSUrl:   res.InternalWSUrl(),
					InternalHTTPUrl: "do we need that?",
				},
			},
		}, nil
	}

	return nil, errors.New("failed to find a valid component")

}

func DeployDons(input *types.DeployCribDonsInput) ([]*types.CapabilitiesAwareNodeSet, error) {
	//ctx := context.Background()
	//
	//nodeset := anvilv1.Component(&nodesetv1.Props{
	//	Namespace: input.Namespace,
	//	// Requirements:
	//	// - image with tag
	//	// - config override files for bt (0-n) and worker nodes (0-n)
	//	// WORKFLOW_DON and GATEWAY_DON
	//	// 3 topologie,
	//	// 1) 1 Don, with one bootstrap which is also a gateway node
	//	//    bootstrap and worker nodes are variable
	//	//  note: impossible to do because of DON_TYPE hardcoded things in the chart
	//	//  this the highest value, as currently it is not possible at all, it is the fastest to deploy and provide a feedback loop to dev
	//	// 2) workflow don
	//	//   separate don which is just a gateway
	//	// 3) same as 2) plus capability DON
	//
	//	// Copying binaries to pod, handle it later, in sandbox there is an image which has all plugins
	//	// it would be nice for feature parity with docker
	//	// low value, high complexity
	//
	//})
	//plan := crib.NewPlan(
	//	"nodesetv1",
	//	crib.Namespace(input.Namespace),
	//	crib.ComponentSet(
	//		nodeset,
	//	),
	//)
	//
	//if err := plan.Apply(ctx); err != nil {
	//	fmt.Printf("Error applying plan: %v\n", err)
	//}
	//fmt.Printf("Successfully applied plan: %s\n", plan.Name())
	//
	//// todo, and now how to get access to the Component props?
	return make([]*types.CapabilitiesAwareNodeSet, 0), nil

}
