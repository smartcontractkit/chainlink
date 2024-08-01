/*
A coordinator provides the ability to filter upkeeps based on some type of
in-flight status.

The report coordinator provides 3 main functions:
IsPending
Accept
IsTransmissionConfirmed

This has 2 purposes:
When an id is accepted using the Accept function, the upkeep id should be
indicated as pending in the IsPending function. This allows an upkeep id to be
filtered out of a list of upkeep keys.

When an upkeep key is accepted using the Accept function, the upkeep key will
return false on IsTransmissionConfirmed until a perform log is identified with
the same key. This allows a coordinated effort on transmit fallbacks.

The report coordinator relies on two log types:
PerformLog - this log type indicates that an upkeep was completed
StaleReportLog - this log type indicates that an upkeep failed and can be
attempted again at a later block height
*/
package coordinator

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink-automation/pkg/util"
	ocr2keepers "github.com/smartcontractkit/chainlink-automation/pkg/v2"
)

var (
	DefaultLockoutWindow  = time.Duration(20) * time.Minute
	DefaultCacheClean     = time.Duration(30) * time.Second
	ErrKeyAlreadyAccepted = fmt.Errorf("key already accepted")
	IndefiniteBlockingKey = ocr2keepers.BlockKey("18446744073709551616") // Higher than possible block numbers (uint64), used to block keys indefintely
)

//go:generate mockery --name LogProvider --srcpkg "github.com/smartcontractkit/chainlink-automation/pkg/v2/coordinator" --case underscore --filename log_provider.generated.go
type LogProvider interface {
	PerformLogs(context.Context) ([]ocr2keepers.PerformLog, error)
	StaleReportLogs(context.Context) ([]ocr2keepers.StaleReportLog, error)
}

//go:generate mockery --name Encoder --srcpkg "github.com/smartcontractkit/chainlink-automation/pkg/v2/coordinator" --case underscore --filename encoder.generated.go
type Encoder interface {
	// SplitUpkeepKey ...
	SplitUpkeepKey(ocr2keepers.UpkeepKey) (ocr2keepers.BlockKey, ocr2keepers.UpkeepIdentifier, error)
	// After a is after b
	After(ocr2keepers.BlockKey, ocr2keepers.BlockKey) (bool, error)
	// Increment
	Increment(ocr2keepers.BlockKey) (ocr2keepers.BlockKey, error)
}

type reportCoordinator struct {
	// injected dependencies
	logger  *log.Logger
	logs    LogProvider
	encoder Encoder

	// initialised by the constructor
	idBlocks       *util.Cache[idBlocker] // should clear out when the next perform with this id occurs
	activeKeys     *util.Cache[bool]
	cacheCleaner   *util.IntervalCacheCleaner[bool]
	idCacheCleaner *util.IntervalCacheCleaner[idBlocker]

	// configurations
	minConfs int

	// run state data
	running atomic.Bool
	chStop  chan struct{}
}

// NewReportCoordinator provides a new report coordinator. The coordinator
// should be started before use.
func NewReportCoordinator(
	lockoutWindow time.Duration,
	cacheClean time.Duration,
	logs LogProvider,
	minConfs int,
	logger *log.Logger,
	encoder Encoder,
) *reportCoordinator {
	if lockoutWindow < 1 {
		lockoutWindow = DefaultLockoutWindow
	}

	if cacheClean < 1 {
		cacheClean = DefaultCacheClean
	}

	c := &reportCoordinator{
		logger:         logger,
		logs:           logs,
		minConfs:       minConfs,
		idBlocks:       util.NewCache[idBlocker](lockoutWindow),
		activeKeys:     util.NewCache[bool](time.Hour), // 1 hour allows the cleanup routine to clear stale data
		idCacheCleaner: util.NewIntervalCacheCleaner[idBlocker](cacheClean),
		cacheCleaner:   util.NewIntervalCacheCleaner[bool](cacheClean),
		chStop:         make(chan struct{}, 1),
		encoder:        encoder,
	}

	return c
}

