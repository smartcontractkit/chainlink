import * as controllers from './controllers'
import { requestWhitelist, logger, errorLogger } from 'express-winston'
import * as winston from 'winston'
import express from 'express'
import http from 'http'
import mime from 'mime-types'
import { bootstrapRealtime } from './realtime'
import seed from './seed'

export const DEFAULT_PORT = parseInt(process.env.SERVER_PORT, 10) || 8080

const LOGGER_WHITELIST = [
  'url',
  'method',
  'httpVersion',
  'originalUrl',
  'query'
]
requestWhitelist.splice(0, requestWhitelist.length, ...LOGGER_WHITELIST)

const addLogging = (app: express.Express) => {
  const consoleTransport = new winston.transports.Console()

  app.use(
    logger({
      expressFormat: true,
      meta: true,
      msg: 'HTTP {{req.method}} {{req.url}}',
      transports: [consoleTransport]
    })
  )

  app.use(
    errorLogger({
      transports: [consoleTransport]
    })
  )
}

const server = (port: number = DEFAULT_PORT) => {
  if (process.env.NODE_ENV === 'development') {
    seed()
  }

  const app = express()
  if (process.env.NODE_ENV !== 'test') {
    addLogging(app)
  }

  app.use(
    express.static('client/build', {
      maxAge: '1y',
      setHeaders: function(res, path) {
        if (mime.lookup(path) === 'text/html') {
          res.setHeader('Cache-Control', 'public, max-age=0')
        }
      }
    })
  )
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
