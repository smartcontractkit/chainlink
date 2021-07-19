package chainlink

import (
	"context"
	// "fmt"

	"github.com/smartcontractkit/chainlink/core/logger"
	"go.uber.org/zap/zapcore"
)

// SetServiceLogger sets the Logger for a given service and stores the setting in the db
func (app *ChainlinkApplication) SetServiceLogger(ctx context.Context, serviceName logger.ServiceName, level zapcore.Level) error {
	// newL, err := app.logger.InitServiceLevelLogger(serviceName, level.String())
	// if err != nil {
	// 	return err
	// }

	// switch serviceName {
	// case logger.HeadTracker:
	// 	app.HeadTracker.SetLogger(newL)
	// case logger.FluxMonitor:
	// 	app.FluxMonitor.SetLogger(newL)
	// case logger.JobSubscriber:
	// 	app.JobSubscriber.SetLogger(newL)
	// case logger.RunQueue:
	// 	app.RunQueue.SetLogger(newL)
	// case logger.BalanceMonitor:
	// 	app.balanceMonitor.SetLogger(newL)
	// case logger.TxManager:
	// 	app.TxManager.SetLogger(newL)
	// case logger.HeadBroadcaster:
	// 	app.HeadBroadcaster.SetLogger(newL)
	// case logger.EventBroadcaster:
	// 	app.EventBroadcaster.SetLogger(newL)
	// case logger.DatabaseBackup:
	// 	if app.databaseBackup != nil {
	// 		app.databaseBackup.SetLogger(newL)
	// 	}
	// case logger.PromReporter:
	// 	if app.promReporter != nil {
	// 		app.promReporter.SetLogger(newL)
	// 	}
	// case logger.SingletonPeerWrapper:
	// 	if app.peerWrapper != nil {
	// 		app.peerWrapper.SetLogger(newL)
	// 	}
	// case logger.OCRContractTracker:
	// case logger.ExplorerClient:
	// 	app.explorerClient.SetLogger(newL)
	// case logger.StatsPusher:
	// 	app.StatsPusher.SetLogger(newL)

	// default:
	// 	return fmt.Errorf("no service found with name: %s", serviceName)
	// }

	// return app.logger.Orm.SetServiceLogLevel(ctx, serviceName, level)
	return nil
}
