import { Server } from 'http'
import { Connection, getCustomRepository } from 'typeorm'
import WebSocket from 'ws'
import { getDb } from '../database'
import { ChainlinkNode, createChainlinkNode } from '../entity/ChainlinkNode'
import { JobRun } from '../entity/JobRun'
import { TaskRun } from '../entity/TaskRun'
import { DEFAULT_TEST_PORT, start, stop } from '../support/server'
import { DEFAULT_TEST_PORT, start as startServer } from '../support/server'
import ethtxFixtureLegacy from './fixtures/JobRunLegacy.ethtx.fixture.json'
import createFixtureLegacy from './fixtures/JobRunLegacy.fixture.json'
import updateFixtureLegacy from './fixtures/JobRunUpdateLegacy.fixture.json'
import ethtxFixture from './fixtures/JobRun.ethtx.fixture.json'
import createFixture from './fixtures/JobRun.fixture.json'
import updateFixture from './fixtures/JobRunUpdate.fixture.json'
import { clearDb } from './testdatabase'
import {
  ACCESS_KEY_HEADER,
  NORMAL_CLOSE,
  SECRET_HEADER,
} from '../utils/constants'
import { JobRunRepository } from '../repositories/JobRunRepository'

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

      ws.send(JSON.stringify(createFixtureLegacy))

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

      ws.send(JSON.stringify(createFixtureLegacy))

      await new Promise(resolve => {
        let responses = 0
        ws.on('message', (data: any) => {
          responses += 1
          const response = JSON.parse(data)

          if (responses === 1) {
            expect(response.status).toEqual(201)
            ws.send(JSON.stringify(updateFixtureLegacy))
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

      ws.send(JSON.stringify(ethtxFixtureLegacy))
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

  describe('when sending messages in JSON-RPC format', () => {
    describe('#upsertJobRun', () => {
      it('can create a job run with valid JSON', async () => {
        expect.assertions(3)

        const ws = await authenticatedNode()

        ws.send(JSON.stringify(createFixture))

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

        ws.send(JSON.stringify(createFixture))

        await new Promise(resolve => {
          ws.on('message', (data: any) => {
            const response = JSON.parse(data)

            if (response.id === createFixture.id) {
              expect(response.result).toEqual('success')
              ws.send(JSON.stringify(updateFixture))
            }

            if (response.id === updateFixture.id) {
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

      it('rejects malformed json events with code -32602', async (done: any) => {
        expect.assertions(2)

        const ws = await authenticatedNode()

        const invalidJSON = {
          jsonrpc: '2.0',
          method: 'upsertJobRun',
          id: 1,
          params: {
            invalid: 'json',
          },
        }

        ws.send(JSON.stringify(invalidJSON))

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

    it('rejects non-existing methods with code -32601', async (done: any) => {
      expect.assertions(2)

      const ws = await authenticatedNode()

      const request = {
        jsonrpc: '2.0',
        method: 'doesNotExist',
        id: 1,
      }

      ws.send(JSON.stringify(request))

      ws.on('message', async (data: any) => {
        const response = JSON.parse(data)
        expect(response.error.code).toEqual(-32601)

        const count = await db.manager.count(JobRun)
        expect(count).toEqual(0)

        ws.close()
        done()
      })
    })

    // this test depends on the presence of "jsonrpc" in the message
    // otherwise, the server will attempt to process the message as a
    // legacy message and will respond with { status: 422 }.
    // This test will be more appropriate once the legacy format is removed.
    it('rejects malformed json with code -32700', async (done: any) => {
      expect.assertions(2)

      const ws = await authenticatedNode()

      const request = 'jsonrpc invalid'

      ws.send(request)

      ws.on('message', async (data: any) => {
        const response = JSON.parse(data)
        expect(response.error.code).toEqual(-32700)

        const count = await db.manager.count(JobRun)
        expect(count).toEqual(0)

        ws.close()
        done()
      })
    })

    it('rejects invalid rpc requests with code -32600', async (done: any) => {
      expect.assertions(2)

      const ws = await authenticatedNode()

      const request = {
        jsonrpc: '2.0',
        function: 'foo',
        id: 1,
      }

      ws.send(JSON.stringify(request))

      ws.on('message', async (data: any) => {
        const response = JSON.parse(data)
        expect(response.error.code).toEqual(-32600)

        const count = await db.manager.count(JobRun)
        expect(count).toEqual(0)

        ws.close()
        done()
      })
    })
  })

  it('rejects invalid authentication', async (done: any) => {
    expect.assertions(1)

    newChainlinkNode(ENDPOINT, chainlinkNode.accessKey, 'lol-no').catch(
      error => {
        expect(error).toBeDefined()
        done()
      },
    )
  })

  it('rejects multiple connections from single node', async done => {
    expect.assertions(8)

    // eslint-disable-next-line prefer-const
    let ws1: WebSocket, ws2: WebSocket, ws3: WebSocket

    // eslint-disable-next-line prefer-const
    ws1 = await newChainlinkNode(ENDPOINT, chainlinkNode.accessKey, secret)

    ws1.addEventListener('close', (event: WebSocket.CloseEvent) => {
      expect(ws1.readyState).toBe(WebSocket.CLOSED)
      expect(ws2.readyState).toBe(WebSocket.OPEN)
      expect(event.code).toBe(NORMAL_CLOSE)
      expect(event.reason).toEqual('Duplicate connection opened')
    })

    ws2 = await newChainlinkNode(ENDPOINT, chainlinkNode.accessKey, secret)

    ws2.addEventListener('close', (event: WebSocket.CloseEvent) => {
      expect(ws2.readyState).toBe(WebSocket.CLOSED)
      expect(ws3.readyState).toBe(WebSocket.OPEN)
      expect(event.code).toBe(NORMAL_CLOSE)
      expect(event.reason).toEqual('Duplicate connection opened')
      ws3.close()
      done()
    })

    ws3 = await newChainlinkNode(ENDPOINT, chainlinkNode.accessKey, secret)
  })
})
