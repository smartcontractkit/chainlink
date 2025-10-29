package aggregation

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
)

func createTestMetrics(t *testing.T) *metrics.Metrics {
	m, err := metrics.NewMetrics()
	require.NoError(t, err)
	return m
}

func TestWorkflowMetadataAggregator_StartStop(t *testing.T) {
	lggr := logger.Test(t)
	testMetrics := createTestMetrics(t)
	agg := NewWorkflowMetadataAggregator(lggr, 2, 100*time.Millisecond, testMetrics)

	ctx := testutils.Context(t)

	err := agg.Start(ctx)
	require.NoError(t, err)
	require.Equal(t, "Started", agg.State())

	// Test that starting again returns error
	err = agg.Start(ctx)
	require.Error(t, err)

	err = agg.Close()
	require.NoError(t, err)
	require.Equal(t, "Stopped", agg.State())

	// Test that closing again returns Error
	err = agg.Close()
	require.Error(t, err)
}

func getRandomECDSAPublicKey(t *testing.T) string {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	return crypto.PubkeyToAddress(key.PublicKey).Hex()
}

// createTestWorkflowMetadata creates a test WorkflowMetadata observation
func createTestWorkflowMetadata(workflowID, workflowName, workflowOwner, workflowTag string, authorizedKeys []gateway_common.AuthorizedKey) *gateway_common.WorkflowMetadata {
	return &gateway_common.WorkflowMetadata{
		WorkflowSelector: gateway_common.WorkflowSelector{
			WorkflowID:    workflowID,
			WorkflowName:  workflowName,
			WorkflowOwner: workflowOwner,
			WorkflowTag:   workflowTag,
		},
		AuthorizedKeys: authorizedKeys,
	}
}

func TestWorkflowMetadataAggregator_Collect(t *testing.T) {
	lggr := logger.Test(t)
	testMetrics := createTestMetrics(t)
	agg := NewWorkflowMetadataAggregator(lggr, 2, 10*time.Second, testMetrics)

	authorizedKey := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: getRandomECDSAPublicKey(t),
	}

	observation := createTestWorkflowMetadata("workflowID", "workflowName", "workflowOwner", "workflowTag", []gateway_common.AuthorizedKey{authorizedKey})

	err := agg.Collect(observation, "node1")
	require.NoError(t, err)
	require.Len(t, agg.observations, 1)
	require.Len(t, agg.observedAt, 1)
	require.Len(t, agg.observedAt["node1"], 1)

	digest, err := observation.Digest()
	require.NoError(t, err)
	nodeObs, exists := agg.observations[digest]
	require.True(t, exists)
	require.Equal(t, observation, nodeObs.observation)
	require.True(t, nodeObs.nodes.Contains("node1"))
	require.Len(t, nodeObs.nodes, 1)
	timestamp1, ok := agg.observedAt["node1"][digest]
	require.True(t, ok)

	// Test collecting from second node with same observation
	err = agg.Collect(observation, "node2")
	require.NoError(t, err)
	require.Len(t, agg.observations, 1)
	require.Len(t, agg.observedAt, 2)

	nodeObs, exists = agg.observations[digest]
	require.True(t, exists)
	require.True(t, nodeObs.nodes.Contains("node1"))
	require.True(t, nodeObs.nodes.Contains("node2"))
	require.Len(t, nodeObs.nodes, 2)

	// Test collecting from same node again (should update timestamp)
	time.Sleep(10 * time.Millisecond) // Small delay to ensure timestamp difference
	err = agg.Collect(observation, "node1")
	require.NoError(t, err)
	require.Len(t, agg.observations, 1)
	require.Len(t, agg.observedAt, 2)

	nodeObs, exists = agg.observations[digest]
	require.True(t, exists)
	require.Len(t, nodeObs.nodes, 2)

	digests, ok := agg.observedAt["node1"]
	require.True(t, ok)
	timestamp2, ok := digests[digest]
	require.True(t, ok)
	require.NotEqual(t, timestamp1, timestamp2)
}

func TestWorkflowMetadataAggregator_CollectDifferentObservations(t *testing.T) {
	lggr := logger.Test(t)
	testMetrics := createTestMetrics(t)
	agg := NewWorkflowMetadataAggregator(lggr, 2, 10*time.Second, testMetrics)

	authorizedKey1 := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: getRandomECDSAPublicKey(t),
	}

	authorizedKey2 := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: getRandomECDSAPublicKey(t),
	}

	observation1 := createTestWorkflowMetadata("workflowID1", "workflowName1", "workflowOwner1", "workflowTag1", []gateway_common.AuthorizedKey{authorizedKey1})
	observation2 := createTestWorkflowMetadata("workflowID2", "workflowName2", "workflowOwner2", "workflowTag2", []gateway_common.AuthorizedKey{authorizedKey2})

	// Collect different observations
	err := agg.Collect(observation1, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation2, "node2")
	require.NoError(t, err)

	// Should have 2 different observations
	require.Len(t, agg.observations, 2)
	require.Len(t, agg.observedAt, 2)

	digest1, err := observation1.Digest()
	require.NoError(t, err)
	digest2, err := observation2.Digest()
	require.NoError(t, err)

	nodeObs1, exists := agg.observations[digest1]
	require.True(t, exists)
	require.Equal(t, observation1, nodeObs1.observation)
	require.True(t, nodeObs1.nodes.Contains("node1"))

	nodeObs2, exists := agg.observations[digest2]
	require.True(t, exists)
	require.Equal(t, observation2, nodeObs2.observation)
	require.True(t, nodeObs2.nodes.Contains("node2"))
}

