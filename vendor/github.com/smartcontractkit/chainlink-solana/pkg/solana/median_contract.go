package solana

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type MedianContract struct {
	stateCache         *StateCache
	transmissionsCache *TransmissionsCache
}

func (c *MedianContract) LatestTransmissionDetails(
	ctx context.Context,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	latestAnswer *big.Int,
	latestTimestamp time.Time,
	err error,
) {
	state, err := c.stateCache.ReadState()
	if err != nil {
		return configDigest, epoch, round, latestAnswer, latestTimestamp, err
	}
	answer, err := c.transmissionsCache.ReadAnswer()
	if err != nil {
		return configDigest, epoch, round, latestAnswer, latestTimestamp, err
	}

	configDigest = state.Config.LatestConfigDigest
	epoch = state.Config.Epoch
	round = state.Config.Round
	latestAnswer = answer.Data
	latestTimestamp = time.Unix(int64(answer.Timestamp), 0)
	return configDigest, epoch, round, latestAnswer, latestTimestamp, nil
}

// LatestRoundRequested returns the configDigest, epoch, and round from the latest
// RoundRequested event emitted by the contract. LatestRoundRequested may or may not
// return a result if the latest such event was emitted in a block b such that
// b.timestamp < tip.timestamp - lookback.
//
// If no event is found, LatestRoundRequested should return zero values, not an error.
// An error should only be returned if an actual error occurred during execution,
// e.g. because there was an error querying the blockchain or the database.
//
// As an optimization, this function may also return zero values, if no
// RoundRequested event has been emitted after the latest NewTransmission event.
func (c *MedianContract) LatestRoundRequested(
	ctx context.Context,
	lookback time.Duration,
) (
	configDigest types.ConfigDigest,
	epoch uint32,
	round uint8,
	err error,
) {
	state, err := c.stateCache.ReadState()
	return state.Config.LatestConfigDigest, 0, 0, err
}
