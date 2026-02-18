package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

// MigrateAddressRefsToDefaultQualifierInput is the input for the migration pipeline.
// ChainSelectors lists the chain selectors whose MCMS address refs should get qualifier "default".
// Run this once per environment (e.g. prod_mainnet) to add refs with qualifier "default" so
// FindQualifierWithFullMCMSSet and batch_native/set_whitelist/transfer_erc20 lookups work.
type MigrateAddressRefsToDefaultQualifierInput struct {
	ChainSelectors []uint64 `json:"chain_selectors"`
}

// MigrateAddressRefsToDefaultQualifier reads the current datastore and outputs a new DataStore
// containing one ref per MCMS contract (Bypasser, Canceller, Proposer, CallProxy, RBACTimelock)
// on the given chains with qualifier DefaultTimelockQualifier ("default"). When the framework
// merges this into the environment, lookups by qualifier "default" will find these refs.
// Run once per environment; safe to run multiple times (adds refs, may duplicate if already migrated).
func MigrateAddressRefsToDefaultQualifier(env cldf.Environment, input MigrateAddressRefsToDefaultQualifierInput) (cldf.ChangesetOutput, error) {
	if env.DataStore == nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("datastore is required for migration")
	}
	if len(input.ChainSelectors) == 0 {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain_selectors must be non-empty")
	}

	store := env.DataStore.Addresses()
	ds := datastore.NewMemoryDataStore()
	mcmsTypes := map[datastore.ContractType]bool{
		datastore.ContractType(types.BypasserManyChainMultisig):  true,
		datastore.ContractType(types.CancellerManyChainMultisig): true,
		datastore.ContractType(types.ProposerManyChainMultisig):  true,
		datastore.ContractType(types.CallProxy):                  true,
		datastore.ContractType(types.RBACTimelock):               true,
	}

	var added int
	for _, chainSelector := range input.ChainSelectors {
		refs := store.Filter(datastore.AddressRefByChainSelector(chainSelector))
		for _, ref := range refs {
			if !mcmsTypes[datastore.ContractType(ref.Type)] {
				continue
			}
			migrated := ref
			migrated.Qualifier = DefaultTimelockQualifier
			if err := ds.Addresses().Add(migrated); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("chain %d: add migrated ref: %w", chainSelector, err)
			}
			added++
		}
	}

	if added == 0 {
		return cldf.ChangesetOutput{}, fmt.Errorf("no MCMS address refs found for the given chain_selectors; nothing to migrate")
	}

	_ = state.RequiredMCMSContractTypes // keep state package in use so FindQualifierWithFullMCMSSet stays consistent
	return cldf.ChangesetOutput{DataStore: ds}, nil
}
