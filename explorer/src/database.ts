import 'reflect-metadata'
import { Connection, createConnection } from 'typeorm'
import { PostgresConnectionOptions } from 'typeorm/driver/postgres/PostgresConnectionOptions'
import ormconfig from '../ormconfig.json'
import { TypeOrmLogger } from './logging'

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
  for (const config of ormconfig) {
    if (isEnvEqual(config.name, env)) {
      delete config.name
      return config
    }
  }
  throw Error(`env ${env} not found in options from ormconfig.json`)
}

// Loads the following ENV vars, giving them precedence.
// i.e. TYPEORM_PORT will replace "port" in ormconfig.json.
const mergeOptions = (): PostgresConnectionOptions => {
  const envOptions: Record<string, string> = {}
  for (const v of overridableKeys) {
    const envVar = process.env[`TYPEORM_${v.toUpperCase()}`]
    if (envVar) {
      envOptions[v] = envVar
    }
  }

  const connectionOpts = {
    ...loadOptions(),
    ...envOptions,
    logger: new TypeOrmLogger(),
  } as PostgresConnectionOptions

  return connectionOpts
}

export const openDbConnection = async (): Promise<Connection> => {
  return createConnection(mergeOptions())
}
