package bindings

import (
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
	hybrid_external "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/hybrid_with_external_minter_fast_transfer_token_pool"
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

// FastTransferPool defines the common interface for all fast transfer pool types
type FastTransferPool interface {
	GetDestChainConfig(opts *bind.CallOpts, remoteChainSelector uint64) (DestChainConfig, []common.Address, error)
	UpdateDestChainConfig(opts *bind.TransactOpts, updates []DestChainConfigUpdateArgs) (*types.Transaction, error)
	GetAllowedFillers(opts *bind.CallOpts) ([]common.Address, error)
	GetToken(opts *bind.CallOpts) (common.Address, error)
	UpdateFillerAllowList(opts *bind.TransactOpts, fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error)
	IsAllowedFiller(opts *bind.CallOpts, filler common.Address) (bool, error)
	CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error)
	GetCcipSendTokenFee(opts *bind.CallOpts, destinationChainSelector uint64, amount *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (Quote, error)
	FilterFastTransferRequested(opts *bind.FilterOpts, destinationChainSelector []uint64, fillID [][32]byte, settlementID [][32]byte) (*FastTransferRequestedIterator, error)
	GetAccumulatedPoolFees(opts *bind.CallOpts) (*big.Int, error)
	WithdrawPoolFees(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error)
}

// Conversion functions for type compatibility
func convertDestChainConfigFromBurnMint(config fast_transfer_token_pool.FastTransferTokenPoolAbstractDestChainConfig) DestChainConfig {
	return DestChainConfig{
		MaxFillAmountPerRequest:  config.MaxFillAmountPerRequest,
		FillerAllowlistEnabled:   config.FillerAllowlistEnabled,
		FastTransferFillerFeeBps: config.FastTransferFillerFeeBps,
		FastTransferPoolFeeBps:   config.FastTransferPoolFeeBps,
		SettlementOverheadGas:    config.SettlementOverheadGas,
		DestinationPool:          config.DestinationPool,
		CustomExtraArgs:          config.CustomExtraArgs,
	}
}

func convertDestChainConfigFromHybrid(config hybrid_external.FastTransferTokenPoolAbstractDestChainConfig) DestChainConfig {
	return DestChainConfig{
		MaxFillAmountPerRequest:  config.MaxFillAmountPerRequest,
		FillerAllowlistEnabled:   config.FillerAllowlistEnabled,
		FastTransferFillerFeeBps: config.FastTransferFillerFeeBps,
		FastTransferPoolFeeBps:   config.FastTransferPoolFeeBps,
		SettlementOverheadGas:    config.SettlementOverheadGas,
		DestinationPool:          config.DestinationPool,
		CustomExtraArgs:          config.CustomExtraArgs,
	}
}

func convertUpdateArgsToBurnMint(updates []DestChainConfigUpdateArgs) []fast_transfer_token_pool.FastTransferTokenPoolAbstractDestChainConfigUpdateArgs {
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
	return convertedUpdates
}

func convertUpdateArgsToHybrid(updates []DestChainConfigUpdateArgs) []hybrid_external.FastTransferTokenPoolAbstractDestChainConfigUpdateArgs {
	convertedUpdates := make([]hybrid_external.FastTransferTokenPoolAbstractDestChainConfigUpdateArgs, len(updates))
	for i, update := range updates {
		convertedUpdates[i] = hybrid_external.FastTransferTokenPoolAbstractDestChainConfigUpdateArgs{
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
	return convertedUpdates
}

func convertQuoteFromBurnMint(quote fast_transfer_token_pool.IFastTransferPoolQuote) Quote {
	return Quote{
		CcipSettlementFee: quote.CcipSettlementFee,
		FastTransferFee:   quote.FastTransferFee,
	}
}

func convertQuoteFromHybrid(quote hybrid_external.IFastTransferPoolQuote) Quote {
	return Quote{
		CcipSettlementFee: quote.CcipSettlementFee,
		FastTransferFee:   quote.FastTransferFee,
	}
}

// burnMintPoolAdapter adapts BurnMintFastTransferTokenPool to the FastTransferPool interface
type burnMintPoolAdapter struct {
	pool *fast_transfer_token_pool.BurnMintFastTransferTokenPool
}

// burnMintExternalPoolAdapter adapts BurnMintWithExternalMinterFastTransferTokenPool to the FastTransferPool interface
type burnMintExternalPoolAdapter struct {
	pool *burn_mint_external.BurnMintWithExternalMinterFastTransferTokenPool
}

// hybridExternalPoolAdapter adapts HybridWithExternalMinterFastTransferTokenPool to the FastTransferPool interface
type hybridExternalPoolAdapter struct {
	pool *hybrid_external.HybridWithExternalMinterFastTransferTokenPool
}

// FastTransferTokenPoolWrapper provides a unified interface for
// BurnMintFastTransferTokenPool, BurnMintWithExternalMinterFastTransferTokenPool, and HybridWithExternalMinterFastTransferTokenPool
type FastTransferTokenPoolWrapper struct {
	contractType cldf.ContractType
	address      common.Address
	pool         FastTransferPool
}

// burnMintPoolAdapter only implements methods that need type conversion
func (a *burnMintPoolAdapter) GetDestChainConfig(opts *bind.CallOpts, remoteChainSelector uint64) (DestChainConfig, []common.Address, error) {
	config, addresses, err := a.pool.GetDestChainConfig(opts, remoteChainSelector)
	if err != nil {
		return DestChainConfig{}, nil, err
	}
	return convertDestChainConfigFromBurnMint(config), addresses, nil
}

func (a *burnMintPoolAdapter) UpdateDestChainConfig(opts *bind.TransactOpts, updates []DestChainConfigUpdateArgs) (*types.Transaction, error) {
	return a.pool.UpdateDestChainConfig(opts, convertUpdateArgsToBurnMint(updates))
}

func (a *burnMintPoolAdapter) GetCcipSendTokenFee(opts *bind.CallOpts, destinationChainSelector uint64, amount *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (Quote, error) {
	quote, err := a.pool.GetCcipSendTokenFee(opts, destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)
	if err != nil {
		return Quote{}, err
	}
	return convertQuoteFromBurnMint(quote), nil
}

func (a *burnMintPoolAdapter) FilterFastTransferRequested(opts *bind.FilterOpts, destinationChainSelector []uint64, fillID [][32]byte, settlementID [][32]byte) (*FastTransferRequestedIterator, error) {
	iter, err := a.pool.FilterFastTransferRequested(opts, destinationChainSelector, fillID, settlementID)
	if err != nil {
		return nil, err
	}
	return &FastTransferRequestedIterator{iter: newBurnMintIteratorWrapper(iter)}, nil
}

// Direct delegation methods for burnMintPoolAdapter
func (a *burnMintPoolAdapter) GetAllowedFillers(opts *bind.CallOpts) ([]common.Address, error) {
	return a.pool.GetAllowedFillers(opts)
}
func (a *burnMintPoolAdapter) GetToken(opts *bind.CallOpts) (common.Address, error) {
	return a.pool.GetToken(opts)
}
func (a *burnMintPoolAdapter) UpdateFillerAllowList(opts *bind.TransactOpts, fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error) {
	return a.pool.UpdateFillerAllowList(opts, fillersToAdd, fillersToRemove)
}
func (a *burnMintPoolAdapter) IsAllowedFiller(opts *bind.CallOpts, filler common.Address) (bool, error) {
	return a.pool.IsAllowedFiller(opts, filler)
}
func (a *burnMintPoolAdapter) CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return a.pool.CcipSendToken(opts, destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
}
func (a *burnMintPoolAdapter) GetAccumulatedPoolFees(opts *bind.CallOpts) (*big.Int, error) {
	return a.pool.GetAccumulatedPoolFees(opts)
}
func (a *burnMintPoolAdapter) WithdrawPoolFees(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return a.pool.WithdrawPoolFees(opts, recipient)
}

// burnMintExternalPoolAdapter - no conversion needed, types already match
func (a *burnMintExternalPoolAdapter) GetDestChainConfig(opts *bind.CallOpts, remoteChainSelector uint64) (DestChainConfig, []common.Address, error) {
	return a.pool.GetDestChainConfig(opts, remoteChainSelector)
}
func (a *burnMintExternalPoolAdapter) UpdateDestChainConfig(opts *bind.TransactOpts, updates []DestChainConfigUpdateArgs) (*types.Transaction, error) {
	return a.pool.UpdateDestChainConfig(opts, updates)
}
func (a *burnMintExternalPoolAdapter) GetAllowedFillers(opts *bind.CallOpts) ([]common.Address, error) {
	return a.pool.GetAllowedFillers(opts)
}
func (a *burnMintExternalPoolAdapter) GetToken(opts *bind.CallOpts) (common.Address, error) {
	return a.pool.GetToken(opts)
}
func (a *burnMintExternalPoolAdapter) UpdateFillerAllowList(opts *bind.TransactOpts, fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error) {
	return a.pool.UpdateFillerAllowList(opts, fillersToAdd, fillersToRemove)
}
func (a *burnMintExternalPoolAdapter) IsAllowedFiller(opts *bind.CallOpts, filler common.Address) (bool, error) {
	return a.pool.IsAllowedFiller(opts, filler)
}
func (a *burnMintExternalPoolAdapter) CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return a.pool.CcipSendToken(opts, destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
}
func (a *burnMintExternalPoolAdapter) GetCcipSendTokenFee(opts *bind.CallOpts, destinationChainSelector uint64, amount *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (Quote, error) {
	return a.pool.GetCcipSendTokenFee(opts, destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)
}
func (a *burnMintExternalPoolAdapter) GetAccumulatedPoolFees(opts *bind.CallOpts) (*big.Int, error) {
	return a.pool.GetAccumulatedPoolFees(opts)
}
func (a *burnMintExternalPoolAdapter) WithdrawPoolFees(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return a.pool.WithdrawPoolFees(opts, recipient)
}

func (a *burnMintExternalPoolAdapter) FilterFastTransferRequested(opts *bind.FilterOpts, destinationChainSelector []uint64, fillID [][32]byte, settlementID [][32]byte) (*FastTransferRequestedIterator, error) {
	iter, err := a.pool.FilterFastTransferRequested(opts, destinationChainSelector, fillID, settlementID)
	if err != nil {
		return nil, err
	}
	return &FastTransferRequestedIterator{iter: newBurnMintExternalIteratorWrapper(iter)}, nil
}

// hybridExternalPoolAdapter only implements methods that need type conversion
func (a *hybridExternalPoolAdapter) GetDestChainConfig(opts *bind.CallOpts, remoteChainSelector uint64) (DestChainConfig, []common.Address, error) {
	config, addresses, err := a.pool.GetDestChainConfig(opts, remoteChainSelector)
	if err != nil {
		return DestChainConfig{}, nil, err
	}
	return convertDestChainConfigFromHybrid(config), addresses, nil
}

func (a *hybridExternalPoolAdapter) UpdateDestChainConfig(opts *bind.TransactOpts, updates []DestChainConfigUpdateArgs) (*types.Transaction, error) {
	return a.pool.UpdateDestChainConfig(opts, convertUpdateArgsToHybrid(updates))
}

func (a *hybridExternalPoolAdapter) GetCcipSendTokenFee(opts *bind.CallOpts, destinationChainSelector uint64, amount *big.Int, receiver []byte, settlementFeeToken common.Address, extraArgs []byte) (Quote, error) {
	quote, err := a.pool.GetCcipSendTokenFee(opts, destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)
	if err != nil {
		return Quote{}, err
	}
	return convertQuoteFromHybrid(quote), nil
}

func (a *hybridExternalPoolAdapter) FilterFastTransferRequested(opts *bind.FilterOpts, destinationChainSelector []uint64, fillID [][32]byte, settlementID [][32]byte) (*FastTransferRequestedIterator, error) {
	iter, err := a.pool.FilterFastTransferRequested(opts, destinationChainSelector, fillID, settlementID)
	if err != nil {
		return nil, err
	}
	return &FastTransferRequestedIterator{iter: newHybridExternalIteratorWrapper(iter)}, nil
}

// Direct delegation methods for hybridExternalPoolAdapter
func (a *hybridExternalPoolAdapter) GetAllowedFillers(opts *bind.CallOpts) ([]common.Address, error) {
	return a.pool.GetAllowedFillers(opts)
}
func (a *hybridExternalPoolAdapter) GetToken(opts *bind.CallOpts) (common.Address, error) {
	return a.pool.GetToken(opts)
}
func (a *hybridExternalPoolAdapter) UpdateFillerAllowList(opts *bind.TransactOpts, fillersToAdd []common.Address, fillersToRemove []common.Address) (*types.Transaction, error) {
	return a.pool.UpdateFillerAllowList(opts, fillersToAdd, fillersToRemove)
}
func (a *hybridExternalPoolAdapter) IsAllowedFiller(opts *bind.CallOpts, filler common.Address) (bool, error) {
	return a.pool.IsAllowedFiller(opts, filler)
}
func (a *hybridExternalPoolAdapter) CcipSendToken(opts *bind.TransactOpts, destinationChainSelector uint64, amount *big.Int, maxFastTransferFee *big.Int, receiver []byte, feeToken common.Address, extraArgs []byte) (*types.Transaction, error) {
	return a.pool.CcipSendToken(opts, destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
}
func (a *hybridExternalPoolAdapter) GetAccumulatedPoolFees(opts *bind.CallOpts) (*big.Int, error) {
	return a.pool.GetAccumulatedPoolFees(opts)
}
func (a *hybridExternalPoolAdapter) WithdrawPoolFees(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return a.pool.WithdrawPoolFees(opts, recipient)
}

// AdapterFactory defines the interface for creating pool adapters
type AdapterFactory func(address common.Address, backend bind.ContractBackend) (FastTransferPool, error)

// adapterRegistry maps contract types to their factory functions
var adapterRegistry = map[cldf.ContractType]AdapterFactory{
	shared.BurnMintFastTransferTokenPool: func(address common.Address, backend bind.ContractBackend) (FastTransferPool, error) {
		pool, err := fast_transfer_token_pool.NewBurnMintFastTransferTokenPool(address, backend)
		if err != nil {
			return nil, err
		}
		return &burnMintPoolAdapter{pool: pool}, nil
	},
	shared.BurnMintWithExternalMinterFastTransferTokenPool: func(address common.Address, backend bind.ContractBackend) (FastTransferPool, error) {
		pool, err := burn_mint_external.NewBurnMintWithExternalMinterFastTransferTokenPool(address, backend)
		if err != nil {
			return nil, err
		}
		return &burnMintExternalPoolAdapter{pool: pool}, nil
	},
	shared.HybridWithExternalMinterFastTransferTokenPool: func(address common.Address, backend bind.ContractBackend) (FastTransferPool, error) {
		pool, err := hybrid_external.NewHybridWithExternalMinterFastTransferTokenPool(address, backend)
		if err != nil {
			return nil, err
		}
		return &hybridExternalPoolAdapter{pool: pool}, nil
	},
}

// NewFastTransferTokenPoolWrapper creates a new wrapper instance using the factory pattern
func NewFastTransferTokenPoolWrapper(
	address common.Address,
	backend bind.ContractBackend,
	contractType cldf.ContractType,
) (*FastTransferTokenPoolWrapper, error) {
	factory, exists := adapterRegistry[contractType]
	if !exists {
		return nil, fmt.Errorf("unsupported contract type: %s", contractType)
	}

	pool, err := factory(address, backend)
	if err != nil {
		return nil, err
	}

	return &FastTransferTokenPoolWrapper{
		contractType: contractType,
		address:      address,
		pool:         pool,
	}, nil
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
	return w.pool.GetDestChainConfig(opts, remoteChainSelector)
}

// UpdateDestChainConfig updates destination chain configurations
func (w *FastTransferTokenPoolWrapper) UpdateDestChainConfig(
	opts *bind.TransactOpts,
	updates []DestChainConfigUpdateArgs,
) (*types.Transaction, error) {
	return w.pool.UpdateDestChainConfig(opts, updates)
}

// GetAllowedFillers retrieves the list of allowed filler addresses
func (w *FastTransferTokenPoolWrapper) GetAllowedFillers(opts *bind.CallOpts) ([]common.Address, error) {
	return w.pool.GetAllowedFillers(opts)
}

// GetToken returns the token address associated with the pool
func (w *FastTransferTokenPoolWrapper) GetToken(opts *bind.CallOpts) (common.Address, error) {
	return w.pool.GetToken(opts)
}

// UpdateFillerAllowList updates the filler allowlist
func (w *FastTransferTokenPoolWrapper) UpdateFillerAllowList(
	opts *bind.TransactOpts,
	fillersToAdd []common.Address,
	fillersToRemove []common.Address,
) (*types.Transaction, error) {
	return w.pool.UpdateFillerAllowList(opts, fillersToAdd, fillersToRemove)
}

// IsAllowedFiller checks if an address is an allowed filler
func (w *FastTransferTokenPoolWrapper) IsAllowedFiller(opts *bind.CallOpts, filler common.Address) (bool, error) {
	return w.pool.IsAllowedFiller(opts, filler)
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
	return w.pool.CcipSendToken(opts, destinationChainSelector, amount, maxFastTransferFee, receiver, feeToken, extraArgs)
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
	return w.pool.GetCcipSendTokenFee(opts, destinationChainSelector, amount, receiver, settlementFeeToken, extraArgs)
}

// FilterFastTransferRequested filters FastTransferRequested events (required for e2e tests)
func (w *FastTransferTokenPoolWrapper) FilterFastTransferRequested(
	opts *bind.FilterOpts,
	destinationChainSelector []uint64,
	fillID [][32]byte,
	settlementID [][32]byte,
) (*FastTransferRequestedIterator, error) {
	return w.pool.FilterFastTransferRequested(opts, destinationChainSelector, fillID, settlementID)
}

// eventIterator defines the common interface for all iterator types
type eventIterator interface {
	Next() bool
	Error() error
	Close() error
	GetEvent() *FastTransferRequestedEvent
}

// FastTransferEvent represents the common event structure across different iterator types
type FastTransferEvent interface {
	GetDestinationChainSelector() uint64
	GetFillId() [32]byte
	GetSettlementId() [32]byte
	GetSourceAmountNetFee() *big.Int
	GetSourceDecimals() uint8
	GetFillerFee() *big.Int
	GetPoolFee() *big.Int
	GetDestinationPool() []byte
	GetReceiver() []byte
	GetRaw() types.Log
}

// simpleIteratorWrapper provides a non-generic wrapper that works with function closures
type simpleIteratorWrapper struct {
	nextFunc  func() bool
	errorFunc func() error
	closeFunc func() error
	eventFunc func() *FastTransferRequestedEvent
}

func (s *simpleIteratorWrapper) Next() bool                            { return s.nextFunc() }
func (s *simpleIteratorWrapper) Error() error                          { return s.errorFunc() }
func (s *simpleIteratorWrapper) Close() error                          { return s.closeFunc() }
func (s *simpleIteratorWrapper) GetEvent() *FastTransferRequestedEvent { return s.eventFunc() }

// Factory functions for creating specific iterator wrappers using closures
func newBurnMintIteratorWrapper(iter *fast_transfer_token_pool.BurnMintFastTransferTokenPoolFastTransferRequestedIterator) eventIterator {
	return &simpleIteratorWrapper{
		nextFunc:  func() bool { return iter.Next() },
		errorFunc: func() error { return iter.Error() },
		closeFunc: func() error { return iter.Close() },
		eventFunc: func() *FastTransferRequestedEvent {
			if iter.Event == nil {
				return nil
			}
			return &FastTransferRequestedEvent{
				DestinationChainSelector: iter.Event.DestinationChainSelector,
				FillID:                   iter.Event.FillId,
				SettlementID:             iter.Event.SettlementId,
				SourceAmountNetFee:       iter.Event.SourceAmountNetFee,
				SourceDecimals:           iter.Event.SourceDecimals,
				FillerFee:                iter.Event.FillerFee,
				PoolFee:                  iter.Event.PoolFee,
				DestinationPool:          iter.Event.DestinationPool,
				Receiver:                 iter.Event.Receiver,
				Raw:                      iter.Event.Raw,
			}
		},
	}
}

func newBurnMintExternalIteratorWrapper(iter *burn_mint_external.BurnMintWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator) eventIterator {
	return &simpleIteratorWrapper{
		nextFunc:  func() bool { return iter.Next() },
		errorFunc: func() error { return iter.Error() },
		closeFunc: func() error { return iter.Close() },
		eventFunc: func() *FastTransferRequestedEvent {
			if iter.Event == nil {
				return nil
			}
			return &FastTransferRequestedEvent{
				DestinationChainSelector: iter.Event.DestinationChainSelector,
				FillID:                   iter.Event.FillId,
				SettlementID:             iter.Event.SettlementId,
				SourceAmountNetFee:       iter.Event.SourceAmountNetFee,
				SourceDecimals:           iter.Event.SourceDecimals,
				FillerFee:                iter.Event.FillerFee,
				PoolFee:                  iter.Event.PoolFee,
				DestinationPool:          iter.Event.DestinationPool,
				Receiver:                 iter.Event.Receiver,
				Raw:                      iter.Event.Raw,
			}
		},
	}
}

func newHybridExternalIteratorWrapper(iter *hybrid_external.HybridWithExternalMinterFastTransferTokenPoolFastTransferRequestedIterator) eventIterator {
	return &simpleIteratorWrapper{
		nextFunc:  func() bool { return iter.Next() },
		errorFunc: func() error { return iter.Error() },
		closeFunc: func() error { return iter.Close() },
		eventFunc: func() *FastTransferRequestedEvent {
			if iter.Event == nil {
				return nil
			}
			return &FastTransferRequestedEvent{
				DestinationChainSelector: iter.Event.DestinationChainSelector,
				FillID:                   iter.Event.FillId,
				SettlementID:             iter.Event.SettlementId,
				SourceAmountNetFee:       iter.Event.SourceAmountNetFee,
				SourceDecimals:           iter.Event.SourceDecimals,
				FillerFee:                iter.Event.FillerFee,
				PoolFee:                  iter.Event.PoolFee,
				DestinationPool:          iter.Event.DestinationPool,
				Receiver:                 iter.Event.Receiver,
				Raw:                      iter.Event.Raw,
			}
		},
	}
}

// FastTransferRequestedIterator wraps all iterator types using a common interface
type FastTransferRequestedIterator struct {
	iter eventIterator
}

// Next advances the iterator
func (it *FastTransferRequestedIterator) Next() bool {
	return it.iter.Next()
}

// Error returns the error (if any)
func (it *FastTransferRequestedIterator) Error() error {
	return it.iter.Error()
}

// Close closes the iterator
func (it *FastTransferRequestedIterator) Close() error {
	return it.iter.Close()
}

// Event returns the current event data
func (it *FastTransferRequestedIterator) Event() *FastTransferRequestedEvent {
	return it.iter.GetEvent()
}

// FastTransferRequestedEvent represents a unified event structure
type FastTransferRequestedEvent struct {
	DestinationChainSelector uint64
	FillID                   [32]byte
	SettlementID             [32]byte
	SourceAmountNetFee       *big.Int
	SourceDecimals           uint8
	FillerFee                *big.Int
	PoolFee                  *big.Int
	DestinationPool          []byte
	Receiver                 []byte
	Raw                      types.Log
}

// GetAccumulatedPoolFees retrieves the accumulated pool fees
func (w *FastTransferTokenPoolWrapper) GetAccumulatedPoolFees(opts *bind.CallOpts) (*big.Int, error) {
	return w.pool.GetAccumulatedPoolFees(opts)
}

// WithdrawPoolFees withdraws accumulated pool fees to the specified recipient
func (w *FastTransferTokenPoolWrapper) WithdrawPoolFees(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return w.pool.WithdrawPoolFees(opts, recipient)
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
	case shared.HybridWithExternalMinterFastTransferTokenPool:
		pool, ok := chainState.HybridWithExternalMinterFastTransferTokenPools[tokenSymbol][contractVersion]
		if !ok {
			return nil, fmt.Errorf("hybrid with external minter fast transfer token pool for token %s and version %s not found on chain %s", tokenSymbol, contractVersion, chain.String())
		}
		return NewFastTransferTokenPoolWrapper(pool.Address(), env.BlockChains.EVMChains()[chainSelector].Client, contractType)
	default:
		return nil, fmt.Errorf("unsupported contract type %s for fast transfer token pools", contractType)
	}
}
