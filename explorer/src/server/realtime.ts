import http from 'http'
import { logger } from '../logging'
import WebSocket from 'ws'
import { getDb } from '../database'
import { authenticate } from '../sessions'
import { closeSession, Session } from '../entity/Session'
import { handleMessage } from './handleMessage'
import {
  ACCESS_KEY_HEADER,
  NORMAL_CLOSE,
  SECRET_HEADER,
} from '../utils/constants'

export const bootstrapRealtime = async (server: http.Server) => {
  const db = await getDb()
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
      logger.debug('websocket connection attempt')

      const accessKey = info.req.headers[ACCESS_KEY_HEADER]
      const secret = info.req.headers[SECRET_HEADER]

      if (typeof accessKey !== 'string' || typeof secret !== 'string') {
        logger.info('client rejected, invalid authentication request')
        return
      }

      authenticate(db, accessKey, secret).then((session: Session | null) => {
        if (session === null) {
          logger.info('client rejected, failed authentication')
          callback(false, 401)
          return
        }

        logger.debug(
          `websocket client successfully authenticated, new session for node ${session.chainlinkNodeId}`,
        )
        sessions.set(accessKey, session)
        const existingConnection = connections.get(accessKey)
        if (existingConnection) {
          existingConnection.close(NORMAL_CLOSE, 'Duplicate connection opened')
        }
        callback(true, 200)
      })
    },
  })

  wss.on('connection', (ws: WebSocket, request: http.IncomingMessage) => {
    // accessKey type already validated in verifyClient()
    const accessKey = request.headers[ACCESS_KEY_HEADER].toString()
    connections.set(accessKey, ws)
    clnodeCount = clnodeCount + 1

    logger.info(
      `websocket connected, total chainlink nodes connected: ${clnodeCount}`,
    )

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
        closeSession(db, session)
        sessions.delete(accessKey)
      }
      if (ws === existingConnection) {
        connections.delete(accessKey)
      }
      clnodeCount = clnodeCount - 1
      logger.info(
        `websocket disconnected, total chainlink nodes connected: ${clnodeCount}`,
      )
    })
  })
}
