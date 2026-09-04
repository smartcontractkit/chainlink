package changeset

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// mergeTestEnvironment builds a minimal environment for merge tests: a memory address book, a
// sealed memory datastore, and a test logger.
func mergeTestEnvironment(t *testing.T) cldf.Environment {
	t.Helper()

	return cldf.Environment{
		Logger:            logger.Test(t),
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		DataStore:         datastore.NewMemoryDataStore().Seal(),
	}
}

func routerRef(chainSelector uint64, addr string) datastore.AddressRef {
	return datastore.AddressRef{
		ChainSelector: chainSelector,
		Address:       addr,
		Type:          datastore.ContractType("Router"),
		Version:       semver.MustParse("1.6.0"),
	}
}

func TestMergeChangesetOutputMergesDataStoreAndReports(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	dest := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, dest.DataStore.Addresses().Add(routerRef(5009297550715157269, "0xaaa")))

	src := cldf.ChangesetOutput{
		DataStore: datastore.NewMemoryDataStore(),
		Reports:   []operations.Report[any, any]{{ID: "report-1"}},
	}
	require.NoError(t, src.DataStore.Addresses().Add(datastore.AddressRef{
		ChainSelector: 5009297550715157269,
		Address:       "0xbbb",
		Type:          datastore.ContractType("FeeQuoter"),
		Version:       semver.MustParse("1.6.0"),
	}))

	require.NoError(t, MergeChangesetOutput(e, &dest, src))

	got, err := dest.DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, got, 2, "source datastore refs must be merged into the destination")
	require.Len(t, dest.Reports, 1, "source reports must be appended to the destination")
	require.Equal(t, "report-1", dest.Reports[0].ID)
}

func TestMergeChangesetOutputDoesNotAliasSourceDataStore(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	first := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, first.DataStore.Addresses().Add(routerRef(5009297550715157269, "0xaaa")))

	second := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, second.DataStore.Addresses().Add(datastore.AddressRef{
		ChainSelector: 5009297550715157269,
		Address:       "0xbbb",
		Type:          datastore.ContractType("FeeQuoter"),
		Version:       semver.MustParse("1.6.0"),
	}))

	var dest cldf.ChangesetOutput
	require.NoError(t, MergeChangesetOutput(e, &dest, first))
	require.NoError(t, MergeChangesetOutput(e, &dest, second))

	destRefs, err := dest.DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, destRefs, 2)

	firstRefs, err := first.DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, firstRefs, 1, "the first sub-changeset's own output must be untouched")
	require.Equal(t, "0xaaa", firstRefs[0].Address)
}

func TestMergeChangesetOutputRejectsConflictingRefs(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	ref := func(addr string) datastore.AddressRef {
		return datastore.AddressRef{
			ChainSelector: 5009297550715157269,
			Address:       addr,
			Type:          datastore.ContractType("BurnMintTokenPool"),
			Version:       semver.MustParse("1.5.0"),
			Qualifier:     "",
		}
	}

	dest := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, dest.DataStore.Addresses().Add(ref("0xaaa")))

	src := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, src.DataStore.Addresses().Add(ref("0xbbb")))

	err := MergeChangesetOutput(e, &dest, src)
	require.ErrorContains(t, err, "both claim")
	require.ErrorContains(t, err, "distinct qualifiers")

	// The same address under the same key is not a conflict, only a repeated write.
	same := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, same.DataStore.Addresses().Add(ref("0xaaa")))
	require.NoError(t, MergeChangesetOutput(e, &dest, same))
}

func TestMergeChangesetOutputPropagatesDataStoreToEnv(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	src := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, src.DataStore.Addresses().Add(routerRef(5009297550715157269, "0xaaa")))

	var dest cldf.ChangesetOutput
	require.NoError(t, MergeChangesetOutput(e, &dest, src))

	got, err := e.DataStore.Addresses().Get(
		datastore.NewAddressRefKey(5009297550715157269, datastore.ContractType("Router"), semver.MustParse("1.6.0"), ""),
	)
	require.NoError(t, err)
	require.Equal(t, "0xaaa", got.Address)
}

