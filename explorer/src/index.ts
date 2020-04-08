import { Connection } from 'typeorm'
import { getDb } from './database'
import { retireSessions } from './entity/Session'
import { logger } from './logging'
import server from './server'

const cleanup = (conn: Connection) => {
  logger.info('Cleaning up sessions...')
  retireSessions(conn)
}

const start = () => {
  logger.info('Starting Explorer Node')
  server()
}

getDb()
  .then(cleanup)
  .then(start)
  .catch(e => {
    logger.error({
      msg: `Exception during startup: ${e.message}`,
      stack: e.stack,
    })
  })
