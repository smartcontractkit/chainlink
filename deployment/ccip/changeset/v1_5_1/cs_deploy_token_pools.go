package v1_5_1

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/ccip-contract-examples/chains/evm/gobindings/generated/1_6_1/burn_mint_with_external_minter_token_pool"
	"github.com/smartcontractkit/ccip-contract-examples/chains/evm/gobindings/generated/1_6_1/hybrid_with_external_minter_token_pool"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_from_mint_token_pool"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_with_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fast_transfer_token_pool"
	burn_mint_token_pool_v1_6_1 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_1/burn_mint_token_pool"
	lock_release_token_pool_v1_6_1 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_1/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/erc20"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/burn_mint_with_external_minter_fast_transfer_token_pool"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/hybrid_with_external_minter_fast_transfer_token_pool"
)

var _ cldf.ChangeSet[DeployTokenPoolContractsConfig] = DeployTokenPoolContractsChangeset

// DeployTokenPoolInput defines all information required of the user to deploy a new token pool contract.
type DeployTokenPoolInput struct {
	// Type is the type of token pool that must be deployed.
	Type cldf.ContractType
	// Version is the version of the token pool that must be deployed.
	Version semver.Version
	// TokenAddress is the address of the token for which we are deploying a pool.
	TokenAddress common.Address
	// TokenType is the type of token that is being deployed. This is used to determine if we should grant burn and mint
	// permissions to the token pool contract (BurnMintERC20).
	TokenType cldf.ContractType
	// AllowList is the optional list of addresses permitted to initiate a token transfer.
	// If omitted, all addresses will be permitted to transfer the token.
	AllowList []common.Address
	// LocalTokenDecimals is the number of decimals used by the token at tokenAddress.
	LocalTokenDecimals uint8
	// AcceptLiquidity indicates whether or not the new pool can accept liquidity from a rebalancer address (lock-release only).
	AcceptLiquidity *bool
	// ExternalMinter only for burn-mint fast transfer pools with external minting.
	ExternalMinter common.Address
	// CCIPAdmin is the address of the CCIP admin for the token and will have default admin role. This is specifically
	// for BurnMintERC20 token.
	CCIPAdmin common.Address
	// TokenGovernor is the address of the token governor contract. This is specifically for BurnMintWithExternalMinterTokenPool
	// and HybridWithExternalMinterTokenPool token pools.
	TokenGovernor common.Address
}

func (i DeployTokenPoolInput) Validate(ctx context.Context, chain cldf_evm.Chain, state evm.CCIPChainState, tokenSymbol shared.TokenSymbol) error {
	// Ensure that required fields are populated
	if i.TokenAddress == utils.ZeroAddress {
		return errors.New("token address must be defined")
	}
	if i.Type == cldf.ContractType("") {
		return errors.New("type must be defined")
	}

	// Validate that the type is known
	if _, ok := shared.TokenPoolTypes[i.Type]; !ok {
		return fmt.Errorf("requested token pool type %s is unknown", i.Type)
	}

	// Validate the token exists and matches the expected symbol
	token, err := erc20.NewERC20(i.TokenAddress, chain.Client)
	if err != nil {
		return fmt.Errorf("failed to connect address %s with erc20 bindings: %w", i.TokenAddress, err)
	}
	symbol, err := token.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to fetch symbol from token with address %s: %w", i.TokenAddress, err)
	}
	if symbol != string(tokenSymbol) {
		return fmt.Errorf("symbol of token with address %s (%s) does not match expected symbol (%s)", i.TokenAddress, symbol, tokenSymbol)
	}

	// Validate localTokenDecimals against the decimals value on the token contract
	decimals, err := token.Decimals(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to fetch decimals from token with address %s: %w", i.TokenAddress, err)
	}
	if decimals != i.LocalTokenDecimals {
		return fmt.Errorf("decimals of token with address %s (%d) does not match localTokenDecimals (%d)", i.TokenAddress, decimals, i.LocalTokenDecimals)
	}

	// Validate acceptLiquidity based on requested pool type
	if i.Type == shared.LockReleaseTokenPool && i.AcceptLiquidity == nil {
		return errors.New("accept liquidity must be defined for lock release pools")
	}
	if i.Type != shared.LockReleaseTokenPool && i.AcceptLiquidity != nil {
		return errors.New("accept liquidity must be nil for burn mint pools")
	}
	if i.Type == shared.BurnMintWithExternalMinterFastTransferTokenPool && i.ExternalMinter == utils.ZeroAddress {
		// TODO: Validate that the external minter token match the token in input
		return errors.New("external minter must be defined for burn mint with external minter fast transfer pools")
	}
	if i.Type == shared.HybridWithExternalMinterFastTransferTokenPool && i.ExternalMinter == utils.ZeroAddress {
		// TODO: Validate that the external minter token match the token in input
		return errors.New("external minter must be defined for hybrid with external minter fast transfer pools")
	}

	// We should check if a token pool with this type, version, and symbol already exists
	_, ok := GetTokenPoolAddressFromSymbolTypeAndVersion(state, chain, tokenSymbol, i.Type, shared.CurrentTokenPoolVersion)
	if ok {
		return fmt.Errorf("token pool with type %s and version %s already exists for %s on %s", i.Type, shared.CurrentTokenPoolVersion, tokenSymbol, chain)
	}

	return nil
}

