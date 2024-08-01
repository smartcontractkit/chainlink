package mercury

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

type OffchainConfig struct {
	ExpirationWindow uint32          `json:"expirationWindow"` // Integer number of seconds
	BaseUSDFee       decimal.Decimal `json:"baseUSDFee"`       // Base USD fee
}

func DecodeOffchainConfig(b []byte) (o OffchainConfig, err error) {
	err = json.Unmarshal(b, &o)
	if err != nil {
		return o, fmt.Errorf("failed to decode offchain config: must be valid JSON (got: 0x%x); %w", b, err)
	}
	return
}

func (c OffchainConfig) Encode() ([]byte, error) {
	return json.Marshal(c)
}
