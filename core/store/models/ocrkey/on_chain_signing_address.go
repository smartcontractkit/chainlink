package ocrkey

import (
	"database/sql/driver"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type OnChainSigningAddress common.Address

func (address OnChainSigningAddress) Value() (driver.Value, error) {
	byteArray := [common.AddressLength]byte(address)
	return byteArray[:], nil
}

func (address *OnChainSigningAddress) Scan(value interface{}) error {
	switch typed := value.(type) {
	case []byte:
		if len(typed) != common.AddressLength {
			return errors.New("wrong number of bytes to scan into address")
		}
		copy(address[:], typed)
		return nil
	default:
		return errors.Errorf(`unable to convert %v of %T to OnChainSigningAddress`, value, value)
	}
}
