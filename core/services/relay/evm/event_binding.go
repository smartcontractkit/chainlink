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
	// indexedTopicsTypes has type info about hashed indexed topics.
	indexedTopicsTypes types.CodecEntry
	// eventModifiers only has a modifier for indexed topic filtering, but data words can also be added if needed.
	eventModifiers map[string]codec.Modifier
	// topics map a generic topic name (key) to topic data
	topics map[string]topicDetail
	// dataWords key is the generic dataWordName.
	dataWords            map[string]dataWordDetail
	confirmationsMapping map[primitives.ConfidenceLevel]evmtypes.Confirmations
	// TODO add lggr and detailed logs for all internal errors BCI-4130
}

type topicDetail struct {
	abi.Argument
	Index uint64
}

// dataWordDetail contains all the information about a single evm Data word.
// For e.g. first evm data word(32bytes) of USDC log event is uint256 var called value.
type dataWordDetail struct {
	index uint8
	abi.Argument
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

	confs, err := confidenceToConfirmations(e.confirmationsMapping, confidenceLevel)
	if err != nil {
		return err
	}

	topicTypeID := WrapItemType(e.contractName, e.eventName, true)

	onChainTypedVal, err := e.toNativeOnChainType(topicTypeID, params)
	if err != nil {
		return fmt.Errorf("failed to convert params to native on chain types: %w", err)
	}

	filterTopics, err := e.extractFilterTopics(topicTypeID, onChainTypedVal)
	if err != nil {
		return err
	}

	var log *logpoller.Log
	if len(filterTopics) != 0 {
		hashedTopics, err := e.hashTopics(topicTypeID, filterTopics)
		if err != nil {
			return err
		}

		if log, err = e.getLatestLog(ctx, confs, hashedTopics); err != nil {
			return err
		}
	} else {
		if log, err = e.lp.LatestLogByEventSigWithConfs(ctx, e.hash, e.address, confs); err != nil {
			return wrapInternalErr(err)
		}
	}

	return e.decodeLog(ctx, log, into)
}

