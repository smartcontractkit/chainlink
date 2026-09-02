package shared

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// TestPopulateDataStoreQualifiersDryRun exercises the address-book -> datastore path
// (PopulateDataStore with explicit caller-supplied qualifiers) against representative data, so
// we can see exactly what qualifiers the resulting AddressRefs carry before any on-chain
// deployment. It is the offline "dry run" for the qualifier convention: the caller already
// knows the qualifier for every multi-instance address it deploys and passes it in.
func TestPopulateDataStoreQualifiersDryRun(t *testing.T) {
	ab := cldf.NewMemoryAddressBook()
	sel := uint64(5009297550715157269) // ethereum mainnet

	// Singletons: empty qualifier.
	mustSave(t, ab, sel, "0x1111111111111111111111111111111111111111", "Router", "1.6.0")
	mustSave(t, ab, sel, "0x2222222222222222222222222222222222222222", "FeeQuoter", "1.6.0")
	mustSave(t, ab, sel, "0x3333333333333333333333333333333333333333", "NonceManager", "1.6.0")

	// Multi-instance: two token pools for two different tokens on the same chain.
	mustSave(t, ab, sel, "0x4444444444444444444444444444444444444444", "BurnMintTokenPool", "1.6.0")
	mustSave(t, ab, sel, "0x5555555555555555555555555555555555555555", "LockReleaseTokenPool", "1.6.0")

	// The caller knows each multi-instance contract's qualifier (the token symbol).
	qualifiers := map[AddressKey]string{
		{ChainSelector: sel, Address: "0x4444444444444444444444444444444444444444"}: "LINK",
		{ChainSelector: sel, Address: "0x5555555555555555555555555555555555555555"}: "USDC",
	}
	ds, err := PopulateDataStore(ab, qualifiers)
	require.NoError(t, err)

	refs, err := ds.Addresses().Fetch()
	require.NoError(t, err)

	// Singletons must be present with the empty qualifier and findable.
	byKey := func(tp string) datastore.AddressRef {
		for _, r := range refs {
			if string(r.Type) == tp {
				return r
			}
		}
		return datastore.AddressRef{}
	}
	for _, tp := range []string{"Router", "FeeQuoter", "NonceManager"} {
		r := byKey(tp)
		require.NotEmpty(t, r.Address, "singleton %s missing from datastore", tp)
		require.Empty(t, r.Qualifier, "singleton %s must have empty qualifier", tp)
	}

	// Both pools are qualified with the token symbol the caller supplied and are
	// distinguishable by their qualifier.
	gotQualifier := map[string]string{}
	for _, r := range refs {
		switch string(r.Type) {
		case "BurnMintTokenPool":
			gotQualifier["BurnMintTokenPool"] = r.Qualifier
		case "LockReleaseTokenPool":
			gotQualifier["LockReleaseTokenPool"] = r.Qualifier
		}
	}
	require.Equal(t, "LINK", gotQualifier["BurnMintTokenPool"])
	require.Equal(t, "USDC", gotQualifier["LockReleaseTokenPool"])

	// Dry-run output: the exact AddressRefs the datastore would carry.
	for _, r := range refs {
		t.Logf("DRY-RUN ref: chain=%d type=%s version=%s qualifier=%q address=%s", r.ChainSelector, r.Type, versionStr(r.Version), r.Qualifier, r.Address)
	}
}

func mustSave(t *testing.T, ab cldf.AddressBook, sel uint64, addr, tp, ver string) {
	t.Helper()
	tv := cldf.NewTypeAndVersion(cldf.ContractType(tp), *semver.MustParse(ver))
	require.NoError(t, ab.Save(sel, addr, tv))
}

func versionStr(v *semver.Version) string {
	if v == nil {
		return ""
	}
	return v.String()
}

// TestPopulateDataStoreRejectsUnqualifiedDuplicate verifies that PopulateDataStore fails loudly
// when two refs map to the same datastore key (chain, type, version, empty qualifier), rather
// than silently keeping an arbitrary one.
func TestPopulateDataStoreRejectsUnqualifiedDuplicate(t *testing.T) {
	sel := uint64(5009297550715157269)
	ver := "1.6.0"

	ab := cldf.NewMemoryAddressBook()
	// Two different addresses, same chain/type/version, no qualifier provided → same datastore key.
	require.NoError(t, ab.Save(sel, "0x1111111111111111111111111111111111111111", cldf.NewTypeAndVersion("BurnMintTokenPool", *semver.MustParse(ver))))
	require.NoError(t, ab.Save(sel, "0x2222222222222222222222222222222222222222", cldf.NewTypeAndVersion("BurnMintTokenPool", *semver.MustParse(ver))))

	_, err := PopulateDataStore(ab, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "BurnMintTokenPool")
}

// TestTokenPoolLookupTableQualifier verifies the composite LUT qualifier is unique per
// (token mint, pool type, metadata).
func TestTokenPoolLookupTableQualifier(t *testing.T) {
	mint := "So11111111111111111111111111111111111111112"
	a := TokenPoolLookupTableQualifier(mint, "BurnMintTokenPool", "CLL")
	b := TokenPoolLookupTableQualifier(mint, "LockReleaseTokenPool", "CLL")
	c := TokenPoolLookupTableQualifier(mint, "BurnMintTokenPool", "customPool9")
	require.NotEqual(t, a, b)
	require.NotEqual(t, a, c)
	require.NotEqual(t, b, c)
}
