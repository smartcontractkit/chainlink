package ratelimit

import (
	"math"
	"time"
)

type MillitokensPerSecond uint64

// A token bucket rate limiter. Can store at most UINT32_MAX tokens before it
// saturates.
// The zero TokenBucket{} is a valid value.
//
// NOT thread-safe
type TokenBucket struct {
	rate     MillitokensPerSecond
	capacity uint32

	nanotokens uint64
	updated    time.Time
}

func NewTokenBucket(rate MillitokensPerSecond, capacity uint32, full bool) *TokenBucket {
	tb := &TokenBucket{}
	tb.SetRate(rate)
	tb.SetCapacity(capacity)
	if full {
		tb.AddTokens(capacity)
	}
	return tb
}

func (tb *TokenBucket) update(now time.Time) {
	if now.Before(tb.updated) { // we assume that time moves forward monotonically
		now = tb.updated
	}

	// We round down the time difference to the nearest microsecond
	// and store the remainder in the updated timestamp
	timeDiff := now.Sub(tb.updated)
	timeDiffMicroseconds := uint64(timeDiff / time.Microsecond)
	timeDiffRemainder := timeDiff % time.Microsecond
	tb.updated = now.Add(-timeDiffRemainder)

	// millitokens per second x microseconds yields nanotokens
	nanotokensDiff := timeDiffMicroseconds * uint64(tb.rate)
	if timeDiffMicroseconds != 0 && nanotokensDiff/uint64(timeDiffMicroseconds) != uint64(tb.rate) {
		// multiplication overflow
		tb.nanotokens = math.MaxUint64
	} else {
		newNanotokens := tb.nanotokens + nanotokensDiff
		if newNanotokens < tb.nanotokens {
			// addition overflow
			tb.nanotokens = math.MaxUint64
		} else {
			tb.nanotokens = newNanotokens
		}
	}

	capacityNanotokens := uint64(tb.capacity) * 1_000_000_000
	if tb.nanotokens > capacityNanotokens {
		tb.nanotokens = capacityNanotokens
	}
}

// Adds n tokens to the bucket.
func (tb *TokenBucket) AddTokens(n uint32) {
	newNanotokens := tb.nanotokens + uint64(n)*1_000_000_000
	if newNanotokens < tb.nanotokens {
		// addition overflow
		tb.nanotokens = math.MaxUint64
	} else {
		tb.nanotokens = newNanotokens
	}
}

func (tb *TokenBucket) removeTokens(now time.Time, n uint32) bool {
	tb.update(now)
	if tb.nanotokens >= uint64(n)*1_000_000_000 {
		tb.nanotokens -= uint64(n) * 1_000_000_000
		return true
	} else {
		tb.nanotokens = 0
		return false
	}
}

// Removes n tokens from the bucket. If the bucket contained at least n tokens,
// return true. Otherwise, returns false and sets bucket to contain zero tokens.
func (tb *TokenBucket) RemoveTokens(n uint32) bool {
	return tb.removeTokens(time.Now(), n)
}

func (tb *TokenBucket) setRate(now time.Time, rate MillitokensPerSecond) {
	tb.update(now)
	tb.rate = rate
}

// Sets the rate at which the bucket fills with tokens.
func (tb *TokenBucket) SetRate(rate MillitokensPerSecond) {
	tb.setRate(time.Now(), rate)
}

func (tb *TokenBucket) Rate() MillitokensPerSecond {
	return tb.rate
}

// Sets the bucket's capacity.
func (tb *TokenBucket) SetCapacity(capacity uint32) {
	tb.capacity = capacity
}

func (tb *TokenBucket) Capacity() uint32 {
	return tb.capacity
}
