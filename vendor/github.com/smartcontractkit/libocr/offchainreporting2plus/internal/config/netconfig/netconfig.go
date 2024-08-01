package netconfig

import (
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr2config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type NetConfig struct {
	ConfigDigest types.ConfigDigest
	F            int
	PeerIDs      []string
}

func NetConfigFromContractConfig(contractConfig types.ContractConfig) (NetConfig, error) {
	switch contractConfig.OffchainConfigVersion {
	case config.OCR2OffchainConfigVersion:
		publicConfig, err := ocr2config.PublicConfigFromContractConfig(true, contractConfig)
		if err != nil {
			return NetConfig{}, err
		}
		return NetConfig{
			publicConfig.ConfigDigest,
			publicConfig.F,
			peerIDs(publicConfig.OracleIdentities),
		}, nil
	case config.OCR3OffchainConfigVersion:
		publicConfig, err := ocr3config.PublicConfigFromContractConfig(true, contractConfig)
		if err != nil {
			return NetConfig{}, err
		}
		return NetConfig{
			publicConfig.ConfigDigest,
			publicConfig.F,
			peerIDs(publicConfig.OracleIdentities),
		}, nil
	default:
		return NetConfig{}, fmt.Errorf("NetConfigFromContractConfig received OffchainConfigVersion %v", contractConfig.OffchainConfigVersion)
	}
}

func peerIDs(identities []config.OracleIdentity) []string {
	var peerIDs []string
	for _, identity := range identities {
		peerIDs = append(peerIDs, identity.PeerID)
	}
	return peerIDs
}
