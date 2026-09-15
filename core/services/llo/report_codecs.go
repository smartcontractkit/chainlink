package llo

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
	lloreportcodec "github.com/smartcontractkit/chainlink-data-streams/llo/reportcodec"
	"github.com/smartcontractkit/chainlink-data-streams/llo/reportcodec/evm"
)

// NOTE: All supported codecs must be specified here
func NewReportCodecs(lggr logger.Logger, donID uint32) map[llotypes.ReportFormat]lloprotocol.ReportCodec {
	codecs := make(map[llotypes.ReportFormat]lloprotocol.ReportCodec)

	codecs[llotypes.ReportFormatJSON] = lloreportcodec.JSONReportCodec{}
	codecs[llotypes.ReportFormatEVMPremiumLegacy] = evm.NewReportCodecPremiumLegacy(lggr, donID)
	codecs[llotypes.ReportFormatEVMABIEncodeUnpacked] = evm.NewReportCodecEVMABIEncodeUnpacked(lggr, donID)
	codecs[llotypes.ReportFormatEVMABIEncodeUnpackedExpr] = evm.NewReportCodecEVMABIEncodeUnpackedExpr(lggr, donID)
	codecs[llotypes.ReportFormatEVMStreamlined] = evm.NewReportCodecStreamlined(lggr)
	codecs[llotypes.ReportFormatHistoryBackfill] = lloprotocol.ReportCodecHistoryBackfill{}

	return codecs
}
