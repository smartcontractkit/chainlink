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

// InitializeVerifierChangeset sets registers an expected Verifier contract on the proxy contract
var InitializeVerifierChangeset deployment.ChangeSetV2[VerifierProxyInitializeVerifierConfig] = &verifierProxyInitializeVerifier{}

type verifierProxyInitializeVerifier struct{}

type VerifierProxyInitializeVerifierConfig struct {
	ConfigPerChain map[uint64][]InitializeVerifierConfig
	MCMSConfig     *changeset.MCMSConfig
}

type InitializeVerifierConfig struct {
	ContractAddress common.Address
	VerifierAddress common.Address
}

func (v verifierProxyInitializeVerifier) Apply(e deployment.Environment, cfg VerifierProxyInitializeVerifierConfig) (deployment.ChangesetOutput, error) {
	txs, err := GetInitializeVerifierTxs(e, cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	if cfg.MCMSConfig != nil {
		proposal, err := mcmsutil.CreateMCMSProposal(e, txs, cfg.MCMSConfig.MinDelay, "InitializeVerifier proposal")
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

// GetInitializeVerifierTxs - returns the transactions to set a verifier on the verifier proxy.
// Does not sign the TXs
func GetInitializeVerifierTxs(e deployment.Environment, cfg VerifierProxyInitializeVerifierConfig) ([]*txutil.PreparedTx, error) {
	var preparedTxs []*txutil.PreparedTx
	for chainSelector, configs := range cfg.ConfigPerChain {
		for _, config := range configs {
			state, err := maybeLoadVerifierProxyState(e, chainSelector, config.ContractAddress.String())
			if err != nil {
				return nil, fmt.Errorf("failed to load verifier proxy state: %w", err)
			}
			tx, err := state.VerifierProxy.InitializeVerifier(deployment.SimTransactOpts(), config.VerifierAddress)
			if err != nil {
				return nil, fmt.Errorf("failed to create InitializeVerifier transaction: %w", err)
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

func (v verifierProxyInitializeVerifier) VerifyPreconditions(e deployment.Environment, cfg VerifierProxyInitializeVerifierConfig) error {
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
