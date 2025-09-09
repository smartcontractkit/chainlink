package testhelpers

import (
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/burn_mint_token_pool"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
)

func ConfigureLBTCTokenPools(
	lggr logger.Logger,
	chains map[uint64]cldf_evm.Chain,
	src, dst uint64,
	state stateview.CCIPOnChainState,
) (srcToken *burn_mint_erc677.BurnMintERC677, dstToken *burn_mint_erc677.BurnMintERC677, err error) {
	srcToken = state.MustGetEVMChainState(src).BurnMintTokens677[shared.LBTCSymbol]
	dstToken = state.MustGetEVMChainState(dst).BurnMintTokens677[shared.LBTCSymbol]
	srcPool := state.MustGetEVMChainState(src).BurnMintTokenPools[shared.LBTCSymbol][deployment.Version1_5_1]
	dstPool := state.MustGetEVMChainState(dst).BurnMintTokenPools[shared.LBTCSymbol][deployment.Version1_5_1]

	args := []struct {
		sourceChain cldf_evm.Chain
		dstChainSel uint64
		state       evm.CCIPChainState
		srcToken    *burn_mint_erc677.BurnMintERC677
		srcPool     *burn_mint_token_pool.BurnMintTokenPool
		dstToken    *burn_mint_erc677.BurnMintERC677
		dstPool     *burn_mint_token_pool.BurnMintTokenPool
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
		configurePoolGrp.Go(configureSingleChainForLBTC(lggr, arg.sourceChain, arg.dstChainSel, arg.state, arg.srcToken, arg.srcPool, arg.dstToken, arg.dstPool))
	}
	if err = configurePoolGrp.Wait(); err != nil {
		return nil, nil, err
	}
	return srcToken, dstToken, nil
}

func configureSingleChainForLBTC(
	lggr logger.Logger,
	sourceChain cldf_evm.Chain,
	dstChainSel uint64,
	state evm.CCIPChainState,
	srcToken *burn_mint_erc677.BurnMintERC677,
	srcPool *burn_mint_token_pool.BurnMintTokenPool,
	dstToken *burn_mint_erc677.BurnMintERC677,
	dstPool *burn_mint_token_pool.BurnMintTokenPool,
) func() error {
	return func() error {
		if err := attachTokenToTheRegistry(sourceChain, state, sourceChain.DeployerKey, srcToken.Address(), srcPool.Address()); err != nil {
			lggr.Errorw("Failed to attach token to the registry", "err", err, "token", srcToken.Address(), "pool", srcPool.Address())
			return err
		}
		if err := setTokenPoolCounterPart(sourceChain, srcPool, sourceChain.DeployerKey, dstChainSel, dstToken.Address().Bytes(), dstPool.Address().Bytes()); err != nil {
			lggr.Errorw("Failed to set counter part", "err", err, "srcPool", srcPool.Address(), "dstPool", dstPool.Address())
			return err
		}
		if err := grantMintBurnPermissions(lggr, sourceChain, srcToken, sourceChain.DeployerKey, srcPool.Address()); err != nil {
			lggr.Errorw("Failed to grant mint/burn permissions", "err", err, "token", srcToken.Address(), "address", srcPool.Address())
			return err
		}

		return nil
	}
}
