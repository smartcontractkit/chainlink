package v1_6_2

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ccipopsv1_5_1 "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_5_1"
	ccipopsv1_6_2 "github.com/smartcontractkit/chainlink/deployment/ccip/operation/evm/v1_6_2"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"golang.org/x/sync/errgroup"
)

type TokenAndPoolContractParams struct {
	Decimals        uint8
	AllowList       []common.Address
	AcceptLiquidity bool
	MaxSupply       *big.Int
	PreMint         *big.Int
	NewOwner        common.Address
}

func (c TokenAndPoolContractParams) Validate(selector uint64) error {
	if c.NewOwner == (common.Address{}) {
		return fmt.Errorf("new owner address cannot be empty for selector %d", selector)
	}

	return nil
}

type DeployFactoryBurnMintERC20PoolsConfig struct {
	TokenAndPoolContractParamsPerChain map[uint64]TokenAndPoolContractParams
	GasBoostConfigPerChain             map[uint64]commontypes.GasBoostConfig
}

func (c DeployFactoryBurnMintERC20PoolsConfig) Validate() error {
	for cs, args := range c.TokenAndPoolContractParamsPerChain {
		if err := deployment.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
		if err := args.Validate(cs); err != nil {
			return fmt.Errorf("invalid contract args for chain %d: %w", cs, err)
		}
	}
	return nil
}

type CCIPContractAddresses struct {
	RMNProxyAddress common.Address
	RouterAddress   common.Address

	// New addresses deployed in this sequence
	FactoryBurnMintERC20Address      common.Address
	BurnMintTokenPoolAddress         common.Address
	BurnFromMintTokenPoolAddress     common.Address
	BurnWithFromMintTokenPoolAddress common.Address
	LockReleaseTokenPoolAddress      common.Address
}

func (c CCIPContractAddresses) Validate(selector uint64) error {
	if c.RMNProxyAddress == (common.Address{}) {
		return fmt.Errorf("rmn proxy address is not defined for chain with selector %d, deploy the ccip contracts first", selector)
	}
	if c.RouterAddress == (common.Address{}) {
		return fmt.Errorf("router address is not defined for chain with selector %d, deploy the ccip contracts first", selector)
	}

	return nil
}

type DeployFactoryBurnMintERC20PoolsSeqConfig struct {
	DeployFactoryBurnMintERC20PoolsConfig
	AddressesPerChain      map[uint64]CCIPContractAddresses
	GasBoostConfigPerChain map[uint64]commontypes.GasBoostConfig
}

func (c DeployFactoryBurnMintERC20PoolsSeqConfig) Validate() error {
	for chainSelector, addresses := range c.AddressesPerChain {
		if err := addresses.Validate(chainSelector); err != nil {
			return fmt.Errorf("invalid addresses for chain %d: %w", chainSelector, err)
		}
	}

	return nil
}

