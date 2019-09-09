import express from 'express'
import helmet from 'helmet'
import http from 'http'
import mime from 'mime-types'
import * as controllers from './controllers'
import { addRequestLogging, logger } from './logging'
import { bootstrapRealtime } from './realtime'
import seed from './seed'

export const DEFAULT_PORT = parseInt(process.env.SERVER_PORT, 10) || 8080

const server = (port: number = DEFAULT_PORT) => {
  if (process.env.NODE_ENV === 'development') {
    seed()
  }

  const app = express()
  addRequestLogging(app)

  app.use(helmet())
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
  app.use('/api/v1', controllers.jobRuns)

  app.get('/*', (_, res) => {
    res.sendFile(`${__dirname}/public/index.html`)
  })

  const httpServer = new http.Server(app)
  bootstrapRealtime(httpServer)
  return httpServer.listen(port, () => {
    logger.info(`server started, listening on port ${port}`)
  })
}

export default server
