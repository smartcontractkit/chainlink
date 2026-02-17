package v1_6_1

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/ccip-contract-examples/chains/evm/gobindings/generated/1_6_1/transparent_upgradeable_proxy"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/1_5_0/burn_mint_erc20_transparent"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	_ cldf.ChangeSet[TransparentUpgradeableProxyChangesetConfig]          = DeployTransparentUpgradeableProxy
	_ cldf.ChangeSet[TransparentUpgradeableProxyChangesetConfig]          = SaveProxyAdmin
	_ cldf.ChangeSet[TransparentUpgradeableProxyChangesetConfig]          = BeginDefaultAdminTransferTransparentUpgradeableProxy
	_ cldf.ChangeSet[TransparentUpgradeableProxyChangesetConfig]          = AcceptDefaultAdminTransferTransparentUpgradeableProxy
	_ cldf.ChangeSet[TransparentUpgradeableProxyChangesetConfig]          = SetCCIPAdminTransferTransparentUpgradeableProxy
	_ cldf.ChangeSet[TransparentUpgradeableProxyGrantRoleChangesetConfig] = GrantRoleTransparentUpgradeableProxy
)

type TransparentUpgradeableProxy struct {
	BurnMintERC20Transparent common.Address
	Symbol                   string
	Decimals                 uint8
	MaxSupply                *big.Int
	PreMint                  *big.Int
	NewAdmin                 common.Address
	Initialize               bool
}

type TransparentUpgradeableProxyChangesetConfig struct {
	Tokens map[uint64]map[string]TransparentUpgradeableProxy
	MCMS   *proposalutils.TimelockConfig
}

type TransparentUpgradeableProxyRole string

const (
	TransparentUpgradeableProxyBurnerRole TransparentUpgradeableProxyRole = "BURNER_ROLE"
	TransparentUpgradeableProxyMinterRole TransparentUpgradeableProxyRole = "MINTER_ROLE"
)

type TransparentUpgradeableProxyGrantRole struct {
	Symbol  string
	Role    TransparentUpgradeableProxyRole
	Account common.Address
}

type TransparentUpgradeableProxyGrantRoleChangesetConfig struct {
	Tokens map[uint64]map[string][]TransparentUpgradeableProxyGrantRole
	MCMS   *proposalutils.TimelockConfig
}

func (c TransparentUpgradeableProxyChangesetConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokens := range c.Tokens {
		for token, config := range tokens {
			chain, ok := e.BlockChains.EVMChains()[chainSelector]
			if !ok {
				return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
			}

			if config.Initialize {
				if config.BurnMintERC20Transparent == (common.Address{}) {
					return fmt.Errorf("BurnMintERC20Transparent address is required for %s token on chain %d", token, chainSelector)
				}

				implementation, err := burn_mint_erc20_transparent.NewBurnMintERC20Transparent(config.BurnMintERC20Transparent, chain.Client)
				if err != nil {
					return fmt.Errorf("failed to instantiate BurnMintERC20Transparent at %s for %s token on %s: %w", config.BurnMintERC20Transparent, token, chain, err)
				}

				symbol, err := implementation.Symbol(&bind.CallOpts{Context: e.GetContext()})
				if err != nil {
					return fmt.Errorf("failed to get symbol from BurnMintERC20Transparent at %s for %s token on %s: %w", config.BurnMintERC20Transparent, token, chain, err)
				}

				if symbol != "" {
					return fmt.Errorf("BurnMintERC20Transparent at %s is already initialized for %s token on %s", config.BurnMintERC20Transparent, token, chain)
				}
			} else {
				chainState, ok := state.EVMChainState(chainSelector)
				if !ok {
					return fmt.Errorf("%s does not exist in state", chain.Name())
				}

				if _, ok := chainState.TransparentUpgradeableProxy[shared.TokenSymbol(config.Symbol)]; !ok {
					return fmt.Errorf("TransparentUpgradeableProxy does not exist for %s token symbol on %s", config.Symbol, chain.Name())
				}
			}
		}
	}

	return nil
}

