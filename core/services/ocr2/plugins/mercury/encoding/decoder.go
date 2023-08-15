package encoding

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	mercuryv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1"
	mercuryv2 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v2"
	mercuryv3 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/pkg/errors"
)

func DecodeV1(report ocrtypes.Report, lggr logger.Logger) (*mercuryv1.Report, error) {
	feedID, err := mercury.FeedIDFromReport(report)
	if err != nil {
		return nil, err
	}

	if !feedID.IsV1() {
		return nil, errors.Errorf("invalid schema version: %d", feedID.Version())
	}

	reportCodec := mercuryv1.NewReportCodec(feedID, lggr)
	return reportCodec.Decode(report)
}

func DecodeV2(report ocrtypes.Report, lggr logger.Logger) (*mercuryv2.Report, error) {
	feedID, err := mercury.FeedIDFromReport(report)
	if err != nil {
		return nil, err
	}

	if !feedID.IsV2() {
		return nil, errors.Errorf("invalid schema version: %d", feedID.Version())
	}

	reportCodec := mercuryv2.NewReportCodec(feedID, lggr)
	return reportCodec.Decode(report)
}

func DecodeV3(report ocrtypes.Report, lggr logger.Logger) (*mercuryv3.Report, error) {
	feedID, err := mercury.FeedIDFromReport(report)
	if err != nil {
		return nil, err
	}

	if !feedID.IsV3() {
		return nil, errors.Errorf("invalid schema version: %d", feedID.Version())
	}

	reportCodec := mercuryv3.NewReportCodec(feedID, lggr)
	return reportCodec.Decode(report)
}
