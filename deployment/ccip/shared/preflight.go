package shared

import (
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// AddressRefConflict is two planned address refs that resolve to the same datastore key.
type AddressRefConflict struct {
	// Key is the datastore key both refs resolve to.
	Key string
	// First and Second are the two refs claiming it.
	First  datastore.AddressRef
	Second datastore.AddressRef
}

// AddressRefValidationError reports planned address refs that cannot be written to a datastore.
// Returning it from VerifyPreconditions stops the changeset before it deploys anything.
type AddressRefValidationError struct {
	// Conflicts are pairs of planned refs sharing a key. Only one of the two addresses could
	// survive being written, so the contracts need distinct qualifiers.
	Conflicts []AddressRefConflict
	// Versionless are planned refs with no version. A datastore key is built from the version,
	// so such a ref has no key.
	Versionless []datastore.AddressRef
}

// Error summarizes the problem. The individual refs are on the struct, for the caller to render.
func (e *AddressRefValidationError) Error() string {
	parts := make([]string, 0, 2)
	if n := len(e.Conflicts); n > 0 {
		parts = append(parts, fmt.Sprintf("%d pairs of contracts share a datastore key", n))
	}
	if n := len(e.Versionless); n > 0 {
		parts = append(parts, fmt.Sprintf("%d refs have no version", n))
	}

	return "address refs cannot be written to the datastore: " + strings.Join(parts, ", ")
}

// ValidateAddressRefs checks that planned address refs can be written to env's datastore, and
// is meant to be called from a changeset's VerifyPreconditions with the refs the changeset is
// about to produce.
//
// The datastore keys contracts on (chain, type, version, qualifier), so two contracts of one
// type and version on one chain need distinct qualifiers. Nothing in the framework can work
// out those qualifiers for a changeset, and nothing else discovers the clash until the
// changeset has already deployed: the refs are only written once Apply returns. A changeset
// deploying more than one instance of a type should therefore declare its refs here first, so
// a missing qualifier costs a failed precondition rather than orphaned contracts.
//
// A planned ref taking over a key that env's datastore already holds is not an error. The
// environment is the state before this run, so replacing it is what a redeploy means; it is
// logged so the replacement is visible. Use ValidateAddressRefsStrict to reject it instead.
func ValidateAddressRefs(env deployment.Environment, planned []datastore.AddressRef) error {
	return validatePlannedAddressRefs(env, planned, false)
}

// ValidateAddressRefsStrict is ValidateAddressRefs, and additionally rejects a planned ref that
// would replace a different address already held under the same key in env's datastore. Use it
// for a changeset that is only ever meant to deploy contracts that do not exist yet.
func ValidateAddressRefsStrict(env deployment.Environment, planned []datastore.AddressRef) error {
	return validatePlannedAddressRefs(env, planned, true)
}

// sameContract reports whether two refs sharing a key are one contract declared twice rather
// than two contracts competing for the key.
//
// Only a known address can establish that. A changeset calling this from VerifyPreconditions has
// not deployed yet, so its refs carry no address; two of them are then indistinguishable, and
// treating them as one contract would let exactly the case this check exists for — two instances
// of a type with no distinguishing qualifier — pass validation and fail after the deploy.
func sameContract(a, b datastore.AddressRef) bool {
	return a.Address != "" && AddressesEqual(a.ChainSelector, a.Address, b.Address)
}

func validatePlannedAddressRefs(env deployment.Environment, planned []datastore.AddressRef, strict bool) error {
	problems := &AddressRefValidationError{}

	claimed := make(map[string]datastore.AddressRef, len(planned))
	for _, ref := range planned {
		if ref.Version == nil {
			problems.Versionless = append(problems.Versionless, ref)
			continue
		}

		key := ref.Key().String()
		if first, taken := claimed[key]; taken && !sameContract(first, ref) {
			problems.Conflicts = append(problems.Conflicts, AddressRefConflict{
				Key: key, First: first, Second: ref,
			})

			continue
		}

		claimed[key] = ref
	}

	if len(problems.Conflicts) > 0 || len(problems.Versionless) > 0 {
		return fmt.Errorf("%w: %w", deployment.ErrInvalidConfig, problems)
	}

	if isNilStore(env.DataStore) {
		return nil
	}

	envRefs, err := env.DataStore.Addresses().Fetch()
	if err != nil {
		return fmt.Errorf("failed to read environment data store: %w", err)
	}

	held := make(map[string]datastore.AddressRef, len(envRefs))
	for _, ref := range envRefs {
		if ref.Version == nil {
			continue
		}
		held[ref.Key().String()] = ref
	}

	for key, ref := range claimed {
		existing, taken := held[key]
		if !taken || sameContract(existing, ref) {
			continue
		}

		if strict {
			return fmt.Errorf("%w: %s already holds %s, and the changeset would replace it with %s",
				deployment.ErrInvalidEnvironment, key, existing.Address, ref.Address)
		}

		if env.Logger != nil {
			env.Logger.Infow("Changeset will supersede an existing address ref",
				"key", key, "previousAddress", existing.Address, "newAddress", ref.Address)
		}
	}

	return nil
}

// AsAddressRefValidationError extracts the detail from an error returned by ValidateAddressRefs.
func AsAddressRefValidationError(err error) (*AddressRefValidationError, bool) {
	var target *AddressRefValidationError
	ok := errors.As(err, &target)

	return target, ok
}
