package shared

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTokenPoolLookupTableQualifier verifies the composite LUT qualifier is unique per
// (token mint, pool type, metadata).
func TestTokenPoolLookupTableQualifier(t *testing.T) {
	t.Parallel()
	mint := "So11111111111111111111111111111111111111112"
	a := TokenPoolLookupTableQualifier(mint, "BurnMintTokenPool", "CLL")
	b := TokenPoolLookupTableQualifier(mint, "LockReleaseTokenPool", "CLL")
	c := TokenPoolLookupTableQualifier(mint, "BurnMintTokenPool", "customPool9")
	require.NotEqual(t, a, b)
	require.NotEqual(t, a, c)
	require.NotEqual(t, b, c)
}

// TestTokenPoolLookupTableQualifierIsUnambiguous verifies that a separator appearing inside a
// caller-supplied component cannot make two distinct (mint, poolType, metadata) triples collide.
func TestTokenPoolLookupTableQualifierIsUnambiguous(t *testing.T) {
	t.Parallel()
	mint := "So11111111111111111111111111111111111111112"
	require.NotEqual(t,
		TokenPoolLookupTableQualifier(mint, "BurnMintTokenPool", "a/b"),
		TokenPoolLookupTableQualifier(mint, "BurnMintTokenPool/a", "b"),
	)
}
