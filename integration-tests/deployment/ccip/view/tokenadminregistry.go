package view

import (
	"github.com/ethereum/go-ethereum/common"
)

type TokenAdminRegistry struct {
	Contract
	Tokens []common.Address `json:"tokens"`
}
