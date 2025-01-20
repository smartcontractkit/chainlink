package changeset_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

type dummyDeployerGroupChangesetConfig struct {
	selector uint64
	address  common.Address
	mints    []*big.Int
	MCMS     *changeset.MCMSConfig
}

func dummyDeployerGroupGrantMintChangeset(e deployment.Environment, cfg dummyDeployerGroupChangesetConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	token := state.Chains[cfg.selector].LinkToken

	group := changeset.NewDeployerGroup(e, state, cfg.MCMS)
	deployer, err := group.GetDeployer(cfg.selector)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	_, err = token.GrantMintRole(deployer, deployer.From)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return group.Enact("Grant mint role")
}

func dummyDeployerGroupMintChangeset(e deployment.Environment, cfg dummyDeployerGroupChangesetConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	token := state.Chains[cfg.selector].LinkToken

	group := changeset.NewDeployerGroup(e, state, cfg.MCMS)
	deployer, err := group.GetDeployer(cfg.selector)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	for _, mint := range cfg.mints {
		_, err = token.Mint(deployer, cfg.address, mint)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
	}

	return group.Enact("Mint tokens")
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

				state, err := changeset.LoadOnchainState(e.Env)
				require.NoError(t, err)

				token := state.Chains[e.HomeChainSel].LinkToken

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
			if tc.expectError {
				t.Skip("skipping test because it's not possible to verify error when using MCMS since we are explicitly failing the test in ApplyChangesets")
			}

			e, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(2))

			tc.cfg.selector = e.HomeChainSel
			tc.cfg.MCMS = &changeset.MCMSConfig{
				MinDelay: 0,
			}
			state, err := changeset.LoadOnchainState(e.Env)
			require.NoError(t, err)

			timelocksPerChain := changeset.BuildTimelockPerChain(e.Env, state)

			contractsByChain := make(map[uint64][]common.Address)
			contractsByChain[e.HomeChainSel] = []common.Address{state.Chains[e.HomeChainSel].LinkToken.Address()}

			_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(commonchangeset.TransferToMCMSWithTimelock),
					Config: commonchangeset.TransferToMCMSWithTimelockConfig{
						ContractsByChain: contractsByChain,
						MinDelay:         0,
					},
				},
			})
			require.NoError(t, err)

			_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(dummyDeployerGroupGrantMintChangeset),
					Config:    tc.cfg,
				},
			})
			require.NoError(t, err)

			_, err = commonchangeset.ApplyChangesets(t, e.Env, timelocksPerChain, []commonchangeset.ChangesetApplication{
				{
					Changeset: commonchangeset.WrapChangeSet(dummyDeployerGroupMintChangeset),
					Config:    tc.cfg,
				},
			})
			require.NoError(t, err)

			state, err = changeset.LoadOnchainState(e.Env)
			require.NoError(t, err)

			token := state.Chains[e.HomeChainSel].LinkToken

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
