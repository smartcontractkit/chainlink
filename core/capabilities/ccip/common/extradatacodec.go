package common

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipsolana"
)

type ExtraDataCodec struct{}

func NewExtraDataCodec() ExtraDataCodec {
	return ExtraDataCodec{}
}

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
		return ccipevm.DecodeExtraArgsToMap(extraArgs)

	case chainsel.FamilySolana:
		return ccipsolana.DecodeExtraArgsToMap(extraArgs)

	default:
		return nil, fmt.Errorf("unsupported family for extra args type %s", family)
	}
}

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
		return ccipevm.DecodeDestExecDataToMap(destExecData)

	case chainsel.FamilySolana:
		return ccipsolana.DecodeDestExecDataToMap(destExecData)

	default:
		return nil, fmt.Errorf("unsupported family for extra args type %s", family)
	}
}
