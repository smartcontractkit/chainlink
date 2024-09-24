package v1_2

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip/view/types"
	commit_store "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_1_2_0"
)

type CommitStoreView struct {
	types.ContractMetaData
	DynamicConfig              commit_store.CommitStoreDynamicConfig   `json:"dynamicConfig"`
	ExpectedNextSequenceNumber uint64                                  `json:"expectedNextSequenceNumber"`
	LatestPriceEpochAndRound   uint64                                  `json:"latestPriceEpochAndRound"`
	StaticConfig               commit_store.CommitStoreStaticConfig    `json:"staticConfig"`
	Transmitters               []common.Address                        `json:"transmitters"`
	IsARMHealthy               bool                                    `json:"isARMHealthy"`
	IsUnpausedAndARMHealthy    bool                                    `json:"isUnpausedAndARMHealthy"`
	LatestConfigDetails        commit_store.LatestConfigDetails        `json:"latestConfigDetails"`
	LatestConfigDigestAndEpoch commit_store.LatestConfigDigestAndEpoch `json:"latestConfigDigestAndEpoch"`
	Paused                     bool                                    `json:"paused"`
}

func GenerateCommitStoreView(c *commit_store.CommitStore) (CommitStoreView, error) {
	if c == nil {
		return CommitStoreView{}, fmt.Errorf("cannot generate view for nil CommitStore")
	}
	meta, err := types.NewContractMetaData(c, c.Address())
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to generate contract metadata for CommitStore: %w", err)
	}
	dynamicConfig, err := c.GetDynamicConfig(nil)
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to get dynamic config: %w", err)
	}
	expectedNextSequenceNumber, err := c.GetExpectedNextSequenceNumber(nil)
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to get expected next sequence number: %w", err)
	}
	latestPriceEpochAndRound, err := c.GetLatestPriceEpochAndRound(nil)
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to get latest price epoch and round: %w", err)
	}
	staticConfig, err := c.GetStaticConfig(nil)
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to get static config: %w", err)
	}
	transmitters, err := c.GetTransmitters(nil)
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to get transmitters: %w", err)
	}
	isARMHealthy, err := c.IsARMHealthy(nil)
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to get is ARM healthy: %w", err)
	}
	isUnpausedAndARMHealthy, err := c.IsUnpausedAndARMHealthy(nil)
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to get is unpaused and ARM healthy: %w", err)
	}
	latestConfigDetails, err := c.LatestConfigDetails(nil)
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to get latest config details: %w", err)
	}
	latestConfigDigestAndEpoch, err := c.LatestConfigDigestAndEpoch(nil)
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to get latest config digest and epoch: %w", err)
	}
	paused, err := c.Paused(nil)
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to get paused: %w", err)
	}
	return CommitStoreView{
		ContractMetaData:           meta,
		DynamicConfig:              dynamicConfig,
		ExpectedNextSequenceNumber: expectedNextSequenceNumber,
		LatestPriceEpochAndRound:   latestPriceEpochAndRound,
		StaticConfig:               staticConfig,
		Transmitters:               transmitters,
		IsARMHealthy:               isARMHealthy,
		IsUnpausedAndARMHealthy:    isUnpausedAndARMHealthy,
		LatestConfigDetails:        latestConfigDetails,
		LatestConfigDigestAndEpoch: latestConfigDigestAndEpoch,
		Paused:                     paused,
	}, nil
}
