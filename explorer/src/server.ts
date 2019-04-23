import express from 'express'
import http from 'http'
import * as controllers from './controllers'
import { bootstrapRealtime } from './realtime'

export const DEFAULT_PORT = parseInt(process.env.SERVER_PORT, 10) || 8080

const server = (port: number = DEFAULT_PORT) => {
  const app = express()
  app.use(express.static('client/build'))
  app.use('/api/v1', controllers.jobRuns)

  app.get('/*', (_, res) => {
    res.sendFile(`${__dirname}/public/index.html`)
  })

  const server = new http.Server(app)
  bootstrapRealtime(server)
  return server.listen(port, () => {
    console.log(`server started, listening on port ${port}`)
  })
}

export default server
