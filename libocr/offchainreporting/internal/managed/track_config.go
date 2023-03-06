package managed

import (
	"context"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

type trackConfigState struct {
	ctx context.Context
	// in
	configTracker types.ContractConfigTracker
	localConfig   types.LocalConfig
	logger        loghelper.LoggerWithContext
	// out
	chChanges chan<- types.ContractConfig
	// local
	subprocesses subprocesses.Subprocesses
	configDigest types.ConfigDigest
}

func (state *trackConfigState) run() {
	// Check immediately after startup
	tCheckLatestConfigDetails := time.After(0)
	tResubscribe := time.After(0)

	var subscription types.ContractConfigSubscription
	var chSubscription <-chan types.ContractConfig

	for {
		select {
		case _, ok := <-chSubscription:
			if ok {
				// Check immediately for new config
				tCheckLatestConfigDetails = time.After(0 * time.Second)
				state.logger.Info("TrackConfig: subscription fired", nil)
			} else {
				chSubscription = nil
				subscription.Close()
				state.logger.Warn("TrackConfig: subscription was closed", nil)
				tResubscribe = time.After(0)
			}
		case <-tCheckLatestConfigDetails:
			change, awaitingConfirmation := state.checkLatestConfigDetails()
			state.logger.Debug("TrackConfig: checking latestConfigDetails", nil)

			// poll more rapidly if we're awaiting confirmation
			if awaitingConfirmation {
				wait := 15 * time.Second
				if state.localConfig.ContractConfigTrackerPollInterval < wait {
					wait = state.localConfig.ContractConfigTrackerPollInterval
				}
				tCheckLatestConfigDetails = time.After(wait)
				state.logger.Info("TrackConfig: awaiting confirmation of new config", commontypes.LogFields{
					"wait": wait,
				})
			} else {
				tCheckLatestConfigDetails = time.After(state.localConfig.ContractConfigTrackerPollInterval)
			}

			if change != nil {
				state.configDigest = change.ConfigDigest
				state.logger.Info("TrackConfig: returning config", commontypes.LogFields{
					"configDigest": change.ConfigDigest.Hex(),
				})
				select {
				case state.chChanges <- *change:
				case <-state.ctx.Done():
				}
			}
		case <-tResubscribe:
			subscribeCtx, subscribeCancel := context.WithTimeout(state.ctx, state.localConfig.BlockchainTimeout)
			var err error
			subscription, err = state.configTracker.SubscribeToNewConfigs(subscribeCtx)
			subscribeCancel()
			if err != nil {
				state.logger.ErrorIfNotCanceled(
					"TrackConfig: failed to SubscribeToNewConfigs. Retrying later",
					subscribeCtx,
					commontypes.LogFields{
						"error":                                  err,
						"ContractConfigTrackerSubscribeInterval": state.localConfig.ContractConfigTrackerSubscribeInterval,
					},
				)
				tResubscribe = time.After(state.localConfig.ContractConfigTrackerSubscribeInterval)
			} else {
				chSubscription = subscription.Configs()
			}
		case <-state.ctx.Done():
		}

		// ensure prompt exit
		select {
		case <-state.ctx.Done():
			state.logger.Debug("TrackConfig: winding down", nil)
			if subscription != nil {
				subscription.Close()
			}
			state.subprocesses.Wait()
			state.logger.Debug("TrackConfig: exiting", nil)
			return
		default:
		}
	}
}

func (state *trackConfigState) checkLatestConfigDetails() (
	latestConfigDetails *types.ContractConfig, awaitingConfirmation bool,
) {
	bhCtx, bhCancel := context.WithTimeout(state.ctx, state.localConfig.BlockchainTimeout)
	defer bhCancel()
	blockheight, err := state.configTracker.LatestBlockHeight(bhCtx)
	if err != nil {
		state.logger.ErrorIfNotCanceled("TrackConfig: error during LatestBlockHeight()", bhCtx, commontypes.LogFields{
			"error": err,
		})
		return nil, false
	}

	detailsCtx, detailsCancel := context.WithTimeout(state.ctx, state.localConfig.BlockchainTimeout)
	defer detailsCancel()
	changedInBlock, latestConfigDigest, err := state.configTracker.LatestConfigDetails(detailsCtx)
	if err != nil {
		state.logger.ErrorIfNotCanceled("TrackConfig: error during LatestConfigDetails()", detailsCtx, commontypes.LogFields{
			"error": err,
		})
		return nil, false
	}
	if latestConfigDigest == (types.ConfigDigest{}) {
		state.logger.Error("TrackConfig: LatestConfigDetails() returned a zero configDigest. Looks like the contract has not been configured", commontypes.LogFields{
			"configDigest": latestConfigDigest,
		})
		return nil, false
	}
	if state.configDigest == latestConfigDigest {
		return nil, false
	}
	if !state.localConfig.SkipContractConfigConfirmations && blockheight < changedInBlock+uint64(state.localConfig.ContractConfigConfirmations)-1 {
		return nil, true
	}
	configCtx, configCancel := context.WithTimeout(state.ctx, state.localConfig.BlockchainTimeout)
	defer configCancel()
	contractConfig, err := state.configTracker.ConfigFromLogs(configCtx, changedInBlock)
	if err != nil {
		state.logger.ErrorIfNotCanceled("TrackConfig: error during LatestConfigDetails()", configCtx, commontypes.LogFields{
			"error": err,
		})
		return nil, true
	}
	if contractConfig.EncodedConfigVersion != config.EncodedConfigVersion {
		state.logger.Error("TrackConfig: received config change with unknown EncodedConfigVersion",
			commontypes.LogFields{"versionReceived": contractConfig.EncodedConfigVersion})
		return nil, false
	}
	return &contractConfig, false
}

func TrackConfig(
	ctx context.Context,

	configTracker types.ContractConfigTracker,
	initialConfigDigest types.ConfigDigest,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,

	chChanges chan<- types.ContractConfig,
) {
	state := trackConfigState{
		ctx,
		// in
		configTracker,
		localConfig,
		logger,
		//out
		chChanges,
		// local
		subprocesses.Subprocesses{},
		initialConfigDigest,
	}
	state.run()
}
