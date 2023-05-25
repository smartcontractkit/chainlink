package s4

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"google.golang.org/protobuf/proto"
)

func MarshalRows(rows []*Row, addressRange *s4.AddressRange) ([]byte, error) {
	var protoAddressRange *AddressRange
	if addressRange != nil {
		minAddressStr, err := addressRange.MinAddress.MarshalText()
		if err != nil {
			return nil, err
		}
		maxAddressStr, err := addressRange.MaxAddress.MarshalText()
		if err != nil {
			return nil, err
		}
		protoAddressRange = &AddressRange{
			MinAddress: string(minAddressStr),
			MaxAddress: string(maxAddressStr),
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
	var addressRange *s4.AddressRange
	if rows.AddressRange != nil {
		minAddress := new(utils.Big)
		maxAddress := new(utils.Big)
		if rows.AddressRange.MinAddress != "" {
			if err := minAddress.UnmarshalText([]byte(rows.AddressRange.MinAddress)); err != nil {
				return nil, nil, err
			}
		}
		if rows.AddressRange.MaxAddress != "" {
			if err := maxAddress.UnmarshalText([]byte(rows.AddressRange.MaxAddress)); err != nil {
				return nil, nil, err
			}
		}
		addressRange = &s4.AddressRange{
			MinAddress: minAddress,
			MaxAddress: maxAddress,
		}
	}
	return rows.Rows, addressRange, nil
}
