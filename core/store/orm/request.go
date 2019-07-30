package orm

import "github.com/smartcontractkit/chainlink/core/store/models"

// Request represents the values that can be changed in the configuration
type Request struct {
	EthGasPriceDefault models.Big `json:"ethGasPriceDefault"`
}
