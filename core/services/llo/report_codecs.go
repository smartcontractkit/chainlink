package llo

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink-data-streams/llo/reportcodecs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/cre"
)

// NOTE: All supported codecs must be specified here
func NewReportCodecs(lggr logger.Logger, donID uint32) map[llotypes.ReportFormat]llo.ReportCodec {
	codecs := make(map[llotypes.ReportFormat]llo.ReportCodec)

	codecs[llotypes.ReportFormatJSON] = llo.JSONReportCodec{}
	codecs[llotypes.ReportFormatEVMPremiumLegacy] = evm.NewReportCodecPremiumLegacy(lggr, donID)
	codecs[llotypes.ReportFormatEVMABIEncodeUnpacked] = evm.NewReportCodecEVMABIEncodeUnpacked(lggr, donID)
	codecs[llotypes.ReportFormatCapabilityTrigger] = cre.NewReportCodecCapabilityTrigger(lggr, donID)
	codecs[llotypes.ReportFormatEVMStreamlined] = evm.NewReportCodecStreamlined()

	return codecs
}
