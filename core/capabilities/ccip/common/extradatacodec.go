package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// RealExtraDataCodec is a concrete implementation of ExtraDataCodec
type RealExtraDataCodec struct {
	EVMExtraDataDecoder    ChainSpecificExtraDataDecoder
	SolanaExtraDataDecoder ChainSpecificExtraDataDecoder
}

// NewExtraDataCodec is a constructor for RealExtraDataCodec
func NewExtraDataCodec(evmExtraDataDecoder, solanaExtraDataDecoder ChainSpecificExtraDataDecoder) RealExtraDataCodec {
	return RealExtraDataCodec{
		EVMExtraDataDecoder:    evmExtraDataDecoder,
		SolanaExtraDataDecoder: solanaExtraDataDecoder,
	}
}

// DecodeExtraArgs reformats bytes into a chain agnostic map[string]any representation for extra args
func (c RealExtraDataCodec) DecodeExtraArgs(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	if len(extraArgs) == 0 {
		// return empty map if extraArgs is empty
		return nil, nil
	}

	family, err := chainsel.GetSelectorFamily(uint64(sourceChainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for selector %d: %w", sourceChainSelector, err)
	}

	switch family {
	case chainsel.FamilyEVM:
		return c.EVMExtraDataDecoder.DecodeExtraArgsToMap(extraArgs)

	case chainsel.FamilySolana:
		return c.SolanaExtraDataDecoder.DecodeExtraArgsToMap(extraArgs)

	default:
		return nil, fmt.Errorf("unsupported family for extra args type %s", family)
	}
}

// DecodeTokenAmountDestExecData reformats bytes to chain-agnostic map[string]any for tokenAmount DestExecData field
func (c RealExtraDataCodec) DecodeTokenAmountDestExecData(destExecData cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	if len(destExecData) == 0 {
		// return empty map if destExecData is empty
		return nil, nil
	}

	family, err := chainsel.GetSelectorFamily(uint64(sourceChainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for selector %d: %w", sourceChainSelector, err)
	}

	switch family {
	case chainsel.FamilyEVM:
		return c.EVMExtraDataDecoder.DecodeDestExecDataToMap(destExecData)

	case chainsel.FamilySolana:
		return c.SolanaExtraDataDecoder.DecodeDestExecDataToMap(destExecData)

	default:
		return nil, fmt.Errorf("unsupported family for extra args type %s", family)
	}
}
