package plugin

import (
	"context"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle"
	oracletypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/quorumhelper"

	ocrtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var _ ocr3types.ReportingPlugin[[]byte] = (*reportingPlugin)(nil)

type reportingPlugin struct {
	store *requests.Store[*oracle.ConsensusRequest]

	f int
	n int

	// outcomeExpirySeqNrSpan is the duration, expressed as a seq number span, after which a request outcome will be pruned from the plugins outcome
	outcomeExpirySeqNrSpan uint64

	config  *ocrtypes.ReportingPluginConfig
	metrics *metrics.Metrics

	// defaultKeyBundleIDForConsensusFailure is the key bundle ID to be used when reporting consensus failures before consensus is reached on request metadata
	defaultKeyBundleIDForConsensusFailure string
	maxRequestOutcomeSize                 int
	maxReportLengthBytes                  int
	maxNumberOfReports                    int

	lggr logger.Logger
}

// NewReportingPlugin creates a new reporting plugin for the OCR3 capability
func NewReportingPlugin(lggr logger.Logger, metrics *metrics.Metrics, f int, n int, store *requests.Store[*oracle.ConsensusRequest],
	configProto *ocrtypes.ReportingPluginConfig, defaultKeyBundleIDForConsensusFailure string,
	maxRequestOutcomeSize int) (*reportingPlugin, error) {
	return &reportingPlugin{
		store:                                 store,
		f:                                     f,
		n:                                     n,
		outcomeExpirySeqNrSpan:                configProto.HistoricalOutcomeExpirySeqNrSpan,
		lggr:                                  logger.Named(lggr, "CapabilityConsensusReportingPlugin"),
		config:                                configProto,
		metrics:                               metrics,
		defaultKeyBundleIDForConsensusFailure: defaultKeyBundleIDForConsensusFailure,
		maxRequestOutcomeSize:                 maxRequestOutcomeSize, // NOTE: can we simply use configProto.MaxOutcomeLengthBytes here?
		maxReportLengthBytes:                  int(configProto.MaxReportLengthBytes),
		maxNumberOfReports:                    int(configProto.MaxReportCount),
	}, nil
}

func ToRequestMetaData(metadata oracle.ConsensusRequestMetadata) *oracletypes.RequestMetaData {
	return &oracletypes.RequestMetaData{
		RequestId:                metadata.RequestID(),
		WorkflowExecutionId:      metadata.WorkflowExecutionID,
		WorkflowStepReference:    metadata.ReferenceID,
		WorkflowId:               metadata.WorkflowID,
		WorkflowOwner:            metadata.WorkflowOwner,
		WorkflowName:             metadata.WorkflowName,
		WorkflowDonId:            metadata.WorkflowDonID,
		WorkflowDonConfigVersion: metadata.WorkflowDonConfigVersion,
		ReportId:                 metadata.ReportID,
		KeyBundleId:              metadata.KeyBundleID,
		RequestType:              metadata.RequestType,
	}
}

func (r *reportingPlugin) ValidateObservation(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query, ao types.AttributedObservation) error {
	return nil
}

func (r *reportingPlugin) ObservationQuorum(ctx context.Context, outctx ocr3types.OutcomeContext, query types.Query, aos []types.AttributedObservation) (bool, error) {
	return quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumTwoFPlusOne, r.n, r.f, aos), nil
}

func (r *reportingPlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, rwi ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	// True because we always want to transmit a report
	return true, nil
}

func (r *reportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, rwi ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	// True because we always want to transmit a report
	return true, nil
}

func (r *reportingPlugin) Close() error {
	return nil
}