// DeployTokenPoolContractsConfig defines the token pool contracts that need to be deployed on each chain.
type DeployTokenPoolContractsConfig struct {
	// Symbol is the symbol of the token for which we are deploying a pool.
	TokenSymbol shared.TokenSymbol
	// NewPools defines the per-chain configuration of each new pool
	NewPools map[uint64]DeployTokenPoolInput
	// IsTestRouter indicates whether or not the test router should be used.
	IsTestRouter bool
	// ForceDatastoreOverwrite permits this changeset to overwrite pool refs that already exist in
	// the environment datastore. Without it, a redeploy over an existing (chain, type, version,
	// TokenSymbol) entry is rejected during validation rather than silently replacing it.
	ForceDatastoreOverwrite bool
}

// deployedTypeAndVersion returns the TypeAndVersion the pool will actually be deployed under.
//
// The version is not simply poolConfig.Version: several pool types are pinned to a version of
// their own regardless of what the config asked for, and two more honour an explicit 1.6.1. The
// datastore key is built from this version, so the rule has to live in one function that both the
// key and the deployment read. Deriving the key from poolConfig.Version instead would key some
// pools under a version they were never deployed at.
func deployedTypeAndVersion(poolConfig DeployTokenPoolInput) cldf.TypeAndVersion {
	version := shared.CurrentTokenPoolVersion

	switch poolConfig.Type {
	case shared.BurnMintTokenPool, shared.LockReleaseTokenPool:
		if poolConfig.Version == deployment.Version1_6_1 {
			version = deployment.Version1_6_1
		}
	case shared.BurnMintFastTransferTokenPool:
		version = deployment.Version1_6_3Dev
	case shared.BurnMintWithExternalMinterFastTransferTokenPool,
		shared.HybridWithExternalMinterFastTransferTokenPool,
		shared.BurnMintWithExternalMinterTokenPool,
		shared.HybridWithExternalMinterTokenPool:
		version = deployment.Version1_6_0
	}

	return cldf.NewTypeAndVersion(poolConfig.Type, version)
}

// plannedRefs returns the datastore ref for every pool this changeset will deploy.
//
// A ref is keyed on (chain, type, version, qualifier) and none of that depends on the address, so
// the whole set is known before a single transaction is sent. This is the only place these keys
// are derived: Validate checks them and reserveRefs claims them, so the check and the write are
// the same computation rather than two that can drift apart.
func (c DeployTokenPoolContractsConfig) plannedRefs() []datastore.AddressRef {
	refs := make([]datastore.AddressRef, 0, len(c.NewPools))
	for chainSelector, poolConfig := range c.NewPools {
		tv := deployedTypeAndVersion(poolConfig)
		version := tv.Version
		refs = append(refs, datastore.AddressRef{
			ChainSelector: chainSelector,
			Type:          datastore.ContractType(tv.Type),
			Version:       &version,
			Qualifier:     string(c.TokenSymbol),
		})
	}

	return refs
}

// reserveRefs claims every planned key. The deployment then fills each one in as it confirms, so
// the store is complete by the time the changeset returns and there is no derivation step left
// that could fail once the pools already exist on chain.
func (c DeployTokenPoolContractsConfig) reserveRefs() (*datastore.MemoryDataStore, error) {
	ds, err := shared.ReserveRefs(c.plannedRefs())
	if err != nil {
		return nil, fmt.Errorf("%s token pools: %w", c.TokenSymbol, err)
	}

	return ds, nil
}

