package s4

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"

	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
)

func MarshalQuery(rows []*SnapshotRow, addressRange *s4.AddressRange) ([]byte, error) {
	rr := &Query{
		AddressRange: &AddressRange{
			MinAddress: addressRange.MinAddress.Bytes(),
			MaxAddress: addressRange.MaxAddress.Bytes(),
		},
		Rows: rows,
	}
	return proto.Marshal(rr)
}

func UnmarshalQuery(data []byte) ([]*SnapshotRow, *s4.AddressRange, error) {
	addressRange := s4.NewFullAddressRange()
	query := &Query{}
	if err := proto.Unmarshal(data, query); err != nil {
		return nil, nil, err
	}
	if query.Rows == nil {
		query.Rows = make([]*SnapshotRow, 0)
	}
	if query.AddressRange != nil {
		addressRange = &s4.AddressRange{
			MinAddress: UnmarshalAddress(query.AddressRange.MinAddress),
			MaxAddress: UnmarshalAddress(query.AddressRange.MaxAddress),
		}
	}
	return query.Rows, addressRange, nil
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

func UnmarshalAddress(address []byte) *ubig.Big {
	return ubig.New(new(big.Int).SetBytes(address))
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
	if !bytes.Equal(signer.Bytes(), address.Bytes()) {
		return s4.ErrWrongSignature
	}
	return nil
}
