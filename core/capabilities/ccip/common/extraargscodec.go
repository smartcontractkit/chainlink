package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

type ExtraArgsCodec struct{}

func NewExtraArgsCodec() ExtraArgsCodec {
	return ExtraArgsCodec{}
}

func (ExtraArgsCodec) DecodeExtraData(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	family, err := chainsel.GetSelectorFamily(uint64(sourceChainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain family for selector %d: %w", sourceChainSelector, err)
	}

	switch family {
	case chainsel.FamilyEVM:
		return ccipevm.DecodeExtraArgsToMap(extraArgs)

	case chainsel.FamilySolana:
		return ccipsolana.DecodeExtraArgsToMap(extraArgs)

	default:
		return nil, fmt.Errorf("unsupported family for extra args type %s", family)
	}
}