func (c TransparentUpgradeableProxyGrantRoleChangesetConfig) Validate(e cldf.Environment) error {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, tokens := range c.Tokens {
		for _, roles := range tokens {
			chain, ok := e.BlockChains.EVMChains()[chainSelector]
			if !ok {
				return fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
			}

			chainState, ok := state.EVMChainState(chainSelector)
			if !ok {
				return fmt.Errorf("%s does not exist in state", chain.Name())
			}

			for _, config := range roles {
				if _, ok := chainState.TransparentUpgradeableProxy[shared.TokenSymbol(config.Symbol)]; !ok {
					return fmt.Errorf("TransparentUpgradeableProxy does not exist for %s token symbol on %s", config.Symbol, chain.Name())
				}

				if config.Role != TransparentUpgradeableProxyBurnerRole && config.Role != TransparentUpgradeableProxyMinterRole {
					return fmt.Errorf("invalid role %s for %s token symbol on %s", config.Role, config.Symbol, chain.Name())
				}
			}
		}
	}

	return nil
}

// DeployTransparentUpgradeableProxy deploys TransparentUpgradeableProxy contracts for the specified tokens on the specified chains.
func DeployTransparentUpgradeableProxy(e cldf.Environment, c TransparentUpgradeableProxyChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TransparentUpgradeableProxyChangesetConfig: %w", err)
	}

	addressBook := cldf.NewMemoryAddressBook()

	for chainSelector, tokens := range c.Tokens {
		chain := e.BlockChains.EVMChains()[chainSelector]

		for token, config := range tokens {
			_, err := cldf.DeployContract(e.Logger, chain, addressBook,
				func(chain cldf_evm.Chain) cldf.ContractDeploy[*transparent_upgradeable_proxy.TransparentUpgradeableProxy] {
					var errs []error

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
					if err != nil {
						errs = append(errs, fmt.Errorf("failed to parse ABI: %w", err))
					}

					initData, err := parsedABI.Pack(
						"initialize",
						token,
						config.Symbol,
						config.Decimals,
						config.MaxSupply,
						config.PreMint,
						chain.DeployerKey.From,
					)

					if err != nil {
						errs = append(errs, fmt.Errorf("failed to pack initialize data: %w", err))
					}

					address, tx, proxy, err := transparent_upgradeable_proxy.DeployTransparentUpgradeableProxy(
						chain.DeployerKey, chain.Client, config.BurnMintERC20Transparent, chain.DeployerKey.From, initData,
					)
					errs = append(errs, err)

					return cldf.ContractDeploy[*transparent_upgradeable_proxy.TransparentUpgradeableProxy]{
						Address:  address,
						Contract: proxy,
						Tv:       cldf.NewTypeAndVersion(shared.TransparentUpgradeableProxy, deployment.Version1_6_1),
						Tx:       tx,
						Err:      errors.Join(errs...),
					}
				},
			)

			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy TransparentUpgradeableProxy for %s token on %s: %w", token, chain.Name(), err)
			}
		}
	}

	ds, err := shared.PopulateDataStore(addressBook)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate in-memory DataStore: %w", err)
	}

	return cldf.ChangesetOutput{
		AddressBook: addressBook,
		DataStore:   ds,
	}, nil
}

// SaveProxyAdmin reads and saves the ProxyAdmin addresses from the TransparentUpgradeableProxy contracts for the specified tokens on the specified chains.
func SaveProxyAdmin(e cldf.Environment, c TransparentUpgradeableProxyChangesetConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	addressBook := cldf.NewMemoryAddressBook()

	for chainSelector, tokens := range c.Tokens {
		chain := e.BlockChains.EVMChains()[chainSelector]

		chainState, ok := state.EVMChainState(chainSelector)
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("%s does not exist in state", chain.Name())
		}

		for _, config := range tokens {
			proxy, ok := chainState.TransparentUpgradeableProxy[shared.TokenSymbol(config.Symbol)]
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("TransparentUpgradeableProxy does not exist for %s token on %s", config.Symbol, chain.Name())
			}

			storageBytes, err := chain.Client.StorageAt(e.GetContext(), proxy.Address(), shared.AdminSlot, nil)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get storage at slot %s for TransparentUpgradeableProxy at %s for %s token on %s: %w", shared.AdminSlot, proxy.Address(), config.Symbol, chain, err)
			}

			proxyAdmin := common.BytesToAddress(storageBytes)
			if err := addressBook.Save(chainSelector, proxyAdmin.String(), cldf.NewTypeAndVersion(shared.ProxyAdmin, deployment.Version1_6_1)); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ProxyAdmin at %s for %s token on %s: %w", proxyAdmin, config.Symbol, chain.Name(), err)
			}
		}
	}

	ds, err := shared.PopulateDataStore(addressBook)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate in-memory DataStore: %w", err)
	}

	return cldf.ChangesetOutput{
		AddressBook: addressBook,
		DataStore:   ds,
	}, nil
}

