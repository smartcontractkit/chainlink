package plugin

import (
	"context"
	"encoding/hex"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

const (
	defaultMaxPhaseOutputBytes              = uint32(ocr3_1types.MaxMaxObservationBytes)
	defaultMaxReportLengthBytes             = uint32(1_000_000) // 1 MB
	defaultMaxReportCount                   = 100
	defaultRequestExpiry                    = 20 * time.Second
	defaultHistoricalOutcomeExpirySeqNrSpan = uint64(4)
)

type SetRequestTimeout func(timeout time.Duration)

type factory struct {
	store *requests.Store[*oracle.ConsensusRequest]

	// Request timeout is set by the plugin factory and used by the reporting plugin to set the timeout for requests
	// created in the capability
	setRequestTimeout SetRequestTimeout
	lggr              logger.Logger
	metrics           *metrics.Metrics

	defaultKeyBundleIDForConsensusFailure string
	maxRequestOutcomeSize                 int

	services.StateMachine
}

func NewReportingPluginFactory(lggr logger.Logger, metrics *metrics.Metrics, s *requests.Store[*oracle.ConsensusRequest],
	setRequestTimeout SetRequestTimeout, defaultKeyBundleIDForConsensusFailure string,
	maxRequestOutcomeSize int,
) (*factory, error) {
	return &factory{
		store:                                 s,
		setRequestTimeout:                     setRequestTimeout,
		lggr:                                  logger.Named(lggr, "ConsensusCapabilityPluginFactory"),
		metrics:                               metrics,
		defaultKeyBundleIDForConsensusFailure: defaultKeyBundleIDForConsensusFailure,
		maxRequestOutcomeSize:                 maxRequestOutcomeSize,
	}, nil
}

func (o *factory) NewReportingPlugin(_ context.Context, config ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[[]byte], ocr3types.ReportingPluginInfo, error) {
	var configProto types.ReportingPluginConfig
	err := proto.Unmarshal(config.OffchainConfig, &configProto)
	if err != nil {
		// an empty byte array will be unmarshalled into zero values without error
		return nil, ocr3types.ReportingPluginInfo{}, err
	}
	if configProto.MaxQueryLengthBytes <= 0 {
		configProto.MaxQueryLengthBytes = defaultMaxPhaseOutputBytes
	}
	if configProto.MaxObservationLengthBytes <= 0 {
		configProto.MaxObservationLengthBytes = defaultMaxPhaseOutputBytes
	}
	if configProto.MaxOutcomeLengthBytes <= 0 {
		configProto.MaxOutcomeLengthBytes = defaultMaxPhaseOutputBytes
	}
	if configProto.MaxReportLengthBytes <= 0 {
		configProto.MaxReportLengthBytes = defaultMaxReportLengthBytes
	}
	if configProto.MaxReportCount <= 0 {
		configProto.MaxReportCount = defaultMaxReportCount
	}
	if configProto.RequestTimeout == nil {
		configProto.RequestTimeout = durationpb.New(defaultRequestExpiry)
	}
	o.setRequestTimeout(configProto.RequestTimeout.AsDuration())

	// Historical outcomes prevent the reporting plugin from processing a request more than once by allowing the plugin
	// to determine if an outcome has already been reported for a given request in the recent past.  The recent past i.e. 'size' of history to retain (the span), needs
	// to take into consideration the likelihood of a request being re-submitted after a certain number of rounds versus the storage cost of retaining history in the
	// plugins outcome.

	// An example of a scenario that could result in a request being resubmitted for reprocessing if there is insufficient history to prevent is:
	//
	// Due to the async nature of the transmit call, it is possible, in the case where a second round occurs before 2f+1
	// nodes have received the transmit call for a given round, for a second possibly different outcome to be generated, which if a node
	// misses the previous transmit call altogether, will result in it seeing a different outcome.

	if configProto.HistoricalOutcomeExpirySeqNrSpan == 0 {
		configProto.HistoricalOutcomeExpirySeqNrSpan = defaultHistoricalOutcomeExpirySeqNrSpan
	}

	rp, err := NewReportingPlugin(o.lggr, o.metrics, config.F, config.N, o.store, &configProto, o.defaultKeyBundleIDForConsensusFailure,
		o.maxRequestOutcomeSize)
	limits31 := reportingPluginLimitsOCR31(&configProto)
	rpInfo31 := ocr3_1types.ReportingPluginInfo1{
		Name:   "Consensus Capability Plugin",
		Limits: limits31,
	}
	if vErr := rpInfo31.Validate(); vErr != nil {
		o.lggr.Errorw("invalid OCR 3.1 reporting plugin limits", "error", vErr)
		return nil, ocr3types.ReportingPluginInfo{}, vErr
	}
	rpInfo := reportingPluginInfoOCR3From31(rpInfo31)
	o.lggr.Infow("Created OCR3 consensus capability reporting plugin with config",
		// "OffchainConfig" fields - internal to our plugin
		"maxQueryLengthBytes", configProto.MaxQueryLengthBytes,
		"maxObservationLengthBytes", configProto.MaxObservationLengthBytes,
		"maxOutcomeLengthBytes", configProto.MaxOutcomeLengthBytes,
		"maxReportLengthBytes", configProto.MaxReportLengthBytes,
		"maxReportCount", configProto.MaxReportCount,
		"maxBatchSize(UNUSED)", configProto.MaxBatchSize, // UNUSED - batch size is now determined by max byte length limits
		"outcomePruningThreshold", configProto.OutcomePruningThreshold,
		"requestTimeout", configProto.RequestTimeout.AsDuration(),
		"historicalOutcomeExpirySeqNrSpan", configProto.HistoricalOutcomeExpirySeqNrSpan,
		// top-level OCR config fields (ocr3types.ReportingPluginConfig)
		"configDigest", hex.EncodeToString(config.ConfigDigest[:]),
		"oracleID", config.OracleID,
		"n", config.N,
		"f", config.F,
		"estimatedRoundInterval", config.EstimatedRoundInterval,
		"maxDurationQuery", config.MaxDurationQuery,
		"maxDurationObservation", config.MaxDurationObservation,
		"maxDurationShouldAcceptAttestedReport", config.MaxDurationShouldAcceptAttestedReport,
		"maxDurationShouldTransmitAcceptedReport", config.MaxDurationShouldTransmitAcceptedReport,
		// extra fields from capability job spec config
		"maxRequestOutcomeSize", o.maxRequestOutcomeSize, // NOTE: use limits instead of a dedicated job spec field?
		"defaultKeyBundleIDForConsensusFailure", o.defaultKeyBundleIDForConsensusFailure,
	)

	return AsOCR3ReportingPlugin(rp), rpInfo, err
}

func (o *factory) Start(ctx context.Context) error {
	return o.StartOnce("ConsensusCapabilityReportingPlugin", func() error {
		return nil
	})
}

func (o *factory) Close() error {
	return o.StopOnce("ConsensusCapabilityReportingPlugin", func() error {
		return nil
	})
}

func (o *factory) Name() string { return o.lggr.Name() }

func (o *factory) HealthReport() map[string]error {
	return map[string]error{o.Name(): o.Healthy()}
}
