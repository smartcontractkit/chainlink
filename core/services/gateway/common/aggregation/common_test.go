package aggregation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringSet(t *testing.T) {
	t.Run("add and contains", func(t *testing.T) {
		s := make(StringSet)

		require.False(t, s.Contains("test"))

		s.Add("test")
		require.True(t, s.Contains("test"))
		require.False(t, s.Contains("other"))

		s.Add("other")
		require.True(t, s.Contains("test"))
		require.True(t, s.Contains("other"))

		// Adding duplicate should not change anything
		s.Add("test")
		require.True(t, s.Contains("test"))

		require.Len(t, s.Values(), 2)
	})

	t.Run("remove", func(t *testing.T) {
		s := make(StringSet)
		s.Add("test")
		s.Add("other")

		require.True(t, s.Contains("test"))
		require.True(t, s.Contains("other"))

		s.Remove("test")
		require.False(t, s.Contains("test"))
		require.True(t, s.Contains("other"))

		// Removing non-existent element should not panic
		s.Remove("nonexistent")
		require.True(t, s.Contains("other"))

		require.Len(t, s.Values(), 1)
	})

	t.Run("values", func(t *testing.T) {
		s := make(StringSet)

		// Empty set
		values := s.Values()
		require.Empty(t, values)

		// Add elements
		s.Add("test")
		s.Add("other")
		s.Add("third")

		values = s.Values()
		require.Len(t, values, 3)
		require.Contains(t, values, "test")
		require.Contains(t, values, "other")
		require.Contains(t, values, "third")
	})
}
