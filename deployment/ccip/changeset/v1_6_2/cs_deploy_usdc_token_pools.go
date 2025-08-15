package v1_6_2

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_2/hybrid_lock_release_usdc_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_2/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/erc20"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

var (
	DeployUSDCTokenPoolNew       = cldf.CreateChangeSet(deployUSDCTokenPoolContractsLogic, deployUSDCTokenPoolContractsPrecondition)
	USDCTokenPoolSentinelAddress = common.HexToAddress("0x0000000000000000000000000000000123456789")
)

// DeployUSDCTokenPoolInput defines all information required of the user to deploy a new USDC token pool contract.
type DeployUSDCTokenPoolInput struct {
	// PreviousPoolAddress is the address of the previous USDC token pool contract, inflight messages
	// are redirected to the previous pool when needed.
	PreviousPoolAddress common.Address
	// TokenMessenger is the address of the USDC token messenger contract.
	TokenMessenger common.Address
	// USDCTokenAddress is the address of the USDC token for which we are deploying a token pool.
	TokenAddress common.Address
	// PoolType is used to determine which type of USDC token pool to deploy.
	PoolType cldf.ContractType
	// AllowList is the optional list of addresses permitted to initiate a token transfer.
	// If omitted, all addresses will be permitted to transfer the token.
	AllowList []common.Address
}

func (i DeployUSDCTokenPoolInput) Validate(ctx context.Context, chain cldf_evm.Chain, state evm.CCIPChainState) error {
	// Ensure that required fields are populated
	if i.TokenAddress == utils.ZeroAddress {
		return errors.New("token address must be defined")
	}
	if i.TokenMessenger == utils.ZeroAddress {
		return errors.New("token messenger must be defined")
	}
	if _, ok := state.CCTPMessageTransmitterProxies[deployment.Version1_6_2]; !ok {
		// Note: This could be deployed automatically if it doesn't exist.
		return fmt.Errorf("CCTP message transmitter proxy for version %s not found on %s", deployment.Version1_6_2, chain)
	}
	if i.PreviousPoolAddress == utils.ZeroAddress {
		if len(state.USDCTokenPools) == 0 && len(state.USDCTokenPoolsV1_6) == 0 {
			return fmt.Errorf("unable to find a previous pool address, specify address or use USDCTokenPoolSentinelAddress (%s) if this is the first USDC token pool", USDCTokenPoolSentinelAddress)
		}
	}

	// Validate the token exists and matches the USDC symbol
	token, err := erc20.NewERC20(i.TokenAddress, chain.Client)
	if err != nil {
		return fmt.Errorf("failed to connect address %s with erc20 bindings: %w", i.TokenAddress, err)
	}
	symbol, err := token.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to fetch symbol from token with address %s: %w", i.TokenAddress, err)
	}
	if symbol != string(shared.USDCSymbol) {
		return fmt.Errorf("symbol of token with address %s (%s) is not USDC", i.TokenAddress, symbol)
	}

	// Check if a USDC token pool with the given version already exists
	if _, ok := state.USDCTokenPoolsV1_6[deployment.Version1_6_2]; ok {
		return fmt.Errorf("USDC token pool with version %s already exists on %s", deployment.Version1_6_2, chain)
	}

	// Perform USDC checks (i.e. make sure we can call the required functions)
	// LocalMessageTransmitter and MessageBodyVersion are called in the contract constructor:
	// https://github.com/smartcontractkit/chainlink/blob/f52a57762643b9cdc8e9241737e13501a4278716/contracts/src/v0.8/ccip/pools/USDC/USDCTokenPool.sol#L83
	messenger, err := mock_usdc_token_messenger.NewMockE2EUSDCTokenMessenger(i.TokenMessenger, chain.Client)
	if err != nil {
		return fmt.Errorf("failed to connect address %s on %s with token messenger bindings: %w", i.TokenMessenger, chain, err)
	}
	_, err = messenger.LocalMessageTransmitter(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to fetch local message transmitter from address %s on %s: %w", i.TokenMessenger, chain, err)
	}
	_, err = messenger.MessageBodyVersion(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to fetch message body version from address %s on %s: %w", i.TokenMessenger, chain, err)
	}

	if i.PoolType != shared.USDCTokenPool && i.PoolType != shared.HybridLockReleaseUSDCTokenPool {
		return fmt.Errorf("unsupported pool type %s", i.PoolType)
	}

	return nil
}

