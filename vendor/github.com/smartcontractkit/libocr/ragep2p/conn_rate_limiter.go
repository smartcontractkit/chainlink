package ragep2p

import (
	"math"
	"sync"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/ragep2p/internal/ratelimit"
)

const tlsFactor = 1.5

type connRateLimiter struct {
	logger loghelper.LoggerWithContext

	mutex        sync.Mutex
	limiter      *ratelimit.TokenBucket
	infiniteRate bool
}

func newConnRateLimiter(logger loghelper.LoggerWithContext) *connRateLimiter {
	return &connRateLimiter{
		logger,

		sync.Mutex{},
		&ratelimit.TokenBucket{},
		false,
	}
}

func (crl *connRateLimiter) Allow(n int) bool {
	crl.mutex.Lock()
	defer crl.mutex.Unlock()

	return crl.infiniteRate || (uint64(n) <= uint64(math.MaxUint32) && crl.limiter.RemoveTokens(uint32(n)))
}

func (crl *connRateLimiter) AddStream(messagesLimit TokenBucketParams, bytesLimit TokenBucketParams) {
	crl.mutex.Lock()
	defer crl.mutex.Unlock()

	crl.addRemoveStream(true, messagesLimit, bytesLimit)
}

func (crl *connRateLimiter) RemoveStream(messagesLimit TokenBucketParams, bytesLimit TokenBucketParams) {
	crl.mutex.Lock()
	defer crl.mutex.Unlock()

	crl.addRemoveStream(false, messagesLimit, bytesLimit)
}

func (crl *connRateLimiter) AddTokens(n uint32) {
	crl.mutex.Lock()
	defer crl.mutex.Unlock()

	crl.limiter.AddTokens(n)
}

func (crl *connRateLimiter) TokenBucketParams() TokenBucketParams {
	crl.mutex.Lock()
	defer crl.mutex.Unlock()

	return TokenBucketParams{
		float64(crl.limiter.Rate()) / 1000, // millitokens per second -> tokens per second
		crl.limiter.Capacity(),
	}
}

func (crl *connRateLimiter) addRemoveStream(add bool, messagesLimit TokenBucketParams, bytesLimit TokenBucketParams) {
	if crl.infiniteRate {
		// we're already in absorbing overflow state, nothing to be done
		return
	}

	deltaRate, deltaCapacity, ok := delta(messagesLimit, bytesLimit)
	if !ok {
		crl.logger.Warn("connRateLimiter entered overflow state after delta() indicated a problem", commontypes.LogFields{
			"messagesLimit": messagesLimit,
			"bytesLimit":    bytesLimit,
		})
		// enter overflow state
		crl.infiniteRate = true
		return
	}

	if !add {
		deltaRate, deltaCapacity = -deltaRate, -deltaCapacity
	}

	oldRate := crl.limiter.Rate()
	newRate := oldRate + deltaRate

	oldCapacity := crl.limiter.Capacity()
	newCapacity := oldCapacity + deltaCapacity

	if add && (newRate < oldRate || newCapacity < oldCapacity) {
		crl.logger.Warn("connRateLimiter entered overflow state after rate or capacity overflow", commontypes.LogFields{
			"messagesLimit": messagesLimit,
			"bytesLimit":    bytesLimit,
			"newRate":       newRate,
			"oldRate":       oldRate,
			"newCapacity":   newCapacity,
			"oldCapacity":   oldCapacity,
		})
		// enter overflow state
		crl.infiniteRate = true
		return
	}

	if !add && (newRate > oldRate || newCapacity > oldCapacity) {
		crl.logger.Warn("connRateLimiter entered overflow state after rate or capacity underflow", commontypes.LogFields{
			"messagesLimit": messagesLimit,
			"bytesLimit":    bytesLimit,
			"newRate":       newRate,
			"oldRate":       oldRate,
			"newCapacity":   newCapacity,
			"oldCapacity":   oldCapacity,
		})
		// enter overflow state
		crl.infiniteRate = true
		return
	}

	crl.limiter.SetRate(newRate)
	crl.limiter.SetCapacity(newCapacity)
}

func delta(messagesLimit TokenBucketParams, bytesLimit TokenBucketParams) (deltaRate ratelimit.MillitokensPerSecond, deltaCapacity uint32, ok bool) {
	var deltaRateF float64
	deltaRateF = messagesLimit.Rate * frameHeaderEncodedSize
	deltaRateF += bytesLimit.Rate
	deltaRateF *= tlsFactor
	deltaRateF = math.Ceil(deltaRateF * 1000)
	deltaRate = ratelimit.MillitokensPerSecond(deltaRateF)

	var deltaCapacityF float64
	deltaCapacityF = float64(messagesLimit.Capacity) * frameHeaderEncodedSize
	deltaCapacityF += float64(bytesLimit.Capacity)
	deltaCapacityF *= tlsFactor
	deltaCapacityF = math.Ceil(deltaCapacityF)
	deltaCapacity = uint32(deltaCapacityF)

	ok = float64(deltaRate) == deltaRateF && float64(deltaCapacity) == deltaCapacityF
	return
}
