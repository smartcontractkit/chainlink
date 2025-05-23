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

// SetAccessControllerChangeset sets the access controller contract on the proxy contract
var SetAccessControllerChangeset cldf.ChangeSetV2[VerifierProxySetAccessControllerConfig] = &verifierProxySetAccessController{}

type verifierProxySetAccessController struct{}

type VerifierProxySetAccessControllerConfig struct {
	ConfigPerChain map[uint64][]SetAccessControllerConfig
	MCMSConfig     *types.MCMSConfig
}

type SetAccessControllerConfig struct {
	ContractAddress         common.Address
	AccessControllerAddress common.Address
}

func (v verifierProxySetAccessController) Apply(e cldf.Environment, cfg VerifierProxySetAccessControllerConfig) (cldf.ChangesetOutput, error) {
	txs, err := GetSetAccessControllerTxs(e, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if cfg.MCMSConfig != nil {
		proposal, err := mcmsutil.CreateMCMSProposal(e, txs, cfg.MCMSConfig.MinDelay, "Set Access Controller proposal")
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

// GetSetAccessControllerTxs - returns the transactions to set the access controller on the verifier proxy.
// Does not sign the TXs
func GetSetAccessControllerTxs(e cldf.Environment, cfg VerifierProxySetAccessControllerConfig) ([]*txutil.PreparedTx, error) {
	var preparedTxs []*txutil.PreparedTx
	for chainSelector, configs := range cfg.ConfigPerChain {
		for _, config := range configs {
			state, err := maybeLoadVerifierProxyState(e, chainSelector, config.ContractAddress.String())
			if err != nil {
				return nil, fmt.Errorf("failed to load verifier proxy state: %w", err)
			}
			tx, err := state.VerifierProxy.SetAccessController(cldf.SimTransactOpts(), config.AccessControllerAddress)

			if err != nil {
				return nil, fmt.Errorf("failed to create SetAccessController transaction: %w", err)
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

func (v verifierProxySetAccessController) VerifyPreconditions(e cldf.Environment, cfg VerifierProxySetAccessControllerConfig) error {
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
