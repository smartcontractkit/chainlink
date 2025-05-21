package contracts

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/curve25519"

	ocrconfighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"

	commit_store_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/commit_store"
	evm_2_evm_offramp_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/evm_2_evm_offramp"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/token_admin_registry"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/commit_store"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/token_pool"
	evm_2_evm_onramp_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/evm_2_evm_onramp"
	price_registry_1_2_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/price_registry"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	lock_release_token_pool_1_4_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_4_0/lock_release_token_pool"
	token_pool_1_4_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_4_0/token_pool"
	usdc_token_pool_1_4_0 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_4_0/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/mock_lbtc_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/type_and_version"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/weth9"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/wrappers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers/testhelpers_1_4_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

// MatchContractVersionsOrAbove checks if the current contract versions for the test match or exceed the provided contract versions
func MatchContractVersionsOrAbove(requiredContractVersions map[Name]Version) error {
	for contractName, r := range requiredContractVersions {
		required := r
		if contractVersion, ok := VersionMap[contractName]; !ok {
			return fmt.Errorf("contract %s not found in version map", contractName)
		} else if contractVersion.Compare(&required.Version) < 0 {
			return fmt.Errorf("contract %s version %s is less than required version %s", contractName, contractVersion, required.Version)
		}
	}
	return nil
}

// NeedTokenAdminRegistry checks if token admin registry is needed for the current version of ccip
// if the version is less than 1.5.0, then token admin registry is not needed
func NeedTokenAdminRegistry() bool {
	return MatchContractVersionsOrAbove(map[Name]Version{
		TokenPoolContract: V1_5_0,
	}) == nil
}

// CCIPContractsDeployer provides the implementations for deploying CCIP ETH contracts
type CCIPContractsDeployer struct {
	evmClient blockchain.EVMClient
	logger    *zerolog.Logger
}

// NewCCIPContractsDeployer returns an instance of a contract deployer for CCIP
func NewCCIPContractsDeployer(logger *zerolog.Logger, bcClient blockchain.EVMClient) (*CCIPContractsDeployer, error) {
	return &CCIPContractsDeployer{
		evmClient: bcClient,
		logger:    logger,
	}, nil
}

func (e *CCIPContractsDeployer) Client() blockchain.EVMClient {
	return e.evmClient
}

func (e *CCIPContractsDeployer) DeployMultiCallContract() (common.Address, error) {
	multiCallABI, err := abi.JSON(strings.NewReader(MultiCallABI))
	if err != nil {
		return common.Address{}, err
	}
	address, tx, _, err := e.evmClient.DeployContract("MultiCall Contract", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		address, tx, contract, err := bind.DeployContract(auth, multiCallABI, common.FromHex(MultiCallBIN), wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return common.Address{}, nil, nil, err
		}
		return address, tx, contract, err
	})
	if err != nil {
		return common.Address{}, err
	}
	r, err := bind.WaitMined(context.Background(), e.evmClient.DeployBackend(), tx)
	if err != nil {
		return common.Address{}, err
	}
	if r.Status != types.ReceiptStatusSuccessful {
		return common.Address{}, fmt.Errorf("deploy multicall failed")
	}
	return *address, nil
}

func (e *CCIPContractsDeployer) DeployTokenMessenger(tokenTransmitter common.Address) (*common.Address, error) {
	address, _, _, err := e.evmClient.DeployContract("Mock Token Messenger", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		address, tx, contract, err := mock_usdc_token_messenger.DeployMockE2EUSDCTokenMessenger(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil), 0, tokenTransmitter)
		if err != nil {
			return common.Address{}, nil, nil, err
		}
		return address, tx, contract, err
	})

	return address, err
}

func (e *CCIPContractsDeployer) NewTokenTransmitter(addr common.Address) (*TokenTransmitter, error) {
	transmitter, err := mock_usdc_token_transmitter.NewMockE2EUSDCTransmitter(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))

	if err != nil {
		return nil, err
	}
	e.logger.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "Mock USDC Token Transmitter").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &TokenTransmitter{
		client:          e.evmClient,
		instance:        transmitter,
		ContractAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) DeployTokenTransmitter(domain uint32, usdcToken common.Address) (*TokenTransmitter, error) {
	address, _, instance, err := e.evmClient.DeployContract("Mock Token Transmitter", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		address, tx, contract, err := mock_usdc_token_transmitter.DeployMockE2EUSDCTransmitter(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil), 0, domain, usdcToken)
		if err != nil {
			return common.Address{}, nil, nil, err
		}
		return address, tx, contract, err
	})

	if err != nil {
		return nil, fmt.Errorf("error in deploying usdc token transmitter: %w", err)
	}

	return &TokenTransmitter{
		client:          e.evmClient,
		instance:        instance.(*mock_usdc_token_transmitter.MockE2EUSDCTransmitter),
		ContractAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) DeployLinkTokenContract() (*LinkToken, error) {
	address, _, instance, err := e.evmClient.DeployContract("Link Token", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return link_token_interface.DeployLinkToken(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
	})

	if err != nil {
		return nil, err
	}
	return &LinkToken{
		client:     e.evmClient,
		logger:     e.logger,
		instance:   instance.(*link_token_interface.LinkToken),
		EthAddress: *address,
	}, err
}

// DeployBurnMintERC677 deploys a BurnMintERC677 contract, mints given amount ( if provided) to the owner address and returns the ERC20Token wrapper instance
func (e *CCIPContractsDeployer) DeployBurnMintERC677(ownerMintingAmount *big.Int) (*ERC677Token, error) {
	address, _, instance, err := e.evmClient.DeployContract("Burn Mint ERC 677", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return burn_mint_erc677.DeployBurnMintERC677(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil), "Test Token ERC677", "TERC677", 6, new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e9)))
	})
	if err != nil {
		return nil, err
	}

	token := &ERC677Token{
		client:          e.evmClient,
		logger:          e.logger,
		ContractAddress: *address,
		instance:        instance.(*burn_mint_erc677.BurnMintERC677),
		OwnerAddress:    common.HexToAddress(e.evmClient.GetDefaultWallet().Address()),
		OwnerWallet:     e.evmClient.GetDefaultWallet(),
	}
	if ownerMintingAmount != nil {
		// grant minter role to owner and mint tokens
		err = token.GrantMintRole(common.HexToAddress(e.evmClient.GetDefaultWallet().Address()))
		if err != nil {
			return token, fmt.Errorf("granting minter role to owner shouldn't fail %w", err)
		}
		err = e.evmClient.WaitForEvents()
		if err != nil {
			return token, fmt.Errorf("error in waiting for granting mint role %w", err)
		}
		err = token.Mint(common.HexToAddress(e.evmClient.GetDefaultWallet().Address()), ownerMintingAmount)
		if err != nil {
			return token, fmt.Errorf("minting tokens shouldn't fail %w", err)
		}
	}
	return token, err
}

