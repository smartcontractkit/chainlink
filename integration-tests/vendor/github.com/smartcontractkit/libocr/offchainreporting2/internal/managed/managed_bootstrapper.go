package managed

import (
	"context"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// RunManagedBootstrapper runs a "managed" bootstrapper. It handles
// configuration updates on the contract.
func RunManagedBootstrapper(
	ctx context.Context,

	bootstrapperFactory types.BootstrapperFactory,
	v2bootstrappers []commontypes.BootstrapperLocator,
	contractConfigTracker types.ContractConfigTracker,
	database types.ConfigDatabase,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	offchainConfigDigester types.OffchainConfigDigester,
) {
	runWithContractConfig(
		ctx,

		contractConfigTracker,
		database,
		func(ctx context.Context, contractConfig types.ContractConfig, logger loghelper.LoggerWithContext) {
			config, err := config.PublicConfigFromContractConfig(true, contractConfig)
			if err != nil {
				logger.Error("ManagedBootstrapper: error while decoding ContractConfig", commontypes.LogFields{
					"error": err,
				})
				return
			}

			peerIDs := []string{}
			for _, pcKey := range config.OracleIdentities {
				peerIDs = append(peerIDs, pcKey.PeerID)
			}

			bootstrapper, err := bootstrapperFactory.NewBootstrapper(config.ConfigDigest, peerIDs, v2bootstrappers, config.F)
			if err != nil {
				logger.Error("ManagedBootstrapper: error during NewBootstrapper", commontypes.LogFields{
					"error":           err,
					"peerIDs":         peerIDs,
					"v2bootstrappers": v2bootstrappers,
				})
				return
			}

			if err := bootstrapper.Start(); err != nil {
				logger.Error("ManagedBootstrapper: error during bootstrapper.Start()", commontypes.LogFields{
					"error": err,
				})
				return
			}
			defer loghelper.CloseLogError(
				bootstrapper,
				logger,
				"ManagedBootstrapper: error during bootstrapper.Close()",
			)

			<-ctx.Done()
		},
		localConfig,
		logger,
		offchainConfigDigester,
	)
}
