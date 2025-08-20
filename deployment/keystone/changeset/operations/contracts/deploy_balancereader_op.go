package contracts

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type DeployBalanceReaderOpDeps struct {
	Env *cldf.Environment
}

type DeployBalanceReaderOpInput struct {
	ChainSelector uint64
}

type DeployBalanceReaderOpOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook cldf.AddressBook
}

// DeployBalanceReaderOp is an operation that deploys the BalanceReader contract.
var DeployBalanceReaderOp = operations.NewOperation(
	"deploy-balance-reader-op",
	semver.MustParse("1.0.0"),
	"Deploy BalanceReader Contract",
	func(b operations.Bundle, deps DeployBalanceReaderOpDeps, input DeployBalanceReaderOpInput) (DeployBalanceReaderOpOutput, error) {
		balanceReaderOutput, err := changeset.DeployBalanceReaderV2(*deps.Env, &changeset.DeployRequestV2{ChainSel: input.ChainSelector})
		if err != nil {
			return DeployBalanceReaderOpOutput{}, err
		}

		return DeployBalanceReaderOpOutput{
			Addresses:   balanceReaderOutput.DataStore.Addresses(),
			AddressBook: balanceReaderOutput.AddressBook, //nolint:staticcheck // keeping the address book since not everything has been migrated to datastore
		}, nil
	},
)
