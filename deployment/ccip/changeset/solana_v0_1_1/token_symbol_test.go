package solana

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

func randomMint(t *testing.T) solana.PublicKey {
	t.Helper()
	key, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	return key.PublicKey()
}

func TestValidateTokenSymbol(t *testing.T) {
	t.Parallel()

	mint := randomMint(t).String()

	require.NoError(t, validateTokenSymbol("TEST", mint))
	require.Error(t, validateTokenSymbol("", mint))
	require.Error(t, validateTokenSymbol(mint+"-token", mint))
}

func TestRecordTokenMultisig(t *testing.T) {
	t.Parallel()

	chainSelector := chain_selectors.SOLANA_TESTNET.Selector
	mint := randomMint(t)
	multisig := randomMint(t).String()
	ab := cldf.NewMemoryAddressBook()
	ds := datastore.NewMemoryDataStore()

	require.NoError(t, recordTokenMultisig(ab, ds, chainSelector, multisig, "customer", mint, "TEST_TOKEN"))

	addresses, err := ab.Addresses() //nolint:staticcheck
	require.NoError(t, err)
	legacyTV := addresses[chainSelector][multisig]
	require.Equal(t, cldf.ContractType("TokenMultisig"), legacyTV.Type)
	require.Equal(t, deployment.Version1_0_0, legacyTV.Version)

	refs, err := ds.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, refs, 1)
	require.Equal(t, datastore.ContractType("TOKEN_MULTISIG"), refs[0].Type)
	require.Equal(t, deployment.Version1_6_0, *refs[0].Version)
	require.Equal(t, "TEST_TOKEN", refs[0].Qualifier)
	require.Equal(t, datastore.NewLabelSet(mint.String()), refs[0].Labels)
}

func TestRecordOnboardedTokenMint(t *testing.T) {
	t.Parallel()

	solanaTestnet := chain_selectors.SOLANA_TESTNET.Selector

	newConfig := func(mint solana.PublicKey, symbol string) OnboardTokenPoolConfig {
		return OnboardTokenPoolConfig{
			TokenMint:        mint,
			TokenProgramName: shared.SPLTokens,
			PoolType:         shared.BurnMintTokenPool,
			Metadata:         shared.CLLMetadata,
			TokenSymbol:      symbol,
		}
	}
	refsByQualifier := func(ds datastore.MutableDataStore) map[string]datastore.AddressRef {
		refs, err := ds.Addresses().Fetch()
		require.NoError(t, err)
		out := make(map[string]datastore.AddressRef, len(refs))
		for _, ref := range refs {
			out[ref.Qualifier] = ref
		}

		return out
	}

	t.Run("two tokens sharing metadata and pool type get distinct symbol keys", func(t *testing.T) {
		t.Parallel()
		ab := cldf.NewMemoryAddressBook()
		ds := datastore.NewMemoryDataStore()
		lnr := newConfig(randomMint(t), "LnR")
		bnm := newConfig(randomMint(t), "BnM")

		require.NoError(t, recordOnboardedTokenMint(solanaTestnet, ab, ds, nil, lnr, "PoolProgram"))
		require.NoError(t, recordOnboardedTokenMint(solanaTestnet, ab, ds, nil, bnm, "PoolProgram"))

		refs := refsByQualifier(ds)
		require.Len(t, refs, 2)
		require.Equal(t, lnr.TokenMint.String(), refs["LnR"].Address)
		require.Equal(t, bnm.TokenMint.String(), refs["BnM"].Address)

		addresses, err := ab.Addresses() //nolint:staticcheck
		require.NoError(t, err)
		require.Len(t, addresses[solanaTestnet], 2)
	})

	t.Run("a known mint reuses its address-book metadata", func(t *testing.T) {
		t.Parallel()
		ab := cldf.NewMemoryAddressBook()
		ds := datastore.NewMemoryDataStore()
		cfg := newConfig(randomMint(t), "LnR")
		envAddresses := map[string]cldf.TypeAndVersion{
			cfg.TokenMint.String(): cldf.NewTypeAndVersion(shared.SPLTokens, deployment.Version1_0_0),
		}

		require.NoError(t, recordOnboardedTokenMint(solanaTestnet, ab, ds, envAddresses, cfg, "PoolProgram"))

		refs := refsByQualifier(ds)
		require.Len(t, refs, 1, "the datastore ref is still backfilled")
		require.Equal(t, "LnR", refs["LnR"].Qualifier)

		addresses, err := ab.Addresses() //nolint:staticcheck
		require.NoError(t, err)
		require.Len(t, addresses[solanaTestnet], 1)
		require.Equal(t, envAddresses[cfg.TokenMint.String()], addresses[solanaTestnet][cfg.TokenMint.String()])
	})

	t.Run("re-recording is idempotent", func(t *testing.T) {
		t.Parallel()
		ab := cldf.NewMemoryAddressBook()
		ds := datastore.NewMemoryDataStore()
		cfg := newConfig(randomMint(t), "LnR")

		require.NoError(t, recordOnboardedTokenMint(solanaTestnet, ab, ds, nil, cfg, "PoolProgram"))
		require.NoError(t, recordOnboardedTokenMint(solanaTestnet, ab, ds, nil, cfg, "PoolProgram"))

		require.Len(t, refsByQualifier(ds), 1)
	})

	t.Run("rejects an empty symbol and a symbol embedding the mint", func(t *testing.T) {
		t.Parallel()
		mint := randomMint(t)
		ab := cldf.NewMemoryAddressBook()
		ds := datastore.NewMemoryDataStore()

		empty := newConfig(mint, "")
		require.Error(t, recordOnboardedTokenMint(solanaTestnet, ab, ds, nil, empty, "PoolProgram"))

		addressDerived := newConfig(mint, mint.String()+"-token")
		require.Error(t, recordOnboardedTokenMint(solanaTestnet, ab, ds, nil, addressDerived, "PoolProgram"))

		refs, err := ds.Addresses().Fetch()
		require.NoError(t, err)
		require.Empty(t, refs)
	})
}
