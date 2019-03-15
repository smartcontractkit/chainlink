import express from 'express'
import http from "http"
import socketio from "socket.io";
import { JobRun } from "./entity/JobRun"
import { Connection } from "typeorm"

const PORT = process.env.SERVER_PORT || 8080

const server = (dbConnection: Connection) => {
  let connections = 0

  const app = express()
  app.set("port", PORT);

  const server = new http.Server(app)
  const io = socketio(server)

  server.listen(PORT, () => {
    console.log(`server started, listening on port ${PORT}`)
  })

  app.use(express.static('client/build'))

  app.get('/api/v1/job_runs', async (req, res) => {
    const jobRuns = await dbConnection.manager.find(JobRun)

    return res.send(jobRuns)
  })

  io.on('connection', (socket) => {
    connections = connections + 1
    console.log(`websocket connected, total connections: ${connections}`);

    socket.emit('connectionCount', connections)
    socket.broadcast.emit('connectionCount', connections)

    socket.on('disconnect', () => {
      connections = connections - 1
      console.log(`websocket disconnected, total connections: ${connections}`);

      socket.broadcast.emit('connectionCount', connections)
    })
  })
}

export default server
