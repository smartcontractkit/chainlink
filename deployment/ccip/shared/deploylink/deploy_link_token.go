package deploylink

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	linkchangesets "github.com/smartcontractkit/cld-changesets/tokens/link/changesets"
)

// DeployLinkTokenChangeset wraps the upstream DeployLinkTokenChangeset and
// additionally writes deployed addresses to AddressBook for backward compatibility.
type DeployLinkTokenChangeset struct{} //nolint:revive // intentional name match with upstream

var _ cldf.ChangeSetV2[linkchangesets.DeployLinkTokenInput] = DeployLinkTokenChangeset{}

func (DeployLinkTokenChangeset) VerifyPreconditions(e cldf.Environment, input linkchangesets.DeployLinkTokenInput) error {
	return (linkchangesets.DeployLinkTokenChangeset{}).VerifyPreconditions(e, input)
}

func (DeployLinkTokenChangeset) Apply(e cldf.Environment, input linkchangesets.DeployLinkTokenInput) (cldf.ChangesetOutput, error) {
	out, err := (linkchangesets.DeployLinkTokenChangeset{}).Apply(e, input)
	if err != nil {
		return out, err
	}

	if out.DataStore != nil {
		ab := cldf.NewMemoryAddressBook()
		refs, fetchErr := out.DataStore.Addresses().Fetch()
		if fetchErr != nil {
			return out, fmt.Errorf("failed to fetch addresses from datastore: %w", fetchErr)
		}
		for _, ref := range refs {
			if ref.Version == nil {
				continue
			}
			tv := cldf.NewTypeAndVersion(cldf.ContractType(ref.Type), *ref.Version)
			if addErr := ab.Save(ref.ChainSelector, ref.Address, tv); addErr != nil {
				return out, fmt.Errorf("failed to save address to address book: %w", addErr)
			}
		}
		out.AddressBook = ab //nolint:staticcheck // intentional use of deprecated AddressBook for backward compat
	}

	return out, nil
}