func (e *CCIPContractsDeployer) DeployCustomBurnMintERC677Token(name, symbol string, decimals uint8, ownerMintingAmount *big.Int) (*ERC677Token, error) {
	address, _, instance, err := e.evmClient.DeployContract("Burn Mint ERC 677", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return burn_mint_erc677.DeployBurnMintERC677(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil), name, symbol, decimals, new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e9)))
	})
	if err != nil {
		return nil, err
	}

	token := &ERC677Token{
		client:          e.evmClient,
		logger:          e.logger,
		ContractAddress: *address,
		instance:        instance.(*burn_mint_erc677.BurnMintERC677),
		OwnerAddress:    common.HexToAddress(e.evmClient.GetDefaultWallet().Address()),
		OwnerWallet:     e.evmClient.GetDefaultWallet(),
	}
	if ownerMintingAmount != nil {
		// grant minter role to owner and mint tokens
		err = token.GrantMintRole(common.HexToAddress(e.evmClient.GetDefaultWallet().Address()))
		if err != nil {
			return token, fmt.Errorf("granting minter role to owner shouldn't fail %w", err)
		}
		err = e.evmClient.WaitForEvents()
		if err != nil {
			return token, fmt.Errorf("error in waiting for granting mint role %w", err)
		}
		err = token.Mint(common.HexToAddress(e.evmClient.GetDefaultWallet().Address()), ownerMintingAmount)
		if err != nil {
			return token, fmt.Errorf("minting tokens shouldn't fail %w", err)
		}
	}
	return token, err
}

func (e *CCIPContractsDeployer) DeployERC20TokenContract(deployerFn blockchain.ContractDeployer) (*ERC20Token, error) {
	address, _, _, err := e.evmClient.DeployContract("Custom ERC20 Token", deployerFn)
	if err != nil {
		return nil, err
	}
	err = e.evmClient.WaitForEvents()
	if err != nil {
		return nil, err
	}
	return e.NewERC20TokenContract(*address)
}

func (e *CCIPContractsDeployer) NewLinkTokenContract(addr common.Address) (*LinkToken, error) {
	token, err := link_token_interface.NewLinkToken(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))

	if err != nil {
		return nil, err
	}
	e.logger.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "Link Token").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &LinkToken{
		client:     e.evmClient,
		logger:     e.logger,
		instance:   token,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) NewERC20TokenContract(addr common.Address) (*ERC20Token, error) {
	token, err := erc20.NewERC20(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))

	if err != nil {
		return nil, err
	}
	e.logger.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "ERC20 Token").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &ERC20Token{
		client:          e.evmClient,
		logger:          e.logger,
		instance:        token,
		ContractAddress: addr,
		OwnerAddress:    common.HexToAddress(e.evmClient.GetDefaultWallet().Address()),
		OwnerWallet:     e.evmClient.GetDefaultWallet(),
	}, err
}

func (e *CCIPContractsDeployer) NewLockReleaseTokenPoolContract(addr common.Address) (
	*TokenPool,
	error,
) {
	version := VersionMap[TokenPoolContract]
	e.logger.Info().Str("Version", version.String()).Msg("New LockRelease Token Pool")
	switch version {
	case Latest:
		pool, err := lock_release_token_pool.NewLockReleaseTokenPool(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))

		if err != nil {
			return nil, err
		}
		e.logger.Info().
			Str("Contract Address", addr.Hex()).
			Str("Contract Name", "Native Token Pool").
			Str("From", e.evmClient.GetDefaultWallet().Address()).
			Str("Network Name", e.evmClient.GetNetworkConfig().Name).
			Msg("New contract")
		poolInstance, err := token_pool.NewTokenPool(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, err
		}
		return &TokenPool{
			client: e.evmClient,
			logger: e.logger,
			Instance: &TokenPoolWrapper{
				Latest: &LatestPool{
					PoolInterface:   poolInstance,
					LockReleasePool: pool,
				},
			},
			EthAddress:   addr,
			OwnerAddress: common.HexToAddress(e.evmClient.GetDefaultWallet().Address()),
			OwnerWallet:  e.evmClient.GetDefaultWallet(),
		}, err
	case V1_4_0:
		pool, err := lock_release_token_pool_1_4_0.NewLockReleaseTokenPool(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, err
		}
		e.logger.Info().
			Str("Contract Address", addr.Hex()).
			Str("Contract Name", "Native Token Pool").
			Str("From", e.evmClient.GetDefaultWallet().Address()).
			Str("Network Name", e.evmClient.GetNetworkConfig().Name).
			Msg("New contract")
		poolInstance, err := token_pool_1_4_0.NewTokenPool(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, err
		}
		return &TokenPool{
			client: e.evmClient,
			logger: e.logger,
			Instance: &TokenPoolWrapper{
				V1_4_0: &V1_4_0Pool{
					PoolInterface:   poolInstance,
					LockReleasePool: pool,
				},
			},
			EthAddress:   addr,
			OwnerAddress: common.HexToAddress(e.evmClient.GetDefaultWallet().Address()),
			OwnerWallet:  e.evmClient.GetDefaultWallet(),
		}, err
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
}

func (e *CCIPContractsDeployer) NewUSDCTokenPoolContract(addr common.Address) (
	*TokenPool,
	error,
) {
	version := VersionMap[TokenPoolContract]
	e.logger.Info().Str("Version", version.String()).Msg("New USDC Token Pool")
	switch version {
	case Latest:
		pool, err := usdc_token_pool.NewUSDCTokenPool(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))

		if err != nil {
			return nil, err
		}
		e.logger.Info().
			Str("Contract Address", addr.Hex()).
			Str("Contract Name", "USDC Token Pool").
			Str("From", e.evmClient.GetDefaultWallet().Address()).
			Str("Network Name", e.evmClient.GetNetworkConfig().Name).
			Msg("New contract")
		poolInterface, err := token_pool.NewTokenPool(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, err
		}
		return &TokenPool{
			client: e.evmClient,
			logger: e.logger,
			Instance: &TokenPoolWrapper{
				Latest: &LatestPool{
					PoolInterface: poolInterface,
					USDCPool:      pool,
				},
			},
			EthAddress:   addr,
			OwnerAddress: common.HexToAddress(e.evmClient.GetDefaultWallet().Address()),
			OwnerWallet:  e.evmClient.GetDefaultWallet(),
		}, err
	case V1_4_0:
		pool, err := usdc_token_pool_1_4_0.NewUSDCTokenPool(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))

		if err != nil {
			return nil, err
		}
		e.logger.Info().
			Str("Contract Address", addr.Hex()).
			Str("Contract Name", "USDC Token Pool").
			Str("From", e.evmClient.GetDefaultWallet().Address()).
			Str("Network Name", e.evmClient.GetNetworkConfig().Name).
			Msg("New contract")
		poolInterface, err := token_pool_1_4_0.NewTokenPool(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, err
		}
		return &TokenPool{
			client: e.evmClient,
			logger: e.logger,
			Instance: &TokenPoolWrapper{
				V1_4_0: &V1_4_0Pool{
					PoolInterface: poolInterface,
					USDCPool:      pool,
				},
			},
			EthAddress:   addr,
			OwnerAddress: common.HexToAddress(e.evmClient.GetDefaultWallet().Address()),
			OwnerWallet:  e.evmClient.GetDefaultWallet(),
		}, err
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}

}

