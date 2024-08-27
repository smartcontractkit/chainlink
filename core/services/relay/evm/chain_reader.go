package evm

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	commoncodec "github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
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
	client   evmclient.Client
	parsed   *codec.ParsedTypes
	bindings *read.BindingsRegistry
	codec    commontypes.RemoteCodec
	commonservices.StateMachine
}

var _ ChainReaderService = (*chainReader)(nil)
var _ commontypes.ContractTypeProvider = &chainReader{}

// NewChainReaderService is a constructor for ChainReader, returns nil if there is any error
// Note that the ChainReaderService returned does not support anonymous events.
func NewChainReaderService(ctx context.Context, lggr logger.Logger, lp logpoller.LogPoller, ht logpoller.HeadTracker, client evmclient.Client, config types.ChainReaderConfig) (ChainReaderService, error) {
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

		if err := cr.bindings.SetFilter(contractName, chainContractReader.PollingFilter.ToLPFilter(eventSigsForContractFilter)); err != nil {
			return err
		}
	}
	return nil
}

func (cr *chainReader) Name() string { return cr.lggr.Name() }

// Start registers polling filters if contracts are already bound.
func (cr *chainReader) Start(ctx context.Context) error {
	return cr.StartOnce("ChainReader", func() error {
		return cr.bindings.RegisterAll(ctx, cr.lp)
	})
}

// Close unregisters polling filters for bound contracts.
func (cr *chainReader) Close() error {
	return cr.StopOnce("ChainReader", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		return cr.bindings.UnregisterAll(ctx, cr.lp)
	})
}

func (cr *chainReader) Ready() error { return nil }

func (cr *chainReader) HealthReport() map[string]error {
	return map[string]error{cr.Name(): nil}
}

func (cr *chainReader) GetLatestValue(ctx context.Context, readName string, confidenceLevel primitives.ConfidenceLevel, params any, returnVal any) error {
	binding, address, err := cr.bindings.GetReader(readName)
	if err != nil {
		return err
	}

	return binding.GetLatestValue(ctx, common.HexToAddress(address), confidenceLevel, params, returnVal)
}

func (cr *chainReader) BatchGetLatestValues(ctx context.Context, request commontypes.BatchGetLatestValuesRequest) (commontypes.BatchGetLatestValuesResult, error) {
	return cr.bindings.BatchGetLatestValues(ctx, request)
}

func (cr *chainReader) Bind(ctx context.Context, bindings []commontypes.BoundContract) error {
	return cr.bindings.Bind(ctx, cr.lp, bindings)
}

func (cr *chainReader) Unbind(ctx context.Context, bindings []commontypes.BoundContract) error {
	return cr.bindings.Unbind(ctx, cr.lp, bindings)
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

	return binding.QueryKey(ctx, common.HexToAddress(address), filter, limitAndSort, sequenceDataType)
}

func (cr *chainReader) CreateContractType(readIdentifier string, forEncoding bool) (any, error) {
	return cr.codec.CreateType(cr.bindings.ReadTypeIdentifier(readIdentifier, forEncoding), forEncoding)
}

func WrapItemType(contractName, itemType string, isParams bool) string {
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

	confirmations, err := ConfirmationsFromConfig(chainReaderDefinition.ConfidenceConfirmations)
	if err != nil {
		return err
	}

	cr.bindings.AddReader(contractName, methodName, read.NewMethodBinding(contractName, methodName, cr.client, cr.ht, confirmations, cr.lggr))

	if err = cr.addEncoderDef(contractName, methodName, method.Inputs, method.ID, chainReaderDefinition.InputModifications); err != nil {
		return err
	}

	return cr.addDecoderDef(contractName, methodName, method.Outputs, chainReaderDefinition.OutputModifications)
}

