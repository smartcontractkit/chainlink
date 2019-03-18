import "reflect-metadata"
import { createConnection, ConnectionOptions, ConnectionOptionsReader } from "typeorm"
import seed from "./seed"
import server from "./server"

const nodeEnv = process.env.NODE_ENV || 'development'
const ormconfigPath = `ormconfig.${nodeEnv}.json`

const loadOptions = async (configName: string): Promise<ConnectionOptions> => {
  const reader = new ConnectionOptionsReader({configName})
  return reader.get("default")
}

loadOptions(ormconfigPath).then(options => {
  createConnection(options).then(async dbConnection => {
    seed(dbConnection)
    server(dbConnection)
  })
}).catch(error => console.log(error))
