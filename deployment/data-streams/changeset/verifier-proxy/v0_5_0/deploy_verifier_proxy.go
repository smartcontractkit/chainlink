package v0_5_0

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
)

var DeployVerifierProxyChangeset = deployment.CreateChangeSet(deployVerifierProxyLogic, deployVerifierProxyPrecondition)

type DeployVerifierProxy struct {
	VerifierProxyAddress common.Address
}

type DeployVerifierProxyConfig struct {
	ChainsToDeploy map[uint64]DeployVerifierProxy
}

func (cc DeployVerifierProxyConfig) Validate() error {
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

func deployVerifierProxyLogic(e deployment.Environment, cc DeployVerifierProxyConfig) (deployment.ChangesetOutput, error) {
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

func deployVerifierProxyPrecondition(_ deployment.Environment, cc DeployVerifierProxyConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierProxyConfig: %w", err)
	}

	return nil
}

func deployVerifier(e deployment.Environment, ab deployment.AddressBook, cc DeployVerifierProxyConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierProxyConfig: %w", err)
	}

	for chainSel := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		deployVerifier := cc.ChainsToDeploy[chainSel]
		_, err := changeset.DeployContract[*verifier_proxy_v0_5_0.VerifierProxy](e, ab, chain, VerifierDeployFn(deployVerifier.VerifierProxyAddress))
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
		if len(chainState.VerifierProxys) == 0 {
			errNoCCS := errors.New("no VerifierProxy on chain")
			e.Logger.Error(errNoCCS)
			return errNoCCS
		}
	}

	return nil
}

// VerifierDeployFn returns a function that deploys a Verifier contract.
func VerifierDeployFn(verifierProxyAddress common.Address) changeset.ContractDeployFn[*verifier_proxy_v0_5_0.VerifierProxy] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy] {
		ccsAddr, ccsTx, ccs, err := verifier_proxy_v0_5_0.DeployVerifierProxy(
			chain.DeployerKey,
			chain.Client,
			verifierProxyAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
