import http from 'http'
import WebSocket from 'ws'

const CLNODE_COUNT_EVENT = 'clnodeCount'

export const bootstrapRealtime = (server: http.Server) => {
  let clnodeCount = 0

  const wss = new WebSocket.Server({ server, perMessageDeflate: false })
  wss.on('connection', function connection(ws) {
    clnodeCount = clnodeCount + 1
    console.log(
      `websocket connected, total chainlink nodes connected: ${clnodeCount}`
    )
    ws.on('message', function incoming(message) {
      console.log('received: %s', message)
    })

    ws.on('close', () => {
      clnodeCount = clnodeCount - 1
      console.log(
        `websocket disconnected, total chainlink nodes connected: ${clnodeCount}`
      )
    })
  })
}
