import server from './server'
import { getDb } from './database'
import { logger } from './logging'
import { retireSessions } from './entity/Session'

const cleanup = () => {
  logger.info('Cleaning up sessions...')
  getDb().then(retireSessions)
}
cleanup()

const start = async () => {
  logger.info('Starting Explorer Node')
  server()
}

start().catch(logger.error)
