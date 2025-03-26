package verification

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
)

// SetFeeManagerChangeset sets the active FeeManager contract on the proxy contract
var SetFeeManagerChangeset deployment.ChangeSetV2[VerifierProxySetFeeManagerConfig] = &verifierProxySetFeeManager{}

type verifierProxySetFeeManager struct{}

type VerifierProxySetFeeManagerConfig struct {
	ConfigPerChain map[uint64][]SetFeeManagerConfig
	MCMSConfig     *changeset.MCMSConfig
}

type SetFeeManagerConfig struct {
	ContractAddress   common.Address
	FeeManagerAddress common.Address
}

func (v verifierProxySetFeeManager) Apply(e deployment.Environment, cfg VerifierProxySetFeeManagerConfig) (deployment.ChangesetOutput, error) {
	txs, err := GetSetFeeManagerTxs(e, cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	if cfg.MCMSConfig != nil {
		proposal, err := mcmsutil.CreateMCMSProposal(e, txs, cfg.MCMSConfig.MinDelay, "Set FeeManager proposal")
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
		}, nil
	}

	_, err = txutil.SignAndExecute(e, txs)
	return deployment.ChangesetOutput{}, err
}

// GetSetFeeManagerTxs - returns the transactions to set fee manager on the verifier proxy.
// Does not sign the TXs
func GetSetFeeManagerTxs(e deployment.Environment, cfg VerifierProxySetFeeManagerConfig) ([]*txutil.PreparedTx, error) {
	var preparedTxs []*txutil.PreparedTx
	for chainSelector, configs := range cfg.ConfigPerChain {
		for _, config := range configs {
			state, err := maybeLoadVerifierProxyState(e, chainSelector, config.ContractAddress.String())
			if err != nil {
				return nil, fmt.Errorf("failed to load verifier proxy state: %w", err)
			}
			tx, err := state.VerifierProxy.SetFeeManager(deployment.SimTransactOpts(), config.FeeManagerAddress)
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

func (v verifierProxySetFeeManager) VerifyPreconditions(e deployment.Environment, cfg VerifierProxySetFeeManagerConfig) error {
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
