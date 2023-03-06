package managed

import (
	"context"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// RunManagedBootstrapNode runs a "managed" bootstrap node. It handles
// configuration updates on the contract.
func RunManagedBootstrapNode(
	ctx context.Context,

	bootstrapperFactory types.BootstrapperFactory,
	v1bootstrappers []string,
	v2bootstrappers []commontypes.BootstrapperLocator,
	contractConfigTracker types.ContractConfigTracker,
	database types.Database,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
) {
	mb := managedBootstrapNodeState{
		ctx: ctx,

		bootstrapperFactory: bootstrapperFactory,
		v1bootstrappers:     v1bootstrappers,
		v2bootstrappers:     v2bootstrappers,
		configTracker:       contractConfigTracker,
		database:            database,
		localConfig:         localConfig,
		logger:              logger,
	}
	mb.run()
}

type managedBootstrapNodeState struct {
	ctx context.Context

	v1bootstrappers     []string
	v2bootstrappers     []commontypes.BootstrapperLocator
	bootstrapperFactory types.BootstrapperFactory
	configTracker       types.ContractConfigTracker
	database            types.Database
	localConfig         types.LocalConfig
	logger              loghelper.LoggerWithContext

	bootstrapper commontypes.Bootstrapper
	config       config.PublicConfig
}

func (mb *managedBootstrapNodeState) run() {
	var subprocesses subprocesses.Subprocesses

	// Restore config from database, so that we can run even if the ethereum node
	// isn't working.
	{
		var cc *types.ContractConfig
		ok := subprocesses.BlockForAtMost(
			mb.ctx,
			mb.localConfig.DatabaseTimeout,
			func(ctx context.Context) {
				cc = loadConfigFromDatabase(ctx, mb.database, mb.logger)
			},
		)
		if !ok {
			mb.logger.Error("ManagedBootstrapNode: database timed out while attempting to restore configuration", commontypes.LogFields{
				"timeout": mb.localConfig.DatabaseTimeout,
			})
		} else if cc != nil {
			mb.configChanged(*cc)
		}
	}

	chNewConfig := make(chan types.ContractConfig, 5)
	subprocesses.Go(func() {
		TrackConfig(mb.ctx, mb.configTracker, mb.config.ConfigDigest, mb.localConfig, mb.logger, chNewConfig)
	})

	for {
		select {
		case cc := <-chNewConfig:
			mb.logger.Info("ManagedBootstrapNode: Switching between configs", commontypes.LogFields{
				"oldConfigDigest": mb.config.ConfigDigest.Hex(),
				"newConfigDigest": cc.ConfigDigest.Hex(),
			})
			mb.configChanged(cc)
		case <-mb.ctx.Done():
			mb.logger.Debug("ManagedBootstrapNode: winding down ", nil)
			mb.closeBootstrapper()
			subprocesses.Wait()
			mb.logger.Debug("ManagedBootstrapNode: exiting", nil)
			return
		}
	}
}

func (mb *managedBootstrapNodeState) closeBootstrapper() {
	if mb.bootstrapper != nil {
		err := mb.bootstrapper.Close()
		// there's not much we can do apart from logging the error and praying
		if err != nil {
			mb.logger.Error("ManagedBootstrapNode: error while closing bootstrapper", commontypes.LogFields{
				"error": err,
			})
		}
		mb.bootstrapper = nil
	}
}

func (mb *managedBootstrapNodeState) configChanged(cc types.ContractConfig) {
	// Cease any operation from earlier configs
	mb.closeBootstrapper()

	var err error
	// We're okay to skip chain-specific checks here. A bootstrap node does not
	// use any chain-specific parameters, since it doesn't participate in the
	// actual OCR protocol. It just hangs out on the P2P network and helps other
	// nodes find each other.
	mb.config, err = config.PublicConfigFromContractConfig(nil, true, cc)
	if err != nil {
		mb.logger.Error("ManagedBootstrapNode: error while decoding ContractConfig", commontypes.LogFields{
			"error": err,
		})
		return
	}

	peerIDs := []string{}
	for _, pcKey := range mb.config.OracleIdentities {
		peerIDs = append(peerIDs, pcKey.PeerID)
	}

	bootstrapper, err := mb.bootstrapperFactory.NewBootstrapper(mb.config.ConfigDigest, peerIDs, mb.v1bootstrappers, mb.v2bootstrappers, mb.config.F)
	if err != nil {
		mb.logger.Error("ManagedBootstrapNode: error during NewBootstrapper", commontypes.LogFields{
			"error":           err,
			"configDigest":    mb.config.ConfigDigest,
			"peerIDs":         peerIDs,
			"v1bootstrappers": mb.v1bootstrappers,
		})
		return
	}
	err = bootstrapper.Start()
	if err != nil {
		mb.logger.Error("ManagedBootstrapNode: error during bootstrapper.Start()", commontypes.LogFields{
			"error":        err,
			"configDigest": mb.config.ConfigDigest,
		})
		return
	}

	mb.bootstrapper = bootstrapper

	childCtx, childCancel := context.WithTimeout(mb.ctx, mb.localConfig.DatabaseTimeout)
	defer childCancel()
	if err := mb.database.WriteConfig(childCtx, cc); err != nil {
		mb.logger.ErrorIfNotCanceled("ManagedBootstrapNode: error writing new config to database", childCtx, commontypes.LogFields{
			"config": cc,
			"error":  err,
		})
		// We can keep running even without storing the config in the database
	}
}
