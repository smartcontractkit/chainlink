/*
The report coordinator provides 3 main functions:
Filter
Accept
IsTransmissionConfirmed

This has 2 purposes:
When an id is accepted using the Accept function, the upkeep id should be
included in the Filter function. This allows an upkeep id to be filtered out
of a list of upkeep keys.

When an upkeep key is accepted using the Accept function, the upkeep key will
return false on IsTransmissionConfirmed until a perform log is identified with
the same key. This allows a coordinated effort on transmit fallbacks.
*/
package keepers

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/smartcontractkit/ocr2keepers/pkg/chain"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"
	"github.com/smartcontractkit/ocr2keepers/pkg/util"
)

var (
	ErrKeyAlreadyAccepted = fmt.Errorf("key alredy accepted")
	IndefiniteBlockingKey = chain.BlockKey("18446744073709551616") // Higher than possible block numbers (uint64), used to block keys indefintely
)

type reportCoordinator struct {
	logger         *log.Logger
	registry       types.Registry
	logs           types.PerformLogProvider
	minConfs       int
	idBlocks       *util.Cache[idBlocker] // should clear out when the next perform with this id occurs
	activeKeys     *util.Cache[bool]
	cacheCleaner   *util.IntervalCacheCleaner[bool]
	idCacheCleaner *util.IntervalCacheCleaner[idBlocker]
	starter        sync.Once
	chStop         chan struct{}
}

func newReportCoordinator(r types.Registry, s time.Duration, cacheClean time.Duration, logs types.PerformLogProvider, minConfs int, logger *log.Logger) *reportCoordinator {
	c := &reportCoordinator{
		logger:         logger,
		registry:       r,
		logs:           logs,
		minConfs:       minConfs,
		idBlocks:       util.NewCache[idBlocker](s),
		activeKeys:     util.NewCache[bool](time.Hour), // 1 hour allows the cleanup routine to clear stale data
		idCacheCleaner: util.NewIntervalCacheCleaner[idBlocker](cacheClean),
		cacheCleaner:   util.NewIntervalCacheCleaner[bool](cacheClean),
		chStop:         make(chan struct{}),
	}

	runtime.SetFinalizer(c, func(srv *reportCoordinator) { srv.stop() })

	c.start()

	return c
}

// Filter returns a filter function that removes upkeep keys that apply to this
// filter. Returns false if a key should be filtered out.
func (rc *reportCoordinator) Filter() func(types.UpkeepKey) bool {
	return func(key types.UpkeepKey) bool {
		blockKey, id, err := key.BlockKeyAndUpkeepID()
		if err != nil {
			return false
		}

		// only apply filter if key id is registered in the cache
		if bl, ok := rc.idBlocks.Get(string(id)); ok {
			isAfter, err := blockKey.After(bl.TransmitBlockNumber)
			if err != nil {
				return false
			}

			// do not filter the key out if key block is after block in cache
			return isAfter
		}

		return true
	}
}

func (rc *reportCoordinator) Accept(key types.UpkeepKey) error {
	blockKey, id, err := key.BlockKeyAndUpkeepID()
	if err != nil {
		return err
	}

	// If a key is already active then don't update filters, but also not throw errors as
	// there might be other keys in the same report which can get accepted
	if _, ok := rc.activeKeys.Get(key.String()); !ok {
		// Set the key as accepted within activeKeys
		rc.activeKeys.Set(key.String(), false, util.DefaultCacheExpiration)

		// Set idBlocks with the key as checkBlockNumber and IndefiniteBlockingKey as TransmitBlockNumber
		rc.updateIdBlock(string(id), idBlocker{
			CheckBlockNumber:    blockKey,
			TransmitBlockNumber: IndefiniteBlockingKey,
		})
	}

	return nil
}

func (rc *reportCoordinator) IsTransmissionConfirmed(key types.UpkeepKey) bool {
	// key is confirmed if it both exists and has been confirmed by the log
	// poller
	confirmed, ok := rc.activeKeys.Get(key.String())
	return !ok || (ok && confirmed)
}