func (cr *chainReader) addEvent(contractName, eventName string, a abi.ABI, chainReaderDefinition types.ChainReaderDefinition) error {
	event, eventExists := a.Events[chainReaderDefinition.ChainSpecificName]
	if !eventExists {
		return fmt.Errorf("%w: event %s doesn't exist", commontypes.ErrInvalidConfig, chainReaderDefinition.ChainSpecificName)
	}

	indexedAsUnIndexedTypes, indexedTopicsCodecInfo, eventDws := getEventTypes(event)
	if err := indexedTopicsCodecInfo.Init(); err != nil {
		return err
	}

	// Encoder defs codec won't be used for encoding, but for storing caller filtering params which won't be hashed.
	err := cr.addEncoderDef(contractName, eventName, indexedAsUnIndexedTypes, nil, chainReaderDefinition.InputModifications)
	if err != nil {
		return err
	}

	filtersTypes, filtersModifiers := make(map[string]types.CodecEntry), make(map[string]codec.Modifier)
	topicsInputIdentifier := WrapItemType(contractName, eventName, true)
	filtersTypes[topicsInputIdentifier], filtersModifiers[topicsInputIdentifier], err = cr.getEventTopicsInput(topicsInputIdentifier, chainReaderDefinition.InputModifications)
	if err != nil {
		return err
	}

	confirmations, err := ConfirmationsFromConfig(chainReaderDefinition.ConfidenceConfirmations)
	if err != nil {
		return err
	}

	eb := read.NewEventBinding(contractName, eventName, cr.lp, event.ID, indexedTopicsCodecInfo, confirmations)
	if eventDefinitions := chainReaderDefinition.EventDefinitions; eventDefinitions != nil {
		if eventDefinitions.PollingFilter != nil {
			eb.SetFilter(eventDefinitions.PollingFilter.ToLPFilter(evmtypes.HashArray{a.Events[event.Name].ID}))
		}

		topics, err := cr.initTopicQuerying(contractName, eventName, event.Inputs, eventDefinitions.GenericTopicNames, chainReaderDefinition.InputModifications)
		if err != nil {
			return err
		}

		dwsInfo, dWSCodecTypeInfo, err := cr.initDWQuerying(contractName, eventName, eventDws, eventDefinitions.GenericDataWordNames)
		if err != nil {
			return err
		}

		eb.topics = topics
		eb.dataWordsInfo = dwsInfo
		maps.Copy(filtersTypes, dWSCodecTypeInfo)
	}

	eb.eventTypes = filtersTypes
	eb.eventModifiers = filtersModifiers

	cr.bindings.AddReader(contractName, eventName, eb)
	return cr.addDecoderDef(contractName, eventName, event.Inputs, chainReaderDefinition.OutputModifications)
}

// initTopicQuerying adds codec types for topics that are used for typing over the wire bytes used in value comparator QueryKey filters.
func (cr *chainReader) initTopicQuerying(contractName, eventName string, eventInputs abi.Arguments, genericTopicNames map[string]string, inputModifications codec.ModifiersConfig) (map[string]topicDetail, error) {
	topicDetails := make(map[string]topicDetail)
	for topicIndex, topic := range eventInputs {
		genericTopicName, ok := genericTopicNames[topic.Name]
		if ok {
			// TODO how did this work before with topicIndex not having a +1
			topicDetails[genericTopicName] = topicDetail{Argument: topic, Index: uint64(topicIndex + 1)}
			// Encoder defs codec won't be used for encoding, but for storing caller filtering params which won't be hashed.
			if err := cr.addEncoderDef(contractName, eventName+"."+genericTopicName, abi.Arguments{{Type: topic.Type}}, nil, inputModifications); err != nil {
				return nil, err
			}
		}
	}
	return topicDetails, nil
}

// initDWQuerying adds codec types for data words that are used for typing over the wire bytes used in value comparator QueryKey filters.
func (cr *chainReader) initDWQuerying(contractName, eventName string, dataWordsInfo []dataWordInfo, genericDataWordNames map[string]uint8) (map[string]dataWordInfo, map[string]types.CodecEntry, error) {
	dwsCodecEntries := make(map[string]types.CodecEntry)
	dwsInfo := make(map[string]dataWordInfo)

	for genericDataWordName, dwIndex := range genericDataWordNames {
		founDW := false
		for _, dwInfo := range dataWordsInfo {
			if dwInfo.index != dwIndex {
				continue
			}

			founDW = true
			dwName := eventName + "." + genericDataWordName

			if err := cr.addEncoderDef(contractName, dwName, abi.Arguments{dataWordsInfo[dwIndex].typ}, nil, nil); err != nil {
				return nil, nil, fmt.Errorf("%w: failed to init codec for data word %s on index %d querying for event: %q", err, dwName, dwIndex, eventName)
			}

			dwsInfo[genericDataWordName] = dataWordsInfo[dwIndex]
			dwsCodecEntries[genericDataWordName] = cr.parsed.EncoderDefs[WrapItemType(contractName, dwName, true)]
			break
		}
		if !founDW {
			return nil, nil, fmt.Errorf("failed to find data word with index %d for event: %q, its either out of bounds or can't be searched for", dwIndex, eventName)
		}
	}
	return dwsInfo, dwsCodecEntries, nil
}

