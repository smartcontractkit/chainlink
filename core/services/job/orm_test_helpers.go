package job

func GetORMAdvisoryLockClassID(oi ORM) int32 {
	return oi.(*orm).advisoryLockClassID
}

func GetORMClaimedJobs(oi ORM) (claimedJobs []Job) {
	o, _ := oi.(*orm)
	o.claimedJobsMu.RLock()
	defer o.claimedJobsMu.RUnlock()
	claimedJobs = make([]Job, 0)
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

func SetORMClaimedJobs(oi ORM, jobs []Job) {
	o, _ := oi.(*orm)
	var claimedJobs = make(map[int32]Job)
	for _, job := range jobs {
		claimedJobs[job.ID] = job
	}

	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()
	o.claimedJobs = claimedJobs
}
