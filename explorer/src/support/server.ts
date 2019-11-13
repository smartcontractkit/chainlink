import { getDb, closeDbConnection } from '../database'
import server from '../server'
import http from 'http'

export const DEFAULT_TEST_PORT =
  parseInt(process.env.TEST_SERVER_PORT, 10) || 8081

/**
 * Start database then initialize the server on the specified port
 */
export async function start(port: number = DEFAULT_TEST_PORT) {
  await getDb()
  return server(port)
}

/**
 * Stop the server then close the database connection
 */
export function stop(server: http.Server, done: jest.DoneCallback): void {
  server.close(() => closeDbConnection().then(done))
}
