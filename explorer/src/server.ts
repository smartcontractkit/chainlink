import cookieSession from 'cookie-session'
import express from 'express'
import helmet from 'helmet'
import http from 'http'
import mime from 'mime-types'
import * as controllers from './controllers'
import { addRequestLogging, logger } from './logging'
import adminAuth from './middleware/adminAuth'
import seed from './seed'
import { bootstrapRealtime } from './server/realtime'

export const DEFAULT_PORT = parseInt(process.env.SERVER_PORT, 10) || 8080
export const COOKIE_EXPIRATION_MS = 86400000 // 1 day in ms

const server = (port: number = DEFAULT_PORT): http.Server => {
  if (process.env.NODE_ENV === 'development') {
    seed()
  }

  const app = express()
  addRequestLogging(app)

  app.use(helmet())
  if (process.env.EXPLORER_DEV) {
    // eslint-disable-next-line @typescript-eslint/no-var-requires
    const cors: typeof import('cors') = require('cors')

    app.use(
      cors({
        origin: [process.env.EXPLORER_CLIENT_ORIGIN],
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
      maxAge: COOKIE_EXPIRATION_MS,
      keys: ['key1', 'key2'],
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

  return httpServer.listen(port, () => {
    logger.info(`Server started, listening on port ${port}`)
  })
}

export default server
