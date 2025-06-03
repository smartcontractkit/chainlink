package deployergroup_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type mintConfig struct {
	selectorIndex uint64
	amount        *big.Int
}

type dummyMultiChainDeployerGroupChangesetConfig struct {
	address common.Address
	mints   []mintConfig
	MCMS    *proposalutils.TimelockConfig
}

type dummyDeployerGroupChangesetConfig struct {
	selector uint64
	address  common.Address
	mints    []*big.Int
	MCMS     *proposalutils.TimelockConfig
}

type dummyEmptyBatchChangesetConfig struct {
	MCMS *proposalutils.TimelockConfig
}

func dummyEmptyBatchChangeset(e cldf.Environment, cfg dummyEmptyBatchChangesetConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	group := deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("empty batch")
	return group.Enact()
}

func dummyDeployerGroupGrantMintChangeset(e cldf.Environment, cfg dummyDeployerGroupChangesetConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	token := state.MustGetEVMChainState(cfg.selector).LinkToken

	group := deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("grant mint role")
	deployer, err := group.GetDeployer(cfg.selector)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	_, err = token.GrantMintRole(deployer, deployer.From)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return group.Enact()
}

func dummyDeployerGroupMintChangeset(e cldf.Environment, cfg dummyDeployerGroupChangesetConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	token := state.MustGetEVMChainState(cfg.selector).LinkToken

	group := deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("mint tokens")
	deployer, err := group.GetDeployer(cfg.selector)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	for _, mint := range cfg.mints {
		_, err = token.Mint(deployer, cfg.address, mint)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}

	return group.Enact()
}

func dummyDeployerGroupGrantMintMultiChainChangeset(e cldf.Environment, cfg dummyMultiChainDeployerGroupChangesetConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	group := deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("grant mint role")
	for _, mint := range cfg.mints {
		selector := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[mint.selectorIndex]
		token := state.MustGetEVMChainState(selector).LinkToken

		deployer, err := group.GetDeployer(selector)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		_, err = token.GrantMintRole(deployer, deployer.From)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}

	return group.Enact()
}

func dummyDeployerGroupMintMultiDeploymentContextChangeset(e cldf.Environment, cfg dummyMultiChainDeployerGroupChangesetConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	var group *deployergroup.DeployerGroup
	var deployer *bind.TransactOpts

	for i, mint := range cfg.mints {
		selector := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[mint.selectorIndex]
		token := state.MustGetEVMChainState(selector).LinkToken

		if group == nil {
			group = deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext(fmt.Sprintf("mint tokens %d", i+1))
		} else {
			group = group.WithDeploymentContext(fmt.Sprintf("mint tokens %d", i+1))
		}
		deployer, err = group.GetDeployer(selector)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		_, err = token.Mint(deployer, cfg.address, mint.amount)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}

	return group.Enact()
}

type deployerGroupTestCase struct {
	name        string
	cfg         dummyDeployerGroupChangesetConfig
	expectError bool
}

var deployerGroupTestCases = []deployerGroupTestCase{
	{
		name: "happy path",
		cfg: dummyDeployerGroupChangesetConfig{
			mints:   []*big.Int{big.NewInt(1), big.NewInt(2)},
			address: common.HexToAddress("0x455E5AA18469bC6ccEF49594645666C587A3a71B"),
		},
	},
	{
		name: "error",
		cfg: dummyDeployerGroupChangesetConfig{
			mints:   []*big.Int{big.NewInt(-1)},
			address: common.HexToAddress("0x455E5AA18469bC6ccEF49594645666C587A3a71B"),
		},
		expectError: true,
	},
}

func TestDeployerGroup(t *testing.T) {
	for _, tc := range deployerGroupTestCases {
		t.Run(tc.name, func(t *testing.T) {
			e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))

			tc.cfg.selector = e.HomeChainSel
			tc.cfg.MCMS = nil

			_, err := dummyDeployerGroupGrantMintChangeset(e.Env, tc.cfg)
			require.NoError(t, err)

			_, err = dummyDeployerGroupMintChangeset(e.Env, tc.cfg)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				state, err := stateview.LoadOnchainState(e.Env)
				require.NoError(t, err)

				token := state.MustGetEVMChainState(e.HomeChainSel).LinkToken

				amount, err := token.BalanceOf(nil, tc.cfg.address)
				require.NoError(t, err)

				sumOfMints := big.NewInt(0)
				for _, mint := range tc.cfg.mints {
					sumOfMints = sumOfMints.Add(sumOfMints, mint)
				}

				require.Equal(t, sumOfMints, amount)
			}
		})
	}
}

func TestDeployerGroupMCMS(t *testing.T) {
	for _, tc := range deployerGroupTestCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "happy path" {
				tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/DX-405")
			}
			if tc.expectError {
				t.Skip("skipping test because it's not possible to verify error when using MCMS since we are explicitly failing the test in ApplyChangesets")
			}

			e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))

			tc.cfg.selector = e.HomeChainSel
			tc.cfg.MCMS = &proposalutils.TimelockConfig{
				MinDelay: 0,
			}
			state, err := stateview.LoadOnchainState(e.Env)
			require.NoError(t, err)

			contractsByChain := make(map[uint64][]common.Address)
			contractsByChain[e.HomeChainSel] = []common.Address{state.MustGetEVMChainState(e.HomeChainSel).LinkToken.Address()}

			_, err = commonchangeset.Apply(t, e.Env,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
					commonchangeset.TransferToMCMSWithTimelockConfig{
						ContractsByChain: contractsByChain,
						MCMSConfig: proposalutils.TimelockConfig{
							MinDelay: 0,
						},
					},
				),
			)
			require.NoError(t, err)

			_, err = commonchangeset.Apply(t, e.Env,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(dummyDeployerGroupGrantMintChangeset),
					tc.cfg,
				),
			)
			require.NoError(t, err)

			_, err = commonchangeset.Apply(t, e.Env,
				commonchangeset.Configure(
					cldf.CreateLegacyChangeSet(dummyDeployerGroupMintChangeset),
					tc.cfg,
				),
			)
			require.NoError(t, err)

			state, err = stateview.LoadOnchainState(e.Env)
			require.NoError(t, err)

			token := state.MustGetEVMChainState(e.HomeChainSel).LinkToken

			amount, err := token.BalanceOf(nil, tc.cfg.address)
			require.NoError(t, err)

			sumOfMints := big.NewInt(0)
			for _, mint := range tc.cfg.mints {
				sumOfMints = sumOfMints.Add(sumOfMints, mint)
			}

			require.Equal(t, sumOfMints, amount)
		})
	}
}

