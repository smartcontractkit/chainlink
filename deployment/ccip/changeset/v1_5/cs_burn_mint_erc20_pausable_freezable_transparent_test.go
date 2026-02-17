package v1_5

import (
	"math/big"
	"strings"
	"testing"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/ccip-contract-examples/chains/evm/gobindings/generated/1_6_1/transparent_upgradeable_proxy"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/1_5_0/burn_mint_erc20_pausable_freezable_transparent"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

func TestBurnMintERC20PausableFreezableTransparentValidate(t *testing.T) {
	t.Parallel()

	var err error

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain := evmSelectors[0]

	const initABIString = `[{
		"name": "initialize",
		"type": "function",
		"inputs": [
			{"name": "name", "type": "string"},
			{"name": "symbol", "type": "string"},
			{"name": "decimals", "type": "uint8"},
			{"name": "maxSupply", "type": "uint256"},
			{"name": "preMint", "type": "uint256"},
			{"name": "defaultAdmin", "type": "address"}
		]
	}]`

	parsedABI, err := abi.JSON(strings.NewReader(initABIString))
	require.NoError(t, err)

	tests := []struct {
		name    string
		cfg     BurnMintERC20PausableFreezableTransparentChangesetConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: BurnMintERC20PausableFreezableTransparentChangesetConfig{
				Tokens: map[uint64][]string{
					chain: {"TEST1", "TEST2", "TEST3"},
				},
			},
		},
		{
			name: "token already exists in state",
			cfg: BurnMintERC20PausableFreezableTransparentChangesetConfig{
				Tokens: map[uint64][]string{
					chain: {"TEST1"},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		for chainSelector, tokens := range tc.cfg.Tokens {
			chain := e.Env.BlockChains.EVMChains()[chainSelector]

			err := tc.cfg.Validate(e.Env)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if !tc.wantErr {
				for _, token := range tokens {
					_, tx, implementation, err := burn_mint_erc20_pausable_freezable_transparent.DeployBurnMintERC20PausableFreezableTransparent(chain.DeployerKey, chain.Client)
					require.NoError(t, err)

					_, err = e.Env.BlockChains.EVMChains()[chainSelector].Confirm(tx)
					require.NoError(t, err)

					err = e.Env.ExistingAddresses.Save(chainSelector, implementation.Address().String(), cldf.NewTypeAndVersion(shared.BurnMintERC20PausableFreezableTransparentToken, deployment.Version1_5_0))
					require.NoError(t, err)

					initData, err := parsedABI.Pack(
						"initialize",
						token,
						token,
						uint8(18),
						big.NewInt(0),
						big.NewInt(0),
						e.Env.BlockChains.EVMChains()[chainSelector].DeployerKey.From,
					)
					require.NoError(t, err)

					_, tx, proxy, err := transparent_upgradeable_proxy.DeployTransparentUpgradeableProxy(
						chain.DeployerKey, chain.Client, implementation.Address(), chain.DeployerKey.From, initData,
					)
					require.NoError(t, err)

					_, err = e.Env.BlockChains.EVMChains()[chainSelector].Confirm(tx)
					require.NoError(t, err)

					err = e.Env.ExistingAddresses.Save(chainSelector, proxy.Address().String(), cldf.NewTypeAndVersion(shared.TransparentUpgradeableProxy, deployment.Version1_6_1))
					require.NoError(t, err)
				}
			}
		}
	}
}

func TestBurnMintERC20PausableFreezableTransparentDeploy(t *testing.T) {
	t.Parallel()

	var err error

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1, chain2 := evmSelectors[0], evmSelectors[1]

	e.Env, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(DeployBurnMintERC20PausableFreezableTransparent),
			BurnMintERC20PausableFreezableTransparentChangesetConfig{
				Tokens: map[uint64][]string{
					chain1: {"TEST1", "TEST2", "TEST3"},
					chain2: {"TEST1", "TEST2", "TEST3"},
				},
			},
		),
	)

	require.NoError(t, err)
}
