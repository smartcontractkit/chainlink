package stateview

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	mcmsolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	mcmscontracts "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/contracts/mcms"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	"github.com/smartcontractkit/ccip-owner-contracts/gethwrappers"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipshared "github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

func testRef(chainSel uint64, addr string, t cldf.ContractType, ver *semver.Version, qualifier string, labels ...string) datastore.AddressRef {
	ref := datastore.AddressRef{
		ChainSelector: chainSel,
		Address:       addr,
		Type:          datastore.ContractType(t),
		Qualifier:     qualifier,
	}
	if ver != nil {
		v := *ver
		ref.Version = &v
	}
	if len(labels) > 0 {
		ref.Labels = datastore.NewLabelSet(labels...)
	}
	return ref
}

func TestAddressRefsToTypeVersions(t *testing.T) {
	v16 := deployment.Version1_6_0
	refs := []datastore.AddressRef{
		testRef(1, "0xA", "OnRamp", &v16, "123", "someLabel"),
		testRef(1, "0xB", "WETH9", nil, ""), // versionless: skipped
	}
	got := AddressRefsToTypeVersions(refs)
	require.Len(t, got, 1)
	tvs := got["0xA"]
	require.Len(t, tvs, 1)
	require.Equal(t, "1.6.0", tvs[0].Version.String())
	require.True(t, tvs[0].Labels.Contains("someLabel"))
}

func TestDataStoreTypeVersionsForChain_FiltersAndUniqueness(t *testing.T) {
	v16 := deployment.Version1_6_0
	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(testRef(1, "0xA", "Router", &v16, "", "legacy")))
	// the superseded twin lives under its pre-rekey qualifier (the store keys by
	// (chain, type, version, qualifier))
	require.NoError(t, ds.Addresses().Add(testRef(1, "0xB", "Router", &v16, "old", SupersededLabel)))
	require.NoError(t, ds.Addresses().Add(testRef(2, "0xD", "Router", &v16, ""))) // other chain

	got, err := DataStoreTypeVersionsForChain(cldf.Environment{DataStore: ds.Seal()}, 1)
	require.NoError(t, err)
	require.Contains(t, got, "0xA")
	require.NotContains(t, got, "0xB") // superseded
	require.NotContains(t, got, "0xD") // other chain
	// versionless refs cannot even be Added here; AddressRefsToTypeVersions skips them

	// two active refs claiming one (type, version, qualifier) identity with different
	// addresses — unreachable through this store's Add, possible in hand-loaded stores
	err = ccipshared.CheckRefUniqueness([]datastore.AddressRef{
		testRef(1, "0xA", "Router", &v16, ""),
		testRef(1, "0xB", "Router", &v16, ""),
	})
	require.ErrorContains(t, err, "0xA")
	require.ErrorContains(t, err, "0xB")

	// same (type, version) under different qualifiers is legitimate
	ds3 := datastore.NewMemoryDataStore()
	require.NoError(t, ds3.Addresses().Add(testRef(1, "0xA", "FeeQuoter", &v16, "instA")))
	require.NoError(t, ds3.Addresses().Add(testRef(1, "0xB", "FeeQuoter", &v16, "instB")))
	got, err = DataStoreTypeVersionsForChain(cldf.Environment{DataStore: ds3.Seal()}, 1)
	require.NoError(t, err)
	require.Len(t, got, 2)
}

