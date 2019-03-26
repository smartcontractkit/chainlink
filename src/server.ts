import express from 'express'
import http from 'http'
import { Connection } from 'typeorm'
import * as controllers from './controllers'
import { bootstrapRealtime } from './realtime'

const PORT = process.env.SERVER_PORT || 8080

const server = () => {
  const app = express()
  app.use(express.static('client/build'))
  app.use('/api/v1', controllers.jobRuns)

  const server = new http.Server(app)
  bootstrapRealtime(server)
  return server.listen(PORT, () => {
    console.log(`server started, listening on port ${PORT}`)
  })
}

export default server
