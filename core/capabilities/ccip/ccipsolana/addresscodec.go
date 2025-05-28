package ccipsolana

import (
	"crypto/rand"
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

func (a AddressCodec) RandomAddressBytes() ([]byte, error) {
	addr := make([]byte, solana.PublicKeyLength)
	_, err := rand.Read(addr)
	if err != nil {
		return nil, err
	}
	// TODO: is it alright if the pub key is off the curve?
	return solana.PublicKeyFromBytes(addr).Bytes(), nil
}
