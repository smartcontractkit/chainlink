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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
)

// DeployVerifierProxyChangeset deploys VerifierProxy to the chains specified in the config.
var DeployVerifierProxyChangeset deployment.ChangeSetV2[DeployVerifierProxyConfig] = &verifierProxyDeploy{}

type verifierProxyDeploy struct{}
type DeployVerifierProxyConfig struct {
	AccessControllerAddress common.Address
	// ChainsToDeploy is a list of chain selectors to deploy the contract to.
	ChainsToDeploy []uint64
	Version        semver.Version
}

func (cfg DeployVerifierProxyConfig) Validate() error {
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

func (v *verifierProxyDeploy) Apply(e deployment.Environment, cc DeployVerifierProxyConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	err := deploy(e, ab, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy VerifierProxy", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{AddressBook: ab}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

func (v *verifierProxyDeploy) VerifyPreconditions(_ deployment.Environment, cc DeployVerifierProxyConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierProxyConfig: %w", err)
	}
	return nil
}

func deploy(e deployment.Environment, ab deployment.AddressBook, cfg DeployVerifierProxyConfig) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierProxyConfig: %w", err)
	}

	for _, chainSel := range cfg.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		_, err := changeset.DeployContract[*verifier_proxy_v0_5_0.VerifierProxy](e, ab, chain, verifyProxyDeployFn(cfg))
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
		if chainState.VerifierProxys == nil || len(chainState.VerifierProxys[chain.Selector]) == 0 {
			errNoCCS := errors.New("no VerifierProxy on chain")
			e.Logger.Error(errNoCCS)
			return errNoCCS
		}
	}

	return nil
}

// verifyProxyDeployFn returns a function that deploys a VerifyProxy contract.
func verifyProxyDeployFn(cfg DeployVerifierProxyConfig) changeset.ContractDeployFn[*verifier_proxy_v0_5_0.VerifierProxy] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy] {
		addr, tx, contract, err := verifier_proxy_v0_5_0.DeployVerifierProxy(
			chain.DeployerKey,
			chain.Client,
			cfg.AccessControllerAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy]{
			Address:  addr,
			Contract: contract,
			Tx:       tx,
			Tv:       deployment.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