func TestWorkflowMetadataAggregator_Aggregate(t *testing.T) {
	lggr := logger.Test(t)
	threshold := 2
	testMetrics := createTestMetrics(t)
	agg := NewWorkflowMetadataAggregator(lggr, threshold, 10*time.Second, testMetrics)

	publicKey1 := getRandomECDSAPublicKey(t)
	publicKey2 := getRandomECDSAPublicKey(t)
	publicKey3 := getRandomECDSAPublicKey(t)

	authorizedKey1 := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: publicKey1,
	}

	authorizedKey2 := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: publicKey2,
	}

	authorizedKey3 := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: publicKey3,
	}

	observation1 := createTestWorkflowMetadata("workflowID1", "workflowName1", "workflowOwner1", "workflowTag1", []gateway_common.AuthorizedKey{authorizedKey1})
	observation2 := createTestWorkflowMetadata("workflowID1", "workflowName1", "workflowOwner1", "workflowTag1", []gateway_common.AuthorizedKey{authorizedKey2})
	observation3 := createTestWorkflowMetadata("workflowID2", "workflowName2", "workflowOwner2", "workflowTag2", []gateway_common.AuthorizedKey{authorizedKey3})

	// Test aggregation with no observations
	result, err := agg.Aggregate()
	require.NoError(t, err)
	require.Empty(t, result)

	// Add observations below threshold
	err = agg.Collect(observation1, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation2, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation3, "node1")
	require.NoError(t, err)

	result, err = agg.Aggregate()
	require.NoError(t, err)
	require.Empty(t, result)

	// Add observations to reach threshold for workflowID1
	err = agg.Collect(observation1, "node2")
	require.NoError(t, err)
	err = agg.Collect(observation2, "node2")
	require.NoError(t, err)

	result, err = agg.Aggregate()
	require.NoError(t, err)
	require.Len(t, result, 2) // observation1 and observation2 reach threshold

	// Add observation to reach threshold for workflowID2
	err = agg.Collect(observation3, "node2")
	require.NoError(t, err)

	result, err = agg.Aggregate()
	require.NoError(t, err)
	require.Len(t, result, 3) // observation3 reaches threshold
}

func TestWorkflowMetadataAggregator_Aggregate_ChronologicalOrder(t *testing.T) {
	lggr := logger.Test(t)
	threshold := 2
	testMetrics := createTestMetrics(t)
	agg := NewWorkflowMetadataAggregator(lggr, threshold, 10*time.Second, testMetrics)

	observation1 := createTestWorkflowMetadata("workflowID1", "workflowName1", "workflowOwner1", "workflowTag1", []gateway_common.AuthorizedKey{
		{KeyType: gateway_common.KeyTypeECDSAEVM, PublicKey: getRandomECDSAPublicKey(t)},
	})
	observation2 := createTestWorkflowMetadata("workflowID2", "workflowName2", "workflowOwner2", "workflowTag2", []gateway_common.AuthorizedKey{
		{KeyType: gateway_common.KeyTypeECDSAEVM, PublicKey: getRandomECDSAPublicKey(t)},
	})
	observation3 := createTestWorkflowMetadata("workflowID3", "workflowName3", "workflowOwner3", "workflowTag3", []gateway_common.AuthorizedKey{
		{KeyType: gateway_common.KeyTypeECDSAEVM, PublicKey: getRandomECDSAPublicKey(t)},
	})

	// Collect observation1 first
	err := agg.Collect(observation1, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation1, "node2")
	require.NoError(t, err)

	// Collect observation2 second
	time.Sleep(10 * time.Millisecond)
	err = agg.Collect(observation2, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation2, "node2")
	require.NoError(t, err)

	// Collect observation3 third (most recent)
	time.Sleep(10 * time.Millisecond)
	err = agg.Collect(observation3, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation3, "node2")
	require.NoError(t, err)

	// All observations should reach threshold
	result, err := agg.Aggregate()
	require.NoError(t, err)
	require.Len(t, result, 3)

	// Verify chronological order
	require.Equal(t, "workflowID3", result[0].WorkflowSelector.WorkflowID)
	require.Equal(t, "workflowName3", result[0].WorkflowSelector.WorkflowName)

	require.Equal(t, "workflowID2", result[1].WorkflowSelector.WorkflowID)
	require.Equal(t, "workflowName2", result[1].WorkflowSelector.WorkflowName)

	require.Equal(t, "workflowID1", result[2].WorkflowSelector.WorkflowID)
	require.Equal(t, "workflowName1", result[2].WorkflowSelector.WorkflowName)
}

