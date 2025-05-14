package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	fee_manager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"
)

// SubscriberDiscountUpdatedIteratorAdapter adapts the concrete iterator to our interface
type SubscriberDiscountUpdatedIteratorAdapter struct {
	*fee_manager.FeeManagerSubscriberDiscountUpdatedIterator
}

func (a *SubscriberDiscountUpdatedIteratorAdapter) GetEvent() *fee_manager.FeeManagerSubscriberDiscountUpdated {
	return a.Event
}

type FeeManagerReader struct {
	*fee_manager.FeeManagerCaller
	*fee_manager.FeeManagerFilterer
}

// FilterSubscriberDiscountUpdated wraps the original method to return our interface
func (w *FeeManagerReader) FilterSubscriberDiscountUpdated(
	opts *bind.FilterOpts,
	subscriber []common.Address,
	feedID [][32]byte) (LogIterator[fee_manager.FeeManagerSubscriberDiscountUpdated], error) {
	iter, err := w.FeeManagerFilterer.FilterSubscriberDiscountUpdated(opts, subscriber, feedID)
	if err != nil {
		return nil, err
	}

	return &SubscriberDiscountUpdatedIteratorAdapter{iter}, nil
}

func NewFeeManagerReader(contract *fee_manager.FeeManager) *FeeManagerReader {
	return &FeeManagerReader{
		FeeManagerFilterer: &contract.FeeManagerFilterer,
		FeeManagerCaller:   &contract.FeeManagerCaller,
	}
}
