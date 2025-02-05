package testhelpers

import (
	"math/big"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	testlogger "github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	LocalTokenDecimals                       = 18
	TestTokenSymbol    changeset.TokenSymbol = "TEST"
)

// SetupTwoChainEnvironmentWithTokens preps the environment for token pool deployment testing.
func SetupTwoChainEnvironmentWithTokens(
	t *testing.T,
	lggr logger.Logger,
	transferToTimelock bool,
) (deployment.Environment, uint64, uint64, map[uint64]*deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677], map[uint64]*proposalutils.TimelockExecutionContracts) {
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 2,
	})
	selectors := e.AllChainSelectors()

	addressBook := deployment.NewMemoryAddressBook()
	prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, len(selectors))
	for i, selector := range selectors {
		prereqCfg[i] = changeset.DeployPrerequisiteConfigPerChain{
			ChainSelector: selector,
		}
	}

	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, selector := range selectors {
		mcmsCfg[selector] = proposalutils.SingleGroupTimelockConfig(t)
	}

	// Deploy one burn-mint token per chain to use in the tests
	tokens := make(map[uint64]*deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677])
	for _, selector := range selectors {
		token, err := deployment.DeployContract(e.Logger, e.Chains[selector], addressBook,
			func(chain deployment.Chain) deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
				tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
					e.Chains[selector].DeployerKey,
					e.Chains[selector].Client,
					string(TestTokenSymbol),
					string(TestTokenSymbol),
					LocalTokenDecimals,
					big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
				)
				return deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
					Address:  tokenAddress,
					Contract: token,
					Tv:       deployment.NewTypeAndVersion(changeset.BurnMintToken, deployment.Version1_0_0),
					Tx:       tx,
					Err:      err,
				}
			},
		)
		require.NoError(t, err)
		tokens[selector] = token
	}

	// Deploy MCMS setup & prerequisite contracts
	e, err := commoncs.ApplyChangesets(t, e, nil, []commoncs.ChangesetApplication{
		{
			Changeset: commoncs.WrapChangeSet(changeset.DeployPrerequisitesChangeset),
			Config: changeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			},
		},
		{
			Changeset: commoncs.WrapChangeSet(commoncs.DeployMCMSWithTimelock),
			Config:    mcmsCfg,
		},
	})
	require.NoError(t, err)

	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	// We only need the token admin registry to be owned by the timelock in these tests
	timelockOwnedContractsByChain := make(map[uint64][]common.Address)
	for _, selector := range selectors {
		timelockOwnedContractsByChain[selector] = []common.Address{state.Chains[selector].TokenAdminRegistry.Address()}
	}

	// Assemble map of addresses required for Timelock scheduling & execution
	timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts)
	for _, selector := range selectors {
		timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
			Timelock:  state.Chains[selector].Timelock,
			CallProxy: state.Chains[selector].CallProxy,
		}
	}

	if transferToTimelock {
		// Transfer ownership of token admin registry to the Timelock
		e, err = commoncs.ApplyChangesets(t, e, timelockContracts, []commoncs.ChangesetApplication{
			{
				Changeset: commoncs.WrapChangeSet(commoncs.TransferToMCMSWithTimelock),
				Config: commoncs.TransferToMCMSWithTimelockConfig{
					ContractsByChain: timelockOwnedContractsByChain,
					MinDelay:         0,
				},
			},
		})
		require.NoError(t, err)
	}

	return e, selectors[0], selectors[1], tokens, timelockContracts
}

