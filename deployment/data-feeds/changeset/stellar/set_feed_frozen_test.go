package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetFeedFrozenChangeset(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedCacheRef(t, &env, testContractID, "test", "1.0.0")

	req := &SetFeedFrozenRequest{
		ChainSel:  testChainSel,
		Qualifier: "test",
		Version:   "1.0.0",
		Admin:     testAdmin,
		DataIDs:   []string{"0x018e16c39e00032000000"},
		Frozen:    true,
	}
	require.NoError(t, SetFeedFrozen{}.VerifyPreconditions(env, req))

	_, err := SetFeedFrozen{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "set_feed_frozen", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)

	// empty DataIDs must fail preconditions
	empty := *req
	empty.DataIDs = nil
	require.Error(t, SetFeedFrozen{}.VerifyPreconditions(env, &empty))

	// invalid admin address must fail preconditions
	badAdmin := *req
	badAdmin.Admin = "not-a-key"
	require.Error(t, SetFeedFrozen{}.VerifyPreconditions(env, &badAdmin))

	// missing cache ref must fail preconditions
	badQualifier := *req
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, SetFeedFrozen{}.VerifyPreconditions(env, &badQualifier))
}
