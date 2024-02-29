package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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

type VRFKeyRegistrationConfig struct {
	VRFKeyUncompressedPubKey string
	RegisterAgainstAddress   string
}

type CoordinatorJobSpecConfig struct {
	BatchFulfillmentEnabled       bool
	BatchFulfillmentGasMultiplier float64
	EstimateGasMultiplier         float64
	PollPeriod                    string
	RequestTimeout                string
	RevertsPipelineEnabled        bool
}

type BHSJobSpecConfig struct {
	RunTimeout     string
	WaitBlocks     int
	LookBackBlocks int
	PollPeriod     string
	RequestTimeout string
}
