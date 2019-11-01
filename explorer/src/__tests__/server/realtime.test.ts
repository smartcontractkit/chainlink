import { Server } from 'http'
import { Connection } from 'typeorm'
import WebSocket from 'ws'
import { getDb } from '../../database'
import { ChainlinkNode, createChainlinkNode } from '../../entity/ChainlinkNode'
import { JobRun } from '../../entity/JobRun'
import { start, stop } from '../../support/server'
import { clearDb } from '../testdatabase'
import { NORMAL_CLOSE } from '../../utils/constants'
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

  const newAuthenticatedNode = async () =>
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

  describe('when sending messages in JSON-RPC format', () => {
    it('rejects non-existing methods with code -32601', async (done: any) => {
      expect.assertions(2)

      const ws = await newAuthenticatedNode()
      const request = createRPCRequest('doesNotExist')
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

      const ws = await newAuthenticatedNode()
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

      const ws = await newAuthenticatedNode()

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