// IsPending returns true if a key should be filtered out.
func (rc *reportCoordinator) IsPending(key ocr2keepers.UpkeepKey) (bool, error) {
	blockKey, id, err := rc.encoder.SplitUpkeepKey(key)
	if err != nil {
		return true, fmt.Errorf("%w: key parse error", err)
	}

	// only apply filter if key id is registered in the cache
	if bl, ok := rc.idBlocks.Get(string(id)); ok {
		isAfter, err := rc.encoder.After(blockKey, bl.TransmitBlockNumber)
		if err != nil {
			return true, fmt.Errorf("%w: not after transmit number", err)
		}

		// do not filter the key out if key block is after block in cache
		return !isAfter, nil
	}

	return false, nil
}

// Accept sets the pending status for a key
func (rc *reportCoordinator) Accept(key ocr2keepers.UpkeepKey) error {
	blockKey, id, err := rc.encoder.SplitUpkeepKey(key)
	if err != nil {
		return err
	}

	// If a key is already active then don't update filters, but also not throw errors as
	// there might be other keys in the same report which can get accepted
	if _, ok := rc.activeKeys.Get(string(key)); !ok {
		// Set the key as accepted within activeKeys
		rc.activeKeys.Set(string(key), false, util.DefaultCacheExpiration)

		// Set idBlocks with the key as checkBlockNumber and IndefiniteBlockingKey as TransmitBlockNumber
		rc.updateIdBlock(string(id), idBlocker{
			CheckBlockNumber:    blockKey,
			TransmitBlockNumber: IndefiniteBlockingKey,
		})
	}

	return nil
}

// IsTransmissionConfirmed returns whether the upkeep was successfully
// completed or not
func (rc *reportCoordinator) IsTransmissionConfirmed(key ocr2keepers.UpkeepKey) bool {
	// key is confirmed if it both exists and has been confirmed by the log
	// poller
	confirmed, ok := rc.activeKeys.Get(string(key))
	return !ok || (ok && confirmed)
}

