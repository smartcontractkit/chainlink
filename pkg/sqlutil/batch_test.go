package sqlutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBatching(t *testing.T) {
	const defaultBatchSize uint = 1000

	t.Run("Verify batching", func(t *testing.T) {
		var callCount uint
		var maxCalls uint = 100
		err := Batch(func(offset, limit uint) (count uint, err error) {
			require.Equal(t, callCount*defaultBatchSize, offset)
			callCount++
			if callCount == maxCalls {
				return 0, nil
			}
			return limit, nil
		}, defaultBatchSize)
		require.NoError(t, err)
	})

	t.Run("Handle batch errors", func(t *testing.T) {
		handleErr := "error during batching"
		err := Batch(func(offset, limit uint) (count uint, err error) {
			return 0, errors.New(handleErr)
		}, defaultBatchSize)
		require.EqualError(t, err, handleErr)
	})

	t.Run("Invalid batch size", func(t *testing.T) {
		err := Batch(func(offset, limit uint) (count uint, err error) {
			return 0, nil
		}, 0)
		require.EqualError(t, err, batchSizeErr)
	})
}
