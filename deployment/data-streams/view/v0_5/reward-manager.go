package v0_5

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
)

type RewardManagerView struct {
	TypeAndVersion string         `json:"typeAndVersion,omitempty"`
	Address        common.Address `json:"address,omitempty"`
	Owner          common.Address `json:"owner,omitempty"`
}

// GenerateRewardManagerView generates a RewardManagerView from a RewardManager contract.
func GenerateRewardManagerView(rm *rewardManager.RewardManager) (RewardManagerView, error) {
	if rm == nil {
		return RewardManagerView{}, errors.New("cannot generate view for nil RewardManager")
	}

	owner, err := rm.Owner(nil)
	if err != nil {
		return RewardManagerView{}, fmt.Errorf("failed to get owner for RewardManager: %w", err)
	}

	return RewardManagerView{
		Address:        rm.Address(),
		Owner:          owner,
		TypeAndVersion: "RewardManager 0.5.0",
	}, nil
}
