package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/registry_module_owner_custom"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/multicall3"
)

var (
	_ deployment.ChangeSet[DeployPrerequisiteConfig] = DeployPrerequisites
)

// DeployPrerequisites deploys the pre-requisite contracts for CCIP
// pre-requisite contracts are the contracts which can be reused from previous versions of CCIP
// Or the contracts which are already deployed on the chain ( for example, tokens, feeds, etc)
// Caller should update the environment's address book with the returned addresses.
func DeployPrerequisites(env deployment.Environment, cfg DeployPrerequisiteConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Wrapf(deployment.ErrInvalidConfig, "%v", err)
	}
	ab := deployment.NewMemoryAddressBook()
	err = deployPrerequisiteChainContracts(env, ab, cfg.ChainSelectors, cfg.Opts...)
	if err != nil {
		env.Logger.Errorw("Failed to deploy prerequisite contracts", "err", err, "addressBook", ab)
		return deployment.ChangesetOutput{
			AddressBook: ab,
		}, fmt.Errorf("failed to deploy prerequisite contracts: %w", err)
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: ab,
		JobSpecs:    nil,
	}, nil
}

type DeployPrerequisiteContractsOpts struct {
	USDCEnabledChains []uint64
	Multicall3Enabled bool
}

type DeployPrerequisiteConfig struct {
	ChainSelectors []uint64
	Opts           []PrerequisiteOpt
	// TODO handle tokens and feeds in prerequisite config
	Tokens map[TokenSymbol]common.Address
	Feeds  map[TokenSymbol]common.Address
}

func (c DeployPrerequisiteConfig) Validate() error {
	mapAllChainSelectors := make(map[uint64]struct{})
	for _, cs := range c.ChainSelectors {
		mapAllChainSelectors[cs] = struct{}{}
		if err := deployment.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
	}
	return nil
}

type PrerequisiteOpt func(o *DeployPrerequisiteContractsOpts)

func WithUSDCChains(chains []uint64) PrerequisiteOpt {
	return func(o *DeployPrerequisiteContractsOpts) {
		o.USDCEnabledChains = chains
	}
}

func WithMulticall3(enabled bool) PrerequisiteOpt {
	return func(o *DeployPrerequisiteContractsOpts) {
		o.Multicall3Enabled = enabled
	}
}

func deployPrerequisiteChainContracts(e deployment.Environment, ab deployment.AddressBook, selectors []uint64, opts ...PrerequisiteOpt) error {
	state, err := LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err")
		return err
	}
	deployGrp := errgroup.Group{}
	for _, sel := range selectors {
		chain := e.Chains[sel]
		deployGrp.Go(func() error {
			err := deployPrerequisiteContracts(e, ab, state, chain, opts...)
			if err != nil {
				e.Logger.Errorw("Failed to deploy prerequisite contracts", "chain", sel, "err", err)
				return err
			}
			return nil
		})
	}
	return deployGrp.Wait()
}

