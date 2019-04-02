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

router.get('/job_runs/:id', async (req, res) => {
  const id = req.params.id
  const params = { where: { id }, relations: ['taskRuns'] }
  const jobRun = await getDb().manager.findOne(JobRun, params)

  if (jobRun) {
    return res.send(jobRun)
  }

  return res.sendStatus(404)
})

export default router
