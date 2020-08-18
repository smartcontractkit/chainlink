import http from 'http'
import WebSocket from 'ws'
import { closeSession, Session } from '../entity/Session'
import { logger } from '../logging'
import { authenticate } from '../sessions'
import {
  ACCESS_KEY_HEADER,
  CORE_SHA_HEADER,
  CORE_VERSION_HEADER,
  NORMAL_CLOSE,
  SECRET_HEADER,
} from '../utils/constants'
import { handleMessage } from './handleMessage'

export type AuthInfo = {
  accessKey?: string
  secret?: string
  coreVersion?: string
  coreSHA?: string
}

export const bootstrapRealtime = (server: http.Server) => {
  let clnodeCount = 0
  const sessions = new Map<string, Session>()
  const connections = new Map<string, WebSocket>()

  // NOTE: This relies on the subtle detail that info.req is the same request
  // as passed in to wss.on to key a session
  const wss = new WebSocket.Server({
    server,
    perMessageDeflate: false,
    verifyClient: (
      info: { origin: string; secure: boolean; req: http.IncomingMessage },
      callback: (
        res: boolean,
        code?: number,
        message?: string,
        headers?: http.OutgoingHttpHeaders,
      ) => void,
    ) => {
      const remote = remoteDetails(info.req)
      logger.debug({ msg: 'websocket connection attempt', remote })

      const authInfo = extractAuthInfo(info.req.headers)
      const { accessKey, secret } = authInfo

      if (typeof accessKey !== 'string' || typeof secret !== 'string') {
        logger.warn({
          msg: 'client rejected, invalid authentication request',
          origin: info.origin,
          ...remote,
        })
        return
      }

      authenticate(authInfo).then((session: Session | null) => {
        if (session === null) {
          logger.warn({
            msg: 'client rejected, failed authentication',
            accessKey,
            origin: info.origin,
            ...remote,
          })
          callback(false, 401)
          return
        }

        logger.info({
          msg: `websocket client successfully authenticated`,
          nodeID: session.chainlinkNodeId,
          origin: info.origin,
          ...remote,
        })

        sessions.set(accessKey, session)
        const existingConnection = connections.get(accessKey)
        if (existingConnection) {
          existingConnection.close(NORMAL_CLOSE, 'Duplicate connection opened')
          logger.warn({
            msg: 'Duplicated connection opened',
            nodeID: session.chainlinkNodeId,
            origin: info.origin,
            ...remote,
          })
        }
        callback(true, 200)
      })
    },
  })

  wss.on('connection', (ws: WebSocket, request: http.IncomingMessage) => {
    const remote = remoteDetails(request)

    // accessKey type already validated in verifyClient()
    const accessKey = request.headers[ACCESS_KEY_HEADER].toString()
    connections.set(accessKey, ws)
    clnodeCount = clnodeCount + 1

    logger.info({
      msg: 'websocket connected',
      nodeCount: clnodeCount,
      ...remote,
    })

    ws.on('message', async (message: WebSocket.Data) => {
      const session = sessions.get(accessKey)
      if (session == null) {
        ws.close()
        return
      }

      const result = await handleMessage(message as string, {
        chainlinkNodeId: session.chainlinkNodeId,
      })

      ws.send(JSON.stringify(result))
    })

    ws.on('close', () => {
      const session = sessions.get(accessKey)
      const existingConnection = connections.get(accessKey)

      if (session != null) {
        closeSession(session)
        sessions.delete(accessKey)
      }
      if (ws === existingConnection) {
        connections.delete(accessKey)
      }
      clnodeCount = clnodeCount - 1
      logger.info({
        msg: 'websocket disconnected',
        nodeCount: clnodeCount,
        ...remote,
      })
    })
  })
}

function remoteDetails(
  req: http.IncomingMessage,
): Record<string, string | number | null> {
  return {
    remotePort: req.socket.remotePort,
    remoteAddress: req.socket.remoteAddress,
  }
}

function extractAuthInfo(headers: http.IncomingHttpHeaders): AuthInfo {
  return {
    accessKey: stringOrUndefined(headers[ACCESS_KEY_HEADER]),
    secret: stringOrUndefined(headers[SECRET_HEADER]),
    coreVersion: stringOrUndefined(headers[CORE_VERSION_HEADER]),
    coreSHA: stringOrUndefined(headers[CORE_SHA_HEADER]),
  }
}

function stringOrUndefined(
  key: string | string[] | undefined,
): string | undefined {
  return typeof key === 'string' ? key : undefined
}
