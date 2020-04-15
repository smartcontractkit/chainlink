import { getConnection, Connection } from 'typeorm'

if (process.env.NODE_ENV !== 'test') {
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

export const getDb = () : Connection => {
  return getConnection()
}

export const clearDb = async () => {
  return getDb().query(`TRUNCATE TABLE ${TRUNCATE_TABLES.join(',')} CASCADE`)
}
