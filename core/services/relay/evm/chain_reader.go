package evm

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type ChainReaderService interface {
	services.ServiceCtx
	commontypes.ContractReader
}

type chainReader struct {
	lggr             logger.Logger
	lp               logpoller.LogPoller
	client           evmclient.Client
	contractBindings bindings
	parsed           *parsedTypes
	codec            commontypes.RemoteCodec
	commonservices.StateMachine
}

var _ ChainReaderService = (*chainReader)(nil)

// NewChainReaderService is a constructor for ChainReader, returns nil if there is any error
// Note that the ChainReaderService returned does not support anonymous events.
func NewChainReaderService(ctx context.Context, lggr logger.Logger, lp logpoller.LogPoller, client evmclient.Client, config types.ChainReaderConfig) (ChainReaderService, error) {
	cr := &chainReader{
		lggr:             lggr.Named("ChainReader"),
		lp:               lp,
		client:           client,
		contractBindings: bindings{},
		parsed:           &parsedTypes{encoderDefs: map[string]types.CodecEntry{}, decoderDefs: map[string]types.CodecEntry{}},
	}

	var err error
	if err = cr.init(config.Contracts); err != nil {
		return nil, err
	}

	if cr.codec, err = cr.parsed.toCodec(); err != nil {
		return nil, err
	}

	err = cr.contractBindings.ForEach(ctx, func(c context.Context, rbs *contractBindings) error {
		for _, rb := range rbs.readBindings {
			rb.SetCodec(cr.codec)
		}
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
	return cr.contractBindings.Bind(ctx, cr.lp, bindings)
}

func (cr *chainReader) QueryKey(ctx context.Context, contractName string, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]commontypes.Sequence, error) {
	b, err := cr.contractBindings.GetReadBinding(contractName, filter.Key)
	if err != nil {
		return nil, err
	}

	return b.QueryKey(ctx, filter, limitAndSort, sequenceDataType)
}

func (cr *chainReader) init(chainContractReaders map[string]types.ChainContractReader) error {
	for contractName, chainContractReader := range chainContractReaders {
		contractAbi, err := abi.JSON(strings.NewReader(chainContractReader.ContractABI))
		if err != nil {
			return err
		}

		contractFilterEvents := chainContractReader.ContractPollingFilter.GenericEventNames
		var eventSigsForContractFilter evmtypes.HashArray
		for typeName, chainReaderDefinition := range chainContractReader.Configs {
			switch chainReaderDefinition.ReadType {
			case types.Method:
				err = cr.addMethod(contractName, typeName, contractAbi, *chainReaderDefinition)
			case types.Event:
				partOfContractFilter := slices.Contains(contractFilterEvents, typeName)
				if !partOfContractFilter && !chainReaderDefinition.HasPollingFilter() {
					return fmt.Errorf(
						"%w: chain reader has no polling filter defined for contract: %s, event: %s",
						commontypes.ErrInvalidConfig, contractName, typeName)
				}

				eventOverridesContractFilter := chainReaderDefinition.HasPollingFilter()
				if eventOverridesContractFilter && partOfContractFilter {
					return fmt.Errorf(
						"%w: conflicting chain reader polling filter definitions for contract: %s event: %s, can't have polling filter defined both on contract and event level",
						commontypes.ErrInvalidConfig, contractName, typeName)
				}

				if !eventOverridesContractFilter {
					eventSigsForContractFilter = append(eventSigsForContractFilter, contractAbi.Events[chainReaderDefinition.ChainSpecificName].ID)
				}

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
		cr.contractBindings[contractName].pollingFilter = chainContractReader.PollingFilter.ToLPFilter(eventSigsForContractFilter)
	}
	return nil
}

func (cr *chainReader) Start(ctx context.Context) error {
	return cr.StartOnce("ChainReader", func() error {
		return cr.contractBindings.ForEach(ctx, func(c context.Context, rbs *contractBindings) error {
			return rbs.Register(ctx, cr.lp)
		})
	})
}

func (cr *chainReader) Close() error {
	return cr.StopOnce("ChainReader", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return cr.contractBindings.ForEach(ctx, func(c context.Context, rbs *contractBindings) error {
			return rbs.Unregister(ctx, cr.lp)
		})
	})
}

func (cr *chainReader) Ready() error { return nil }
func (cr *chainReader) HealthReport() map[string]error {
	return map[string]error{cr.Name(): nil}
}

func (cr *chainReader) CreateContractType(contractName, itemType string, forEncoding bool) (any, error) {
	return cr.codec.CreateType(wrapItemType(contractName, itemType, forEncoding), forEncoding)
}

func wrapItemType(contractName, itemType string, isParams bool) string {
	if isParams {
		return fmt.Sprintf("params.%s.%s", contractName, itemType)
	}
	return fmt.Sprintf("return.%s.%s", contractName, itemType)
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

	var inputFields []string
	if chainReaderDefinition.EventDefinitions != nil {
		inputFields = chainReaderDefinition.EventDefinitions.InputFields
	}

	filterArgs, codecTopicInfo, indexArgNames := setupEventInput(event, inputFields)
	if err := verifyEventInputsUsed(eventName, inputFields, indexArgNames); err != nil {
		return err
	}

	if err := codecTopicInfo.Init(); err != nil {
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

	confirmations, err := confirmationsFromConfig(chainReaderDefinition.ConfidenceConfirmations)
	if err != nil {
		return err
	}

	eb := &eventBinding{
		contractName:         contractName,
		eventName:            eventName,
		lp:                   cr.lp,
		hash:                 event.ID,
		inputInfo:            inputInfo,
		inputModifier:        inputModifier,
		codecTopicInfo:       codecTopicInfo,
		topics:               make(map[string]topicDetail),
		eventDataWords:       make(map[string]uint8),
		confirmationsMapping: confirmations,
	}

	if eventDefinitions := chainReaderDefinition.EventDefinitions; eventDefinitions != nil {
		if eventDefinitions.PollingFilter != nil {
			eb.FilterRegisterer = &FilterRegisterer{
				pollingFilter: eventDefinitions.PollingFilter.ToLPFilter(evmtypes.HashArray{a.Events[event.Name].ID}),
				filterLock:    sync.Mutex{},
			}
		}

		if eventDefinitions.GenericDataWordNames != nil {
			eb.eventDataWords = eventDefinitions.GenericDataWordNames
		}

		cr.addQueryingReadBindings(contractName, eventDefinitions.GenericTopicNames, event.Inputs, eb)
	}

	cr.contractBindings.AddReadBinding(contractName, eventName, eb)

	return cr.addDecoderDef(contractName, eventName, event.Inputs, chainReaderDefinition)
}

// addQueryingReadBindings reuses the eventBinding and maps it to topic and dataWord keys used for QueryKey.
func (cr *chainReader) addQueryingReadBindings(contractName string, genericTopicNames map[string]string, eventInputs abi.Arguments, eb *eventBinding) {
	// add topic read readBindings for QueryKey
	for topicIndex, topic := range eventInputs {
		genericTopicName, ok := genericTopicNames[topic.Name]
		if ok {
			eb.topics[genericTopicName] = topicDetail{
				Argument: topic,
				Index:    uint64(topicIndex),
			}
		}
		cr.contractBindings.AddReadBinding(contractName, genericTopicName, eb)
	}

	// add data word read readBindings for QueryKey
	for genericDataWordName := range eb.eventDataWords {
		cr.contractBindings.AddReadBinding(contractName, genericDataWordName, eb)
	}
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

func verifyEventInputsUsed(eventName string, inputFields []string, indexArgNames map[string]bool) error {
	for _, value := range inputFields {
		if !indexArgNames[abi.ToCamelCase(value)] {
			return fmt.Errorf("%w: %s is not an indexed argument of event %s", commontypes.ErrInvalidConfig, value, eventName)
		}
	}
	return nil
}

func (cr *chainReader) addEncoderDef(contractName, itemType string, args abi.Arguments, prefix []byte, chainReaderDefinition types.ChainReaderDefinition) error {
	// ABI.Pack prepends the method.ID to the encodings, we'll need the encoder to do the same.
	inputMod, err := chainReaderDefinition.InputModifications.ToModifier(evmDecoderHooks...)
	if err != nil {
		return err
	}
	input := types.NewCodecEntry(args, prefix, inputMod)

	if err = input.Init(); err != nil {
		return err
	}

	cr.parsed.encoderDefs[wrapItemType(contractName, itemType, true)] = input
	return nil
}

func (cr *chainReader) addDecoderDef(contractName, itemType string, outputs abi.Arguments, def types.ChainReaderDefinition) error {
	mod, err := def.OutputModifications.ToModifier(evmDecoderHooks...)
	if err != nil {
		return err
	}
	output := types.NewCodecEntry(outputs, nil, mod)
	cr.parsed.decoderDefs[wrapItemType(contractName, itemType, false)] = output
	return output.Init()
}

func setupEventInput(event abi.Event, inputFields []string) ([]abi.Argument, types.CodecEntry, map[string]bool) {
	topicFieldDefs := map[string]bool{}
	for _, value := range inputFields {
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

func confirmationsFromConfig(values map[string]int) (map[primitives.ConfidenceLevel]evmtypes.Confirmations, error) {
	mappings := map[primitives.ConfidenceLevel]evmtypes.Confirmations{
		primitives.Unconfirmed: evmtypes.Unconfirmed,
		primitives.Finalized:   evmtypes.Finalized,
	}

	if values == nil {
		return mappings, nil
	}

	for key, mapped := range values {
		mappings[primitives.ConfidenceLevel(key)] = evmtypes.Confirmations(mapped)
	}

	if mappings[primitives.Finalized] != evmtypes.Finalized &&
		mappings[primitives.Finalized] > mappings[primitives.Unconfirmed] {
		return nil, errors.New("finalized confidence level should map to -1 or a higher value than 0")
	}

	return mappings, nil
}
