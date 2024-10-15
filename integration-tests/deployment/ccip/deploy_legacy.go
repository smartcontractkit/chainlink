package ccipdeployment

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
)

var (
	MockRMN              deployment.ContractType = "MockRMN"
	LockReleaseTokenPool deployment.ContractType = "LockReleaseTokenPool"
)

type LegacyContracts interface {
	*mock_rmn_contract.MockRMNContract |
		*weth9.WETH9 |
		*router.Router |
		*link_token_interface.LinkToken |
		*lock_release_token_pool.LockReleaseTokenPool |
		*token_admin_registry.TokenAdminRegistry
}

type LegacyContractDeploy[C LegacyContracts] struct {
	// We just keep all the deploy return values
	// since some will be empty if there's an error.
	Address  common.Address
	Contract C
	Tx       *types.Transaction
	Tv       deployment.TypeAndVersion
	Err      error
}

// TODO: pull up to general deployment pkg somehow
// without exposing all product specific contracts?
func deployLegacyContract[C LegacyContracts](
	lggr logger.Logger,
	chain deployment.Chain,
	addressBook deployment.AddressBook,
	deploy func(chain deployment.Chain) LegacyContractDeploy[C],
) (*LegacyContractDeploy[C], error) {
	contractDeploy := deploy(chain)
	if contractDeploy.Err != nil {
		lggr.Errorw("Failed to deploy contract", "err", contractDeploy.Err)
		return nil, contractDeploy.Err
	}
	_, err := chain.Confirm(contractDeploy.Tx)
	if err != nil {
		lggr.Errorw("Failed to confirm deployment", "err", err)
		return nil, err
	}
	err = addressBook.Save(chain.Selector, contractDeploy.Address.String(), contractDeploy.Tv)
	if err != nil {
		lggr.Errorw("Failed to save contract address", "err", err)
		return nil, err
	}
	return &contractDeploy, nil
}

func DeployMockRMN(
	e deployment.Environment,
	chain deployment.Chain,
	ab deployment.AddressBook,
) (*LegacyContractDeploy[*mock_rmn_contract.MockRMNContract], error) {
	return deployLegacyContract(e.Logger, chain, ab,
		func(chain deployment.Chain) LegacyContractDeploy[*mock_rmn_contract.MockRMNContract] {
			rmnRemoteAddr, tx, rmnRemote, err2 := mock_rmn_contract.DeployMockRMNContract(
				chain.DeployerKey,
				chain.Client,
			)
			return LegacyContractDeploy[*mock_rmn_contract.MockRMNContract]{
				rmnRemoteAddr, rmnRemote, tx, deployment.NewTypeAndVersion(MockRMN, deployment.Version1_5_0), err2,
			}
		})
}

func DeployWrappedNative(
	e deployment.Environment,
	chain deployment.Chain,
	ab deployment.AddressBook,
) (*LegacyContractDeploy[*weth9.WETH9], error) {
	return deployLegacyContract(e.Logger, chain, ab,
		func(chain deployment.Chain) LegacyContractDeploy[*weth9.WETH9] {
			wrappedNativeAddr, tx, wrappedNative, err2 := weth9.DeployWETH9(
				chain.DeployerKey,
				chain.Client,
			)
			return LegacyContractDeploy[*weth9.WETH9]{
				wrappedNativeAddr, wrappedNative, tx, deployment.NewTypeAndVersion(WETH9, deployment.Version1_0_0), err2,
			}
		})
}

func DeployRouter(
	e deployment.Environment,
	chain deployment.Chain,
	ab deployment.AddressBook,
	wrappedNative common.Address,
	armProxy common.Address,
) (*LegacyContractDeploy[*router.Router], error) {
	return deployLegacyContract(e.Logger, chain, ab,
		func(chain deployment.Chain) LegacyContractDeploy[*router.Router] {
			routerAddr, tx, routerContract, err2 := router.DeployRouter(
				chain.DeployerKey,
				chain.Client,
				wrappedNative,
				armProxy,
			)
			return LegacyContractDeploy[*router.Router]{
				routerAddr, routerContract, tx, deployment.NewTypeAndVersion(Router, deployment.Version1_2_0), err2,
			}
		})
}

