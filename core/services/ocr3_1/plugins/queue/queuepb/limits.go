package queuepb

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
)

// Decode parses the consensus-queue plugin config from the OCR offchain config bytes.
// The node and the deployment tooling both call it, so a config the tooling writes is one the node accepts.
//
// Returns an error if b is empty, cannot be unmarshaled, or holds invalid limits (see OCRLimits).
func Decode(b []byte) (*ReportingPluginConfig, error) {
	if len(b) == 0 {
		// proto.Unmarshal accepts empty input and yields zero limits, so it has to be rejected here.
		return nil, errors.New("empty consensus-queue reporting plugin config")
	}

	cfg := &ReportingPluginConfig{}
	if err := proto.Unmarshal(b, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal consensus-queue reporting plugin config: %w", err)
	}
	if _, err := cfg.OCRLimits(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// OCRLimits converts the config's limits to libocr's ReportingPluginLimits.
// There are no defaults. Every value comes from the config so all oracles agree on it.
//
// Returns an error if the limits are missing, if MaxObservationBytes, MaxKeyValueModifiedKeys or
// MaxKeyValueModifiedKeysPlusValuesBytes is zero, or if any limit exceeds libocr's maximum.
func (c *ReportingPluginConfig) OCRLimits() (ocr3_1types.ReportingPluginLimits, error) {
	l := c.GetLimits()
	if l == nil {
		return ocr3_1types.ReportingPluginLimits{}, errors.New("consensus-queue reporting plugin config has no limits")
	}

	out := ocr3_1types.ReportingPluginLimits{
		MaxQueryBytes:                                   int(l.GetMaxQueryBytes()),
		MaxObservationBytes:                             int(l.GetMaxObservationBytes()),
		MaxReportsPlusPrecursorBytes:                    int(l.GetMaxReportsPlusPrecursorBytes()),
		MaxReportBytes:                                  int(l.GetMaxReportBytes()),
		MaxReportCount:                                  int(l.GetMaxReportCount()),
		MaxKeyValueModifiedKeys:                         int(l.GetMaxKeyValueModifiedKeys()),
		MaxKeyValueModifiedKeysPlusValuesBytes:          int(l.GetMaxKeyValueModifiedKeysPlusValuesBytes()),
		MaxBlobPayloadBytes:                             int(l.GetMaxBlobPayloadBytes()),
		MaxPerOracleUnexpiredBlobCumulativePayloadBytes: int(l.GetMaxPerOracleUnexpiredBlobCumulativePayloadBytes()),
		MaxPerOracleUnexpiredBlobCount:                  int(l.GetMaxPerOracleUnexpiredBlobCount()),
	}

	var err error
	if out.MaxObservationBytes == 0 {
		err = errors.Join(err, errors.New("MaxObservationBytes must be at least 1"))
	}
	if out.MaxKeyValueModifiedKeys == 0 {
		err = errors.Join(err, errors.New("MaxKeyValueModifiedKeys must be at least 1"))
	}
	if out.MaxKeyValueModifiedKeysPlusValuesBytes == 0 {
		err = errors.Join(err, errors.New("MaxKeyValueModifiedKeysPlusValuesBytes must be at least 1"))
	}
	if err = errors.Join(err, out.Validate()); err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("invalid consensus-queue limits: %w", err)
	}
	return out, nil
}
