import { getDb } from '../database'
import { JobRun } from '../entity/JobRun'
import { Router, Request, Response } from 'express'
import { search, count, ISearchParams } from '../queries/search'
import jobRunsSerializer from '../serializers/jobRunsSerializer'
import jobRunSerializer from '../serializers/jobRunSerializer'

const router = Router()

const DEFAULT_PAGE = 1
const DEFAULT_SIZE = 10

const searchParams = (req: Request): ISearchParams => {
  const page = parseInt(req.query.page, 10) || DEFAULT_PAGE
  const size = parseInt(req.query.size, 10) || DEFAULT_SIZE

  return {
    searchQuery: req.query.query,
    page: page,
    limit: size
  }
}

router.get('/job_runs', async (req: Request, res: Response) => {
  const params = searchParams(req)
  const db = await getDb()
  const runs = await search(db, params)
  const runCount = await count(db, params)
  const json = jobRunsSerializer(runs, runCount)
  return res.send(json)
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
    const json = jobRunSerializer(jobRun)
    return res.send(json)
  }

  return res.sendStatus(404)
})

export default router
