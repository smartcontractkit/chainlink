package aggregation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
)

func TestIdenticalNodeResponseAggregator_CollectAndAggregate(t *testing.T) {
	t.Run("single node response below threshold", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(2)
		require.NoError(t, err)

		rawMsg := json.RawMessage(`{"result": "success"}`)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-id",
			Result: &rawMsg,
		}

		result, err := agg.CollectAndAggregate(resp, "node1")
		require.NoError(t, err)
		require.Nil(t, result) // Should not return response until threshold is met
	})

	t.Run("threshold reached with identical responses", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(2)
		require.NoError(t, err)

		rawMsg := json.RawMessage(`{"result": "success"}`)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-id",
			Result: &rawMsg,
		}

		// First response - should return nil
		result, err := agg.CollectAndAggregate(resp, "node1")
		require.NoError(t, err)
		require.Nil(t, result)

		// Second response with same content - should return aggregated result
		result, err = agg.CollectAndAggregate(resp, "node2")
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, resp.ID, result.ID)
		require.Equal(t, resp.Result, result.Result)
	})

	t.Run("threshold reached with same node sending multiple times", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(2)
		require.NoError(t, err)

		rawMsg := json.RawMessage(`{"result": "success"}`)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-id",
			Result: &rawMsg,
		}
		result, err := agg.CollectAndAggregate(resp, "node1")
		require.NoError(t, err)
		require.Nil(t, result)

		// Same node sends again - should still not nil
		result, err = agg.CollectAndAggregate(resp, "node1")
		require.NoError(t, err)
		require.Nil(t, result)

		// Different node sends - should return aggregated result
		result, err = agg.CollectAndAggregate(resp, "node2")
		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("different responses do not aggregate", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(2)
		require.NoError(t, err)

		rawMsg1 := json.RawMessage(`{"result": "success"}`)
		resp1 := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-id",
			Result: &rawMsg1,
		}

		rawMsg2 := json.RawMessage(`{"result": "failure"}`)
		resp2 := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-id",
			Result: &rawMsg2,
		}

		result, err := agg.CollectAndAggregate(resp1, "node1")
		require.NoError(t, err)
		require.Nil(t, result)

		// Different response content - should not return aggregated result
		result, err = agg.CollectAndAggregate(resp2, "node2")
		require.NoError(t, err)
		require.Nil(t, result)
	})

	t.Run("threshold 1 immediately returns response", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(1)
		require.NoError(t, err)

		rawMsg := json.RawMessage(`{"result": "success"}`)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-id",
			Result: &rawMsg,
		}

		result, err := agg.CollectAndAggregate(resp, "node1")
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, resp.ID, result.ID)
		require.Equal(t, resp.Result, result.Result)
	})

	t.Run("higher threshold requires more nodes", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(3)
		require.NoError(t, err)

		rawMsg := json.RawMessage(`{"result": "success"}`)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-id",
			Result: &rawMsg,
		}

		result, err := agg.CollectAndAggregate(resp, "node1")
		require.NoError(t, err)
		require.Nil(t, result)

		result, err = agg.CollectAndAggregate(resp, "node2")
		require.NoError(t, err)
		require.Nil(t, result)

		result, err = agg.CollectAndAggregate(resp, "node3")
		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("mixed responses with different digests", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(2)
		require.NoError(t, err)

		rawMsg1 := json.RawMessage(`{"result": "success"}`)
		resp1 := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-id",
			Result: &rawMsg1,
		}

		rawMsg2 := json.RawMessage(`{"result": "failure"}`)
		resp2 := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-id",
			Result: &rawMsg2,
		}

		rawMsg3 := json.RawMessage(`{"result": "success"}`)
		resp3 := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-id",
			Result: &rawMsg3, // Same as resp1
		}

		result, err := agg.CollectAndAggregate(resp1, "node1")
		require.NoError(t, err)
		require.Nil(t, result)

		result, err = agg.CollectAndAggregate(resp2, "node2")
		require.NoError(t, err)
		require.Nil(t, result)

		result, err = agg.CollectAndAggregate(resp3, "node3")
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, resp1.Result, result.Result)
	})

	t.Run("error responses are handled correctly", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(2)
		require.NoError(t, err)

		resp := &jsonrpc.Response[json.RawMessage]{
			ID:    "test-id",
			Error: &jsonrpc.WireError{Code: 500, Message: "Internal error"},
		}

		result, err := agg.CollectAndAggregate(resp, "node1")
		require.NoError(t, err)
		require.Nil(t, result)

		result, err = agg.CollectAndAggregate(resp, "node2")
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, resp.Error.Code, result.Error.Code)
		require.Equal(t, resp.Error.Message, result.Error.Message)
	})
}

