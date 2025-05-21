package read

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
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

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type EventBinding struct {
	// read-only properties
	contractName string
	eventName    string
	hash         common.Hash
	// eventTypes has all the types for GetLatestValue unHashed indexed topics params and for QueryKey data words or unHashed indexed topics value comparators.
	eventTypes map[string]types.CodecEntry
	// indexedTopicsTypes has type info about hashed indexed topics.
	indexedTopicsTypes types.CodecEntry
	// eventModifiers only has a modifier for indexed topic filtering, but data words can also be added if needed.
	eventModifiers map[string]commoncodec.Modifier

	// dependencies
	// filterRegistrar in EventBinding is to be used as an override for lp filter defined in the contract binding.
	// If filterRegisterer is nil, this event should be registered with the lp filter defined in the contract binding.
	registrar      *syncedFilter
	registerCalled bool
	lp             logpoller.LogPoller

	// internal properties / state
	codec commontypes.RemoteCodec
	bound map[common.Address]bool // bound determines if address is set to the contract binding.
	mu    sync.RWMutex
	// topics map a generic topic name (key) to topic data
	topics map[string]TopicDetail
	// dataWords key is the generic dataWordNamb.
	dataWords            map[string]DataWordDetail
	confirmationsMapping map[primitives.ConfidenceLevel]evmtypes.Confirmations
}

func NewEventBinding(
	contract, event string,
	poller logpoller.LogPoller,
	hash common.Hash,
	indexedTopicsTypes types.CodecEntry,
	confirmations map[primitives.ConfidenceLevel]evmtypes.Confirmations,
) *EventBinding {
	return &EventBinding{
		contractName:         contract,
		eventName:            event,
		lp:                   poller,
		hash:                 hash,
		indexedTopicsTypes:   indexedTopicsTypes,
		confirmationsMapping: confirmations,
		topics:               make(map[string]TopicDetail),
		dataWords:            make(map[string]DataWordDetail),
		bound:                make(map[common.Address]bool),
	}
}

type TopicDetail struct {
	abi.Argument
	Index uint64
}

// DataWordDetail contains all the information about a single evm Data word.
// For e.g. first evm data word(32bytes) of USDC log event is uint256 var called value.
type DataWordDetail struct {
	Index int
	abi.Argument
}

var _ Reader = &EventBinding{}

func (b *EventBinding) SetCodec(codec commontypes.RemoteCodec) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.codec = codec
}

func (b *EventBinding) SetFilter(filter logpoller.Filter) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.registrar = newSyncedFilter()
	b.registrar.SetFilter(filter)
}

func (b *EventBinding) SetCodecTypesAndModifiers(types map[string]types.CodecEntry, modifiers map[string]commoncodec.Modifier) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.eventTypes = types
	b.eventModifiers = modifiers
}

func (b *EventBinding) SetDataWordsDetails(dwDetail map[string]DataWordDetail) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.dataWords = dwDetail
}

func (b *EventBinding) SetTopicDetails(topicDetails map[string]TopicDetail) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.topics = topicDetails
}

func (b *EventBinding) GetDataWords() map[string]DataWordDetail {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.dataWords
}

