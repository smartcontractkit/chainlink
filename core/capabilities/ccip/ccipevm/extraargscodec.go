package ccipevm

import (
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type ExtraArgsCodec struct{}

func NewExtraArgsCodec() ExtraArgsCodec {
	return ExtraArgsCodec{}
}

func (ExtraArgsCodec) DecodeExtraData(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	// Not implemented and not return error
	return nil, nil
}
