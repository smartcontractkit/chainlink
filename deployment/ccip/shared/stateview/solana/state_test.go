package solana

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

func TestSelectSolanaRefsSelectsNewestVersion(t *testing.T) {
	v10 := semver.MustParse("1.0.0")
	v11 := semver.MustParse("1.1.0")
	refs := []datastore.AddressRef{
		{Address: "router-v1", Type: datastore.ContractType(shared.Router), Version: v10},
		{Address: "router-v2", Type: datastore.ContractType(shared.Router), Version: v11},
	}

	selected, err := selectSolanaRefs(refs)
	require.NoError(t, err)
	require.Len(t, selected, 1)
	require.Equal(t, "router-v2", selected[0].Address)
}

func TestSelectSolanaRefsRejectsEqualVersionAmbiguity(t *testing.T) {
	v10 := semver.MustParse("1.0.0")
	refs := []datastore.AddressRef{
		{Address: "router-a", Type: datastore.ContractType(shared.Router), Version: v10},
		{Address: "router-b", Type: datastore.ContractType(shared.Router), Version: v10},
	}

	_, err := selectSolanaRefs(refs)
	require.Error(t, err)
}

func TestSelectSolanaRefsTreatsReceiverSpellingsAsOneSlot(t *testing.T) {
	v10 := semver.MustParse("1.0.0")
	refs := []datastore.AddressRef{
		{Address: "receiver-a", Type: datastore.ContractType(shared.Receiver), Version: v10},
		{Address: "receiver-b", Type: datastore.ContractType(shared.TestReceiver), Version: v10},
	}

	_, err := selectSolanaRefs(refs)
	require.Error(t, err)
	require.ErrorContains(t, err, string(shared.Receiver))
}

func TestSelectSolanaRefsKeepsQualifiedInstances(t *testing.T) {
	v10 := semver.MustParse("1.0.0")
	refs := []datastore.AddressRef{
		{Address: "token-a", Type: datastore.ContractType(shared.SPLTokens), Version: v10, Qualifier: "a"},
		{Address: "token-b", Type: datastore.ContractType(shared.SPLTokens), Version: v10, Qualifier: "b"},
	}

	selected, err := selectSolanaRefs(refs)
	require.NoError(t, err)
	require.Len(t, selected, 2)
}

func TestLoadChainStateSolanaDecodesLookupTableQualifier(t *testing.T) {
	mint := solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	lookupTable := solana.MustPublicKeyFromBase58("11111111111111111111111111111111")
	metadata := "pool/with/slashes"
	v10 := semver.MustParse("1.0.0")

	state, err := LoadChainStateSolana(cldf_solana.Chain{Selector: 1}, []datastore.AddressRef{{
		ChainSelector: 1,
		Address:       lookupTable.String(),
		Type:          datastore.ContractType(shared.TokenPoolLookupTable),
		Version:       v10,
		Qualifier:     shared.TokenPoolLookupTableQualifier(mint.String(), shared.BurnMintTokenPool.String(), metadata),
	}})
	require.NoError(t, err)
	require.Equal(t, lookupTable, state.TokenPoolLookupTable[mint][shared.BurnMintTokenPool][metadata])
}
