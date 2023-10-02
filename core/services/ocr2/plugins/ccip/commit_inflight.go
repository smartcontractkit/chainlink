package ccip

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	report    commit_store.CommitStoreCommitReport
	createdAt time.Time
}

type InflightPriceUpdate struct {
	priceUpdates  commit_store.InternalPriceUpdates
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
	var max uint64
	for _, report := range c.inFlight {
		if report.report.Interval.Max >= max {
			max = report.report.Interval.Max
		}
	}
	return max
}

// getLatestInflightGasPriceUpdate returns the latest inflight gas price update, and bool flag on if update exists.
func (c *inflightCommitReportsContainer) getLatestInflightGasPriceUpdate() (update, bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	updateFound := false
	latestGasPriceUpdate := update{}
	var latestEpochAndRound uint64
	for _, inflight := range c.inFlightPriceUpdates {
		if inflight.priceUpdates.DestChainSelector == 0 {
			// Price updates did not include a gas price
			continue
		}
		if !updateFound || inflight.epochAndRound > latestEpochAndRound {
			// First price found or found later update, set it
			updateFound = true
			latestGasPriceUpdate = update{
				timestamp: inflight.createdAt,
				value:     inflight.priceUpdates.UsdPerUnitGas,
			}
			latestEpochAndRound = inflight.epochAndRound
			continue
		}
	}
	return latestGasPriceUpdate, updateFound
}

// latestInflightTokenPriceUpdates returns a map of the latest token price updates
func (c *inflightCommitReportsContainer) latestInflightTokenPriceUpdates() map[common.Address]update {
	c.locker.RLock()
	defer c.locker.RUnlock()
	latestTokenPriceUpdates := make(map[common.Address]update)
	latestEpochAndRounds := make(map[common.Address]uint64)
	for _, inflight := range c.inFlightPriceUpdates {
		for _, inflightTokenUpdate := range inflight.priceUpdates.TokenPriceUpdates {
			if _, ok := latestTokenPriceUpdates[inflightTokenUpdate.SourceToken]; !ok {
				latestTokenPriceUpdates[inflightTokenUpdate.SourceToken] = update{
					value:     inflightTokenUpdate.UsdPerToken,
					timestamp: inflight.createdAt,
				}
				latestEpochAndRounds[inflightTokenUpdate.SourceToken] = inflight.epochAndRound
			}
			if inflight.epochAndRound > latestEpochAndRounds[inflightTokenUpdate.SourceToken] {
				latestTokenPriceUpdates[inflightTokenUpdate.SourceToken] = update{
					value:     inflightTokenUpdate.UsdPerToken,
					timestamp: inflight.createdAt,
				}
				latestEpochAndRounds[inflightTokenUpdate.SourceToken] = inflight.epochAndRound
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
			lggr.Infow("Inflight price update expired", "updates", inFlightFeeUpdate.priceUpdates)
		} else {
			// If the update is still valid, we keep it in the inflight list.
			stillInflight = append(stillInflight, inFlightFeeUpdate)
		}
	}
	c.inFlightPriceUpdates = stillInflight
}

func (c *inflightCommitReportsContainer) add(lggr logger.Logger, report commit_store.CommitStoreCommitReport, epochAndRound uint64) error {
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

	if report.PriceUpdates.DestChainSelector != 0 || len(report.PriceUpdates.TokenPriceUpdates) != 0 {
		lggr.Infow("Adding to inflight fee updates", "priceUpdates", report.PriceUpdates)
		c.inFlightPriceUpdates = append(c.inFlightPriceUpdates, InflightPriceUpdate{
			priceUpdates:  report.PriceUpdates,
			createdAt:     time.Now(),
			epochAndRound: epochAndRound,
		})
	}
	return nil
}
