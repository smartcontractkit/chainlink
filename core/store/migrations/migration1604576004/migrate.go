package migration1604576004

import "github.com/jinzhu/gorm"

// Migrate makes key deletion into a soft delete.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	---
	--- Notify the Chainlink node when a job is deleted
	---

	CREATE OR REPLACE FUNCTION notifyJobDeleted() RETURNS TRIGGER AS $_$
	BEGIN
		PERFORM pg_notify('delete_from_jobs', OLD.id::text);
		RETURN OLD;
	END
	$_$ LANGUAGE 'plpgsql';

	CREATE TRIGGER notify_job_deleted
	AFTER DELETE ON jobs
	FOR EACH ROW EXECUTE PROCEDURE notifyJobDeleted();
    `).Error
}