func (rc *reportCoordinator) checkLogs(ctx context.Context) error {

	var (
		performLogs     []ocr2keepers.PerformLog
		staleReportLogs []ocr2keepers.StaleReportLog
		err             error
	)

	performLogs, err = rc.logs.PerformLogs(ctx)
	if err != nil {
		return err
	}
	// Perform log entries indicate that a perform exists on chain in some
	// capacity. the existance of an entry means that the transaction
	// was broadcast by at least one node. reorgs can still happen
	// causing performs to get moved to a later block or change to reorg logs.
	// Higher minConfirmations setting reduces the chances of this happening.
	//
	// We do two things upon receiving a perform log
	// - Mark the upkeep key responsible for the perform as 'transmitted', so that this node does not
	//   waste gas trying to transmit the same report again
	// - Unblock the upkeep from idBlocks so that it can be observed and reported on again.
	for _, l := range performLogs {
		if l.Confirmations < int64(rc.minConfs) {
			rc.logger.Printf("Skipping perform log in transaction %s as confirmations (%d) is less than min confirmations (%d)", l.TransactionHash, l.Confirmations, rc.minConfs)
			continue
		}

		logCheckBlockKey, id, err := rc.encoder.SplitUpkeepKey(l.Key)
		if err != nil {
			continue
		}

		if confirmed, ok := rc.activeKeys.Get(string(l.Key)); ok {
			if !confirmed {
				// Process log if the key hasn't been confirmed yet
				rc.logger.Printf("Perform log found for key %s in transaction %s at block %s, with confirmations %d", l.Key, l.TransactionHash, l.TransmitBlock, l.Confirmations)

				// set state of key to indicate that the report was transmitted
				rc.activeKeys.Set(string(l.Key), true, util.DefaultCacheExpiration)

				rc.updateIdBlock(string(id), idBlocker{
					CheckBlockNumber:    logCheckBlockKey,
					TransmitBlockNumber: l.TransmitBlock, // Removes the id from filters from higher blocks
				})
			}

			if confirmed {
				// This can happen if we get a perform log for the same key again on a newer block in case of reorgs
				// In this case, no change to activeKeys is needed, but idBlocks is updated to the newer BlockNumber
				idBlock, ok := rc.idBlocks.Get(string(id))
				if ok && string(idBlock.CheckBlockNumber) == string(logCheckBlockKey) &&
					string(idBlock.TransmitBlockNumber) != string(l.TransmitBlock) {

					rc.logger.Printf("Got a re-orged perform log for key %s in transaction %s at block %s, with confirmations %d", l.Key, l.TransactionHash, l.TransmitBlock, l.Confirmations)

					rc.updateIdBlock(string(id), idBlocker{
						CheckBlockNumber:    logCheckBlockKey,
						TransmitBlockNumber: l.TransmitBlock,
					})
				}
			}
		}
	}

	staleReportLogs, err = rc.logs.StaleReportLogs(context.Background())
	// It can happen that in between the time the report is generated and it gets
	// confirmed on chain something changes and it becomes stale. Current scenarios are:
	//    - Another report for the upkeep is transmitted making this report stale
	//    - Reorg happens which changes the checkBlockHash making this reorged report as it was checked on a different chain
	//    - There's a massive gas spike and upkeep does not have sufficient funds when report gets on chain
	// In such cases the upkeep is not performed and the contract emits a log indicating the staleness reason
	// instead of UpkeepPerformed log. We don't have different behaviours for different staleness
	// reasons and just want to unlock the upkeep when we receive such log.
	//
	// For these logs we do not have the exact key which generated this log. Hence we
	// are not able to mark the key responsible as transmitted which will result in some wasted
	// gas if this node tries to transmit it again, however we prioritise the upkeep performance
	// and clear the idBlocks for this upkeep.
	for _, l := range staleReportLogs {
		if l.Confirmations < int64(rc.minConfs) {
			rc.logger.Printf("Skipping stale report log in transaction %s as confirmations (%d) is less than min confirmations (%d)", l.TransactionHash, l.Confirmations, rc.minConfs)
			continue
		}

		logCheckBlockKey, id, err := rc.encoder.SplitUpkeepKey(l.Key)
		if err != nil {
			continue
		}
		nextKey, err := rc.encoder.Increment(logCheckBlockKey)
		if err != nil {
			continue
		}

		if confirmed, ok := rc.activeKeys.Get(string(l.Key)); ok {
			if !confirmed {
				// Process log if the key hasn't been confirmed yet
				rc.logger.Printf("Stale report log found for key %s in transaction %s at block %s, with confirmations %d", l.Key, l.TransactionHash, l.TransmitBlock, l.Confirmations)
				// set state of key to indicate that the report was transmitted
				rc.activeKeys.Set(string(l.Key), true, util.DefaultCacheExpiration)

				rc.updateIdBlock(string(id), idBlocker{
					CheckBlockNumber:    logCheckBlockKey,
					TransmitBlockNumber: nextKey, // Removes the id from filters after logCheckBlockKey+1
					// We add one here as this filter is applied on RPC checkBlockNumber (which will be at least logCheckBlockKey+1+1)
					// resulting in at least report checkBlockNumber of logCheckBlockKey+1
				})
			}

			if confirmed {
				// This can happen if we get a stale log for the same key again on a newer block or in case
				// the key was unblocked due to a performLog which later got reorged into a stale log
				idBlock, ok := rc.idBlocks.Get(string(id))
				if ok && string(idBlock.CheckBlockNumber) == string(logCheckBlockKey) &&
					string(idBlock.TransmitBlockNumber) != string(nextKey) {

					rc.logger.Printf("Got a stale report log for previously accepted key %s in transaction %s at block %s, with confirmations %d", l.Key, l.TransactionHash, l.TransmitBlock, l.Confirmations)

					rc.updateIdBlock(string(id), idBlocker{
						CheckBlockNumber:    logCheckBlockKey,
						TransmitBlockNumber: nextKey,
					})
				}
			}
		}
	}

	return err
}

