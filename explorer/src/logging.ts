import PinoHttp from 'express-pino-logger'
import pino from 'pino'
import express from 'express'

const options: any = {
  name: 'Explorer',
  level: 'warn',
}
if (process.env.EXPLORER_DEV) {
  options['prettyPrint'] = { colorize: true }
  options['level'] = 'debug'
} else if (process.env.NODE_ENV === 'test') {
  options['level'] = 'silent'
}
export const logger = pino(options)

export const addRequestLogging = (app: express.Express) => {
  app.use(PinoHttp({ logger }))
}
