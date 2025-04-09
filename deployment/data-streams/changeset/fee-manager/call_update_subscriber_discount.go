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

// UpdateSubscriberDiscountChangeset sets the discount for a subscriber
var UpdateSubscriberDiscountChangeset deployment.ChangeSetV2[UpdateSubscriberDiscountConfig] = &discount{}

type discount struct{}

type UpdateSubscriberDiscountConfig struct {
	ConfigPerChain map[uint64][]UpdateSubscriberDiscount
	MCMSConfig     *types.MCMSConfig
}

type UpdateSubscriberDiscount struct {
	FeeManagerAddress common.Address
	SubscriberAddress common.Address
	FeedID            [32]byte
	TokenAddress      common.Address
	Discount          uint64
}

func (a UpdateSubscriberDiscount) GetContractAddress() common.Address {
	return a.FeeManagerAddress
}

func (cs discount) Apply(e deployment.Environment, cfg UpdateSubscriberDiscountConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.FeeManager.String(),
		cfg.ConfigPerChain,
		LoadFeeManagerState,
		doUpdateSubscriberDiscount,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building UpdateSubscriberDiscount txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "UpdateSubscriberDiscount proposal")
}

func (cs discount) VerifyPreconditions(e deployment.Environment, cfg UpdateSubscriberDiscountConfig) error {
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

func doUpdateSubscriberDiscount(
	fm *fee_manager_v0_5_0.FeeManager,
	c UpdateSubscriberDiscount,
) (*goEthTypes.Transaction, error) {
	return fm.UpdateSubscriberDiscount(
		deployment.SimTransactOpts(),
		c.SubscriberAddress,
		c.FeedID,
		c.TokenAddress,
		c.Discount)
}
