package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
)

// RewardRecipientsUpdatedIteratorAdapter adapts the concrete iterator to our interface
type RewardRecipientsUpdatedIteratorAdapter struct {
	*reward_manager_v0_5_0.RewardManagerRewardRecipientsUpdatedIterator
}

// GetEvent returns the current event
func (a *RewardRecipientsUpdatedIteratorAdapter) GetEvent() *reward_manager_v0_5_0.RewardManagerRewardRecipientsUpdated {
	return a.Event
}

// RewardManagerReader wraps the actual contract to return our interface
type RewardManagerReader struct {
	// Embed the real contract so its methods are promoted
	*reward_manager_v0_5_0.RewardManagerFilterer
	*reward_manager_v0_5_0.RewardManagerCaller
}

// FilterRewardRecipientsUpdated wraps the original method to return our interface
func (r *RewardManagerReader) FilterRewardRecipientsUpdated(
	opts *bind.FilterOpts,
	poolID [][32]byte,
) (LogIterator[reward_manager_v0_5_0.RewardManagerRewardRecipientsUpdated], error) {
	iter, err := r.RewardManagerFilterer.FilterRewardRecipientsUpdated(opts, poolID)
	if err != nil {
		return nil, err
	}

	return &RewardRecipientsUpdatedIteratorAdapter{iter}, nil
}

// NewRewardManagerReader creates a new wrapper for the RewardManager contract
func NewRewardManagerReader(contract *reward_manager_v0_5_0.RewardManager) *RewardManagerReader {
	return &RewardManagerReader{
		RewardManagerCaller:   &contract.RewardManagerCaller,
		RewardManagerFilterer: &contract.RewardManagerFilterer,
	}
}
