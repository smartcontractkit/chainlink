package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFeedAdminChangesets(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedCacheRef(t, &env, testContractID, "test", "1.0.0")

	req := &FeedAdminRequest{
		ChainSel:  testChainSel,
		Qualifier: "test",
		Version:   "1.0.0",
		Admin:     testAdmin,
	}

	require.NoError(t, AddFeedAdmin{}.VerifyPreconditions(env, req))
	_, err := AddFeedAdmin{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "add_feed_admin", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)

	require.NoError(t, RemoveFeedAdmin{}.VerifyPreconditions(env, req))
	_, err = RemoveFeedAdmin{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 2)
	require.Equal(t, "remove_feed_admin", inv.calls[1].Function)
	require.Equal(t, testContractID, inv.calls[1].ContractID)

	// invalid admin address must fail preconditions for both changesets
	badAdmin := *req
	badAdmin.Admin = "not-a-key"
	require.Error(t, AddFeedAdmin{}.VerifyPreconditions(env, &badAdmin))
	require.Error(t, RemoveFeedAdmin{}.VerifyPreconditions(env, &badAdmin))

	// missing cache ref must fail preconditions for both changesets
	badQualifier := *req
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, AddFeedAdmin{}.VerifyPreconditions(env, &badQualifier))
	require.Error(t, RemoveFeedAdmin{}.VerifyPreconditions(env, &badQualifier))
}