func TestMergeChangesetOutputHandlesTypedNilDataStore(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	var nilStore *datastore.MemoryDataStore

	dest := cldf.ChangesetOutput{DataStore: nilStore}
	src := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, src.DataStore.Addresses().Add(routerRef(5009297550715157269, "0xaaa")))

	require.NoError(t, MergeChangesetOutput(e, &dest, src))

	destRefs, err := dest.DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, destRefs, 1)

	// A typed-nil source is a no-op rather than a panic.
	require.NoError(t, MergeChangesetOutput(e, &dest, cldf.ChangesetOutput{DataStore: nilStore}))
}

func TestMergeChangesetOutputRejectsRefWithoutVersion(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	src := datastore.NewMemoryDataStore()
	src.AddressRefStore.Records = append(src.AddressRefStore.Records, datastore.AddressRef{
		ChainSelector: 5009297550715157269,
		Address:       "0xaaa",
		Type:          datastore.ContractType("Router"),
		Version:       nil,
	})

	dest := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}

	err := MergeChangesetOutput(e, &dest, cldf.ChangesetOutput{DataStore: src})
	require.ErrorIs(t, err, datastore.ErrAddressRefVersionRequired)
	require.ErrorContains(t, err, "0xaaa")
}

func TestMergeChangesetOutputRejectsEmptyAddressRef(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	// A reserved-but-never-filled key (as ReserveRefs leaves behind when a changeset returns
	// without deploying everything it planned) is a pending claim, not a record.
	src := datastore.NewMemoryDataStore()
	require.NoError(t, src.Addresses().Add(datastore.AddressRef{
		ChainSelector: 5009297550715157269,
		Address:       "",
		Type:          datastore.ContractType("Router"),
		Version:       semver.MustParse("1.6.0"),
		Qualifier:     "",
	}))

	dest := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}

	err := MergeChangesetOutput(e, &dest, cldf.ChangesetOutput{DataStore: src})
	require.ErrorContains(t, err, "has no address")

	envRefs, envErr := e.DataStore.Addresses().Fetch()
	require.NoError(t, envErr)
	require.Empty(t, envRefs, "the reservation must not reach the environment")

	destRefs, destErr := dest.DataStore.Addresses().Fetch()
	require.NoError(t, destErr)
	require.Empty(t, destRefs, "the reservation must not reach the destination")
}

func TestMergeChangesetOutputComparesAddressesExactly(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	ref := func(addr string) datastore.AddressRef {
		return datastore.AddressRef{
			ChainSelector: 124615329519749607,
			Address:       addr,
			Type:          datastore.ContractType("TokenPool"),
			Version:       semver.MustParse("1.0.0"),
		}
	}

	dest := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, dest.DataStore.Addresses().Add(ref("AbCdEf")))

	src := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, src.DataStore.Addresses().Add(ref("aBcDeF")))

	require.ErrorContains(t, MergeChangesetOutput(e, &dest, src), "both claim")
}

func TestMergeChangesetOutputRejectsConflictingMetadata(t *testing.T) {
	t.Parallel()
	t.Run("chain metadata", func(t *testing.T) {
		t.Parallel()
		e := mergeTestEnvironment(t)

		dest := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
		require.NoError(t, dest.DataStore.ChainMetadata().Add(datastore.ChainMetadata{
			ChainSelector: 5009297550715157269, Metadata: map[string]any{"owner": "a"},
		}))

		src := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
		require.NoError(t, src.DataStore.ChainMetadata().Add(datastore.ChainMetadata{
			ChainSelector: 5009297550715157269, Metadata: map[string]any{"owner": "b"},
		}))

		require.ErrorContains(t, MergeChangesetOutput(e, &dest, src), "conflicting chain metadata")
	})

	t.Run("contract metadata", func(t *testing.T) {
		t.Parallel()
		e := mergeTestEnvironment(t)

		dest := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
		require.NoError(t, dest.DataStore.ContractMetadata().Add(datastore.ContractMetadata{
			ChainSelector: 5009297550715157269, Address: "0xaaa", Metadata: map[string]any{"n": "1"},
		}))

		src := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
		require.NoError(t, src.DataStore.ContractMetadata().Add(datastore.ContractMetadata{
			ChainSelector: 5009297550715157269, Address: "0xaaa", Metadata: map[string]any{"n": "2"},
		}))

		require.ErrorContains(t, MergeChangesetOutput(e, &dest, src), "conflicting contract metadata")
	})

	t.Run("environment metadata", func(t *testing.T) {
		t.Parallel()
		e := mergeTestEnvironment(t)

		dest := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
		require.NoError(t, dest.DataStore.EnvMetadata().Set(datastore.EnvMetadata{
			Metadata: map[string]any{"release": "a"},
		}))

		src := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
		require.NoError(t, src.DataStore.EnvMetadata().Set(datastore.EnvMetadata{
			Metadata: map[string]any{"release": "b"},
		}))

		require.ErrorContains(t, MergeChangesetOutput(e, &dest, src), "conflicting environment metadata")
	})

	t.Run("identical metadata merges", func(t *testing.T) {
		t.Parallel()
		e := mergeTestEnvironment(t)

		meta := map[string]any{"owner": "a"}

		dest := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
		require.NoError(t, dest.DataStore.ChainMetadata().Add(datastore.ChainMetadata{
			ChainSelector: 5009297550715157269, Metadata: meta,
		}))

		src := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
		require.NoError(t, src.DataStore.ChainMetadata().Add(datastore.ChainMetadata{
			ChainSelector: 5009297550715157269, Metadata: meta,
		}))

		require.NoError(t, MergeChangesetOutput(e, &dest, src))
	})
}