func TestMCMSBundleRefs_Isolation(t *testing.T) {
	v10 := deployment.Version1_0_0
	refs := []datastore.AddressRef{
		testRef(1, "0xQualifiedTimelock", mcmscontracts.RBACTimelock, &v10, DefaultMCMSQualifier),
		testRef(1, "0xQualifiedProposer", mcmscontracts.ProposerManyChainMultisig, &v10, DefaultMCMSQualifier),
		testRef(1, "0xEmptyTimelock", mcmscontracts.RBACTimelock, &v10, ""),
		testRef(1, "0xEmptyCallProxy", mcmscontracts.CallProxy, &v10, ""),
		testRef(1, "0xVersionless", mcmscontracts.RBACTimelock, nil, DefaultMCMSQualifier),
		testRef(1, "0xSuperseded", mcmscontracts.CallProxy, &v10, DefaultMCMSQualifier, SupersededLabel),
	}

	bundle, err := ccipshared.MCMSBundleRefs(refs, 1, DefaultMCMSQualifier)
	require.NoError(t, err)
	require.Len(t, bundle, 2) // qualified rows only; versionless and superseded dropped
	for _, ref := range bundle {
		require.Equal(t, DefaultMCMSQualifier, ref.Qualifier)
	}

	// custom qualifier with no rows fails closed — no empty-qualifier fallback
	_, err = ccipshared.MCMSBundleRefs(refs, 1, "RMNMCMS")
	require.ErrorContains(t, err, "no mcms refs for chain 1")

	// default qualifier with no rows yields an empty bundle (MaybeLoad semantics)
	bundle, err = ccipshared.MCMSBundleRefs(refs, 7, DefaultMCMSQualifier)
	require.NoError(t, err)
	require.Empty(t, bundle)

	// An unqualified bundle IS the legacy fallback for the default qualifier
	// (singletons and DeployMCMSWithTimelockV2 test deployments are empty-qualified);
	// dedicated bundles (RMNMCMS etc.) never fall back.
	emptyOnly := []datastore.AddressRef{testRef(1, "0xEmptyOnly", mcmscontracts.RBACTimelock, &v10, "")}
	bundle, err = ccipshared.MCMSBundleRefs(emptyOnly, 1, DefaultMCMSQualifier)
	require.NoError(t, err)
	require.Len(t, bundle, 1)
	require.Empty(t, bundle[0].Qualifier)

	// the fallback never fires for a custom qualifier
	_, err = ccipshared.MCMSBundleRefs(emptyOnly, 1, "RMNMCMS")
	require.ErrorContains(t, err, "no mcms refs for chain 1")

	// two active refs, one identity
	dup := append([]datastore.AddressRef(nil), refs...)
	dup = append(dup, testRef(1, "0xQualifiedTimelock2", mcmscontracts.RBACTimelock, &v10, DefaultMCMSQualifier))
	_, err = ccipshared.MCMSBundleRefs(dup, 1, DefaultMCMSQualifier)
	require.ErrorContains(t, err, "both")
}

func TestLoadChainState_LabeledRefDispatch(t *testing.T) {
	// A labeled ref must dispatch like its unlabeled twin; String() would fold the label in.
	chain := cldf_evm.Chain{Selector: 1}
	addr := "0x00000000000000000000000000000000000000Ab"
	addresses := map[string][]cldf.TypeAndVersion{
		addr: {
			{Type: commontypes.RBACTimelock, Version: deployment.Version1_0_0, Labels: cldf.NewLabelSet("someLabel")},
		},
	}
	state, err := LoadChainState(t.Context(), chain, addresses, WithMCMSQualifier(DefaultMCMSQualifier))
	require.NoError(t, err)
	require.Equal(t, gethwrappers.RBACTimelockABI, state.ABIByAddress[addr])
}

func TestLoadChainState_SingularAmbiguity(t *testing.T) {
	chain := cldf_evm.Chain{Selector: 1}
	v12 := deployment.Version1_2_0
	refs := []datastore.AddressRef{
		{Address: "0x00000000000000000000000000000000000000A1", ChainSelector: 1, Type: datastore.ContractType("Router"), Version: &v12},
		{Address: "0x00000000000000000000000000000000000000B2", ChainSelector: 1, Type: datastore.ContractType("Router"), Version: &v12, Labels: datastore.NewLabelSet("other")},
	}
	_, err := loadChainStateFromDataStore(t.Context(), chain, refs, WithMCMSQualifier(DefaultMCMSQualifier))
	require.ErrorContains(t, err, "datastore is ambiguous")
}

func TestLoadChainState_UnmodeledDuplicatesSkipped(t *testing.T) {
	// types the loader does not model are skipped in datastore mode, duplicates included —
	// they load nowhere, so multiplicity cannot corrupt state
	chain := cldf_evm.Chain{Selector: 1}
	addresses := map[string][]cldf.TypeAndVersion{
		"0x00000000000000000000000000000000000000A1": {{Type: "BurnMintERC20WithDripToken", Version: deployment.Version1_0_0}},
		"0x00000000000000000000000000000000000000B2": {{Type: "BurnMintERC20WithDripToken", Version: deployment.Version1_0_0}},
		"0x00000000000000000000000000000000000000C3": {{Type: "BurnMintERC20WithDripToken", Version: deployment.Version1_0_0}},
	}
	_, err := LoadChainState(t.Context(), chain, addresses, WithMCMSQualifier(DefaultMCMSQualifier), WithTolerateUnknownContractTypes())
	require.NoError(t, err)

	// without tolerance the unknown type still errors
	_, err = LoadChainState(t.Context(), chain, addresses, WithMCMSQualifier(DefaultMCMSQualifier))
	require.ErrorContains(t, err, "unknown contract")
}

