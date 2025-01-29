package common

import (
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type ExtraArgsCodec struct{}

func NewExtraArgsCodec() ExtraArgsCodec {
	return ExtraArgsCodec{}
}

func (ExtraArgsCodec) DecodeExtraArgs(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	// Not implemented and not return error
	return nil, nil
}

func (c ExtraArgsCodec) DecodeTokenAmountDestExecData(destExecData cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	//TODO implement me
	panic("implement me")
}
