package requests

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

	s := NewStore()
	rid := uuid.New().String()
	req := &Request{
		WorkflowExecutionID: rid,
		ExpiresAt:           n.Add(10 * time.Second),
	}

	t.Run("Add", func(t *testing.T) {
		err := s.Add(req)
		require.NoError(t, err)
	})

	t.Run("add duplicate", func(t *testing.T) {
		err := s.Add(req)
		require.Error(t, err)
	})

	t.Run("evict", func(t *testing.T) {
		_, wasPresent := s.evict(rid)
		assert.True(t, wasPresent)
		assert.Len(t, s.requests, 0)

		// evicting doesn't remove from the list of requestIDs
		assert.Len(t, s.requestIDs, 1)
	})

	t.Run("firstN, evicts removed items", func(t *testing.T) {
		r, err := s.FirstN(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, r, 0)
	})

	t.Run("firstN, zero batch size", func(t *testing.T) {
		_, err := s.FirstN(ctx, 0)
		assert.ErrorContains(t, err, "batchsize cannot be 0")
	})

	t.Run("firstN, batchSize larger than queue", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			err := s.Add(&Request{WorkflowExecutionID: uuid.New().String(), ExpiresAt: n.Add(1 * time.Hour)})
			require.NoError(t, err)
		}
		items, err := s.FirstN(ctx, 100)
		require.NoError(t, err)
		assert.Len(t, items, 10)
	})

	t.Run("getN", func(t *testing.T) {
		rid2 := uuid.New().String()
		err := s.Add(req)
		require.NoError(t, err)
		reqs := s.GetN(ctx, []string{rid, rid2})
		require.Equal(t, 1, len(reqs))
	})
}