func (c DeployTokenPoolContractsConfig) Validate(env cldf.Environment) error {
	// Ensure that required fields are populated
	if c.TokenSymbol == shared.TokenSymbol("") {
		return errors.New("token symbol must be defined")
	}

	// Reserve the keys to prove the config can be recorded, then check them against the
	// environment. Both run before anything is deployed.
	if _, err := c.reserveRefs(); err != nil {
		return err
	}
	refs := c.plannedRefs()
	// Taking over a key the environment already holds is a redeploy. It is rejected unless the
	// caller asked for it, in which case it is logged rather than passing unremarked.
	validateRefs := shared.ValidateAddressRefsStrict
	if c.ForceDatastoreOverwrite {
		validateRefs = shared.ValidateAddressRefs
	}
	if err := validateRefs(env, refs); err != nil {
		return fmt.Errorf("token pool datastore refs conflict: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, poolConfig := range c.NewPools {
		err := cldf.IsValidChainSelector(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to validate chain selector %d: %w", chainSelector, err)
		}
		chain, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}
		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in state", chainSelector)
		}
		if c.IsTestRouter {
			if chainState.TestRouter == nil {
				return fmt.Errorf("missing test router on %s", chain.String())
			}
		} else {
			if chainState.Router == nil {
				return fmt.Errorf("missing router on %s", chain.String())
			}
		}
		if rmnProxy := chainState.RMNProxy; rmnProxy == nil {
			return fmt.Errorf("missing rmnProxy on %s", chain.String())
		}
		err = poolConfig.Validate(env.GetContext(), chain, chainState, c.TokenSymbol)
		if err != nil {
			return fmt.Errorf("failed to validate token pool config for chain selector %d: %w", chainSelector, err)
		}
	}
	return nil
}

// DeployTokenPoolContractsChangeset deploys new pools for a given token across multiple chains.
func DeployTokenPoolContractsChangeset(env cldf.Environment, c DeployTokenPoolContractsConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid DeployTokenPoolContractsConfig: %w", err)
	}
	newAddresses := cldf.NewMemoryAddressBook()

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Claim every datastore key before deploying. After this the store is only ever filled in,
	// never re-keyed, so there is no post-deploy datastore step that can fail.
	ds, err := c.reserveRefs()
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	deployGrp := errgroup.Group{}

	for chainSelector, poolConfig := range c.NewPools {
		deployGrp.Go(func() error {
			chain := env.BlockChains.EVMChains()[chainSelector]
			chainState := state.Chains[chainSelector]
			// The ref is written by deployTokenPool at the moment the deployment confirms, onto
			// the key this chain reserved above. Nothing is left to record afterwards.
			contract, err := deployTokenPool(env.Logger, chain, chainState, newAddresses, ds, string(c.TokenSymbol), deployedTypeAndVersion(poolConfig), poolConfig, c.IsTestRouter)
			if err != nil {
				return fmt.Errorf("failed to deploy token pool contract: %w", err)
			}
			if poolConfig.TokenType == shared.BurnMintERC20Token {
				if err := addMinterAndBurnerForBurnMintERC20Token(env, chain.Selector, poolConfig.TokenAddress, contract.Address); err != nil {
					return fmt.Errorf("failed to add minter and burner for BurnMintERC20 token %s on %s: %w",
						poolConfig.TokenAddress, chain, err)
				}
			}
			if poolConfig.TokenType == shared.ERC677TokenHelper || poolConfig.TokenType == commontypes.LinkToken {
				if err := addMinterForERC677Token(env, chain, poolConfig.TokenAddress, contract.Address); err != nil {
					return fmt.Errorf("failed to add Token pool as minter and burner for ERC677 token %s on %s: %w",
						poolConfig.TokenAddress, chain, err)
				}
			}

			return nil
		})
	}

	if err := deployGrp.Wait(); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy %s token pool on %w",
			c.TokenSymbol, err)
	}

	// No datastore is derived here. Every ref was keyed before the first transaction and filled
	// in as each pool was confirmed, so by this point the store is already complete and there is
	// nothing left that can fail.
	return cldf.ChangesetOutput{
		AddressBook: newAddresses,
		DataStore:   ds,
	}, nil
}

