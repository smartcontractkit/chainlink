package fee_manager

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	goEthTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
)

// WithdrawChangeset will withdraw funds from the FeeManager contract to a recipient
var WithdrawChangeset deployment.ChangeSetV2[FeeManagerWithdrawConfig] = &withdraw{}

type withdraw struct{}

type FeeManagerWithdrawConfig struct {
	ConfigPerChain map[uint64][]Withdraw
	MCMSConfig     *types.MCMSConfig
}

type Withdraw struct {
	FeeManagerAddress common.Address
	AssetAddress      common.Address
	RecipientAddress  common.Address
	Quantity          *big.Int
}

func (a Withdraw) GetContractAddress() common.Address {
	return a.FeeManagerAddress
}

func (cs withdraw) Apply(e deployment.Environment, cfg FeeManagerWithdrawConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.FeeManager.String(),
		cfg.ConfigPerChain,
		LoadFeeManagerState,
		doWithdraw,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building Withdraw txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "Withdraw proposal")
}

func (cs withdraw) VerifyPreconditions(e deployment.Environment, cfg FeeManagerWithdrawConfig) error {
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

func doWithdraw(
	fm *fee_manager_v0_5_0.FeeManager,
	c Withdraw,
) (*goEthTypes.Transaction, error) {
	return fm.Withdraw(
		deployment.SimTransactOpts(),
		c.AssetAddress,
		c.RecipientAddress,
		c.Quantity)
}
