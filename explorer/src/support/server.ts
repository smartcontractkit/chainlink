import { randomBytes } from 'crypto'
import http from 'http'
import { getConfig } from '../config'
import { closeDbConnection, getDb } from '../database'
import server from '../server'

export const DEFAULT_TEST_PORT =
  parseInt(process.env.EXPLORER_TEST_SERVER_PORT, 10) || 8081

/**
 * Start database then initialize the server on the specified port
 */
export async function start() {
  Object.assign(process.env, {
    EXPLORER_SERVER_PORT: `${DEFAULT_TEST_PORT}`,
    EXPLORER_COOKIE_SECRET: randomBytes(32).toString('hex'),
  })
  const conf = getConfig()
  await getDb()
  return await server(conf)
}

/**
 * Stop the server then close the database connection
 */
export function stop(server: http.Server, done: jest.DoneCallback): void {
  server.close(() => closeDbConnection().then(done))
}
