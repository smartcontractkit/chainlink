-- +goose Up
-- +goose StatementBegin

ALTER TABLE upkeep_registrations ALTER COLUMN upkeep_id TYPE numeric(78,0);

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
perform_upkeep_tx        [type=ethtx
                          minConfirmations=0
                          to="$(jobSpec.contractAddress)"
                          from="[$(jobSpec.fromAddress)]"
                          evmChainID="$(jobSpec.evmChainID)"
                          data="$(encode_perform_upkeep_tx)"
                          gasLimit="$(jobSpec.performUpkeepGasLimit)"
                          txMeta="{\"jobID\":$(jobSpec.jobID),\"upkeepID\":$(jobSpec.prettyID)}"]
encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> perform_upkeep_tx'
WHERE id IN (
    SELECT pipeline_spec_id
    FROM jobs
    WHERE type = 'keeper' AND schema_version = 3
);

UPDATE jobs
SET schema_version = 4
WHERE type = 'keeper' AND schema_version = 3;

-- uint256_to_bit converts a uint256 to a bit string
CREATE OR REPLACE FUNCTION uint256_to_bit(num NUMERIC)
  RETURNS BIT VARYING AS $$
DECLARE
  -- 1 + largest positive INT --
  max_int32 NUMERIC := '4294967296' :: NUMERIC(10);
  result BIT VARYING;
BEGIN
  ASSERT num <= 115792089237316195423570985008687907853269984665640564039457584007913129639935 AND num >= 0, "num outside uint256 range";
  -- break num into 32 bit chunks
  WITH chunks (exponent, chunk) AS (
    SELECT
      exponent,
      floor(num::NUMERIC(178,100) / (max_int32 ^ exponent) % max_int32)::BIGINT from generate_series(0,7) exponent
  )
  -- concat 32 bit chunks together
  SELECT bit_or(chunk::bit(256) << (32*(exponent)))
  FROM chunks INTO result;
  RETURN result;
END;
$$ LANGUAGE plpgsql;

-- least_significant selects the least significant n bits of a bit string
CREATE OR REPLACE FUNCTION least_significant(bits BIT VARYING, n integer)
  RETURNS BIT VARYING AS $$
BEGIN
  ASSERT length(bits) >= n, "slice is larger than input";
  RETURN substring(bits from length(bits) - n + 1 for n);
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE upkeep_registrations ALTER COLUMN upkeep_id TYPE bigint;

UPDATE jobs
SET schema_version = 3
WHERE type = 'keeper' AND schema_version = 4;

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
perform_upkeep_tx        [type=ethtx
                          minConfirmations=0
                          to="$(jobSpec.contractAddress)"
                          from="[$(jobSpec.fromAddress)]"
                          evmChainID="$(jobSpec.evmChainID)"
                          data="$(encode_perform_upkeep_tx)"
                          gasLimit="$(jobSpec.performUpkeepGasLimit)"
                          txMeta="{\"jobID\":$(jobSpec.jobID),\"upkeepID\":$(jobSpec.upkeepID)}"]
encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> perform_upkeep_tx'
WHERE id IN (
    SELECT pipeline_spec_id
    FROM jobs
    WHERE type = 'keeper' AND schema_version = 3
);


-- +goose StatementEnd
