package ocrkey

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

const onChainSigningAddressPrefix = "ocrsad_"

type OnChainSigningAddress ocrtypes.OnChainSigningAddress

func (ocsa OnChainSigningAddress) String() string {
	return fmt.Sprintf("%s%s", onChainSigningAddressPrefix, hexutil.Encode(ocsa[:]))
}

func (ocsa OnChainSigningAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(ocsa.String())
}

func (ocsa *OnChainSigningAddress) UnmarshalJSON(input []byte) error {
	var hexString string
	if err := json.Unmarshal(input, &hexString); err != nil {
		return err
	}
	return ocsa.UnmarshalText([]byte(hexString))
}

func (ocsa *OnChainSigningAddress) UnmarshalText(bs []byte) error {
	input := string(bs)
	if strings.HasPrefix(input, onChainSigningAddressPrefix) {
		input = string(bs[len(onChainSigningAddressPrefix):])
	}

	result, err := hexutil.Decode(input)
	if err != nil {
		return err
	}

	var onChainSigningAddress common.Address
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