func TestWorkflowMetadataAggregator_Aggregate_ChronologicalOrder_SameWorkflowNameOwnerTag(t *testing.T) {
	lggr := logger.Test(t)
	threshold := 2
	testMetrics := createTestMetrics(t)
	agg := NewWorkflowMetadataAggregator(lggr, threshold, 10*time.Second, testMetrics)

	// Create four observations
	observation1 := createTestWorkflowMetadata("workflowID1", "workflowName1", "workflowOwner1", "workflowTag1", []gateway_common.AuthorizedKey{
		{KeyType: gateway_common.KeyTypeECDSAEVM, PublicKey: getRandomECDSAPublicKey(t)},
	})
	observation2 := createTestWorkflowMetadata("workflowID2", "workflowName1", "workflowOwner1", "workflowTag1", []gateway_common.AuthorizedKey{
		{KeyType: gateway_common.KeyTypeECDSAEVM, PublicKey: getRandomECDSAPublicKey(t)},
	})
	observation3 := createTestWorkflowMetadata("workflowID3", "workflowName1", "workflowOwner1", "workflowTag1", []gateway_common.AuthorizedKey{
		{KeyType: gateway_common.KeyTypeECDSAEVM, PublicKey: getRandomECDSAPublicKey(t)},
	})
	observation4 := createTestWorkflowMetadata("workflowID4", "workflowName1", "workflowOwner1", "workflowTag1", []gateway_common.AuthorizedKey{
		{KeyType: gateway_common.KeyTypeECDSAEVM, PublicKey: getRandomECDSAPublicKey(t)},
	})

	// Collect observation1 (oldest, reaches threshold)
	err := agg.Collect(observation1, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation1, "node2")
	require.NoError(t, err)

	// Collect observation2 (doesn't reach threshold initially)
	time.Sleep(10 * time.Millisecond)
	err = agg.Collect(observation2, "node1")
	require.NoError(t, err)

	// Collect observation3 (reaches threshold)
	time.Sleep(10 * time.Millisecond)
	err = agg.Collect(observation3, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation3, "node2")
	require.NoError(t, err)

	// Collect observation4 (most recent, doesn't reach threshold)
	time.Sleep(10 * time.Millisecond)
	err = agg.Collect(observation4, "node1")
	require.NoError(t, err)

	// Only observations that reached threshold should be returned
	result, err := agg.Aggregate()
	require.NoError(t, err)
	require.Len(t, result, 2)

	// Verify order: observation3 (newer) before observation1 (older)
	require.Equal(t, "workflowID3", result[0].WorkflowSelector.WorkflowID)
	require.Equal(t, "workflowID1", result[1].WorkflowSelector.WorkflowID)

	// Now make observation2 reach threshold (it was collected before observation3)
	time.Sleep(10 * time.Millisecond)
	err = agg.Collect(observation2, "node2")
	require.NoError(t, err)

	result, err = agg.Aggregate()
	require.NoError(t, err)
	require.Len(t, result, 3)

	// Verify order: observation3 (newest), observation2 (middle), observation1 (oldest)
	require.Equal(t, "workflowID3", result[0].WorkflowSelector.WorkflowID)
	require.Equal(t, "workflowID2", result[1].WorkflowSelector.WorkflowID)
	require.Equal(t, "workflowID1", result[2].WorkflowSelector.WorkflowID)
}

func TestWorkflowMetadataAggregator_ReapObservations(t *testing.T) {
	lggr := logger.Test(t)
	cleanupInterval := 1 * time.Second
	testMetrics := createTestMetrics(t)
	agg := NewWorkflowMetadataAggregator(lggr, 2, cleanupInterval, testMetrics)

	authorizedKey1 := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: "test-public-key-1",
	}

	authorizedKey2 := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: "test-public-key-2",
	}

	observation1 := createTestWorkflowMetadata("workflowID1", "workflowName1", "workflowOwner1", "workflowTag1", []gateway_common.AuthorizedKey{authorizedKey1})
	observation2 := createTestWorkflowMetadata("workflowID2", "workflowName2", "workflowOwner2", "workflowTag2", []gateway_common.AuthorizedKey{authorizedKey2})

	err := agg.Collect(observation1, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation1, "node2")
	require.NoError(t, err)
	err = agg.Collect(observation2, "node1")
	require.NoError(t, err)
	require.Len(t, agg.observations, 2)
	require.Len(t, agg.observedAt, 2)

	err = agg.Start(testutils.Context(t))
	require.NoError(t, err)
	// Wait for cleanup interval to pass
	time.Sleep(cleanupInterval + 100*time.Millisecond)
	require.Empty(t, agg.observations)
}

