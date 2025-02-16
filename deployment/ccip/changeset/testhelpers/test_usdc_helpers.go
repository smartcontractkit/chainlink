package testhelpers

import (
	"math/big"
	"testing"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/latest/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_1/usdc_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
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
	lggr logger.Logger,
	chain deployment.Chain,
	state changeset.CCIPChainState,
	dstChain uint64,
	usdcToken *burn_mint_erc677.BurnMintERC677,
) error {
	config := []fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{
		{
			DestChainSelector: dstChain,
			TokenTransferFeeConfigs: []fee_quoter.FeeQuoterTokenTransferFeeConfigSingleTokenArgs{
				{
					Token: usdcToken.Address(),
					TokenTransferFeeConfig: fee_quoter.FeeQuoterTokenTransferFeeConfig{
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

func DeployUSDCPrerequisites(
	t *testing.T,
	logger logger.Logger,
	chain deployment.Chain,
	addressBook deployment.AddressBook,
) (*deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677], *deployment.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger]) {
	usdcToken, err := deployment.DeployContract(logger, chain, addressBook,
		func(chain deployment.Chain) deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
				chain.DeployerKey,
				chain.Client,
				"USDC",
				"USDC",
				6,
				big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
			)
			return deployment.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
				Address:  tokenAddress,
				Contract: token,
				Tv:       deployment.NewTypeAndVersion(changeset.USDCTokenPool, deployment.Version1_5_1),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	transmitter, err := deployment.DeployContract(logger, chain, addressBook,
		func(chain deployment.Chain) deployment.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter] {
			transmitterAddress, tx, transmitter, err := mock_usdc_token_transmitter.DeployMockE2EUSDCTransmitter(chain.DeployerKey, chain.Client, 0, 1, usdcToken.Address)
			return deployment.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter]{
				Address:  transmitterAddress,
				Contract: transmitter,
				Tv:       deployment.NewTypeAndVersion(changeset.USDCMockTransmitter, deployment.Version1_0_0),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	messenger, err := deployment.DeployContract(logger, chain, addressBook,
		func(chain deployment.Chain) deployment.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger] {
			messengerAddress, tx, messenger, err := mock_usdc_token_messenger.DeployMockE2EUSDCTokenMessenger(chain.DeployerKey, chain.Client, 0, transmitter.Address)
			return deployment.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger]{
				Address:  messengerAddress,
				Contract: messenger,
				Tv:       deployment.NewTypeAndVersion(changeset.USDCTokenMessenger, deployment.Version1_0_0),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	return usdcToken, messenger
}
