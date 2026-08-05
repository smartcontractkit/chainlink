package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveFeedConfigsChangeset(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedCacheRef(t, &env, testContractID, "test", "1.0.0")
	seedContractMetadata(t, &env, testContractID, ContractMetadata{
		Feeds: map[string]FeedMetadata{
			"0x018e16c39e0003320000000000000000": {Description: "BTC/USD"},
			"0x01c50f0e2106d5fd0000000000000000": {Description: "SOL/USD"},
		},
	})

	req := &RemoveFeedConfigsRequest{
		ChainSel:  testChainSel,
		Qualifier: "test",
		Version:   "1.0.0",
		Admin:     testAdmin,
		DataIDs:   []string{"0x018e16c39e0003320000000000000000"},
	}
	require.NoError(t, RemoveFeedConfigs{}.VerifyPreconditions(env, req))

	out, err := RemoveFeedConfigs{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "remove_feed_configs", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)

	// only the removed feed leaves the metadata mirror
	meta := outputMetadata(t, out, testContractID)
	require.NotContains(t, meta.Feeds, "0x018e16c39e0003320000000000000000")
	require.Equal(t, "SOL/USD", meta.Feeds["0x01c50f0e2106d5fd0000000000000000"].Description)

	// empty DataIDs must fail preconditions
	empty := *req
	empty.DataIDs = nil
	require.Error(t, RemoveFeedConfigs{}.VerifyPreconditions(env, &empty))

	// invalid admin address must fail preconditions
	badAdmin := *req
	badAdmin.Admin = "not-a-key"
	require.Error(t, RemoveFeedConfigs{}.VerifyPreconditions(env, &badAdmin))

	// missing cache ref must fail preconditions
	badQualifier := *req
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, RemoveFeedConfigs{}.VerifyPreconditions(env, &badQualifier))
}