var (
	DeployFactoryBurnMintERC20PoolsSeq = operations.NewSequence(
		"DeployChainContractsSeq",
		semver.MustParse("1.0.0"),
		"Deploys the FactoryBurnMintERC20 token and a BurnMintTokenPool, a BurnFromMintTokenPool, a BurnWithFromMintTokenPool, and a LockReleaseTokenPool for the specified evm chain(s)",
		func(b operations.Bundle, deps map[uint64]cldf_evm.Chain, input DeployFactoryBurnMintERC20PoolsSeqConfig) (map[uint64]map[string]string, error) {
			err := input.Validate()
			if err != nil {
				return nil, fmt.Errorf("invalid DeployFactoryBurnMintERC20PoolsSeqConfig: %w", err)
			}

			gasBoostConfigs := opsutil.GasBoostConfigsForChainMap(input.TokenAndPoolContractParamsPerChain, input.GasBoostConfigPerChain)
			out := make(map[uint64]map[string]string)
			grp := errgroup.Group{}
			mu := sync.Mutex{}
			for chainSelector, contractParams := range input.TokenAndPoolContractParamsPerChain {
				newAddresses := make(map[string]string)
				chainSelector := chainSelector
				contractParams := contractParams
				chainAddresses := input.AddressesPerChain[chainSelector]
				chain, ok := deps[chainSelector]
				if !ok {
					return nil, fmt.Errorf("chain with selector %d not defined in dependencies", chainSelector)
				}

				var factoryBurnMintERC20Address common.Address
				var burnMintTokenPoolAddress common.Address
				var burnFromMintTokenPoolAddress common.Address
				var burnWithFromMintTokenPoolAddress common.Address
				var lockReleaseTokenPoolAddress common.Address

				grp.Go(func() error {
					// if FactoryBurnMintERC20 is not already deployed
					if chainAddresses.FactoryBurnMintERC20Address == (common.Address{}) {
						report, err := operations.ExecuteOperation(b, ccipopsv1_6_2.DeployFactoryBurnMintERC20TokenOp, chain, opsutil.EVMDeployInput[ccipopsv1_6_2.DeployFactoryBurnMintERC20TokenInput]{
							ChainSelector: chainSelector,
							DeployInput: ccipopsv1_6_2.DeployFactoryBurnMintERC20TokenInput{
								Name:      shared.FactoryBurnMintERC20Symbol.String(),
								Symbol:    shared.FactoryBurnMintERC20Symbol.String(),
								Decimals:  contractParams.Decimals,
								MaxSupply: contractParams.MaxSupply,
								PreMint:   contractParams.PreMint,
								NewOwner:  contractParams.NewOwner,
							},
						}, opsutil.RetryDeploymentWithGasBoost[ccipopsv1_6_2.DeployFactoryBurnMintERC20TokenInput](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to deploy FactoryBurnMintERC20 for %s: %w", chain, err)
						}
						factoryBurnMintERC20Address = report.Output.Address
						newAddresses[factoryBurnMintERC20Address.Hex()] = report.Output.TypeAndVersion
					} else {
						factoryBurnMintERC20Address = chainAddresses.FactoryBurnMintERC20Address
					}
					// if BurnMintTokenPool is not already deployed
					if chainAddresses.BurnMintTokenPoolAddress == (common.Address{}) {
						report, err := operations.ExecuteOperation(b, ccipopsv1_5_1.DeployBurnMintTokenPoolOp, chain, opsutil.EVMDeployInput[ccipopsv1_5_1.DeployBurnMintTokenPoolInput]{
							ChainSelector: chainSelector,
							DeployInput: ccipopsv1_5_1.DeployBurnMintTokenPoolInput{
								TokenAddress:    factoryBurnMintERC20Address,
								Decimals:        contractParams.Decimals,
								AllowList:       contractParams.AllowList,
								RMNProxyAddress: chainAddresses.RMNProxyAddress,
								RouterAddress:   chainAddresses.RouterAddress,
							},
						}, opsutil.RetryDeploymentWithGasBoost[ccipopsv1_5_1.DeployBurnMintTokenPoolInput](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to deploy BurnMintTokenPool for %s: %w", chain, err)
						}
						burnMintTokenPoolAddress = report.Output.Address
						newAddresses[burnMintTokenPoolAddress.Hex()] = report.Output.TypeAndVersion
					} else {
						burnMintTokenPoolAddress = chainAddresses.BurnMintTokenPoolAddress
					}
					// if BurnFromMintTokenPool is not already deployed
					if chainAddresses.BurnFromMintTokenPoolAddress == (common.Address{}) {
						report, err := operations.ExecuteOperation(b, ccipopsv1_5_1.DeployBurnFromMintTokenPoolOp, chain, opsutil.EVMDeployInput[ccipopsv1_5_1.DeployBurnMintTokenPoolInput]{
							ChainSelector: chainSelector,
							DeployInput: ccipopsv1_5_1.DeployBurnMintTokenPoolInput{
								TokenAddress:    factoryBurnMintERC20Address,
								Decimals:        contractParams.Decimals,
								AllowList:       contractParams.AllowList,
								RMNProxyAddress: chainAddresses.RMNProxyAddress,
								RouterAddress:   chainAddresses.RouterAddress,
							},
						}, opsutil.RetryDeploymentWithGasBoost[ccipopsv1_5_1.DeployBurnMintTokenPoolInput](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to deploy BurnFromMintTokenPool for %s: %w", chain, err)
						}
						burnFromMintTokenPoolAddress = report.Output.Address
						newAddresses[burnFromMintTokenPoolAddress.Hex()] = report.Output.TypeAndVersion
					} else {
						burnFromMintTokenPoolAddress = chainAddresses.BurnFromMintTokenPoolAddress
					}
					// if BurnWithFromMintTokenPool is not already deployed
					if chainAddresses.BurnWithFromMintTokenPoolAddress == (common.Address{}) {
						report, err := operations.ExecuteOperation(b, ccipopsv1_5_1.DeployBurnWithFromMintTokenPoolOp, chain, opsutil.EVMDeployInput[ccipopsv1_5_1.DeployBurnMintTokenPoolInput]{
							ChainSelector: chainSelector,
							DeployInput: ccipopsv1_5_1.DeployBurnMintTokenPoolInput{
								TokenAddress:    factoryBurnMintERC20Address,
								Decimals:        contractParams.Decimals,
								AllowList:       contractParams.AllowList,
								RMNProxyAddress: chainAddresses.RMNProxyAddress,
								RouterAddress:   chainAddresses.RouterAddress,
							},
						}, opsutil.RetryDeploymentWithGasBoost[ccipopsv1_5_1.DeployBurnMintTokenPoolInput](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to deploy BurnWithFromMintTokenPool for %s: %w", chain, err)
						}
						burnWithFromMintTokenPoolAddress = report.Output.Address
						newAddresses[burnWithFromMintTokenPoolAddress.Hex()] = report.Output.TypeAndVersion
					} else {
						burnWithFromMintTokenPoolAddress = chainAddresses.BurnWithFromMintTokenPoolAddress
					}
					// if LockReleaseTokenPool is not already deployed
					if chainAddresses.LockReleaseTokenPoolAddress == (common.Address{}) {
						report, err := operations.ExecuteOperation(b, ccipopsv1_5_1.DeployLockReleaseTokenPoolOp, chain, opsutil.EVMDeployInput[ccipopsv1_5_1.DeployLockReleaseTokenPoolInput]{
							ChainSelector: chainSelector,
							DeployInput: ccipopsv1_5_1.DeployLockReleaseTokenPoolInput{
								TokenAddress:    factoryBurnMintERC20Address,
								Decimals:        contractParams.Decimals,
								AllowList:       contractParams.AllowList,
								RMNProxyAddress: chainAddresses.RMNProxyAddress,
								AcceptLiquidity: contractParams.AcceptLiquidity,
								RouterAddress:   chainAddresses.RouterAddress,
							},
						}, opsutil.RetryDeploymentWithGasBoost[ccipopsv1_5_1.DeployLockReleaseTokenPoolInput](gasBoostConfigs[chainSelector]))
						if err != nil {
							return fmt.Errorf("failed to deploy LockReleaseTokenPool for %s: %w", chain, err)
						}
						lockReleaseTokenPoolAddress = report.Output.Address
						newAddresses[lockReleaseTokenPoolAddress.Hex()] = report.Output.TypeAndVersion
					} else {
						lockReleaseTokenPoolAddress = chainAddresses.LockReleaseTokenPoolAddress
					}

					mu.Lock()
					out[chainSelector] = newAddresses
					mu.Unlock()

					return nil
				})
			}
			if err := grp.Wait(); err != nil {
				return nil, fmt.Errorf("failed to deploy FactoryBurnMintERC20 and TokenPool contracts: %w", err)
			}
			return out, nil
		})
)