func TestMergeChangesetOutputLeavesNothingHalfMerged(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	chainSelector := chainsel.TEST_90000001.Selector
	addrA := "0x5B5BBb15ECE0a4Ed8cDab22F902e83F66aBe848f"
	addrB := "0x6619Bad7fadbc282B1EF2F6cC078fCbE61478792"
	tv := cldf.NewTypeAndVersion("Router", *semver.MustParse("1.6.0"))

	ref := func(addr string) datastore.AddressRef {
		return datastore.AddressRef{
			ChainSelector: chainSelector,
			Address:       addr,
			Type:          datastore.ContractType("Router"),
			Version:       semver.MustParse("1.6.0"),
		}
	}

	// The destination already holds one Router; the source deploys a different one under the
	// same key, which is the conflict that will reject the merge.
	dest := cldf.ChangesetOutput{
		AddressBook: cldf.NewMemoryAddressBook(), //nolint:staticcheck // Phase 1 merge still stages the address book
		DataStore:   datastore.NewMemoryDataStore(),
	}
	require.NoError(t, dest.AddressBook.Save(chainSelector, addrA, tv)) //nolint:staticcheck // Phase 1 merge still stages the address book
	require.NoError(t, dest.DataStore.Addresses().Add(ref(addrA)))

	src := cldf.ChangesetOutput{
		AddressBook: cldf.NewMemoryAddressBook(), //nolint:staticcheck // Phase 1 merge still stages the address book
		DataStore:   datastore.NewMemoryDataStore(),
	}
	require.NoError(t, src.AddressBook.Save(chainSelector, addrB, tv)) //nolint:staticcheck // Phase 1 merge still stages the address book
	require.NoError(t, src.DataStore.Addresses().Add(ref(addrB)))

	require.Error(t, MergeChangesetOutput(e, &dest, src))

	// The address book merge would have succeeded on its own, and used to run first. It must
	// not have been applied.
	destAddrs, err := dest.AddressBook.Addresses() //nolint:staticcheck // Phase 1 merge still stages the address book
	require.NoError(t, err)
	require.Len(t, destAddrs[chainSelector], 1, "destination address book must be unchanged")
	require.Contains(t, destAddrs[chainSelector], addrA)

	envAddrs, err := e.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Empty(t, envAddrs, "environment address book must be unchanged")

	envRefs, err := e.DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Empty(t, envRefs, "environment data store must be unchanged")

	destRefs, err := dest.DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, destRefs, 1)
	require.Equal(t, addrA, destRefs[0].Address)
}

