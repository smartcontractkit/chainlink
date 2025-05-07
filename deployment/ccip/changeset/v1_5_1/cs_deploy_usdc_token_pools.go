package v1_5_1

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
)

var _ deployment.ChangeSet[DeployUSDCTokenPoolContractsConfig] = DeployUSDCTokenPoolContractsChangeset

// DeployUSDCTokenPoolInput defines all information required of the user to deploy a new USDC token pool contract.
type DeployUSDCTokenPoolInput struct {
	// TokenMessenger is the address of the USDC token messenger contract.
	TokenMessenger common.Address
	// USDCTokenAddress is the address of the USDC token for which we are deploying a token pool.
	TokenAddress common.Address
	// AllowList is the optional list of addresses permitted to initiate a token transfer.
	// If omitted, all addresses will be permitted to transfer the token.
	AllowList []common.Address
}

func (i DeployUSDCTokenPoolInput) Validate(ctx context.Context, chain deployment.Chain, state changeset.CCIPChainState) error {
	// Ensure that required fields are populated
	if i.TokenAddress == utils.ZeroAddress {
		return errors.New("token address must be defined")
	}
	if i.TokenMessenger == utils.ZeroAddress {
		return errors.New("token messenger must be defined")
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
	if symbol != string(changeset.USDCSymbol) {
		return fmt.Errorf("symbol of token with address %s (%s) is not USDC", i.TokenAddress, symbol)
	}

	// Check if a USDC token pool with the given version already exists
	fmt.Println(state.USDCTokenPools[deployment.Version1_5_1])
	if _, ok := state.USDCTokenPools[deployment.Version1_5_1]; ok {
		return fmt.Errorf("USDC token pool with version %s already exists on %s", deployment.Version1_5_1, chain)
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

	return nil
}

// DeployUSDCTokenPoolContractsConfig defines the USDC token pool contracts that need to be deployed on each chain.
type DeployUSDCTokenPoolContractsConfig struct {
	// USDCPools defines the per-chain configuration of each new USDC pool.
	USDCPools    map[uint64]DeployUSDCTokenPoolInput
	IsTestRouter bool
}

func (c DeployUSDCTokenPoolContractsConfig) Validate(env deployment.Environment) error {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for chainSelector, poolConfig := range c.USDCPools {
		err := deployment.IsValidChainSelector(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to validate chain selector %d: %w", chainSelector, err)
		}
		chain, ok := env.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
		}
		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("chain with selector %d does not exist in state", chainSelector)
		}
		if !c.IsTestRouter && chainState.Router == nil {
			return fmt.Errorf("missing router on %s", chain)
		}
		if c.IsTestRouter && chainState.TestRouter == nil {
			return fmt.Errorf("missing test router on %s", chain)
		}
		if chainState.RMNProxy == nil {
			return fmt.Errorf("missing rmnProxy on %s", chain)
		}
		err = poolConfig.Validate(env.GetContext(), chain, chainState)
		if err != nil {
			return fmt.Errorf("failed to validate USDC token pool config for chain selector %d: %w", chainSelector, err)
		}
	}
	return nil
}

// DeployUSDCTokenPoolContractsChangeset deploys new USDC pools across multiple chains.
func DeployUSDCTokenPoolContractsChangeset(env deployment.Environment, c DeployUSDCTokenPoolContractsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid DeployUSDCTokenPoolContractsConfig: %w", err)
	}
	newAddresses := deployment.NewMemoryAddressBook()

	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, poolConfig := range c.USDCPools {
		chain := env.Chains[chainSelector]
		chainState := state.Chains[chainSelector]
		router := chainState.Router
		if c.IsTestRouter {
			router = chainState.TestRouter
		}
		_, err := cldf.DeployContract(env.Logger, chain, newAddresses,
			func(chain deployment.Chain) cldf.ContractDeploy[*usdc_token_pool.USDCTokenPool] {
				poolAddress, tx, usdcTokenPool, err := usdc_token_pool.DeployUSDCTokenPool(
					chain.DeployerKey, chain.Client, poolConfig.TokenMessenger, poolConfig.TokenAddress,
					poolConfig.AllowList, chainState.RMNProxy.Address(), router.Address(),
				)
				return cldf.ContractDeploy[*usdc_token_pool.USDCTokenPool]{
					Address:  poolAddress,
					Contract: usdcTokenPool,
					Tv:       deployment.NewTypeAndVersion(changeset.USDCTokenPool, deployment.Version1_5_1),
					Tx:       tx,
					Err:      err,
				}
			},
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy USDC token pool on %s: %w", chain, err)
		}
	}

	return deployment.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}
