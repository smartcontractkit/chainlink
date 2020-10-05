package managed

import (
	"context"
	"math/rand"
	"time"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

const collectInterval = 10 * time.Minute
const olderThan = 24 * time.Hour

func collectGarbage(
	ctx context.Context,
	database types.Database,
	localConfig types.LocalConfig,
	logger types.Logger,
) {
	for {
		wait := collectInterval + time.Duration(rand.Float64()*5.0*60.0)*time.Second
		logger.Info("collectGarbage: going to sleep", types.LogFields{
			"duration": wait,
		})
		select {
		case <-time.After(wait):
			logger.Info("collectGarbage: starting collection of old transmissions", types.LogFields{
				"olderThan": olderThan,
			})
			childCtx, childCancel := context.WithTimeout(ctx, localConfig.DatabaseTimeout)
			defer childCancel()
			err := database.DeletePendingTransmissionsOlderThan(childCtx, time.Now().Add(-olderThan))
			if err != nil {
				logger.Info("collectGarbage: error in DeletePendingTransmissionsOlderThan", types.LogFields{
					"error":     err,
					"olderThan": olderThan,
				})
			} else {
				logger.Info("collectGarbage: finished collection", nil)
			}
		case <-ctx.Done():
			logger.Info("collectGarbage: exiting", nil)
			return
		}
	}
}
