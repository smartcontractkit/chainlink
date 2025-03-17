package ccipsolana

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	solanaconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/solana"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// PluginConfig is a struct that contains the configuration for a plugin.
type PluginConfig struct {
	extraDataCodec ccipcommon.ExtraDataCodec
}

// NewPluginConfig returns a new PluginConfig.
func NewPluginConfig(extraDataCodec ccipcommon.ExtraDataCodec) *PluginConfig {
	return &PluginConfig{
		extraDataCodec: extraDataCodec,
	}
}

// InitializePluginConfig returns a PluginConfig for Solana chains.
func (p PluginConfig) InitializePluginConfig() ccipcommon.PluginConfig {
	return ccipcommon.PluginConfig{
		CommitPluginCodec:  NewCommitPluginCodecV1(),
		ExecutePluginCodec: NewExecutePluginCodecV1(p.extraDataCodec),
		MessageHasher: func(lggr logger.Logger) cciptypes.MessageHasher {
			return NewMessageHasherV1(lggr, p.extraDataCodec)
		},
		TokenDataEncoder:           NewSolanaTokenDataEncoder(),
		GasEstimateProvider:        NewGasEstimateProvider(),
		RMNCrypto:                  func(lggr logger.Logger) cciptypes.RMNCrypto { return nil },
		ContractTransmitterFactory: ocrimpls.NewSVMContractTransmitterFactory(p.extraDataCodec),
		GetChainReaderWriter:       GetCRCW{},
	}
}

// GetCRCW is a struct that implements the GetChainReaderWriter interface for Solana chains.
type GetCRCW struct{}

// GetChainWriter returns a new ContractWriter for Solana chains.
func (g GetCRCW) GetChainWriter(ctx context.Context, pararms ccipcommon.GetChainWriterParams) (types.ContractWriter, error) {
	if solana.PublicKeyLength != len(pararms.OfframpProgramAddress) {
		return nil, fmt.Errorf("invalid offrampProgramAddress length: %d", len(pararms.OfframpProgramAddress))
	}

	offrampAddress := solana.PublicKeyFromBytes(pararms.OfframpProgramAddress)
	transmitter := pararms.Transmitters[types.NewRelayID(pararms.ChainFamily, pararms.ChainID)]
	solConfig, err := solanaconfig.GetSolanaChainWriterConfig(offrampAddress.String(), transmitter[0])
	if err == nil {
		return nil, fmt.Errorf("failed to get Solana chain writer config: %w", err)
	}
	chainWriterConfig, err := json.Marshal(solConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Solana chain writer config: %w", err)
	}

	cw, err := pararms.Relayer.NewContractWriter(ctx, chainWriterConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain writer for chain %s: %w", pararms.ChainID, err)
	}

	return cw, nil
}

// GetChainReader returns a new ContractReader for Solana chains.
func (g GetCRCW) GetChainReader(ctx context.Context, params ccipcommon.GetChainReaderParams) (types.ContractReader, error) {
	var err error
	var cfg config.ContractReader
	if params.ChainID == params.DestChainID {
		cfg, err = solanaconfig.DestContractReaderConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get Solana dest contract reader config: %w", err)
		}
	} else {
		cfg, err = solanaconfig.SourceContractReaderConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get Solana source contract reader config: %w", err)
		}
	}

	marshaledConfig, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chain reader config: %w", err)
	}

	cr, err := params.Relayer.NewContractReader(ctx, marshaledConfig)
	if err != nil {
		return nil, err
	}

	return cr, nil
}