type idBlocker struct {
	CheckBlockNumber    ocr2keepers.BlockKey
	TransmitBlockNumber ocr2keepers.BlockKey
}

// idBlock should only be updated if checkBlockNumber is set higher
// or checkBlockNumber is the same and transmitBlockNumber is higher
// (with a special case for IndefiniteBlockingKey being considered lowest).
//
// For a sequence of updates, updateIdBlock can be called in any order
// on different nodes, but by maintaining this invariant it results in
// an eventually consistent value across nodes.
func (b idBlocker) shouldUpdate(val idBlocker, e Encoder) (bool, error) {
	isAfter, err := e.After(val.CheckBlockNumber, b.CheckBlockNumber)
	if err != nil {
		return false, err
	}

	if isAfter {
		// val has higher checkBlockNumber
		return true, nil
	}

	isAfter, err = e.After(b.CheckBlockNumber, val.CheckBlockNumber)
	if err != nil {
		return false, err
	}

	if isAfter {
		// b has higher checkBlockNumber
		return false, nil
	}

	// Now the checkBlockNumber should be same

	// If b has an IndefiniteBlockingKey, then update
	if string(b.TransmitBlockNumber) == string(IndefiniteBlockingKey) {
		return true, nil
	}

	// If val has an IndefiniteBlockingKey, then don't update
	if string(val.TransmitBlockNumber) == string(IndefiniteBlockingKey) {
		return false, nil
	}

	// return true if val.TransmitBlockNumber is higher
	return e.After(val.TransmitBlockNumber, b.TransmitBlockNumber)
}

func (rc *reportCoordinator) updateIdBlock(key string, val idBlocker) {
	idBlock, ok := rc.idBlocks.Get(key)
	if ok {
		shouldUpdate, err := idBlock.shouldUpdate(val, rc.encoder)
		if err != nil {
			// Don't update on errors
			return
		}
		if !shouldUpdate {
			rc.logger.Printf("updateIdBlock for key %s: Not updating idBlocks (%+v) to new val (%+v)", key, idBlock, val)
			return
		}
	}

	rc.logger.Printf("updateIdBlock for key %s: value updated to %+v", key, val)
	rc.idBlocks.Set(key, val, util.DefaultCacheExpiration)
}

// Start starts all subprocesses
func (rc *reportCoordinator) Start() {
	if !rc.running.Load() {
		go rc.run()
		go rc.idCacheCleaner.Run(rc.idBlocks)
		go rc.cacheCleaner.Run(rc.activeKeys)

		rc.running.Swap(true)
	}
}

// Close terminates all subprocesses
func (rc *reportCoordinator) Close() error {
	if rc.running.Load() {
		rc.chStop <- struct{}{}
		rc.idCacheCleaner.Stop()
		rc.cacheCleaner.Stop()
		rc.running.Swap(false)
	}

	return nil
}

func (rc *reportCoordinator) run() {
	cadence := time.Second
	timer := time.NewTimer(cadence)

	for {
		select {
		case <-timer.C:
			startTime := time.Now()

			if err := rc.checkLogs(context.Background()); err != nil {
				rc.logger.Printf("failed to check perform and stale report logs: %s", err)
			}

			// attempt to adhere to a cadence of at least every second
			// a slow DB will cause the cadence to increase. these cases are logged
			diff := time.Since(startTime)
			if diff > cadence {
				rc.logger.Printf("log poll took %dms to complete; expected cadence is %dms; check database indexes and other performance improvements", diff/time.Millisecond, cadence/time.Millisecond)
				// start again immediately
				timer.Reset(time.Microsecond)
			} else {
				// wait the difference between the cadence and the time taken
				timer.Reset(cadence - diff)
			}
		case <-rc.chStop:
			timer.Stop()
			return
		}
	}
}
