package v1_5

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
)

type CommitStoreView struct {
	types.ContractMetaData
	DynamicConfig              commit_store.CommitStoreDynamicConfig   `json:"dynamicConfig"`
	ExpectedNextSequenceNumber uint64                                  `json:"expectedNextSequenceNumber"`
	LatestPriceEpochAndRound   uint64                                  `json:"latestPriceEpochAndRound"`
	StaticConfig               commit_store.CommitStoreStaticConfig    `json:"staticConfig"`
	Transmitters               []common.Address                        `json:"transmitters"`
	IsUnpausedAndNotCursed     bool                                    `json:"isUnpausedAndNotCursed"`
	LatestConfigDetails        commit_store.LatestConfigDetails        `json:"latestConfigDetails"`
	LatestConfigDigestAndEpoch commit_store.LatestConfigDigestAndEpoch `json:"latestConfigDigestAndEpoch"`
	Paused                     bool                                    `json:"paused"`
}
