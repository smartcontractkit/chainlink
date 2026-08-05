package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetFeedConfigsChangeset(t *testing.T) {
	env, inv, _ := newTestEnv(t)
	seedCacheRef(t, &env, testContractID, "test", "1.0.0")
	seedContractMetadata(t, &env, testContractID, ContractMetadata{
		Feeds: map[string]FeedMetadata{
			"0x01c50f0e2106d5fd0000000000000000": {Description: "SOL/USD"},
		},
	})

	req := &SetFeedConfigsRequest{
		ChainSel:     testChainSel,
		Qualifier:    "test",
		Version:      "1.0.0",
		Admin:        testAdmin,
		DataIDs:      []string{"0x018e16c39e0003320000000000000000", "0x018e16c39e0003990000000000000000"},
		Descriptions: []string{"BTC/USD", "ETH/USD"},
		Permissions: []FeedPermission{{
			AllowedSender:        testContractID, // the forwarder
			AllowedWorkflowOwner: "0x0102030405060708090a0b0c0d0e0f1011121314",
			AllowedWorkflowName:  "abc",
		}},
	}
	require.NoError(t, SetFeedConfigs{}.VerifyPreconditions(env, req))

	out, err := SetFeedConfigs{}.Apply(env, req)
	require.NoError(t, err)
	require.Len(t, inv.calls, 1)
	require.Equal(t, "set_feed_configs", inv.calls[0].Function)
	require.Equal(t, testContractID, inv.calls[0].ContractID)

	// the metadata mirror keeps the pre-existing feed and adds both new ones,
	// each with its own description and id-derived decimals
	meta := outputMetadata(t, out, testContractID)
	require.Len(t, meta.Feeds, 3)
	require.Equal(t, "SOL/USD", meta.Feeds["0x01c50f0e2106d5fd0000000000000000"].Description)
	btc := meta.Feeds["0x018e16c39e0003320000000000000000"]
	require.Equal(t, "BTC/USD", btc.Description)
	require.Equal(t, uint32(18), btc.Decimals)
	require.Equal(t, req.Permissions, btc.Permissions)
	eth := meta.Feeds["0x018e16c39e0003990000000000000000"]
	require.Equal(t, "ETH/USD", eth.Description)
	require.Equal(t, uint32(0), eth.Decimals) // byte 7 outside the decimals range

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
