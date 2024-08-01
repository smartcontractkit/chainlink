package managed

import (
	"context"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/serialization/protobuf"
	"github.com/smartcontractkit/libocr/offchainreporting/internal/shim"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/subprocesses"
)

// RunManagedOracle runs a "managed" version of protocol.RunOracle. It handles
// configuration updates and translating from commontypes.BinaryNetworkEndpoint to
// protocol.NetworkEndpoint.
func RunManagedOracle(
	ctx context.Context,

	v2bootstrappers []commontypes.BootstrapperLocator,
	configOverrider types.ConfigOverrider,
	configTracker types.ContractConfigTracker,
	contractTransmitter types.ContractTransmitter,
	database types.Database,
	datasource types.DataSource,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	monitoringEndpoint commontypes.MonitoringEndpoint,
	netEndpointFactory types.BinaryNetworkEndpointFactory,
	privateKeys types.PrivateKeys,
) {
	mo := managedOracleState{
		ctx: ctx,

		v2bootstrappers:     v2bootstrappers,
		configOverrider:     configOverrider,
		configTracker:       configTracker,
		contractTransmitter: contractTransmitter,
		database:            database,
		datasource:          datasource,
		localConfig:         localConfig,
		logger:              logger,
		monitoringEndpoint:  monitoringEndpoint,
		netEndpointFactory:  netEndpointFactory,
		privateKeys:         privateKeys,
	}
	mo.run()
}

type managedOracleState struct {
	ctx context.Context

	v2bootstrappers     []commontypes.BootstrapperLocator
	config              config.SharedConfig
	configOverrider     types.ConfigOverrider
	configTracker       types.ContractConfigTracker
	contractTransmitter types.ContractTransmitter
	database            types.Database
	datasource          types.DataSource
	localConfig         types.LocalConfig
	logger              loghelper.LoggerWithContext
	monitoringEndpoint  commontypes.MonitoringEndpoint
	netEndpointFactory  types.BinaryNetworkEndpointFactory
	privateKeys         types.PrivateKeys

	chTelemetry        chan<- *protobuf.TelemetryWrapper
	netEndpoint        *shim.SerializingEndpoint
	oracleCancel       context.CancelFunc
	oracleSubprocesses subprocesses.Subprocesses
	otherSubprocesses  subprocesses.Subprocesses
}

func (mo *managedOracleState) run() {
	{
		chTelemetry := make(chan *protobuf.TelemetryWrapper, 100)
		mo.chTelemetry = chTelemetry
		mo.otherSubprocesses.Go(func() {
			forwardTelemetry(mo.ctx, mo.logger, mo.monitoringEndpoint, chTelemetry)
		})
	}

	mo.otherSubprocesses.Go(func() {
		collectGarbage(mo.ctx, mo.database, mo.localConfig, mo.logger)
	})

	// Restore config from database, so that we can run even if the ethereum node
	// isn't working.
	{
		var cc *types.ContractConfig
		ok := mo.otherSubprocesses.BlockForAtMost(
			mo.ctx,
			mo.localConfig.DatabaseTimeout,
			func(ctx context.Context) {
				cc = loadConfigFromDatabase(ctx, mo.database, mo.logger)
			},
		)
		if !ok {
			mo.logger.Error("ManagedOracle: database timed out while attempting to restore configuration", commontypes.LogFields{
				"timeout": mo.localConfig.DatabaseTimeout,
			})
		} else if cc != nil {
			mo.configChanged(*cc)
		}
	}

	// Only start tracking config after we attempted to load config from db
	chNewConfig := make(chan types.ContractConfig, 5)
	mo.otherSubprocesses.Go(func() {
		TrackConfig(mo.ctx, mo.configTracker, mo.config.ConfigDigest, mo.localConfig, mo.logger, chNewConfig)
	})

	for {
		select {
		case change := <-chNewConfig:
			mo.logger.Info("ManagedOracle: switching between configs", commontypes.LogFields{
				"oldConfigDigest": mo.config.ConfigDigest.Hex(),
				"newConfigDigest": change.ConfigDigest.Hex(),
			})
			mo.configChanged(change)
		case <-mo.ctx.Done():
			mo.logger.Info("ManagedOracle: winding down", nil)
			mo.closeOracle()
			mo.otherSubprocesses.Wait()
			mo.logger.Info("ManagedOracle: exiting", nil)
			return // Exit ManagedOracle event loop altogether
		}
	}
}

