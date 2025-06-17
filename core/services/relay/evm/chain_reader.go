package evm

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	commoncodec "github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/read"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type ChainReaderService interface {
	services.ServiceCtx
	commontypes.ContractReader
}

type chainReader struct {
	commontypes.UnimplementedContractReader
	lggr     logger.Logger
	ht       logpoller.HeadTracker
	lp       logpoller.LogPoller
	client   EVMClient
	parsed   *codec.ParsedTypes
	bindings *read.BindingsRegistry
	codec    commontypes.RemoteCodec
	commonservices.StateMachine
}

type EVMClient interface {
	read.EVMBatchCaller
	read.EVMMethodClient
}

var _ ChainReaderService = (*chainReader)(nil)
var _ commontypes.ContractTypeProvider = &chainReader{}

// NewChainReaderService is a constructor for ChainReader, returns nil if there is any error
// Note that the ChainReaderService returned does not support anonymous events.
func NewChainReaderService(_ context.Context, lggr logger.Logger, lp logpoller.LogPoller, ht logpoller.HeadTracker, client EVMClient, config types.ChainReaderConfig) (ChainReaderService, error) {
	cr := &chainReader{
		lggr:     logger.Named(lggr, "ChainReader"),
		ht:       ht,
		lp:       lp,
		client:   client,
		bindings: read.NewBindingsRegistry(),
		parsed:   &codec.ParsedTypes{EncoderDefs: map[string]types.CodecEntry{}, DecoderDefs: map[string]types.CodecEntry{}},
	}

	var err error
	if err = cr.init(config.Contracts); err != nil {
		return nil, err
	}

	if cr.codec, err = cr.parsed.ToCodec(); err != nil {
		return nil, err
	}

	cr.bindings.SetBatchCaller(read.NewDynamicLimitedBatchCaller(
		cr.lggr,
		cr.codec,
		cr.client,
		read.DefaultRpcBatchSizeLimit,
		read.DefaultRpcBatchBackOffMultiplier,
		read.DefaultMaxParallelRpcCalls,
	))

	cr.bindings.SetCodecAll(cr.codec)

	return cr, err
}

