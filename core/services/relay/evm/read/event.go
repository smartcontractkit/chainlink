package read

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	commoncodec "github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type EventBinding struct {
	// read-only properties
	contractName   string
	eventName      string
	hash           common.Hash
	inputInfo      types.CodecEntry
	inputModifier  commoncodec.Modifier
	codecTopicInfo types.CodecEntry

	// dependencies
	// filterRegistrar in EventBinding is to be used as an override for lp filter defined in the contract binding.
	// If filterRegisterer is nil, this event should be registered with the lp filter defined in the contract binding.
	registrar      *syncedFilter
	registerCalled bool
	lp             logpoller.LogPoller

	// internal properties / state
	codec  commontypes.RemoteCodec
	bound  map[common.Address]bool // bound determines if address is set to the contract binding.
	mu     sync.RWMutex
	topics map[string]topicDetail // topics maps a generic topic name (key) to topic data
	// eventDataWords maps a generic name to a word index
	// key is a predefined generic name for evm log event data word
	// for e.g. first evm data word(32bytes) of USDC log event is value so the key can be called value
	eventDataWords       map[string]uint8
	confirmationsMapping map[primitives.ConfidenceLevel]evmtypes.Confirmations
}

func NewEventBinding(
	contract, event string,
	poller logpoller.LogPoller,
	hash common.Hash,
	inputInfo types.CodecEntry,
	inputModifier commoncodec.Modifier,
	codecTopicInfo types.CodecEntry,
	confirmations map[primitives.ConfidenceLevel]evmtypes.Confirmations,
) *EventBinding {
	return &EventBinding{
		contractName:         contract,
		eventName:            event,
		lp:                   poller,
		hash:                 hash,
		inputInfo:            inputInfo,
		inputModifier:        inputModifier,
		codecTopicInfo:       codecTopicInfo,
		confirmationsMapping: confirmations,
		topics:               make(map[string]topicDetail),
		eventDataWords:       make(map[string]uint8),
		bound:                make(map[common.Address]bool),
	}
}

type topicDetail struct {
	abi.Argument
	Index uint64
}

var _ Reader = &EventBinding{}

func (b *EventBinding) SetCodec(codec commontypes.RemoteCodec) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.codec = codec
}

func (b *EventBinding) Bind(ctx context.Context, bindings ...common.Address) error {
	if b.hasBindings() {
		// we are changing contract address reference, so we need to unregister old filter if it exists
		if err := b.Unregister(ctx); err != nil {
			return err
		}
	}

	// filterRegisterer isn't required here because the event can also be polled for by the contractBinding common filter.
	if b.registrar != nil {
		b.registrar.SetName(logpoller.FilterName(fmt.Sprintf("%s.%s.%s", b.contractName, b.eventName, uuid.NewString())))
	}

	for _, binding := range bindings {
		if b.isBound(binding) {
			continue
		}

		if b.registrar != nil {
			b.registrar.AddAddress(binding)
		}

		b.addBinding(binding)
	}

	if b.registered() {
		return b.Register(ctx)
	}

	return nil
}

func (b *EventBinding) Unbind(ctx context.Context, bindings ...common.Address) error {
	for _, binding := range bindings {
		if !b.isBound(binding) {
			continue
		}

		if b.registrar != nil {
			b.registrar.RemoveAddress(binding)
		}

		b.removeBinding(binding)
	}

	if err := b.Unregister(ctx); err != nil {
		return err
	}

	// we are changing contract address reference, so we need to unregister old filter or re-register existing filter
	if b.registrar != nil {
		if b.registrar.Count() == 0 {
			b.registrar.SetName("")

			return b.Unregister(ctx)
		} else if b.registered() {
			return b.Register(ctx)
		}
	}

	return nil
}

func (b *EventBinding) Register(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.registrar == nil {
		return nil
	}

	b.registerCalled = true

	// can't be true before filters params are set so there is no race with a bad filter outcome
	if len(b.bound) == 0 {
		return nil
	}

	return b.registrar.Register(ctx, b.lp)
}