func (e *CCIPContractsDeployer) DeployUSDCTokenPoolContract(tokenAddr string, tokenMessenger, rmnProxy common.Address, router common.Address) (
	*TokenPool,
	error,
) {
	version := VersionMap[TokenPoolContract]
	e.logger.Debug().Str("Token", tokenAddr).Msg("Deploying USDC token pool")
	token := common.HexToAddress(tokenAddr)
	switch version {
	case Latest:
		address, _, _, err := e.evmClient.DeployContract("USDC Token Pool", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return usdc_token_pool.DeployUSDCTokenPool(
				auth,
				wrappers.MustNewWrappedContractBackend(e.evmClient, nil),
				tokenMessenger,
				token,
				[]common.Address{},
				rmnProxy,
				router,
			)
		})

		if err != nil {
			return nil, err
		}
		return e.NewUSDCTokenPoolContract(*address)
	case V1_4_0:
		address, _, _, err := e.evmClient.DeployContract("USDC Token Pool", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return usdc_token_pool_1_4_0.DeployUSDCTokenPool(
				auth,
				wrappers.MustNewWrappedContractBackend(e.evmClient, nil),
				tokenMessenger,
				token,
				[]common.Address{},
				rmnProxy,
				router,
			)
		})

		if err != nil {
			return nil, err
		}
		return e.NewUSDCTokenPoolContract(*address)
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
}

func (e *CCIPContractsDeployer) NewMockLBTCTokenPoolContract(addr common.Address) (
	*TokenPool,
	error,
) {
	version := VersionMap[TokenPoolContract]
	e.logger.Info().Str("Version", version.String()).Msg("New Mock LBTC Token Pool")
	switch version {
	case Latest:
		pool, err := mock_lbtc_token_pool.NewMockLBTCTokenPool(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))

		if err != nil {
			return nil, err
		}
		e.logger.Info().
			Str("Contract Address", addr.Hex()).
			Str("Contract Name", "Mock LBTC Token Pool").
			Str("From", e.evmClient.GetDefaultWallet().Address()).
			Str("Network Name", e.evmClient.GetNetworkConfig().Name).
			Msg("New contract")
		poolInterface, err := token_pool.NewTokenPool(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, err
		}
		return &TokenPool{
			client: e.evmClient,
			logger: e.logger,
			Instance: &TokenPoolWrapper{
				Latest: &LatestPool{
					PoolInterface: poolInterface,
					MockLBTCPool:  pool,
				},
			},
			EthAddress:   addr,
			OwnerAddress: common.HexToAddress(e.evmClient.GetDefaultWallet().Address()),
			OwnerWallet:  e.evmClient.GetDefaultWallet(),
		}, err
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
}

func (e *CCIPContractsDeployer) DeployMockLBTCTokenPoolContract(tokenAddr string, rmnProxy common.Address, router common.Address, destPoolData []byte) (
	*TokenPool,
	error,
) {
	e.logger.Println("In DeployMockLBTCTokenPoolContract")
	version := VersionMap[TokenPoolContract]
	e.logger.Debug().Str("Token", tokenAddr).Msg("Deploying Mock LBTC token pool")
	token := common.HexToAddress(tokenAddr)
	switch version {
	case Latest:
		address, _, _, err := e.evmClient.DeployContract("Mock LBTC Token Pool", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return mock_lbtc_token_pool.DeployMockLBTCTokenPool(
				auth,
				wrappers.MustNewWrappedContractBackend(e.evmClient, nil),
				token,
				[]common.Address{},
				rmnProxy,
				router,
				destPoolData,
			)
		})

		if err != nil {
			return nil, err
		}
		return e.NewMockLBTCTokenPoolContract(*address)
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
}

func (e *CCIPContractsDeployer) DeployLockReleaseTokenPoolContract(tokenAddr string, rmnProxy common.Address, router common.Address) (
	*TokenPool,
	error,
) {
	version := VersionMap[TokenPoolContract]
	e.logger.Info().Str("Version", version.String()).Msg("Deploying LockRelease Token Pool")
	token := common.HexToAddress(tokenAddr)
	switch version {
	case Latest:
		address, _, _, err := e.evmClient.DeployContract("LockRelease Token Pool", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return lock_release_token_pool.DeployLockReleaseTokenPool(
				auth,
				wrappers.MustNewWrappedContractBackend(e.evmClient, nil),
				token,
				testhelpers.TokenDecimals,
				[]common.Address{},
				rmnProxy,
				true,
				router,
			)
		})

		if err != nil {
			return nil, err
		}
		return e.NewLockReleaseTokenPoolContract(*address)
	case V1_4_0:
		address, _, _, err := e.evmClient.DeployContract("LockRelease Token Pool", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return lock_release_token_pool_1_4_0.DeployLockReleaseTokenPool(
				auth,
				wrappers.MustNewWrappedContractBackend(e.evmClient, nil),
				token,
				[]common.Address{},
				rmnProxy,
				true,
				router,
			)
		})

		if err != nil {
			return nil, err
		}
		return e.NewLockReleaseTokenPoolContract(*address)
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
}

func (e *CCIPContractsDeployer) DeployMockRMNContract() (*common.Address, error) {
	address, _, _, err := e.evmClient.DeployContract("Mock ARM Contract", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return mock_rmn_contract.DeployMockRMNContract(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
	})
	return address, err
}

