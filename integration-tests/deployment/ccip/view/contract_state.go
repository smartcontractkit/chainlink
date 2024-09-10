package view

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type ContractStatus string

const (
	Active          ContractStatus = "active"
	Inactive        ContractStatus = "inactive"
	Decommissioning ContractStatus = "decommissioning"
	Dead            ContractStatus = "dead"
)

// TODO : Should this denote blue-green state?
var ContractStatusLookup = map[string]ContractStatus{
	"active":          Active,
	"inactive":        Inactive,
	"decommissioning": Decommissioning,
	"dead":            Dead,
}

type Contract struct {
	TypeAndVersion string `json:"typeAndVersion,omitempty"`
	Address        string `json:"address,omitempty"`
}

type ContractState interface {
	TypeAndVersion(opts *bind.CallOpts) (string, error)
	Address() common.Address
}
