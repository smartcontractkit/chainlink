package v1_5

import (
	"errors"
	"fmt"

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

func GenerateCommitStoreView(c *commit_store.CommitStore) (CommitStoreView, error) {
	if c == nil {
		return CommitStoreView{}, errors.New("cannot generate view for nil CommitStore")
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
	isUnpausedAndNotCursed, err := c.IsUnpausedAndNotCursed(nil)
	if err != nil {
		return CommitStoreView{}, fmt.Errorf("failed to get is unpaused and not cursed: %w", err)
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
		IsUnpausedAndNotCursed:     isUnpausedAndNotCursed,
		LatestConfigDetails:        latestConfigDetails,
		LatestConfigDigestAndEpoch: latestConfigDigestAndEpoch,
		Paused:                     paused,
	}, nil
}
