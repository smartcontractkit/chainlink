package aptos

import (
	"errors"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"golang.org/x/sync/errgroup"

	module_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp/offramp"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/regulated_token_pool"
	aptosHelpers "github.com/smartcontractkit/chainlink-aptos/bindings/helpers"
	"github.com/smartcontractkit/chainlink-aptos/bindings/managed_token"
	"github.com/smartcontractkit/chainlink-aptos/bindings/regulated_token"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	"github.com/smartcontractkit/chainlink/deployment/ccip/view"
	aptosview "github.com/smartcontractkit/chainlink/deployment/ccip/view/aptos"
	"github.com/smartcontractkit/chainlink/deployment/helpers"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/managed_token_pool"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

type CCIPChainState struct {
	MCMSAddress aptos.AccountAddress
	CCIPAddress aptos.AccountAddress

	LinkTokenAddress aptos.AccountAddress
	ManagedTokens    map[shared.TokenSymbol]aptos.AccountAddress

	AptosManagedTokenPools map[aptos.AccountAddress]aptos.AccountAddress // TokenAddress -> TokenPoolAddress
	RegulatedTokenPools    map[aptos.AccountAddress]aptos.AccountAddress // TokenAddress -> TokenPoolAddress
	BurnMintTokenPools     map[aptos.AccountAddress]aptos.AccountAddress // TokenAddress -> TokenPoolAddress
	LockReleaseTokenPools  map[aptos.AccountAddress]aptos.AccountAddress // TokenAddress -> TokenPoolAddress

	// Test contracts
	TestRouterAddress aptos.AccountAddress
	ReceiverAddress   aptos.AccountAddress
}

// LoadOnchainStateAptos loads chain state for Aptos chains from env
func LoadOnchainStateAptos(env cldf.Environment) (map[uint64]CCIPChainState, error) {
	aptosChains := make(map[uint64]CCIPChainState)
	for chainSelector := range env.BlockChains.AptosChains() {
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, cldf.ErrChainNotFound) {
				return aptosChains, err
			}
			addresses = make(map[string]cldf.TypeAndVersion)
		}
		chainState, err := loadAptosChainStateFromAddresses(addresses, env.BlockChains.AptosChains()[chainSelector].Client)
		if err != nil {
			return aptosChains, err
		}
		aptosChains[chainSelector] = chainState
	}
	return aptosChains, nil
}

