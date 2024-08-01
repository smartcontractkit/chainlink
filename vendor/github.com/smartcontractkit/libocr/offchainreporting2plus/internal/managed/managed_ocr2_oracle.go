package managed

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr2config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/managed/limits"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr2/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr2/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/shim"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
	"go.uber.org/multierr"
)

// RunManagedOCR2Oracle runs a "managed" version of protocol.RunOracle. It handles
// setting up telemetry, garbage collection, configuration updates, translating
// from commontypes.BinaryNetworkEndpoint to protocol.NetworkEndpoint, and
// creation/teardown of reporting plugins.
func RunManagedOCR2Oracle(
	ctx context.Context,

	v2bootstrappers []commontypes.BootstrapperLocator,
	configTracker types.ContractConfigTracker,
	contractTransmitter types.ContractTransmitter,
	database types.Database,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
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

	metricsRegistererWrapper := metricshelper.NewPrometheusRegistererWrapper(metricsRegisterer, logger)

	subs.Go(func() {
		collectGarbage(ctx, database, localConfig, logger)
	})

	runWithContractConfig(
		ctx,

		configTracker,
		database,
		func(ctx context.Context, contractConfig types.ContractConfig, logger loghelper.LoggerWithContext) {
			skipResourceExhaustionChecks := localConfig.DevelopmentMode == types.EnableDangerousDevelopmentMode

			fromAccount, err := contractTransmitter.FromAccount()
			if err != nil {
				logger.Error("ManagedOCR2Oracle: error getting FromAccount", commontypes.LogFields{
					"error": err,
				})
				return
			}

			sharedConfig, oid, err := ocr2config.SharedConfigFromContractConfig(
				skipResourceExhaustionChecks,
				contractConfig,
				offchainKeyring,
				onchainKeyring,
				netEndpointFactory.PeerID(),
				fromAccount,
			)
			if err != nil {
				logger.Error("ManagedOCR2Oracle: error while updating config", commontypes.LogFields{
					"error": err,
				})
				return
			}

			registerer := prometheus.WrapRegistererWith(
				prometheus.Labels{
					// disambiguate different protocol instances by configDigest
					"config_digest": sharedConfig.ConfigDigest.String(),
					// disambiguate different oracle instances by offchainPublicKey
					"offchain_public_key": fmt.Sprintf("%x", offchainKeyring.OffchainPublicKey()),
				},
				metricsRegistererWrapper,
			)

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
				logger.Error("ManagedOCR2Oracle: error during NewReportingPlugin()", commontypes.LogFields{
					"error": err,
				})
				return
			}
			defer loghelper.CloseLogError(
				reportingPlugin,
				logger,
				"ManagedOCR2Oracle: error during reportingPlugin.Close()",
			)
			if err := validateReportingPluginLimits(reportingPluginInfo.Limits); err != nil {
				logger.Error("ManagedOCR2Oracle: invalid ReportingPluginInfo", commontypes.LogFields{
					"error":               err,
					"reportingPluginInfo": reportingPluginInfo,
				})
				return
			}

			maxSigLen := onchainKeyring.MaxSignatureLength()
			lims, err := limits.OCR2Limits(sharedConfig.PublicConfig, reportingPluginInfo.Limits, maxSigLen)
			if err != nil {
				logger.Error("ManagedOCR2Oracle: error during limits", commontypes.LogFields{
					"error":               err,
					"publicConfig":        sharedConfig.PublicConfig,
					"reportingPluginInfo": reportingPluginInfo,
					"maxSigLen":           maxSigLen,
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
				logger.Error("ManagedOCR2Oracle: error during NewEndpoint", commontypes.LogFields{
					"error":           err,
					"peerIDs":         peerIDs,
					"v2bootstrappers": v2bootstrappers,
				})
				return
			}

			// No need to binNetEndpoint.Start/Close since netEndpoint will handle that for us

			netEndpoint := shim.NewOCR2SerializingEndpoint(
				chTelemetrySend,
				sharedConfig.ConfigDigest,
				binNetEndpoint,
				childLogger,
				reportingPluginInfo.Limits,
			)
			if err := netEndpoint.Start(); err != nil {
				logger.Error("ManagedOCR2Oracle: error during netEndpoint.Start()", commontypes.LogFields{
					"error":        err,
					"configDigest": sharedConfig.ConfigDigest,
				})
				return
			}
			defer loghelper.CloseLogError(
				netEndpoint,
				logger,
				"ManagedOCR2Oracle: error during netEndpoint.Close()",
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
				registerer,
				netEndpoint,
				offchainKeyring,
				onchainKeyring,
				shim.LimitCheckReportingPlugin{reportingPlugin, reportingPluginInfo.Limits},
				reportQuorum,
				shim.MakeOCR2TelemetrySender(chTelemetrySend, childLogger),
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
