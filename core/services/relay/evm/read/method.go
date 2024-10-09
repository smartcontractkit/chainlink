package read

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

type MethodBinding struct {
	// read-only properties
	contractName string
	method       string

	// dependencies
	client               evmclient.Client
	ht                   logpoller.HeadTracker
	lggr                 logger.Logger
	confirmationsMapping map[primitives.ConfidenceLevel]evmtypes.Confirmations

	// internal state properties
	codec    commontypes.Codec
	bindings map[common.Address]struct{}
	mu       sync.RWMutex
}

func NewMethodBinding(
	name, method string,
	client evmclient.Client,
	heads logpoller.HeadTracker,
	confs map[primitives.ConfidenceLevel]evmtypes.Confirmations,
	lggr logger.Logger,
) *MethodBinding {
	return &MethodBinding{
		contractName:         name,
		method:               method,
		client:               client,
		ht:                   heads,
		lggr:                 lggr,
		confirmationsMapping: confs,
		bindings:             make(map[common.Address]struct{}),
	}
}

var _ Reader = &MethodBinding{}

func (b *MethodBinding) Bind(ctx context.Context, bindings ...common.Address) error {
	for _, binding := range bindings {
		if b.isBound(binding) {
			continue
		}

		// check for contract byte code at the latest block and provided address
		byteCode, err := b.client.CodeAt(ctx, binding, nil)
		if err != nil {
			return ErrRead{
				Err: fmt.Errorf("%w: code at call failure: %s", commontypes.ErrInternal, err.Error()),
				Detail: &readDetail{
					Address:  binding.Hex(),
					Contract: b.contractName,
					Params:   nil,
					RetVal:   nil,
				},
			}
		}

		if len(byteCode) == 0 {
			return NoContractExistsError{Err: commontypes.ErrInternal, Address: binding}
		}

		b.setBinding(binding)
	}

	return nil
}

func (b *MethodBinding) Unbind(ctx context.Context, bindings ...common.Address) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, binding := range bindings {
		delete(b.bindings, binding)
	}

	return nil
}

func (b *MethodBinding) SetCodec(codec commontypes.RemoteCodec) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.codec = codec
}

func (b *MethodBinding) BatchCall(address common.Address, params, retVal any) (Call, error) {
	if !b.isBound(address) {
		return Call{}, fmt.Errorf("%w: %w", commontypes.ErrInvalidConfig, newUnboundAddressErr(address.Hex(), b.contractName, b.method))
	}

	return Call{
		ContractAddress: address,
		ContractName:    b.contractName,
		ReadName:        b.method,
		Params:          params,
		ReturnVal:       retVal,
	}, nil
}

func (b *MethodBinding) GetLatestValue(ctx context.Context, addr common.Address, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	if !b.isBound(addr) {
		return fmt.Errorf("%w: %w", commontypes.ErrInvalidConfig, newUnboundAddressErr(addr.Hex(), b.contractName, b.method))
	}

	block, err := b.blockNumberFromConfidence(ctx, confidenceLevel)
	if err != nil {
		return err
	}

	data, err := b.codec.Encode(ctx, params, codec.WrapItemType(b.contractName, b.method, true))
	if err != nil {
		callErr := newErrorFromCall(
			fmt.Errorf("%w: encoding params: %s", commontypes.ErrInvalidType, err.Error()),
			Call{
				ContractAddress: addr,
				ContractName:    b.contractName,
				ReadName:        b.method,
				Params:          params,
				ReturnVal:       returnVal,
			}, block.String(), false)

		return callErr
	}

	callMsg := ethereum.CallMsg{
		To:   &addr,
		From: addr,
		Data: data,
	}

	bytes, err := b.client.CallContract(ctx, callMsg, block)
	if err != nil {
		callErr := newErrorFromCall(
			fmt.Errorf("%w: contract call: %s", commontypes.ErrInvalidType, err.Error()),
			Call{
				ContractAddress: addr,
				ContractName:    b.contractName,
				ReadName:        b.method,
				Params:          params,
				ReturnVal:       returnVal,
			}, block.String(), false)

		return callErr
	}

	if err = b.codec.Decode(ctx, bytes, returnVal, codec.WrapItemType(b.contractName, b.method, false)); err != nil {
		callErr := newErrorFromCall(
			fmt.Errorf("%w: decode return data: %s", commontypes.ErrInvalidType, err.Error()),
			Call{
				ContractAddress: addr,
				ContractName:    b.contractName,
				ReadName:        b.method,
				Params:          params,
				ReturnVal:       returnVal,
			}, block.String(), false)

		strResult := hexutil.Encode(bytes)
		callErr.Result = &strResult

		return callErr
	}

	return nil
}

func (b *MethodBinding) QueryKey(
	_ context.Context,
	_ common.Address,
	_ query.KeyFilter,
	_ query.LimitAndSort,
	_ any,
) ([]commontypes.Sequence, error) {
	return nil, nil
}

func (b *MethodBinding) Register(_ context.Context) error   { return nil }
func (b *MethodBinding) Unregister(_ context.Context) error { return nil }

func (b *MethodBinding) blockNumberFromConfidence(ctx context.Context, confidenceLevel primitives.ConfidenceLevel) (*big.Int, error) {
	confirmations, err := confidenceToConfirmations(b.confirmationsMapping, confidenceLevel)
	if err != nil {
		err = fmt.Errorf("%w: contract: %s; method: %s;", err, b.contractName, b.method)
		if confidenceLevel == primitives.Unconfirmed {
			b.lggr.Debugw("Falling back to default contract call behaviour that calls latest state", "contract", b.contractName, "method", b.method, "err", err)

			return nil, nil
		}

		return nil, err
	}

	_, finalized, err := b.ht.LatestAndFinalizedBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: head tracker: %w", commontypes.ErrInternal, err)
	}

	if confirmations == evmtypes.Finalized {
		return big.NewInt(finalized.Number), nil
	} else if confirmations == evmtypes.Unconfirmed {
		return nil, nil
	}

	return nil, fmt.Errorf("%w: [unknown evm confirmations]: %v; contract: %s; method: %s;", commontypes.ErrInvalidConfig, confirmations, b.contractName, b.method)
}

func (b *MethodBinding) isBound(binding common.Address) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, exists := b.bindings[binding]

	return exists
}

func (b *MethodBinding) setBinding(binding common.Address) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.bindings[binding] = struct{}{}
}
