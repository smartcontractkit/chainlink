package ocr3

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestOCR3Store(t *testing.T) {
	ctx := tests.Context(t)
	n := time.Now()

	s := newStore()
	rid := uuid.New().String()
	req := &request{
		WorkflowExecutionID: rid,
		ExpiresAt:           n.Add(10 * time.Second),
	}

	t.Run("add", func(t *testing.T) {
		err := s.add(ctx, req)
		require.NoError(t, err)
	})

	t.Run("add duplicate", func(t *testing.T) {
		err := s.add(ctx, req)
		require.Error(t, err)
	})

	t.Run("evict", func(t *testing.T) {
		_, wasPresent := s.evict(ctx, rid)
		assert.True(t, wasPresent)
		assert.Len(t, s.requests, 0)

		// evicting doesn't remove from the list of requestIDs
		assert.Len(t, s.requestIDs, 1)
	})

	t.Run("firstN, evicts removed items", func(t *testing.T) {
		r, err := s.firstN(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, r, 0)
	})

	t.Run("firstN, zero batch size", func(t *testing.T) {
		_, err := s.firstN(ctx, 0)
		assert.ErrorContains(t, err, "batchsize cannot be 0")
	})

	t.Run("firstN, batchSize larger than queue", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			err := s.add(ctx, &request{WorkflowExecutionID: uuid.New().String(), ExpiresAt: n.Add(1 * time.Hour)})
			require.NoError(t, err)
		}
		items, err := s.firstN(ctx, 100)
		require.NoError(t, err)
		assert.Len(t, items, 10)
	})
}
