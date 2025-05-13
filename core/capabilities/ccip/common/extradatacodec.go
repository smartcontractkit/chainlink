package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// ExtraDataCodec is a struct that holds the chain specific extra data codec
type ExtraDataCodec struct {
	EVMExtraDataCodec    SourceChainExtraDataCodec
	SolanaExtraDataCodec SourceChainExtraDataCodec
}

// NewExtraDataCodec is a constructor for ExtraDataCodec
func NewExtraDataCodec(evmExtraDataCodec, solanaExtraDataCodec SourceChainExtraDataCodec) ExtraDataCodec {
	return ExtraDataCodec{
		EVMExtraDataCodec:    evmExtraDataCodec,
		SolanaExtraDataCodec: solanaExtraDataCodec,
	}
}

// DecodeExtraArgs reformats bytes into a chain agnostic map[string]any representation for extra args
func (c ExtraDataCodec) DecodeExtraArgs(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
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
		return c.EVMExtraDataCodec.DecodeExtraArgsToMap(extraArgs)

	case chainsel.FamilySolana:
		return c.SolanaExtraDataCodec.DecodeExtraArgsToMap(extraArgs)

	default:
		return nil, fmt.Errorf("unsupported family for extra args type %s", family)
	}
}

// DecodeTokenAmountDestExecData reformats bytes to chain-agnostic map[string]any for tokenAmount DestExecData field
func (c ExtraDataCodec) DecodeTokenAmountDestExecData(destExecData cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
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
		return c.EVMExtraDataCodec.DecodeDestExecDataToMap(destExecData)

	case chainsel.FamilySolana:
		return c.SolanaExtraDataCodec.DecodeDestExecDataToMap(destExecData)

	default:
		return nil, fmt.Errorf("unsupported family for extra args type %s", family)
	}
}
