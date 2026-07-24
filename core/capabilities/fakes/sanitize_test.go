package fakes

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestSanitizeUTF8(t *testing.T) {
	t.Parallel()
	t.Run("returns valid strings unchanged", func(t *testing.T) {
		t.Parallel()
		for _, s := range []string{"", "ascii", "héllo wörld", "日本語", "emoji 🚀"} {
			require.Equal(t, s, sanitizeUTF8(s))
		}
	})

	t.Run("replaces invalid bytes with U+FFFD", func(t *testing.T) {
		t.Parallel()
		got := sanitizeUTF8("a" + string([]byte{0xff}) + "b")
		require.True(t, utf8.ValidString(got))
		require.Equal(t, "a�b", got)
	})
}
