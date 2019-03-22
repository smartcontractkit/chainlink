import 'reflect-metadata'
import { createConnection, Connection } from 'typeorm'
import { PostgresConnectionOptions } from 'typeorm/driver/postgres/PostgresConnectionOptions'
import options from '../ormconfig.json'

const overridableKeys = ['host', 'port', 'username', 'password', 'database']

// Loads the following ENV vars, giving them precedence.
// i.e. TYPEORM_PORT will replace "port" in ormconfig.json.
const loadOptions = (): PostgresConnectionOptions => {
  let envOptions: { [key: string]: string } = {}
  for (let v of overridableKeys) {
    const envVar = process.env[`TYPEORM_${v.toUpperCase()}`]
    if (envVar) {
      envOptions[v] = envVar
    }
  }
  return {
    ...options,
    ...envOptions
  } as PostgresConnectionOptions
}

let db: Connection
export const createDbConnection = async () => {
  db = await createConnection(loadOptions())
}

export const getDb = (): Connection => {
  return db
}