func TestLoadChainState_RegistryModuleOrderDeterministic(t *testing.T) {
	// RegistryModule appends in sorted address order, not map order
	chain := cldf_evm.Chain{Selector: 1}
	addresses := map[string][]cldf.TypeAndVersion{
		"0x00000000000000000000000000000000000000B2": {{Type: "RegistryModuleOwnerCustom", Version: deployment.Version1_6_0}},
		"0x00000000000000000000000000000000000000A1": {{Type: "RegistryModuleOwnerCustom", Version: deployment.Version1_6_0}},
	}
	state, err := LoadChainState(t.Context(), chain, addresses, WithMCMSQualifier(DefaultMCMSQualifier))
	require.NoError(t, err)
	require.Len(t, state.RegistryModules1_6, 2)
	require.Equal(t, common.HexToAddress("0x00000000000000000000000000000000000000A1"), state.RegistryModules1_6[0].Address())
	require.Equal(t, common.HexToAddress("0x00000000000000000000000000000000000000B2"), state.RegistryModules1_6[1].Address())
}

func TestLoadChainState_PerTokenProxyDeterministic(t *testing.T) {
	// USDCTokenPoolProxy is per-token multi-instance in the datastore; the legacy state keeps
	// one deterministically (sorted order — here the address-book's own pick)
	chain := cldf_evm.Chain{Selector: 1}
	v20 := deployment.Version2_0_0
	addresses := map[string][]cldf.TypeAndVersion{
		"0x00000000000000000000000000000000000000B2": {{Type: "USDCTokenPoolProxy", Version: v20}},
		"0x00000000000000000000000000000000000000A1": {{Type: "USDCTokenPoolProxy", Version: v20}},
	}
	state, err := LoadChainState(t.Context(), chain, addresses, WithMCMSQualifier(DefaultMCMSQualifier))
	require.NoError(t, err)
	require.Equal(t, common.HexToAddress("0x00000000000000000000000000000000000000B2"), state.USDCTokenPoolProxies[deployment.Version2_0_0])
}

func TestLoadChainState_FeeQuoterEqualVersionAmbiguity(t *testing.T) {
	chain := cldf_evm.Chain{Selector: 1}
	v16 := deployment.Version1_6_0
	addresses := map[string][]cldf.TypeAndVersion{
		"0x00000000000000000000000000000000000000A1": {{Type: "FeeQuoter", Version: v16}},
		"0x00000000000000000000000000000000000000B2": {{Type: "FeeQuoter", Version: v16}},
	}
	_, err := LoadChainState(t.Context(), chain, addresses, WithMCMSQualifier(DefaultMCMSQualifier))
	require.ErrorContains(t, err, "ambiguous FeeQuoter 1.6.0")

	// different versions still select the highest (MultipleFeeQuoters semantics)
	addresses = map[string][]cldf.TypeAndVersion{
		"0x00000000000000000000000000000000000000A1": {{Type: "FeeQuoter", Version: v16}},
		"0x00000000000000000000000000000000000000B2": {{Type: "FeeQuoter", Version: deployment.Version1_2_0}},
	}
	state, err := LoadChainState(t.Context(), chain, addresses, WithMCMSQualifier(DefaultMCMSQualifier))
	require.NoError(t, err)
	require.Equal(t, common.HexToAddress("0x00000000000000000000000000000000000000A1"), state.FeeQuoter.Address())
}

