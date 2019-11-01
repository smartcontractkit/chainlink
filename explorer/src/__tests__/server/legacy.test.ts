import { Server } from 'http'
import { Connection, getCustomRepository } from 'typeorm'
import WebSocket from 'ws'
import { getDb } from '../../database'
import { ChainlinkNode, createChainlinkNode } from '../../entity/ChainlinkNode'
import { JobRun } from '../../entity/JobRun'
import { TaskRun } from '../../entity/TaskRun'
import { DEFAULT_TEST_PORT, start, stop } from '../../support/server'
import ethtxFixture from '../fixtures/JobRun.ethtx.fixture.json'
import createFixture from '../fixtures/JobRun.fixture.json'
import updateFixture from '../fixtures/JobRunUpdate.fixture.json'
import { clearDb } from '../testdatabase'
import { ACCESS_KEY_HEADER, SECRET_HEADER } from '../../utils/constants'
import { JobRunRepository } from '../../repositories/JobRunRepository'

const ENDPOINT = `ws://localhost:${DEFAULT_TEST_PORT}`

const newChainlinkNode = (
  url: string,
  accessKey: string,
  secret: string,
): Promise<WebSocket> => {
  const ws = new WebSocket(ENDPOINT, {
    headers: {
      [ACCESS_KEY_HEADER]: accessKey,
      [SECRET_HEADER]: secret,
    },
  })

  return new Promise((resolve: (arg0: WebSocket) => void, reject) => {
    ws.on('error', (error: Error) => {
      reject(error)
    })

    ws.on('open', () => resolve(ws))
  })
}

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

  describe('when sending messages in legacy format', () => {
    it('can create a job run with valid JSON', async () => {
      expect.assertions(3)

      const ws = await authenticatedNode()
      ws.send(JSON.stringify(createFixture))

      await new Promise(resolve => {
        ws.on('message', (data: WebSocket.Data) => {
          const response = JSON.parse(data as string)
          expect(response.status).toEqual(201)
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
          expect(response.status).toEqual(201)
          resolve()
        })
      })

      ws.send(JSON.stringify(ethtxFixture))
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

    it('rejects malformed json events with code 422', async (done: any) => {
      expect.assertions(2)

      const ws = await authenticatedNode()
      ws.send('{invalid json}')

      ws.on('message', async (data: any) => {
        const response = JSON.parse(data)
        expect(response.status).toEqual(422)

        const count = await db.manager.count(JobRun)
        expect(count).toEqual(0)

        ws.close()
        done()
      })
    })
  })
})
