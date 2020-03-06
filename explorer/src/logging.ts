import PinoHttp from 'express-pino-logger'
import pino from 'pino'
import express from 'express'

export const addRequestLogging = (app: express.Express) => {
  app.use(PinoHttp)
}

const options: any = {
  name: 'Explorer',
  level: 'warn',
}
if (process.env.EXPLORER_DEV) {
  options['prettyPrint'] = { colorize: true }
  options['level'] = 'debug'
}
export const logger = pino(options)
