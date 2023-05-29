package s4

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"google.golang.org/protobuf/proto"
)

func MarshalRows(rows []*Row, addressRange *s4.AddressRange) ([]byte, error) {
	var protoAddressRange *AddressRange
	if addressRange != nil {
		minAddressStr := MarshalAddress(addressRange.MinAddress)
		maxAddressStr := MarshalAddress(addressRange.MaxAddress)
		protoAddressRange = &AddressRange{
			MinAddress: minAddressStr,
			MaxAddress: maxAddressStr,
		}
	}
	rr := &Rows{
		Rows:         rows,
		AddressRange: protoAddressRange,
	}
	return proto.Marshal(rr)
}

func UnmarshalRows(data []byte) ([]*Row, *s4.AddressRange, error) {
	rows := &Rows{}
	if err := proto.Unmarshal(data, rows); err != nil {
		return nil, nil, err
	}
	addressRange, err := UnmarshalAddressRange(rows.AddressRange)
	if err != nil {
		return nil, nil, err
	}
	if rows.Rows == nil {
		rows.Rows = make([]*Row, 0)
	}
	return rows.Rows, addressRange, nil
}

func UnmarshalQuery(data []byte) ([]*VersionRow, *s4.AddressRange, error) {
	query := &Query{}
	if err := proto.Unmarshal(data, query); err != nil {
		return nil, nil, err
	}
	addressRange, err := UnmarshalAddressRange(query.AddressRange)
	if err != nil {
		return nil, nil, err
	}
	if query.Versions == nil {
		query.Versions = make([]*VersionRow, 0)
	}
	return query.Versions, addressRange, nil
}

func UnmarshalAddressRange(addressRange *AddressRange) (*s4.AddressRange, error) {
	if addressRange == nil {
		return nil, nil
	}

	var ormAddressRange *s4.AddressRange
	minAddress, err := UnmarshalAddress(addressRange.MinAddress)
	if err != nil {
		return nil, err
	}
	maxAddress, err := UnmarshalAddress(addressRange.MaxAddress)
	if err != nil {
		return nil, err
	}
	ormAddressRange = &s4.AddressRange{
		MinAddress: minAddress,
		MaxAddress: maxAddress,
	}
	return ormAddressRange, nil
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