func (b *EventBinding) Bind(ctx context.Context, bindings ...common.Address) error {
	for _, binding := range bindings {
		if b.isBound(binding) {
			continue
		}

		if b.registrar != nil {
			b.registrar.AddAddress(binding)
		}

		b.addBinding(binding)

	}

	if b.registrar == nil || !b.registrar.Dirty() {
		return nil
	}

	if b.registered() {
		return b.Update(ctx)
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

func (b *EventBinding) Update(ctx context.Context) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	name := logpoller.FilterName(fmt.Sprintf("%s.%s.%s", b.contractName, b.eventName, uuid.NewString()))

	if b.registrar == nil {
		return nil
	}

	if len(b.bound) == 0 {
		return nil
	}

	return b.registrar.Update(ctx, b.lp, name)
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

func (b *EventBinding) GetLatestValueWithHeadData(ctx context.Context, address common.Address, confidenceLevel primitives.ConfidenceLevel, params, into any) (head *commontypes.Head, err error) {
	var (
		confs  evmtypes.Confirmations
		result *string
	)

	defer func() {
		if err != nil {
			callErr := newErrorFromCall(err, Call{
				ContractAddress: address,
				ContractName:    b.contractName,
				ReadName:        b.eventName,
				Params:          params,
				ReturnVal:       into,
			}, strconv.Itoa(int(confs)), eventReadType)

			callErr.Result = result

			err = callErr
		}
	}()

	if err = b.validateBound(address); err != nil {
		return nil, err
	}

	confs, err = confidenceToConfirmations(b.confirmationsMapping, confidenceLevel)
	if err != nil {
		return nil, err
	}

	topicTypeID := codec.WrapItemType(b.contractName, b.eventName, true)

	onChainTypedVal, err := b.toNativeOnChainType(topicTypeID, params)
	if err != nil {
		return nil, err
	}

	filterTopics, err := b.extractFilterTopics(topicTypeID, onChainTypedVal)
	if err != nil {
		return nil, err
	}

	var log *logpoller.Log
	if len(filterTopics) != 0 {
		var hashedTopics []common.Hash
		hashedTopics, err = b.hashTopics(topicTypeID, filterTopics)
		if err != nil {
			return nil, err
		}

		if log, err = b.getLatestLog(ctx, address, confs, hashedTopics); err != nil {
			return nil, err
		}
	} else {
		if log, err = b.lp.LatestLogByEventSigWithConfs(ctx, b.hash, address, confs); err != nil {
			return nil, wrapInternalErr(err)
		}
	}

	if err = b.decodeLog(ctx, log, into); err != nil {
		encoded := hex.EncodeToString(log.Data)
		result = &encoded
		return nil, err
	}

	return &commontypes.Head{
		Height: strconv.FormatInt(log.BlockNumber, 10),
		Hash:   log.BlockHash.Bytes(),
		//nolint:gosec // G115
		Timestamp: uint64(log.BlockTimestamp.Unix()),
	}, nil
}

func (b *EventBinding) QueryKey(ctx context.Context, address common.Address, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) (sequences []commontypes.Sequence, err error) {
	defer func() {
		if err != nil {
			err = newErrorFromCall(err, Call{
				ContractAddress: address,
				ContractName:    b.contractName,
				ReadName:        b.eventName,
				ReturnVal:       sequenceDataType,
			}, "", eventReadType)
		}
	}()

	if err = b.validateBound(address); err != nil {
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

	logs, err := b.lp.FilteredLogs(ctx, remapped.Expressions, limitAndSort, b.contractName+"-"+address.String()+"-"+b.eventName)
	if err != nil {
		return nil, wrapInternalErr(err)
	}

	// no need to return an error. an empty list is fine
	if len(logs) == 0 {
		return []commontypes.Sequence{}, nil
	}

	sequences, err = b.decodeLogsIntoSequences(ctx, logs, sequenceDataType)
	if err != nil {
		return nil, err
	}

	return sequences, nil
}

func (b *EventBinding) getLatestLog(ctx context.Context, address common.Address, confs evmtypes.Confirmations, hashedTopics []common.Hash) (*logpoller.Log, error) {
	// Create limiter and filter for the query.
	limiter := query.NewLimitAndSort(query.CountLimit(1), query.NewSortBySequence(query.Desc))

	topicFilters, err := createTopicFilters(hashedTopics)
	if err != nil {
		return nil, err
	}

	filter, err := logpoller.Where(
		topicFilters,
		logpoller.NewAddressFilter(address),
		logpoller.NewEventSigFilter(b.hash),
		logpoller.NewConfirmationsFilter(confs),
	)
	if err != nil {
		return nil, wrapInternalErr(err)
	}

	// Gets the latest log that matches the filter and limiter.
	logs, err := b.lp.FilteredLogs(ctx, filter, limiter, b.contractName+"-"+address.String()+"-"+b.eventName)
	if err != nil {
		return nil, wrapInternalErr(err)
	}

	if len(logs) == 0 {
		return nil, fmt.Errorf("%w: no events found", commontypes.ErrNotFound)
	}

	return &logs[0], err
}

func (b *EventBinding) decodeLogsIntoSequences(ctx context.Context, logs []logpoller.Log, into any) ([]commontypes.Sequence, error) {
	sequences := make([]commontypes.Sequence, len(logs))

	for idx := range logs {
		sequences[idx] = commontypes.Sequence{
			Cursor: logpoller.FormatContractReaderCursor(logs[idx]),
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

// extractFilterTopics extracts filter topics from input params and returns them as a slice of any.
// returned slice will retain the order of the topics and fill in missing topics with nil, if all values are nil, empty slice is returned.
func (b *EventBinding) extractFilterTopics(topicTypeID string, value any) (filterTopics []any, err error) {
	item := reflect.ValueOf(value)

	switch item.Kind() {
	case reflect.Array, reflect.Slice:
		var native any
		native, err = codec.RepresentArray(item, b.eventTypes[topicTypeID])
		if err != nil {
			return nil, fmt.Errorf("%w: error converting params to log topics: %s", commontypes.ErrInternal, err.Error())
		}

		filterTopics = []any{native}
	case reflect.Struct, reflect.Map:
		if filterTopics, err = codec.UnrollItem(item, b.eventTypes[topicTypeID]); err != nil {
			return nil, fmt.Errorf("%w: error unrolling params into log topics: %s", commontypes.ErrInternal, err.Error())
		}
	default:
		return nil, fmt.Errorf("%w: cannot encode kind %v", commontypes.ErrInvalidType, item.Kind())
	}

	// check if at least one topic filter is present
	for _, filterVal := range derefValues(filterTopics) {
		if filterVal != nil {
			return filterTopics, nil
		}
	}

	return []any{}, nil
}

// hashTopics hashes topic filters values to match on chain indexed topics.
func (b *EventBinding) hashTopics(topicTypeID string, topics []any) ([]common.Hash, error) {
	var hashableTopics []any
	for i, topic := range derefValues(topics) {
		if topic == nil {
			continue
		}

		// make topic value for non-fixed bytes array manually because geth MakeTopics doesn't support it
		topicTyp, exists := b.eventTypes[topicTypeID]
		if !exists {
			return nil, fmt.Errorf("%w: cannot find event type entry for topic: %s", commontypes.ErrNotFound, topicTypeID)
		}

		if abiArg := topicTyp.Args()[i]; abiArg.Type.T == abi.ArrayTy && (abiArg.Type.Elem != nil && abiArg.Type.Elem.T == abi.UintTy) {
			packed, err := abi.Arguments{abiArg}.Pack(topic)
			if err != nil {
				return nil, fmt.Errorf("%w: err failed to abi pack topics: %s", commontypes.ErrInternal, err.Error())
			}

			topic = crypto.Keccak256Hash(packed)
		}

		hashableTopics = append(hashableTopics, topic)
	}

	hashes, err := abi.MakeTopics(hashableTopics)
	if err != nil {
		return nil, wrapInternalErr(err)
	}

	if len(hashes) != 1 {
		return nil, fmt.Errorf("%w: expected 1 filter set, got %d", commontypes.ErrInternal, len(hashes))
	}

	return hashes[0], nil
}

func (b *EventBinding) decodeLog(ctx context.Context, log *logpoller.Log, into any) error {
	// decode non indexed topics and apply output modifiers
	if err := b.codec.Decode(ctx, log.Data, into, codec.WrapItemType(b.contractName, b.eventName, false)); err != nil {
		return fmt.Errorf("%w: failed to decode log data: %s", commontypes.ErrInvalidType, err.Error())
	}

	// decode indexed topics which is rarely useful since most indexed topic types get Keccak256 hashed and should be just used for log filtering.
	topics := make([]common.Hash, len(b.indexedTopicsTypes.Args()))
	if len(log.Topics) < len(topics)+1 {
		return fmt.Errorf("%w: not enough topics to decode", commontypes.ErrInvalidType)
	}

	for i := 0; i < len(topics); i++ {
		topics[i] = common.Hash(log.Topics[i+1])
	}

	topicsInto := map[string]any{}
	if err := abi.ParseTopicsIntoMap(topicsInto, b.indexedTopicsTypes.Args(), topics); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
	}

	if err := codec.MapstructureDecode(topicsInto, into); err != nil {
		return fmt.Errorf("%w: failed to decode log data: %s", commontypes.ErrInvalidEncoding, err.Error())
	}

	return nil
}

// remap chain agnostic primitives to chain specific logPoller primitives.
func (b *EventBinding) remap(filter query.KeyFilter) (remapped query.KeyFilter, err error) {
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

	return b.remapPrimitive(expression)
}

func (b *EventBinding) remapPrimitive(expression query.Expression) (query.Expression, error) {
	switch primitive := expression.Primitive.(type) {
	case *primitives.Comparator:
		hashedValComps, err := b.encodeComparator(primitive)
		if err != nil {
			return query.Expression{}, fmt.Errorf("failed to encode comparator %q: %w", primitive.Name, err)
		}
		return hashedValComps, nil
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

func (b *EventBinding) valueCmpToHashedCmp(itemType string, isDW bool, valueCmp primitives.ValueComparator) (logpoller.HashedValueComparator, error) {
	var valuesToEncode []any
	switch typedValue := valueCmp.Value.(type) {
	case primitives.AnyOperator:
		valuesToEncode = typedValue
	default:
		valuesToEncode = []any{typedValue}
	}

	result := logpoller.HashedValueComparator{
		Operator: valueCmp.Operator,
		Values:   make([]common.Hash, len(valuesToEncode)),
	}

	for i := range valuesToEncode {
		valueToEncode := valuesToEncode[i]
		encodedValue, err := b.encodeValue(itemType, isDW, valueToEncode)
		if err != nil {
			return logpoller.HashedValueComparator{}, fmt.Errorf("failed to encode value %v: %w", valueToEncode, err)
		}
		result.Values[i] = encodedValue
	}

	return result, nil
}

func (b *EventBinding) encodeValue(itemType string, isDW bool, value any) (common.Hash, error) {
	onChainTypedVal, err := b.toNativeOnChainType(itemType, value)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to convert comparator value to native on chain type: %w", err)
	}

	if isDW {
		return b.encodeValComparatorDataWord(itemType, onChainTypedVal)
	}

	return b.encodeValComparatorTopic(itemType, onChainTypedVal)
}

func (b *EventBinding) encodeComparator(comparator *primitives.Comparator) (query.Expression, error) {
	dwInfo, isDW := b.dataWords[comparator.Name]
	if !isDW {
		if _, exists := b.topics[comparator.Name]; !exists {
			return query.Expression{}, fmt.Errorf("%w: comparator name doesn't match any of the indexed topics or data words", commontypes.ErrInvalidConfig)
		}
	}

	var hashedValComps []logpoller.HashedValueComparator
	itemType := codec.WrapItemType(b.contractName, b.eventName+"."+comparator.Name, true)
	for _, valComp := range comparator.ValueComparators {
		hashedValComp, err := b.valueCmpToHashedCmp(itemType, isDW, valComp)
		if err != nil {
			return query.Expression{}, err
		}
		hashedValComps = append(hashedValComps, hashedValComp)
	}

	if isDW {
		return logpoller.NewEventByWordFilter(dwInfo.Index, hashedValComps), nil
	}

	return logpoller.NewEventByTopicFilter(b.topics[comparator.Name].Index, hashedValComps), nil
}

func (b *EventBinding) encodeValComparatorDataWord(dwTypeID string, value any) (hash common.Hash, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: cannot encode %s data word comparator. Recovered from panic: %v", commontypes.ErrInvalidType, dwTypeID, r)
		}
	}()

	dwTypes, exists := b.eventTypes[dwTypeID]
	if !exists {
		return common.Hash{}, fmt.Errorf("%w: cannot find data word type for %s", commontypes.ErrInvalidConfig, dwTypeID)
	}

	packedArgs, err := dwTypes.Args().Pack(value)
	if err != nil {
		return common.Hash{}, fmt.Errorf("%w: failed to pack values: %w", commontypes.ErrInternal, err)
	}

	return common.BytesToHash(packedArgs), nil
}

func (b *EventBinding) encodeValComparatorTopic(topicTypeID string, value any) (hash common.Hash, err error) {
	hashedTopics, err := b.hashTopics(topicTypeID, []any{value})
	if err != nil {
		return common.Hash{}, err
	}

	return hashedTopics[0], nil
}

// toNativeOnChainType converts value into its on chain version by applying codec modifiers, map structure hooks and abi typing.
func (b *EventBinding) toNativeOnChainType(itemType string, value any) (any, error) {
	offChain, err := b.codec.CreateType(itemType, true)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create type: %w", commontypes.ErrInvalidType, err)
	}

	// apply map struct evm hooks to correct incoming values
	if err = codec.MapstructureDecode(value, offChain); err != nil {
		return nil, fmt.Errorf("%w: failed to decode offChain value: %s", commontypes.ErrInternal, err.Error())
	}

	// apply modifiers if present
	onChain := offChain
	if modifier, exists := b.eventModifiers[itemType]; exists {
		onChain, err = modifier.TransformToOnChain(offChain, "" /* unused */)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to apply modifiers to offchain type %T: %w", commontypes.ErrInvalidType, onChain, err)
		}
	}

	typ, exists := b.eventTypes[itemType]
	if !exists {
		return query.Expression{}, fmt.Errorf("%w: cannot find event type entry for %s", commontypes.ErrInvalidType, itemType)
	}

	native, err := typ.ToNative(reflect.ValueOf(onChain))
	if err != nil {
		return query.Expression{}, fmt.Errorf("%w: codec to native: %s", commontypes.ErrInvalidType, err.Error())
	}

	return native.Interface(), nil
}

func (b *EventBinding) validateBound(address common.Address) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	bound, exists := b.bound[address]
	if !exists || !bound {
		return fmt.Errorf("%w: %w", commontypes.ErrInvalidConfig, newUnboundAddressErr(address.String(), b.contractName, b.eventName))
	}

	return nil
}

func createTopicFilters(hashedTopics []common.Hash) (query.Expression, error) {
	var expressions []query.Expression
	for topicID, hash := range hashedTopics {
		expressions = append(expressions, logpoller.NewEventByTopicFilter(
			// adding 1 to skip even signature
			uint64(topicID+1), //nolint:gosec // G115
			[]logpoller.HashedValueComparator{{Values: []common.Hash{hash}, Operator: primitives.Eq}},
		))
	}

	if len(expressions) == 0 {
		return query.Expression{}, fmt.Errorf("%w: no topic filters found during query creation", commontypes.ErrInternal)
	}

	return query.And(expressions...), nil
}

// derefValues dereferences pointers to nil values to nil.
func derefValues(topics []any) []any {
	for i, topic := range topics {
		rTopic := reflect.ValueOf(topic)
		if rTopic.Kind() == reflect.Pointer {
			if rTopic.IsNil() {
				topics[i] = nil
			} else {
				topics[i] = rTopic.Elem().Interface()
			}
		}
	}

	return topics
}

func wrapInternalErr(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()
	if strings.Contains(errStr, "not found") || strings.Contains(errStr, "no rows") {
		return fmt.Errorf("%w: %w", commontypes.ErrNotFound, err)
	}

	return fmt.Errorf("%w: %s", commontypes.ErrInternal, err.Error())
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
