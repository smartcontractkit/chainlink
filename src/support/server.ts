import { createDbConnection } from '../database'
import seed from '../seed'
import server from '../server'

export const start = async () => {
  await createDbConnection()
  return server()
}

export const startAndSeed = async () => {
  const server = await start()
  seed()
  return server
}
