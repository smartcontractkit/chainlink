package testhelpers

import (
	"math/big"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_5_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

const (
	LocalTokenDecimals                    = 18
	TestTokenSymbol    shared.TokenSymbol = "TEST"
)

// CreateSymmetricRateLimits is a utility to quickly create a rate limiter config with equal inbound and outbound values.
func CreateSymmetricRateLimits(rate int64, capacity int64) v1_5_1.RateLimiterConfig {
	return v1_5_1.RateLimiterConfig{
		Inbound: token_pool.RateLimiterConfig{
			IsEnabled: rate != 0 || capacity != 0,
			Rate:      big.NewInt(rate),
			Capacity:  big.NewInt(capacity),
		},
		Outbound: token_pool.RateLimiterConfig{
			IsEnabled: rate != 0 || capacity != 0,
			Rate:      big.NewInt(rate),
			Capacity:  big.NewInt(capacity),
		},
	}
}

// SetupTwoChainEnvironmentWithTokens preps the environment for token pool deployment testing.
func SetupTwoChainEnvironmentWithTokens(
	t *testing.T,
	lggr logger.Logger,
	transferToTimelock bool,
) (env cldf.Environment, sel1 uint64, sel2 uint64, ercmap map[uint64]*cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677], contractmap map[uint64]*proposalutils.TimelockExecutionContracts) {
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})
	selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))

	addressBook := cldf.NewMemoryAddressBook()
	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, len(selectors))
	for i, selector := range selectors {
		prereqCfg[i] = changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: selector,
		}
	}

	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfigV2)
	for _, selector := range selectors {
		mcmsCfg[selector] = proposalutils.SingleGroupTimelockConfigV2(t)
	}

	// Deploy one burn-mint token per chain to use in the tests
	tokens := make(map[uint64]*cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677])
	for _, selector := range selectors {
		token, err := cldf.DeployContract(e.Logger, e.BlockChains.EVMChains()[selector], addressBook,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
				tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
					e.BlockChains.EVMChains()[selector].DeployerKey,
					e.BlockChains.EVMChains()[selector].Client,
					string(TestTokenSymbol),
					string(TestTokenSymbol),
					LocalTokenDecimals,
					big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
				)
				return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
					Address:  tokenAddress,
					Contract: token,
					Tv:       cldf.NewTypeAndVersion(shared.BurnMintToken, deployment.Version1_0_0),
					Tx:       tx,
					Err:      err,
				}
			},
		)
		require.NoError(t, err)
		tokens[selector] = token
	}

	// Deploy MCMS setup & prerequisite contracts
	e, err := commoncs.Apply(t, e, nil,
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset),
			changeset.DeployPrerequisiteConfig{Configs: prereqCfg},
		),
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(commoncs.DeployMCMSWithTimelockV2),
			mcmsCfg,
		),
	)
	require.NoError(t, err)

	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	// We only need the token admin registry to be owned by the timelock in these tests
	timelockOwnedContractsByChain := make(map[uint64][]common.Address)
	for _, selector := range selectors {
		timelockOwnedContractsByChain[selector] = []common.Address{state.MustGetEVMChainState(selector).TokenAdminRegistry.Address()}
	}

	// Assemble map of addresses required for Timelock scheduling & execution
	timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	for _, selector := range selectors {
		timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.MustGetEVMChainState(selector).Timelock,
			CallProxy: state.MustGetEVMChainState(selector).CallProxy,
		}
	}

	if transferToTimelock {
		// Transfer ownership of token admin registry to the Timelock
		e, err = commoncs.Apply(t, e, timelockContracts,
			commoncs.Configure(
				cldf.CreateLegacyChangeSet(commoncs.TransferToMCMSWithTimelockV2),
				commoncs.TransferToMCMSWithTimelockConfig{
					ContractsByChain: timelockOwnedContractsByChain,
					MCMSConfig: proposalutils.TimelockConfig{
						MinDelay: 0 * time.Second,
					},
				},
			),
		)
		require.NoError(t, err)
	}

	return e, selectors[0], selectors[1], tokens, timelockContracts
}

// getPoolsOwnedByDeployer returns any pools that need to be transferred to timelock.
func getPoolsOwnedByDeployer[T commonchangeset.Ownable](t *testing.T, contracts map[semver.Version]T, chain cldf_evm.Chain) []common.Address {
	var addresses []common.Address
	for _, contract := range contracts {
		owner, err := contract.Owner(nil)
		require.NoError(t, err)
		if owner == chain.DeployerKey.From {
			addresses = append(addresses, contract.Address())
		}
	}
	return addresses
}

// DeployTestTokenPools deploys token pools tied for the TEST token across multiple chains.
func DeployTestTokenPools(
	t *testing.T,
	e cldf.Environment,
	newPools map[uint64]v1_5_1.DeployTokenPoolInput,
	transferToTimelock bool,
) cldf.Environment {
	selectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))

	e, err := commonchangeset.Apply(t, e, nil,
		commoncs.Configure(
			cldf.CreateLegacyChangeSet(v1_5_1.DeployTokenPoolContractsChangeset),
			v1_5_1.DeployTokenPoolContractsConfig{
				TokenSymbol: TestTokenSymbol,
				NewPools:    newPools,
			},
		),
	)
	require.NoError(t, err)

	state, err := stateview.LoadOnchainState(e)
	require.NoError(t, err)

	if transferToTimelock {
		// Assemble map of addresses required for Timelock scheduling & execution
		timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts)
		for _, selector := range selectors {
			timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
				Timelock:  state.MustGetEVMChainState(selector).Timelock,
				CallProxy: state.MustGetEVMChainState(selector).CallProxy,
			}
		}

		timelockOwnedContractsByChain := make(map[uint64][]common.Address)
		for _, selector := range selectors {
			if newPool, ok := newPools[selector]; ok {
				switch newPool.Type {
				case shared.BurnFromMintTokenPool:
					timelockOwnedContractsByChain[selector] = getPoolsOwnedByDeployer(t, state.MustGetEVMChainState(selector).BurnFromMintTokenPools[TestTokenSymbol], e.BlockChains.EVMChains()[selector])
				case shared.BurnWithFromMintTokenPool:
					timelockOwnedContractsByChain[selector] = getPoolsOwnedByDeployer(t, state.MustGetEVMChainState(selector).BurnWithFromMintTokenPools[TestTokenSymbol], e.BlockChains.EVMChains()[selector])
				case shared.BurnMintTokenPool:
					timelockOwnedContractsByChain[selector] = getPoolsOwnedByDeployer(t, state.MustGetEVMChainState(selector).BurnMintTokenPools[TestTokenSymbol], e.BlockChains.EVMChains()[selector])
				case shared.LockReleaseTokenPool:
					timelockOwnedContractsByChain[selector] = getPoolsOwnedByDeployer(t, state.MustGetEVMChainState(selector).LockReleaseTokenPools[TestTokenSymbol], e.BlockChains.EVMChains()[selector])
				}
			}
		}

		// Transfer ownership of token admin registry to the Timelock
		e, err = commoncs.Apply(t, e, timelockContracts,
			commoncs.Configure(
				cldf.CreateLegacyChangeSet(commoncs.TransferToMCMSWithTimelockV2),
				commoncs.TransferToMCMSWithTimelockConfig{
					ContractsByChain: timelockOwnedContractsByChain,
					MCMSConfig: proposalutils.TimelockConfig{
						MinDelay: 0 * time.Second,
					},
				},
			),
		)
		require.NoError(t, err)
	}

	return e
}
