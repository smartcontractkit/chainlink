package llo

import (
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/llo/evm"
)

// NOTE: All supported codecs must be specified here
func NewCodecs() map[llotypes.ReportFormat]llo.ReportCodec {
	codecs := make(map[llotypes.ReportFormat]llo.ReportCodec)

	codecs[llotypes.ReportFormatJSON] = llo.JSONReportCodec{}
	codecs[llotypes.ReportFormatEVMPremiumLegacy] = evm.ReportCodecPremiumLegacy{}

	return codecs
}
