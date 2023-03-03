package persist

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type persistPacemakerState struct {
	ctx context.Context

	chPersist       <-chan types.PersistentState
	configDigest    types.ConfigDigest
	database        types.Database
	databaseTimeout time.Duration
	logger          loghelper.LoggerWithContext

	writtenState *types.PersistentState
}

// PersistPacemaker receives states from the pacemaker protocol it should
// persist to the db through chPersist and writes them to database.
func PersistPacemaker(
	ctx context.Context,
	chPersist <-chan types.PersistentState,
	configDigest types.ConfigDigest,
	database types.Database,
	databaseTimeout time.Duration,
	logger loghelper.LoggerWithContext,
) {
	ps := persistPacemakerState{
		ctx,

		chPersist,
		configDigest,
		database,
		databaseTimeout,
		logger,

		nil,
	}
	ps.run()
}

// run gets updates from the outside (through chPersist) in a loop, drains
// chPersist so that it can ignore all but the latest state, and writes the
// latest state to the database if it's new, i.e. differs from the previously
// written state.
func (ps *persistPacemakerState) run() {
	for {
		select {
		case state, ok := <-ps.chPersist:
			if !ok {
				ps.logger.Error("Persist: chPersist closed unexpectedly, can no longer persist state. This should *not* happen.", commontypes.LogFields{
					"lastWrittenState": ps.writtenState,
				})
				return
			}
		DrainChannel:
			for {
				select {
				case state, ok = <-ps.chPersist:
					if !ok {
						ps.logger.Error("Persist: chPersist closed unexpectedly, can no longer persist state. This should *not* happen.", commontypes.LogFields{
							"lastWrittenState": ps.writtenState,
						})
						return
					}
				default:
					break DrainChannel
				}
			}
			ps.writeIfNew(state)

		case <-ps.ctx.Done():
			ps.logger.Debug("Persist: exiting", nil)
			return
		}
	}
}

// writeIfNew writes pendingState to the database, iff pendingState differs from
// the last written state.
func (ps *persistPacemakerState) writeIfNew(pendingState types.PersistentState) {
	if ps.writtenState != nil && pendingState.Equal(*ps.writtenState) {
		return
	}

	writeCtx, writeCancel := context.WithTimeout(ps.ctx, ps.databaseTimeout)
	defer writeCancel()
	err := ps.database.WriteState(
		writeCtx,
		ps.configDigest,
		pendingState,
	)
	if err != nil {
		ps.logger.ErrorIfNotCanceled("Persist: unexpected error while persisting state to database", writeCtx, commontypes.LogFields{
			"error": err,
		})
		return
	}

	ps.writtenState = &pendingState
}
