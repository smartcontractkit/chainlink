package encoding

import (
	"encoding/binary"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	mercuryv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1"
	mercuryv2 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v2"
	mercuryv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/pkg/errors"
)

type FeedIDPrefix uint16

const (
	_         FeedIDPrefix = 0 // reserved to prevent errors where a zero-default creeps through somewhere
	REPORT_V1 FeedIDPrefix = 1
	REPORT_V2 FeedIDPrefix = 2
	REPORT_V3 FeedIDPrefix = 3
	_         FeedIDPrefix = 0xFFFF // reserved for future use
)

func DecodeSchemaVersionFromFeedId(feedID [32]byte) (FeedIDPrefix, error) {
	schemaVersion := FeedIDPrefix(binary.BigEndian.Uint16(feedID[:2]))
	if schemaVersion != REPORT_V1 && schemaVersion != REPORT_V2 && schemaVersion != REPORT_V3 {
		return 0, errors.Errorf("invalid schema version: %d", schemaVersion)
	}
	return schemaVersion, nil
}

type ReportDecoder interface {
	GetSchemaVersion() FeedIDPrefix

	DecodeAsV1() (*mercuryv1.Report, error)
	DecodeAsV2() (*mercuryv2.Report, error)
	DecodeAsV3() (*mercuryv3.Report, error)
}

type reportDecoder struct {
	report        ocrtypes.Report
	feedId        [32]byte
	schemaVersion FeedIDPrefix

	lggr logger.Logger
}

var _ ReportDecoder = (*reportDecoder)(nil)

func NewReportDecoder(report ocrtypes.Report, lggr logger.Logger) (ReportDecoder, error) {
	var feedId [32]byte
	if n := copy(feedId[:], report); n != 32 {
		return &reportDecoder{}, errors.Errorf("invalid length for report: %d", len(report))
	}

	schemaVersion, err := DecodeSchemaVersionFromFeedId(feedId)
	if err != nil {
		return nil, err
	}

	switch schemaVersion {
	case REPORT_V1:
		return &reportDecoder{
			report:        report,
			feedId:        feedId,
			schemaVersion: schemaVersion,
		}, nil
	case REPORT_V2:
		return &reportDecoder{
			report:        report,
			feedId:        feedId,
			schemaVersion: schemaVersion,
		}, nil
	case REPORT_V3:
		return &reportDecoder{
			report:        report,
			feedId:        feedId,
			schemaVersion: schemaVersion,
		}, nil
	default:
		return &reportDecoder{}, errors.Errorf("invalid schema version: %d", schemaVersion)
	}
}

func (d *reportDecoder) DecodeAsV1() (*mercuryv1.Report, error) {
	if d.schemaVersion != REPORT_V1 {
		return nil, errors.Errorf("invalid schema version: %d", d.schemaVersion)
	}

	reportCodec := mercuryv1.NewReportCodec(d.feedId, d.lggr)
	return reportCodec.Decode(d.report)
}

func (d *reportDecoder) DecodeAsV2() (*mercuryv2.Report, error) {
	if d.schemaVersion != REPORT_V2 {
		return nil, errors.Errorf("invalid schema version: %d", d.schemaVersion)
	}

	reportCodec := mercuryv2.NewReportCodec(d.feedId, d.lggr)
	return reportCodec.Decode(d.report)
}

func (d *reportDecoder) DecodeAsV3() (*mercuryv3.Report, error) {
	if d.schemaVersion != REPORT_V3 {
		return nil, errors.Errorf("invalid schema version: %d", d.schemaVersion)
	}

	reportCodec := mercuryv3.NewReportCodec(d.feedId, d.lggr)
	return reportCodec.Decode(d.report)
}

func (d *reportDecoder) GetSchemaVersion() FeedIDPrefix {
	return d.schemaVersion
}
