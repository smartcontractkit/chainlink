import {
  requestWhitelist,
  logger as loggerConfig,
  errorLogger
} from 'express-winston'
import * as winston from 'winston'
import express from 'express'

const LOGGER_WHITELIST = [
  'url',
  'method',
  'httpVersion',
  'originalUrl',
  'query'
]
requestWhitelist.splice(0, requestWhitelist.length, ...LOGGER_WHITELIST)

export const addRequestLogging = (app: express.Express) => {
  const consoleTransport = new winston.transports.Console()

  app.use(
    loggerConfig({
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

const transports = {
  console: new winston.transports.Console({ level: 'info' })
}

export const logger = winston.createLogger({
  transports: [transports.console]
})