func (e *CCIPContractsDeployer) NewRMNContract(addr common.Address) (*ARM, error) {
	arm, err := rmn_contract.NewRMNContract(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
	if err != nil {
		return nil, err
	}
	e.logger.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "Mock ARM Contract").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")

	return &ARM{
		client:     e.evmClient,
		Instance:   arm,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) NewCommitStore(addr common.Address) (
	*CommitStore,
	error,
) {
	version := VersionMap[CommitStoreContract]
	e.logger.Info().Str("Version", version.String()).Msg("New CommitStore")
	switch version {
	case Latest:
		ins, err := commit_store.NewCommitStore(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		e.logger.Info().
			Str("Contract Address", addr.Hex()).
			Str("Contract Name", "CommitStore").
			Str("From", e.evmClient.GetDefaultWallet().Address()).
			Str("Network Name", e.evmClient.GetNetworkConfig().Name).
			Msg("New contract")
		return &CommitStore{
			client: e.evmClient,
			logger: e.logger,
			Instance: &CommitStoreWrapper{
				Latest: ins,
			},
			EthAddress: addr,
		}, err
	case V1_2_0:
		ins, err := commit_store_1_2_0.NewCommitStore(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		e.logger.Info().
			Str("Contract Address", addr.Hex()).
			Str("Contract Name", "CommitStore").
			Str("From", e.evmClient.GetDefaultWallet().Address()).
			Str("Network Name", e.evmClient.GetNetworkConfig().Name).
			Msg("New contract")
		return &CommitStore{
			client: e.evmClient,
			logger: e.logger,
			Instance: &CommitStoreWrapper{
				V1_2_0: ins,
			},
			EthAddress: addr,
		}, err
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
}

func (e *CCIPContractsDeployer) DeployCommitStore(sourceChainSelector, destChainSelector uint64, onRamp common.Address, armProxy common.Address) (*CommitStore, error) {
	version, ok := VersionMap[CommitStoreContract]
	if !ok {
		return nil, fmt.Errorf("versioning not supported: %s", version)
	}
	e.logger.Info().Str("Version", version.String()).Msg("Deploying CommitStore")
	switch version {
	case Latest:
		address, _, instance, err := e.evmClient.DeployContract("CommitStore Contract", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return commit_store.DeployCommitStore(
				auth,
				wrappers.MustNewWrappedContractBackend(e.evmClient, nil),
				commit_store.CommitStoreStaticConfig{
					ChainSelector:       destChainSelector,
					SourceChainSelector: sourceChainSelector,
					OnRamp:              onRamp,
					RmnProxy:            armProxy,
				},
			)
		})
		if err != nil {
			return nil, err
		}
		return &CommitStore{
			client: e.evmClient,
			logger: e.logger,
			Instance: &CommitStoreWrapper{
				Latest: instance.(*commit_store.CommitStore),
			},
			EthAddress: *address,
		}, err
	case V1_2_0:
		address, _, instance, err := e.evmClient.DeployContract("CommitStore Contract", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return commit_store_1_2_0.DeployCommitStore(
				auth,
				wrappers.MustNewWrappedContractBackend(e.evmClient, nil),
				commit_store_1_2_0.CommitStoreStaticConfig{
					ChainSelector:       destChainSelector,
					SourceChainSelector: sourceChainSelector,
					OnRamp:              onRamp,
					ArmProxy:            armProxy,
				},
			)
		})
		if err != nil {
			return nil, err
		}
		return &CommitStore{
			client: e.evmClient,
			logger: e.logger,
			Instance: &CommitStoreWrapper{
				V1_2_0: instance.(*commit_store_1_2_0.CommitStore),
			},
			EthAddress: *address,
		}, err
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
}

func (e *CCIPContractsDeployer) DeployReceiverDapp(revert bool) (
	*ReceiverDapp,
	error,
) {
	address, _, instance, err := e.evmClient.DeployContract("ReceiverDapp", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return maybe_revert_message_receiver.DeployMaybeRevertMessageReceiver(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil), revert)
	})
	if err != nil {
		return nil, err
	}
	return &ReceiverDapp{
		client:     e.evmClient,
		logger:     e.logger,
		Instance:   instance.(*maybe_revert_message_receiver.MaybeRevertMessageReceiver),
		EthAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) NewReceiverDapp(addr common.Address) (
	*ReceiverDapp,
	error,
) {
	ins, err := maybe_revert_message_receiver.NewMaybeRevertMessageReceiver(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
	e.logger.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "ReceiverDapp").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &ReceiverDapp{
		client:     e.evmClient,
		logger:     e.logger,
		Instance:   ins,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) DeployRouter(wrappedNative common.Address, armAddress common.Address) (
	*Router,
	error,
) {
	address, _, instance, err := e.evmClient.DeployContract("Router", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return router.DeployRouter(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil), wrappedNative, armAddress)
	})
	if err != nil {
		return nil, err
	}
	return &Router{
		client:     e.evmClient,
		logger:     e.logger,
		Instance:   instance.(*router.Router),
		EthAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) NewRouter(addr common.Address) (
	*Router,
	error,
) {
	r, err := router.NewRouter(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
	e.logger.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "Router").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	if err != nil {
		return nil, err
	}
	return &Router{
		client:     e.evmClient,
		logger:     e.logger,
		Instance:   r,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) NewPriceRegistry(addr common.Address) (
	*PriceRegistry,
	error,
) {
	var wrapper *PriceRegistryWrapper
	version := VersionMap[PriceRegistryContract]
	e.logger.Info().Str("Version", version.String()).Msg("New PriceRegistry")
	switch version {
	case Latest:
		ins, err := price_registry_1_2_0.NewPriceRegistry(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, fmt.Errorf("error in creating price registry instance: %w", err)
		}
		wrapper = &PriceRegistryWrapper{
			V1_2_0: ins,
		}
	case V1_2_0:
		ins, err := price_registry_1_2_0.NewPriceRegistry(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, fmt.Errorf("error in creating price registry instance: %w", err)
		}
		wrapper = &PriceRegistryWrapper{
			V1_2_0: ins,
		}
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
	e.logger.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "PriceRegistry").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &PriceRegistry{
		client:     e.evmClient,
		logger:     e.logger,
		Instance:   wrapper,
		EthAddress: addr,
	}, nil
}

func (e *CCIPContractsDeployer) DeployPriceRegistry(tokens []common.Address) (*PriceRegistry, error) {
	var address *common.Address
	var wrapper *PriceRegistryWrapper
	var err error
	var instance interface{}
	version := VersionMap[PriceRegistryContract]
	e.logger.Info().Str("Version", version.String()).Msg("Deploying PriceRegistry")
	switch version {
	case Latest:
		address, _, instance, err = e.evmClient.DeployContract("PriceRegistry", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return price_registry_1_2_0.DeployPriceRegistry(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil), nil, tokens, 60*60*24*14)
		})
		if err != nil {
			return nil, err
		}
		wrapper = &PriceRegistryWrapper{
			V1_2_0: instance.(*price_registry_1_2_0.PriceRegistry),
		}
	case V1_2_0:
		address, _, instance, err = e.evmClient.DeployContract("PriceRegistry", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return price_registry_1_2_0.DeployPriceRegistry(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil), nil, tokens, 60*60*24*14)
		})
		if err != nil {
			return nil, err
		}
		wrapper = &PriceRegistryWrapper{
			V1_2_0: instance.(*price_registry_1_2_0.PriceRegistry),
		}
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
	reg := &PriceRegistry{
		client:     e.evmClient,
		logger:     e.logger,
		EthAddress: *address,
		Instance:   wrapper,
	}
	return reg, err
}

