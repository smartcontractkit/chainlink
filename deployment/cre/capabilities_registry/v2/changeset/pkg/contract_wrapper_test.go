package pkg

import (
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaginate(t *testing.T) {
	fakeContract := func(total int, calls *int) func(start, limit *big.Int) ([]int, error) {
		return func(start, limit *big.Int) ([]int, error) {
			*calls++
			var page []int
			for i := start.Int64(); i < int64(total) && i < start.Int64()+limit.Int64(); i++ {
				page = append(page, int(i))
			}
			return page, nil
		}
	}

	for _, tc := range []struct {
		total     int
		wantCalls int
	}{
		{total: 0, wantCalls: 1},
		{total: 1, wantCalls: 1},
		{total: PageSize - 1, wantCalls: 1},
		{total: PageSize, wantCalls: 2},
		{total: PageSize + 1, wantCalls: 2},
		{total: 3*PageSize + 7, wantCalls: 4},
	} {
		calls := 0
		got, err := paginate(fakeContract(tc.total, &calls))
		require.NoError(t, err)
		require.Len(t, got, tc.total)
		for i, v := range got {
			assert.Equal(t, i, v)
		}
		assert.Equal(t, tc.wantCalls, calls, "total=%d", tc.total)
	}

	t.Run("error", func(t *testing.T) {
		_, err := paginate(func(start, limit *big.Int) ([]int, error) {
			return nil, errors.New("boom")
		})
		require.ErrorContains(t, err, "boom")
	})
}
