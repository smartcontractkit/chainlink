package evm

import (
	"context"
	"fmt"
	"reflect"
	"strings"

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
	parsed           *parsedTypes
	codec            commontypes.RemoteCodec
	commonservices.StateMachine
}

// NewChainReaderService is a constructor for ChainReader, returns nil if there is any error
func NewChainReaderService(lggr logger.Logger, lp logpoller.LogPoller, chain legacyevm.Chain, config types.ChainReaderConfig) (ChainReaderService, error) {
	cr := &chainReader{
		lggr:             lggr.Named("ChainReader"),
		lp:               lp,
		client:           chain.Client(),
		contractBindings: contractBindings{},
		parsed:           &parsedTypes{encoderDefs: map[string]types.CodecEntry{}, decoderDefs: map[string]types.CodecEntry{}},
	}

	var err error
	if err = cr.init(config.Contracts); err != nil {
		return nil, err
	}

	if cr.codec, err = cr.parsed.toCodec(); err != nil {
		return nil, err
	}

	err = cr.contractBindings.ForEach(func(b readBinding) error {
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

func (cr *chainReader) Bind(_ context.Context, bindings []commontypes.BoundContract) error {
	return cr.contractBindings.Bind(bindings)
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

func (cr *chainReader) Start(_ context.Context) error {
	return cr.StartOnce("ChainReader", func() error {
		return cr.contractBindings.ForEach(readBinding.Register)
	})
}

func (cr *chainReader) Close() error {
	return cr.StopOnce("ChainReader", func() error {
		return cr.contractBindings.ForEach(readBinding.Unregister)
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

	cr.contractBindings.AddReadBinding(contractName, eventName, &eventBinding{
		contractName:  contractName,
		eventName:     eventName,
		lp:            cr.lp,
		hash:          event.ID,
		inputInfo:     inputInfo,
		inputModifier: inputModifier,
		topicInfo:     topicInfo,
		id:            wrapItemType(contractName, eventName, false) + uuid.NewString(),
	})

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
