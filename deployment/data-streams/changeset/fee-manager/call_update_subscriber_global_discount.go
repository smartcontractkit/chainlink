package fee_manager

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	goEthTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
)

// UpdateSubscriberGlobalDiscountChangeset sets the global discount for a subscriber
var UpdateSubscriberGlobalDiscountChangeset deployment.ChangeSetV2[UpdateSubscriberGlobalDiscountConfig] = &globalDiscount{}

type globalDiscount struct{}

type UpdateSubscriberGlobalDiscountConfig struct {
	ConfigPerChain map[uint64][]UpdateSubscriberGlobalDiscount
	MCMSConfig     *types.MCMSConfig
}

type UpdateSubscriberGlobalDiscount struct {
	FeeManagerAddress common.Address
	SubscriberAddress common.Address
	TokenAddress      common.Address
	Discount          uint64
}

func (a UpdateSubscriberGlobalDiscount) GetContractAddress() common.Address {
	return a.FeeManagerAddress
}

func (cs globalDiscount) Apply(e deployment.Environment, cfg UpdateSubscriberGlobalDiscountConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.FeeManager.String(),
		cfg.ConfigPerChain,
		LoadFeeManagerState,
		doUpdateSubscriberGlobalDiscount,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building UpdateSubscriberGlobalDiscount txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "UpdateSubscriberGlobalDiscount proposal")
}

func (cs globalDiscount) VerifyPreconditions(e deployment.Environment, cfg UpdateSubscriberGlobalDiscountConfig) error {
	if len(cfg.ConfigPerChain) == 0 {
		return errors.New("ConfigPerChain is empty")
	}
	for cs := range cfg.ConfigPerChain {
		if err := deployment.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}
	return nil
}

func doUpdateSubscriberGlobalDiscount(
	fm *fee_manager_v0_5_0.FeeManager,
	c UpdateSubscriberGlobalDiscount,
) (*goEthTypes.Transaction, error) {
	return fm.UpdateSubscriberGlobalDiscount(
		deployment.SimTransactOpts(),
		c.SubscriberAddress,
		c.TokenAddress,
		c.Discount)
}
