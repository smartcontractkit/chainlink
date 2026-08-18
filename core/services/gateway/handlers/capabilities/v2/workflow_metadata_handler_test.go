package v2

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	testWorkflowID1      = "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"
	testWorkflowID2      = "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	testWorkflowNameHex1 = "0x1234567890abcdef1234"
	testWorkflowNameHex2 = "0xabcdef1234567890abcd"
	testWorkflowOwner1   = "0x1234567890abcdef1234567890abcdef12345678"
	testWorkflowOwner2   = "0xabcdef1234567890abcdef1234567890abcdef12"
	testWorkflowTag1     = "workflowTag1"
	testWorkflowTag2     = "workflowTag2"
	testPublicKey1       = "0x1234567890abcdef1234567890abcdef12345678"
	testPublicKey2       = "0xabcdef1234567890abcdef1234567890abcdef12"
	testPublicKey3       = "0xabcdef1234567890abcdef1234567890abcdefab"
)

func createTestWorkflowMetadataHandler(t *testing.T) (*WorkflowMetadataHandler, *mocks.DON, *config.DONConfig) {
	lggr := logger.Test(t)
	mockDon := mocks.NewDON(t)

	donConfig := &config.DONConfig{
		F: 1,
		Members: []config.NodeConfig{
			{Address: "node1"},
			{Address: "node2"},
			{Address: "node3"},
		},
	}

	cfg := WithDefaults(ServiceConfig{})
	testMetrics, err := metrics.NewMetrics(donConfig.Members)
	require.NoError(t, err)
	shards, nodeAddrToShard := singleShardEndpoints(t, donConfig, mockDon)
	handler := NewWorkflowMetadataHandler(lggr, cfg, shards, nodeAddrToShard, testMetrics)
	return handler, mockDon, donConfig
}

// singleShardEndpoints builds a one-DON one-shard endpoint matrix from a legacy
// DONConfig for tests.
func singleShardEndpoints(t *testing.T, donConfig *config.DONConfig, mockDon *mocks.DON) ([]*shardEndpoint, map[string]*shardEndpoint) {
	shardedDONs := []config.ShardedDONConfig{
		{DonName: donConfig.DonId, F: donConfig.F, Shards: []config.Shard{{Nodes: donConfig.Members}}},
	}
	shards, nodeAddrToShard, err := buildShardEndpoints(shardedDONs, [][]handlers.DON{{mockDon}})
	require.NoError(t, err)
	return shards, nodeAddrToShard
}

func TestSyncMetadata(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)

	// Test when aggregator has no data
	handler.syncMetadata(t.Context())
	require.Empty(t, handler.authorizedKeys)

	// Start the aggregator to enable data collection
	ctx := t.Context()
	err := handler.aggs[handler.shards[0].donID].Start(ctx)
	require.NoError(t, err)
	defer handler.aggs[handler.shards[0].donID].Close()

	// Add some test data to aggregator
	key := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: "key1",
	}
	observation := gateway_common.WorkflowMetadata{
		WorkflowSelector: gateway_common.WorkflowSelector{
			WorkflowID:    testWorkflowID1,
			WorkflowName:  testWorkflowNameHex1,
			WorkflowOwner: testWorkflowOwner1,
			WorkflowTag:   testWorkflowTag1,
		},
		AuthorizedKeys: []gateway_common.AuthorizedKey{key},
	}

	// Collect enough observations to meet threshold (F+1 = 2)
	err = handler.aggs[handler.shards[0].donID].Collect(&observation, "node1")
	require.NoError(t, err)
	err = handler.aggs[handler.shards[0].donID].Collect(&observation, "node2")
	require.NoError(t, err)
	handler.syncMetadata(t.Context())

	workflowKeys, exists := handler.authorizedKeys[testWorkflowID1]
	require.True(t, exists)
	_, exists = workflowKeys[key]
	require.True(t, exists)
	require.Len(t, workflowKeys, 1)
	ref, exists := handler.workflowIDToRef[testWorkflowID1]
	require.True(t, exists)
	expectedRef := workflowReference{
		workflowName:  testWorkflowNameHex1,
		workflowOwner: testWorkflowOwner1,
		workflowTag:   testWorkflowTag1,
	}
	require.Equal(t, expectedRef, ref)
	workflowID, exists := handler.workflowRefToID[expectedRef]
	require.True(t, exists)
	require.Equal(t, testWorkflowID1, workflowID)
}

