package vault

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

func deletePendingItem(t *testing.T, id string) *vaultcommon.StoredPendingQueueItem {
	t.Helper()
	req := &vaultcommon.DeleteSecretsRequest{
		RequestId:     id,
		OrgId:         "org",
		WorkflowOwner: "wf",
		Ids: []*vaultcommon.SecretIdentifier{
			{Key: "k", Owner: "owner", Namespace: "ns"},
		},
	}
	a, err := anypb.New(req)
	require.NoError(t, err)
	return &vaultcommon.StoredPendingQueueItem{Id: id, Item: a}
}

func TestPackPendingQueueItems_Empty(t *testing.T) {
	ctx := t.Context()
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	packed, truncated := packPendingQueueItemsToByteBudgets(ctx, nil, r.cfg, r.onchainCfg.F, r.lggr)
	require.Empty(t, packed)
	require.False(t, truncated)
}

func TestPackPendingQueueItems_AllPackedNoTruncation(t *testing.T) {
	ctx := t.Context()
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	items := []*vaultcommon.StoredPendingQueueItem{
		deletePendingItem(t, "a"),
		deletePendingItem(t, "b"),
	}
	packed, truncated := packPendingQueueItemsToByteBudgets(ctx, items, r.cfg, r.onchainCfg.F, r.lggr)
	require.Len(t, packed, 2)
	require.False(t, truncated)
}

func TestPackPendingQueueItems_TruncatedWhenBudgetExceeded(t *testing.T) {
	ctx := t.Context()
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	items := []*vaultcommon.StoredPendingQueueItem{
		deletePendingItem(t, "a"),
		deletePendingItem(t, "b"),
		deletePendingItem(t, "c"),
	}
	c1, err := pendingQueueItemObservationBytes(ctx, items[0], r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	c2, err := pendingQueueItemObservationBytes(ctx, items[1], r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	_, err = pendingQueueItemObservationBytes(ctx, items[2], r.cfg, r.onchainCfg.F)
	require.NoError(t, err)

	// Allow first two observations but not enough headroom for the third.
	r.cfg.ObsArrayBudgetBytes = c1 + c2 + 1
	r.cfg.PrecursorArrayBudgetBytes = 10_000_000

	packed, truncated := packPendingQueueItemsToByteBudgets(ctx, items, r.cfg, r.onchainCfg.F, r.lggr)
	require.Len(t, packed, 2)
	require.True(t, truncated)
}

func TestPackPendingQueueItems_ForcedOversizeFirstStillPackedWithoutTruncationFlag(t *testing.T) {
	ctx := t.Context()
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	item := deletePendingItem(t, "solo")
	cost, err := pendingQueueItemObservationBytes(ctx, item, r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	r.cfg.ObsArrayBudgetBytes = cost - 1
	r.cfg.PrecursorArrayBudgetBytes = 10_000_000

	packed, truncated := packPendingQueueItemsToByteBudgets(ctx, []*vaultcommon.StoredPendingQueueItem{item}, r.cfg, r.onchainCfg.F, r.lggr)
	require.Len(t, packed, 1)
	require.False(t, truncated)
}

func TestPackPendingQueueItems_TruncatedWhenPrecursorBudgetExceeded(t *testing.T) {
	ctx := t.Context()
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	items := []*vaultcommon.StoredPendingQueueItem{
		deletePendingItem(t, "a"),
		deletePendingItem(t, "b"),
		deletePendingItem(t, "c"),
	}
	p1, err := pendingQueueItemOutcomeBytes(ctx, items[0], r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	p2, err := pendingQueueItemOutcomeBytes(ctx, items[1], r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	_, err = pendingQueueItemOutcomeBytes(ctx, items[2], r.cfg, r.onchainCfg.F)
	require.NoError(t, err)

	r.cfg.ObsArrayBudgetBytes = 10_000_000
	r.cfg.PrecursorArrayBudgetBytes = p1 + p2 + 1

	packed, truncated := packPendingQueueItemsToByteBudgets(ctx, items, r.cfg, r.onchainCfg.F, r.lggr)
	require.Len(t, packed, 2)
	require.True(t, truncated)
}
