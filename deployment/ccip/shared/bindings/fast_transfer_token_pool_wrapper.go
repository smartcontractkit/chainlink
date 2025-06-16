package bindings

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	fast_transfer_token_pool "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/fast_transfer_token_pool"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	burn_mint_external "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/burn_mint_with_external_minter_fast_transfer_token_pool"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// Re-exported types to provide a clean API boundary
type (
	// DestChainConfig represents destination chain configuration
	DestChainConfig = burn_mint_external.FastTransferTokenPoolAbstractDestChainConfig

	// DestChainConfigUpdateArgs represents arguments for updating destination chain configuration
	DestChainConfigUpdateArgs = burn_mint_external.FastTransferTokenPoolAbstractDestChainConfigUpdateArgs

	// Quote represents a fee quote for fast transfer operations
	Quote = burn_mint_external.IFastTransferPoolQuote
)

// FastTransferTokenPoolWrapper provides a unified interface for both
// BurnMintFastTransferTokenPool and BurnMintWithExternalMinterFastTransferTokenPool
type FastTransferTokenPoolWrapper struct {
	contractType cldf.ContractType
	address      common.Address

	// Underlying contract instances (only one will be non-nil)
	burnMintPool         *fast_transfer_token_pool.BurnMintFastTransferTokenPool
	burnMintExternalPool *burn_mint_external.BurnMintWithExternalMinterFastTransferTokenPool
}

// NewFastTransferTokenPoolWrapper creates a new wrapper instance
func NewFastTransferTokenPoolWrapper(
	address common.Address,
	backend bind.ContractBackend,
	contractType cldf.ContractType,
) (*FastTransferTokenPoolWrapper, error) {
	wrapper := &FastTransferTokenPoolWrapper{
		contractType: contractType,
		address:      address,
	}

	switch contractType {
	case shared.BurnMintFastTransferTokenPool:
		pool, err := fast_transfer_token_pool.NewBurnMintFastTransferTokenPool(address, backend)
		if err != nil {
			return nil, err
		}
		wrapper.burnMintPool = pool
	case shared.BurnMintWithExternalMinterFastTransferTokenPool:
		pool, err := burn_mint_external.NewBurnMintWithExternalMinterFastTransferTokenPool(address, backend)
		if err != nil {
			return nil, err
		}
		wrapper.burnMintExternalPool = pool
	default:
		return nil, errors.New("unsupported contract type")
	}

	return wrapper, nil
}

// Address returns the contract address
func (w *FastTransferTokenPoolWrapper) Address() common.Address {
	return w.address
}

// ContractType returns the underlying contract type
func (w *FastTransferTokenPoolWrapper) ContractType() cldf.ContractType {
	return w.contractType
}

// GetDestChainConfig retrieves destination chain configuration
func (w *FastTransferTokenPoolWrapper) GetDestChainConfig(
	opts *bind.CallOpts,
	remoteChainSelector uint64,
) (DestChainConfig, []common.Address, error) {
	switch w.contractType {
	case shared.BurnMintFastTransferTokenPool:
		config, addresses, err := w.burnMintPool.GetDestChainConfig(opts, remoteChainSelector)
		if err != nil {
			return DestChainConfig{}, nil, err
		}
		// Convert from fast_transfer_token_pool to burn_mint_external types
		convertedConfig := DestChainConfig{
			MaxFillAmountPerRequest:  config.MaxFillAmountPerRequest,
			FillerAllowlistEnabled:   config.FillerAllowlistEnabled,
			FastTransferFillerFeeBps: config.FastTransferFillerFeeBps,
			FastTransferPoolFeeBps:   config.FastTransferPoolFeeBps,
			SettlementOverheadGas:    config.SettlementOverheadGas,
			DestinationPool:          config.DestinationPool,
			CustomExtraArgs:          config.CustomExtraArgs,
		}
		return convertedConfig, addresses, nil
	case shared.BurnMintWithExternalMinterFastTransferTokenPool:
		return w.burnMintExternalPool.GetDestChainConfig(opts, remoteChainSelector)
	default:
		return DestChainConfig{}, nil, errors.New("unsupported contract type")
	}
}

