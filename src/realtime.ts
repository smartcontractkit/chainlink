import { getDb } from './database'
import http from 'http'
import { fromString } from './entity/JobRun'
import WebSocket from 'ws'

const CLNODE_COUNT_EVENT = 'clnodeCount'

export const bootstrapRealtime = (server: http.Server) => {
  const db = getDb()
  let clnodeCount = 0

  const wss = new WebSocket.Server({ server, perMessageDeflate: false })
  wss.on('connection', function connection(ws) {
    clnodeCount = clnodeCount + 1
    console.log(
      `websocket connected, total chainlink nodes connected: ${clnodeCount}`
    )
    ws.on('message', function incoming(message: WebSocket.Data) {
      let result

      console.log('received: %s', message)
      try {
        const jobRun = fromString(message as string)
        db.manager
          .save(jobRun)
          .then(entity => {
            console.log('saved job run %s', entity.id)
          })
          .catch(console.error)
        result = { status: 201 }
      } catch (e) {
        console.error(e)
        result = { status: 422 }
      }

      ws.send(JSON.stringify(result))
    })

    ws.on('close', () => {
      clnodeCount = clnodeCount - 1
      console.log(
        `websocket disconnected, total chainlink nodes connected: ${clnodeCount}`
      )
    })
  })
}
