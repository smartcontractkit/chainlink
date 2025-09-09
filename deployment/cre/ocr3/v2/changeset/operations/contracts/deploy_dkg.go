package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	dkg "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/dkg"
)

type DeployDKGDeps struct {
	Env *cldf.Environment
}

type DeployDKGInput struct {
	ChainSelector uint64
	Qualifier     string
}

type DeployDKGOutput struct {
	Address       string
	ChainSelector uint64
	Qualifier     string
	Type          string
	Version       string
	Labels        []string
	Datastore     datastore.DataStore
	AddressBook   cldf.AddressBook // backward compatibility, to be removed in CRE-742
}

// DeployDKG is an operation that deploys the DKG contract.
// This atomic operation performs the single side effect of deploying and registering the contract.
var DeployDKG = operations.NewOperation[DeployDKGInput, DeployDKGOutput, DeployDKGDeps](
	"deploy-dkg-op",
	semver.MustParse("1.0.0"),
	"Deploy DKG Contract",
	func(b operations.Bundle, deps DeployDKGDeps, input DeployDKGInput) (DeployDKGOutput, error) {
		lggr := deps.Env.Logger

		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return DeployDKGOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Deploy the DKG contract
		dkgAddr, tx, dkg, err := dkg.DeployDKG(
			chain.DeployerKey,
			chain.Client,
		)
		if err != nil {
			return DeployDKGOutput{}, fmt.Errorf("failed to deploy DKG: %w", err)
		}

		// Wait for deployment confirmation
		_, err = chain.Confirm(tx)
		if err != nil {
			return DeployDKGOutput{}, fmt.Errorf("failed to confirm DKG deployment: %w", err)
		}

		// Get type and version from the deployed contract
		tvStr, err := dkg.TypeAndVersion(&bind.CallOpts{})
		if err != nil {
			return DeployDKGOutput{}, fmt.Errorf("failed to get type and version: %w", err)
		}

		tv, err := cldf.TypeAndVersionFromString(tvStr)
		if err != nil {
			return DeployDKGOutput{}, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
		}

		// Create labels from the operation output
		labels := datastore.NewLabelSet()
		for _, label := range tv.Labels.List() {
			labels.Add(label)
		}

		addressRef := datastore.AddressRef{
			ChainSelector: chain.Selector,
			Address:       dkgAddr.Hex(),
			Type:          datastore.ContractType(tv.Type),
			Version:       &tv.Version,
			Labels:        labels,
			Qualifier:     input.Qualifier,
		}

		// Create a mutable datastore in order to be able to add the ocr3 address and access it from the configure step
		ds := datastore.NewMemoryDataStore()
		err = ds.Merge(deps.Env.DataStore)
		if err != nil {
			return DeployDKGOutput{}, fmt.Errorf("failed to merge datastore: %w", err)
		}

		if err := ds.AddressRefStore.Add(addressRef); err != nil {
			return DeployDKGOutput{}, fmt.Errorf("failed to add DKG address %v to datastore: %w", addressRef, err)
		}

		lggr.Infof("Deployed %s on chain selector %d at address %s", tv.String(), chain.Selector, dkgAddr.String())

		ab := cldf.NewMemoryAddressBook()
		err = ab.Save(chain.Selector, dkgAddr.String(), tv)
		if err != nil {
			return DeployDKGOutput{}, fmt.Errorf("failed to save address to address book: %w", err)
		}
		return DeployDKGOutput{
			Address:       dkgAddr.String(),
			ChainSelector: input.ChainSelector,
			Qualifier:     input.Qualifier,
			Type:          string(tv.Type),
			Version:       tv.Version.String(),
			Labels:        tv.Labels.List(),
			Datastore:     ds.Seal(),
			AddressBook:   ab, // TODO: CRE-742 remove AddressBook
		}, nil
	},
)
