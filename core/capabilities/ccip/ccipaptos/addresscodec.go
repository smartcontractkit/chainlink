package ccipaptos

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
)

type AddressCodec struct{}

func (a AddressCodec) AddressBytesToString(addr []byte) (string, error) {
	return addressBytesToString(addr)
}

func (a AddressCodec) AddressStringToBytes(addr string) ([]byte, error) {
	return addressStringToBytes(addr)
}

func addressBytesToString(addr []byte) (string, error) {
	if len(addr) != 32 {
		return "", fmt.Errorf("invalid Aptos address length, expected 32, got %d", len(addr))
	}

	accAddress := aptos.AccountAddress(addr)
	return accAddress.String(), nil
}

func addressStringToBytes(addr string) ([]byte, error) {
	var accAddress aptos.AccountAddress
	err := accAddress.ParseStringRelaxed(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Aptos address '%s': %w", addr, err)
	}
	return accAddress[:], nil
}

func addressStringToBytes32(addr string) ([32]byte, error) {
	var accAddress aptos.AccountAddress
	err := accAddress.ParseStringRelaxed(addr)
	if err != nil {
		return accAddress, fmt.Errorf("failed to decode Aptos address '%s': %w", addr, err)
	}
	return accAddress, nil
}

func addressBytesToBytes32(addr []byte) ([32]byte, error) {
	if len(addr) > 32 {
		return [32]byte{}, fmt.Errorf("invalid Aptos address length, expected 32, got %d", len(addr))
	}
	var result [32]byte
	// Left pad by copying to the end of the 32 byte array
	copy(result[32-len(addr):], addr)
	return result, nil
}

// takes a valid Aptos address string and converts it into canonical format.
func addressStringToString(addr string) (string, error) {
	var accAddress aptos.AccountAddress
	err := accAddress.ParseStringRelaxed(addr)
	if err != nil {
		return "", fmt.Errorf("failed to decode Aptos address '%s': %w", addr, err)
	}
	return accAddress.String(), nil
}

func addressIsValid(addr string) bool {
	_, err := addressStringToBytes(addr)
	return err == nil
}
