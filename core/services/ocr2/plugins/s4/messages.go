package s4

import (
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"google.golang.org/protobuf/proto"
)

func MarshalRows(rows []*Row) ([]byte, error) {
	rr := &Rows{
		Rows: rows,
	}
	return proto.Marshal(rr)
}

func MarshalQuery(versions []*VersionRow) ([]byte, error) {
	rr := &Query{
		Versions: versions,
	}
	return proto.Marshal(rr)
}

func UnmarshalRows(data []byte) ([]*Row, error) {
	rows := &Rows{}
	if err := proto.Unmarshal(data, rows); err != nil {
		return nil, err
	}
	if rows.Rows == nil {
		rows.Rows = make([]*Row, 0)
	}
	return rows.Rows, nil
}

func UnmarshalQuery(data []byte) ([]*VersionRow, error) {
	query := &Query{}
	if err := proto.Unmarshal(data, query); err != nil {
		return nil, err
	}
	if query.Versions == nil {
		query.Versions = make([]*VersionRow, 0)
	}
	return query.Versions, nil
}

func MarshalAddress(address *utils.Big) string {
	return address.Hex()
}

func UnmarshalAddress(address string) (*utils.Big, error) {
	decoded, err := hexutil.DecodeBig(address)
	if err != nil {
		return nil, err
	}
	return utils.NewBig(decoded), nil
}
