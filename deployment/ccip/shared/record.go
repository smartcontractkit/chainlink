package shared

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// DeployContractAndRecord is deployment.DeployContract, and additionally writes the datastore
// ref for the contract as part of the same step.
//
// Recording the ref here rather than deriving it from the address book afterwards is what stops a
// deployment succeeding and then failing to be recorded. A derived-at-the-end datastore is a
// second computation over the results, so it can disagree with what was deployed, and it can only
// fail once every contract already exists on chain at which point there is nothing useful left
// to do about it. Written here, the ref lands at the same moment the address does.
//
// It also means the store is filled in incrementally, so a run that dies partway through still
// has a datastore record of what it managed to deploy.
//
// The key is claimed BEFORE the transaction is sent, which is the whole point of the ordering
// here. The datastore key is (chain, type, version, qualifier) and none of that depends on the
// deployed address, so it can be established up front, because a conflict found
// afterwards would leave a contract deployed and written to the address book with no datastore
// ref, which is the inconsistency this helper exists to prevent. That is why tv and qualifier are
// parameters rather than being read back off the deploy closure's result: the closure performs the
// deployment, so anything it returns is only knowable too late.
//
// The qualifier says which instance of tv this is, since only the caller knows that. A chain
// singleton passes "", anything deployed more than once on a chain e.g. (token pool or a lane) needs
// a value, or the instances collide on one key.
//
// A caller that reserved its keys earlier (see ValidateAddressRefs) finds its own reservation here
// and takes it. A caller that did not claims the key now. Either way a second contract arriving at
// an occupied key fails before it is deployed, so nothing is ever stranded and the contract
// already recorded there is never silently dropped.
func DeployContractAndRecord[C any](
	lggr logger.Logger,
	chain cldf_evm.Chain,
	addressBook deployment.AddressBook,
	ds datastore.MutableDataStore,
	tv deployment.TypeAndVersion,
	qualifier string,
	deploy func(chain cldf_evm.Chain) deployment.ContractDeploy[C],
) (*deployment.ContractDeploy[C], error) {
	// A nil store is a programming error rather than a shorthand for "skip the datastore":
	// silently degrading to DeployContract would drop the ref this function exists to write.
	if isNilStore(ds) {
		return nil, errors.New("DeployContractAndRecord requires a datastore; use DeployContract to write only the address book")
	}

	version := tv.Version
	ref := datastore.AddressRef{
		ChainSelector: chain.Selector,
		Type:          datastore.ContractType(tv.Type),
		Version:       &version,
		Qualifier:     qualifier,
	}
	if !tv.Labels.IsEmpty() {
		ref.Labels = datastore.NewLabelSet(tv.Labels.List()...)
	}

	claimed, err := claimAddressRef(ds, ref)
	if err != nil {
		lggr.Errorw("Refusing to deploy: datastore key is already taken",
			"Contract", tv.String(), "chain", chain.String(), "err", err)

		return nil, err
	}

	// A failure from here on releases the reservation again but only one this call made: a
	// reservation found already in place belongs to whoever made it. A leftover empty-address
	// ref is not a record of progress (the address is the part that matters), and it would be
	// merged and persisted as if it were one.
	release := func() {
		if !claimed {
			return
		}
		if delErr := ds.Addresses().Delete(ref.Key()); delErr != nil {
			lggr.Warnw("Failed to release the datastore key reservation after a failed deployment",
				"Contract", tv.String(), "chain", chain.String(), "err", delErr)
		}
	}

	// The deploy-and-confirm sequence of DeployContract is inlined rather than delegated to,
	// because the declared type and version must be checked against what was deployed BEFORE
	// the address book is written: a mismatch found afterwards would leave the address saved
	// under the wrong type while the datastore key was reserved for the declared one.
	contractDeploy := deploy(chain)
	if contractDeploy.Err != nil {
		lggr.Errorw("Failed to deploy contract", "chain", chain.String(), "err", contractDeploy.Err)
		release()

		return nil, contractDeploy.Err
	}
	if !chain.IsZkSyncVM {
		if _, err = chain.Confirm(contractDeploy.Tx); err != nil {
			lggr.Errorw("Failed to confirm deployment", "chain", chain.String(), "Contract", contractDeploy.Tv.String(), "err", err)
			release()

			return nil, err
		}
	}
	lggr.Infow("Deployed contract", "Contract", contractDeploy.Tv.String(), "addr", contractDeploy.Address, "chain", chain.String())

	if contractDeploy.Tv.String() != tv.String() {
		release()

		return nil, fmt.Errorf(
			"deployed contract is %s but %s was declared to DeployContractAndRecord; the datastore ref was reserved for the declared type",
			contractDeploy.Tv.String(), tv.String())
	}

	if err = addressBook.Save(chain.Selector, contractDeploy.Address.String(), contractDeploy.Tv); err != nil {
		lggr.Errorw("Failed to save contract address", "Contract", contractDeploy.Tv.String(), "addr", contractDeploy.Address, "chain", chain.String(), "err", err)
		release()

		return nil, err
	}

	// The key is already held by this call, so filling in the address cannot conflict.
	ref.Address = contractDeploy.Address.String()
	if err := ds.Addresses().Upsert(ref); err != nil {
		lggr.Errorw("Failed to record contract address ref",
			"Contract", tv.String(), "addr", contractDeploy.Address,
			"chain", chain.String(), "err", err)
		release()

		return nil, err
	}

	return &contractDeploy, nil
}