// deployTokenPool deploys a token pool contract based on a given type & configuration.
func deployTokenPool(
	logger logger.Logger,
	chain cldf_evm.Chain,
	chainState evm.CCIPChainState,
	addressBook cldf.AddressBook,
	ds datastore.MutableDataStore,
	qualifier string,
	tv cldf.TypeAndVersion,
	poolConfig DeployTokenPoolInput,
	isTestRouter bool,
) (*cldf.ContractDeploy[*token_pool.TokenPool], error) {
	router := chainState.Router
	if isTestRouter {
		router = chainState.TestRouter
	}
	rmnProxy := chainState.RMNProxy

	return shared.DeployContractAndRecord(logger, chain, addressBook, ds, tv, qualifier,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*token_pool.TokenPool] {
			var tpAddr common.Address
			var tx *types.Transaction
			var err error
			switch poolConfig.Type {
			case shared.BurnMintTokenPool:
				if poolConfig.Version == deployment.Version1_6_1 {
					tpAddr, tx, _, err = burn_mint_token_pool_v1_6_1.DeployBurnMintTokenPool(
						chain.DeployerKey, chain.Client, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
						poolConfig.AllowList, rmnProxy.Address(), router.Address(),
					)
				} else {
					tpAddr, tx, _, err = burn_mint_token_pool.DeployBurnMintTokenPool(
						chain.DeployerKey, chain.Client, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
						poolConfig.AllowList, rmnProxy.Address(), router.Address(),
					)
				}
			case shared.BurnWithFromMintTokenPool:
				tpAddr, tx, _, err = burn_with_from_mint_token_pool.DeployBurnWithFromMintTokenPool(
					chain.DeployerKey, chain.Client, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
					poolConfig.AllowList, rmnProxy.Address(), router.Address(),
				)
			case shared.BurnFromMintTokenPool:
				tpAddr, tx, _, err = burn_from_mint_token_pool.DeployBurnFromMintTokenPool(
					chain.DeployerKey, chain.Client, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
					poolConfig.AllowList, rmnProxy.Address(), router.Address(),
				)
			case shared.LockReleaseTokenPool:
				if poolConfig.Version == deployment.Version1_6_1 {
					tpAddr, tx, _, err = lock_release_token_pool_v1_6_1.DeployLockReleaseTokenPool(
						chain.DeployerKey, chain.Client, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
						poolConfig.AllowList, rmnProxy.Address(), router.Address(),
					)
				} else {
					tpAddr, tx, _, err = lock_release_token_pool.DeployLockReleaseTokenPool(
						chain.DeployerKey, chain.Client, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
						poolConfig.AllowList, rmnProxy.Address(), *poolConfig.AcceptLiquidity, router.Address(),
					)
				}
			case shared.BurnMintFastTransferTokenPool:
				tpAddr, tx, _, err = fast_transfer_token_pool.DeployBurnMintFastTransferTokenPool(
					chain.DeployerKey, chain.Client, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
					poolConfig.AllowList, rmnProxy.Address(), router.Address(), chain.Selector,
				)
			case shared.BurnMintWithExternalMinterFastTransferTokenPool:
				tpAddr, tx, _, err = burn_mint_with_external_minter_fast_transfer_token_pool.DeployBurnMintWithExternalMinterFastTransferTokenPool(
					chain.DeployerKey, chain.Client, poolConfig.ExternalMinter, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
					poolConfig.AllowList, rmnProxy.Address(), router.Address(),
				)
			case shared.HybridWithExternalMinterFastTransferTokenPool:
				tpAddr, tx, _, err = hybrid_with_external_minter_fast_transfer_token_pool.DeployHybridWithExternalMinterFastTransferTokenPool(
					chain.DeployerKey, chain.Client, poolConfig.ExternalMinter, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
					poolConfig.AllowList, rmnProxy.Address(), router.Address(),
				)
			case shared.BurnMintWithExternalMinterTokenPool:
				tpAddr, tx, _, err = burn_mint_with_external_minter_token_pool.DeployBurnMintWithExternalMinterTokenPool(
					chain.DeployerKey, chain.Client, poolConfig.TokenGovernor, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
					poolConfig.AllowList, rmnProxy.Address(), router.Address(),
				)
			case shared.HybridWithExternalMinterTokenPool:
				tpAddr, tx, _, err = hybrid_with_external_minter_token_pool.DeployHybridWithExternalMinterTokenPool(
					chain.DeployerKey, chain.Client, poolConfig.TokenGovernor, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
					poolConfig.AllowList, rmnProxy.Address(), router.Address(),
				)
			}
			var tp *token_pool.TokenPool
			if err == nil { // prevents overwriting the error (also, if there were an error with deployment, converting to an abstract token pool wouldn't be useful)
				tp, err = token_pool.NewTokenPool(tpAddr, chain.Client)
			}
			return cldf.ContractDeploy[*token_pool.TokenPool]{
				Address:  tpAddr,
				Contract: tp,
				Tv:       tv,
				Tx:       tx,
				Err:      err,
			}
		},
	)
}
