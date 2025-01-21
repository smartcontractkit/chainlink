package ccipevm

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

const (
	EVMExtraArgsKey     = "gasLimit"
	SVMComputeUnitKey   = "ComputeUnits"
	SVMAccountBitmapKey = "AccountIsWritableBitmap"
	SVMAllowOOEKey      = "AllowOutOfOrderExecution"
	SVMTokenReceiverKey = "TokenReceiver"
	SVMAccountListKey   = "Accounts"
)

type ExtraArgsCodec struct{}

func NewExtraArgsCodec() ExtraArgsCodec {
	return ExtraArgsCodec{}
}

func (ExtraArgsCodec) DecodeExtraData(extraArgs cciptypes.Bytes, sourceChainSelector cciptypes.ChainSelector) (map[string]any, error) {
	family, err := chainsel.GetSelectorFamily(uint64(sourceChainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to decode extra data, %w", err)
	}

	switch family {
	case chainsel.FamilyEVM:
		gas, err1 := decodeExtraArgsV1V2(extraArgs)
		if err1 != nil {
			return nil, fmt.Errorf("failed to decode EVM extra data, %w", err)
		}

		return map[string]any{
			EVMExtraArgsKey: gas,
		}, nil

	case chainsel.FamilySolana:
		v1, err1 := decodeExtraArgsSVMV1(extraArgs)
		if err1 != nil {
			return nil, fmt.Errorf("failed to decode SVM extra data, %w", err)
		}

		return map[string]any{
			SVMComputeUnitKey:   v1.ComputeUnits,
			SVMAccountBitmapKey: v1.AccountIsWritableBitmap,
			SVMAllowOOEKey:      v1.AllowOutOfOrderExecution,
			SVMTokenReceiverKey: v1.TokenReceiver,
			SVMAccountListKey:   v1.Accounts,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported family for extra args type %s", family)
	}
}
