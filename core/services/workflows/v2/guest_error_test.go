package v2

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestTruncateGuestError(t *testing.T) {
	t.Parallel()

	t.Run("short error is unchanged", func(t *testing.T) {
		t.Parallel()
		require.Empty(t, truncateGuestError(""))
		require.Equal(t, "boom", truncateGuestError("boom"))
		exact := strings.Repeat("a", maxGuestErrorBytes)
		require.Equal(t, exact, truncateGuestError(exact))
	})

	t.Run("long error is cut and records its size", func(t *testing.T) {
		t.Parallel()
		long := strings.Repeat("a", 40<<20)
		got := truncateGuestError(long)
		require.Equal(t, strings.Repeat("a", maxGuestErrorBytes)+"... [truncated, 41943040 bytes total]", got)
	})

	t.Run("cut never splits a multi-byte rune", func(t *testing.T) {
		t.Parallel()
		// "é" is 2 bytes; the odd prefix puts a rune across the cut point.
		long := "x" + strings.Repeat("é", maxGuestErrorBytes)
		got := truncateGuestError(long)
		require.True(t, utf8.ValidString(got))
		require.LessOrEqual(t, len(got), maxGuestErrorBytes+len("... [truncated, 8193 bytes total]"))
	})
}
