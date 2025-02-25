package ccip

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"time"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-integrations/evm/gas"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// The disjunct methods in IncompleteSourceCommitStoreReader and IncompleteDestCommitStoreReader satisfy the full
// CommitStoreReader iface in Union
var _ cciptypes.CommitStoreReader = (*ProviderProxyCommitStoreReader)(nil)

// ProviderProxyCommitStoreReader is a CommitStoreReader that proxies to two custom provider grpc backed implementations
// of a CommitStoreReader.
// [ProviderProxyCommitStoreReader] lives in the memory space of the reporting plugin factory and reporting plugin, and should have no chain-specific details.
// Why? Historical implementations of a commit store consumed in reporting plugins mixed usage of a gas estimator from
// the source relayer and contract read and write abilities to a dest relayer. This is not valid in LOOP world.
type ProviderProxyCommitStoreReader struct {
	srcCommitStoreReader IncompleteSourceCommitStoreReader
	dstCommitStoreReader IncompleteDestCommitStoreReader
}

// IncompleteSourceCommitStoreReader contains only the methods of CommitStoreReader that are serviced by the source chain/relayer.
type IncompleteSourceCommitStoreReader interface {
	ChangeConfig(ctx context.Context, onchainConfig []byte, offchainConfig []byte) (cciptypes.Address, error)
	GasPriceEstimator(ctx context.Context) (cciptypes.GasPriceEstimatorCommit, error)
	OffchainConfig(ctx context.Context) (cciptypes.CommitOffchainConfig, error)
	io.Closer
}

// IncompleteDestCommitStoreReader contains only the methods of CommitStoreReader that are serviced by the dest chain/relayer.
type IncompleteDestCommitStoreReader interface {
	DecodeCommitReport(ctx context.Context, report []byte) (cciptypes.CommitStoreReport, error)
	EncodeCommitReport(ctx context.Context, report cciptypes.CommitStoreReport) ([]byte, error)
	GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confirmations int) ([]cciptypes.CommitStoreReportWithTxMeta, error)
	GetCommitReportMatchingSeqNum(ctx context.Context, seqNum uint64, confirmations int) ([]cciptypes.CommitStoreReportWithTxMeta, error)
	GetCommitStoreStaticConfig(ctx context.Context) (cciptypes.CommitStoreStaticConfig, error)
	GetExpectedNextSequenceNumber(ctx context.Context) (uint64, error)
	GetLatestPriceEpochAndRound(ctx context.Context) (uint64, error)
	IsBlessed(ctx context.Context, root [32]byte) (bool, error)
	IsDestChainHealthy(ctx context.Context) (bool, error)
	IsDown(ctx context.Context) (bool, error)
	VerifyExecutionReport(ctx context.Context, report cciptypes.ExecReport) (bool, error)
	io.Closer
}

func NewProviderProxyCommitStoreReader(srcReader cciptypes.CommitStoreReader, dstReader cciptypes.CommitStoreReader) *ProviderProxyCommitStoreReader {
	return &ProviderProxyCommitStoreReader{
		srcCommitStoreReader: srcReader,
		dstCommitStoreReader: dstReader,
	}
}

// ChangeConfig updates the offchainConfig values for the source relayer gas estimator by calling ChangeConfig
// on the source relayer. Once this is called, GasPriceEstimator and OffchainConfig can be called.
func (p *ProviderProxyCommitStoreReader) ChangeConfig(ctx context.Context, onchainConfig []byte, offchainConfig []byte) (cciptypes.Address, error) {
	return p.srcCommitStoreReader.ChangeConfig(ctx, onchainConfig, offchainConfig)
}

func (p *ProviderProxyCommitStoreReader) DecodeCommitReport(ctx context.Context, report []byte) (cciptypes.CommitStoreReport, error) {
	return p.dstCommitStoreReader.DecodeCommitReport(ctx, report)
}

