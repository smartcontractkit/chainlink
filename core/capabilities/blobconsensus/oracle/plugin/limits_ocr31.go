package plugin

import (
	ocrtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

func reportingPluginLimitsOCR31(configProto *ocrtypes.ReportingPluginConfig) ocr3_1types.ReportingPluginLimits {
	maxPrec := int(configProto.MaxOutcomeLengthBytes)
	kvPlus := maxPrec + len(prevOutcomeStateKey) + 256
	if kvPlus > ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes {
		kvPlus = ocr3_1types.MaxMaxKeyValueModifiedKeysPlusValuesBytes
	}
	return ocr3_1types.ReportingPluginLimits{
		MaxQueryBytes:                int(configProto.MaxQueryLengthBytes),
		MaxObservationBytes:          int(configProto.MaxObservationLengthBytes),
		MaxReportsPlusPrecursorBytes: maxPrec,
		MaxReportBytes:               int(configProto.MaxReportLengthBytes),
		MaxReportCount:               int(configProto.MaxReportCount),
		MaxKeyValueModifiedKeys:      2,
		MaxKeyValueModifiedKeysPlusValuesBytes:          kvPlus,
		MaxBlobPayloadBytes:                             ocr3_1types.MaxMaxBlobPayloadBytes,
		MaxPerOracleUnexpiredBlobCumulativePayloadBytes: 8 * ocr3_1types.MaxMaxObservationBytes,
		MaxPerOracleUnexpiredBlobCount:                  256,
	}
}

func reportingPluginInfoOCR3From31(info ocr3_1types.ReportingPluginInfo1) ocr3types.ReportingPluginInfo {
	l := info.Limits
	return ocr3types.ReportingPluginInfo{
		Name: info.Name,
		Limits: ocr3types.ReportingPluginLimits{
			MaxQueryLength:       l.MaxQueryBytes,
			MaxObservationLength: l.MaxObservationBytes,
			MaxOutcomeLength:     l.MaxReportsPlusPrecursorBytes,
			MaxReportLength:      l.MaxReportBytes,
			MaxReportCount:       l.MaxReportCount,
		},
	}
}
