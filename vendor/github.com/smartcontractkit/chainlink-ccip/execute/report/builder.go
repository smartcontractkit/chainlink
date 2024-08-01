package report

import (
	"context"
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-ccip/execute/types"
	"github.com/smartcontractkit/chainlink-ccip/plugintypes"
)

var _ ExecReportBuilder = &execReportBuilder{}

type ExecReportBuilder interface {
	Add(report plugintypes.ExecutePluginCommitData) (plugintypes.ExecutePluginCommitData, error)
	Build() ([]cciptypes.ExecutePluginReportSingleChain, error)
}

func NewBuilder(
	ctx context.Context,
	logger logger.Logger,
	hasher cciptypes.MessageHasher,
	tokenDataReader types.TokenDataReader,
	encoder cciptypes.ExecutePluginCodec,
	maxReportSizeBytes uint64,
	maxGas uint64,
) ExecReportBuilder {
	return &execReportBuilder{
		ctx:  ctx,
		lggr: logger,

		tokenDataReader: tokenDataReader,
		encoder:         encoder,
		hasher:          hasher,

		maxReportSizeBytes: maxReportSizeBytes,
		maxGas:             maxGas,
	}
}

// validationMetadata contains all metadata needed to accumulate results across multiple reports and messages.
type validationMetadata struct {
	encodedSizeBytes uint64

	// TODO: gas limit
	//gas             uint64
}

type execReportBuilder struct {
	ctx  context.Context // TODO: remove context from builder so that it can be pure?
	lggr logger.Logger

	// Providers
	tokenDataReader types.TokenDataReader
	encoder         cciptypes.ExecutePluginCodec
	hasher          cciptypes.MessageHasher

	// Config
	maxReportSizeBytes uint64
	maxGas             uint64

	// State
	accumulated validationMetadata

	// Result
	execReports []cciptypes.ExecutePluginReportSingleChain
}

func (b *execReportBuilder) Add(
	commitReport plugintypes.ExecutePluginCommitData,
) (plugintypes.ExecutePluginCommitData, error) {
	execReport, updatedReport, err := b.buildSingleChainReport(b.ctx, commitReport)

	// No messages fit into the report, move to next report
	if errors.Is(err, ErrEmptyReport) {
		return commitReport, nil
	}
	if err != nil {
		return commitReport, fmt.Errorf("unable to add a single chain report: %w", err)
	}

	b.execReports = append(b.execReports, execReport)

	return updatedReport, nil
}

func (b *execReportBuilder) Build() ([]cciptypes.ExecutePluginReportSingleChain, error) {
	b.lggr.Infow(
		"selected commit reports for execution report",
		"numReports", len(b.execReports),
		"sizeBytes", b.accumulated.encodedSizeBytes,
		"maxSize", b.maxReportSizeBytes)
	return b.execReports, nil
}