// UpdateDestChainConfig updates destination chain configurations
func (w *FastTransferTokenPoolWrapper) UpdateDestChainConfig(
	opts *bind.TransactOpts,
	updates []DestChainConfigUpdateArgs,
) (*types.Transaction, error) {
	switch w.contractType {
	case shared.BurnMintFastTransferTokenPool:
		// Convert from burn_mint_external to fast_transfer_token_pool types
		convertedUpdates := make([]fast_transfer_token_pool.FastTransferTokenPoolAbstractDestChainConfigUpdateArgs, len(updates))
		for i, update := range updates {
			convertedUpdates[i] = fast_transfer_token_pool.FastTransferTokenPoolAbstractDestChainConfigUpdateArgs{
				FillerAllowlistEnabled:   update.FillerAllowlistEnabled,
				FastTransferFillerFeeBps: update.FastTransferFillerFeeBps,
				FastTransferPoolFeeBps:   update.FastTransferPoolFeeBps,
				SettlementOverheadGas:    update.SettlementOverheadGas,
				RemoteChainSelector:      update.RemoteChainSelector,
				ChainFamilySelector:      update.ChainFamilySelector,
				MaxFillAmountPerRequest:  update.MaxFillAmountPerRequest,
				DestinationPool:          update.DestinationPool,
				CustomExtraArgs:          update.CustomExtraArgs,
			}
		}
		return w.burnMintPool.UpdateDestChainConfig(opts, convertedUpdates)
	case shared.BurnMintWithExternalMinterFastTransferTokenPool:
		return w.burnMintExternalPool.UpdateDestChainConfig(opts, updates)
	default:
		return nil, errors.New("unsupported contract type")
	}
}

// GetAllowedFillers retrieves the list of allowed filler addresses
func (w *FastTransferTokenPoolWrapper) GetAllowedFillers(opts *bind.CallOpts) ([]common.Address, error) {
	switch w.contractType {
	case shared.BurnMintFastTransferTokenPool:
		return w.burnMintPool.GetAllowedFillers(opts)
	case shared.BurnMintWithExternalMinterFastTransferTokenPool:
		return w.burnMintExternalPool.GetAllowedFillers(opts)
	default:
		return nil, errors.New("unsupported contract type")
	}
}

// UpdateFillerAllowList updates the filler allowlist
func (w *FastTransferTokenPoolWrapper) UpdateFillerAllowList(
	opts *bind.TransactOpts,
	fillersToAdd []common.Address,
	fillersToRemove []common.Address,
) (*types.Transaction, error) {
	switch w.contractType {
	case shared.BurnMintFastTransferTokenPool:
		return w.burnMintPool.UpdateFillerAllowList(opts, fillersToAdd, fillersToRemove)
	case shared.BurnMintWithExternalMinterFastTransferTokenPool:
		return w.burnMintExternalPool.UpdateFillerAllowList(opts, fillersToAdd, fillersToRemove)
	default:
		return nil, errors.New("unsupported contract type")
	}
}

// IsAllowedFiller checks if an address is an allowed filler
func (w *FastTransferTokenPoolWrapper) IsAllowedFiller(opts *bind.CallOpts, filler common.Address) (bool, error) {
	switch w.contractType {
	case shared.BurnMintFastTransferTokenPool:
		return w.burnMintPool.IsAllowedFiller(opts, filler)
	case shared.BurnMintWithExternalMinterFastTransferTokenPool:
		return w.burnMintExternalPool.IsAllowedFiller(opts, filler)
	default:
		return false, errors.New("unsupported contract type")
	}
}

