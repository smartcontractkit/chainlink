package view

import (
	"github.com/ethereum/go-ethereum/common"
)

type ContractStatus string

const (
	Active          ContractStatus = "active"
	Inactive        ContractStatus = "inactive"
	Decommissioning ContractStatus = "decommissioning"
	Dead            ContractStatus = "dead"
)

var ContractStatusLookup = map[string]ContractStatus{
	"active":          Active,
	"inactive":        Inactive,
	"decommissioning": Decommissioning,
	"dead":            Dead,
}

type Contract struct {
	TypeAndVersion string         `json:"typeAndVersion"`
	Address        common.Address `json:"address"`
}
