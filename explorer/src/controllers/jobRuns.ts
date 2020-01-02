import { getDb } from '../database'
import { Router, Request, Response } from 'express'
import { search, count, SearchParams } from '../queries/search'
import jobRunsSerializer from '../serializers/jobRunsSerializer'
import jobRunSerializer from '../serializers/jobRunSerializer'
import { getCustomRepository } from 'typeorm'
import { JobRunRepository } from '../repositories/JobRunRepository'

const router = Router()

const DEFAULT_PAGE = 1
const DEFAULT_SIZE = 10

const searchParams = (req: Request): SearchParams => {
  const page = parseInt(req.query.page, 10) || DEFAULT_PAGE
  const size = parseInt(req.query.size, 10) || DEFAULT_SIZE

  return {
    searchQuery: req.query.query,
    page,
    limit: size,
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
  const jobRunRepository = getCustomRepository(JobRunRepository, db.name)
  const jobRun = await jobRunRepository.findById(id)

  if (jobRun) {
    const json = jobRunSerializer(jobRun)
    return res.send(json)
  }

  return res.sendStatus(404)
})

export default router
