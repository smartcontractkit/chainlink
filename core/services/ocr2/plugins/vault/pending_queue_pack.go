package vault

import (
	"context"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// packPendingQueueItemsToByteBudgets greedily keeps the prefix of sortedCandidates
// that fits under ObsArrayBudgetBytes (next-round observation wire) and
// PrecursorArrayBudgetBytes (next-round precursor wire), using pendingQueueItem*Bytes.
//
// truncated is true when packing stopped early because adding the next item would exceed a
// byte budget (observation or precursor), after at least one item was already accepted. It is
// false when every candidate is packed, when there are no candidates, when only the first item
// is forced in despite exceeding budgets, or when remaining items are only skipped due to
// size-estimation errors (not a byte-capacity backlog signal).
func packPendingQueueItemsToByteBudgets(
	ctx context.Context,
	sortedCandidates []*vaultcommon.StoredPendingQueueItem,
	cfg *ReportingPluginConfig,
	f int,
	lggr logger.Logger,
) (packed []*vaultcommon.StoredPendingQueueItem, truncated bool) {
	candidateCount := len(sortedCandidates)
	var obsAccum, precursorAccum int
	packed = sortedCandidates[:0]
	stoppedEarly := false
	for _, item := range sortedCandidates {
		obsCost, err := pendingQueueItemObservationBytes(ctx, item, cfg, f)
		if err != nil {
			lggr.Errorw("could not estimate observation bytes for pending queue item, skipping", "id", item.Id, "error", err)
			continue
		}
		precursorCost, err := pendingQueueItemOutcomeBytes(ctx, item, cfg, f)
		if err != nil {
			lggr.Errorw("could not estimate outcome bytes for pending queue item, skipping", "id", item.Id, "error", err)
			continue
		}

		exceedsObs := obsAccum+obsCost > cfg.ObsArrayBudgetBytes
		exceedsPrecursor := precursorAccum+precursorCost > cfg.PrecursorArrayBudgetBytes
		if exceedsObs || exceedsPrecursor {
			if len(packed) != 0 {
				lggr.Warnw("pending queue truncated due to byte budget",
					"kept", len(packed), "totalCandidates", candidateCount)
				stoppedEarly = true
				break
			}
			lggr.Errorw("single pending queue item exceeds byte budget; check limit configuration",
				"id", item.Id,
				"obsCost", obsCost, "obsArrayBudget", cfg.ObsArrayBudgetBytes,
				"precursorCost", precursorCost, "precursorArrayBudget", cfg.PrecursorArrayBudgetBytes,
			)
		}
		obsAccum += obsCost
		precursorAccum += precursorCost
		packed = append(packed, item)
	}
	truncated = stoppedEarly
	return packed, truncated
}
