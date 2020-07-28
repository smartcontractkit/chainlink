import http from 'http'
import { logger } from '../logging'
import WebSocket from 'ws'
import { authenticate } from '../sessions'
import { closeSession, Session } from '../entity/Session'
import { handleMessage } from './handleMessage'
import {
  ACCESS_KEY_HEADER,
  NORMAL_CLOSE,
  SECRET_HEADER,
} from '../utils/constants'

export const bootstrapRealtime = async (server: http.Server) => {
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

      const accessKey = info.req.headers[ACCESS_KEY_HEADER]
      const secret = info.req.headers[SECRET_HEADER]

      if (typeof accessKey !== 'string' || typeof secret !== 'string') {
        logger.warn({
          msg: 'client rejected, invalid authentication request',
          origin: info.origin,
          ...remote,
        })
        return
      }

      authenticate(accessKey, secret).then((session: Session | null) => {
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
