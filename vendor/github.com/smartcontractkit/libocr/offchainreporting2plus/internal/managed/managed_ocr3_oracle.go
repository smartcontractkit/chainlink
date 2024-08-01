package managed

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"github.com/smartcontractkit/libocr/internal/metricshelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/managed/limits"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/protocol"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/ocr3/serialization"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/shim"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/subprocesses"
	"go.uber.org/multierr"
)

// RunManagedOCR3Oracle runs a "managed" version of protocol.RunOracle. It handles
// setting up telemetry, garbage collection, configuration updates, translating
// from commontypes.BinaryNetworkEndpoint to protocol.NetworkEndpoint, and
// creation/teardown of reporting plugins.
func RunManagedOCR3Oracle[RI any](
	ctx context.Context,

	v2bootstrappers []commontypes.BootstrapperLocator,
	configTracker types.ContractConfigTracker,
	contractTransmitter ocr3types.ContractTransmitter[RI],
	database ocr3types.Database,
	localConfig types.LocalConfig,
	logger loghelper.LoggerWithContext,
	metricsRegisterer prometheus.Registerer,
	monitoringEndpoint commontypes.MonitoringEndpoint,
	netEndpointFactory types.BinaryNetworkEndpointFactory,
	offchainConfigDigester types.OffchainConfigDigester,
	offchainKeyring types.OffchainKeyring,
	onchainKeyring ocr3types.OnchainKeyring[RI],
	reportingPluginFactory ocr3types.ReportingPluginFactory[RI],
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

	runWithContractConfig(
		ctx,

		configTracker,
		database,
		func(ctx context.Context, contractConfig types.ContractConfig, logger loghelper.LoggerWithContext) {
			skipResourceExhaustionChecks := localConfig.DevelopmentMode == types.EnableDangerousDevelopmentMode

			fromAccount, err := contractTransmitter.FromAccount()
			if err != nil {
				logger.Error("ManagedOCR3Oracle: error getting FromAccount", commontypes.LogFields{
					"error": err,
				})
				return
			}

			sharedConfig, oid, err := ocr3config.SharedConfigFromContractConfig(
				skipResourceExhaustionChecks,
				contractConfig,
				offchainKeyring,
				onchainKeyring,
				netEndpointFactory.PeerID(),
				fromAccount,
			)
			if err != nil {
				logger.Error("ManagedOCR3Oracle: error while updating config", commontypes.LogFields{
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

			reportingPlugin, reportingPluginInfo, err := reportingPluginFactory.NewReportingPlugin(ocr3types.ReportingPluginConfig{
				sharedConfig.ConfigDigest,
				oid,
				sharedConfig.N(),
				sharedConfig.F,
				sharedConfig.OnchainConfig,
				sharedConfig.ReportingPluginConfig,
				sharedConfig.DeltaRound,
				sharedConfig.MaxDurationQuery,
				sharedConfig.MaxDurationObservation,
				sharedConfig.MaxDurationShouldAcceptAttestedReport,
				sharedConfig.MaxDurationShouldTransmitAcceptedReport,
			})

			if err != nil {
				logger.Error("ManagedOCR3Oracle: error during NewReportingPlugin()", commontypes.LogFields{
					"error": err,
				})
				return
			}
			defer loghelper.CloseLogError(
				reportingPlugin,
				logger,
				"ManagedOCR3Oracle: error during reportingPlugin.Close()",
			)

			if err := validateOCR3ReportingPluginLimits(reportingPluginInfo.Limits); err != nil {
				logger.Error("ManagedOCR3Oracle: invalid ReportingPluginInfo", commontypes.LogFields{
					"error":               err,
					"reportingPluginInfo": reportingPluginInfo,
				})
				return
			}

			maxSigLen := onchainKeyring.MaxSignatureLength()
			lims, err := limits.OCR3Limits(sharedConfig.PublicConfig, reportingPluginInfo.Limits, maxSigLen)
			if err != nil {
				logger.Error("ManagedOCR3Oracle: error during limits", commontypes.LogFields{
					"error":                 err,
					"publicConfig":          sharedConfig.PublicConfig,
					"reportingPluginLimits": reportingPluginInfo.Limits,
					"maxSigLen":             maxSigLen,
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
				logger.Error("ManagedOCR3Oracle: error during NewEndpoint", commontypes.LogFields{
					"error":           err,
					"peerIDs":         peerIDs,
					"v2bootstrappers": v2bootstrappers,
				})
				return
			}

			// No need to binNetEndpoint.Start/Close since netEndpoint will handle that for us

			netEndpoint := shim.NewOCR3SerializingEndpoint[RI](
				chTelemetrySend,
				sharedConfig.ConfigDigest,
				binNetEndpoint,
				maxSigLen,
				childLogger,
				reportingPluginInfo.Limits,
				sharedConfig.N(),
				sharedConfig.F,
			)
			if err := netEndpoint.Start(); err != nil {
				logger.Error("ManagedOCR3Oracle: error during netEndpoint.Start()", commontypes.LogFields{
					"error":        err,
					"configDigest": sharedConfig.ConfigDigest,
				})
				return
			}
			defer loghelper.CloseLogError(
				netEndpoint,
				logger,
				"ManagedOCR3Oracle: error during netEndpoint.Close()",
			)

			protocol.RunOracle[RI](
				ctx,
				sharedConfig,
				contractTransmitter,
				&shim.SerializingOCR3Database{database},
				oid,
				localConfig,
				childLogger,
				registerer,
				netEndpoint,
				offchainKeyring,
				onchainKeyring,
				shim.LimitCheckOCR3ReportingPlugin[RI]{reportingPlugin, reportingPluginInfo.Limits},
				shim.MakeOCR3TelemetrySender(chTelemetrySend, childLogger),
			)
		},
		localConfig,
		logger,
		offchainConfigDigester,
	)
}

func validateOCR3ReportingPluginLimits(limits ocr3types.ReportingPluginLimits) error {
	var err error
	if !(0 <= limits.MaxQueryLength && limits.MaxQueryLength <= ocr3types.MaxMaxQueryLength) {
		err = multierr.Append(err, fmt.Errorf("MaxQueryLength (%v) out of range. Should be between 0 and %v", limits.MaxQueryLength, ocr3types.MaxMaxQueryLength))
	}
	if !(0 <= limits.MaxObservationLength && limits.MaxObservationLength <= ocr3types.MaxMaxObservationLength) {
		err = multierr.Append(err, fmt.Errorf("MaxObservationLength (%v) out of range. Should be between 0 and %v", limits.MaxObservationLength, ocr3types.MaxMaxObservationLength))
	}
	if !(0 <= limits.MaxOutcomeLength && limits.MaxOutcomeLength <= ocr3types.MaxMaxOutcomeLength) {
		err = multierr.Append(err, fmt.Errorf("MaxOutcomeLength (%v) out of range. Should be between 0 and %v", limits.MaxOutcomeLength, ocr3types.MaxMaxOutcomeLength))
	}
	if !(0 <= limits.MaxReportLength && limits.MaxReportLength <= ocr3types.MaxMaxReportLength) {
		err = multierr.Append(err, fmt.Errorf("MaxReportLength (%v) out of range. Should be between 0 and %v", limits.MaxReportLength, ocr3types.MaxMaxReportLength))
	}
	if !(0 <= limits.MaxReportCount && limits.MaxReportCount <= ocr3types.MaxMaxReportCount) {
		err = multierr.Append(err, fmt.Errorf("MaxReportCount (%v) out of range. Should be between 0 and %v", limits.MaxReportCount, ocr3types.MaxMaxReportCount))
	}
	return err
}
