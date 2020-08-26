import { getConnection } from 'typeorm'
import { Config, Environment } from '../config'

if (Config.env() !== Environment.TEST) {
  throw Error(
    'trying to load test database in a non test db environment is not supported!',
  )
}

const TRUNCATE_TABLES: string[] = [
  'admin',
  'chainlink_node',
  'ethereum_head',
  'ethereum_log',
]

export const clearDb = (): Promise<any> =>
  getConnection().query(`TRUNCATE TABLE ${TRUNCATE_TABLES.join(',')} CASCADE`)
