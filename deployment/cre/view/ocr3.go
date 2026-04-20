package view

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"

	capocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	evmcapocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/consensus/ocr3/types"
	dontimepb "github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime/pb"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

var ErrOCR3NotConfigured = errors.New("OCR3 not configured")

type OCR3ConfigView struct {
	Signers               []string            `json:"signers" yaml:"signers"`
	Transmitters          []ocr2types.Account `json:"transmitters" yaml:"transmitters"`
	F                     uint8               `json:"f" yaml:"f"`
	OnchainConfig         []byte              `json:"onchainConfig,omitempty" yaml:"onchainConfig,omitempty"`
	OffchainConfigVersion uint64              `json:"offchainConfigVersion" yaml:"offchainConfigVersion"`
	OffchainConfig        ocr3.OracleConfig   `json:"offchainConfig" yaml:"offchainConfig"`
}

type OCR3ConfigViewLegacy struct {
	Signers               []string            `json:"signers"`
	Transmitters          []ocr2types.Account `json:"transmitters"`
	F                     uint8               `json:"f"`
	OnchainConfig         []byte              `json:"onchainConfig"`
	OffchainConfigVersion uint64              `json:"offchainConfigVersion"`
	OffchainConfig        OracleConfigLegacy  `json:"offchainConfig"`
}

type OracleConfigLegacy struct {
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

	MaxQueryLengthBytes       uint32
	MaxObservationLengthBytes uint32
	MaxReportLengthBytes      uint32
	MaxOutcomeLengthBytes     uint32
	MaxReportCount            uint32
	MaxBatchSize              uint32
	OutcomePruningThreshold   uint64
	RequestTimeout            time.Duration
}

