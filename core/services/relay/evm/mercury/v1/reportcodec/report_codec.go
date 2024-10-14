package reportcodec

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	pkgerrors "github.com/pkg/errors"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/utils"
	reporttypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1/types"
)

// NOTE:
// This report codec is based on the original median evmreportcodec
// here:
// https://github.com/smartcontractkit/offchain-reporting/blob/master/lib/offchainreporting2/reportingplugin/median/evmreportcodec/reportcodec.go
var ReportTypes = reporttypes.GetSchema()
var maxReportLength = 32 * len(ReportTypes) // each arg is 256 bit EVM word

var _ v1.ReportCodec = &ReportCodec{}

type ReportCodec struct {
	logger logger.Logger
	feedID utils.FeedID
}

func NewReportCodec(feedID [32]byte, lggr logger.Logger) *ReportCodec {
	return &ReportCodec{lggr, feedID}
}

func (r *ReportCodec) BuildReport(ctx context.Context, rf v1.ReportFields) (ocrtypes.Report, error) {
	var merr error
	if rf.BenchmarkPrice == nil {
		merr = errors.Join(merr, errors.New("benchmarkPrice may not be nil"))
	}
	if rf.Bid == nil {
		merr = errors.Join(merr, errors.New("bid may not be nil"))
	}
	if rf.Ask == nil {
		merr = errors.Join(merr, errors.New("ask may not be nil"))
	}
	if len(rf.CurrentBlockHash) != 32 {
		merr = errors.Join(merr, fmt.Errorf("invalid length for currentBlockHash, expected: 32, got: %d", len(rf.CurrentBlockHash)))
	}
	if merr != nil {
		return nil, merr
	}
	var currentBlockHash common.Hash
	copy(currentBlockHash[:], rf.CurrentBlockHash)

	reportBytes, err := ReportTypes.Pack(r.feedID, rf.Timestamp, rf.BenchmarkPrice, rf.Bid, rf.Ask, uint64(rf.CurrentBlockNum), currentBlockHash, uint64(rf.ValidFromBlockNum), rf.CurrentBlockTimestamp)
	return ocrtypes.Report(reportBytes), pkgerrors.Wrap(err, "failed to pack report blob")
}

// Maximum length in bytes of Report returned by BuildReport. Used for
// defending against spam attacks.
func (r *ReportCodec) MaxReportLength(ctx context.Context, n int) (int, error) {
	return maxReportLength, nil
}

func (r *ReportCodec) CurrentBlockNumFromReport(ctx context.Context, report ocrtypes.Report) (int64, error) {
	decoded, err := r.Decode(report)
	if err != nil {
		return 0, err
	}
	if decoded.CurrentBlockNum > math.MaxInt64 {
		return 0, fmt.Errorf("CurrentBlockNum=%d overflows max int64", decoded.CurrentBlockNum)
	}
	return int64(decoded.CurrentBlockNum), nil
}

func (r *ReportCodec) Decode(report ocrtypes.Report) (*reporttypes.Report, error) {
	return reporttypes.Decode(report)
}

func (r *ReportCodec) BenchmarkPriceFromReport(ctx context.Context, report ocrtypes.Report) (*big.Int, error) {
	decoded, err := r.Decode(report)
	if err != nil {
		return nil, err
	}
	return decoded.BenchmarkPrice, nil
}

func (r *ReportCodec) ObservationTimestampFromReport(ctx context.Context, report ocrtypes.Report) (uint32, error) {
	decoded, err := r.Decode(report)
	if err != nil {
		return 0, err
	}
	return decoded.ObservationsTimestamp, nil
}
