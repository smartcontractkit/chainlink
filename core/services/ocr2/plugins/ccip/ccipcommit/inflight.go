package ccipcommit

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
)

const (
	PRICE_EXPIRY_MULTIPLIER = 3 // Keep price update cache longer and use it as source of truth
)

// InflightCommitReport represents a commit report which has been submitted
// to the transaction manager and we expect to be included in the chain.
// By keeping track of the inflight reports, we are able to build subsequent
// reports "on top" of the inflight ones for improved throughput - for example
// if seqNrs=[1,2] are inflight, we can build and send [3,4] while [1,2] is still confirming
// and optimistically assume they will complete in order. If for whatever reason (re-org or
// RPC timing) leads to [3,4] arriving before [1,2], we'll revert onchain. Once the cache
// expires we'll then build from the onchain state again and retries. In this manner,
// we are able to obtain high throughput during happy path yet still naturally recover
// if a reorg or issue causes onchain reverts.
type InflightCommitReport struct {
	report    cciptypes.CommitStoreReport
	createdAt time.Time
}

type InflightPriceUpdate struct {
	gasPrices     []cciptypes.GasPrice
	tokenPrices   []cciptypes.TokenPrice
	createdAt     time.Time
	epochAndRound uint64
}

// inflightExecReportsContainer holds existing inflight reports.
// it provides a thread-safe access as it is called from multiple goroutines,
// e.g. reporting and transmission protocols.
type inflightCommitReportsContainer struct {
	locker               sync.RWMutex
	inFlight             map[[32]byte]InflightCommitReport
	inFlightPriceUpdates []InflightPriceUpdate
	cacheExpiry          time.Duration
}

func newInflightCommitReportsContainer(inflightCacheExpiry time.Duration) *inflightCommitReportsContainer {
	return &inflightCommitReportsContainer{
		locker:               sync.RWMutex{},
		inFlight:             make(map[[32]byte]InflightCommitReport),
		inFlightPriceUpdates: []InflightPriceUpdate{},
		cacheExpiry:          inflightCacheExpiry,
	}
}

func (c *inflightCommitReportsContainer) maxInflightSeqNr() uint64 {
	c.locker.RLock()
	defer c.locker.RUnlock()
	var maxSeqNr uint64
	for _, report := range c.inFlight {
		if report.report.Interval.Max >= maxSeqNr {
			maxSeqNr = report.report.Interval.Max
		}
	}
	return maxSeqNr
}

// latestInflightGasPriceUpdates returns a map of the latest gas price updates.
func (c *inflightCommitReportsContainer) latestInflightGasPriceUpdates() map[uint64]update {
	c.locker.RLock()
	defer c.locker.RUnlock()
	latestGasPriceUpdates := make(map[uint64]update)
	latestEpochAndRounds := make(map[uint64]uint64)

	for _, inflight := range c.inFlightPriceUpdates {
		for _, inflightGasUpdate := range inflight.gasPrices {
			_, ok := latestGasPriceUpdates[inflightGasUpdate.DestChainSelector]
			if !ok || inflight.epochAndRound > latestEpochAndRounds[inflightGasUpdate.DestChainSelector] {
				latestGasPriceUpdates[inflightGasUpdate.DestChainSelector] = update{
					value:     inflightGasUpdate.Value,
					timestamp: inflight.createdAt,
				}
				latestEpochAndRounds[inflightGasUpdate.DestChainSelector] = inflight.epochAndRound
			}
		}
	}

	return latestGasPriceUpdates
}

