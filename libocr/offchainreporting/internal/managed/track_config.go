package managed

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"github.com/smartcontractkit/chainlink/libocr/subprocesses"
)

type trackConfigState struct {
	ctx           context.Context
	configTracker types.ContractConfigTracker
	localConfig   types.LocalConfig
	logger        types.Logger
	chChanges     chan<- types.ContractConfig
	subprocesses  subprocesses.Subprocesses
	configDigest  types.ConfigDigest
}

func (state *trackConfigState) run() {
	tCheckLatestConfigDetails := time.After(0)
	tResubscribe := time.After(0)

	var subscription types.ContractConfigSubscription
	var chSubscription <-chan types.ContractConfig

	for {
		select {
		case _, ok := <-chSubscription:
			if ok {
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

			if awaitingConfirmation {
				wait := 15 * time.Second
				if state.localConfig.ContractConfigTrackerPollInterval < wait {
					wait = state.localConfig.ContractConfigTrackerPollInterval
				}
				tCheckLatestConfigDetails = time.After(wait)
				state.logger.Info("TrackConfig: awaiting confirmation of new config", types.LogFields{
					"wait": wait,
				})
			} else {
				tCheckLatestConfigDetails = time.After(state.localConfig.ContractConfigTrackerPollInterval)
			}

			if change != nil {
				state.configDigest = change.ConfigDigest
				state.logger.Info("TrackConfig: returning config", types.LogFields{
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
				state.logger.Error("TrackConfig: failed to SubscribeToNewConfigs. Retrying later", types.LogFields{
					"error":                                  err,
					"ContractConfigTrackerSubscribeInterval": state.localConfig.ContractConfigTrackerSubscribeInterval,
				})
				tResubscribe = time.After(state.localConfig.ContractConfigTrackerSubscribeInterval)
			} else {
				chSubscription = subscription.Configs()
			}
		case <-state.ctx.Done():
		}

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
		state.logger.Error("TrackConfig: error during LatestBlockHeight()", types.LogFields{
			"error": err,
		})
		return nil, false
	}

	detailsCtx, detailsCancel := context.WithTimeout(state.ctx, state.localConfig.BlockchainTimeout)
	defer detailsCancel()
	changedInBlock, latestConfigDigest, err := state.configTracker.LatestConfigDetails(detailsCtx)
	if err != nil {
		state.logger.Error("TrackConfig: error during LatestConfigDetails()", types.LogFields{
			"error": err,
		})
		return nil, false
	}
	if state.configDigest == latestConfigDigest {
		return nil, false
	}
	if blockheight < changedInBlock+uint64(state.localConfig.ContractConfigConfirmations) {
		return nil, true
	}
	configCtx, configCancel := context.WithTimeout(state.ctx, state.localConfig.BlockchainTimeout)
	defer configCancel()
	contractConfig, err := state.configTracker.ConfigFromLogs(configCtx, changedInBlock)
	if err != nil {
		state.logger.Error("TrackConfig: error during LatestConfigDetails()", types.LogFields{
			"error": err,
		})
		return nil, true
	}
	if contractConfig.EncodedConfigVersion != config.EncodedConfigVersion {
		state.logger.Error("TrackConfig: received config change with unknown EncodedConfigVersion",
			types.LogFields{"versionReceived": contractConfig.EncodedConfigVersion})
		return nil, false
	}
	return &contractConfig, false
}

func TrackConfig(
	ctx context.Context,

	configTracker types.ContractConfigTracker,
	localConfig types.LocalConfig,
	logger types.Logger,

	chChanges chan<- types.ContractConfig,
) {
	state := trackConfigState{
		ctx,
		configTracker,
		localConfig,
		logger,
		chChanges,
		subprocesses.Subprocesses{},
		types.ConfigDigest{},
	}
	state.run()
}
