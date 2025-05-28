package ccipsolana

import (
	"encoding/binary"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

type AddressCodec struct{}

func (a AddressCodec) AddressBytesToString(addr []byte) (string, error) {
	if len(addr) != solana.PublicKeyLength {
		return "", fmt.Errorf("invalid SVM address length, expected %v, got %d", solana.PublicKeyLength, len(addr))
	}
	return solana.PublicKeyFromBytes(addr).String(), nil
}

func (a AddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	pk, err := solana.PublicKeyFromBase58(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode SVM address '%s': %w", addr, err)
	}
	return pk.Bytes(), nil
}

func (a AddressCodec) OracleIDAsAddressBytes(oracleID uint8) ([]byte, error) {
	addr := make([]byte, solana.PublicKeyLength)

	// write oracleID into addr in little endian, since solana is little endian
	binary.LittleEndian.PutUint32(addr, uint32(oracleID))

	// TODO: is it alright if the pub key is off the curve?
	return solana.PublicKeyFromBytes(addr).Bytes(), nil
}
