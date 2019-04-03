import { getDb } from '../database'
import { JobRun, search } from '../entity/JobRun'
import { Router, Request, Response } from 'express'

const router = Router()

router.get('/job_runs', async (req: Request, res: Response) => {
  const jobRuns = await search(getDb(), req.query.query)
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

export default router