func TestSyncMetadataMultipleWorkflows(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)

	ctx := t.Context()
	err := handler.aggs[handler.shards[0].donID].Start(ctx)
	require.NoError(t, err)
	defer handler.aggs[handler.shards[0].donID].Close()

	// Add observations for multiple workflows
	workflows := []string{"workflow1", "workflow2"}
	keys := []string{"key1", "key2", "key3"}

	for _, workflowID := range workflows {
		for _, key := range keys {
			observation := gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    workflowID,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: key,
					},
				},
			}
			err = handler.aggs[handler.shards[0].donID].Collect(&observation, "node1")
			require.NoError(t, err)
			err = handler.aggs[handler.shards[0].donID].Collect(&observation, "node2")
			require.NoError(t, err)
		}
	}
	handler.syncMetadata(t.Context())

	expectedRef := workflowReference{
		workflowName:  testWorkflowNameHex1,
		workflowOwner: testWorkflowOwner1,
		workflowTag:   testWorkflowTag1,
	}
	// Both "workflow1" and "workflow2" share the same workflow reference. The
	// workflow reference is unique (enforced by the on-chain registry), so
	// syncMetadata deduplicates by reference: only the first observed workflow ID
	// wins. Aggregate() returns observations newest-first, so "workflow2" (the
	// last collected, highest sequence number) is processed first and wins the
	// reference; "workflow1" is dropped as a duplicate reference.
	require.Len(t, handler.authorizedKeys, 1)
	require.Contains(t, handler.authorizedKeys, "workflow2")
	workflowKeys := handler.authorizedKeys["workflow2"]
	require.Len(t, workflowKeys, 1)

	ref, exists := handler.workflowIDToRef["workflow2"]
	require.True(t, exists)
	require.Equal(t, expectedRef, ref)
	winningWorkflowID, exists := handler.workflowRefToID[expectedRef]
	require.True(t, exists)
	require.Equal(t, "workflow2", winningWorkflowID)
	// "workflow1" must have been dropped by the reference dedup.
	_, exists = handler.workflowIDToRef["workflow1"]
	require.False(t, exists)
}

func TestSendMetadataPullRequest(t *testing.T) {
	handler, mockDon, donConfig := createTestWorkflowMetadataHandler(t)
	for _, member := range donConfig.Members {
		mockDon.EXPECT().SendToNode(mock.Anything, member.Address, mock.Anything).Return(nil).Once()
	}

	err := handler.sendMetadataPullRequest()
	require.NoError(t, err)
	mockDon.AssertExpectations(t)
}

func TestSendMetadataPullRequestWithErrors(t *testing.T) {
	handler, mockDon, donConfig := createTestWorkflowMetadataHandler(t)

	// Mock errors for some nodes
	expectedErrors := []error{
		errors.New("connection failed"),
		nil,
		errors.New("timeout"),
	}

	for i, member := range donConfig.Members {
		mockDon.EXPECT().SendToNode(mock.Anything, member.Address, mock.Anything).Return(expectedErrors[i]).Once()
	}

	err := handler.sendMetadataPullRequest()
	require.Error(t, err)
	require.Contains(t, err.Error(), "connection failed")
	require.Contains(t, err.Error(), "timeout")
	require.NotContains(t, err.Error(), "node2")
	mockDon.AssertExpectations(t)
}