// deployPrerequisiteContracts deploys the contracts that can be ported from previous CCIP version to the new one.
// This is only required for staging and test environments where the contracts are not already deployed.
func deployPrerequisiteContracts(e deployment.Environment, ab deployment.AddressBook, state CCIPOnChainState, chain deployment.Chain, opts ...PrerequisiteOpt) error {
	deployOpts := &DeployPrerequisiteContractsOpts{}
	for _, opt := range opts {
		if opt != nil {
			opt(deployOpts)
		}
	}
	var isUSDC bool
	for _, sel := range deployOpts.USDCEnabledChains {
		if sel == chain.Selector {
			isUSDC = true
			break
		}
	}
	lggr := e.Logger
	chainState, chainExists := state.Chains[chain.Selector]
	var weth9Contract *weth9.WETH9
	var tokenAdminReg *token_admin_registry.TokenAdminRegistry
	var registryModule *registry_module_owner_custom.RegistryModuleOwnerCustom
	var rmnProxy *rmn_proxy_contract.RMNProxyContract
	var r *router.Router
	var mc3 *multicall3.Multicall3
	if chainExists {
		weth9Contract = chainState.Weth9
		tokenAdminReg = chainState.TokenAdminRegistry
		registryModule = chainState.RegistryModule
		rmnProxy = chainState.RMNProxy
		r = chainState.Router
		mc3 = chainState.Multicall3
	}
	if rmnProxy == nil {
		rmn, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*mock_rmn_contract.MockRMNContract] {
				rmnAddr, tx2, rmn, err2 := mock_rmn_contract.DeployMockRMNContract(
					chain.DeployerKey,
					chain.Client,
				)
				return deployment.ContractDeploy[*mock_rmn_contract.MockRMNContract]{
					Address: rmnAddr, Contract: rmn, Tx: tx2, Tv: deployment.NewTypeAndVersion(MockRMN, deployment.Version1_0_0), Err: err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy mock RMN", "chain", chain.String(), "err", err)
			return err
		}
		rmnProxyContract, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*rmn_proxy_contract.RMNProxyContract] {
				rmnProxyAddr, tx2, rmnProxy, err2 := rmn_proxy_contract.DeployRMNProxyContract(
					chain.DeployerKey,
					chain.Client,
					rmn.Address,
				)
				return deployment.ContractDeploy[*rmn_proxy_contract.RMNProxyContract]{
					rmnProxyAddr, rmnProxy, tx2, deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0), err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy RMNProxyExisting", "chain", chain.String(), "err", err)
			return err
		}
		rmnProxy = rmnProxyContract.Contract
	}
	if tokenAdminReg == nil {
		tokenAdminRegistry, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*token_admin_registry.TokenAdminRegistry] {
				tokenAdminRegistryAddr, tx2, tokenAdminRegistry, err2 := token_admin_registry.DeployTokenAdminRegistry(
					chain.DeployerKey,
					chain.Client)
				return deployment.ContractDeploy[*token_admin_registry.TokenAdminRegistry]{
					tokenAdminRegistryAddr, tokenAdminRegistry, tx2, deployment.NewTypeAndVersion(TokenAdminRegistry, deployment.Version1_5_0), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy token admin registry", "chain", chain.String(), "err", err)
			return err
		}
		tokenAdminReg = tokenAdminRegistry.Contract
	} else {
		e.Logger.Infow("tokenAdminRegistry already deployed", "chain", chain.String(), "addr", tokenAdminReg.Address)
	}
	if registryModule == nil {
		customRegistryModule, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom] {
				regModAddr, tx2, regMod, err2 := registry_module_owner_custom.DeployRegistryModuleOwnerCustom(
					chain.DeployerKey,
					chain.Client,
					tokenAdminReg.Address())
				return deployment.ContractDeploy[*registry_module_owner_custom.RegistryModuleOwnerCustom]{
					regModAddr, regMod, tx2, deployment.NewTypeAndVersion(RegistryModule, deployment.Version1_5_0), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy custom registry module", "chain", chain.String(), "err", err)
			return err
		}
		registryModule = customRegistryModule.Contract
	} else {
		e.Logger.Infow("custom registry module already deployed", "chain", chain.String(), "addr", registryModule.Address)
	}
	isRegistryAdded, err := tokenAdminReg.IsRegistryModule(nil, registryModule.Address())
	if err != nil {
		e.Logger.Errorw("Failed to check if registry module is added on token admin registry", "chain", chain.String(), "err", err)
		return fmt.Errorf("failed to check if registry module is added on token admin registry: %w", err)
	}
	if !isRegistryAdded {
		tx, err := tokenAdminReg.AddRegistryModule(chain.DeployerKey, registryModule.Address())
		if err != nil {
			e.Logger.Errorw("Failed to assign registry module on token admin registry", "chain", chain.String(), "err", err)
			return fmt.Errorf("failed to assign registry module on token admin registry: %w", err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			e.Logger.Errorw("Failed to confirm assign registry module on token admin registry", "chain", chain.String(), "err", err)
			return fmt.Errorf("failed to confirm assign registry module on token admin registry: %w", err)
		}
		e.Logger.Infow("assigned registry module on token admin registry")
	}
	if weth9Contract == nil {
		weth, err := deployment.DeployContract(lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*weth9.WETH9] {
				weth9Addr, tx2, weth9c, err2 := weth9.DeployWETH9(
					chain.DeployerKey,
					chain.Client,
				)
				return deployment.ContractDeploy[*weth9.WETH9]{
					weth9Addr, weth9c, tx2, deployment.NewTypeAndVersion(WETH9, deployment.Version1_0_0), err2,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy weth9", "chain", chain.String(), "err", err)
			return err
		}
		weth9Contract = weth.Contract
	} else {
		lggr.Infow("weth9 already deployed", "addr", weth9Contract.Address)
	}
	// if router is not already deployed, we deploy it
	if r == nil {
		routerContract, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*router.Router] {
				routerAddr, tx2, routerC, err2 := router.DeployRouter(
					chain.DeployerKey,
					chain.Client,
					weth9Contract.Address(),
					rmnProxy.Address(),
				)
				return deployment.ContractDeploy[*router.Router]{
					routerAddr, routerC, tx2, deployment.NewTypeAndVersion(Router, deployment.Version1_2_0), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy router", "chain", chain.String(), "err", err)
			return err
		}

		r = routerContract.Contract
	} else {
		e.Logger.Infow("router already deployed", "chain", chain.String(), "addr", chainState.Router.Address)
	}
	if deployOpts.Multicall3Enabled && mc3 == nil {
		_, err := deployment.DeployContract(e.Logger, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*multicall3.Multicall3] {
				multicall3Addr, tx2, multicall3Wrapper, err2 := multicall3.DeployMulticall3(
					chain.DeployerKey,
					chain.Client,
				)
				return deployment.ContractDeploy[*multicall3.Multicall3]{
					multicall3Addr, multicall3Wrapper, tx2, deployment.NewTypeAndVersion(Multicall3, deployment.Version1_0_0), err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy ccip multicall", "chain", chain.String(), "err", err)
			return err
		}
	} else {
		if mc3 != nil {
			e.Logger.Info("ccip multicall already deployed", "chain", chain.String(), "addr", mc3.Address)
		}
	}
	if isUSDC {
		token, pool, messenger, transmitter, err1 := DeployUSDC(e.Logger, chain, ab, rmnProxy.Address(), r.Address())
		if err1 != nil {
			return err1
		}
		e.Logger.Infow("Deployed USDC contracts",
			"chain", chain.String(),
			"token", token.Address(),
			"pool", pool.Address(),
			"transmitter", transmitter.Address(),
			"messenger", messenger.Address(),
		)
	}
	return nil
}
