package v1_5

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestBurnMintERC20TransparentValidate(t *testing.T) {
	t.Parallel()

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain := evmSelectors[0]

	tests := []struct {
		name    string
		cfg     BurnMintERC20TransparentChangesetConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: BurnMintERC20TransparentChangesetConfig{
				Tokens: map[uint64]map[string]BurnMintERC20Transparent{
					chain: {
						"TEST_1": {
							MaxSupply: big.NewInt(1000),
							PreMint:   big.NewInt(100),
						},
					},
				},
			},
		},
		{
			name: "missing max supply",
			cfg: BurnMintERC20TransparentChangesetConfig{
				Tokens: map[uint64]map[string]BurnMintERC20Transparent{
					chain: {
						"TEST_2": {
							PreMint: big.NewInt(100),
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.cfg.Validate(e.Env)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBurnMintERC20TransparentDeploy(t *testing.T) {
	t.Parallel()

	var err error

	e, _ := testhelpers.NewMemoryEnvironment(t)
	evmSelectors := e.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain := evmSelectors[0]

	e.Env, err = commonchangeset.Apply(t, e.Env,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(DeployBurnMintERC20Transparent),
			BurnMintERC20TransparentChangesetConfig{
				Tokens: map[uint64]map[string]BurnMintERC20Transparent{
					chain: {
						"TEST": {
							MaxSupply: big.NewInt(0),
							PreMint:   big.NewInt(0),
						},
					},
				},
			},
		),
	)

	require.NoError(t, err)
}
