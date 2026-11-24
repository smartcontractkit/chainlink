package ocr3_1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3_1/plugins/queue/queuelimits"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3_1/plugins/queue/queuepb"
)

type pluginConfig interface {
	Marshal() ([]byte, error)
}

type DKGOffchainConfig struct {
	DealerPublicKeys    []string `yaml:"dealerPublicKeys" json:"dealerPublicKeys"`
	RecipientPublicKeys []string `yaml:"recipientPublicKeys" json:"recipientPublicKeys"`
	T                   int      `yaml:"t" json:"t"`
	PreviousInstanceID  *string  `yaml:"previousInstanceID" json:"previousInstanceID"`
}

func hexArrToParticipantArr(stringKeys []string) ([]dkgocrtypes.P256ParticipantPublicKey, error) {
	keys := []dkgocrtypes.P256ParticipantPublicKey{}
	for _, dk := range stringKeys {
		k, err := hex.DecodeString(dk)
		if err != nil {
			return nil, err
		}

		keys = append(keys, k)
	}

	return keys, nil
}

func (c *DKGOffchainConfig) Marshal() ([]byte, error) {
	dealerKeys, err := hexArrToParticipantArr(c.DealerPublicKeys)
	if err != nil {
		return nil, fmt.Errorf("could not decode dealer keys: %w", err)
	}
	recipientKeys, err := hexArrToParticipantArr(c.RecipientPublicKeys)
	if err != nil {
		return nil, fmt.Errorf("could not decode recipient keys: %w", err)
	}
	cfg := dkgocrtypes.ReportingPluginConfig{
		T:                   c.T,
		DealerPublicKeys:    dealerKeys,
		RecipientPublicKeys: recipientKeys,
	}
	if c.PreviousInstanceID != nil {
		pid := dkgocrtypes.InstanceID(*c.PreviousInstanceID)
		cfg.PreviousInstanceID = &pid
	}
	return cfg.MarshalBinary()
}

// VaultOffchainConfig mirrors vault.ReportingPluginConfig fields for embedding in V3_1OracleConfig.
// DKGInstanceID is computed at runtime and must be set by the caller before the generator is invoked.
type VaultOffchainConfig struct {
	BatchSize                         int32   `yaml:"batchSize" json:"batchSize"`
	MaxSecretsPerOwner                int32   `yaml:"maxSecretsPerOwner" json:"maxSecretsPerOwner"`
	MaxCiphertextLengthBytes          int32   `yaml:"maxCiphertextLengthBytes" json:"maxCiphertextLengthBytes"`
	MaxIdentifierKeyLengthBytes       int32   `yaml:"maxIdentifierKeyLengthBytes" json:"maxIdentifierKeyLengthBytes"`
	MaxIdentifierOwnerLengthBytes     int32   `yaml:"maxIdentifierOwnerLengthBytes" json:"maxIdentifierOwnerLengthBytes"`
	MaxIdentifierNamespaceLengthBytes int32   `yaml:"maxIdentifierNamespaceLengthBytes" json:"maxIdentifierNamespaceLengthBytes"`
	DKGInstanceID                     *string `yaml:"dkgInstanceID,omitempty" json:"dkgInstanceID,omitempty"`

	LimitsMaxQueryLength                                  int32 `yaml:"limitsMaxQueryLength" json:"limitsMaxQueryLength"`
	LimitsMaxObservationLength                            int32 `yaml:"limitsMaxObservationLength" json:"limitsMaxObservationLength"`
	LimitsMaxReportsPlusPrecursorLength                   int32 `yaml:"limitsMaxReportsPlusPrecursorLength" json:"limitsMaxReportsPlusPrecursorLength"`
	LimitsMaxReportLength                                 int32 `yaml:"limitsMaxReportLength" json:"limitsMaxReportLength"`
	LimitsMaxReportCount                                  int32 `yaml:"limitsMaxReportCount" json:"limitsMaxReportCount"`
	LimitsMaxKeyValueModifiedKeysPlusValuesLength         int32 `yaml:"limitsMaxKeyValueModifiedKeysPlusValuesLength" json:"limitsMaxKeyValueModifiedKeysPlusValuesLength"`
	LimitsMaxBlobPayloadLength                            int32 `yaml:"limitsMaxBlobPayloadLength" json:"limitsMaxBlobPayloadLength"`
	LimitsMaxKeyValueModifiedKeys                         int32 `yaml:"limitsMaxKeyValueModifiedKeys" json:"limitsMaxKeyValueModifiedKeys"`
	LimitsMaxPerOracleUnexpiredBlobCumulativePayloadBytes int32 `yaml:"limitsMaxPerOracleUnexpiredBlobCumulativePayloadBytes" json:"limitsMaxPerOracleUnexpiredBlobCumulativePayloadBytes"`
	LimitsMaxPerOracleUnexpiredBlobCount                  int32 `yaml:"limitsMaxPerOracleUnexpiredBlobCount" json:"limitsMaxPerOracleUnexpiredBlobCount"`
}

