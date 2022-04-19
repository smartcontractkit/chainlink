-- +goose Up
-- +goose StatementBegin

ALTER TABLE upkeep_registrations ALTER COLUMN upkeep_id TYPE numeric(78,0);

CREATE OR REPLACE FUNCTION uint256_to_bit(NUMERIC)
  RETURNS BIT VARYING AS $$
DECLARE
  num ALIAS FOR $1;
  -- 1 + largest positive BIGINT --
  max_bigint NUMERIC := '9223372036854775808' :: NUMERIC(19, 0);
  result BIT VARYING;
BEGIN
  ASSERT num <= 115792089237316195423570985008687907853269984665640564039457584007913129639935, "value larger than max uint256";
  WITH
      chunks (exponent, chunk) AS (
        SELECT
          exponent,
          floor((num / (max_bigint ^ exponent) :: NUMERIC(256, 20)) % max_bigint) :: BIGINT
        FROM generate_series(0, 5) exponent
    )
  SELECT bit_or(chunk :: BIT(256) :: BIT VARYING << (63 * (exponent))) :: BIT VARYING
  FROM chunks INTO result;
  RETURN result;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION least_significant(bits BIT VARYING, num integer)
  RETURNS BIT VARYING AS $$
BEGIN
  ASSERT length(bits) >= num, "slice is larger than input";
  RETURN substring(bits from length(bits) - num + 1 for num);
END;
$$ LANGUAGE plpgsql;


-- CREATE UNIQUE INDEX idx_upkeep_registrations_unique_upkeep_ids_per_keeper ON upkeep_registrations(registry_id, upkeep_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE upkeep_registrations ALTER COLUMN upkeep_id TYPE bigint;

-- +goose StatementEnd
