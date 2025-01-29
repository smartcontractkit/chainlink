package common

import (
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type ExtraDataCodec struct{}

func NewExtraDataCodec() ExtraDataCodec {
	return ExtraDataCodec{}
}

func (c ExtraDataCodec) DecodeExtraData(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	// Not implemented and not return error
	return nil, nil
}

func (c ExtraDataCodec) DecodeExtraArgs(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	// Not implemented and not return error
	return nil, nil
}

func (c ExtraDataCodec) DecodeTokenAmountDestExecData(destExecData cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	// Not implemented and not return error
	return nil, nil
}
