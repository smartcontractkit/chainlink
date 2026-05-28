package v1_6

import (
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
)

// PublicConfigFromCCIPHome unwraps the OCR3 OffchainConfig payload carried in a CCIPHome
// versioned config by running ocr3confighelper.PublicConfigFromContractConfig. This is the
// shared assembly used by the view layer (populateDecodedOCRParams) and any caller that
// needs the typed plugin offchain config via DecodeCommitOffchain / DecodeExecuteOffchain.
// Returns the zero PublicConfig (no error) when the versioned config is empty.
func PublicConfigFromCCIPHome(v ccip_home.CCIPHomeVersionedConfig) (ocr3confighelper.PublicConfig, error) {
	if v.ConfigDigest == [32]byte{} || len(v.Config.OffchainConfig) == 0 {
		return ocr3confighelper.PublicConfig{}, nil
	}
	signers := make([]ocrtypes.OnchainPublicKey, 0, len(v.Config.Nodes))
	transmitters := make([]ocrtypes.Account, 0, len(v.Config.Nodes))
	for _, node := range v.Config.Nodes {
		signers = append(signers, node.SignerKey)
		transmitters = append(transmitters, ocrtypes.Account(node.TransmitterKey))
	}
	publicConfig, err := ocr3confighelper.PublicConfigFromContractConfig(false, ocrtypes.ContractConfig{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     v.Config.FRoleDON,
		OnchainConfig:         []byte{},
		OffchainConfigVersion: v.Config.OffchainConfigVersion,
		OffchainConfig:        v.Config.OffchainConfig,
	})
	if err != nil {
		return ocr3confighelper.PublicConfig{}, fmt.Errorf("PublicConfigFromContractConfig: %w", err)
	}

	return publicConfig, nil
}

// DecodeCommitOffchain returns the typed CommitOffchainConfig embedded in a CCIPHome
// versioned OCR3 config. Returns the zero value (no error) when the versioned config or
// its ReportingPluginConfig is empty.
func DecodeCommitOffchain(v ccip_home.CCIPHomeVersionedConfig) (pluginconfig.CommitOffchainConfig, error) {
	publicConfig, err := PublicConfigFromCCIPHome(v)
	if err != nil {
		return pluginconfig.CommitOffchainConfig{}, err
	}
	if len(publicConfig.ReportingPluginConfig) == 0 {
		return pluginconfig.CommitOffchainConfig{}, nil
	}
	out, err := pluginconfig.DecodeCommitOffchainConfig(publicConfig.ReportingPluginConfig)
	if err != nil {
		return pluginconfig.CommitOffchainConfig{}, fmt.Errorf("decode CommitOffchainConfig: %w", err)
	}

	return out, nil
}

// DecodeExecuteOffchain is the exec counterpart of DecodeCommitOffchain.
func DecodeExecuteOffchain(v ccip_home.CCIPHomeVersionedConfig) (pluginconfig.ExecuteOffchainConfig, error) {
	publicConfig, err := PublicConfigFromCCIPHome(v)
	if err != nil {
		return pluginconfig.ExecuteOffchainConfig{}, err
	}
	if len(publicConfig.ReportingPluginConfig) == 0 {
		return pluginconfig.ExecuteOffchainConfig{}, nil
	}
	out, err := pluginconfig.DecodeExecuteOffchainConfig(publicConfig.ReportingPluginConfig)
	if err != nil {
		return pluginconfig.ExecuteOffchainConfig{}, fmt.Errorf("decode ExecuteOffchainConfig: %w", err)
	}

	return out, nil
}
