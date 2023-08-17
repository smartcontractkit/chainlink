package mercury

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateValidFromTimestamp(t *testing.T) {
	t.Run("succeeds when observationTimestamp is >= validFromTimestamp", func(t *testing.T) {
		err := ValidateValidFromTimestamp(456, 123)
		assert.NoError(t, err)
	})
	t.Run("fails when observationTimestamp is < validFromTimestamp", func(t *testing.T) {
		err := ValidateValidFromTimestamp(111, 112)
		assert.EqualError(t, err, "observationTimestamp (111) must be >= validFromTimestamp (112)")
	})
}

func TestValidateExpiresAt(t *testing.T) {
	t.Run("succeeds when no overflow occurs", func(t *testing.T) {
		err := ValidateExpiresAt(456, 123)
		assert.NoError(t, err)
	})

	t.Run("fails when overflow occurs", func(t *testing.T) {
		err := ValidateExpiresAt(math.MaxUint32, 1)
		assert.EqualError(t, err, "timestamp 4294967295 + expiration window 1 overflows uint32")
	})
}
