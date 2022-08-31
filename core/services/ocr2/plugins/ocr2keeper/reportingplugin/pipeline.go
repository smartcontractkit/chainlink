package reportingplugin

const queryObservationSource = `
    encode_check_upkeep_tx   [type=ethabiencode
                              abi="checkUpkeep(uint256 id)"
                              data="{\"id\":$(jobSpec.upkeepID)}"]
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
                              abi="bool upkeepNeeded, bytes memory performData, uint8 upkeepFailureReason, uint256 gasUsed"]
    encode_perform_upkeep_tx [type=ethabiencode
                              abi="performUpkeep(uint256 id, bytes calldata performData)"
                              data="{\"id\": $(jobSpec.upkeepID),\"performData\":$(decode_check_upkeep_tx.performData)}"]
    encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx
`
