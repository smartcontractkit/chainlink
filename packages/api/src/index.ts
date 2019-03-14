import "reflect-metadata"
import { createConnection } from "typeorm"
import * as express from 'express'
import { JobRun } from "./entity/JobRun"

const PORT = process.env.SERVER_PORT || 8080

createConnection().then(async connection => {
  console.log("Inserting a new job run into the database...")
  const jobRunA = new JobRun()
  jobRunA.requestId = "66eb9365-6c0c-487c-9297-7b1b44d87711"
  await connection.manager.save(jobRunA)
  console.log("Saved a new job run with id: " + jobRunA.id)

  console.log("Inserting a new job run into the database...")
  const jobRunB = new JobRun()
  jobRunB.requestId = "66eb9365-6c0c-487c-9297-7b1b44d87711"
  await connection.manager.save(jobRunB)
  console.log("Saved a new job run with id: " + jobRunB.id)

  const app = express()

  app.use(express.static('public'))

  app.get('/api/v1/job_runs', async (req, res) => {
    const jobRuns = await connection.manager.find(JobRun)

    return res.send(jobRuns)
  })

  app.listen(PORT, () => {
    console.log(`server started, listening on port ${PORT}`)
  })

}).catch(error => console.log(error))
