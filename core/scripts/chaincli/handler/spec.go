package handler

import "fmt"

// KeeperSpecParams struct
type KeeperSpecParams struct {
	Name                     string
	ContractAddress          string
	FromAddress              string
	EvmChainID               int
	MinIncomingConfirmations int
}

// KeeperSpec struct
type KeeperSpec struct {
	KeeperSpecParams
	toml string
}

// Toml to TOML
func (os KeeperSpec) Toml() string {
	return os.toml
}

// GenerateKeeperSpec generate repeatable keeper job spec
func GenerateKeeperSpec(params KeeperSpecParams) KeeperSpec {
	template := `
type            		 	= "keeper"
schemaVersion   		 	= 3
name            		 	= "%s"
contractAddress 		 	= "%s"
fromAddress     		 	= "%s"
evmChainID      		 	= %d
minIncomingConfirmations	= %d


observationSource = """
encode_check_upkeep_tx   [type=ethabiencode
                          abi="checkUpkeep(uint256 id, address from)"
                          data="{\\"id\\":$(jobSpec.upkeepID),\\"from\\":$(jobSpec.fromAddress)}"]
check_upkeep_tx          [type=ethcall
                          failEarly=true
                          extractRevertReason=true
                          evmChainID="$(jobSpec.evmChainID)"
                          contract="$(jobSpec.contractAddress)"
                          gas="$(jobSpec.checkUpkeepGasLimit)"
                          gasPrice="$(jobSpec.gasPrice)"
                          gasTipCap="$(jobSpec.gasTipCap)"
                          gasFeeCap="$(jobSpec.gasFeeCap)"
                          data="$(encode_check_upkeep_tx)"]
decode_check_upkeep_tx   [type=ethabidecode
                          abi="bytes memory performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth"]
encode_perform_upkeep_tx [type=ethabiencode
                          abi="performUpkeep(uint256 id, bytes calldata performData)"
                          data="{\\"id\\": $(jobSpec.upkeepID),\\"performData\\":$(decode_check_upkeep_tx.performData)}"]
perform_upkeep_tx        [type=ethtx
                          minConfirmations=0
                          to="$(jobSpec.contractAddress)"
                          from="[$(jobSpec.fromAddress)]"
                          evmChainID="$(jobSpec.evmChainID)"
                          data="$(encode_perform_upkeep_tx)"
                          gasLimit="$(jobSpec.performUpkeepGasLimit)"
                          txMeta="{\\"jobID\\":$(jobSpec.jobID)}"]
encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> perform_upkeep_tx
"""
`
	return KeeperSpec{
		KeeperSpecParams: params,
		toml:             fmt.Sprintf(template, params.Name, params.ContractAddress, params.FromAddress, params.EvmChainID, params.MinIncomingConfirmations),
	}
}
