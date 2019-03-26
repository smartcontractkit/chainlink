import { getDb } from '../database'
import { JobRun } from '../entity/JobRun'
import { Router } from 'express'

const router = Router()
router.get('/job_runs', async (req, res) => {
  const searchQuery = req.query.query
  let params = {}
  if (searchQuery) {
    params = { where: { jobId: searchQuery } }
  }
  const jobRuns = await getDb().manager.find(JobRun, params)

  return res.send(jobRuns)
})

export default router
