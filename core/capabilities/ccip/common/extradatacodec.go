package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// ExtraDataCodec is an interface for decoding extra args and dest exec data into a chain-agnostic map[string]any representation
type ExtraDataCodec interface {
	// DecodeExtraArgs reformat bytes into a chain agnostic map[string]any representation for extra args
	DecodeExtraArgs(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error)
	// DecodeTokenAmountDestExecData reformat bytes to chain-agnostic map[string]any for tokenAmount DestExecData field
	DecodeTokenAmountDestExecData(destExecData cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error)
}

// ExtraDataDecoder is an interface for decoding extra args and dest exec data into a map[string]any representation
type ExtraDataDecoder interface {
	DecodeExtraArgsToMap(extraArgs cciptypes.Bytes) (map[string]any, error)
	DecodeDestExecDataToMap(destExecData cciptypes.Bytes) (map[string]any, error)
}

// RealExtraDataCodec is a concrete implementation of ExtraDataCodec
type RealExtraDataCodec struct {
	EVMExtraDataDecoder    ExtraDataDecoder
	SolanaExtraDataDecoder ExtraDataDecoder
}

// ExtraDataCodecParams is a struct that holds the parameters for creating a RealExtraDataCodec
type ExtraDataCodecParams struct {
	evmExtraDataDecoder    ExtraDataDecoder
	solanaExtraDataDecoder ExtraDataDecoder
}

// NewExtraDataCodecParams is a constructor for ExtraDataCodecParams
func NewExtraDataCodecParams(evmDecoder ExtraDataDecoder, solanaDecoder ExtraDataDecoder) ExtraDataCodecParams {
	return ExtraDataCodecParams{
		evmExtraDataDecoder:    evmDecoder,
		solanaExtraDataDecoder: solanaDecoder,
	}
}

// NewExtraDataCodec is a constructor for RealExtraDataCodec
func NewExtraDataCodec(params ExtraDataCodecParams) RealExtraDataCodec {
	return RealExtraDataCodec{
		EVMExtraDataDecoder:    params.evmExtraDataDecoder,
		SolanaExtraDataDecoder: params.solanaExtraDataDecoder,
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