func loadAptosChainStateFromAddresses(addresses map[string]cldf.TypeAndVersion, client aptos.AptosRpcClient) (CCIPChainState, error) {
	chainState := CCIPChainState{
		ManagedTokens:          make(map[shared.TokenSymbol]aptos.AccountAddress),
		AptosManagedTokenPools: make(map[aptos.AccountAddress]aptos.AccountAddress),
		RegulatedTokenPools:    make(map[aptos.AccountAddress]aptos.AccountAddress),
		BurnMintTokenPools:     make(map[aptos.AccountAddress]aptos.AccountAddress),
		LockReleaseTokenPools:  make(map[aptos.AccountAddress]aptos.AccountAddress),
	}

	for addrStr, typeAndVersion := range addresses {
		// Parse address
		address := &aptos.AccountAddress{}
		err := address.ParseStringRelaxed(addrStr)
		if err != nil {
			return chainState, fmt.Errorf("failed to parse address %s for %s: %w", addrStr, typeAndVersion.Type, err)
		}
		// Set address based on type
		switch typeAndVersion.Type {
		case shared.AptosMCMSType:
			chainState.MCMSAddress = *address
		case shared.AptosCCIPType:
			chainState.CCIPAddress = *address
		case types.LinkToken:
			chainState.LinkTokenAddress = *address
		case shared.AptosReceiverType:
			chainState.ReceiverAddress = *address
		case shared.AptosManagedTokenType:
			noLabel := typeAndVersion.Labels.IsEmpty()
			symbol := shared.TokenSymbol("")
			if noLabel {
				token := managed_token.Bind(*address, client)
				metadataAddress, err := token.ManagedToken().TokenMetadata(nil)
				if err != nil {
					return chainState, fmt.Errorf("failed to get FA Metadata for ManagedToken %s: %w", addrStr, err)
				}
				metadata, err := aptosHelpers.GetFungibleAssetMetadata(client, metadataAddress)
				if err != nil {
					return chainState, fmt.Errorf("failed to get Fungible Asset Metadata for Managed Token %s: %w", addrStr, err)
				}
				symbol = shared.TokenSymbol(metadata.Symbol)
			} else {
				symbol = shared.TokenSymbol(typeAndVersion.Labels.List()[0])
			}
			chainState.ManagedTokens[symbol] = *address
		case shared.AptosRegulatedTokenType:
			noLabel := typeAndVersion.Labels.IsEmpty()
			symbol := shared.TokenSymbol("")
			if noLabel {
				token := regulated_token.Bind(*address, client)
				metadataAddress, err := token.RegulatedToken().TokenMetadata(nil)
				if err != nil {
					return chainState, fmt.Errorf("failed to get FA Metadata for RegulatedToken %s: %w", addrStr, err)
				}
				metadata, err := aptosHelpers.GetFungibleAssetMetadata(client, metadataAddress)
				if err != nil {
					return chainState, fmt.Errorf("failed to get Fungible Asset Metadata for RegulatedToken %s: %w", addrStr, err)
				}
				symbol = shared.TokenSymbol(metadata.Symbol)
			} else {
				symbol = shared.TokenSymbol(typeAndVersion.Labels.List()[0])
			}
			chainState.ManagedTokens[symbol] = *address
		case shared.AptosManagedTokenPoolType:
			noLabel := typeAndVersion.Labels.IsEmpty()
			token := aptos.AccountAddress{}
			if noLabel {
				pool := managed_token_pool.Bind(*address, client)
				t, err := pool.ManagedTokenPool().GetToken(nil)
				if err != nil {
					return chainState, fmt.Errorf("failed to get token for ManagedTokenPool %s: %w", addrStr, err)
				}
				token = t
			} else {
				labels := typeAndVersion.Labels.List()
				tokenStr := labels[0]
				err := token.ParseStringRelaxed(tokenStr)
				if err != nil {
					return chainState, fmt.Errorf("failed to parse token address %s for ManagedTokenPool %s: %w", tokenStr, addrStr, err)
				}
			}
			chainState.AptosManagedTokenPools[token] = *address
		case shared.AptosRegulatedTokenPoolType:
			noLabel := typeAndVersion.Labels.IsEmpty()
			token := aptos.AccountAddress{}
			if noLabel {
				pool := regulated_token_pool.Bind(*address, client)
				t, err := pool.RegulatedTokenPool().GetToken(nil)
				if err != nil {
					return chainState, fmt.Errorf("failed to get token for RegulatedTokenPool: %s: %w", addrStr, err)
				}
				token = t
			} else {
				labels := typeAndVersion.Labels.List()
				tokenStr := labels[0]
				err := token.ParseStringRelaxed(tokenStr)
				if err != nil {
					return chainState, fmt.Errorf("failed to parse token address %s for RegulatedTokenPool %s: %w", tokenStr, addrStr, err)
				}
			}
			chainState.RegulatedTokenPools[token] = *address
		case shared.BurnMintTokenPool:
			noLabel := typeAndVersion.Labels.IsEmpty()
			token := aptos.AccountAddress{}
			if noLabel {
				pool := burn_mint_token_pool.Bind(*address, client)
				t, err := pool.BurnMintTokenPool().GetToken(nil)
				if err != nil {
					return chainState, fmt.Errorf("failed to get token for BurnMintTokenPool %s: %w", addrStr, err)
				}
				token = t
			} else {
				labels := typeAndVersion.Labels.List()
				tokenStr := labels[0]
				err := token.ParseStringRelaxed(tokenStr)
				if err != nil {
					return chainState, fmt.Errorf("failed to parse token address %s for BurnMintTokenPool %s: %w", tokenStr, addrStr, err)
				}
			}
			chainState.BurnMintTokenPools[token] = *address
		case shared.LockReleaseTokenPool:
			noLabel := typeAndVersion.Labels.IsEmpty()
			token := aptos.AccountAddress{}
			if noLabel {
				pool := lock_release_token_pool.Bind(*address, client)
				t, err := pool.LockReleaseTokenPool().GetToken(nil)
				if err != nil {
					return chainState, fmt.Errorf("failed to get token for LockReleaseTokenPool %s: %w", addrStr, err)
				}
				token = t
			} else {
				labels := typeAndVersion.Labels.List()
				tokenStr := labels[0]
				err := token.ParseStringRelaxed(tokenStr)
				if err != nil {
					return chainState, fmt.Errorf("failed to parse token address %s for LockReleaseTokenPool %s: %w", tokenStr, addrStr, err)
				}
			}
			chainState.LockReleaseTokenPools[token] = *address
		}
	}
	return chainState, nil
}

