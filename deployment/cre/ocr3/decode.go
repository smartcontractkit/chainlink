package ocr3

import (
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"

	capocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	evmcapocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/consensus/ocr3/types"
	dontimepb "github.com/smartcontractkit/chainlink-common/pkg/workflows/dontime/pb"
)

// DecodeCapRegOCR3Configs decodes the OCR3 offchain configs from a serialized
// CapabilityConfig protobuf (as stored in the CapabilitiesRegistry contract).
// Returns a map from OCR3 instance name (e.g. "__default__") to a decoded OracleConfig
// containing human-readable timing parameters and plugin-specific offchain config.
func DecodeCapRegOCR3Configs(rawCapCfg []byte) (map[string]*OracleConfig, error) {
	if len(rawCapCfg) == 0 {
		return nil, nil
	}

	pbCfg := &capabilitiespb.CapabilityConfig{}
	if err := proto.Unmarshal(rawCapCfg, pbCfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CapabilityConfig proto: %w", err)
	}

	if len(pbCfg.Ocr3Configs) == 0 {
		return nil, nil
	}

	result := make(map[string]*OracleConfig, len(pbCfg.Ocr3Configs))
	for instanceName, ocr3Cfg := range pbCfg.Ocr3Configs {
		if ocr3Cfg == nil || len(ocr3Cfg.OffchainConfig) == 0 {
			continue
		}
		decoded, err := decodeCapRegOCR3Config(ocr3Cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to decode OCR3 config for instance %q: %w", instanceName, err)
		}
		result[instanceName] = decoded
	}

	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

// decodeCapRegOCR3Config decodes a single OCR3Config proto into an OracleConfig.
// The offchain_config bytes are libocr-encoded; PublicConfigFromContractConfig is used
// to extract timing parameters and the embedded reporting plugin config.
func decodeCapRegOCR3Config(pbOCR3Cfg *capabilitiespb.OCR3Config) (*OracleConfig, error) {
	signers := make([]ocr2types.OnchainPublicKey, len(pbOCR3Cfg.Signers))
	for i, s := range pbOCR3Cfg.Signers {
		signers[i] = s
	}

	transmitters := make([]ocr2types.Account, len(pbOCR3Cfg.Transmitters))
	for i, t := range pbOCR3Cfg.Transmitters {
		transmitters[i] = ocr2types.Account(common.BytesToAddress(t).String())
	}

	var zeroDigest ocr2types.ConfigDigest
	publicConfig, err := ocr3confighelper.PublicConfigFromContractConfig(true, ocr2types.ContractConfig{
		ConfigDigest:          zeroDigest,
		ConfigCount:           pbOCR3Cfg.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     uint8(pbOCR3Cfg.F), //nolint:gosec // F is bounded by the DON node count
		OnchainConfig:         pbOCR3Cfg.OnchainConfig,
		OffchainConfigVersion: pbOCR3Cfg.OffchainConfigVersion,
		OffchainConfig:        pbOCR3Cfg.OffchainConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to decode libocr offchain config: %w", err)
	}

	pluginCfg, err := decodeCapRegReportingPluginConfig(publicConfig.ReportingPluginConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode reporting plugin config: %w", err)
	}

	oracleConfig := buildCapRegOracleConfig(publicConfig, pluginCfg)
	return &oracleConfig, nil
}

// capRegPluginResult holds the decoded plugin-specific sub-config extracted from
// the OCR3 ReportingPluginConfig bytes stored in the CapabilitiesRegistry.
type capRegPluginResult struct {
	consensus *ConsensusCapOffchainConfig
	dontime   *DontimeOffchainConfig
	chainCap  *ChainCapOffchainConfig
}

// decodeCapRegReportingPluginConfig auto-detects the plugin type from binary
// ReportingPluginConfig bytes using the same heuristics as the OCR3 contract view decoder.
//
// Detection heuristics (applied in order):
//  1. Dontime:   executionRemovalTime (field 8) decodes to a duration >= 1 minute.
//  2. Consensus: outcomePruningThreshold (field 7) or requestTimeout (field 8) is non-zero.
//  3. Chain-cap: maxQueryLengthBytes (field 1) is non-zero but no fields 7–9 are set.
func decodeCapRegReportingPluginConfig(data []byte) (capRegPluginResult, error) {
	if len(data) == 0 {
		return capRegPluginResult{}, nil
	}

	var dCfg dontimepb.Config
	if err := proto.Unmarshal(data, &dCfg); err == nil &&
		dCfg.ExecutionRemovalTime != nil &&
		dCfg.ExecutionRemovalTime.AsDuration() >= time.Minute {
		var execRemoval time.Duration
		if dCfg.ExecutionRemovalTime != nil {
			execRemoval = dCfg.ExecutionRemovalTime.AsDuration()
		}
		return capRegPluginResult{dontime: &DontimeOffchainConfig{
			MaxQueryLengthBytes:       dCfg.MaxQueryLengthBytes,
			MaxObservationLengthBytes: dCfg.MaxObservationLengthBytes,
			MaxOutcomeLengthBytes:     dCfg.MaxOutcomeLengthBytes,
			MaxReportLengthBytes:      dCfg.MaxReportLengthBytes,
			MaxReportCount:            dCfg.MaxReportCount,
			MaxBatchSize:              dCfg.MaxBatchSize,
			MinTimeIncrease:           dCfg.MinTimeIncrease,
			ExecutionRemovalTime:      execRemoval,
		}}, nil
	}

	var cCfg capocr3types.ReportingPluginConfig
	if err := proto.Unmarshal(data, &cCfg); err == nil &&
		(cCfg.OutcomePruningThreshold != 0 || cCfg.RequestTimeout != nil) {
		var reqTimeout time.Duration
		if cCfg.RequestTimeout != nil {
			reqTimeout = cCfg.RequestTimeout.AsDuration()
		}
		return capRegPluginResult{consensus: &ConsensusCapOffchainConfig{
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
	}

	var eCfg evmcapocr3types.ReportingPluginConfig
	if err := proto.Unmarshal(data, &eCfg); err == nil && eCfg.MaxQueryLengthBytes != 0 {
		return capRegPluginResult{chainCap: &ChainCapOffchainConfig{
			MaxQueryLengthBytes:       eCfg.MaxQueryLengthBytes,
			MaxObservationLengthBytes: eCfg.MaxObservationLengthBytes,
			MaxReportLengthBytes:      eCfg.MaxReportLengthBytes,
			MaxOutcomeLengthBytes:     eCfg.MaxOutcomeLengthBytes,
			MaxReportCount:            eCfg.MaxReportCount,
			MaxBatchSize:              eCfg.MaxBatchSize,
		}}, nil
	}

	return capRegPluginResult{}, nil
}

// buildCapRegOracleConfig assembles an OracleConfig from libocr PublicConfig timing fields
// and the decoded plugin sub-config.
func buildCapRegOracleConfig(publicConfig ocr3confighelper.PublicConfig, pluginCfg capRegPluginResult) OracleConfig {
	return OracleConfig{
		ConsensusCapOffchainConfig: pluginCfg.consensus,
		DontimeOffchainConfig:      pluginCfg.dontime,
		ChainCapOffchainConfig:     pluginCfg.chainCap,
		UniqueReports:              true,

		DeltaProgressMillis:               msToUint32(publicConfig.DeltaProgress),
		DeltaResendMillis:                 msToUint32(publicConfig.DeltaResend),
		DeltaInitialMillis:                msToUint32(publicConfig.DeltaInitial),
		DeltaRoundMillis:                  msToUint32(publicConfig.DeltaRound),
		DeltaGraceMillis:                  msToUint32(publicConfig.DeltaGrace),
		DeltaCertifiedCommitRequestMillis: msToUint32(publicConfig.DeltaCertifiedCommitRequest),
		DeltaStageMillis:                  msToUint32(publicConfig.DeltaStage),
		MaxRoundsPerEpoch:                 publicConfig.RMax,
		TransmissionSchedule:              publicConfig.S,

		MaxDurationQueryMillis:          msToUint32(publicConfig.MaxDurationQuery),
		MaxDurationObservationMillis:    msToUint32(publicConfig.MaxDurationObservation),
		MaxDurationShouldAcceptMillis:   msToUint32(publicConfig.MaxDurationShouldAcceptAttestedReport),
		MaxDurationShouldTransmitMillis: msToUint32(publicConfig.MaxDurationShouldTransmitAcceptedReport),

		MaxFaultyOracles: publicConfig.F,
	}
}

func msToUint32(dur time.Duration) uint32 {
	ms := dur.Milliseconds()
	if ms > int64(math.MaxUint32) {
		return math.MaxUint32
	}
	//nolint:gosec // disable G115 as it is practically impossible to overflow here
	return uint32(ms)
}
