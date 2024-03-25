package chainreader_test

import (
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	. "github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
)

var errorTypes = []error{
	types.ErrInvalidEncoding,
	types.ErrInvalidType,
	types.ErrFieldNotFound,
	types.ErrSliceWrongLen,
	types.ErrNotASlice,
	types.ErrNotFound,
}

type cannotEncode struct{}

func (*cannotEncode) MarshalCBOR() ([]byte, error) {
	return nil, errors.New("nope")
}

func (*cannotEncode) UnmarshalCBOR([]byte) error {
	return errors.New("nope")
}

func (*cannotEncode) MarshalText() ([]byte, error) {
	return nil, errors.New("nope")
}

func (*cannotEncode) UnmarshalText() error {
	return errors.New("nope")
}

type interfaceTesterBase struct{}

var anyAccountBytes = []byte{1, 2, 3}

func (it *interfaceTesterBase) GetAccountBytes(_ int) []byte {
	return anyAccountBytes
}

func (it *interfaceTesterBase) Name() string {
	return "relay client"
}

type fakeTypeProvider struct{}

func (f fakeTypeProvider) CreateType(itemType string, isEncode bool) (any, error) {
	return f.CreateContractType("", itemType, isEncode)
}

func (fakeTypeProvider) CreateContractType(_, itemType string, isEncode bool) (any, error) {
	switch itemType {
	case NilType:
		return &struct{}{}, nil
	case TestItemType:
		return &TestStruct{}, nil
	case TestItemSliceType:
		return &[]TestStruct{}, nil
	case TestItemArray2Type:
		return &[2]TestStruct{}, nil
	case TestItemArray1Type:
		return &[1]TestStruct{}, nil
	case MethodTakingLatestParamsReturningTestStruct:
		if isEncode {
			return &LatestParams{}, nil
		}
		return &TestStruct{}, nil
	case MethodReturningUint64, DifferentMethodReturningUint64:
		tmp := uint64(0)
		return &tmp, nil
	case MethodReturningUint64Slice:
		var tmp []uint64
		return &tmp, nil
	case MethodReturningSeenStruct, TestItemWithConfigExtra:
		if isEncode {
			return &TestStruct{}, nil
		}
		return &TestStructWithExtraField{}, nil
	case EventName, EventWithFilterName:
		if isEncode {
			return &FilterEventParams{}, nil
		}
		return &TestStruct{}, nil
	}

	return nil, types.ErrInvalidType
}