func TestDeployerGroupGenerateMultipleProposals(t *testing.T) {
	tc := dummyMultiChainDeployerGroupChangesetConfig{
		address: common.HexToAddress("0x455E5AA18469bC6ccEF49594645666C587A3a71B"),
		mints: []mintConfig{
			{
				selectorIndex: 0,
				amount:        big.NewInt(1),
			},
			{
				selectorIndex: 0,
				amount:        big.NewInt(2),
			},
			{
				selectorIndex: 1,
				amount:        big.NewInt(4),
			},
		},
		MCMS: &proposalutils.TimelockConfig{
			MinDelay: 0,
		},
	}
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	contractsByChain := make(map[uint64][]common.Address)
	for _, chain := range e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)) {
		contractsByChain[chain] = []common.Address{state.MustGetEVMChainState(chain).LinkToken.Address()}
	}

	_, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: contractsByChain,
				MCMSConfig: proposalutils.TimelockConfig{
					MinDelay: 0,
				},
			},
		),
	)
	require.NoError(t, err)

	_, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(dummyDeployerGroupGrantMintMultiChainChangeset),
			tc,
		),
	)
	require.NoError(t, err)

	cs, err := dummyDeployerGroupMintMultiDeploymentContextChangeset(e.Env, tc)
	require.NoError(t, err)
	require.Len(t, cs.MCMSTimelockProposals, len(tc.mints))
	require.Equal(t, "mint tokens 1", cs.MCMSTimelockProposals[0].Description)
	require.Equal(t, "mint tokens 2", cs.MCMSTimelockProposals[1].Description)
	require.Equal(t, "mint tokens 3", cs.MCMSTimelockProposals[2].Description)
	require.Equal(t, uint64(2), cs.MCMSTimelockProposals[0].ChainMetadata[mcmstypes.ChainSelector(e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[tc.mints[0].selectorIndex])].StartingOpCount)
	require.Equal(t, uint64(3), cs.MCMSTimelockProposals[1].ChainMetadata[mcmstypes.ChainSelector(e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[tc.mints[1].selectorIndex])].StartingOpCount)
	require.Equal(t, uint64(2), cs.MCMSTimelockProposals[2].ChainMetadata[mcmstypes.ChainSelector(e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[tc.mints[2].selectorIndex])].StartingOpCount)
	require.Len(t, cs.DescribedTimelockProposals, len(tc.mints))
	require.NotEmpty(t, cs.DescribedTimelockProposals[0])
	require.NotEmpty(t, cs.DescribedTimelockProposals[1])
	require.NotEmpty(t, cs.DescribedTimelockProposals[2])
}

func TestDeployerGroupMultipleProposalsMCMS(t *testing.T) {
	cfg := dummyMultiChainDeployerGroupChangesetConfig{
		address: common.HexToAddress("0x455E5AA18469bC6ccEF49594645666C587A3a71B"),
		mints: []mintConfig{
			{
				selectorIndex: 0,
				amount:        big.NewInt(1),
			},
			{
				selectorIndex: 0,
				amount:        big.NewInt(2),
			},
		},
		MCMS: &proposalutils.TimelockConfig{
			MinDelay: 0,
		},
	}

	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))

	currentState, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	contractsByChain := make(map[uint64][]common.Address)
	for _, chain := range e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM)) {
		contractsByChain[chain] = []common.Address{currentState.MustGetEVMChainState(chain).LinkToken.Address()}
	}

	_, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(commonchangeset.TransferToMCMSWithTimelockV2),
			commonchangeset.TransferToMCMSWithTimelockConfig{
				ContractsByChain: contractsByChain,
				MCMSConfig: proposalutils.TimelockConfig{
					MinDelay: 0,
				},
			},
		),
	)
	require.NoError(t, err)

	_, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(dummyDeployerGroupGrantMintMultiChainChangeset),
			cfg,
		),
	)
	require.NoError(t, err)

	_, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(dummyDeployerGroupMintMultiDeploymentContextChangeset),
			cfg,
		),
	)
	require.NoError(t, err)

	currentState, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	token := currentState.MustGetEVMChainState(e.HomeChainSel).LinkToken

	amount, err := token.BalanceOf(nil, cfg.address)
	require.NoError(t, err)

	sumOfMints := big.NewInt(0)
	for _, mint := range cfg.mints {
		sumOfMints = sumOfMints.Add(sumOfMints, mint.amount)
	}

	require.Equal(t, sumOfMints, amount)
}

func TestEmptyBatch(t *testing.T) {
	e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))

	cfg := dummyEmptyBatchChangesetConfig{
		MCMS: &proposalutils.TimelockConfig{
			MinDelay: 0,
		},
	}

	result, err := dummyEmptyBatchChangeset(e.Env, cfg)
	require.NoError(t, err)
	require.Empty(t, result.MCMSTimelockProposals)
}
