package types

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
)

// MedianProvider provides all components needed for a median OCR2 plugin.
type MedianProvider interface {
	PluginProvider
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
	OnchainConfigCodec() median.OnchainConfigCodec
}
