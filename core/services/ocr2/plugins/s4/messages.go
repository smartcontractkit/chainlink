package s4

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"google.golang.org/protobuf/proto"
)

func MarshalRows(rows []*Row, addressRange *s4.AddressRange) ([]byte, error) {
	var protoAddressRange *AddressRange
	if addressRange != nil {
		minAddressStr, err := MarshalAddress(addressRange.MinAddress)
		if err != nil {
			return nil, err
		}
		maxAddressStr, err := MarshalAddress(addressRange.MaxAddress)
		if err != nil {
			return nil, err
		}
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
	var addressRange *s4.AddressRange
	if rows.AddressRange != nil {
		minAddress, err := UnmarshalAddress(rows.AddressRange.MinAddress)
		if err != nil {
			return nil, nil, err
		}
		maxAddress, err := UnmarshalAddress(rows.AddressRange.MaxAddress)
		if err != nil {
			return nil, nil, err
		}
		addressRange = &s4.AddressRange{
			MinAddress: minAddress,
			MaxAddress: maxAddress,
		}
	}
	return rows.Rows, addressRange, nil
}

func MarshalAddress(address *utils.Big) (string, error) {
	addressStr, err := address.MarshalText()
	if err != nil {
		return "", err
	}
	return string(addressStr), nil
}

func UnmarshalAddress(address string) (*utils.Big, error) {
	bigAddress := new(utils.Big)
	if err := bigAddress.UnmarshalText([]byte(address)); err != nil {
		return nil, err
	}
	return bigAddress, nil
}