func (c *VaultOffchainConfig) Marshal() ([]byte, error) {
	if c.DKGInstanceID == nil {
		return nil, errors.New("must provide a DKGInstanceID")
	}
	pb := &vault.ReportingPluginConfig{
		BatchSize:                         c.BatchSize,
		MaxSecretsPerOwner:                c.MaxSecretsPerOwner,
		MaxCiphertextLengthBytes:          c.MaxCiphertextLengthBytes,
		MaxIdentifierKeyLengthBytes:       c.MaxIdentifierKeyLengthBytes,
		MaxIdentifierOwnerLengthBytes:     c.MaxIdentifierOwnerLengthBytes,
		MaxIdentifierNamespaceLengthBytes: c.MaxIdentifierNamespaceLengthBytes,
		DKGInstanceID:                     c.DKGInstanceID,

		LimitsMaxQueryLength:                                  c.LimitsMaxQueryLength,
		LimitsMaxObservationLength:                            c.LimitsMaxObservationLength,
		LimitsMaxReportsPlusPrecursorLength:                   c.LimitsMaxReportsPlusPrecursorLength,
		LimitsMaxReportLength:                                 c.LimitsMaxReportLength,
		LimitsMaxReportCount:                                  c.LimitsMaxReportCount,
		LimitsMaxKeyValueModifiedKeysPlusValuesLength:         c.LimitsMaxKeyValueModifiedKeysPlusValuesLength,
		LimitsMaxBlobPayloadLength:                            c.LimitsMaxBlobPayloadLength,
		LimitsMaxKeyValueModifiedKeys:                         c.LimitsMaxKeyValueModifiedKeys,
		LimitsMaxPerOracleUnexpiredBlobCumulativePayloadBytes: c.LimitsMaxPerOracleUnexpiredBlobCumulativePayloadBytes,
		LimitsMaxPerOracleUnexpiredBlobCount:                  c.LimitsMaxPerOracleUnexpiredBlobCount,
	}
	return proto.Marshal(pb)
}

// ConsensusQueueOffchainConfig is the consensus-queue plugin config written into the OCR offchain config.
// The node reads its OCR limits only from here, so the config digest covers them.
type ConsensusQueueOffchainConfig struct {
	Limits ConsensusQueueLimits `yaml:"limits" json:"limits"`
}

// ConsensusQueueLimits mirrors libocr's ocr3_1types.ReportingPluginLimits.
type ConsensusQueueLimits struct {
	MaxQueryBytes                                   uint32 `yaml:"maxQueryBytes" json:"maxQueryBytes"`
	MaxObservationBytes                             uint32 `yaml:"maxObservationBytes" json:"maxObservationBytes"`
	MaxReportsPlusPrecursorBytes                    uint32 `yaml:"maxReportsPlusPrecursorBytes" json:"maxReportsPlusPrecursorBytes"`
	MaxReportBytes                                  uint32 `yaml:"maxReportBytes" json:"maxReportBytes"`
	MaxReportCount                                  uint32 `yaml:"maxReportCount" json:"maxReportCount"`
	MaxKeyValueModifiedKeys                         uint32 `yaml:"maxKeyValueModifiedKeys" json:"maxKeyValueModifiedKeys"`
	MaxKeyValueModifiedKeysPlusValuesBytes          uint32 `yaml:"maxKeyValueModifiedKeysPlusValuesBytes" json:"maxKeyValueModifiedKeysPlusValuesBytes"`
	MaxBlobPayloadBytes                             uint32 `yaml:"maxBlobPayloadBytes" json:"maxBlobPayloadBytes"`
	MaxPerOracleUnexpiredBlobCumulativePayloadBytes uint32 `yaml:"maxPerOracleUnexpiredBlobCumulativePayloadBytes" json:"maxPerOracleUnexpiredBlobCumulativePayloadBytes"`
	MaxPerOracleUnexpiredBlobCount                  uint32 `yaml:"maxPerOracleUnexpiredBlobCount" json:"maxPerOracleUnexpiredBlobCount"`
}

// consensusQueueReferenceBlobPayloadPerRound is the sustained blob payload the default per-oracle quota is sized
// for: a full observation of the calibration's ~1 KiB (EVM log) events each round. The config field is a uint32, so
// the plugin's worst case (MaxBlobHandlesPerObservation full blobs every round for the whole event lifetime) cannot
// be expressed. DONs that sustain larger trigger payloads must raise MaxPerOracleUnexpiredBlobCumulativePayloadBytes.
const consensusQueueReferenceBlobPayloadPerRound = 1 << 20

