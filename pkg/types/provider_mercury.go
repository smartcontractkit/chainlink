package types

import (
	"github.com/smartcontractkit/chainlink-common/pkg/reportingplugins/mercury"
	v1 "github.com/smartcontractkit/chainlink-common/pkg/reportingplugins/mercury/v1"
	v2 "github.com/smartcontractkit/chainlink-common/pkg/reportingplugins/mercury/v2"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/reportingplugins/mercury/v3"
)

// MercuryProvider provides components needed for a mercury OCR2 plugin.
// Mercury requires config tracking but does not transmit on-chain.
type MercuryProvider interface {
	PluginProvider

	ReportCodecV1() v1.ReportCodec
	ReportCodecV2() v2.ReportCodec
	ReportCodecV3() v3.ReportCodec
	OnchainConfigCodec() mercury.OnchainConfigCodec
	MercuryServerFetcher() mercury.MercuryServerFetcher
	ChainReader() mercury.ChainReader
}
