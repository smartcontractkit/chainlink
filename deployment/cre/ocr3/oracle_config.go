package ocr3

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	capocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	evmcapocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/consensus/ocr3/types"
)

type OracleConfig struct {
	// Excluded from JSON to maintain backward compatibility with previous versions where ReportingPluginConfig was embedded
	OffchainConfig     OffchainConfig  `json:"-"`
	RawOffchainConfig  json.RawMessage `json:"OffchainConfig"`
	OffchainConfigType OffchainConfigType

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
}

func (oc *OracleConfig) UnmarshalJSON(data []byte) error {
	type aliasT OracleConfig
	err := json.Unmarshal(data, (*aliasT)(oc))
	if err != nil {
		return fmt.Errorf("failed to unmarshal OracleConfig: %w", err)
	}

	switch oc.OffchainConfigType {
	case "", OffchainConfigTypeConsensusCap:
		oc.OffchainConfig = &ConsensusCapOffchainConfig{}
	case OffchainConfigTypeChainCap:
		oc.OffchainConfig = &ChainCapOffchainConfig{}
	default:
		return fmt.Errorf("unsupported OffchainConfigType: %s", oc.OffchainConfigType)
	}

	// if offchain_config is empty, try to use previous version, where OffchainConfig was embedded
	rawOffchainConfig := oc.RawOffchainConfig
	if len(rawOffchainConfig) == 0 {
		rawOffchainConfig = data
	}
	// try to use previous version, where OffchainConfig was embedded
	err = json.Unmarshal(rawOffchainConfig, &oc.OffchainConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal OffchainConfig: %w", err)
	}

	oc.RawOffchainConfig = nil // clear raw data after successful unmarshalling
	return nil
}

func (oc OracleConfig) MarshalJSON() ([]byte, error) {
	// ensure that caller did not forget to set OffchainConfigType
	if oc.OffchainConfigType == "" && oc.OffchainConfig != nil {
		_, ok := oc.OffchainConfig.(*ConsensusCapOffchainConfig)
		if !ok {
			return nil, errors.New("OffchainConfigType must be set when OffchainConfig is set")
		}
	}

	offchainConfigAsJSON, err := json.Marshal(oc.OffchainConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OffchainConfig: %w", err)
	}

	cfgToMarshal := oc
	cfgToMarshal.RawOffchainConfig = offchainConfigAsJSON

	type aliasT OracleConfig
	asJSON, err := json.Marshal((aliasT)(cfgToMarshal))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OracleConfig: %w", err)
	}

	return asJSON, nil
}

type OffchainConfig interface {
	ToProto() (proto.Message, error)
	isOffchainConfig()
}

type OffchainConfigType string

func (t OffchainConfigType) String() string {
	return string(t)
}

const (
	OffchainConfigTypeConsensusCap OffchainConfigType = "consensus-cap"
	OffchainConfigTypeChainCap     OffchainConfigType = "chain-cap"
)

func ChainCapChainSelectorLabel(chainSelector uint64) string {
	return fmt.Sprintf("chain-selector-%d", chainSelector)
}

func NewOffchainConfigFromProto(cfgType OffchainConfigType, raw []byte) (OffchainConfig, error) {
	switch cfgType {
	case OffchainConfigTypeConsensusCap:
		cfg := &capocr3types.ReportingPluginConfig{}
		err := proto.Unmarshal(raw, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal ConsensusCap OffchainConfig from proto: %w", err)
		}

		return &ConsensusCapOffchainConfig{
			MaxQueryLengthBytes:       cfg.MaxQueryLengthBytes,
			MaxObservationLengthBytes: cfg.MaxObservationLengthBytes,
			MaxReportLengthBytes:      cfg.MaxReportLengthBytes,
			MaxOutcomeLengthBytes:     cfg.MaxOutcomeLengthBytes,
			MaxReportCount:            cfg.MaxReportCount,
			MaxBatchSize:              cfg.MaxBatchSize,
			OutcomePruningThreshold:   cfg.OutcomePruningThreshold,
			RequestTimeout:            cfg.RequestTimeout.AsDuration(),
		}, nil
	case OffchainConfigTypeChainCap:
		cfg := &evmcapocr3types.ReportingPluginConfig{}
		err := proto.Unmarshal(raw, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal ChainCap OffchainConfig from proto: %w", err)
		}
		return &ChainCapOffchainConfig{
			MaxQueryLengthBytes:       cfg.MaxQueryLengthBytes,
			MaxObservationLengthBytes: cfg.MaxObservationLengthBytes,
			MaxReportLengthBytes:      cfg.MaxReportLengthBytes,
			MaxOutcomeLengthBytes:     cfg.MaxOutcomeLengthBytes,
			MaxReportCount:            cfg.MaxReportCount,
			MaxBatchSize:              cfg.MaxBatchSize,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported OffchainConfigType: %s", cfgType)
	}
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

func (*ConsensusCapOffchainConfig) isOffchainConfig() {}

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

func (*ChainCapOffchainConfig) isOffchainConfig() {}
