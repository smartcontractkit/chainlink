package customendpoint

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// This returns the latest observations that were last uploaded to the endpoint.
// All transmissions are saved in the contractTracker, during call to Transmit().
func (c *contractTracker) LatestTransmissionDetails(
	ctx context.Context,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	latestAnswer *big.Int,
	latestTimestamp time.Time,
	err error,
) {
	digester, err := c.digester.configDigest()
	storedAnswer := c.getLastTransmittedAnswer()
	return digester,
		storedAnswer.epoch,
		storedAnswer.round,
		storedAnswer.Data,
		storedAnswer.Timestamp,
		err
}

// It is safe to return 0 values here.
func (c *contractTracker) LatestRoundRequested(
	ctx context.Context,
	lookback time.Duration,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	err error,
) {
	digester, err := c.digester.configDigest()
	return digester, 0, 0, err
}
