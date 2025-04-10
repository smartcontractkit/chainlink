package ccipsolana

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	solanaconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/solana"
)

// GetCRCW is a struct that implements the GetChainReaderWriter interface for Solana chains.
type GetCRCW struct{}

// GetChainWriter returns a new ContractWriter for Solana chains.
func (g GetCRCW) GetChainWriter(ctx context.Context, pararms ccipcommon.GetChainWriterParams) (types.ContractWriter, error) {
	var offrampProgramAddress solana.PublicKey
	// NOTE: this function can still be called with EVM inputs, and PublicKeyFromBytes will panic on addresses with len=20
	// technically we only need the writer to do fee estimation so this doesn't matter and we can use a zero address
	if len(pararms.OfframpProgramAddress) == solana.PublicKeyLength {
		offrampProgramAddress = solana.PublicKeyFromBytes(pararms.OfframpProgramAddress)
	}

	transmitter := pararms.Transmitters[types.NewRelayID(pararms.ChainFamily, pararms.ChainID)]
	solConfig, err := solanaconfig.GetSolanaChainWriterConfig(offrampProgramAddress.String(), transmitter[0])
	if err != nil {
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
