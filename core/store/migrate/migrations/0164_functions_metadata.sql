-- +goose Up
ALTER TABLE ocr2dr_requests
ADD COLUMN aggregation_method int8 DEFAULT 1;

-- Populate aggregation_method with the default value of median = 1
UPDATE ocr2dr_requests
SET aggregation_method = 1 where aggregation_method is NULL

-- +goose Down
ALTER TABLE ocr2dr_requests
DROP COLUMN aggregation_method;