func (e *eventBinding) QueryKey(ctx context.Context, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]commontypes.Sequence, error) {
	if err := e.validateBound(); err != nil {
		return nil, err
	}

	remapped, err := e.remap(filter)
	if err != nil {
		return nil, err
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

func (e *eventBinding) getLatestLog(ctx context.Context, confs evmtypes.Confirmations, hashedTopics []common.Hash) (*logpoller.Log, error) {
	// Create limiter and filter for the query.
	limiter := query.NewLimitAndSort(query.CountLimit(1), query.NewSortBySequence(query.Desc))
	topicFilters, err := createTopicFilters(hashedTopics)
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
		return nil, wrapInternalErr(err)
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

// extractFilterTopics extracts filter topics from input params and returns them as a slice of any.
// returned slice will retain the order of the topics and fill in missing topics with nil, if all values are nil, empty slice is returned.
func (e *eventBinding) extractFilterTopics(topicTypeID string, value any) (filterTopics []any, err error) {
	item := reflect.ValueOf(value)
	switch item.Kind() {
	case reflect.Array, reflect.Slice:
		native, err := representArray(item, e.eventTypes[topicTypeID])
		if err != nil {
			return nil, err
		}
		filterTopics = []any{native}
	case reflect.Struct, reflect.Map:
		if filterTopics, err = unrollItem(item, e.eventTypes[topicTypeID]); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w: cannot encode kind %v", commontypes.ErrInvalidType, item.Kind())
	}

	noValues := true
	for _, topic := range filterTopics {
		if topic != nil {
			noValues = false
		}
	}
	if noValues {
		return []any{}, nil
	}
	return filterTopics, nil
}

// hashTopics hashes topic filters values to match on chain indexed topics.
func (e *eventBinding) hashTopics(topicTypeID string, topics []any) ([]common.Hash, error) {
	var hashableTopics []any
	for i, topic := range derefValues(topics) {
		if topic == nil {
			continue
		}

		// make topic value for non-fixed bytes array manually because geth MakeTopics doesn't support it
		topicTyp, exists := e.eventTypes[topicTypeID]
		if !exists {
			return nil, fmt.Errorf("cannot find event type entry")
		}

		if abiArg := topicTyp.Args()[i]; abiArg.Type.T == abi.ArrayTy && (abiArg.Type.Elem != nil && abiArg.Type.Elem.T == abi.UintTy) {
			packed, err := abi.Arguments{abiArg}.Pack(topic)
			if err != nil {
				return nil, err
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

func (e *eventBinding) decodeLog(ctx context.Context, log *logpoller.Log, into any) error {
	// decode non indexed topics and apply output modifiers
	if err := e.codec.Decode(ctx, log.Data, into, WrapItemType(e.contractName, e.eventName, false)); err != nil {
		return err
	}

	// decode indexed topics which is rarely useful since most indexed topic types get Keccak256 hashed and should be just used for log filtering.
	topics := make([]common.Hash, len(e.indexedTopicsTypes.Args()))
	if len(log.Topics) < len(topics)+1 {
		return fmt.Errorf("%w: not enough topics to decode", commontypes.ErrInvalidType)
	}

	for i := 0; i < len(topics); i++ {
		topics[i] = common.Hash(log.Topics[i+1])
	}

	topicsInto := map[string]any{}
	if err := abi.ParseTopicsIntoMap(topicsInto, e.indexedTopicsTypes.Args(), topics); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
	}

	return mapstructureDecode(topicsInto, into)
}

// remap chain agnostic primitives to chain specific logPoller primitives.
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

func (e *eventBinding) remapPrimitive(expression query.Expression) (query.Expression, error) {
	switch primitive := expression.Primitive.(type) {
	case *primitives.Comparator:
		hashedValComps, err := e.encodeComparator(primitive)
		if err != nil {
			return query.Expression{}, fmt.Errorf("failed to encode comparator %q: %w", primitive.Name, err)
		}
		return hashedValComps, nil
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

func (e *eventBinding) encodeComparator(comparator *primitives.Comparator) (query.Expression, error) {
	dwInfo, isDW := e.dataWords[comparator.Name]
	if !isDW {
		if _, exists := e.topics[comparator.Name]; !exists {
			return query.Expression{}, fmt.Errorf("comparator name doesn't match any of the indexed topics or data words")
		}
	}

	var hashedValComps []logpoller.HashedValueComparator
	itemType := WrapItemType(e.contractName, e.eventName+"."+comparator.Name, true)
	for _, valComp := range comparator.ValueComparators {
		onChainTypedVal, err := e.toNativeOnChainType(itemType, valComp.Value)
		if err != nil {
			return query.Expression{}, fmt.Errorf("failed to convert comparator value to native on chain type: %w", err)
		}

		hashedValComp := logpoller.HashedValueComparator{Operator: valComp.Operator}
		if isDW {
			hashedValComp.Value, err = e.encodeValComparatorDataWord(itemType, onChainTypedVal)
		} else {
			hashedValComp.Value, err = e.encodeValComparatorTopic(itemType, onChainTypedVal)
		}
		if err != nil {
			return query.Expression{}, err
		}
		hashedValComps = append(hashedValComps, hashedValComp)
	}

	if isDW {
		return logpoller.NewEventByWordFilter(dwInfo.index, hashedValComps), nil
	}

	return logpoller.NewEventByTopicFilter(e.topics[comparator.Name].Index, hashedValComps), nil
}

func (e *eventBinding) encodeValComparatorDataWord(dwTypeID string, value any) (hash common.Hash, err error) {
	dwTypes, exists := e.eventTypes[dwTypeID]
	if !exists {
		return common.Hash{}, fmt.Errorf("cannot find data word type for %s", dwTypeID)
	}

	packedArgs, err := dwTypes.Args().Pack(value)
	if err != nil {
		return common.Hash{}, err
	}

	return common.BytesToHash(packedArgs), nil
}

func (e *eventBinding) encodeValComparatorTopic(topicTypeID string, value any) (hash common.Hash, err error) {
	hashedTopics, err := e.hashTopics(topicTypeID, []any{value})
	if err != nil {
		return common.Hash{}, err
	}

	return hashedTopics[0], nil
}

// toNativeOnChainType converts value into its on chain version by applying codec modifiers, map structure hooks and abi typing.
func (e *eventBinding) toNativeOnChainType(itemType string, value any) (any, error) {
	offChain, err := e.codec.CreateType(itemType, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create type: %w", err)
	}

	// apply map struct evm hooks to correct incoming values
	if err = mapstructureDecode(value, offChain); err != nil {
		return nil, err
	}

	// apply modifiers if present
	onChain := offChain
	if modifier, exists := e.eventModifiers[itemType]; exists {
		onChain, err = modifier.TransformToOnChain(offChain, "" /* unused */)
		if err != nil {
			return nil, fmt.Errorf("failed to apply modifiers to offchain type %T: %w", onChain, err)
		}
	}

	typ, exists := e.eventTypes[itemType]
	if !exists {
		return query.Expression{}, fmt.Errorf("cannot find event type entry")
	}

	native, err := typ.ToNative(reflect.ValueOf(onChain))
	if err != nil {
		return query.Expression{}, err
	}

	for native.Kind() == reflect.Pointer {
		native = reflect.Indirect(native)
	}

	return native.Interface(), nil
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

func createTopicFilters(hashedTopics []common.Hash) (query.Expression, error) {
	var expressions []query.Expression
	for topicID, hash := range hashedTopics {
		// first topic index is 1-based, so we add 1.
		expressions = append(expressions, logpoller.NewEventByTopicFilter(
			uint64(topicID+1), []logpoller.HashedValueComparator{{Value: hash, Operator: primitives.Eq}},
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
	return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
}
