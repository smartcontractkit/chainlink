package ocr3

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"

	libocrtypes "github.com/smartcontractkit/libocr/ragep2p/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment"
)

const (
	// OCR3ConfigDefaultKey is the default key for single-instance OCR capabilities
	// in the ocr3Configs map. Mirrors the constant in chainlink-common/pkg/capabilities/pb.
	OCR3ConfigDefaultKey = "__default__"
)

// ComputeOCR3Config generates a full OCR3 configuration (signers, transmitters,
// offchain config, etc.) from oracle config parameters and node information
// fetched from Job Distributor.
func ComputeOCR3Config(
	env cldf.Environment,
	registryChainSel uint64,
	registryNodes []capabilities_registry_v2.INodeInfoProviderNodeInfo,
	oracleOffchainConfig OracleConfig,
	reportingPluginConfig []byte,
	extraSignerFamilies []string,
) (*OCR2OracleConfig, error) {
	p2pIDStrings := RegistryNodesToP2PIDs(registryNodes)

	jdNodes, err := deployment.NodeInfo(p2pIDStrings, env.Offchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get node info from Job Distributor: %w", err)
	}

	config, err := GenerateOCR3ConfigFromNodes(
		oracleOffchainConfig,
		jdNodes,
		registryChainSel,
		env.OCRSecrets,
		reportingPluginConfig,
		extraSignerFamilies,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OCR3 config: %w", err)
	}

	return &config, nil
}

// RegistryNodesToP2PIDs converts registry node P2P IDs ([32]byte) to string
// format ("p2p_<base58>") for Job Distributor lookups.
func RegistryNodesToP2PIDs(nodes []capabilities_registry_v2.INodeInfoProviderNodeInfo) []string {
	result := make([]string, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, "p2p_"+libocrtypes.PeerID(node.P2pId).String())
	}
	return result
}

// ValidateOCR3Config validates OCR3 config fields according to the protocol rules.
func ValidateOCR3Config(signers [][]byte, transmitters [][]byte, f uint32) error {
	if len(signers) != len(transmitters) {
		return fmt.Errorf("signers count (%d) != transmitters count (%d)", len(signers), len(transmitters))
	}

	if f == 0 {
		return errors.New("f must be positive")
	}

	n := uint32(len(signers)) //nolint:gosec // G115 - len(signers) is bounded by the number of DON nodes, never overflows uint32
	if n <= 3*f {
		return fmt.Errorf("not enough nodes for f=%d: need at least %d, got %d (3f+1 rule)", f, 3*f+1, n)
	}

	for i, t := range transmitters {
		if len(t) == 0 {
			return fmt.Errorf("transmitter %d is empty", i)
		}
	}

	return nil
}

// ocr3ConfigsJSON is a helper struct for parsing the ocr3Configs field from
// a CapabilityConfig protobuf via JSON roundtrip.
type ocr3ConfigsJSON struct {
	Ocr3Configs map[string]struct {
		ConfigCount uint64 `json:"configCount"`
	} `json:"ocr3Configs"`
}

// GetCurrentOCR3ConfigCount reads the current OCR3 config count for a
// capability/DON from the on-chain registry. Returns 0 if no config exists.
func GetCurrentOCR3ConfigCount(
	capReg *capabilities_registry_v2.CapabilitiesRegistry,
	donName string,
	capabilityID string,
	ocrConfigKey string,
) (uint64, error) {
	don, err := capReg.GetDONByName(&bind.CallOpts{}, donName)
	if err != nil {
		return 0, fmt.Errorf("failed to get DON by name %q: %w", donName, err)
	}

	for _, capCfg := range don.CapabilityConfigurations {
		if capCfg.CapabilityId == capabilityID {
			return extractOCR3ConfigCount(capCfg.Config, ocrConfigKey)
		}
	}
	return 0, nil
}

func extractOCR3ConfigCount(rawConfig []byte, ocrConfigKey string) (uint64, error) {
	if len(rawConfig) == 0 {
		return 0, nil
	}
	pbCfg := &capabilitiespb.CapabilityConfig{}
	if err := proto.Unmarshal(rawConfig, pbCfg); err != nil {
		return 0, fmt.Errorf("failed to unmarshal capability config: %w", err)
	}
	if pbCfg.Ocr3Configs == nil || pbCfg.Ocr3Configs[ocrConfigKey] == nil {
		return 0, nil
	}
	return pbCfg.Ocr3Configs[ocrConfigKey].ConfigCount, nil
}

// OCR2OracleConfigToMap converts an OCR2OracleConfig to a protojson-compatible
// map suitable for injection into a capability config's ocr3Configs field.
// For example, []byte fields are base64-encoded.
func OCR2OracleConfigToMap(config *OCR2OracleConfig, configCount uint64) map[string]any {
	signers := make([]string, len(config.Signers))
	for i, s := range config.Signers {
		signers[i] = base64.StdEncoding.EncodeToString(s)
	}

	transmitters := make([]string, len(config.Transmitters))
	for i, t := range config.Transmitters {
		transmitters[i] = base64.StdEncoding.EncodeToString(t.Bytes())
	}

	result := map[string]any{
		"signers":               signers,
		"transmitters":          transmitters,
		"f":                     uint32(config.F),
		"offchainConfigVersion": config.OffchainConfigVersion,
		"configCount":           configCount,
	}

	if len(config.OnchainConfig) > 0 {
		result["onchainConfig"] = base64.StdEncoding.EncodeToString(config.OnchainConfig)
	}
	if len(config.OffchainConfig) > 0 {
		result["offchainConfig"] = base64.StdEncoding.EncodeToString(config.OffchainConfig)
	}

	return result
}

// ValidateOCR2OracleConfig validates a generated OCR2OracleConfig.
func ValidateOCR2OracleConfig(config *OCR2OracleConfig) error {
	transmitterBytes := make([][]byte, len(config.Transmitters))
	for i, t := range config.Transmitters {
		if t == (common.Address{}) {
			return fmt.Errorf("transmitter %d is zero address", i)
		}
		transmitterBytes[i] = t.Bytes()
	}
	return ValidateOCR3Config(config.Signers, transmitterBytes, uint32(config.F))
}
