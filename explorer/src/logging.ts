import PinoHttp from 'express-pino-logger'
import pino from 'pino'
import express from 'express'
import { Logger } from 'typeorm'

const options: Parameters<typeof pino>[0] = {
  name: 'Explorer',
  level: 'warn',
  redact: {
    paths: ['req.headers', 'res.headers'],
  },
}
if (process.env.EXPLORER_DEV) {
  options.prettyPrint = { colorize: true }
  options.level = 'debug'
} else if (process.env.NODE_ENV === 'test') {
  options.level = 'silent'
}
export const logger = pino(options)

export const addRequestLogging = (app: express.Express) => {
  app.use(PinoHttp({ logger }))
}

export class TypeOrmLogger implements Logger {
  public logQuery(
    query: string,
    parameters?: any[],
    queryRunner?: QueryRunner): any {
    logger.trace({msg: query, parameters})
  }

  public logQueryError(
    error: string,
    query: string,
    parameters?: any[],
    queryRunner?: QueryRunner): any {
    logger.error({msg: 'DB query failed', error, query, parameters})
  }

  public logQuerySlow(
    time: number,
    query: string,
    parameters?: any[],
    queryRunner?: QueryRunner): any {
    logger.warn({msg: 'Slow DB query detected', duration: time, query, parameters})
  }

  public logSchemaBuild(
    message: string,
    queryRunner?: QueryRunner): any {
    logger.trace({msg: message})
  }

  public logMigration(message: string, queryRunner?: QueryRunner): any {
    logger.info({msg: message})
  }

  public log(
    level: 'log' | 'info' | 'warn',
    message: any,
    queryRunner?: QueryRunner): any {
    logger[level]({msg: `TypeORM: ${message}`})
  }
}
