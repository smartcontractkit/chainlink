package txm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

func TestTimeBasedDetection(t *testing.T) {
	t.Parallel()

	t.Run("returns false if transaction is not stuck", func(t *testing.T) {
		config := StuckTxDetectorConfig{
			BlockTime:             10 * time.Second,
			StuckTxBlockThreshold: 5,
		}
		fromAddress := testutils.NewAddress()
		s := NewStuckTxDetector(logger.Test(t), "", config)

		// No previous broadcast
		tx := &types.Transaction{
			ID:              1,
			LastBroadcastAt: nil,
			FromAddress:     fromAddress,
		}
		assert.False(t, s.timeBasedDetection(tx))
		// Not enough time has passed since last broadcast
		now := time.Now()
		tx.LastBroadcastAt = &now
		assert.False(t, s.timeBasedDetection(tx))
		// Not enough time has passed since last purge
		tx.LastBroadcastAt = &time.Time{}
		s.lastPurgeMap[fromAddress] = now
		assert.False(t, s.timeBasedDetection(tx))
	})

	t.Run("returns true if transaction is stuck", func(t *testing.T) {
		config := StuckTxDetectorConfig{
			BlockTime:             10 * time.Second,
			StuckTxBlockThreshold: 5,
		}
		fromAddress := testutils.NewAddress()
		s := NewStuckTxDetector(logger.Test(t), "", config)

		tx := &types.Transaction{
			ID:              1,
			LastBroadcastAt: &time.Time{},
			FromAddress:     fromAddress,
		}
		assert.True(t, s.timeBasedDetection(tx))
	})

	t.Run("marks first tx as stuck, updates purge time for address, and returns false for the second tx with the same broadcast time", func(t *testing.T) {
		config := StuckTxDetectorConfig{
			BlockTime:             1 * time.Second,
			StuckTxBlockThreshold: 10,
		}
		fromAddress := testutils.NewAddress()
		s := NewStuckTxDetector(logger.Test(t), "", config)

		tx1 := &types.Transaction{
			ID:              1,
			LastBroadcastAt: &time.Time{},
			FromAddress:     fromAddress,
		}
		tx2 := &types.Transaction{
			ID:              2,
			LastBroadcastAt: &time.Time{},
			FromAddress:     fromAddress,
		}
		assert.True(t, s.timeBasedDetection(tx1))
		assert.False(t, s.timeBasedDetection(tx2))
	})
}
