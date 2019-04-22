package migration1554131520

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// Migrate adds the run request table
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&models.RunRequest{}).Error; err != nil {
		return err
	}

	if err := tx.AutoMigrate(&jobRun{}).Error; err != nil {
		return err
	}

	if err := backfillRunRequests(tx, "runat", models.RunRequest{}); err != nil {
		return err
	}
	if err := backfillRunRequests(tx, "cron", models.RunRequest{}); err != nil {
		return err
	}
	if err := backfillRunRequests(tx, "web", models.RunRequest{}); err != nil {
		return err
	}
	txhash := common.HexToHash("0xeeeeeeeeeeeeeeee")
	if err := backfillRunRequests(tx, "ethlog", models.RunRequest{TxHash: &txhash}); err != nil {
		return err
	}
	requester := common.HexToAddress("0xdeadbeef")
	requestID := "BACKFILLED_FAKE"
	runlogrr := models.RunRequest{TxHash: &txhash, Requester: &requester, RequestID: &requestID}
	if err := backfillRunRequests(tx, "runlog", runlogrr); err != nil {
		return err
	}

	return nil
}

func backfillRunRequests(tx *gorm.DB, initrType string, rr models.RunRequest) error {
	results, err := runIdsFor(tx, initrType)
	if err != nil {
		return err
	}
	for _, jrid := range results {
		if err := replaceRunRequest(tx, jrid, rr); err != nil {
			return err
		}
	}
	return nil
}

func runIdsFor(tx *gorm.DB, initrType string) ([]string, error) {
	var results []string
	err := tx.Unscoped().
		Table("job_runs").
		Joins("inner join initiators on job_runs.initiator_id = initiators.id").
		Where("job_runs.run_request_id IS NULL AND initiators.type = ?", initrType).
		Pluck("job_runs.id", &results).Error
	return results, err
}

func replaceRunRequest(tx *gorm.DB, jrid string, rr models.RunRequest) error {
	if err := tx.Create(&rr).Error; err != nil {
		return err
	}

	if err := tx.Exec(`UPDATE "job_runs" SET run_request_id = ? WHERE id = ?`, rr.ID, jrid).Error; err != nil {
		return err
	}
	return nil
}

// jobRun private type here to capture and isolate the single change necessary
// to job run, insulating this migration from future changes to models.JobRun.
type jobRun struct {
	ID           string `gorm:"primary_key"`
	RunRequestID uint
}
