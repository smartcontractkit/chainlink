package binary

import (
	"fmt"
	"reflect"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/codec/encodings"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type OracleID Uint8

var _ encodings.TypeCodec = &OracleID{}

func (o *OracleID) Encode(value any, into []byte) ([]byte, error) {
	if oid, ok := value.(commontypes.OracleID); ok {
		return (*Uint8)(o).Encode(uint8(oid), into)
	}

	return nil, fmt.Errorf("%w: %v", types.ErrInvalidType, reflect.TypeOf(value))
}

func (o *OracleID) Decode(encoded []byte) (any, []byte, error) {
	decoded, remaining, err := (*Uint8)(o).Decode(encoded)
	return commontypes.OracleID(decoded.(uint8)), remaining, err
}

func (o *OracleID) GetType() reflect.Type {
	return reflect.TypeOf(commontypes.OracleID(0))
}

func (o *OracleID) Size(_ int) (int, error) {
	return 1, nil
}

func (o *OracleID) FixedSize() (int, error) {
	return 1, nil
}
