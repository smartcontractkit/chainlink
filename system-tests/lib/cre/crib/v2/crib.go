package cribv2

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/crib-sdk/crib"
	anvilv1 "github.com/smartcontractkit/crib-sdk/crib/composite/anvil/v1"
)

func DeployBlockchain(input *types.DeployCribBlockchainInput) (*blockchain.Output, error) {
	ctx := context.Background()
	cribNamespace := "crib-ns-tbd"

	plan := crib.NewPlan(
		"anvilv1",
		crib.Namespace(cribNamespace),
		crib.ComponentSet(
			anvilv1.Component(&anvilv1.Props{
				Namespace: cribNamespace,
				ChainID:   "",
			}),
		),
	)

	if err := plan.Apply(ctx); err != nil {
		fmt.Printf("Error applying plan: %v\n", err)
	}
	fmt.Printf("Successfully applied plan: %s\n", plan.Name())

	// todo, and now how to get access to the Component props?

}
