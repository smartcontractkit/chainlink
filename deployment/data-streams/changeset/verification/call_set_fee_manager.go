package verification

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
)

// SetFeeManagerChangeset sets the active FeeManager contract on the proxy contract
var SetFeeManagerChangeset = cldf.CreateChangeSet(verifierProxySetFeeManagerLogic, verifierProxySetFeeManagerPrecondition)

type VerifierProxySetFeeManagerConfig struct {
	ConfigPerChain map[uint64][]SetFeeManagerConfig
	MCMSConfig     *types.MCMSConfig
}

type SetFeeManagerConfig struct {
	VerifierProxyAddress common.Address
	FeeManagerAddress    common.Address
}

func verifierProxySetFeeManagerLogic(e cldf.Environment, cfg VerifierProxySetFeeManagerConfig) (cldf.ChangesetOutput, error) {
	txs, err := GetSetFeeManagerTxs(e, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if cfg.MCMSConfig != nil {
		proposal, err := mcmsutil.CreateMCMSProposal(e, txs, cfg.MCMSConfig.MinDelay, "Set FeeManager proposal")
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
		}, nil
	}

	_, err = txutil.SignAndExecute(e, txs)
	return cldf.ChangesetOutput{}, err
}

// GetSetFeeManagerTxs - returns the transactions to set fee manager on the verifier proxy.
// Does not sign the TXs
func GetSetFeeManagerTxs(e cldf.Environment, cfg VerifierProxySetFeeManagerConfig) ([]*txutil.PreparedTx, error) {
	var preparedTxs []*txutil.PreparedTx
	for chainSelector, configs := range cfg.ConfigPerChain {
		for _, config := range configs {
			state, err := maybeLoadVerifierProxyState(e, chainSelector, config.VerifierProxyAddress.String())
			if err != nil {
				return nil, fmt.Errorf("failed to load verifier proxy state: %w", err)
			}
			tx, err := state.VerifierProxy.SetFeeManager(cldf.SimTransactOpts(), config.FeeManagerAddress)
			if err != nil {
				return nil, fmt.Errorf("failed to create SetFeeManager transaction: %w", err)
			}
			preparedTx := txutil.PreparedTx{
				Tx:            tx,
				ChainSelector: chainSelector,
				ContractType:  types.VerifierProxy.String(),
			}
			preparedTxs = append(preparedTxs, &preparedTx)
		}
	}
	return preparedTxs, nil
}

func verifierProxySetFeeManagerPrecondition(e cldf.Environment, cfg VerifierProxySetFeeManagerConfig) error {
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
