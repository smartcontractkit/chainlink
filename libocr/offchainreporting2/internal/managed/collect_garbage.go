package managed

import (
	"context"
	"math/rand"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

const collectInterval = 10 * time.Minute
const olderThan = 24 * time.Hour

// collectGarbage periodically collects garbage left by old transmission protocol instances
func collectGarbage(
	ctx context.Context,
	database types.Database,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
) {
	for {
		wait := collectInterval + time.Duration(rand.Float64()*5.0*60.0)*time.Second
		logger.Info("collectGarbage: going to sleep", commontypes.LogFields{
			"duration": wait,
		})
		select {
		case <-time.After(wait):
			logger.Info("collectGarbage: starting collection of old transmissions", commontypes.LogFields{
				"olderThan": olderThan,
			})
			// To make sure the context is not leaked we are wrapping the database query.
			func() {
				childCtx, childCancel := context.WithTimeout(ctx, localConfig.DatabaseTimeout)
				defer childCancel()
				err := database.DeletePendingTransmissionsOlderThan(childCtx, time.Now().Add(-olderThan))
				if err != nil {
					logger.ErrorIfNotCanceled(
						"collectGarbage: error in DeletePendingTransmissionsOlderThan",
						childCtx,
						commontypes.LogFields{
							"error":     err,
							"olderThan": olderThan,
						},
					)
				} else {
					logger.Info("collectGarbage: finished collection", nil)
				}
			}()
		case <-ctx.Done():
			logger.Info("collectGarbage: exiting", nil)
			return
		}
	}
}
