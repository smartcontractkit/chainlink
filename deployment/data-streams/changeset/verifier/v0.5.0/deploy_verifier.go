package v0_5_0

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"

	"github.com/smartcontractkit/chainlink/deployment"
	datastreams "github.com/smartcontractkit/chainlink/deployment/data-streams"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
)

// DeployVerifierChangeset deploys Verifier to the chains specified in the config.
var DeployVerifierChangeset deployment.ChangeSetV2[DeployVerifierConfig] = &verifierDeploy{}

type verifierDeploy struct{}
type DeployVerifierConfig struct {
	VerifierProxyAddress common.Address
	// ChainsToDeploy is a list of chain selectors to deploy the contract to.
	ChainsToDeploy []uint64
	Version        semver.Version
}

func (cfg DeployVerifierConfig) Validate() error {
	switch cfg.Version {
	case deployment.Version0_5_0:
		// no-op
	default:
		return fmt.Errorf("unsupported contract version %s", cfg.Version)
	}
	if len(cfg.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for _, chain := range cfg.ChainsToDeploy {
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}

func (v *verifierDeploy) Apply(e deployment.Environment, cc DeployVerifierConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	err := deploy(e, ab, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy Verifier", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{AddressBook: ab}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

func (v *verifierDeploy) VerifyPreconditions(_ deployment.Environment, cc DeployVerifierConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierConfig: %w", err)
	}
	return nil
}

func deploy(e deployment.Environment, ab deployment.AddressBook, cfg DeployVerifierConfig) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierConfig: %w", err)
	}

	for _, chainSel := range cfg.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		_, err := changeset.DeployContract[*verifier_v0_5_0.Verifier](e, ab, chain, deployFn(cfg))
		if err != nil {
			return err
		}
		chainAddresses, err := ab.AddressesForChain(chain.Selector)
		if err != nil {
			e.Logger.Errorw("Failed to get chain addresses", "err", err)
			return err
		}
		chainState, err := datastreams.LoadChainConfig(chain, chainAddresses)
		if err != nil {
			e.Logger.Errorw("Failed to load chain state", "err", err)
			return err
		}
		if chainState.Verifiers == nil || len(chainState.Verifiers[chain.Selector]) == 0 {
			errNoCCS := errors.New("no Verifier on chain")
			e.Logger.Error(errNoCCS)
			return errNoCCS
		}
	}

	return nil
}

// deployFn returns a function that deploys a Verifier contract.
func deployFn(cfg DeployVerifierConfig) changeset.ContractDeployFn[*verifier_v0_5_0.Verifier] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*verifier_v0_5_0.Verifier] {
		addr, tx, contract, err := verifier_v0_5_0.DeployVerifier(
			chain.DeployerKey,
			chain.Client,
			cfg.VerifierProxyAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*verifier_v0_5_0.Verifier]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*verifier_v0_5_0.Verifier]{
			Address:  addr,
			Contract: contract,
			Tx:       tx,
			Tv:       deployment.NewTypeAndVersion(types.Verifier, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
