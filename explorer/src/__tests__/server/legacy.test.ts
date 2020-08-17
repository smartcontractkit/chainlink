import { Server } from 'http'
import { getCustomRepository, getRepository } from 'typeorm'
import WebSocket from 'ws'
import { ChainlinkNode, createChainlinkNode } from '../../entity/ChainlinkNode'
import { JobRun } from '../../entity/JobRun'
import { TaskRun } from '../../entity/TaskRun'
import { JobRunRepository } from '../../repositories/JobRunRepository'
import { newChainlinkNode, sendSingleMessage } from '../../support/client'
import { start, stop } from '../../support/server'
import ethtxFixture from '../fixtures/JobRun.ethtx.fixture.json'
import createFixture from '../fixtures/JobRun.fixture.json'
import updateFixture from '../fixtures/JobRunUpdate.fixture.json'
import { clearDb } from '../testdatabase'

describe('realtime', () => {
  let server: Server
  let chainlinkNode: ChainlinkNode
  let secret: string
  let ws: WebSocket

  function closeWebsocket(): Promise<void> {
    ws?.close()
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => {
        reject('[closeWebsocket] Timed out waiting.')
      }, 3000)

      ws?.on('close', () => {
        clearTimeout(timer)
        resolve()
      })
    })
  }

  beforeAll(async () => {
    server = await start()
  })

  beforeEach(async () => {
    clearDb()
    ;[chainlinkNode, secret] = await createChainlinkNode(
      'legacy test chainlinkNode',
    )
    ws = await newChainlinkNode(chainlinkNode.accessKey, secret)
  })

  afterEach(async () => {
    await closeWebsocket()
  })

  afterAll(done => stop(server, done))

  describe('when sending messages in legacy format', () => {
    it('can create a job run with valid JSON', async () => {
      expect.assertions(3)

      const response = await sendSingleMessage(ws, createFixture)
      expect(response.status).toEqual(201)

      const jobRunCount = await getRepository(JobRun).count()
      expect(jobRunCount).toEqual(1)

      const taskRunCount = await getRepository(TaskRun).count()
      expect(taskRunCount).toEqual(1)
    })

    it('can create and update a job run and task runs', async () => {
      expect.assertions(6)

      ws.send(JSON.stringify(createFixture))

      await new Promise(resolve => {
        let responses = 0
        ws.on('message', (data: any) => {
          responses += 1
          const response = JSON.parse(data)

          if (responses === 1) {
            expect(response.status).toEqual(201)
            ws.send(JSON.stringify(updateFixture))
          }

          if (responses === 2) {
            expect(response.status).toEqual(201)
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

      const response = await sendSingleMessage(ws, ethtxFixture)
      expect(response.status).toEqual(201)

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

    it('rejects malformed json events with code 422', async () => {
      expect.assertions(2)
      const request = '{invalid json}'
      const response = await sendSingleMessage(ws, request)
      expect(response.status).toEqual(422)
      const count = await getRepository(JobRun).count()
      expect(count).toEqual(0)
    })
  })
})
