package v1_6

import (
	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipoperations "github.com/smartcontractkit/chainlink/deployment/ccip/operations"
)

var (
	DeployRMNRemoteOp = operations.NewOperation(
		"DeployRMNRemote",
		semver.MustParse("1.0.0"),
		"Deploys RMNRemote 1.6 contract on the specified evm chain(s)",
		func(b operations.Bundle, deps ccipoperations.OpDependencies, _ any) (deployment.AddressBook, error) {

		})
)
