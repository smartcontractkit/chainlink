package view

import (
	"github.com/ethereum/go-ethereum/common"
)

type NonceManager struct {
	Contract
	AuthorizedCallers []common.Address `json:"authorizedCallers"`
}
