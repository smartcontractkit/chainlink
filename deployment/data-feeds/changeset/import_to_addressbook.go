package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/shared"
)

// ImportToAddressbookChangeset is a changeset that reads already deployed contract addresses from input file
// and saves them to the address book. Returns a new addressbook with the imported addresses.
var ImportToAddressbookChangeset = deployment.CreateChangeSet(importToAddressbookLogic, importToAddressbookPrecondition)

type AddressesSchema struct {
	Address        string                    `json:"address"`
	TypeAndVersion deployment.TypeAndVersion `json:"typeAndVersion"`
	Label          string                    `json:"label"`
}

func importToAddressbookLogic(env deployment.Environment, c types.ImportToAddressbookConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()

	addresses, _ := shared.LoadJSON[[]*AddressesSchema](c.InputFileName, c.InputFS)

	for _, address := range addresses {
		address.TypeAndVersion.AddLabel(address.Label)
		err := ab.Save(
			c.ChainSelector,
			address.Address,
			address.TypeAndVersion,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to save address %s: %w", address.Address, err)
		}
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

func importToAddressbookPrecondition(env deployment.Environment, c types.ImportToAddressbookConfig) error {
	_, ok := env.Chains[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if c.InputFileName == "" {
		return errors.New("input file name is required")
	}

	_, err := shared.LoadJSON[[]*AddressesSchema](c.InputFileName, c.InputFS)
	if err != nil {
		return fmt.Errorf("failed to load addresses input file: %w", err)
	}

	return nil
}
