import cookieSession from 'cookie-session'
import express from 'express'
import helmet from 'helmet'
import http from 'http'
import mime from 'mime-types'
import { Environment, ExplorerConfig } from './config'
import * as controllers from './controllers'
import { addRequestLogging, logger } from './logging'
import adminAuth from './middleware/adminAuth'
import seed from './seed'
import { bootstrapRealtime } from './server/realtime'

export default function server(conf: ExplorerConfig): Promise<http.Server> {
  if (conf.env === Environment.DEV) {
    seed()
  }

  const app = express()
  addRequestLogging(app)

  app.use(helmet())
  if (conf.env === Environment.DEV) {
    // eslint-disable-next-line @typescript-eslint/no-var-requires
    const cors: typeof import('cors') = require('cors')

    app.use(
      cors({
        origin: [conf.clientOrigin],
        methods: 'GET,HEAD,PUT,PATCH,POST,DELETE',
        preflightContinue: false,
        optionsSuccessStatus: 204,
        credentials: true,
      }),
    )
  }

  app.use(
    cookieSession({
      name: 'explorer',
      maxAge: conf.cookieExpirationMs,
      secret: conf.cookieSecret,
    }),
  )
  app.use(express.json())
  app.use(
    express.static('client/build', {
      maxAge: '1y',
      setHeaders(res, path) {
        if (mime.lookup(path) === 'text/html') {
          res.setHeader('Cache-Control', 'public, max-age=0')
        }
      },
    }),
  )

  app.use('/api/v1/admin/*', adminAuth)
  const ADMIN_CONTROLLERS = [
    controllers.adminLogin,
    controllers.adminLogout,
    controllers.adminNodes,
    controllers.adminHeads,
  ]
  ADMIN_CONTROLLERS.forEach(c => app.use('/api/v1/admin', c))

  app.use('/api/v1', controllers.jobRuns)

  app.get('/*', (_, res) => {
    res.sendFile(`${__dirname}/public/index.html`)
  })

  const httpServer = new http.Server(app)
  bootstrapRealtime(httpServer)

  return new Promise(resolve => {
    const server = httpServer.listen(conf.port, async () => {
      logger.info(`Server started, listening on port ${conf.port}`)
      resolve(server)
    })
  })
}
