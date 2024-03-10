-- +goose Up

-- +goose StatementBegin
DROP FUNCTION IF EXISTS PUBLIC.check_terra_msg_state_transition;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

CREATE OR REPLACE FUNCTION PUBLIC.check_terra_msg_state_transition() RETURNS TRIGGER AS $$
DECLARE
state_transition_map jsonb := json_build_object(
        'unstarted', json_build_object('errored', true, 'started', true),
        'started', json_build_object('errored', true, 'broadcasted', true),
        'broadcasted', json_build_object('errored', true, 'confirmed', true));
BEGIN
    IF NOT state_transition_map ? OLD.state THEN
        RAISE EXCEPTION 'Invalid from state %. Valid from states %', OLD.state, state_transition_map;
END IF;
    IF NOT state_transition_map->OLD.state ? NEW.state THEN
        RAISE EXCEPTION 'Invalid state transition from % to %. Valid to states %', OLD.state, NEW.state, state_transition_map->OLD.state;
END IF;
RETURN NEW;
END
$$ LANGUAGE plpgsql;

-- +goose StatementEnd