package s4_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	bigmath "github.com/smartcontractkit/chainlink-common/pkg/utils/big_math"

	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
)

func TestAddressRange_NewFullAddressRange(t *testing.T) {
	t.Parallel()

	full := s4.NewFullAddressRange()
	assert.Equal(t, s4.MinAddress, full.MinAddress)
	assert.Equal(t, s4.MaxAddress, full.MaxAddress)

	t.Run("advance has no effect", func(t *testing.T) {
		full.Advance()
		assert.Equal(t, s4.MinAddress, full.MinAddress)
		assert.Equal(t, s4.MaxAddress, full.MaxAddress)
	})
}

func TestAddressRange_NewSingleAddressRange(t *testing.T) {
	t.Parallel()

	addr := sqlutil.NewI(0x123)
	sar, err := s4.NewSingleAddressRange(addr)
	assert.NoError(t, err)
	assert.Equal(t, addr, sar.MinAddress)
	assert.Equal(t, addr, sar.MaxAddress)
	assert.True(t, sar.Contains(addr))
	assert.Equal(t, int64(1), sar.Interval().ToInt().Int64())

	sar.Advance()
	assert.False(t, sar.Contains(addr))
}

func TestAddressRange_NewInitialAddressRangeForIntervals(t *testing.T) {
	t.Parallel()

	t.Run("invalid intervals", func(t *testing.T) {
		_, err := s4.NewInitialAddressRangeForIntervals(0)
		assert.ErrorIs(t, err, s4.ErrInvalidIntervals)

		_, err = s4.NewInitialAddressRangeForIntervals(3)
		assert.ErrorIs(t, err, s4.ErrInvalidIntervals)
	})

	t.Run("full range for one interval", func(t *testing.T) {
		r, err := s4.NewInitialAddressRangeForIntervals(1)
		assert.NoError(t, err)
		assert.Equal(t, s4.NewFullAddressRange(), r)
	})

	t.Run("initial range for 256 intervals", func(t *testing.T) {
		r, err := s4.NewInitialAddressRangeForIntervals(256)
		assert.NoError(t, err)
		assert.Equal(t, "0x0", hex(r.MinAddress))
		assert.Equal(t, "0xffffffffffffffffffffffffffffffffffffff", hex(r.MaxAddress))
	})

	t.Run("advance for 256 intervals", func(t *testing.T) {
		r, err := s4.NewInitialAddressRangeForIntervals(256)
		assert.NoError(t, err)

		r.Advance()
		assert.Equal(t, "0x100000000000000000000000000000000000000", hex(r.MinAddress))
		assert.Equal(t, "0x1ffffffffffffffffffffffffffffffffffffff", hex(r.MaxAddress))

		r.Advance()
		assert.Equal(t, "0x200000000000000000000000000000000000000", hex(r.MinAddress))
		assert.Equal(t, "0x2ffffffffffffffffffffffffffffffffffffff", hex(r.MaxAddress))

		for range 253 {
			r.Advance()
		}
		assert.Equal(t, "0xff00000000000000000000000000000000000000", hex(r.MinAddress))
		assert.Equal(t, "0xffffffffffffffffffffffffffffffffffffffff", hex(r.MaxAddress))

		// initial
		r.Advance()
		assert.Equal(t, s4.MinAddress, r.MinAddress)
		assert.Equal(t, "0xffffffffffffffffffffffffffffffffffffff", hex(r.MaxAddress))
	})
}

func TestAddressRange_Contains(t *testing.T) {
	t.Parallel()

	r, err := s4.NewInitialAddressRangeForIntervals(256)
	assert.NoError(t, err)
	assert.True(t, r.Contains(r.MinAddress))
	assert.True(t, r.Contains(r.MaxAddress))
	assert.False(t, r.Contains(sqlutil.New(bigmath.Add(r.MaxAddress.ToInt(), big.NewInt(1)))))

	r.Advance()
	assert.True(t, r.Contains(r.MinAddress))
	assert.True(t, r.Contains(r.MaxAddress))
	assert.False(t, r.Contains(sqlutil.New(bigmath.Sub(r.MinAddress.ToInt(), big.NewInt(1)))))
}

func hex(b *sqlutil.Big) string {
	return hexutil.EncodeBig(b.ToInt())
}
