package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func GetMaxSizeFormEntry(n int, entry *types.CodecEntry) (int, error) {
	if entry == nil {
		return 0, relaytypes.InvalidEncodingError{}
	}
	return GetMaxSize(n, entry.Args)
}

func GetMaxSize(n int, args abi.Arguments) (int, error) {
	size := 0
	for _, arg := range args {
		argSize, err := getTypeSize(n, &arg.Type)
		if err != nil {
			return 0, err
		}
		size += argSize
	}

	return size, nil
}

const noSizeAllowed = -1

func getTypeSize(n int, t *abi.Type) (int, error) {
	// See https://docs.soliditylang.org/en/latest/abi-spec.html#formal-specification-of-the-encoding
	switch t.T {
	case abi.ArrayTy:
		elmSize, err := getTypeSize(noSizeAllowed, t.Elem)
		return t.Size * elmSize, err
	case abi.SliceTy:
		if noSizeAllowed == n {
			return 0, relaytypes.InvalidTypeError{}
		}
		elmSize, err := getTypeSize(noSizeAllowed, t.Elem)
		return 32 /*header*/ + 32 /*footer*/ + elmSize*n, err
	case abi.TupleTy:
		// No header or footer, because if the tuple is dynamically sized we would need to know the inner slice sizes
		// so it would return error for that element.
		size := 0
		for _, elm := range t.TupleElems {
			argSize, err := getTypeSize(noSizeAllowed, elm)
			if err != nil {
				return 0, err
			}
			size += argSize
		}
		return size, nil
	default:
		// types are padded to 32 bytes
		return 32, nil
	}
}
