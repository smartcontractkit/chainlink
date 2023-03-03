package persist

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type TransmissionDBUpdate struct {
	Timestamp           types.ReportTimestamp
	PendingTransmission *types.PendingTransmission
}

// Persists state from the transmission protocol to the database to allow for recovery
// after restarts
func PersistTransmission(
	ctx context.Context,
	chPersist <-chan TransmissionDBUpdate,
	db types.Database,
	dbTimeout time.Duration,
	logger loghelper.LoggerWithContext,
) {
	for {
		select {
		case update, ok := <-chPersist:
			if !ok {
				logger.Error("PersistTransmission: chPersist closed unexpectedly, exiting", nil)
				return
			}

			func() {
				dbCtx, dbCancel := context.WithTimeout(ctx, dbTimeout)
				defer dbCancel()

				store := update.PendingTransmission != nil
				var err error
				if store {
					err = db.StorePendingTransmission(dbCtx, update.Timestamp, *update.PendingTransmission)
				} else {
					err = db.DeletePendingTransmission(dbCtx, update.Timestamp)
				}
				if err != nil {
					logger.ErrorIfNotCanceled(
						"PersistTransmission: error updating database",
						dbCtx,
						commontypes.LogFields{"error": err, "store": store},
					)
					return
				}
			}()

		case <-ctx.Done():
			logger.Info("PersistTransmission: exiting", nil)
			return
		}
	}
}
