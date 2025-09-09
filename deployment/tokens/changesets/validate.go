package changesets

import (
	"fmt"
	"maps"
	"slices"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/ops"
)

// validateNoDupeSelectors checks that there are no duplicate chain selectors in
// the slice.
func validateNoDupeSelectors(csels []uint64) error {
	// Validate there are no duplicate chain selectors
	seen := make(map[uint64]bool)
	for _, csel := range csels {
		if seen[csel] {
			return fmt.Errorf("duplicate chain selector found: %d", csel)
		}
		seen[csel] = true
	}

	return nil
}

// validateNoExistingLinkToken checks that there are no existing Link Token
// contracts in the address book or datastore for the given chain selectors.
func validateNoExistingLinkToken(
	csels []uint64,
	qualifier string,
	addrBook cldf.AddressBook,
	ds datastore.DataStore,
) error {
	// Validate there is no existing Link Token contract in the address book
	for _, csel := range csels {
		addrBookByChain, err := addrBook.AddressesForChain(csel)
		if err != nil {
			// If the chain selector is not found in the address book, it means that
			// no contract record exists. This is not an error and we can continue.
			continue
		}

		if incl := slices.ContainsFunc(
			slices.Collect(maps.Values(addrBookByChain)),
			func(tv cldf.TypeAndVersion) bool {
				return tv.Equal(ops.LinkTokenTypeAndVersion1)
			},
		); incl {
			return fmt.Errorf("link token contract already exists for chain selector %d in address book", csel)
		}
	}

	// Validate there is no existing token contract in the datastore
	for _, csel := range csels {
		if _, err := ds.Addresses().Get(datastore.NewAddressRefKey(
			csel,
			datastore.ContractType(ops.LinkTokenTypeAndVersion1.Type.String()),
			&ops.LinkTokenTypeAndVersion1.Version,
			qualifier,
		)); err == nil {
			return fmt.Errorf("link token contract already exists for chain selector %d in datastore", csel)
		}
	}

	return nil
}

// validateChainSelectorFamily checks that all chain selectors are in the provided family.
func validateChainSelectorsFamily(csels []uint64, fam string) error {
	for _, csel := range csels {
		cselFam, err := chain_selectors.GetSelectorFamily(csel)
		if err != nil {
			return fmt.Errorf("failed to get family for chain selector %d: %w", csel, err)
		}

		if cselFam != fam {
			return fmt.Errorf("chain selector %d is not in the %s family", csel, fam)
		}
	}

	return nil
}

// validateSelectorsInEnvironment checks that all chain selectors are valid and
// available in the environment's blockchains.
func validateSelectorsInEnvironment(
	blockchains cldf_chain.BlockChains, selectors []uint64,
) error {
	for _, sel := range selectors {
		if !blockchains.Exists(sel) {
			return fmt.Errorf("chain %d not found in environment", sel)
		}
	}

	return nil
}