func TestMergeChangesetOutputPublishesFirstMergeToEnv(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	chainSelector := chainsel.TEST_90000001.Selector
	addr := "0x5B5BBb15ECE0a4Ed8cDab22F902e83F66aBe848f"
	tv := cldf.NewTypeAndVersion("Router", *semver.MustParse("1.6.0"))

	src := cldf.ChangesetOutput{AddressBook: cldf.NewMemoryAddressBook()} //nolint:staticcheck // Phase 1 merge still stages the address book
	require.NoError(t, src.AddressBook.Save(chainSelector, addr, tv))     //nolint:staticcheck // Phase 1 merge still stages the address book

	var dest cldf.ChangesetOutput
	require.NoError(t, MergeChangesetOutput(e, &dest, src))

	envAddrs, err := e.ExistingAddresses.Addresses()
	require.NoError(t, err)
	require.Contains(t, envAddrs[chainSelector], addr)

	// The destination is a copy, not the source's own book: merging again must not corrupt the
	// first sub-changeset's output.
	require.NotSame(t, src.AddressBook, dest.AddressBook) //nolint:staticcheck // Phase 1 merge still stages the address book
}

func TestMergeChangesetOutputToleratesKnownAddresses(t *testing.T) {
	t.Parallel()
	chainSelector := chainsel.TEST_90000001.Selector
	addr := "0x5B5BBb15ECE0a4Ed8cDab22F902e83F66aBe848f"
	tv := cldf.NewTypeAndVersion("Router", *semver.MustParse("1.6.0"))

	t.Run("same type and version", func(t *testing.T) {
		t.Parallel()
		e := mergeTestEnvironment(t)
		require.NoError(t, e.ExistingAddresses.Save(chainSelector, addr, tv))

		src := cldf.ChangesetOutput{AddressBook: cldf.NewMemoryAddressBook()} //nolint:staticcheck // Phase 1 merge still stages the address book
		require.NoError(t, src.AddressBook.Save(chainSelector, addr, tv))     //nolint:staticcheck // Phase 1 merge still stages the address book

		var dest cldf.ChangesetOutput
		require.NoError(t, MergeChangesetOutput(e, &dest, src))
	})

	t.Run("different type", func(t *testing.T) {
		t.Parallel()
		e := mergeTestEnvironment(t)
		require.NoError(t, e.ExistingAddresses.Save(chainSelector, addr, tv))

		src := cldf.ChangesetOutput{AddressBook: cldf.NewMemoryAddressBook()} //nolint:staticcheck // Phase 1 merge still stages the address book
		require.NoError(t, src.AddressBook.Save(chainSelector, addr,          //nolint:staticcheck // Phase 1 merge still stages the address book
			cldf.NewTypeAndVersion("FeeQuoter", *semver.MustParse("1.6.0"))))

		var dest cldf.ChangesetOutput
		err := MergeChangesetOutput(e, &dest, src)
		require.ErrorContains(t, err, "the environment records")
	})
}

func TestPlanDataStoreSupersessionUsesFamilyCasingRules(t *testing.T) {
	t.Parallel()
	chainSelector := chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector
	checksummed := "0x5B5BBb15ECE0a4Ed8cDab22F902e83F66aBe848f"

	srcWith := func(addr string) datastore.MutableDataStore {
		src := datastore.NewMemoryDataStore()
		require.NoError(t, src.Addresses().Add(datastore.AddressRef{
			ChainSelector: chainSelector,
			Address:       addr,
			Type:          datastore.ContractType("Router"),
			Version:       semver.MustParse("1.6.0"),
		}))

		return src
	}

	newEnv := func(t *testing.T) cldf.Environment {
		t.Helper()

		e := mergeTestEnvironment(t)
		mutable, ok := e.DataStore.Addresses().(datastore.MutableAddressRefStore)
		require.True(t, ok)
		require.NoError(t, mutable.Upsert(datastore.AddressRef{
			ChainSelector: chainSelector,
			Address:       checksummed,
			Type:          datastore.ContractType("Router"),
			Version:       semver.MustParse("1.6.0"),
		}))

		return e
	}

	t.Run("same contract in the other casing is not a supersession", func(t *testing.T) {
		t.Parallel()
		plan := &changesetMergePlan{}
		require.NoError(t, plan.planDataStore(newEnv(t), nil, srcWith("0x5b5bbb15ece0a4ed8cdab22f902e83f66abe848f")))
		require.Empty(t, plan.superseded)
	})

	t.Run("a different address is a supersession", func(t *testing.T) {
		t.Parallel()
		plan := &changesetMergePlan{}
		require.NoError(t, plan.planDataStore(newEnv(t), nil, srcWith("0x6619Bad7fadbc282B1EF2F6cC078fCbE61478792")))
		require.Len(t, plan.superseded, 1)
		require.Equal(t, checksummed, plan.superseded[0].Address)
	})
}