func (b *EventBinding) Unregister(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.registrar == nil {
		return nil
	}

	if len(b.bound) == 0 {
		return nil
	}

	return b.registrar.Unregister(ctx, b.lp)
}

func (b *EventBinding) BatchCall(_ common.Address, _, _ any) (Call, error) {
	return Call{}, fmt.Errorf("%w: events are not yet supported in batch get latest values", commontypes.ErrInvalidType)
}

func (b *EventBinding) GetLatestValue(ctx context.Context, address common.Address, confidenceLevel primitives.ConfidenceLevel, params, into any) error {
	if err := b.validateBound(address); err != nil {
		return err
	}

	confirmations, err := confidenceToConfirmations(b.confirmationsMapping, confidenceLevel)
	if err != nil {
		return err
	}

	if len(b.inputInfo.Args()) == 0 {
		return b.getLatestValueWithoutFilters(ctx, address, confirmations, into)
	}

	return b.getLatestValueWithFilters(ctx, address, confirmations, params, into)
}

func (b *EventBinding) QueryKey(
	ctx context.Context,
	address common.Address,
	filter query.KeyFilter,
	limitAndSort query.LimitAndSort,
	sequenceDataType any,
) ([]commontypes.Sequence, error) {
	if err := b.validateBound(address); err != nil {
		return nil, err
	}

	remapped, err := b.remap(filter)
	if err != nil {
		return nil, err
	}

	// filter should always use the address and event sig
	defaultExpressions := []query.Expression{
		logpoller.NewAddressFilter(address),
		logpoller.NewEventSigFilter(b.hash),
	}

	remapped.Expressions = append(defaultExpressions, remapped.Expressions...)

	logs, err := b.lp.FilteredLogs(ctx, remapped.Expressions, limitAndSort, b.contractName+"-"+address.String()+""+b.eventName)
	if err != nil {
		return nil, err
	}

	// no need to return an error. an empty list is fine
	if len(logs) == 0 {
		return []commontypes.Sequence{}, nil
	}

	return b.decodeLogsIntoSequences(ctx, logs, sequenceDataType)
}

func (b *EventBinding) SetFilter(filter logpoller.Filter) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.registrar = newSyncedFilter()
	b.registrar.SetFilter(filter)
}

func (b *EventBinding) WithTopic(name string, topic abi.Argument, index uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.topics[name] = topicDetail{
		Argument: topic,
		Index:    index,
	}
}

func (b *EventBinding) SetDataWords(eventDataWords map[string]uint8) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.eventDataWords = eventDataWords
}

func (b *EventBinding) GetDataWords() map[string]uint8 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.eventDataWords
}

func (b *EventBinding) validateBound(address common.Address) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	bound, exists := b.bound[address]
	if !exists || !bound {
		return fmt.Errorf(
			"%w: event %s that belongs to contract: %s, not bound",
			commontypes.ErrInvalidType,
			b.eventName,
			b.contractName,
		)
	}

	return nil
}

func (b *EventBinding) getLatestValueWithoutFilters(
	ctx context.Context,
	address common.Address,
	confs evmtypes.Confirmations,
	into any,
) error {
	log, err := b.lp.LatestLogByEventSigWithConfs(ctx, b.hash, address, confs)
	if err = wrapInternalErr(err); err != nil {
		return err
	}

	return b.decodeLog(ctx, log, into)
}

func (b *EventBinding) getLatestValueWithFilters(
	ctx context.Context,
	address common.Address,
	confs evmtypes.Confirmations,
	params, into any,
) error {
	offChain, err := b.convertToOffChainType(params)
	if err != nil {
		return err
	}

	checkedParams, err := b.inputModifier.TransformToOnChain(offChain, "" /* unused */)
	if err != nil {
		return err
	}

	nativeParams, err := b.inputInfo.ToNative(reflect.ValueOf(checkedParams))
	if err != nil {
		return err
	}

	filtersAndIndices, err := b.encodeParams(nativeParams)
	if err != nil {
		return err
	}

	remainingFilters := filtersAndIndices[1:]

	// Create limiter and filter for the query.
	limiter := query.NewLimitAndSort(query.CountLimit(1), query.NewSortBySequence(query.Desc))
	filter, err := query.Where(
		"",
		logpoller.NewAddressFilter(address),
		logpoller.NewEventSigFilter(b.hash),
		logpoller.NewConfirmationsFilter(confs),
		createTopicFilters(filtersAndIndices),
	)
	if err != nil {
		return wrapInternalErr(err)
	}

	// Gets the latest log that matches the filter and limiter.
	logs, err := b.lp.FilteredLogs(ctx, filter.Expressions, limiter, b.contractName+"-"+address.String()+"-"+b.eventName)
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

	return b.decodeLog(ctx, logToUse, into)
}

