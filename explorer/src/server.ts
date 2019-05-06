import * as controllers from './controllers'
import * as expressWinston from 'express-winston'
import * as winston from 'winston'
import express from 'express'
import http from 'http'
import { bootstrapRealtime } from './realtime'
import { ChainlinkNode, createChainlinkNode } from './entity/ChainlinkNode'
import seed from './seed'

export const DEFAULT_PORT = parseInt(process.env.SERVER_PORT, 10) || 8080

const server = (port: number = DEFAULT_PORT) => {
  if (process.env.NODE_ENV === 'development') {
    seed()
  }

  const app = express()

  const consoleTransport = new winston.transports.Console()

  app.use(
    expressWinston.logger({
      expressFormat: true,
      meta: true,
      msg: 'HTTP {{req.method}} {{req.url}}',
      transports: [consoleTransport]
    })
  )

  app.use(
    expressWinston.errorLogger({
      transports: [consoleTransport]
    })
  )

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
