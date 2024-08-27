package evm

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

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type eventBinding struct {
	address      common.Address
	contractName string
	eventName    string
	lp           logpoller.LogPoller
	// filterRegisterer in eventBinding is to be used as an override for lp filter defined in the contract binding.
	// If filterRegisterer is nil, this event should be registered with the lp filter defined in the contract binding.
	*filterRegisterer
	hash  common.Hash
	codec commontypes.RemoteCodec
	// bound determines if address is set to the contract binding.
	bound    bool
	bindLock sync.Mutex
	// eventTypes has all the types for GetLatestValue unHashed indexed topics params and for QueryKey data words or unHashed indexed topics value comparators.
	eventTypes map[string]types.CodecEntry
	// indexedTopicsCodecInfo has type info about hashed indexed topics.
	indexedTopicsCodecInfo types.CodecEntry
	// eventModifiers only has a modifier for indexed topic filtering, but data words can also be added if needed.
	eventModifiers map[string]codec.Modifier
	// topics map a generic topic name (key) to topic data
	topics map[string]topicDetail
	// dataWordsInfo key is eventName.dataWordName which maps to data word info
	dataWordsInfo        map[string]dataWordInfo
	confirmationsMapping map[primitives.ConfidenceLevel]evmtypes.Confirmations
}

type topicDetail struct {
	abi.Argument
	Index uint64
}

var _ readBinding = &eventBinding{}

func (e *eventBinding) SetCodec(codec commontypes.RemoteCodec) {
	e.codec = codec
}

func (e *eventBinding) Bind(ctx context.Context, binding commontypes.BoundContract) error {
	// it's enough to just lock bound here since Register/Unregister are only called from here and from Start/Close
	// even if they somehow happen at the same time it will be fine because of filter lock and hasFilter check
	e.bindLock.Lock()
	defer e.bindLock.Unlock()

	if e.bound {
		// we are changing contract address reference, so we need to unregister old filter it exists
		if err := e.Unregister(ctx); err != nil {
			return err
		}
	}

	e.address = common.HexToAddress(binding.Address)

	// filterRegisterer isn't required here because the event can also be polled for by the contractBinding common filter.
	if e.filterRegisterer != nil {
		id := fmt.Sprintf("%s.%s.%s", e.contractName, e.eventName, uuid.NewString())
		e.pollingFilter.Name = logpoller.FilterName(id, e.address)
		e.pollingFilter.Addresses = evmtypes.AddressArray{e.address}
		e.bound = true
		if e.registerCalled {
			return e.Register(ctx)
		}
	}
	e.bound = true
	return nil
}

func (e *eventBinding) Register(ctx context.Context) error {
	if e.filterRegisterer == nil {
		return nil
	}

	e.filterLock.Lock()
	defer e.filterLock.Unlock()

	e.registerCalled = true
	// can't be true before filters params are set so there is no race with a bad filter outcome
	if !e.bound {
		return nil
	}

	if e.lp.HasFilter(e.pollingFilter.Name) {
		return nil
	}

	if err := e.lp.RegisterFilter(ctx, e.pollingFilter); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	return nil
}

func (e *eventBinding) Unregister(ctx context.Context) error {
	if e.filterRegisterer == nil {
		return nil
	}

	e.filterLock.Lock()
	defer e.filterLock.Unlock()

	if !e.bound {
		return nil
	}

	if !e.lp.HasFilter(e.pollingFilter.Name) {
		return nil
	}

	if err := e.lp.UnregisterFilter(ctx, e.pollingFilter.Name); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	return nil
}

func (e *eventBinding) GetLatestValue(ctx context.Context, confidenceLevel primitives.ConfidenceLevel, params, into any) error {
	if err := e.validateBound(); err != nil {
		return err
	}

	confirmations, err := confidenceToConfirmations(e.confirmationsMapping, confidenceLevel)
	if err != nil {
		return err
	}

	return e.getLatestValue(ctx, confirmations, params, into)
}

func (e *eventBinding) QueryKey(ctx context.Context, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]commontypes.Sequence, error) {
	if err := e.validateBound(); err != nil {
		return nil, err
	}

	remapped, err := e.remap(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to remap chain agnostic querying filters to chain specific filters: %w", err)
	}

	// filter should always use the address and event sig
	defaultExpressions := []query.Expression{
		logpoller.NewAddressFilter(e.address),
		logpoller.NewEventSigFilter(e.hash),
	}
	remapped.Expressions = append(defaultExpressions, remapped.Expressions...)

	logs, err := e.lp.FilteredLogs(ctx, remapped.Expressions, limitAndSort, e.contractName+"-"+e.address.String()+"-"+e.eventName)
	if err != nil {
		return nil, err
	}

	// no need to return an error. an empty list is fine
	if len(logs) == 0 {
		return []commontypes.Sequence{}, nil
	}

	return e.decodeLogsIntoSequences(ctx, logs, sequenceDataType)
}

