package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/shared"
)

var _ deployment.ChangeSet[types.ImportToAddressbookConfig] = ImportToAddressbookChangeset

type AddressbookSchema struct {
	Address        string                    `json:"address"`
	TypeAndVersion deployment.TypeAndVersion `json:"typeAndVersion"`
	Label          string                    `json:"label"`
}

func ImportToAddressbookChangeset(env deployment.Environment, c types.ImportToAddressbookConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()

	addresses, err := shared.LoadJSON[[]*AddressbookSchema](c.InputFileName, c.InputFS)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load addresses input file: %w", err)
	}

	for _, address := range addresses {
		address.TypeAndVersion.AddLabel(address.Label)
		err = ab.Save(
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