func (b *EventBinding) convertToOffChainType(params any) (any, error) {
	offChain, err := b.codec.CreateType(WrapItemType(b.contractName, b.eventName, true), true)
	if err != nil {
		return nil, err
	}

	if err = codec.MapstructureDecode(params, offChain); err != nil {
		return nil, err
	}

	return offChain, nil
}

func (b *EventBinding) encodeParams(item reflect.Value) ([]common.Hash, error) {
	for item.Kind() == reflect.Pointer {
		item = reflect.Indirect(item)
	}

	var params []any
	switch item.Kind() {
	case reflect.Array, reflect.Slice:
		native, err := codec.RepresentArray(item, b.inputInfo)
		if err != nil {
			return nil, err
		}
		params = []any{native}
	case reflect.Struct, reflect.Map:
		var err error
		if params, err = codec.UnrollItem(item, b.inputInfo); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w: cannot encode kind %v", commontypes.ErrInvalidType, item.Kind())
	}

	// abi params allow you to Pack a pointers, but MakeTopics doesn't work with pointers.
	if err := b.derefTopics(params); err != nil {
		return nil, err
	}

	return b.makeTopics(params)
}

func (b *EventBinding) derefTopics(topics []any) error {
	for i, topic := range topics {
		rTopic := reflect.ValueOf(topic)
		if rTopic.Kind() == reflect.Pointer {
			if rTopic.IsNil() {
				return fmt.Errorf(
					"%w: input topic %s cannot be nil", commontypes.ErrInvalidType, b.inputInfo.Args()[i].Name)
			}

			topics[i] = rTopic.Elem().Interface()
		}
	}

	return nil
}

// makeTopics encodes and hashes params filtering values to match onchain indexed topics.
func (b *EventBinding) makeTopics(params []any) ([]common.Hash, error) {
	// make topic value for non-fixed bytes array manually because geth MakeTopics doesn't support it
	for i, topic := range params {
		if abiArg := b.inputInfo.Args()[i]; abiArg.Type.T == abi.ArrayTy && (abiArg.Type.Elem != nil && abiArg.Type.Elem.T == abi.UintTy) {
			packed, err := abi.Arguments{abiArg}.Pack(topic)
			if err != nil {
				return nil, err
			}
			params[i] = crypto.Keccak256Hash(packed)
		}
	}

	hashes, err := abi.MakeTopics(params)
	if err != nil {
		return nil, wrapInternalErr(err)
	}

	if len(hashes) != 1 {
		return nil, fmt.Errorf("%w: expected 1 filter set, got %d", commontypes.ErrInternal, len(hashes))
	}

	return hashes[0], nil
}

func (b *EventBinding) decodeLog(ctx context.Context, log *logpoller.Log, into any) error {
	if err := b.codec.Decode(ctx, log.Data, into, WrapItemType(b.contractName, b.eventName, false)); err != nil {
		return err
	}

	topics := make([]common.Hash, len(b.codecTopicInfo.Args()))
	if len(log.Topics) < len(topics)+1 {
		return fmt.Errorf("%w: not enough topics to decode", commontypes.ErrInvalidType)
	}

	for i := 0; i < len(topics); i++ {
		topics[i] = common.Hash(log.Topics[i+1])
	}

	topicsInto := map[string]any{}
	if err := abi.ParseTopicsIntoMap(topicsInto, b.codecTopicInfo.Args(), topics); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
	}

	return codec.MapstructureDecode(topicsInto, into)
}

