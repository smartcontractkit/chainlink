package ccipdeployment

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"math/big"
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

	// Attach token pools to registry
	if err := attachTokenToTheRegistry(chains[src], state.Chains[src], chains[src].DeployerKey, srcToken.Address(), srcPool.Address()); err != nil {
		lggr.Errorw("Failed to attach token to the registry", "err", err, "token", srcToken.Address(), "pool", srcPool.Address())
		return nil, nil, err
	}

	if err := attachTokenToTheRegistry(chains[dst], state.Chains[dst], chains[dst].DeployerKey, dstToken.Address(), dstPool.Address()); err != nil {
		lggr.Errorw("Failed to attach token to the registry", "err", err, "token", dstToken.Address(), "pool", dstPool.Address())
		return nil, nil, err
	}

	// Connect pool to each other
	if err := setUSDCTokenPoolCounterPart(chains[src], srcPool, dst, dstToken.Address(), dstPool.Address()); err != nil {
		lggr.Errorw("Failed to set counter part", "err", err, "srcPool", srcPool.Address(), "dstPool", dstPool.Address())
		return nil, nil, err
	}

	if err := setUSDCTokenPoolCounterPart(chains[dst], dstPool, src, srcToken.Address(), srcPool.Address()); err != nil {
		lggr.Errorw("Failed to set counter part", "err", err, "srcPool", dstPool.Address(), "dstPool", srcPool.Address())
		return nil, nil, err
	}

	// Add burn/mint permissions for source
	for _, addr := range []common.Address{
		srcPool.Address(),
		state.Chains[src].MockUSDCTokenMessenger.Address(),
		state.Chains[src].MockUSDCTransmitter.Address(),
	} {
		if err := grantMintBurnPermissions(lggr, chains[src], srcToken, addr); err != nil {
			lggr.Errorw("Failed to grant mint/burn permissions", "err", err, "token", srcToken.Address(), "minter", addr)
			return nil, nil, err
		}
	}

	// Add burn/mint permissions for dest
	for _, addr := range []common.Address{
		dstPool.Address(),
		state.Chains[dst].MockUSDCTokenMessenger.Address(),
		state.Chains[dst].MockUSDCTransmitter.Address(),
	} {
		if err := grantMintBurnPermissions(lggr, chains[dst], dstToken, addr); err != nil {
			lggr.Errorw("Failed to grant mint/burn permissions", "err", err, "token", dstToken.Address(), "minter", addr)
			return nil, nil, err
		}
	}

	return srcToken, dstToken, nil
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
	state CCIPChainState,
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
				"USDC Token",
				"USDC",
				uint8(18),
				big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
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
		lggr.Errorw("Failed to deploy USDC token", "err", err)
		return nil, nil, nil, nil, err
	}

	tx, err := token.Contract.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	if err != nil {
		lggr.Errorw("Failed to grant mint role", "token", token.Contract.Address(), "err", err)
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
		lggr.Errorw("Failed to deploy mock USDC transmitter", "err", err)
		return nil, nil, nil, nil, err
	}

	lggr.Infow("deployed mock USDC transmitter", "addr", transmitter.Address)

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
		lggr.Errorw("Failed to deploy USDC token messenger", "err", err)
		return nil, nil, nil, nil, err
	}
	lggr.Infow("deployed mock USDC token messenger", "addr", messenger.Address)

	tokenPool, err := deployment.DeployContract(lggr, chain, addresses,
		func(chain deployment.Chain) deployment.ContractDeploy[*usdc_token_pool.USDCTokenPool] {
			tokenPoolAddress, tx, tokenPoolContract, err2 := usdc_token_pool.DeployUSDCTokenPool(
				chain.DeployerKey,
				chain.Client,
				messenger.Address,
				token.Address,
				[]common.Address{},
				state.RMNProxyExisting.Address(),
				state.Router.Address(),
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
		lggr.Errorw("Failed to deploy USDC token pool", "err", err)
		return nil, nil, nil, nil, err
	}
	lggr.Infow("deployed USDC token pool", "addr", tokenPool.Address)

	return token.Contract, tokenPool.Contract, messenger.Contract, transmitter.Contract, nil
}