// RecordAddress records a contract whose address the caller already has, into both the address
// book and the datastore, under a declared qualifier.
//
// It is the counterpart to DeployContractAndRecord for addresses that are not produced by an EVM
// deployment: a Solana PDA or a Sui object is derived rather than deployed, and an imported
// contract already exists. There is no transaction to order against here, so unlike the deploy
// path this cannot fail "too late" — but it enforces the same rule, that a key already holding a
// different contract (or an unowned empty reservation) is an error rather than an overwrite.
//
// Recording the same address twice is a no-op, so a re-run is safe. Addresses are compared with
// the case rules of the chain's family, so a checksummed EVM address and its lower-case form are
// recognised as one contract. A nil address book is an explicit opt-out of the legacy registry.
func RecordAddress(
	ab deployment.AddressBook,
	ds datastore.MutableDataStore,
	chainSelector uint64,
	address string,
	tv deployment.TypeAndVersion,
	qualifier string,
) error {
	if isNilStore(ds) {
		return errors.New("RecordAddress requires a datastore")
	}

	version := tv.Version
	ref := datastore.AddressRef{
		ChainSelector: chainSelector,
		Address:       address,
		Type:          datastore.ContractType(tv.Type),
		Version:       &version,
		Qualifier:     qualifier,
	}
	if !tv.Labels.IsEmpty() {
		ref.Labels = datastore.NewLabelSet(tv.Labels.List()...)
	}

	// The key check comes before anything is written: a key already holding a different
	// contract is an error, and raising it only after one registry was written would leave the
	// two disagreeing. An empty-address reservation is a pending claim in the caller's private
	// store (typically created by ReserveRefs), so it is this call's to fill.
	existing, getErr := ds.Addresses().Get(ref.Key())
	notFound := errors.Is(getErr, datastore.ErrAddressRefNotFound)
	switch {
	case notFound:
		// The key is free; it is added once the address book has accepted the address.
	case getErr != nil:
		return fmt.Errorf("failed to read address ref %s: %w", ref.Key().String(), getErr)

	case existing.Address == "" || AddressesEqual(chainSelector, existing.Address, address):
		// An empty address is a reservation this run made (via ReserveRefs) to fill; an equal
		// address is the same contract recorded twice. Either way the key is ours to write.

	default:
		// A different contract is genuinely occupied.
		return fmt.Errorf(
			"%w: %s already holds %s, so %s cannot be recorded under it; give the two contracts distinct qualifiers",
			datastore.ErrAddressRefExists, ref.Key().String(), existing.Address, address)
	}

	// The address book write comes first: it is the write that can still fail here, and with
	// the datastore untouched a failure leaves nothing to roll back.
	if !isNilStore(ab) {
		if err := saveAddressBookRecord(ab, chainSelector, address, tv); err != nil {
			return err
		}
	}

	if notFound {
		return ds.Addresses().Add(ref)
	}

	return ds.Addresses().Upsert(ref)
}

// saveAddressBookRecord writes the address to the address book unless it is already there under
// the same type and version, which is what makes a repeated RecordAddress a no-op rather than an
// error: AddressBookMap.Save rejects a repeated address outright. The same address recorded
// under a different type and version is a genuine disagreement, and is rejected before anything
// is written.
func saveAddressBookRecord(ab deployment.AddressBook, chainSelector uint64, address string, tv deployment.TypeAndVersion) error {
	chainAddrs, err := ab.AddressesForChain(chainSelector)
	if err != nil && !errors.Is(err, deployment.ErrChainNotFound) {
		return fmt.Errorf("failed to read address book for chain %d: %w", chainSelector, err)
	}

	for recorded, known := range chainAddrs {
		if !AddressesEqual(chainSelector, recorded, address) {
			continue
		}
		if !known.Equal(tv) {
			return fmt.Errorf(
				"the address book records %s on chain %d as %s, so it cannot be recorded as %s",
				recorded, chainSelector, known.String(), tv.String())
		}

		return nil
	}

	return ab.Save(chainSelector, address, tv)
}

// claimAddressRef takes ownership of ref's key ahead of the deployment. ref carries no address
// yet, by design: the claim has to happen before the contract exists, so the only thing that can
// be established here is the key. It reports whether it created the reservation, so the caller
// can remove it again if the deployment goes on to fail.
//
// A key that already holds a deployed address is occupied and a second claim is refused before
// the transaction is sent. An empty-address reservation is a pending claim: either one this call
// just made, or one made earlier in the same run via ReserveRefs, which validates the whole key
// set up front (no duplicate keys, every multi-instance member qualified) and hands the caller a
// private store. Both are the caller's to fill.
func claimAddressRef(ds datastore.MutableDataStore, ref datastore.AddressRef) (bool, error) {
	existing, err := ds.Addresses().Get(ref.Key())
	switch {
	case errors.Is(err, datastore.ErrAddressRefNotFound):
		if addErr := ds.Addresses().Add(ref); addErr != nil {
			return false, addErr
		}

		return true, nil

	case err != nil:
		return false, fmt.Errorf("failed to read address ref %s: %w", ref.Key().String(), err)

	// An empty address is a reservation in this private store (possibly from ReserveRefs), which
	// is this caller's to fill rather than release.
	case existing.Address == "":
		return false, nil

	default:
		return false, fmt.Errorf(
			"%w: %s already holds %s; give the two contracts distinct qualifiers",
			datastore.ErrAddressRefExists, ref.Key().String(), existing.Address)
	}
}

// isNilStore reports whether v is unset, including the case of a non-nil interface holding a
// nil pointer. A changeset returning ChangesetOutput{DataStore: (*datastore.MemoryDataStore)(nil)}
// passes an ordinary != nil check and then panics on first use.
func isNilStore(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	return rv.Kind() == reflect.Pointer && rv.IsNil()
}
