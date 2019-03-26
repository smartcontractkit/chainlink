import express from 'express'
import http from 'http'
import socketio, { Socket } from 'socket.io'
import { Connection } from 'typeorm'
import * as controllers from './controllers'

const PORT = process.env.SERVER_PORT || 8080
const CLNODE_COUNT_EVENT = 'clnodeCount'

const server = () => {
  let clnodeCount = 0

  const app = express()
  app.use(express.static('client/build'))
  app.use('/api/v1', controllers.jobRuns)

  const server = new http.Server(app)
  server.listen(PORT, () => {
    console.log(`server started, listening on port ${PORT}`)
  })

  const statsclientio: socketio.Server = socketio(server, { path: '/client' })
  statsclientio.on('connection', (socket: Socket) => {
    socket.emit(CLNODE_COUNT_EVENT, clnodeCount)
  })

  const clnodeio: socketio.Server = socketio(server, { path: '/clnode' })
  clnodeio.on('connection', (socket: Socket) => {
    clnodeCount = clnodeCount + 1
    console.log(
      `websocket connected, total chainlink nodes connected: ${clnodeCount}`
    )
    statsclientio.emit(CLNODE_COUNT_EVENT, clnodeCount)

    socket.on('disconnect', () => {
      clnodeCount = clnodeCount - 1
      console.log(
        `websocket disconnected, total chainlink nodes connected: ${clnodeCount}`
      )
      statsclientio.emit(CLNODE_COUNT_EVENT, clnodeCount)
    })
  })

  return server
}

export default server