func TestIdenticalNodeResponseAggregator_EdgeCases(t *testing.T) {
	t.Run("empty response", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(1)
		require.NoError(t, err)

		resp := &jsonrpc.Response[json.RawMessage]{
			ID: "test-id",
		}

		result, err := agg.CollectAndAggregate(resp, "node1")
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, resp.ID, result.ID)
	})

	t.Run("invalid threshold", func(t *testing.T) {
		_, err := NewIdenticalNodeResponseAggregator(0)
		require.Error(t, err)
	})

	t.Run("nil response", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(1)
		require.NoError(t, err)

		_, err = agg.CollectAndAggregate(nil, "node1")
		require.Error(t, err)
	})

	t.Run("empty node address", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(1)
		require.NoError(t, err)

		rawMsg := json.RawMessage(`{"result": "success"}`)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-id",
			Result: &rawMsg,
		}

		// Empty node address should not work
		_, err = agg.CollectAndAggregate(resp, "")
		require.Error(t, err)
	})
}

func TestIdenticalNodeResponseAggregator_NodeChangesResponse(t *testing.T) {
	t.Run("node changes response and reaches threshold", func(t *testing.T) {
		agg, err := NewIdenticalNodeResponseAggregator(2)
		require.NoError(t, err)

		// Create two different responses
		rawRes1 := json.RawMessage([]byte(`{"key":"value1"}`))
		resp1 := &jsonrpc.Response[json.RawMessage]{
			ID:     "test",
			Result: &rawRes1,
		}

		rawRes2 := json.RawMessage([]byte(`{"key":"value2"}`))
		resp2 := &jsonrpc.Response[json.RawMessage]{
			ID:     "test",
			Result: &rawRes2,
		}

		// Node1 submits response1
		result, err := agg.CollectAndAggregate(resp1, "node1")
		require.NoError(t, err)
		require.Nil(t, result) // Not enough nodes yet

		// Node2 submits response2
		result, err = agg.CollectAndAggregate(resp2, "node2")
		require.NoError(t, err)
		require.Nil(t, result) // Different responses, threshold not reached

		// Node1 changes to response2 - this should reach threshold
		result, err = agg.CollectAndAggregate(resp2, "node1")
		require.NoError(t, err)
		require.NotNil(t, result)
		require.JSONEq(t, string(*resp2.Result), string(*result.Result))

		// Generate keys
		key1, err := resp1.Digest()
		require.NoError(t, err)
		key2, err := resp2.Digest()
		require.NoError(t, err)

		// Both nodes should be associated with key2
		require.Equal(t, key2, agg.nodeToResponse["node1"])
		require.Equal(t, key2, agg.nodeToResponse["node2"])

		// No nodes should be in resp1 group
		if nodes, exists := agg.responses[key1]; exists {
			require.Empty(t, nodes)
		}

		// Both nodes should be in resp2 group
		nodes, exists := agg.responses[key2]
		require.True(t, exists)
		require.True(t, nodes.Contains("node1"))
		require.True(t, nodes.Contains("node2"))
	})
}
