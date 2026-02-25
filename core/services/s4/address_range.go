package s4

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	bigmath "github.com/smartcontractkit/chainlink-common/pkg/utils/big_math"
)

// AddressRange represents a range of Ethereum addresses.
type AddressRange struct {
	// MinAddress (inclusive).
	MinAddress *sqlutil.Big
	// MaxAddress (inclusive).
	MaxAddress *sqlutil.Big
}

var (
	ErrInvalidIntervals = errors.New("invalid intervals value")
	MinAddress          = sqlutil.New(common.BytesToAddress(bytes.Repeat([]byte{0x00}, common.AddressLength)).Big())
	MaxAddress          = sqlutil.New(common.BytesToAddress(bytes.Repeat([]byte{0xff}, common.AddressLength)).Big())
)

// NewFullAddressRange creates AddressRange for all address space: 0x00..-0xFF..
func NewFullAddressRange() *AddressRange {
	return &AddressRange{
		MinAddress: MinAddress,
		MaxAddress: MaxAddress,
	}
}

// NewSingleAddressRange creates AddressRange for a single address.
func NewSingleAddressRange(address *sqlutil.Big) (*AddressRange, error) {
	if address == nil || address.ToInt().Cmp(MinAddress.ToInt()) < 0 || address.ToInt().Cmp(MaxAddress.ToInt()) > 0 {
		return nil, errors.New("invalid address")
	}
	return &AddressRange{
		MinAddress: address,
		MaxAddress: address,
	}, nil
}

// NewInitialAddressRangeForIntervals splits the full address space with intervals,
// and returns a range for the first interval.
// Number of intervals must be > 0 and a power of 2.
func NewInitialAddressRangeForIntervals(intervals uint) (*AddressRange, error) {
	if intervals == 0 || (intervals&(intervals-1) != 0) {
		return nil, ErrInvalidIntervals
	}

	if intervals == 1 {
		return NewFullAddressRange(), nil
	}

	divisor := big.NewInt(int64(intervals))
	maxPlusOne := bigmath.Add(MaxAddress.ToInt(), big.NewInt(1))
	interval := bigmath.Div(maxPlusOne, divisor)

	return &AddressRange{
		MinAddress: MinAddress,
		MaxAddress: sub(add(MinAddress, sqlutil.New(interval)), sqlutil.NewI(1)),
	}, nil
}

// Advances the AddressRange by r.Interval. Has no effect for NewFullAddressRange().
// When it reaches the end of the address space, it resets to the initial state,
// returned by NewAddressRangeForFirstInterval().
func (r *AddressRange) Advance() {
	if r == nil {
		return
	}

	interval := r.Interval()

	r.MinAddress = add(r.MinAddress, interval)
	r.MaxAddress = add(r.MaxAddress, interval)

	if r.MinAddress.ToInt().Cmp(MaxAddress.ToInt()) >= 0 {
		r.MinAddress = MinAddress
		r.MaxAddress = sub(add(MinAddress, interval), sqlutil.NewI(1))
	}

	if r.MaxAddress.ToInt().Cmp(MaxAddress.ToInt()) > 0 {
		r.MaxAddress = MaxAddress
	}
}

// Contains returns true if the given address belongs to the range.
func (r *AddressRange) Contains(address *sqlutil.Big) bool {
	if r == nil {
		return false
	}
	return r.MinAddress.ToInt().Cmp(address.ToInt()) <= 0 && r.MaxAddress.ToInt().Cmp(address.ToInt()) >= 0
}

// Interval returns the interval between max and min address plus one.
func (r *AddressRange) Interval() *sqlutil.Big {
	if r == nil {
		return nil
	}
	return add(sub(r.MaxAddress, r.MinAddress), sqlutil.NewI(1))
}

func sub(a, b *sqlutil.Big) *sqlutil.Big {
	return sqlutil.New(bigmath.Sub(a.ToInt(), b.ToInt()))
}

func add(a, b *sqlutil.Big) *sqlutil.Big {
	return sqlutil.New(bigmath.Add(a.ToInt(), b.ToInt()))
}
