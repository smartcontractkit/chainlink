package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type ExtraDataCodec interface {
	// DecodeExtraArgs reformat bytes into a chain agnostic map[string]any representation for extra args
	DecodeExtraArgs(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error)
	// DecodeTokenAmountDestExecData reformat bytes to chain-agnostic map[string]any for tokenAmount DestExecData field
	DecodeTokenAmountDestExecData(destExecData cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error)
}

type ExtraDataDecoder interface {
	DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error)
	DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]any, error)
}

type RealExtraDataCodec struct {
	EVMExtraDataDecoder    ExtraDataDecoder
	SolanaExtraDataDecoder ExtraDataDecoder
}

func NewExtraDataCodec(evmDecoder ExtraDataDecoder, solanaDecoder ExtraDataDecoder) RealExtraDataCodec {
	return RealExtraDataCodec{
		EVMExtraDataDecoder:    evmDecoder,
		SolanaExtraDataDecoder: solanaDecoder,
	}
}

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