func TestWorkflowMetadataAggregator_ReapObservations_UnexpiredObservation(t *testing.T) {
	lggr := logger.Test(t)
	cleanupInterval := 1 * time.Second
	testMetrics := createTestMetrics(t)
	agg := NewWorkflowMetadataAggregator(lggr, 2, cleanupInterval, testMetrics)

	authorizedKey := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: "test-public-key-1",
	}

	observation := createTestWorkflowMetadata("workflowID", "workflowName", "workflowOwner", "workflowTag", []gateway_common.AuthorizedKey{authorizedKey})

	err := agg.Collect(observation, "node1")
	require.NoError(t, err)
	err = agg.Collect(observation, "node2")
	require.NoError(t, err)
	require.Len(t, agg.observations, 1)
	require.Len(t, agg.observedAt, 2)

	// Wait for cleanup interval to pass
	time.Sleep(cleanupInterval + 100*time.Millisecond)
	// Add observation from node3 (fresh)
	err = agg.Collect(observation, "node3")
	require.NoError(t, err)
	// Manually trigger cleanup
	agg.reapObservations(context.Background())

	require.Len(t, agg.observations, 1)
	digest, err := observation.Digest()
	require.NoError(t, err)
	o, ok := agg.observations[digest]
	require.True(t, ok)
	require.Equal(t, observation, o.observation)
	require.True(t, o.nodes.Contains("node3"))
}

func TestWorkflowMetadataAggregator_Collect_EdgeCases(t *testing.T) {
	lggr := logger.Test(t)

	t.Run("empty workflow ID", func(t *testing.T) {
		testMetrics := createTestMetrics(t)
		agg := NewWorkflowMetadataAggregator(lggr, 1, 10*time.Second, testMetrics)

		authorizedKey := gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSAEVM,
			PublicKey: getRandomECDSAPublicKey(t),
		}

		observation := createTestWorkflowMetadata("", "workflowName", "workflowOwner", "workflowTag", []gateway_common.AuthorizedKey{authorizedKey})

		err := agg.Collect(observation, "node1")
		require.Error(t, err)
	})

	t.Run("empty workflow name", func(t *testing.T) {
		testMetrics := createTestMetrics(t)
		agg := NewWorkflowMetadataAggregator(lggr, 1, 10*time.Second, testMetrics)

		authorizedKey := gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSAEVM,
			PublicKey: getRandomECDSAPublicKey(t),
		}

		observation := createTestWorkflowMetadata("workflowID", "", "workflowOwner", "workflowTag", []gateway_common.AuthorizedKey{authorizedKey})

		err := agg.Collect(observation, "node1")
		require.Error(t, err)
	})

	t.Run("empty workflow owner", func(t *testing.T) {
		testMetrics := createTestMetrics(t)
		agg := NewWorkflowMetadataAggregator(lggr, 1, 10*time.Second, testMetrics)

		authorizedKey := gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSAEVM,
			PublicKey: getRandomECDSAPublicKey(t),
		}

		observation := createTestWorkflowMetadata("workflowID", "workflowName", "", "workflowTag", []gateway_common.AuthorizedKey{authorizedKey})

		err := agg.Collect(observation, "node1")
		require.Error(t, err)
	})

	t.Run("empty workflow tag", func(t *testing.T) {
		testMetrics := createTestMetrics(t)
		agg := NewWorkflowMetadataAggregator(lggr, 1, 10*time.Second, testMetrics)

		authorizedKey := gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSAEVM,
			PublicKey: getRandomECDSAPublicKey(t),
		}

		observation := createTestWorkflowMetadata("workflowID", "workflowName", "workflowOwner", "", []gateway_common.AuthorizedKey{authorizedKey})

		err := agg.Collect(observation, "node1")
		require.Error(t, err)
	})

	t.Run("empty node address", func(t *testing.T) {
		testMetrics := createTestMetrics(t)
		agg := NewWorkflowMetadataAggregator(lggr, 1, 10*time.Second, testMetrics)

		authorizedKey := gateway_common.AuthorizedKey{
			KeyType:   gateway_common.KeyTypeECDSAEVM,
			PublicKey: getRandomECDSAPublicKey(t),
		}

		observation := createTestWorkflowMetadata("workflowID", "workflowName", "workflowOwner", "workflowTag", []gateway_common.AuthorizedKey{authorizedKey})

		err := agg.Collect(observation, "")
		require.Error(t, err)
	})
}
