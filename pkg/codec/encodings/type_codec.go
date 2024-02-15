package encodings

import (
	"context"
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type TypeCodec interface {
	Encode(value any, into []byte) ([]byte, error)
	Decode(encoded []byte) (any, []byte, error)
	GetType() reflect.Type

	// Size returns the size of the encoded value in bytes if there are N reports.
	// As such, any nested elements should be called with FixedSize to determine their size, unless it's implicit how
	// how many nested items there will be.
	// As an example, a struct { A: []int } should return the size of numItems ints, but a struct { A: [][]int }
	// should return an error, as each report could have a different number of elements in their slice.
	Size(numItems int) (int, error)

	// FixedSize returns the size of the encoded value, without providing a count of elements.
	// If a count of elements is required to know the size, an error must be returned.
	FixedSize() (int, error)
}

// TopLevelCodec is a TypeCodec that can be encoded at the top level of a report.
// This allows each member to be called with Size(numItems) when SiteAtTopLevel(numItems) is called.
type TopLevelCodec interface {
	TypeCodec
	SizeAtTopLevel(numItems int) (int, error)
}

// CodecFromTypeCodec maps TypeCodec to types.RemoteCodec, using the key as the itemType
// If the TypeCodec is a TopLevelCodec, GetMaxEncodingSize and GetMaxDecodingSize will call SizeAtTopLevel instead of SIze.
type CodecFromTypeCodec map[string]TypeCodec

var _ types.RemoteCodec = &CodecFromTypeCodec{}

func (c CodecFromTypeCodec) CreateType(itemType string, _ bool) (any, error) {
	ntcwt, ok := c[itemType]
	if !ok {
		return nil, fmt.Errorf("%w: cannot find type %s", types.ErrInvalidType, itemType)
	}

	tpe := ntcwt.GetType()
	if tpe.Kind() == reflect.Pointer {
		tpe = tpe.Elem()
	}

	return reflect.New(tpe).Interface(), nil
}

func (c CodecFromTypeCodec) Encode(_ context.Context, item any, itemType string) ([]byte, error) {
	ntcwt, ok := c[itemType]
	if !ok {
		return nil, fmt.Errorf("%w: cannot find type %s", types.ErrInvalidType, itemType)
	}

	if item != nil {
		rItem := reflect.ValueOf(item)
		myType := ntcwt.GetType()
		if rItem.Kind() == reflect.Pointer && myType.Kind() != reflect.Pointer {
			rItem = reflect.Indirect(rItem)
		}

		if !rItem.IsZero() && rItem.Type() != myType {
			tmp := reflect.New(myType)
			if err := codec.Convert(rItem, tmp, nil); err != nil {
				return nil, err
			}
			item = tmp.Elem().Interface()
		} else {
			item = rItem.Interface()
		}
	}

	return ntcwt.Encode(item, nil)
}

func (c CodecFromTypeCodec) GetMaxEncodingSize(_ context.Context, n int, itemType string) (int, error) {
	ntcwt, ok := c[itemType]
	if !ok {
		return 0, fmt.Errorf("%w: cannot find type %s", types.ErrInvalidType, itemType)
	}

	if lp, ok := ntcwt.(TopLevelCodec); ok {
		return lp.SizeAtTopLevel(n)
	}
	return ntcwt.Size(n)
}

func (c CodecFromTypeCodec) Decode(_ context.Context, raw []byte, into any, itemType string) error {
	ntcwt, ok := c[itemType]
	if !ok {
		return fmt.Errorf("%w: cannot find type %s", types.ErrInvalidType, itemType)
	}
	val, remaining, err := ntcwt.Decode(raw)
	if err != nil {
		return err
	}

	if len(remaining) != 0 {
		return fmt.Errorf("%w: remaining bytes after decoding %s", types.ErrInvalidEncoding, itemType)
	}

	if err = codec.Convert(reflect.ValueOf(val), reflect.ValueOf(into), nil); err != nil {
		return fmt.Errorf("%w: %v", types.ErrInvalidType, err)
	}

	return nil
}

func (c CodecFromTypeCodec) GetMaxDecodingSize(ctx context.Context, n int, itemType string) (int, error) {
	return c.GetMaxEncodingSize(ctx, n, itemType)
}
