package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	type_and_version "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/type_and_version_interface_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/confirmed_owner_with_proposal"
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

const (
	RouterTypeAndVersionV1_2             = "Router 1.2.0"
	TokenAdminRegistryTypeAndVersionV1_5 = "TokenAdminRegistry 1.5.0"
	FEEQuoterTypeAndVersionV1_6          = "FeeQuoter 1.6.0-dev"
)

type ContractMetaData struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Address        common.Address `json:"address,omitempty"`
	Owner          common.Address `json:"owner,omitempty"`
}

func (c *ContractMetaData) SetOwner(owner common.Address) {
	c.Owner = owner
}

func (c *ContractMetaData) Validate() error {
	if c.TypeAndVersion == "" {
		return fmt.Errorf("type and version is required")
	}
	if c.Address == (common.Address{}) {
		return fmt.Errorf("address is required")
	}
	return nil
}

func NewContractMetaData(address common.Address, client bind.ContractBackend) (ContractMetaData, error) {
	tv, err := type_and_version.NewTypeAndVersionInterface(address, client)
	if err != nil {
		return ContractMetaData{}, fmt.Errorf("failed to get type and version for contract %s: %w", address, err)
	}
	tvStr, err := tv.TypeAndVersion(nil)
	if err != nil {
		return ContractMetaData{}, err
	}

	co, err := confirmed_owner_with_proposal.NewConfirmedOwnerWithProposal(address, client)
	if err != nil {
		return ContractMetaData{}, fmt.Errorf("failed to get owner for contract %s: %w", address, err)
	}
	ownerAddr, err := co.Owner(nil)
	if err != nil {
		return ContractMetaData{}, fmt.Errorf("failed to call owner method for contract %s: %w", address, err)
	}

	return ContractMetaData{
		TypeAndVersion: tvStr,
		Address:        address,
		Owner:          ownerAddr,
	}, nil
}

type Snapshotter interface {
	Snapshot(contractMeta ContractMetaData, dependenciesMeta []ContractMetaData, client bind.ContractBackend) error
}