func TestSendMetadataPullRequestVerifyPayload(t *testing.T) {
	handler, mockDon, donConfig := createTestWorkflowMetadataHandler(t)
	// Capture the request payload
	var capturedReq *jsonrpc.Request[json.RawMessage]
	mockDon.On("SendToNode", mock.Anything, mock.AnythingOfType("string"), mock.Anything).
		Run(func(args mock.Arguments) {
			capturedReq = args.Get(2).(*jsonrpc.Request[json.RawMessage])
		}).Return(nil)

	err := handler.sendMetadataPullRequest()
	require.NoError(t, err)

	require.Equal(t, jsonrpc.JsonRpcVersion, capturedReq.Version)
	require.Equal(t, gateway_common.MethodPullWorkflowMetadata, capturedReq.Method)
	require.NotEmpty(t, capturedReq.ID)

	mockDon.AssertNumberOfCalls(t, "SendToNode", len(donConfig.Members))
}

// createMultiShardMetadataHandler builds a metadata handler over two shards of
// one DON: shard 0 = {node1,node2}, shard 1 = {node3,node4}, F=0 so threshold=1.
func createMultiShardMetadataHandler(t *testing.T) (*WorkflowMetadataHandler, *mocks.DON, *mocks.DON) {
	lggr := logger.Test(t)
	mockDon0 := mocks.NewDON(t)
	mockDon1 := mocks.NewDON(t)
	shardedDONs := []config.ShardedDONConfig{
		{DonName: "don", F: 0, Shards: []config.Shard{
			{Nodes: []config.NodeConfig{{Address: "node1"}, {Address: "node2"}}},
			{Nodes: []config.NodeConfig{{Address: "node3"}, {Address: "node4"}}},
		}},
	}
	shards, nodeAddrToShard, err := buildShardEndpoints(shardedDONs, [][]handlers.DON{{mockDon0, mockDon1}})
	require.NoError(t, err)
	testMetrics, err := metrics.NewMetrics(allMembers(shards))
	require.NoError(t, err)
	cfg := WithDefaults(ServiceConfig{})
	return NewWorkflowMetadataHandler(lggr, cfg, shards, nodeAddrToShard, testMetrics), mockDon0, mockDon1
}

func TestMultiShardAssignment(t *testing.T) {
	t.Parallel()

	handler, _, _ := createMultiShardMetadataHandler(t)
	ctx := t.Context()
	for _, shard := range handler.shards {
		require.NoError(t, handler.aggs[shard.donID].Start(ctx))
		t.Cleanup(func() { _ = handler.aggs[shard.donID].Close() })
	}

	obs := gateway_common.WorkflowMetadata{
		WorkflowSelector: gateway_common.WorkflowSelector{
			WorkflowID:    testWorkflowID1,
			WorkflowName:  testWorkflowNameHex1,
			WorkflowOwner: testWorkflowOwner1,
			WorkflowTag:   testWorkflowTag1,
		},
		AuthorizedKeys: []gateway_common.AuthorizedKey{{KeyType: gateway_common.KeyTypeECDSAEVM, PublicKey: testPublicKey1}},
	}

	// Initially unassigned.
	require.Empty(t, handler.WorkflowShards(testWorkflowID1))

	// Quorum (F+1=1) reached only in shard 1 (node3).
	require.NoError(t, handler.aggs[handler.shards[1].donID].Collect(&obs, "node3"))
	handler.syncMetadata(ctx)
	assigned := handler.WorkflowShards(testWorkflowID1)
	require.Len(t, assigned, 1)
	require.Equal(t, handler.shards[1].donID, assigned[0].donID)

	// Later quorum reached in shard 0 too (node1) => assigned to both shards.
	require.NoError(t, handler.aggs[handler.shards[0].donID].Collect(&obs, "node1"))
	handler.syncMetadata(ctx)
	assigned = handler.WorkflowShards(testWorkflowID1)
	require.Len(t, assigned, 2)

	// An observation from an unknown node is rejected.
	_, err := handler.aggForNode("unknown-node")
	require.Error(t, err)
}

