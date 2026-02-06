package fakes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	confidentialhttp "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialhttp"
)

func TestHasEncryptionSecret(t *testing.T) {
	t.Run("returns true when magic key exists", func(t *testing.T) {
		secrets := []*confidentialhttp.SecretIdentifier{
			{Key: "other-key"},
			{Key: AESGCMEncryptionKeyName},
		}
		assert.True(t, hasEncryptionSecret(secrets))
	})

	t.Run("returns false when magic key does not exist", func(t *testing.T) {
		secrets := []*confidentialhttp.SecretIdentifier{
			{Key: "other-key"},
			{Key: "another-key"},
		}
		assert.False(t, hasEncryptionSecret(secrets))
	})

	t.Run("returns false for empty secrets", func(t *testing.T) {
		assert.False(t, hasEncryptionSecret(nil))
		assert.False(t, hasEncryptionSecret([]*confidentialhttp.SecretIdentifier{}))
	})
}
