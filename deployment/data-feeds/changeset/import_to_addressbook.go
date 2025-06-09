package changeset

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// ImportToAddressbookChangeset is a changeset that reads already deployed contract addresses from input file
// and saves them to the address book. Returns a new addressbook with the imported addresses.
var ImportToAddressbookChangeset = cldf.CreateChangeSet(importToAddressbookLogic, importToAddressbookPrecondition)

type AddressesSchema struct {
	Address        string              `json:"address"`
	TypeAndVersion cldf.TypeAndVersion `json:"typeAndVersion"`
	Label          string              `json:"label"`
}

func importToAddressbookLogic(env cldf.Environment, c types.ImportAddressesConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()

	addresses, _ := LoadJSON[[]*AddressesSchema](c.InputFileName, c.InputFS)

	for _, address := range addresses {
		address.TypeAndVersion.AddLabel(address.Label)
		err := ab.Save(
			c.ChainSelector,
			address.Address,
			address.TypeAndVersion,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save address %s: %w", address.Address, err)
		}
	}

	return cldf.ChangesetOutput{AddressBook: ab}, nil
}

func importToAddressbookPrecondition(env cldf.Environment, c types.ImportAddressesConfig) error {
	if !env.BlockChains.Exists(c.ChainSelector) {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if c.InputFileName == "" {
		return errors.New("input file name is required")
	}

	_, err := LoadJSON[[]*AddressesSchema](c.InputFileName, c.InputFS)
	if err != nil {
		return fmt.Errorf("failed to load addresses input file: %w", err)
	}

	return nil
}
