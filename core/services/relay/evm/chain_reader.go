package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type ClosableChainReader interface {
	relaytypes.ChainReader
	Start(ctx context.Context) error
	Close() error
}

// NewChainReader is a constructor for ChainReader, returns nil if there is any error
func NewChainReader(lggr logger.Logger, chain evm.Chain, ropts *types.RelayOpts) (ClosableChainReader, error) {
	relayConfig, err := ropts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed parsing RelayConfig: %w", err)
	}

	if relayConfig.ChainReader == nil {
		return nil, relaytypes.ErrorChainReaderUnsupported{}
	}

	crConfig := *relayConfig.ChainReader

	parsed := &parsedTypes{
		encoderDefs: map[string]*CodecEntry{},
		decoderDefs: map[string]*CodecEntry{},
	}
	for k, v := range crConfig.ChainCodecConfigs {
		args := abi.Arguments{}
		if err = json.Unmarshal(([]byte)(v.TypeAbi), &args); err != nil {
			return nil, err
		}

		item := &CodecEntry{Args: args}
		if err = item.Init(); err != nil {
			return nil, err
		}

		parsed.encoderDefs[k] = item
		parsed.decoderDefs[k] = item
	}

	if err = addTypes(crConfig.ChainContractReaders, parsed); err != nil {
		return nil, err
	}

	enc := &encoder{Definitions: parsed.encoderDefs}
	dec := &decoder{Definitions: parsed.decoderDefs}

	return &chainReader{
		lggr:    lggr.Named("ChainReader"),
		lp:      chain.LogPoller(),
		Encoder: enc,
		Decoder: dec,
		client:  chain.Client(),
		types:   parsed,
	}, nil
}

func addTypes(chainContractReaders map[string]types.ChainContractReader, parsed *parsedTypes) error {
	for contractName, chainContractReader := range chainContractReaders {
		contractAbi, err := abi.JSON(strings.NewReader(chainContractReader.ContractABI))
		if err != nil {
			return err
		}

		for chainReadingDefinitionName, chainReaderDefinition := range chainContractReader.ChainReaderDefinitions {
			switch chainReaderDefinition.ReadType {
			case types.Method:
				err = addMethods(chainReadingDefinitionName, contractAbi, chainReaderDefinition, parsed)
			case types.Event:
				err = addEventTypes(chainReadingDefinitionName, contractAbi, chainReaderDefinition, parsed)
			default:
				return fmt.Errorf("invalid chain reader definition read type: %d", chainReaderDefinition.ReadType)
			}
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("invalid chain reader config for contract: %q chain reading definition: %q", contractName, chainReadingDefinitionName))
			}
		}
	}

	return nil
}

func addEventTypes(name string, contractABI abi.ABI, chainReaderDefinition types.ChainReaderDefinition, parsed *parsedTypes) error {
	event, methodExists := contractABI.Events[chainReaderDefinition.ChainSpecificName]
	if !methodExists {
		return fmt.Errorf("method: %s doesn't exist", chainReaderDefinition.ChainSpecificName)
	}

	if err := addOverrides(chainReaderDefinition, event.Inputs); err != nil {
		return err
	}

	return addDecoderDef(name, event.Inputs, parsed)
}

func addMethods(name string, abi abi.ABI, chainReaderDefinition types.ChainReaderDefinition, parsed *parsedTypes) error {
	method, methodExists := abi.Methods[chainReaderDefinition.ChainSpecificName]
	if !methodExists {
		return fmt.Errorf("method: %q doesn't exist", chainReaderDefinition.ChainSpecificName)
	}

	if err := addOverrides(chainReaderDefinition, method.Inputs); err != nil {
		return err
	}

	// ABI.Pack prepends the method.ID to the encodings, we'll need the encoder to do the same.
	input := &CodecEntry{Args: method.Inputs, encodingPrefix: method.ID}
	if err := input.Init(); err != nil {
		return err
	}

	parsed.encoderDefs[name] = input
	return addDecoderDef(name, method.Outputs, parsed)
}

func addDecoderDef(name string, outputs abi.Arguments, parsed *parsedTypes) error {
	output := &CodecEntry{Args: outputs}
	parsed.decoderDefs[name] = output
	return output.Init()
}

func (cr *chainReader) initialize() error {
	// Initialize chain reader, start cache polling loop, etc.
	return nil
}

func addOverrides(chainReaderDefinition types.ChainReaderDefinition, inputs abi.Arguments) error {
	// TODO add transforms to add params artificially
paramsLoop:
	for argName, param := range chainReaderDefinition.Params {
		// TODO add type check too
		_ = param
		for _, input := range inputs {
			if argName == input.Name {
				continue paramsLoop
			}
		}
		return fmt.Errorf("cannot find parameter %v in %v", argName, chainReaderDefinition.ChainSpecificName)
	}

	return nil
}

type ChainReaderService interface {
	services.ServiceCtx
	relaytypes.ChainReader
}

type chainReader struct {
	lggr logger.Logger
	lp   logpoller.LogPoller
	relaytypes.Encoder
	relaytypes.Decoder
	types  *parsedTypes
	client evmclient.Client
}

var _ relaytypes.RemoteCodec = &chainReader{}

func (cr *chainReader) GetLatestValue(ctx context.Context, bc relaytypes.BoundContract, method string, params any, returnVal any) error {
	data, err := cr.Encode(ctx, params, method)
	if err != nil {
		return err
	}

	address := common.HexToAddress(bc.Address)
	callMsg := ethereum.CallMsg{
		To:   &address,
		From: address,
		Data: data,
	}

	output, err := cr.client.CallContract(ctx, callMsg, nil)

	if err != nil {
		return err
	}

	return cr.Decode(ctx, output, returnVal, method)
}

func (cr *chainReader) Start(ctx context.Context) error {
	if err := cr.initialize(); err != nil {
		return fmt.Errorf("Failed to initialize ChainReader: %w", err)
	}
	return nil
}
func (cr *chainReader) Close() error { return nil }

func (cr *chainReader) Ready() error { return nil }
func (cr *chainReader) HealthReport() map[string]error {
	return map[string]error{cr.Name(): nil}
}
func (cr *chainReader) Name() string { return cr.lggr.Name() }

func (cr *chainReader) CreateType(itemType string, forEncoding bool) (any, error) {
	var itemTypes map[string]*CodecEntry
	if forEncoding {
		itemTypes = cr.types.encoderDefs
	} else {
		itemTypes = cr.types.decoderDefs
	}

	def, ok := itemTypes[itemType]
	if !ok {
		return nil, relaytypes.InvalidTypeError{}
	}

	return def.checkedType, nil
}

type parsedTypes struct {
	encoderDefs map[string]*CodecEntry
	decoderDefs map[string]*CodecEntry
}
