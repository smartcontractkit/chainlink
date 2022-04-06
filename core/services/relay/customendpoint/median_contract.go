package customendpoint

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// This returns the latest observations that were last uploaded to the endpoint.
// All transmissions are saved in the ContractTracker, during call to Transmit().
func (c *ContractTracker) LatestTransmissionDetails(
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
	answer := c.GetLastTransmittedAnswer()
	return digester,
		answer.epoch,
		answer.round,
		answer.Data,
		answer.Timestamp,
		err
}

// It is safe to return 0 values here.
func (c *ContractTracker) LatestRoundRequested(
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