func (e *eventBinding) decodeLogsIntoSequences(ctx context.Context, logs []logpoller.Log, into any) ([]commontypes.Sequence, error) {
	sequences := make([]commontypes.Sequence, len(logs))

	for idx := range logs {
		sequences[idx] = commontypes.Sequence{
			Cursor: fmt.Sprintf("%s-%s-%d", logs[idx].BlockHash, logs[idx].TxHash, logs[idx].LogIndex),
			Head: commontypes.Head{
				Identifier: fmt.Sprint(logs[idx].BlockNumber),
				Hash:       logs[idx].BlockHash.Bytes(),
				Timestamp:  uint64(logs[idx].BlockTimestamp.Unix()),
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

		if err := e.decodeLog(ctx, &logs[idx], sequences[idx].Data); err != nil {
			return nil, err
		}
	}
	return sequences, nil
}

func (e *eventBinding) getLatestValue(ctx context.Context, confs evmtypes.Confirmations, params, into any) error {
	checkedValues, err := e.toChecked(WrapItemType(e.contractName, e.eventName, true), params)
	if err != nil {
		return err
	}

	filtersAndIndices, err := e.encodeCheckedTopicFilters(reflect.ValueOf(checkedValues))
	if err != nil {
		return err
	}

	hasFilter := false
	for _, filters := range filtersAndIndices {
		if filters != nil {
			hasFilter = true
		}
	}

	var log *logpoller.Log
	if hasFilter {
		if log, err = e.getLatestLog(ctx, confs, filtersAndIndices); err != nil {
			return err
		}
	} else {
		if log, err = e.lp.LatestLogByEventSigWithConfs(ctx, e.hash, e.address, confs); err != nil {
			return wrapInternalErr(err)
		}
	}

	return e.decodeLog(ctx, log, into)
}

func (e *eventBinding) getLatestLog(ctx context.Context, confs evmtypes.Confirmations, filtersAndIndices []*common.Hash) (*logpoller.Log, error) {
	// Create limiter and filter for the query.
	limiter := query.NewLimitAndSort(query.CountLimit(1), query.NewSortBySequence(query.Desc))
	topicFilters, err := createTopicFilters(filtersAndIndices)
	if err != nil {
		return nil, err
	}

	filter, err := logpoller.Where(
		topicFilters,
		logpoller.NewAddressFilter(e.address),
		logpoller.NewEventSigFilter(e.hash),
		logpoller.NewConfirmationsFilter(confs),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}

	// Gets the latest log that matches the filter and limiter.
	logs, err := e.lp.FilteredLogs(ctx, filter, limiter, e.contractName+"-"+e.address.String()+"-"+e.eventName)
	if err != nil {
		return nil, wrapInternalErr(err)
	}

	if len(logs) == 0 {
		return nil, fmt.Errorf("%w: no events found", commontypes.ErrNotFound)
	}
	return &logs[0], err
}

func createTopicFilters(filtersAndIndices []*common.Hash) (query.Expression, error) {
	var expressions []query.Expression
	for topicID, fai := range filtersAndIndices {
		if fai == nil {
			continue
		}
		// first topic index is 1-based, so we add 1.
		expressions = append(expressions, logpoller.NewEventByTopicFilter(
			uint64(topicID+1), []logpoller.HashedValueComparator{{Value: *fai, Operator: primitives.Eq}},
		))
	}

	if len(expressions) == 0 {
		return query.Expression{}, fmt.Errorf("%w: no topc filters found during query creation", commontypes.ErrInternal)
	}

	return query.And(expressions...), nil
}

// encodeCheckedTopicFilters accepts checked chain types and encodes them to match onchain topics.
func (e *eventBinding) encodeCheckedTopicFilters(checkedTypes reflect.Value) ([]*common.Hash, error) {
	// convert onChain params to native types similarly to generated abi wrappers, for e.g. fixed bytes32 abi type to [32]uint8.
	nativeParams, err := e.eventTypes[e.eventName].ToNative(checkedTypes)
	if err != nil {
		return nil, err
	}

	for nativeParams.Kind() == reflect.Pointer {
		nativeParams = reflect.Indirect(nativeParams)
	}

	var params []any
	switch nativeParams.Kind() {
	case reflect.Array, reflect.Slice:
		native, err := representArray(nativeParams, e.eventTypes[e.eventName])
		if err != nil {
			return nil, err
		}
		params = []any{native}
	case reflect.Struct, reflect.Map:
		if params, err = unrollItem(nativeParams, e.eventTypes[e.eventName]); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w: cannot encode kind %v", commontypes.ErrInvalidType, nativeParams.Kind())
	}

	// abi params allow you to Pack a pointers, but makeTopics doesn't work with pointers.
	return e.makeTopics(derefValues(params))
}

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

// makeTopics encodes and hashes params filtering values to match onchain indexed topics.
func (e *eventBinding) makeTopics(topics []any) ([]*common.Hash, error) {
	var hashableTopics []any
	nonHashableTopics := make(map[int]bool)
	for i, topic := range topics {
		if topic == nil {
			nonHashableTopics[i] = true
			continue
		}

		// make topic value for non-fixed bytes array manually because geth MakeTopics doesn't support it
		if abiArg := e.eventTypes[e.eventName].Args()[i]; abiArg.Type.T == abi.ArrayTy && (abiArg.Type.Elem != nil && abiArg.Type.Elem.T == abi.UintTy) {
			packed, err := abi.Arguments{abiArg}.Pack(topic)
			if err != nil {
				return nil, err
			}
			topics[i] = crypto.Keccak256Hash(packed)
		}

		hashableTopics = append(hashableTopics, topic)
		nonHashableTopics[i] = false
	}

	hashes, err := abi.MakeTopics(hashableTopics)
	if err != nil {
		return nil, wrapInternalErr(err)
	}

	if len(hashes) != 1 {
		return nil, fmt.Errorf("%w: expected 1 filter set, got %d", commontypes.ErrInternal, len(hashes))
	}

	var ptrHashes []*common.Hash
	for i, hash := range hashes[0] {
		if nonHashableTopics[i] {
			ptrHashes = append(ptrHashes, nil)
		}
		ptrHashes = append(ptrHashes, &hash)
	}

	return ptrHashes, nil
}

// derefValues dereferences pointers to nil values to nil.

func (e *eventBinding) decodeLog(ctx context.Context, log *logpoller.Log, into any) error {
	// decode non indexed topics and apply output modifiers
	if err := e.codec.Decode(ctx, log.Data, into, WrapItemType(e.contractName, e.eventName, false)); err != nil {
		return err
	}

	// decode indexed topics which is rarely useful since most indexed topic types get Keccak256 hashed and should be just used for log filtering.
	topics := make([]common.Hash, len(e.indexedTopicsCodecInfo.Args()))
	if len(log.Topics) < len(topics)+1 {
		return fmt.Errorf("%w: not enough topics to decode", commontypes.ErrInvalidType)
	}

	for i := 0; i < len(topics); i++ {
		topics[i] = common.Hash(log.Topics[i+1])
	}

	topicsInto := map[string]any{}
	if err := abi.ParseTopicsIntoMap(topicsInto, e.indexedTopicsCodecInfo.Args(), topics); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
	}

	return mapstructureDecode(topicsInto, into)
}

func (e *eventBinding) remap(filter query.KeyFilter) (remapped query.KeyFilter, err error) {
	for _, expression := range filter.Expressions {
		remappedExpression, err := e.remapExpression(filter.Key, expression)
		if err != nil {
			return query.KeyFilter{}, err
		}

		remapped.Expressions = append(remapped.Expressions, remappedExpression)
	}

	return remapped, nil
}

func (e *eventBinding) remapExpression(key string, expression query.Expression) (query.Expression, error) {
	if !expression.IsPrimitive() {
		remappedBoolExpressions := make([]query.Expression, len(expression.BoolExpression.Expressions))
		for i := range expression.BoolExpression.Expressions {
			remapped, err := e.remapExpression(key, expression.BoolExpression.Expressions[i])
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

	return e.remapPrimitive(expression)
}

// remap chain agnostic primitives to chain specific
func (e *eventBinding) remapPrimitive(expression query.Expression) (query.Expression, error) {
	switch primitive := expression.Primitive.(type) {
	case *primitives.Comparator:
		dwInfo, isDW := e.dataWordsInfo[primitive.Name]
		hashedValComps, err := e.encodeComparator(*primitive, isDW)
		if err != nil {
			return query.Expression{}, fmt.Errorf("failed to encode comparator %q: %w", primitive.Name, err)
		}
		if isDW {
			// TODO fix dw indexing to start from 1, or not
			return logpoller.NewEventByWordFilter(dwInfo.index, hashedValComps), nil
		}
		return logpoller.NewEventByTopicFilter(e.topics[primitive.Name].Index, hashedValComps), nil
	case *primitives.Confidence:
		confirmations, err := confidenceToConfirmations(e.confirmationsMapping, primitive.ConfidenceLevel)
		if err != nil {
			return query.Expression{}, err
		}
		return logpoller.NewConfirmationsFilter(confirmations), nil
	default:
		return expression, nil
	}
}

func (e *eventBinding) encodeComparator(comparator primitives.Comparator, isDW bool) ([]logpoller.HashedValueComparator, error) {
	var hashedValComps []logpoller.HashedValueComparator
	for _, valComp := range comparator.ValueComparators {
		valChecked, err := e.toChecked(WrapItemType(e.contractName, e.eventName+"."+comparator.Name, true), valComp.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert value to checked type: %w", err)
		}

		hashedValComp := logpoller.HashedValueComparator{Operator: valComp.Operator}
		if isDW {
			hashedValComp.Value, err = e.encodeCheckedValComparatorDataWord(comparator.Name, valChecked)
		} else {
			hashedValComp.Value, err = e.encodeCheckedValComparatorTopic(comparator.Name, valChecked)
		}
		if err != nil {
			return nil, err
		}
		hashedValComps = append(hashedValComps, hashedValComp)
	}

	return hashedValComps, nil
}

func (e *eventBinding) encodeCheckedValComparatorDataWord(valCompName string, checkedVal any) (hash common.Hash, err error) {
	defer func() {
		// shouldn't happen, but reflection can panic
		if r := recover(); r != nil {
			err = fmt.Errorf("event %q dataword %q panicked with %v, while trying to encode val %v of type %T", e.eventName, valCompName, r, checkedVal, checkedVal)
		}
	}()

	nativeParams, err := e.eventTypes[valCompName].ToNative(reflect.ValueOf(checkedVal))
	if err != nil {
		return common.Hash{}, err
	}

	for nativeParams.Kind() == reflect.Pointer {
		nativeParams = reflect.Indirect(nativeParams)
	}

	// TODO recalculate dwIndex for fields inside of nested structs, this will only work with static types that don't have a dynamic type before them
	dwInfo, ok := e.dataWordsInfo[valCompName]
	if !ok {
		return common.Hash{}, fmt.Errorf("cannot find data word maping for %s", valCompName)
	}

	packedArgs, err := abi.Arguments{dwInfo.typ}.Pack(nativeParams.Interface())
	if err != nil {
		return common.Hash{}, err
	}

	return common.BytesToHash(packedArgs), nil
}

func (e *eventBinding) encodeCheckedValComparatorTopic(valCompName string, checkedVal any) (hash common.Hash, err error) {
	defer func() {
		// shouldn't happen, but reflection can panic
		if r := recover(); r != nil {
			err = fmt.Errorf("event %q topic %q panicked with %v, while trying to encode val %v of type %T", e.eventName, valCompName, r, checkedVal, checkedVal)
		}
	}()

	// topic encoding is always done on a struct that represents all indexed topics, this makes codec typing same for val comparators and GetLatestValue.
	var filterIndex int
	var topicFields []reflect.StructField
	for i, arg := range e.eventTypes[e.eventName].Args() {
		if arg.Name == e.topics[valCompName].Name {
			filterIndex = i
			topicFields = append(topicFields, reflect.StructField{Name: valCompName, Type: reflect.TypeOf(checkedVal)})
		} else {
			topicFields = append(topicFields, reflect.StructField{Name: arg.Name, Type: reflect.TypeOf(struct{}{})})
		}
	}

	topics := reflect.New(reflect.StructOf(topicFields))
	topics.Elem().FieldByName(valCompName).Set(reflect.ValueOf(checkedVal))

	hashedTopics, err := e.encodeCheckedTopicFilters(topics)
	if err != nil {
		return common.Hash{}, err
	}

	if hashedTopic := hashedTopics[filterIndex]; hashedTopic != nil {
		return *hashedTopic, nil
	}

	return common.Hash{}, fmt.Errorf("event %q topic %q is nil in hashed values %v on index %d", e.eventName, valCompName, hashedTopics, filterIndex)
}

// toChecked injects value into a type that matches onchain types.
func (e *eventBinding) toChecked(itemType string, value any) (any, error) {
	offChain, err := e.codec.CreateType(itemType, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create type %s: %w", itemType, err)
	}

	// apply map struct evm hooks to correct incoming values
	if err = mapstructureDecode(value, offChain); err != nil {
		return nil, err
	}

	// convert caller chain agnostic params types to types representing onchain abi types, for e.g. bytes32. and apply modifiers
	if modifier, exists := e.eventModifiers[itemType]; exists {
		value, err = modifier.TransformToOnChain(offChain, "" /* unused */)
		if err != nil {
			return nil, fmt.Errorf("failed to transform %s to onchain: %w", itemType, err)
		}
	}

	// if no modifiers are present for this type then skip them
	return value, nil
}

func (e *eventBinding) validateBound() error {
	if !e.bound {
		return fmt.Errorf(
			"%w: event %s that belongs to contract: %s, not bound",
			commontypes.ErrInvalidType,
			e.eventName,
			e.contractName,
		)
	}
	return nil
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
