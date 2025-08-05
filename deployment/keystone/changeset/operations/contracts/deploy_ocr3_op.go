package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type DeployOCR3OpDeps struct {
	Env *cldf.Environment
}

type DeployOCR3OpInput struct {
	ChainSelector uint64
	Qualifier     string
}

type DeployOCR3OpOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook cldf.AddressBook
}

// DeployOCR3Op is an operation that deploys the OCR3 contract.
var DeployOCR3Op = operations.NewOperation[DeployOCR3OpInput, DeployOCR3OpOutput, DeployOCR3OpDeps](
	"deploy-ocr3-op",
	semver.MustParse("1.0.0"),
	"Deploy OCR3 Contract",
	func(b operations.Bundle, deps DeployOCR3OpDeps, input DeployOCR3OpInput) (DeployOCR3OpOutput, error) {
		ocr3Output, err := changeset.DeployOCR3V2(*deps.Env, &changeset.DeployRequestV2{
			ChainSel:  input.ChainSelector,
			Qualifier: input.Qualifier,
		})
		if err != nil {
			return DeployOCR3OpOutput{}, fmt.Errorf("DeployOCR3Op error: failed to deploy OCR3 contract: %w", err)
		}

		return DeployOCR3OpOutput{
			Addresses: ocr3Output.DataStore.Addresses(), AddressBook: ocr3Output.AddressBook, //nolint:staticcheck // keeping the address book since not everything has been migrated to datastore
		}, nil
	},
)
