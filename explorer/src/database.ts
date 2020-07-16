import 'reflect-metadata'
import { Connection, createConnection } from 'typeorm'
import { PostgresConnectionOptions } from 'typeorm/driver/postgres/PostgresConnectionOptions'
import { TypeOrmLogger } from './logging'
import { Config } from './config'

const overridableKeys = ['host', 'port', 'username', 'password', 'database']

// Loads the following ENV vars, giving them precedence.
// i.e. TYPEORM_PORT will replace "port" in ormconfig.json.
const mergeOptions = (
  loadedOptions: PostgresConnectionOptions,
): PostgresConnectionOptions => {
  const envOptions: Record<string, string> = {}
  for (const v of overridableKeys) {
    const envVar = process.env[`TYPEORM_${v.toUpperCase()}`]
    if (envVar) {
      envOptions[v] = envVar
    }
  }

  const connectionOpts = {
    ...loadedOptions,
    ...envOptions,
    logger: new TypeOrmLogger(),
  } as PostgresConnectionOptions

  return connectionOpts
}

export const openDbConnection = async (): Promise<Connection> => {
  const options = await import(`../ormconfig/${Config.typeorm()}.json`)
  return createConnection({
    ...mergeOptions(options),
    logger: new TypeOrmLogger(),
  })
}