func (p *ProviderProxyCommitStoreReader) EncodeCommitReport(ctx context.Context, report cciptypes.CommitStoreReport) ([]byte, error) {
	return p.dstCommitStoreReader.EncodeCommitReport(ctx, report)
}

// GasPriceEstimator constructs a gas price estimator on the source relayer
func (p *ProviderProxyCommitStoreReader) GasPriceEstimator(ctx context.Context) (cciptypes.GasPriceEstimatorCommit, error) {
	return p.srcCommitStoreReader.GasPriceEstimator(ctx)
}

func (p *ProviderProxyCommitStoreReader) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confirmations int) ([]cciptypes.CommitStoreReportWithTxMeta, error) {
	return p.dstCommitStoreReader.GetAcceptedCommitReportsGteTimestamp(ctx, ts, confirmations)
}

func (p *ProviderProxyCommitStoreReader) GetCommitReportMatchingSeqNum(ctx context.Context, seqNum uint64, confirmations int) ([]cciptypes.CommitStoreReportWithTxMeta, error) {
	return p.dstCommitStoreReader.GetCommitReportMatchingSeqNum(ctx, seqNum, confirmations)
}

func (p *ProviderProxyCommitStoreReader) GetCommitStoreStaticConfig(ctx context.Context) (cciptypes.CommitStoreStaticConfig, error) {
	return p.dstCommitStoreReader.GetCommitStoreStaticConfig(ctx)
}

func (p *ProviderProxyCommitStoreReader) GetExpectedNextSequenceNumber(ctx context.Context) (uint64, error) {
	return p.dstCommitStoreReader.GetExpectedNextSequenceNumber(ctx)
}

func (p *ProviderProxyCommitStoreReader) GetLatestPriceEpochAndRound(ctx context.Context) (uint64, error) {
	return p.dstCommitStoreReader.GetLatestPriceEpochAndRound(ctx)
}

func (p *ProviderProxyCommitStoreReader) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	return p.dstCommitStoreReader.IsBlessed(ctx, root)
}

func (p *ProviderProxyCommitStoreReader) IsDestChainHealthy(ctx context.Context) (bool, error) {
	return p.dstCommitStoreReader.IsDestChainHealthy(ctx)
}

func (p *ProviderProxyCommitStoreReader) IsDown(ctx context.Context) (bool, error) {
	return p.dstCommitStoreReader.IsDown(ctx)
}

func (p *ProviderProxyCommitStoreReader) OffchainConfig(ctx context.Context) (cciptypes.CommitOffchainConfig, error) {
	return p.srcCommitStoreReader.OffchainConfig(ctx)
}

func (p *ProviderProxyCommitStoreReader) VerifyExecutionReport(ctx context.Context, report cciptypes.ExecReport) (bool, error) {
	return p.dstCommitStoreReader.VerifyExecutionReport(ctx, report)
}

// SetGasEstimator is invalid on ProviderProxyCommitStoreReader. The provider based impl's do not have SetGasEstimator
// defined, so this serves no purpose other than satisfying an interface.
func (p *ProviderProxyCommitStoreReader) SetGasEstimator(ctx context.Context, gpe gas.EvmFeeEstimator) error {
	return fmt.Errorf("invalid usage of ProviderProxyCommitStoreReader")
}

// SetSourceMaxGasPrice is invalid on ProviderProxyCommitStoreReader. The provider based impl's do not have SetSourceMaxGasPrice
// defined, so this serves no purpose other than satisfying an interface.
func (p *ProviderProxyCommitStoreReader) SetSourceMaxGasPrice(ctx context.Context, sourceMaxGasPrice *big.Int) error {
	return fmt.Errorf("invalid usage of ProviderProxyCommitStoreReader")
}

func (p *ProviderProxyCommitStoreReader) Close() error {
	return multierr.Append(p.srcCommitStoreReader.Close(), p.dstCommitStoreReader.Close())
}