// latestInflightTokenPriceUpdates returns a map of the latest token price updates
func (c *inflightCommitReportsContainer) latestInflightTokenPriceUpdates() map[cciptypes.Address]update {
	c.locker.RLock()
	defer c.locker.RUnlock()
	latestTokenPriceUpdates := make(map[cciptypes.Address]update)
	latestEpochAndRounds := make(map[cciptypes.Address]uint64)
	for _, inflight := range c.inFlightPriceUpdates {
		for _, inflightTokenUpdate := range inflight.tokenPrices {
			_, ok := latestTokenPriceUpdates[inflightTokenUpdate.Token]
			if !ok || inflight.epochAndRound > latestEpochAndRounds[inflightTokenUpdate.Token] {
				latestTokenPriceUpdates[inflightTokenUpdate.Token] = update{
					value:     inflightTokenUpdate.Value,
					timestamp: inflight.createdAt,
				}
				latestEpochAndRounds[inflightTokenUpdate.Token] = inflight.epochAndRound
			}
		}
	}
	return latestTokenPriceUpdates
}

func (c *inflightCommitReportsContainer) reset(lggr logger.Logger) {
	lggr.Infow("Inflight report reset")
	c.locker.Lock()
	defer c.locker.Unlock()
	c.inFlight = make(map[[32]byte]InflightCommitReport)
	c.inFlightPriceUpdates = []InflightPriceUpdate{}
}

func (c *inflightCommitReportsContainer) expire(lggr logger.Logger) {
	c.locker.Lock()
	defer c.locker.Unlock()
	// Reap any expired entries from inflight.
	for root, inFlightReport := range c.inFlight {
		if time.Since(inFlightReport.createdAt) > c.cacheExpiry {
			// Happy path: inflight report was successfully transmitted onchain, we remove it from inflight and onchain state reflects inflight.
			// Sad path: inflight report reverts onchain, we remove it from inflight, onchain state does not reflect the chains so we retry.
			lggr.Infow("Inflight report expired", "rootOfRoots", hexutil.Encode(inFlightReport.report.MerkleRoot[:]))
			delete(c.inFlight, root)
		}
	}

	lggr.Infow("Inflight expire with price count", "count", len(c.inFlightPriceUpdates))

	var stillInflight []InflightPriceUpdate
	for _, inFlightFeeUpdate := range c.inFlightPriceUpdates {
		timeSinceUpdate := time.Since(inFlightFeeUpdate.createdAt)
		// If time passed since the price update is greater than multiplier * cache expiry, we remove it from the inflight list.
		if timeSinceUpdate > c.cacheExpiry*PRICE_EXPIRY_MULTIPLIER {
			// Happy path: inflight report was successfully transmitted onchain, we remove it from inflight and onchain state reflects inflight.
			// Sad path: inflight report reverts onchain, we remove it from inflight, onchain state does not reflect the chains, so we retry.
			lggr.Infow("Inflight price update expired", "gasPrices", inFlightFeeUpdate.gasPrices, "tokenPrices", inFlightFeeUpdate.tokenPrices)
		} else {
			// If the update is still valid, we keep it in the inflight list.
			stillInflight = append(stillInflight, inFlightFeeUpdate)
		}
	}
	c.inFlightPriceUpdates = stillInflight
}

func (c *inflightCommitReportsContainer) add(lggr logger.Logger, report cciptypes.CommitStoreReport, epochAndRound uint64) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	if report.MerkleRoot != [32]byte{} {
		// Set new inflight ones as pending
		lggr.Infow("Adding to inflight report", "rootOfRoots", hexutil.Encode(report.MerkleRoot[:]))
		c.inFlight[report.MerkleRoot] = InflightCommitReport{
			report:    report,
			createdAt: time.Now(),
		}
	}

	if len(report.GasPrices) != 0 || len(report.TokenPrices) != 0 {
		lggr.Infow("Adding to inflight fee updates", "gasPrices", report.GasPrices, "tokenPrices", report.TokenPrices)
		c.inFlightPriceUpdates = append(c.inFlightPriceUpdates, InflightPriceUpdate{
			gasPrices:     report.GasPrices,
			tokenPrices:   report.TokenPrices,
			createdAt:     time.Now(),
			epochAndRound: epochAndRound,
		})
	}
	return nil
}
