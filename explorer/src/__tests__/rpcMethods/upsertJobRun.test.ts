import { Server } from 'http'
import { Connection, getCustomRepository } from 'typeorm'
import WebSocket from 'ws'
import { getDb } from '../../database'
import { ChainlinkNode, createChainlinkNode } from '../../entity/ChainlinkNode'
import { JobRun } from '../../entity/JobRun'
import { TaskRun } from '../../entity/TaskRun'
import { start, stop } from '../../support/server'
import ethtxFixture from '../fixtures/JobRun.ethtx.fixture.json'
import createFixture from '../fixtures/JobRun.fixture.json'
import updateFixture from '../fixtures/JobRunUpdate.fixture.json'
import { clearDb } from '../testdatabase'
import { JobRunRepository } from '../../repositories/JobRunRepository'
import {
  ENDPOINT,
  createRPCRequest,
  newChainlinkNode,
} from '../../support/client'

describe('realtime', () => {
  let server: Server
  let db: Connection
  let chainlinkNode: ChainlinkNode
  let secret: string

  const authenticatedNode = async () =>
    newChainlinkNode(ENDPOINT, chainlinkNode.accessKey, secret)

  beforeAll(async () => {
    server = await start()
    db = await getDb()
  })

  beforeEach(async () => {
    clearDb()
    ;[chainlinkNode, secret] = await createChainlinkNode(
      db,
      'explore realtime test chainlinkNode',
    )
  })

  afterAll(done => stop(server, done))

  describe('#upsertJobRun', () => {
    it('can create a job run with valid JSON', async () => {
      expect.assertions(3)

      const ws = await authenticatedNode()
      const request = createRPCRequest('upsertJobRun', createFixture)
      ws.send(JSON.stringify(request))

      await new Promise(resolve => {
        ws.on('message', (data: WebSocket.Data) => {
          const response = JSON.parse(data as string)
          expect(response.result).toEqual('success')
          ws.close()
          resolve()
        })
      })

      const jobRunCount = await db.manager.count(JobRun)
      expect(jobRunCount).toEqual(1)

      const taskRunCount = await db.manager.count(TaskRun)
      expect(taskRunCount).toEqual(1)
    })

    it('can create and update a job run and task runs', async () => {
      expect.assertions(6)

      const ws = await authenticatedNode()
      const createRequest = createRPCRequest('upsertJobRun', createFixture)
      const updateRequest = createRPCRequest('upsertJobRun', updateFixture)
      ws.send(JSON.stringify(createRequest))

      await new Promise(resolve => {
        ws.on('message', (data: any) => {
          const response = JSON.parse(data)
          if (response.id === createRequest.id) {
            expect(response.result).toEqual('success')
            ws.send(JSON.stringify(updateRequest))
          }
          if (response.id === updateRequest.id) {
            expect(response.result).toEqual('success')
            ws.close()
            resolve()
          }
        })
      })

      const jobRunCount = await db.manager.count(JobRun)
      expect(jobRunCount).toEqual(1)

      const taskRunCount = await db.manager.count(TaskRun)
      expect(taskRunCount).toEqual(1)

      const jr = await db.manager.findOne(JobRun)
      expect(jr.status).toEqual('completed')

      const tr = jr.taskRuns[0]
      expect(tr.status).toEqual('completed')
    })

    it('can create a task run with transactionHash and status', async () => {
      expect.assertions(10)

      const ws = await authenticatedNode()

      const messageReceived = new Promise(resolve => {
        ws.on('message', (data: any) => {
          const response = JSON.parse(data)
          expect(response.result).toEqual('success')
          resolve()
        })
      })

      const request = createRPCRequest('upsertJobRun', ethtxFixture)
      ws.send(JSON.stringify(request))

      await messageReceived

      const jobRunCount = await db.manager.count(JobRun)
      expect(jobRunCount).toEqual(1)

      const taskRunCount = await db.manager.count(TaskRun)
      expect(taskRunCount).toEqual(4)

      const jobRunRepository = getCustomRepository(JobRunRepository, db.name)
      const jr = await jobRunRepository.getFirst()

      expect(jr.status).toEqual('completed')

      const tr = jr.taskRuns[3]
      expect(tr.status).toEqual('completed')
      expect(tr.transactionHash).toEqual(
        '0x1111111111111111111111111111111111111111111111111111111111111111',
      )
      expect(tr.timestamp).toEqual(new Date('2018-01-08T18:12:01.103Z'))
      expect(tr.blockHeight).toEqual('3735928559')
      expect(tr.blockHash).toEqual('0xbadc0de5')
      expect(tr.transactionStatus).toEqual('fulfilledRunLog')
      ws.close()
    })

    it('rejects invalid params with code -32602', async (done: any) => {
      expect.assertions(2)

      const ws = await authenticatedNode()
      const request = createRPCRequest('upsertJobRun', { invalid: 'params' })
      ws.send(JSON.stringify(request))

      ws.on('message', async (data: any) => {
        const response = JSON.parse(data)
        expect(response.error.code).toEqual(-32602)

        const count = await db.manager.count(JobRun)
        expect(count).toEqual(0)

        ws.close()
        done()
      })
    })
  })
})