func (e *CCIPContractsDeployer) DeployTokenAdminRegistry() (*TokenAdminRegistry, error) {
	address, _, instance, err := e.evmClient.DeployContract("TokenAdminRegistry", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return token_admin_registry.DeployTokenAdminRegistry(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
	})
	if err != nil {
		return nil, err
	}
	return &TokenAdminRegistry{
		client:     e.evmClient,
		logger:     e.logger,
		Instance:   instance.(*token_admin_registry.TokenAdminRegistry),
		EthAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) NewTokenAdminRegistry(addr common.Address) (
	*TokenAdminRegistry,
	error,
) {
	ins, err := token_admin_registry.NewTokenAdminRegistry(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
	e.logger.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "TokenAdminRegistry").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &TokenAdminRegistry{
		client:     e.evmClient,
		logger:     e.logger,
		Instance:   ins,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) NewOnRamp(addr common.Address) (
	*OnRamp,
	error,
) {
	version := VersionMap[OnRampContract]
	e.logger.Info().Str("Version", version.String()).Msg("New OnRamp")
	e.logger.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "OnRamp").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	switch version {
	case V1_2_0:
		ins, err := evm_2_evm_onramp_1_2_0.NewEVM2EVMOnRamp(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, err
		}
		return &OnRamp{
			client:     e.evmClient,
			logger:     e.logger,
			Instance:   &OnRampWrapper{V1_2_0: ins},
			EthAddress: addr,
		}, err
	case Latest:
		ins, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, err
		}
		return &OnRamp{
			client:     e.evmClient,
			logger:     e.logger,
			Instance:   &OnRampWrapper{Latest: ins},
			EthAddress: addr,
		}, nil
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
}

func (e *CCIPContractsDeployer) DeployOnRamp(
	sourceChainSelector, destChainSelector uint64,
	tokensAndPools []evm_2_evm_onramp_1_2_0.InternalPoolUpdate,
	rmn,
	router,
	priceRegistry,
	tokenAdminRegistry common.Address,
	opts RateLimiterConfig,
	feeTokenConfig []evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs,
	tokenTransferFeeConfig []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs,
	linkTokenAddress common.Address,
) (*OnRamp, error) {
	version := VersionMap[OnRampContract]
	e.logger.Info().Str("Version", version.String()).Msg("Deploying OnRamp")
	switch version {
	case V1_2_0:
		feeTokenConfigV1_2_0 := make([]evm_2_evm_onramp_1_2_0.EVM2EVMOnRampFeeTokenConfigArgs, len(feeTokenConfig))
		for i, f := range feeTokenConfig {
			feeTokenConfigV1_2_0[i] = evm_2_evm_onramp_1_2_0.EVM2EVMOnRampFeeTokenConfigArgs{
				Token:                      f.Token,
				NetworkFeeUSDCents:         f.NetworkFeeUSDCents,
				GasMultiplierWeiPerEth:     f.GasMultiplierWeiPerEth,
				PremiumMultiplierWeiPerEth: f.PremiumMultiplierWeiPerEth,
				Enabled:                    f.Enabled,
			}
		}
		tokenTransferFeeConfigV1_2_0 := make([]evm_2_evm_onramp_1_2_0.EVM2EVMOnRampTokenTransferFeeConfigArgs, len(tokenTransferFeeConfig))
		for i, f := range tokenTransferFeeConfig {
			tokenTransferFeeConfigV1_2_0[i] = evm_2_evm_onramp_1_2_0.EVM2EVMOnRampTokenTransferFeeConfigArgs{
				Token:             f.Token,
				MinFeeUSDCents:    f.MinFeeUSDCents,
				MaxFeeUSDCents:    f.MaxFeeUSDCents,
				DeciBps:           f.DeciBps,
				DestGasOverhead:   f.DestGasOverhead,
				DestBytesOverhead: f.DestBytesOverhead,
			}
		}
		address, _, instance, err := e.evmClient.DeployContract("OnRamp", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return evm_2_evm_onramp_1_2_0.DeployEVM2EVMOnRamp(
				auth,
				wrappers.MustNewWrappedContractBackend(e.evmClient, nil),
				evm_2_evm_onramp_1_2_0.EVM2EVMOnRampStaticConfig{
					LinkToken:         linkTokenAddress,
					ChainSelector:     sourceChainSelector, // source chain id
					DestChainSelector: destChainSelector,   // destinationChainSelector
					DefaultTxGasLimit: 200_000,
					MaxNopFeesJuels:   big.NewInt(0).Mul(big.NewInt(100_000_000), big.NewInt(1e18)),
					PrevOnRamp:        common.HexToAddress(""),
					ArmProxy:          rmn,
				},
				evm_2_evm_onramp_1_2_0.EVM2EVMOnRampDynamicConfig{
					Router:                            router,
					MaxNumberOfTokensPerMsg:           50,
					DestGasOverhead:                   350_000,
					DestGasPerPayloadByte:             16,
					DestDataAvailabilityOverheadGas:   33_596,
					DestGasPerDataAvailabilityByte:    16,
					DestDataAvailabilityMultiplierBps: 6840, // 0.684
					PriceRegistry:                     priceRegistry,
					MaxDataBytes:                      50000,
					MaxPerMsgGasLimit:                 4_000_000,
				},
				tokensAndPools,
				evm_2_evm_onramp_1_2_0.RateLimiterConfig{
					Capacity: opts.Capacity,
					Rate:     opts.Rate,
				},
				feeTokenConfigV1_2_0,
				tokenTransferFeeConfigV1_2_0,
				[]evm_2_evm_onramp_1_2_0.EVM2EVMOnRampNopAndWeight{},
			)
		})
		if err != nil {
			return nil, err
		}
		return &OnRamp{
			client: e.evmClient,
			logger: e.logger,
			Instance: &OnRampWrapper{
				V1_2_0: instance.(*evm_2_evm_onramp_1_2_0.EVM2EVMOnRamp),
			},
			EthAddress: *address,
		}, nil
	case Latest:
		address, _, instance, err := e.evmClient.DeployContract("OnRamp", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return evm_2_evm_onramp.DeployEVM2EVMOnRamp(
				auth,
				wrappers.MustNewWrappedContractBackend(e.evmClient, nil),
				evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
					LinkToken:          linkTokenAddress,
					ChainSelector:      sourceChainSelector, // source chain id
					DestChainSelector:  destChainSelector,   // destinationChainSelector
					DefaultTxGasLimit:  200_000,
					MaxNopFeesJuels:    big.NewInt(0).Mul(big.NewInt(100_000_000), big.NewInt(1e18)),
					PrevOnRamp:         common.HexToAddress(""),
					RmnProxy:           rmn,
					TokenAdminRegistry: tokenAdminRegistry,
				},
				evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
					Router:                            router,
					MaxNumberOfTokensPerMsg:           50,
					DestGasOverhead:                   350_000,
					DestGasPerPayloadByte:             16,
					DestDataAvailabilityOverheadGas:   33_596,
					DestGasPerDataAvailabilityByte:    16,
					DestDataAvailabilityMultiplierBps: 6840, // 0.684
					PriceRegistry:                     priceRegistry,
					MaxDataBytes:                      50000,
					MaxPerMsgGasLimit:                 4_000_000,
					DefaultTokenFeeUSDCents:           50,
					DefaultTokenDestGasOverhead:       125_000,
					EnforceOutOfOrder:                 false,
				},
				evm_2_evm_onramp.RateLimiterConfig{
					Capacity: opts.Capacity,
					Rate:     opts.Rate,
				},
				feeTokenConfig,
				tokenTransferFeeConfig,
				[]evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{},
			)
		})
		if err != nil {
			return nil, err
		}
		return &OnRamp{
			client: e.evmClient,
			logger: e.logger,
			Instance: &OnRampWrapper{
				Latest: instance.(*evm_2_evm_onramp.EVM2EVMOnRamp),
			},
			EthAddress: *address,
		}, err
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
}