func TestSendMetadataPullRequestMultiShard(t *testing.T) {
	t.Parallel()

	handler, mockDon0, mockDon1 := createMultiShardMetadataHandler(t)
	for _, addr := range []string{"node1", "node2"} {
		mockDon0.EXPECT().SendToNode(mock.Anything, addr, mock.Anything).Return(nil).Once()
	}
	for _, addr := range []string{"node3", "node4"} {
		mockDon1.EXPECT().SendToNode(mock.Anything, addr, mock.Anything).Return(nil).Once()
	}
	require.NoError(t, handler.sendMetadataPullRequest())
	mockDon0.AssertExpectations(t)
	mockDon1.AssertExpectations(t)
}

func TestOnMetadataPush(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)
	ctx := t.Context()

	err := handler.aggs[handler.shards[0].donID].Start(ctx)
	require.NoError(t, err)
	defer handler.aggs[handler.shards[0].donID].Close()

	metadata := gateway_common.WorkflowMetadata{
		WorkflowSelector: gateway_common.WorkflowSelector{
			WorkflowID:    testWorkflowID1,
			WorkflowName:  testWorkflowNameHex1,
			WorkflowOwner: testWorkflowOwner1,
			WorkflowTag:   testWorkflowTag1,
		},
		AuthorizedKeys: []gateway_common.AuthorizedKey{
			{
				KeyType:   gateway_common.KeyTypeECDSAEVM,
				PublicKey: testWorkflowOwner1,
			},
			{
				KeyType:   gateway_common.KeyTypeECDSAEVM,
				PublicKey: testWorkflowOwner2,
			},
		},
	}

	result, err := json.Marshal(metadata)
	require.NoError(t, err)

	rawResult := json.RawMessage(result)
	resp := &jsonrpc.Response[json.RawMessage]{
		Result: &rawResult,
	}

	err = handler.OnMetadataPush(ctx, resp, "node1")
	require.NoError(t, err)

	handler.syncMetadata(t.Context())
	require.Empty(t, handler.authorizedKeys)
	require.Empty(t, handler.workflowIDToRef)
	require.Empty(t, handler.workflowRefToID)
}

func TestOnMetadataPushInvalidJSON(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)
	ctx := t.Context()

	invalidJSON := json.RawMessage(`{"invalid": json}`)
	resp := &jsonrpc.Response[json.RawMessage]{
		Result: &invalidJSON,
	}

	err := handler.OnMetadataPush(ctx, resp, "node1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to unmarshal metadata")
}

func TestOnMetadataPullResponse(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)
	ctx := t.Context()

	err := handler.aggs[handler.shards[0].donID].Start(ctx)
	require.NoError(t, err)
	defer handler.aggs[handler.shards[0].donID].Close()

	key1 := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: testPublicKey1,
	}
	key2 := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: testPublicKey2,
	}
	key3 := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: testPublicKey3,
	}
	metadata := []gateway_common.WorkflowMetadata{
		{
			WorkflowSelector: gateway_common.WorkflowSelector{
				WorkflowID:    testWorkflowID1,
				WorkflowName:  testWorkflowNameHex1,
				WorkflowOwner: testWorkflowOwner1,
				WorkflowTag:   testWorkflowTag1,
			},
			AuthorizedKeys: []gateway_common.AuthorizedKey{key1},
		},
		{
			WorkflowSelector: gateway_common.WorkflowSelector{
				WorkflowID:    testWorkflowID2,
				WorkflowName:  testWorkflowNameHex2,
				WorkflowOwner: testWorkflowOwner2,
				WorkflowTag:   testWorkflowTag2,
			},
			AuthorizedKeys: []gateway_common.AuthorizedKey{key2, key3},
		},
	}

	result, err := json.Marshal(metadata)
	require.NoError(t, err)

	rawResult := json.RawMessage(result)
	resp := &jsonrpc.Response[json.RawMessage]{
		Result: &rawResult,
	}

	err = handler.OnMetadataPullResponse(ctx, resp, "node1")
	require.NoError(t, err)
	handler.syncMetadata(t.Context())
	require.Empty(t, handler.authorizedKeys)
	require.Empty(t, handler.workflowIDToRef)
	require.Empty(t, handler.workflowRefToID)

	// node2 responds with the same payload so observations should be aggregated because f=1
	err = handler.OnMetadataPullResponse(ctx, resp, "node2")
	require.NoError(t, err)
	handler.syncMetadata(t.Context())
	require.Len(t, handler.authorizedKeys, 2)
	keys, exists := handler.authorizedKeys[testWorkflowID1]
	require.True(t, exists)
	require.Len(t, keys, 1)
	_, exists = keys[key1]
	require.True(t, exists)
	keys, exists = handler.authorizedKeys[testWorkflowID2]
	require.True(t, exists)
	require.Len(t, keys, 2)
	_, exists = keys[key2]
	require.True(t, exists)
	_, exists = keys[key3]
	require.True(t, exists)
	ref1 := workflowReference{
		workflowOwner: testWorkflowOwner1,
		workflowName:  testWorkflowNameHex1,
		workflowTag:   testWorkflowTag1,
	}
	ref2 := workflowReference{
		workflowName:  testWorkflowNameHex2,
		workflowOwner: testWorkflowOwner2,
		workflowTag:   testWorkflowTag2,
	}
	id, exists := handler.workflowRefToID[ref1]
	require.True(t, exists)
	require.Equal(t, testWorkflowID1, id)
	id, exists = handler.workflowRefToID[ref2]
	require.True(t, exists)
	require.Equal(t, testWorkflowID2, id)
	r1, exists := handler.workflowIDToRef[testWorkflowID1]
	require.True(t, exists)
	require.Equal(t, ref1, r1)
	r2, exists := handler.workflowIDToRef[testWorkflowID2]
	require.True(t, exists)
	require.Equal(t, ref2, r2)
}

