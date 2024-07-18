package evm

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

// SendEtherRequest represents a request to transfer ETH.
type SendEtherRequest struct {
	DestinationAddress common.Address `json:"address"`
	FromAddress        common.Address `json:"from"`
	Amount             assets.Eth     `json:"amount"`
	EVMChainID         *big.Big       `json:"evmChainID"`
	AllowHigherAmounts bool           `json:"allowHigherAmounts"`
	SkipWaitTxAttempt  bool           `json:"skipWaitTxAttempt"`
	WaitAttemptTimeout *time.Duration `json:"waitAttemptTimeout"`
}

// AddressCollection is an array of common.Address
// serializable to and from a database.
type AddressCollection []common.Address

// ToStrings returns this address collection as an array of strings.
func (r AddressCollection) ToStrings() []string {
	// Unable to convert copy-free without unsafe:
	// https://stackoverflow.com/a/48554123/639773
	converted := make([]string, len(r))
	for i, e := range r {
		converted[i] = e.Hex()
	}
	return converted
}

// Value returns the string value to be written to the database.
func (r AddressCollection) Value() (driver.Value, error) {
	return strings.Join(r.ToStrings(), ","), nil
}

// Scan parses the database value as a string.
func (r *AddressCollection) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("unable to convert %v of %T to AddressCollection", value, value)
	}

	if len(str) == 0 {
		return nil
	}

	arr := strings.Split(str, ",")
	collection := make(AddressCollection, len(arr))
	for i, a := range arr {
		collection[i] = common.HexToAddress(a)
	}
	*r = collection
	return nil
}
