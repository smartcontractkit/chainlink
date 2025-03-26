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

// UnsetVerifierChangeset unsets a registered verifier on the proxy contract
var UnsetVerifierChangeset deployment.ChangeSetV2[VerifierProxyUnsetVerifierConfig] = &verifierProxyUnsetVerifier{}

type verifierProxyUnsetVerifier struct{}

type VerifierProxyUnsetVerifierConfig struct {
	ConfigPerChain map[uint64][]UnsetVerifierConfig
	MCMSConfig     *changeset.MCMSConfig
}

type UnsetVerifierConfig struct {
	ContractAddress common.Address
	ConfigDigest    [32]byte
}

func (v verifierProxyUnsetVerifier) Apply(e deployment.Environment, cfg VerifierProxyUnsetVerifierConfig) (deployment.ChangesetOutput, error) {
	txs, err := GetUnsetVerifierTxs(e, cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	if cfg.MCMSConfig != nil {
		proposal, err := mcmsutil.CreateMCMSProposal(e, txs, cfg.MCMSConfig.MinDelay, "UnsetVerifier proposal")
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

// GetUnsetVerifierTxs - returns the transactions to run the operation on the verifier proxy.
// Does not sign the TXs
func GetUnsetVerifierTxs(e deployment.Environment, cfg VerifierProxyUnsetVerifierConfig) ([]*txutil.PreparedTx, error) {
	var preparedTxs []*txutil.PreparedTx
	for chainSelector, configs := range cfg.ConfigPerChain {
		for _, config := range configs {
			state, err := maybeLoadVerifierProxyState(e, chainSelector, config.ContractAddress.String())
			if err != nil {
				return nil, fmt.Errorf("failed to load verifier proxy state: %w", err)
			}
			tx, err := state.VerifierProxy.UnsetVerifier(deployment.SimTransactOpts(), config.ConfigDigest)
			if err != nil {
				return nil, fmt.Errorf("failed to create UnsetVerifier transaction: %w", err)
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

func (v verifierProxyUnsetVerifier) VerifyPreconditions(e deployment.Environment, cfg VerifierProxyUnsetVerifierConfig) error {
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
