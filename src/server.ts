import { JobRun } from "./entity/JobRun"

const PORT = process.env.SERVER_PORT || 8080

const server = (dbConnection: any) => {
  let connections = 0

  const express = require('express')
  const app = express()
  app.set("port", PORT);

  const server = require("http").Server(app)
  const io = require("socket.io")(server)

  server.listen(PORT, () => {
    console.log(`server started, listening on port ${PORT}`)
  })

  app.use(express.static('public'))

  app.get('/api/v1/job_runs', async (req, res) => {
    const jobRuns = await dbConnection.manager.find(JobRun)

    return res.send(jobRuns)
  })

  io.on('connection', (socket: any) => {
    connections = connections + 1
    console.log(`websocket connected, total connections: ${connections}`);

    socket.on('disconnect', () => {
      connections = connections - 1
      console.log(`websocket disconnected, total connections: ${connections}`);
    })
  })
}

export default server
