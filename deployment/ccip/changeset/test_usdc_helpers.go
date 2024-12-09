package changeset

import (
	"math/big"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
)

func ConfigureUSDCTokenPools(
	lggr logger.Logger,
	chains map[uint64]deployment.Chain,
	src, dst uint64,
	state CCIPOnChainState,
) (*burn_mint_erc677.BurnMintERC677, *burn_mint_erc677.BurnMintERC677, error) {
	srcToken := state.Chains[src].BurnMintTokens677[USDCSymbol]
	dstToken := state.Chains[dst].BurnMintTokens677[USDCSymbol]
	srcPool := state.Chains[src].USDCTokenPool
	dstPool := state.Chains[dst].USDCTokenPool

	args := []struct {
		sourceChain deployment.Chain
		dstChainSel uint64
		state       CCIPChainState
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
	state CCIPChainState,
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
	lggr logger.Logger,
	chain deployment.Chain,
	state CCIPChainState,
	dstChain uint64,
	usdcToken *burn_mint_erc677.BurnMintERC677,
) error {
	config := []fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{
		{
			DestChainSelector: dstChain,
			TokenTransferFeeConfigs: []fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs{
				{
					usdcToken.Address(),
					fee_quoter.FeeQuoterTokenTransferFeeConfig{
						MinFeeUSDCents:    50,
						MaxFeeUSDCents:    50_000,
						DeciBps:           0,
						DestGasOverhead:   180_000,
						DestBytesOverhead: 640,
						IsEnabled:         true,
					},
				},
			},
		},
	}

	tx, err := state.FeeQuoter.ApplyTokenTransferFeeConfigUpdates(
		chain.DeployerKey,
		config,
		[]fee_quoter.FeeQuoterTokenTransferFeeConfigRemoveArgs{},
	)
	if err != nil {
		lggr.Errorw("Failed to apply token transfer fee config updates", "err", err, "config", config)
		return err
	}

	_, err = chain.Confirm(tx)
	return err
}

func DeployUSDC(
	lggr logger.Logger,
	chain deployment.Chain,
	addresses deployment.AddressBook,
	rmnProxy common.Address,
	router common.Address,
) (
	*burn_mint_erc677.BurnMintERC677,
	*usdc_token_pool.USDCTokenPool,
	*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger,
	*mock_usdc_token_transmitter.MockE2EUSDCTransmitter,
	error,
) {
	token, err := deployment.DeployContract(lggr, chain, addresses,
		func(chain deployment.Chain) deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			tokenAddress, tx, tokenContract, err2 := burn_mint_erc677.DeployBurnMintERC677(
				chain.DeployerKey,
				chain.Client,
				USDCName,
				string(USDCSymbol),
				UsdcDecimals,
				big.NewInt(0),
			)
			return deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
				Address:  tokenAddress,
				Contract: tokenContract,
				Tx:       tx,
				Tv:       deployment.NewTypeAndVersion(USDCToken, deployment.Version1_0_0),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy USDC token", "chain", chain.String(), "err", err)
		return nil, nil, nil, nil, err
	}

	tx, err := token.Contract.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	if err != nil {
		lggr.Errorw("Failed to grant mint role", "chain", chain.String(), "token", token.Contract.Address(), "err", err)
		return nil, nil, nil, nil, err
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	transmitter, err := deployment.DeployContract(lggr, chain, addresses,
		func(chain deployment.Chain) deployment.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter] {
			transmitterAddress, tx, transmitterContract, err2 := mock_usdc_token_transmitter.DeployMockE2EUSDCTransmitter(
				chain.DeployerKey,
				chain.Client,
				0,
				reader.AllAvailableDomains()[chain.Selector],
				token.Address,
			)
			return deployment.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter]{
				Address:  transmitterAddress,
				Contract: transmitterContract,
				Tx:       tx,
				Tv:       deployment.NewTypeAndVersion(USDCMockTransmitter, deployment.Version1_0_0),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy mock USDC transmitter", "chain", chain.String(), "err", err)
		return nil, nil, nil, nil, err
	}

	messenger, err := deployment.DeployContract(lggr, chain, addresses,
		func(chain deployment.Chain) deployment.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger] {
			messengerAddress, tx, messengerContract, err2 := mock_usdc_token_messenger.DeployMockE2EUSDCTokenMessenger(
				chain.DeployerKey,
				chain.Client,
				0,
				transmitter.Address,
			)
			return deployment.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger]{
				Address:  messengerAddress,
				Contract: messengerContract,
				Tx:       tx,
				Tv:       deployment.NewTypeAndVersion(USDCTokenMessenger, deployment.Version1_0_0),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy USDC token messenger", "chain", chain.String(), "err", err)
		return nil, nil, nil, nil, err
	}

	tokenPool, err := deployment.DeployContract(lggr, chain, addresses,
		func(chain deployment.Chain) deployment.ContractDeploy[*usdc_token_pool.USDCTokenPool] {
			tokenPoolAddress, tx, tokenPoolContract, err2 := usdc_token_pool.DeployUSDCTokenPool(
				chain.DeployerKey,
				chain.Client,
				messenger.Address,
				token.Address,
				[]common.Address{},
				rmnProxy,
				router,
			)
			return deployment.ContractDeploy[*usdc_token_pool.USDCTokenPool]{
				Address:  tokenPoolAddress,
				Contract: tokenPoolContract,
				Tx:       tx,
				Tv:       deployment.NewTypeAndVersion(USDCTokenPool, deployment.Version1_0_0),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy USDC token pool", "chain", chain.String(), "err", err)
		return nil, nil, nil, nil, err
	}

	return token.Contract, tokenPool.Contract, messenger.Contract, transmitter.Contract, nil
}
