package s4

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"
)

func MarshalQuery(rows []*SnapshotRow) ([]byte, error) {
	rr := &Query{
		Rows: rows,
	}
	return proto.Marshal(rr)
}

func UnmarshalQuery(data []byte) ([]*SnapshotRow, error) {
	query := &Query{}
	if err := proto.Unmarshal(data, query); err != nil {
		return nil, err
	}
	if query.Rows == nil {
		query.Rows = make([]*SnapshotRow, 0)
	}
	return query.Rows, nil
}

func MarshalRows(rows []*Row) ([]byte, error) {
	rr := &Rows{
		Rows: rows,
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

func UnmarshalAddress(address []byte) *utils.Big {
	return utils.NewBig(new(big.Int).SetBytes(address))
}

func (row *Row) VerifySignature() error {
	address := common.BytesToAddress(row.Address)
	e := &s4.Envelope{
		Address:    address.Bytes(),
		SlotID:     uint(row.Slotid),
		Payload:    row.Payload,
		Version:    row.Version,
		Expiration: row.Expiration,
	}
	signer, err := e.GetSignerAddress(row.Signature)
	if err != nil {
		return err
	}
	if signer != address {
		return s4.ErrWrongSignature
	}
	return nil
}
