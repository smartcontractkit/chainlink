import http from 'http'
import { getRepository } from 'typeorm'
import request from 'supertest'
import { ChainlinkNode, createChainlinkNode } from '../../entity/ChainlinkNode'
import { JobRun } from '../../entity/JobRun'
import { TaskRun } from '../../entity/TaskRun'
import { createJobRun } from '../../factories'
import { start, stop } from '../../support/server'

let server: http.Server
beforeAll(async () => {
  server = await start()
})
afterAll(done => stop(server, done))

describe('#index', () => {
  describe('with no runs', () => {
    it('returns empty', async () => {
      const response = await request(server).get('/api/v1/job_runs')
      expect(response.status).toEqual(200)
    })
  })

  describe('with runs', () => {
    let jobRun: JobRun

    beforeEach(async () => {
      const [node] = await createChainlinkNode('jobRunsIndexTestChainlinkNode')
      jobRun = await createJobRun(node)
    })

    it('returns runs with chainlink node names', async () => {
      const response = await request(server).get(
        `/api/v1/job_runs?query=${jobRun.runId}`,
      )
      expect(response.status).toEqual(200)

      const chainlinkNode = response.body.included[0]
      expect(chainlinkNode.attributes.name).toBeDefined()
      expect(chainlinkNode.attributes.accessKey).not.toBeDefined()
      expect(chainlinkNode.attributes.salt).not.toBeDefined()
      expect(chainlinkNode.attributes.hashedSecret).not.toBeDefined()
    })
  })
})

describe('#show', () => {
  let node: ChainlinkNode

  beforeEach(async () => {
    ;[node] = await createChainlinkNode('jobRunsShowTestChainlinkNode')
  })

  it('returns the job run with task runs', async () => {
    const jobRun = await createJobRun(node)
    const response = await request(server).get(`/api/v1/job_runs/${jobRun.id}`)
    expect(response.status).toEqual(200)
    expect(response.body.data.id).toEqual(jobRun.id)
    expect(response.body.data.attributes.runId).toEqual(jobRun.runId)
    expect(response.body.data.relationships.taskRuns.data.length).toEqual(1)
  })

  describe('with out of order task runs', () => {
    let jobRunId: string
    beforeEach(async () => {
      const [chainlinkNode] = await createChainlinkNode(
        'testOutOfOrderTaskRuns',
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
      await getRepository(JobRun).save(jobRun)
      jobRunId = jobRun.id

      const taskRunB = new TaskRun()
      taskRunB.jobRun = jobRun
      taskRunB.index = 1
      taskRunB.status = ''
      taskRunB.type = 'jsonparse'
      await getRepository(TaskRun).save(taskRunB)

      const taskRunA = new TaskRun()
      taskRunA.jobRun = jobRun
      taskRunA.index = 0
      taskRunA.status = 'in_progress'
      taskRunA.type = 'httpget'
      await getRepository(TaskRun).save(taskRunA)
    })

    it('returns ordered task runs', async () => {
      const response = await request(server).get(`/api/v1/job_runs/${jobRunId}`)
      expect(response.status).toEqual(200)
      expect(response.body.data.relationships.taskRuns.data.length).toEqual(2)

      const taskRun1 = response.body.included[1]
      const taskRun2 = response.body.included[2]
      expect(taskRun1.attributes.index).toEqual(0)
      expect(taskRun2.attributes.index).toEqual(1)
    })
  })

  it('returns the job run with only public chainlink node information', async () => {
    const jobRun = await createJobRun(node)

    const response = await request(server).get(`/api/v1/job_runs/${jobRun.id}`)
    expect(response.status).toEqual(200)

    const clnode = response.body.included[0]
    expect(clnode).toBeDefined()
    expect(clnode.id).toBeDefined()
    expect(clnode.attributes.name).toEqual('jobRunsShowTestChainlinkNode')
    expect(clnode.attributes.accessKey).not.toBeDefined()
    expect(clnode.attributes.hashedSecret).not.toBeDefined()
    expect(clnode.attributes.salt).not.toBeDefined()
  })

  it('returns a 404', async () => {
    const response = await request(server).get('/api/v1/job_runs/1')
    expect(response.status).toEqual(404)
  })
})