func TestOnMetadataPullResponseInvalidJSON(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)
	ctx := t.Context()

	invalidJSON := json.RawMessage(`[{"invalid": json}]`)
	resp := &jsonrpc.Response[json.RawMessage]{
		Result: &invalidJSON,
	}

	err := handler.OnMetadataPullResponse(ctx, resp, "node1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to unmarshal metadata pull response")
}

func TestStartAndClose(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)
	ctx := t.Context()

	err := handler.Start(ctx)
	require.NoError(t, err)
	require.NoError(t, handler.Ready())
	err = handler.Start(ctx) // Should error on second start
	require.Error(t, err)

	err = handler.Close()
	require.NoError(t, err)
	require.Error(t, handler.Ready())
	err = handler.Close() // Should error on second close
	require.Error(t, err)
}

func TestValidateAuthMetadata(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)

	tests := []struct {
		name        string
		metadata    gateway_common.WorkflowMetadata
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid metadata",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: false,
		},
		{
			name: "empty workflow ID",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    "",
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid workflow ID",
		},
		{
			name: "empty workflow name",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  "",
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid workflow name",
		},
		{
			name: "empty workflow owner",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: "",
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid workflow owner",
		},
		{
			name: "empty workflow tag",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   "",
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid workflow tag",
		},
		{
			name: "no authorized keys",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{},
			},
			expectError: true,
			errorMsg:    "no authorized keys",
		},
		{
			name: "invalid key type",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   "invalid",
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid key type",
		},
		{
			name: "empty public key",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: "",
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid public key:",
		},
		{
			name: "public key without 0x prefix",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: "1234567890abcdef1234567890abcdef12345678",
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid public key:",
		},
		{
			name: "public key too short",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: "0x123456789",
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid public key:",
		},
		{
			name: "public key too long",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: "0x1234567890abcdef1234567890abcdef123456789",
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid public key:",
		},
		{
			name: "public key not lowercase",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: "0x1234567890ABCDEF1234567890abcdef12345678",
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid public key: must be all lowercase",
		},
		{
			name: "multiple valid keys",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey2,
					},
				},
			},
			expectError: false,
		},
		{
			name: "workflow ID too short",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef12345", // 65 chars
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid workflow ID",
		},
		{
			name: "workflow ID too long",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef1234567", // 67 chars
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid workflow ID",
		},
		{
			name: "workflow owner too short",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: "0x1234567890abcdef1234567890abcdef1234567", // 41 chars
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid workflow owner",
		},
		{
			name: "workflow owner too long",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: "0x1234567890abcdef1234567890abcdef123456789", // 43 chars
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid workflow owner",
		},
		{
			name: "workflow name too short",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  "0x1234567890abcdef123", // 21 chars
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid workflow name",
		},
		{
			name: "workflow name too long",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  "0x1234567890abcdef12345", // 23 chars
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid workflow name",
		},
		{
			name: "workflow tag too long",
			metadata: gateway_common.WorkflowMetadata{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   "this_is_a_very_long_workflow_tag_that_exceeds_the_maximum_length", // 65 chars
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testPublicKey1,
					},
				},
			},
			expectError: true,
			errorMsg:    "invalid workflow tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateAuthMetadata(tt.metadata)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOnMetadataPushWithValidation(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)
	ctx := t.Context()

	err := handler.aggs[handler.shards[0].donID].Start(ctx)
	require.NoError(t, err)
	defer handler.aggs[handler.shards[0].donID].Close()

	t.Run("valid metadata passes validation", func(t *testing.T) {
		metadata := gateway_common.WorkflowMetadata{
			WorkflowSelector: gateway_common.WorkflowSelector{
				WorkflowID:    testWorkflowID1,
				WorkflowName:  testWorkflowNameHex1,
				WorkflowOwner: testWorkflowOwner1,
				WorkflowTag:   testWorkflowTag1,
			},
			AuthorizedKeys: []gateway_common.AuthorizedKey{
				{
					KeyType:   gateway_common.KeyTypeECDSAEVM,
					PublicKey: testPublicKey1,
				},
			},
		}

		result, err := json.Marshal(metadata)
		require.NoError(t, err)

		rawResult := json.RawMessage(result)
		resp := &jsonrpc.Response[json.RawMessage]{
			Result: &rawResult,
		}

		err = handler.OnMetadataPush(ctx, resp, "node1")
		require.NoError(t, err)
	})

	t.Run("invalid metadata fails validation", func(t *testing.T) {
		metadata := gateway_common.WorkflowMetadata{
			WorkflowSelector: gateway_common.WorkflowSelector{
				WorkflowID:    "", // Invalid: empty workflow ID
				WorkflowName:  testWorkflowNameHex1,
				WorkflowOwner: testWorkflowOwner1,
				WorkflowTag:   testWorkflowTag1,
			},
			AuthorizedKeys: []gateway_common.AuthorizedKey{
				{
					KeyType:   gateway_common.KeyTypeECDSAEVM,
					PublicKey: testWorkflowOwner1,
				},
			},
		}

		result, err := json.Marshal(metadata)
		require.NoError(t, err)

		rawResult := json.RawMessage(result)
		resp := &jsonrpc.Response[json.RawMessage]{
			Result: &rawResult,
		}

		err = handler.OnMetadataPush(ctx, resp, "node1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid workflow ID")
	})
}

