import { getDb } from '../database'
import { JobRun, search } from '../entity/JobRun'
import { Router, Request, Response } from 'express'

const router = Router()

const DEFAULT_PAGE = 1
const DEFAULT_SIZE = 10

router.get('/job_runs', async (req: Request, res: Response) => {
  const page = parseInt(req.query.page, 10) || DEFAULT_PAGE
  const size = parseInt(req.query.size, 10) || DEFAULT_SIZE
  const searchParams = {
    searchQuery: req.query.query,
    page: page,
    limit: size
  }
  const jobRuns = await search(getDb(), searchParams)
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