func (rc *reportCoordinator) checkLogs() {
	performLogs, _ := rc.logs.PerformLogs(context.Background())
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

		logCheckBlockKey, id, err := l.Key.BlockKeyAndUpkeepID()
		if err != nil {
			continue
		}

		if confirmed, ok := rc.activeKeys.Get(l.Key.String()); ok {
			if !confirmed {
				// Process log if the key hasn't been confirmed yet
				rc.logger.Printf("Perform log found for key %s in transaction %s at block %s, with confirmations %d", l.Key, l.TransactionHash, l.TransmitBlock, l.Confirmations)

				// set state of key to indicate that the report was transmitted
				rc.activeKeys.Set(l.Key.String(), true, util.DefaultCacheExpiration)

				rc.updateIdBlock(string(id), idBlocker{
					CheckBlockNumber:    logCheckBlockKey,
					TransmitBlockNumber: l.TransmitBlock, // Removes the id from filters from higher blocks
				})
			}

			if confirmed {
				// This can happen if we get a perform log for the same key again on a newer block in case of reorgs
				// In this case, no change to activeKeys is needed, but idBlocks is updated to the newer BlockNumber
				idBlock, ok := rc.idBlocks.Get(string(id))
				if ok && idBlock.CheckBlockNumber.String() == logCheckBlockKey.String() &&
					idBlock.TransmitBlockNumber.String() != l.TransmitBlock.String() {

					rc.logger.Printf("Got a re-orged perform log for key %s in transaction %s at block %s, with confirmations %d", l.Key, l.TransactionHash, l.TransmitBlock, l.Confirmations)

					rc.updateIdBlock(string(id), idBlocker{
						CheckBlockNumber:    logCheckBlockKey,
						TransmitBlockNumber: l.TransmitBlock,
					})
				}
			}
		}
	}

	staleReportLogs, _ := rc.logs.StaleReportLogs(context.Background())
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

		logCheckBlockKey, id, err := l.Key.BlockKeyAndUpkeepID()
		if err != nil {
			continue
		}
		nextKey, err := logCheckBlockKey.Next()
		if err != nil {
			continue
		}

		if confirmed, ok := rc.activeKeys.Get(l.Key.String()); ok {
			if !confirmed {
				// Process log if the key hasn't been confirmed yet
				rc.logger.Printf("Stale report log found for key %s in transaction %s at block %s, with confirmations %d", l.Key, l.TransactionHash, l.TransmitBlock, l.Confirmations)
				// set state of key to indicate that the report was transmitted
				rc.activeKeys.Set(l.Key.String(), true, util.DefaultCacheExpiration)

				rc.updateIdBlock(string(id), idBlocker{
					CheckBlockNumber:    logCheckBlockKey,
					TransmitBlockNumber: nextKey, // Removes the id from filters after logCheckBlockKey+1
					// We add one here as this filter is applied on RPC checkBlockNumber (which will be atleast logCheckBlockKey+1+1)
					// resulting in atleast report checkBlockNumber of logCheckBlockKey+1
				})
			}

			if confirmed {
				// This can happen if we get a stale log for the same key again on a newer block or in case
				// the key was unblocked due to a performLog which later got reorged into a stale log
				idBlock, ok := rc.idBlocks.Get(string(id))
				if ok && idBlock.CheckBlockNumber.String() == logCheckBlockKey.String() &&
					idBlock.TransmitBlockNumber.String() != nextKey.String() {

					rc.logger.Printf("Got a stale report log for previously accepted key %s in transaction %s at block %s, with confirmations %d", l.Key, l.TransactionHash, l.TransmitBlock, l.Confirmations)

					rc.updateIdBlock(string(id), idBlocker{
						CheckBlockNumber:    logCheckBlockKey,
						TransmitBlockNumber: nextKey,
					})
				}
			}
		}
	}
}

type idBlocker struct {
	CheckBlockNumber    types.BlockKey
	TransmitBlockNumber types.BlockKey
}

// idBlock should only be updated if checkBlockNumber is set higher
// or checkBlockNumber is the same and transmitBlockNumber is higher
// (with a special case for IndefiniteBlockingKey being considered lowest).
//
// For a sequence of updates, updateIdBlock can be called in any order
// on different nodes, but by maintaining this invariant it results in
// an eventually consistent value across nodes.
func (b idBlocker) shouldUpdate(val idBlocker) (bool, error) {
	isAfter, err := val.CheckBlockNumber.After(b.CheckBlockNumber)
	if err != nil {
		return false, err
	}
	if isAfter {
		// val has higher checkBlockNumber
		return true, nil
	}

	isAfter, err = b.CheckBlockNumber.After(val.CheckBlockNumber)
	if err != nil {
		return false, err
	}
	if isAfter {
		// b has higher checkBlockNumber
		return false, nil
	}

	// Now the checkBlockNumber should be same

	// If b has an IndefiniteBlockingKey, then update
	if b.TransmitBlockNumber.String() == IndefiniteBlockingKey.String() {
		return true, nil
	}
	// If val has an IndefiniteBlockingKey, then don't update
	if val.TransmitBlockNumber.String() == IndefiniteBlockingKey.String() {
		return false, nil
	}
	// return true if val.TransmitBlockNumber is higher
	return val.TransmitBlockNumber.After(b.TransmitBlockNumber)
}

func (rc *reportCoordinator) updateIdBlock(key string, val idBlocker) {
	idBlock, ok := rc.idBlocks.Get(key)
	if ok {
		shouldUpdate, err := idBlock.shouldUpdate(val)
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

func (rc *reportCoordinator) start() {
	rc.starter.Do(func() {
		go rc.run()
		go rc.idCacheCleaner.Run(rc.idBlocks)
		go rc.cacheCleaner.Run(rc.activeKeys)
	})
}

func (rc *reportCoordinator) stop() {
	rc.chStop <- struct{}{}
	rc.idCacheCleaner.Stop()
	rc.cacheCleaner.Stop()
}

func (rc *reportCoordinator) run() {
	// TODO: handle panics by restarting this process

	cadence := time.Second
	timer := time.NewTimer(cadence)

	for {
		select {
		case <-timer.C:
			startTime := time.Now()

			rc.checkLogs()

			// attempt to ahere to a cadence of at least every second
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