func TestValidateSolanaTimelockConfig(t *testing.T) {
	program := solana.MustPublicKeyFromBase58("9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM")
	valid := mcmsolana.ContractAddress(program, mcmsolana.PDASeed{1, 2, 3})
	v16 := deployment.Version1_6_0

	// configured qualifier resolves its own bundle, ignoring the default-qualified decoy
	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(testRef(1, valid, mcmscontracts.RBACTimelock, &v16, "RMNMCMS")))
	require.NoError(t, ds.Addresses().Add(testRef(1, valid, mcmscontracts.ProposerManyChainMultisig, &v16, "RMNMCMS")))
	require.NoError(t, ds.Addresses().Add(testRef(1, "not-a-program-address", mcmscontracts.RBACTimelock, &v16, DefaultMCMSQualifier)))
	tc := &cldfproposalutils.TimelockConfig{TimelockQualifierPerChain: map[uint64]string{1: "RMNMCMS"}}
	require.NoError(t, ValidateSolanaTimelockConfig(cldf.Environment{DataStore: ds.Seal()}, 1, tc))

	// default qualifier, action-specific contract; empty action defaults in place
	// (framework validateCommon contract)
	ds2 := datastore.NewMemoryDataStore()
	require.NoError(t, ds2.Addresses().Add(testRef(1, valid, mcmscontracts.RBACTimelock, &v16, DefaultMCMSQualifier)))
	require.NoError(t, ds2.Addresses().Add(testRef(1, valid, mcmscontracts.BypasserManyChainMultisig, &v16, DefaultMCMSQualifier)))
	require.NoError(t, ds2.Addresses().Add(testRef(1, valid, mcmscontracts.ProposerManyChainMultisig, &v16, DefaultMCMSQualifier)))
	tc2 := &cldfproposalutils.TimelockConfig{MCMSAction: mcmstypes.TimelockActionBypass}
	require.NoError(t, ValidateSolanaTimelockConfig(cldf.Environment{DataStore: ds2.Seal()}, 1, tc2))

	tcDefault := &cldfproposalutils.TimelockConfig{}
	require.NoError(t, ValidateSolanaTimelockConfig(cldf.Environment{DataStore: ds2.Seal()}, 1, tcDefault))
	require.Equal(t, mcmstypes.TimelockActionSchedule, tcDefault.MCMSAction)

	// The default bundle may be stored without a qualifier by legacy deployments.
	dsEmptyDefault := datastore.NewMemoryDataStore()
	require.NoError(t, dsEmptyDefault.Addresses().Add(testRef(1, valid, mcmscontracts.RBACTimelock, &v16, "")))
	require.NoError(t, dsEmptyDefault.Addresses().Add(testRef(1, valid, mcmscontracts.ProposerManyChainMultisig, &v16, "")))
	require.NoError(t, ValidateSolanaTimelockConfig(cldf.Environment{DataStore: dsEmptyDefault.Seal()}, 1, &cldfproposalutils.TimelockConfig{}))

	// Custom qualifiers remain strict and do not fall back to the unqualified bundle.
	tcEmptyCustom := &cldfproposalutils.TimelockConfig{TimelockQualifierPerChain: map[uint64]string{1: "RMNMCMS"}}
	require.Error(t, ValidateSolanaTimelockConfig(cldf.Environment{DataStore: dsEmptyDefault.Seal()}, 1, tcEmptyCustom))

	// missing contract
	_, err := dataStoreSolanaContractAddress(cldf.Environment{DataStore: ds2.Seal()}, 1, mcmscontracts.CancellerManyChainMultisig, DefaultMCMSQualifier)
	require.ErrorContains(t, err, "no CancellerManyChainMultiSig ref")

	// superseded refs are invisible
	ds3 := datastore.NewMemoryDataStore()
	require.NoError(t, ds3.Addresses().Add(testRef(1, valid, mcmscontracts.RBACTimelock, &v16, DefaultMCMSQualifier, SupersededLabel)))
	_, err = dataStoreSolanaContractAddress(cldf.Environment{DataStore: ds3.Seal()}, 1, mcmscontracts.RBACTimelock, DefaultMCMSQualifier)
	require.ErrorContains(t, err, "no RBACTimelock ref")

	// nil config and invalid action
	require.Error(t, ValidateSolanaTimelockConfig(cldf.Environment{DataStore: ds2.Seal()}, 1, nil))
	tcBad := &cldfproposalutils.TimelockConfig{MCMSAction: "nope"}
	require.ErrorContains(t, ValidateSolanaTimelockConfig(cldf.Environment{DataStore: ds2.Seal()}, 1, tcBad), "invalid MCMS action")

	// configured qualifier that holds nothing errors (strict, no fallback)
	tc4 := &cldfproposalutils.TimelockConfig{TimelockQualifierPerChain: map[uint64]string{1: "OTHER"}}
	require.Error(t, ValidateSolanaTimelockConfig(cldf.Environment{DataStore: ds2.Seal()}, 1, tc4))
}
