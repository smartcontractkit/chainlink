package job

import "github.com/smartcontractkit/chainlink/core/store/models"

func GetORMAdvisoryLockClassID(oi ORM) int32 {
	return oi.(*orm).advisoryLockClassID
}

func GetORMClaimedJobs(oi ORM) (claimedJobs []models.JobSpecV2) {
	o := oi.(*orm)
	o.claimedJobsMu.RLock()
	defer o.claimedJobsMu.RUnlock()
	claimedJobs = make([]models.JobSpecV2, 0)
	for _, job := range o.claimedJobs {
		claimedJobs = append(claimedJobs, job)
	}
	return claimedJobs
}

func GetORMClaimedJobIDs(oi ORM) (ids []int32) {
	for _, j := range GetORMClaimedJobs(oi) {
		ids = append(ids, j.ID)
	}
	return
}

func SetORMClaimedJobs(oi ORM, jobs []models.JobSpecV2) {
	o := oi.(*orm)
	var claimedJobs = make(map[int32]models.JobSpecV2)
	for _, job := range jobs {
		claimedJobs[job.ID] = job
	}

	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()
	o.claimedJobs = claimedJobs
}
