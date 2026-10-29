package shared

import (
	"fmt"
	"sort"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	deployutil "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
)

// DefaultMCMSQualifier qualifies the default CCIP MCMS bundle; further bundles (RMN,
// ultra-fast-curse, operator) resolve by their own qualifier.
const DefaultMCMSQualifier = deployutil.CLLQualifier

// MCMSBundleRefs selects the active MCMS bundle rows for one chain and qualifier.
// For the default qualifier only, an empty-qualified bundle is accepted as a legacy
// fallback: singletons (and MCMS deployed by DeployMCMSWithTimelockV2 in test
// environments) are written with the empty qualifier. A non-empty selection under the
// requested qualifier always wins; the fallback never masks a qualifier mismatch for
// dedicated bundles (RMN, ultra-fast-curse, operator).
func MCMSBundleRefs(refs []datastore.AddressRef, selector uint64, qualifier string) ([]datastore.AddressRef, error) {
	active := make([]datastore.AddressRef, 0, len(refs))
	for _, ref := range refs {
		if ref.ChainSelector == selector && ref.Version != nil && !ref.Labels.Contains(SupersededLabel) {
			active = append(active, ref)
		}
	}
	qualifiers := []string{qualifier}
	if qualifier == DefaultMCMSQualifier {
		qualifiers = append(qualifiers, "")
	}
	for _, q := range qualifiers {
		var bundle []datastore.AddressRef
		for _, ref := range active {
			if ref.Qualifier == q {
				bundle = append(bundle, ref)
			}
		}
		if len(bundle) > 0 {
			sort.Slice(bundle, func(i, j int) bool {
				a, b := bundle[i], bundle[j]
				if a.Type != b.Type {
					return a.Type < b.Type
				}
				if a.Qualifier != b.Qualifier {
					return a.Qualifier < b.Qualifier
				}
				return a.Address < b.Address
			})
			if err := CheckRefUniqueness(bundle); err != nil {
				return nil, err
			}
			return bundle, nil
		}
	}
	if qualifier != DefaultMCMSQualifier {
		return nil, fmt.Errorf("no mcms refs for chain %d with qualifier %q", selector, qualifier)
	}
	return nil, nil
}

// SupersededLabel marks rows an operator declared historical; they never load as active contracts.
const SupersededLabel = "superseded"

// CheckRefUniqueness rejects two active refs claiming one (type, version, qualifier) identity
// with different addresses.
func CheckRefUniqueness(refs []datastore.AddressRef) error {
	seen := make(map[string]string)
	for _, ref := range refs {
		key := fmt.Sprintf("%s %s %s", ref.Type, ref.Version, ref.Qualifier)
		if prev, dup := seen[key]; dup && prev != ref.Address {
			return fmt.Errorf("%s %s qualified %q resolves to both %s and %s", ref.Type, ref.Version, ref.Qualifier, prev, ref.Address)
		}
		seen[key] = ref.Address
	}
	return nil
}