// GetRoleFromTransparentUpgradeableProxy retrieves the specified role from the TransparentUpgradeableProxy contract.
func GetRoleFromTransparentUpgradeableProxy(ctx context.Context, transparent *burn_mint_erc20_transparent.BurnMintERC20Transparent, role TransparentUpgradeableProxyRole) ([32]byte, error) {
	if transparent == nil {
		return [32]byte{}, errors.New("proxy is nil")
	}

	switch role {
	case TransparentUpgradeableProxyBurnerRole:
		r, err := transparent.BURNERROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to get BURNER_ROLE: %w", err)
		}
		return r, nil
	case TransparentUpgradeableProxyMinterRole:
		r, err := transparent.MINTERROLE(&bind.CallOpts{Context: ctx})
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to get MINTER_ROLE: %w", err)
		}
		return r, nil
	}

	return [32]byte{}, fmt.Errorf("invalid role: %s", role)
}

// GrantRoleTransparentUpgradeableProxy grants roles on TransparentUpgradeableProxy contracts for the specified tokens on the specified chains.
func GrantRoleTransparentUpgradeableProxy(e cldf.Environment, c TransparentUpgradeableProxyGrantRoleChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TransparentUpgradeableProxyGrantRoleChangesetConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(e, state, c.MCMS).WithDeploymentContext("grant roles to TransparentUpgradeableProxy")

	for chainSelector, tokens := range c.Tokens {
		chain := e.BlockChains.EVMChains()[chainSelector]

		chainState, ok := state.EVMChainState(chainSelector)
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("%s does not exist in state", chain)
		}

		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
		}

		for token, roles := range tokens {
			for _, config := range roles {
				proxy := chainState.TransparentUpgradeableProxy[shared.TokenSymbol(config.Symbol)]

				transparent, err := burn_mint_erc20_transparent.NewBurnMintERC20Transparent(proxy.Address(), chain.Client)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to instantiate BurnMintERC20Transparent at %s for %s token on %s: %w", proxy.Address(), token, chain, err)
				}

				r, err := GetRoleFromTransparentUpgradeableProxy(e.GetContext(), transparent, config.Role)
				if err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to get role from TransparentUpgradeableProxy at %s for %s token on %s: %w", proxy.Address(), token, chain.Name(), err)
				}

				if hasRole, err := transparent.HasRole(&bind.CallOpts{Context: e.GetContext()}, r, config.Account); err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to check if account %s has role %s on TransparentUpgradeableProxy at %s for %s token on %s: %w", config.Account, config.Role, proxy.Address(), token, chain.Name(), err)
				} else if hasRole {
					e.Logger.Infof("Account %s already has role %s on TransparentUpgradeableProxy at %s for %s token on %s, skipping grantRole", config.Account, config.Role, proxy.Address(), token, chain.Name())
					continue
				}

				if _, err := transparent.GrantRole(opts, r, config.Account); err != nil {
					return cldf.ChangesetOutput{}, fmt.Errorf("failed to create grantRole transaction for role %s to %s on TransparentUpgradeableProxy at %s for %s token on %s: %w", config.Role, config.Account, proxy.Address(), token, chain.Name(), err)
				}
			}
		}
	}

	return deployerGroup.Enact()
}

