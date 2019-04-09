import express from 'express'
import http from 'http'
import request from 'supertest'
import { Connection } from 'typeorm'
import { JobRun } from '../../entity/JobRun'
import jobRuns from '../../controllers/jobRuns'
import seed, { JOB_RUN_B_ID } from '../../seed'
import { createDbConnection, closeDbConnection } from '../../database'
import { clearDb } from '../testdatabase'

const controller = express()
controller.use('/api/v1', jobRuns)

let server: http.Server
let connection: Connection
beforeAll(async () => {
  connection = await createDbConnection()
  server = controller.listen(null)
})
afterAll(async () => {
  if (server) {
    server.close()
    await closeDbConnection()
  }
})
afterEach(async () => clearDb())

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
    const jobRun = await connection.manager.findOne(JobRun, {
      where: { runId: JOB_RUN_B_ID }
    })
    const response = await request(server).get(`/api/v1/job_runs/${jobRun.id}`)
    expect(response.status).toEqual(200)
    expect(response.body.id).toEqual(jobRun.id)
    expect(response.body.runId).toEqual(JOB_RUN_B_ID)
    expect(response.body.taskRuns.length).toEqual(1)
    expect(response.body.initiator.id).toBeDefined()
  })

  it('returns a 404', async () => {
    const response = await request(server).get(`/api/v1/job_runs/-1`)
    expect(response.status).toEqual(404)
  })
})
