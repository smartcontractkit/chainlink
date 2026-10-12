package shared

import (
	"errors"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

const recordTestChainSelector = 5009297550715157269

// recordTestRef builds a ref as claimAddressRef sees one: keyed, but with no address, because the
// claim happens before the contract is deployed.
func recordTestRef(qualifier string) datastore.AddressRef {
	return datastore.AddressRef{
		ChainSelector: recordTestChainSelector,
		Type:          datastore.ContractType("BurnMintTokenPool"),
		Version:       semver.MustParse("1.5.0"),
		Qualifier:     qualifier,
	}
}

func recordTestDeployed(qualifier, address string) datastore.AddressRef {
	ref := recordTestRef(qualifier)
	ref.Address = address

	return ref
}

// A key nobody holds is claimed by the deployment that reaches it first.
func TestClaimAddressRef_ClaimsFreeKey(t *testing.T) {
	t.Parallel()
	ds := datastore.NewMemoryDataStore()

	ref := recordTestRef("LINK")
	claimed, err := claimAddressRef(ds, ref)
	require.NoError(t, err)
	assert.True(t, claimed, "a claim on a free key creates the reservation")

	got, err := ds.Addresses().Get(ref.Key())
	require.NoError(t, err)
	assert.Empty(t, got.Address, "a claim records the key, not an address")
}

// An empty-address reservation (e.g. made earlier by ReserveRefs in this private store) is this
// caller's to fill, not a new claim and not a conflict.
func TestClaimAddressRef_AcceptsOwnReservation(t *testing.T) {
	t.Parallel()
	ds := datastore.NewMemoryDataStore()

	ref := recordTestRef("LINK")
	require.NoError(t, ds.Addresses().Add(ref))

	claimed, err := claimAddressRef(ds, ref)
	require.NoError(t, err)
	assert.False(t, claimed, "an existing reservation is not this call's to release")

	refs, err := ds.Addresses().Fetch()
	require.NoError(t, err)
	assert.Len(t, refs, 1)
}

// A key already holding a deployed address cannot take a second contract. Crucially this is
// raised before anything is deployed, so nothing is stranded by it.
func TestClaimAddressRef_RejectsKeyHoldingADeployedContract(t *testing.T) {
	t.Parallel()
	ds := datastore.NewMemoryDataStore()

	require.NoError(t, ds.Addresses().Add(recordTestDeployed("", "0xaaa")))

	_, err := claimAddressRef(ds, recordTestRef(""))
	require.ErrorIs(t, err, datastore.ErrAddressRefExists)
	require.ErrorContains(t, err, "0xaaa")

	// The contract already recorded is untouched.
	got, getErr := ds.Addresses().Get(recordTestRef("").Key())
	require.NoError(t, getErr)
	assert.Equal(t, "0xaaa", got.Address)
}

// Distinct qualifiers are what let two instances of one type coexist.
func TestClaimAddressRef_DistinctQualifiersCoexist(t *testing.T) {
	t.Parallel()
	ds := datastore.NewMemoryDataStore()

	_, err := claimAddressRef(ds, recordTestRef("LINK"))
	require.NoError(t, err)
	_, err = claimAddressRef(ds, recordTestRef("USDC"))
	require.NoError(t, err)

	refs, err := ds.Addresses().Fetch()
	require.NoError(t, err)
	assert.Len(t, refs, 2)
}

// AddressesEqual has to fold case on EVM, where the same contract is routinely written both
// checksummed and lower-case, and must not fold it on Solana, where case is significant.
func TestAddressesEqual(t *testing.T) {
	t.Parallel()
	evm := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	solana := chainsel.SOLANA_DEVNET.Selector

	checksummed := "0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846"
	lowered := "0x0b9d5d9136855f6fec3c0993fee6e9ce8a297846"

	assert.True(t, AddressesEqual(evm, checksummed, lowered),
		"the same EVM contract written in two casings is one contract")
	assert.False(t, AddressesEqual(evm, checksummed, "0x1234567890123456789012345678901234567890"))

	assert.False(t, AddressesEqual(solana,
		"So11111111111111111111111111111111111111112",
		"so11111111111111111111111111111111111111112"),
		"Solana addresses are case-sensitive, so these are two different accounts")

	// An unrecognised selector keeps the distinction rather than erasing it.
	assert.False(t, AddressesEqual(0, checksummed, lowered))
}

// The record tests run on a real EVM chain so the address book's validation applies.
var (
	recordSepoliaSelector = chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	recordSepoliaAddr     = "0x5B5BBb15ECE0a4Ed8cDab22F902e83F66aBe848f"
)

func recordSepoliaTV() deployment.TypeAndVersion {
	return deployment.NewTypeAndVersion(deployment.ContractType("BurnMintTokenPool"), *semver.MustParse("1.5.0"))
}

// recordSepoliaQualifier is the qualifier every record test writes: a token pool for the LINK
// token, so the token symbol is the qualifier.
const recordSepoliaQualifier = "LINK"

func recordSepoliaKey(tv deployment.TypeAndVersion) datastore.AddressRefKey {
	return datastore.NewAddressRefKey(
		recordSepoliaSelector, datastore.ContractType(tv.Type), &tv.Version, recordSepoliaQualifier)
}

// The common case: one call, both registries hold the contract.
func TestRecordAddress_WritesBothRegistries(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()

	require.NoError(t, RecordAddress(ab, ds, recordSepoliaSelector, recordSepoliaAddr, tv, "LINK"))

	addrs, err := ab.AddressesForChain(recordSepoliaSelector)
	require.NoError(t, err)
	assert.Contains(t, addrs, recordSepoliaAddr)

	ref, err := ds.Addresses().Get(recordSepoliaKey(tv))
	require.NoError(t, err)
	assert.Equal(t, recordSepoliaAddr, ref.Address)
}

// The doc comment promises a re-run is a no-op. AddressBookMap.Save rejects a repeated address
// outright, so the second call has to recognise the entry as its own rather than save again.
func TestRecordAddress_RepeatedRecordingIsNoOp(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()

	require.NoError(t, RecordAddress(ab, ds, recordSepoliaSelector, recordSepoliaAddr, tv, "LINK"))
	require.NoError(t, RecordAddress(ab, ds, recordSepoliaSelector, recordSepoliaAddr, tv, "LINK"),
		"recording the same contract twice must not fail")

	addrs, err := ab.AddressesForChain(recordSepoliaSelector)
	require.NoError(t, err)
	assert.Len(t, addrs, 1)

	refs, err := ds.Addresses().Fetch()
	require.NoError(t, err)
	assert.Len(t, refs, 1)
}

// A checksummed address and its lower-case form are one contract on EVM, so re-recording under
// the other casing is the same no-op.
func TestRecordAddress_CaseVariantIsTheSameContract(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()
	lowered := "0x5b5bbb15ece0a4ed8cdab22f902e83f66abe848f"

	require.NoError(t, RecordAddress(ab, ds, recordSepoliaSelector, recordSepoliaAddr, tv, "LINK"))
	require.NoError(t, RecordAddress(ab, ds, recordSepoliaSelector, lowered, tv, "LINK"))

	addrs, err := ab.AddressesForChain(recordSepoliaSelector)
	require.NoError(t, err)
	assert.Len(t, addrs, 1)

	refs, err := ds.Addresses().Fetch()
	require.NoError(t, err)
	assert.Len(t, refs, 1)
}

// A key already holding a different contract is refused before anything is written, so the
// address book stays untouched.
func TestRecordAddress_RefusingAConflictWritesNothing(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()

	other := recordSepoliaDeployed(tv, "0x6619Bad7fadbc282B1EF2F6cC078fCbE61478792")
	require.NoError(t, ds.Addresses().Add(other))

	err := RecordAddress(ab, ds, recordSepoliaSelector, recordSepoliaAddr, tv, "LINK")
	require.ErrorIs(t, err, datastore.ErrAddressRefExists)

	_, err = ab.AddressesForChain(recordSepoliaSelector)
	require.ErrorIs(t, err, deployment.ErrChainNotFound, "the address book must be untouched")
}

func recordSepoliaDeployed(tv deployment.TypeAndVersion, address string) datastore.AddressRef {
	return datastore.AddressRef{
		ChainSelector: recordSepoliaSelector,
		Address:       address,
		Type:          datastore.ContractType(tv.Type),
		Version:       &tv.Version,
		Qualifier:     recordSepoliaQualifier,
	}
}

// The address book already holding the address under a different type is a disagreement, not a
// re-run, and is rejected before the datastore is written.
func TestRecordAddress_AddressBookDisagreementWritesNothing(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()

	require.NoError(t, ab.Save(recordSepoliaSelector, recordSepoliaAddr,
		deployment.NewTypeAndVersion(deployment.ContractType("Router"), *semver.MustParse("1.6.0"))))

	err := RecordAddress(ab, ds, recordSepoliaSelector, recordSepoliaAddr, tv, "LINK")
	require.ErrorContains(t, err, "the address book records")

	refs, err := ds.Addresses().Fetch()
	require.NoError(t, err)
	assert.Empty(t, refs, "the datastore must be untouched")
}

// An empty-address reservation made earlier in the run (via ReserveRefs) belongs to this caller
// and gets filled.
func TestRecordAddress_FillsOwnReservation(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()

	require.NoError(t, ds.Addresses().Add(recordSepoliaDeployed(tv, "")))
	require.NoError(t, RecordAddress(ab, ds, recordSepoliaSelector, recordSepoliaAddr, tv, "LINK"))

	ref, err := ds.Addresses().Get(recordSepoliaKey(tv))
	require.NoError(t, err)
	assert.Equal(t, recordSepoliaAddr, ref.Address)

	addrs, err := ab.AddressesForChain(recordSepoliaSelector)
	require.NoError(t, err)
	assert.Contains(t, addrs, recordSepoliaAddr)
}

// A nil address book is an explicit opt-out of the legacy registry, not an error.
func TestRecordAddress_NilAddressBookWritesOnlyTheDataStore(t *testing.T) {
	t.Parallel()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()

	require.NoError(t, RecordAddress(nil, ds, recordSepoliaSelector, recordSepoliaAddr, tv, "LINK"))

	ref, err := ds.Addresses().Get(recordSepoliaKey(tv))
	require.NoError(t, err)
	assert.Equal(t, recordSepoliaAddr, ref.Address)
}

func TestRecordAddress_RequiresDataStore(t *testing.T) {
	t.Parallel()
	err := RecordAddress(deployment.NewMemoryAddressBook(), nil, recordSepoliaSelector, recordSepoliaAddr, recordSepoliaTV(), "LINK")
	require.ErrorContains(t, err, "requires a datastore")
}

// recordDeployChain is a zkSync-flavoured test chain: the IsZkSyncVM flag skips transaction
// confirmation, so deployment tests need no client.
func recordDeployChain() cldf_evm.Chain {
	return cldf_evm.Chain{Selector: recordSepoliaSelector, IsZkSyncVM: true}
}

// The happy path records the contract in both registries under the declared qualifier.
func TestDeployContractAndRecord_WritesBothRegistries(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()

	got, err := DeployContractAndRecord(logger.Test(t), recordDeployChain(), ab, ds, tv, "LINK",
		func(cldf_evm.Chain) deployment.ContractDeploy[any] {
			return deployment.ContractDeploy[any]{Address: common.HexToAddress(recordSepoliaAddr), Tv: tv}
		})
	require.NoError(t, err)
	require.NotNil(t, got)

	addrs, err := ab.AddressesForChain(recordSepoliaSelector)
	require.NoError(t, err)
	assert.Contains(t, addrs, recordSepoliaAddr)

	ref, err := ds.Addresses().Get(recordSepoliaKey(tv))
	require.NoError(t, err)
	assert.Equal(t, recordSepoliaAddr, ref.Address)
}

// A key already holding a deployed contract refuses the deploy before the closure runs: the
// contract is never even sent.
func TestDeployContractAndRecord_RefusesATakenKeyBeforeDeploying(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()

	require.NoError(t, ds.Addresses().Add(recordSepoliaDeployed(tv, "0x6619Bad7fadbc282B1EF2F6cC078fCbE61478792")))

	deployed := false
	_, err := DeployContractAndRecord(logger.Test(t), recordDeployChain(), ab, ds, tv, "LINK",
		func(cldf_evm.Chain) deployment.ContractDeploy[any] {
			deployed = true

			return deployment.ContractDeploy[any]{Address: common.HexToAddress(recordSepoliaAddr), Tv: tv}
		})
	require.ErrorIs(t, err, datastore.ErrAddressRefExists)
	assert.False(t, deployed, "the deploy closure must not run when the key is taken")
}

// A failed deployment takes its reservation with it: an empty-address ref left behind would be
// merged and persisted as if it were a record of progress, and it is not one.
func TestDeployContractAndRecord_DeployFailureReleasesTheReservation(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()

	_, err := DeployContractAndRecord(logger.Test(t), recordDeployChain(), ab, ds, tv, "LINK",
		func(cldf_evm.Chain) deployment.ContractDeploy[any] {
			return deployment.ContractDeploy[any]{Err: errors.New("deploy boom")}
		})
	require.ErrorContains(t, err, "deploy boom")

	refs, err := ds.Addresses().Fetch()
	require.NoError(t, err)
	assert.Empty(t, refs, "the reservation must not survive the failed deployment")

	_, err = ab.AddressesForChain(recordSepoliaSelector)
	require.ErrorIs(t, err, deployment.ErrChainNotFound)
}

// A reservation this call did not make (e.g. from ReserveRefs) is not its to release: a failed
// deployment must leave the reserved key in place.
func TestDeployContractAndRecord_FailureKeepsAPreExistingReservation(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()

	require.NoError(t, ds.Addresses().Add(recordSepoliaDeployed(tv, "")))

	_, err := DeployContractAndRecord(logger.Test(t), recordDeployChain(), ab, ds, tv, "LINK",
		func(cldf_evm.Chain) deployment.ContractDeploy[any] {
			return deployment.ContractDeploy[any]{Err: errors.New("deploy boom")}
		})
	require.ErrorContains(t, err, "deploy boom")

	got, err := ds.Addresses().Get(recordSepoliaKey(tv))
	require.NoError(t, err, "a reservation made earlier in the run is not this call's to release")
	assert.Empty(t, got.Address)
}

// The deployed type and version are checked before the address book is written: a mismatch
// leaves neither registry touched.
func TestDeployContractAndRecord_TypeMismatchWritesNothing(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()
	deployed := deployment.NewTypeAndVersion(deployment.ContractType("Router"), *semver.MustParse("1.6.0"))

	_, err := DeployContractAndRecord(logger.Test(t), recordDeployChain(), ab, ds, tv, "LINK",
		func(cldf_evm.Chain) deployment.ContractDeploy[any] {
			return deployment.ContractDeploy[any]{Address: common.HexToAddress(recordSepoliaAddr), Tv: deployed}
		})
	require.ErrorContains(t, err, "was declared to DeployContractAndRecord")

	_, abErr := ab.AddressesForChain(recordSepoliaSelector)
	require.ErrorIs(t, abErr, deployment.ErrChainNotFound, "the type check must come before the address book write")

	refs, err := ds.Addresses().Fetch()
	require.NoError(t, err)
	assert.Empty(t, refs, "the reservation must not survive the mismatch")
}

// An address book rejection (here: the zero address) likewise leaves no datastore reservation.
func TestDeployContractAndRecord_AddressBookFailureReleasesTheReservation(t *testing.T) {
	t.Parallel()
	ab := deployment.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()
	tv := recordSepoliaTV()

	_, err := DeployContractAndRecord(logger.Test(t), recordDeployChain(), ab, ds, tv, "LINK",
		func(cldf_evm.Chain) deployment.ContractDeploy[any] {
			return deployment.ContractDeploy[any]{Address: common.Address{}, Tv: tv}
		})
	require.Error(t, err)

	refs, err := ds.Addresses().Fetch()
	require.NoError(t, err)
	assert.Empty(t, refs, "the reservation must not survive the address book failure")
}