func (cr *chainReader) init(chainContractReaders map[string]types.ChainContractReader) error {
	for contractName, chainContractReader := range chainContractReaders {
		contractAbi, err := abi.JSON(strings.NewReader(chainContractReader.ContractABI))
		if err != nil {
			return fmt.Errorf("failed to parse abi for contract: %s, err: %w", contractName, err)
		}

		var eventSigsForContractFilter evmtypes.HashArray
		for typeName, chainReaderDefinition := range chainContractReader.Configs {
			injectEVMSpecificCodecModifiers(chainReaderDefinition)

			switch chainReaderDefinition.ReadType {
			case types.Method:
				err = cr.addMethod(contractName, typeName, contractAbi, *chainReaderDefinition)
			case types.Event:
				partOfContractCommonFilter := slices.Contains(chainContractReader.GenericEventNames, typeName)
				if !partOfContractCommonFilter && !chainReaderDefinition.HasPollingFilter() {
					return fmt.Errorf(
						"%w: chain reader has no polling filter defined for contract: %s, event: %s",
						commontypes.ErrInvalidConfig, contractName, typeName)
				}

				eventOverridesContractFilter := chainReaderDefinition.HasPollingFilter()
				if eventOverridesContractFilter && partOfContractCommonFilter {
					return fmt.Errorf(
						"%w: conflicting chain reader polling filter definitions for contract: %s event: %s, can't have polling filter defined both on contract and event level",
						commontypes.ErrInvalidConfig, contractName, typeName)
				}

				if !eventOverridesContractFilter &&
					!slices.Contains(eventSigsForContractFilter, contractAbi.Events[chainReaderDefinition.ChainSpecificName].ID) {
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

		if !cr.bindings.HasContractBinding(contractName) {
			return fmt.Errorf("%w: no read bindings added for contract: %s", commontypes.ErrInvalidConfig, contractName)
		}

		if err = cr.bindings.SetFilter(contractName, chainContractReader.PollingFilter.ToLPFilter(eventSigsForContractFilter)); err != nil {
			return err
		}
	}
	return nil
}

// injectEVMSpecificCodecModifiers injects an AddressModifier into Input/OutputModifications of a ChainReaderDefinition.
func injectEVMSpecificCodecModifiers(chainReaderDefinition *types.ChainReaderDefinition) {
	for i, modConfig := range chainReaderDefinition.InputModifications {
		if addrModifierConfig, ok := modConfig.(*commoncodec.AddressBytesToStringModifierConfig); ok {
			addrModifierConfig.Modifier = codec.EVMAddressModifier{}
			chainReaderDefinition.InputModifications[i] = addrModifierConfig
		}
	}

	for i, modConfig := range chainReaderDefinition.OutputModifications {
		if addrModifierConfig, ok := modConfig.(*commoncodec.AddressBytesToStringModifierConfig); ok {
			addrModifierConfig.Modifier = codec.EVMAddressModifier{}
			chainReaderDefinition.OutputModifications[i] = addrModifierConfig
		}
	}
}

func (cr *chainReader) Name() string { return cr.lggr.Name() }

// Start registers polling filters if contracts are already bound.
func (cr *chainReader) Start(ctx context.Context) error {
	return cr.StartOnce("ChainReader", func() error {
		return cr.bindings.RegisterAll(ctx, cr.lp)
	})
}

func (cr *chainReader) Close() error {
	return cr.StopOnce("ChainReader", func() error {
		return nil
	})
}

func (cr *chainReader) Ready() error { return nil }

func (cr *chainReader) HealthReport() map[string]error {
	report := map[string]error{
		cr.Name(): cr.Healthy(),
	}

	commonservices.CopyHealth(report, cr.lp.HealthReport())
	commonservices.CopyHealth(report, cr.ht.HealthReport())
	return report
}

func (cr *chainReader) Bind(ctx context.Context, bindings []commontypes.BoundContract) error {
	return cr.bindings.Bind(ctx, cr.lp, bindings)
}

func (cr *chainReader) Unbind(ctx context.Context, bindings []commontypes.BoundContract) error {
	return cr.bindings.Unbind(ctx, cr.lp, bindings)
}

func (cr *chainReader) GetLatestValue(ctx context.Context, readName string, confidenceLevel primitives.ConfidenceLevel, params any, returnVal any) error {
	binding, address, err := cr.bindings.GetReader(readName)
	if err != nil {
		return err
	}

	ptrToValue, isValue := returnVal.(*values.Value)
	if !isValue {
		_, err = binding.GetLatestValueWithHeadData(ctx, common.HexToAddress(address), confidenceLevel, params, returnVal)
		if err != nil {
			return err
		}

		return nil
	}

	contractType, err := cr.CreateContractType(readName, false)
	if err != nil {
		return err
	}

	if err = cr.GetLatestValue(ctx, readName, confidenceLevel, params, contractType); err != nil {
		return err
	}

	value, err := values.Wrap(contractType)
	if err != nil {
		return err
	}

	*ptrToValue = value

	return nil
}

func (cr *chainReader) GetLatestValueWithHeadData(ctx context.Context, readName string, confidenceLevel primitives.ConfidenceLevel, params any, returnVal any) (head *commontypes.Head, err error) {
	binding, address, err := cr.bindings.GetReader(readName)
	if err != nil {
		return nil, err
	}

	ptrToValue, isValue := returnVal.(*values.Value)
	if !isValue {
		return binding.GetLatestValueWithHeadData(ctx, common.HexToAddress(address), confidenceLevel, params, returnVal)
	}

	contractType, err := cr.CreateContractType(readName, false)
	if err != nil {
		return nil, err
	}

	head, err = cr.GetLatestValueWithHeadData(ctx, readName, confidenceLevel, params, contractType)
	if err != nil {
		return nil, err
	}

	value, err := values.Wrap(contractType)
	if err != nil {
		return nil, err
	}

	*ptrToValue = value

	return head, nil
}

func (cr *chainReader) BatchGetLatestValues(ctx context.Context, request commontypes.BatchGetLatestValuesRequest) (commontypes.BatchGetLatestValuesResult, error) {
	return cr.bindings.BatchGetLatestValues(ctx, request)
}

func (cr *chainReader) QueryKey(
	ctx context.Context,
	contract commontypes.BoundContract,
	filter query.KeyFilter,
	limitAndSort query.LimitAndSort,
	sequenceDataType any,
) ([]commontypes.Sequence, error) {
	binding, address, err := cr.bindings.GetReader(contract.ReadIdentifier(filter.Key))
	if err != nil {
		return nil, err
	}

	_, isValuePtr := sequenceDataType.(*values.Value)
	if !isValuePtr {
		return binding.QueryKey(ctx, common.HexToAddress(address), filter, limitAndSort, sequenceDataType)
	}

	dataTypeFromReadIdentifier, err := cr.CreateContractType(contract.ReadIdentifier(filter.Key), false)
	if err != nil {
		return nil, err
	}

	sequence, err := binding.QueryKey(ctx, common.HexToAddress(address), filter, limitAndSort, dataTypeFromReadIdentifier)
	if err != nil {
		return nil, err
	}

	sequenceOfValues := make([]commontypes.Sequence, len(sequence))
	for idx, entry := range sequence {
		value, err := values.Wrap(entry.Data)
		if err != nil {
			return nil, err
		}
		sequenceOfValues[idx] = commontypes.Sequence{
			Cursor: entry.Cursor,
			Head:   entry.Head,
			Data:   &value,
		}
	}

	return sequenceOfValues, nil
}

func (cr *chainReader) QueryKeys(ctx context.Context, filters []commontypes.ContractKeyFilter,
	limitAndSort query.LimitAndSort) (iter.Seq2[string, commontypes.Sequence], error) {
	eventQueries := make([]read.EventQuery, 0, len(filters))
	for _, filter := range filters {
		binding, address, err := cr.bindings.GetReader(filter.Contract.ReadIdentifier(filter.KeyFilter.Key))
		if err != nil {
			return nil, err
		}

		sequenceDataType := filter.SequenceDataType
		_, isValuePtr := filter.SequenceDataType.(*values.Value)
		if isValuePtr {
			sequenceDataType, err = cr.CreateContractType(filter.Contract.ReadIdentifier(filter.Key), false)
			if err != nil {
				return nil, err
			}
		}

		eventBinding, ok := binding.(*read.EventBinding)
		if !ok {
			return nil, fmt.Errorf("query key %s is not an event", filter.KeyFilter.Key)
		}

		eventQueries = append(eventQueries, read.EventQuery{
			Filter:           filter.KeyFilter,
			SequenceDataType: sequenceDataType,
			IsValuePtr:       isValuePtr,
			EventBinding:     eventBinding,
			Address:          common.HexToAddress(address),
		})
	}

	return read.MultiEventTypeQuery(ctx, cr.lp, eventQueries, limitAndSort)
}

func (cr *chainReader) CreateContractType(readIdentifier string, forEncoding bool) (any, error) {
	return cr.codec.CreateType(cr.bindings.ReadTypeIdentifier(readIdentifier, forEncoding), forEncoding)
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

	confirmations, err := ConfirmationsFromConfig(chainReaderDefinition.ConfidenceConfirmations)
	if err != nil {
		return err
	}

	if err = cr.bindings.AddReader(contractName, methodName, read.NewMethodBinding(contractName, methodName, cr.client, cr.ht, confirmations, cr.lggr)); err != nil {
		return err
	}

	if err = cr.addEncoderDef(contractName, methodName, method.Inputs, method.ID, chainReaderDefinition.InputModifications); err != nil {
		return err
	}

	return cr.addDecoderDef(contractName, methodName, method.Outputs, chainReaderDefinition.OutputModifications)
}

func (cr *chainReader) addEvent(contractName, eventName string, a abi.ABI, chainReaderDefinition types.ChainReaderDefinition) error {
	event, eventExists := a.Events[chainReaderDefinition.ChainSpecificName]
	if !eventExists {
		return fmt.Errorf("%w: event %q doesn't exist", commontypes.ErrInvalidConfig, chainReaderDefinition.ChainSpecificName)
	}

	indexedAsUnIndexedABITypes, indexedTopicsCodecTypes, dwsDetails := getEventTypes(event)
	if err := indexedTopicsCodecTypes.Init(); err != nil {
		return err
	}

	// Encoder defs codec won't be used for encoding, but for storing caller filtering params which won't be hashed.
	err := cr.addEncoderDef(contractName, eventName, indexedAsUnIndexedABITypes, nil, chainReaderDefinition.InputModifications)
	if err != nil {
		return err
	}

	codecTypes, codecModifiers := make(map[string]types.CodecEntry), make(map[string]commoncodec.Modifier)
	topicTypeID := codec.WrapItemType(contractName, eventName, true)
	codecTypes[topicTypeID], codecModifiers[topicTypeID] = cr.getEventItemTypeAndModifier(topicTypeID)

	confirmations, err := ConfirmationsFromConfig(chainReaderDefinition.ConfidenceConfirmations)
	if err != nil {
		return err
	}

	eb := read.NewEventBinding(contractName, eventName, cr.lp, event.ID, indexedTopicsCodecTypes, confirmations)
	if eventDefinitions := chainReaderDefinition.EventDefinitions; eventDefinitions != nil {
		if eventDefinitions.PollingFilter != nil {
			eb.SetFilter(eventDefinitions.PollingFilter.ToLPFilter(evmtypes.HashArray{a.Events[event.Name].ID}))
		}

		topicsDetails, topicsCodecTypeInfo, topicsModifiers, initQueryingErr := cr.initTopicQuerying(contractName, eventName, event.Inputs, eventDefinitions.GenericTopicNames, chainReaderDefinition.InputModifications)
		if initQueryingErr != nil {
			return initQueryingErr
		}
		maps.Copy(codecTypes, topicsCodecTypeInfo)
		// TODO BCFR-44 reused GetLatestValue params modifiers, probably can be left like this
		maps.Copy(codecModifiers, topicsModifiers)

		// TODO BCFR-44 no dw modifier for now
		dataWordsDetails, dWSCodecTypeInfo, initDWQueryingErr := cr.initDWQuerying(contractName, eventName, dwsDetails, eventDefinitions.GenericDataWordDetails)
		if initDWQueryingErr != nil {
			return fmt.Errorf("failed to init dw querying for event: %q, err: %w", eventName, initDWQueryingErr)
		}
		maps.Copy(codecTypes, dWSCodecTypeInfo)

		eb.SetTopicDetails(topicsDetails)
		eb.SetDataWordsDetails(dataWordsDetails)
	}

	eb.SetCodecTypesAndModifiers(codecTypes, codecModifiers)
	if err = cr.bindings.AddReader(contractName, eventName, eb); err != nil {
		return err
	}

	return cr.addDecoderDef(contractName, eventName, event.Inputs, chainReaderDefinition.OutputModifications)
}

// initTopicQuerying registers codec types and modifiers for topics to be used for typing value comparator QueryKey filters.
func (cr *chainReader) initTopicQuerying(contractName, eventName string, eventInputs abi.Arguments, genericTopicNames map[string]string, inputModifications commoncodec.ModifiersConfig) (map[string]read.TopicDetail, map[string]types.CodecEntry, map[string]commoncodec.Modifier, error) {
	topicsDetails := make(map[string]read.TopicDetail)
	topicsTypes := make(map[string]types.CodecEntry)
	topicsModifiers := make(map[string]commoncodec.Modifier)
	for topicIndex, topic := range eventInputs {
		genericTopicName, ok := genericTopicNames[topic.Name]
		if ok {
			topicsDetails[genericTopicName] = read.TopicDetail{Argument: topic, Index: uint64(topicIndex + 1)}

			topicTypeID := eventName + "." + genericTopicName
			if err := cr.addEncoderDef(contractName, topicTypeID, abi.Arguments{{Type: topic.Type}}, nil, inputModifications); err != nil {
				return nil, nil, nil, err
			}

			topicCodecTypeID := codec.WrapItemType(contractName, topicTypeID, true)
			topicsTypes[topicCodecTypeID], topicsModifiers[topicCodecTypeID] = cr.getEventItemTypeAndModifier(topicCodecTypeID)
		}
	}
	return topicsDetails, topicsTypes, topicsModifiers, nil
}

// initDWQuerying registers codec types for evm data words to be used for typing value comparator QueryKey filters.
func (cr *chainReader) initDWQuerying(contractName, eventName string, abiDWsDetails map[string]read.DataWordDetail, cfgDWsDetails map[string]types.DataWordDetail) (map[string]read.DataWordDetail, map[string]types.CodecEntry, error) {
	dWsDetail, err := cr.constructDWDetails(cfgDWsDetails, abiDWsDetails)
	if err != nil {
		return nil, nil, err
	}

	dwsCodecTypeInfo := make(map[string]types.CodecEntry)
	for genericName := range cfgDWsDetails {
		dwDetail, exists := dWsDetail[genericName]
		if !exists {
			return nil, nil, fmt.Errorf("failed to find data word: %q, it either doesn't exist or can't be searched for", genericName)
		}

		dwTypeID := eventName + "." + genericName
		if err = cr.addEncoderDef(contractName, dwTypeID, abi.Arguments{abi.Argument{Type: dwDetail.Type}}, nil, nil); err != nil {
			return nil, nil, fmt.Errorf("failed to init codec for data word: %q on index: %d, err: %w", genericName, dwDetail.Index, err)
		}

		dwCodecTypeID := codec.WrapItemType(contractName, dwTypeID, true)
		dwsCodecTypeInfo[dwCodecTypeID] = cr.parsed.EncoderDefs[dwCodecTypeID]
	}

	return dWsDetail, dwsCodecTypeInfo, nil
}

// constructDWDetails combines data word details from config and abi.
func (cr *chainReader) constructDWDetails(cfgDWsDetails map[string]types.DataWordDetail, abiDWsDetails map[string]read.DataWordDetail) (map[string]read.DataWordDetail, error) {
	dWsDetail := make(map[string]read.DataWordDetail)
	for genericName, cfgDWDetail := range cfgDWsDetails {
		for eventID, dWDetail := range abiDWsDetails {
			// Extract field name in this manner to account for nested fields
			fieldName := strings.Join(strings.Split(eventID, ".")[1:], ".")
			if fieldName == cfgDWDetail.Name {
				dWsDetail[genericName] = dWDetail
				break
			}
		}
	}

	// if dw detail isn't set, the index can't be programmatically determined, so we get index and type from cfg
	for genericName, cfgDWDetail := range cfgDWsDetails {
		dwDetail, exists := dWsDetail[genericName]
		if exists && cfgDWDetail.Index != nil {
			return nil, fmt.Errorf("data word: %q at index: %d details, were calculated automatically and shouldn't be manully overridden by cfg", genericName, dwDetail.Index)
		}

		if cfgDWDetail.Index != nil {
			abiTyp, err := abi.NewType(cfgDWDetail.Type, "", nil)
			if err != nil {
				return nil, fmt.Errorf("bad abi type: %q provided for data word: %q at index: %d in config", cfgDWDetail.Type, genericName, *cfgDWDetail.Index)
			}
			dWsDetail[genericName] = read.DataWordDetail{Argument: abi.Argument{Type: abiTyp}, Index: *cfgDWDetail.Index}
		}
	}
	return dWsDetail, nil
}

// getEventItemTypeAndModifier returns codec entry for expected incoming event item and the modifier.
func (cr *chainReader) getEventItemTypeAndModifier(itemType string) (types.CodecEntry, commoncodec.Modifier) {
	inputTypeInfo := cr.parsed.EncoderDefs[itemType]
	return inputTypeInfo, inputTypeInfo.Modifier()
}

func (cr *chainReader) addEncoderDef(contractName, itemType string, args abi.Arguments, prefix []byte, inputModifications commoncodec.ModifiersConfig) error {
	// ABI.Pack prepends the method.ID to the encodings, we'll need the encoder to do the same.
	inputMod, err := inputModifications.ToModifier(codec.DecoderHooks...)
	if err != nil {
		return err
	}
	input := types.NewCodecEntry(args, prefix, inputMod)

	if err = input.Init(); err != nil {
		return err
	}

	cr.parsed.EncoderDefs[codec.WrapItemType(contractName, itemType, true)] = input
	return nil
}

func (cr *chainReader) addDecoderDef(contractName, itemType string, outputs abi.Arguments, outputModifications commoncodec.ModifiersConfig) error {
	mod, err := outputModifications.ToModifier(codec.DecoderHooks...)
	if err != nil {
		return err
	}
	output := types.NewCodecEntry(outputs, nil, mod)
	cr.parsed.DecoderDefs[codec.WrapItemType(contractName, itemType, false)] = output
	return output.Init()
}

// getEventTypes returns abi args where indexed flag is set to false because we expect caller to filter with params that aren't hashed,
// codecEntry where expected on chain types are set, for e.g. indexed topics of type string or uint8[32] array are expected as common.Hash onchain,
// and un-indexed data info in form of evm indexed 32 byte data words.
func getEventTypes(event abi.Event) ([]abi.Argument, types.CodecEntry, map[string]read.DataWordDetail) {
	indexedAsUnIndexedTypes := make([]abi.Argument, 0, types.MaxTopicFields)
	indexedTypes := make([]abi.Argument, 0, len(event.Inputs))
	for _, input := range event.Inputs {
		if input.Indexed {
			indexedAsUnIndexed := input
			indexedAsUnIndexed.Indexed = false
			// when presenting the filter off-chain, the caller will provide the unHashed version of the input and CR will hash topics when needed.
			indexedAsUnIndexedTypes = append(indexedAsUnIndexedTypes, indexedAsUnIndexed)
			indexedTypes = append(indexedTypes, input)
		}
	}

	return indexedAsUnIndexedTypes, types.NewCodecEntry(indexedTypes, nil, nil), getDWIndexesWithTypes(event.Name, event.Inputs)
}

func getDWIndexesWithTypes(eventName string, eventInputs abi.Arguments) map[string]read.DataWordDetail {
	var dwIndexOffset int
	dataWords := make(map[string]read.DataWordDetail)
	dynamicQueue := make([]abi.Argument, 0)

	for _, input := range eventInputs.NonIndexed() {
		// each dynamic field has an extra field that stores the dwIndexOffset that points to the start of the dynamic data.
		if isDynamic(input.Type) {
			dynamicQueue = append(dynamicQueue, input)
			dwIndexOffset++
		} else {
			dwIndexOffset = processDWStaticField(input.Type, eventName+"."+input.Name, dataWords, dwIndexOffset)
		}
	}

	processDWDynamicFields(eventName, dataWords, dynamicQueue, dwIndexOffset)
	return dataWords
}

func isDynamic(fieldType abi.Type) bool {
	switch fieldType.T {
	case abi.StringTy, abi.SliceTy, abi.BytesTy:
		return true
	case abi.TupleTy:
		// If one element in a struct is dynamic, the whole struct is treated as dynamic.
		for _, elem := range fieldType.TupleElems {
			if isDynamic(*elem) {
				return true
			}
		}
	}
	return false
}

func processDWStaticField(fieldType abi.Type, parentFieldPath string, dataWords map[string]read.DataWordDetail, index int) int {
	switch fieldType.T {
	case abi.TupleTy:
		// Recursively process tuple elements
		for i, tupleElem := range fieldType.TupleElems {
			fieldName := fieldType.TupleRawNames[i]
			fullFieldPath := fmt.Sprintf("%s.%s", parentFieldPath, fieldName)
			index = processDWStaticField(*tupleElem, fullFieldPath, dataWords, index)
		}
		return index
	case abi.ArrayTy:
		// Static arrays are not searchable, however, we can reliably calculate their size so that the fields
		// after them can be searched.
		return index + fieldType.Size
	default:
		dataWords[parentFieldPath] = read.DataWordDetail{
			Index:    index,
			Argument: abi.Argument{Type: fieldType},
		}
		return index + 1
	}
}

// processDWDynamicFields indexes static fields in dynamic structs.
// These fields come first after the static fields in the event encoding, so we can calculate their indices.
func processDWDynamicFields(eventName string, dataWords map[string]read.DataWordDetail, dynamicQueue []abi.Argument, dwIndex int) {
	for _, arg := range dynamicQueue {
		switch arg.Type.T {
		case abi.TupleTy:
			for i, tupleElem := range arg.Type.TupleElems {
				// any field after a dynamic field can't be predictably indexed.
				if isDynamic(*tupleElem) {
					return
				}
				dwIndex = processDWStaticField(*tupleElem, fmt.Sprintf("%s.%s.%s", eventName, arg.Name, arg.Type.TupleRawNames[i]), dataWords, dwIndex)
			}
		default:
			// exit if we see a dynamic field, as we can't predict the index of fields after it.
			return
		}
	}
}

// ConfirmationsFromConfig maps chain agnostic confidence levels defined in config to predefined EVM finality.
func ConfirmationsFromConfig(values map[string]int) (map[primitives.ConfidenceLevel]evmtypes.Confirmations, error) {
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
