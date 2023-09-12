-- +goose Up

ALTER TABLE "s4".shared ALTER COLUMN version TYPE BIGINT;

-- +goose Down

ALTER TABLE "s4".shared ALTER COLUMN version TYPE INT USING version::integer;
