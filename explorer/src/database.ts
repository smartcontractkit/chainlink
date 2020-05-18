import 'reflect-metadata'
import { Connection, createConnection } from 'typeorm'
import { TypeOrmLogger } from './logging'
import { Config } from './config'

export const openDbConnection = async (): Promise<Connection> => {
  const options = await import(`../ormconfig/${Config.typeorm()}.json`)
  return createConnection({ ...options, logger: new TypeOrmLogger() })
}
