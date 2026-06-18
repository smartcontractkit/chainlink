package types

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArtifactFetchError(t *testing.T) {
	inner := errors.New("connection refused")
	fetchErr := &ArtifactFetchError{
		ArtifactType: "binary",
		URL:          "https://storage.example.com/artifacts/abc123/binary.wasm?Expires=123&Signature=xyz",
		Err:          inner,
	}

	t.Run("Error preserves full URL for internal debugging", func(t *testing.T) {
		assert.Contains(t, fetchErr.Error(), "Expires=123&Signature=xyz")
		assert.Contains(t, fetchErr.Error(), "binary.wasm")
		assert.Contains(t, fetchErr.Error(), "connection refused")
	})

	t.Run("CustomerError is deterministic and omits URL details", func(t *testing.T) {
		msg := fetchErr.CustomerError()
		assert.Equal(t, "Internal error: failed to fetch workflow binary from storage. Contact support if this persists.", msg)
		assert.NotContains(t, msg, "Expires")
		assert.NotContains(t, msg, "Signature")
		assert.NotContains(t, msg, "example.com")
	})

	t.Run("CustomerError reflects artifact type", func(t *testing.T) {
		configErr := &ArtifactFetchError{ArtifactType: "config", URL: "https://x.com/c?s=1", Err: inner}
		assert.Contains(t, configErr.CustomerError(), "workflow config")
	})

	t.Run("Unwrap returns inner error", func(t *testing.T) {
		require.ErrorIs(t, fetchErr, inner)
	})

	t.Run("errors.As matches through wrapping", func(t *testing.T) {
		wrapped := fmt.Errorf("outer: %w", fetchErr)
		var target *ArtifactFetchError
		require.ErrorAs(t, wrapped, &target)
		assert.Equal(t, "binary", target.ArtifactType)
	})
}