// BeginDefaultAdminTransferTransparentUpgradeableProxy begins the default admin transfer on TransparentUpgradeableProxy contracts for the specified tokens on the specified chains.
func BeginDefaultAdminTransferTransparentUpgradeableProxy(e cldf.Environment, c TransparentUpgradeableProxyChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TransparentUpgradeableProxyChangesetConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(e, state, c.MCMS).WithDeploymentContext("begin default admin transfer of TransparentUpgradeableProxy")

	for chainSelector, tokens := range c.Tokens {
		chain := e.BlockChains.EVMChains()[chainSelector]

		chainState, ok := state.EVMChainState(chainSelector)
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("%s does not exist in state", chain)
		}

		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
		}

		for token, config := range tokens {
			if config.NewAdmin == (common.Address{}) {
				return cldf.ChangesetOutput{}, fmt.Errorf("NewAdmin address is required for %s token on chain %s", token, chain.Name())
			}

			proxy, ok := chainState.TransparentUpgradeableProxy[shared.TokenSymbol(config.Symbol)]
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("TransparentUpgradeableProxy does not exist for %s token on %s", token, chain.Name())
			}

			implementation, err := burn_mint_erc20_transparent.NewBurnMintERC20Transparent(proxy.Address(), chain.Client)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to instantiate BurnMintERC20Transparent at %s for %s token on %s: %w", config.BurnMintERC20Transparent, token, chain.Name(), err)
			}

			if _, err := implementation.BeginDefaultAdminTransfer(opts, config.NewAdmin); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create beginDefaultAdminTransfer transaction to %s for BurnMintERC20Transparent at %s for %s token on %s: %w", config.NewAdmin, config.BurnMintERC20Transparent, token, chain.Name(), err)
			}
		}
	}

	return deployerGroup.Enact()
}

// AcceptDefaultAdminTransferTransparentUpgradeableProxy accepts the default admin transfer on TransparentUpgradeableProxy contracts for the specified tokens on the specified chains.
func AcceptDefaultAdminTransferTransparentUpgradeableProxy(e cldf.Environment, c TransparentUpgradeableProxyChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TransparentUpgradeableProxyChangesetConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(e, state, c.MCMS).WithDeploymentContext("accept default admin transfer of TransparentUpgradeableProxy")

	for chainSelector, tokens := range c.Tokens {
		chain := e.BlockChains.EVMChains()[chainSelector]

		chainState, ok := state.EVMChainState(chainSelector)
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("%s does not exist in state", chain)
		}

		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
		}

		for token, config := range tokens {
			proxy, ok := chainState.TransparentUpgradeableProxy[shared.TokenSymbol(config.Symbol)]
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("TransparentUpgradeableProxy does not exist for %s token on %s", token, chain.Name())
			}

			implementation, err := burn_mint_erc20_transparent.NewBurnMintERC20Transparent(proxy.Address(), chain.Client)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to instantiate BurnMintERC20Transparent at %s for %s token on %s: %w", proxy.Address(), token, chain.Name(), err)
			}

			if _, err := implementation.AcceptDefaultAdminTransfer(opts); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create acceptDefaultAdminTransfer transaction for BurnMintERC20Transparent at %s for %s token on %s: %w", proxy.Address(), token, chain.Name(), err)
			}
		}
	}

	return deployerGroup.Enact()
}

// SetCCIPAdminTransferTransparentUpgradeableProxy sets the CCIP admin on TransparentUpgradeableProxy contracts for the specified tokens on the specified chains.
func SetCCIPAdminTransferTransparentUpgradeableProxy(e cldf.Environment, c TransparentUpgradeableProxyChangesetConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid TransparentUpgradeableProxyChangesetConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(e, state, c.MCMS).WithDeploymentContext("set CCIP admin for TransparentUpgradeableProxy")

	for chainSelector, tokens := range c.Tokens {
		chain := e.BlockChains.EVMChains()[chainSelector]

		chainState, ok := state.EVMChainState(chainSelector)
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("%s does not exist in state", chain)
		}

		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for %s", chain)
		}

		for token, config := range tokens {
			if config.NewAdmin == (common.Address{}) {
				return cldf.ChangesetOutput{}, fmt.Errorf("NewAdmin address is required for %s token on chain %s", token, chain.Name())
			}

			proxy, ok := chainState.TransparentUpgradeableProxy[shared.TokenSymbol(config.Symbol)]
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("TransparentUpgradeableProxy does not exist for %s token on %s", token, chain.Name())
			}

			implementation, err := burn_mint_erc20_transparent.NewBurnMintERC20Transparent(proxy.Address(), chain.Client)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to instantiate BurnMintERC20Transparent at %s for %s token on %s: %w", proxy.Address(), token, chain.Name(), err)
			}

			if _, err := implementation.SetCCIPAdmin(opts, config.NewAdmin); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create setCCIPadmin transaction to %s for BurnMintERC20Transparent at %s for %s token on %s: %w", config.NewAdmin, proxy.Address(), token, chain.Name(), err)
			}
		}
	}

	return deployerGroup.Enact()
}
