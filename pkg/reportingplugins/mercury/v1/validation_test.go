package mercury_v1 //nolint:revive

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	rf := ReportFields{
		CurrentBlockHash: make([]byte, 32),
	}

	t.Run("ValidateCurrentBlock", func(t *testing.T) {
		t.Run("succeeds when validFromBlockNum < current block num", func(t *testing.T) {
			rf.ValidFromBlockNum = 16634363
			rf.CurrentBlockNum = 16634364
			err := ValidateCurrentBlock(rf)
			assert.NoError(t, err)
		})
		t.Run("succeeds when validFromBlockNum is equal to current block number", func(t *testing.T) {
			rf.ValidFromBlockNum = 16634364
			rf.CurrentBlockNum = 16634364
			err := ValidateCurrentBlock(rf)
			assert.NoError(t, err)
		})
		t.Run("zero is ok", func(t *testing.T) {
			rf.ValidFromBlockNum = 0
			rf.CurrentBlockNum = 0
			err := ValidateCurrentBlock(rf)
			assert.NoError(t, err)
		})
		t.Run("errors when validFromBlockNum number < 0", func(t *testing.T) {
			rf.ValidFromBlockNum = -1
			rf.CurrentBlockNum = -1
			err := ValidateCurrentBlock(rf)
			assert.EqualError(t, err, "validFromBlockNum must be >= 0 (got: -1)")
		})
		t.Run("errors when validFrom > block number", func(t *testing.T) {
			rf.CurrentBlockNum = 1
			rf.ValidFromBlockNum = 16634366
			err := ValidateCurrentBlock(rf)
			assert.EqualError(t, err, "validFromBlockNum (Value: 16634366) must be less than or equal to CurrentBlockNum (Value: 1)")
		})
		t.Run("errors when validFrom < 0", func(t *testing.T) {
			rf.ValidFromBlockNum = -1
			err := ValidateCurrentBlock(rf)
			assert.EqualError(t, err, "validFromBlockNum must be >= 0 (got: -1)")
		})
		t.Run("errors when hash has incorrect length", func(t *testing.T) {
			rf.ValidFromBlockNum = 16634363
			rf.CurrentBlockNum = 16634364
			rf.CurrentBlockHash = []byte{}
			err := ValidateCurrentBlock(rf)
			assert.EqualError(t, err, "invalid length for hash; expected 32 (got: 0)")
			rf.CurrentBlockHash = make([]byte, 64)
			err = ValidateCurrentBlock(rf)
			assert.EqualError(t, err, "invalid length for hash; expected 32 (got: 64)")
		})
	})
}
