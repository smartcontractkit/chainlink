package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
)

// chainCapabilityVersion is the version of the registry-launched chain capabilities.
const chainCapabilityVersion = "1.0.0"

// ChainCapabilityEntry is one chain capability of a DON in the offchain capabilities registry
// payload: its capability ID and the config the nodes launch it with. The config carries the shared
// (DON-wide) spec_config and an unset localOnly, which is false and matches the onchain value.
type ChainCapabilityEntry struct {
	CapabilityID string
	Config       *capabilitiespb.CapabilityConfig
}

// ChainCapabilityID returns the capability ID of a chain capability, for example
// "evm:ChainSelector:3379446385462418246@1.0.0". family is "evm", "aptos" or "stellar".
func ChainCapabilityID(family string, chainSelector uint64) string {
	return fmt.Sprintf("%s:ChainSelector:%d@%s", family, chainSelector, chainCapabilityVersion)
}

// EVMCapabilityEntryInput holds what the EVM job-spec changeset takes for one chain, minus the nodes:
// the registry payload carries one shared config per DON, so the per-node nodeAddress is not part of it
// (the node fills it from its OCR3 transmitter at launch).
type EVMCapabilityEntryInput struct {
	ChainSelector    uint64
	ForwarderAddress string
	DeltaStage       time.Duration
	// ForwarderLookbackBlocks is the DON-wide default; a non-zero Override value wins.
	ForwarderLookbackBlocks int64
	Override                OverrideDefaultCfg
}

// EVMCapabilityEntry builds the payload entry of an EVM chain, with the same derivation and defaults as
// ProposeEVMCapJobSpec.
func EVMCapabilityEntry(in EVMCapabilityEntryInput) (ChainCapabilityEntry, error) {
	chainID, err := chainIDFromSelector(in.ChainSelector)
	if err != nil {
		return ChainCapabilityEntry{}, err
	}
	numericChainID, err := strconv.ParseUint(chainID, 10, 64)
	if err != nil {
		return ChainCapabilityEntry{}, fmt.Errorf("failed to parse chain ID %s: %w", chainID, err)
	}
	if in.ForwarderAddress == "" {
		return ChainCapabilityEntry{}, fmt.Errorf("evm chain selector %d: forwarder address is required", in.ChainSelector)
	}
	cfg := evmCapConfig(in.Override, numericChainID, in.ForwarderAddress, in.DeltaStage, in.ForwarderLookbackBlocks)
	cfg.NodeAddress = "" // per node; filled at launch
	return chainCapabilityEntry("evm", in.ChainSelector, cfg)
}

// AptosCapabilityEntryInput holds what the Aptos job-spec changeset takes for one chain, minus the nodes.
type AptosCapabilityEntryInput struct {
	ChainSelector          uint64
	ForwarderAddress       string
	DeltaStage             time.Duration
	TxSearchStartingBuffer time.Duration
	P2PToTransmitterMap    map[string]string
	Override               AptosOverrideDefaultCfg
}

// AptosCapabilityEntry builds the payload entry of an Aptos chain, with the same derivation as
// ProposeAptosCapJobSpec.
func AptosCapabilityEntry(in AptosCapabilityEntryInput) (ChainCapabilityEntry, error) {
	chainID, err := chainIDFromSelector(in.ChainSelector)
	if err != nil {
		return ChainCapabilityEntry{}, err
	}
	if in.ForwarderAddress == "" {
		return ChainCapabilityEntry{}, fmt.Errorf("aptos chain selector %d: forwarder address is required", in.ChainSelector)
	}
	cfg := aptosCapConfig(in.Override, chainID, in.ForwarderAddress, in.P2PToTransmitterMap, in.DeltaStage, in.TxSearchStartingBuffer)
	return chainCapabilityEntry("aptos", in.ChainSelector, cfg)
}

// StellarCapabilityEntryInput holds what the Stellar job-spec changeset takes for one chain, minus the nodes.
type StellarCapabilityEntryInput struct {
	ChainSelector    uint64
	ForwarderAddress string
	DeltaStage       time.Duration
	// ForwarderLookbackLedgers is the DON-wide default; a non-zero Override value wins.
	ForwarderLookbackLedgers int64
	Override                 StellarOverrideDefaultCfg
}

// StellarCapabilityEntry builds the payload entry of a Stellar chain, with the same derivation as
// ProposeStellarCapJobSpec.
func StellarCapabilityEntry(in StellarCapabilityEntryInput) (ChainCapabilityEntry, error) {
	chainID, err := chainIDFromSelector(in.ChainSelector)
	if err != nil {
		return ChainCapabilityEntry{}, err
	}
	if in.ForwarderAddress == "" {
		return ChainCapabilityEntry{}, fmt.Errorf("stellar chain selector %d: forwarder address is required", in.ChainSelector)
	}
	cfg := stellarCapConfig(in.Override, chainID, in.ForwarderAddress, in.DeltaStage, in.ForwarderLookbackLedgers)
	return chainCapabilityEntry("stellar", in.ChainSelector, cfg)
}

func chainIDFromSelector(chainSelector uint64) (string, error) {
	chainID, err := chainselectors.GetChainIDFromSelector(chainSelector)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID from selector %d: %w", chainSelector, err)
	}
	return chainID, nil
}

func chainCapabilityEntry(family string, chainSelector uint64, cfg any) (ChainCapabilityEntry, error) {
	specConfig, err := specConfigFromStruct(cfg)
	if err != nil {
		return ChainCapabilityEntry{}, fmt.Errorf("%s chain selector %d: %w", family, chainSelector, err)
	}
	return ChainCapabilityEntry{
		CapabilityID: ChainCapabilityID(family, chainSelector),
		Config:       &capabilitiespb.CapabilityConfig{SpecConfig: specConfig},
	}, nil
}

// specConfigFromStruct converts a capability config struct to registry spec_config. It goes through the
// struct's JSON form, so the keys and omitted-when-empty fields are exactly those of the job-spec config
// the node used to receive. Integral numbers stay integral (int64) and are not widened through float64.
func specConfigFromStruct(cfg any) (*valuespb.Map, error) {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capability config: %w", err)
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var generic map[string]any
	if err = dec.Decode(&generic); err != nil {
		return nil, fmt.Errorf("failed to decode capability config: %w", err)
	}
	for key, v := range generic {
		generic[key], err = numberValue(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", key, err)
		}
	}
	wrapped, err := values.NewMap(generic)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap capability config: %w", err)
	}
	return values.ProtoMap(wrapped), nil
}

func numberValue(v any) (any, error) {
	switch n := v.(type) {
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return i, nil
		}
		f, err := n.Float64()
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", n.String(), err)
		}
		return f, nil
	case map[string]any:
		for k, el := range n {
			converted, err := numberValue(el)
			if err != nil {
				return nil, err
			}
			n[k] = converted
		}
		return n, nil
	default:
		return v, nil
	}
}
