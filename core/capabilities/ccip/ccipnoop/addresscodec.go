package ccipnoop

type AddressCodec struct{}

func (n AddressCodec) AddressBytesToString(addr []byte) (string, error) {
	return string(addr), nil
}

func (n AddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	return []byte(addr), nil
}

func (n AddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	return []byte{oracleID}, nil
}

func (n AddressCodec) TransmitterBytesToString(addr []byte) (string, error) {
	return string(addr), nil
}
