package mercury

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	min := big.NewInt(0)
	max := big.NewInt(10_000)

	badMin := big.NewInt(9_000)
	badMax := big.NewInt(10)

	t.Run("ValidateValidFromTimestamp", func(t *testing.T) {
		t.Run("succeeds when observationTimestamp is >= validFromTimestamp", func(t *testing.T) {
			err := ValidateValidFromTimestamp(456, 123)
			assert.NoError(t, err)
			err = ValidateValidFromTimestamp(123, 123)
			assert.NoError(t, err)
		})
		t.Run("fails when observationTimestamp is < validFromTimestamp", func(t *testing.T) {
			err := ValidateValidFromTimestamp(111, 112)
			assert.EqualError(t, err, "observationTimestamp (Value: 111) must be >= validFromTimestamp (Value: 112)")
		})
	})
	t.Run("ValidateExpiresAt", func(t *testing.T) {
		t.Run("succeeds when observationTimestamp <= expiresAt", func(t *testing.T) {
			err := ValidateExpiresAt(123, 456)
			assert.NoError(t, err)
			err = ValidateExpiresAt(123, 123)
			assert.NoError(t, err)
		})

		t.Run("fails when observationTimestamp > expiresAt", func(t *testing.T) {
			err := ValidateExpiresAt(112, 111)
			assert.EqualError(t, err, "expiresAt (Value: 111) must be ahead of observation timestamp (Value: 112)")
		})
	})
	t.Run("ValidateBetween", func(t *testing.T) {
		bm := big.NewInt(346)
		err := ValidateBetween("test foo", bm, min, max)
		assert.NoError(t, err)

		err = ValidateBetween("test bar", bm, min, badMax)
		assert.EqualError(t, err, "test bar (Value: 346) is outside of allowable range (Min: 0, Max: 10)")
		err = ValidateBetween("test baz", bm, badMin, max)
		assert.EqualError(t, err, "test baz (Value: 346) is outside of allowable range (Min: 9000, Max: 10000)")
	})
}
