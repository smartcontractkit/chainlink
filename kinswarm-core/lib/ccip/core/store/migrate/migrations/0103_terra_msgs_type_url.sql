-- +goose Up
-- +goose StatementBegin
ALTER TABLE terra_msgs ADD COLUMN type text NOT NULL DEFAULT '/terra.wasm.v1beta1.MsgExecuteContract';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE terra_msgs DROP COLUMN type;
-- +goose StatementEnd
