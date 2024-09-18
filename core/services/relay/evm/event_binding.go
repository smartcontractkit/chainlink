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
	bound          bool
	bindLock       sync.Mutex
	inputInfo      types.CodecEntry
	inputModifier  codec.Modifier
	codecTopicInfo types.CodecEntry
	// topics maps a generic topic name (key) to topic data
	topics        map[string]topicDetail
	dataWordsInfo eventDataWords
	// dataWordsMapping maps a generic name to a word index
	// key is a predefined generic name for evm log event data word
	// for e.g. first evm data word(32bytes) of USDC log event is value so the key can be called value
	dataWordsMapping     map[string]uint8
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

	if len(e.inputInfo.Args()) == 0 {
		return e.getLatestValueWithoutFilters(ctx, confirmations, into)
	}

	return e.getLatestValueWithFilters(ctx, confirmations, params, into)
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

func (e *eventBinding) getLatestValueWithoutFilters(ctx context.Context, confs evmtypes.Confirmations, into any) error {
	log, err := e.lp.LatestLogByEventSigWithConfs(ctx, e.hash, e.address, confs)
	if err = wrapInternalErr(err); err != nil {
		return err
	}

	return e.decodeLog(ctx, log, into)
}

func (e *eventBinding) getLatestValueWithFilters(
	ctx context.Context, confs evmtypes.Confirmations, params, into any) error {
	checkedValues, err := e.toChecked(WrapItemType(e.contractName, e.eventName, true), params)
	if err != nil {
		return err
	}

	filtersAndIndices, err := e.encodeParams(reflect.ValueOf(checkedValues))
	if err != nil {
		return err
	}

	// Create limiter and filter for the query.
	limiter := query.NewLimitAndSort(query.CountLimit(1), query.NewSortBySequence(query.Desc))
	filter, err := logpoller.Where(
		logpoller.NewAddressFilter(e.address),
		logpoller.NewEventSigFilter(e.hash),
		logpoller.NewConfirmationsFilter(confs),
		createTopicFilters(filtersAndIndices),
	)
	if err != nil {
		return wrapInternalErr(err)
	}

	// Gets the latest log that matches the filter and limiter.
	logs, err := e.lp.FilteredLogs(ctx, filter, limiter, e.contractName+"-"+e.address.String()+"-"+e.eventName)
	if err != nil {
		return wrapInternalErr(err)
	}

	if len(logs) == 0 {
		return fmt.Errorf("%w: no events found", commontypes.ErrNotFound)
	}

	return e.decodeLog(ctx, &logs[0], into)
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

// toChecked injects value into a type that matches onchain types.
func (e *eventBinding) toChecked(itemType string, value any) (any, error) {
	offChain, err := e.codec.CreateType(itemType, true)
	if err != nil {
		return nil, err
	}

	// apply map struct evm hooks to correct incoming values
	if err = mapstructureDecode(value, offChain); err != nil {
		return nil, err
	}

	// convert caller chain agnostic params types to types representing onchain abi types, for e.g. bytes32. and apply modifiers
	return e.inputModifier.TransformToOnChain(offChain, "" /* unused */)
}

// TODO unspaghetti
// encodeParams accepts chain types and encodes them to match onchain topics.
func (e *eventBinding) encodeParams(checkedTypes reflect.Value) ([]common.Hash, error) {
	// convert onChain params to native types similarly to generated abi wrappers, for e.g. fixed bytes32 abi type to [32]uint8.
	nativeParams, err := e.inputInfo.ToNative(checkedTypes)
	if err != nil {
		return nil, err
	}

	for nativeParams.Kind() == reflect.Pointer {
		nativeParams = reflect.Indirect(nativeParams)
	}

	var values []any
	switch nativeParams.Kind() {
	case reflect.Array, reflect.Slice:
		native, err := representArray(nativeParams, e.inputInfo)
		if err != nil {
			return nil, err
		}
		values = []any{native}
	case reflect.Struct, reflect.Map:
		if values, err = unrollItem(nativeParams, e.inputInfo); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w: cannot encode kind %v", commontypes.ErrInvalidType, nativeParams.Kind())
	}

	// abi params allow you to Pack a pointers, but makeTopics doesn't work with pointers.
	// TODO extract this from here for remaps and deref manually because of error referencing topics by indexes
	if err = e.derefTopics(values); err != nil {
		return nil, err
	}

	return e.makeTopics(values)
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

// makeTopics encodes and hashes params filtering values to match onchain indexed topics.
func (e *eventBinding) makeTopics(params []any) ([]common.Hash, error) {
	// make topic value for non-fixed bytes array manually because geth MakeTopics doesn't support it
	for i, topic := range params {
		// TODO if you didn't add input info in config, but have a topic for QueryKey, this will panic, which is not good
		if abiArg := e.inputInfo.Args()[i]; abiArg.Type.T == abi.ArrayTy && (abiArg.Type.Elem != nil && abiArg.Type.Elem.T == abi.UintTy) {
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

func (e *eventBinding) decodeLog(ctx context.Context, log *logpoller.Log, into any) error {
	// decode non indexed topics and apply output modifiers
	if err := e.codec.Decode(ctx, log.Data, into, WrapItemType(e.contractName, e.eventName, false)); err != nil {
		return err
	}

	// decode indexed topics which is rarely useful since most indexed topic types get Keccak256 hashed and should be just used for log filtering.
	topics := make([]common.Hash, len(e.codecTopicInfo.Args()))
	if len(log.Topics) < len(topics)+1 {
		return fmt.Errorf("%w: not enough topics to decode", commontypes.ErrInvalidType)
	}

	for i := 0; i < len(topics); i++ {
		topics[i] = common.Hash(log.Topics[i+1])
	}

	topicsInto := map[string]any{}
	if err := abi.ParseTopicsIntoMap(topicsInto, e.codecTopicInfo.Args(), topics); err != nil {
		return fmt.Errorf("%w: %w", commontypes.ErrInvalidType, err)
	}

	return mapstructureDecode(topicsInto, into)
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
		var hashedValComps []logpoller.HashedValueComparator
		dwIndex, isDW := e.dataWordsMapping[primitive.Name]
		for _, valComp := range primitive.ValueComparators {
			valChecked, err := e.toChecked(WrapItemType(e.contractName, e.eventName+"."+primitive.Name, true), valComp.Value)
			if err != nil {
				return query.Expression{}, err
			}

			hashedValComp := logpoller.HashedValueComparator{Operator: valComp.Operator}
			if isDW {
				// TODO
				return query.Expression{}, nil
			} else {
				wrappedVal := reflect.New(reflect.StructOf([]reflect.StructField{{Name: primitive.Name, Type: reflect.TypeOf(valChecked)}}))
				wrappedVal.Elem().FieldByName(primitive.Name).Set(reflect.ValueOf(valChecked))
				topics, err := e.encodeParams(wrappedVal)
				if err != nil {
					return query.Expression{}, err
				}
				hashedValComp.Value = topics[0]
			}
			hashedValComps = append(hashedValComps, hashedValComp)
		}
		if isDW {
			return logpoller.NewEventByWordFilter(dwIndex, hashedValComps), nil
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

// TODO
func (e *eventBinding) encodeDataWord(nativeValue reflect.Value, primitive primitives.Comparator) (common.Hash, error) {
	// TODO recalculate dwIndex if value that we are searching for is after a dynamic type like string, tuple, array...
	// this will only work with static types that don't have a dynamic type before them
	dwIndex, ok := e.dataWordsMapping[primitive.Name]
	if !ok {
		return common.Hash{}, fmt.Errorf("cannot find data word maping for %s", primitive.Name)
	}

	if len(e.dataWordsInfo) <= int(dwIndex) {
		return common.Hash{}, fmt.Errorf("data word index is out of bounds %d for data word  %s", dwIndex, primitive.Name)
	}

	packedArgs, err := abi.Arguments{e.dataWordsInfo[dwIndex].Argument}.Pack(nativeValue)
	if err != nil {
		return common.Hash{}, err
	}

	return common.BytesToHash(packedArgs), nil
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
