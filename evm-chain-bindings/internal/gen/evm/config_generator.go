package evm

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/smart-contract-spec/internal/utils"
	"math/big"
)

func GenerateChainReaderChainWriterConfig(contracts map[string]CodeDetails) (types.ChainReaderConfig, types.ChainWriterConfig, error) {
	chainContractReaders := map[string]types.ChainContractReader{}
	chainWriterContractConfig := map[string]*types.ContractConfig{}

	for contractName, contractDetails := range contracts {
		chainReaderDefinitions := map[string]*types.ChainReaderDefinition{}
		chainWriterDefinitions := map[string]*types.ChainWriterDefinition{}

		for _, function := range contractDetails.Functions {
			if function.RequiresTransaction {
				chainWriterDefinitions[utils.CapitalizeFirstLetter(function.Name)] = &types.ChainWriterDefinition{
					ChainSpecificName: function.SolidityName,
					//TODO review proper value
					Checker: "simulate",
				}
			} else {
				chainReaderDefinitions[utils.CapitalizeFirstLetter(function.Name)] = &types.ChainReaderDefinition{
					//TODO: properly configure
					CacheEnabled:      false,
					ChainSpecificName: function.SolidityName,
					//TODO: properly configure
					InputModifications: nil,
					//TODO: properly configure
					OutputModifications: nil,
					//TODO: properly configure
					EventDefinitions: nil,
					//TODO: properly configure
					ConfidenceConfirmations: nil,
				}
			}
		}

		contractReader := types.ChainContractReader{
			ContractABI: contractDetails.ABI,
			Configs:     chainReaderDefinitions,
		}

		contractConfig := types.ContractConfig{
			ContractABI: contractDetails.ABI,
			Configs:     chainWriterDefinitions,
		}

		chainContractReaders[utils.CapitalizeFirstLetter(contractName)] = contractReader
		chainWriterContractConfig[utils.CapitalizeFirstLetter(contractName)] = &contractConfig
	}

	chainReaderConfig := types.ChainReaderConfig{
		Contracts: chainContractReaders,
	}

	maxGasPrice := assets.NewWei(big.NewInt(10000000))
	chainWriterConfig := types.ChainWriterConfig{
		Contracts: chainWriterContractConfig,
		//TODO fix, assets.Wei cannot be treated as a generic type with public fields
		MaxGasPrice: maxGasPrice,
	}

	return chainReaderConfig, chainWriterConfig, nil
}

//func getParamModifications(params []code.Field, input bool) (codec.ModifiersConfig, error) {
//	//TODO implement missing modifiers
//	modifiersConfig := codec.ModifiersConfig{}
//	renameFields := map[string]string{}
//	for _, param := range params {
//		if param.Name != param.GoName {
//			if input {
//				renameFields[param.GoName] = param.Name
//			} else {
//				renameFields[param.Name] = param.GoName
//			}
//		}
//	}
//	if len(renameFields) != 0 {
//		modifiersConfig = append(modifiersConfig, &codec.RenameModifierConfig{renameFields})
//	}
//	return modifiersConfig, nil
//}
