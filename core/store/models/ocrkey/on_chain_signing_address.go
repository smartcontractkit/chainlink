package ocrkey

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

type OnChainSigningAddress ocrtypes.OnChainSigningAddress

func (ocsa OnChainSigningAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(hexutil.Encode(ocsa[:]))
}

func (ocsa *OnChainSigningAddress) UnmarshalJSON(input []byte) error {
	var hexString string
	var onChainSigningAddress common.Address
	if err := json.Unmarshal(input, &hexString); err != nil {
		return err
	}

	result, err := hexutil.Decode(hexString)
	if err != nil {
		return err
	}

	copy(onChainSigningAddress[:], result[:common.AddressLength])
	*ocsa = OnChainSigningAddress(onChainSigningAddress)
	return nil
}

func (ocsa OnChainSigningAddress) Value() (driver.Value, error) {
	byteArray := [common.AddressLength]byte(ocsa)
	return byteArray[:], nil
}

func (ocsa *OnChainSigningAddress) Scan(value interface{}) error {
	switch typed := value.(type) {
	case []byte:
		if len(typed) != common.AddressLength {
			return errors.New("wrong number of bytes to scan into address")
		}
		copy(ocsa[:], typed)
		return nil
	default:
		return errors.Errorf(`unable to convert %v of %T to OnChainSigningAddress`, value, value)
	}
}
