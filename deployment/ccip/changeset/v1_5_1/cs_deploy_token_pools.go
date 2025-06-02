package v1_5_1

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_from_mint_token_pool"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_with_from_mint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"
)

var _ cldf.ChangeSet[DeployTokenPoolContractsConfig] = DeployTokenPoolContractsChangeset

// DeployTokenPoolInput defines all information required of the user to deploy a new token pool contract.
type DeployTokenPoolInput struct {
	// Type is the type of token pool that must be deployed.
	Type cldf.ContractType
	// TokenAddress is the address of the token for which we are deploying a pool.
	TokenAddress common.Address
	// AllowList is the optional list of addresses permitted to initiate a token transfer.
	// If omitted, all addresses will be permitted to transfer the token.
	AllowList []common.Address
	// LocalTokenDecimals is the number of decimals used by the token at tokenAddress.
	LocalTokenDecimals uint8
	// AcceptLiquidity indicates whether or not the new pool can accept liquidity from a rebalancer address (lock-release only).
	AcceptLiquidity *bool
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
}

func (c DeployTokenPoolContractsConfig) Validate(env cldf.Environment) error {
	// Ensure that required fields are populated
	if c.TokenSymbol == shared.TokenSymbol("") {
		return errors.New("token symbol must be defined")
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

	deployGrp := errgroup.Group{}

	for chainSelector, poolConfig := range c.NewPools {
		chainSelector, poolConfig := chainSelector, poolConfig
		deployGrp.Go(func() error {
			chain := env.BlockChains.EVMChains()[chainSelector]
			chainState := state.Chains[chainSelector]
			_, err := deployTokenPool(env.Logger, chain, chainState, newAddresses, poolConfig, c.IsTestRouter)
			if err != nil {
				return fmt.Errorf("failed to deploy token pool contract: %w", err)
			}
			return nil
		})
	}

	if err := deployGrp.Wait(); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy %s token pool on %w",
			c.TokenSymbol, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

// deployTokenPool deploys a token pool contract based on a given type & configuration.
func deployTokenPool(
	logger logger.Logger,
	chain cldf_evm.Chain,
	chainState evm.CCIPChainState,
	addressBook cldf.AddressBook,
	poolConfig DeployTokenPoolInput,
	isTestRouter bool,
) (*cldf.ContractDeploy[*token_pool.TokenPool], error) {
	router := chainState.Router
	if isTestRouter {
		router = chainState.TestRouter
	}
	rmnProxy := chainState.RMNProxy

	return cldf.DeployContract(logger, chain, addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*token_pool.TokenPool] {
			var tpAddr common.Address
			var tx *types.Transaction
			var err error
			switch poolConfig.Type {
			case shared.BurnMintTokenPool:
				tpAddr, tx, _, err = burn_mint_token_pool.DeployBurnMintTokenPool(
					chain.DeployerKey, chain.Client, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
					poolConfig.AllowList, rmnProxy.Address(), router.Address(),
				)
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
				tpAddr, tx, _, err = lock_release_token_pool.DeployLockReleaseTokenPool(
					chain.DeployerKey, chain.Client, poolConfig.TokenAddress, poolConfig.LocalTokenDecimals,
					poolConfig.AllowList, rmnProxy.Address(), *poolConfig.AcceptLiquidity, router.Address(),
				)
			}
			var tp *token_pool.TokenPool
			if err == nil { // prevents overwriting the error (also, if there were an error with deployment, converting to an abstract token pool wouldn't be useful)
				tp, err = token_pool.NewTokenPool(tpAddr, chain.Client)
			}
			return cldf.ContractDeploy[*token_pool.TokenPool]{
				Address:  tpAddr,
				Contract: tp,
				Tv:       cldf.NewTypeAndVersion(poolConfig.Type, shared.CurrentTokenPoolVersion),
				Tx:       tx,
				Err:      err,
			}
		},
	)
}