// DefaultConsensusQueueOffchainConfigForRoundInterval returns starting limits for a consensus-queue DON whose OCR
// DeltaRound is roundDur.
//
// The plugin broadcasts trigger events as blobs and keeps them until their events expire, which is the workflow's
// TriggerEventQueueTimeout plus the consensus-queue expiration grace. The blob quotas must therefore cover every blob
// a node can broadcast over that lifetime: a fixed default is wrong for a DON whose rounds are faster or slower than
// the one it was sized for. The plugin has no queries or reports, so those limits are zero, and it writes a fixed
// handful of KV keys per round.
//
// The sizes use the default TriggerEventQueueTimeout and ConsensusQueueExpirationGraceInterval. A node or workflow
// that raises either keeps blobs longer, so a DON using larger values needs a proportionally larger quota.
//
// The node never applies defaults itself. Different node versions would otherwise disagree under the same digest.
func DefaultConsensusQueueOffchainConfigForRoundInterval(roundDur time.Duration) *ConsensusQueueOffchainConfig {
	blobCount, blobPayloadBytes := queuelimits.BlobQuota(
		roundDur,
		cresettings.Default.PerWorkflow.TriggerEventQueueTimeout.DefaultValue,
		cresettings.Default.ConsensusQueueExpirationGraceInterval.DefaultValue,
		consensusQueueReferenceBlobPayloadPerRound,
	)
	return &ConsensusQueueOffchainConfig{
		Limits: ConsensusQueueLimits{
			MaxObservationBytes:                             512 * 1024,
			MaxKeyValueModifiedKeys:                         4,
			MaxKeyValueModifiedKeysPlusValuesBytes:          10 * 1024 * 1024,
			MaxBlobPayloadBytes:                             5 * 1024 * 1024,
			MaxPerOracleUnexpiredBlobCumulativePayloadBytes: blobPayloadBytes,
			MaxPerOracleUnexpiredBlobCount:                  blobCount,
		},
	}
}

// DefaultConsensusQueueOffchainConfig returns the starting limits for a one second DeltaRound, a common default.
// Use DefaultConsensusQueueOffchainConfigForRoundInterval when the DON's DeltaRound is known to be different.
func DefaultConsensusQueueOffchainConfig() *ConsensusQueueOffchainConfig {
	return DefaultConsensusQueueOffchainConfigForRoundInterval(time.Second)
}

// Marshal encodes the config as a queuepb.ReportingPluginConfig.
// Returns an error if the limits would be rejected by the node (see queuepb.Decode).
func (c *ConsensusQueueOffchainConfig) Marshal() ([]byte, error) {
	l := c.Limits
	b, err := proto.Marshal(&queuepb.ReportingPluginConfig{
		Limits: &queuepb.ReportingPluginLimits{
			MaxQueryBytes:                                   l.MaxQueryBytes,
			MaxObservationBytes:                             l.MaxObservationBytes,
			MaxReportsPlusPrecursorBytes:                    l.MaxReportsPlusPrecursorBytes,
			MaxReportBytes:                                  l.MaxReportBytes,
			MaxReportCount:                                  l.MaxReportCount,
			MaxKeyValueModifiedKeys:                         l.MaxKeyValueModifiedKeys,
			MaxKeyValueModifiedKeysPlusValuesBytes:          l.MaxKeyValueModifiedKeysPlusValuesBytes,
			MaxBlobPayloadBytes:                             l.MaxBlobPayloadBytes,
			MaxPerOracleUnexpiredBlobCumulativePayloadBytes: l.MaxPerOracleUnexpiredBlobCumulativePayloadBytes,
			MaxPerOracleUnexpiredBlobCount:                  l.MaxPerOracleUnexpiredBlobCount,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal consensus-queue offchain config: %w", err)
	}
	if _, err := queuepb.Decode(b); err != nil {
		return nil, err
	}
	return b, nil
}

func getPluginConfig(cfg V3_1OracleConfig) (pluginConfig, error) {
	var set []pluginConfig
	if cfg.DKGOffchainConfig != nil {
		set = append(set, cfg.DKGOffchainConfig)
	}
	if cfg.VaultOffchainConfig != nil {
		set = append(set, cfg.VaultOffchainConfig)
	}
	if cfg.ConsensusQueueOffchainConfig != nil {
		set = append(set, cfg.ConsensusQueueOffchainConfig)
	}

	switch len(set) {
	case 0:
		return nil, nil
	case 1:
		return set[0], nil
	default:
		return nil, errors.New("multiple offchain configs specified in V3_1OracleConfig: only one of DKGOffchainConfig, VaultOffchainConfig or ConsensusQueueOffchainConfig may be set")
	}
}