func (e *CCIPContractsDeployer) NewOffRamp(addr common.Address) (
	*OffRamp,
	error,
) {
	version := VersionMap[OffRampContract]
	e.logger.Info().Str("Version", version.String()).Msg("New OffRamp")
	switch version {
	case V1_2_0:
		ins, err := evm_2_evm_offramp_1_2_0.NewEVM2EVMOffRamp(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, err
		}
		e.logger.Info().
			Str("Contract Address", addr.Hex()).
			Str("Contract Name", "OffRamp").
			Str("From", e.evmClient.GetDefaultWallet().Address()).
			Str("Network Name", e.evmClient.GetNetworkConfig().Name).
			Msg("New contract")
		return &OffRamp{
			client:     e.evmClient,
			logger:     e.logger,
			Instance:   &OffRampWrapper{V1_2_0: ins},
			EthAddress: addr,
		}, err
	case Latest:
		ins, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
		if err != nil {
			return nil, err
		}
		e.logger.Info().
			Str("Contract Address", addr.Hex()).
			Str("Contract Name", "OffRamp").
			Str("From", e.evmClient.GetDefaultWallet().Address()).
			Str("Network Name", e.evmClient.GetNetworkConfig().Name).
			Msg("New contract")
		return &OffRamp{
			client:     e.evmClient,
			logger:     e.logger,
			Instance:   &OffRampWrapper{Latest: ins},
			EthAddress: addr,
		}, err
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
}

func (e *CCIPContractsDeployer) DeployOffRamp(
	sourceChainSelector, destChainSelector uint64,
	commitStore, onRamp common.Address,
	opts RateLimiterConfig,
	sourceTokens, pools []common.Address,
	rmnProxy common.Address,
	tokenAdminRegistry common.Address,
) (*OffRamp, error) {
	version := VersionMap[OffRampContract]
	e.logger.Info().Str("Version", version.String()).Msg("Deploying OffRamp")
	switch version {
	case V1_2_0:
		address, _, instance, err := e.evmClient.DeployContract("OffRamp Contract", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return evm_2_evm_offramp_1_2_0.DeployEVM2EVMOffRamp(
				auth,
				wrappers.MustNewWrappedContractBackend(e.evmClient, nil),
				evm_2_evm_offramp_1_2_0.EVM2EVMOffRampStaticConfig{
					CommitStore:         commitStore,
					ChainSelector:       destChainSelector,
					SourceChainSelector: sourceChainSelector,
					OnRamp:              onRamp,
					PrevOffRamp:         common.Address{},
					ArmProxy:            rmnProxy,
				},
				sourceTokens,
				pools,
				evm_2_evm_offramp_1_2_0.RateLimiterConfig{
					IsEnabled: true,
					Capacity:  opts.Capacity,
					Rate:      opts.Rate,
				},
			)
		})
		if err != nil {
			return nil, err
		}
		return &OffRamp{
			client: e.evmClient,
			logger: e.logger,
			Instance: &OffRampWrapper{
				V1_2_0: instance.(*evm_2_evm_offramp_1_2_0.EVM2EVMOffRamp),
			},
			EthAddress: *address,
		}, err
	case Latest:
		staticConfig := evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
			CommitStore:         commitStore,
			ChainSelector:       destChainSelector,
			SourceChainSelector: sourceChainSelector,
			OnRamp:              onRamp,
			PrevOffRamp:         common.Address{},
			RmnProxy:            rmnProxy,
			TokenAdminRegistry:  tokenAdminRegistry,
		}
		address, _, instance, err := e.evmClient.DeployContract("OffRamp Contract", func(
			auth *bind.TransactOpts,
			_ bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return evm_2_evm_offramp.DeployEVM2EVMOffRamp(
				auth,
				wrappers.MustNewWrappedContractBackend(e.evmClient, nil),
				staticConfig,
				evm_2_evm_offramp.RateLimiterConfig{
					IsEnabled: true,
					Capacity:  opts.Capacity,
					Rate:      opts.Rate,
				},
			)
		})
		e.logger.Info().Msg(fmt.Sprintf("deploying offramp with static config: %+v", staticConfig))

		if err != nil {
			return nil, err
		}
		return &OffRamp{
			client: e.evmClient,
			logger: e.logger,
			Instance: &OffRampWrapper{
				Latest: instance.(*evm_2_evm_offramp.EVM2EVMOffRamp),
			},
			EthAddress: *address,
		}, err
	default:
		return nil, fmt.Errorf("version not supported: %s", version)
	}
}

func (e *CCIPContractsDeployer) DeployWrappedNative() (*common.Address, error) {
	address, _, _, err := e.evmClient.DeployContract("WrappedNative", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return weth9.DeployWETH9(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
	})
	if err != nil {
		return nil, err
	}
	return address, err
}

func (e *CCIPContractsDeployer) DeployMockAggregator(decimals uint8, initialAns *big.Int) (*MockAggregator, error) {
	address, _, instance, err := e.evmClient.DeployContract("MockAggregator", func(
		auth *bind.TransactOpts,
		_ bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return mock_v3_aggregator_contract.DeployMockV3Aggregator(auth, wrappers.MustNewWrappedContractBackend(e.evmClient, nil), decimals, initialAns)
	})
	if err != nil {
		return nil, fmt.Errorf("deploying mock aggregator: %w", err)
	}
	e.logger.Info().
		Str("Contract Address", address.Hex()).
		Str("Contract Name", "MockAggregator").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &MockAggregator{
		client:          e.evmClient,
		logger:          e.logger,
		Instance:        instance.(*mock_v3_aggregator_contract.MockV3Aggregator),
		ContractAddress: *address,
	}, nil
}