func TestMergeChangesetOutputSupersedesEnvRef(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	key := datastore.NewAddressRefKey(
		5009297550715157269, datastore.ContractType("Router"), semver.MustParse("1.6.0"), "",
	)

	mutable, ok := e.DataStore.Addresses().(datastore.MutableAddressRefStore)
	require.True(t, ok)
	require.NoError(t, mutable.Upsert(datastore.AddressRef{
		ChainSelector: 5009297550715157269,
		Address:       "0xold",
		Type:          datastore.ContractType("Router"),
		Version:       semver.MustParse("1.6.0"),
	}))

	src := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, src.DataStore.Addresses().Add(datastore.AddressRef{
		ChainSelector: 5009297550715157269,
		Address:       "0xnew",
		Type:          datastore.ContractType("Router"),
		Version:       semver.MustParse("1.6.0"),
	}))

	var dest cldf.ChangesetOutput
	require.NoError(t, MergeChangesetOutput(e, &dest, src))

	got, err := e.DataStore.Addresses().Get(key)
	require.NoError(t, err)
	require.Equal(t, "0xnew", got.Address)
}

func TestMergeChangesetOutputPropagatesMetadataToEnv(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	chainMeta := datastore.ChainMetadata{ChainSelector: 5009297550715157269, Metadata: map[string]any{"owner": "ccip"}}
	contractMeta := datastore.ContractMetadata{
		ChainSelector: 5009297550715157269, Address: "0xaaa", Metadata: map[string]any{"n": "1"},
	}
	envMeta := datastore.EnvMetadata{Metadata: map[string]any{"release": "2026"}}

	src := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, src.DataStore.ChainMetadata().Add(chainMeta))
	require.NoError(t, src.DataStore.ContractMetadata().Add(contractMeta))
	require.NoError(t, src.DataStore.EnvMetadata().Set(envMeta))

	var dest cldf.ChangesetOutput
	require.NoError(t, MergeChangesetOutput(e, &dest, src))

	gotChain, err := e.DataStore.ChainMetadata().Fetch()
	require.NoError(t, err)
	require.Len(t, gotChain, 1)
	require.Equal(t, map[string]any{"owner": "ccip"}, gotChain[0].Metadata)

	gotContract, err := e.DataStore.ContractMetadata().Fetch()
	require.NoError(t, err)
	require.Len(t, gotContract, 1)
	require.Equal(t, "0xaaa", gotContract[0].Address)

	gotEnv, err := e.DataStore.EnvMetadata().Get()
	require.NoError(t, err)
	require.Equal(t, map[string]any{"release": "2026"}, gotEnv.Metadata)
}

func TestMergeChangesetOutputPropagatesDeletionsToEnv(t *testing.T) {
	t.Parallel()
	e := mergeTestEnvironment(t)

	key := datastore.NewAddressRefKey(
		5009297550715157269, datastore.ContractType("Router"), semver.MustParse("1.6.0"), "",
	)

	// A contract already in the environment that a sub-changeset deletes.
	mutable, ok := e.DataStore.Addresses().(datastore.MutableAddressRefStore)
	require.True(t, ok)
	require.NoError(t, mutable.Upsert(datastore.AddressRef{
		ChainSelector: 5009297550715157269,
		Address:       "0xold",
		Type:          datastore.ContractType("Router"),
		Version:       semver.MustParse("1.6.0"),
	}))

	src := cldf.ChangesetOutput{DataStore: datastore.NewMemoryDataStore()}
	require.NoError(t, src.DataStore.Addresses().Add(datastore.AddressRef{
		ChainSelector: 5009297550715157269,
		Address:       "0xold",
		Type:          datastore.ContractType("Router"),
		Version:       semver.MustParse("1.6.0"),
	}))
	// Stage the deletion so the merge propagates it, as MemoryDataStore.Merge does.
	require.NoError(t, src.DataStore.Addresses().RemoteDelete(key))

	var dest cldf.ChangesetOutput
	require.NoError(t, MergeChangesetOutput(e, &dest, src))

	_, err := e.DataStore.Addresses().Get(key)
	require.ErrorIs(t, err, datastore.ErrAddressRefNotFound, "the deleted ref must be gone from the environment")
}
