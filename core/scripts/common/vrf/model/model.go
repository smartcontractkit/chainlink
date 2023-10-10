package model

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	VRFPrimaryNodeName = "vrf-primary-node"
	VRFBackupNodeName  = "vrf-backup-node"
	BHSNodeName        = "bhs-node"
	BHSBackupNodeName  = "bhs-backup-node"
	BHFNodeName        = "bhf-node"
)

type Node struct {
	URL                         string
	CredsFile                   string
	SendingKeys                 []SendingKey
	NumberOfSendingKeysToCreate int
	SendingKeyFundingAmount     *big.Int
	VrfKeys                     []string
	jobSpec                     string
}

type SendingKey struct {
	Address    string
	BalanceEth *big.Int
}

type JobSpecs struct {
	VRFPrimaryNode string
	VRFBackupyNode string
	BHSNode        string
	BHSBackupNode  string
	BHFNode        string
}

type ContractAddresses struct {
	LinkAddress             string
	LinkEthAddress          string
	BhsContractAddress      common.Address
	BatchBHSAddress         common.Address
	CoordinatorAddress      common.Address
	BatchCoordinatorAddress common.Address
}