func (b *EventBinding) decodeLogsIntoSequences(ctx context.Context, logs []logpoller.Log, into any) ([]commontypes.Sequence, error) {
	sequences := make([]commontypes.Sequence, len(logs))

	for idx := range logs {
		sequences[idx] = commontypes.Sequence{
			Cursor: fmt.Sprintf("%s-%s-%d", logs[idx].BlockHash, logs[idx].TxHash, logs[idx].LogIndex),
			Head: commontypes.Head{
				Height:    fmt.Sprint(logs[idx].BlockNumber),
				Hash:      logs[idx].BlockHash.Bytes(),
				Timestamp: uint64(logs[idx].BlockTimestamp.Unix()),
			},
		}

		var typeVal reflect.Value

		typeInto := reflect.TypeOf(into)
		if typeInto.Kind() == reflect.Pointer {
			typeVal = reflect.New(typeInto.Elem())
		} else {
			typeVal = reflect.Indirect(reflect.New(typeInto))
		}

		// create a new value of the same type as 'into' for the data to be extracted to
		sequences[idx].Data = typeVal.Interface()

		if err := b.decodeLog(ctx, &logs[idx], sequences[idx].Data); err != nil {
			return nil, err
		}
	}

	return sequences, nil
}

func (b *EventBinding) remap(filter query.KeyFilter) (query.KeyFilter, error) {
	remapped := query.KeyFilter{}

	for _, expression := range filter.Expressions {
		remappedExpression, err := b.remapExpression(filter.Key, expression)
		if err != nil {
			return query.KeyFilter{}, err
		}

		remapped.Expressions = append(remapped.Expressions, remappedExpression)
	}

	return remapped, nil
}

func (b *EventBinding) remapExpression(key string, expression query.Expression) (query.Expression, error) {
	if !expression.IsPrimitive() {
		remappedBoolExpressions := make([]query.Expression, len(expression.BoolExpression.Expressions))

		for i := range expression.BoolExpression.Expressions {
			remapped, err := b.remapExpression(key, expression.BoolExpression.Expressions[i])
			if err != nil {
				return query.Expression{}, err
			}

			remappedBoolExpressions[i] = remapped
		}

		if expression.BoolExpression.BoolOperator == query.AND {
			return query.And(remappedBoolExpressions...), nil
		}

		return query.Or(remappedBoolExpressions...), nil
	}

	return b.remapPrimitive(key, expression)
}

// remap chain agnostic primitives to chain specific
func (b *EventBinding) remapPrimitive(key string, expression query.Expression) (query.Expression, error) {
	switch primitive := expression.Primitive.(type) {
	case *primitives.Comparator:
		if val, ok := b.eventDataWords[primitive.Name]; ok {
			return logpoller.NewEventByWordFilter(b.hash, val, primitive.ValueComparators), nil
		}

		return logpoller.NewEventByTopicFilter(b.topics[key].Index, primitive.ValueComparators), nil
	case *primitives.Confidence:
		confirmations, err := confidenceToConfirmations(b.confirmationsMapping, primitive.ConfidenceLevel)
		if err != nil {
			return query.Expression{}, err
		}

		return logpoller.NewConfirmationsFilter(confirmations), nil
	default:
		return expression, nil
	}
}

func (b *EventBinding) hasBindings() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.bound) > 0
}

func (b *EventBinding) isBound(binding common.Address) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	_, exists := b.bound[binding]

	return exists
}

func (b *EventBinding) addBinding(binding common.Address) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.bound[binding] = true
}

func (b *EventBinding) removeBinding(binding common.Address) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.bound, binding)
}

func (b *EventBinding) registered() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.registerCalled
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

func createTopicFilters(filtersAndIndices []common.Hash) query.Expression {
	var expressions []query.Expression
	for topicID, fai := range filtersAndIndices {
		// first topic index is 1-based, so we add 1.
		expressions = append(expressions, logpoller.NewEventByTopicFilter(
			uint64(topicID+1), []primitives.ValueComparator{{Value: fai.Hex(), Operator: primitives.Eq}},
		))
	}
	return query.And(expressions...)
}
