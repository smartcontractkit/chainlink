package migrations

import (
	"gorm.io/gorm"
)

const up42 = `
CREATE TYPE job_proposal_status AS ENUM ('pending', 'approved', 'rejected');
CREATE TABLE job_proposals (
	id BIGSERIAL PRIMARY KEY,
	spec TEXT NOT NULL,
	status job_proposal_status NOT NULL,
	job_id uuid REFERENCES jobs (external_job_id) DEFERRABLE INITIALLY IMMEDIATE,
	feeds_manager_id int NOT NULL,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone NOT NULL,
	CONSTRAINT fk_feeds_manager FOREIGN KEY(feeds_manager_id) REFERENCES feeds_managers(id) DEFERRABLE INITIALLY IMMEDIATE,
	CONSTRAINT chk_job_proposals_status_fsm CHECK (
		(status = 'pending' AND job_id IS NULL) OR
		(status = 'approved' AND job_id IS NOT NULL) OR
		(status = 'rejected' AND job_id IS NULL)
	)
);
CREATE UNIQUE INDEX idx_job_proposals_job_id on job_proposals (job_id);
CREATE INDEX idx_job_proposals_feeds_manager_id on job_proposals (feeds_manager_id);
`
const down42 = `
	DROP TABLE job_proposals
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0042_create_job_proposals",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up42).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down42).Error
		},
	})
}
