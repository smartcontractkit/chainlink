package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
)

func GetMaxSizeFormEntry(n int, entry *CodecEntry) (int, error) {
	if entry == nil {
		return 0, relaytypes.InvalidTypeError{}
	}
	return GetMaxSize(n, entry.Args)
}

func GetMaxSize(n int, args abi.Arguments) (int, error) {
	size := 0
	for _, arg := range args {
		tmp := arg.Type
		argSize, _, err := getTypeSize(n, &tmp, true, false)
		if err != nil {
			return 0, err
		}
		size += argSize
	}

	return size, nil
}

func getTypeSize(n int, t *abi.Type, dynamicTypeAllowed bool, isNested bool) (int, bool, error) {
	// See https://docs.soliditylang.org/en/latest/abi-spec.html#formal-specification-of-the-encoding
	switch t.T {
	case abi.ArrayTy:
		elmSize, _, err := getTypeSize(n, t.Elem, false, true)
		return t.Size * elmSize, false, err
	case abi.SliceTy:
		if !dynamicTypeAllowed {
			return 0, false, relaytypes.InvalidTypeError{}
		}
		elmSize, _, err := getTypeSize(n, t.Elem, false, true)
		return 32 /*header*/ + 32 /*footer*/ + elmSize*n, true, err
	case abi.BytesTy, abi.StringTy:
		if !dynamicTypeAllowed {
			return 0, false, relaytypes.InvalidTypeError{}
		}
		totalSize := (n + 31) / 32 * 32 // strings and bytes are padded to 32 bytes
		return 32 /*header*/ + 32 /*footer*/ + totalSize, true, nil
	case abi.TupleTy:
		// No header or footer, because if the tuple is dynamically sized we would need to know the inner slice sizes
		// so it would return error for that element.
		size := 0
		dynamic := false
		for _, elm := range t.TupleElems {
			argSize, dynamicArg, err := getTypeSize(n, elm, !isNested, true)
			if err != nil {
				return 0, false, err
			}
			dynamic = dynamic || dynamicArg
			size += argSize
		}

		if dynamic {
			// offset for the element needs to be included there are dynamic elements
			size += 32
		}

		return size, dynamic, nil
	default:
		// types are padded to 32 bytes
		return 32, false, nil
	}
}
