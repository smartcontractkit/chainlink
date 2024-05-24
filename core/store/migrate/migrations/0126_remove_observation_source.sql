-- +goose Up
UPDATE pipeline_specs
SET dot_dag_source = ''
WHERE id IN (
    SELECT pipeline_spec_id
    FROM jobs
    WHERE type = 'keeper'
);

-- +goose Down
UPDATE pipeline_specs
SET dot_dag_source = 'encode_check_upkeep_tx   [type=ethabiencode
                          abi="checkUpkeep(uint256 id, address from)"
                          data="{\"id\":$(jobSpec.upkeepID),\"from\":$(jobSpec.fromAddress)}"]
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
                          data="{\"id\": $(jobSpec.upkeepID),\"performData\":$(decode_check_upkeep_tx.performData)}"]
simulate_perform_upkeep_tx  [type=ethcall
                          extractRevertReason=true
                          evmChainID="$(jobSpec.evmChainID)"
                          contract="$(jobSpec.contractAddress)"
                          from="$(jobSpec.fromAddress)"
                          gas="$(jobSpec.performUpkeepGasLimit)"
                          data="$(encode_perform_upkeep_tx)"]
decode_check_perform_tx  [type=ethabidecode
                          abi="bool success"]
check_success            [type=conditional
                          failEarly=true
                          data="$(decode_check_perform_tx.success)"]
perform_upkeep_tx        [type=ethtx
                          minConfirmations=0
                          to="$(jobSpec.contractAddress)"
                          from="[$(jobSpec.fromAddress)]"
                          evmChainID="$(jobSpec.evmChainID)"
                          data="$(encode_perform_upkeep_tx)"
                          gasLimit="$(jobSpec.performUpkeepGasLimit)"
                          txMeta="{\"jobID\":$(jobSpec.jobID),\"upkeepID\":$(jobSpec.prettyID)}"]
encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> simulate_perform_upkeep_tx -> decode_check_perform_tx -> check_success -> perform_upkeep_tx'
WHERE id IN (
    SELECT pipeline_spec_id
    FROM jobs
    WHERE type = 'keeper'
);
