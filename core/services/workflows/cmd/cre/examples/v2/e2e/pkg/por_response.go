package pkg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type RawReserveInfo struct {
	LastUpdated  time.Time       `json:"lastUpdated"`
	TotalReserve decimal.Decimal `json:"totalReserve"`
}

type ReserveInfo struct {
	LastUpdated  int64           `consensus:"median"`
	TotalReserve decimal.Decimal `consensus:"median"`
}

type PorResponse struct {
	DataSignature string `json:"dataSignature"`
	Ripcord       bool   `json:"ripcord"`
	Data          string `json:"data"`
}

func UnescapeJSONString(input string) (string, error) {
	var unescaped string
	err := json.Unmarshal([]byte(`"`+input+`"`), &unescaped)
	if err != nil {
		return "", fmt.Errorf("failed to unescape string: %w", err)
	}
	return unescaped, nil
}
