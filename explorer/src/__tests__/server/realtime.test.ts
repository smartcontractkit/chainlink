import { Server } from 'http'
import { Connection } from 'typeorm'
import WebSocket from 'ws'
import { getDb } from '../../database'
import { ChainlinkNode, createChainlinkNode } from '../../entity/ChainlinkNode'
import { DEFAULT_TEST_PORT, start, stop } from '../../support/server'
import { clearDb } from '../testdatabase'
import {
  ACCESS_KEY_HEADER,
  NORMAL_CLOSE,
  SECRET_HEADER,
} from '../../utils/constants'

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
