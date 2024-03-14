package evm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"

	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type ChainReaderService interface {
	services.ServiceCtx
	commontypes.ChainReader
}

type chainReader struct {
	lggr             logger.Logger
	lp               logpoller.LogPoller
	client           evmclient.Client
	contractBindings contractBindings
	// TODO should keys that are bound just to event sigs also live here? Probably yes
	eventIndexBindings EventIndexBindings
	parsed             *parsedTypes
	codec              commontypes.RemoteCodec
	commonservices.StateMachine
}

var _ ChainReaderService = (*chainReader)(nil)

// NewChainReaderService is a constructor for ChainReader, returns nil if there is any error
func NewChainReaderService(ctx context.Context, lggr logger.Logger, lp logpoller.LogPoller, chain legacyevm.Chain, config types.ChainReaderConfig) (ChainReaderService, error) {
	cr := &chainReader{
		lggr:               lggr.Named("ChainReader"),
		lp:                 lp,
		client:             chain.Client(),
		contractBindings:   contractBindings{},
		eventIndexBindings: EventIndexBindings{},
		parsed:             &parsedTypes{encoderDefs: map[string]types.CodecEntry{}, decoderDefs: map[string]types.CodecEntry{}},
	}

	var err error
	if err = cr.init(config.Contracts); err != nil {
		return nil, err
	}

	if cr.codec, err = cr.parsed.toCodec(); err != nil {
		return nil, err
	}

	err = cr.contractBindings.ForEach(ctx, func(b readBinding, c context.Context) error {
		b.SetCodec(cr.codec)
		return nil
	})

	return cr, err
}

func (cr *chainReader) Name() string { return cr.lggr.Name() }

var _ commontypes.ContractTypeProvider = &chainReader{}

func (cr *chainReader) GetLatestValue(ctx context.Context, contractName, method string, params any, returnVal any) error {
	b, err := cr.contractBindings.GetReadBinding(contractName, method)
	if err != nil {
		return err
	}

	return b.GetLatestValue(ctx, params, returnVal)
}

func (cr *chainReader) Bind(ctx context.Context, bindings []commontypes.BoundContract) error {
	return cr.contractBindings.Bind(ctx, bindings)
}

func (cr *chainReader) QueryKey(ctx context.Context, key string, queryFilters []commontypes.QueryFilter, limitAndSort commontypes.LimitAndSort) ([]commontypes.Sequence, error) {
	logFilters, err := remapQueryKeyFilters(key, cr.eventIndexBindings, queryFilters)
	if err != nil {
		return nil, err
	}
	_, err = cr.lp.FilteredLogs(logFilters, limitAndSort)
	return nil, err
}

// TODO if slice of keys then matrix of queryFilters?
func (cr *chainReader) QueryKeys(ctx context.Context, keys []string, queryFilters []commontypes.QueryFilter, limitAndSort commontypes.LimitAndSort) ([][]commontypes.Sequence, error) {
	var logFilters []commontypes.QueryFilter
	for _, key := range keys {
		filters, err := remapQueryKeyFilters(key, cr.eventIndexBindings, queryFilters)
		if err != nil {
			return nil, err
		}
		logFilters = append(logFilters, filters...)
	}
	_, err := cr.lp.FilteredLogs(logFilters, limitAndSort)
	return nil, err
}

func (cr *chainReader) QueryKeyByValues(ctx context.Context, key string, values []string, queryFilters []commontypes.QueryFilter, limitAndSort commontypes.LimitAndSort) ([]commontypes.Sequence, error) {
	logFilters, err := remapQueryKeyByValuesFilters(key, values, cr.eventIndexBindings, queryFilters)
	if err != nil {
		return nil, err
	}
	_, err = cr.lp.FilteredLogs(logFilters, limitAndSort)
	return nil, err
}

// TODO if slice of keys then matrix of queryFilters?
// TODO values shouldn't be string?
func (cr *chainReader) QueryKeysByValues(ctx context.Context, keys []string, values [][]string, queryFilters []commontypes.QueryFilter, limitAndSort commontypes.LimitAndSort) ([][]commontypes.Sequence, error) {
	var logFilters []commontypes.QueryFilter
	for i, key := range keys {
		// TODO should this be an error? how to handle.
		if len(values) < i {
			return nil, fmt.Errorf("")
		}
		filters, err := remapQueryKeyByValuesFilters(key, values[i], cr.eventIndexBindings, queryFilters)
		if err != nil {
			return nil, err
		}
		logFilters = append(logFilters, filters...)
	}
	_, err := cr.lp.FilteredLogs(logFilters, limitAndSort)
	return nil, err
}

