package fee_manager

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	goEthTypes "github.com/ethereum/go-ethereum/core/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
)

// SetNativeSurchargeChangeset sets the native surcharge on the FeeManager contract
var SetNativeSurchargeChangeset = cldf.CreateChangeSet(setNativeSurchargeLogic, setNativeSurchargePrecondition)

type SetNativeSurchargeConfig struct {
	ConfigPerChain map[uint64][]SetNativeSurcharge
	MCMSConfig     *types.MCMSConfig
}

type SetNativeSurcharge struct {
	FeeManagerAddress common.Address
	Surcharge         uint64
}

func (a SetNativeSurcharge) GetContractAddress() common.Address {
	return a.FeeManagerAddress
}

func setNativeSurchargeLogic(e cldf.Environment, cfg SetNativeSurchargeConfig) (cldf.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.FeeManager.String(),
		cfg.ConfigPerChain,
		LoadFeeManagerState,
		doSetNativeSurcharge,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed building SetNativeSurcharge txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "SetNativeSurcharge proposal")
}

func setNativeSurchargePrecondition(e cldf.Environment, cfg SetNativeSurchargeConfig) error {
	if len(cfg.ConfigPerChain) == 0 {
		return errors.New("ConfigPerChain is empty")
	}
	for cs := range cfg.ConfigPerChain {
		if err := cldf.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}
	return nil
}

func doSetNativeSurcharge(
	fm *fee_manager_v0_5_0.FeeManager,
	c SetNativeSurcharge,
) (*goEthTypes.Transaction, error) {
	return fm.SetNativeSurcharge(
		cldf.SimTransactOpts(),
		c.Surcharge)
}
