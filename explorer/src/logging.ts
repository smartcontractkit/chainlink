import PinoHttp from 'express-pino-logger'
import pino from 'pino'
import express from 'express'

export const addRequestLogging = (app: express.Express) => {
  app.use(PinoHttp)
}

export const logger = pino({
  name: 'Explorer',
  level: 'debug',
  prettyPrint: { colorize: true },
})