func TestOnMetadataPullResponseWithValidation(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)
	ctx := t.Context()

	err := handler.aggs[handler.shards[0].donID].Start(ctx)
	require.NoError(t, err)
	defer handler.aggs[handler.shards[0].donID].Close()

	t.Run("valid metadata array passes validation", func(t *testing.T) {
		metadata := []gateway_common.WorkflowMetadata{
			{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID1,
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testWorkflowOwner1,
					},
				},
			},
			{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    testWorkflowID2,
					WorkflowName:  testWorkflowNameHex2,
					WorkflowOwner: testWorkflowOwner2,
					WorkflowTag:   testWorkflowTag2,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testWorkflowOwner2,
					},
				},
			},
		}

		result, err := json.Marshal(metadata)
		require.NoError(t, err)

		rawResult := json.RawMessage(result)
		resp := &jsonrpc.Response[json.RawMessage]{
			Result: &rawResult,
		}

		err = handler.OnMetadataPullResponse(ctx, resp, "node1")
		require.NoError(t, err)
	})

	t.Run("invalid metadata in array fails validation", func(t *testing.T) {
		metadata := []gateway_common.WorkflowMetadata{
			{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    "workflowID1",
					WorkflowName:  "workflowName1",
					WorkflowOwner: "workflowOwner1",
					WorkflowTag:   testWorkflowTag1,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testWorkflowOwner1,
					},
				},
			},
			{
				WorkflowSelector: gateway_common.WorkflowSelector{
					WorkflowID:    "", // Invalid: empty workflow ID
					WorkflowName:  testWorkflowNameHex1,
					WorkflowOwner: testWorkflowOwner1,
					WorkflowTag:   testWorkflowTag2,
				},
				AuthorizedKeys: []gateway_common.AuthorizedKey{
					{
						KeyType:   gateway_common.KeyTypeECDSAEVM,
						PublicKey: testWorkflowOwner2,
					},
				},
			},
		}

		result, err := json.Marshal(metadata)
		require.NoError(t, err)

		rawResult := json.RawMessage(result)
		resp := &jsonrpc.Response[json.RawMessage]{
			Result: &rawResult,
		}

		err = handler.OnMetadataPullResponse(ctx, resp, "node1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid workflow ID")
	})
}

