import { getDb } from '../database'
import { JobRun, present, search } from '../entity/JobRun'
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
  const db = await getDb()
  const jobRuns = await search(db, searchParams)
  return res.send(jobRuns.map(jr => present(jr)))
})

router.get('/job_runs/:id', async (req: Request, res: Response) => {
  const id = req.params.id
  const db = await getDb()
  const jobRun = await db
    .getRepository(JobRun)
    .createQueryBuilder('job_run')
    .leftJoinAndSelect('job_run.taskRuns', 'task_run')
    .leftJoinAndSelect('job_run.chainlinkNode', 'chainlink_node')
    .orderBy('job_run.createdAt, task_run.index', 'ASC')
    .where('job_run.id = :id', { id })
    .getOne()

  if (jobRun) {
    return res.send(jobRun)
  }

  return res.sendStatus(404)
})

export default router
