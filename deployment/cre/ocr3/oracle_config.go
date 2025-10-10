package ocr3

import (
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"gopkg.in/yaml.v3"

	capocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	evmcapocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/consensus/ocr3/types"
)

type OracleConfig struct {
	UniqueReports                     bool
	DeltaProgressMillis               uint32
	DeltaResendMillis                 uint32
	DeltaInitialMillis                uint32
	DeltaRoundMillis                  uint32
	DeltaGraceMillis                  uint32
	DeltaCertifiedCommitRequestMillis uint32
	DeltaStageMillis                  uint32
	MaxRoundsPerEpoch                 uint64
	TransmissionSchedule              []int

	MaxDurationQueryMillis          uint32
	MaxDurationObservationMillis    uint32
	MaxDurationShouldAcceptMillis   uint32
	MaxDurationShouldTransmitMillis uint32

	MaxFaultyOracles int

	ConsensusCapOffchainConfig *ConsensusCapOffchainConfig
	ChainCapOffchainConfig     *ChainCapOffchainConfig
}

func (oc *OracleConfig) UnmarshalJSON(data []byte) error {
	// ensure that caller migrated to new OracleConfig structure, where ConsensusCapOffchainConfig is not embedded
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal OracleConfig into map[string]interface{}: %w", err)
	}

	var legacyOffchainConfigFields = []string{"MaxQueryLengthBytes", "MaxObservationLengthBytes", "MaxReportLengthBytes", "MaxOutcomeLengthBytes", "MaxReportCount", "MaxBatchSize", "OutcomePruningThreshold", "RequestTimeout"}
	err := ensureNoLegacyFields(legacyOffchainConfigFields, raw)
	if err != nil {
		return err
	}

	type aliasT OracleConfig
	err = json.Unmarshal(data, (*aliasT)(oc))
	return err
}

func ensureNoLegacyFields(legacyFields []string, raw map[string]any) error {
	for _, f := range legacyFields {
		if _, exists := raw[f]; exists {
			return fmt.Errorf("not supported config format detected: field %s is not supported. All %v must be moved into ConsensusCapOffchainConfig", f, legacyFields)
		}
	}

	return nil
}

func (oc *OracleConfig) UnmarshalYAML(value *yaml.Node) error {
	// ensure that caller migrated to new OracleConfig structure, where ConsensusCapOffchainConfig is not embedded
	var raw map[string]any
	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("failed to decode OracleConfig into map[string]interface{}: %w", err)
	}

	var legacyOffchainConfigFields = []string{"maxQueryLengthBytes", "maxObservationLengthBytes", "maxReportLengthBytes", "maxOutcomeLengthBytes", "maxReportCount", "maxBatchSize", "outcomePruningThreshold", "requestTimeout"}
	err := ensureNoLegacyFields(legacyOffchainConfigFields, raw)
	if err != nil {
		return err
	}

	type aliasT OracleConfig
	return value.Decode((*aliasT)(oc))
}

type offchainConfig interface {
	ToProto() (proto.Message, error)
}

type ConsensusCapOffchainConfig struct {
	MaxQueryLengthBytes       uint32
	MaxObservationLengthBytes uint32
	MaxReportLengthBytes      uint32
	MaxOutcomeLengthBytes     uint32
	MaxReportCount            uint32
	MaxBatchSize              uint32
	OutcomePruningThreshold   uint64
	RequestTimeout            time.Duration
}

func (oc *ConsensusCapOffchainConfig) UnmarshalJSON(data []byte) error {
	type aliasT ConsensusCapOffchainConfig
	temp := &struct {
		RequestTimeout string `json:"RequestTimeout"`
		*aliasT
	}{
		aliasT: (*aliasT)(oc),
	}
	if err := json.Unmarshal(data, temp); err != nil {
		return fmt.Errorf("failed to unmarshal OracleConfig: %w", err)
	}

	if temp.RequestTimeout == "" {
		oc.RequestTimeout = 0
	} else {
		requestTimeout, err := time.ParseDuration(temp.RequestTimeout)
		if err != nil {
			return fmt.Errorf("failed to parse RequestTimeout: %w", err)
		}
		oc.RequestTimeout = requestTimeout
	}

	return nil
}

func (oc *ConsensusCapOffchainConfig) MarshalJSON() ([]byte, error) {
	type aliasT ConsensusCapOffchainConfig
	return json.Marshal(&struct {
		RequestTimeout string `json:"RequestTimeout"`
		*aliasT
	}{
		RequestTimeout: oc.RequestTimeout.String(),
		aliasT:         (*aliasT)(oc),
	})
}

func (oc *ConsensusCapOffchainConfig) ToProto() (proto.Message, error) {
	// let's keep reqTimeout as nil if it's 0, so we can use the default value within `chainlink-common`.
	// See: https://github.com/smartcontractkit/chainlink-common/blob/main/pkg/capabilities/consensus/ocr3/factory.go#L73
	var reqTimeout *durationpb.Duration
	if oc.RequestTimeout > 0 {
		reqTimeout = durationpb.New(oc.RequestTimeout)
	}
	return &capocr3types.ReportingPluginConfig{
		MaxQueryLengthBytes:       oc.MaxQueryLengthBytes,
		MaxObservationLengthBytes: oc.MaxObservationLengthBytes,
		MaxReportLengthBytes:      oc.MaxReportLengthBytes,
		MaxOutcomeLengthBytes:     oc.MaxOutcomeLengthBytes,
		MaxReportCount:            oc.MaxReportCount,
		MaxBatchSize:              oc.MaxBatchSize,
		OutcomePruningThreshold:   oc.OutcomePruningThreshold,
		RequestTimeout:            reqTimeout,
	}, nil
}

type ChainCapOffchainConfig struct {
	MaxQueryLengthBytes       uint32
	MaxObservationLengthBytes uint32
	MaxReportLengthBytes      uint32
	MaxOutcomeLengthBytes     uint32
	MaxReportCount            uint32
	MaxBatchSize              uint32
}

func (oc *ChainCapOffchainConfig) ToProto() (proto.Message, error) {
	return &evmcapocr3types.ReportingPluginConfig{
		MaxQueryLengthBytes:       oc.MaxQueryLengthBytes,
		MaxObservationLengthBytes: oc.MaxObservationLengthBytes,
		MaxReportLengthBytes:      oc.MaxReportLengthBytes,
		MaxOutcomeLengthBytes:     oc.MaxOutcomeLengthBytes,
		MaxReportCount:            oc.MaxReportCount,
		MaxBatchSize:              oc.MaxBatchSize,
	}, nil
}
