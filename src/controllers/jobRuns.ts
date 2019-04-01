import { getDb } from '../database'
import { JobRun, search } from '../entity/JobRun'
import { Router, Request, Response } from 'express'

const router = Router()

router.get('/job_runs', async (req: Request, res: Response) => {
  const searchQuery = req.query.query
  let params = {}
  if (searchQuery) {
    params = { where: { jobId: searchQuery } }
  }
  const jobRuns = await getDb().manager.find(JobRun, params)

  return res.send(jobRuns)
})

router.get('/job_runs/:id', async (req: Request, res: Response) => {
  const id = req.params.id
  const params = { where: { id }, relations: ['taskRuns'] }
  const jobRun = await getDb().manager.findOne(JobRun, params)

  if (jobRun) {
    return res.send(jobRun)
  }

  return res.sendStatus(404)
})

router.get('/job_runs/search', async (req: Request, res: Response) => {
  const searchTokens = req.query.query.split(/\s+/)
  const jobRuns = await search(getDb(), searchTokens)
  return res.send(jobRuns)
})

export default router
