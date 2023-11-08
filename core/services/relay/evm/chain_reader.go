package evm

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

// newChainReader validates config and initializes chainReader, returns nil if there is any error.
func newChainReader(lggr logger.Logger, chain evm.Chain, ropts *types.RelayOpts) (*chainReader, error) {
	relayConfig, err := ropts.RelayConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed parsing RelayConfig: %w", err)
	}

	if relayConfig.ChainReader == nil {
		return nil, relaytypes.ErrorChainReaderUnsupported{}
	}

	if err = parseChainContractReadersABIs(relayConfig.ChainReader.ChainContractReaders); err != nil {
		return nil, err
	}

	if err = validateChainReaderConfig(*relayConfig.ChainReader); err != nil {
		return nil, fmt.Errorf("invalid ChainReader configuration: %w", err)
	}

	return &chainReader{lggr.Named("ChainReader"), chain.LogPoller(), relayConfig.ChainReader.ChainContractReaders, chain}, nil
}

func parseChainContractReadersABIs(chainContractReaders map[string]types.ChainContractReader) error {
	for key, chainContractReader := range chainContractReaders {
		parsedABI, err := abi.JSON(strings.NewReader(chainContractReader.ContractABI))
		if err != nil {
			return fmt.Errorf("falied to parse contract:%s abi:%s, err:%w", key, chainContractReader.ContractABI, err)
		}
		chainContractReader.ParsedContractABI = &parsedABI
		chainContractReaders[key] = chainContractReader
	}
	return nil
}

func validateChainReaderConfig(cfg types.ChainReaderConfig) (err error) {
	for contractName, chainContractReader := range cfg.ChainContractReaders {
		if chainContractReader.ParsedContractABI == nil {
			return fmt.Errorf("contract: %s ABI is not parsed", contractName)
		}

		for chainReadingDefinitionName, chainReaderDefinition := range chainContractReader.ChainReaderDefinitions {
			switch chainReaderDefinition.ReadType {
			case types.Method:
				err = validateMethods(*chainContractReader.ParsedContractABI, chainReaderDefinition)
			case types.Event:
				err = validateEvents(*chainContractReader.ParsedContractABI, chainReaderDefinition)
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

func validateEvents(contractABI abi.ABI, chainReaderDefinition types.ChainReaderDefinition) error {
	event, methodExists := contractABI.Events[chainReaderDefinition.ChainSpecificName]
	if !methodExists {
		return fmt.Errorf("method: %s doesn't exist", chainReaderDefinition.ChainSpecificName)
	}

	if !areChainReaderArgumentsValid(event.Inputs, chainReaderDefinition.ReturnValues) {
		var abiEventInputsNames []string
		for _, input := range event.Inputs {
			abiEventInputsNames = append(abiEventInputsNames, input.Name)
		}
		return fmt.Errorf("return values: [%s] don't match abi event inputs: [%s]", strings.Join(chainReaderDefinition.ReturnValues, ","), strings.Join(abiEventInputsNames, ","))
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

	if !areChainReaderArgumentsValid(method.Outputs, chainReaderDefinition.ReturnValues) {
		var abiMethodOutputs []string
		for _, input := range method.Outputs {
			abiMethodOutputs = append(abiMethodOutputs, input.Name)
		}
		return fmt.Errorf("return values: [%s] don't match abi method outputs: [%s]", strings.Join(chainReaderDefinition.ReturnValues, ","), strings.Join(abiMethodOutputs, ","))
	}

	return nil
}

func areChainReaderArgumentsValid(contractArgs []abi.Argument, chainReaderArgs []string) bool {
	for _, chArgName := range chainReaderArgs {
		found := false
		for _, contractArg := range contractArgs {
			if chArgName == contractArg.Name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (cr *chainReader) initialize() error {
	// Initialize chain reader, start cache polling loop, etc.
	return nil
}

type ChainReaderService interface {
	services.ServiceCtx
	relaytypes.ChainReader
}

type chainReader struct {
	lggr logger.Logger
	lp   logpoller.LogPoller
	// key being contract name
	chainContractReaders map[string]types.ChainContractReader
	chain                evm.Chain
}

// GetLatestValue calls given contract method and returns current value.
func (cr *chainReader) GetLatestValue(ctx context.Context, bc relaytypes.BoundContract, method string, params any, returnVal any) (err error) {
	chainContractReader := cr.chainContractReaders[bc.Name]
	chainReadingDefinition := chainContractReader.ChainReaderDefinitions[method]
	chainSpecificName := chainReadingDefinition.ChainSpecificName
	contractAddr := common.HexToAddress(bc.Address)

	if chainReadingDefinition.ReadType == types.Method {
		var callData []byte
		if params != nil {
			callData, err = chainContractReader.ParsedContractABI.Pack(chainSpecificName, params)
			if err != nil {
				return err
			}
		} else {
			callData, err = chainContractReader.ParsedContractABI.Pack(chainSpecificName)
			if err != nil {
				return err
			}
		}

		ethCallMsg := ethereum.CallMsg{
			From: common.Address{},
			To:   &contractAddr,
			Data: callData,
		}

		response, err := cr.chain.Client().CallContract(ctx, ethCallMsg, nil)
		if err != nil {
			return err
		}

		if err = chainContractReader.ParsedContractABI.UnpackIntoInterface(returnVal, chainSpecificName, response); err != nil {
			return errors.Wrap(err, "failed to unpack response")
		}
	} //else {
	//	event := chainContractReader.ParsedContractABI.Events[chainSpecificName]
	//	if err != nil {
	//		return err
	//	}
	//	for _, input := range event.Inputs {
	//		overrideVal := chainReadingDefinition.Params[input.Name]
	//		if overrideVal != nil {
	//			params = overrideVal
	//		}
	//	}
	//	cr.lp.IndexedLogs(event.ID, common.HexToAddress(), contractAddr, params)
	//}

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
func (cr *chainReader) Name() string { return cr.lggr.Name() }
