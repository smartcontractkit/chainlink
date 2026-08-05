package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetFeedConfigsChangeset(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedCacheRef(t, &env, testContractID, "test", "1.0.0")

	req := &SetFeedConfigsRequest{
		ChainSel:     testChainSel,
		Qualifier:    "test",
		Version:      "1.0.0",
		Admin:        testAdmin,
		DataIDs:      []string{"0x018e16c39e00032000000"},
		Descriptions: []string{"BTC/USD"},
		Permissions: []FeedPermission{{
			AllowedSender:        testContractID, // the forwarder
			AllowedWorkflowOwner: "0x0102030405060708090a0b0c0d0e0f1011121314",
			AllowedWorkflowName:  "abc",
		}},
	}
	require.NoError(t, SetFeedConfigs{}.VerifyPreconditions(env, req))

	_, err := SetFeedConfigs{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "set_feed_configs", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)

	// Permissions apply to every feed in the batch: a single permission set
	// with two DataIDs must still pass length validation.
	multi := *req
	multi.DataIDs = []string{"0x01", "0x02"}
	multi.Descriptions = []string{"BTC/USD", "ETH/USD"}
	require.NoError(t, SetFeedConfigs{}.VerifyPreconditions(env, &multi))

	// mismatched lengths must fail preconditions
	bad := *req
	bad.Descriptions = nil
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &bad))

	// empty DataIDs must fail preconditions
	empty := *req
	empty.DataIDs = nil
	empty.Descriptions = nil
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &empty))

	// invalid admin address must fail preconditions
	badAdmin := *req
	badAdmin.Admin = "not-a-key"
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &badAdmin))

	// missing cache ref must fail preconditions
	badQualifier := *req
	badQualifier.Qualifier = "does-not-exist"
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &badQualifier))

	// invalid permission fields must fail preconditions
	badPerm := *req
	badPerm.Permissions = []FeedPermission{{
		AllowedSender:        "not-a-key",
		AllowedWorkflowOwner: "0x0102030405060708090a0b0c0d0e0f1011121314",
		AllowedWorkflowName:  "abc",
	}}
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &badPerm))

	// empty Permissions must fail preconditions
	emptyPerm := *req
	emptyPerm.Permissions = nil
	require.Error(t, SetFeedConfigs{}.VerifyPreconditions(env, &emptyPerm))
}
