import "reflect-metadata"
import { createConnection } from "typeorm"
import express from "express"
import { JobRun } from "./entity/JobRun"

const PORT = process.env.SERVER_PORT || 8080

const startServer = (dbConnection: any) => {
  const app = express()
  app.set("port", PORT);

  const server = require("http").Server(app)
  const io = require("socket.io")(server)

  server.listen(PORT, () => {
    console.log(`server started, listening on port ${PORT}`)
  })

  app.use(express.static('client/build'))

  app.get('/api/v1/job_runs', async (req, res) => {
    const jobRuns = await dbConnection.manager.find(JobRun)

    return res.send(jobRuns)
  })

  io.on('connection', (socket: any) => {
    console.log('a user connected')

    socket.on('disconnect', () => {
      console.log('a user disconnected');
    })
  })
}

const seedDb = async (dbConnection: any) => {
  console.log("Inserting a new job run into the database...")
  const jobRunA = new JobRun()
  jobRunA.requestId = "66eb9365-6c0c-487c-9297-7b1b44d87711"
  await dbConnection.manager.save(jobRunA)
  console.log("Saved a new job run with id: " + jobRunA.id)

  console.log("Inserting a new job run into the database...")
  const jobRunB = new JobRun()
  jobRunB.requestId = "66eb9365-6c0c-487c-9297-7b1b44d87711"
  await dbConnection.manager.save(jobRunB)
  console.log("Saved a new job run with id: " + jobRunB.id)
}

createConnection().then(async dbConnection => {
  seedDb(dbConnection)
  startServer(dbConnection)
}).catch(error => console.log(error))
