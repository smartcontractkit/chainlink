package evm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type eventBinding struct {
	address        common.Address
	contractName   string
	eventName      string
	lp             logpoller.LogPoller
	hash           common.Hash
	codec          commontypes.RemoteCodec
	pending        bool
	bound          bool
	registerCalled bool
	lock           sync.Mutex
	inputInfo      types.CodecEntry
	inputModifier  codec.Modifier
	topicInfo      types.CodecEntry
	// used to allow Register and Unregister to be unique in case two bindings have the same event.
	// otherwise, if one unregisters, it'll unregister both with the LogPoller.
	id string
}

var _ readBinding = &eventBinding{}

func (e *eventBinding) SetCodec(codec commontypes.RemoteCodec) {
	e.codec = codec
}

func (e *eventBinding) Register() error {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.registerCalled = true
	if !e.bound || e.lp.HasFilter(e.id) {
		return nil
	}

	if err := e.lp.RegisterFilter(logpoller.Filter{
		Name:      e.id,
		EventSigs: evmtypes.HashArray{e.hash},
		Addresses: evmtypes.AddressArray{e.address},
	}); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}
	return nil
}

func (e *eventBinding) Unregister() error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if !e.lp.HasFilter(e.id) {
		return nil
	}

	if err := e.lp.UnregisterFilter(e.id); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}
	return nil
}

func (e *eventBinding) GetLatestValue(ctx context.Context, params, into any) error {
	if !e.bound {
		return fmt.Errorf("%w: event not bound", commontypes.ErrInvalidType)
	}

	confs := logpoller.Finalized
	if e.pending {
		confs = logpoller.Unconfirmed
	}

	if len(e.inputInfo.Args()) == 0 {
		return e.getLatestValueWithoutFilters(ctx, confs, into)
	}

	return e.getLatestValueWithFilters(ctx, confs, params, into)
}

func (e *eventBinding) Bind(binding commontypes.BoundContract) error {
	if err := e.Unregister(); err != nil {
		return err
	}

	e.address = common.HexToAddress(binding.Address)
	e.pending = binding.Pending
	e.bound = true

	if e.registerCalled {
		return e.Register()
	}
	return nil
}

func (e *eventBinding) getLatestValueWithoutFilters(ctx context.Context, confs logpoller.Confirmations, into any) error {
	log, err := e.lp.LatestLogByEventSigWithConfs(e.hash, e.address, confs)
	if err = wrapInternalErr(err); err != nil {
		return err
	}

	return e.decodeLog(ctx, log, into)
}

func (e *eventBinding) getLatestValueWithFilters(
	ctx context.Context, confs logpoller.Confirmations, params, into any) error {
	offChain, err := e.convertToOffChainType(params)
	if err != nil {
		return err
	}

	checkedParams, err := e.inputModifier.TransformToOnChain(offChain, "" /* unused */)
	if err != nil {
		return err
	}

	nativeParams, err := e.inputInfo.ToNative(reflect.ValueOf(checkedParams))
	if err != nil {
		return err
	}

	filtersAndIndices, err := e.encodeParams(nativeParams)
	if err != nil {
		return err
	}

	fai := filtersAndIndices[0]
	remainingFilters := filtersAndIndices[1:]

	logs, err := e.lp.IndexedLogs(e.hash, e.address, 1, []common.Hash{fai}, confs)
	if err != nil {
		return wrapInternalErr(err)
	}

	// TODO: there should be a better way to ask log poller to filter these
	// First, you should be able to ask for as many topics to match
	// Second, you should be able to get the latest only
	var logToUse *logpoller.Log
	for _, log := range logs {
		tmp := log
		if compareLogs(&tmp, logToUse) > 0 && matchesRemainingFilters(&tmp, remainingFilters) {
			// copy so that it's not pointing to the changing variable
			logToUse = &tmp
		}
	}

	if logToUse == nil {
		return fmt.Errorf("%w: no events found", commontypes.ErrNotFound)
	}

	return e.decodeLog(ctx, logToUse, into)
}

func (e *eventBinding) convertToOffChainType(params any) (any, error) {
	itemType := wrapItemType(e.contractName, e.eventName, true)
	offChain, err := e.codec.CreateType(itemType, true)
	if err != nil {
		return nil, err
	}

	if err = mapstructureDecode(params, offChain); err != nil {
		return nil, err
	}

	return offChain, nil
}

func compareLogs(log, use *logpoller.Log) int64 {
	if use == nil {
		return 1
	}

	if log.BlockNumber != use.BlockNumber {
		return log.BlockNumber - use.BlockNumber
	}

	return log.LogIndex - use.LogIndex
}

func matchesRemainingFilters(log *logpoller.Log, filters []common.Hash) bool {
	for i, rfai := range filters {
		if !reflect.DeepEqual(rfai[:], log.Topics[i+2]) {
			return false
		}
	}

	return true
}

func (e *eventBinding) encodeParams(item reflect.Value) ([]common.Hash, error) {
	for item.Kind() == reflect.Pointer {
		item = reflect.Indirect(item)
	}

	var topics []any
	switch item.Kind() {
	case reflect.Array, reflect.Slice:
		native, err := representArray(item, e.inputInfo)
		if err != nil {
			return nil, err
		}
		topics = []any{native}
	case reflect.Struct, reflect.Map:
		var err error
		if topics, err = unrollItem(item, e.inputInfo); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w: cannot encode kind %v", commontypes.ErrInvalidType, item.Kind())
	}

	// abi params allow you to Pack a pointers, but MakeTopics doesn't work with pointers.
	if err := e.derefTopics(topics); err != nil {
		return nil, err
	}

	hashes, err := abi.MakeTopics(topics)
	if err != nil {
		return nil, wrapInternalErr(err)
	}

	if len(hashes) != 1 {
		return nil, fmt.Errorf("%w: expected 1 filter set, got %d", commontypes.ErrInternal, len(hashes))
	}

	return hashes[0], nil
}

func (e *eventBinding) derefTopics(topics []any) error {
	for i, topic := range topics {
		rTopic := reflect.ValueOf(topic)
		if rTopic.Kind() == reflect.Pointer {
			if rTopic.IsNil() {
				return fmt.Errorf(
					"%w: input topic %s cannot be nil", commontypes.ErrInvalidType, e.inputInfo.Args()[i].Name)
			}
			topics[i] = rTopic.Elem().Interface()
		}
	}
	return nil
}

func (e *eventBinding) decodeLog(ctx context.Context, log *logpoller.Log, into any) error {
	dataType := wrapItemType(e.contractName, e.eventName, false)
	if err := e.codec.Decode(ctx, log.Data, into, dataType); err != nil {
		return err
	}

	topics := make([]common.Hash, len(e.topicInfo.Args()))
	if len(log.Topics) < len(topics)+1 {
		return fmt.Errorf("%w: not enough topics to decode", commontypes.ErrInvalidType)
	}

	for i := 0; i < len(topics); i++ {
		topics[i] = common.Hash(log.Topics[i+1])
	}

	topicsInto := map[string]any{}
	if err := abi.ParseTopicsIntoMap(topicsInto, e.topicInfo.Args(), topics); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
	}

	return mapstructureDecode(topicsInto, into)
}

func wrapInternalErr(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()
	if strings.Contains(errStr, "not found") || strings.Contains(errStr, "no rows") {
		return fmt.Errorf("%w: %w", commontypes.ErrNotFound, err)
	}
	return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
}
