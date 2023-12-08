package evm

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

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
	lggr       logger.Logger
	lp         logpoller.LogPoller
	codec      commontypes.RemoteCodec
	client     evmclient.Client
	contractID common.Address
	events     map[string]common.Hash
	commonservices.StateMachine
}

// NewChainReaderService is a constructor for ChainReader, returns nil if there is any error
func NewChainReaderService(lggr logger.Logger, lp logpoller.LogPoller, contractID common.Address, chain legacyevm.Chain, config types.ChainReaderConfig) (ChainReaderService, error) {
	parsed := &parsedTypes{
		encoderDefs: map[string]*codecEntry{},
		decoderDefs: map[string]*codecEntry{},
	}

	events, err := addTypes(config.ChainContractReaders, parsed)
	if err != nil {
		return nil, err
	}

	c, err := parsed.toCodec()

	return &chainReader{
		lggr:       lggr.Named("ChainReader"),
		lp:         lp,
		codec:      c,
		client:     chain.Client(),
		contractID: contractID,
		events:     events,
	}, err
}

func (cr *chainReader) Name() string { return cr.lggr.Name() }

var _ commontypes.TypeProvider = &chainReader{}

func (cr *chainReader) GetLatestValue(ctx context.Context, bc commontypes.BoundContract, method string, params any, returnVal any) error {
	if hash, ok := cr.events[method]; ok {
		return cr.getLatestValueFromLogPoller(ctx, bc, method, hash, returnVal)
	}
	return cr.getLatestValueFromContract(ctx, bc, method, params, returnVal)
}

func (cr *chainReader) getLatestValueFromLogPoller(ctx context.Context, bc commontypes.BoundContract, method string, hash common.Hash, returnVal any) error {
	contractAddr := common.HexToAddress(bc.Address)
	log, err := cr.lp.LatestLogByEventSigWithConfs(hash, contractAddr, logpoller.Finalized)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("%w: %w", commontypes.ErrNotFound, err)
		}
		return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
	}
	return cr.codec.Decode(ctx, log.Data, returnVal, wrapItemType(method, false))
}

func (cr *chainReader) getLatestValueFromContract(ctx context.Context, bc commontypes.BoundContract, method string, params any, returnVal any) error {
	data, err := cr.codec.Encode(ctx, params, wrapItemType(method, true))
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

	return cr.codec.Decode(ctx, output, returnVal, wrapItemType(method, false))
}

func (cr *chainReader) Start(_ context.Context) error {
	return cr.StartOnce("ChainReader", func() error {
		for name, eventId := range cr.events {
			if err := cr.lp.RegisterFilter(logpoller.Filter{
				Name:      name,
				EventSigs: evmtypes.HashArray{eventId},
				Addresses: evmtypes.AddressArray{cr.contractID},
			}); err != nil {
				return fmt.Errorf("%w: %w", commontypes.ErrInternal, err)
			}
		}
		return nil
	})
}
func (cr *chainReader) Close() error {
	return nil
}

func (cr *chainReader) Ready() error { return nil }
func (cr *chainReader) HealthReport() map[string]error {
	return map[string]error{cr.Name(): nil}
}

func (cr *chainReader) CreateType(itemType string, forEncoding bool) (any, error) {
	return cr.codec.CreateType(wrapItemType(itemType, forEncoding), forEncoding)
}

func addEventTypes(name string, contractABI abi.ABI, chainReaderDefinition types.ChainReaderDefinition, parsed *parsedTypes) (common.Hash, error) {
	event, methodExists := contractABI.Events[chainReaderDefinition.ChainSpecificName]
	if !methodExists {
		return common.Hash{}, fmt.Errorf("method: %s doesn't exist", chainReaderDefinition.ChainSpecificName)
	}

	return event.ID, addDecoderDef(name, event.Inputs, parsed, chainReaderDefinition)
}

func addMethods(name string, abi abi.ABI, chainReaderDefinition types.ChainReaderDefinition, parsed *parsedTypes) error {
	method, methodExists := abi.Methods[chainReaderDefinition.ChainSpecificName]
	if !methodExists {
		return fmt.Errorf("method: %q doesn't exist", chainReaderDefinition.ChainSpecificName)
	}

	if err := addEncoderDef(name, method, parsed, chainReaderDefinition); err != nil {
		return err
	}

	return addDecoderDef(name, method.Outputs, parsed, chainReaderDefinition)
}

func addEncoderDef(name string, method abi.Method, parsed *parsedTypes, chainReaderDefinition types.ChainReaderDefinition) error {
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
	parsed.encoderDefs[wrapItemType(name, true)] = input
	return nil
}

func addDecoderDef(name string, outputs abi.Arguments, parsed *parsedTypes, def types.ChainReaderDefinition) error {
	output := &codecEntry{Args: outputs}
	mod, err := def.OutputModifications.ToModifier(evmDecoderHooks...)
	if err != nil {
		return err
	}
	output.mod = mod
	parsed.decoderDefs[wrapItemType(name, false)] = output
	return output.Init()
}

func addTypes(chainContractReaders map[string]types.ChainContractReader, parsed *parsedTypes) (map[string]common.Hash, error) {
	events := map[string]common.Hash{}
	for contractName, chainContractReader := range chainContractReaders {
		contractAbi, err := abi.JSON(strings.NewReader(chainContractReader.ContractABI))
		if err != nil {
			return nil, err
		}

		for chainReadingDefinitionName, chainReaderDefinition := range chainContractReader.ChainReaderDefinitions {
			switch chainReaderDefinition.ReadType {
			case types.Method:
				err = addMethods(chainReadingDefinitionName, contractAbi, chainReaderDefinition, parsed)
			case types.Event:
				var hash common.Hash
				hash, err = addEventTypes(chainReadingDefinitionName, contractAbi, chainReaderDefinition, parsed)
				events[chainReadingDefinitionName] = hash
			default:
				return nil, fmt.Errorf("invalid chain reader definition read type: %d", chainReaderDefinition.ReadType)
			}

			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("invalid chain reader config for contract: %q chain reading definition: %q", contractName, chainReadingDefinitionName))
			}
		}
	}

	return events, nil
}

func wrapItemType(itemType string, isParams bool) string {
	if isParams {
		return fmt.Sprintf("params.%s", itemType)
	}
	return fmt.Sprintf("return.%s", itemType)
}