// DeployUSDCTokenPoolContractsConfig defines the USDC token pool contracts that need to be deployed on each chain.
type DeployUSDCTokenPoolContractsConfig struct {
	// USDCPools defines the per-chain configuration of each new USDC pool.
	USDCPools    map[uint64]DeployUSDCTokenPoolInput
	IsTestRouter bool
}

func deployUSDCTokenPoolContractsPrecondition(env cldf.Environment, c DeployUSDCTokenPoolContractsConfig) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, poolConfig := range c.USDCPools {
		chain, chainState, err := state.GetEVMChainState(env, chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get EVM chain state for chain selector %d: %w", chainSelector, err)
		}
		if !c.IsTestRouter && chainState.Router == nil {
			return fmt.Errorf("missing router on %s", chain)
		}
		if c.IsTestRouter && chainState.TestRouter == nil {
			return fmt.Errorf("missing test router on %s", chain)
		}
		err = poolConfig.Validate(env.GetContext(), chain, chainState)
		if err != nil {
			return fmt.Errorf("failed to validate USDC token pool config for chain selector %d: %w", chainSelector, err)
		}
	}
	return nil
}

// DeployUSDCTokenPoolContractsChangeset deploys new USDC pools across multiple chains.
func deployUSDCTokenPoolContractsLogic(env cldf.Environment, c DeployUSDCTokenPoolContractsConfig) (cldf.ChangesetOutput, error) {
	if err := deployUSDCTokenPoolContractsPrecondition(env, c); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid DeployUSDCTokenPoolContractsConfig: %w", err)
	}
	newAddresses := cldf.NewMemoryAddressBook()

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, poolConfig := range c.USDCPools {
		chain, chainState, err := state.GetEVMChainState(env, chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get EVM chain state for chain selector %d: %w", chainSelector, err)
		}
		router := chainState.Router
		if c.IsTestRouter {
			router = chainState.TestRouter
		}

		var deployErr error
		switch poolConfig.PoolType {
		case shared.USDCTokenPool:
			deployErr = deployUSDCTokenPool(env.Logger, chain, newAddresses, poolConfig, chainState, router.Address())
		case shared.HybridLockReleaseUSDCTokenPool:
			deployErr = deployHybridLockReleaseUSDCTokenPool(env.Logger, chain, newAddresses, poolConfig, chainState, router.Address())
		default:
			return cldf.ChangesetOutput{},
				fmt.Errorf("failed to deploy %s on %s: unknown pool type", poolConfig.PoolType, chain)
		}
		if deployErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy USDC token pool on %s: %w", chain, deployErr)
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: newAddresses, // TODO: this is deprecated, how do I use the DataStore instead?
	}, nil
}

