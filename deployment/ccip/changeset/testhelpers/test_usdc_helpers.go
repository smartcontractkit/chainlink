package testhelpers

import (
	"testing"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func ConfigureUSDCTokenPools(
	lggr logger.Logger,
	chains map[uint64]deployment.Chain,
	src, dst uint64,
	state changeset.CCIPOnChainState,
) (*burn_mint_erc677.BurnMintERC677, *burn_mint_erc677.BurnMintERC677, error) {
	srcToken := state.Chains[src].BurnMintTokens677[changeset.USDCSymbol]
	dstToken := state.Chains[dst].BurnMintTokens677[changeset.USDCSymbol]
	srcPool := state.Chains[src].USDCTokenPools[deployment.Version1_5_1]
	dstPool := state.Chains[dst].USDCTokenPools[deployment.Version1_5_1]

	args := []struct {
		sourceChain deployment.Chain
		dstChainSel uint64
		state       changeset.CCIPChainState
		srcToken    *burn_mint_erc677.BurnMintERC677
		srcPool     *usdc_token_pool.USDCTokenPool
		dstToken    *burn_mint_erc677.BurnMintERC677
		dstPool     *usdc_token_pool.USDCTokenPool
	}{
		{
			chains[src],
			dst,
			state.Chains[src],
			srcToken,
			srcPool,
			dstToken,
			dstPool,
		},
		{
			chains[dst],
			src,
			state.Chains[dst],
			dstToken,
			dstPool,
			srcToken,
			srcPool,
		},
	}

	configurePoolGrp := errgroup.Group{}
	for _, arg := range args {
		configurePoolGrp.Go(configureSingleChain(lggr, arg.sourceChain, arg.dstChainSel, arg.state, arg.srcToken, arg.srcPool, arg.dstToken, arg.dstPool))
	}
	if err := configurePoolGrp.Wait(); err != nil {
		return nil, nil, err
	}
	return srcToken, dstToken, nil
}

func configureSingleChain(
	lggr logger.Logger,
	sourceChain deployment.Chain,
	dstChainSel uint64,
	state changeset.CCIPChainState,
	srcToken *burn_mint_erc677.BurnMintERC677,
	srcPool *usdc_token_pool.USDCTokenPool,
	dstToken *burn_mint_erc677.BurnMintERC677,
	dstPool *usdc_token_pool.USDCTokenPool,
) func() error {
	return func() error {
		if err := attachTokenToTheRegistry(sourceChain, state, sourceChain.DeployerKey, srcToken.Address(), srcPool.Address()); err != nil {
			lggr.Errorw("Failed to attach token to the registry", "err", err, "token", srcToken.Address(), "pool", srcPool.Address())
			return err
		}

		if err := setUSDCTokenPoolCounterPart(sourceChain, srcPool, dstChainSel, sourceChain.DeployerKey, dstToken.Address(), dstPool.Address()); err != nil {
			lggr.Errorw("Failed to set counter part", "err", err, "srcPool", srcPool.Address(), "dstPool", dstPool.Address())
			return err
		}

		for _, addr := range []common.Address{
			srcPool.Address(),
			state.MockUSDCTokenMessenger.Address(),
			state.MockUSDCTransmitter.Address(),
		} {
			if err := grantMintBurnPermissions(lggr, sourceChain, srcToken, sourceChain.DeployerKey, addr); err != nil {
				lggr.Errorw("Failed to grant mint/burn permissions", "err", err, "token", srcToken.Address(), "address", addr)
				return err
			}
		}
		return nil
	}
}

func UpdateFeeQuoterForUSDC(
	t *testing.T,
	e deployment.Environment,
	lggr logger.Logger,
	chain deployment.Chain,
	dstChain uint64,
) error {
	config := fee_quoter.FeeQuoterTokenTransferFeeConfig{
		MinFeeUSDCents:    50,
		MaxFeeUSDCents:    50_000,
		DeciBps:           0,
		DestGasOverhead:   180_000,
		DestBytesOverhead: 640,
		IsEnabled:         true,
	}
	_, err := commonchangeset.Apply(t, e, nil,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.ApplyTokenTransferFeeConfigUpdatesFeeQuoterChangeset),
			v1_6.ApplyTokenTransferFeeConfigUpdatesConfig{
				UpdatesByChain: map[uint64]v1_6.ApplyTokenTransferFeeConfigUpdatesConfigPerChain{
					chain.Selector: {
						TokenTransferFeeConfigArgs: []v1_6.TokenTransferFeeConfigArg{
							{
								DestChain: dstChain,
								TokenTransferFeeConfigPerToken: map[changeset.TokenSymbol]fee_quoter.FeeQuoterTokenTransferFeeConfig{
									changeset.USDCSymbol: config,
								},
							},
						},
					},
				},
			}),
	)

	if err != nil {
		lggr.Errorw("Failed to apply token transfer fee config updates", "err", err, "config", config)
		return err
	}
	return nil
}
