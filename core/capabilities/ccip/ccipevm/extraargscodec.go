package ccipevm

import (
	"fmt"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

const evmExtraArgsKey = "gasLimit"

type ExtraArgsCodec struct{}

func NewExtraArgsCodec() ExtraArgsCodec {
	return ExtraArgsCodec{}
}

func (ExtraArgsCodec) DecodeExtraData(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	family, err := chain_selectors.GetSelectorFamily(uint64(sourceChainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to decode extra data, %w", err)
	}

	extraArgsMap := make(map[string]any)
	switch family {
	case chain_selectors.FamilyEVM:
		v2, err1 := decodeExtraArgsV1V2(extraArgs)
		if err1 != nil {
			return nil, err1
		}

		extraArgsMap[evmExtraArgsKey] = v2
	case chain_selectors.FamilySolana:
		// TODO add svm args type decoding logic once on-chain work is finished
	}

	return extraArgsMap, nil
}
