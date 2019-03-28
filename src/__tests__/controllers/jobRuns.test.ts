import express from 'express'
import http from 'http'
import request from 'supertest'
import jobRuns from '../../controllers/jobRuns'
import seed, { JOB_RUN_B_ID } from '../../seed'
import { createDbConnection, closeDbConnection } from '../../database'
import { clearDb } from '../testdatabase'

const controller = express()
controller.use('/api/v1', jobRuns)

let server: http.Server
beforeAll(async () => {
  await createDbConnection()
  server = controller.listen(null)
})
afterAll(async () => {
  if (server) {
    server.close()
    await closeDbConnection()
  }
})
beforeEach(async () => {
  await clearDb()
})

describe('#index', () => {
  describe('with no runs', () => {
    it('returns empty', async () => {
      const response = await request(server).get(`/api/v1/job_runs`)
      expect(response.status).toEqual(200)
    })
  })

  describe('with runs', () => {
    beforeEach(async () => {
      await seed()
    })

    it('returns runs', async () => {
      const response = await request(server).get(`/api/v1/job_runs`)
      expect(response.status).toEqual(200)
    })
  })
})

describe('#show', () => {
  beforeEach(async () => {
    await seed()
  })

  it('returns the job run with task runs', async () => {
    const response = await request(server).get(
      `/api/v1/job_runs/${JOB_RUN_B_ID}`
    )
    expect(response.status).toEqual(200)
    expect(response.body.id).toEqual(JOB_RUN_B_ID)
    expect(response.body.taskRuns.length).toEqual(1)
  })

  it('returns a 404', async () => {
    const response = await request(server).get(`/api/v1/job_runs/not-found`)
    expect(response.status).toEqual(404)
  })
})
