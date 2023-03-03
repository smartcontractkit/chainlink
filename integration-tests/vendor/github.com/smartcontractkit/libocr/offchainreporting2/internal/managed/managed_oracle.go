package managed

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2/internal/shim"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/subprocesses"
	"go.uber.org/multierr"
)

// RunManagedOracle runs a "managed" version of protocol.RunOracle. It handles
// setting up telemetry, garbage collection, configuration updates, translating
// from commontypes.BinaryNetworkEndpoint to protocol.NetworkEndpoint, and
// creation/teardown of reporting plugins.
func RunManagedOracle(
	ctx context.Context,

	v2bootstrappers []commontypes.BootstrapperLocator,
	configTracker types.ContractConfigTracker,
	contractTransmitter types.ContractTransmitter,
	database types.Database,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	monitoringEndpoint commontypes.MonitoringEndpoint,
	netEndpointFactory types.BinaryNetworkEndpointFactory,
	offchainConfigDigester types.OffchainConfigDigester,
	offchainKeyring types.OffchainKeyring,
	onchainKeyring types.OnchainKeyring,
	reportingPluginFactory types.ReportingPluginFactory,
) {
	subs := subprocesses.Subprocesses{}
	defer subs.Wait()

	var chTelemetrySend chan<- *serialization.TelemetryWrapper
	{
		chTelemetry := make(chan *serialization.TelemetryWrapper, 100)
		chTelemetrySend = chTelemetry
		subs.Go(func() {
			forwardTelemetry(ctx, logger, monitoringEndpoint, chTelemetry)
		})
	}

	subs.Go(func() {
		collectGarbage(ctx, database, localConfig, logger)
	})

	runWithContractConfig(
		ctx,

		configTracker,
		database,
		func(ctx context.Context, contractConfig types.ContractConfig, logger loghelper.LoggerWithContext) {
			skipResourceExhaustionChecks := localConfig.DevelopmentMode == types.EnableDangerousDevelopmentMode
			sharedConfig, oid, err := config.SharedConfigFromContractConfig(
				skipResourceExhaustionChecks,
				contractConfig,
				offchainKeyring,
				onchainKeyring,
				netEndpointFactory.PeerID(),
				contractTransmitter.FromAccount(),
			)
			if err != nil {
				logger.Error("ManagedOracle: error while updating config", commontypes.LogFields{
					"error": err,
				})
				return
			}

			// Run with new config
			peerIDs := []string{}
			for _, identity := range sharedConfig.OracleIdentities {
				peerIDs = append(peerIDs, identity.PeerID)
			}

			childLogger := logger.MakeChild(commontypes.LogFields{
				"oid": oid,
			})

			reportingPlugin, reportingPluginInfo, err := reportingPluginFactory.NewReportingPlugin(types.ReportingPluginConfig{
				sharedConfig.ConfigDigest,
				oid,
				sharedConfig.N(),
				sharedConfig.F,
				sharedConfig.OnchainConfig,
				sharedConfig.ReportingPluginConfig,
				sharedConfig.DeltaRound,
				sharedConfig.MaxDurationQuery,
				sharedConfig.MaxDurationObservation,
				sharedConfig.MaxDurationReport,
				sharedConfig.MaxDurationShouldAcceptFinalizedReport,
				sharedConfig.MaxDurationShouldTransmitAcceptedReport,
			})
			if err != nil {
				logger.Error("ManagedOracle: error during NewReportingPlugin()", commontypes.LogFields{
					"error": err,
				})
				return
			}
			defer loghelper.CloseLogError(
				reportingPlugin,
				logger,
				"ManagedOracle: error during reportingPlugin.Close()",
			)
			if err := validateReportingPluginLimits(reportingPluginInfo.Limits); err != nil {
				logger.Error("ManagedOracle: invalid ReportingPluginInfo", commontypes.LogFields{
					"error":               err,
					"reportingPluginInfo": reportingPluginInfo,
				})
				return
			}

			lims, err := limits(sharedConfig.PublicConfig, reportingPluginInfo.Limits, onchainKeyring.MaxSignatureLength())
			if err != nil {
				logger.Error("ManagedOracle: error during limits", commontypes.LogFields{
					"error":               err,
					"publicConfig":        sharedConfig.PublicConfig,
					"reportingPluginInfo": reportingPluginInfo,
					"maxSigLen":           onchainKeyring.MaxSignatureLength(),
				})
				return
			}
			binNetEndpoint, err := netEndpointFactory.NewEndpoint(
				sharedConfig.ConfigDigest,
				peerIDs,
				v2bootstrappers,
				sharedConfig.F,
				lims,
			)
			if err != nil {
				logger.Error("ManagedOracle: error during NewEndpoint", commontypes.LogFields{
					"error":           err,
					"peerIDs":         peerIDs,
					"v2bootstrappers": v2bootstrappers,
				})
				return
			}

			// No need to binNetEndpoint.Start/Close since netEndpoint will handle that for us

			netEndpoint := shim.NewSerializingEndpoint(
				chTelemetrySend,
				sharedConfig.ConfigDigest,
				binNetEndpoint,
				childLogger,
				reportingPluginInfo.Limits,
			)
			if err := netEndpoint.Start(); err != nil {
				logger.Error("ManagedOracle: error during netEndpoint.Start()", commontypes.LogFields{
					"error":        err,
					"configDigest": sharedConfig.ConfigDigest,
				})
				return
			}
			defer loghelper.CloseLogError(
				netEndpoint,
				logger,
				"ManagedOracle: error during netEndpoint.Close()",
			)

			var reportQuorum int
			if reportingPluginInfo.UniqueReports {
				// We require greater than (n+f)/2 signatures to reach a byzantine
				// quorum. This ensures unique reports since each honest node will sign
				// at most one report for any given (epoch, round).
				//
				// Argument:
				//
				// (n+f)/2 = ((n-f)+f+f)/2 = (n-f)/2 + f
				//
				// There are (n-f) honest nodes, so to get two reports for an (epoch,
				// round) to reach  quorum, we'd need an honest node to sign two reports
				// which contradicts the assumption that an honest node will sign at
				// most one report for any given (epoch, round).
				reportQuorum = (sharedConfig.N()+sharedConfig.F)/2 + 1
			} else {
				reportQuorum = sharedConfig.F + 1
			}

			protocol.RunOracle(
				ctx,
				sharedConfig,
				contractTransmitter,
				database,
				oid,
				localConfig,
				childLogger,
				netEndpoint,
				offchainKeyring,
				onchainKeyring,
				shim.LimitCheckReportingPlugin{reportingPlugin, reportingPluginInfo.Limits},
				reportQuorum,
				shim.MakeTelemetrySender(chTelemetrySend, childLogger),
			)
		},
		localConfig,
		logger,
		offchainConfigDigester,
	)
}

func validateReportingPluginLimits(limits types.ReportingPluginLimits) error {
	var err error
	if !(0 <= limits.MaxQueryLength && limits.MaxQueryLength <= types.MaxMaxQueryLength) {
		err = multierr.Append(err, fmt.Errorf("MaxQueryLength (%v) out of range. Should be between 0 and %v", limits.MaxQueryLength, types.MaxMaxQueryLength))
	}
	if !(0 <= limits.MaxObservationLength && limits.MaxObservationLength <= types.MaxMaxObservationLength) {
		err = multierr.Append(err, fmt.Errorf("MaxObservationLength (%v) out of range. Should be between 0 and %v", limits.MaxObservationLength, types.MaxMaxObservationLength))
	}
	if !(0 <= limits.MaxReportLength && limits.MaxReportLength <= types.MaxMaxReportLength) {
		err = multierr.Append(err, fmt.Errorf("MaxReportLength (%v) out of range. Should be between 0 and %v", limits.MaxReportLength, types.MaxMaxReportLength))
	}
	return err
}
