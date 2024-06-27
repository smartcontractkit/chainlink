package s4

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/ethereum/go-ethereum/common"
)

// AddressRange represents a range of Ethereum addresses.
type AddressRange struct {
	// MinAddress (inclusive).
	MinAddress *utils.Big
	// MaxAddress (inclusive).
	MaxAddress *utils.Big
}

var (
	ErrInvalidIntervals = errors.New("invalid intervals value")
	MinAddress          = utils.NewBig(common.BytesToAddress(bytes.Repeat([]byte{0x00}, common.AddressLength)).Big())
	MaxAddress          = utils.NewBig(common.BytesToAddress(bytes.Repeat([]byte{0xff}, common.AddressLength)).Big())
)

// NewFullAddressRange creates AddressRange for all address space: 0x00..-0xFF..
func NewFullAddressRange() *AddressRange {
	return &AddressRange{
		MinAddress: MinAddress,
		MaxAddress: MaxAddress,
	}
}

// NewSingleAddressRange creates AddressRange for a single address.
func NewSingleAddressRange(address *utils.Big) (*AddressRange, error) {
	if address == nil || address.Cmp(MinAddress) < 0 || address.Cmp(MaxAddress) > 0 {
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
	maxPlusOne := MaxAddress.Add(utils.NewBigI(1))
	interval := utils.NewBig(new(big.Int).Div(maxPlusOne.ToInt(), divisor))

	return &AddressRange{
		MinAddress: MinAddress,
		MaxAddress: MinAddress.Add(interval).Sub(utils.NewBigI(1)),
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

	r.MinAddress = r.MinAddress.Add(interval)
	r.MaxAddress = r.MaxAddress.Add(interval)

	if r.MinAddress.Cmp(MaxAddress) >= 0 {
		r.MinAddress = MinAddress
		r.MaxAddress = MinAddress.Add(interval).Sub(utils.NewBigI(1))
	}

	if r.MaxAddress.Cmp(MaxAddress) > 0 {
		r.MaxAddress = MaxAddress
	}
}

// Contains returns true if the given address belongs to the range.
func (r *AddressRange) Contains(address *utils.Big) bool {
	if r == nil {
		return false
	}
	return r.MinAddress.Cmp(address) <= 0 && r.MaxAddress.Cmp(address) >= 0
}

// Interval returns the interval between max and min address plus one.
func (r *AddressRange) Interval() *utils.Big {
	if r == nil {
		return nil
	}
	return r.MaxAddress.Sub(r.MinAddress).Add(utils.NewBigI(1))
}
