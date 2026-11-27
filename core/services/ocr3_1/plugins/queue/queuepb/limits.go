package queuepb

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3_1/plugins/queue/queuelimits"
)

// Config floors. A value below one of these does not merely degrade the plugin: it cannot observe events or commit a
// state transition, permanently and DON-wide.
const (
	// MinEventPayloadBytes is a conservative floor for one encoded trigger event: a 64-character workflow ID plus the
	// trigger response. No event can be observed or broadcast below it.
	MinEventPayloadBytes = 256
	// MinObservationBytes is the least observation budget beyond the blob handle reservation: one event's metadata,
	// the expiration and capacity fields, and the workflow capacities map.
	MinObservationBytes = 1024
	// MinKeyValueModifiedKeysPlusValuesBytes is the least budget that can hold every KV key with an empty value. A
	// smaller budget makes libocr reject every state transition, so no round ever commits.
	MinKeyValueModifiedKeysPlusValuesBytes = len(queuelimits.ObservedEventsKey) + len(queuelimits.QueuedEventsKey) +
		len(queuelimits.ConsensusEventsKey) + len(queuelimits.BlobHandlesKey)
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
// Returns an error if the limits are missing, if they cannot support the plugin's writes and observations (see
// ValidateConfig), or if any limit exceeds libocr's maximum.
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

	if err := errors.Join(ValidateConfig(out), out.Validate()); err != nil {
		return ocr3_1types.ReportingPluginLimits{}, fmt.Errorf("invalid consensus-queue limits: %w", err)
	}
	return out, nil
}

// ValidateConfig checks the configuration-only floors the plugin relies on. A value below one of these does not
// merely degrade: the plugin cannot observe events or commit a state transition, permanently and DON-wide.
func ValidateConfig(l ocr3_1types.ReportingPluginLimits) error {
	var err error
	if l.MaxObservationBytes < MinObservationBytes {
		err = errors.Join(err, fmt.Errorf(
			"MaxObservationBytes must be at least %d to carry one event, but got %d", MinObservationBytes, l.MaxObservationBytes))
	}
	if l.MaxKeyValueModifiedKeys < queuelimits.NumKVKeys {
		err = errors.Join(err, fmt.Errorf(
			"MaxKeyValueModifiedKeys must be at least %d, the keys written every round, but got %d", queuelimits.NumKVKeys, l.MaxKeyValueModifiedKeys))
	}
	if l.MaxKeyValueModifiedKeysPlusValuesBytes < MinKeyValueModifiedKeysPlusValuesBytes {
		err = errors.Join(err, fmt.Errorf(
			"MaxKeyValueModifiedKeysPlusValuesBytes must be at least %d, the keys written every round, but got %d",
			MinKeyValueModifiedKeysPlusValuesBytes, l.MaxKeyValueModifiedKeysPlusValuesBytes))
	}
	minBlob := queuelimits.MaxEventFramingOverhead(queuelimits.CoordinatedEventsFieldNumber) + MinEventPayloadBytes
	if l.MaxBlobPayloadBytes < minBlob {
		err = errors.Join(err, fmt.Errorf(
			"MaxBlobPayloadBytes must be at least %d to carry one event, but got %d", minBlob, l.MaxBlobPayloadBytes))
	}
	if l.MaxPerOracleUnexpiredBlobCount <= 0 {
		err = errors.Join(err, fmt.Errorf(
			"MaxPerOracleUnexpiredBlobCount must be positive, but got %d", l.MaxPerOracleUnexpiredBlobCount))
	}
	if l.MaxPerOracleUnexpiredBlobCumulativePayloadBytes <= 0 {
		err = errors.Join(err, fmt.Errorf(
			"MaxPerOracleUnexpiredBlobCumulativePayloadBytes must be positive, but got %d", l.MaxPerOracleUnexpiredBlobCumulativePayloadBytes))
	}
	return err
}

// ValidateForDon checks the limits that depend on the DON's size. Observation reserves space for
// queuelimits.MaxBlobHandlesPerObservation blob handles up front, and their wire size grows with N and F, so a config
// that is valid in general can still be too small for a particular DON.
func ValidateForDon(l ocr3_1types.ReportingPluginLimits, n, f int) error {
	handleBytes := ocr3_1types.BlobHandleMarshalledBytesUpperBound(n, f)
	minBytes := queuelimits.MaxBlobHandlesPerObservation*handleBytes + MinObservationBytes
	if l.MaxObservationBytes < minBytes {
		return fmt.Errorf(
			"MaxObservationBytes must be at least %d for N=%d F=%d, which reserves space for %d blob handles, but got %d",
			minBytes, n, f, queuelimits.MaxBlobHandlesPerObservation, l.MaxObservationBytes)
	}
	return nil
}