func DeployLinkToken(
	e deployment.Environment,
	chain deployment.Chain,
	ab deployment.AddressBook,
) (*LegacyContractDeploy[*link_token_interface.LinkToken], error) {
	return deployLegacyContract(e.Logger, chain, ab,
		func(chain deployment.Chain) LegacyContractDeploy[*link_token_interface.LinkToken] {
			linkTokenAddr, tx, linkToken, err2 := link_token_interface.DeployLinkToken(
				chain.DeployerKey,
				chain.Client,
			)
			return LegacyContractDeploy[*link_token_interface.LinkToken]{
				linkTokenAddr, linkToken, tx, deployment.NewTypeAndVersion(LinkToken, deployment.Version1_0_0), err2,
			}
		})
}

func DeployLockReleaseTokenPool(
	e deployment.Environment,
	chain deployment.Chain,
	ab deployment.AddressBook,
	token common.Address,
	allowlist []common.Address,
	rmnProxy common.Address,
	acceptLiquidity bool,
	router common.Address,
) (*LegacyContractDeploy[*lock_release_token_pool.LockReleaseTokenPool], error) {
	return deployLegacyContract(e.Logger, chain, ab,
		func(chain deployment.Chain) LegacyContractDeploy[*lock_release_token_pool.LockReleaseTokenPool] {
			poolAddr, tx, pool, err2 := lock_release_token_pool.DeployLockReleaseTokenPool(
				chain.DeployerKey,
				chain.Client,
				token,
				allowlist,
				rmnProxy,
				acceptLiquidity,
				router,
			)
			return LegacyContractDeploy[*lock_release_token_pool.LockReleaseTokenPool]{
				poolAddr, pool, tx, deployment.NewTypeAndVersion(LockReleaseTokenPool, deployment.Version1_5_0), err2,
			}
		})
}

func DeployTokenAdminRegistry(
	e deployment.Environment,
	chain deployment.Chain,
	ab deployment.AddressBook,
) (*LegacyContractDeploy[*token_admin_registry.TokenAdminRegistry], error) {
	return deployLegacyContract(e.Logger, chain, ab,
		func(chain deployment.Chain) LegacyContractDeploy[*token_admin_registry.TokenAdminRegistry] {
			registryAddr, tx, registry, err2 := token_admin_registry.DeployTokenAdminRegistry(
				chain.DeployerKey,
				chain.Client,
			)
			return LegacyContractDeploy[*token_admin_registry.TokenAdminRegistry]{
				registryAddr, registry, tx, deployment.NewTypeAndVersion(TokenAdminRegistry, deployment.Version1_5_0), err2,
			}
		})
}

func DeployLegacyContracts(
	e deployment.Environment,
	chain deployment.Chain,
	ab deployment.AddressBook,
) error {
	rmn, err := DeployMockRMN(e, chain, ab)
	if err != nil {
		return err
	}
	weth, err := DeployWrappedNative(e, chain, ab)
	if err != nil {
		return err
	}
	router, err := DeployRouter(e, chain, ab, weth.Address, rmn.Address)
	if err != nil {
		return err
	}
	link, err := DeployLinkToken(e, chain, ab)
	if err != nil {
		return err
	}

	_, err = DeployLockReleaseTokenPool(e, chain, ab, link.Address, []common.Address{},
		rmn.Address,
		true,
		router.Address,
	)
	if err != nil {
		return err
	}

	registry, err := DeployTokenAdminRegistry(e, chain, ab)
	if err != nil {
		return err
	}

	callOpts := &bind.CallOpts{}
	registryOwner, err := registry.Contract.Owner(callOpts)
	if err != nil {
		return err
	}
	fmt.Println(registryOwner)

	return nil
}
