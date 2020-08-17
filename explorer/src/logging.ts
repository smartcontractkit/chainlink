import express from 'express'
import PinoHttp from 'express-pino-logger'
import pino from 'pino'
import { Logger } from 'typeorm'
import { Config, Environment } from './config'

const options: Parameters<typeof pino>[0] = {
  name: 'Explorer',
  level: 'info',
  redact: {
    paths: ['req.headers', 'res.headers'],
  },
}

switch (Config.env()) {
  case Environment.PROD:
    break
  case Environment.TEST:
    options.level = 'silent'
    break
  default:
    options.prettyPrint = { colorize: true }
    options.level = 'trace'
}
export const logger = pino(options)

export const addRequestLogging = (app: express.Express) => {
  app.use(PinoHttp({ logger }))
}

export class TypeOrmLogger implements Logger {
  private logger = logger.child({ module: 'TypeORM' })

  public logQuery(query: string, parameters?: any[]): any {
    this.logger.trace(query, { parameters })
  }

  public logQueryError(error: string, query: string, parameters?: any[]): any {
    this.logger.error('DB query failed', { error, query, parameters })
  }

  public logQuerySlow(time: number, query: string, parameters?: any[]): any {
    this.logger.warn('Slow DB query detected', {
      duration: time,
      query,
      parameters,
    })
  }

  public logSchemaBuild(message: string): any {
    this.logger.trace(message)
  }

  public logMigration(message: string): any {
    this.logger.info(message)
  }

  public log(level: 'log' | 'info' | 'warn', message: any): any {
    this.logger[level](message)
  }
}
