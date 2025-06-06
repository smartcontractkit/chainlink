package testhelpers

import (
	"testing"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func ConfigureUSDCTokenPools(
	lggr logger.Logger,
	chains map[uint64]cldf_evm.Chain,
	src, dst uint64,
	state stateview.CCIPOnChainState,
) (*burn_mint_erc677.BurnMintERC677, *burn_mint_erc677.BurnMintERC677, error) {
	srcToken := state.MustGetEVMChainState(src).BurnMintTokens677[shared.USDCSymbol]
	dstToken := state.MustGetEVMChainState(dst).BurnMintTokens677[shared.USDCSymbol]
	srcPool := state.MustGetEVMChainState(src).USDCTokenPools[deployment.Version1_5_1]
	dstPool := state.MustGetEVMChainState(dst).USDCTokenPools[deployment.Version1_5_1]

	args := []struct {
		sourceChain cldf_evm.Chain
		dstChainSel uint64
		state       evm.CCIPChainState
		srcToken    *burn_mint_erc677.BurnMintERC677
		srcPool     *usdc_token_pool.USDCTokenPool
		dstToken    *burn_mint_erc677.BurnMintERC677
		dstPool     *usdc_token_pool.USDCTokenPool
	}{
		{
			chains[src],
			dst,
			state.MustGetEVMChainState(src),
			srcToken,
			srcPool,
			dstToken,
			dstPool,
		},
		{
			chains[dst],
			src,
			state.MustGetEVMChainState(dst),
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
	sourceChain cldf_evm.Chain,
	dstChainSel uint64,
	state evm.CCIPChainState,
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
	e cldf.Environment,
	lggr logger.Logger,
	chain cldf_evm.Chain,
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
	_, err := commonchangeset.Apply(t, e,
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(v1_6.ApplyTokenTransferFeeConfigUpdatesFeeQuoterChangeset),
			v1_6.ApplyTokenTransferFeeConfigUpdatesConfig{
				UpdatesByChain: map[uint64]v1_6.ApplyTokenTransferFeeConfigUpdatesConfigPerChain{
					chain.Selector: {
						TokenTransferFeeConfigArgs: []v1_6.TokenTransferFeeConfigArg{
							{
								DestChain: dstChain,
								TokenTransferFeeConfigPerToken: map[shared.TokenSymbol]fee_quoter.FeeQuoterTokenTransferFeeConfig{
									shared.USDCSymbol: config,
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