func (oc *OracleConfigLegacy) UnmarshalJSON(data []byte) error {
	type aliasT OracleConfigLegacy
	temp := &struct {
		RequestTimeout string `json:"RequestTimeout"`
		*aliasT
	}{
		aliasT: (*aliasT)(oc),
	}
	if err := json.Unmarshal(data, temp); err != nil {
		return fmt.Errorf("failed to unmarshal OracleConfigLegacy: %w", err)
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

func GenerateOCR3ConfigView(ctx context.Context, ocr3Cap ocr3_capability.OCR3Capability) (OCR3ConfigView, error) {
	details, err := ocr3Cap.LatestConfigDetails(nil)
	if err != nil {
		return OCR3ConfigView{}, err
	}

	blockNumber := uint64(details.BlockNumber)
	configIterator, err := ocr3Cap.FilterConfigSet(&bind.FilterOpts{
		Start:   blockNumber,
		End:     &blockNumber,
		Context: ctx,
	})
	if err != nil {
		return OCR3ConfigView{}, err
	}
	var config *ocr3_capability.OCR3CapabilityConfigSet
	for configIterator.Next() {
		// We wait for the iterator to receive an event
		if configIterator.Event == nil {
			return OCR3ConfigView{}, ErrOCR3NotConfigured
		}

		config = configIterator.Event
	}
	if config == nil {
		return OCR3ConfigView{}, ErrOCR3NotConfigured
	}

	var signers []ocr2types.OnchainPublicKey
	var readableSigners []string
	for _, s := range config.Signers {
		signers = append(signers, s)
		readableSigners = append(readableSigners, hex.EncodeToString(s))
	}
	var transmitters []ocr2types.Account
	for _, t := range config.Transmitters {
		transmitters = append(transmitters, ocr2types.Account(t.String()))
	}
	// `PublicConfigFromContractConfig` returns the `ocr2types.PublicConfig` that contains all the `OracleConfig` fields we need, including the
	// report plugin config.
	publicConfig, err := ocr3confighelper.PublicConfigFromContractConfig(true, ocr2types.ContractConfig{
		ConfigDigest:          config.ConfigDigest,
		ConfigCount:           config.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     config.F,
		OnchainConfig:         nil, // empty onChain config, currently we always use a nil onchain config when calling SetConfig
		OffchainConfigVersion: config.OffchainConfigVersion,
		OffchainConfig:        config.OffchainConfig,
	})
	if err != nil {
		return OCR3ConfigView{}, err
	}

	pluginCfg, err := decodeReportingPluginConfig(publicConfig.ReportingPluginConfig)
	if err != nil {
		return OCR3ConfigView{}, err
	}

	return OCR3ConfigView{
		Signers:               readableSigners,
		Transmitters:          transmitters,
		F:                     config.F,
		OnchainConfig:         nil, // empty onChain config
		OffchainConfigVersion: config.OffchainConfigVersion,
		OffchainConfig:        buildOracleConfig(publicConfig, pluginCfg),
	}, nil
}

// PluginType identifies the reporting-plugin schema embedded in OCR3 ReportingPluginConfig bytes.
type PluginType string

const (
	PluginTypeConsensus PluginType = "consensus"
	PluginTypeDontime   PluginType = "dontime"
	PluginTypeChainCap  PluginType = "chain-cap"
)

// GenerateOCR3ConfigViewForPlugin is like GenerateOCR3ConfigView but decodes the
// ReportingPluginConfig bytes using the explicitly supplied plugin type rather than
// applying a heuristic. Use this when the plugin type is known (e.g. from the contract
// qualifier in address_refs.json).
func GenerateOCR3ConfigViewForPlugin(ctx context.Context, ocr3Cap ocr3_capability.OCR3Capability, pluginType PluginType) (OCR3ConfigView, error) {
	details, err := ocr3Cap.LatestConfigDetails(nil)
	if err != nil {
		return OCR3ConfigView{}, err
	}

	blockNumber := uint64(details.BlockNumber)
	configIterator, err := ocr3Cap.FilterConfigSet(&bind.FilterOpts{
		Start:   blockNumber,
		End:     &blockNumber,
		Context: ctx,
	})
	if err != nil {
		return OCR3ConfigView{}, err
	}
	var config *ocr3_capability.OCR3CapabilityConfigSet
	for configIterator.Next() {
		if configIterator.Event == nil {
			return OCR3ConfigView{}, ErrOCR3NotConfigured
		}
		config = configIterator.Event
	}
	if config == nil {
		return OCR3ConfigView{}, ErrOCR3NotConfigured
	}

	var signers []ocr2types.OnchainPublicKey
	var readableSigners []string
	for _, s := range config.Signers {
		signers = append(signers, s)
		readableSigners = append(readableSigners, hex.EncodeToString(s))
	}
	var transmitters []ocr2types.Account
	for _, t := range config.Transmitters {
		transmitters = append(transmitters, ocr2types.Account(t.String()))
	}

	publicConfig, err := ocr3confighelper.PublicConfigFromContractConfig(true, ocr2types.ContractConfig{
		ConfigDigest:          config.ConfigDigest,
		ConfigCount:           config.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     config.F,
		OnchainConfig:         nil,
		OffchainConfigVersion: config.OffchainConfigVersion,
		OffchainConfig:        config.OffchainConfig,
	})
	if err != nil {
		return OCR3ConfigView{}, err
	}

	pluginCfg, err := decodeReportingPluginConfigForType(publicConfig.ReportingPluginConfig, pluginType)
	if err != nil {
		return OCR3ConfigView{}, err
	}

	oracleConfig := buildOracleConfig(publicConfig, pluginCfg)
	return OCR3ConfigView{
		Signers:               readableSigners,
		Transmitters:          transmitters,
		F:                     config.F,
		OnchainConfig:         nil,
		OffchainConfigVersion: config.OffchainConfigVersion,
		OffchainConfig:        oracleConfig,
	}, nil
}

// reportingPluginResult holds the decoded plugin-specific sub-config.
// Exactly one field will be non-nil after a successful decode; all may be nil when the
// plugin config is empty or unrecognised (e.g. DKG vault contracts).
type reportingPluginResult struct {
	consensus *ocr3.ConsensusCapOffchainConfig
	dontime   *ocr3.DontimeOffchainConfig
	chainCap  *ocr3.ChainCapOffchainConfig
}

// decodeReportingPluginConfigForType decodes ReportingPluginConfig bytes using the
// supplied plugin type — no heuristic is applied.
func decodeReportingPluginConfigForType(data []byte, pluginType PluginType) (reportingPluginResult, error) {
	if len(data) == 0 {
		return reportingPluginResult{}, nil
	}
	switch pluginType {
	case PluginTypeConsensus:
		var cCfg capocr3types.ReportingPluginConfig
		if err := proto.Unmarshal(data, &cCfg); err != nil {
			return reportingPluginResult{}, fmt.Errorf("unmarshal consensus ReportingPluginConfig: %w", err)
		}
		var reqTimeout time.Duration
		if cCfg.RequestTimeout != nil {
			reqTimeout = cCfg.RequestTimeout.AsDuration()
		}
		return reportingPluginResult{consensus: &ocr3.ConsensusCapOffchainConfig{
			MaxQueryLengthBytes:       cCfg.MaxQueryLengthBytes,
			MaxObservationLengthBytes: cCfg.MaxObservationLengthBytes,
			MaxReportLengthBytes:      cCfg.MaxReportLengthBytes,
			MaxOutcomeLengthBytes:     cCfg.MaxOutcomeLengthBytes,
			MaxReportCount:            cCfg.MaxReportCount,
			// NOTE: MaxBatchSize is not used by the consensus plugin v2 (but still used by v1)
			MaxBatchSize:            cCfg.MaxBatchSize,
			OutcomePruningThreshold: cCfg.OutcomePruningThreshold,
			RequestTimeout:          reqTimeout,
		}}, nil

	case PluginTypeDontime:
		var dCfg dontimepb.Config
		if err := proto.Unmarshal(data, &dCfg); err != nil {
			return reportingPluginResult{}, fmt.Errorf("unmarshal dontime ReportingPluginConfig: %w", err)
		}
		var execRemoval time.Duration
		if dCfg.ExecutionRemovalTime != nil {
			execRemoval = dCfg.ExecutionRemovalTime.AsDuration()
		}
		return reportingPluginResult{dontime: &ocr3.DontimeOffchainConfig{
			MaxQueryLengthBytes:       dCfg.MaxQueryLengthBytes,
			MaxObservationLengthBytes: dCfg.MaxObservationLengthBytes,
			MaxOutcomeLengthBytes:     dCfg.MaxOutcomeLengthBytes,
			MaxReportLengthBytes:      dCfg.MaxReportLengthBytes,
			MaxReportCount:            dCfg.MaxReportCount,
			MaxBatchSize:              dCfg.MaxBatchSize,
			MinTimeIncrease:           dCfg.MinTimeIncrease,
			ExecutionRemovalTime:      execRemoval,
		}}, nil

	case PluginTypeChainCap:
		var eCfg evmcapocr3types.ReportingPluginConfig
		if err := proto.Unmarshal(data, &eCfg); err != nil {
			return reportingPluginResult{}, fmt.Errorf("unmarshal chain-cap ReportingPluginConfig: %w", err)
		}
		return reportingPluginResult{chainCap: &ocr3.ChainCapOffchainConfig{
			MaxQueryLengthBytes:       eCfg.MaxQueryLengthBytes,
			MaxObservationLengthBytes: eCfg.MaxObservationLengthBytes,
			MaxReportLengthBytes:      eCfg.MaxReportLengthBytes,
			MaxOutcomeLengthBytes:     eCfg.MaxOutcomeLengthBytes,
			MaxReportCount:            eCfg.MaxReportCount,
			MaxBatchSize:              eCfg.MaxBatchSize,
		}}, nil

	default:
		return reportingPluginResult{}, fmt.Errorf("unknown plugin type %q: must be one of %q, %q, %q",
			pluginType, PluginTypeConsensus, PluginTypeDontime, PluginTypeChainCap)
	}
}

// decodeReportingPluginConfig auto-detects the plugin type from binary ReportingPluginConfig bytes
// by probing distinctive fields in each proto schema.
//
// Detection heuristics (applied in order):
//
//  1. Dontime:   executionRemovalTime (field 8) decodes to a duration >= 1 minute. Consensus
//     requestTimeout lives at the same field number but is always a short timeout (< 1 minute
//     in all known CRE deployments).
//  2. Consensus: outcomePruningThreshold (field 7) or requestTimeout (field 8) is non-zero.
//  3. Chain-cap: maxQueryLengthBytes (field 1) is non-zero but no fields 7–9 are set.
//
// Prefer GenerateOCR3ConfigViewForPlugin when the plugin type is known.
func decodeReportingPluginConfig(data []byte) (reportingPluginResult, error) {
	if len(data) == 0 {
		return reportingPluginResult{}, nil
	}

	var dCfg dontimepb.Config
	if err := proto.Unmarshal(data, &dCfg); err == nil &&
		dCfg.ExecutionRemovalTime != nil &&
		dCfg.ExecutionRemovalTime.AsDuration() >= time.Minute {
		return decodeReportingPluginConfigForType(data, PluginTypeDontime)
	}

	var cCfg capocr3types.ReportingPluginConfig
	if err := proto.Unmarshal(data, &cCfg); err == nil &&
		(cCfg.OutcomePruningThreshold != 0 || cCfg.RequestTimeout != nil) {
		return decodeReportingPluginConfigForType(data, PluginTypeConsensus)
	}

	var eCfg evmcapocr3types.ReportingPluginConfig
	if err := proto.Unmarshal(data, &eCfg); err == nil && eCfg.MaxQueryLengthBytes != 0 {
		return decodeReportingPluginConfigForType(data, PluginTypeChainCap)
	}

	return reportingPluginResult{}, nil
}

// buildOracleConfig assembles an OracleConfig from the libocr PublicConfig timing fields
// and the decoded plugin sub-config.
func buildOracleConfig(publicConfig ocr3confighelper.PublicConfig, pluginCfg reportingPluginResult) ocr3.OracleConfig {
	return ocr3.OracleConfig{
		ConsensusCapOffchainConfig: pluginCfg.consensus,
		DontimeOffchainConfig:      pluginCfg.dontime,
		ChainCapOffchainConfig:     pluginCfg.chainCap,
		UniqueReports:              true,

		DeltaProgressMillis:               millisecondsToUint32(publicConfig.DeltaProgress),
		DeltaResendMillis:                 millisecondsToUint32(publicConfig.DeltaResend),
		DeltaInitialMillis:                millisecondsToUint32(publicConfig.DeltaInitial),
		DeltaRoundMillis:                  millisecondsToUint32(publicConfig.DeltaRound),
		DeltaGraceMillis:                  millisecondsToUint32(publicConfig.DeltaGrace),
		DeltaCertifiedCommitRequestMillis: millisecondsToUint32(publicConfig.DeltaCertifiedCommitRequest),
		DeltaStageMillis:                  millisecondsToUint32(publicConfig.DeltaStage),
		MaxRoundsPerEpoch:                 publicConfig.RMax,
		TransmissionSchedule:              publicConfig.S,

		MaxDurationQueryMillis:          millisecondsToUint32(publicConfig.MaxDurationQuery),
		MaxDurationObservationMillis:    millisecondsToUint32(publicConfig.MaxDurationObservation),
		MaxDurationShouldAcceptMillis:   millisecondsToUint32(publicConfig.MaxDurationShouldAcceptAttestedReport),
		MaxDurationShouldTransmitMillis: millisecondsToUint32(publicConfig.MaxDurationShouldTransmitAcceptedReport),

		MaxFaultyOracles: publicConfig.F,
	}
}

func millisecondsToUint32(dur time.Duration) uint32 {
	ms := dur.Milliseconds()
	if ms > int64(math.MaxUint32) {
		return math.MaxUint32
	}
	//nolint:gosec // disable G115 as it is practically impossible to overflow here
	return uint32(ms)
}
