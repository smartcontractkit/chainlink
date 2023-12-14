package evm

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
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
	lggr                       logger.Logger
	lp                         logpoller.LogPoller
	codec                      commontypes.RemoteCodec
	client                     evmclient.Client
	contractReadingDefinitions types.ContractReaders
	commonservices.StateMachine
}

// NewChainReaderService is a constructor for ChainReader, returns nil if there is any error
func NewChainReaderService(lggr logger.Logger, lp logpoller.LogPoller, b Bindings, chain legacyevm.Chain, config types.ChainReaderConfig) (ChainReaderService, error) {
	parsed := &parsedTypes{
		encoderDefs: map[string]*codecEntry{},
		decoderDefs: map[string]*codecEntry{},
	}

	if err := addTypes(config.ContractReaders, parsed); err != nil {
		return nil, err
	}

	contractReaders := config.ContractReaders
	if err := bindReadingDefinitionToContractAddresses(b, contractReaders); err != nil {
		return nil, err
	}

	c, err := parsed.toCodec()
	if err != nil {
		return nil, err
	}

	return &chainReader{
		lggr:                       lggr.Named("ChainReader"),
		lp:                         lp,
		codec:                      c,
		client:                     chain.Client(),
		contractReadingDefinitions: contractReaders,
	}, err
}

func (cr *chainReader) Name() string { return cr.lggr.Name() }

var _ commontypes.ContractTypeProvider = &chainReader{}

func (cr *chainReader) GetLatestValue(ctx context.Context, contractName, readName string, params any, returnVal any) error {
	contractAddress, err := cr.contractReadingDefinitions.GetReadingDefinitionContractAddress(contractName, readName)
	if err != nil {
		return err
	}

	readType, err := cr.contractReadingDefinitions.GetReadType(contractName, readName)
	if err != nil {
		return err
	}

	if readType == types.Event {
		return cr.getLatestValueFromLogPoller(ctx, contractAddress, contractName, readName, returnVal)
	}

	return cr.getLatestValueFromContract(ctx, contractAddress, contractName, readName, params, returnVal)
}

func (cr *chainReader) getLatestValueFromLogPoller(ctx context.Context, contractAddress common.Address, contractName, eventName string, returnVal any) error {
	eventHash, err := cr.contractReadingDefinitions.GetEventHash(contractName, eventName)
	if err != nil {
		return err
	}

	log, err := cr.lp.LatestLogByEventSigWithConfs(eventHash, contractAddress, logpoller.Finalized)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("%w: %w", commontypes.ErrNotFound, err)
		}
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}
	return cr.codec.Decode(ctx, log.Data, returnVal, wrapItemType(contractName, eventName, false))
}

func (cr *chainReader) getLatestValueFromContract(ctx context.Context, contractAddress common.Address, contractName, methodName string, params any, returnVal any) error {
	data, err := cr.codec.Encode(ctx, params, wrapItemType(contractName, methodName, true))
	if err != nil {
		return err
	}

	callMsg := ethereum.CallMsg{
		From: contractAddress,
		Data: data,
	}

	output, err := cr.client.CallContract(ctx, callMsg, nil)
	if err != nil {
		return err
	}

	return cr.codec.Decode(ctx, output, returnVal, wrapItemType(contractName, methodName, false))
}

func (cr *chainReader) Start(_ context.Context) error {
	return cr.StartOnce("ChainReader", func() error {
		for contractName, contractReadingDef := range cr.contractReadingDefinitions {
			for readName, _ := range contractReadingDef.ReadingDefinitions {
				contractAddress, err := cr.contractReadingDefinitions.GetReadingDefinitionContractAddress(contractName, readName)
				if err != nil {
					return err
				}
				if hash, err := cr.contractReadingDefinitions.GetEventHash(contractName, readName); err == nil {
					if err = cr.lp.RegisterFilter(logpoller.Filter{
						Name:      wrapItemType(contractName, readName, false),
						EventSigs: evmtypes.HashArray{hash},
						Addresses: evmtypes.AddressArray{contractAddress},
					}); err != nil {
						return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
					}
				}
			}
		}
		return nil
	})
}
func (cr *chainReader) Close() error {
	return cr.StopOnce("ChainReader", func() error {
		for contractName, contractReadingDef := range cr.contractReadingDefinitions {
			for readName, _ := range contractReadingDef.ReadingDefinitions {
				if _, err := cr.contractReadingDefinitions.GetEventHash(contractName, readName); err == nil {
					if err = cr.lp.UnregisterFilter(wrapItemType(contractName, readName, false)); err != nil {
						return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
					}
				}
			}
		}
		return nil
	})
}

