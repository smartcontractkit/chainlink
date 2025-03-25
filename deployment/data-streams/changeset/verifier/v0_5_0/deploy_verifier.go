package v0_5_0

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	verifier "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
)

var DeployVerifierChangeset = deployment.CreateChangeSet(deployVerifierLogic, deployVerifierPrecondition)

type DeployVerifier struct {
	VerifierProxyAddress common.Address
}

type DeployVerifierConfig struct {
	ChainsToDeploy map[uint64]DeployVerifier
}

func (cc DeployVerifierConfig) Validate() error {
	if len(cc.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for chain := range cc.ChainsToDeploy {
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}

func deployVerifierLogic(e deployment.Environment, cc DeployVerifierConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	err := deployVerifier(e, ab, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy Verifier", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{AddressBook: ab}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

func deployVerifierPrecondition(_ deployment.Environment, cc DeployVerifierConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierConfig: %w", err)
	}

	return nil
}

func deployVerifier(e deployment.Environment, ab deployment.AddressBook, cc DeployVerifierConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierConfig: %w", err)
	}

	for chainSel := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		deployVerifier := cc.ChainsToDeploy[chainSel]
		_, err := changeset.DeployContract(e, ab, chain, VerifierDeployFn(deployVerifier.VerifierProxyAddress))
		if err != nil {
			return err
		}
		chainAddresses, err := ab.AddressesForChain(chain.Selector)
		if err != nil {
			e.Logger.Errorw("Failed to get chain addresses", "err", err)
			return err
		}
		chainState, err := changeset.LoadChainState(e.Logger, chain, chainAddresses)
		if err != nil {
			e.Logger.Errorw("Failed to load chain state", "err", err)
			return err
		}
		if len(chainState.Verifiers) == 0 {
			errNoCCS := errors.New("no Verifier on chain")
			e.Logger.Error(errNoCCS)
			return errNoCCS
		}
	}

	return nil
}

func VerifierDeployFn(verifierProxyAddress common.Address) changeset.ContractDeployFn[*verifier.Verifier] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*verifier.Verifier] {
		ccsAddr, ccsTx, ccs, err := verifier.DeployVerifier(
			chain.DeployerKey,
			chain.Client,
			verifierProxyAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*verifier.Verifier]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*verifier.Verifier]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.Verifier, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
