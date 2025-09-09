package ccipton

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/sigurn/crc16"
	"github.com/xssnick/tonutils-go/address"
)

type AddressCodec struct{}

func (a AddressCodec) AddressBytesToString(addr []byte) (string, error) {
	decoded := base64.RawURLEncoding.EncodeToString(addr)
	// verify that the TON address string is valid
	tonAddr, err := address.ParseAddr(decoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode TVM address bytes: %w", err)
	}

	return tonAddr.String(), nil
}

func (a AddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	// underneath implementation of TON address is base64 Raw URL encoding, reference: tonutils-go/address ParseAddr function
	decodeString, err := base64.RawURLEncoding.DecodeString(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode TVM address bytes: %w", err)
	}

	if len(decodeString) != 36 {
		return nil, fmt.Errorf("invalid address length: expected 36 bytes, got %d", len(decodeString))
	}

	checksum := decodeString[len(decodeString)-2:]
	if crc16.Checksum(decodeString[:len(decodeString)-2], crc16.MakeTable(crc16.CRC16_XMODEM)) != binary.BigEndian.Uint16(checksum) {
		return nil, fmt.Errorf("invalid checksum for address: %s", addr)
	}

	return decodeString, nil
}

func (a AddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	addr := make([]byte, 32)
	// write oracleID into addr in big endian
	binary.BigEndian.PutUint32(addr, uint32(oracleID))
	tonAddr := address.NewAddress(0, 0, addr)
	decodeString, err := base64.RawURLEncoding.DecodeString(tonAddr.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode TVM address bytes: %w", err)
	}

	return decodeString, nil
}

func (a AddressCodec) TransmitterBytesToString(addr []byte) (string, error) {
	// Transmitter accounts are addresses
	return a.AddressBytesToString(addr)
}