// CcipSendToken initiates a fast transfer (required for e2e tests)
func (w *FastTransferTokenPoolWrapper) CcipSendToken(
	opts *bind.TransactOpts,
	destinationChainSelector uint64,
	amount *big.Int,
	maxFastTransferFee *big.Int,
	receiver []byte,
	feeToken common.Address,
	extraArgs []byte,
) (*types.Transaction, error) {
	switch w.contractType {
	case shared.BurnMintFastTransferTokenPool:
		return w.burnMintPool.CcipSendToken(opts, destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
	case shared.BurnMintWithExternalMinterFastTransferTokenPool:
		return w.burnMintExternalPool.CcipSendToken(opts, destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
	default:
		return nil, errors.New("unsupported contract type")
	}
}

// GetCcipSendTokenFee calculates fees for sending tokens (required for e2e tests)
func (w *FastTransferTokenPoolWrapper) GetCcipSendTokenFee(
	opts *bind.CallOpts,
	destinationChainSelector uint64,
	amount *big.Int,
	receiver []byte,
	settlementFeeToken common.Address,
	extraArgs []byte,
) (Quote, error) {
	switch w.contractType {
	case shared.BurnMintFastTransferTokenPool:
		quote, err := w.burnMintPool.GetCcipSendTokenFee(opts, destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)
		if err != nil {
			return Quote{}, err
		}
		// Convert from fast_transfer_token_pool to burn_mint_external types
		convertedQuote := Quote{
			CcipSettlementFee: quote.CcipSettlementFee,
			FastTransferFee:   quote.FastTransferFee,
		}
		return convertedQuote, nil
	case shared.BurnMintWithExternalMinterFastTransferTokenPool:
		return w.burnMintExternalPool.GetCcipSendTokenFee(opts, destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)
	default:
		return Quote{}, errors.New("unsupported contract type")
	}
}

// FilterFastTransferRequested filters FastTransferRequested events (required for e2e tests)
func (w *FastTransferTokenPoolWrapper) FilterFastTransferRequested(
	opts *bind.FilterOpts,
	destinationChainSelector []uint64,
	fillID [][32]byte,
	settlementID [][32]byte,
) (*FastTransferRequestedIterator, error) {
	switch w.contractType {
	case shared.BurnMintFastTransferTokenPool:
		iter, err := w.burnMintPool.FilterFastTransferRequested(opts, destinationChainSelector, fillID, settlementID)
		if err != nil {
			return nil, err
		}
		return &FastTransferRequestedIterator{burnMintIter: iter}, nil
	case shared.BurnMintWithExternalMinterFastTransferTokenPool:
		iter, err := w.burnMintExternalPool.FilterFastTransferRequested(opts, destinationChainSelector, fillID, settlementID)
		if err != nil {
			return nil, err
		}
		return &FastTransferRequestedIterator{burnMintExternalIter: iter}, nil
	default:
		return nil, errors.New("unsupported contract type")
	}
}

// FastTransferRequestedIterator wraps both iterator types
type FastTransferRequestedIterator struct {
	burnMintIter         *fast_transfer_token_pool.BurnMintFastTransferTokenPoolFastTransferRequestedIterator
	burnMintExternalIter *burn_mint_external.BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator
}

// Next advances the iterator
func (it *FastTransferRequestedIterator) Next() bool {
	if it.burnMintIter != nil {
		return it.burnMintIter.Next()
	}
	if it.burnMintExternalIter != nil {
		return it.burnMintExternalIter.Next()
	}
	return false
}

// Error returns the error (if any)
func (it *FastTransferRequestedIterator) Error() error {
	if it.burnMintIter != nil {
		return it.burnMintIter.Error()
	}
	if it.burnMintExternalIter != nil {
		return it.burnMintExternalIter.Error()
	}
	return nil
}

// Close closes the iterator
func (it *FastTransferRequestedIterator) Close() error {
	if it.burnMintIter != nil {
		return it.burnMintIter.Close()
	}
	if it.burnMintExternalIter != nil {
		return it.burnMintExternalIter.Close()
	}
	return nil
}

// Event returns the current event data
func (it *FastTransferRequestedIterator) Event() *FastTransferRequestedEvent {
	if it.burnMintIter != nil && it.burnMintIter.Event != nil {
		return &FastTransferRequestedEvent{
			DestinationChainSelector: it.burnMintIter.Event.DestinationChainSelector,
			FillID:                   it.burnMintIter.Event.FillId,
			SettlementID:             it.burnMintIter.Event.SettlementId,
			SourceAmountNetFee:       it.burnMintIter.Event.SourceAmountNetFee,
			SourceDecimals:           it.burnMintIter.Event.SourceDecimals,
			FastTransferFee:          it.burnMintIter.Event.FastTransferFee,
			Receiver:                 it.burnMintIter.Event.Receiver,
			Raw:                      it.burnMintIter.Event.Raw,
		}
	}
	if it.burnMintExternalIter != nil && it.burnMintExternalIter.Event != nil {
		return &FastTransferRequestedEvent{
			DestinationChainSelector: it.burnMintExternalIter.Event.DestinationChainSelector,
			FillID:                   it.burnMintExternalIter.Event.FillId,
			SettlementID:             it.burnMintExternalIter.Event.SettlementId,
			SourceAmountNetFee:       it.burnMintExternalIter.Event.SourceAmountNetFee,
			SourceDecimals:           it.burnMintExternalIter.Event.SourceDecimals,
			FastTransferFee:          it.burnMintExternalIter.Event.FastTransferFee,
			Receiver:                 it.burnMintExternalIter.Event.Receiver,
			Raw:                      it.burnMintExternalIter.Event.Raw,
		}
	}
	return nil
}

// FastTransferRequestedEvent represents a unified event structure
type FastTransferRequestedEvent struct {
	DestinationChainSelector uint64
	FillID                   [32]byte
	SettlementID             [32]byte
	SourceAmountNetFee       *big.Int
	SourceDecimals           uint8
	FastTransferFee          *big.Int
	Receiver                 []byte
	Raw                      types.Log
}

func GetFastTransferTokenPoolContract(env cldf.Environment, tokenSymbol shared.TokenSymbol, contractType cldf.ContractType, contractVersion semver.Version, chainSelector uint64) (*FastTransferTokenPoolWrapper, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chain, ok := env.BlockChains.EVMChains()[chainSelector]
	if !ok {
		return nil, fmt.Errorf("chain with selector %d does not exist in environment", chainSelector)
	}

	chainState, ok := state.Chains[chainSelector]
	if !ok {
		return nil, fmt.Errorf("%s does not exist in state", chain.String())
	}

	switch contractType {
	case shared.BurnMintFastTransferTokenPool:
		pool, ok := chainState.BurnMintFastTransferTokenPools[tokenSymbol][contractVersion]
		if !ok {
			return nil, fmt.Errorf("burn mint fast transfer token pool for token %s and version %s not found on chain %s", tokenSymbol, contractVersion, chain.String())
		}
		return NewFastTransferTokenPoolWrapper(pool.Address(), env.BlockChains.EVMChains()[chainSelector].Client, contractType)
	case shared.BurnMintWithExternalMinterFastTransferTokenPool:
		pool, ok := chainState.BurnMintWithExternalMinterFastTransferTokenPools[tokenSymbol][contractVersion]
		if !ok {
			return nil, fmt.Errorf("burn mint with external minter fast transfer token pool for token %s and version %s not found on chain %s", tokenSymbol, contractVersion, chain.String())
		}
		return NewFastTransferTokenPoolWrapper(pool.Address(), env.BlockChains.EVMChains()[chainSelector].Client, contractType)
	default:
		return nil, fmt.Errorf("unsupported contract type %s for fast transfer token pools", contractType)
	}
}
