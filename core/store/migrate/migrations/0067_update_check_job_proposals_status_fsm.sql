-- +goose Up
-- +goose StatementBegin

ALTER TABLE job_proposals
DROP CONSTRAINT chk_job_proposals_status_fsm;

ALTER TABLE job_proposals
ADD CONSTRAINT chk_job_proposals_status_fsm CHECK (
	(status = 'pending' AND external_job_id IS NULL) OR
	(status = 'approved' AND external_job_id IS NOT NULL) OR
	(status = 'rejected' AND external_job_id IS NULL) OR
	(status = 'cancelled' AND external_job_id IS NULL)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE job_proposals
DROP CONSTRAINT chk_job_proposals_status_fsm;

ALTER TABLE job_proposals
ADD CONSTRAINT chk_job_proposals_status_fsm CHECK (
	(status = 'pending' AND external_job_id IS NULL) OR
	(status = 'approved' AND external_job_id IS NOT NULL) OR
	(status = 'rejected' AND external_job_id IS NULL)
);

-- +goose StatementEnd
