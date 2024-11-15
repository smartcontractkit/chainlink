package ccipdeployment

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	sel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/burn_mint_erc677"
	"math/big"
)

func ConfigureUSDCTokenPools(lggr logger.Logger, chains map[uint64]deployment.Chain, src, dst uint64, state CCIPOnChainState) (*burn_mint_erc677.BurnMintERC677, *burn_mint_erc677.BurnMintERC677, error) {
	srcToken := state.Chains[src].BurnMintTokens677[USDCSymbol]
	dstToken := state.Chains[dst].BurnMintTokens677[USDCSymbol]
	srcPool := state.Chains[src].USDCTokenPool
	dstPool := state.Chains[dst].USDCTokenPool
	srcTransmitter := state.Chains[src].MockUSDCTokenMessenger
	dstTransmitter := state.Chains[dst].MockUSDCTokenMessenger

	// Attach token pools to registry
	if err := attachTokenToTheRegistry(chains[src], state.Chains[src], chains[src].DeployerKey, srcToken.Address(), srcPool.Address()); err != nil {
		return nil, nil, err
	}

	if err := attachTokenToTheRegistry(chains[dst], state.Chains[dst], chains[dst].DeployerKey, dstToken.Address(), dstPool.Address()); err != nil {
		return nil, nil, err
	}

	// Connect pool to each other
	if err := setUSDCTokenPoolCounterPart(chains[src], srcPool, dst, dstToken.Address(), dstPool.Address()); err != nil {
		return nil, nil, err
	}

	if err := setUSDCTokenPoolCounterPart(chains[dst], dstPool, src, srcToken.Address(), srcPool.Address()); err != nil {
		return nil, nil, err
	}

	// Add burn/mint permissions
	if err := grantMintBurnPermissions(lggr, chains[src], srcToken, srcPool.Address()); err != nil {
		return nil, nil, err
	}

	if err := grantMintBurnPermissions(lggr, chains[dst], dstToken, dstPool.Address()); err != nil {
		return nil, nil, err
	}

	if err := grantMintBurnPermissions(lggr, chains[src], srcToken, srcTransmitter.Address()); err != nil {
		return nil, nil, err
	}

	if err := grantMintBurnPermissions(lggr, chains[dst], dstToken, dstTransmitter.Address()); err != nil {
		return nil, nil, err
	}

	return srcToken, dstToken, nil
}

