import { Server } from 'http'
import jayson from 'jayson'
import WebSocket from 'ws'
import { ChainlinkNode, createChainlinkNode } from '../../entity/ChainlinkNode'
import {
  createRPCRequest,
  newChainlinkNode,
  sendSingleMessage,
} from '../../support/client'
import { start, stop } from '../../support/server'
import { NORMAL_CLOSE } from '../../utils/constants'
import { clearDb } from '../testdatabase'
import { getRepository } from 'typeorm'

const { PARSE_ERROR, INVALID_REQUEST, METHOD_NOT_FOUND } = jayson.Server.errors

describe('realtime', () => {
  let server: Server
  let chainlinkNode: ChainlinkNode
  let secret: string

  const newAuthenticatedNode = async () =>
    newChainlinkNode(chainlinkNode.accessKey, secret)

  beforeAll(async () => {
    server = await start()
  })

  beforeEach(async () => {
    await clearDb()
    ;[chainlinkNode, secret] = await createChainlinkNode(
      'realtime test chainlinkNode',
    )
  })

  afterAll(async done => {
    stop(server, done)
  })

  describe('when sending messages in JSON-RPC format', () => {
    let ws: WebSocket

    beforeEach(async () => {
      ws = await newAuthenticatedNode()
    })

    afterEach(async () => {
      ws.close()
    })

    it(`rejects non-existing methods with code ${METHOD_NOT_FOUND}`, async () => {
      expect.assertions(1)
      const request = createRPCRequest('doesNotExist')
      const response = await sendSingleMessage(ws, request)
      expect(response.error.code).toEqual(METHOD_NOT_FOUND)
    })

    // this test depends on the presence of "jsonrpc" in the message
    // otherwise, the server will attempt to process the message as a
    // legacy message and will respond with { status: 422 }.
    // This test will be more appropriate once the legacy format is removed.
    it(`rejects malformed json with code ${PARSE_ERROR}`, async () => {
      expect.assertions(1)
      const request = 'jsonrpc invalid'
      const response = await sendSingleMessage(ws, request)
      expect(response.error.code).toEqual(PARSE_ERROR)
    })

    it(`rejects invalid rpc requests with code ${INVALID_REQUEST}`, async () => {
      expect.assertions(1)
      const request = {
        jsonrpc: '2.0',
        function: 'foo',
        id: 1,
      }
      const response = await sendSingleMessage(ws, request)
      expect(response.error.code).toEqual(INVALID_REQUEST)
    })
  })

  it('rejects invalid authentication', async done => {
    expect.assertions(1)
    newChainlinkNode(chainlinkNode.accessKey, 'lol-no').catch(error => {
      expect(error).toBeDefined()
      done()
    })
  })

  it('rejects multiple connections from single node', async done => {
    expect.assertions(8)

    // eslint-disable-next-line prefer-const
    let ws1: WebSocket, ws2: WebSocket, ws3: WebSocket

    // eslint-disable-next-line prefer-const
    ws1 = await newAuthenticatedNode()

    ws1.addEventListener('close', (event: WebSocket.CloseEvent) => {
      expect(ws1.readyState).toBe(WebSocket.CLOSED)
      expect(ws2.readyState).toBe(WebSocket.OPEN)
      expect(event.code).toBe(NORMAL_CLOSE)
      expect(event.reason).toEqual('Duplicate connection opened')
    })

    ws2 = await newAuthenticatedNode()

    ws2.addEventListener('close', (event: WebSocket.CloseEvent) => {
      expect(ws2.readyState).toBe(WebSocket.CLOSED)
      expect(ws3.readyState).toBe(WebSocket.OPEN)
      expect(event.code).toBe(NORMAL_CLOSE)
      expect(event.reason).toEqual('Duplicate connection opened')
      ws3.close()
      done()
    })

    ws3 = await newAuthenticatedNode()
  })

  it('should record core version and SHA info, if provided', async () => {
    expect.assertions(2)

    const version = '1.0.1-rc666'
    const sha = 'BADC0FF33'
    const ws = await newChainlinkNode(
      chainlinkNode.accessKey,
      secret,
      version,
      sha,
    )
    const node = await getRepository(ChainlinkNode).findOne(chainlinkNode.id)

    expect(node.coreVersion).toEqual(version)
    expect(node.coreSHA).toEqual(sha)

    ws.close()
  })
})