func deployUSDCTokenPool(lggr logger.Logger, chain cldf_evm.Chain, newAddresses *cldf.AddressBookMap, poolConfig DeployUSDCTokenPoolInput, chainState evm.CCIPChainState, routerAddr common.Address) error {
	_, err := cldf.DeployContract(lggr, chain, newAddresses,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*usdc_token_pool.USDCTokenPool] {
			previousPoolAddress := poolConfig.PreviousPoolAddress

			switch previousPoolAddress {
			case USDCTokenPoolSentinelAddress:
				// If the previous pool address is USDCTokenPoolSentinelAddress, this is the first usdc token
				// pool and the address should be set to the ZeroAddress.
				// set the previous address to zero address.
				previousPoolAddress = utils.ZeroAddress

			case utils.ZeroAddress:
				// If the previous pool address is not set, we try to find the latest deployed pool address
				var err error
				previousPoolAddress, err = getPreviousPoolAddress(chainState, chain.Name())
				if err != nil {
					return cldf.ContractDeploy[*usdc_token_pool.USDCTokenPool]{Err: err}
				}
			}

			poolAddress, tx, usdcTokenPool, err := usdc_token_pool.DeployUSDCTokenPool(chain.DeployerKey,
				chain.Client, poolConfig.TokenMessenger,
				chainState.CCTPMessageTransmitterProxies[deployment.Version1_6_2].Address(),
				poolConfig.TokenAddress, poolConfig.AllowList, chainState.RMNProxy.Address(), routerAddr,
				previousPoolAddress)
			return cldf.ContractDeploy[*usdc_token_pool.USDCTokenPool]{
				Address:  poolAddress,
				Contract: usdcTokenPool,
				Tv:       cldf.NewTypeAndVersion(shared.USDCTokenPool, deployment.Version1_6_2),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	return err
}

func deployHybridLockReleaseUSDCTokenPool(lggr logger.Logger, chain cldf_evm.Chain, newAddresses *cldf.AddressBookMap, poolConfig DeployUSDCTokenPoolInput, chainState evm.CCIPChainState, routerAddr common.Address) error {
	_, err := cldf.DeployContract(lggr, chain, newAddresses,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*hybrid_lock_release_usdc_token_pool.HybridLockReleaseUSDCTokenPool] {
			previousPoolAddress := poolConfig.PreviousPoolAddress

			switch previousPoolAddress {
			case USDCTokenPoolSentinelAddress:
				// If the previous pool address is USDCTokenPoolSentinelAddress, this is the first usdc token
				// pool and the address should be set to the ZeroAddress.
				// set the previous address to zero address.
				previousPoolAddress = utils.ZeroAddress

			case utils.ZeroAddress:
				// If the previous pool address is not set, we try to find the latest deployed pool address
				var err error
				previousPoolAddress, err = getPreviousPoolAddress(chainState, chain.Name())
				if err != nil {
					return cldf.ContractDeploy[*hybrid_lock_release_usdc_token_pool.HybridLockReleaseUSDCTokenPool]{Err: err}
				}
			}

			poolAddress, tx, usdcTokenPool, err := hybrid_lock_release_usdc_token_pool.DeployHybridLockReleaseUSDCTokenPool(chain.DeployerKey,
				chain.Client, poolConfig.TokenMessenger,
				chainState.CCTPMessageTransmitterProxies[deployment.Version1_6_2].Address(),
				poolConfig.TokenAddress, poolConfig.AllowList, chainState.RMNProxy.Address(), routerAddr,
				previousPoolAddress)
			return cldf.ContractDeploy[*hybrid_lock_release_usdc_token_pool.HybridLockReleaseUSDCTokenPool]{
				Address:  poolAddress,
				Contract: usdcTokenPool,
				Tv:       cldf.NewTypeAndVersion(shared.HybridLockReleaseUSDCTokenPool, deployment.Version1_6_2),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	return err
}

func getPreviousPoolAddress(chainState evm.CCIPChainState, chainName string) (common.Address, error) {
	var previousPoolAddress common.Address
	switch {
	case chainState.USDCTokenPoolsV1_6[deployment.Version1_6_2] == nil:
		previousPoolAddress = chainState.USDCTokenPoolsV1_6[deployment.Version1_6_2].Address()
	case chainState.USDCTokenPools[deployment.Version1_5_1] == nil:
		previousPoolAddress = chainState.USDCTokenPools[deployment.Version1_5_1].Address()
	case chainState.USDCTokenPools[deployment.Version1_5_0] == nil:
		previousPoolAddress = chainState.USDCTokenPools[deployment.Version1_5_0].Address()
	default:
		return common.Address{}, fmt.Errorf("previous USDC pool address (%s) not found on %s", previousPoolAddress.Hex(), chainName)
	}
	return previousPoolAddress, nil
}