func TestWorkflowMetadataHandler_Authorize(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	signerAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

	workflowID := testWorkflowID1
	authorizedKey := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: strings.ToLower(signerAddr.Hex()),
	}
	handler.authorizedKeys = map[string]map[gateway_common.AuthorizedKey]struct{}{
		workflowID: {authorizedKey: {}},
	}

	t.Run("successful authorization", func(t *testing.T) {
		params := json.RawMessage(`{"test": "data"}`)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &params,
		}

		token, err := utils.CreateRequestJWT(*req)
		require.NoError(t, err)

		tokenString, err := token.SignedString(privateKey)
		require.NoError(t, err)

		key, err := handler.Authorize(workflowID, tokenString, req)
		require.NoError(t, err)
		require.NotNil(t, key)
		require.Equal(t, authorizedKey.KeyType, key.KeyType)
		require.Equal(t, authorizedKey.PublicKey, key.PublicKey)
	})

	t.Run("invalid JWT token", func(t *testing.T) {
		params := json.RawMessage(`{"test": "data"}`)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-2",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &params,
		}

		key, err := handler.Authorize(workflowID, "invalid.jwt.token", req)
		require.Error(t, err)
		require.Nil(t, key)
	})

	t.Run("workflow not found in authorized keys", func(t *testing.T) {
		nonExistentWorkflowID := "0x123456"

		params := json.RawMessage(`{"test": "data"}`)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-3",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &params,
		}

		token, err := utils.CreateRequestJWT(*req)
		require.NoError(t, err)

		tokenString, err := token.SignedString(privateKey)
		require.NoError(t, err)

		key, err := handler.Authorize(nonExistentWorkflowID, tokenString, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found")
		require.Nil(t, key)
	})

	t.Run("unauthorized signer", func(t *testing.T) {
		unauthorizedKey, err := crypto.GenerateKey()
		require.NoError(t, err)

		params := json.RawMessage(`{"test": "data"}`)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-4",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &params,
		}

		token, err := utils.CreateRequestJWT(*req)
		require.NoError(t, err)

		tokenString, err := token.SignedString(unauthorizedKey)
		require.NoError(t, err)

		key, err := handler.Authorize(workflowID, tokenString, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "is not authorized for workflow")
		require.Nil(t, key)
	})

	t.Run("JWT digest mismatch", func(t *testing.T) {
		params := json.RawMessage(`{"test": "data"}`)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-5",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &params,
		}

		differentParams := json.RawMessage(`{"different": "data"}`)
		differentReq := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "different-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &differentParams,
		}

		token, err := utils.CreateRequestJWT(*differentReq)
		require.NoError(t, err)

		tokenString, err := token.SignedString(privateKey)
		require.NoError(t, err)

		key, err := handler.Authorize(workflowID, tokenString, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "does not match calculated request digest")
		require.Nil(t, key)
	})

	t.Run("JWT replay protection", func(t *testing.T) {
		params := json.RawMessage(`{"test": "data"}`)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-replay",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &params,
		}

		token, err := utils.CreateRequestJWT(*req)
		require.NoError(t, err)

		tokenString, err := token.SignedString(privateKey)
		require.NoError(t, err)

		key, err := handler.Authorize(workflowID, tokenString, req)
		require.NoError(t, err)
		require.NotNil(t, key)

		// Second authorization with same JWT should fail (replay attack)
		key, err = handler.Authorize(workflowID, tokenString, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "JWT token has already been used. Please generate a new one with new id (jti)")
		require.Nil(t, key)
	})

	t.Run("different JWT IDs should work", func(t *testing.T) {
		params := json.RawMessage(`{"test": "data"}`)
		req1 := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-1",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &params,
		}

		req2 := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-2",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &params,
		}

		token1, err := utils.CreateRequestJWT(*req1)
		require.NoError(t, err)
		tokenString1, err := token1.SignedString(privateKey)
		require.NoError(t, err)

		key1, err := handler.Authorize(workflowID, tokenString1, req1)
		require.NoError(t, err)
		require.NotNil(t, key1)

		token2, err := utils.CreateRequestJWT(*req2)
		require.NoError(t, err)
		tokenString2, err := token2.SignedString(privateKey)
		require.NoError(t, err)

		key2, err := handler.Authorize(workflowID, tokenString2, req2)
		require.NoError(t, err)
		require.NotNil(t, key2)
	})
}

