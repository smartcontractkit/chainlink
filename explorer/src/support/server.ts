import { randomBytes } from 'crypto'
import http from 'http'
import server from '../server'
import { Server } from 'http'
import { Config } from '../config'

/**
 * Start database then initialize the server on the specified port
 */

export async function start(): Promise<Server> {
  Config.setEnv('EXPLORER_SERVER_PORT', `${Config.testPort()}`)
  Config.setEnv('EXPLORER_COOKIE_SECRET', randomBytes(32).toString('hex'))
  Config.setEnv('GA_ID', 'GA-123')
  return server()
}

/**
 * Stop the server
 */
export function stop(server: http.Server, done: jest.DoneCallback): void {
  server.close(done)
}
