import { getDb } from '../database'
import { JobRun, search } from '../entity/JobRun'
import { Router, Request, Response } from 'express'
import { Option } from 'prelude-ts'

const router = Router()

router.get('/job_runs', async (req: Request, res: Response) => {
  const jobRuns = await search(getDb(), Option.of(req.query.query))
  return res.send(jobRuns)
})

router.get('/job_runs/:id', async (req: Request, res: Response) => {
  const runId = req.params.id
  const params = { where: { runId }, relations: ['taskRuns'] }
  const jobRun = await getDb().manager.findOne(JobRun, params)

  if (jobRun) {
    return res.send(jobRun)
  }

  return res.sendStatus(404)
})

export default router