func GetOfframpDynamicConfig(c cldf_aptos.Chain, ccipAddress aptos.AccountAddress) (module_offramp.DynamicConfig, error) {
	offrampBind := ccip_offramp.Bind(ccipAddress, c.Client)
	return offrampBind.Offramp().GetDynamicConfig(&bind.CallOpts{})
}

func FindAptosAddress(tv cldf.TypeAndVersion, addresses map[string]cldf.TypeAndVersion) aptos.AccountAddress {
	for address, tvStr := range addresses {
		if tv.String() == tvStr.String() {
			addr := aptos.AccountAddress{}
			_ = addr.ParseStringRelaxed(address)
			return addr
		}
	}
	return aptos.AccountAddress{}
}

func (s CCIPChainState) GenerateView(e *cldf.Environment, selector uint64, chainName string) (view.AptosChainView, error) {
	lggr := e.Logger
	chain := e.BlockChains.AptosChains()[selector]
	chainView := view.NewAptosChainView()
	errGroup := errgroup.Group{}
	lggr.Infow("generating Aptos chain view",
		"chain", chain.Name,
		"selector", selector)

	// Tokens
	errGroup.Go(func() error {
		for symbol, address := range s.ManagedTokens {
			tokenView, err := aptosview.GenerateTokenView(chain, address)
			if err != nil {
				return fmt.Errorf("failed to generate token view for managed token (%s) %s: %w", symbol, address.StringLong(), err)
			}
			chainView.UpdateMu.Lock()
			if symbol == shared.LinkSymbol {
				chainView.LinkToken = tokenView
			} else {
				chainView.Tokens[symbol.String()] = tokenView
			}
			chainView.UpdateMu.Unlock()
			lggr.Infow("generated token view", "tokenAddress", address.StringLong(), "symbol", symbol, "chain", chainName)
		}
		return nil
	})

	// MCMS
	errGroup.Go(func() error {
		mcmsView, err := aptosview.GenerateMCMSWithTimelockView(chain, s.MCMSAddress)
		if err != nil {
			return fmt.Errorf("failed to generate MCMS with timelock view for MCMS with timelock %s: %w", s.MCMSAddress.StringLong(), err)
		}
		chainView.UpdateMu.Lock()
		chainView.MCMSWithTimelock = mcmsView
		chainView.UpdateMu.Unlock()
		lggr.Infow("generated MCMS with timelock view", "MCMSAddress", s.MCMSAddress.StringLong(), "chain", chainName)
		return nil
	})

	// CCIP
	errGroup.Go(func() error {
		ccipView, err := aptosview.GenerateCCIPView(chain, s.CCIPAddress, s.CCIPAddress)
		if err != nil {
			return fmt.Errorf("failed to generate CCIP view for CCIP contract %s: %w", s.CCIPAddress.StringLong(), err)
		}
		chainView.UpdateMu.Lock()
		chainView.CCIP = ccipView
		chainView.UpdateMu.Unlock()
		lggr.Infow("generated CCIP view", "CCIPAddress", s.CCIPAddress.StringLong(), "chain", chainName)
		return nil
	})

	errGroup.Go(func() error {
		routerView, err := aptosview.GenerateRouterView(chain, s.CCIPAddress, []aptos.AccountAddress{s.CCIPAddress}, false)
		if err != nil {
			return fmt.Errorf("failed to generate router view for router %s: %w", s.CCIPAddress.StringLong(), err)
		}
		chainView.UpdateMu.Lock()
		chainView.Router[s.CCIPAddress.StringLong()] = routerView
		chainView.UpdateMu.Unlock()
		lggr.Infow("generated router view", "routerAddress", s.CCIPAddress.StringLong(), "chain", chainName)
		return nil
	})

	errGroup.Go(func() error {
		onRampView, err := aptosview.GenerateOnRampView(chain, s.CCIPAddress, s.CCIPAddress, s.CCIPAddress)
		if err != nil {
			return fmt.Errorf("failed to generate OnRamp view for OnRamp contract %s: %w", s.CCIPAddress.StringLong(), err)
		}
		chainView.UpdateMu.Lock()
		chainView.OnRamp[s.CCIPAddress.StringLong()] = onRampView
		chainView.UpdateMu.Unlock()
		lggr.Infow("generated onRamp view", "onRampAddress", s.CCIPAddress.StringLong(), "chain", chainName)
		return nil
	})

	errGroup.Go(func() error {
		offRampView, err := aptosview.GenerateOffRampView(chain, s.CCIPAddress, s.CCIPAddress)
		if err != nil {
			return fmt.Errorf("failed to generate OffRamp view for OffRamp contract %s: %w", s.CCIPAddress.StringLong(), err)
		}
		chainView.UpdateMu.Lock()
		chainView.OffRamp[s.CCIPAddress.StringLong()] = offRampView
		chainView.UpdateMu.Unlock()
		lggr.Infow("gneerated offRamp view", "offRampAddress", s.CCIPAddress.StringLong(), "chain", chainName)
		return nil
	})

	// Token pools
	errGroup.Go(func() error {
		for tokenAddress, tokenPoolAddress := range s.AptosManagedTokenPools {
			faMetadata, err := aptosHelpers.GetFungibleAssetMetadata(chain.Client, tokenAddress)
			if err != nil {
				return fmt.Errorf("failed to get fungible asset metadata for token %s: %w", tokenAddress, err)
			}
			contract := managed_token_pool.Bind(tokenPoolAddress, chain.Client)
			tokenPoolView, err := aptosview.GenerateTokenPoolView(chain, tokenPoolAddress, contract.ManagedTokenPool())
			if err != nil {
				return fmt.Errorf("failed to generate token pool view for ManagedTokenPool %s: %w", tokenPoolAddress.StringLong(), err)
			}
			chainView.UpdateMu.Lock()
			chainView.TokenPools = helpers.AddValueToNestedMap(chainView.TokenPools, faMetadata.Symbol, tokenPoolAddress.StringLong(), tokenPoolView)
			chainView.UpdateMu.Unlock()
		}
		return nil
	})
	errGroup.Go(func() error {
		for tokenAddress, tokenPoolAddress := range s.RegulatedTokenPools {
			faMetadata, err := aptosHelpers.GetFungibleAssetMetadata(chain.Client, tokenAddress)
			if err != nil {
				return fmt.Errorf("failed to get fungible asset metadata for token %s: %w", tokenAddress, err)
			}
			contract := regulated_token_pool.Bind(tokenPoolAddress, chain.Client)
			tokenPoolView, err := aptosview.GenerateTokenPoolView(chain, tokenPoolAddress, contract.RegulatedTokenPool())
			if err != nil {
				return fmt.Errorf("failed to generate token pool view for RegulatedTokenPool %s: %w", tokenPoolAddress.StringLong(), err)
			}
			chainView.UpdateMu.Lock()
			chainView.TokenPools = helpers.AddValueToNestedMap(chainView.TokenPools, faMetadata.Symbol, tokenPoolAddress.StringLong(), tokenPoolView)
			chainView.UpdateMu.Unlock()
		}
		return nil
	})
	errGroup.Go(func() error {
		for tokenAddress, tokenPoolAddress := range s.BurnMintTokenPools {
			faMetadata, err := aptosHelpers.GetFungibleAssetMetadata(chain.Client, tokenAddress)
			if err != nil {
				return fmt.Errorf("failed to get fungible asset metadata for token %s: %w", tokenAddress, err)
			}
			contract := burn_mint_token_pool.Bind(tokenPoolAddress, chain.Client)
			tokenPoolView, err := aptosview.GenerateTokenPoolView(chain, tokenPoolAddress, contract.BurnMintTokenPool())
			if err != nil {
				return fmt.Errorf("failed to generate token pool view for BurnMintTokenPool %s: %w", tokenPoolAddress.StringLong(), err)
			}
			chainView.UpdateMu.Lock()
			chainView.TokenPools = helpers.AddValueToNestedMap(chainView.TokenPools, faMetadata.Symbol, tokenPoolAddress.StringLong(), tokenPoolView)
			chainView.UpdateMu.Unlock()
		}
		return nil
	})
	errGroup.Go(func() error {
		for tokenAddress, tokenPoolAddress := range s.LockReleaseTokenPools {
			faMetadata, err := aptosHelpers.GetFungibleAssetMetadata(chain.Client, tokenAddress)
			if err != nil {
				return fmt.Errorf("failed to get fungible asset metadata for token %s: %w", tokenAddress, err)
			}
			contract := lock_release_token_pool.Bind(tokenPoolAddress, chain.Client)
			tokenPoolView, err := aptosview.GenerateTokenPoolView(chain, tokenPoolAddress, contract.LockReleaseTokenPool())
			if err != nil {
				return fmt.Errorf("failed to generate token pool view for LockReleaseTokenPool %s: %w", tokenPoolAddress.StringLong(), err)
			}
			chainView.UpdateMu.Lock()
			chainView.TokenPools = helpers.AddValueToNestedMap(chainView.TokenPools, faMetadata.Symbol, tokenPoolAddress.StringLong(), tokenPoolView)
			chainView.UpdateMu.Unlock()
		}
		return nil
	})

	return chainView, errGroup.Wait()
}