func (mo *managedOracleState) closeOracle() {
	if mo.oracleCancel != nil {
		mo.oracleCancel()
		mo.oracleSubprocesses.Wait()
		err := mo.netEndpoint.Close()
		if err != nil {
			mo.logger.Error("ManagedOracle: error while closing BinaryNetworkEndpoint", commontypes.LogFields{
				"error": err,
			})
			// nothing to be done about it, let's try to carry on.
		}
		mo.oracleCancel = nil
		mo.netEndpoint = nil
	}
}

func (mo *managedOracleState) configChanged(contractConfig types.ContractConfig) {
	// Cease any operation from earlier configs
	mo.closeOracle()

	// Decode contractConfig
	skipChainSpecificChecks := mo.localConfig.DevelopmentMode == types.EnableDangerousDevelopmentMode
	var err error
	var oid commontypes.OracleID
	mo.config, oid, err = config.SharedConfigFromContractConfig(
		mo.contractTransmitter.ChainID(),
		skipChainSpecificChecks,
		contractConfig,
		mo.privateKeys,
		mo.netEndpointFactory.PeerID(),
		mo.contractTransmitter.FromAddress(),
	)
	if err != nil {
		mo.logger.Error("ManagedOracle: error while updating config", commontypes.LogFields{
			"error": err,
		})
		return
	}

	// Run with new config
	peerIDs := []string{}
	for _, identity := range mo.config.OracleIdentities {
		peerIDs = append(peerIDs, identity.PeerID)
	}

	childLogger := mo.logger.MakeChild(commontypes.LogFields{
		"configDigest": fmt.Sprintf("%x", mo.config.ConfigDigest),
		"oid":          oid,
	})

	binNetEndpoint, err := mo.netEndpointFactory.NewEndpoint(mo.config.ConfigDigest, peerIDs,
		mo.v2bootstrappers, mo.config.F, computeTokenBucketRefillRate(mo.config.PublicConfig),
		computeTokenBucketSize())
	if err != nil {
		mo.logger.Error("ManagedOracle: error during NewEndpoint", commontypes.LogFields{
			"error":           err,
			"configDigest":    mo.config.ConfigDigest,
			"peerIDs":         peerIDs,
			"v2bootstrappers": mo.v2bootstrappers,
		})
		return
	}

	netEndpoint := shim.NewSerializingEndpoint(
		mo.chTelemetry,
		mo.config.ConfigDigest,
		binNetEndpoint,
		childLogger,
	)

	if err := netEndpoint.Start(); err != nil {
		mo.logger.Error("ManagedOracle: error during netEndpoint.Start()", commontypes.LogFields{
			"error":        err,
			"configDigest": mo.config.ConfigDigest,
		})
		return
	}

	mo.netEndpoint = netEndpoint
	oracleCtx, oracleCancel := context.WithCancel(mo.ctx)
	mo.oracleCancel = oracleCancel
	mo.oracleSubprocesses.Go(func() {
		defer oracleCancel()
		protocol.RunOracle(
			oracleCtx,
			mo.config,
			ConfigOverriderWrapper{mo.configOverrider},
			mo.contractTransmitter,
			mo.database,
			mo.datasource,
			oid,
			mo.privateKeys,
			mo.localConfig,
			childLogger,
			mo.netEndpoint,
			shim.MakeTelemetrySender(mo.chTelemetry, childLogger),
		)
	})

	childCtx, childCancel := context.WithTimeout(mo.ctx, mo.localConfig.DatabaseTimeout)
	defer childCancel()
	if err := mo.database.WriteConfig(childCtx, contractConfig); err != nil {
		mo.logger.ErrorIfNotCanceled("ManagedOracle: error writing new config to database", childCtx, commontypes.LogFields{
			"configDigest": mo.config.ConfigDigest,
			"config":       contractConfig,
			"error":        err,
		})
	}
}

func computeTokenBucketRefillRate(cfg config.PublicConfig) float64 {
	return (1.0*float64(time.Second)/float64(cfg.DeltaResend) +
		1.0*float64(time.Second)/float64(cfg.DeltaProgress) +
		1.0*float64(time.Second)/float64(cfg.DeltaRound) +
		3.0*float64(time.Second)/float64(cfg.DeltaRound) +
		2.0*float64(time.Second)/float64(cfg.DeltaRound)) * 2.0
}

func computeTokenBucketSize() int {
	return (2 + 6) * 2
}
