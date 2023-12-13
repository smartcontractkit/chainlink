package evm

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

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
	contractID common.Address
	lp         logpoller.LogPoller
}

// NewChainReaderService constructor for ChainReader
func NewChainReaderService(lggr logger.Logger, lp logpoller.LogPoller, contractID common.Address, config types.ChainReaderConfig) (*chainReader, error) {
	if err := validateChainReaderConfig(config); err != nil {
		return nil, fmt.Errorf("%w: %w", commontypes.ErrInvalidConfig, err)
	}

	// TODO BCF-2814 implement initialisation of chain reading definitions and pass them into chainReader
	return &chainReader{lggr.Named("ChainReader"), contractID, lp}, nil
}

func (cr *chainReader) Name() string { return cr.lggr.Name() }

func (cr *chainReader) initialize() error {
	// Initialize chain reader, start cache polling loop, etc.
	return nil
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

func (cr *chainReader) GetLatestValue(ctx context.Context, bc commontypes.BoundContract, method string, params any, returnVal any) error {
	return commontypes.UnimplementedError("Unimplemented method GetLatestValue called")
}

func validateChainReaderConfig(cfg types.ChainReaderConfig) error {
	if len(cfg.ChainContractReaders) == 0 {
		return fmt.Errorf("%w: no contract readers defined", commontypes.ErrInvalidConfig)
	}

	for contractName, chainContractReader := range cfg.ChainContractReaders {
		abi, err := abi.JSON(strings.NewReader(chainContractReader.ContractABI))
		if err != nil {
			return fmt.Errorf("invalid abi: %w", err)
		}

		for chainReadingDefinitionName, chainReaderDefinition := range chainContractReader.ChainReaderDefinitions {
			switch chainReaderDefinition.ReadType {
			case types.Method:
				err = validateMethods(abi, chainReaderDefinition)
			case types.Event:
				err = validateEvents(abi, chainReaderDefinition)
			default:
				return fmt.Errorf("%w: invalid chainreading definition read type: %d for contract: %q", commontypes.ErrInvalidConfig, chainReaderDefinition.ReadType, contractName)
			}
			if err != nil {
				return fmt.Errorf("%w: invalid chainreading definition: %q for contract: %q, err: %w", commontypes.ErrInvalidConfig, chainReadingDefinitionName, contractName, err)
			}
		}
	}

	return nil
}

func validateEvents(contractABI abi.ABI, chainReaderDefinition types.ChainReaderDefinition) error {
	event, methodExists := contractABI.Events[chainReaderDefinition.ChainSpecificName]
	if !methodExists {
		return fmt.Errorf("event: %s doesn't exist", chainReaderDefinition.ChainSpecificName)
	}

	var abiEventIndexedInputs []abi.Argument
	for _, eventInput := range event.Inputs {
		if eventInput.Indexed {
			abiEventIndexedInputs = append(abiEventIndexedInputs, eventInput)
		}
	}

	var chainReaderEventParams []string
	for chainReaderEventParam := range chainReaderDefinition.Params {
		chainReaderEventParams = append(chainReaderEventParams, chainReaderEventParam)
	}

	if !areChainReaderArgumentsValid(abiEventIndexedInputs, chainReaderEventParams) {
		var abiEventIndexedInputsNames []string
		for _, abiEventIndexedInput := range abiEventIndexedInputs {
			abiEventIndexedInputsNames = append(abiEventIndexedInputsNames, abiEventIndexedInput.Name)
		}
		return fmt.Errorf("params: [%s] don't match abi event indexed inputs: [%s]", strings.Join(chainReaderEventParams, ","), strings.Join(abiEventIndexedInputsNames, ","))
	}
	return nil
}

func validateMethods(abi abi.ABI, chainReaderDefinition types.ChainReaderDefinition) error {
	method, methodExists := abi.Methods[chainReaderDefinition.ChainSpecificName]
	if !methodExists {
		return fmt.Errorf("method: %q doesn't exist", chainReaderDefinition.ChainSpecificName)
	}

	var methodNames []string
	for methodName := range chainReaderDefinition.Params {
		methodNames = append(methodNames, methodName)
	}

	if !areChainReaderArgumentsValid(method.Inputs, methodNames) {
		var abiMethodInputs []string
		for _, input := range method.Inputs {
			abiMethodInputs = append(abiMethodInputs, input.Name)
		}
		return fmt.Errorf("params: [%s] don't match abi method inputs: [%s]", strings.Join(methodNames, ","), strings.Join(abiMethodInputs, ","))
	}

	return nil
}

func areChainReaderArgumentsValid(contractArgs []abi.Argument, chainReaderArgs []string) bool {
	for _, contractArg := range contractArgs {
		if !slices.Contains(chainReaderArgs, contractArg.Name) {
			return false
		}
	}

	return true
}