func DeployUSDCToken(
	lggr logger.Logger,
	chains map[uint64]deployment.Chain,
	src, dst uint64,
	state CCIPOnChainState,
	addresses deployment.AddressBook,
) (
	*burn_mint_erc677.BurnMintERC677,
	*usdc_token_pool.USDCTokenPool,
	*burn_mint_erc677.BurnMintERC677,
	*usdc_token_pool.USDCTokenPool,
	error,
) {
	srcToken, srcPool, srcTransmitter, _, err := deployUSDCTokenOneEnd(lggr, chains[src], addresses)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	dstToken, dstPool, dstTransmitter, _, err := deployUSDCTokenOneEnd(lggr, chains[dst], addresses)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Attach token pools to registry
	if err := attachTokenToTheRegistry(chains[src], state.Chains[src], chains[src].DeployerKey, srcToken.Address(), srcPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := attachTokenToTheRegistry(chains[dst], state.Chains[dst], chains[dst].DeployerKey, dstToken.Address(), dstPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	// Connect pool to each other
	if err := setUSDCTokenPoolCounterPart(chains[src], srcPool, dst, dstToken.Address(), dstPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := setUSDCTokenPoolCounterPart(chains[dst], dstPool, src, srcToken.Address(), srcPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	// Add burn/mint permissions
	if err := grantMintBurnPermissions(lggr, chains[src], srcToken, srcPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := grantMintBurnPermissions(lggr, chains[dst], dstToken, dstPool.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := grantMintBurnPermissions(lggr, chains[src], srcToken, srcTransmitter.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := grantMintBurnPermissions(lggr, chains[dst], dstToken, dstTransmitter.Address()); err != nil {
		return nil, nil, nil, nil, err
	}

	return srcToken, dstPool, dstToken, srcPool, nil
}

func UpdateFeeQuoterForUSDC(
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
		return err
	}

	_, err = chain.Confirm(tx)
	return err
}

func deployUSDCTokenOneEnd(
	lggr logger.Logger,
	chain deployment.Chain,
	addresses deployment.AddressBook,
) (
	*burn_mint_erc677.BurnMintERC677,
	*usdc_token_pool.USDCTokenPool,
	*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger,
	*mock_usdc_token_transmitter.MockE2EUSDCTransmitter,
	error,
) {
	token, err := deployContract(lggr, chain, addresses,
		func(chain deployment.Chain) ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			USDCTokenAddr, tx, token, err2 := burn_mint_erc677.DeployBurnMintERC677(
				chain.DeployerKey,
				chain.Client,
				"USDC Token",
				"USDC",
				uint8(18),
				big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
			)
			return ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
				USDCTokenAddr, token, tx, deployment.NewTypeAndVersion(USDCToken, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tx, err := token.Contract.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	chainId, _ := sel.ChainIdFromSelector(chain.Selector)

	transmitter, err := deployContract(lggr, chain, addresses,
		func(chain deployment.Chain) ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter] {
			transmitterAddress, tx, mockTransmitterContract, err2 := mock_usdc_token_transmitter.DeployMockE2EUSDCTransmitter(
				chain.DeployerKey,
				chain.Client,
				0,
				uint32(chainId),
				token.Address,
			)
			return ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter]{
				transmitterAddress, mockTransmitterContract, tx, deployment.NewTypeAndVersion(USDCMockTransmitter, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy mock USDC transmitter", "err", err)
		return nil, nil, nil, nil, err
	}

	lggr.Infow("deployed mock USDC transmitter", "addr", transmitter.Address)

	messenger, err := deployContract(lggr, chain, addresses,
		func(chain deployment.Chain) ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger] {
			tokenMessengerAddress, tx, tokenMessengerContract, err2 := mock_usdc_token_messenger.DeployMockE2EUSDCTokenMessenger(
				chain.DeployerKey,
				chain.Client,
				0,
				transmitter.Address,
			)
			return ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger]{
				tokenMessengerAddress, tokenMessengerContract, tx, deployment.NewTypeAndVersion(USDCTokenMessenger, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy USDC token messenger", "err", err)
		return nil, nil, nil, nil, err
	}
	lggr.Infow("deployed mock USDC token messenger", "addr", messenger.Address)
	chainAddr, err := addresses.AddressesForChain(chain.Selector)
	if err != nil {
		lggr.Errorw("Failed to get addresses of chain", "err", err)
		return nil, nil, nil, nil, err
	}

	var rmnAddress, routerAddress string
	for address, v := range chainAddr {
		if deployment.NewTypeAndVersion(ARMProxy, deployment.Version1_0_0) == v {
			rmnAddress = address
		}
		if deployment.NewTypeAndVersion(Router, deployment.Version1_2_0) == v {
			routerAddress = address
		}
		if rmnAddress != "" && routerAddress != "" {
			break
		}
	}

	tokenPool, err := deployContract(lggr, chain, addresses,
		func(chain deployment.Chain) ContractDeploy[*usdc_token_pool.USDCTokenPool] {
			tokenPoolAddress, tx, tokenPoolContract, err2 := usdc_token_pool.DeployUSDCTokenPool(
				chain.DeployerKey,
				chain.Client,
				messenger.Address,
				token.Address,
				[]common.Address{},
				common.HexToAddress(rmnAddress),
				common.HexToAddress(routerAddress),
			)
			return ContractDeploy[*usdc_token_pool.USDCTokenPool]{
				tokenPoolAddress, tokenPoolContract, tx, deployment.NewTypeAndVersion(USDCTokenPool, deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy USDC token pool", "err", err)
		return nil, nil, nil, nil, err
	}
	lggr.Infow("deployed USDC token pool", "addr", tokenPool.Address)

	domainIdentifier, err := tokenPool.Contract.ILocalDomainIdentifier(nil)
	fmt.Println("LocalDomainIdentifier", domainIdentifier)

	return token.Contract, tokenPool.Contract, messenger.Contract, transmitter.Contract, nil
}
