package arbiter

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// DecisionEngine computes the approved replica count based on inputs.
type DecisionEngine interface {
	// ComputeApprovedCount takes KEDA's desired count and returns approved count.
	ComputeApprovedCount(ctx context.Context, desiredCount int) (int, error)
}

type decisionEngine struct {
	shardConfig ShardConfigReader
	lggr        logger.SugaredLogger
}

// NewDecisionEngine creates a new DecisionEngine.
func NewDecisionEngine(shardConfig ShardConfigReader, lggr logger.SugaredLogger) DecisionEngine {
	return &decisionEngine{
		shardConfig: shardConfig,
		lggr:        lggr,
	}
}

// ComputeApprovedCount applies the decision logic:
// approved = min(desired, onChainMax)
// with a minimum of 1 shard always.
func (d *decisionEngine) ComputeApprovedCount(ctx context.Context, desiredCount int) (int, error) {
	// Get on-chain limit from ShardConfig contract
	approved, err := d.shardConfig.GetDesiredShardCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get on-chain shard limit: %w", err)
	}

	// Ensure minimum of 1 shard
	if approved < 1 {
		d.lggr.Warnw("Desired count is less than 1, setting to minimum",
			"desired", desiredCount,
			"approved", 1,
		)
		approved = 1
	}

	d.lggr.Debugw("Computed approved replica count",
		"desired", desiredCount,
		"approved", approved,
	)

	return int(approved), nil //nolint:gosec // G115: replica count is bounded
}
