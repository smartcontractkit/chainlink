import { Server } from 'http'
import jayson from 'jayson'
import { getRepository, getCustomRepository } from 'typeorm'
import WebSocket from 'ws'
import { ChainlinkNode, createChainlinkNode } from '../../entity/ChainlinkNode'
import { JobRun } from '../../entity/JobRun'
import { TaskRun } from '../../entity/TaskRun'
import { JobRunRepository } from '../../repositories/JobRunRepository'
import {
  createRPCRequest,
  newChainlinkNode,
  sendSingleMessage,
} from '../../support/client'
import { start, stop } from '../../support/server'
import ethtxFixture from '../fixtures/JobRun.ethtx.fixture.json'
import createFixture from '../fixtures/JobRun.fixture.json'
import updateFixture from '../fixtures/JobRunUpdate.fixture.json'
import { clearDb } from '../testdatabase'

const { INVALID_PARAMS } = jayson.Server.errors

describe('realtime', () => {
  let server: Server
  let chainlinkNode: ChainlinkNode
  let secret: string
  let ws: WebSocket

  beforeAll(async () => {
    server = await start()
  })

  beforeEach(async () => {
    clearDb()
    ;[chainlinkNode, secret] = await createChainlinkNode(
      'upsertJobRun test chainlinkNode',
    )
    ws = await newChainlinkNode(chainlinkNode.accessKey, secret)
  })

  afterEach(async () => {
    ws.close()
  })

  afterAll(done => stop(server, done))

  describe('#upsertJobRun', () => {
    it('can create a job run with valid JSON', async () => {
      expect.assertions(3)

      const request = createRPCRequest('upsertJobRun', createFixture)
      const response = await sendSingleMessage(ws, request)
      expect(response.result).toEqual('success')

      const jobRunCount = await getRepository(JobRun).count()
      expect(jobRunCount).toEqual(1)

      const taskRunCount = await getRepository(TaskRun).count()
      expect(taskRunCount).toEqual(1)
    })

    it('can create and update a job run and task runs', async () => {
      expect.assertions(6)

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

      const jobRunCount = await getRepository(JobRun).count()
      expect(jobRunCount).toEqual(1)

      const taskRunCount = await getRepository(TaskRun).count()
      expect(taskRunCount).toEqual(1)

      const jr = await getRepository(JobRun).findOne()
      expect(jr.status).toEqual('completed')

      const tr = jr.taskRuns[0]
      expect(tr.status).toEqual('completed')
    })

    it('can create a task run with transactionHash and status', async () => {
      expect.assertions(10)

      const request = createRPCRequest('upsertJobRun', ethtxFixture)
      const response = await sendSingleMessage(ws, request)
      expect(response.result).toEqual('success')

      const jobRunCount = await getRepository(JobRun).count()
      expect(jobRunCount).toEqual(1)

      const taskRunCount = await getRepository(TaskRun).count()
      expect(taskRunCount).toEqual(4)

      const jobRunRepository = getCustomRepository(JobRunRepository)
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
    })

    it(`rejects invalid params with code ${INVALID_PARAMS}`, async () => {
      expect.assertions(2)
      const request = createRPCRequest('upsertJobRun', { invalid: 'params' })
      const response = await sendSingleMessage(ws, request)
      expect(response.error.code).toEqual(INVALID_PARAMS)
      const count = await getRepository(JobRun).count()
      expect(count).toEqual(0)
    })
  })
})
