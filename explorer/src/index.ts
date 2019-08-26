import server from './server'
import { getDb } from './database'
import { logger } from './logging'
import { retireSessions } from './entity/Session'

const start = async () => {
  logger.info('Starting Explorer Node')
  server()
}

const cleanup = async () => {
  const db = await getDb()
  await retireSessions(db)
}

process.on('exit', cleanup)
process.on('SIGINT', cleanup)
process.on('SIGUSR1', cleanup)
process.on('SIGUSR2', cleanup)
process.on('uncaughtException', cleanup)

start().catch(logger.error)
