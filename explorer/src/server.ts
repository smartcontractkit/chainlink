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
  app.use(express.static('client/build'))
  app.use('/api/v1', controllers.jobRuns)
  app.use(
    expressWinston.logger({
      expressFormat: true,
      meta: true,
      msg: 'HTTP {{req.method}} {{req.url}}',
      transports: [new winston.transports.Console()]
    })
  )

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