func (e *CCIPContractsDeployer) NewMockAggregator(addr common.Address) (*MockAggregator, error) {
	ins, err := mock_v3_aggregator_contract.NewMockV3Aggregator(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
	if err != nil {
		return nil, fmt.Errorf("creating mock aggregator: %w", err)
	}
	return &MockAggregator{
		client:          e.evmClient,
		logger:          e.logger,
		Instance:        ins,
		ContractAddress: addr,
	}, nil
}

func (e *CCIPContractsDeployer) TypeAndVersion(addr common.Address) (string, error) {
	tv, err := type_and_version.NewITypeAndVersion(addr, wrappers.MustNewWrappedContractBackend(e.evmClient, nil))
	if err != nil {
		return "", err
	}
	tvStr, err := tv.TypeAndVersion(nil)
	if err != nil {
		return "", fmt.Errorf("error calling typeAndVersion on addr: %s %w", addr.Hex(), err)
	}
	e.logger.Info().
		Str("TypeAndVersion", tvStr).
		Str("Contract Address", addr.Hex()).
		Msg("TypeAndVersion")

	_, versionStr, err := ccipconfig.ParseTypeAndVersion(tvStr)
	if err != nil {
		return versionStr, err
	}
	v, err := semver.NewVersion(versionStr)
	if err != nil {
		return "", fmt.Errorf("failed parsing version %s: %w", versionStr, err)
	}
	return v.String(), nil
}

// OCR2ParamsForCommit and OCR2ParamsForExec -
// These functions return the default OCR2 parameters for Commit and Exec respectively.
// Refer to CommitOCRParams and ExecOCRParams in CCIPTestConfig located in testconfig/ccip.go to override these values with custom param values.
func OCR2ParamsForCommit(blockTime time.Duration) contracts.OffChainAggregatorV2Config {
	// slow blocktime chains like Ethereum
	if blockTime >= 10*time.Second {
		return contracts.OffChainAggregatorV2Config{
			DeltaProgress:                           config.MustNewDuration(2 * time.Minute),
			DeltaResend:                             config.MustNewDuration(5 * time.Second),
			DeltaRound:                              config.MustNewDuration(90 * time.Second),
			DeltaGrace:                              config.MustNewDuration(5 * time.Second),
			DeltaStage:                              config.MustNewDuration(60 * time.Second),
			MaxDurationQuery:                        config.MustNewDuration(100 * time.Millisecond),
			MaxDurationObservation:                  config.MustNewDuration(35 * time.Second),
			MaxDurationReport:                       config.MustNewDuration(10 * time.Second),
			MaxDurationShouldAcceptFinalizedReport:  config.MustNewDuration(5 * time.Second),
			MaxDurationShouldTransmitAcceptedReport: config.MustNewDuration(10 * time.Second),
		}
	}
	// fast blocktime chains like Avalanche
	return contracts.OffChainAggregatorV2Config{
		DeltaProgress:                           config.MustNewDuration(2 * time.Minute),
		DeltaResend:                             config.MustNewDuration(5 * time.Second),
		DeltaRound:                              config.MustNewDuration(60 * time.Second),
		DeltaGrace:                              config.MustNewDuration(5 * time.Second),
		DeltaStage:                              config.MustNewDuration(25 * time.Second),
		MaxDurationQuery:                        config.MustNewDuration(100 * time.Millisecond),
		MaxDurationObservation:                  config.MustNewDuration(35 * time.Second),
		MaxDurationReport:                       config.MustNewDuration(10 * time.Second),
		MaxDurationShouldAcceptFinalizedReport:  config.MustNewDuration(5 * time.Second),
		MaxDurationShouldTransmitAcceptedReport: config.MustNewDuration(10 * time.Second),
	}
}

func OCR2ParamsForExec(blockTime time.Duration) contracts.OffChainAggregatorV2Config {
	// slow blocktime chains like Ethereum
	if blockTime >= 10*time.Second {
		return contracts.OffChainAggregatorV2Config{
			DeltaProgress:                           config.MustNewDuration(2 * time.Minute),
			DeltaResend:                             config.MustNewDuration(5 * time.Second),
			DeltaRound:                              config.MustNewDuration(90 * time.Second),
			DeltaGrace:                              config.MustNewDuration(5 * time.Second),
			DeltaStage:                              config.MustNewDuration(60 * time.Second),
			MaxDurationQuery:                        config.MustNewDuration(100 * time.Millisecond),
			MaxDurationObservation:                  config.MustNewDuration(35 * time.Second),
			MaxDurationReport:                       config.MustNewDuration(10 * time.Second),
			MaxDurationShouldAcceptFinalizedReport:  config.MustNewDuration(5 * time.Second),
			MaxDurationShouldTransmitAcceptedReport: config.MustNewDuration(10 * time.Second),
		}
	}
	// fast blocktime chains like Avalanche
	return contracts.OffChainAggregatorV2Config{
		DeltaProgress:                           config.MustNewDuration(120 * time.Second),
		DeltaResend:                             config.MustNewDuration(5 * time.Second),
		DeltaRound:                              config.MustNewDuration(30 * time.Second),
		DeltaGrace:                              config.MustNewDuration(5 * time.Second),
		DeltaStage:                              config.MustNewDuration(10 * time.Second),
		MaxDurationQuery:                        config.MustNewDuration(100 * time.Millisecond),
		MaxDurationObservation:                  config.MustNewDuration(35 * time.Second),
		MaxDurationReport:                       config.MustNewDuration(10 * time.Second),
		MaxDurationShouldAcceptFinalizedReport:  config.MustNewDuration(5 * time.Second),
		MaxDurationShouldTransmitAcceptedReport: config.MustNewDuration(10 * time.Second),
	}
}

func OffChainAggregatorV2ConfigWithNodes(numberNodes int, inflightExpiry time.Duration, cfg contracts.OffChainAggregatorV2Config) contracts.OffChainAggregatorV2Config {
	if numberNodes <= 4 {
		log.Err(fmt.Errorf("insufficient number of nodes (%d) supplied for OCR, need at least 5", numberNodes)).
			Int("Number Chainlink Nodes", numberNodes).
			Msg("You likely need more chainlink nodes to properly configure OCR, try 5 or more.")
	}
	s := make([]int, 0)
	for i := 0; i < numberNodes; i++ {
		s = append(s, 1)
	}
	faultyNodes := 0
	if numberNodes > 1 {
		faultyNodes = (numberNodes - 1) / 3
	}
	if faultyNodes == 0 {
		faultyNodes = 1
	}
	if cfg.DeltaStage == nil {
		cfg.DeltaStage = config.MustNewDuration(inflightExpiry)
	}
	return contracts.OffChainAggregatorV2Config{
		DeltaProgress:                           cfg.DeltaProgress,
		DeltaResend:                             cfg.DeltaResend,
		DeltaRound:                              cfg.DeltaRound,
		DeltaGrace:                              cfg.DeltaGrace,
		DeltaStage:                              cfg.DeltaStage,
		RMax:                                    3,
		S:                                       s,
		F:                                       faultyNodes,
		Oracles:                                 []ocrconfighelper2.OracleIdentityExtra{},
		MaxDurationQuery:                        cfg.MaxDurationQuery,
		MaxDurationObservation:                  cfg.MaxDurationObservation,
		MaxDurationReport:                       cfg.MaxDurationReport,
		MaxDurationShouldAcceptFinalizedReport:  cfg.MaxDurationShouldAcceptFinalizedReport,
		MaxDurationShouldTransmitAcceptedReport: cfg.MaxDurationShouldTransmitAcceptedReport,
		OnchainConfig:                           []byte{},
	}
}

func stripKeyPrefix(key string) string {
	chunks := strings.Split(key, "_")
	if len(chunks) == 3 {
		return chunks[2]
	}
	return key
}

func NewCommitOffchainConfig(
	GasPriceHeartBeat config.Duration,
	DAGasPriceDeviationPPB uint32,
	ExecGasPriceDeviationPPB uint32,
	TokenPriceHeartBeat config.Duration,
	TokenPriceDeviationPPB uint32,
	InflightCacheExpiry config.Duration,
	priceReportingDisabled bool) (ccipconfig.OffchainConfig, error) {
	switch VersionMap[CommitStoreContract] {
	case Latest:
		return testhelpers.NewCommitOffchainConfig(
			GasPriceHeartBeat,
			DAGasPriceDeviationPPB,
			ExecGasPriceDeviationPPB,
			TokenPriceHeartBeat,
			TokenPriceDeviationPPB,
			InflightCacheExpiry,
			priceReportingDisabled,
		), nil
	case V1_2_0:
		return testhelpers_1_4_0.NewCommitOffchainConfig(
			GasPriceHeartBeat,
			DAGasPriceDeviationPPB,
			ExecGasPriceDeviationPPB,
			TokenPriceHeartBeat,
			TokenPriceDeviationPPB,
			InflightCacheExpiry,
			priceReportingDisabled,
		), nil
	default:
		return nil, fmt.Errorf("version not supported: %s", VersionMap[CommitStoreContract])
	}
}

func NewCommitOnchainConfig(
	PriceRegistry common.Address,
) (abihelpers.AbiDefined, error) {
	switch VersionMap[CommitStoreContract] {
	case Latest:
		return testhelpers.NewCommitOnchainConfig(PriceRegistry), nil
	case V1_2_0:
		return testhelpers_1_4_0.NewCommitOnchainConfig(PriceRegistry), nil
	default:
		return nil, fmt.Errorf("version not supported: %s", VersionMap[CommitStoreContract])
	}
}

func NewExecOnchainConfig(
	PermissionLessExecutionThresholdSeconds uint32,
	Router common.Address,
	PriceRegistry common.Address,
	MaxNumberOfTokensPerMsg uint16,
	MaxDataBytes uint32,
	MaxPoolReleaseOrMintGas uint32,
) (abihelpers.AbiDefined, error) {
	switch VersionMap[OffRampContract] {
	case Latest:
		return testhelpers.NewExecOnchainConfig(PermissionLessExecutionThresholdSeconds, Router, PriceRegistry, MaxNumberOfTokensPerMsg, MaxDataBytes), nil
	case V1_2_0:
		return testhelpers_1_4_0.NewExecOnchainConfig(
			PermissionLessExecutionThresholdSeconds,
			Router,
			PriceRegistry,
			MaxNumberOfTokensPerMsg,
			MaxDataBytes,
			MaxPoolReleaseOrMintGas,
		), nil
	default:
		return nil, fmt.Errorf("version not supported: %s", VersionMap[OffRampContract])
	}
}

// NewExecOffchainConfig creates a config for the OffChain portion of how CCIP operates
func NewExecOffchainConfig(
	destOptimisticConfirmations uint32,
	batchGasLimit uint32,
	relativeBoostPerWaitHour float64,
	inflightCacheExpiry config.Duration,
	rootSnoozeTime config.Duration,
	batchingStrategyID uint32, // See ccipexec package
) (ccipconfig.OffchainConfig, error) {
	switch VersionMap[OffRampContract] {
	case Latest:
		return testhelpers.NewExecOffchainConfig(
			destOptimisticConfirmations,
			batchGasLimit,
			relativeBoostPerWaitHour,
			inflightCacheExpiry,
			rootSnoozeTime,
			batchingStrategyID,
		), nil
	case V1_2_0:
		return testhelpers_1_4_0.NewExecOffchainConfig(
			destOptimisticConfirmations,
			batchGasLimit,
			relativeBoostPerWaitHour,
			inflightCacheExpiry,
			rootSnoozeTime,
			batchingStrategyID,
		), nil
	default:
		return nil, fmt.Errorf("version not supported: %s", VersionMap[OffRampContract])
	}
}

func NewOffChainAggregatorV2ConfigForCCIPPlugin[T ccipconfig.OffchainConfig](
	nodes []*nodeclient.CLNodesWithKeys,
	offchainCfg T,
	onchainCfg abihelpers.AbiDefined,
	ocr2Params contracts.OffChainAggregatorV2Config,
	inflightExpiry time.Duration,
) (
	signers []common.Address,
	transmitters []common.Address,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	oracleIdentities := make([]ocrconfighelper2.OracleIdentityExtra, 0)
	ocrConfig := OffChainAggregatorV2ConfigWithNodes(len(nodes), inflightExpiry, ocr2Params)
	var onChainKeys []ocrtypes2.OnchainPublicKey
	for i, nodeWithKeys := range nodes {
		ocr2Key := nodeWithKeys.KeysBundle.OCR2Key.Data
		offChainPubKeyTemp, err := hex.DecodeString(stripKeyPrefix(ocr2Key.Attributes.OffChainPublicKey))
		if err != nil {
			return nil, nil, 0, nil, 0, nil, err
		}
		formattedOnChainPubKey := stripKeyPrefix(ocr2Key.Attributes.OnChainPublicKey)
		cfgPubKeyTemp, err := hex.DecodeString(stripKeyPrefix(ocr2Key.Attributes.ConfigPublicKey))
		if err != nil {
			return nil, nil, 0, nil, 0, nil, err
		}
		cfgPubKeyBytes := [ed25519.PublicKeySize]byte{}
		copy(cfgPubKeyBytes[:], cfgPubKeyTemp)
		offChainPubKey := [curve25519.PointSize]byte{}
		copy(offChainPubKey[:], offChainPubKeyTemp)
		ethAddress := nodeWithKeys.KeysBundle.EthAddress
		p2pKeys := nodeWithKeys.KeysBundle.P2PKeys
		peerID := p2pKeys.Data[0].Attributes.PeerID
		oracleIdentities = append(oracleIdentities, ocrconfighelper2.OracleIdentityExtra{
			OracleIdentity: ocrconfighelper2.OracleIdentity{
				OffchainPublicKey: offChainPubKey,
				OnchainPublicKey:  common.HexToAddress(formattedOnChainPubKey).Bytes(),
				PeerID:            peerID,
				TransmitAccount:   ocrtypes2.Account(ethAddress),
			},
			ConfigEncryptionPublicKey: cfgPubKeyBytes,
		})
		onChainKeys = append(onChainKeys, oracleIdentities[i].OnchainPublicKey)
		transmitters = append(transmitters, common.HexToAddress(ethAddress))
	}
	signers, err = evm.OnchainPublicKeyToAddress(onChainKeys)
	if err != nil {
		return nil, nil, 0, nil, 0, nil, err
	}
	ocrConfig.Oracles = oracleIdentities
	ocrConfig.ReportingPluginConfig, err = ccipconfig.EncodeOffchainConfig(offchainCfg)
	if err != nil {
		return nil, nil, 0, nil, 0, nil, err
	}
	ocrConfig.OnchainConfig, err = abihelpers.EncodeAbiStruct(onchainCfg)
	if err != nil {
		return nil, nil, 0, nil, 0, nil, err
	}

	_, _, f_, onchainConfig_, offchainConfigVersion, offchainConfig, err = ocrconfighelper2.ContractSetConfigArgsForTests(
		ocrConfig.DeltaProgress.Duration(),
		ocrConfig.DeltaResend.Duration(),
		ocrConfig.DeltaRound.Duration(),
		ocrConfig.DeltaGrace.Duration(),
		ocrConfig.DeltaStage.Duration(),
		ocrConfig.RMax,
		ocrConfig.S,
		ocrConfig.Oracles,
		ocrConfig.ReportingPluginConfig,
		nil,
		ocrConfig.MaxDurationQuery.Duration(),
		ocrConfig.MaxDurationObservation.Duration(),
		ocrConfig.MaxDurationReport.Duration(),
		ocrConfig.MaxDurationShouldAcceptFinalizedReport.Duration(),
		ocrConfig.MaxDurationShouldTransmitAcceptedReport.Duration(),
		ocrConfig.F,
		ocrConfig.OnchainConfig,
	)
	return
}
