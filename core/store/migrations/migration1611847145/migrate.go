package migration1611847145

import "github.com/jinzhu/gorm"

// Migrate changes trigger to only notify on unfinished runs
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	CREATE OR REPLACE FUNCTION notifyPipelineRunStarted() RETURNS TRIGGER AS $_$
	BEGIN
		IF NEW.finished_at IS NULL THEN
			PERFORM pg_notify('pipeline_run_started', NEW.id::text);
		END IF;
		RETURN NEW;
	END
	$_$ LANGUAGE 'plpgsql';
	`).Error
}
