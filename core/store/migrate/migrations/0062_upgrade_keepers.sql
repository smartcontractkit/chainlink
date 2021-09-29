-- +goose Up
UPDATE pipeline_specs
SET dot_dag_source = 'encode_check_upkeep_tx   [type=ethabiencode
                          abi="checkUpkeep(uint256 id, address from)"
                          data="{\"id\":$(jobSpec.upkeepID),\"from\":$(jobSpec.fromAddress)}"]
check_upkeep_tx          [type=ethcall
                          failEarly=true
                          extractRevertReason=true
                          contract="$(jobSpec.contractAddress)"
                          gas="$(jobSpec.checkUpkeepGasLimit)"
                          data="$(encode_check_upkeep_tx)"]
decode_check_upkeep_tx   [type=ethabidecode
                          abi="bytes memory performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth"]
encode_perform_upkeep_tx [type=ethabiencode
                          abi="performUpkeep(uint256 id, bytes calldata performData)"
                          data="{\"id\": $(jobSpec.upkeepID),\"performData\":$(decode_check_upkeep_tx.performData)}"]
perform_upkeep_tx        [type=ethtx
                          minConfirmations=0
                          to="$(jobSpec.contractAddress)"
                          data="$(encode_perform_upkeep_tx)"
                          gasLimit="$(jobSpec.performUpkeepGasLimit)"
                          txMeta="{\"jobID\":$(jobSpec.jobID)}"]
encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> perform_upkeep_tx'
WHERE id IN (
    SELECT pipeline_spec_id
    FROM jobs
    WHERE type = 'keeper' AND schema_version = 1
);

UPDATE jobs
SET schema_version = 2
WHERE type = 'keeper' AND schema_version = 1;

-- +goose Down
UPDATE jobs
SET schema_version = 1
WHERE type = 'keeper' AND schema_version = 2;

UPDATE pipeline_specs
SET dot_dag_source = ''
WHERE id IN (
    SELECT pipeline_spec_id
    FROM jobs
    WHERE type = 'keeper' AND schema_version = 1
);