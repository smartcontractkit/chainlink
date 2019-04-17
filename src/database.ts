import 'reflect-metadata'
import { createConnection, Connection } from 'typeorm'
import { PostgresConnectionOptions } from 'typeorm/driver/postgres/PostgresConnectionOptions'
import options from '../ormconfig.json'

const overridableKeys = ['host', 'port', 'username', 'password', 'database']

// isEnvEqual returns true if the option name is the same as the env second paramter
// with the following exception that development == default.
const isEnvEqual = (optionName: string, env: string): boolean => {
  return (
    env === optionName || (env === 'development' && optionName === 'default')
  )
}

const loadOptions = (env?: string) => {
  env = env || process.env.TYPEORM_NAME || process.env.NODE_ENV || 'default'
  for (let option of options) {
    if (isEnvEqual(option.name, env)) {
      return option
    }
  }
  throw Error(`env ${env} not found in options from ormconfig.json`)
}

// Loads the following ENV vars, giving them precedence.
// i.e. TYPEORM_PORT will replace "port" in ormconfig.json.
const mergeOptions = (): PostgresConnectionOptions => {
  let envOptions: { [key: string]: string } = {}
  for (let v of overridableKeys) {
    const envVar = process.env[`TYPEORM_${v.toUpperCase()}`]
    if (envVar) {
      envOptions[v] = envVar
    }
  }
  return {
    ...loadOptions(),
    ...envOptions
  } as PostgresConnectionOptions
}

let db: Connection | undefined

export const getDb = async (): Promise<Connection> => {
  if (db == null) {
    db = await createConnection(mergeOptions())
  }
  if (db == null) {
    throw new Error('no db connection returned')
  }
  return db
}

export const closeDbConnection = async (): Promise<void> => {
  const saveDb = db
  db = null
  return saveDb.close()
}
