import { getDb } from '../database'
import server from '../server'

export const DEFAULT_TEST_PORT =
  parseInt(process.env.TEST_SERVER_PORT, 10) || 8081

export const start = async (port: number = DEFAULT_TEST_PORT) => {
  await getDb()
  return server(port)
}
