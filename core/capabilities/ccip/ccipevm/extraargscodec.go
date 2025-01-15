package ccipevm

import (
	"errors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type ExtraArgsCodec struct{}

func NewExtraArgsCodec() ExtraArgsCodec {
	return ExtraArgsCodec{}
}

func (ExtraArgsCodec) DecodeExtraData(ExtraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	return nil, errors.New("not implemented")
}