func (cr *chainReader) init(chainContractReaders map[string]types.ChainContractReader) error {
	for contractName, chainContractReader := range chainContractReaders {
		contractAbi, err := abi.JSON(strings.NewReader(chainContractReader.ContractABI))
		if err != nil {
			return err
		}

		for typeName, chainReaderDefinition := range chainContractReader.Configs {
			switch chainReaderDefinition.ReadType {
			case types.Method:
				err = cr.addMethod(contractName, typeName, contractAbi, *chainReaderDefinition)
			case types.Event:
				err = cr.addEvent(contractName, typeName, contractAbi, *chainReaderDefinition)
			default:
				return fmt.Errorf(
					"%w: invalid chain reader definition read type: %s",
					commontypes.ErrInvalidConfig,
					chainReaderDefinition.ReadType)
			}

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (cr *chainReader) Start(ctx context.Context) error {
	return cr.StartOnce("ChainReader", func() error {
		return cr.contractBindings.ForEach(ctx, readBinding.Register)
	})
}

func (cr *chainReader) Close() error {
	return cr.StopOnce("ChainReader", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return cr.contractBindings.ForEach(ctx, readBinding.Unregister)
	})
}

func (cr *chainReader) Ready() error { return nil }
func (cr *chainReader) HealthReport() map[string]error {
	return map[string]error{cr.Name(): nil}
}

func (cr *chainReader) CreateContractType(contractName, methodName string, forEncoding bool) (any, error) {
	return cr.codec.CreateType(wrapItemType(contractName, methodName, forEncoding), forEncoding)
}

func wrapItemType(contractName, methodName string, isParams bool) string {
	if isParams {
		return fmt.Sprintf("params.%s.%s", contractName, methodName)
	}
	return fmt.Sprintf("return.%s.%s", contractName, methodName)
}

func (cr *chainReader) addMethod(
	contractName,
	methodName string,
	abi abi.ABI,
	chainReaderDefinition types.ChainReaderDefinition) error {
	method, methodExists := abi.Methods[chainReaderDefinition.ChainSpecificName]
	if !methodExists {
		return fmt.Errorf("%w: method %s doesn't exist", commontypes.ErrInvalidConfig, chainReaderDefinition.ChainSpecificName)
	}

	if len(chainReaderDefinition.EventInputFields) != 0 {
		return fmt.Errorf(
			"%w: method %s has event topic fields defined, but is not an event",
			commontypes.ErrInvalidConfig,
			chainReaderDefinition.ChainSpecificName)
	}

	cr.contractBindings.AddReadBinding(contractName, methodName, &methodBinding{
		contractName: contractName,
		method:       methodName,
		client:       cr.client,
	})

	if err := cr.addEncoderDef(contractName, methodName, method.Inputs, method.ID, chainReaderDefinition); err != nil {
		return err
	}

	return cr.addDecoderDef(contractName, methodName, method.Outputs, chainReaderDefinition)
}

func (cr *chainReader) addEvent(contractName, eventName string, a abi.ABI, chainReaderDefinition types.ChainReaderDefinition) error {
	event, eventExists := a.Events[chainReaderDefinition.ChainSpecificName]
	if !eventExists {
		return fmt.Errorf("%w: event %s doesn't exist", commontypes.ErrInvalidConfig, chainReaderDefinition.ChainSpecificName)
	}

	filterArgs, topicInfo, indexArgNames := setupEventInput(event, chainReaderDefinition)
	if err := verifyEventInputsUsed(chainReaderDefinition, indexArgNames); err != nil {
		return err
	}

	if err := topicInfo.Init(); err != nil {
		return err
	}

	// Encoder def's codec won't be used to encode, only for its type as input for GetLatestValue
	if err := cr.addEncoderDef(contractName, eventName, filterArgs, nil, chainReaderDefinition); err != nil {
		return err
	}

	inputInfo, inputModifier, err := cr.getEventInput(chainReaderDefinition, contractName, eventName)
	if err != nil {
		return err
	}

	eb := &eventBinding{
		contractName:  contractName,
		eventName:     eventName,
		lp:            cr.lp,
		hash:          event.ID,
		inputInfo:     inputInfo,
		inputModifier: inputModifier,
		topicInfo:     topicInfo,
		id:            wrapItemType(contractName, eventName, false) + uuid.NewString(),
	}

	// set key binding for QueryKeys
	for index, topics := range event.Inputs {
		genericTopicName, ok := chainReaderDefinition.GenericTopicNames[topics.Name]
		if ok {
			cr.eventIndexBindings.Bind(eb, eventName+"-"+genericTopicName, index)
		}

	}

	cr.contractBindings.AddReadBinding(contractName, eventName, eb)

	return cr.addDecoderDef(contractName, eventName, event.Inputs, chainReaderDefinition)
}

func (cr *chainReader) getEventInput(def types.ChainReaderDefinition, contractName, eventName string) (
	types.CodecEntry, codec.Modifier, error) {
	inputInfo := cr.parsed.encoderDefs[wrapItemType(contractName, eventName, true)]
	inMod, err := def.InputModifications.ToModifier(evmDecoderHooks...)
	if err != nil {
		return nil, nil, err
	}

	// initialize the modification
	if _, err = inMod.RetypeToOffChain(reflect.PointerTo(inputInfo.CheckedType()), ""); err != nil {
		return nil, nil, err
	}

	return inputInfo, inMod, nil
}

func verifyEventInputsUsed(chainReaderDefinition types.ChainReaderDefinition, indexArgNames map[string]bool) error {
	for _, value := range chainReaderDefinition.EventInputFields {
		if !indexArgNames[abi.ToCamelCase(value)] {
			return fmt.Errorf("%w: %s is not an indexed argument of event %s", commontypes.ErrInvalidConfig, value, chainReaderDefinition.ChainSpecificName)
		}
	}
	return nil
}

func (cr *chainReader) addEncoderDef(contractName, methodName string, args abi.Arguments, prefix []byte, chainReaderDefinition types.ChainReaderDefinition) error {
	// ABI.Pack prepends the method.ID to the encodings, we'll need the encoder to do the same.
	inputMod, err := chainReaderDefinition.InputModifications.ToModifier(evmDecoderHooks...)
	if err != nil {
		return err
	}
	input := types.NewCodecEntry(args, prefix, inputMod)

	if err := input.Init(); err != nil {
		return err
	}

	cr.parsed.encoderDefs[wrapItemType(contractName, methodName, true)] = input
	return nil
}

func (cr *chainReader) addDecoderDef(contractName, methodName string, outputs abi.Arguments, def types.ChainReaderDefinition) error {
	mod, err := def.OutputModifications.ToModifier(evmDecoderHooks...)
	if err != nil {
		return err
	}
	output := types.NewCodecEntry(outputs, nil, mod)
	cr.parsed.decoderDefs[wrapItemType(contractName, methodName, false)] = output
	return output.Init()
}

func remapQueryKeyFilters(key string, eventIndexBindings EventIndexBindings, queryFilters []commontypes.QueryFilter) ([]commontypes.QueryFilter, error) {
	eventSig, address, _, err := eventIndexBindings.Get(key)
	if err != nil {
		return nil, err
	}

	var logFilters []commontypes.QueryFilter
	for _, queryFilter := range queryFilters {
		filter, err := remapQueryFilter(queryFilter)
		if err != nil {
			return nil, err
		}
		logFilters = append(logFilters, filter)
	}

	return append(logFilters, NewEventFilter(address, eventSig)), nil
}

// TODO values can't be string?
func remapQueryKeyByValuesFilters(key string, values []string, eventIndexBindings EventIndexBindings, queryFilters []commontypes.QueryFilter) ([]commontypes.QueryFilter, error) {
	eventSig, address, topicIndex, err := eventIndexBindings.Get(key)
	if err != nil {
		return nil, err
	}

	var logFilters []commontypes.QueryFilter
	for _, queryFilter := range queryFilters {
		filter, err := remapQueryFilter(queryFilter)
		if err != nil {
			return nil, err
		}
		logFilters = append(logFilters, filter)
	}

	return append(logFilters, NewEventTopicsByValueFilter(address, values, eventSig, topicIndex)), nil
}

// remapQueryFilter, changes some chain agnostic filters to match evm specific filters.
func remapQueryFilter(queryFilter commontypes.QueryFilter) (commontypes.QueryFilter, error) {
	switch filter := queryFilter.(type) {
	case *commontypes.AndFilter:
		var remappedFilters []commontypes.QueryFilter
		for _, f := range filter.Filters {
			remappedFilter, err := remapQueryFilter(f)
			if err != nil {
				return nil, err
			}
			remappedFilters = append(remappedFilters, remappedFilter)
		}
		return &commontypes.AndFilter{Filters: remappedFilters}, nil
	case *commontypes.ConfirmationsFilter:
		return NewFinalityFilter(filter)
	default:
		return filter, nil
	}
}

func setupEventInput(event abi.Event, def types.ChainReaderDefinition) ([]abi.Argument, types.CodecEntry, map[string]bool) {
	topicFieldDefs := map[string]bool{}
	for _, value := range def.EventInputFields {
		capFirstValue := abi.ToCamelCase(value)
		topicFieldDefs[capFirstValue] = true
	}

	filterArgs := make([]abi.Argument, 0, types.MaxTopicFields)
	inputArgs := make([]abi.Argument, 0, len(event.Inputs))
	indexArgNames := map[string]bool{}

	for _, input := range event.Inputs {
		if !input.Indexed {
			continue
		}

		filterWith := topicFieldDefs[abi.ToCamelCase(input.Name)]
		if filterWith {
			// When presenting the filter off-chain,
			// the user will provide the unhashed version of the input
			// The reader will hash topics if needed.
			inputUnindexed := input
			inputUnindexed.Indexed = false
			filterArgs = append(filterArgs, inputUnindexed)
		}

		inputArgs = append(inputArgs, input)
		indexArgNames[abi.ToCamelCase(input.Name)] = true
	}

	return filterArgs, types.NewCodecEntry(inputArgs, nil, nil), indexArgNames
}
