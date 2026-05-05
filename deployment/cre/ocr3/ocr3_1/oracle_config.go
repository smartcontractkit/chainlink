package ocr3_1

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

type pluginConfig interface {
	Marshal() ([]byte, error)
}

// DKGOffchainConfig wraps dkgocrtypes.ReportingPluginConfig for embedding in V3_1OracleConfig.
// Field names in YAML/JSON use Go defaults (no tags on the underlying type).
type DKGOffchainConfig struct {
	dkgocrtypes.ReportingPluginConfig `yaml:",inline"`
	DealerPublicKeys                  []string `yaml:"dealerPublicKeys" json:"dealerPublicKeys"`
	RecipientPublicKeys               []string `yaml:"recipientPublicKeys" json:"recipientPublicKeys"`
	T                                 int      `yaml:"t" json:"t"`
	PreviousInstanceID                *string  `yaml:"previousInstanceID" json:"previousInstanceID"`
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
	var iid dkgocrtypes.InstanceID
	if c.PreviousInstanceID != nil {
		iid = dkgocrtypes.InstanceID(*c.PreviousInstanceID)
	}
	cfg := dkgocrtypes.ReportingPluginConfig{
		T:                   c.T,
		PreviousInstanceID:  &iid,
		DealerPublicKeys:    dealerKeys,
		RecipientPublicKeys: recipientKeys,
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

func getPluginConfig(cfg V3_1OracleConfig) (pluginConfig, error) {
	var result pluginConfig
	if cfg.DKGOffchainConfig != nil {
		result = cfg.DKGOffchainConfig
	}
	if cfg.VaultOffchainConfig != nil {
		if result != nil {
			return nil, errors.New("multiple offchain configs specified in V3_1OracleConfig: only one of DKGOffchainConfig or VaultOffchainConfig may be set")
		}
		result = cfg.VaultOffchainConfig
	}
	return result, nil
}
