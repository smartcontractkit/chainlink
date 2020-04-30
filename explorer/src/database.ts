import 'reflect-metadata'
import { Connection, createConnection } from 'typeorm'
import { TypeOrmLogger } from './logging'

const getEnv = (): string => {
  return process.env.TYPEORM_NAME || process.env.NODE_ENV || 'development'
}

export const openDbConnection = async (): Promise<Connection> => {
  const env = getEnv()
  const options = await import(`../ormconfig/${env}.json`)
  return createConnection({ ...options, logger: new TypeOrmLogger() })
}
