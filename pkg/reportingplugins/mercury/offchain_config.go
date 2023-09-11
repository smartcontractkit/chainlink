package mercury

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type OffchainConfig struct {
	ExpirationWindow uint32          // Integer number of seconds
	BaseUSDFee       decimal.Decimal // Base USD fee
}

func DecodeOffchainConfig(b []byte) (o OffchainConfig, err error) {
	err = json.Unmarshal(b, &o)
	return
}

func (c OffchainConfig) Encode() ([]byte, error) {
	return json.Marshal(c)
}
