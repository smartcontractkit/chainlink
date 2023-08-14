package encoding

import (
	"encoding/binary"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
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

func SchemaVersionFromFeedId(feedID [32]byte) FeedIDPrefix {
	return FeedIDPrefix(binary.BigEndian.Uint16(feedID[:2]))
}

func DecodeV1(report ocrtypes.Report, lggr logger.Logger) (*mercuryv1.Report, error) {
	feedId, err := mercury.FeedIDFromReport(report)
	if err != nil {
		return nil, err
	}

	schemaVersion := SchemaVersionFromFeedId(feedId)
	if schemaVersion != REPORT_V1 {
		return nil, errors.Errorf("invalid schema version: %d", schemaVersion)
	}

	reportCodec := mercuryv1.NewReportCodec(feedId, lggr)
	return reportCodec.Decode(report)
}

func DecodeV2(report ocrtypes.Report, lggr logger.Logger) (*mercuryv2.Report, error) {
	feedId, err := mercury.FeedIDFromReport(report)
	if err != nil {
		return nil, err
	}

	schemaVersion := SchemaVersionFromFeedId(feedId)
	if schemaVersion != REPORT_V2 {
		return nil, errors.Errorf("invalid schema version: %d", schemaVersion)
	}

	reportCodec := mercuryv2.NewReportCodec(feedId, lggr)
	return reportCodec.Decode(report)
}

func DecodeV3(report ocrtypes.Report, lggr logger.Logger) (*mercuryv3.Report, error) {
	feedId, err := mercury.FeedIDFromReport(report)
	if err != nil {
		return nil, err
	}

	schemaVersion := SchemaVersionFromFeedId(feedId)
	if schemaVersion != REPORT_V3 {
		return nil, errors.Errorf("invalid schema version: %d", schemaVersion)
	}

	reportCodec := mercuryv3.NewReportCodec(feedId, lggr)
	return reportCodec.Decode(report)
}
