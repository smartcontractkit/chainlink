import express from 'express'
import http from 'http'
import request from 'supertest'
import { Connection } from 'typeorm'
import { JobRun } from '../../entity/JobRun'
import { TaskRun } from '../../entity/TaskRun'
import jobRuns from '../../controllers/jobRuns'
import { ChainlinkNode, createChainlinkNode } from '../../entity/ChainlinkNode'
import seed, { JOB_RUN_B_ID } from '../../seed'
import { closeDbConnection, getDb } from '../../database'

const controller = express()
controller.use('/api/v1', jobRuns)

let server: http.Server
let connection: Connection
beforeAll(async () => {
  connection = await getDb()
  server = controller.listen(null)
})
afterAll(async () => {
  if (server) {
    server.close()
    await closeDbConnection()
  }
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

    it('returns runs with chainlink node names', async () => {
      const response = await request(server).get(`/api/v1/job_runs`)
      expect(response.status).toEqual(200)
      const jr = response.body[0]
      expect(jr.publicChainlinkNode.name).toBeDefined()
      expect(jr.publicChainlinkNode.accessKey).not.toBeDefined()
      expect(jr.publicChainlinkNode.salt).not.toBeDefined()
      expect(jr.publicChainlinkNode.hashedSecret).not.toBeDefined()
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
  })

  describe('with out of order task runs', () => {
    let jobRunId: number
    beforeEach(async () => {
      const [chainlinkNode, _] = await createChainlinkNode(
        connection,
        'testOutOfOrderTaskRuns'
      )
      const jobRun = new JobRun()
      jobRun.chainlinkNodeId = chainlinkNode.id
      jobRun.runId = 'OutOfOrderTaskRuns'
      jobRun.jobId = 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
      jobRun.status = 'in_progress'
      jobRun.type = 'runlog'
      jobRun.txHash = 'txA'
      jobRun.requestId = 'requestIdA'
      jobRun.requester = 'requesterA'
      jobRun.createdAt = new Date('2019-04-08T01:00:00.000Z')
      await connection.manager.save(jobRun)
      jobRunId = jobRun.id

      const taskRunB = new TaskRun()
      taskRunB.jobRun = jobRun
      taskRunB.index = 1
      taskRunB.status = ''
      taskRunB.type = 'jsonparse'
      await connection.manager.save(taskRunB)

      const taskRunA = new TaskRun()
      taskRunA.jobRun = jobRun
      taskRunA.index = 0
      taskRunA.status = 'in_progress'
      taskRunA.type = 'httpget'
      await connection.manager.save(taskRunA)
    })

    it('returns ordered task runs', async () => {
      const response = await request(server).get(`/api/v1/job_runs/${jobRunId}`)
      expect(response.status).toEqual(200)
      expect(response.body.taskRuns.length).toEqual(2)
      const jr = JSON.parse(response.text)
      expect(jr.taskRuns[0].index).toEqual(0)
      expect(jr.taskRuns[1].index).toEqual(1)
    })
  })

  it('returns a 404', async () => {
    const response = await request(server).get(`/api/v1/job_runs/-1`)
    expect(response.status).toEqual(404)
  })
})