// getEventTopicsInput returns codec entry for expected incoming event topic and the modifier to be applied to the params.
func (cr *chainReader) getEventTopicsInput(itemType string, inputMod codec.ModifiersConfig) (types.CodecEntry, codec.Modifier, error) {
	inputInfo := cr.parsed.EncoderDefs[itemType]
	// TODO can this be simplified? Isn't this same as inputInfo.Modifier()? BCI-3909
	inMod, err := inputMod.ToModifier(codec.DecoderHooks...)
	if err != nil {
		return nil, nil, err
	}

	// initialize the modification
	if _, err = inMod.RetypeToOffChain(reflect.PointerTo(inputInfo.CheckedType()), ""); err != nil {
		return nil, nil, err
	}

	return inputInfo, inMod, nil
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

	cr.parsed.EncoderDefs[WrapItemType(contractName, itemType, true)] = input
	return nil
}

func (cr *chainReader) addDecoderDef(contractName, itemType string, outputs abi.Arguments, outputModifications commoncodec.ModifiersConfig) error {
	mod, err := outputModifications.ToModifier(codec.DecoderHooks...)
	if err != nil {
		return err
	}
	output := types.NewCodecEntry(outputs, nil, mod)
	cr.parsed.DecoderDefs[read.WrapItemType(contractName, itemType, false)] = output
	return output.Init()
}

// getEventTypes returns abi args where indexed flag is set to false because we expect caller to filter with params that aren't hashed,
// codecEntry where expected onchain types are set, for e.g. indexed topics of type string or uint8[32] array are expected as common.Hash onchain,
// and unIndexed data info in form of evm indexed 32 byte data words.
func getEventTypes(event abi.Event) ([]abi.Argument, types.CodecEntry, eventDataWords) {
	indexedAsUnIndexedTypes := make([]abi.Argument, 0, types.MaxTopicFields)
	indexedTypes := make([]abi.Argument, 0, len(event.Inputs))
	dataWords := eventDataWords{}
	hadDynamicType := false
	dwIndex := 1

	fmt.Println("event inputs ", event.Inputs)
	for _, input := range event.Inputs {
		if !input.Indexed {
			// there are some cases where we can calculate the exact data word index even if there was a dynamic type before, but it is complex and probably not needed.
			if input.Type.T == abi.TupleTy || input.Type.T == abi.SliceTy || input.Type.T == abi.StringTy || input.Type.T == abi.BytesTy {
				hadDynamicType = true
			}
			if hadDynamicType {
				continue
			}

			fmt.Println("add dataWord ", input.Name)
			dataWords = append(dataWords, dataWordInfo{
				index: uint8(dwIndex),
				typ:   abi.Argument{Name: input.Name, Type: input.Type},
			})
			dwIndex++
			continue
		}

		indexedAsUnIndexed := input
		indexedAsUnIndexed.Indexed = false
		// when presenting the filter off-chain, the caller will provide the unHashed version of the input and CR will hash topics when needed.
		indexedAsUnIndexedTypes = append(indexedAsUnIndexedTypes, indexedAsUnIndexed)
		indexedTypes = append(indexedTypes, input)
	}

	return indexedAsUnIndexedTypes, types.NewCodecEntry(indexedTypes, nil, nil), dataWords
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

// eventDataWords maps event hashes to an ordered list of data word info.
// if a data word has a preceding data word that is dynamic, the exact data word index can't be known in advance and has to be calculated using the actual log data.
type eventDataWords []dataWordInfo

// For e.g. first evm data word(32bytes) of USDC log event is value so the key can be called value
// dataWordInfo contains all the information about a single evm Data word.
type dataWordInfo struct {
	// Name can be a nested field name, e.g. "foo.bar.baz"
	index uint8
	typ   abi.Argument
}