// getPoolsOwnedByDeployer returns any pools that need to be transferred to timelock.
func getPoolsOwnedByDeployer[T commonchangeset.Ownable](t *testing.T, contracts map[semver.Version]T, chain deployment.Chain) []common.Address {
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
	e deployment.Environment,
	newPools map[uint64]changeset.DeployTokenPoolInput,
	transferToTimelock bool,
) deployment.Environment {
	selectors := e.AllChainSelectors()

	e, err := commonchangeset.ApplyChangesets(t, e, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployTokenPoolContractsChangeset),
			Config: changeset.DeployTokenPoolContractsConfig{
				TokenSymbol: TestTokenSymbol,
				NewPools:    newPools,
			},
		},
	})
	require.NoError(t, err)

	state, err := changeset.LoadOnchainState(e)
	require.NoError(t, err)

	if transferToTimelock {
		// Assemble map of addresses required for Timelock scheduling & execution
		timelockContracts := make(map[uint64]*proposalutils.TimelockExecutionContracts)
		for _, selector := range selectors {
			timelockContracts[selector] = &proposalutils.TimelockExecutionContracts{
				Timelock:  state.Chains[selector].Timelock,
				CallProxy: state.Chains[selector].CallProxy,
			}
		}

		timelockOwnedContractsByChain := make(map[uint64][]common.Address)
		for _, selector := range selectors {
			if newPool, ok := newPools[selector]; ok {
				switch newPool.Type {
				case changeset.BurnFromMintTokenPool:
					timelockOwnedContractsByChain[selector] = getPoolsOwnedByDeployer(t, state.Chains[selector].BurnFromMintTokenPools[TestTokenSymbol], e.Chains[selector])
				case changeset.BurnWithFromMintTokenPool:
					timelockOwnedContractsByChain[selector] = getPoolsOwnedByDeployer(t, state.Chains[selector].BurnWithFromMintTokenPools[TestTokenSymbol], e.Chains[selector])
				case changeset.BurnMintTokenPool:
					timelockOwnedContractsByChain[selector] = getPoolsOwnedByDeployer(t, state.Chains[selector].BurnMintTokenPools[TestTokenSymbol], e.Chains[selector])
				case changeset.LockReleaseTokenPool:
					timelockOwnedContractsByChain[selector] = getPoolsOwnedByDeployer(t, state.Chains[selector].LockReleaseTokenPools[TestTokenSymbol], e.Chains[selector])
				}
			}
		}

		// Transfer ownership of token admin registry to the Timelock
		e, err = commoncs.ApplyChangesets(t, e, timelockContracts, []commoncs.ChangesetApplication{
			{
				Changeset: commoncs.WrapChangeSet(commoncs.TransferToMCMSWithTimelock),
				Config: commoncs.TransferToMCMSWithTimelockConfig{
					ContractsByChain: timelockOwnedContractsByChain,
					MinDelay:         0,
				},
			},
		})
		require.NoError(t, err)
	}

	return e
}

// SetupTokens deploys transferable tokens on the source and dest, mints tokens for the source and dest, and
// approves the router to spend the tokens
func SetupTransferableTokens(
	t *testing.T,
	state changeset.CCIPOnChainState,
	tenv DeployedEnv,
	src, dest uint64,
	transferTokenMintAmount,
	feeTokenMintAmount *big.Int,
) (
	srcToken *burn_mint_erc677.BurnMintERC677,
	dstToken *burn_mint_erc677.BurnMintERC677,
) {
	lggr := testlogger.TestLogger(t)
	e := tenv.Env

	// Deploy the token to test transferring
	srcToken, _, dstToken, _, err := DeployTransferableToken(
		lggr,
		tenv.Env.Chains,
		src,
		dest,
		tenv.Env.Chains[src].DeployerKey,
		tenv.Env.Chains[dest].DeployerKey,
		state,
		tenv.Env.ExistingAddresses,
		"MY_TOKEN",
	)
	require.NoError(t, err)

	linkToken := state.Chains[src].LinkToken

	tx, err := srcToken.Mint(
		e.Chains[src].DeployerKey,
		e.Chains[src].DeployerKey.From,
		transferTokenMintAmount,
	)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	// Mint a destination token
	tx, err = dstToken.Mint(
		e.Chains[dest].DeployerKey,
		e.Chains[dest].DeployerKey.From,
		transferTokenMintAmount,
	)
	_, err = deployment.ConfirmIfNoError(e.Chains[dest], tx, err)
	require.NoError(t, err)

	// Approve the router to spend the tokens and confirm the tx's
	// To prevent having to approve the router for every transfer, we approve a sufficiently large amount
	tx, err = srcToken.Approve(e.Chains[src].DeployerKey, state.Chains[src].Router.Address(), math.MaxBig256)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	tx, err = dstToken.Approve(e.Chains[dest].DeployerKey, state.Chains[dest].Router.Address(), math.MaxBig256)
	_, err = deployment.ConfirmIfNoError(e.Chains[dest], tx, err)
	require.NoError(t, err)

	// Grant mint and burn roles to the deployer key for the newly deployed linkToken
	// Since those roles are not granted automatically
	tx, err = linkToken.GrantMintAndBurnRoles(e.Chains[src].DeployerKey, e.Chains[src].DeployerKey.From)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	// Mint link token and confirm the tx
	tx, err = linkToken.Mint(
		e.Chains[src].DeployerKey,
		e.Chains[src].DeployerKey.From,
		feeTokenMintAmount,
	)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	require.NoError(t, err)

	return srcToken, dstToken
}
