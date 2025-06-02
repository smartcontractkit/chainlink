package ccipnoop

type NoopAddressCodec struct{}

func (n NoopAddressCodec) AddressBytesToString(addr []byte) (string, error) {
	return string(addr), nil
}

func (n NoopAddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	return []byte(addr), nil
}

func (n NoopAddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	return []byte{oracleID}, nil
}
