import { Server } from 'http'
import WebSocket from 'ws'
import { start as startServer, DEFAULT_TEST_PORT } from '../support/server'
import { Connection } from 'typeorm'
import { closeDbConnection, getDb } from '../database'
import ethtxFixture from './fixtures/JobRun.ethtx.fixture.json'
import createFixture from './fixtures/JobRun.fixture.json'
import updateFixture from './fixtures/JobRunUpdate.fixture.json'
import { JobRun } from '../entity/JobRun'
import { TaskRun } from '../entity/TaskRun'
import { ChainlinkNode, createChainlinkNode } from '../entity/ChainlinkNode'
import { clearDb } from './testdatabase'
import {
  ACCESS_KEY_HEADER,
  CLOSED_CONNECTION,
  NORMAL_CLOSE,
  OPEN_CONNECTION,
  SECRET_HEADER,
} from '../utils/constants'

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

  beforeAll(async () => {
    server = await startServer()
    db = await getDb()
  })

  beforeEach(async () => {
    clearDb()
    ;[chainlinkNode, secret] = await createChainlinkNode(
      db,
      'explore realtime test chainlinkNode',
    )
  })

  afterAll(async () => {
    return Promise.all([server.close(), closeDbConnection()])
  })

  it('create a job run for valid JSON', async () => {
    expect.assertions(3)

    const ws = await newChainlinkNode(ENDPOINT, chainlinkNode.accessKey, secret)

    ws.send(JSON.stringify(createFixture))

    await new Promise(resolve => {
      ws.on('message', (data: WebSocket.Data) => {
        const result = JSON.parse(data as string)
        expect(result.status).toEqual(201)
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

    const ws = await newChainlinkNode(ENDPOINT, chainlinkNode.accessKey, secret)

    ws.send(JSON.stringify(createFixture))

    await new Promise(resolve => {
      let responses = 0
      ws.on('message', (data: any) => {
        responses += 1
        const result = JSON.parse(data)

        if (responses === 1) {
          expect(result.status).toEqual(201)
          ws.send(JSON.stringify(updateFixture))
        }

        if (responses === 2) {
          expect(result.status).toEqual(201)
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
    expect.assertions(7)

    const ws = await newChainlinkNode(ENDPOINT, chainlinkNode.accessKey, secret)

    const messageReceived = new Promise(resolve => {
      ws.on('message', (data: any) => {
        const result = JSON.parse(data)
        expect(result.status).toEqual(201)
        resolve()
      })
    })

    ws.send(JSON.stringify(ethtxFixture))
    await messageReceived

    const jobRunCount = await db.manager.count(JobRun)
    expect(jobRunCount).toEqual(1)

    const taskRunCount = await db.manager.count(TaskRun)
    expect(taskRunCount).toEqual(4)

    const jr = await db.manager.findOne(JobRun)
    expect(jr.status).toEqual('completed')

    const tr = jr.taskRuns[3]
    expect(tr.status).toEqual('completed')
    expect(tr.transactionHash).toEqual(
      '0x1111111111111111111111111111111111111111111111111111111111111111',
    )
    expect(tr.transactionStatus).toEqual('fulfilledRunLog')
    ws.close()
  })

  it('rejects malformed json events with code 422', async (done: any) => {
    expect.assertions(2)

    const ws = await newChainlinkNode(ENDPOINT, chainlinkNode.accessKey, secret)

    ws.send('{invalid json}')

    ws.on('message', async (data: any) => {
      const result = JSON.parse(data)
      expect(result.status).toEqual(422)

      const count = await db.manager.count(JobRun)
      expect(count).toEqual(0)

      ws.close()
      done()
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
    expect.assertions(4)

    let ws1: WebSocket
    let ws2: WebSocket

    const onCloseCallback = (event: WebSocket.CloseEvent) => {
      expect(ws1.readyState).toBe(CLOSED_CONNECTION)
      expect(ws2.readyState).toBe(OPEN_CONNECTION)
      expect(event.code).toBe(NORMAL_CLOSE)
      expect(event.reason).toEqual('Duplicate connection opened')
      ws2.close()
      done()
    }

    ws1 = await newChainlinkNode(ENDPOINT, chainlinkNode.accessKey, secret)
    ws1.addEventListener('close', onCloseCallback)
    ws2 = await newChainlinkNode(ENDPOINT, chainlinkNode.accessKey, secret)
  })
})
