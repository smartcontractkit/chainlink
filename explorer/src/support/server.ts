import { Connection } from 'typeorm'
import { randomBytes } from 'crypto'
import http from 'http'
import { getConfig } from '../config'
import server from '../server'
import { Server } from 'http'

export const DEFAULT_TEST_PORT =
  parseInt(process.env.EXPLORER_TEST_SERVER_PORT, 10) || 8081

/**
 * Start database then initialize the server on the specified port
 */
export async function start(db: Connection): Promise<Server> {
  Object.assign(process.env, {
    EXPLORER_SERVER_PORT: `${DEFAULT_TEST_PORT}`,
    EXPLORER_COOKIE_SECRET: randomBytes(32).toString('hex'),
  })

  const conf = getConfig()
  return server(conf, db)
}

/**
 * Stop the server then close the database connection
 */
export function stop(server: http.Server, done: jest.DoneCallback): void {
  server.close(done)
}