func (cr *chainReader) Ready() error { return nil }
func (cr *chainReader) HealthReport() map[string]error {
	return map[string]error{cr.Name(): nil}
}

func bindReadingDefinitionToContractAddresses(b Bindings, contractReadingDefinitions types.ContractReaders) (err error) {
	for contractName, contractBindings := range b {
		for readName, contractAddress := range contractBindings {
			readingDefinition, err := contractReadingDefinitions.GetReadingDefinition(contractName, readName)
			if err != nil {
				return err
			}
			readingDefinition.ContractAddress = contractAddress
		}
	}
	return nil
}

func (cr *chainReader) CreateContractType(contractName, methodName string, forEncoding bool) (any, error) {
	return cr.codec.CreateType(wrapItemType(contractName, methodName, forEncoding), forEncoding)
}

func addEvent(contractName, eventName string, contractABI abi.ABI, chainReaderDefinition types.ReadingDefinition, parsed *parsedTypes) error {
	event, methodExists := contractABI.Events[chainReaderDefinition.ChainSpecificName]
	if !methodExists {
		return fmt.Errorf("%w: method %s doesn't exist", commontypes.ErrInvalidConfig, chainReaderDefinition.ChainSpecificName)
	}

	return addDecoderDef(contractName, eventName, event.Inputs, parsed, chainReaderDefinition)
}

func addMethod(
	contractName, methodName string, abi abi.ABI, chainReaderDefinition types.ReadingDefinition, parsed *parsedTypes) error {
	abiMethod, methodExists := abi.Methods[chainReaderDefinition.ChainSpecificName]
	if !methodExists {
		return fmt.Errorf("method: %q doesn't exist", chainReaderDefinition.ChainSpecificName)
	}

	if err := addEncoderDef(contractName, methodName, abiMethod, parsed, chainReaderDefinition); err != nil {
		return err
	}

	return addDecoderDef(contractName, methodName, abiMethod.Outputs, parsed, chainReaderDefinition)
}

func addEncoderDef(contractName, methodName string, method abi.Method, parsed *parsedTypes, chainReaderDefinition types.ReadingDefinition) error {
	// ABI.Pack prepends the method.ID to the encodings, we'll need the encoder to do the same.
	input := &codecEntry{Args: method.Inputs, encodingPrefix: method.ID}

	if err := input.Init(); err != nil {
		return err
	}

	inputMod, err := chainReaderDefinition.InputModifications.ToModifier(evmDecoderHooks...)
	if err != nil {
		return err
	}
	input.mod = inputMod
	parsed.encoderDefs[wrapItemType(contractName, methodName, true)] = input
	return nil
}

func addDecoderDef(contractName, readName string, outputs abi.Arguments, parsed *parsedTypes, def types.ReadingDefinition) error {
	output := &codecEntry{Args: outputs}
	mod, err := def.OutputModifications.ToModifier(evmDecoderHooks...)
	if err != nil {
		return err
	}
	output.mod = mod
	parsed.decoderDefs[wrapItemType(contractName, readName, false)] = output
	return output.Init()
}

func addTypes(contractReaders map[string]types.ContractReader, parsed *parsedTypes) error {
	for contractName, contractReader := range contractReaders {
		contractAbi, err := abi.JSON(strings.NewReader(contractReader.ContractABI))
		if err != nil {
			return err
		}

		for typeName, chainReaderDefinition := range contractReader.ReadingDefinitions {
			switch chainReaderDefinition.ReadType {
			case types.Method:
				err = addMethod(contractName, typeName, contractAbi, chainReaderDefinition, parsed)
			case types.Event:
				err = addEvent(contractName, typeName, contractAbi, chainReaderDefinition, parsed)
			default:
				return fmt.Errorf(
					"%w: invalid chain reader definition read type: %d",
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

func wrapItemType(contractName, methodName string, isParams bool) string {
	if isParams {
		return fmt.Sprintf("params.%s.%s", contractName, methodName)
	}
	return fmt.Sprintf("return.%s.%s", contractName, methodName)
}