func TestWorkflowMetadataHandler_GetWorkflowID(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)

	workflowOwner := testWorkflowOwner1
	workflowName := "test-workflow"
	workflowNameHash := "0x" + hex.EncodeToString([]byte(workflows.HashTruncateName(workflowName)))
	workflowTag := "v1.0"
	workflowID := testWorkflowID1

	workflowRef := workflowReference{
		workflowOwner: workflowOwner,
		workflowName:  workflowNameHash,
		workflowTag:   workflowTag,
	}
	handler.workflowRefToID = map[workflowReference]string{
		workflowRef: workflowID,
	}

	t.Run("successful workflow lookup", func(t *testing.T) {
		id, found := handler.GetWorkflowID(workflowOwner, workflowNameHash, workflowTag)
		require.True(t, found)
		require.Equal(t, workflowID, id)
	})

	t.Run("workflow not found", func(t *testing.T) {
		id, found := handler.GetWorkflowID(workflowOwner, "nonexistent-workflow", workflowTag)
		require.False(t, found)
		require.Empty(t, id)
	})

	t.Run("workflow not found - different owner", func(t *testing.T) {
		id, found := handler.GetWorkflowID("0xdifferentowner", workflowName, workflowTag)
		require.False(t, found)
		require.Empty(t, id)
	})

	t.Run("workflow not found - different tag", func(t *testing.T) {
		id, found := handler.GetWorkflowID(workflowOwner, workflowName, "v2.0")
		require.False(t, found)
		require.Empty(t, id)
	})
}

func TestWorkflowMetadataHandler_GetWorkflowReference(t *testing.T) {
	handler, _, _ := createTestWorkflowMetadataHandler(t)

	workflowOwner := testWorkflowOwner1
	workflowName := "test-workflow"
	workflowTag := "v1.0"
	workflowID := testWorkflowID1

	expectedRef := workflowReference{
		workflowOwner: workflowOwner,
		workflowName:  "0x" + hex.EncodeToString([]byte(workflows.HashTruncateName(workflowName))),
		workflowTag:   workflowTag,
	}
	handler.workflowIDToRef = map[string]workflowReference{
		workflowID: expectedRef,
	}

	t.Run("successful reference lookup", func(t *testing.T) {
		ref, found := handler.GetWorkflowReference(workflowID)
		require.True(t, found)
		require.Equal(t, expectedRef, ref)
	})

	t.Run("reference not found", func(t *testing.T) {
		nonExistentID := "0x123456"
		_, found := handler.GetWorkflowReference(nonExistentID)
		require.False(t, found)
	})
}
