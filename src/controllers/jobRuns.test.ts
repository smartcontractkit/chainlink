import express from 'express'
import http from 'http'
import request from 'supertest'
import jobRuns from './jobRuns'
import seed from '../seed'
import { createDbConnection, closeDbConnection } from '../database'
import { clearDb } from '../testdatabase'

const controller = express()
controller.use('/api/v1', jobRuns)

let server: http.Server
beforeAll(async () => {
  await createDbConnection()
  server = controller.listen(null)
})
afterAll(async () => {
  server.close()
  await closeDbConnection()
})
beforeEach(async () => {
  await clearDb()
})

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
